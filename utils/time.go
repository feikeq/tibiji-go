package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// 将时间戳转换为日期格式（年-月-日 时:分:秒）
// utils.FormatTimestamp(time.Now().Unix()) 获取当前秒时间戳
func FormatTimestamp(timestamp int64) (time.Time, string) {
	// 时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）
	// 也可以不用此函数直接在程序里使用 time.Now().Format("2006-01-02 15:04:05")
	theTime := time.Unix(timestamp, 0)
	return theTime, theTime.Format("2006-01-02 15:04:05")
}

// 字符串时间转时间 返回时间和时间戳
// utils.ParseTimeToTimestamp("2023-05-12 15:44:34") 严格按照模板来
func ParseTimeToTimestamp(timeStr string) (time.Time, int64) {
	// println("字符串时间转时间戳=>", timeStr)
	// formatTime, _ := time.Parse("2006-01-02 15:04:05", timeStr)

	loc, _ := time.LoadLocation("Local") // 当前时区
	// time.LoadLocation("Local") 和 time.Local 都返回本地时区。
	// 但是，它们之间的区别在于，前者返回一个新的 *time.Location 对象，
	// 而后者返回一个全局变量，该变量在程序中的任何位置都可以使用。
	// 如果您需要在程序中的多个位置使用本地时区，则应使用 time.LoadLocation()。
	// 如果您只需要在程序中的一个位置使用本地时区，则可以使用 time.Local。

	// ParseInLocation 将时间字符串转换为本地时间类型
	date, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	// 返回时间和时间戳 date.Year(), date.Month(), date.Day() ,date.Hour(), date.Minute(), date.Second()
	return date, date.Unix()
}

// RFC3339字符串时间按格式转日期字符串
// utils.RFC3339ToString("2023-05-15T10:21:46+08:00")
func RFC3339ToString(timeStr string, num ...int) string {
	// RFC3339是一种日期和时间格式，它是ISO 8601的一种日期显示格式。 RFC 3339需要完整的日期和时间表示（只有小数秒是可选的）。
	// RFC 3339格式：“{年}-{月}-{日}T{时}:{分}:{秒}.{毫秒}{时区}”；其中的年要用零补齐为4位，月日时分秒则补齐为2位。毫秒部分是可选的。最后一部分是时区。
	thetime, _ := time.Parse(time.RFC3339, timeStr)
	hideTime := 3
	if len(num) > 0 {
		hideTime = num[0]
	}
	if hideTime == 0 {
		return thetime.Format("2006-01-02")
	} else if hideTime == 2 {
		return thetime.Format("2006-01-02 15:04")
	} else {
		return thetime.Format("2006-01-02 15:04:05")
	}
}

// 标准英文时间字符串 转为 年月日时分秒 中文时间字符串
// utils.DateTimeToNYRSFMTime("2023-05-15T10:21:46+08:00")
func DateTimeToYRNTime(timeStr string, showYear ...bool) string {
	// RFC3339是一种日期和时间格式，它是ISO 8601的一种日期显示格式。 RFC 3339需要完整的日期和时间表示（只有小数秒是可选的）。
	// RFC 3339格式：“{年}-{月}-{日}T{时}:{分}:{秒}.{毫秒}{时区}”；其中的年要用零补齐为4位，月日时分秒则补齐为2位。毫秒部分是可选的。最后一部分是时区。
	thetime, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	hideYear := true
	if len(showYear) > 0 {
		hideYear = !showYear[0]
	}
	if hideYear {
		return thetime.Format("01月02日")
	} else {
		return thetime.Format("01月02日（2006年）")
	}
}

// RFC3339字符串时间按格式转时间
// utils.RFC3339ToTime("2023-05-15T10:21:46+08:00")
func RFC3339ToTime(timeStr string) time.Time {
	// println("RFC3339字符串时间按格式出来字符串")
	thetime, _ := time.Parse(time.RFC3339, timeStr)

	return thetime
}

// 计算两个日期之间相差天数
// utils.DaysBetweenDates("2023-05-12 00:00:01", "2023-05-22 23:59:59")
func DaysBetweenDates(date1 string, date2 string) int {
	layout := "2006-01-02 15:04:05"

	t1, _ := time.Parse(layout, date1)
	t2, _ := time.Parse(layout, date2)

	// 只计算日期部分的天数差异，而忽略时分秒的影响 （其实也不需要因为不需要四舍五入）
	// time.Date最后一个参数 t1.Location() 的返回 t1 时间的时区信息的时区创建的实例
	// 如果您使用 time.Parse() 解析一个时间字符串，它将默认为 UTC 时区。如果您想要使用本地时区，您可以使用 time.ParseInLocation() 并将本地时区作为第三个参数传递。例如：
	// 如果您想要创建一个与本地时区相关的时间值，可以使用 time.Local 作为最后一个参数。
	// 如果您想要创建一个与协调世界时 (UTC) 相关的时间值，可以使用 time.UTC 作为最后一个参数。
	// t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	// t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())

	duration := t2.Sub(t1)
	// fmt.Println(duration)
	days := int(duration.Hours() / 24)
	return days
}

