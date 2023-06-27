package models

import (
	"fmt"
	"strings"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type ContactModel struct {
	DB        *sqlx.DB
	TableName string
}

func NewContactModel(db *sqlx.DB) *ContactModel {
	return &ContactModel{
		DB:        db,
		TableName: "tbj_contact",
	}
}

// 联系人结构体
type ContactInfo struct {
	Cid      int64  `db:"cid" json:"cid" description:"联系人ID"`
	UID      int64  `db:"uid" json:"uid" description:"用户ID"`
	Fullname string `db:"fullname" json:"fullname" description:"姓名"`
	Pinyin   string `db:"pinyin" json:"pinyin" description:"拼音"`
	NickName string `db:"nickname" json:"nickname" description:"昵称绰号"`
	Picture  string `db:"picture" json:"picture" description:"相片照片"`
	Phone    string `db:"phone" json:"phone" description:"电话与传真"`
	Mail     string `db:"mail" json:"mail" description:"邮箱"`
	IM       string `db:"im" json:"im" description:"聊天帐号"`
	Http     string `db:"http" json:"http" description:"网址"`
	Company  string `db:"company" json:"company" description:"公司部门"`
	Position string `db:"position" json:"position" description:"职位头衔"`
	Address  string `db:"address" json:"address" description:"地址"`
	Gender   int    `db:"gender" json:"gender" description:"性别"`
	Birthday string `db:"birthday" json:"birthday" description:"生日时间"`
	Lunar    int    `db:"lunar" json:"lunar" description:"是否为农历"`
	GroupTag string `db:"grouptag" json:"grouptag" description:"分组"`
	Remind   string `db:"remind" json:"remind" description:"提醒方式"`
	Relation string `db:"relation" json:"relation" description:"关系"`
	Family   string `db:"family" json:"family" description:"家庭户主"`
	Note     string `db:"note" json:"note" description:"备注"`
	State    int    `db:"state" json:"state" description:"状态"`
	Intime   string `db:"intime" json:"intime" description:"创建时间"`
	Uptime   string `db:"uptime" json:"uptime" description:"更新时间"`
}

func (m *ContactModel) Create(data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(ContactInfo{}, data)

	// 生成入库ID，防止自增让人猜出平台使用量
	data["cid"] = utils.GenerateTimerID(9999) // 四位随机数

	// 判断是否存在字段 "fullname"
	if _, ok := data["fullname"]; ok {
		if data["fullname"] != "" {
			data["pinyin"] = utils.PinYin(data["fullname"].(string))
		}
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
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", m.TableName, strings.Join(fields, ","), strings.Join(values, ","))
	// println("\r\n", sql)                            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// utils.PrintExtSql(sql, args) // 打印最终执行的SQL语句

	// 执行数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		// println("NamedExec failed: ", err.Error())
		// 处理数据库操作错误
		// ctx.StatusCode(iris.StatusInternalServerError)
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

func (m *ContactModel) Update(uid, id int64, data map[string]interface{}) (int64, error) {
	// 实现更新数据库中用户信息的逻辑

	// 按结构体映射提交字段
	data = utils.StructAssigMap(ContactInfo{}, data)

	// 没有更新数据项直接返回
	if len(data) == 0 {
		return 0, nil
	}

	// 判断是否存在字段 "fullname"
	if _, ok := data["fullname"]; ok {
		if data["fullname"] != "" {
			data["pinyin"] = utils.PinYin(data["fullname"].(string))
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
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `cid`=%d AND `uid`=%d ", m.TableName, strings.Join(fields, ","), id, uid)
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行插入数据库的操作
	database, err := m.DB.NamedExec(sql, args)
	if err != nil {
		println("NamedExec failed: ", err.Error())
		// 处理数据库操作错误
		// ctx.StatusCode(iris.StatusInternalServerError)
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

func (m *ContactModel) Delete(uid, id int64) (int64, error) {

	// 非物理删除更新数据
	_, upErr := m.Update(uid, id, map[string]interface{}{
		"state": 0,
	})
	if upErr != nil {
		println("m.Update Error: ", upErr.Error())
		return 0, upErr
	}
	return id, nil

}

// 检查是否存在
func (m *ContactModel) Check(fieldName string, fieldValue string) bool {
	// existCell := c.Models.Check("cell", "13838389438")

	// 拼接 GET 的 select 查询语句
	sql := fmt.Sprintf("SELECT `uid` FROM `%s` WHERE `%s`=? LIMIT 1", m.TableName, fieldName)

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

func (m *ContactModel) List(uid int64, filters map[string]interface{}, pageNumber, pageSize int64, pageOrder, pageField string) ([]ContactInfo, int64, error) {
	// 按结构体映射提交字段
	filters = utils.StructAssigMap(ContactInfo{}, filters)

	var infos []ContactInfo
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
	query := fmt.Sprintf("FROM `%s` WHERE `uid` = %d AND state != 0 %s ", m.TableName, uid, strings.Join(where_arr, " "))
	orderANDlimit := fmt.Sprintf("%s LIMIT %d,%d", order, starNum, endNum)
	// 拼接 GET 的 select 查询语句
	sql := "SELECT `cid`,`uid`,`fullname`,`pinyin`,`nickname`,`picture`,`phone`,`mail`,`im`,`http`,`company`,`position`,`address`,`gender`,`birthday`,`lunar`,`grouptag`,`remind`,`relation`,`family`,`note`,`state`,`intime`,`uptime` " + query + orderANDlimit
	count_sql := "SELECT COUNT(*) AS `count` " + query
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)
	// println("\r\n", count_sql) // 打印 count_sql

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return infos, total, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var info ContactInfo
		if err := rows.StructScan(&info); err != nil {
			return nil, total, err
		}
		// 处理字段

		// 数据处理 - 转换时间格式 不用脱敏 直接拿到秒级精确时间
		info.Birthday = utils.RFC3339ToString(info.Birthday, 2)
		info.Intime = utils.RFC3339ToString(info.Intime, 2)
		info.Uptime = utils.RFC3339ToString(info.Uptime, 2)
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

// 分组对象
func (m *ContactModel) Groups(uid int64) ([]KVObject, error) {
	sql := fmt.Sprintf("SELECT `grouptag`AS 'str',`grouptag` AS 'val' FROM `%s` ", m.TableName)
	sql += "WHERE `uid`=? AND `grouptag` !='' GROUP BY `grouptag`"
	var data []KVObject
	err := m.DB.Select(&data, sql, uid)
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}
	return data, nil
}
