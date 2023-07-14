package models

// CRUD 常用操作 Create（创建）、Read（读取）、Update（更新）、Delete（删除）、Insert（插入）、Change（修改）
import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

// 用户模型的结构体
type UserModel struct {
	DB                *sqlx.DB
	TableName         string
	OauthTableName    string
	MaterialTableName string
	LogsTableName     string
}

// 创建一个新的模型实例
func NewUserModel(db *sqlx.DB) *UserModel {
	return &UserModel{
		DB:                db,
		TableName:         "sys_user",
		OauthTableName:    "sys_oauth",
		MaterialTableName: "sys_material",
		LogsTableName:     "sys_logs",
	}
}

// 其它引用数据表名常量 - 项目初始化的时候会用所所以不删除
const (
	UserTableName = "sys_user" // 虽然结构体定义了但模型上其实用不到
)

// 用户结构体
type UserInfo struct {
	UID          *int64  `db:"uid" json:"uid" description:"用户ID"`
	UserName     *string `db:"username" json:"username" description:"帐号"`
	Ciphers      *string `db:"ciphers" json:"ciphers" description:"密码"`
	Email        *string `db:"email" json:"email" description:"邮箱"`
	Cell         *string `db:"cell" json:"cell" description:"电话"`
	NickName     *string `db:"nickname" json:"nickname" description:"昵称"`
	Headimg      *string `db:"headimg" json:"headimg" description:"头像"`
	Sex          *int    `db:"sex" json:"sex" description:"性别"`
	Birthday     *string `db:"birthday" json:"birthday" description:"出生日期"`
	Company      *string `db:"company" json:"company" description:"公司"`
	Address      *string `db:"address" json:"address" description:"地址"`
	City         *string `db:"city" json:"city" description:"城市"`
	Province     *string `db:"province" json:"province" description:"省份"`
	Country      *string `db:"country" json:"country" description:"国家"`
	Regip        *string `db:"regip" json:"regip" description:"注册IP地址"`
	Referer      *string `db:"referer" json:"referer" description:"用户来源"`
	Inviter      *int64  `db:"inviter" json:"inviter" description:"邀请者"`
	Userclan     *string `db:"userclan" json:"userclan" description:"用户拓谱图"`
	FName        *string `db:"fname" json:"fname" description:"真实姓名"`
	Bankcard     *string `db:"bankcard" json:"bankcard" description:"银行卡"`
	IdentityCard *string `db:"identity_card" json:"identity_card" description:"身份证"`
	GroupTag     *string `db:"grouptag" json:"grouptag" description:"用户组"`
	Remark       *string `db:"remark" json:"remark" description:"备注"`
	Object       *string `db:"object" json:"object" description:"预留字段"`
	State        *int    `db:"state" json:"state" description:"状态"`
	Intime       *string `db:"intime" json:"intime" description:"注册时间"`
	Uptime       *string `db:"uptime" json:"uptime" description:"更新时间"`
}

// // 使用结构体 将密码字段转换为布尔类型
// // 结构体可以这样定义
// // {
// // 	Ciphers     string  `db:"ciphers" json:"-"`
// // 	HasPassword  bool    `db:"-" json:"ciphers"`
// // }
// user.HasPassword = (user.Ciphers != "")
// // 数据库时如果有字段为NULL，而当为NULL时程序报错 converting NULL to string is unsupported
// // 字段为NULL时可以结构体类型前面加星号 *string 指向字符串类型的指针一劳永逸，但可能会稍微影响性能和增加复杂性
// // 也可以是使用database/sql包里的 sql.NullString 类型来接收查询结果中有NULL的字段
// // 而基于gorm库的不需要考虑这种为NULL情况，因为gorm库本身已经做了这种兼容处理。

