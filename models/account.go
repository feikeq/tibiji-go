package models

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
)

type AccountModel struct {
	DB               *sqlx.DB
	TableName        string
	ContactTableName string
}

func NewAccountModel(db *sqlx.DB) *AccountModel {
	return &AccountModel{
		DB:               db,
		TableName:        "tbj_account",
		ContactTableName: "tbj_contact",
	}
}

// 帐单结构体
type AccountInfo struct {
	AID      int64   `db:"aid" json:"aid" description:"自动id"`
	UID      int64   `db:"uid" json:"uid" description:"用户ID"`
	Item     string  `db:"item" json:"item" description:"操作项目"`
	Class    string  `db:"class" json:"class" description:"操作项目"`
	Sort     string  `db:"sort" json:"sort" description:"子类别"`
	CID      int64   `db:"cid" json:"cid" description:"收支对象ID"`
	Object   string  `db:"object" json:"object" description:"收支对象"`
	Accounts string  `db:"accounts" json:"accounts" description:"操作账户"`
	Money    float64 `db:"money" json:"money" description:"金额"`
	Note     string  `db:"note" json:"note" description:"备注说明"`
	Btime    string  `db:"btime" json:"btime" description:"帐单时间"`
	Intime   string  `db:"intime" json:"intime" description:"入库时间"`
	Uptime   string  `db:"uptime" json:"uptime" description:"更新时间"`
}

// 帐单日历结构体
type AccountCalendar struct {
	Year  int     `db:"year" json:"year" description:"年份"`
	Month int     `db:"month" json:"month" description:"月份"`
	Inc   float64 `db:"inc" json:"inc" description:"收入"`
	Out   float64 `db:"out" json:"out" description:"支出"`
	Oth   float64 `db:"oth" json:"oth" description:"其它"`
}

// 分组结构体
type AccountGroup struct {
	Day    string        `json:"day" description:"分组标题"`
	Moment string        `json:"moment" description:"时间片刻"`
	Inc    float64       `json:"inc" description:"收入"`
	Out    float64       `json:"out" description:"支出"`
	Oth    float64       `json:"oth" description:"其他"`
	List   []AccountInfo `json:"list" description:"帐单列表"`
}

// 分类统计结构体
type AccountCount struct {
	Sort  string  `json:"sort" description:"分类"`
	Total float64 `json:"total" description:"合计"`
}

// 常用KV对象结构体
type KVObject struct {
	Key int64  `json:"key" description:"键"`
	Val string `json:"val" description:"值"`
	Str string `json:"str" description:"项"`
}

