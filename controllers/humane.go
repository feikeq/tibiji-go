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

// 人情列表 GET:/humane
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

	items, err := c.Models.Items(tkUid, allData, pageOrder, pageField)
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

// 人情详情 GET:/humane/{cid}
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

// 统计比例 GET:/humane/report/ratio
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
		if env != "" {
			println("Models.ReportRatio Error: ", err.Error())
			ctx.JSON(iris.Map{"data": tkUid, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 礼金分类回报排名 GET:/humane/report/top
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
			ctx.JSON(iris.Map{"data": tkUid, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}
