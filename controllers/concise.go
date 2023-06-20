package controllers

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

// 定义一个结构体
type ConciseController struct {
	DB *sqlx.DB // 控制器 func (c *ConciseController) 里使用c.DB访问数据库连接
}

// Concise 其它任务的控制器 - 不需要模型简单赋值调用

func (c *ConciseController) Get(ctx iris.Context) interface{} {
	// 打印模块名
	println("---------------------------------------------------------")
	println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
	println("---------------------------------------------------------")

	// 在 Get 方法中调用 Post 方法
	// 如果在同一个结构体中的方法之间可以相互调用。所以在 Get 方法中您可以通过直接调用 c.Post(ctx) 来调用 Post 方法：
	response := c.Post(ctx)

	fmt.Printf("Type: %T , Data: %v\n", response, response)

	var uid int64

	if data, ok := response.(map[string]interface{}); ok {
		code := data["code"].(int)
		dataValue := data["data"].(int)

		// 打印获取的值
		println("Code:", code)
		println("Data:", dataValue)

		uid = int64(dataValue)
	}

	// // 调取别的模块下的模型或控制器

	// // 创建 XXXModel 实例
	// XXXModel := models.NewXXXModel(c.DB)
	// // 使用 XXXModel 进行操作
	// XXX, _ := XXXModel.Read(uid)
	// // 直接赋值 XXXController 中的结构体
	// XXXController := &XXXController{DB:c.DB,Models:models.NewXXXModel(c.DB)}
	// // 或者调用 XXXController 中的工厂函数
	// XXXController := NewXXXController(c.DB)

	return iris.Map{"data": uid, "code": iris.StatusOK, "msg": ""} // 返回提交结果
}

func (c *ConciseController) Post(ctx iris.Context) interface{} {
	// 打印模块名
	println("---------------------------------------------------------")
	println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
	println("---------------------------------------------------------")

	id := 1684200935383956

	// 返回成功响应

	// ctx.StatusCode(iris.StatusOK)
	// ctx.JSON(iris.Map{"data":id, "code": iris.StatusOK,"msg":""})
	return iris.Map{"data": id, "code": iris.StatusOK, "msg": ""} // 返回结果
}

func (c *ConciseController) GetBy(ctx iris.Context, uid int64) interface{} {
	// 打印模块名
	println("---------------------------------------------------------")
	println(ctx.GetCurrentRoute().MainHandlerName() + " [" + ctx.GetCurrentRoute().Path() + "] " + ctx.Method())
	println("---------------------------------------------------------")

	// ctx.JSON(user)
	return iris.Map{"data": 0, "code": 0, "msg": ""} // 返回结果

}
