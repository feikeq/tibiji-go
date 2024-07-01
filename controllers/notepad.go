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
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
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

	// 检查当前用户状态
	statu := c.UserModel.CheckStatus(tkUid)
	if statu == 2 {
		ctx.JSON(iris.Map{"code": config.ErrNoActivate, "msg": config.ErrMsgs[config.ErrNoActivate]})
		return
	} else if statu == 0 {
		ctx.JSON(iris.Map{"code": config.ErrUserDisabled, "msg": config.ErrMsgs[config.ErrUserDisabled]})
		return
	}

	// 获取配置 - 每个人最大云纸张(记事本)数量
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	max := otherCfg["SERV_NOTEPAD_MAX"].(int64)
	// println("系统配置总纸张数：", max)

	// 获取总纸张数量
	num := c.Models.CheckNum(tkUid)
	// println("此用户总纸张数：", num)

	// 如果超出最大数量 -  默认5记事本，vip不限制数目记事本
	if num >= max {
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

	// 客户端IP地址
	allData["ip"] = utils.GetRealIP(ctx)

	// 获取总纸张数量
	total := c.Models.Check(allData["url"].(string))

	if total > 0 {
		ctx.JSON(iris.Map{"code": config.ErrResExists, "msg": config.ErrMsgs[config.ErrResExists]})
		return
	}

	// 调取创建用户模型 - 返回新插入数据的id
	uid, err := c.Models.Create(allData)
	if err != nil {
		if env != "" {
			println("Models.Create Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 返回成功响应
	ctx.JSON(iris.Map{"data": uid, "code": 0, "msg": ""})
}

// 更新纸张  PUT:/Notepad/{nid}
func (c *NotepadController) PutBy(nid int64) {
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

	// 客户端IP地址
	allData["ip"] = utils.GetRealIP(ctx)

	// 删除不能修改的字段
	delete(allData, "uid")     // 删除 用户ID
	delete(allData, "intime")  // 删除 创建时间
	delete(allData, "share")   // 共享地址(区分大小写)
	delete(allData, "referer") // 纸张来源

	// // 权限不够则删除 -  用户自己可以锁定这个记事本不被分享
	// delete(allData, "state") //状态

	notepad, err := c.Models.Detail(nid)
	if err != nil {
		if env != "" {
			println("Models.Find Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": notepad, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果不是自己的纸直接返回
	if tkUid != notepad.UID {
		ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
		return
	}

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Update(nid, allData)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 删除纸张 DELETE:/Notepad/{nid}
func (c *NotepadController) DeleteBy(nid int64) {
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
	row, err := c.Models.Delete(tkUid, nid)
	if err != nil {
		if env != "" {
			println("Models.Delete Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": nid, "_debug_err": err.Error()})
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
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": url, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果这页纸没有主人请去认领
	if notepad.UID == 0 {
		// 闲置空资源待认领
		ctx.JSON(iris.Map{"code": config.ErrEmptyIdle, "msg": config.ErrMsgs[config.ErrEmptyIdle]})
		return
	}

	// 检查纸张的用户状态
	statu := c.UserModel.CheckStatus(notepad.UID)
	if statu == 2 {
		ctx.JSON(iris.Map{"code": config.ErrNoActivate, "msg": config.ErrMsgs[config.ErrNoActivate]})
		return
	} else if statu == 0 {
		ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
		return
	}

	// 数据处理 - 转换时间格式
	notepad.Intime = utils.RFC3339ToString(notepad.Intime, 2) //防止拿到秒级精确时间
	notepad.Uptime = utils.RFC3339ToString(notepad.Uptime, 2) //防止拿到秒级精确时间

	// 判断是否为VIP
	_, v_err := c.UserModel.IsVip(tkUid)
	if v_err != nil {
		notepad.Share = ""
	}

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

	notepad, err := c.Models.Read(uuid)
	if err != nil {
		println("Models.Read Error: ", err.Error())
		ctx.JSON(iris.Map{"code": config.ErrNotFound, "msg": config.ErrMsgs[config.ErrNotFound]})
		return
	}

	// 如果这页纸没有主人请去认领
	if notepad.UID == 0 {
		// 闲置空资源待认领
		ctx.JSON(iris.Map{"code": config.ErrEmptyIdle, "msg": config.ErrMsgs[config.ErrEmptyIdle]})
		return
	}

	// 检查纸张的用户状态
	statu := c.UserModel.CheckStatus(notepad.UID)
	if statu == 2 {
		ctx.JSON(iris.Map{"code": config.ErrNoActivate, "msg": config.ErrMsgs[config.ErrNoActivate]})
		return
	} else if statu == 0 {
		ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
		return
	}

	ctx.JSON(iris.Map{"data": notepad.Content, "code": 0, "msg": ""})
}

// 领取绑定 PATCH:/notepad/{url}
func (c *NotepadController) PatchBy(url string) {
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
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": url, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果这页纸已有主人
	if notepad.UID > 0 {
		if env != "" {
			println("Models.Find Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission], "_debug_carry": notepad, "_debug_err": "这页纸已有主人"})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
		}
		return
	}

	// 如果有密码
	if notepad.Pwd != "" {
		if pwd != notepad.Pwd {
			ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
			return
		}
	}

	// 设置此页纸张主人为本人
	allData["uid"] = tkUid

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Update(notepad.Nid, allData)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}
