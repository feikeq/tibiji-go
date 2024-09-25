package controllers

import (
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type AccountController struct {
	DB     *sqlx.DB
	Models *models.AccountModel
	CTX    iris.Context
}

func NewAccountController(db *sqlx.DB) *AccountController {
	// 返回一个结构体指针
	return &AccountController{
		DB:     db,
		Models: models.NewAccountModel(db),
	}
}

// 帐单列表 GET:/account/
// pageNumber 分页页码
// pageSize 每页条数
// pageOrder 排序方向
// pageField 排序字段
func (c *AccountController) Get() {
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

	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// 获取每页条数配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	pageSize := otherCfg["SERV_LIST_SIZE"].(int64)

	// 分页参数
	var pageNumber int64 = 1
	if pageSize < 1 {
		pageSize = 20 // 默认单页条数
	}

	var pageOrder, pageField string = "desc", "aid" //  pageOrder 也支持这种写法 ascend descend
	var searchOR string                             // 模糊查询字段

	// 判断是否存在字段 "page"
	if _, ok := allData["pageNumber"]; ok {
		pageNumber = utils.ParseInt64(allData["pageNumber"]) // 任意数据转int64数字
	}
	// 判断是否存在字段 "pageSize"
	if _, ok := allData["pageSize"]; ok {
		pageSize = utils.ParseInt64(allData["pageSize"]) // 任意数据转int64数字
	}

	// 判断是否存在字段 "pageOrder"
	if _, ok := allData["pageOrder"]; ok {
		pageOrder = allData["pageOrder"].(string)
	}
	// 判断是否存在字段 "pageField"
	if _, ok := allData["pageField"]; ok {
		pageField = allData["pageField"].(string)
	}

	// 判断是否存在字段 "searchOR"
	if val, ok := allData["searchOR"]; ok {
		searchOR = val.(string)
	}

	list, total, err := c.Models.List(tkUid, allData, pageNumber, pageSize, pageOrder, pageField, searchOR)
	if err != nil {
		if env != "" {
			println("Models.List Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	type temp struct {
		Total      int64                `json:"total" description:"总条数"`
		PageNumber int64                `json:"pageNumber" description:"分页页码"`
		PageSize   int64                `json:"pageSize" description:"每页条数"`
		PageOrder  string               `json:"pageOrder" description:"排序方向"`
		PageField  string               `json:"pageField" description:"排序字段"`
		List       []models.AccountInfo `json:"list" description:"列表数据"`
	}
	data := temp{total, pageNumber, pageSize, pageOrder, pageField, list}
	ctx.JSON(iris.Map{"data": data, "code": 0, "msg": ""})

}

// 添加帐单 POST:/account/
func (c *AccountController) Post() {
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

	allData := utils.AllDataToMap(ctx) // ctx.FormValues() 等同于 ctx.Request().Form

	// 判断是否存在字段 "money"
	if _, ok := allData["money"]; !ok {
		println("帐单金额不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["money"] == "" {
			println("帐单金额不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
	}

	// 判断是否存在字段 "class"
	if _, ok := allData["class"]; !ok {
		println("帐单主分类不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["class"] == "" {
			println("帐单主分类不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
	}

	// 判断是否存在字段 "sort"
	if _, ok := allData["sort"]; !ok {
		println("帐单子类别不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["sort"] == "" {
			println("帐单子类别不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
	}

	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)
	// 添加用户ID
	allData["uid"] = tkUid

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
	// ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"data": uid, "code": 0, "msg": ""})
}

// 帐单详情  PUT:/account/{aid}/
func (c *AccountController) GetBy(id int64) {
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

	// 联系人调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Read(id)
	if err != nil {
		if env != "" {
			println("Models.Update Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 如果不是本人的联系人
	if row.UID != tkUid {
		// // 如果不是管理员
		// if !c.UserModels.IsAdmin(tkUid) {

		// }
		// 详情功能必须是自己创建的才可以，管理员也无没权获取他人联系人
		ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 修改帐单  PUT:/account/{aid}/
func (c *AccountController) PutBy(id int64) {
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

	// 删除不能修改的字段
	delete(allData, "uid")    // 删除 用户ID
	delete(allData, "intime") // 删除 注册时间
	delete(allData, "uptime") // 删除 更新时间

	// 调取模型 - 根据ID更新数据库中的信息
	row, err := c.Models.Update(tkUid, id, allData)
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

// 帐单类目 GET:/account/type
func (c *AccountController) GetType() {
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

	types := c.Models.Types()
	// fmt.Printf("types %T -> %v", types, types)
	ctx.JSON(iris.Map{"data": types, "code": 0, "msg": ""})
}

// 月份帐单 GET:/account/month/
func (c *AccountController) GetMonth() {
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

	allData := utils.AllDataToMap(ctx) // ctx.FormValues() 等同于 ctx.Request().Form

	var year, month string
	// 判断是否存在字段 "year"
	if _, ok := allData["year"]; !ok {
		println("年份不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["year"] == "" {
			println("年份不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		} else {
			year = allData["year"].(string)
		}
	}

	// 判断是否存在字段 "month"
	if _, ok := allData["month"]; !ok {
		println("月份不能为空")
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["month"] == "" {
			println("月份不能为空")
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		} else {
			month = allData["month"].(string)
		}
	}

	// 调取模型 - 根据参数获取数据库中的信息
	monthList, err := c.Models.Month(tkUid, year, month)
	if err != nil {
		ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		return
	}

	// fmt.Printf("types %T -> %v", types, types)
	ctx.JSON(iris.Map{"data": monthList, "code": 0, "msg": ""})
}

// 帐单日历 GET:/account/calendar/
func (c *AccountController) GetCalendar() {
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

	// 调取模型 - 根据ID获取数据库中的信息
	calendarList, err := c.Models.Calendar(tkUid)
	if err != nil {
		ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		return
	}

	// fmt.Printf("types %T -> %v", types, types)
	ctx.JSON(iris.Map{"data": calendarList, "code": 0, "msg": ""})
}

// 删除帐单 DELETE:/account/{aid}
func (c *AccountController) DeleteBy(id int64) {
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
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 收支比例 GET:/account/report/ratio/
func (c *AccountController) GetReportRatio() {
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

	var year, month int

	// 判断是否存在字段 "year"
	if val, ok := allData["year"]; ok {
		year = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}
	// 判断是否存在字段 "month"
	if val, ok := allData["month"]; ok {
		month = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	row, err := c.Models.ReportRatio(tkUid, year, month)
	if err != nil {
		if env != "" {
			println("Models.ReportRatio Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 近年统计 - 收支列表 ( 近一年或半年的统计 ) GET:/account/report/ratios/
func (c *AccountController) GetReportRatios() {
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

	var limit int

	// 判断是否存在字段 "limit"
	if _, ok := allData["limit"]; ok {
		limit = utils.ParseInt(allData["limit"]) // 任意数据数字（字符串转数字，字符转数字）
	}

	row, err := c.Models.ReportRatios(tkUid, limit)
	if err != nil {
		if env != "" {
			println("Models.ReportRatio Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 收支明细 GET:/account/report/details/
func (c *AccountController) GetReportDetails() {
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

	var year, month int
	var item string

	// 判断是否存在字段 "item"
	if _, ok := allData["item"]; ok {
		item = allData["item"].(string)
	} else {
		item = "支出"
	}

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "year"
	if val, ok := allData["year"]; ok {
		year = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}
	// 判断是否存在字段 "month"
	if val, ok := allData["month"]; ok {
		month = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	row, err := c.Models.ReportDetails(tkUid, item, year, month)
	if err != nil {
		if env != "" {
			println("Models.ReportDetails Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 流水账户 GET:/account/report/accounts/
func (c *AccountController) GetReportAccounts() {
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

	var year, month int

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "year"
	if val, ok := allData["year"]; ok {
		year = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}
	// 判断是否存在字段 "month"
	if val, ok := allData["month"]; ok {
		month = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	row, err := c.Models.ReportAccounts(tkUid, year, month)
	if err != nil {
		if env != "" {
			println("Models.ReportAccounts Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 收支对象 GET:/account/objects/
func (c *AccountController) GetObjects() {
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

	row, err := c.Models.Objects(tkUid)
	if err != nil {
		if env != "" {
			println("Models.Objects Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}
