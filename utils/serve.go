// 服务器相关函数定义

package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/teris-io/shortid"
)

// 生成时间唯一ID
// ( 默认randomSource随机数据是999，可以指定随机数 )
// utils.GenerateTimerID() 默认16位
// utils.GenerateTimerID(88888)) 这种就是13+5位置随机数是18位
func GenerateTimerID(randomSource ...int) string {
	num := 999
	if len(randomSource) > 0 {
		num = randomSource[0]
	}
	randomLength := len(strconv.Itoa(num)) // 默认随机数长度为 3

	// 获取时间，该时间带有时区等信息，获取的为当前地区所用时区的时间
	now := time.Now()                                // 2022-12-01 15:52:47.310419 +0800 CST m=+5.643172070
	timestamp := now.UnixMilli()                     // 获取当前毫秒时间戳
	timestampStr := strconv.FormatInt(timestamp, 10) // 转为10进制数字的字符串 数字转字符串

	// 随机数据这里不采用真随机 "crypto/rand" 性能太差，我们这里采用伪随机 “crypto/rand” 性能高10倍！
	// 使用不同的seed来确保每次启动都产生新的随机数
	rand.Seed(now.UnixNano()) // 获取当前纳秒时间戳

	// 前置补零(前导零)，保证随机数始终是 randomLength 长度的位数
	randomNum := fmt.Sprintf("%0"+strconv.Itoa(randomLength)+"d", rand.Intn(num))

	// 拼接时间戳和随机数，生成16位用户ID
	temID := timestampStr + randomNum

	return temID
}

// 生成Token密钥
// token, _ := utils.GenerateToken(1, 86400, "secret")
func GenerateToken(uid int64, second int64, secret string) (string, error) {
	// TOKEN = AES( uid + MD5( uid + MD5(secret) + ( unxi + second) ) + ( unxi + second) , MD5(secret))
	// uid用户ID、secret:key配置密钥+UA、unxi当前时间、second配置过期秒、MD5加密、AES加密

	// println(" ------- 生成Token密钥 ------")
	// println("uid", uid)
	// println("second", second)
	// println("secret", secret)

	// 获取过期时间
	expiration := time.Now().Add(time.Second * time.Duration(second)).Unix()
	// 计算密钥
	key := CalculateMD5(secret)                  // 做MD5处理保证32位字符
	uidStr := strconv.FormatInt(uid, 10)         // 用户ID转为10进制数字的字符串 数字转字符串
	expStr := strconv.FormatInt(expiration, 10)  // 过期时间转为10进制数字的字符串 数字转字符串
	valid := CalculateMD5(uidStr + key + expStr) // 验证串
	token := uidStr + valid + expStr             // token明文
	// nowStr := time.Now().Unix()                  // 获取当时间戳秒
	// println("nowStr", nowStr)                    // nowStr 1684300013
	// println("expStr", expStr)                    // expStr 1684472813
	// println("token:", token)                     // token: 16842021669879880cd0949e4af83201e8f2bd8f98104de31684472813

	// 生成签名字符串
	tokenString, err := EncryptAES(token, key)
	if err != nil {
		return "", fmt.Errorf("无法生成 token: %v", err)
	}
	// println("tokenAES:", tokenString) // 2Q2eHwGaZXKLaN2zNcc3tdNZdh9qpSBHBADukhhjRF9felQ9I4FB4LMiUZeJKce13KSbEutAEBeXzGBbQnf1eKLNi1wvdWlgmDI=
	return tokenString, nil
}