// 平台结构体
type UserOAuth struct {
	OID           int64  `db:"oid" json:"oid" description:"自增ID"`
	UID           int64  `db:"uid" json:"uid" description:"用户ID"`
	Platfrom      string `db:"platfrom" json:"platfrom" description:"外接平台名"`
	Openid        string `db:"openid" json:"openid" description:"外接平台身份ID"`
	Headimg       string `db:"headimg" json:"headimg" description:"头像"`
	Unionid       string `db:"unionid" json:"unionid" description:"外接唯一标识"`
	NickName      string `db:"nickname" json:"nickname" description:"昵称"`
	Sex           int    `db:"sex" json:"sex" description:"性别"`
	City          string `db:"city" json:"city" description:"城市"`
	Province      string `db:"province" json:"province" description:"省份"`
	Country       string `db:"country" json:"country" description:"国家"`
	Language      string `db:"language" json:"language" description:"语言"`
	Privilege     string `db:"privilege" json:"privilege" description:"用户特权信息"`
	Token         string `db:"token" json:"token" description:"平台接口授权凭证"`
	Expires       string `db:"expires" json:"expires" description:"平台授权失效时间"`
	Refresh       string `db:"refresh" json:"refresh" description:"平台刷新token"`
	Scope         string `db:"scope" json:"scope" description:"用户授权的作用域"`
	Subscribe     int    `db:"subscribe" json:"subscribe" description:"是否关注"`
	SubscribeTime string `db:"subscribetime" json:"subscribetime" description:"用户拓谱图"`
	GroupTag      string `db:"grouptag" json:"grouptag" description:"用户组"`
	Tidings       string `db:"tidings" json:"tidings" description:"用户动态"`
	Remark        string `db:"remark" json:"remark" description:"备注"`
	Object        string `db:"object" json:"object" description:"预留字段"`
	Intime        string `db:"intime" json:"intime" description:"入库时间"`
	Uptime        string `db:"uptime" json:"uptime" description:"更新时间"`
}

// 附属资料结构体
type UserMaterial struct {
	UID     int64   `db:"uid" json:"uid" description:"用户ID"`
	CID     int64   `db:"cid" json:"cid" description:"联系人ID"`
	Balance float64 `db:"balance" json:"balance" description:"余额"`
	Vip     int     `db:"vip" json:"vip" description:"会员"`
	Exptime string  `db:"exptime" json:"exptime" description:"会员到期时间"`
	Manage  int     `db:"manage" json:"manage" description:"管理权限"`
	Tag     string  `db:"tag" json:"tag" description:"权限标识"`
	Remark  string  `db:"remark" json:"remark" description:"备注"`
	Object  string  `db:"object" json:"object" description:"预留字段"`
	Uptime  string  `db:"uptime" json:"uptime" description:"更新时间"`
}

// 用户日志结构体
type UserLogs struct {
	ID     int64  `db:"id" json:"id" description:"日志ID"`
	UID    int64  `db:"uid" json:"uid" description:"用户ID"`
	Action string `db:"action" json:"action" description:"操作标识"`
	Note   string `db:"note" json:"note" description:"操作说明"`
	Actip  string `db:"actip" json:"actip" description:"操作IP地址"`
	Ua     string `db:"ua" json:"ua" description:"设备信息"`
	Intime string `db:"intime" json:"intime" description:"更新时间"`
}

