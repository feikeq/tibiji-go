/*
在工程化的Go语言开发项目中，Go语言的源码复用是建立在包（package）基础之上的。
一个文件夹下面直接包含的文件只能归属一个package,同样一个package的文件不能再多个文件夹下。
所有的包名都应该使用小写字母。包名可以不和文件夹的名字一样（包名不能包含-符号）包名和文件夹名字不一样时不会自动引包
如果想要构建一个程序，则包和包内的文件都必须以正确的顺序进行编译。包的依赖关系决定了其构建顺序。
属于同一个包的源文件必须全部被一起编译，一个包即是编译时的一个单元，因此根据惯例，每个目录都只包含一个包。
*/
// 通过 package 这种方式在 Go 中将要创建的应用程序定义为可执行程序（可以运行的文件）
package main //要创建的包名
/*包含的文件只能归属一个包，同一个包的文件不能在多个文件夹下。包名为main的包是应用程序的入口包，这种包编译后会得到一个可执行文件。*/

// 通过 import 语句可从其他包中的其他代码访问程序
// 一个 Go 程序是通过 import 关键字将一组包链接在一起。
import (
	// import这里不是指文件名而是package-name(代码内的 package 包名)

	// Go的标准库包含了大量的包（如：fmt 和 os）

	"os"

	// 本项目自己的包
	"tibiji-go/config"
	"tibiji-go/controllers"
	"tibiji-go/middle"

	// 第三方开源的包 （ “_”是特殊标识符用来忽略结果，只是引用该包仅仅是为了调用init()函数所以无法通过包名来调用包中的其他函数 ）
	_ "github.com/go-sql-driver/mysql"               // sqlx文档 http://jmoiron.github.io/sqlx/
	"github.com/kataras/iris/v12"                    // iris包
	"github.com/kataras/iris/v12/middleware/logger"  // iris的日志
	"github.com/kataras/iris/v12/middleware/recover" // 能够在崩溃时记录和恢复
	"github.com/kataras/iris/v12/mvc"                // iris的mcv
)

