package models

import (
	"fmt"
	"sort"
	"strings"
	"tibiji-go/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type RemindModel struct {
	DB                *sqlx.DB
	TableName         string
	UserTableName     string
	MaterialTableName string
}

func NewRemindModel(db *sqlx.DB) *RemindModel {
	return &RemindModel{
		DB:                db,
		TableName:         "tbj_contact",
		UserTableName:     "sys_user",
		MaterialTableName: "sys_material",
	}
}

// 提醒结构体
type RemindInfo struct {
	Cid      int64  `db:"cid" json:"cid" description:"用户ID"`
	UID      int64  `db:"uid" json:"uid" description:"用户ID"`
	Fullname string `db:"fullname" json:"fullname" description:"姓名"`
	Pinyin   string `db:"pinyin" json:"pinyin" description:"拼音"`
	NickName string `db:"nickname" json:"nickname" description:"昵称绰号"`
	Picture  string `db:"picture" json:"picture" description:"相片照片"`
	Gender   int    `db:"gender" json:"gender" description:"性别"`
	Birthday string `db:"birthday" json:"birthday" description:"出生时间"`
	Lunar    int    `db:"lunar" json:"lunar" description:"是否为农历"`
	GroupTag string `db:"grouptag" json:"grouptag" description:"分组"`
	Remind   string `db:"remind" json:"remind" description:"提醒方式"`
	Relation string `db:"relation" json:"relation" description:"关系"`
	Note     string `db:"note" json:"note" description:"备注"`
	State    int    `db:"state" json:"state" description:"状态类别"`

	// 数据加工处理生成的附加字段
	RemindNum      int    `json:"remind_num" description:"剩余天数"`
	RemindDate     string `json:"remind_date" description:"生日的日期 YYYY-MM-DD"`
	RemindCNDate   string `json:"remind_cndate" description:"生日的日期 X月X日"`
	RemindYear     int    `json:"remind_year" description:"年龄或周年"`
	RemindStar     string `json:"remind_star" description:"星座"`
	RemindZodiac   string `json:"remind_zodiac" description:"生肖"`
	BirthdayCNDate string `json:"birthday_cndate" description:"出生中文日期"`
	RemindDay      string `json:"remind_day" description:"今天明天后天N天后"`
	RemindWeek     string `json:"remind_week" description:"星期几"`
}

// 提醒模式结构体
type RemindOBJ struct {
	Email   []string `json:"email" description:"邮件"`
	Phone   []string `json:"phone" description:"手机"`
	Notice  []string `json:"notice" description:"消息推送"`
	Message []string `json:"message" description:"聊天消息"`
}

// 提醒数据结构体
type RemindInfoTask struct {
	RemindInfo
	RemindType RemindOBJ `json:"remind_type" description:"提醒方式对象"`
	// 用户数据
	UserFname    string  `db:"user_fname" json:"user_fname" description:"真实姓名"`
	UserNickname string  `db:"user_nickname" json:"user_nickname" description:"昵称"`
	UserUsername string  `db:"user_username" json:"user_username" description:"帐号"`
	UserEmail    string  `db:"user_email" json:"user_email" description:"邮箱"`
	UserCell     string  `db:"user_cell" json:"user_cell" description:"电话"`
	UserBalance  float64 `db:"user_balance" json:"user_balance" description:"余额"`
	UserVip      int     `db:"user_vip" json:"user_vip" description:"会员"`
	UserExptime  string  `db:"user_exptime" json:"user_exptime" description:"会员到期时间"`
}

// 提醒队列结构体
type RemindQueues struct {
	Email   []RemindInfoTask `json:"email" description:"邮件队列"`
	Phone   []RemindInfoTask `json:"phone" description:"手机队列"`
	Notice  []RemindInfoTask `json:"notice" description:"消息推送队列"`
	Message []RemindInfoTask `json:"message" description:"聊天消息队列"`
}

func (m *RemindModel) RemindToOBJ(str string) RemindOBJ {
	// 提醒方式字段格式化成对象数组 "email::7,1,0||phone::7,1,0"
	arr_type := strings.Split(str, "||") // 分割字符串 - 拆分(表达式:值)
	var obj RemindOBJ
	if len(arr_type) > 0 {
		// 循环遍历数据
		for _, value := range arr_type {
			arr_item := strings.Split(value, "::") // 分割字符串 - 拆分(表达式:值)
			// fmt.Println(arr_item)
			if len(arr_type) > 0 {
				switch arr_item[0] {
				case "email":
					obj.Email = strings.Split(arr_item[1], ",")
				case "phone":
					obj.Phone = strings.Split(arr_item[1], ",")
				case "notice":
					obj.Notice = strings.Split(arr_item[1], ",")
				default:
					obj.Message = strings.Split(arr_item[1], ",")
				}
			}

		}
	}
	return obj
}