// Create（创建）在数据库中创建一个新的用户
func (m *UserModel) Create(data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(UserInfo{}, data)

	// 设置初始值 初始化入库数据 (手动校验或添加其他表单项字段）

	// 判断是否存在字段 "cell"
	if _, ok := data["cell"]; ok {
		if data["cell"] != "" {
			if m.Check("cell", data["cell"].(string)) {
				return 0, fmt.Errorf("用户cell电话已存在")
			}
		}
	}

	// 判断是否存在字段 "email"
	if _, ok := data["email"]; ok {
		if data["email"] != "" {
			if m.Check("email", data["email"].(string)) {
				return 0, fmt.Errorf("用户email邮箱已存在")
			}
		}
	}

	// 判断是否存在字段 "identity_card"
	if _, ok := data["identity_card"]; ok {
		if data["identity_card"] != "" {
			if m.Check("identity_card", data["identity_card"].(string)) {
				return 0, fmt.Errorf("用户identity_card身份证已存在")
			}
		}
	}

	// 判断是否存在字段 "uid"
	if _, ok := data["uid"]; ok {
		if data["uid"] == "" {
			data["uid"] = utils.GenerateTimerID()
		}
	} else {
		data["uid"] = utils.GenerateTimerID()
	}

	// 入库时间
	_, intime := utils.FormatTimestamp(time.Now().Unix())
	data["intime"] = intime
	data["uptime"] = intime

	// 判断是否存在字段 "ciphers"
	if _, ok := data["ciphers"]; ok {
		if data["ciphers"] != "" {
			data["ciphers"] = utils.HashPassword(data["ciphers"].(string), data["intime"].(string))
		}
	}

	// 判断是否存在字段 "object"
	if _, ok := data["object"]; !ok {
		data["object"] = ""
	}

	// 判断是否存在字段 "sex"
	if _, ok := data["sex"]; ok {
		if data["sex"] == "" {
			delete(data, "sex") // 删除性别字段
		}
	}

	// 判断是否存在字段 "inviter"
	if _, ok := data["inviter"]; ok {
		if data["inviter"] == "" {
			delete(data, "inviter") // 删除邀请字段
		} else {
			intv := utils.ParseInt64(data["inviter"]) // 任意数据转int64数字
			// 生成用户拓谱图
			clan := m.GenerateUserClan(intv)
			data["userclan"] = strings.Join(clan, ",")
		}
	}

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 循环遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`")
		values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", m.TableName, strings.Join(fields, ","), strings.Join(values, ","))
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// utils.PrintExtSql(sql, args) // 打印最终执行的SQL语句

	// 执行数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return 0, err
	}

	// 获取插入结果
	id, err := database.LastInsertId() // 新插入数据的id
	if err != nil {
		println("LastInsertId failed: ", err.Error())
		return 0, err
	}

	return id, nil // 返回结果
}

// Update（更新）根据用户ID更新数据库中的用户信息
func (m *UserModel) Update(id int64, data map[string]interface{}) (int64, error) {
	// 实现更新数据库中用户信息的逻辑

	// 按结构体映射提交字段
	data = utils.StructAssigMap(UserInfo{}, data)

	// 没有更新数据项直接返回
	if len(data) == 0 {
		return 0, nil
	}

	// 判断是否存在字段 "cell"
	if _, ok := data["cell"]; ok {
		if data["cell"] != "" {
			if m.Check("cell", data["cell"].(string)) {
				return 0, fmt.Errorf("用户cell电话已存在")
			}
		}
	}

	// 判断是否存在字段 "email"
	if _, ok := data["email"]; ok {
		if data["email"] != "" {
			if m.Check("email", data["email"].(string)) {
				return 0, fmt.Errorf("用户email邮箱已存在")
			}
		}
	}

	// 判断是否存在字段 "identity_card"
	if _, ok := data["identity_card"]; ok {
		if data["identity_card"] != "" {
			if m.Check("identity_card", data["identity_card"].(string)) {
				return 0, fmt.Errorf("用户identity_card身份证已存在")
			}
		}
	}

	// 更新时间
	_, uptime := utils.FormatTimestamp(time.Now().Unix())
	data["uptime"] = uptime

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	// values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`=:"+key)
		// values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `uid`=%d", m.TableName, strings.Join(fields, ","), id)
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行插入数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return 0, err
	}

	// 获取插入结果
	row, err := database.RowsAffected() // 更新行数
	if err != nil {
		println("LastInsertId failed: ", err.Error())
		return 0, err
	}

	return row, nil
}

// Delete（删除）根据用户ID从数据库中删除用户
// data返回当前用户ID则是逻辑删除、返回1则是物理删除
func (m *UserModel) Delete(id int64) (int64, error) {

	// 先查询当前此用户
	user, err := m.Read(id)
	if err != nil {
		println("m.Read Error: ", err.Error())
		return 0, err
	}

	// 判断当前用户状态(0停用 1正常 2待激活)
	if *user.State == 0 {
		// 物理删除  delete 删除
		sql := fmt.Sprintf("DELETE FROM `%s` WHERE `uid` = ?", m.TableName)
		resDelete, errDelete := m.DB.Exec(sql, id)
		if errDelete != nil {
			println("数据库delete删除失败，", errDelete.Error())
			return 0, errDelete
		}
		rowDelete, errDelete := resDelete.RowsAffected() // 操作影响的行数
		if errDelete != nil {
			println("rows failed, ", errDelete.Error())
			return 0, errDelete
		}
		return rowDelete, nil

	} else {
		// 更新数据
		_, upErr := m.Update(id, map[string]interface{}{
			"state": 0,
		})
		if upErr != nil {
			println("m.Update Error: ", upErr.Error())
			return 0, upErr
		}
		return id, nil
	}
}

// Read（读取）获取用户信息 - 根据用户ID从数据库中获取用户
func (m *UserModel) Read(id int64) (UserInfo, error) {
	// SQL注入问题：我们任何时候都不应该自己拼接SQL语句！

	// 拼接 GET 的 select 查询语句
	fields := "`uid`,`username`,`email`,`cell`,`nickname`,`headimg`,`sex`,`birthday`,`company`,`address`,`city`,`province`,`country`,`grouptag`,`state`,`intime`,`uptime`,`ciphers`,`regip`,`referer`,`inviter`,`userclan`,`fname`,`bankcard`,`identity_card`,`remark`,`object`"
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `uid`=? LIMIT 1", fields, m.TableName)
	// println("\r\n",sql)                    // 打印sql

	var user UserInfo
	err := m.DB.Get(&user, sql, id) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return user, err
	}

	// println("数据查询成功：", user)
	return user, nil
}

