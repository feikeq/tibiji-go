package models

import (
	"fmt"
	"sort"
	"strings"
	"tibiji-go/utils"

	"github.com/jmoiron/sqlx"
)

type HumaneModel struct {
	DB               *sqlx.DB
	TableName        string
	ContactTableName string
}

func NewHumaneModel(db *sqlx.DB) *HumaneModel {
	return &HumaneModel{
		DB:               db,
		TableName:        "tbj_account",
		ContactTableName: "tbj_contact",
	}
}

// 人情帐单结构体
type HumaneInfo struct {
	CID      int64  `db:"cid" json:"cid" description:"收支对象ID"`
	Fullname string `db:"fullname" json:"fullname" description:"姓名"`
	Pinyin   string `db:"pinyin" json:"pinyin" description:"拼音"`
	Picture  string `db:"picture" json:"picture" description:"相片照片"`
	HumaneItems
}

// 简单帐单结构体
type HumaneItems struct {
	Item     string  `db:"item" json:"item" description:"操作项目"`
	Class    string  `db:"class" json:"class" description:"主分类"`
	Sort     string  `db:"sort" json:"sort" description:"子类别"`
	Object   string  `db:"object" json:"object" description:"收支对象"`
	Accounts string  `db:"accounts" json:"accounts" description:"操作账户"`
	Money    float64 `db:"money" json:"money" description:"金额"`
	Note     string  `db:"note" json:"note" description:"备注说明"`
	Btime    string  `db:"btime" json:"btime" description:"帐单时间"`
}

type HumaneGroup struct {
	CID      int64         `db:"cid" json:"cid" description:"收支对象ID"`
	Fullname string        `db:"fullname" json:"fullname" description:"姓名"`
	Pinyin   string        `db:"pinyin" json:"pinyin" description:"拼音"`
	Picture  string        `db:"picture" json:"picture" description:"相片照片"`
	Inc      float64       `json:"inc" description:"对方回报"`
	Out      float64       `json:"out" description:"我的付出"`
	IncRatio float64       `json:"inc_ratio" description:"收入占比"`
	OutRatio float64       `json:"out_ratio" description:"支出占比"`
	Coef     float64       `json:"coef" description:"人情系数"`
	Btime    int64         `json:"btime" description:"最后往来时间"`
	List     []HumaneItems `json:"list" description:"帐单列表"`
}

type HumaneRatio struct {
	Sort  string  `db:"sort" json:"sort" description:"子类别"`
	Money float64 `db:"money" json:"money" description:"金额"`
	Class string  `db:"class" json:"class" description:"主分类"`
}

