package controllers

import (
	"encoding/json"
	"fmt"
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

// 接入用户 POST:/oauth
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

	// 接入一个新的的第三方平台用户 - 返回新插入数据的id
	id, err := c.Models.CreateOAuth(allData)
	if err != nil {
		if env != "" {
			println("Models.CreateOAuth Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 返回成功响应
	//向数据模板传值 当然也可以绑定其他值
	// ctx.ViewData("", mapData)
	// ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"data": id, "code": 0, "msg": ""})
}

// 微信接入 - 网页授权任意请求类型访问 POST:/oauth/wx/{code}
func (c *OauthController) PostWxBy(code string) {
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

	oauthData := make(map[string]interface{})
	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range allData {
		oauthData[key] = value
	}

	// 获取配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	wx_appid := otherCfg["WX_APPID"].(string)         // 公众号的唯一标识
	wx_appsecret := otherCfg["WX_APPSECRET"].(string) // 公众号的appsecret

	// 判断是否存在字段 "appid"
	if _, ok := allData["appid"]; ok {
		wx_appid = allData["appid"].(string)
	}
	// 判断是否存在字段 "appid"
	if _, ok := allData["appsecret"]; ok {
		wx_appsecret = allData["appsecret"].(string)
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
