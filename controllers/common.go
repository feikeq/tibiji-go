// 通用功能
package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tibiji-go/config"
	"tibiji-go/utils"
	"time"

	"github.com/kataras/iris/v12"
)

// 全局变量必须以关键字var开头，然后是变量名和类型。
var AssetsPath string = "./assets" // 资源文件夹

// CommonController 其它通用控制器 - 不需要模型简单赋值调用
// 定义一个结构体
type CommonController struct {
	// DB *sqlx.DB // 控制器 func (c *CommonController) 里使用c.DB访问数据库连接
}

// 上传文件 POST:/common/upload/
func (c *CommonController) PostUpload(ctx iris.Context) {

	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	upPath := otherCfg["UPLOAD_PATH"].(string)   // 上传目录
	upField := otherCfg["UPLOAD_FIELD"].(string) // 表单名字段名
	println("upPath:", upPath)
	println("upField:", upField)

	// Get the max post value size passed via iris.WithPostMaxMemory.
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	// println("maxSize", maxSize)
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		// ctx.StopWithError(iris.StatusInternalServerError, err)
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	}

	form := ctx.Request().MultipartForm

	fmt.Printf("form=> %T = %v", form, form)
	println()

	// <input type="file" name="upfile" size="30" accept="audio/*,video/*,image/*" capture="camera" multiple>
	files := form.File[upField]
	num := len(files)
	if num == 0 {
		files = form.File[upField+"[]"]
		num = len(files)
	}

	// 如果文件数依旧为0
	if num == 0 {
		// ctx.StopWithError(iris.StatusInternalServerError, err)
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": "请使用 name=upfile 做为file字段名"})
		return
	}

	println("您上传的文件个数：", len(files))
	thePath := AssetsPath + upPath

	// 判断文件是否存在 - 检查主目录是否存在，如果不存在则创建目录
	if _, err := os.Stat(thePath); os.IsNotExist(err) {
		os.MkdirAll(thePath, os.ModePerm) // 自动创建文件夹
	}

	// 生成年月的子目件夹
	folderName := upPath + "/" + time.Now().Format("20060102") + "/"

	theFolderName := AssetsPath + folderName // 真实地址
	// 检查子目录是否存在，如果不存在则创建目录
	if _, err := os.Stat(theFolderName); os.IsNotExist(err) {
		os.MkdirAll(theFolderName, os.ModePerm) // 自动创建文件夹
	}

	failures := ""
	okList := []string{}
	for _, file := range files {
		// 获取后缀名
		ext := filepath.Ext(file.Filename)
		// 生成新的文件名
		newFilename := utils.GenerateTimerID(88888) + ext // （13位时间戳+5随机尾数最大到5个8）
		// 存数据库时不知为何后缀会变小写
		newFilename = strings.ToLower(newFilename) // 将文件名转为小写，以免费有后缀出现大写.JPG就访问不到了

		_realUrl := theFolderName + newFilename
		// println(failures, "=>", _realUrl) // 输出真实地址
		_, err = ctx.SaveFormFile(file, _realUrl)
		if err != nil {
			failures += file.Filename + "，" // 上传文件名
		}
		otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
		uploadDomain := otherCfg["UPLOAD_DOMAIN"].(string)        // 上传文件的访问域名(没有/结尾)
		_virtualURL := folderName + newFilename                   // 输出虚拟地址
		_virtualURL = utils.GetFullURL(uploadDomain, _virtualURL) // 生成完整的文件 URL
		okList = append(okList, _virtualURL)                      // 返回全量地址
	}
	errTXT := ""
	if failures != "" {
		errTXT = "有" + failures + "这些文件上传失败"
	}
	ctx.JSON(iris.Map{"data": okList, "code": 0, "msg": errTXT})
}