// 人情往来列表
func (m *HumaneModel) Items(uid int64, filters map[string]interface{}, pageOrder, pageField string) ([]HumaneGroup, error) {
	/*
		SELECT a.`cid`,c.`fullname`,c.`picture`,c.`pinyin`,a.`item`, a.`class`, a.`sort`, a.`object`, a.`accounts`,a.`note`, a.`btime`, a.`money`
		FROM `tbj_account` a,`tbj_contact` c
		WHERE a.`uid` = 6 AND a.`cid`!=0 AND a.cid = c.cid
		ORDER BY a.`intime` DESC


		检验结果
		SELECT a.`item`,a.`cid`,SUM(a.`money`)
		FROM `tbj_account` a,`tbj_contact` c WHERE a.`uid` = 6 AND a.`cid`=5534 AND a.cid = c.cid
		GROUP BY a.item
		ORDER BY a.`intime` DESC
	*/

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
	}

	// 按结构体映射提交字段
	filters = utils.StructAssigMap(HumaneInfo{}, filters)
	fmt.Printf("Type: %T , Data: %v\n", filters, filters) // 打印 args 映射的内容

	// 获取过滤条件数据（表达式::值）
	where_arr, args := utils.GetWhereArgs(filters)
	println(strings.Join(where_arr, " "))
	println("==============================")
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 拼接 NamedQuery 的 select 查询语句 (注意 0000-00-00 00:00:00 在sqlx库中的NamedQuery函数不支持冒号（:）作为命名参数的一部分，因为它会将其解释为命名参数的分隔符。解决方法是使用问号（?）作为命名参数的占位符，然后将命名参数作为第二个参数传递给函数。)
	query := fmt.Sprintf("FROM `%s` a,`%s` c ", m.TableName, m.ContactTableName)
	query += fmt.Sprintf("WHERE a.`uid` = %d AND a.`cid`!=0 AND a.cid = c.cid %s ", uid, strings.Join(where_arr, " "))
	query += "ORDER BY a.`intime` DESC "

	// 拼接 GET 的 select 查询语句
	sql := "SELECT a.`cid`,c.`fullname`,c.`picture`,c.`pinyin`,a.`item`, a.`class`, a.`sort`, a.`object`, a.`accounts`,a.`note`, a.`btime`, a.`money` " + query
	println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// 执行数据库的查询操作 也可以进行结构体 -> 数据库映射，所以结构字段是小写的，并且`db`标签被考虑在内。
	rows, err := m.DB.NamedQuery(sql, args)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return nil, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池

	// 使用结构体拼装数据项组装数据遍历进行数据组装 - 根据KEY值放入map对象进行分组
	// 映射是一个指向值的指针，方便之后使用时检查键值是否存在。
	list := make(map[int64]*HumaneGroup) // 这里如果不使用指针会遇到过错误“无法分配给地图中的结构字段”

	// 遍历查询结果每一行
	for rows.Next() {
		var info HumaneInfo
		if err := rows.StructScan(&info); err != nil {
			return nil, err
		}

		// 数据处理 -
		lastTime := utils.RFC3339ToTime(info.Btime)
		lastUnix := lastTime.Unix() // 获取帐单的秒时间戳
		// 转换时间格式 脱敏
		info.Btime = utils.RFC3339ToString(info.Btime, 2)

		// infos = append(infos, info)

		// 拼装数据组装数据项分组key定义
		key := info.CID
		// fmt.Printf("Type: %T , key: %v\n", key, key)

		// 检查键值是否存在
		if _, ok := list[key]; !ok {
			list[key] = &HumaneGroup{
				CID:      info.CID,
				Fullname: info.Fullname,
				Pinyin:   info.Pinyin,
				Picture:  info.Picture,
				Inc:      0.0,
				Out:      0.0,
				Btime:    lastUnix,
				List:     make([]HumaneItems, 0),
			}
		}
		tmp := HumaneItems{
			Item:     info.Item,
			Class:    info.Class,
			Sort:     info.Sort,
			Object:   info.Object,
			Accounts: info.Accounts,
			Money:    info.Money,
			Note:     info.Note,
			Btime:    info.Btime,
		}

		if list[key].Btime < lastUnix {
			list[key].Btime = lastUnix // 更新为最近的时间
		}

		// 我们可以访问并重新分配 Group 结构对象的值
		// 进行必要的调整后，可以将此 thisProduct 分配回同一个指针。
		list[key].List = append(list[key].List, tmp) // 往MAP里添加数据

		// 将float64浮点数保留两位小数并四舍五入：
		if info.Item == "收入" {
			list[key].Inc = utils.Round2(list[key].Inc + info.Money)
		} else if info.Item == "支出" {
			list[key].Out = utils.Round2(list[key].Out + info.Money)
		}

		// // 占比 - 计算收入和支出的百分比
		// list[key].IncRatio = utils.Round2(list[key].Inc / (list[key].Inc + list[key].Out) * 100) // 收入占比
		// list[key].OutRatio = utils.Round2(100 - list[key].IncRatio)                              // 支出占比

		// // 人情系数 =  收入或支出占比(取最小值的 inc_ratio 或 out_ratio ) x 往来总次数(list.length)
		// if list[key].OutRatio < list[key].IncRatio {
		// 	list[key].Coef = utils.Round2(list[key].OutRatio * float64(len(list[key].List)))
		// } else {
		// 	list[key].Coef = utils.Round2(list[key].IncRatio * float64(len(list[key].List)))
		// }
	}

	// 重新把数据变为 [Group{},Group{}]
	groups := make([]HumaneGroup, 0)
	for _, v := range list {
		// 只最后做一次 占比 和 人情系数 的计算，减少性能消耗

		// 占比 - 计算收入和支出的百分比
		v.IncRatio = utils.Round2(v.Inc / (v.Inc + v.Out) * 100) // 收入占比
		v.OutRatio = utils.Round2(100 - v.IncRatio)              // 支出占比

		// 人情系数 =  收入或支出占比(取最小值的 inc_ratio 或 out_ratio ) x 往来总次数(list.length)
		if v.OutRatio < v.IncRatio {
			v.Coef = utils.Round2(v.OutRatio * float64(len(v.List)))
		} else {
			v.Coef = utils.Round2(v.IncRatio * float64(len(v.List)))
		}

		groups = append(groups, *v)
	}

	// 对切片进行排序 - 根据指定字段从小到大正序或反序列排列groups切片
	sort.Slice(groups, func(i, j int) bool {
		// 排序字段(默认 btime 往来时间、out 我的付出、inc 对方回报、list 往来次数 )
		if pageField == "inc" { // 对方回报
			if pageOrder == "DESC" {
				return groups[i].Inc > groups[j].Inc
			} else {
				return groups[i].Inc < groups[j].Inc
			}
		} else if pageField == "out" { // 我的付出
			if pageOrder == "DESC" {
				return groups[i].Out > groups[j].Out
			} else {
				return groups[i].Out < groups[j].Out
			}
		} else if pageField == "list" { // 往来次数
			if pageOrder == "DESC" {
				return len(groups[i].List) > len(groups[j].List)
			} else {
				return len(groups[i].List) < len(groups[j].List)
			}
		} else if pageField == "coef" { // 人情系数
			if pageOrder == "DESC" {
				return groups[i].Coef > groups[j].Coef
			} else {
				return groups[i].Coef < groups[j].Coef
			}
		} else { // 往来时间
			if pageOrder == "DESC" {
				return groups[i].Btime > groups[j].Btime
			} else {
				return groups[i].Btime < groups[j].Btime
			}
		}
	})

	return groups, nil
}

