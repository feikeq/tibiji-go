package controllers

import (
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type NotepadController struct {
	DB        *sqlx.DB
	Models    *models.NotepadModel
	UserModel *models.UserModel
	CTX       iris.Context
}

func NewNotepadController(db *sqlx.DB, cfg map[string]interface{}) *NotepadController {
	// 返回一个结构体指针
	return &NotepadController{
		DB:        db,
		Models:    models.NewNotepadModel(db),
		UserModel: models.NewUserModel(db, cfg),
	}
}

// 纸张列表 GET:/notepad/
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

// 添加纸张 POST:/notepad/
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

// 更新纸张  PUT:/notepad/{nid}/
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

	var pwdKey string // 操作密码

	// 判断是否存在字段 "key"
	if val, ok := allData["key"]; ok {
		if val != nil {
			pwdKey = allData["key"].(string)
		} else {
			// 处理传过来 null 的情况
			pwdKey = ""
		}
	}

	// // 权限不够则删除 -  用户自己可以锁定这个记事本不被分享
	// delete(allData, "state") //状态

	notepad, err := c.Models.Detail(nid)
	if err != nil {
		if env != "" {
			println("Models.Detail Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": notepad, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果不是自己的纸
	if tkUid != notepad.UID {
		// 如果有密码
		if notepad.Pwd != "" {
			// 如果密码不相符直接返回
			if pwdKey != notepad.Pwd {
				ctx.JSON(iris.Map{"code": config.ErrNoPermission, "msg": config.ErrMsgs[config.ErrNoPermission]})
				return
			}

		}
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

// 删除纸张 DELETE:/notepad/{nid}/
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

// 获取纸张 GET:/notepad/{url}/
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

	// 2分钟内可访问4次每次间隔30秒
	{
		// 限制访问频率 - 接口访问频率限制： 利用 Go 标准库中的 sync.Mutex 和 time 包来实现基于 IP 地址的请求限流

		// 定义限流时间(秒)时间间隔是否大于limiting秒(相当于30秒内可以发生多少次 limitmax 判断)
		limiting := 30 * time.Second
		// 定义限流时间(秒)内最大请求次数(为0时不判断 就只走clearting规则，非0时相当于limitmax+1次)
		limitmax := 3
		// 定义超过时间(分钟)未请求视为过期(相当于5分钟可以发生多少次 limiting 判断)
		clearting := 2 * time.Minute

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
	vip, v_err := c.UserModel.IsVip(notepad.UID)
	if v_err != nil {
		notepad.Share = ""
	}
	println(notepad.UID, "  vip =", vip)

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

// 分享纸张 GET:/notepad/share/{uuid}/
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

	// 判断是否为VIP
	_, v_err := c.UserModel.IsVip(notepad.UID)
	if v_err != nil {
		ctx.JSON(iris.Map{"code": config.ErrNotVip, "msg": config.ErrMsgs[config.ErrNotVip]})
		return
	}

	ctx.JSON(iris.Map{"data": notepad.Content, "code": 0, "msg": ""})
}

// 领取绑定 PATCH:/notepad/{url}/
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
		// 如果本身就是这页纸的主人则返回正常
		if tkUid == notepad.UID {
			ctx.JSON(iris.Map{"data": notepad.Nid, "code": 0, "msg": ""})
		} else {
			if env != "" {
				println("Models.Find Error: ", "如果这页纸已有主人")
				ctx.JSON(iris.Map{"code": config.ErrResExists, "msg": config.ErrMsgs[config.ErrResExists], "_debug_carry": notepad, "_debug_err": "这页纸已有主人"})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrResExists, "msg": config.ErrMsgs[config.ErrResExists]})
			}
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
