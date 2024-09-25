package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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

func NewOauthController(db *sqlx.DB, cfg map[string]interface{}) *OauthController {
	// 返回一个结构体指针
	return &OauthController{
		DB:     db,
		Models: models.NewUserModel(db, cfg),
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

// 绑定查询 GET:/oauth/
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

	// 调取模型 - 根据参数获取数据库中的信息
	oauths, err := c.Models.FindOAuth(tkUid)
	if err != nil {
		ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		return
	}

	// fmt.Printf("types %T -> %v", types, types)
	ctx.JSON(iris.Map{"data": oauths, "code": 0, "msg": ""})
}

// 接入查询 GET:/oauth/{openid_unionid}/
func (c *OauthController) GetBy(openid_unionid string) {
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

	var platfrom string

	// 判断是否存在字段 "platfrom"
	if _, ok := allData["platfrom"]; !ok {
		println("外接平台身份ID或唯一标识不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["platfrom"] == "" {
			println("外接平台身份ID或唯一标识不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
		platfrom = allData["platfrom"].(string)
	}

	oUser, err := c.Models.FindOAuthOpenid(platfrom, openid_unionid)
	if err == nil {
		ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})
		return
	}

	// 如果找到 unionid 相符的用户
	oUser, _ = c.Models.FindOAuthUnionid(platfrom, openid_unionid)

	ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})

}

// 手动接入 POST:/oauth/
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
	oUser, err := c.Models.CreateOAuth(allData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户返回
	if oUser.UID > 0 {

		// 获取用户信息
		user, err := c.Models.Read(oUser.UID)
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
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
				println("Models.SetLogs Error: ", log.Error())
				ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
				return
			}
		}
		// 返回登录状态
		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
		return
	}

	// 返回开放平台数据
	ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})
}