func (m *AccountModel) Create(data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(AccountInfo{}, data)

	// 入库时间
	_, intime := utils.FormatTimestamp(time.Now().Unix())
	data["intime"] = intime
	data["uptime"] = intime
	// 生成入库ID，防止自增让人猜出平台使用量
	data["aid"] = utils.GenerateTimerID(9999) // 四位随机数

	// 判断是否存在字段 "btime"
	if _, ok := data["btime"]; ok {
		if data["btime"] == "" {
			data["btime"] = intime
		}
	} else {
		data["btime"] = intime
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
	// println("\r\n", sql)            // 打印sql
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

func (m *AccountModel) Update(uid int64, id int64, data map[string]interface{}) (int64, error) {
	// 按结构体映射提交字段
	data = utils.StructAssigMap(AccountInfo{}, data)

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
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `aid`=%d AND `uid`=%d", m.TableName, strings.Join(fields, ","), uid, id)
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

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

// 查找记录的用户ID
func (m *AccountModel) Find(id int64) (int64, error) {

	// 拼接 GET 的 select 查询语句
	sql := "SELECT `uid` FROM `" + m.TableName + "` WHERE `aid`=? LIMIT 1"

	var uid int64
	err := m.DB.Get(&uid, sql, id) // 查询单行数据，也可以用 NamedQuery
	if err != nil {
		println("Err: ", err.Error())
		return uid, err
	}

	return uid, nil // 返回结果
}

// 获取列表
func (m *AccountModel) List(uid int64, filters map[string]interface{}, pageNumber, pageSize int64, pageOrder, pageField string) ([]AccountInfo, int64, error) {
	// 按结构体映射提交字段
	filters = utils.StructAssigMap(AccountInfo{}, filters)

	var accounts []AccountInfo
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
	query := fmt.Sprintf("FROM `%s` WHERE `uid`=%d %s", m.TableName, uid, strings.Join(where_arr, " "))
	orderANDlimit := fmt.Sprintf("%s LIMIT %d,%d", order, starNum, endNum)
	// 拼接 GET 的 select 查询语句
	sql := "SELECT `aid`,`uid`,`item`,`class`,`sort`,`cid`,`object`,`accounts`,`money`,`note`,`btime`,`intime`,`uptime` " + query + orderANDlimit
	count_sql := "SELECT COUNT(*) AS `count` " + query
	// println("\r\n", sql)            // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)
	// println("\r\n", count_sql) // 打印 count_sql

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return accounts, total, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var account AccountInfo
		if err := rows.StructScan(&account); err != nil {
			return nil, total, err
		}
		// 处理字段

		// 数据处理 - 转换时间格式 进行脱敏
		account.Btime = utils.RFC3339ToString(account.Btime, 2)
		account.Intime = utils.RFC3339ToString(account.Intime, 2)
		account.Uptime = utils.RFC3339ToString(account.Uptime, 2)
		accounts = append(accounts, account)
	}

	// 当为第一页时并且总行数据小于分页条数 直接返回总数
	total = int64(len(accounts))
	if pageNumber == 0 && total < pageSize {
		return accounts, total, nil
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
	return accounts, total, nil
}

// 帐单类目
func (m *AccountModel) Types() map[string]string {
	// 定义帐单类目
	var accountTypes = map[string]string{
		"早午晚餐": "基本生活",
		"水果零食": "基本生活",
		"日常用品": "基本生活",
		"柴米油盐": "基本生活",
		"水电煤气": "基本生活",
		"房租物业": "基本生活",
		"医药保健": "基本生活",

		"服饰装扮": "衣服饰品",
		"鞋帽手套": "衣服饰品",
		"饰品包包": "衣服饰品",
		"化妆美容": "衣服饰品",

		"公共交通": "交通通讯",
		"打车租车": "交通通讯",
		"私车供养": "交通通讯",
		"话费网费": "交通通讯",
		"邮递快递": "交通通讯",

		"旅游度假": "休闲娱乐",
		"休闲玩乐": "休闲娱乐",
		"朋友聚会": "休闲娱乐",
		"运动健身": "休闲娱乐",
		"宠物宝贝": "休闲娱乐",
		"博彩彩票": "休闲娱乐",

		"书报影音": "文化进修",
		"数码装备": "文化进修",
		"教育培训": "文化进修",

		"人际往来": "人情往来",
		"孝敬长辈": "人情往来",
		"婚嫁礼金": "人情往来",
		"生日礼金": "人情往来",
		"节日礼金": "人情往来",
		"礼品礼金": "人情往来",
		"慈善捐助": "人情往来",

		"投资理财": "其他杂项",
		"金融保险": "其他杂项",
		"电器家居": "其他杂项",
		"家政服务": "其他杂项",
		"房屋房产": "其他杂项",
		"车产船产": "其他杂项",
		"其他杂项": "其他杂项",

		"工资收入": "职业收入",
		"补助津贴": "职业收入",
		"加班收入": "职业收入",
		"公务报销": "职业收入",
		"奖金收入": "职业收入",
		"兼职收入": "职业收入",
		"投资收入": "职业收入",
		"经营收入": "职业收入",

		"利息收入": "其他收入",
		"中奖收入": "其他收入",
		"意外来财": "其他收入",
		"其他收益": "其他收入",
	}
	return accountTypes
}

// 月份帐单
func (m *AccountModel) Month(uid int64, year, month string) ([]AccountGroup, error) {
	/*
		[
			{
				"day": "2012-09-24",
				"moment": "今天 2012-09-24/ 昨天 2012-09-24 / 2012-09-24星期几",
				"inc": 0.0, // 收入
				"out": 4.6, // 支出
				"oth": 0.0, // 其他
				"list": [
					{
						"aid": 16856698287568482,
						"uid": 8,
						"item": "收入",
						"class": "人情往来",
						"sort": "节日礼金",
						"cid": 0,
						"object": "张三",
						"accounts": "银行卡",
						"money": 889,
						"note": "过生日呀",
						"btime": "2020-11-12 18:18",
						"intime": "2023-06-02 09:37",
						"uptime": "2023-06-02 09:37"
					}
				]
			}
		]
	*/

	sql := "SELECT `aid`,`uid`,`item`,`class`,`sort`,`cid`,`object`,`accounts`,`money`,`note`,`btime`,`intime`,`uptime` "
	sql += "FROM `" + m.TableName + "` WHERE `uid` = ? AND YEAR(`intime`)= ? AND MONTH(`intime`) = ? "
	sql += "ORDER BY 'btime' DESC;"
	// println(sql)

	var data []AccountInfo
	err := m.DB.Select(&data, sql, uid, year, month) // 查询多行数据
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return nil, err
	}

	// 使用结构体拼装数据项组装数据遍历进行数据组装 - 根据KEY值放入map对象进行分组

	// 映射是一个指向值的指针，方便之后使用时检查键值是否存在。
	list := make(map[string]*AccountGroup) // 这里如果不使用指针会遇到过错误“无法分配给地图中的结构字段”

	// 遍历数据
	for _, value := range data {

		// fmt.Printf("===key: %+v\n", key)
		// fmt.Printf("value: %+v\n", value)

		// 数据处理 - 转换时间格式 进行脱敏
		bTime := value.Btime // 缓存变量之后还需要用
		value.Btime = utils.RFC3339ToString(bTime, 2)
		value.Intime = utils.RFC3339ToString(value.Intime, 2)
		value.Uptime = utils.RFC3339ToString(value.Uptime, 2)

		// 分组标题
		title := utils.RFC3339ToString(bTime, 0)
		// 友好时间
		dayStr, dateStr, weekStr := utils.FriendlyDate(bTime)
		momentStr := dateStr + " " + weekStr
		if dayStr != "" {
			momentStr = dayStr + " " + momentStr
		}
		// 检查键值是否存在
		if _, ok := list[title]; !ok {
			list[title] = &AccountGroup{
				Day:    title,
				Moment: momentStr,
				Inc:    0.0,
				Out:    0.0,
				Oth:    0.0,
				List:   make([]AccountInfo, 0),
			}
		}

		// 我们可以访问并重新分配 Group 结构对象的值
		// 进行必要的调整后，可以将此 thisProduct 分配回同一个指针。
		list[title].List = append(list[title].List, value) // 往MAP里添加数据

		if value.Item == "收入" {
			list[title].Inc += value.Money
		} else if value.Item == "支出" {
			list[title].Out += value.Money
		} else {
			list[title].Oth += value.Money
		}
	}

	// 拼成的数据是 {"2022-12-12":Group{},"2022-12-13":Group{}}
	// fmt.Printf("list: %+v\n", list)

	// 重新把数据变为 [Group{},Group{}]
	groups := make([]AccountGroup, 0)
	for _, v := range list {
		groups = append(groups, *v)
	}

	// 对切片进行排序 - 根据Day字段倒序排列groups切片
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Day > groups[j].Day
	})

	return groups, nil
}

