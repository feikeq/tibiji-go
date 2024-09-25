package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type RemindController struct {
	DB         *sqlx.DB
	Models     *models.RemindModel // 模型
	UserModels *models.UserModel   // 用户模型

	// 在 Iris 框架中，控制器的 iris.Context 字段会被自动注入为控制器方法的参数。
	// 如果在控制器结构体中定义了名为 CTX 的 iris.Context 字段，Iris 框架会自动将上下文对象赋值给该字段。
	CTX iris.Context // 以便访问 ctx iris.Context （好像也可以不用初始化就能用好奇怪）
}

func NewRemindController(db *sqlx.DB, cfg map[string]interface{}) *RemindController {
	// 返回一个结构体指针
	return &RemindController{
		DB:         db,
		Models:     models.NewRemindModel(db),
		UserModels: models.NewUserModel(db, cfg),
	}
}

// 获取列表 GET:/remind/
func (c *RemindController) Get() {
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
	// fmt.Printf("args: %+v\n", allData) // 打印

	items, err := c.Models.Items(tkUid, allData)
	if err != nil {
		if env != "" {
			println("Models.Items Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": items, "code": 0, "msg": ""})
}

// 任务队列 GET:/remind/task/
func (c *RemindController) GetTask() {
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
	/*
			// 队列文件会保存到 /assets/queue.json 网址直接访问 http://xxx.xxx.xxx/queue.json
			{
			"email": [
				{
					"cid": 1205,
					"uid": 20,
					"fullname": "刘黄芳",
					"pinyin": "liu huang fang ",
					"nickname": "",
					"picture": "",
					"gender": 2,
					"birthday": "1966-08-04 00:00",
					"lunar": 1,
					"grouptag": "",
					"remind": "email::7,1,0||phone::7,1,0",
					"relation": "",
					"note": "",
					"state": 1,
					"remind_num": 7,
					"remind_date": "2024-09-06",
					"remind_cndate": "09月06日",
					"remind_year": 58,
					"remind_star": "处女",
					"remind_zodiac": "马",
					"birthday_cndate": "一九六六年八月初四",
					"remind_day": "7天后",
					"remind_week": "星期五",
					"remind_type": {
						"email": [
							"7",
							"1",
							"0"
						],
						"phone": [
							"7",
							"1",
							"0"
						],
						"notice": null,
						"message": null
					},
					"user_fname": "",
					"user_nickname": "",
					"user_username": "刘建华",
					"user_email": "94xxxxxxxxx89@qq.com",
					"user_cell": "180xxxx6",
					"user_balance": 0,
					"user_vip": 0,
					"user_exptime": "0001-01-01 00:00"
				}
			],
			"phone": [
				{
					"cid": 2356121,
					"uid": 122010,
					"fullname": "妈",
					"pinyin": "ma ",
					"nickname": "",
					"picture": "",
					"gender": 2,
					"birthday": "1960-08-31 00:00",
					"lunar": 0,
					"grouptag": "",
					"remind": "email::7,1,0||phone::7,1,0",
					"relation": "",
					"note": "",
					"state": 1,
					"remind_num": 1,
					"remind_date": "2024-08-31",
					"remind_cndate": "08月31日",
					"remind_year": 64,
					"remind_star": "处女",
					"remind_zodiac": "鼠",
					"birthday_cndate": "一九六〇年七月初十",
					"remind_day": "明天",
					"remind_week": "星期六",
					"remind_type": {
						"email": [
							"7",
							"1",
							"0"
						],
						"phone": [
							"7",
							"1",
							"0"
						],
						"notice": null,
						"message": null
					},
					"user_fname": "",
					"user_nickname": "",
					"user_username": "喵了个咪",
					"user_email": "aiyxxxxx7@163.com",
					"user_cell": "",
					"user_balance": 0,
					"user_vip": 0,
					"user_exptime": "0001-01-01 00:00"
				}
			],
			"notice": null,
			"message": null
		}
	*/

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

	items, err := c.Models.Task()
	if err != nil {
		if env != "" {
			println("Models.Task Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		println("json.Marshal Error: ", err.Error())
	}

	// 将队列数据写入文件中
	file, err := os.Create(AssetsPath + "/queue.json")
	if err != nil {
		println("os.Create Error: ", err.Error())
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(string(data))
	if err != nil {
		println("file.WriteString Error: ", err.Error())
		panic(err)
	}

	ctx.JSON(iris.Map{"data": len(items.Email) + len(items.Phone) + len(items.Notice) + len(items.Message), "code": 0, "msg": "已成功生成队列任务"})
}

// 发送提醒 GET:/remind/queue/
func (c *RemindController) GetQueue() {
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

	// 获取配置项
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	servKeySecret := otherCfg["SERV_KEY_SECRET"].(string)   // API高级密钥
	smsTemplateIds := otherCfg["SMS_TEMPLATE_IDS"].(string) // 短信模版ID 模板类别(0其它 1生日 2纪念日 3闹铃)

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

	// 队列文件地址
	filePath := AssetsPath + "/queue.json"

	// 判断文件是否存在 - 检查主目录是否存在，如果不存在则创建目录
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(iris.Map{"data": 0, "code": 0, "msg": "无任务队列"})
		return
	}

	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		println("os.Open Error: ", err.Error())
		panic(err)
	}
	defer file.Close()

	// // 使用Go语言自带的json库来读取json文件并转为map[string]interface{}类型。
	// decoder := json.NewDecoder(file)
	// var infos map[string]interface{}
	// err = decoder.Decode(&infos)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// 从文件中读取队列数据
	txt, err := io.ReadAll(file)
	if err != nil {
		println("io.ReadAll Error: ", err.Error())
		panic(err)
	}

	var infos, newQueue models.RemindQueues
	err = json.Unmarshal(txt, &infos)
	if err != nil {
		fmt.Println("Unmarshal err,", err)
		ctx.JSON(iris.Map{"code": config.ErrUnknown, "msg": "解析队列失败"})
		return
	}

	// fmt.Printf("data: %+v\n", infos) // 打印
	println("Email 队列总条数：", len(infos.Email))
	println("Phone 队列总条数：", len(infos.Phone))
	println("Notice 队列总条数：", len(infos.Notice))
	println("Message 队列总条数：", len(infos.Message))

	println("从队列中取出短信和邮件并发送")

	//通道（channel）可用于两个 goroutine 之间通过传递一个指定类型的值来同步运行和通讯。
	// 操作符 <- 用于指定通道的方向，发送或接收。如果未指定方向，则为双向通道。
	// emailResult := make(chan error) // 使用chan关键字声明一个通道
	// smsResult := make(chan error)

	go func() {

		if len(infos.Email) > 0 {
			// println("--- Email ---")
			for _, value := range infos.Email {
				// 判断是否有邮箱号
				if value.UserEmail != "" {
					name := value.UserFname
					if name == "" {
						name = value.UserNickname
					}
					if name == "" {
						name = value.UserUsername
					}

					// 获取通知的主题和内容
					subject, body := utils.GetSubjectAndBody(ctx, value.State, value.Lunar, value.RemindYear, name, value.Fullname, value.RemindDay, value.RemindDate, value.Birthday, value.BirthdayCNDate, value.RemindCNDate, value.RemindWeek)
					// println(subject)
					// println(body)
					email := value.UserEmail

					// 邮件发送
					errEmail := utils.SendEmail(ctx, email, subject, body)
					if errEmail != nil {
						println("SendEmail ERR:", errEmail.Error())
						// fmt.Printf("value: %+v\n", value) // 打印
						newQueue.Email = append(newQueue.Email, value) // 添加到失败队列方便下次重试
					}

				}
			}
		}

		if len(infos.Phone) > 0 {
			// println("--- Phone ---")
			for _, value := range infos.Phone {

				phone := value.UserCell
				phone = utils.FormatMobile(phone) // 格式化手机格式

				println(phone)

				// 判断是否有手机号
				if phone != "" {

					// 获取附属资料调取模型 - 根据ID读取数据库中的信息
					material, _ := c.UserModels.ReadMaterial(value.UID)
					// println(value.UID, material.Balance)

					// 数据库中确实有大于一条短信的余额
					if material.Balance > 0.1 {

						age := fmt.Sprintf("%d", value.RemindYear)
						// smsTemplate := [4]string{"0", "1815541", "1815543", "1815721"}
						smsTemplate := strings.Split(smsTemplateIds, ",") // 取配置中的SMS_TEMPLATE_IDS短信模版ID
						_parms := []string{value.RemindDay, value.Fullname, age}
						if value.State == 3 {
							_parms = []string{value.RemindDay, value.Fullname} //变成两个参数
						}
						//短信发送
						errSms := utils.SendSMS(ctx, phone, smsTemplate[value.State], _parms)
						if errSms != nil {
							println("SendSMS ERR:", errSms.Error())
							newQueue.Phone = append(newQueue.Phone, value) // 添加到失败队列方便下次重试
						} else {
							// 进行扣费
							optData := map[string]interface{}{
								"balance": material.Balance - 0.1,
							}
							// 更新用户附属资料表调取模型 - 根据ID更新数据库中的信息
							_, err := c.UserModels.UpdateMaterial(value.UID, optData)
							if err != nil {
								println("Models.UpdateMaterial Error: ", err.Error())
							}

						}

					}

				}
			}
		}

		println("---=-=-=-=--=-=-=-=-=------")
		fmt.Println(newQueue)
		// 格式化数据方便存入缓存文件
		data, err := json.Marshal(newQueue)
		if err != nil {
			println("json.Marshal Error: ", err.Error())
		}

		// 打开队列数据文件
		file, err := os.Create(filePath)
		if err != nil {
			println("os.Create Error: ", err.Error())
			panic(err)
		}
		defer file.Close()
		// 将队列数据写入文件中
		_, err = file.WriteString(string(data))
		if err != nil {
			println("file.WriteString Error: ", err.Error())
			panic(err)
		}

	}()

	ctx.JSON(iris.Map{"data": len(infos.Email) + len(infos.Phone) + len(infos.Notice) + len(infos.Message), "code": 0, "msg": "已推送队列消息"})

}
