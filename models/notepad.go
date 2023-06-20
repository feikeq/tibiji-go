package models

import (
	"fmt"
	"strings"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type NotepadModel struct {
	DB        *sqlx.DB
	TableName string
}

func NewNotepadModel(db *sqlx.DB) *NotepadModel {
	return &NotepadModel{
		DB:        db,
		TableName: "tbj_notepad",
	}
}

// 记事本结构体
type NotepadInfo struct {
	Nid     int64  `db:"nid" json:"nid" description:"纸张ID"`
	UID     int64  `db:"uid" json:"uid" description:"用户ID"`
	Url     string `db:"url" json:"url" description:"访问地址"`
	Share   string `db:"share" json:"share" description:"共享地址"`
	Content string `db:"content" json:"content" description:"内容"`
	Pwd     string `db:"pwd" json:"pwd" description:"密码"`
	Caret   int    `db:"caret" json:"caret" description:"光标位置"`
	Scroll  int    `db:"scroll" json:"scroll" description:"滚动位置"`
	IP      string `db:"ip" json:"ip" description:"IP"`
	Referer string `db:"referer" json:"referer" description:"纸张来源"`
	State   int    `db:"state" json:"state" description:"状态"`
	Intime  string `db:"intime" json:"intime" description:"创建时间"`
	Uptime  string `db:"uptime" json:"uptime" description:"更新时间"`
}

func (m *NotepadModel) Create(data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(NotepadInfo{}, data)

	// 判断是否存在字段 "content"
	if _, ok := data["content"]; !ok {
		data["content"] = ""
	}

	// 入库时间
	_, intime := utils.FormatTimestamp(time.Now().Unix())
	data["intime"] = intime
	data["uptime"] = intime
	// 生成入库ID，防止自增让人猜出平台使用量
	data["nid"] = utils.GenerateTimerID(9999) // 四位随机数
	data["share"] = utils.GenerateShortId()   // 生成独特的非连续短ID

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

func (m *NotepadModel) Update(uid int64, id int64, data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(NotepadInfo{}, data)

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
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `nid`=%d AND `uid`=%d", m.TableName, strings.Join(fields, ","), id, uid)
	println("\r\n", sql)                            // 打印sql
	fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行数据库的操作
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

// 物理删除
func (m *NotepadModel) Delete(uid int64, id int64) (int64, error) {
	// 物理删除  delete 删除
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE `nid` = ? AND `uid`=?", m.TableName)
	resDelete, errDelete := m.DB.Exec(sql, id, uid)
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
}

// 所有纸张
func (m *NotepadModel) List(id int64) ([]NotepadInfo, error) {

	var infos []NotepadInfo

	// 拼接 GET 的 select 查询语句
	fields := "`nid`,`uid`,`url`,`share`,`pwd`,`ip`,`referer`,`state`,`intime`,`uptime` "
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `uid`= ? ", fields, m.TableName)
	// println("\r\n", sql) // 打印sql

	// err := m.DB.Select(&infos, sql, id) // 查询单行数据，也可以用 NamedQuery
	rows, err := m.DB.Queryx(sql, id) // 支持?号的方式
	if err != nil {
		println("Err: ", err.Error())
		return infos, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var info NotepadInfo
		if err := rows.StructScan(&info); err != nil {
			return nil, err
		}
		// 处理字段

		// 数据处理 - 转换时间格式脱敏时间
		info.Intime = utils.RFC3339ToString(info.Intime, 2)
		info.Uptime = utils.RFC3339ToString(info.Uptime, 2)
		infos = append(infos, info)
	}
	return infos, nil
}

// 查找云纸张
func (m *NotepadModel) Find(url string) (NotepadInfo, error) {
	// SQL注入问题：我们任何时候都不应该自己拼接SQL语句！

	// 拼接 GET 的 select 查询语句
	fields := "`nid`,`uid`,`url`,`share`,`content`,`pwd`,`caret`,`scroll`,`ip`,`referer`,`state`,`intime`,`uptime` "
	sql := fmt.Sprintf("SELECT %s FROM `%s` WHERE `url` = ? LIMIT 1", fields, m.TableName)
	println("\r\n", sql) // 打印sql

	var user NotepadInfo
	err := m.DB.Get(&user, sql, url) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return user, err
	}

	return user, nil
}

// Read（读取）云纸张
func (m *NotepadModel) Read(url string) (string, error) {

	// 拼接 GET 的 select 查询语句
	sql := fmt.Sprintf("SELECT `content` FROM `%s` WHERE `share`=? AND `state` = 1 LIMIT 1", m.TableName)
	// println("\r\n", sql) // 打印sql

	var txt string
	err := m.DB.Get(&txt, sql, url) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return txt, err
	}
	return txt, nil
}

// 检查云纸张数量
func (m *NotepadModel) Check(uid int64) int64 {
	// 拼接 GET 的 select 查询语句
	sql := fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `uid` =? ", m.TableName)
	// println("\r\n", sql) // 打印sql
	var total int64
	err := m.DB.Get(&total, sql, uid) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return total
	}
	return total
}
