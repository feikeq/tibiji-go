// 公共函数定义

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mozillazg/go-pinyin"
)

// 计算给定数据的MD5散列值(32位小写)
// utils.CalculateMD5("MD5 Hash")
func CalculateMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

// 计算提供的输入字符串的SHA1哈希值(40位十六进制数字串)
// utils.CalculateSHA1("SHA-1 Hash")
func CalculateSHA1(input string) string {
	// 将输入字符串转换为字节数组
	inputBytes := []byte(input)
	// 计算SHA1哈希
	hashSum := sha1.Sum(inputBytes)
	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashSum[:])
	return hashString
}

// 使用AES对提供的明文进行加密，加密模式采用的是GCM（Galois/Counter Mode）且不带随机向量
// utils.EncryptAES("TEXT", 32个字符) 密钥长度 !=32 时可以使用 utils.CalculateMD5()
func EncryptAES(plainText string, key string) (string, error) {
	keybyte := []byte(key)

	// 创建一个新的AES加密块
	block, err := aes.NewCipher(keybyte)
	if err != nil {
		return "", fmt.Errorf("无法创建AES加密块：%v", err)
	}

	// 创建一个GCM模式的加密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("无法创建GCM加密器：%v", err)
	}

	// // 生成随机的初始化向量
	// // 是通过make函数生成一个指定大小的随机nonce，‌其大小由gcm.NonceSize()决定，‌通常用于加密过程中，‌保证每次加密的唯一性。‌
	// nonce := make([]byte, gcm.NonceSize())
	// if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
	// 	return "", fmt.Errorf("无法生成随机向量：%v", err)
	// }

	// 初始向量的作用是在加密过程中为每个密码块提供一个独特的起始值
	// 固定初始化向量 ( 减少 "io" 和 "crypto/rand" 包的引用 )
	nonce, _ := base64.StdEncoding.DecodeString("TIBIJIwwwFK68net")
	// 将一个Base64编码的字符串解码成字节数组作为nonce。‌这种方式通常用于解密过程，‌或者当需要固定nonce进行某些操作时。‌

	// 对明文进行加密
	ciphertext := gcm.Seal(nil, nonce, []byte(plainText), nil)

	// // 将加密结果和随机向量进行合并
	// encrypted := append(nonce, ciphertext...)
	// 不合并
	encrypted := ciphertext

	// 对加密结果进行Base64编码
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

// 使用AES对提供的密文进行解密，解密模式采用的是GCM（Galois/Counter Mode）且自带固定向量
// utils.DecryptAES("XXX=", 32个字符) 密钥长度 !=32 时可以使用 utils.CalculateMD5()
func DecryptAES(encryptedText string, key string) (string, error) {
	keybyte := []byte(key)

	// 对Base64编码的密文进行解码
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("无法解码密文：%v", err)
	}

	// 创建一个新的AES解密块
	block, err := aes.NewCipher(keybyte)
	if err != nil {
		return "", fmt.Errorf("无法创建AES解密块：%v", err)
	}

	// 创建一个GCM模式的解密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("无法创建GCM解密器：%v", err)
	}

	// // 获取随机向量的长度
	// nonceSize := gcm.NonceSize()
	// // 提取随机向量和密文
	// nonce := encrypted[:nonceSize]
	// ciphertext := encrypted[nonceSize:]

	// 初始向量的作用是在加密过程中为每个密码块提供一个独特的起始值  - 不用提取直接赋值
	nonce, _ := base64.StdEncoding.DecodeString("TIBIJIwwwFK68net")
	// 将一个Base64编码的字符串解码成字节数组作为nonce。‌这种方式通常用于解密过程，‌或者当需要固定nonce进行某些操作时。‌
	// 在实际应用中，‌应根据加密或解密的需求选择合适的nonce生成方式。‌

	ciphertext := encrypted

	// 对密文进行解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("无法解密密文：%v", err)
	}

	return string(plaintext), nil
}

// 将Form表单json结构和地址栏params统一归纳到 map[string]interface{}
// utils.AllDataToMap(ctx)
func AllDataToMap(ctx iris.Context) map[string]interface{} {

	// 创建一个用于存储所有数据的 map
	allData := make(map[string]interface{})

	// 获取所有表单数据
	formValues := ctx.FormValues() // ctx.FormValues() 等同于 ctx.Request().Form
	// get可以拿到form-data和raw-json的数据但 x-www-form-urlencoded 提交上来的数据拿不到，只有post能全部拿到
	// 将表单数据添加到 map 中
	for key, vals := range formValues {
		if len(vals) > 0 {
			allData[key] = vals[0]
		}
	}

	// 定义目标类型
	var jsonData map[string]interface{}
	// 读取 JSON 数据
	err := ctx.ReadJSON(&jsonData)
	if err != nil {
		// 处理读取 JSON 数据失败的情况
		//  invalid character '-' in numeric literal
		// println("处理读取 JSON 数据失败：", err.Error())
		return allData
	}
	// 将JSON数据添加到 map 中
	for key, value := range jsonData {
		allData[key] = value
	}
	// 返回所有数据
	return allData
}