// 验证Token并返回UID
// uid, err := utils.VerifyToken("2QyXH1aZNnuKZtG6Ns0/u4ILIUs2qCEbAVfpmE9gR14CcwQ3ddFOs78qUFje2Z3DT9gQCgxCkc55ywQ=", "secret")
func VerifyToken(tokenString string, secret string) (int64, error) {
	// TOKEN = AES( uid + MD5( uid + MD5(key) + ( unxi + second) ) + ( unxi + second) , MD5(key))
	// uid用户ID、key配置密钥、unxi当前时间、second配置过期秒、MD5加密、AES加密

	// 2QyXH1aZNnuKZtG6Ns0/u4ILIUs2qCEbAVfpmE9gR14CcwQ3ddFOs78qUFje2Z3DT9gQCgxCkc55ywQ=

	key := CalculateMD5(secret) // 做MD5处理保证32位字符
	// 解析 Token
	token, err := DecryptAES(tokenString, key)
	if err != nil {
		return 0, fmt.Errorf("无法解析 token: %v", err)
	}

	// token字符长度
	tokeLen := len(token)
	if tokeLen < 42 {
		// 如果长度小于42不是合法的toke验证失败
		return 0, fmt.Errorf("非法 token")
	}

	// 拆解TOKEN
	expStr := token[tokeLen-10:]                 // 最后10位是秒的时间戳
	valid := token[tokeLen-42 : tokeLen-10]      // 验证串是中间32位置
	uidStr := token[:tokeLen-42]                 // 用户ID转为10进制数字的字符串
	check := CalculateMD5(uidStr + key + expStr) // 生成检查串

	// println("expStr", expStr)
	// println("valid", valid)
	// println("check", check)
	// println("uidStr", uidStr)
	uid, _ := strconv.ParseInt(uidStr, 10, 64) // 字符转数字
	exp, _ := strconv.ParseInt(expStr, 10, 64) // 字符转数字
	// println("exp", exp)

	// 校验 Token 是否有效
	if valid != check {
		return 0, fmt.Errorf("token 无效")
	}

	// 检查过期时间是否超时（当前时间是否大于token有效时间）
	if time.Now().Unix() > exp {
		return 0, fmt.Errorf("token 已过期")
	}

	// println("uid", uid)
	return uid, nil
}

// 密码脱敏 密码 = MD5(注册时间格式字符串+MD5(明文密码串))
// utils.HashPassword("pwdStr","2000-10-10 10:20:30")
func HashPassword(password string, intime string) string {
	hash := CalculateMD5(password)
	hashedPassword := CalculateMD5(intime + hash)
	return hashedPassword
}

// 手机号脱敏函数
// utils.MaskPhoneNumber("13838389438")
func MaskPhoneNumber(phoneNumber string) string {
	if phoneNumber == "" {
		return ""
	}
	// 获取手机号的长度
	phoneLength := len(phoneNumber)
	if phoneLength < 7 {
		// 如果手机号长度小于7，直接返回原始手机号
		return phoneNumber
	}

	// 保留前三位和后四位，其他部分用 * 替代
	masked := phoneNumber[:3] + "****" + phoneNumber[phoneLength-4:]
	return masked
}

// 银行卡脱敏函数
// utils.MaskBankCardNumber("12121221121212")
func MaskBankCardNumber(cardNumber string) string {
	if cardNumber == "" {
		return ""
	}
	// 保留后四位，其他部分用 * 替代
	masked := []rune(cardNumber)
	for i := 0; i < len(masked)-4; i++ {
		masked[i] = '*'
	}
	return string(masked)
}

// 真实姓名脱敏
// utils.MaskRealName("刘德华")
func MaskRealName(name string) string {
	if name == "" {
		return ""
	}
	// 获取姓氏的首个字符
	firstChar := string([]rune(name)[0])

	// 生成星号或其他符号的字符串
	mask := strings.Repeat("*", len([]rune(name))-1)

	// 拼接姓氏首字符和脱敏字符串
	return firstChar + mask

}

// 身份证号脱敏函数
// utils.MaskIDCardNumber("123456789012345678")
func MaskIDCardNumber(idCardNumber string) string {
	if idCardNumber == "" {
		return ""
	}
	// 保留前六位和后四位，其他部分用 * 替代
	masked := []rune(idCardNumber)
	for i := 6; i < len(masked)-4; i++ {
		masked[i] = '*'
	}
	return string(masked)
}