// 人情往来列表
func (m *HumaneModel) Read(uid, cid int64) (HumaneGroup, error) {
	/*
		SELECT a.`cid`,c.`fullname`,c.`picture`,c.`pinyin`,a.`item`, a.`class`, a.`sort`, a.`object`, a.`accounts`,a.`note`, a.`btime`, a.`money`
		FROM `tbj_account` a,`tbj_contact` c
		WHERE a.`uid` = 6 AND a.`cid`= 123 AND a.cid = c.cid
		ORDER BY a.`intime` DESC
	*/
	var infoGroup HumaneGroup

	// 拼接 NamedQuery 的 select 查询语句 (注意 0000-00-00 00:00:00 在sqlx库中的NamedQuery函数不支持冒号（:）作为命名参数的一部分，因为它会将其解释为命名参数的分隔符。解决方法是使用问号（?）作为命名参数的占位符，然后将命名参数作为第二个参数传递给函数。)
	query := fmt.Sprintf("FROM `%s` a,`%s` c ", m.TableName, m.ContactTableName)
	query += "WHERE a.`uid` = ? AND a.`cid`= ? AND a.cid = c.cid "
	query += "ORDER BY a.`intime` DESC "

	// 拼接 GET 的 select 查询语句
	sql := "SELECT a.`cid`,c.`fullname`,c.`picture`,c.`pinyin`,a.`item`, a.`class`, a.`sort`, a.`object`, a.`accounts`,a.`note`, a.`btime`, a.`money` " + query
	// println("\r\n", sql) // 打印sql
	// fmt.Printf("Type: %T , Data: %v\n", args, args) // 打印 args 映射的内容

	// db.Queryx：这个方法与 db.Query 类似，但它返回的是一个 sqlx.Rows 对象，该对象具有比标准库 sql.Rows 更强大的功能。sqlx.Rows 支持结构体映射、更方便的字段访问方法等，能够简化数据提取过程。适用于需要更灵活的结果处理和数据提取的场景。
	rows, err := m.DB.Queryx(sql, uid, cid)
	if err != nil {
		// println("NamedQuery failed: ", err.Error())
		return infoGroup, err
	}
	defer rows.Close() // 程序结束后释放资源给连接池

	infoGroup = HumaneGroup{
		Inc:  0.0,
		Out:  0.0,
		List: make([]HumaneItems, 0),
	}

	// 遍历查询结果每一行
	for rows.Next() {
		var info HumaneInfo
		if err := rows.StructScan(&info); err != nil {
			return infoGroup, err
		}

		// 数据处理
		lastTime := utils.RFC3339ToTime(info.Btime)
		lastUnix := lastTime.Unix() // 获取帐单的秒时间戳

		if infoGroup.CID == 0 {
			infoGroup.CID = info.CID
			infoGroup.Fullname = info.Fullname
			infoGroup.Pinyin = info.Pinyin
			infoGroup.Picture = info.Picture
			infoGroup.Btime = lastUnix
		}

		// 转换时间格式 脱敏
		info.Btime = utils.RFC3339ToString(info.Btime, 2)

		tmp := HumaneItems{
			Item:     info.Item,
			Class:    info.Class,
			Sort:     info.Sort,
			Object:   info.Object,
			Accounts: info.Accounts,
			Money:    info.Money,
			Note:     info.Note,
			Btime:    info.Btime,
		}

		if infoGroup.Btime < lastUnix {
			infoGroup.Btime = lastUnix // 更新为最近的时间
		}

		// 我们可以访问并重新分配 Group 结构对象的值
		// 进行必要的调整后，可以将此 thisProduct 分配回同一个指针。
		infoGroup.List = append(infoGroup.List, tmp) // 往MAP里添加数据

		// 将float64浮点数保留两位小数并四舍五入：
		if info.Item == "收入" {
			infoGroup.Inc = utils.Round2(infoGroup.Inc + info.Money)
		} else if info.Item == "支出" {
			infoGroup.Out = utils.Round2(infoGroup.Out + info.Money)
		}

		// // 占比 - 计算收入和支出的百分比
		// infoGroup.IncRatio = utils.Round2(infoGroup.Inc / (infoGroup.Inc + infoGroup.Out) * 100) // 收入占比
		// infoGroup.OutRatio = utils.Round2(100 - infoGroup.IncRatio)                              // 支出占比

		// // 人情系数 =  收入或支出占比(取最小值的 inc_ratio 或 out_ratio ) x 往来总次数(list.length)
		// if infoGroup.OutRatio < infoGroup.IncRatio {
		// 	infoGroup.Coef = utils.Round2(infoGroup.OutRatio * float64(len(infoGroup.List)))
		// } else {
		// 	infoGroup.Coef = utils.Round2(infoGroup.IncRatio * float64(len(infoGroup.List)))
		// }
	}

	// 只最后做一次 占比 和 人情系数 的计算，减少性能消耗
	// 占比 - 计算收入和支出的百分比
	infoGroup.IncRatio = utils.Round2(infoGroup.Inc / (infoGroup.Inc + infoGroup.Out) * 100) // 收入占比
	infoGroup.OutRatio = utils.Round2(100 - infoGroup.IncRatio)                              // 支出占比

	// 人情系数 =  收入或支出占比(取最小值的 inc_ratio 或 out_ratio ) x 往来总次数(list.length)
	if infoGroup.OutRatio < infoGroup.IncRatio {
		infoGroup.Coef = utils.Round2(infoGroup.OutRatio * float64(len(infoGroup.List)))
	} else {
		infoGroup.Coef = utils.Round2(infoGroup.IncRatio * float64(len(infoGroup.List)))
	}

	return infoGroup, nil
}