// 将结构体按提交数据映射为 map[string]interface{} 做数据库SQl拼接时只接受单层结构
// result := utils.StructAssigMap(person)
// for key, value := range result { fmt.Printf("%s: %v\n", key, value)}
func StructAssigMap(theStruct interface{}, theData map[string]interface{}) map[string]interface{} {
	// 创建一个空的 map[string]interface{} 用于存储结果
	result := make(map[string]interface{})

	// 使用反射获取传入变量的反射值
	value := reflect.ValueOf(theStruct)

	// fmt.Printf("变量类型value: %T, 变量的值value: %v\n", value, value)
	// fmt.Printf("变量类型value.Kind(): %T, 变量的值value.Kind(): %v\n", value.Kind(), value.Kind())

	// 如果传入的是指针类型，则获取其指向的值
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 如果值的类型不是结构体，则返回空结果
	if value.Kind() != reflect.Struct {
		return result
	}

	// 获取结构体的类型信息
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		// 获取字段的信息
		field := typ.Field(i)
		f_name := field.Name // 字段名
		// println(f_name)
		// fmt.Printf("获取字段的信息 %v \n", field)
		// fmt.Printf("获取字段的Tag %v \n", field.Tag)
		// println("Tag里的json", field.Tag.Get("json"))
		// println("Tag里的description", field.Tag.Get("description"))
		if field.Tag != "" {
			db_name := field.Tag.Get("db")   // 数据库字段名
			js_name := field.Tag.Get("json") // json字段名
			if db_name != "" {
				f_name = db_name
			} else if field.Tag.Get("json") != "" {
				f_name = js_name
			}
		}
		// println("---------------" + f_name + "----------------")
		if f_name != "-" {
			// 递归多层结构
			if field.Type.Kind() == reflect.Struct {
				// println("==> Struct")
				// result[f_name] = StructAssigMap(field, theData)

				// println("---------------" + f_name + "----------------")

				nestedResult := StructAssigMap(value.Field(i).Interface(), theData)
				for k, v := range nestedResult {
					// println(k, v)
					result[k] = v
				}

			} else {
				// println("==> Setval")
				// 判断是否存在字段 f_name
				if f_value, ok := theData[f_name]; ok {
					result[f_name] = f_value
				}
			}
		}

		// 获取字段的值
		// fieldValue := value.Field(i).Interface()
		// fieldValue := value.Field(i)
		// fmt.Printf("获取字段的值  %v \n", fieldValue)

		// // 转换字段值为字符串形式
		// var stringValue string
		// switch fieldValue.Kind() {
		// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 	stringValue = strconv.FormatInt(fieldValue.Int(), 10)
		// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 	stringValue = strconv.FormatUint(fieldValue.Uint(), 10)
		// case reflect.String:
		// 	stringValue = fieldValue.String()
		// }

		// 将字段名作为键，字段值作为值，添加到结果 map 中
		// result[field.Name] = []string{stringValue}
		// result[field.Name] = fieldValue
	}
	// 返回转换后的 map[string]interface{}
	return result
}