// 个人的提醒列表
func (m *RemindModel) Items(uid int64, filters map[string]interface{}) ([]RemindInfo, error) {

	// 按结构体映射提交字段
	filters = utils.StructAssigMap(ContactInfo{}, filters)

	var infos []RemindInfo

	// 获取过滤条件数据（表达式::值）
	where_arr, args := utils.GetWhereArgs(filters)
	// println(strings.Join(where_arr, " "))
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 拼接 NamedQuery 的 select 查询语句 (注意 0000-00-00 00:00:00 在sqlx库中的NamedQuery函数不支持冒号（:）作为命名参数的一部分，因为它会将其解释为命名参数的分隔符。解决方法是使用问号（?）作为命名参数的占位符，然后将命名参数作为第二个参数传递给函数。)
	query := fmt.Sprintf("FROM `%s` WHERE `uid`=%d AND `state` != 0 AND `birthday` > '0000-00-00' %s", m.TableName, uid, strings.Join(where_arr, " "))

	// 拼接 GET 的 select 查询语句
	sql := "SELECT `cid`,`uid`,`fullname`,`pinyin`,`nickname`,`picture`,`gender`,`birthday`,`lunar`,`grouptag`,`remind`,`relation`,`note`,`state` " + query
	// println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return infos, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var info RemindInfo
		if err := rows.StructScan(&info); err != nil {
			return nil, err
		}
		info.RemindNum, info.RemindDay, info.RemindDate, info.RemindWeek, info.RemindYear, info.RemindStar, info.RemindZodiac, info.RemindCNDate, info.BirthdayCNDate = utils.BirthdayAndDaysStarZodiac(utils.RFC3339ToTime(info.Birthday), info.Lunar)
		// println(info.RemindNum, info.RemindDay, info.RemindDate, info.RemindWeek, info.RemindYear, info.RemindStar, info.RemindZodiac, info.RemindCNDate, info.BirthdayCNDate)

		// 数据处理 - 转换时间格式 脱敏
		info.Birthday = utils.RFC3339ToString(info.Birthday, 2)

		infos = append(infos, info)
	}

	// 对切片进行排序 - 根据Day字段从小到大正序排列infos切片
	sort.Slice(infos, func(i, j int) bool {
		// 多字段排序 它将按 RemindNum 从低到高排，然后按 RemindYear 年龄从大到小排。
		if infos[i].RemindNum == infos[j].RemindNum {
			return infos[i].RemindYear > infos[j].RemindYear
		}
		return infos[i].RemindNum < infos[j].RemindNum
	})

	return infos, nil
}