// 邮箱脱敏函数
// utils.MaskEmail("xxxx@ccav.tv")
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	// 找到 @ 符号的位置
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		// 如果找不到 @ 符号，返回原始邮箱
		return email
	}

	// 获取用户名部分和域名部分
	username := email[:atIndex]
	domain := email[atIndex+1:]

	// 获取用户名的长度
	usernameLength := len(username)
	if usernameLength <= 2 {
		// 如果用户名长度小于等于2，直接返回第一个字符和最后一个字符，并拼接上域名
		return username[:1] + "****" + username[usernameLength-1:] + "@" + domain
	}

	// 取用户名的前两个字符和最后两个字符，其余部分用 * 替代
	maskedUsername := username[:1] + strings.Repeat("*", usernameLength-2) + username[usernameLength-1:]

	// 拼接脱敏后的用户名和域名
	return maskedUsername + "@" + domain
}

// 检查邮箱格式
// utils.CheckEmail("13838389438@139.com")
func CheckEmail(email string) bool {
	// pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 检查手机格式
// utils.CheckMobile("13838389438")
func CheckMobile(cell string) bool {
	// pattern := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	// pattern := "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
	pattern := "^1[3-9]{1}\\d{9}$" // 第一位必为1的第二位不能为2和1的十一位数字
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(cell)
}

// 格式化手机格式
// utils.FormatMobile("1-380-260-8878")
func FormatMobile(cell string) string {
	// phoneNumbers := []string{"1-380-260-8878", "138-3838-9438", "+86a1383838943", "182 1080 1989"}
	reg := regexp.MustCompile(`[^+\d]+`)
	return reg.ReplaceAllString(cell, "")
}

// 检查身份证号
// utils.CheckIdCard("440308199901101512")
func CheckIdCard(idcard string) bool {
	// 18位身份证匹配规则 ^(\d{17})([0-9]|X)$
	// 15位身份证 (^\d{15}$) // 18位身份证 (^\d{18}$)
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	pattern := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"
	// 正则调用规则
	reg := regexp.MustCompile(pattern)
	// 返回 MatchString 是否匹配
	return reg.MatchString(idcard)
}

// 获取用户真实的 IP 地址
// utils.GetRealIP(ctx)
func GetRealIP(ctx iris.Context) string {
	// 在 Iris 框架中，可以通过 ctx.RemoteAddr() 方法来获取用户的 IP 地址。
	// 然而，当使用 NGINX 代理层时，RemoteAddr() 方法返回的是代理服务器的 IP 地址，而不是用户的真实 IP 地址。
	// 要获取 NGINX 代理一层后用户的真实 IP 地址，可以使用 X-Real-IP 或 X-Forwarded-For 头字段。
	// 这些头字段通常由代理服务器设置，用于传递真实的客户端 IP 地址。

	// 优先获取 X-Real-IP 头字段
	// 如果使用代理，则尝试获取 X-Real-IP 或 X-Forwarded-For 头字段
	if forwardedIP := ctx.GetHeader("X-Real-IP"); forwardedIP != "" {
		return forwardedIP

		// 如果 X-Real-IP 头字段不存在，则尝试获取 X-Forwarded-For 头字段
	} else if forwardedIPs := ctx.GetHeader("X-Forwarded-For"); forwardedIPs != "" {
		// 多个 IP 地址可能以逗号分隔，取第一个 IP 地址
		return strings.Split(forwardedIPs, ",")[0] // 分割字符串
	}

	// 如果以上头字段都不存在，则返回 RemoteAddr
	return ctx.RemoteAddr()
}

// 获取过滤条件数据（表达式::值）
// utils.GetWhereArgs(map)
func GetWhereArgs(filters map[string]interface{}) ([]string, map[string]interface{}) {
	/*
		条件过滤字段筛选数据样本
		"sex":"!=::0","headimg":"!=::","username":"LIKE::fk%","state":1,"intime":"<::2023-05-23 07:00:00","country":"NOT LIKE::%美国%"
		生成的sql条件是
		AND `headimg` != '' AND `sex` != 0 AND `country` NOT LIKE '%美国%' AND `state` = 1 AND `intime` < '2023-05-23 07:00:00' AND `username` LIKE 'fk%'
		注意：因为MAP是无序的所以生成的筛选条件顺序也完全是随机的
	*/

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	where := make([]string, 0)
	// values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range filters {

		var exp, val, str string

		// 使用一个switch语句来检查每个值的类型进行类型断言 判断类型
		switch strVal := value.(type) {
		case int:
			str = fmt.Sprintf("%d", strVal)
			// println(key, "int")
		case float64:
			str = fmt.Sprintf("%d", int(strVal)) // 类型强制转换后变成字符
			// println(key, "float64")
		default:
			// 处理其他类型
			str = strVal.(string)
			// println(key, "default")
		}

		// 拆分 （表达式:值）
		temp := strings.Split(str, "::") // 分割字符串

		if len(temp) > 1 {
			exp = temp[0]
			val = temp[1]
		} else {
			exp = "="
			val = temp[0]

		}

		where = append(where, "AND `"+key+"` "+exp+" :"+key)
		// values = append(values, ":"+key)
		args[key] = val

		// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", val, val)
	}

	return where, args
}