// 结构体转MAP
// result := utils.StructToMap(person,"json")
func StructToMap(theStruct interface{}, tagName ...string) map[string]interface{} {
	tag := ""
	if len(tagName) > 0 {
		tag = tagName[0]
	}

	// 创建一个空的 map[string]interface{} 用于存储结果
	result := make(map[string]interface{})

	// 使用反射获取传入变量的反射值
	value := reflect.ValueOf(theStruct)

	// fmt.Printf("变量类型value: %T, 变量的值value: %v\n", value, value)
	// fmt.Printf("变量类型value.Kind(): %T, 变量的值value.Kind(): %v\n", value.Kind(), value.Kind())

	// 如果传入的是指针类型，则获取其指向的值
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 如果值的类型不是结构体，则返回空结果
	if value.Kind() != reflect.Struct {
		return result
	}

	// 获取结构体的类型信息
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		// 获取字段的信息
		field := typ.Field(i)
		f_name := ""
		// println(f_name)
		// fmt.Printf("获取字段的信息 %v \n", field)
		// fmt.Printf("获取字段的Tag %v \n", field.Tag)
		// println("Tag里的json", field.Tag.Get("json"))
		// println("Tag里的description", field.Tag.Get("description"))
		if tag != "" {
			if field.Tag != "" {
				tag_name := field.Tag.Get(tag)
				if tag_name != "" {
					f_name = tag_name
				}
			}
		} else {
			f_name = field.Name // 字段名
		}

		// println("---------------"+f_name+"----------------", field.Type.Kind())
		if f_name != "-" {
			// 递归多层结构 f_name 有可能为空
			if field.Type.Kind() == reflect.Struct {
				// println("递归 ==> Struct", f_name, value.Field(i).Interface())

				// result[f_name] = StructToMap(field, tag)

				// //在递归调用 StructToMap 时，使用 value.Field(i).Interface() 来传入结构体的值，这样在处理嵌套结构体时才能获取其内部字段并正确转换成 map。
				// result[f_name] = StructToMap(value.Field(i).Interface(), tag) // 这样会变成两层数据

				nestedResult := StructToMap(value.Field(i).Interface(), tag)
				for k, v := range nestedResult {
					result[k] = v
				}

			} else {
				// println("==> Setval")
				result[f_name] = value.Field(i).Interface()

			}
		}

		// 获取字段的值
		// fieldValue := value.Field(i).Interface()
		// fieldValue := value.Field(i)
		// fmt.Printf("获取字段的值  %v \n", fieldValue)

		// // 转换字段值为字符串形式
		// var stringValue string
		// switch fieldValue.Kind() {
		// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 	stringValue = strconv.FormatInt(fieldValue.Int(), 10)
		// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 	stringValue = strconv.FormatUint(fieldValue.Uint(), 10)
		// case reflect.String:
		// 	stringValue = fieldValue.String()
		// }

		// 将字段名作为键，字段值作为值，添加到结果 map 中
		// result[field.Name] = []string{stringValue}
		// result[field.Name] = fieldValue
	}
	// 返回转换后的 map[string]interface{}
	return result
}

// 将数任意据转为INT64 字符转数字
// num := utils.ParseInt64(allData["pageNumber"]) // 任意数据数字（字符串转数字，字符转数字）
func ParseInt64(theNum interface{}) int64 {
	var intv int64

	// strconv.Atoi()和strconv.ParseInt()是两个不同的Golang整数转换函数，它们给出不同的结果是因为它们有不同的参数。
	// strconv.Atoi()只接受字符串参数，而strconv.ParseInt()可以接受字符串和进制参数。
	// 因此，当使用strconv.Atoi()时，它只会将字符串转换为十进制整数，
	// 而strconv.ParseInt()可以将字符串转换为任意进制的整数。
	// strconv.Atoi 的全称，即：ASCII to integer。性能上 Atoi 会略微快一些。

	// 使用类型断言来检查这些值的类型
	// if str, ok := theNum.(string); ok {
	// 	// 处理字符串为整型
	// 	intv, _ = strconv.ParseInt(str, 10, 64) // 字符转数字
	// } else {
	// 	intv = theNum.(int64)
	// }

	// 使用一个switch语句来检查每个值的类型进行类型断言
	// 断言类型 - 使用一个switch语句来检查每个值的类型进行类型断言
	switch inviter := theNum.(type) {
	case string:
		// 处理字符串
		intv, _ = strconv.ParseInt(inviter, 10, 64) // 字符转数字 - strconv.ParseInt()可以将字符串转换为任意进制的整数
	case float64:
		// 处理整数
		intv = int64(inviter) // 将 float64 转换为 int64 (浮点转整型)
	case int:
		// 处理 int
		intv = int64(inviter)
	case int64:
		// 处理 int64
		intv = inviter
	default:
		// 处理其他类型
		intv = inviter.(int64) // 将 int 转换为 int64
	}
	return intv
}

// 任意类类型转整型int
// num := utils.ParseInt(allData["pageNumber"]) // 任意数据数字（字符串转数字，字符转数字）
func ParseInt(theNum interface{}) int {
	var intv int
	// 断言类型 - 使用一个switch语句来检查每个值的类型进行类型断言
	switch inviter := theNum.(type) {
	case int:
		intv = inviter
	case string:
		intv, _ = strconv.Atoi(inviter) // 字符转数字
	case float64:
		// 处理整数
		intv = int(inviter) // 将 float64 转换为 int64 (浮点转整型)
	default:
		// 处理其他类型
		intv = inviter.(int) // 将 int 转换为 int64
	}
	return intv
}

