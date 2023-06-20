package controllers

import (
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type NotepadController struct {
	DB        *sqlx.DB
	Models    *models.NotepadModel
	UserModel *models.UserModel
	CTX       iris.Context
}

func NewNotepadController(db *sqlx.DB) *NotepadController {
	// 返回一个结构体指针
	return &NotepadController{
		DB:        db,
		Models:    models.NewNotepadModel(db),
		UserModel: models.NewUserModel(db),
	}
}

// 纸张列表 GET:/notepad
func (c *NotepadController) Get() {
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

	list, err := c.Models.List(tkUid)
	if err != nil {
		if env != "" {
			println("Models.List Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": list, "code": 0, "msg": ""})
}

// 添加纸张 POST:/notepad
func (c *NotepadController) Post() {
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

	// 获取配置 - 每个人最大云纸张(记事本)数量
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	max := otherCfg["SERV_NOTEPAD_MAX"].(int64)
	// println("系统配置总纸张数：", max)

	// 获取总纸张数量
	total := c.Models.Check(tkUid)
	// println("此用户总纸张数：", total)

	// 如果超出最大数量 -  默认5记事本，vip不限制数目记事本
	if total > max {
		// 判断是否为VIP
		_, v_err := c.UserModel.IsVip(tkUid)
		if v_err != nil {
			ctx.JSON(iris.Map{"code": config.ErrNotVip, "msg": config.ErrMsgs[config.ErrNotVip]})
			return
		}
	}

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)

	// 添加用户ID
	allData["uid"] = tkUid
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// 调取创建用户模型 - 返回新插入数据的id
	uid, err := c.Models.Create(allData)
	if err != nil {
		if env != "" {
			println("Models.Create Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 返回成功响应
	ctx.JSON(iris.Map{"data": uid, "code": 0, "msg": ""})
}

// 更新纸张  PUT:/Notepad/{cid}
func (c *NotepadController) PutBy(id int64) {
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

	// 删除不能修改的字段
	delete(allData, "uid")    // 删除 用户ID
	delete(allData, "intime") // 删除 创建时间
	delete(allData, "share")  // 共享地址(区分大小写)

	// 权限不够则删除
	delete(allData, "state") // 管理员才能修改 状态

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Update(tkUid, id, allData)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 删除纸张 DELETE:/Notepad/{cid}
func (c *NotepadController) DeleteBy(id int64) {
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

	// // 拿所有提交数据
	// allData := utils.AllDataToMap(ctx)

	// 调取模型 - 根据ID删除数据库中的信息
	row, err := c.Models.Delete(tkUid, id)
	if err != nil {
		if env != "" {
			println("Models.Delete Error: ", err.Error())
			ctx.JSON(iris.Map{"data": id, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 获取纸张 GET:/notepad/{url}
func (c *NotepadController) GetBy(url string) {
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

	var pwd string

	// 判断是否存在字段 "pwd"
	if _, ok := allData["pwd"]; ok {
		pwd = allData["pwd"].(string)
	}

	notepad, err := c.Models.Find(url)
	if err != nil {
		if env != "" {
			println("Models.Find Error: ", err.Error())
			ctx.JSON(iris.Map{"data": url, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 数据处理 - 转换时间格式
	notepad.Intime = utils.RFC3339ToString(notepad.Intime, 2) //防止拿到秒级精确时间
	notepad.Uptime = utils.RFC3339ToString(notepad.Uptime, 2) //防止拿到秒级精确时间

	// 如果没有密码直接返回
	if notepad.Pwd == "" {
		ctx.JSON(iris.Map{"data": notepad, "code": 0, "msg": ""})
		return
	}

	// 如果密码相符直接返回
	if pwd == notepad.Pwd {
		ctx.JSON(iris.Map{"data": notepad, "code": 0, "msg": ""})
		return
	}

	// 如果是自己的纸直接返回
	if tkUid == notepad.UID {
		ctx.JSON(iris.Map{"data": notepad, "code": 0, "msg": ""})
		return
	}

	ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
}

// 分享纸张 GET:/notepad/share/{uuid}
func (c *NotepadController) GetShareBy(uuid string) {
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

	txt, err := c.Models.Read(uuid)
	if err != nil {
		println("Models.Read Error: ", err.Error())
		ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		return
	}

	ctx.JSON(iris.Map{"data": txt, "code": 0, "msg": ""})
}