// 返回生日时、天数、周年、生肖、星座
// utils.BirthdayAndDays(info.Birthday, info.Lunar)
func BirthdayAndDaysStarZodiac(bTime time.Time, isLunar int) (int, string, string, string, int, string, string, string, string) {
	// println(bTime.Format("2006-01-02 15:04"))

	// 在Go中，可以使用recover函数捕获panic异常，然后正常处理，从而恢复正常代码的执行.
	// 通常情况下，Go中可以抛出一个panic的异常，然后在defer中通过recover捕获这个异常
	// 使用defer语句用于捕获panic异常:
	defer func() {
		if err := recover(); err != nil {
			println("BirthdayAndDaysStarZodiac ERR:")
			fmt.Println(err)
		}
	}()

	// 当前年份
	theYear := time.Now().Year()
	// 当前时间
	nowStr := fmt.Sprintf("%04d-%02d-%02d 00:00:00", theYear, int(time.Now().Month()), time.Now().Day())
	var newStr string
	var num int       // 剩余天数
	var age int       // 年龄或周年
	var star string   // 星座
	var zodiac string // 生肖
	var cndate string // 中文日期
	nextYear := false // 下一年

	// 先判断是否为农历
	if isLunar == 1 {
		// println("农历（阴历）")
		// 阳历>阴历 将当前日期阳历转换为阴历 （主要方便取得阴历年份，以免出现当日期为13年1月时，农历为12年12[腊月]7号时算到农历13年）
		// 取当前时间转为农历(阴历)
		lunarDate, _ := DateToLunar(time.Now()) // 公历转农历(阳历转阴历)

		// 前面添0补0 前置补零(前导零)，保证年份四位置数 月和日始终是2位数的长度
		// timeStr := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", lunarDate.Year(), int(bTime.Month()), bTime.Day(), bTime.Hour(), bTime.Minute(), bTime.Second())
		timeStr := fmt.Sprintf("%04d-%02d-%02d 00:00:00", lunarDate.Year(), int(bTime.Month()), bTime.Day()) // 把 提醒时间的 月-日 取出来
		// 用当前lunarDate农历的年拼接提醒的农历月日生成新的时间格式
		newTime, _ := ParseTimeToTimestamp(timeStr)

		//进行 农历(阴历) 转 公历(阳历) 来求的农历几月几号在今年的哪一天
		_, newStr = LunarToDate(newTime)

		// 计算两个日期之间相差天数
		num = DaysBetweenDates(nowStr, newStr)

		// 计算年龄
		age = lunarDate.Year() - bTime.Year()

		// 如果已经过了生日则转至下一年
		if num < 0 {
			timeStr = fmt.Sprintf("%04d-%02d-%02d 00:00:00", lunarDate.Year()+1, int(bTime.Month()), bTime.Day())
			newTime, _ = ParseTimeToTimestamp(timeStr)
			_, newStr = LunarToDate(newTime)
			num = DaysBetweenDates(nowStr, newStr)
			age += 1
			nextYear = true
		}

	} else {
		// println("公历（阳历）")
		// 把 提醒时间的 年-月-日 取出来
		newStr = fmt.Sprintf("%04d-%02d-%02d 00:00:00", theYear, int(bTime.Month()), bTime.Day())
		// 计算两个日期之间相差天数
		num = DaysBetweenDates(nowStr, newStr)
		// 计算年龄
		age = theYear - bTime.Year()

		// 如果已经过了生日则转至下一年
		if num < 0 {
			// 将时间加一年
			newStr = fmt.Sprintf("%04d-%02d-%02d 00:00:00", theYear+1, int(bTime.Month()), bTime.Day())
			num = DaysBetweenDates(nowStr, newStr)
			age += 1
			nextYear = true
		}
	}

	// 星座
	star = GetXingZuo(bTime, isLunar)
	// 生肖
	zodiac = GetYearShengXiao(bTime, isLunar)
	// 中文日期
	cndate = GetChineseYYYMMDD(bTime, isLunar)

	// 今天明天后天N天 星期几
	dayStr, dateStr, weekStr := FriendlyDate(newStr, false) // 友好时间

	if dayStr == "" {
		dayStr = fmt.Sprintf("%d天后", num)
	}

	// // 返回下次生日的具体日期
	// date := strings.Split(newStr, " ")[0] // 分割字符串 - 拆分(表达式:值)

	// 获取 X月X日（2023年）
	yrn := DateTimeToYRNTime(newStr, nextYear)

	//    int, string,string,string, int, string, string, string
	return num, dayStr, dateStr, weekStr, age, star, zodiac, yrn, cndate
}

// 获取星座
func GetXingZuo(date time.Time, isLunar int) string {
	if isLunar == 1 {
		date, _ = LunarToDate(date)
	}
	solar := calendar.NewSolarFromDate(date)
	return solar.GetXingZuo()
}