// 将float64浮点数保留两位小数位数并四舍五入
// utils.Round2(3.1415926)
func Round2(num float64) float64 {
	// 四舍六入
	// 第三位为5且5之后有有效数字满足五入
	// 第二位为奇数则进位，第二位为偶数则舍
	// str := fmt.Sprintf("%.2f", num)
	// // 3 -> 3.00 , 3.816 -> 3.82 , 3.815 -> 3.81 , 3.8151 -> 3.82

	// 此方法与  fmt.Sprintf("%.2f", num) 一至
	str := strconv.FormatFloat(num, 'f', 2, 64)

	// 转为 float64
	val, _ := strconv.ParseFloat(str, 64)
	return val

}

// Golang内置net/http包中http.Client结构用于实现HTTP客户端，因此无需 curl 就可以直接使用GET或POST方式发起HTTP请求
// resp, err := utils.HTTPDo("GET", url, map[string][]string{},map[string]string{})
func HTTPDo(method string, url string, values url.Values, headers map[string]string) (map[string]interface{}, error) {
	body := strings.NewReader(values.Encode())

	if headers == nil {
		headers = make(map[string]string) // 创建一个空的 map
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		println("1", err.Error())
		return nil, err
	}

	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", cookie)
	//req.Header.Set("Connection", "keep-alive")
	//req.Header.Add("x-requested-with", "XMLHttpRequest") //AJAX
	// req.Header.Set("User-Agent", "UA信息")

	/*
		使用指定头信息传参时
		var headers = map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}
	*/

	// if headers != nil { // 在 Go 语言中，对于 map 类型的变量，即使它是 nil，也可以使用 range 迭代器进行遍历，并且不会引发运行时错误。因此，在上面的代码中，可以直接删除对 headers 变量是否为 nil 的检查。
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// }

	client := &http.Client{}
	// resp	*http.Response	响应对象
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code %v", resp.StatusCode)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 汉语拼音转换
// utils.PinYin("中国人")
func PinYin(str string) string {
	// go get -u github.com/mozillazg/go-pinyin
	var py_arr []string
	args := pinyin.NewArgs() // 默认

	// 默认情况下会忽略没有拼音的字符（可以通过自定义 Fallback 参数的值来自定义如何处理没有拼音的字符）。
	args.Fallback = func(r rune, a pinyin.Args) []string {
		// 中国人yysd  处理成 [[zhong] [guo] [ren] [y] [y] [s] [d]]
		return []string{string(r)}
	}
	arr := pinyin.Pinyin(str, args) // str = 中国人
	// fmt.Println(arr)                // [["zhong"] ["guo"] ["ren"]]

	// args.Style = pinyin.Tone              // 包含声调
	// fmt.Println(pinyin.Pinyin(str, args)) // [[zhōng] [guó] [rén]]

	// args.Style = pinyin.Tone2             // 声调用数字表示
	// fmt.Println(pinyin.Pinyin(str, args)) // [[zho1ng] [guo2] [re2n]]

	for _, val := range arr {
		// println(key, val[0]) // 0 zhong / 1 guo / 2 ren
		py_arr = append(py_arr, val[0])
	}
	// fmt.Println(py_arr) //[zhong guo ren]

	return strings.Join(py_arr, " ")
}

// 判断数据是否在数组中
// utils.isInArray(arr,2)
func InArray(arr []string, num int) bool {
	val := fmt.Sprintf("%d", num)
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

// 打印最终执行的SQL语句 - 打印最终执行的SQL语句和参数
// utils.PrintExtSql(sql, args)
func PrintExtSql(sql string, args map[string]interface{}) {
	println()
	println()
	println()

	for key, val := range args {
		println(key)

		// 断言类型 - 使用一个switch语句来检查每个值的类型进行类型断言
		switch inviter := val.(type) {
		case string:
			sql = strings.ReplaceAll(sql, ":"+key, "'"+inviter+"'")
		default:
			str := fmt.Sprint(inviter)
			sql = strings.ReplaceAll(sql, ":"+key, str)
		}

	}
	println(sql)

	println()
	println()
	println()
}

// SerializeJSON 序列化结构体为JSON字符串
func SerializeJSON(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// 统一封装返回的JSON数据结构 - 返回约定的JSON数据结构
func ResponseJSON(code int, msg string, data interface{}) string {
	response := struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
