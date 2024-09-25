package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tibiji-go/config"
	"tibiji-go/models"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kataras/iris/v12"
)

type ContactController struct {
	DB         *sqlx.DB
	Models     *models.ContactModel
	UserModels *models.UserModel // 用户模型
	CTX        iris.Context
}

func NewContactController(db *sqlx.DB, cfg map[string]interface{}) *ContactController {
	// 返回一个结构体指针
	return &ContactController{
		DB:         db,
		Models:     models.NewContactModel(db, cfg),
		UserModels: models.NewUserModel(db, cfg),
	}
}

// 联系人列表 GET:/contact/
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

	var pageOrder, pageField string = "desc", "cid" //  pageOrder 也支持这种写法 ascend descend
	var searchOR string                             // 模糊查询字段

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "page"
	if val, ok := allData["pageNumber"]; ok {
		pageNumber = utils.ParseInt64(val) // 任意数据转int64数字（字符串转数字）
	}
	// 判断是否存在字段 "pageSize"
	if val, ok := allData["pageSize"]; ok {
		pageSize = utils.ParseInt64(val) // 任意数据转int64数字（字符串转数字，字符转数字）
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

	// 获取自己的联系人列表
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
		List       []models.ContactInfo `json:"list" description:"列表数据"`
	}
	data := temp{total, pageNumber, pageSize, pageOrder, pageField, list}
	ctx.JSON(iris.Map{"data": data, "code": 0, "msg": ""})
}

// 添加联系人 POST:/contact/
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

	// // 判断是否存在字段 "gender"
	// if gender, ok := allData["gender"]; ok {
	// 	// 强制将字符转为数字 防止出现提交“男”或“女”进入数据库char(1)因为没有使用tinyint(1)
	// 	allData["gender"], _ = strconv.Atoi(gender.(string)) // 字符转数字
	// 	// fmt.Printf("gender %T -> %v", allData["gender"], allData["gender"])
	// }

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "year"
	if val, ok := allData["gender"]; ok {
		allData["gender"] = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	// 添加用户ID
	allData["uid"] = tkUid
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// 调取创建用户模型 - 返回新插入数据的id
	cid, err := c.Models.Create(allData)
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
	ctx.JSON(iris.Map{"data": cid, "code": 0, "msg": ""})
}