// 获取生肖
func GetYearShengXiao(date time.Time, isLunar int) string {
	var lunar *calendar.Lunar
	if isLunar == 1 {
		lunar = calendar.NewLunarFromYmd(date.Year(), int(date.Month()), date.Day())
	} else {
		lunar = calendar.NewLunarFromDate(date) // 公历转农历 （阳历转阴历）
	}
	// println(lunar.String())
	return lunar.GetYearShengXiao()
}

// 获取阴历的中文日期
func GetChineseYYYMMDD(date time.Time, isLunar int) string {
	var lunar *calendar.Lunar
	if isLunar == 1 {
		lunar = calendar.NewLunarFromYmd(date.Year(), int(date.Month()), date.Day())
	} else {
		lunar = calendar.NewLunarFromDate(date) // 公历转农历 （阳历转阴历）
	}
	return lunar.String() // 一九八六年四月廿一
}

// 根据时间返回“前天”、“昨天”、“今天”、“明天”、“后天”这种更友好的字符串
// dayStr,dateStr,weekStr := utils.FriendlyDate(time) // 明天, 2020-12-12, 星期几
func FriendlyDate(timeStr string, isRFC3339 ...bool) (string, string, string) {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	afterTomorrow := today.AddDate(0, 0, 2)
	beforeYesterday := today.AddDate(0, 0, -2)
	layout := "2006-01-02"
	weekdays := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

	isRFC := true
	var thetime time.Time
	if len(isRFC3339) > 0 {
		isRFC = isRFC3339[0]
	}
	if isRFC {
		thetime, _ = time.Parse(time.RFC3339, timeStr)
	} else {
		thetime, _ = time.Parse(layout+" 15:04:05", timeStr)
	}

	weekStr := weekdays[thetime.Weekday()]
	dateStr := thetime.Format(layout)
	dayStr := ""

	if dateStr == today.Format(layout) {
		dateStr = "今天"
	} else if dateStr == yesterday.Format(layout) {
		dateStr = "昨天"
	} else if dateStr == tomorrow.Format(layout) {
		dateStr = "明天"
	} else if dateStr == afterTomorrow.Format(layout) {
		dateStr = "后天"
	} else if dateStr == beforeYesterday.Format(layout) {
		dateStr = "前天"
	}
	return dayStr, dateStr, weekStr
}

// 公历转农历(阳历转阴历)
// date, str := utils.DateToLunar(time.Now())
func DateToLunar(date time.Time) (time.Time, string) {
	// 实例化
	lunar := calendar.NewLunarFromDate(date)

	// println(lunar.String()) // 二〇二三年四月十三
	// println(lunar.ToFullString())
	// 二〇二三年四月十三 癸卯(兔)年 丁巳(蛇)月 己丑(牛)日 辰(龙)时 纳音[金箔金 沙中土 霹雳火 大林木] 星期三 南方朱雀 星宿[轸水蚓](吉) 彭祖百忌[己不破券二比并亡 丑不冠带主不还乡] 喜神方位[艮](东北) 阳贵神方位[坎](正北) 阴贵神方位[坤](西南) 福神方位[坎](正北) 财神方位[坎](正北) 冲[(癸未)羊] 煞[东]

	// 前面添0补0 前置补零(前导零)，保证年份四位置数 月和日始终是2位数的长度
	timeStr := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", lunar.GetYear(), lunar.GetMonth(), lunar.GetDay(), date.Hour(), date.Minute(), date.Second())

	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)

	return formatTime, timeStr
}

// 农历转公历(阴历转阳历)
// date, str := utils.LunarToDate(time.Now())
func LunarToDate(date time.Time) (time.Time, string) {
	// 在Go中，可以使用recover函数捕获panic异常，然后正常处理，从而恢复正常代码的执行.
	// 通常情况下，Go中可以抛出一个panic的异常，然后在defer中通过recover捕获这个异常
	// 使用defer语句用于捕获panic异常:
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		println("LunarToDate Err:")
	// 		fmt.Println(err)
	// 	}
	// }()

	// 实例化
	lunar := calendar.NewLunarFromYmd(date.Year(), int(date.Month()), date.Day())

	// 转阳历
	solar := lunar.GetSolar()
	// println(solar.ToYmdHms()) // 2023-06-29 00:00:00
	// println(solar.ToFullString()) // 2023-06-29 00:00:00 星期四 巨蟹座

	tempArr := strings.Split(solar.ToYmdHms(), " ") // 分割字符串 - 拆分(表达式:值)
	// fmt.Printf("LunarToDate:%v", tempArr)

	// 前面添0补0 前置补零(前导零)，保证年份四位置数 月和日始终是2位数的长度
	timeStr := fmt.Sprintf("%s %02d:%02d:%02d", tempArr[0], date.Hour(), date.Minute(), date.Second())

	formatTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)

	return formatTime, timeStr
}

// 结构体对应字段为 utils.FKTime 按指定格式化时间
type FKTime time.Time

// 提高转化时间字段的效率 (虽然输出方便了但使用StructAssigMap将结构体按提交数据映射为 map[string]interface{} 时会出问题)
func (t FKTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04"))
	return []byte(stamp), nil
}