// 检查是否存在用户信息 - 根据字段名和值从数据库中获取用户
func (m *UserModel) Check(fieldName string, fieldValue string) bool {
	// existCell := c.Models.Check("cell", "13838389438")

	// 拼接 GET 的 select 查询语句
	sql := "SELECT `uid` FROM `" + m.TableName + "` WHERE `" + fieldName + "`=? LIMIT 1"
	// println("\r\n", sql) // 打印sql

	var uid int64
	err := m.DB.Get(&uid, sql, fieldValue) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return false
	}

	// println("数据查询成功：", uid)
	return true
}

// 获取用户列表
func (m *UserModel) List(filters map[string]interface{}, pageNumber, pageSize int64, pageOrder, pageField string) ([]UserInfo, int64, error) {
	// 按结构体映射提交字段
	filters = utils.StructAssigMap(UserInfo{}, filters)

	var infos []UserInfo
	var total int64
	var order string

	// 判断是否存在字段 "pageOrder"
	if pageOrder != "" {
		_temp := strings.ToUpper(pageOrder) // 转大写
		// 检查字符串是否包含
		if strings.Contains(_temp, "DESC") {
			pageOrder = "DESC"
		} else {
			pageOrder = "ASC"
		}
	}
	// 判断是否存在字段 "pageField"
	if pageField != "" {
		pageField = strings.TrimSpace(pageField) // 去除字符串前后空格
		order = fmt.Sprintf(" ORDER BY `%s` %s", pageField, pageOrder)
	}

	pageNumber = pageNumber - 1 // 转为数据库语义，因为数据库limit是从第0行开始的

	// 防止分页页码小于0
	if pageNumber < 0 {
		pageNumber = 0
	}

	starNum := pageNumber * pageSize
	endNum := starNum + pageSize
	// println("pageNumber/pageSize", pageNumber, "/", pageSize)
	// println("pageField/pageOrder", pageField, "/", pageOrder)

	// 获取过滤条件数据（表达式::值）
	where_arr, args := utils.GetWhereArgs(filters)
	// println(strings.Join(where_arr, " "))
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 拼接 NamedQuery 的 select 查询语句
	query := fmt.Sprintf("FROM `%s` WHERE 1=1 %s", m.TableName, strings.Join(where_arr, " "))
	orderANDlimit := fmt.Sprintf("%s LIMIT %d,%d", order, starNum, endNum)
	// 拼接 GET 的 select 查询语句
	sql := "SELECT `uid`,`username`,`email`,`cell`,`nickname`,`headimg`,`sex`,`birthday`,`company`,`address`,`city`,`province`,`country`,`grouptag`,`state`,`intime`,`uptime`,`ciphers`,`regip`,`referer`,`inviter`,`userclan`,`fname`,`bankcard`,`identity_card`,`remark`,`object` " + query + orderANDlimit
	count_sql := "SELECT COUNT(*) AS `count` " + query
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)
	// println("\r\n", count_sql) // 打印 count_sql

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args) // 不支持?号的方式

	// db.Queryx：这个方法与 db.Query 类似，但它返回的是一个 sqlx.Rows 对象，该对象具有比标准库 sql.Rows 更强大的功能。sqlx.Rows 支持结构体映射、更方便的字段访问方法等，能够简化数据提取过程。适用于需要更灵活的结果处理和数据提取的场景。
	// rows, err := m.DB.Queryx(sql, uid, cid) // 支持?号的方式

	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return infos, total, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var info UserInfo
		if err := rows.StructScan(&info); err != nil {
			return nil, total, err
		}
		// 处理字段

		// 数据处理 - 转换时间格式 不用脱敏 直接拿到秒级精确时间
		*info.Birthday = utils.RFC3339ToString(*info.Birthday)
		*info.Intime = utils.RFC3339ToString(*info.Intime)
		*info.Uptime = utils.RFC3339ToString(*info.Uptime)
		infos = append(infos, info)
	}

	// 当为第一页时并且总行数据小于分页条数 直接返回总数
	total = int64(len(infos))
	if pageNumber == 0 && total < pageSize {
		return infos, total, nil
	}

	// 查询总数
	count_rows, count_err := m.DB.NamedQuery(count_sql, args)
	if count_err != nil {
		println("NamedQuery failed: ", count_err.Error())
	}
	defer count_rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for count_rows.Next() {
		if err := count_rows.Scan(&total); err != nil {
			return nil, total, err
		}
	}
	// println("总行数:", total)
	return infos, total, nil
}

