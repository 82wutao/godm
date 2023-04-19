package common

type Runnable func()
type Func[A interface{}, R interface{}] func(A) R                          // [[T], R]
type BinFunc[A1 interface{}, A2 interface{}, R interface{}] func(A1, A2) R //[[T, S], R]
type Consumer[A interface{}] func(A)                                       //[[T], None]
type BinConsumer[A1 interface{}, A2 interface{}] func(A1, A2)              //[[T, S], None]
type Supplier[R interface{}] func() R                                      //[[], R]
type Predicate[A interface{}] func(A) bool                                 //[[T], bool]
type Comparator[A interface{}] func(A, A) int                              //[[T, T], int]
type SortKey[A interface{}, R interface{}] func(A) R                       //[[T], R]