// 帐单日历
func (m *AccountModel) Calendar(uid int64) ([]AccountCalendar, error) {

	/*
		SELECT
			YEAR(`intime`) AS 年份,MONTH(`intime`) AS 月份,
			SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 收入,
			SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 支出,
			SUM(CASE WHEN `item` NOT IN ('收入', '支出') THEN `money` ELSE 0 END) AS 其他
		FROM  `tbj_account` WHERE `uid` = 6  AND YEAR(`intime`) =2022
		GROUP BY YEAR(intime), MONTH(intime)
		ORDER BY 年份, 月份;
	*/

	// SELECT YEAR(`intime`) 'Y' , MONTH(`intime`) 'M' , `item`, sum(`money`) 'SR' FROM tbj_account WHERE `uid`=6  AND YEAR(`intime`) =2022 GROUP BY  YEAR(`intime`) , MONTH(`intime`),`item`  ORDER BY `intime`
	sql := "SELECT YEAR(`btime`) AS 'year',MONTH(`btime`) AS 'month',"
	sql += "SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 'inc', "
	sql += "SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 'out', "
	sql += "SUM(CASE WHEN `item` NOT IN ('收入', '支出') THEN `money` ELSE 0 END) AS 'oth' "
	sql += "FROM `" + m.TableName + "` WHERE `uid` = ? "
	sql += "GROUP BY YEAR(btime), MONTH(btime) ORDER BY 'year', 'month';"
	// println(sql)
	// println(m.TableName, id)

	var data []AccountCalendar
	err := m.DB.Select(&data, sql, uid) // 查询多行数据
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return nil, err
	}

	return data, nil
}