// 根据用户ID获取用户附属资料表
func (m *UserModel) ReadMaterial(uid int64) (UserMaterial, error) {
	// SQL注入问题：我们任何时候都不应该自己拼接SQL语句！

	// 拼接 GET 的 select 查询语句
	fields := "`uid`,`cid`,`balance`,`vip`,`exptime`,`manage`,`tag`,`remark`,`object`,`uptime`"
	sql := "SELECT " + fields + " FROM `" + m.MaterialTableName + "` WHERE `uid`=? LIMIT 1"
	// println("\r\n", sql) // 打印sql

	var material UserMaterial
	err := m.DB.Get(&material, sql, uid) // 查询单行数据
	if err != nil {
		println("Err: ", err.Error())
		return material, err
	}

	// fmt.Println("数据查询成功：%v", material)
	return material, nil
}

// 添加用户附属资料表
func (m *UserModel) CreateMaterial(data map[string]interface{}) error {
	// 打印模块名
	// println("___________CreateMaterial_____________")
	// println("platfrom", data["platfrom"].(string))

	// 判断是否存在字段 "uid"
	if _, ok := data["uid"]; !ok {
		return fmt.Errorf("uid用户ID不能为空")
	}

	// 设置初始值 初始化入库数据 (手动添加其他表单项字段）

	// 判断是否存在字段 "object"
	if _, ok := data["object"]; !ok {
		data["object"] = ""
	}

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`")
		values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", m.MaterialTableName, strings.Join(fields, ","), strings.Join(values, ","))
	// println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行数据库的操作
	_, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return err
	}

	// // 获取插入结果
	//  LastInsertId() 方法依赖于数据库自动生成的主键值。如果您的表没有指定主键字段或主键字段没有设置自增属性，那么 LastInsertId() 方法将无法返回有效的主键ID。
	// id, err := database.LastInsertId() // 新插入数据的id(必须要返回主键？)
	// if err != nil {
	// 	// println("LastInsertId failed: ", err.Error())
	// 	return 0, err
	// }
	// println("-----------", id)

	return nil // 返回空结果
}

// 查找用户 (使用username、email、cell、identity_card查找用户)
func (m *UserModel) Find(name string) (UserInfo, string, error) {
	// 打印模块名
	// println("___________Find_____________")
	// println("platfrom", data["platfrom"].(string))

	var user UserInfo

	typeName := "username"

	if utils.CheckEmail(name) {
		typeName = "email"
	} else if utils.CheckMobile(name) {
		typeName = "cell"
	} else if utils.CheckIdCard(name) {
		typeName = "identity_card"
	}
	// println(where)

	// 拼接 GET 的 select 查询语句
	fields := "`uid`,`username`,`email`,`cell`,`nickname`,`headimg`,`sex`,`birthday`,`company`,`address`,`city`,`province`,`country`,`grouptag`,`state`,`intime`,`uptime`,`ciphers`,`regip`,`referer`,`inviter`,`userclan`,`fname`,`bankcard`,`identity_card`,`remark`,`object`"
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `%s`=? LIMIT 1", fields, m.TableName, typeName)

	// println("\r\n", sql) // 打印sql

	err := m.DB.Get(&user, sql, name) // 查询单行数据，也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return user, typeName, err
	}

	return user, typeName, nil // 返回结果
}

// 记录操作日志
func (m *UserModel) SetLogs(data map[string]interface{}) error {

	// 判断是否存在字段 "uid"
	if _, ok := data["uid"]; !ok {
		return fmt.Errorf("uid用户ID不能为空")
	}

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`")
		values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", m.LogsTableName, strings.Join(fields, ","), strings.Join(values, ","))
	// println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行数据库的操作
	_, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return err
	}

	// // 获取插入结果
	//  LastInsertId() 方法依赖于数据库自动生成的主键值。如果您的表没有指定主键字段或主键字段没有设置自增属性，那么 LastInsertId() 方法将无法返回有效的主键ID。
	// id, err := database.LastInsertId() // 新插入数据的id(必须要返回主键？)
	// if err != nil {
	// 	// println("LastInsertId failed: ", err.Error())
	// 	return 0, err
	// }
	// println("-----------", id)

	return nil // 返回空结果
}