// 手动绑定 PUT:/oauth/
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
	oUser, err := c.Models.FindOAuthOid(oid)
	if err != nil {
		if env != "" {
			println("Models.FindOAuthOid Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户则提示非法操作
	if oUser.UID > 0 {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	allData["uid"] = tkUid

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.UpdateOAuth(oid, allData)
	if err != nil {
		if env != "" {
			println("Models.UpdateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    tkUid,
		"action": "bind",
		"note":   oUser.Platfrom,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
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

// 手动解绑 DELETE:/oauth/{oid}/
func (c *OauthController) DeleteBy(oid int64) {
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

	oUser, err := c.Models.FindOAuthOid(oid)
	if err != nil {
		if env != "" {
			println("Models.FindOAuthOid Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果操作人不是自己
	if tkUid != oUser.UID {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	//  UNIQUE KEY `platfrom_openid` (`platfrom`,`openid`),
	allData["uid"] = 0 // 下次来再绑定即可

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.UpdateOAuth(oid, allData)
	if err != nil {
		if env != "" {
			println("Models.UpdateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 操作写入日志表
	logData := map[string]interface{}{
		"uid":    tkUid,
		"action": "unbind",
		"note":   oUser.Platfrom,
		"actip":  utils.GetRealIP(ctx),
		"ua":     ctx.GetHeader("User-Agent"), // 拿到UA信息User-Agent
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

	var platfrom string

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
	access, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
	if err != nil {
		fmt.Println(url, err)
		// ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode], "data": err.Error()})
		// return
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
	info, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
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
	oauthData["uid"] = tkUid
	oauthData["platfrom"] = platfrom

	// 转义对应字段名
	if _, ok := info["headimgurl"]; ok {
		oauthData["headimg"] = info["headimgurl"].(string)
	}
	if _, ok := access["refresh_token"]; ok {
		oauthData["refresh"] = access["refresh_token"].(string)
	}
	if _, ok := access["expires_in"]; ok {
		// oauthData["expires"] = int64(access["expires_in"].(float64))
		var _tmp = int64(access["expires_in"].(float64))
		oauthData["expires"] = strconv.FormatInt(_tmp, 10) // 过期时间转为10进制数字的字符串 数字转字符串
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
	oUser, err := c.Models.CreateOAuth(oauthData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": oauthData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户返回
	if oUser.UID > 0 {
		// 获取用户信息
		user, err := c.Models.Read(oUser.UID)
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
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
				println("Models.SetLogs Error: ", log.Error())
				ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
				return
			}
		}
		// 返回登录状态
		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
		return
	}

	ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})
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
			access, err := utils.HTTPDo("GET", asurl, map[string][]string{}, map[string]string{})
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
			tk, tkerr := utils.HTTPDo("GET", tkurl, map[string][]string{}, map[string]string{})
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

// 微博接入 - 网页授权任意请求类型访问 POST:/oauth/wb/{code}/
// 微博授权 - 网页授权任意请求类型访问 GET:/oauth/wb/{scope}/
func (c *OauthController) AllWbBy(code string) {
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
	// 新浪微博 https://open.weibo.com/
	// 微博登录文档 https://open.weibo.com/wiki/Connect/login

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
	app_key := otherCfg["WB_APP_KEY"].(string)       // 第三方应用在微博开放平台注册的APPKEY。
	app_secret := otherCfg["WB_APP_SECRET"].(string) // 在微博开放平台注册的应用所对应的AppSecret。
	// 通过 ctx.Request().Referer() 方法获取用户跳转过来的页面域名
	// 以使用ctx.Request().Header.Get("Referer") 来获取用户访问接口时用户当前的网页完整地址包括路径和参数。
	// 通过 ctx.Request().URL.String() 来获取用户访问API接口时的完整地址包括路径和参数
	redirect_uri := ctx.GetHeader("Referer") // 授权回调地址，传的值需与在开放平台网站填写的回调地址一致，设置填写位置：“我的应用>应用信息>高级信息”。
	grant_type := "authorization_code"       // 请求的类型，需填写 authorization_code。
	wb_state := "weibo"                      //用于保持请求和回调的状态，在回调时，会在Query Parameter中回传该参数。开发者可以用这个参数验证请求有效性，也可以记录用户请求授权页前的位置，这个参数可用于防止跨站请求伪造（CSRF）攻击。
	// println("redirect_uri ", redirect_uri)

	platfrom := "新浪微博" // 缺省平台名

	// 判断是否存在字段 "platfrom"
	if _, ok := allData["platfrom"]; ok {
		platfrom = allData["platfrom"].(string)
	}
	// 判断是否存在字段 "appid"
	if _, ok := allData["appid"]; ok {
		app_key = allData["appid"].(string)
	}
	// 判断是否存在字段 "appsecret"
	if _, ok := allData["appsecret"]; ok {
		app_secret = allData["appsecret"].(string)
	}
	// 判断是否存在字段 "state"
	if _, ok := allData["state"]; ok {
		wb_state = allData["state"].(string)
	}

	// 回调地址
	if _, ok := allData["redirect_uri"]; ok {
		redirect_uri = allData["redirect_uri"].(string)
	}

	// 微博好像只有一种 authorization_code
	if code == "snsapi_userinfo" {
		grant_type = "authorization_code"
	} else if code == "snsapi_privateinfo" {
		grant_type = "authorization_code"
	}

	if method == "GET" {
		// 用户同意授权，获取code
		println("ctx.Method()", method)

		// redirect_uri 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
		redirect_uri = url.QueryEscape(redirect_uri)
		// 企业微信授权会多一个 agentid 参数 当oauth2中appid=corpid时，scope为snsapi_userinfo或snsapi_privateinfo时，必须填agentid参数，否则系统会视为snsapi_base，不会返回敏感信息。
		codeUrl := fmt.Sprintf("https://api.weibo.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s", app_key, redirect_uri, wb_state)

		ctx.JSON(iris.Map{"data": codeUrl, "code": 0, "msg": ""})
		return
	}

	// -------------------------- post -------------------------

	/*
		https://open.weibo.com/wiki/Oauth2/access_token
		返回字段说明
		返回字段	字段类型	字段说明
		access_token	string	用户授权的唯一票据，用于调用微博的开放接口，同时也是第三方应用验证微博用户登录的唯一票据，第三方应用应该对该票据进行校验，校验方法为调用 oauth2/get_token_info 接口，对比返回的授权信息中的APPKEY是否正确一致，然后用 access_token 与自己应用内的用户建立唯一影射关系，来识别登录状态，不能只是简单的使用本返回值里的UID字段来做登录识别。
		expires_in	string	access_token 的生命周期，单位是秒数。
		uid	string	授权用户的UID，本字段只是为了方便开发者，减少一次 user/show 接口调用而返回的，第三方应用不能用此字段作为用户登录状态的识别，只有 access_token 才是用户授权的唯一票据。
	*/

	println("OAuth2授权第二步access_token接口")

	// 通过code换取网页授权access_token
	params := url.Values{
		"client_id":     {app_key},
		"client_secret": {app_secret},
		"grant_type":    {grant_type},
		"code":          {code},
		"redirect_uri":  {redirect_uri},
	}

	url := "https://api.weibo.com/oauth2/access_token"
	var headers = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	access, err := utils.HTTPDo("POST", url, params, headers)
	if err != nil {
		fmt.Println(url, err)
		ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode], "data": err.Error()})
		return
	}
	fmt.Printf("access: %T -> %v", access, access)

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
	println("access_token:", access_token)
	if _, ok := access["refresh_token"]; ok {
		oauthData["refresh"] = access["refresh_token"].(string)
	}
	var openid string
	if _, ok := access["uid"]; ok {
		openid = access["uid"].(string)
	}
	if _, ok := access["expires_in"]; ok {
		var _tmp = int64(access["expires_in"].(float64))
		oauthData["expires"] = strconv.FormatInt(_tmp, 10) // 过期时间转为10进制数字的字符串 数字转字符串
	}
	println("openid:", openid)

	// // 获取用户提交的所有表单项字段 遍历数据
	// for key, value := range access {
	// 	oauthData[key] = value
	// }

	// 根据用户UID或昵称获取用户资料 https://open.weibo.com/wiki/2/users/show
	println("根据用户UID或昵称获取用户资料")
	// 拉取用户信息(需scope为 snsapi_userinfo)
	url = fmt.Sprintf("https://api.weibo.com/2/users/show.json?access_token=%s&uid=%s", access_token, openid)
	println("url:", url)
	info, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
	if err != nil {
		fmt.Println(err) //拉不到信息也要往下走因为有openid了
	} else {

		println("拉取用户信息")
		fmt.Printf("info  %T -> %v", info, info)
		var _openid string
		// 判断是否存在字段 "id"
		// if _, ok := info["id"]; ok {
		// 	openid = info["id"].(string)
		// 	if _openid != openid {
		// 		println("openid1:", _openid)
		// 		println("openid2:", openid)
		// 	}
		// }
		if _, ok := info["idstr"]; ok {
			_openid = info["idstr"].(string)
			if _openid != openid {
				openid = _openid
				println("openid1:", openid)
				println("openid2:", _openid)
			}
		}

		if openid == "" {
			ctx.JSON(iris.Map{"data": allData, "code": config.ErrNotFound, "msg": config.ErrMsgs[config.ErrNotFound]})
			return
		}
		println("获取用户提交的所有表单项字段")

		println("遍历数据")
		for key, value := range info {
			oauthData[key] = value
		}

		// 转义对应字段名
		if gender, ok := info["gender"]; ok {
			oauthData["sex"] = gender.(string)
			if oauthData["sex"] == "m" {
				oauthData["sex"] = 1
			} else if oauthData["sex"] == "f" {
				oauthData["sex"] = 2
			} else {
				oauthData["sex"] = 0
			}
		}
		if _, ok := info["profile_image_url"]; ok {
			oauthData["headimg"] = info["profile_image_url"].(string)
		}
		if _, ok := info["screen_name"]; ok {
			oauthData["nickname"] = info["screen_name"].(string)
		}
		if _, ok := info["domain"]; ok {
			oauthData["unionid"] = info["domain"].(string)
		}
		if _, ok := info["description"]; ok {
			oauthData["remark"] = info["description"].(string)
		}
		if _, ok := info["province"]; ok {
			oauthData["province"] = info["province"].(string)
		}
		if _, ok := info["location"]; ok {
			oauthData["city"] = info["location"].(string)
		}
		if _, ok := info["city"]; ok {
			oauthData["city"] = oauthData["city"].(string) + " " + info["city"].(string)
		}

		followers_count := info["followers_count"].(float64)
		friends_count := info["friends_count"].(float64)
		statuses_count := info["statuses_count"].(float64)
		favourites_count := info["favourites_count"].(float64)
		oauthData["privilege"] = map[string]float64{
			"followers_count":  followers_count,
			"friends_count":    friends_count,
			"statuses_count":   statuses_count,
			"favourites_count": favourites_count,
		}

		// 判断是否存在字段 "privilege" 转为字符串
		if _, ok := oauthData["privilege"]; ok {
			// 判断类型
			switch oauthData["privilege"].(type) {
			case string:
				oauthData["privilege"] = oauthData["privilege"].(string)
			default:
				// 将 map 转换为 JSON 字符串
				oauthData["privilege"], _ = json.Marshal(oauthData["privilege"])
			}
		}

	}

	println("添加其他必要字段")
	oauthData["uid"] = tkUid
	oauthData["openid"] = openid
	oauthData["platfrom"] = platfrom

	fmt.Printf("oauthData  %T -> %v", oauthData, oauthData)

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	oUser, err := c.Models.CreateOAuth(oauthData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": oauthData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户返回
	if oUser.UID > 0 {
		// 获取用户信息
		user, err := c.Models.Read(oUser.UID)
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
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
				println("Models.SetLogs Error: ", log.Error())
				ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
				return
			}
		}
		// 返回登录状态
		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
		return
	}
	ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})
}

// QQ接入 - 网页授权任意请求类型访问 POST:/oauth/qq/{code}
// QQ授权 - 网页授权任意请求类型访问 GET:/oauth/qq/{scope}
func (c *OauthController) AllQqBy(code string) {
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
	// QQ登录 https://connect.qq.com/index.html
	// 文档 https://wiki.connect.qq.com/%e4%bd%bf%e7%94%a8authorization_code%e8%8e%b7%e5%8f%96access_token

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
	appid := otherCfg["QQ_APP_ID"].(string)      // 申请QQ登录成功后，分配给应用的appid
	appsecret := otherCfg["QQ_APP_KEY"].(string) // 申请QQ登录成功后，分配给网站的appkey
	state := "STATE"                             // 用于保持请求和回调的状态
	display := ""                                // 用于展示的样式。不传则默认展示为PC下的样式。 如果传入“mobile”，则展示为mobile端下的样式。
	platfrom := "腾讯QQ"                           // 缺省平台名

	// 判断是否存在字段 "platfrom"
	if _, ok := allData["platfrom"]; ok {
		platfrom = allData["platfrom"].(string)
	}
	// 判断是否存在字段 "appid"
	if _, ok := allData["appid"]; ok {
		appid = allData["appid"].(string)
	}
	// 判断是否存在字段 "appsecret"
	if _, ok := allData["appsecret"]; ok {
		appsecret = allData["appsecret"].(string)
	}
	// 判断是否存在字段 "state"
	if _, ok := allData["state"]; ok {
		state = allData["state"].(string)
	}
	// 判断是否存在字段 "display"
	if _, ok := allData["display"]; ok {
		display = allData["display"].(string)
	}
	var redirect_uri string
	if _, ok := allData["redirect_uri"]; ok {
		redirect_uri = allData["redirect_uri"].(string)
	}
	// redirect_uri 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
	redirect_uri = url.QueryEscape(redirect_uri) // * 注册qq的回调地址要严格按网站回调域配置的地址来
	// 如果地址与配置不对则报错： redirect uri is illegal [appid: 100315571] (100010) 通知: QQ互联加强网站应用回调地址校验

	println("ctx.Method()", method)

	if method == "GET" {
		// 用户同意授权，获取code

		scope := "get_user_info" //请求用户授权时向用户显示的可进行授权的列表。 不传则默认请求对接口get_user_info进行授权
		if code != "" {
			scope = code
		}

		// https://wiki.connect.qq.com/%e4%bd%bf%e7%94%a8authorization_code%e8%8e%b7%e5%8f%96access_token
		codeUrl := fmt.Sprintf("https://graph.qq.com/oauth2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&display=%s", appid, redirect_uri, scope, state, display)

		ctx.JSON(iris.Map{"data": codeUrl, "code": 0, "msg": ""})
		return
	}

	// -------------------------- post -------------------------

	/*
		https://wiki.connect.qq.com/%e4%bd%bf%e7%94%a8authorization_code%e8%8e%b7%e5%8f%96access_token

		access_token	授权令牌，Access_Token。
		expires_in	该access token的有效期，单位为秒。
		refresh_token	在授权自动续期步骤中，获取新的Access_Token时需要提供的参数。
		注：refresh_token仅一次有效
	*/

	println("通过Authorization Code获取Access Token") // redirect_uri与上面一步中传入的redirect_uri保持一致。
	url := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&fmt=json&redirect_uri=%s", appid, appsecret, code, redirect_uri)

	// println(url)
	// if url != "" {
	// 	return
	// }
	access, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
	if err != nil {
		fmt.Println(url, err)
		ctx.JSON(iris.Map{"code": config.ErrExpireCode, "msg": config.ErrMsgs[config.ErrExpireCode], "data": err.Error()})
		return
	}
	fmt.Printf("access %T -> %v", access, access)

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
	println("---------- access_token ---------")
	println(access_token)
	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range access {
		oauthData[key] = value
	}

	println("---------- 获取用户OpenID_OAuth2.0 ---------")
	// 获取用户OpenID_OAuth2.0
	url = fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s&fmt=json", access_token)
	me, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("me %T -> %v", me, me)
	var openid string
	// 判断是否存在字段 "appid"
	if _, ok := me["openid"]; ok {
		openid = me["openid"].(string)
	}
	if openid == "" {
		ctx.JSON(iris.Map{"data": allData, "code": me["errcode"], "msg": me["errmsg"]})
		return
	}
	println("---------- openid ---------")
	println(openid)

	println("---------- 获取到Access Token和OpenID后，可通过调用OpenAPI来获取或修改用户个人信息。 ---------")
	// 获取到Access Token和OpenID后，可通过调用OpenAPI来获取或修改用户个人信息。
	url = fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s&fmt=json", access_token, appid, openid)
	info, err := utils.HTTPDo("GET", url, map[string][]string{}, map[string]string{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("info %T -> %v", info, info)

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range info {
		oauthData[key] = value
	}

	// 添加其他必要字段
	oauthData["uid"] = tkUid
	oauthData["openid"] = openid
	oauthData["platfrom"] = platfrom

	// 转义对应字段名
	if gender, ok := info["gender"]; ok {
		oauthData["sex"] = gender.(string)
		if oauthData["sex"] == "男" {
			oauthData["sex"] = 1
		} else if oauthData["sex"] == "女" {
			oauthData["sex"] = 2
		} else {
			oauthData["sex"] = 0
		}
	}
	if _, ok := info["figureurl"]; ok {
		oauthData["headimg"] = info["figureurl"].(string)
	}
	if _, ok := info["nickname"]; ok {
		oauthData["nickname"] = info["nickname"].(string)
	}
	if _, ok := info["province"]; ok {
		oauthData["province"] = info["province"].(string)
	}
	if _, ok := info["city"]; ok {
		oauthData["city"] = info["city"].(string)
	}
	if _, ok := info["msg"]; ok {
		oauthData["remark"] = info["msg"].(string)
	}

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	oUser, err := c.Models.CreateOAuth(oauthData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": oauthData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果已绑定用户返回
	if oUser.UID > 0 {
		// 获取用户信息
		user, err := c.Models.Read(oUser.UID)
		if err != nil {
			if env != "" {
				println("Models.Read Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrNoRecords, "msg": config.ErrMsgs[config.ErrNoRecords]})
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
				println("Models.SetLogs Error: ", log.Error())
				ctx.JSON(iris.Map{"data": result, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
				return
			}
		}
		// 返回登录状态
		ctx.JSON(iris.Map{"data": result, "code": 0, "msg": ""})
		return
	}

	ctx.JSON(iris.Map{"data": oUser, "code": 0, "msg": ""})
}
