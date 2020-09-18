// 函数原型
// 函数声明
// 函数定义
// 函数代理
// 代理获得
// 函数使用

// 函数原型
// 	func (string,string) bool
// 函数声明
// 	func login(name ,password string ) bool
// 函数定义
// 	func login(name ,password string ) bool{ return name == "" && password == ""}
// 函数代理
// 	aop.proxy(func (string,string) bool,login,"func_id",joinpointer...)
// 代理获得
// 	function = aop.get("func_id)
// 函数使用
// 	flag = function("hello","world")

// Step 1: Define Beans factory
// beanFactory := aop.NewClassicBeanFactory()
// beanFactory.RegisterBean("auth", new(Auth))

// Step 2: Define Aspect
// aspect := aop.NewAspect("aspect_1", "auth")
// aspect.SetBeanFactory(beanFactory)

// Step 3: Define Pointcut
// pointcut := aop.NewPointcut("pointcut_1").Execution(`Login()`)
// aspect.AddPointcut(pointcut)

// Step 4: Add Advice
// aspect.AddAdvice(&aop.Advice{Ordering: aop.Before, Method: "Before", PointcutRefID: "pointcut_1"})
// aspect.AddAdvice(&aop.Advice{Ordering: aop.After, Method: "After", PointcutRefID: "pointcut_1"})
// aspect.AddAdvice(&aop.Advice{Ordering: aop.Around, Method: "Around", PointcutRefID: "pointcut_1"})

// Step 5: Create AOP
// gogapAop := aop.NewAOP()
// gogapAop.SetBeanFactory(beanFactory)
// gogapAop.AddAspect(aspect)

// Setp 6: Get Proxy
// proxy, err := gogapAop.GetProxy("auth")

// Last Step: Enjoy
// login := proxy.Method(new(Auth).Login).(func(string, string) bool)("zeal", "gogap")
// fmt.Println("login result:", login)


type Joinpointer struct {
}
type Pointcut interface {
	BeforeAdvice(jp *joinpointer)
	AfterAdvice(jp *joinpointer)
	AroundAdvice(jp *joinpointer) interface{}...
}

func ProxyFunc(function interface{}, pc Pointcut) {

}
