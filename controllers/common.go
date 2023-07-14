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

// 上传文件 POST:/common/upload
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
		newFilename := utils.GenerateTimerID(88888) + ext // 五位随机数据最大到5个8

		// println(failures, ":::=>", theFolderName+newFilename)
		_, err = ctx.SaveFormFile(file, theFolderName+newFilename)
		if err != nil {
			failures += file.Filename + "，"
		}
		okList = append(okList, folderName+newFilename) // 输出虚拟地址
	}
	errTXT := ""
	if failures != "" {
		errTXT = "有" + failures + "这些文件上传失败"
	}
	ctx.JSON(iris.Map{"data": okList, "code": 0, "msg": errTXT})
}

// 邮件发送 GET:/common/email
func (c *CommonController) GetEmail(ctx iris.Context) {

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
	servKeySecret := otherCfg["SERV_KEY_SECRET"].(string) // API高级密钥

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	// fmt.Printf("args: %+v\n", allData) // 打印

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

	smtpName := otherCfg["SMTP_FROM_NAME"].(string)  // 发件人名称
	smtpWebsite := otherCfg["SMTP_WEBSITE"].(string) // 邮件的网站网址

	// println(smtpHost, smtpPort, smtpUser, smtpPass, smtpEmail, smtpName)

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

	subject := fmt.Sprintf("[%s%s提醒]%s%s%s%s%s", smtpName, stra[1], days[0], "刘德华", "60", strb[1], stra[1])
	body := fmt.Sprintf("尊敬的%s用户 %s 您好：<br/>\r\n", smtpName, "username")
	body += fmt.Sprintf("%s是 %s 的%s%s%s%s(%s)，", "03月08日", "刘德华", "60", strb[1], lunar[0], stra[1], "2020-01-02")
	body += fmt.Sprintf("%s~<br/>\r\n<br/>\r\n", bless[0])
	body += fmt.Sprintf("<a href=\"%s\">提笔记 %s</a> &nbsp; <br/>", smtpWebsite, strings.Split(smtpWebsite, "//")[1])

	println(subject, "\r\n", body)

	// 邮件发送
	errEmail := utils.SendEmail(ctx, "feikeq@qq.com", subject, body)
	if errEmail != nil {
		println("SendEmail ERR:", errEmail.Error())
		ctx.JSON(iris.Map{"code": config.ErrUnknown, "msg": config.ErrMsgs[config.ErrUnknown]})
		return
	}

	//短信发送
	errSms := utils.SendSMS(ctx, "13838389438", "1815541", []string{"今天", "张惠妹", "22"})
	if errSms != nil {
		println("SendSMS ERR:", errSms.Error())
		ctx.JSON(iris.Map{"code": config.ErrUnknown, "msg": config.ErrMsgs[config.ErrUnknown]})
		return
	}

	ctx.JSON(iris.Map{"data": 0, "code": 0, "msg": "Email Send ok"})
}
