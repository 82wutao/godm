package stream

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

type student struct {
	name   string
	gender int
	score  int
	weight decimal.Decimal
}

func getData() []*student {
	students := make([]*student, 10)
	students[0] = &student{name: "boy1", gender: 1, score: 99, weight: decimal.NewFromInt(80)}
	students[1] = &student{name: "boy2", gender: 1, score: 97, weight: decimal.NewFromInt(81)}
	students[2] = &student{name: "boy3", gender: 1, score: 95, weight: decimal.NewFromInt(82)}
	students[3] = &student{name: "boy4", gender: 1, score: 93, weight: decimal.NewFromInt(83)}
	students[4] = &student{name: "boy5", gender: 1, score: 91, weight: decimal.NewFromInt(84)}
	students[5] = &student{name: "girl5", gender: 0, score: 90, weight: decimal.NewFromInt(85)}
	students[6] = &student{name: "girl4", gender: 0, score: 92, weight: decimal.NewFromInt(86)}
	students[7] = &student{name: "girl3", gender: 0, score: 94, weight: decimal.NewFromInt(87)}
	students[8] = &student{name: "girl2", gender: 0, score: 96, weight: decimal.NewFromInt(88)}
	students[9] = &student{name: "girl1", gender: 0, score: 98, weight: decimal.NewFromInt(89)}
	return students
}

func Test_Foreach(t *testing.T) {
	students := getData()
	NewStreamFromSlice(students).Foreach(func(src interface{}) {
		src.(*student).score = 100
	})
	for _, s := range students {
		if s.score == 100 {
			continue
		}
		t.Errorf("%s 's score should be %d", s.name, 100)
	}
}

// Reduce 收敛所有元素为一个形态
func Test_Reduce(t *testing.T) {
	students := getData()
	sum := NewStreamFromSlice(students).Filter(func(src interface{}) bool {
		return src.(*student).gender == 0
	}).Reduce(func() interface{} { return 0 }, func(result, src interface{}) interface{} {
		return result.(int) + src.(*student).score
	}).(int)
	if sum != (88+98)*3-88 {
		t.Errorf("sum of girls 's score should be %d", (88+98)*3-88)
	}

	sum2 := NewStreamFromSlice(students).Reduce(func() interface{} {
		return decimal.Zero
	}, func(result, src interface{}) interface{} {
		return result.(decimal.Decimal).Add(src.(*student).weight)
	}).(decimal.Decimal)
	if !sum2.Equal(decimal.NewFromInt(80*5 + 89*5)) {
		t.Errorf("sum of girls 's score should be %d", 80*5+89*5)
	}
}

// Collect 收集元素到一个容器
func Test_Collect(t *testing.T) {
	students := getData()
	girls := NewStreamFromSlice(students).Filter(func(src interface{}) bool {
		return src.(*student).gender == 0
	}).Collect(func() interface{} { return make([]*student, 0) }).([]*student)

	for _, s := range girls {
		if s.gender != 0 {
			t.Errorf("want a girl ,but a boy")
		}
	}

}

// TopK 有序取头部
func Test_TopK(t *testing.T) {
	students := getData()
	top4 := NewStreamFromSlice(students).TopK(4, func() interface{} {
		return make([]*student, 0)
	}, func(a, b interface{}) int {
		sa := a.(*student)
		sb := b.(*student)
		if sa.score < sb.score {
			return 1
		}
		if sa.score > sb.score {
			return -1
		}
		return 0
	}).([]*student)

	if len(top4) != 4 {
		t.Errorf("topk, resultset should be 4")
	}
	for i, s := range top4 {
		if s.score != 99-i {
			t.Errorf("topk %d, score is wrange,expect %d actual %d", i, 99-i, s.score)
		}
	}
}

// Group 对元素进行分组，一般多用于从slice到map
func Test_Group(t *testing.T) {
	students := getData()
	group := NewStreamFromSlice(students).Group(func() interface{} {
		return make(map[int]*student, 0)
	}, func(src interface{}) (interface{}, interface{}) {
		s := src.(*student)
		return s.gender, s
	}).(map[int]*student)

	if len(group) != 2 {
		t.Errorf("group, resultset should be 2")
	}

}

// GroupAndReduce 对元素进行分组，并对同一组的元素进行收敛，
func Test_GroupAndReduce(t *testing.T) {
	students := getData()
	group := NewStreamFromSlice(students).GroupAndReduce(func() interface{} {
		return make(map[int][]*student, 0)
	}, func(src interface{}) (interface{}, interface{}) {
		s := src.(*student)
		return s.gender, s
	}, func() interface{} {
		return make([]*student, 0)
	}, func(result, src interface{}) interface{} {
		return append(result.([]*student), src.(*student))
	}).(map[int][]*student)

	if len(group) != 2 {
		t.Errorf("groupandreduce, resultset should be 2")
	}
	girls, g := group[0]
	if !g {
		t.Errorf("groupandreduce, resultset should have one group girl")
	}
	if len(girls) != 5 {
		t.Errorf("groupandreduce,  girl group should be 5length")
	}
	for i, girl := range girls {
		if girl.gender != 0 || girl.name != fmt.Sprintf("girl%d", 5-i) || girl.score != 90+(i*2) {
			t.Errorf("groupandreduce,  girl wrang")
		}
	}

	boys, b := group[1]
	if !b {
		t.Errorf("groupandreduce, resultset should have one group boy")
	}
	if len(boys) != 5 {
		t.Errorf("groupandreduce,  boy group should be 5length")
	}
	for i, boy := range boys {
		if boy.gender != 1 || boy.name != fmt.Sprintf("boy%d", i+1) || boy.score != 99-(i*2) {
			t.Errorf("groupandreduce,  girl wrang")
		}
	}
}

func Test_FindAny(t *testing.T) {
	students := getData()
	girls := NewStreamFromSlice(students).Filter(func(src interface{}) bool {
		return src.(*student).gender == 0
	}).Limit(2).Collect(func() interface{} {
		return make([]*student, 0)
	}).([]*student)
	if len(girls) != 2 {
		t.Errorf("groupandreduce,  boy group should be 2length")
	}
}

func Test_GroupAndAvg(t *testing.T) {
	students := getData()

	type Aggregate struct {
		sum   int64
		count int64
		avg   decimal.Decimal
	}
	group := NewStreamFromSlice(students).GroupAndReduce(func() interface{} {
		return make(map[int]*Aggregate)
	}, func(src interface{}) (interface{}, interface{}) {
		s := src.(*student)
		return s.gender, s.score
	}, func() interface{} {
		return new(Aggregate)
	}, func(result, src interface{}) interface{} {
		agg := result.(*Aggregate)
		agg.sum += int64(src.(int))
		agg.count++
		agg.avg = decimal.NewFromInt(agg.sum).DivRound(decimal.NewFromInt(agg.count), 2)
		return agg
	}).(map[int]*Aggregate)

	if !group[0].avg.Equal(decimal.RequireFromString("94.00")) {
		t.Error("girl avg not expection")
	}
	if !group[1].avg.Equal(decimal.RequireFromString("95.00")) {
		t.Error("boy avg not expection")
	}
}
