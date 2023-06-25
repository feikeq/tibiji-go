package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type RemindController struct {
	DB     *sqlx.DB
	Models *models.RemindModel // 模型

	// 在 Iris 框架中，控制器的 iris.Context 字段会被自动注入为控制器方法的参数。
	// 如果在控制器结构体中定义了名为 CTX 的 iris.Context 字段，Iris 框架会自动将上下文对象赋值给该字段。
	CTX iris.Context // 以便访问 ctx iris.Context （好像也可以不用初始化就能用好奇怪）
}

func NewRemindController(db *sqlx.DB) *RemindController {
	// 返回一个结构体指针
	return &RemindController{
		DB:     db,
		Models: models.NewRemindModel(db),
	}
}

// 提醒列表 GET:/remind
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
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": items, "code": 0, "msg": ""})
}

// 提醒任务 GET:/remind/task
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
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
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

// 提醒队列 GET:/remind/queue
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
			for _, value := range infos.Email {
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

		if len(infos.Phone) > 0 {
			for _, value := range infos.Phone {

				phone := value.UserCell

				// 如果有余额
				if value.UserBalance > 0 {
					phone = utils.FormatMobile(phone) // 格式化手机格式
				} else {
					phone = ""
				}

				// 判断是否有手机号
				if phone != "" {
					age := fmt.Sprintf("%d", value.RemindYear)
					smsTemplate := [4]string{"0", "1815541", "1815543", "1815721"}
					//短信发送
					errSms := utils.SendSMS(ctx, phone, smsTemplate[value.State], []string{value.RemindDay, value.Fullname, age})
					if errSms != nil {
						println("SendSMS ERR:", errSms.Error())
						newQueue.Phone = append(newQueue.Phone, value) // 添加到失败队列方便下次重试
					}
				}
			}
		}

		println("---=-=-=-=--=-=-=-=-=------")
		fmt.Println(newQueue)
		// 存入缓存文件
		data, err := json.Marshal(newQueue)
		if err != nil {
			println("json.Marshal Error: ", err.Error())
		}

		// 将队列数据写入文件中
		file, err := os.Create(filePath)
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

	}()

	ctx.JSON(iris.Map{"data": len(infos.Email) + len(infos.Phone) + len(infos.Notice) + len(infos.Message), "code": 0, "msg": "已推送队列消息"})

}