func (m *RemindModel) Task() (RemindQueues, error) {
	now := time.Now()
	// ps:全局生成通知列表时可以用当天的农历时间去筛选一次数据库减少样本（因为农历总小于阳历）
	// 取当前时间转为农历(阴历)
	lunarDate, _ := utils.DateToLunar(now) // 公历转农历(阳历转阴历)
	// 前面添0补0 前置补零(前导零)，保证月和日始终是2位数的长度
	timeStartStr := fmt.Sprintf("%02d-%02d", int(lunarDate.Month()), lunarDate.Day())

	// ps:同时只取10天内的数据 - 当前日期加10天
	afterTenDays := now.AddDate(0, 0, 10) // 添加10天，往后十天
	timeEndStr := fmt.Sprintf("%02d-%02d", int(afterTenDays.Month()), afterTenDays.Day())

	var infos RemindQueues

	/*
		SELECT c.`cid`,c.`uid`,c.`fullname`,c.`pinyin`,c.`nickname`,c.`picture`,c.`gender`,c.`birthday`,c.`lunar`,c.`grouptag`,c.`remind`,c.`relation`,c.`note`,c.`state`,
		u.`fname` 'user_fname', u.`nickname` 'user_nickname', u.`username` 'user_username', u.`email` 'user_email',u.`cell` 'user_cell',
		m.`balance` 'user_balance',m.`vip` 'user_vip',m.`exptime` 'user_exptime'
		FROM `tbj_contact` c, `sys_user` u, `sys_material` m
		WHERE c.`state` != 0 AND c.`birthday` > '0000-00-00' AND DATE_FORMAT(c.`birthday`,'%m-%d') >= '12-21'
		AND u.`uid` = c.`uid` AND m.`uid` = c.`uid`
	*/

	// 拼接 NamedQuery 的 select 查询语句 (注意 0000-00-00 00:00:00 在sqlx库中的NamedQuery函数不支持冒号（:）作为命名参数的一部分，因为它会将其解释为命名参数的分隔符。解决方法是使用问号（?）作为命名参数的占位符，然后将命名参数作为第二个参数传递给函数。)
	query := fmt.Sprintf("FROM `%s` c, `%s` u, `%s` m ", m.TableName, m.UserTableName, m.MaterialTableName)
	query += "WHERE c.`state` != 0 AND c.`birthday` > '0000-00-00' "
	query += "AND DATE_FORMAT(c.`birthday`,'%m-%d') >= ? AND DATE_FORMAT(c.`birthday`,'%m-%d') < ? "
	query += "AND u.`uid` = c.`uid` AND m.`uid` = c.`uid` "

	// 拼接 GET 的 select 查询语句
	sql := "SELECT c.`cid`,c.`uid`,c.`fullname`,c.`pinyin`,c.`nickname`,c.`picture`,c.`gender`,c.`birthday`,c.`lunar`,c.`grouptag`,c.`remind`,c.`relation`,c.`note`,c.`state`,"
	sql += "u.`fname` 'user_fname', u.`nickname` 'user_nickname', u.`username` 'user_username', u.`email` 'user_email',u.`cell` 'user_cell',"
	sql += "m.`balance` 'user_balance',m.`vip` 'user_vip',m.`exptime` 'user_exptime' "
	sql += query
	println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容
	// fmt.Printf("Type: %T , Data: %v\n", args, args)
	println("队列区间：", timeStartStr, timeEndStr)

	// db.Queryx：这个方法与 db.Query 类似，但它返回的是一个 sqlx.Rows 对象，该对象具有比标准库 sql.Rows 更强大的功能。sqlx.Rows 支持结构体映射、更方便的字段访问方法等，能够简化数据提取过程。适用于需要更灵活的结果处理和数据提取的场景。
	rows, err := m.DB.Queryx(sql, timeStartStr, timeEndStr)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return infos, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池
	// 遍历查询结果每一行
	for rows.Next() {
		var info RemindInfoTask
		if err := rows.StructScan(&info); err != nil {
			return infos, err
		}
		// println(info.Cid, "|", info.UID, "|", info.Fullname)

		// 预处理字段
		info.RemindNum, info.RemindDay, info.RemindDate, info.RemindWeek, info.RemindYear, info.RemindStar, info.RemindZodiac, info.RemindCNDate, info.BirthdayCNDate = utils.BirthdayAndDaysStarZodiac(utils.RFC3339ToTime(info.Birthday), info.Lunar)
		// println(info.RemindNum, info.RemindDay, info.RemindDate, info.RemindWeek, info.RemindYear, info.RemindStar, info.RemindZodiac, info.RemindCNDate, info.BirthdayCNDate)

		// 数据处理 - 转换时间格式 脱敏
		info.Birthday = utils.RFC3339ToString(info.Birthday, 2)
		info.UserExptime = utils.RFC3339ToString(info.Birthday, 2)

		// println()
		// println(info.RemindNum, info.Birthday)

		// 提醒时间碰撞处理 - 判断剩余天数是否与提醒日期相符
		info.RemindType = m.RemindToOBJ(info.Remind) // 只管理前端自己定义数据结构即可没必要后端来定这格式
		// fmt.Println(info.RemindType)                  // {Email:[7 1 0],Phone:[7 1 0]}
		// Email 邮件队列
		if utils.InArray(info.RemindType.Email, info.RemindNum) {
			println("Email:")
			println(info.RemindNum, info.Birthday)
			println(info.Cid, "|", info.UID, "|", info.Fullname)
			infos.Email = append(infos.Email, info)
		}
		// Phone 手机队列
		if utils.InArray(info.RemindType.Phone, info.RemindNum) {
			println("Phone:")
			println(info.RemindNum, info.Birthday)
			println(info.Cid, "|", info.UID, "|", info.Fullname)
			infos.Phone = append(infos.Phone, info)
		}
		// Notice 消息推送队列
		// infos.Notice = append(infos.Notice, info)
		// Message 聊天消息队列
		// infos.Message = append(infos.Message, info)
	}

	return infos, nil
}
