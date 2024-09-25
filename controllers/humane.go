package controllers

import (
	"fmt"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type HumaneController struct {
	DB     *sqlx.DB
	Models *models.HumaneModel // 模型

	// 在 Iris 框架中，控制器的 iris.Context 字段会被自动注入为控制器方法的参数。
	// 如果在控制器结构体中定义了名为 CTX 的 iris.Context 字段，Iris 框架会自动将上下文对象赋值给该字段。
	CTX iris.Context // 以便访问 ctx iris.Context （好像也可以不用初始化就能用好奇怪）
}

func NewHumaneController(db *sqlx.DB) *HumaneController {
	// 返回一个结构体指针
	return &HumaneController{
		DB:     db,
		Models: models.NewHumaneModel(db),
	}
}

// 人情列表 GET:/humane/
func (c *HumaneController) Get() {
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
	fmt.Printf("Type: %T , Data: %v\n", allData, allData) // 打印

	var pageOrder, pageField string = "desc", "btime" //  pageOrder 也支持这种写法 ascend descend

	// 判断是否存在字段 "pageOrder"
	if _, ok := allData["pageOrder"]; ok {
		pageOrder = allData["pageOrder"].(string)
	}
	// 判断是否存在字段 "pageField"
	if _, ok := allData["pageField"]; ok {
		pageField = allData["pageField"].(string)
	}
	listSize := 4

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "listSize"
	if val, ok := allData["listSize"]; ok {
		listSize = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	items, err := c.Models.Items(tkUid, pageOrder, pageField, listSize)
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

// 人情详情 GET:/humane/{cid}/
func (c *HumaneController) GetBy(cid int64) {
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
	// fmt.Printf("Type: %T , Data: %v\n", allData, allData) // 打印

	items, err := c.Models.Read(tkUid, cid)
	if err != nil {
		ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		return
	}
	ctx.JSON(iris.Map{"data": items, "code": 0, "msg": ""})
}

// 统计比例 GET:/humane/report/ratio/
func (c *HumaneController) GetReportRatio() {
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

	row, err := c.Models.ReportRatio(tkUid)
	if err != nil {
		// if env != "" {
		// 	println("Models.ReportRatio Error: ", err.Error())
		// 	ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		// } else {
		// 	ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		// }

		// 当用户一条记都没有的时候会查出来null结果集导致报钷
		// 因为这个结果结只会有一条记录不会有别的错所以不做错误处理了直接返回空
		ctx.JSON(iris.Map{"data": nil, "code": 0, "msg": ""})

		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 礼金分类回报排名 GET:/humane/report/top/
func (c *HumaneController) GetReportTop() {
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

	row, err := c.Models.ReportTop(tkUid)
	if err != nil {
		if env != "" {
			println("Models.ReportTop Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}