// 删除帐单 (物理删除)
func (m *AccountModel) Delete(uid int64, id int64) (int64, error) {
	// 物理删除  delete 删除
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE `aid` = ? AND `uid`=?", m.TableName)
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

// 收支比例
func (m *AccountModel) ReportRatio(uid int64, year, month int) (AccountCalendar, error) {

	// 构建数据库的SQL语句
	sql := "SELECT "
	sql += "SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 'inc', "
	sql += "SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 'out', "
	sql += "SUM(CASE WHEN `item` NOT IN ('收入', '支出') THEN `money` ELSE 0 END) AS 'oth' "
	sql += fmt.Sprintf(" FROM `%s` WHERE `uid`=? ", m.TableName)

	if year > 0 {
		sql += " AND YEAR(`btime`)=? "

	}
	if month > 0 {
		sql += " AND MONTH(`btime`)=? "
	}
	if year > 0 {
		if month > 0 {
			sql += " GROUP BY YEAR(btime), MONTH(btime) "
		} else {
			sql += " GROUP BY YEAR(btime) "
		}
	}

	// println("::::", sql)

	var err error
	// var data []map[string]interface{}
	var data AccountCalendar
	if year > 0 {
		data.Year = year
		if month > 0 {
			err = m.DB.Get(&data, sql, uid, year, month)
			data.Month = month
		} else {
			err = m.DB.Get(&data, sql, uid, year)
		}
	} else {
		err = m.DB.Get(&data, sql, uid)
	}
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}

	return data, nil
}

// 收支列表
func (m *AccountModel) ReportRatios(uid int64, limit int) ([]AccountCalendar, error) {

	// 构建数据库的SQL语句
	sql := "SELECT YEAR(btime) AS `year`, MONTH(btime) AS `month`,"
	sql += "SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 'inc', "
	sql += "SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 'out', "
	sql += "SUM(CASE WHEN `item` NOT IN ('收入', '支出') THEN `money` ELSE 0 END) AS 'oth' "
	sql += fmt.Sprintf(" FROM `%s` WHERE `uid`=? ", m.TableName)
	sql += "GROUP BY `year`, `month` ORDER BY `year` DESC , `month` DESC "
	if limit > 0 {
		sql += " LIMIT ?"
	}

	// println(limit, "::::", sql)

	var err error

	var data []AccountCalendar
	if limit > 0 {
		err = m.DB.Select(&data, sql, uid, limit)
	} else {
		err = m.DB.Select(&data, sql, uid)
	}
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}

	return data, nil
}