// 礼金收支回报
func (m *HumaneModel) ReportRatio(uid int64) (map[string]interface{}, error) {
	/*
		SELECT SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 'inc',
		SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 'out'
		FROM `tbj_account` WHERE `uid` = 6 AND `cid` != 0
	*/

	// 拼接 GET 的 select 查询语句
	sql := "SELECT SUM(CASE WHEN `item` = '收入' THEN `money` ELSE 0 END) AS 'inc',"
	sql += "SUM(CASE WHEN `item` = '支出' THEN `money` ELSE 0 END) AS 'out' "
	sql += fmt.Sprintf("FROM `%s` WHERE `uid` = ? AND `cid` != 0", m.TableName)
	// println("\r\n", sql) // 打印sql

	type tmepRatio struct {
		Inc float64 `json:"inc" description:"对方回报"`
		Out float64 `json:"out" description:"我的付出"`
	}

	var data tmepRatio
	err := m.DB.Get(&data, sql, uid) // 查询单行数据 ， 也可以用 NamedQuery
	if err != nil {
		// println("Err: ", err.Error())
		return nil, err
	}
	result := utils.StructToMap(data, "json") // 结构体转MAP

	return result, nil
}

// 礼金分类回报排名
func (m *HumaneModel) ReportTop(uid int64) ([]HumaneRatio, error) {
	/*
		SELECT  `sort`, SUM(`money`) 'money', `class`
		FROM `tbj_account` WHERE `uid` = 6 AND `cid` != 0  AND `item` ='收入'
		GROUP BY `sort` ORDER BY `money` DESC
	*/

	// 拼接 GET 的 select 查询语句
	sql := "SELECT  `sort`,  SUM(`money`) 'money', `class` "
	sql += fmt.Sprintf("FROM `%s` WHERE `uid` = ? AND `cid` != 0  AND `item` ='收入' ", m.TableName)
	sql += "GROUP BY `sort` ORDER BY `money` DESC"
	println("\r\n", sql) // 打印sql

	var data []HumaneRatio
	err := m.DB.Select(&data, sql, uid) // 查询多行数据
	if err != nil {
		// println("Err: ", err.Error())
		return nil, err
	}

	return data, nil
}