// main 函数是程序的起点。 因此，整个包中只能有一个 main 函数（在第一行中定义的那个）
// main 函数没有任何参数，并且不返回任何内容。 但这并不意味着其不能从用户读取值，如命令行参数。
// 如要访问 Go 中的命令行参数，可以使用用于保存传递到程序的所有参数的 os 包 和 os.Args 变量来执行操作。
func main() {

	// 新建iris
	app := iris.New()

	assetsPath := controllers.AssetsPath // 资源文件目录

	//  HandleDir 注册一个处理程序从特定目录（系统物理目录或应用程序嵌入目录）提供http请求静态文件服务
	// 还能设置其它缓存showList和IndexName选项 #文档  https://docs.iris-go.com/iris/file-server/introduction
	// 显示静态文件 显示“assets”文件夹以下的文件
	app.HandleDir("/", iris.Dir(assetsPath))
	/* 该HandleDir方法接受第三个可选参数DirOptions：
	app.HandleDir("/uploads", iris.Dir("./uploads"), iris.DirOptions{
		Compress: true,
	})
	*/
	// app.Favicon("./assets/favicon.ico") // 有了HandleDir指定根目录就不需要指定favicon图标了

	// 全局配置文件

	// 当您有两种配置时可一种用于开发另一种用于生产。
	// TOML 的输入字符串参数是“~”，然后它从主目录加载配置，并且可以在许多iris实例之间共享。
	cfg := iris.TOML("./config/cfg.ini")
	// cfg的内容在app.Configure或app.Run或app.Listen引用配置之后可在控制器里拿
	// 获取APP配置 cfg := ctx.Application().ConfigurationReadOnly()

	// // conf 配置文件的使用:
	// app.Configure(iris.WithConfiguration(iris.YAML("myconfig.yml")))
	// app.Configure(iris.WithConfiguration(iris.TOML("myconfig.tml")))
	// app.Configure(iris.WithConfiguration(cfg))

	// 获取配置项的值
	// timeFormat := cfg.GetTimeFormat() // # 时间格式TimeFormat配置项 app.ConfigurationReadOnly().GetTimeFormat()
	// println("时间格式模板：", timeFormat)
	addr := cfg.Other["SERV_ADDR"].(string)
	ver := cfg.Other["SERV_VERSION"].(string)
	name := cfg.Other["SERV_NAME"].(string)
	otherCfg := cfg.GetOther() // 返回Other字段里的所有配置项

	// // println("cfg:::", cfg) // println 函数不支持直接使用 %T 和 %v 来格式化输出
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", cfg, cfg)
	// println("otherCfg[\"SERV_NAME\"]:", otherCfg["SERV_NAME"].(string))
	// println("addr:", addr)

	println("=========================")
	println("       " + name)
	println("=========================")
	println("接口版本:", ver)

	// env := os.Environ() //environ 以 key=value 的形式返回所有环境变量。
	println("当前用户:", os.Getenv("USER"))
	println("当前目录:", os.Getenv("PWD"))
	// 读取环境变量，如果没有设置，默认为正式生产环境
	env := os.Getenv("TIBIJI_SERV_ENV")
	// 你可以通过在终端或命令提示符中使用以下命令来设置环境变量：
	// Unix系统 export TIBIJI_SERV_ENV=development
	// Windows系统 set TIBIJI_SERV_ENV=development
	println("当前环境:", env)

	// 将当前环境添加到上下文中(将值传递给下一个处理程序)
	// ctx.Values().Set() 值是处理程序（或中间件）在彼此之间进行通信的方式，
	// 获取时使用 ctx.Values().Get("xxxx")内存地址,当然也可以 GetString 或 GetInt64 拿实际值
	app.Use(func(ctx iris.Context) {
		ctx.Values().Set("ENV", env)
		ctx.Next() // 为了执行链中的下一个处理程序

		/*
			在 Iris 框架中，iris.Context 类型提供了多种方法来处理请求和响应。
			以下是一些常用的 iris.Context 方法：

			Request 和 Response 相关方法：

			ctx.Request()：获取请求的 *http.Request 对象。
			ctx.Response()：获取响应的 http.ResponseWriter 对象。
			ctx.GetHeader(name string)：获取指定名称的请求头信息。
			ctx.GetContentType()：获取请求的内容类型。
			ctx.GetContentLength()：获取请求的内容长度。
			ctx.RemoteAddr()：获取客户端的 IP 地址。
			ctx.Hostname()：获取主机名。
			ctx.URI()：获取请求的 URI 信息。
			ctx.Path()：获取请求的路径。

			参数相关方法：
			ctx.Params().Get(name string)：获取通过路径参数传递的值。
			ctx.URLParam(name string)：获取 URL 查询参数的值。
			ctx.FormValue(name string)：获取表单字段的值。
			ctx.FormValues()：获取所有表单字段的值。
			ctx.PostValue(name string)：获取 POST 请求的字段值。
			ctx.PostValues()：获取所有 POST 请求的字段值。

			Cookie 相关方法：
			ctx.SetCookie(cookie *http.Cookie)：设置一个 Cookie。
			ctx.GetCookie(name string)：获取指定名称的 Cookie。
			ctx.ReadCookies() []*http.Cookie：读取并返回所有的 Cookie。

			响应相关方法：
			ctx.StatusCode(code int)：设置响应状态码。
			ctx.WriteString(text string)：向响应写入字符串。
			ctx.JSON(v interface{})：以 JSON 格式写入响应。
			ctx.XML(v interface{})：以 XML 格式写入响应。
			ctx.HTML(htmlContents string)：向响应写入 HTML 内容。

			会话相关方法：
			ctx.Session()：获取或创建一个会话对象。
			ctx.Values().Get(key string)：获取上下文级别的值。
			ctx.Values().Set(key string, value interface{})：设置上下文级别的值。
			在Iris框架中，您可以使用ctx.Values()方法来设置和获取全局可访问的变量。但是，这些变量只在当前请求的上下文中有效，并且不能在不同的请求之间共享。
			在您的代码中，当您在GetHahaha()方法中设置wx_ccav变量时，它只在当前请求的上下文中有效。当您在GetHehehe()方法中尝试获取wx_ccav变量时，它将返回一个空字符串，因为它不在当前请求的上下文中。
			如果您想要在不同的请求之间共享变量，您可以考虑使用一个全局变量或者一个缓存库来(Redis或Memcached等)存储这些变量。

			解析相关方法 从请求上下文(ctx)中获取请求体数据
			ctx.FormValues()：该方法用于获取表单数据。它解析并返回请求中的表单数据，可以获取URL查询参数或通过POST请求发送的表单数据。返回的结果是一个url.Values类型的映射，可以通过键名获取对应的值。
			使用场景：当您需要获取通过表单提交的数据时，例如HTML表单中的输入字段值。

			ctx.GetBody()：该方法用于获取请求体的原始字节数据。它返回一个[]byte类型的字节数组，包含请求体的原始内容。
			使用场景：当您需要直接操作请求体的原始字节数据时，例如进行自定义的解析或处理。

			ctx.ReadJSON()：该方法用于将请求体解析为JSON格式，并将结果绑定到指定的结构体或变量上。它会自动解析请求体中的JSON数据，并将解析结果存储到提供的目标结构体或变量中。
			使用场景：当您期望请求体是JSON格式，并且希望将其自动解析为结构化的数据时，例如RESTful API中的请求体。

			ctx.ReadBody()：该方法用于读取和解析请求体中的数据，并根据请求的Content-Type自动选择合适的解析器来解析请求体。
			使用场景：当您需要读取和解析请求体数据，但不确定请求的Content-Type是什么，或者希望自动选择合适的解析器时，例如处理多种类型的请求体数据。

			ctx.ReadForm(dest interface{}) 将请求的表单数据解析为键值对，并将其映射到给定结构体的字段上。（dest：要映射的结构体指针，必须是一个有效的结构体指针）
			*ctx.FormFile(fieldName string) (multipart.FileHeader, error) 获取上传文件的文件头信息。
			*ctx.FormFiles() ([]multipart.FileHeader, error) 获取所有上传文件的文件头信息。
		*/
	})

	// 设置错误等级
	// app.Logger().SetLevel("debug")
	// app.Logger().Info("ver", " : ", ver)

	/*
		Iris 有 7(+1) 种不同的方法来注册中间件，具体取决于您的应用程序的要求。
				app.WrapRouter(路由包装器)
			    app.UseRouter(路由中间件)
			    app.UseGlobal(全局中间件)
			    app.Use(使用中间件)
			    app.UseError(错误中间件)
				app.Done(完成)
				app.DoneGlobal(全局完成)
				// Adding a OnErrorCode(iris.StatusNotFound) causes `.UseGlobal`
			    // to be fired on 404 pages without this,
			    // only `UseError` will be called, and thus should
			    // be used for error pages.
			    app.OnErrorCode(iris.StatusNotFound, notFoundHandler)
	*/

	// 定制CORS中间件-实现服务端跨域(Our custom CORS middleware.)
	app.UseRouter(middle.MiddlewareCORS)

	/*
		使用`UseGlobal`去注册一个中间件，用于在所有子域名中使用 - 前置middleware中间件

		app.Use方法注册的中间件只会在当前路由组中生效，而不会在子路由组中生效。
		app.UseGlobal方法注册的中间件会在所有路由组中生效，包括子路由组。

			Point your browser to http://localhost:8080, and the output should look exactly like this:

			#1 .WrapRouter
			#2 .UseRouter
			#3 .UseGlobal
			#4 .Use
			Main Handler


			app.WrapRouter(routerWrapper)
		    app.UseRouter(routerMiddleware)
		    app.UseGlobal(globalMiddleware)
		    app.Use(useMiddleware)
		    app.UseError(errorMiddleware)
		    app.Done(done)
		    app.DoneGlobal(doneGlobal)
	*/

	// // 注册 "Done" ，在所有路由的处理程序之后调用
	// app.Done(func(ctx iris.Context) {})
	// app.UseGlobal(func(ctx iris.Context) {
	// 	ctx.Next()
	// })

	if env == "" {
		// 使用 app.Use() 函数来添加中间件
		// HTTP 压缩使用公共域压缩算法，如 gzip 和 brotli，在服务器上压缩 HTML、JavaScript、CSS 和其他文本文件。
		app.Use(iris.Compression) //压缩

		//（可选）添加两个内置处理程序
		// recover 和logger 是内建的中间件，帮助在 崩溃时记录和恢复
		app.Use(recover.New()) // 可以从任何http中进程退出时恢复（当服务panic时进程不挂掉）
		app.Use(logger.New())  // 将请求记录到终端。
	}

	// 从“./assets”文件夹加载所有扩展名为“.html”模板并解析它们，使用标准的“html/template”包。
	// 也可以为每上 Party 注册不同的视图引擎 例如  Application/Party.RegisterView(ViewEngine)
	// 设置为返回“assets”目录下的“.html文件”可能会被覆盖并隐藏错误，所以基本上只设置一次，
	// #文档 https://docs.iris-go.com/iris/view-templates/view
	app.RegisterView(iris.HTML(assetsPath+"/views/", ".html"))

	// 注册路由 (不推荐在main函数里写业务逻辑：采用低耦合高内聚的项目架构)
	app.Get("/", func(ctx iris.Context) {
		// // ctx.HTML("<h1>项目名</h1>")

		// // 绑定字符串 "项目名" 到 {{.message}} 变量
		// ctx.ViewData("message", name)
		// // 渲染模板文件 ./assets/views/index.html
		// ctx.View("index") // 或者写成 index.html 也行
		// // 如果有多层 ./assets/views/head/index.html 则指定要返还的view文件路径 ctx.View("head/index")

		// // ctx.JSON(iris.Map{"Author": "肥客泉", "status": iris.StatusOK})
		// // ctx.WriteString("你没有权限")

		// 读取文件内容
		data, err := os.ReadFile("./README.md")
		if err != nil {
			println("读取文件失败：", err.Error())
			return
		}

		// 输出文件内容
		ctx.Markdown(data)

	})

	// 404
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		println("[404] " + ctx.Path())

		// println("iris.StatusNotFound:", iris.StatusNotFound)
		// ctx.HTML("<b>Resource Not found</b>")
		// ctx.WriteString("404")

		// 出现 404 的时候，就跳转到 ./assets/views/404.htm 模板
		ctx.View("[404] " + ctx.Path())
	})
	// 500
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		println("iris.StatusInternalServerError:", iris.StatusInternalServerError)
		ctx.WriteString("500")
	})

	// 注册鉴权中间件 （在这之前的路由不进行权限鉴别）
	app.Use(middle.MiddlewareAuthToken)

	// 初始化数据库
	db, err := config.ConnectDB(otherCfg)
	if err != nil {
		// println("Failed to connect to MySQL database: %v", err)
		// return

		// 在Go语言中，使用多值返回来返回错误。不要用异常代替错误，更不要用来控制流程。
		// Go语言不支持传统的 try…catch…finally 这种异常，因为Go语言的设计者们认为，将异常与控制结构混在一起会很容易使得代码变得混乱。
		panic(err) // 手动触发宕机停止常规的goroutine（panic和recover：用来做错误处理）
	}
	println("数据库连接成功内存地址为:", &db, "->", db)

	// defer是Go语言中的延迟执行语句。
	// 用来添加函数结束时执行的代码，常用于释放某些已分配的资源、关闭数据库连接、断开socket连接、解锁一个加锁的资源。
	defer db.Close() // 应用退出时关闭数据库连接（Go语言机制担保一定会执行defer语句中的代码）

	// 两个{}只是把相同路由组的放在一个区块，没有其他用特殊含义
	{
		// 注册控制器 （创建控制器实例，并将数据库连接对象传递给控制器）

		// 通过调取工厂函数进行初始化实例（简洁和易于扩展），工厂函数它封装了创建实例的逻辑，可以在内部处理控制器的初始化过程
		userController := controllers.NewUserController(db)
		oauthController := controllers.NewOauthController(db)
		accountController := controllers.NewAccountController(db)
		contactController := controllers.NewContactController(db)
		remindController := controllers.NewRemindController(db)
		humaneController := controllers.NewHumaneController(db)
		notepadController := controllers.NewNotepadController(db)

		// 通过赋值了一个xxxxController实例，直接创建实例并手动传入 DB 对象
		commonController := &controllers.CommonController{} // 不需要模型简单赋值调用（当然也可以传值进去像 CommonController{DB: db} ）

		// 使用MVC模式组织路由和控制器

		// .Register(userDB) 对于每个控制器，我们都注册了相应的服务，并将控制器实例作为路由处理函数。（ mvc.应用程序 Register 方法注册一个依赖项）
		// userApp := mvc.New(app.Party("/user"))
		// userApp.Register(userController)
		// userApp.Handle(new(controller.UserController))

		// MVC 应用程序有自己的 Router ，它是一个 iris/router.Party 类型的标准 iris api。
		// 当 iris/router.Party 如期望的那样开始执行处理程序的时候 ， Controllers 可以注册到任意的 Party，包括子域。
		// 基于/user提供mvc服务，如果只填 app 则基于根路由器“/”为控制器提供服务。
		// ( gin的分组路由注册仍需一条条注册，而iris-mvc将子路由路径托管给controller的方法名。在易用上irismvc的方式更加简便 )
		// 将控制器传递给每个路由处理程序
		mvc.New(app.Party("/user")).Handle(userController) // 用户中心
		// // 如果想针对指定模块注册中间件就使用.Use()即可，不过通过控制器内的 BeforeActivation 注册更优雅
		// {
		// 	userApp := app.Party("/user")
		// 	userApp.Use(middle.MiddlewareVerifyAdmin)
		// 	mvc.New(userApp).Handle(userController)
		// }
		// // 也可以像这样设置中间件,将中间件移动到了mvcApp对象上
		// {
		// 	mvcApp := mvc.New(app.Party("/user"))
		// 	mvcApp.Router.Use(userController)
		// }
		mvc.New(app.Party("/common")).Handle(commonController) // 通用功能
		mvc.New(app.Party("/oauth")).Handle(oauthController)   // 平台接入
		mvc.New(app.Party("/account")).Handle(accountController)
		mvc.New(app.Party("/contact")).Handle(contactController)
		mvc.New(app.Party("/humane")).Handle(humaneController)
		mvc.New(app.Party("/notepad")).Handle(notepadController)
		mvc.New(app.Party("/remind")).Handle(remindController)

	}

	println() //只是为了让控制台换一行，不进行实质性的内容输出

	// 启动 HTTP 服务器

	if env != "" {
		// 开发环境和测试环境中可能更倾向于使用 app.Run
		// app.Run 是一个阻塞方法，它会在启动 HTTP 服务器后阻止程序继续向下执行。这对于开发阶段的快速测试和调试非常方便，因为你可以在服务器运行时实时查看结果，并进行必要的更改和调试。
		// app.Run(iris.Addr(":8888"), iris.WithCharset("UTF-8")) // 在Tcp上监听网络地址 0.0.0.0:8888
		devErr := app.Run(iris.Addr(addr), iris.WithConfiguration(cfg))
		if devErr != nil {
			println("devErr", devErr.Error())
			app.Run(iris.Addr(":0"), iris.WithConfiguration(cfg)) // 随机生成一个可用的端口号？
		}

	} else {
		// 生产环境中则更倾向于使用 app.Listen
		// app.Listen 是一个非阻塞方法，它允许你在服务器运行的同时执行其他操作。这对于生产环境非常重要，因为你可能需要进行其他初始化或设置操作，或者并行处理其他任务，而无需等待服务器运行完毕。
		// 允许你显式地指定要绑定的网络地址，并提供更多的服务器选项。这样你可以根据需要进行自定义配置，如绑定到特定的 IP 地址和端口、启用 HTTPS、设置超时时间、启用自定义的中间件等。这对于生产环境和高级配置需求非常重要。
		// Listen 是 app.Run(iris.Addr(":8080")) 的快捷方式。
		serErr := app.Listen(addr, iris.WithConfiguration(cfg))
		if serErr != nil {
			println("serErr", serErr.Error())
			app.Listen(":0", iris.WithConfiguration(cfg)) // 随机生成一个可用的端口号？
		}

	}

}