// 收支明细
func (m *AccountModel) ReportDetails(uid int64, item string, year, month int) ([]AccountCount, error) {
	/*
		SELECT `sort`,SUM(`money`) AS `total` FROM `tbj_account`
		WHERE `uid`=6 AND `item`='收入'
		AND YEAR(`btime`)=2023 AND MONTH(`btime`) = 6
		GROUP BY `sort` ORDER BY `total` DESC;
	*/
	// 构建数据库的SQL语句
	sql := fmt.Sprintf("SELECT `sort`,SUM(`money`) AS `total` FROM `%s` ", m.TableName)
	sql += "WHERE `uid`=? AND `item`=? "

	if year > 0 {
		sql += "AND YEAR(`btime`)=? "
	}
	if month > 0 {
		sql += "AND MONTH(`btime`)=? "
	}
	sql += "GROUP BY `sort` ORDER BY `total` DESC "

	println("::::", sql)

	var err error
	// var data []map[string]interface{}
	var data []AccountCount
	if year > 0 {
		if month > 0 {
			err = m.DB.Select(&data, sql, uid, item, year, month)
		} else {
			err = m.DB.Select(&data, sql, uid, item, year)
		}
	} else {
		err = m.DB.Select(&data, sql, uid, item)
	}
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}

	return data, nil
}

// 流水账户
func (m *AccountModel) ReportAccounts(uid int64, year, month int) ([]AccountCount, error) {
	/*
		SELECT `accounts` AS 'sort' ,SUM(`money`) 'total' FROM `tbj_account`
		WHERE `uid`=6  AND YEAR(`btime`)=2022  AND MONTH(`btime`)=6
		GROUP BY `accounts` ORDER BY `total` DESC;
	*/
	// 构建数据库的SQL语句
	sql := fmt.Sprintf("SELECT `accounts` AS 'sort' ,SUM(`money`) 'total' FROM `%s` ", m.TableName)
	sql += "WHERE `uid`=? "
	if year > 0 {
		sql += "AND YEAR(`btime`)=? "
	}
	if month > 0 {
		sql += "AND MONTH(`btime`)=? "
	}
	sql += "GROUP BY `accounts` ORDER BY `total` DESC "

	println("::::", sql)

	var err error
	// var data []map[string]interface{}
	var data []AccountCount
	if year > 0 {
		if month > 0 {
			err = m.DB.Select(&data, sql, uid, year, month)
		} else {
			err = m.DB.Select(&data, sql, uid, year)
		}
	} else {
		err = m.DB.Select(&data, sql, uid)
	}
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}

	return data, nil
}

// 收支对象
func (m *AccountModel) Objects(uid int64) ([]KVObject, error) {
	/*
		两个表 tbj_account 和 tbj_contact ，tbj_account 表的主键为 aid 外键为 cid ， 而 tbj_contact 表的主键为 cid ，
		现需要用 uid = 6 查询 tbj_contact 的数据并按 tbj_account 表里已有CID的 btime 时间倒序排序，否则按 ，tbj_account 里的ID倒序排序。
		tbj_contact表里的数据没有重复项，而tbj_account里的数据有重复项，这里去要对tbj_account表里的cid去重后再按tbj_account表里的btime时间倒序排序tbj_contact表里的数据。
		也就是一个表中的数据按另一个表中的某个字段来决定顺序：

		SELECT c.`cid`,c.`uid`,c.`fullname` FROM `tbj_contact` c
		LEFT JOIN (
		    SELECT `cid`,`uid`, MAX(`btime`) AS `btime`
		    FROM `tbj_account`
		    WHERE `uid` = 6 AND `cid` !=0
				GROUP BY `cid`
				ORDER BY `btime` DESC
		) AS a ON c.`cid` = a.`cid`
		WHERE c.`uid`=6
	*/

	sql := fmt.Sprintf("SELECT c.`cid` AS `key`,c.`fullname` AS `val` FROM `%s` c ", m.ContactTableName)
	sql += fmt.Sprintf("LEFT JOIN ( SELECT `cid`,`uid`, MAX(`btime`) AS `btime` FROM `%s` ", m.TableName)
	sql += "WHERE `uid` = ? AND `cid` !=0 GROUP BY `cid` ORDER BY `btime` DESC ) "
	sql += "AS a ON c.`cid` = a.`cid` WHERE c.`uid`=? "
	println(":::::", sql)

	var data []KVObject
	err := m.DB.Select(&data, sql, uid, uid)
	if err != nil {
		println("数据库select查询失败，", err.Error())
		return data, err
	}
	return data, nil
}
