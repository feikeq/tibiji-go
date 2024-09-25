/*
在工程化的Go语言开发项目中，Go语言的源码复用是建立在包（package）基础之上的。
一个文件夹下面直接包含的文件只能归属一个package,同样一个package的文件不能再多个文件夹下。
所有的包名都应该使用小写字母。包名可以不和文件夹的名字一样（包名不能包含-符号）包名和文件夹名字不一样时不会自动引包
如果想要构建一个程序，则包和包内的文件都必须以正确的顺序进行编译。包的依赖关系决定了其构建顺序。
属于同一个包的源文件必须全部被一起编译，一个包即是编译时的一个单元，因此根据惯例，每个目录都只包含一个包。
*/
// 通过 package 这种方式在 Go 中将要创建的应用程序定义为可执行程序（可以运行的文件）
package controllers //要创建的包名

// 通过 import 语句可从其他包中的其他代码访问程序
// 一个 Go 程序是通过 import 关键字将一组包链接在一起。
import (
	// import这里不是指文件名而是package-name(代码内的 package 包名)

	// Go的标准库包含了大量的包（如： fmt 和 os）

	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	// 本项目自己的包
	"tibiji-go/config"
	"tibiji-go/middle"
	"tibiji-go/models"
	"tibiji-go/utils"

	// 第三方开源的包 （ “_”是特殊标识符用来忽略结果，只是引用该包仅仅是为了调用init()函数所以无法通过包名来调用包中的其他函数 ）
	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

/*
func (变量名 结构体类型) 方法名(参数列表) 返回类型{
  // 方法体（方法就是一个属于特定类型的函数）
  // 其实就是在函数名称之前引用结构体来为该结构体添加方法。
  // 若要创建一个方法修改结构体实例，该方法必须引用指向结构体的指针:func (变量名 *结构体类型) 方法名(参数列表) 返回类型{}
}
func 函数名(参数列表) 返回类型{
  // 函数体
}

在 Go 中，错误并不会视为异常。
没有 try 或 catch 的机制。作为替代，如果发生错误，需要从函数内返回错误。
Go 支持一个函数返回多个值。如果调用的函数可能会返回错误，你必须检测这个错误是否存在，然后处理这个错误。
func GetName(name string) (string, error) {
  if name == "Bob" {
    return "", fmt.Errorf("Name cannot be Bob")
  }
  return name, nil
}
func main() {
  name, err := GetName("Bob")
  if err != nil {
    fmt.Println("Uh-oh an error has occurred")
  }
}




如果想让一个包中的标识符（如变量、常量、类型、函数等）能被外部的包使用，那么标识符必须是对外可见的（public）。在Go语言中是通过标识符的首字母大/小写来控制标识符的对外可见（public）/不可见（private）的。在一个包内部只有首字母大写的标识符才是对外可见的。

// 如需将某些内容设为专用内容，请以小写字母开始。(小驼峰只能从包内调用方法或变量)
func print(msg string, end string) {
	fmt.Println("-------sys-------")
	fmt.Print(msg, end)
}

// 如需将某些内容设为公共内容，请以大写字母开始。(大驼峰可以从任何位置访问变量或函数，建议你添加注释来描述此函数的用途)
func Log(msg string) {
	print(msg, "\r\n")
}


内部方法与外部方法
在 Go 语言中，函数名的首字母大小写非常重要，它被来实现控制对方法的访问权限
当方法的首字母为大写时，这个方法对于所有包都是Public，其他包可以随意调用
当方法的首字母为小写时，这个方法是Private，其他包是无法访问的

接口就是方法签名（Method Signature）的集合
在 Go 语言中封装和继承是通过 struct 来实现的，而多态则是通过接口(interface)来实现的。
接口指定了一个类型应该具有的方法，并由该类型决定如何实现这些方法
使用 type 关键字来定义接口
type Phone interface {
   call()
}
这里定义了一个电话接口，接口要求必须实现 call 方法
*/

// IPRequestInfo 用于存储每个 IP 地址的请求信息
type IPRequestInfo struct {
	Count       int       // 请求计数器（如果需要统计请求次数或者做其他相关处理）
	LastRequest time.Time // 上次请求时间
}

/*
ipRequestMap 结构：

	{
	    "192.168.172.166": {
	        "Count": 2,//正常访问次数
	        "LastRequest": "2024-06-28T09:46:36.114986+08:00"//正常访问的最后一次时间
	    },
	    "192.168.172.88": {
	        "Count": 2,
	        "LastRequest": "2024-06-28T09:46:43.306207+08:00"
	    }
	}
*/
var (
	ipRequestMap = make(map[string]*IPRequestInfo) // 存储 IP 地址的请求信息
	mu           sync.Mutex                        // 互斥锁，保护 ipRequestMap 的并发访问
)

// 定义一个结构体
/*
struct结构体类型 - 封装多个基本数据类型
Go语言中没有“类”的概念，也不支持“类”的继承等面向对象的概念，通过struct来实现面向对象。
Go语言中通过结构体的内嵌再配合接口比面向对象具有更高的扩展性和灵活性。
匿名结构体 var user struct{Name string; Age int} 在定义一些临时数据结构等场景下使用匿名结构体
*/
type UserController struct {
	DB     *sqlx.DB          // 控制器 func (c *UserController) 里使用c.DB访问数据库连接
	Models *models.UserModel // 模型

	// 在 Iris 框架中，控制器的 iris.Context 字段会被自动注入为控制器方法的参数。
	// 如果在控制器结构体中定义了名为 CTX 的 iris.Context 字段，Iris 框架会自动将上下文对象赋值给该字段。
	CTX iris.Context // CTX iris.Context 是控制器中的一个字段，它不需要初始化。当您在处理请求时，Iris框架会自动将当前请求的上下文信息填充到这个字段中。
}

// 构造函数创建控制器对象并进行初始化
// 通常用于将依赖项注入到控制器中，例如数据库连接、配置对象等。
// 通过使用构造函数，可以确保在创建控制器对象时进行必要的初始化操作，并将所需的依赖项传递给控制器。
func NewUserController(db *sqlx.DB, cfg map[string]interface{}) *UserController {
	// 返回一个结构体指针
	return &UserController{
		DB:     db,
		Models: models.NewUserModel(db, cfg),
	}
}

/*
// 在 UserController 结构体上添加方法
func (u *UserController) area() int {
    return u.width * u.height
}
*/

// 模块初始化函数 import 包时被调用
// go语言中init函数用于包(package)的初始化，该函数是go语言的一个重要特性。
// 在 Go 中，init() 函数有着特定的用途，它被用于初始化包（package）级别的数据和执行特定的初始化逻辑。init() 函数在程序运行时会自动被调用，无需手动调用。
func init() {
	// 当程序启动的时候，init函数会按照它们声明的顺序自动执行。
	/*
			1 init函数是用于程序执行前做包的初始化的函数，比如初始化包里的变量等
		    2 每个包可以拥有多个init函数
		    3 包的每个源文件也可以拥有多个init函数
		    4 同一个包中多个init函数的执行顺序go语言没有明确的定义(说明)
		    5 不同包的init函数按照包导入的依赖关系决定该初始化函数的执行顺序
		    6 init函数不能被其他函数调用，而是在main函数执行之前，自动被调用


		虽然 init() 函数可以用于一些初始化操作，但它并不是构造函数的替代品。主要的区别如下：
			调用时机：init() 函数是在程序启动时自动调用的，而构造函数是在创建对象时显式调用的。
			对象创建：构造函数主要用于创建和初始化对象，而 init() 函数用于初始化包级别的数据和执行特定的初始化逻辑，不直接与对象创建相关。
			调用方式：构造函数需要显式地调用来创建对象，而 init() 函数是自动被调用的，无法手动调用。
		通常情况下，构造函数用于创建对象并进行初始化，而 init() 函数用于执行包级别的初始化逻辑。它们在功能和调用方式上有一定的区别，应根据需要选择适当的方式进行对象初始化和程序初始化。
		需要注意的是，在 Go 中，每个包可以包含多个 init() 函数，它们按照声明的顺序依次执行。在包的初始化过程中，init() 函数的执行顺序是由编译器决定的，并且不能手动调用 init() 函数。
		综上所述，init() 函数和构造函数在 Go 中有不同的用途和调用方式，它们并不是相互替代的关系。
	*/
	// fmt.Println("------------UserController模块加载---------------")
}

// HandleError 捕获控制器的方法和服务时间依赖注入错误（覆盖 mvcApp.HandleError 函数）
// 访问控制器的字段 Context（不需要手动绑定）如 Ctx iris.Context 或者当作方法参数输入，
// 如 func(ctx iris.Context, otherArguments...)。
func (c *UserController) HandleError(ctx iris.Context, err error) {
	// // 在 Get 方法中调用 Post 方法
	// // 如果在同一个结构体中的方法之间可以相互调用。所以在 Get 方法中您可以通过直接调用 c.Post(ctx) 来调用 Post 方法：
	// response := c.Post() // 前提此方法有返回值
	// fmt.Printf("Type: %T , Data: %v\n", response, response)

	// // 调取别的模块下的模型或控制器

	// // 创建 XXXModel 实例
	// XXXModel := models.NewXXXModel(c.DB)
	// // 使用 XXXModel 进行操作
	// XXX, _ := XXXModel.Read(uid)
	// // 直接赋值 XXXController 中的结构体
	// XXXController := &XXXController{DB:c.DB,Models:models.NewXXXModel(c.DB)}
	// // 或者调用 XXXController 中的工厂函数
	// XXXController := NewXXXController(c.DB)

	/*
		在 Iris 框架中，iris.Map 是一个特定类型的映射（map），它是 Iris 框架为了方便处理 HTTP 请求和响应而提供的一种数据结构。iris.Map 类型是 map[string]interface{} 的别名，它允许你以键值对的形式存储和操作数据。
		iris.Map 的设计目的是为了在处理 HTTP 请求和响应时，提供一种简单的方式来操作数据，并能方便地将数据转换为 JSON、XML 等格式。它在处理中间件、控制器以及模板渲染时经常被使用。
		iris.Map：map[string]interface{} 的别名，用于处理键值对形式的数据。
		iris.MapSlice：[]iris.Map 的别名，用于处理具有顺序的键值对数据。
		iris.MapStrings：map[string][]string 的别名，用于处理具有重复键的字符串映射。
		iris.MapAny：map[string]interface{} 的别名，用于处理任意类型的键值对数据。
		这些类型都是为了在 Iris 框架中处理请求和响应时提供更方便的数据结构。你可以根据具体的需求选择合适的类型来处理和操作数据。
		需要注意的是，这些类型都是 Iris 框架特有的类型，并不是 Go 语言标准库中的类型。它们被设计为适用于 Iris 框架的特定场景和需求。
	*/

	// ctx.StopWithError(iris.StatusBadRequest, err)
	// Note that,
	// you can ignore this error and continue by not stopping the execution.
	ctx.HTML(fmt.Sprintf("<i>%s</i>", err.Error()))
}

// User 用户的控制器。
/*

在MVC模式下，控制器方法的返回值取决于您的设计和需求。在Iris Go框架中，控制器方法的返回值通常用于指定响应的内容和状态码。
在GetBy方法中，您可以使用ctx.JSON方法直接将结果作为JSON响应发送给客户端。因此，您不需要显式地返回任何值。相反，您将通过设置适当的HTTP状态码和JSON响应来告知客户端请求的结果。
对于Get和Post方法，它们通常需要返回视图模板、重定向URL或其他数据，以便在处理后继请求时进行使用。因此，您可能需要在这些方法中返回适当的值。
请注意，Iris Go框架中的控制器方法返回值不受限制，您可以根据您的需求自由选择返回值的类型和数量。以上解释仅是一种常见的用法示例。根据您的项目要求，您可以根据需要自定义返回值。


通过控制器方法的输入参数访问动态路径参数，不需要绑定。
mvc.New(app.Party("/user")).Handle(new(user.Controller))：
在控制器中每个以HTTP方法(Get，Post，Put，Delete...) 为前缀的函数，都作为一个 HTTP 端点。
当你使用 Iris 的默认语法从一个控制器中解析处理器，你需要定义方法的后缀为 By。
例如：
func(*UserController) Get() - GET:/user
func(*UserController) Post() - POST:/user
func(*UserController) GetLogin() - GET:/user/login
func(*UserController) PostLogin() - POST:/user/login
func(*UserController) GetProfileFollowers() - GET:/user/profile/followers
func(*UserController) PostProfileFollowers() - POST:/user/profile/followers
func(*UserController) GetBy(id int64) - GET:/user/{param:long}
func(*UserController) PostBy(id int64) - POST:/user/{param:long}

By关键字：建立了没有样板的动态路由， 控制器由此知道怎么处理 GetWelcomeBy 上的 "name" 和 "numTimes"，

*/

// 用户列表 GET:/user/
// pageNumber 分页页码
// pageSize 每页条数
// pageOrder 排序方向
// pageField 排序字段
func (c *UserController) Get() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果不是管理员
	if !c.Models.IsAdmin(tkUid) {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 获取每页条数配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	pageSize := otherCfg["SERV_LIST_SIZE"].(int64)

	// 分页参数
	var pageNumber int64 = 1
	if pageSize < 1 {
		pageSize = 20 // 默认单页条数
	}

	var pageOrder, pageField string = "desc", "uid" //  pageOrder 也支持这种写法 ascend descend
	var searchOR string                             // 模糊查询字段

	// // 判断是否存在字段 "key"
	// if val, ok := allData["key"]; ok {
	// 	// 当ok成立后，表示allData["key"]这个键存在并且对应的值不为nil
	// 	// 在判断val是否为空字符串时，可以直接使用if val != ""来判断，而不需要再进行断言转义为string类型。
	// 	if val != nil {
	// 		val = allData["key"].(string)
	// 	}
	// 	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。
	// }

	// 判断是否存在字段 "page"
	if _, ok := allData["pageNumber"]; ok {
		pageNumber = utils.ParseInt64(allData["pageNumber"]) // 任意数据转int64数字
	}
	// 判断是否存在字段 "pageSize"
	if _, ok := allData["pageSize"]; ok {
		pageSize = utils.ParseInt64(allData["pageSize"]) // 任意数据转int64数字
	}

	// 判断是否存在字段 "pageOrder"
	if _, ok := allData["pageOrder"]; ok {
		pageOrder = allData["pageOrder"].(string)
	}
	// 判断是否存在字段 "pageField"
	if _, ok := allData["pageField"]; ok {
		pageField = allData["pageField"].(string)
	}

	// 判断是否存在字段 "searchOR"
	if val, ok := allData["searchOR"]; ok {
		searchOR = val.(string)
	}

	list, total, err := c.Models.List(allData, pageNumber, pageSize, pageOrder, pageField, searchOR)
	if err != nil {
		if env != "" {
			println("Models.List Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	type temp struct {
		Total      int64              `json:"total" description:"总条数"`
		PageNumber int64              `json:"pageNumber" description:"分页页码"`
		PageSize   int64              `json:"pageSize" description:"每页条数"`
		PageOrder  string             `json:"pageOrder" description:"排序方向"`
		PageField  string             `json:"pageField" description:"排序字段"`
		List       []models.UserUnion `json:"list" description:"列表数据"`
	}
	data := temp{total, pageNumber, pageSize, pageOrder, pageField, list}
	// 返回成功响应
	// 向数据模板传值 当然也可以绑定其他值
	// ctx.ViewData("", mapData)
	// ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"data": data, "code": 0, "msg": ""})
}

// 新增用户 POST:/user/
func (c *UserController) Post() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	allData := utils.AllDataToMap(ctx) // ctx.FormValues() 等同于 ctx.Request().Form
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	var errTxt = ""
	// 判断是否存在字段 "username"
	if _, ok := allData["username"]; !ok {
		errTxt = "用户名username不能为空"
	} else {
		if allData["username"] == "" {
			errTxt = "用户名username不能为空"
		}
	}

	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		}
		return
	}

	// 客户端IP地址
	allData["regip"] = utils.GetRealIP(ctx)

	// 通过用户名判断类型并更新相应的数据字段
	username := allData["username"].(string)
	// 判断类型
	if utils.CheckEmail(username) {
		allData["email"] = username
	} else if utils.CheckMobile(username) {
		allData["cell"] = username
	} else if utils.CheckIdCard(username) {
		allData["identity_card"] = username
	}

	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	isCheck := otherCfg["SERV_OPEN_CHECK"].(bool) // 是否开启验证(登录注册是否验证)
	fmt.Printf("是否开启验证(登录注册是否验证) 类型type: %T, 值value: %v\n", isCheck, isCheck)

	isAdmin := c.Models.IsAdmin(tkUid) //检查是否为管理员
	if isCheck {
		// 如果不是管理员
		if !isAdmin {
			var errTxt = ""
			// 判断是否存在字段 "ticket"
			if _, ok := allData["ticket"]; !ok {
				errTxt = "验证令牌ticket不能为空"
			} else {
				if allData["ticket"] == "" {
					errTxt = "验证令牌ticket不能为空"
				}
			}

			// println("ticket", allData["ticket"])

			// 判断是否存在字段 "code"
			if _, ok := allData["code"]; !ok {
				errTxt = "验证码code不能为空"
			} else {
				if allData["code"] == "" {
					errTxt = "验证码code不能为空"
				}
			}
			if errTxt != "" {
				if env != "" {
					println("errTxt Error: ", errTxt)
					ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

				}
				return
			}
			ticket := allData["ticket"].(string)
			code := allData["code"].(string)
			ua := ctx.GetHeader("User-Agent")                                     // 拿到UA信息User-Agent
			secret := username + otherCfg["SERV_KEY_SECRET"].(string) + code + ua // 验证码的特殊密钥

			println(ticket, code, secret)

			// 进行 Basic Auth 身份认证
			temp_uid, err := utils.VerifyToken(ticket, secret)
			if err != nil {
				println("VerifyTicket Error: ", err.Error())

				switch err.Error() {
				case "0":
					println("code无效")
					ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
				case "1":
					println("code已过期")
					ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode]})
				default:
					println("code其他错误")
					ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
				}

				return
			}

			println("验证通过为0的虚拟ID", temp_uid) // 继续往下执行

		}

	}

	// 如果是管理添加的用户
	if isAdmin {
		allData["inviter"] = tkUid     //邀请者UID
		allData["referer"] = "管理员后台添加" //用户来源
	}

	// 调取创建用户模型 - 返回新插入数据的id
	uid, err := c.Models.Create(allData)
	if err != nil {
		if env != "" {
			println("Models.Create Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 初始化用户附属资料
	errM := c.Models.CreateMaterial(map[string]interface{}{
		"uid": uid,
	})
	if errM != nil {
		if env != "" {
			println("Models.CreateMaterial Error: ", errM.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": uid, "_debug_err": errM.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 返回成功响应
	// ctx.StatusCode(iris.StatusOK)
	// ctx.JSON(iris.Map{"data": uid, "code": 0, "msg": ""})

	// 如果开启验证 并且 不是管理 返回登录信息
	if isCheck && !isAdmin {
		// 返回登录信息
		// 获取用户信息
		user, err := c.Models.Read(uid)
		// user, err := c.Models.ReadMaterial(uid) // 这里好像不用取完整附属资料，因为新注册的cid铁定是0
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
			}
			return
		}

		// 数据处理 - 转换时间格式
		*user.Birthday = utils.RFC3339ToString(*user.Birthday, 0)
		*user.Intime = utils.RFC3339ToString(*user.Intime, 2) //防止拿到秒级精确时间
		*user.Uptime = utils.RFC3339ToString(*user.Uptime, 2) //防止拿到秒级精确时间
		// 对手机号等敏感信息进行脱敏处理
		*user.Cell = utils.MaskPhoneNumber(*user.Cell)
		// 对邮箱进行脱敏处理
		*user.Email = utils.MaskEmail(*user.Email)
		// 对银行卡进行脱敏处理
		*user.Bankcard = utils.MaskBankCardNumber(*user.Bankcard)
		// 对身份证进行脱敏处理
		*user.IdentityCard = utils.MaskIDCardNumber(*user.IdentityCard)
		// 对真实姓名脱敏
		*user.FName = utils.MaskRealName(*user.FName)

		// // 对密码进行脱敏处理
		// if user.Ciphers == "" {
		// 	user.Ciphers = "0"
		// } else {
		// 	user.Ciphers = "1"
		// }

		result := utils.StructToMap(user, "json") // 结构体转MAP

		// 对密码进行类型转换的脱敏处理
		if result["ciphers"] != "" {
			result["password"] = true
		} else {
			result["password"] = false
		}
		delete(result, "ciphers")

		ua := ctx.GetHeader("User-Agent") // 拿到UA信息User-Agent
		exptime := otherCfg["SERV_EXPIRES_TIME"].(int64)
		secret := otherCfg["SERV_KEY_SECRET"].(string) + ua
		// 添加 token
		token, _ := utils.GenerateToken(*user.UID, exptime, secret)
		result["token"] = token

		// 操作写入日志表
		logData := map[string]interface{}{
			"uid":    *user.UID,
			"action": "login",
			"note":   "register",
			"actip":  utils.GetRealIP(ctx),
			"ua":     ua,
		}
		log := c.Models.SetLogs(logData)
		if log != nil {
			if env != "" {
				println("Models.SetLogs Error: ", log.Error())
				ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
				return
			}
		}

		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})

	} else {
		// 返回成功响应
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{"data": uid, "code": 0, "msg": ""})
	}
}

