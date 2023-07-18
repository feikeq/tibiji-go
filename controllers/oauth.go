package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type OauthController struct {
	DB     *sqlx.DB
	Models *models.UserModel
	CTX    iris.Context
}

func NewOauthController(db *sqlx.DB) *OauthController {
	// 返回一个结构体指针
	return &OauthController{
		DB:     db,
		Models: models.NewUserModel(db),
	}
}

// 微信缓存
type WXCach struct {
	AccessToken string `description:"获取到的凭证"`
	TokenExp    int64  `description:"凭证有效时间"`
	JsapiTicket string `description:"JS-SDK接口的签名票据"`
	TicketExp   int64  `description:"临时票据时效"`
}

var wx = WXCach{} // 全局变量

// 接入查询 GET:/oauth
func (c *OauthController) Get() {
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

	var platfrom, openid_unionid string

	// 判断是否存在字段 "platfrom"
	if _, ok := allData["platfrom"]; !ok {
		println("外接平台名不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["platfrom"] == "" {
			println("外接平台名不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		platfrom = allData["platfrom"].(string)
	}
	// 判断是否存在字段 "openid_unionid"
	if _, ok := allData["openid_unionid"]; !ok {
		println("外接平台身份ID或唯一标识不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["openid_unionid"] == "" {
			println("外接平台身份ID或唯一标识不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		openid_unionid = allData["openid_unionid"].(string)
	}

	useroauth, err := c.Models.FindOAuthOpenid(platfrom, openid_unionid)
	if err == nil {
		ctx.JSON(iris.Map{"data": useroauth, "code": 0, "msg": ""})
		return
	}

	// 如果找到 unionid 相符的用户
	useroauth, _ = c.Models.FindOAuthUnionid(platfrom, openid_unionid)

	ctx.JSON(iris.Map{"data": useroauth, "code": 0, "msg": ""})

}

// 手动接入 POST:/oauth
func (c *OauthController) Post() {
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
	// fmt.Printf("allData %T -> %v", allData, allData)

	var platfrom, openid string

	// 判断是否存在字段 "platfrom"
	if _, ok := allData["platfrom"]; !ok {
		println("外接平台名不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["platfrom"] == "" {
			println("外接平台名不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		platfrom = allData["platfrom"].(string)
	}

	// 判断是否存在字段 "openid"
	if _, ok := allData["openid"]; !ok {
		println("外接平台名不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["openid"] == "" {
			println("外接平台名不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		openid = allData["openid"].(string)
	}
	println(platfrom, openid)

	// 删除不能修改的字段
	// delete(allData, "uid")    // 删除 用户ID
	delete(allData, "uptime") // 删除 更新时间
	allData["uid"] = tkUid

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	useroauth, err := c.Models.CreateOAuth(allData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户返回
	if useroauth.UID > 0 {

		// 获取用户信息
		user, err := c.Models.Read(useroauth.UID)
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
			}
			return
		}

		// 在控制器层将结果进行修改和脱敏并得到最终的数据
		if *user.State == 2 {
			if env != "" {
				println("帐号还未激活")
				ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": "帐号还未激活"})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
			}
			return
		} else if *user.State == 0 {
			if env != "" {
				println("帐号已被禁用")
				ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": "帐号已被禁用"})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrUserDisabled, "msg": config.ErrMsgs[config.ErrUserDisabled]})
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

		// 获取配置项
		otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
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
			"note":   platfrom,
			"actip":  utils.GetRealIP(ctx),
			"ua":     ua,
		}
		log := c.Models.SetLogs(logData)
		if log != nil {
			if env != "" {
				println("Models.SetLogs Error: ", err.Error())
				ctx.JSON(iris.Map{"data": logData, "code": "err debug", "msg": err.Error()})
				return
			}
		}
		// 返回登录状态
		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
		return
	}

	// 返回开放平台数据
	ctx.JSON(iris.Map{"data": useroauth, "code": 0, "msg": ""})
}

// 手动绑定 PUT:/oauth
func (c *OauthController) Put() {
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
	// fmt.Printf("allData %T -> %v", allData, allData)

	var oid int64

	// 判断是否存在字段 "oid"
	if _, ok := allData["oid"]; !ok {
		println("oid不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["oid"] == "" {
			println("oid不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		oid = utils.ParseInt64(allData["oid"]) // 任意数据转int64数字
	}

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	useroauth, err := c.Models.FindOAuthOid(oid)
	if err != nil {
		if env != "" {
			println("Models.FindOAuthOid Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户则提示非法操作
	if useroauth.UID > 0 {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	allData["uid"] = tkUid

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.UpdateOAuth(oid, allData)
	if err != nil {
		if env != "" {
			println("Models.UpdateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    tkUid,
		"action": "bind",
		"note":   useroauth.Platfrom,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", err.Error())
			ctx.JSON(iris.Map{"data": logData, "code": "err debug", "msg": err.Error()})
			return
		}
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 手动解绑 DELETE:/oauth
func (c *OauthController) Delete() {
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
	fmt.Printf("allData: %+v\n", allData) // 打印allData

	var oid int64

	// 判断是否存在字段 "oid"
	if _, ok := allData["oid"]; !ok {
		println("oid不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["oid"] == "" {
			println("oid不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		oid = utils.ParseInt64(allData["oid"]) // 任意数据转int64数字
	}

	useroauth, err := c.Models.FindOAuthOid(oid)
	if err != nil {
		if env != "" {
			println("Models.FindOAuthOid Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果操作人不是自己
	if tkUid != useroauth.UID {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	allData["uid"] = 0

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.UpdateOAuth(oid, allData)
	if err != nil {
		if env != "" {
			println("Models.UpdateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    tkUid,
		"action": "unbind",
		"note":   useroauth.Platfrom,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
	}
	log := c.Models.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", err.Error())
			ctx.JSON(iris.Map{"data": logData, "code": "err debug", "msg": err.Error()})
			return
		}
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})

}

// 微信接入 - 网页授权任意请求类型访问 POST:/oauth/wx/{code}
// 微信授权 - 网页授权任意请求类型访问 GET:/oauth/wx/{scope}
func (c *OauthController) AllWxBy(code string) {
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

	// 微信网页登录网页授权文档
	// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html

	method := ctx.Method()

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)

	oauthData := make(map[string]interface{})
	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range allData {
		oauthData[key] = value
	}

	// 获取配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	wx_appid := otherCfg["WX_APPID"].(string)         // 公众号的唯一标识
	wx_appsecret := otherCfg["WX_APPSECRET"].(string) // 公众号的appsecret
	wx_state := "STATE"
	wx_agentid := "" // 企业微信 应用agentid

	// 判断是否存在字段 "appid"
	if _, ok := allData["appid"]; ok {
		wx_appid = allData["appid"].(string)
	}
	// 判断是否存在字段 "appsecret"
	if _, ok := allData["appsecret"]; ok {
		wx_appsecret = allData["appsecret"].(string)
	}
	// 判断是否存在字段 "state"
	if _, ok := allData["state"]; ok {
		wx_state = allData["state"].(string)
	}

	// 判断是否存在字段 "agentid"
	if _, ok := allData["agentid"]; ok {
		wx_agentid = allData["agentid"].(string)
	}

	if method == "GET" {
		// 用户同意授权，获取code
		println("ctx.Method()", method)
		var wx_redirect_uri string
		if _, ok := allData["redirect_uri"]; ok {
			wx_redirect_uri = allData["redirect_uri"].(string)
		}
		scope := "snsapi_base"
		if code == "snsapi_userinfo" {
			scope = "snsapi_userinfo"
		} else if code == "snsapi_privateinfo" {
			scope = "snsapi_privateinfo"
		}
		// redirect_uri 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
		wx_redirect_uri = url.QueryEscape(wx_redirect_uri)
		// 企业微信授权会多一个 agentid 参数 当oauth2中appid=corpid时，scope为snsapi_userinfo或snsapi_privateinfo时，必须填agentid参数，否则系统会视为snsapi_base，不会返回敏感信息。
		codeUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s", wx_appid, wx_redirect_uri, scope, wx_state)
		// 如果有企业微信的应用agentid
		if wx_agentid != "" {
			codeUrl += fmt.Sprintf("&agentid=%s#wechat_redirect", wx_agentid)
		} else {
			codeUrl += "#wechat_redirect"
		}
		ctx.JSON(iris.Map{"data": codeUrl, "code": 0, "msg": ""})
		return
	}

	// 通过code换取网页授权access_token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", wx_appid, wx_appsecret, code)
	access, err := utils.HTTPDo("GET", url, map[string][]string{})
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("access %T -> %v", access, access)

	var access_token string
	// 判断是否存在字段 "access_token"
	if _, ok := access["access_token"]; ok {
		access_token = access["access_token"].(string)
		oauthData["token"] = access_token
	}
	if access_token == "" {
		ctx.JSON(iris.Map{"data": allData, "code": access["errcode"], "msg": access["errmsg"]})
		return
	}
	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range access {
		oauthData[key] = value
	}

	// 拉取用户信息(需scope为 snsapi_userinfo)
	url = fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", access_token, access["openid"].(string))
	info, err := utils.HTTPDo("GET", url, map[string][]string{})
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("info %T -> %v", info, info)
	var openid string
	// 判断是否存在字段 "appid"
	if _, ok := info["openid"]; ok {
		openid = info["openid"].(string)
	}
	if openid == "" {
		ctx.JSON(iris.Map{"data": allData, "code": info["errcode"], "msg": info["errmsg"]})
		return
	}
	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range info {
		oauthData[key] = value
	}

	// 添加其他必要字段
	oauthData["uid"] = 3
	oauthData["platfrom"] = "wx2"

	// 转义对应字段名
	if _, ok := info["headimgurl"]; ok {
		oauthData["headimg"] = info["headimgurl"].(string)
	}
	if _, ok := access["refresh_token"]; ok {
		oauthData["refresh"] = access["refresh_token"].(string)
	}
	if _, ok := access["expires_in"]; ok {
		oauthData["expires"] = int64(access["expires_in"].(float64))
	}

	// 判断是否存在字段 "privilege" 转为字符串
	if _, ok := oauthData["privilege"]; ok {
		// 判断类型
		switch oauthData["privilege"].(type) {
		case string:
			oauthData["privilege"] = oauthData["privilege"].(string)
		default:
			oauthData["privilege"], _ = json.Marshal(oauthData["privilege"])
		}
	}

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	id, err := c.Models.CreateOAuth(oauthData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"data": oauthData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": id, "code": 0, "msg": ""})
}

// 微信签名票据 GET:/oauth/wx/jsticket
func (c *OauthController) GetWxJsticket() {
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

	// 定义获取当前时间
	now := time.Now()

	// 获取配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	wx_appid := otherCfg["WX_APPID"].(string)         //公众号的唯一标识
	wx_appsecret := otherCfg["WX_APPSECRET"].(string) // 公众号的appsecret

	var url string

	// println("wx_token:", wx.AccessToken, "\nwx_expires:", wx.TokenExp, "\nwx_ticket:", wx.JsapiTicket, "\nwx_ticket_exp:", wx.TicketExp)

	// 判断是否存在字段 "appid"
	if _, ok := allData["appid"]; ok {
		wx_appid = allData["appid"].(string)
	}
	// 判断是否存在字段 "appid"
	if _, ok := allData["appsecret"]; ok {
		wx_appsecret = allData["appsecret"].(string)
	}
	// 判断是否存在字段 "appid"
	if _, ok := allData["url"]; ok {
		url = allData["url"].(string)
	}

	now_unix := now.Unix()

	// 检查过期时间是否超时（当前时间是否大于ticket有效时间）
	if now_unix > wx.TicketExp {

		// 检查过期时间是否超时（当前时间是否大于ticket有效时间）
		if now_unix > wx.TokenExp {

			// 获取 Access token
			asurl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", wx_appid, wx_appsecret)
			access, err := utils.HTTPDo("GET", asurl, map[string][]string{})
			if err != nil {
				fmt.Println(err.Error())
			}
			// println(asurl)
			// fmt.Printf("access %T -> %v", access, access)

			// 判断是否存在字段 "access_token"
			if _, ok := access["access_token"]; ok {
				wx.AccessToken = access["access_token"].(string) // 全局缓存
			}
			if wx.AccessToken == "" {
				ctx.JSON(iris.Map{"data": allData, "code": access["errcode"], "msg": access["errmsg"]})
				return
			}

			if _, ok := access["expires_in"]; ok {
				wx.TokenExp = int64(access["expires_in"].(float64))
				wx.TokenExp = now.Add(time.Second * time.Duration(wx.TokenExp)).Unix() // 全局缓存
			}

		} else {
			println("cach_wx_token:", wx.AccessToken, "\ncach_wx_expires:", wx.TokenExp)

		}

		// 获得jsapi_ticket
		{

			// 获得jsapi_ticket（有效期7200秒，开发者必须在自己的服务全局缓存jsapi_ticket）
			tkurl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", wx.AccessToken)
			tk, tkerr := utils.HTTPDo("GET", tkurl, map[string][]string{})
			if tkerr != nil {
				fmt.Println(tkerr.Error())
			}
			// fmt.Printf("ticket %T -> %v", tk, tk)

			// 判断是否存在字段 "ticket"
			if _, ok := tk["ticket"]; ok {
				wx.JsapiTicket = tk["ticket"].(string) // 全局缓存
			}
			if wx.JsapiTicket == "" {
				ctx.JSON(iris.Map{"data": allData, "code": tk["errcode"], "msg": tk["errmsg"]})
				return
			}
			if _, ok := tk["expires_in"]; ok {
				wx.TicketExp = int64(tk["expires_in"].(float64))
				wx.TicketExp = now.Add(time.Second * time.Duration(wx.TicketExp)).Unix() // 全局缓存
			}

		}
	} else {
		println("\ncach_wx_ticket:", wx.JsapiTicket, "\ncach_wx_ticket_exp:", wx.TicketExp)

	}

	// 生成JS-SDK权限验证的签名
	type Sign struct {
		AppID     string `json:"appId" description:"必填，公众号的唯一标识"`
		Timestamp int64  `json:"timestamp" description:"必填，生成签名的时间戳"`
		NonceStr  string `json:"nonceStr" description:"必填，生成签名的随机串"`
		Signature string `json:"signature" description:"必填，签名"`
	}

	signature := Sign{
		AppID:     wx_appid,
		Timestamp: time.Now().Unix(),
		NonceStr:  "FK68TaJiuShiYYSD",
		Signature: "",
	}

	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", wx.JsapiTicket, signature.NonceStr, signature.Timestamp, url)
	// println(str)
	signature.Signature = utils.CalculateSHA1(str)
	ctx.JSON(iris.Map{"data": signature, "code": 0, "msg": ""})
}
