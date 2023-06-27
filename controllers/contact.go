package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type ContactController struct {
	DB     *sqlx.DB
	Models *models.ContactModel
	CTX    iris.Context
}

func NewContactController(db *sqlx.DB) *ContactController {
	// 返回一个结构体指针
	return &ContactController{
		DB:     db,
		Models: models.NewContactModel(db),
	}
}

// 联系人列表 GET:/contact
func (c *ContactController) Get() {
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

	// 获取每页条数配置
	otherCfg := ctx.Application().ConfigurationReadOnly().GetOther()
	pageSize := otherCfg["SERV_LIST_SIZE"].(int64)

	// 分页参数
	var pageNumber int64 = 1
	if pageSize < 1 {
		pageSize = 20 // 默认单页条数
	}

	var pageOrder, pageField string = "desc", "uid" //  pageOrder 也支持这种写法 ascend descend

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

	list, total, err := c.Models.List(tkUid, allData, pageNumber, pageSize, pageOrder, pageField)
	if err != nil {
		if env != "" {
			println("Models.List Error: ", err.Error())
			ctx.JSON(iris.Map{"data": allData, "code": "err debug", "msg": err.Error()})
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
		List       []models.ContactInfo `json:"list" description:"列表数据"`
	}
	data := temp{total, pageNumber, pageSize, pageOrder, pageField, list}
	ctx.JSON(iris.Map{"data": data, "code": total, "msg": ""})
}

// 添加联系人 POST:/contact
func (c *ContactController) Post() {
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

	// 判断是否存在字段 "fullname"
	if _, ok := allData["fullname"]; !ok {
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	} else {
		if allData["fullname"] == "" {
			ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
			return
		}
	}

	// 添加用户ID
	allData["uid"] = tkUid
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// 调取创建用户模型 - 返回新插入数据的id
	cid, err := c.Models.Create(allData)
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
	ctx.JSON(iris.Map{"data": cid, "code": 0, "msg": ""})
}

// 修改联系人  PUT:/contact/{cid}
func (c *ContactController) PutBy(id int64) {
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
	delete(allData, "cid")    // 删除 CID
	delete(allData, "intime") // 删除 创建时间

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

// 删除联系人 DELETE:/contact/{cid}
func (c *ContactController) DeleteBy(id int64) {
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

// 联系人分组 GET:/contact/groups
func (c *ContactController) GetGroups() {
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

	row, err := c.Models.Groups(tkUid)
	if err != nil {
		if env != "" {
			println("Models.Groups Error: ", err.Error())
			ctx.JSON(iris.Map{"data": tkUid, "code": "err debug", "msg": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 导入VCF POST:/contact/vcards
/*
导出 vCard
在 iCloud.com 上的“通讯录”中，点按以选择联系人列表中的所需联系人。如果你要导出多个联系人，请按住 Command 键（Mac 电脑上）或 Ctrl 键（Windows 电脑上），然后点按你要导出的每个联系人。
点按边栏中的 “显示操作菜单”弹出式按钮，然后选取“导出 vCard”。
如果你选择多个联系人，“通讯录”会导出一个包含所有联系人的 vCard。
*/
func (c *ContactController) PostVcards() {
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
	upPath := otherCfg["UPLOAD_PATH"].(string)   // 上传目录
	upField := otherCfg["UPLOAD_FIELD"].(string) // 表单名字段名
	println("upPath:", upPath)
	println("upField:", upField)

	// Get the max post value size passed via iris.WithPostMaxMemory.
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	// println("maxSize", maxSize)
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		// ctx.StopWithError(iris.StatusInternalServerError, err)
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": config.ErrMsgs[config.ErrParamEmpty]})
		return
	}

	form := ctx.Request().MultipartForm

	fmt.Printf("form=> %T = %v", form, form)
	println()

	// <input type="file" name="upfile" size="30" accept="audio/*,video/*,image/*" capture="camera" multiple>
	files := form.File[upField]
	num := len(files)
	if num == 0 {
		files = form.File[upField+"[]"]
		num = len(files)
	}

	// 如果文件数依旧为0
	if num == 0 {
		// ctx.StopWithError(iris.StatusInternalServerError, err)
		ctx.JSON(iris.Map{"code": config.ErrParamEmpty, "msg": "请使用 name=upfile 做为file字段名"})
		return
	}

	println("您上传的文件个数：", len(files))
	thePath := AssetsPath + upPath

	// 判断文件是否存在 - 检查主目录是否存在，如果不存在则创建目录
	if _, err := os.Stat(thePath); os.IsNotExist(err) {
		os.MkdirAll(thePath, os.ModePerm) // 自动创建文件夹
	}

	// 生成年月的子目件夹
	folderName := upPath + "/" + time.Now().Format("20060102") + "/"

	theFolderName := AssetsPath + folderName // 真实地址
	// 检查子目录是否存在，如果不存在则创建目录
	if _, err := os.Stat(theFolderName); os.IsNotExist(err) {
		os.MkdirAll(theFolderName, os.ModePerm) // 自动创建文件夹
	}

	failures := ""
	okList := []string{}
	status := map[string]int{
		"total":   0,
		"success": 0,
	}
	for _, file := range files {
		// 获取后缀名
		ext := filepath.Ext(file.Filename)
		// 生成新的文件名
		newFilename := utils.GenerateTimerID(88888) + ext // 五位随机数据最大到5个8

		// println(failures, ":::=>", theFolderName+newFilename)
		_, err = ctx.SaveFormFile(file, theFolderName+newFilename)
		if err != nil {
			failures += file.Filename + "，"
		}
		okList = append(okList, folderName+newFilename) // 输出虚拟地址

		// 解析联系人数据
		data := utils.ParseVCards(theFolderName+newFilename, AssetsPath)
		// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", data, data)

		// 入库操作
		for _, val := range data {
			status["total"] += 1
			// 添加用户ID
			val["uid"] = tkUid
			// 调取创建用户模型 - 返回新插入数据的id
			cid, _ := c.Models.Create(val)
			// println()
			// println(" ------------- ", cid)
			// println("入库操作")
			if cid > 0 {
				status["success"] += 1
			}
			// fmt.Println(val)

		}
	}
	// errTXT := ""
	// if failures != "" {
	// 	errTXT = "有" + failures + "这些文件上传失败"
	// }
	ctx.JSON(iris.Map{"data": status, "code": 0, "msg": okList})

}
