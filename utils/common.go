// 通用函数定义

package utils

import (
	"crypto/tls"
	"net/smtp"

	"github.com/kataras/iris/v12"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// 邮件发送
// err := utils.SendEmail(ctx,"ccav@123.com","标题","内容")
func SendEmail(ctx iris.Context, to, subject, body string) error {
	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	smtpHost := otherCfg["SMTP_HOST"].(string)        // SMTP服务器
	smtpPort := otherCfg["SMTP_PORT"].(string)        // SMTP服务器端口
	smtpUser := otherCfg["SMTP_USER"].(string)        // SMTP服务器用户名
	smtpPass := otherCfg["SMTP_PASS"].(string)        // SMTP服务器密码
	smtpEmail := otherCfg["SMTP_FROM_EMAIL"].(string) // 发件人EMAIL
	smtpName := otherCfg["SMTP_FROM_NAME"].(string)   // 发件人名称

	// println(smtpHost, smtpPort, smtpUser, smtpPass, smtpEmail, smtpName)

	contentType := "Content-Type: text/html; charset=UTF-8"
	rfc := "To: " + to + "\r\n" +
		"From: " + smtpName + "<" + smtpEmail + ">" + "\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n" +
		"\r\n" +
		body + "\r\n"
	// 要求 msg 参数要符合 RFC 822 电子邮件的标准格式。
	msg := []byte(rfc)
	// println(rfc)

	// 设置认证信息并连接到 SMTP 服务器
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// smtp.SendMail 这种方式根本发不出邮件老是报错EOF
	// err := smtp.SendMail(smtpHost+smtpPort, auth, smtpEmail, to, msg)
	// if err != nil {
	// 	println("SendMail Err:", err.Error())
	// 	ctx.JSON(iris.Map{"data": 0, "code": config.ErrUnknown, "msg": err.Error()})
	// 	return
	// }

	// SSL/TLS Email Example

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// 这里是钥匙，您需要调用 tls.dial 而不是 smtp.dial
	// 用于在需要SSL连接的465上运行的SMTP服务器
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpHost+smtpPort, tlsconfig)
	if err != nil {
		println("Dial ERR:", err.Error())
		return err
	}

	smtp_new, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		println("NewClient ERR:", err.Error())
		return err
	}

	// Auth
	if err = smtp_new.Auth(auth); err != nil {
		println("Auth ERR:", err.Error())
		return err
	}

	// To && From
	if err = smtp_new.Mail(smtpEmail); err != nil {
		println("Mail ERR:", err.Error())
		return err
	}

	if err = smtp_new.Rcpt(to); err != nil {
		println("Rcpt ERR:", err.Error())
		return err
	}

	// Data
	w, err := smtp_new.Data()
	if err != nil {
		println("Data ERR:", err.Error())
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		println("Write ERR:", err.Error())
		return err
	}

	err = w.Close()
	if err != nil {
		println("Close ERR:", err.Error())
		return err
	}

	smtp_new.Quit()

	return nil
}

// 短信发送
// err = utils.SendSMS(ctx,"13838389438","模板ID",[]string{"模板变量1","模板变量2"})
func SendSMS(ctx iris.Context, phone, templateId string, templateArgs []string) error {
	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	smsSecretId := otherCfg["SMS_SECRET_ID"].(string)   // API密钥ID
	smsSecretKey := otherCfg["SMS_SECRET_KEY"].(string) // API密钥KEY
	smsSignName := otherCfg["SMS_SIGN_NAME"].(string)   // 短信签名
	smsSdkAppId := otherCfg["SMS_SDK_APPID"].(string)   // 短信SdkAppId

	/*
		腾讯云发送短信 https://console.cloud.tencent.com/api/explorer?Product=sms&Version=2021-01-11&Action=SendSms
		SDK使用说明
		1.安装公共基础包：
		go get -v -u github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common

		2.安装对应的产品包（如 sms
		go get -v -u github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111

	*/

	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := common.NewCredential(smsSecretId, smsSecretKey)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := sms.NewClient(credential, "ap-beijing", cpf) // 地域参数还支持  ap-guangzhou、ap-nanjing

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := sms.NewSendSmsRequest()

	// 下发手机号码，采用 E.164 标准，格式为+[国家或地区码][手机号]，单次请求最多支持200个手机号且要求全为境内手机号或全为境外手机号。
	// 例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号。
	// 注：发送国内短信格式还支持0086、86或无任何国家或地区码的11位手机号码，前缀默认为+86。
	request.PhoneNumberSet = common.StringPtrs([]string{phone})

	// 短信 SdkAppId，在 短信控制台 添加应用后生成的实际 SdkAppId，示例如1400006666。
	request.SmsSdkAppId = common.StringPtr(smsSdkAppId)

	// 短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名，例如：腾讯云，签名信息可前往 国内短信 或 国际/港澳台短信 的签名管理查看。 发送国内短信该参数必填。
	request.SignName = common.StringPtr(smsSignName)

	// 模板 ID，必须填写已审核通过的模板 ID。模板 ID 可前往 国内短信 或 国际/港澳台短信 的正文模板管理查看，若向境外手机号发送短信，仅支持使用国际/港澳台短信模板。
	request.TemplateId = common.StringPtr(templateId)

	// 模板参数，若无模板参数，则设置为空。（模板参数的个数需要与 TemplateId 对应模板的变量个数保持一致。）
	request.TemplateParamSet = common.StringPtrs(templateArgs)

	// 返回的resp是一个SendSmsResponse的实例，与请求对象对应
	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		// UnauthorizedOperation.SmsSdkAppIdVerifyFail  SmsSdkAppId 校验失败，请检查 SmsSdkAppId 是否属于 云API密钥 的关联账户。
		// FailedOperation.FailResolvePacket	请求包解析失败，通常情况下是由于没有遵守 API 接口说明规范导致的，请参考 请求包体解析1004错误详解。
		// FailedOperation.SignatureIncorrectOrUnapproved	签名未审批或格式错误。（1）可登录 短信控制台，核查签名是否已审批并且审批通过；（2）核查是否符合格式规范，签名只能由中英文、数字组成，要求2 - 12个字，若存在疑问可联系 腾讯云短信小助手。
		return err
	}
	if err != nil {
		println("SendSms ERR:", err.Error())
		return err
	}
	// 输出json格式的字符串回包
	println(response.ToJsonString())
	return nil
}

// 通知推送
// err := utils.SendMessage(ctx)
func SendMessage(ctx iris.Context) error {
	return nil
}