// 测试通知 GET:/common/testmsg/
func (c *CommonController) GetTestmsg(ctx iris.Context) {

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

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	servKeySecret := otherCfg["SERV_KEY_SECRET"].(string) // API高级密钥

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("args: %+v\n", allData) // 打印

	var errTxt = ""

	var email string
	// 判断是否存在字段 "email"
	if _, ok := allData["email"]; ok {
		email = allData["email"].(string)
	}

	var phone string
	// 判断是否存在字段 "phone"
	if _, ok := allData["phone"]; ok {
		phone = allData["phone"].(string)
	}

	var types int
	// 判断是否存在字段 "types"
	if val, ok := allData["types"]; ok {
		types = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	if phone == "" && email == "" {
		errTxt = "手机和邮件地址不能全为空"
	}
	// println("types", types)
	// println("email", email)
	// println("phone", phone)

	if errTxt != "" {
		if env != "" {
			println("errTxt Error: ", errTxt)
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty], "_debug_carry": allData, "_debug_err": errTxt})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})

		}
		return
	}

	var key string
	// 判断是否存在字段 "key"
	if _, ok := allData["key"]; ok {
		key = allData["key"].(string)
	}

	// 鉴别API高级密钥
	if servKeySecret != key {
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	var rest string //返回结果

	if email != "" {

		smtpName := otherCfg["SMTP_FROM_NAME"].(string)  // 发件人名称
		smtpWebsite := otherCfg["SMTP_WEBSITE"].(string) // 邮件的网站网址
		smtpUser := otherCfg["SMTP_USER"].(string)       // SMTP服务器用户名
		// println(smtpWebsite, smtpName, smtpUser)

		/*
			[提笔记生日提醒]1天后刘德华32岁生日 / [提笔记周年提醒]今天结婚10周年纪念 / [提笔记一次性提醒]7天后公司年会闹铃

			尊敬的提笔记用户 xxxx 您好：
			03月08日是 张宸熙 的13岁农历生日(2010-02-17)，可别忘了祝福哦~ / 03月01日是 芭月凉 的34岁公历生日(1989-03-01)，可别忘了祝福哦~ / 01月09日是 结婚 的8周年公历纪念(2015-01-09)，可别忘了祝福哦~  / 11月25日是 周六上午8点金台会议 的一次性公历闹铃(2017-11-25)，可别忘了哦~

			(温馨提示：您的提笔记帐户余额0.00元，为了不影响您手机接收短信提醒，请及时充值！)
			提笔记 www.tibiji.com   2023-03-08 08:30:02


			尊敬的提笔记用户您好：
			您的用户名：xxxx
			请务必在24小时内通过下面这个地址修改您的密码，此链接24小时后失效！
			action=RetakePassword&tk=cfcfc7fe5fd159d296c1b787d47f4ea34021
			提笔记 2023-02-22 14:31:55
		*/
		stra := [4]string{" ", "生日", "纪念", "闹铃"}
		strb := [4]string{" ", "岁", "周年", "一次性"}
		days := [4]string{"今天", "1天后", "2天后", "3天后"}
		lunar := [2]string{"公历", "农历"}
		bless := [2]string{"可别忘了祝福哦", "可别忘了哦"}

		// 在go语言里使用fmt.Sprintf拼接字符串字段
		subject := fmt.Sprintf("[%s%s提醒]%s%s%s%s%s", smtpName, stra[1], days[0], "刘德华", "60", strb[1], stra[1])
		body := fmt.Sprintf("尊敬的%s用户 %s 您好：<br/>\r\n", smtpName, "username")
		body += fmt.Sprintf("%s是 %s 的%s%s%s%s(%s)，", "03月08日", "刘德华", "60", strb[1], lunar[0], stra[1], "2020-01-02")
		body += fmt.Sprintf("%s~<br/>\r\n<br/>\r\n", bless[0])
		body += fmt.Sprintf("<a href=\"%s\">提笔记 %s</a> &nbsp; <br/>", smtpWebsite, strings.Split(smtpWebsite, "//")[1])

		// println(subject, "\r\n", body)

		// 邮件发送
		errEmail := utils.SendEmail(ctx, email, subject, body)
		if errEmail != nil {
			ctx.JSON(iris.Map{"code": config.ErrUnknown, "msg": errEmail.Error(), "data": smtpUser})
			return
		}
		rest += "邮件由" + smtpUser + "发送至" + email + "; "
	}

	if phone != "" {
		smsTemplateIds := otherCfg["SMS_TEMPLATE_IDS"].(string) // 短信模版ID 模板类别(0其它 1生日 2纪念日 3闹铃)
		smsTemplate := strings.Split(smsTemplateIds, ",")       // 取配置中的SMS_TEMPLATE_IDS短信模版ID

		// 定义一个二维切片
		args := [][]string{
			{"123456"},
			{"今天", "张惠妹", "22"},
			{"明天", "结婚", "10"},
			{"5天后", "抢高铁票"},
		}

		//短信发送
		errSms := utils.SendSMS(ctx, phone, smsTemplate[types], args[types])
		if errSms != nil {
			ctx.JSON(iris.Map{"code": config.ErrUnknown, "msg": errSms.Error(), "data": smsTemplate})
			return
		}
		rest += "短信由" + smsTemplate[1] + "模板发送至" + phone + "; "
	}

	ctx.JSON(iris.Map{"data": rest, "code": 0, "msg": "Email or Phone Send OK"})
}

// 测试代码 GET:/common/
func (c *CommonController) Get(ctx iris.Context) {
	env := ctx.Values().GetString("ENV")
	tkUid, _ := ctx.Values().GetInt64("UID")
	if env != "" {
		// 打印模块名
		println("\r\n\r\n", env, tkUid)
		println("---------------------------------------------------------")
		println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
		println("---------------------------------------------------------")
	}

	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	uploadDomain := otherCfg["UPLOAD_DOMAIN"].(string) // 上传文件的访问域名(没有/结尾)
	noParameter := otherCfg["NO_PARAMETER"].(bool)     // 上传文件目录地址不存参数内容

	url := "http://192.168.172.88:8888/uploads/20121212/test_pic.JPG?a=1&b=2#ccc"
	ok := utils.GetDirectoryPath(uploadDomain, url, noParameter) // 从完整的文件 URL 中提取目录路径

	ctx.JSON(iris.Map{"data": ok, "code": 0, "msg": url})

}