// 获取操作日志
func (m *UserModel) GetLogs(filters map[string]interface{}, pageNumber, pageSize int64, pageOrder, pageField string) ([]UserLogs, int64, error) {
	// 按结构体映射提交字段
	filters = utils.StructAssigMap(UserLogs{}, filters)

	var logs []UserLogs
	var total int64
	var order string

	// 判断是否存在字段 "pageOrder"
	if pageOrder != "" {
		_temp := strings.ToUpper(pageOrder) // 转大写
		// 检查字符串是否包含
		if strings.Contains(_temp, "DESC") {
			pageOrder = "DESC"
		} else {
			pageOrder = "ASC"
		}
	}
	// 判断是否存在字段 "pageField"
	if pageField != "" {
		pageField = strings.TrimSpace(pageField) // 去除字符串前后空格
		order = fmt.Sprintf(" ORDER BY `%s` %s", pageField, pageOrder)
	}

	pageNumber = pageNumber - 1 // 转为数据库语义，因为数据库limit是从第0行开始的

	// 防止分页页码小于0
	if pageNumber < 0 {
		pageNumber = 0
	}

	starNum := pageNumber * pageSize
	endNum := starNum + pageSize
	// println("pageNumber/pageSize", pageNumber, "/", pageSize)
	// println("pageField/pageOrder", pageField, "/", pageOrder)

	// 获取过滤条件数据（表达式::值）
	where_arr, args := utils.GetWhereArgs(filters)
	println(strings.Join(where_arr, " "))
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 拼接 NamedQuery 的 select 查询语句
	query := fmt.Sprintf("FROM `%s` WHERE 1=1 %s", m.LogsTableName, strings.Join(where_arr, " "))
	orderANDlimit := fmt.Sprintf("%s LIMIT %d,%d", order, starNum, endNum)
	// 拼接 GET 的 select 查询语句
	sql := "SELECT `id`,`uid`,`action`,`note`,`actip`,`ua`,`intime` " + query + orderANDlimit
	count_sql := "SELECT COUNT(*) AS `count` " + query
	// println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)
	// println("\r\n", count_sql) // 打印 count_sql

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return logs, total, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var log UserLogs
		if err := rows.StructScan(&log); err != nil {
			return nil, total, err
		}

		// 处理字段
		log.Intime = utils.RFC3339ToString(log.Intime) // 转换时间格式

		logs = append(logs, log)
	}

	// 当为第一页时并且总行数据小于分页条数 直接返回总数
	total = int64(len(logs))
	if pageNumber == 0 && total < pageSize {
		return logs, total, nil
	}

	// 查询总数
	count_rows, count_err := m.DB.NamedQuery(count_sql, args)
	if count_err != nil {
		println("NamedQuery failed: ", count_err.Error())
	}
	defer count_rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for count_rows.Next() {
		if err := count_rows.Scan(&total); err != nil {
			return nil, total, err
		}
	}
	// println("总行数:", total)
	return logs, total, nil
}

// 生成用户拓谱图
func (m *UserModel) GenerateUserClan(uid int64) []string {
	// println("====GenerateUserClan===")
	var clan string
	level := 10 // 级别 最多有多少级的上级
	var result []string
	result = append(result, strconv.FormatInt(uid, 10)) // 数字转字符串 并添加到切片中

	sql := fmt.Sprintf("SELECT `userclan` FROM `%s` WHERE `uid`=? LIMIT 1;", m.TableName)
	// println("\r\n", sql)             // 打印sql
	err := m.DB.Get(&clan, sql, uid) // 查询单行数据
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return result
	}

	clans := strings.Split(clan, ",") // 分割字符串
	// strings.Join(clans, ","), // 拼接字符串
	level -= 1
	// 按最高等级数添加拓扑级别
	for key, value := range clans {
		// println("GenerateUserClan:",key, value)
		if key < level {
			result = append(result, value)
		}
	}

	return result
}

/*
公用方法
首先在在别的控制器里添加 UserModel *models.UserModel 的初始化代码，例如：

	type AccountController struct {
		DB        *sqlx.DB
		Models    *models.AccountModel
		CTX       iris.Context
		UserModel *models.UserModel
	}

	func NewAccountController(db *sqlx.DB) *AccountController {
		// 返回一个结构体指针
		return &AccountController{
			DB:        db,
			Models:    models.NewAccountModel(db),
			UserModel: models.NewUserModel(db),
		}
	}

然后再到控制器方法里添加使用，例如：
	// 如果操作人不是自己
	if tkUid != uid {
		// 如果不是管理员
		if !c.UserModel.IsAdmin(tkUid) {
			ctx.JSON(iris.Map{"code": config.ErrUnauthorized, "msg": config.ErrMsgs[config.ErrUnauthorized]})
			return
		}
	}

	// 如果超出最大数量 -  默认3记事本，vip不限制数目记事本
	if total > max {
		// 判断是否为VIP
		vipLeve, v_err := c.UserModel.IsVip(tkUid)
		if v_err != nil {
			ctx.JSON(iris.Map{"code": config.ErrNotVip, "msg": config.ErrMsgs[config.ErrNotVip]})
			return
		}
		println(vipLeve)
	}
*/

// 是否是管理员
func (m *UserModel) IsAdmin(uid int64) bool {
	// println("====IsAdmin===")
	// 调取模型
	material, err := m.ReadMaterial(uid)
	if err != nil {
		println("m.ReadMaterial Error: ", err.Error())
		return false
	}
	if material.Manage == 1 {
		return true
	}
	return false
}