// 获取通知的主题和内容
// subject,body := utils.GetSubjectAndBody(info.Birthday, info.Lunar)
func GetSubjectAndBody(ctx iris.Context, state, lunar, age int, name, fullname, remindDay, remindDate, birthday, birthdayCNdate, remindCNdate, remindWeek string) (string, string) {
	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	smtpName := otherCfg["SMTP_FROM_NAME"].(string)  // 发件人名称
	smtpWebsite := otherCfg["SMTP_WEBSITE"].(string) // 邮件的网站网址
	// println(smtpWebsite, smtpName)

	/*
		[提笔记生日提醒]1天后刘德华32岁生日 / [提笔记周年提醒]今天结婚10周年纪念 / [提笔记一次性提醒]7天后公司年会一次性闹铃

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
	strl := [2]string{"公历", "农历"}
	bless := "可别忘了祝福哦"
	if state != 1 {
		bless = "可别忘了哦"
	}

	age_str := fmt.Sprintf("%d%s", age, strb[state])
	if state == 3 {
		// 当为一次性闹铃时去掉年份字符
		age_str = strb[state]
	}
	date := strings.Split(birthday, " ")[0] // 分割字符串
	if lunar == 1 {
		date = birthdayCNdate
	}

	subject := fmt.Sprintf("[%s%s提醒]%s%s%s%s", smtpName, stra[state], remindDay, fullname, age_str, stra[state])
	// println(subject)
	body := fmt.Sprintf("尊敬的%s用户 %s 您好：<br/>\r\n", smtpName, name)
	body += fmt.Sprintf("%s是 %s 的%s%s%s(%s)，", remindCNdate, fullname, age_str, strl[lunar], stra[state], date)
	body += fmt.Sprintf("%s~<br/>\r\n<br/>\r\n", bless)
	body += fmt.Sprintf("<a href=\"%s\">%s %s</a> &nbsp; <br/>", smtpWebsite, smtpName, strings.Split(smtpWebsite, "//")[1])

	// println(body)

	return subject, body
}

// 独特的非连续短ID生成器
// utils.GenerateShortId()
func GenerateShortId() string {
	// 创建一个新的ShortID生成器
	generator := shortid.MustNew(1, shortid.DefaultABC, 2342)
	// github.com/teris-io/shortid 生成的ID默认是区分大小写的，这意味着 abc 和 ABC 被视为两个不同的ID。
	// 如果您需要生成不区分大小写的ID，可以将 DefaultABC 常量替换为
	// shortid.DefaultABC + "abcdefghijklmnopqrstuvwxyz"，这样生成的ID将仅包含小写字母。

	// 生成一个新的ShortID
	id, err := generator.Generate()
	if err != nil {
		panic(err)
	}
	return id
}
