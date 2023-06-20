package middle

import (
	"tibiji-go/config"
	"tibiji-go/utils"

	"github.com/kataras/iris/v12"
)

const IS_PASSABLE = false // 是否放行开关

// 密钥验证的中间件
func MiddlewareAuthToken(ctx iris.Context) {
	if IS_PASSABLE {
		ctx.Next()
		return
	}
	println("--------- MiddlewareAuthToken ---------")

	// 拿到配置和头信息
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	ua := ctx.GetHeader("User-Agent") // 拿到UA信息User-Agent
	// exptime := otherCfg["SERV_EXPIRES_TIME"].(int64)
	secret := otherCfg["SERV_KEY_SECRET"].(string) + ua
	// println("secret:", secret)
	// println("UA ----")

	// 获取 Basic Auth 认证头信息
	// token := ctx.Request().Header.Get("Authorization")
	token := ctx.GetHeader("Authorization")

	// 进行 Basic Auth 身份认证
	token_uid, token_err := utils.VerifyToken(token, secret)
	// 将当前环境添加到上下文中(将值传递给下一个处理程序)
	// ctx.Values().Set() 值是处理程序（或中间件）在彼此之间进行通信的方式，
	// 获取时使用 ctx.Values().Get("xxxx")内存地址,当然也可以 GetString 或 GetInt64 拿实际值
	ctx.Values().Set("UID", token_uid)

	// 白名单 定义不需要进行鉴权的URL列表
	// excludeLists := []string{"GET:/user", "POST:/user"}
	// currentTag := ctx.Method() + ":" + ctx.Path()

	// 白名单 定义不需要进行鉴权的FUN列表
	excludeLists := []string{
		"controllers.UserController.Post",           // 添加用户
		"controllers.UserController.PostLogin",      // 用户登录
		"controllers.OauthController.Post",          // 接入用户
		"controllers.OauthController.AllWxBy",       // 微信接入
		"controllers.OauthController.AllWxJsticket", // 微信JS-SDK
		"controllers.CommonController.GetEmail",     // 邮件发送测试
		"controllers.UserController.GetPassword",    // 找回密码
		"controllers.UserController.PostPassword",   // 找回密码后设置新密码
		"controllers.AccountController.GetType",     // 帐单类目
		"controllers.RemindController.GetTask",      // 提醒任务
		"controllers.RemindController.GetQueue",     // 提醒队列
		"controllers.NotepadController.GetBy",       // 获取纸张
		"controllers.NotepadController.GetShareBy",  // 分享纸张
	}
	currentTag := ctx.GetCurrentRoute().MainHandlerName()

	// 判断当前请求的URL是否需要进行鉴权
	needAuth := true
	for _, p := range excludeLists {
		// println(p, currentTag)
		if p == currentTag {
			needAuth = false
			break
		}
	}

	// 进行用户校验的逻辑
	if !needAuth {
		println(currentTag, "放行")
		// 如果是白名单列表的接口，则不需要进行鉴权
		ctx.Next()
		return
	} else {
		println(currentTag, "正在进行权鉴...")
	}

	// println(secret, "校验标识：", currentTag)

	// println(token)
	if token == "" {
		println(iris.StatusUnauthorized, config.ErrMsgs[config.ErrHeader])
		// 未提供认证头，返回 401 Unauthorized
		ctx.StatusCode(iris.StatusUnauthorized)
		// ctx.WriteString("unauthorized")
		ctx.JSON(iris.Map{"code": config.ErrHeader, "msg": config.ErrMsgs[config.ErrHeader]})
		return
	}
	// println(token)

	if token_err != nil {
		println("VerifyToken Error: ", token_err.Error())
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"code": config.ErrToken, "msg": config.ErrMsgs[config.ErrToken]})
		return
	}

	// 校验用户名和密码
	if token_uid == 0 {
		// 认证失败，返回 401 Unauthorized
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"code": config.ErrToken, "msg": config.ErrMsgs[config.ErrToken]})
		return
	}

	// 认证通过，继续处理请求
	ctx.Next()
}

// 管理身份验证的中间件
func MiddlewareVerifyAdmin(ctx iris.Context) {
	if IS_PASSABLE {
		ctx.Next()
		return
	}
	println("--------- MiddlewareVerifyAdmin ---------")

	// 认证通过，继续处理请求
	ctx.Next()
}