// 修改联系人  PUT:/contact/{cid}/
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

	// 在 Go 语言中，if 语句的短变量声明 (:=) 中定义的变量具有局部作用域，仅在 if 语句块中有效。
	// 因此，if 语句中的 val 和 ok 变量不会影响其他 if 语句中的同名变量。

	// 判断是否存在字段 "gender"
	if val, ok := allData["gender"]; ok {
		allData["gender"] = utils.ParseInt(val) // 任意数据数字（字符串转数字，字符转数字）
	}

	// 删除不能修改的字段
	delete(allData, "uid")    // 删除 用户ID
	delete(allData, "cid")    // 删除 CID
	delete(allData, "intime") // 删除 创建时间

	// 权限不够则删除
	delete(allData, "state") //  删除状态（管理员也不能修改）

	// 只能修改自己添加的联系人，因为更新条件会自带带当前用户ID

	// 修改联系人修改联系人调取模型 - 根据ID更新数据库中的信息
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

	// 获取附属资料调取模型 - 根据ID读取数据库中的信息
	user, err := c.UserModels.ReadMaterial(tkUid)
	if err != nil {
		if env != "" {
			println("Models.ReadMaterial Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}
	if user.CID == id {
		// 修改联系人时如果是用户绑定的名片则允许更新用户资料
		// 映射数据
		upData := map[string]interface{}{}

		if val, ok := allData["fullname"]; ok {
			if val != "" && *user.FName == "" {
				upData["fname"] = val.(string)
			}
		}
		if val, ok := allData["nickname"]; ok {
			if val != "" && *user.NickName == "" {
				upData["nickname"] = val.(string)
			}
		}
		if val, ok := allData["picture"]; ok {
			if val != "" && *user.Headimg == "" {
				upData["headimg"] = val.(string)
			}
		}
		if val, ok := allData["gender"]; ok {
			if val != 0 && *user.Sex == 0 {
				upData["sex"] = val.(int)
			}
		}
		if val, ok := allData["company"]; ok {
			if val != "" && *user.Company == "" {
				upData["company"] = val.(string)
			}
		}
		if val, ok := allData["address"]; ok {
			if val != "" && *user.Address == "" {
				upData["address"] = val.(string)
			}
		}
		if val, ok := allData["phone"]; ok {
			if val != "" && *user.Cell == "" {
				// upData["cell"] = val.(string)
				arr_type := strings.Split(val.(string), "||") // 分割字符串 - 拆分(表达式||值)
				arr_item := strings.Split(arr_type[0], "::")  // 分割字符串 - 拆分(表达式::值)
				upData["cell"] = arr_item[1]
			}
		}
		if val, ok := allData["mail"]; ok {
			if val != "" && *user.Email == "" {
				println("---- email ----", val.(string))
				// upData["email"] = val.(string)
				arr_type := strings.Split(val.(string), "||") // 分割字符串 - 拆分(表达式||值)
				arr_item := strings.Split(arr_type[0], "::")  // 分割字符串 - 拆分(表达式::值)
				upData["email"] = arr_item[1]

			}
		}
		// 调取模型 - 根据ID更新数据库中的信息
		row2, err2 := c.UserModels.Update(tkUid, upData)
		if err2 != nil {
			if env != "" {
				println("Models.Update Error: ", err2.Error())
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err2.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
			}
			return
		}
		println("更新用户资料影响行数据：", row2)

	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 删除联系人 DELETE:/contact/{cid}/
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
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": id, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	ctx.JSON(iris.Map{"data": row, "code": 0, "msg": ""})
}

// 联系人分组 GET:/contact/groups/
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

	// 获取自己的联系人分组
	row, err := c.Models.Groups(tkUid)
	if err != nil {
		if env != "" {
			println("Models.Groups Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
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
		newFilename := utils.GenerateTimerID(88888) + ext // （13位时间戳+5随机尾数最大到5个8

		// println(failures, "=>", theFolderName+newFilename)
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

// 获取联系人  GET:/contact/{cid}/
func (c *ContactController) GetBy(cid int64) {
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
	row, err := c.Models.Read(cid)
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

// 关联联系人  PATCH:/contact/   - 绑定用户资料和附属资料到联系人
func (c *ContactController) Patch() {
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

	// 获取附属资料调取模型 - 根据ID读取数据库中的信息
	user, err := c.UserModels.ReadMaterial(tkUid)
	if err != nil {
		if env != "" {
			println("Models.ReadMaterial Error: ", err.Error())
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": tkUid, "_debug_err": err.Error()})
		} else {
			ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
		}
		return
	}

	// 使用 fmt.Println 打印结构体数据
	// fmt.Println(user)

	var bindCid = user.CID
	println("bindCid：", bindCid)

	// 拿所有提交数据
	allData := utils.AllDataToMap(ctx)
	if *user.NickName != "" {
		allData["nickname"] = *user.NickName
	}
	if *user.Headimg != "" {
		allData["picture"] = *user.Headimg
	}
	if *user.Sex != 0 {
		allData["gender"] = *user.Sex
	}
	if *user.Birthday != "" {
		allData["birthday"] = utils.RFC3339ToString(*user.Birthday, 2) //防止拿到秒级精确时间
	}
	if *user.Company != "" {
		allData["company"] = *user.Company
	}
	if *user.Address != "" {
		allData["address"] = *user.Address
	}
	if *user.FName != "" {
		allData["fullname"] = *user.FName
	}
	if *user.UserName != "" {
		allData["note"] = *user.UserName
		// 这里不能直接判断 allData["nickname"] 有内存地址在
		if *user.NickName == "" {
			allData["nickname"] = *user.UserName
		}
	}
	if *user.Email != "" {
		// 在go语言里可直接使用+号拼接字符串字段
		allData["mail"] = "EMAIL::" + *user.Email
	}
	if *user.Cell != "" {
		allData["phone"] = "TEL::" + *user.Cell
	}

	// 添加用户ID
	allData["uid"] = tkUid
	// fmt.Printf("变量类型type: %T, 变量的值value: %v\n", allData, allData)

	// fmt.Printf("allData: %+v\n", allData) // 打印allData
	// 如果没有绑定联系人卡片
	if bindCid == 0 {
		println("----- 根据用户资料创建联系人 -----")
		cid, err := c.Models.Create(allData)
		if err != nil {
			if env != "" {
				println("Models.Create Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": user, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
			}
			return
		}
		bindCid = cid

		println("新创建联系人ID:", bindCid)

		upData := map[string]interface{}{
			"cid":    bindCid,
			"remark": "通过绑定联系人添加cid",
		}

		// 更新用户附属资料表调取模型 - 根据ID更新数据库中的信息
		row, err := c.UserModels.UpdateMaterial(tkUid, upData)
		if err != nil {
			if env != "" {
				println("Models.UpdateMaterial Error: ", err.Error())
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": err.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
			}
			return
		}

		if row == 0 {
			ctx.JSON(iris.Map{"code": config.ErrResExists, "msg": config.ErrMsgs[config.ErrResExists]})
			return
		}

	} else {

		// 联系人调取模型 - 根据ID更新数据库中的信息
		row, err := c.Models.Read(bindCid)
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
			// 必须是自己创建的才可以，管理员也无没权获取他人联系人
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}

		// 删除已存在的字段

		if row.NickName != "" {
			delete(allData, "nickname")
		}
		if row.Picture != "" {
			delete(allData, "picture")
		}
		if row.Gender != 0 {
			delete(allData, "gender")
		}
		if row.Birthday != "" {
			delete(allData, "birthday")
		}
		if row.Company != "" {
			delete(allData, "company")
		}
		if row.Address != "" {
			delete(allData, "address")
		}
		if row.Fullname != "" {
			delete(allData, "fullname")
		}
		if row.Note != "" {
			delete(allData, "note")
		}
		if row.Mail != "" {
			delete(allData, "mail")
		}
		if row.Phone != "" {
			delete(allData, "phone")
		}

		delete(allData, "uid")    // 删除 用户ID
		delete(allData, "cid")    // 删除 CID
		delete(allData, "intime") // 删除 创建时间
		// 权限不够则删除
		delete(allData, "state") //  删除状态（管理员也不能修改）

		// 只能修改自己添加的联系人，因为更新条件会自带带当前用户ID

		// 调取模型 - 根据ID更新数据库中的信息
		_, uperr := c.Models.Update(tkUid, bindCid, allData)
		if uperr != nil {
			if env != "" {
				println("Models.Update Error: ", uperr.Error())
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase], "_debug_carry": allData, "_debug_err": uperr.Error()})
			} else {
				ctx.JSON(iris.Map{"code": config.ErrDatabase, "msg": config.ErrMsgs[config.ErrDatabase]})
			}
			return
		}

	}

	// 操作写入日志表
	ua := ctx.GetHeader("User-Agent") // 拿到UA信息User-Agent
	logData := map[string]interface{}{
		"uid":    tkUid,
		"action": "correlated",
		"note":   "contact",
		"actip":  utils.GetRealIP(ctx),
		"ua":     ua,
	}
	log := c.UserModels.SetLogs(logData)
	if log != nil {
		if env != "" {
			println("Models.SetLogs Error: ", log.Error())
			ctx.JSON(iris.Map{"data": bindCid, "code": 0, "msg": "操作成功但日志记录失败", "_debug_carry": logData, "_debug_err": log.Error()})
			return
		}
	}

	// 返回成功响应
	ctx.JSON(iris.Map{"data": bindCid, "code": 0, "msg": ""})

}