// 用户信息 GET:/user/{uid}/
func (c *UserController) GetBy(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}
	// // 拿所有提交数据
	// allData := utils.AllDataToMap(ctx)

	// 如果操作人不是自己
	if tkUid != id {
		// 如果不是管理员
		if !c.Models.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// ctx.Writef("身份验证: %s \n", authentication)
	// fmt.Printf("用户ID： %d ！\n", uid)

	user, err := c.Models.Read(id)
	if err != nil {
		if env != "" {
			println("Models.Read Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// println(user.UID, user.Object)
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", user.Object, user.Object)

	// 在控制器层将结果进行修改和脱敏并得到最终的数据

	// 数据处理 - 转换时间格式
	*user.Birthday = utils.RFC3339ToString(*user.Birthday, 0)
	*user.Intime = utils.RFC3339ToString(*user.Intime, 2) //防止拿到秒级精确时间
	*user.Uptime = utils.RFC3339ToString(*user.Uptime, 2) //防止拿到秒级精确时间
	// 对手机号等敏感信息进行脱敏处理
	*user.Cell = utils.MaskPhoneNumber(*user.Cell)
	// 对邮箱进行脱敏处理
	*user.Email = utils.MaskEmail(*user.Email)
	// 对银行卡进行脱敏处理
	*user.Bankcard = utils.MaskBankCardNumber(*user.Bankcard)
	// 对身份证进行脱敏处理
	*user.IdentityCard = utils.MaskIDCardNumber(*user.IdentityCard)
	// 对真实姓名脱敏
	*user.FName = utils.MaskRealName(*user.FName)

	// // 对密码进行脱敏处理
	// if user.Ciphers == "" {
	// 	user.Ciphers = "0"
	// } else {
	// 	user.Ciphers = "1"
	// }

	result := utils.StructToMap(user, "json") // 结构体转MAP

	// 对密码进行类型转换的脱敏处理
	if result["ciphers"] != "" {
		result["password"] = true
	} else {
		result["password"] = false
	}
	delete(result, "ciphers")

	// 在 Iris 框架中 iris.Map 类型是 map[string]interface{} 的别名，它允许你以键值对的形式存储和操作数据。
	ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""}) // ctx.JSON(iris.Map{"Author": "肥客泉", "userID": id})

}

// 修改用户  PUT:/user/{uid}/
func (c *UserController) PutBy(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果操作人不是自己
	if tkUid != id {
		// 如果不是管理员
		if !c.Models.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 删除不能修改的字段
	delete(allData, "uid") // 删除 用户ID

	delete(allData, "username")      // 删除 帐号
	delete(allData, "ciphers")       // 删除 密码
	delete(allData, "email")         // 删除 邮箱
	delete(allData, "cell")          // 删除 电话
	delete(allData, "identity_card") // 删除 身份证

	delete(allData, "regip")    // 删除 注册IP
	delete(allData, "intime")   // 删除 注册时间
	delete(allData, "uptime")   // 删除 更新时间
	delete(allData, "userclan") // 删除 用户拓谱图

	// 权限不够则删除
	delete(allData, "referer")  // 管理员才能修改 用户来源
	delete(allData, "inviter")  // 管理员才能修改  邀请ID
	delete(allData, "state")    // 管理员才能修改 状态
	delete(allData, "grouptag") // 管理员才能修改 用户组

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Update(id, allData)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 删除用户 DELETE:/user/{uid}/
func (c *UserController) DeleteBy(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果操作人不是自己
	if tkUid != id {
		// 如果不是管理员
		if !c.Models.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// // 拿所有提交数据
	// allData := utils.AllDataToMap(ctx)
	// // fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 调取模型 - 根据ID删除数据库中的信息
	row, err := c.Models.Delete(id)
	if err != nil {
		if env != "" {
			println("Models.Delete Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 获取附属资料 GET:/user/{uid}/material/
func (c *UserController) GetByMaterial(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果操作人不是自己
	if tkUid != id {
		// 如果不是管理员
		if !c.Models.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// 拿所有提交数据
	// allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 获取附属资料调取模型 - 根据ID读取数据库中的信息
	material, err := c.Models.ReadMaterial(id)
	if err != nil {
		if env != "" {
			println("Models.ReadMaterial Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 对密码进行脱敏处理
	if *material.Ciphers == "" {
		*material.Ciphers = "0"
	} else {
		*material.Ciphers = "1"
	}

	// 数据处理 - 转换时间格式 不用脱敏 直接拿到秒级精确时间
	*material.Birthday = utils.RFC3339ToString(*material.Birthday, 0)
	*material.Intime = utils.RFC3339ToString(*material.Intime, 2)
	*material.Uptime = utils.RFC3339ToString(*material.Uptime, 2)

	material.Exptime = utils.RFC3339ToString(material.Exptime)
	material.MaterialUptime = utils.RFC3339ToString(material.MaterialUptime, 2)

	// 如果不是管理员
	if !c.Models.IsAdmin(tkUid) {
		// 对手机号等敏感信息进行脱敏处理
		*material.Cell = utils.MaskPhoneNumber(*material.Cell)
		// 对邮箱进行脱敏处理
		*material.Email = utils.MaskEmail(*material.Email)
		// 对银行卡进行脱敏处理
		*material.Bankcard = utils.MaskBankCardNumber(*material.Bankcard)
		// 对身份证进行脱敏处理
		*material.IdentityCard = utils.MaskIDCardNumber(*material.IdentityCard)
		// 对真实姓名脱敏
		*material.FName = utils.MaskRealName(*material.FName)
	}

	ctx.JSON(iris.Map{"data": material, "code": 0, "msg": ""})
}

// 修改附属资料 PUT:/user/{uid}/material/
func (c *UserController) PutByMaterial(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果不是管理员
	if !c.Models.IsAdmin(tkUid) {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	delete(allData, "material_uptime") // 删除 附属更新时间
	delete(allData, "uid")             // 删除 uid
	delete(allData, "intime")          // 删除 intime
	delete(allData, "uptime")          // 删除 uptime

	if _, ok := allData["ciphers"]; ok {
		if allData["ciphers"] != "" {
			user, err := c.Models.Read(id)
			if err != nil {
				if env != "" {
					println("Models.Read Error: ", err.Error())
					ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
				}
				return
			}
			intime := utils.RFC3339ToString(*user.Intime) // 拿到秒级精确时间
			// 计算出用户密码的脱敏数据
			allData["ciphers"] = utils.HashPassword(allData["ciphers"].(string), intime)
		}
	}

	// 根据ID更新数据库中的信息
	_, u_err := c.Models.Update(id, allData)
	if u_err != nil {
		if env != "" {
			println("Models.Update Error: ", u_err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": u_err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 判断是否存在字段 "material_remark"
	if _, ok := allData["material_remark"]; ok {
		allData["remark"] = allData["material_remark"]
		delete(allData, "material_remark") // 删除 附属备注
	}

	// 判断是否存在字段 "material_object"
	if _, ok := allData["material_object"]; ok {
		allData["object"] = allData["material_object"]
		delete(allData, "material_object") // 删除 附属预留值
	}

	// 更新用户附属资料表调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.UpdateMaterial(id, allData)
	if err != nil {
		if env != "" {
			println("Models.UpdateMaterial Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ua := ctx.GetHeader("User-Agent") // 拿到UA信息User-Agent

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    id,
		"action": "material",
		"note":   utils.SerializeJSON(allData),
		"actip":  utils.GetRealIP(ctx),
		"ua":     ua,
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", log.Error())
			ctx.JSON(iris.Map{"data": row, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
			return
		}
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 更新用户 PATCH:/user/{uid}/
func (c *UserController) PatchBy(id int64) {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果操作人不是自己
	if tkUid != id {
		// 如果不是管理员
		if !c.Models.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 更新指定用户UID的密码、邮箱、电话、用户名、身份证
	// 注意：每次从只能更新一项数据
	errCode := 0

	// 只筛选出一个字段进行更新
	upData := map[string]interface{}{}
	key, val := "", ""
	// 判断是否存在字段
	if _, ok := allData["ciphers"]; ok {
		_ciphers := allData["ciphers"].(string)
		key = "ciphers"
		if _ciphers != "" {
			user, err := c.Models.Read(id)
			if err != nil {
				if env != "" {
					println("Models.Read Error: ", err.Error())
					ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
				}
				return
			}
			// pwd旧密码原来密码
			var _pwd string
			if _, ok := allData["pwd"]; ok {
				_pwd = allData["pwd"].(string)
			} else {
				errCode = config.ErrParamEmpty
			}

			intime := utils.RFC3339ToString(*user.Intime) // 拿到秒级精确时间
			// println(intime)                               // 2023-05-22 09:13:35

			dbPwd := *user.Ciphers
			oldPwd := utils.HashPassword(_pwd, intime) //原密码
			// println("dbPwd", dbPwd)
			// println("odPwd", oldPwd)
			// print(_pwd, " : ", oldPwd == *user.Ciphers)

			// 比较原来密码
			if oldPwd == dbPwd {
				// 计算出用户密码的脱敏数据
				allData["ciphers"] = utils.HashPassword(_ciphers, intime)
			} else {
				// 密码错误
				errCode = config.ErrFormat

			}

		} else {
			errCode = config.ErrParamEmpty
		}
	} else if _, ok := allData["email"]; ok {
		key = "email"
		if allData["email"] != "" {
			if !utils.CheckEmail(allData["email"].(string)) {
				errCode = config.ErrFormat
			}
		} else {
			errCode = config.ErrParamEmpty
		}
	} else if _, ok := allData["cell"]; ok {
		key = "cell"
		if allData["cell"] != "" {
			if !utils.CheckMobile(allData["cell"].(string)) {
				errCode = config.ErrFormat
			}
		} else {
			errCode = config.ErrParamEmpty
		}
	} else if _, ok := allData["identity_card"]; ok {
		key = "identity_card"
		if allData["identity_card"] != "" {
			if !utils.CheckIdCard(allData["identity_card"].(string)) {
				errCode = config.ErrFormat
			}
			if _, ok := allData["fname"]; ok {
				upData["fname"] = allData["fname"]
			}
		} else {
			errCode = config.ErrParamEmpty
		}
	} else if _, ok := allData["username"]; ok {
		key = "username"
		if allData["username"] == "" {
			errCode = config.ErrParamEmpty
		}
	} else {
		errCode = config.ErrParamEmpty
	}

	if errCode != 0 {
		ctx.JSON(iris.Map{"data": key, "code": errCode, "msg": config.ErrMsgs[errCode]})
		return
	}

	// 只筛选出一个字段进行更新
	val = allData[key].(string)
	upData[key] = allData[key]

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	isCheck := otherCfg["SERV_OPEN_CHECK"].(bool) // 是否开启验证(登录注册是否验证)
	fmt.Printf("是否开启验证(登录注册是否验证) 类型type: %T, 值value: %v\n", isCheck, isCheck)
	// 当开启 SERV_OPEN_CHECK 是否开启验证 时 ，且登录方式是用邮箱/手机时 code 和 ticket 才为必填字段，当为 用户名/身份证 登录时可不用验证码就可以登录。
	if isCheck && (key == "cell" || key == "email") {

		var errTxt = ""
		// 判断是否存在字段 "ticket"
		if _, ok := allData["ticket"]; !ok {
			errTxt = "验证令牌ticket不能为空"
		} else {
			if allData["ticket"] == "" {
				errTxt = "验证令牌ticket不能为空"
			}
		}

		println("ticket", allData["ticket"])

		// 判断是否存在字段 "code"
		if _, ok := allData["code"]; !ok {
			errTxt = "验证码code不能为空"
		} else {
			if allData["code"] == "" {
				errTxt = "验证码code不能为空"
			}
		}
		if errTxt != "" {
			if env != "" {
				println("errTxt Error: ", errTxt)
				ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

			}
			return
		}
		ticket := allData["ticket"].(string)
		code := allData["code"].(string)
		ua := ctx.GetHeader("User-Agent")                                // 拿到UA信息User-Agent
		secret := val + otherCfg["SERV_KEY_SECRET"].(string) + code + ua // 验证码的特殊密钥

		println(ticket, code, secret)

		// 进行 Basic Auth 身份认证
		temp_uid, err := utils.VerifyToken(ticket, secret)
		if err != nil {
			println("VerifyTicket Error: ", err.Error())

			switch err.Error() {
			case "0":
				println("code无效")
				ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
			case "1":
				println("code已过期")
				ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode]})
			default:
				println("code其他错误")
				ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
			}

			return
		}
		//因为和验token是一个机制，所以验证code通过时大多是为0的虚拟ID
		println("code验证通过是为0的虚拟ID", temp_uid) // 继续往下执行
	}

	// 调取模型 - 根据ID更新数据库中的信息
	_, err := c.Models.Update(id, upData)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    id,
		"action": "change",
		"note":   key + " > " + val,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", log.Error())
			ctx.JSON(iris.Map{"data": key, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
			return
		}
	}

	ctx.JSON(iris.Map{"data": key, "code": 0, "msg": "ok"})
}

// 用户登录 POST:/user/login/
func (c *UserController) PostLogin() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	var errTxt = ""
	// 判断是否存在字段 "name"
	if _, ok := allData["name"]; !ok {
		errTxt = "用户名name不能为空"
	} else {
		if allData["name"] == "" {
			errTxt = "用户名name不能为空"
		}
	}
	// 判断是否存在字段 "pwd"
	if _, ok := allData["pwd"]; !ok {
		errTxt = "用户密码pwd不能为空"
	} else {
		if allData["pwd"] == "" {
			errTxt = "用户pwd不能为空"
		}
	}

	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

		}
		return
	}

	pwd := allData["pwd"].(string)
	name := allData["name"].(string)

	// 判断类型
	typeName := "username"
	if utils.CheckEmail(name) {
		typeName = "email"
	} else if utils.CheckMobile(name) {
		typeName = "cell"
	} else if utils.CheckIdCard(name) {
		typeName = "identity_card"
	}

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	isCheck := otherCfg["SERV_OPEN_CHECK"].(bool) // 是否开启验证(登录注册是否验证)
	fmt.Printf("是否开启验证(登录注册是否验证) 类型type: %T, 值value: %v\n", isCheck, isCheck)
	// 当开启 SERV_OPEN_CHECK 是否开启验证 时 ，且登录方式是用邮箱/手机时 code 和 ticket 才为必填字段，当为 用户名/身份证 登录时可不用验证码就可以登录。
	if isCheck && (typeName == "cell" || typeName == "email") {

		var errTxt = ""
		// 判断是否存在字段 "ticket"
		if _, ok := allData["ticket"]; !ok {
			errTxt = "验证令牌ticket不能为空"
		} else {
			if allData["ticket"] == "" {
				errTxt = "验证令牌ticket不能为空"
			}
		}

		println("ticket", allData["ticket"])

		// 判断是否存在字段 "code"
		if _, ok := allData["code"]; !ok {
			errTxt = "验证码code不能为空"
		} else {
			if allData["code"] == "" {
				errTxt = "验证码code不能为空"
			}
		}
		if errTxt != "" {
			if env != "" {
				println("errTxt Error: ", errTxt)
				ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

			}
			return
		}
		ticket := allData["ticket"].(string)
		code := allData["code"].(string)
		ua := ctx.GetHeader("User-Agent")                                 // 拿到UA信息User-Agent
		secret := name + otherCfg["SERV_KEY_SECRET"].(string) + code + ua // 验证码的特殊密钥

		println(ticket, code, secret)

		// 进行 Basic Auth 身份认证
		temp_uid, err := utils.VerifyToken(ticket, secret)
		if err != nil {
			println("VerifyTicket Error: ", err.Error())

			switch err.Error() {
			case "0":
				println("code无效")
				ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
			case "1":
				println("code已过期")
				ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode]})
			default:
				println("code其他错误")
				ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
			}

			return
		}
		//因为和验token是一个机制，所以验证code通过时大多是为0的虚拟ID
		println("code验证通过是为0的虚拟ID", temp_uid) // 继续往下执行
	}

	// 查找用户 (使用username、email、cell、identity_card查找用户)
	user_tmp, typeName, err := c.Models.Find(name)
	if err != nil {
		if env != "" {
			println("Models.Find Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrValidation, "msg": config.ErrMsgs[config.ErrValidation], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrValidation, "msg": config.ErrMsgs[config.ErrValidation]})
		}
		return
	}
	// 获取附属资料调取模型 - 根据ID读取数据库中的信息
	user, err := c.Models.ReadMaterial(*user_tmp.UID)
	if err != nil {
		if env != "" {
			println("Models.ReadMaterial Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": user_tmp, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 在控制器层将结果进行修改和脱敏并得到最终的数据
	if *user.State == 2 {
		if env != "" {
			println("帐号还未激活")
			ctx.JSON(iris.Map{"code": config.ErrNoActivate, "msg": config.ErrMsgs[config.ErrNoActivate], "_debug_carry": allData, "_debug_err": "帐号还未激活"})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrNoActivate, "msg": config.ErrMsgs[config.ErrNoActivate]})
		}
		return
	} else if *user.State == 0 {
		if env != "" {
			println("帐号已被禁用")
			ctx.JSON(iris.Map{"code": config.ErrUserDisabled, "msg": config.ErrMsgs[config.ErrUserDisabled], "_debug_carry": allData, "_debug_err": "帐号已被禁用"})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrUserDisabled, "msg": config.ErrMsgs[config.ErrUserDisabled]})
		}
		return
	}

	// 转换时间格式
	inTime := utils.RFC3339ToString(*user.Intime) // 拿到秒级精确时间
	// println(inTime) // 2023-05-22 09:13:35

	dbPwd := *user.Ciphers
	inPwd := utils.HashPassword(pwd, inTime)
	// println("dbPwd", dbPwd)
	// println("inPwd", inPwd)
	// print(pwd, " : ", dbPwd == inPwd)

	if dbPwd != inPwd {
		if env != "" {
			println("密码错误")
			ctx.JSON(iris.Map{"code": config.ErrValidation, "msg": config.ErrMsgs[config.ErrValidation], "_debug_carry": allData, "_debug_err": "密码错误"})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrValidation, "msg": config.ErrMsgs[config.ErrValidation]})
		}
		return
	}

	// 数据处理 - 转换时间格式
	*user.Birthday = utils.RFC3339ToString(*user.Birthday, 0)
	*user.Intime = utils.RFC3339ToString(*user.Intime, 2) //防止拿到秒级精确时间
	*user.Uptime = utils.RFC3339ToString(*user.Uptime, 2) //防止拿到秒级精确时间

	user.Exptime = utils.RFC3339ToString(user.Exptime)
	user.MaterialUptime = utils.RFC3339ToString(user.MaterialUptime, 2)

	// 对手机号等敏感信息进行脱敏处理
	*user.Cell = utils.MaskPhoneNumber(*user.Cell)
	// 对邮箱进行脱敏处理
	*user.Email = utils.MaskEmail(*user.Email)
	// 对银行卡进行脱敏处理
	*user.Bankcard = utils.MaskBankCardNumber(*user.Bankcard)
	// 对身份证进行脱敏处理
	*user.IdentityCard = utils.MaskIDCardNumber(*user.IdentityCard)
	// 对真实姓名脱敏
	*user.FName = utils.MaskRealName(*user.FName)

	// // 对密码进行脱敏处理
	// if user.Ciphers == "" {
	// 	user.Ciphers = "0"
	// } else {
	// 	user.Ciphers = "1"
	// }

	result := utils.StructToMap(user, "json") // 结构体转MAP

	// 对密码进行类型转换的脱敏处理
	if result["ciphers"] != "" {
		result["password"] = true
	} else {
		result["password"] = false
	}
	delete(result, "ciphers")

	ua := ctx.GetHeader("User-Agent") // 拿到UA信息User-Agent
	exptime := otherCfg["SERV_EXPIRES_TIME"].(int64)
	secret := otherCfg["SERV_KEY_SECRET"].(string) + ua
	// 添加 token
	token, _ := utils.GenerateToken(*user.UID, exptime, secret)
	result["token"] = token

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    *user.UID,
		"action": "login",
		"note":   typeName,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ua,
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", log.Error())
			ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
			return
		}
	}

	ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
}

// 用户日志 GET:/user/logs/
func (c *UserController) GetLogs() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	// 如果不是管理员
	if !c.Models.IsAdmin(tkUid) {
		// var aa = utils.ParseInt(tkUid)
		// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", aa, aa)
		// // int64类型的值直接赋给 map[string]interface{} 时，因溢出会发生类型转换错误只能是int，这里和json里的int最大值一个道理。
		// allData["uid"] = utils.ParseInt(tkUid) // 只能查自己的记录
		// allData["uid"] = tkUid  // 因为 utils.GetWhereArgs 反回的值只是字符串，所以不能直接这样赋值
		allData["uid"] = strconv.FormatInt(tkUid, 10) //数字转十进制字符串
		fmt.Printf("allData: %+v\n", allData)         // 打印allData
	}

	// 获取每页条数配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	pageSize := otherCfg["SERV_LIST_SIZE"].(int64)

	// 分页参数
	var pageNumber int64 = 1
	if pageSize < 1 {
		pageSize = 20 // 默认单页条数
	}

	var pageOrder, pageField string = "desc", "id" //  pageOrder 也支持这种写法 ascend descend

	// 判断是否存在字段 "page"
	if _, ok := allData["pageNumber"]; ok {
		pageNumber = utils.ParseInt64(allData["pageNumber"]) // 任意数据转int64数字
	}
	// 判断是否存在字段 "pageSize"
	if _, ok := allData["pageSize"]; ok {
		pageSize = utils.ParseInt64(allData["pageSize"]) // 任意数据转int64数字
	}

	// 判断是否存在字段 "pageOrder"
	if _, ok := allData["pageOrder"]; ok {
		pageOrder = allData["pageOrder"].(string)
	}
	// 判断是否存在字段 "pageField"
	if _, ok := allData["pageField"]; ok {
		pageField = allData["pageField"].(string)
	}

	list, total, err := c.Models.GetLogs(allData, pageNumber, pageSize, pageOrder, pageField)
	if err != nil {
		if env != "" {
			println("Models.GetLogs Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	type temp struct {
		Total      int64             `json:"total" description:"总条数"`
		PageNumber int64             `json:"pageNumber" description:"分页页码"`
		PageSize   int64             `json:"pageSize" description:"每页条数"`
		PageOrder  string            `json:"pageOrder" description:"排序方向"`
		PageField  string            `json:"pageField" description:"排序字段"`
		List       []models.UserLogs `json:"list" description:"列表数据"`
	}
	data := temp{total, pageNumber, pageSize, pageOrder, pageField, list}
	ctx.JSON(iris.Map{"data": data, "code": 0, "msg": ""})
}

// 忘记密码 GET:/user/password/
func (c *UserController) GetPassword() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	{
		// 限制访问频率 - 接口访问频率限制： 利用 Go 标准库中的 sync.Mutex 和 time 包来实现基于 IP 地址的请求限流

		// 定义限流时间(秒)时间间隔是否大于limiting秒(相当于30秒内可以发生多少次 limitmax 判断)
		limiting := 30 * time.Second
		// 定义限流时间(秒)内最大请求次数(为0时不判断 就只走clearting规则，非0时相当于limitmax+1次)
		limitmax := 0
		// 定义超过时间(分钟)未请求视为过期(相当于5分钟可以发生多少次 limiting 判断)
		clearting := 5 * time.Minute

		// 获取客户端 IP 地址
		ip := ctx.RemoteAddr()

		mu.Lock()
		defer mu.Unlock()

		// 获取当前 IP 地址的请求信息，如果不存在则初始化
		info, exists := ipRequestMap[ip]
		if !exists {
			info = &IPRequestInfo{}
			ipRequestMap[ip] = info
		}

		// 清理 ipRequestMap 中过期的记录
		now := time.Now()
		for ip, all := range ipRequestMap {
			// _, _s := utils.FormatTimestamp(all.LastRequest.Unix())
			// println(ip, _s, now.Sub(all.LastRequest), " >", clearting)

			// now.Sub(LastRequest) 计算出从 LastRequest 到当前时间 (now) 的时间间隔，并与 clearting 进行比较。这行代码用来检查LastRequest和 当前时间之间的差。
			// time.Time 类型的 IsZero() 方法可以用来判断时间值是否为零值（即未被赋值的状态）。当调用 IsZero() 时，它会返回一个布尔值：如果时间是零值则返回 true，否则返回 false。
			if !all.LastRequest.IsZero() && now.Sub(all.LastRequest) > clearting {
				// 假设超过时间未请求，则视为过期
				delete(ipRequestMap, ip)
			}
		}

		// 检查上次请求时间和当前时间间隔是否大于limiting秒
		// println(info.Count, ">", limitmax)
		// time.Since()是从某个时间开始，期间返回自某个时间以来经过的时间。
		if time.Since(info.LastRequest) < limiting {
			if limitmax > 0 {
				if info.Count > limitmax {
					ctx.StatusCode(iris.StatusTooManyRequests)
					// ctx.WriteString("请求过于频繁，请稍后再试")
					if env != "" {
						println("ipRequestMap", ipRequestMap)
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
					} else {
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
					}
					return
				}
			} else {

				ctx.StatusCode(iris.StatusTooManyRequests)
				// ctx.WriteString("请求过于频繁，请稍后再试")
				if env != "" {
					println("ipRequestMap", ipRequestMap)
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
				}
				return

			}
		}

		// 更新计数器和上次请求时间
		info.Count++
		info.LastRequest = time.Now()

	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	var errTxt = ""
	// 判断是否存在字段 "name"
	if _, ok := allData["name"]; !ok {
		errTxt = "用户名name不能为空"
	} else {
		if allData["name"] == "" {
			errTxt = "用户名name不能为空"
		}
	}

	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

		}
		return
	}
	name := allData["name"].(string)

	// 查找用户 (使用username、email、cell、identity_card查找用户)
	user, typeName, err := c.Models.Find(name)
	if err != nil {
		if env != "" {
			println("Models.Find Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
		}
		return
	}
	println("typeName", typeName)

	// 对手机号等敏感信息进行脱敏处理
	cell := utils.MaskPhoneNumber(*user.Cell)
	// 对邮箱进行脱敏处理
	email := utils.MaskEmail(*user.Email)
	// // 对身份证进行脱敏处理
	// identityCard := utils.MaskIDCardNumber(*user.IdentityCard)

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	exptime := otherCfg["SERV_SAFE_ETIME"].(int64)                          // 密保时间                         // 设置密保超时时间(秒)
	smtpName := otherCfg["SMTP_FROM_NAME"].(string)                         // 发件人名称                        // 发件人名称
	smtpWebsite := otherCfg["SMTP_WEBSITE"].(string)                        // 邮件的网站网址
	timeFormat := ctx.Application().ConfigurationReadOnly().GetTimeFormat() // # 时间格式TimeFormat配置项

	duration := time.Duration(exptime) * time.Second // 将秒转为小时 不使用 /3600 的方式去计算
	hours := duration.Hours()                        // 多少小时
	now := time.Now()                                // 当前时间
	expStr := now.Add(duration).Format(timeFormat)   // 格式化密保超时时间 2023-06-01 14:35:04
	// // 用纯时间戳（毫秒）生成验证码code
	// milli := fmt.Sprintf("%d", now.UnixMilli()) // 获取时间戳（毫秒） 1670919222532 类似于JS里的 Date.now()
	// code := milli[len(milli)-6:]                // 取最后6位做为code验证码
	// println(milli,"code验证码:", code)

	// 用纯时间戳（毫秒）+5随机数  生成验证码code
	milli := utils.GenerateTimerID(99999) // （13位时间戳+5随机尾数每位最大到9
	code := milli[len(milli)-6:]          // 取最后6位做为code验证码
	println(now.UnixMilli(), milli, "用纯时间戳（毫秒）+5随机数  生成验证码code:", code)
	// 1719305088793 171930508879316643 用纯时间戳（毫秒）+5随机数  生成验证码code: 316643

	secret := otherCfg["SERV_KEY_SECRET"].(string) + code // 验证码的特殊密钥
	// 生成 token 返回给客户
	token, _ := utils.GenerateToken(*user.UID, exptime, secret)

	ctx.JSON(iris.Map{"data": token, "code": 0, "msg": email + " " + cell})

	if email != "" {
		subject := fmt.Sprintf("[%s安全中心]密码找回服务", smtpName)
		body := fmt.Sprintf("尊敬的%s用户您好：<br/>", smtpName)
		body += fmt.Sprintf("您的用户名：%s<br/>", *user.UserName)
		body += fmt.Sprintf("您的邮箱：%s<br/>", email)
		body += fmt.Sprintf("您的电话：%s<br/>", cell)
		// 将float64浮点只保留一位小数
		body += fmt.Sprintf("请务必在<b>%.1f</b>小时内通过下面这个地址修改您的密码，此链接将在%s后失效！<br/><br/>", hours, expStr)

		// url.QueryEscape函数在Go语言中的作用类似于JavaScript中的encodeURIComponent函数。它会对URL中的特殊字符进行编码，以确保URL参数能够正确传递并被解析。这样可以避免出现URL参数中含有特殊字符导致解析错误的情况。
		params := "/web/retrieve/?code=" + code + "&ticket=" + url.QueryEscape(token) + "&name=" + url.QueryEscape(name)
		body += fmt.Sprintf("<a href=\"%s\">%s</a> &nbsp; <br/>", smtpWebsite+params, strings.Split(smtpWebsite, "//")[1]+params)
		body += fmt.Sprintf("<b>您的的验证码: %s </b><br>", code)
		body += fmt.Sprintf("<br/>%s安全中心 %s<br/>", smtpName, now.Format(timeFormat))

		println("邮件发送:", *user.Email)
		println(body)
		// Go 并发线程 - 通过 go 关键字来开启 goroutine 即可
		go utils.SendEmail(ctx, *user.Email, subject, body) // 邮件发送
	}

	if cell != "" {
		smsTemplateIds := otherCfg["SMS_TEMPLATE_IDS"].(string) // 短信模版ID 模板类别(0其它 1生日 2纪念日 3闹铃)
		smsTemplate := strings.Split(smsTemplateIds, ",")       // 取配置中的SMS_TEMPLATE_IDS短信模版ID
		println("短信发送:", *user.Cell)
		println(code)
		//短信发送
		go utils.SendSMS(ctx, *user.Cell, smsTemplate[0], []string{code})

	}

}

// 找回密码后设置新密码 POST:/user/password/
func (c *UserController) PostPassword() {
	// 利用ticket验证令牌 和code验证码 来判断令牌是否正确

	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	{
		// 限制访问频率 - 接口访问频率限制： 利用 Go 标准库中的 sync.Mutex 和 time 包来实现基于 IP 地址的请求限流

		// 定义限流时间(秒)时间间隔是否大于limiting秒(相当于30秒内可以发生多少次 limitmax 判断)
		limiting := 30 * time.Second
		// 定义限流时间(秒)内最大请求次数(为0时不判断 就只走clearting规则，非0时相当于limitmax+1次)
		limitmax := 0
		// 定义超过时间(分钟)未请求视为过期(相当于5分钟可以发生多少次 limiting 判断)
		clearting := 5 * time.Minute

		// 获取客户端 IP 地址
		ip := ctx.RemoteAddr()

		mu.Lock()
		defer mu.Unlock()

		// 获取当前 IP 地址的请求信息，如果不存在则初始化
		info, exists := ipRequestMap[ip]
		if !exists {
			info = &IPRequestInfo{}
			ipRequestMap[ip] = info
		}

		// 清理 ipRequestMap 中过期的记录
		now := time.Now()
		for ip, all := range ipRequestMap {
			// _, _s := utils.FormatTimestamp(all.LastRequest.Unix())
			// println(ip, _s, now.Sub(all.LastRequest), " >", clearting)

			// now.Sub(LastRequest) 计算出从 LastRequest 到当前时间 (now) 的时间间隔，并与 clearting 进行比较。这行代码用来检查LastRequest和 当前时间之间的差。
			// time.Time 类型的 IsZero() 方法可以用来判断时间值是否为零值（即未被赋值的状态）。当调用 IsZero() 时，它会返回一个布尔值：如果时间是零值则返回 true，否则返回 false。
			if !all.LastRequest.IsZero() && now.Sub(all.LastRequest) > clearting {
				// 假设超过时间未请求，则视为过期
				delete(ipRequestMap, ip)
			}
		}

		// 检查上次请求时间和当前时间间隔是否大于limiting秒
		// println(info.Count, ">", limitmax)
		// time.Since()是从某个时间开始，期间返回自某个时间以来经过的时间。
		if time.Since(info.LastRequest) < limiting {
			if limitmax > 0 {
				if info.Count > limitmax {
					ctx.StatusCode(iris.StatusTooManyRequests)
					// ctx.WriteString("请求过于频繁，请稍后再试")
					if env != "" {
						println("ipRequestMap", ipRequestMap)
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
					} else {
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
					}
					return
				}
			} else {

				ctx.StatusCode(iris.StatusTooManyRequests)
				// ctx.WriteString("请求过于频繁，请稍后再试")
				if env != "" {
					println("ipRequestMap", ipRequestMap)
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
				}
				return

			}
		}

		// 更新计数器和上次请求时间
		info.Count++
		info.LastRequest = time.Now()

	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	var errTxt = ""
	// 判断是否存在字段 "ticket"
	if _, ok := allData["ticket"]; !ok {
		errTxt = "验证令牌ticket不能为空"
	} else {
		if allData["ticket"] == "" {
			errTxt = "验证令牌ticket不能为空"
		}
	}

	// 判断是否存在字段 "pwd"
	if _, ok := allData["pwd"]; !ok {
		errTxt = "用户新密码pwd不能为空"
	} else {
		if allData["pwd"] == "" {
			errTxt = "用户新密码pwd不能为空"
		}
	}

	// 判断是否存在字段 "code"
	if _, ok := allData["code"]; !ok {
		errTxt = "验证码code不能为空"
	} else {
		if allData["code"] == "" {
			errTxt = "验证码code不能为空"
		}
	}
	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

		}
		return
	}
	ticket := allData["ticket"].(string)
	pwd := allData["pwd"].(string)
	code := allData["code"].(string)

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	secret := otherCfg["SERV_KEY_SECRET"].(string) + code // 验证码的特殊密钥

	// 进行 Basic Auth 身份认证
	uid, err := utils.VerifyToken(ticket, secret)
	if err != nil {
		println("VerifyTicket Error: ", err.Error())
		// ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"code": config.ErrVerificationCode, "msg": config.ErrMsgs[config.ErrVerificationCode]})
		return
	}

	// 获取用户信息
	user, err := c.Models.Read(uid)
	if err != nil {
		if env != "" {
			println("Models.Read Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
		}
		return
	}

	// 只筛选出一个字段进行更新
	upData := map[string]interface{}{}
	key, val := "ciphers", ""
	// 判断是否存在字段
	intime := utils.RFC3339ToString(*user.Intime) // 拿到秒级精确时间
	// 计算出用户密码的脱敏数据
	val = utils.HashPassword(pwd, intime)
	upData[key] = val

	// 调取模型 - 根据ID更新数据库中的信息
	_, upErr := c.Models.Update(*user.UID, upData)
	if upErr != nil {
		if env != "" {
			println("Models.Update Error: ", upErr.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": upErr.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    *user.UID,
		"action": "forgetpwd",
		"note":   key + " > " + val,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", log.Error())
			ctx.JSON(iris.Map{"code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
			return
		}
	}

	// 对手机号等敏感信息进行脱敏处理
	*user.Cell = utils.MaskPhoneNumber(*user.Cell)
	// 对邮箱进行脱敏处理
	*user.Email = utils.MaskEmail(*user.Email)
	// 对身份证进行脱敏处理
	*user.IdentityCard = utils.MaskIDCardNumber(*user.IdentityCard)

	data := map[string]interface{}{
		"username":      *user.UserName,     // 帐号
		"cell":          *user.Cell,         // 电话
		"email":         *user.Email,        // 邮箱
		"identity_card": *user.IdentityCard, // 身份证
	}

	ctx.JSON(iris.Map{"data": data, "code": 0, "msg": "密码重置成功"})
}

// 获取验证码 PATCH:/user/captcha/
func (c *UserController) PatchCaptcha() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	{
		// 限制访问频率 - 接口访问频率限制： 利用 Go 标准库中的 sync.Mutex 和 time 包来实现基于 IP 地址的请求限流

		// 定义限流时间(秒)时间间隔是否大于limiting秒(相当于30秒内可以发生多少次 limitmax 判断)
		limiting := 30 * time.Second
		// 定义限流时间(秒)内最大请求次数(为0时不判断 就只走clearting规则，非0时相当于limitmax+1次)
		limitmax := 0
		// 定义超过时间(分钟)未请求视为过期(相当于5分钟可以发生多少次 limiting 判断)
		clearting := 5 * time.Minute

		// 获取客户端 IP 地址
		ip := ctx.RemoteAddr()

		mu.Lock()
		defer mu.Unlock()

		// 获取当前 IP 地址的请求信息，如果不存在则初始化
		info, exists := ipRequestMap[ip]
		if !exists {
			info = &IPRequestInfo{}
			ipRequestMap[ip] = info
		}

		// 清理 ipRequestMap 中过期的记录
		now := time.Now()
		for ip, all := range ipRequestMap {
			// _, _s := utils.FormatTimestamp(all.LastRequest.Unix())
			// println(ip, _s, now.Sub(all.LastRequest), " >", clearting)

			// now.Sub(LastRequest) 计算出从 LastRequest 到当前时间 (now) 的时间间隔，并与 clearting 进行比较。这行代码用来检查LastRequest和 当前时间之间的差。
			// time.Time 类型的 IsZero() 方法可以用来判断时间值是否为零值（即未被赋值的状态）。当调用 IsZero() 时，它会返回一个布尔值：如果时间是零值则返回 true，否则返回 false。
			if !all.LastRequest.IsZero() && now.Sub(all.LastRequest) > clearting {
				// 假设超过时间未请求，则视为过期
				delete(ipRequestMap, ip)
			}
		}

		// 检查上次请求时间和当前时间间隔是否大于limiting秒
		// println(info.Count, ">", limitmax)
		// time.Since()是从某个时间开始，期间返回自某个时间以来经过的时间。
		if time.Since(info.LastRequest) < limiting {
			if limitmax > 0 {
				if info.Count > limitmax {
					ctx.StatusCode(iris.StatusTooManyRequests)
					// ctx.WriteString("请求过于频繁，请稍后再试")
					if env != "" {
						println("ipRequestMap", ipRequestMap)
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
					} else {
						ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
					}
					return
				}
			} else {

				ctx.StatusCode(iris.StatusTooManyRequests)
				// ctx.WriteString("请求过于频繁，请稍后再试")
				if env != "" {
					println("ipRequestMap", ipRequestMap)
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
				} else {
					ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
				}
				return

			}
		}

		// 更新计数器和上次请求时间
		info.Count++
		info.LastRequest = time.Now()

	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("allData: %+v\n", allData) // 打印allData

	var errTxt = ""
	// 判断是否存在字段 "name"
	if _, ok := allData["name"]; !ok {
		errTxt = "用户名name不能为空"
	} else {
		if allData["name"] == "" {
			errTxt = "用户名name不能为空"
		}
	}

	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		}
		return
	}
	name := allData["name"].(string) // 用户名/邮箱/手机/身份证

	// 判断类型
	typeName := "username"
	if utils.CheckEmail(name) {
		typeName = "email"
	} else if utils.CheckMobile(name) {
		typeName = "cell"
	} else if utils.CheckIdCard(name) {
		typeName = "identity_card"
	}

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	exptime := otherCfg["SERV_SAFE_ETIME"].(int64)                          // 密保时间                         // 设置密保超时时间(秒)
	smtpName := otherCfg["SMTP_FROM_NAME"].(string)                         // 发件人名称                        // 发件人名称
	timeFormat := ctx.Application().ConfigurationReadOnly().GetTimeFormat() // # 时间格式TimeFormat配置项

	duration := time.Duration(exptime) * time.Second // 将秒转为小时 不使用 /3600 的方式去计算
	hours := duration.Hours()                        // 多少小时
	now := time.Now()                                // 当前时间
	expStr := now.Add(duration).Format(timeFormat)   // 格式化密保超时时间 2023-06-01 14:35:04
	// // 用纯时间戳（毫秒）生成验证码code
	// milli := fmt.Sprintf("%d", now.UnixMilli()) // 获取时间戳（毫秒） 1670919222532 类似于JS里的 Date.now()
	// code := milli[len(milli)-6:]                // 取最后6位做为code验证码
	// println(milli,"code验证码:", code)

	// 用纯时间戳（毫秒）+5随机数  生成验证码code
	milli := utils.GenerateTimerID(99999) // （13位时间戳+5随机尾数每位最大到9
	code := milli[len(milli)-6:]          // 取最后6位做为code验证码
	println(now.UnixMilli(), milli, "用纯时间戳（毫秒）+5随机数  生成验证码code:", code)
	// 1719305088793 171930508879316643 用纯时间戳（毫秒）+5随机数  生成验证码code: 316643

	ua := ctx.GetHeader("User-Agent")                                 // 拿到UA信息User-Agent
	secret := name + otherCfg["SERV_KEY_SECRET"].(string) + code + ua // 验证码的特殊密钥(添加用户名name防止假令牌篡改登录或注册)
	// 生成 token 返回给客户
	token, _ := utils.GenerateToken(0, exptime, secret)
	println(token)
	// ctx.JSON(iris.Map{"data": token, "code": 0, "msg": typeName})

	if typeName == "email" {
		subject := fmt.Sprintf("[%s]验证码", smtpName)
		body := fmt.Sprintf("尊敬的%s用户您好：<br/>", smtpName)
		// 将float64浮点只保留一位小数
		body += fmt.Sprintf("请务必在<b>%.1f</b>小时内通过验证，验证码将在%s后失效！<br/><br/>", hours, expStr)
		body += fmt.Sprintf("<b>您的的验证码: %s </b><br>", code)
		body += fmt.Sprintf("<br/>%s %s<br/>", smtpName, now.Format(timeFormat))

		// println("邮件发送:", *user.Email)
		// Go 并发线程 - 通过 go 关键字来开启 goroutine 即可
		go utils.SendEmail(ctx, name, subject, body) // 邮件发送
		ctx.JSON(iris.Map{"data": token, "code": 0, "msg": typeName})
	} else if typeName == "cell" {
		smsTemplateIds := otherCfg["SMS_TEMPLATE_IDS"].(string) // 短信模版ID 模板类别(0其它 1生日 2纪念日 3闹铃)
		smsTemplate := strings.Split(smsTemplateIds, ",")       // 取配置中的SMS_TEMPLATE_IDS短信模版ID
		//短信发送
		go utils.SendSMS(ctx, name, smsTemplate[0], []string{code})
		ctx.JSON(iris.Map{"data": token, "code": 0, "msg": typeName})
	} else {
		ctx.JSON(iris.Map{"data": token, "code": 0, "msg": typeName})
	}

}

// 用户分组 GET:/contact/groups/
func (c *UserController) GetGroups() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 如果不是管理员
	if !c.Models.IsAdmin(tkUid) {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	// // 拿所有提交数据
	// allData := utils.AllDataToMap(ctx)

	// 获取所有用户分组
	row, err := c.Models.Groups()
	if err != nil {
		if env != "" {
			println("Models.Groups Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 测试功能 - 无操作
func (c *UserController) Put(ctx iris.Context) interface{} {
	fmt.Println("------------ user [PUT]---------------")

	json := ctx.FormValues() // ctx.FormValues() 等同于 ctx.Request().Form

	if json != nil {
		var s []string
		s = append(s, "3")
		json["hahahahah"] = s
	}

	// username := post
	// return "Create by user with username: " + username
	return json
}

// 测试功能 get:/user/test
func (c *UserController) GetTest() {
	ctx := c.CTX
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	{
		// 接口访问频率限制： 利用 Go 标准库中的 sync.Mutex 和 time 包来实现基于 IP 地址的请求限流

		// 定义限流时间(秒)时间间隔是否大于limiting秒
		limiting := 30 * time.Second
		// 定义超过时间(分钟)未请求视为过期
		clearting := 5 * time.Minute

		// 获取客户端 IP 地址
		ip := ctx.RemoteAddr()

		mu.Lock()
		defer mu.Unlock()

		// 获取当前 IP 地址的请求信息，如果不存在则初始化
		info, exists := ipRequestMap[ip]
		if !exists {
			info = &IPRequestInfo{}
			ipRequestMap[ip] = info
		}

		// 检查上次请求时间和当前时间间隔是否大于limiting秒
		// time.Since是从某个时间开始，期间返回自某个时间以来经过的时间。
		if time.Since(info.LastRequest) < limiting {
			ctx.StatusCode(iris.StatusTooManyRequests)
			// ctx.WriteString("请求过于频繁，请稍后再试")
			if env != "" {
				println("ipRequestMap", ipRequestMap)
				ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent], "_debug_carry": ipRequestMap})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrFrequent, "msg": config.ErrMsgs[config.ErrFrequent]})
			}
			return
		}

		// 更新计数器和上次请求时间
		info.Count++
		info.LastRequest = time.Now()

		// 清理 ipRequestMap 中过期的记录
		now := time.Now()
		for ip, info := range ipRequestMap {
			if now.Sub(info.LastRequest) > clearting {
				// 假设超过时间未请求，则视为过期
				delete(ipRequestMap, ip)
			}
		}

	}

}

// 手动指定哪个链接去执行哪个方法  - 自定义匹配
// 在控制器之间共享依赖关系或者在父 MVC 应用程序上注册它们，以及在控制器里通过 BeforeActivation 可选回调事件，都可以修改每个控制器的依赖关系。
// 每个控制器通过 BeforeActivation 自定义事件回调，用来自定义控制器的结构的方法与自定义路径处理程序（甚至有正则表达式路径）。
func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	// BeforeActivation只调用一次，在控制器适应主应用程序之前，做一些不用 mvc 也能做的事情。
	// b.Dependencies().Add/Remove
	// b.Router().Use/UseGlobal/Done // 和你已知的任何标准 API  调用

	// Handle(方法 ，路径 ，控制器函数的名称将被解析未一个处理程序 [ handler ] ，任何应该在 MyCustomHandler 之前运行的处理程序[ handlers ])
	// b.Handle("GET", "/something/{id:long}", "MyCustomHandler", anyMiddleware...)
	// b.Handle("GET", "/query", "UserInfo")

	// 甚至添加基于此控制器路由的全局中间件
	// anyMiddlewareHere := func(ctx iris.Context) {
	// 	println(" ---  anyMiddlewareHere --- ")
	// 	ctx.Next()
	// }

	b.Router().Use(middle.MiddlewareVerifyAdmin)
}

/* Can use more than one, the factory will make sure
that the correct http methods are being registered for each route
for this controller, uncomment these if you want:

func (c *UserController) Post() {}
func (c *UserController) Put() {}
func (c *UserController) Delete() {}
func (c *UserController) Connect() {}
func (c *UserController) Head() {}
func (c *UserController) Patch() {}
func (c *UserController) Options() {}
func (c *UserController) Trace() {}
func (c *UserController) All() {}   /  func (c *UserController) Any() {}




func (c *UserController) GetInfo() mvc.Result{
	// 如果只是写接口不需要视图 那么就用mvc.Result 如果需要视图那就是mvc.View
}

// 激活后，所有依赖项都已设置 - 因此对它们的只读访问
// 但仍然可以添加自定义控制器或简单的标准处理程序。
func (c *UserController) AfterActivation(a mvc.AfterActivation) {}
*/