// 是否是VIP并返回级别
func (m *UserModel) IsVip(uid int64) (int, error) {
	// println("====IsVip===")
	// 调取模型
	material, err := m.ReadMaterial(uid)
	if err != nil {
		println("m.ReadMaterial Error: ", err.Error())
		return 0, err
	}
	// fmt.Println("material:%v", material)
	// println(":::vip", material.Vip)
	// 会员等级
	if material.Vip > 0 {
		// println("是会员")
		// 会员到期时间 - 转换时间格  直接拿到秒级精确时间
		material.Exptime = utils.RFC3339ToString(material.Exptime)
		_, exp := utils.ParseTimeToTimestamp(material.Exptime)
		// println(time.Now().Unix(), exp)
		// println(time.Now().Unix() > exp)
		// 检查过期时间是否超时 1685084839  >  1685109600
		if time.Now().Unix() > exp {
			// println("VIP已于", material.Exptime, "过期")
			return material.Vip, fmt.Errorf("VIP已于 %s 过期", material.Exptime)
		}
		return material.Vip, nil // 返回VIP级别
	}
	return 0, fmt.Errorf("您还不是VIP")
}

// 接入用户 (接入第三方平台用户)
func (m *UserModel) CreateOAuth(data map[string]interface{}) (UserOAuth, error) {
	// 先判断该用户platfrom_openid是否已经在如果不存在去查platfrom_unionID是否存在则更新该用户信息
	// 如果platfrom_openid和platfrom_unionidf都不存在，则将用户信息插入表中并返回用户信息

	var useroauth UserOAuth
	var platfrom, openid, unionid string

	// 按结构体映射提交字段
	data = utils.StructAssigMap(UserOAuth{}, data)

	// 打印模块名
	// println("___________CreateOAuth_____________")
	// println("platfrom", data["platfrom"].(string))

	// // 判断是否存在字段 "uid"
	// if _, ok := data["uid"]; !ok {
	// 	return 0, fmt.Errorf("uid用户ID不能为空")
	// }

	// 判断是否存在字段 "platfrom"
	if _, ok := data["platfrom"]; !ok {
		return useroauth, fmt.Errorf("platfrom外接平台名不能为空")
	} else {
		platfrom = data["platfrom"].(string)
	}

	// 判断是否存在字段 "openid"
	if _, ok := data["openid"]; !ok {
		return useroauth, fmt.Errorf("openid外接平台身份ID不能为空")
	} else {
		openid = data["openid"].(string)
	}

	// 判断是否存在字段 "unionid"
	if _, ok := data["unionid"]; ok {
		unionid = data["unionid"].(string)
	}

	// 设置初始值 初始化入库数据 (手动添加其他表单项字段）

	// 判断是否存在字段 "privilege"
	if _, ok := data["privilege"]; !ok {
		data["privilege"] = ""
	}
	// 判断是否存在字段 "tidings"
	if _, ok := data["tidings"]; !ok {
		data["tidings"] = ""
	}
	// 判断是否存在字段 "object"
	if _, ok := data["object"]; !ok {
		data["object"] = ""
	}

	useroauth, err := m.FindOAuthOpenid(platfrom, openid)
	if err == nil {
		return useroauth, nil
	}

	// 如果找到 unionid 相符的用户
	useroauth, err = m.FindOAuthUnionid(platfrom, unionid)
	if err == nil {
		// 绑定同一用户id并入库
		data["uid"] = useroauth.UID
	}

	// 生成入库ID，防止自增让人猜出平台使用量
	data["oid"] = utils.GenerateTimerID(9999) // 四位随机数

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`")
		values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", m.OauthTableName, strings.Join(fields, ","), strings.Join(values, ","))
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return useroauth, err
	}

	// 获取插入结果
	id, err := database.LastInsertId() // 新插入数据的id
	if err != nil {
		// println("LastInsertId failed: ", err.Error())
		return useroauth, err
	}
	// println("oid:", id)
	data["oid"] = id

	jsonStr, _ := json.Marshal(data)
	var person UserOAuth
	err = json.Unmarshal(jsonStr, &person)
	if err != nil {
		// 处理错误
		return person, err
	}
	return person, nil // 返回结果
}

// 查找平台openid用户
func (m *UserModel) FindOAuthOpenid(platfrom, openid string) (UserOAuth, error) {
	// 先判断该用户platfrom_openid是否已经在如果不存在去查platfrom_unionID是否存在则更新该用户信息
	// 如果platfrom_openid和platfrom_unionidf都不存在，则将用户信息插入表中并返回用户信息

	var useroauth UserOAuth

	// 拼接 GET 的 select 查询语句
	fields := "`oid`,`uid`,`platfrom`,`openid`,`unionid`,`nickname`,`headimg`,`city`,`province`,`country`,`grouptag`,`language`,`intime`,`uptime`,`privilege`,`token`,`expires`,`refresh`,`scope`,`subscribe`,`subscribetime`,`tidings`,`remark`,`object`"
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `openid`=? AND `platfrom`=? LIMIT 1", fields, m.OauthTableName)
	// println("\r\n", sql) // 打印sql

	err := m.DB.Get(&useroauth, sql, openid, platfrom) // 查询单行数据，也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return useroauth, err
	}

	useroauth.Intime = utils.RFC3339ToString(useroauth.Intime, 2)
	useroauth.Uptime = utils.RFC3339ToString(useroauth.Uptime, 2)
	return useroauth, nil // 返回结果
}

// 查找平台unionid用户
func (m *UserModel) FindOAuthUnionid(platfrom, unionid string) (UserOAuth, error) {
	// 先判断该用户platfrom_openid是否已经在如果不存在去查platfrom_unionID是否存在则更新该用户信息
	// 如果platfrom_openid和platfrom_unionidf都不存在，则将用户信息插入表中并返回用户信息

	var useroauth UserOAuth

	// 拼接 GET 的 select 查询语句
	fields := "`oid`,`uid`,`platfrom`,`openid`,`unionid`,`nickname`,`headimg`,`city`,`province`,`country`,`grouptag`,`language`,`intime`,`uptime`,`privilege`,`token`,`expires`,`refresh`,`scope`,`subscribe`,`subscribetime`,`tidings`,`remark`,`object`"
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `unionid`=? AND `platfrom`=? LIMIT 1", fields, m.OauthTableName)
	// println("\r\n", sql) // 打印sql

	err := m.DB.Get(&useroauth, sql, unionid, platfrom) // 查询单行数据，也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return useroauth, err
	}

	useroauth.Intime = utils.RFC3339ToString(useroauth.Intime, 2)
	useroauth.Uptime = utils.RFC3339ToString(useroauth.Uptime, 2)
	return useroauth, nil // 返回结果
}

// 查找oid平台用户
func (m *UserModel) FindOAuthOid(oid int64) (UserOAuth, error) {
	var useroauth UserOAuth

	// 拼接 GET 的 select 查询语句
	fields := "`oid`,`uid`,`platfrom`,`openid`,`unionid`,`nickname`,`headimg`,`city`,`province`,`country`,`grouptag`,`language`,`intime`,`uptime`,`privilege`,`token`,`expires`,`refresh`,`scope`,`subscribe`,`subscribetime`,`tidings`,`remark`,`object`"
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `oid`=? LIMIT 1", fields, m.OauthTableName)
	// println("\r\n", sql) // 打印sql

	err := m.DB.Get(&useroauth, sql, oid) // 查询单行数据，也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return useroauth, err
	}

	useroauth.Intime = utils.RFC3339ToString(useroauth.Intime, 2)
	useroauth.Uptime = utils.RFC3339ToString(useroauth.Uptime, 2)

	return useroauth, nil // 返回结果
}

// Update（更新）根据开放平台OID更新数据库中的用户的第三方信息
func (m *UserModel) UpdateOAuth(oid int64, data map[string]interface{}) (int64, error) {
	// 实现更新数据库中用户信息的逻辑

	// 按结构体映射提交字段
	data = utils.StructAssigMap(UserOAuth{}, data)

	// 没有更新数据项直接返回
	if len(data) == 0 {
		return 0, nil
	}

	// 更新时间
	_, uptime := utils.FormatTimestamp(time.Now().Unix())
	data["uptime"] = uptime

	// 使用 make() 函数来创建切片，Go语言切片是对数组的抽象
	fields := make([]string, 0)
	// values := make([]string, 0)
	args := make(map[string]interface{})

	// 获取用户提交的所有表单项字段 遍历数据
	for key, value := range data {
		fields = append(fields, "`"+key+"`=:"+key)
		// values = append(values, ":"+key)
		args[key] = value
	}

	// 构建数据库的SQL语句
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `oid`=%d", m.OauthTableName, strings.Join(fields, ","), oid)
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行插入数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// 处理数据库操作错误
		// println("NamedExec failed: ", err.Error())
		return 0, err
	}

	// 获取插入结果
	row, err := database.RowsAffected() // 更新行数
	if err != nil {
		println("LastInsertId failed: ", err.Error())
		return 0, err
	}

	return row, nil
}
