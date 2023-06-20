/*
在工程化的Go语言开发项目中，Go语言的源码复用是建立在包（package）基础之上的。
一个文件夹下面直接包含的文件只能归属一个package,同样一个package的文件不能再多个文件夹下。
所有的包名都应该使用小写字母。包名可以不和文件夹的名字一样（包名不能包含-符号）包名和文件夹名字不一样时不会自动引包
如果想要构建一个程序，则包和包内的文件都必须以正确的顺序进行编译。包的依赖关系决定了其构建顺序。
属于同一个包的源文件必须全部被一起编译，一个包即是编译时的一个单元，因此根据惯例，每个目录都只包含一个包。
*/
// 通过 package 这种方式在 Go 中将要创建的应用程序定义为可执行程序（可以运行的文件）
package config //要创建的包名

// 通过 import 语句可从其他包中的其他代码访问程序
// 一个 Go 程序是通过 import 关键字将一组包链接在一起。
import (
	// import这里不是指文件名而是package-name(代码内的 package 包名)

	// Go的标准库包含了大量的包（如：fmt 和 os）
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	// 本项目自己的包
	"tibiji-go/models"
	"tibiji-go/utils"

	// GORM是GoLang中最出色的ORM框架，支持MySQL、PostgreSQL、Sqlite、SQL Server，功能非常强大，也可以直接执行SQL并获取结果集。
	// "github.com/jinzhu/gorm"

	// 当然GoLang中已有标准database/sql引用也可使用 import "database/sql"
	// go语言自身"database/sql"库实现了sql连接的接口，只要改变数据库驱动就可以连接多种数据库
	// _ "github.com/go-sql-driver/mysql"

	// _ "github.com/lib/pq"

	// Sqlx是对GoLang标准database/sql的部分功能做了扩展，让其使用更为方便!

	// 在GORM和Sqlx之间做选择，需要根据团队的编程风格来确定。
	// 如果想干净的使用SQL，就选择Sqlx。如果偏重于ORM，就选择GORM。
	// GORM存储关系数据时，比Sqlx要少写许多代码。
	// 不管使用哪个，在头脑中的概念一定要清楚：在整洁架构中，领域模型是核心，GORM和Sqlx只是外围存储部分。
	"github.com/jmoiron/sqlx" // sqlx文档 http://jmoiron.github.io/sqlx/
	/*
		在 iris-go 中使用 sqlx 操作数据库，有以下几个方法可用：db.Select、db.Get、db.Query和db.Queryx。这些方法是基于 sqlx 包提供的功能进行封装的。

		1. db.Select：这个方法用于执行 SELECT 查询，并将结果映射到给定的切片对象。它期望查询会返回多行结果。使用该方法时，你需要提供一个切片作为参数，结果将被映射到该切片中的每个元素。适用于需要获取多行结果的场景。

		2. db.Get：这个方法用于执行 SELECT 查询，并将结果映射到给定的单个对象。它期望查询只返回一行结果。使用该方法时，你需要提供一个指针参数，结果将被映射到该指针指向的对象中。适用于只需要获取单行结果的场景。

		3. db.Query：这个方法用于执行任意类型的查询语句，包括 SELECT、INSERT、UPDATE、DELETE 等。它返回一个 Rows 对象，你可以通过该对象的方法来遍历查询结果。适用于需要执行复杂查询或执行非查询语句的场景。

		4. db.Queryx：这个方法与 db.Query 类似，但它返回的是一个 sqlx.Rows 对象，该对象具有比标准库 sql.Rows 更强大的功能。sqlx.Rows 支持结构体映射、更方便的字段访问方法等，能够简化数据提取过程。适用于需要更灵活的结果处理和数据提取的场景。

		5. db.QueryRow：这个方法用于执行 SELECT 查询，并返回查询结果的第一行。它期望查询只返回一行结果。使用该方法时配合Scan()方法使用reflect将sql列返回类型映射为Go类型，你需要提供一个或多个参数来接收查询结果的值。适用于只需要获取查询结果的第一行的场景。

		一般来说，如果你需要获取多行结果并将其映射到切片对象中，可以使用 db.Select。如果你只需要获取一行结果并将其映射到单个对象中，可以使用 db.Get。如果你需要执行复杂的查询语句或非查询语句，可以使用 db.Query。如果你需要更灵活的结果处理或数据提取，可以使用 db.Queryx。根据具体的业务需求和结果处理的方式选择适合的方法。
		db.Get：这个方法也用于执行 SELECT 查询，并返回查询结果的第一行。它与 db.QueryRow 的区别在于，你需要提供一个指针参数来接收查询结果的值，而不是传入多个参数。该指针指向的对象将被映射为查询结果的值。适用于只需要获取查询结果的第一行，并将其映射到单个对象的场景。

		一般来说，如果你只需要获取查询结果的第一行，并且需要将其映射到单个对象，可以使用 db.Get。如果你只需要获取查询结果的第一行，但不需要将其映射到单个对象，或者需要获取多个值，可以使用 db.QueryRow。
		使用 db.Get 的场景可以是在你知道查询结果只会返回一行数据，并且你希望将该行数据映射到一个结构体或对象中。
		使用 db.QueryRow 的场景可以是在你只需要查询结果的第一行数据，并且你不需要将其映射到一个结构体或对象中，而是需要将查询结果的各个列值以不同的数据类型返回，或者进行其他一些特定的处理。

		db.Query 可以用于执行任意类型的查询语句，包括 SELECT、INSERT、UPDATE、DELETE 等，并返回查询结果的多行数据。
		db.Exec 主要用于执行非查询语句，如 INSERT、UPDATE、DELETE，返回执行结果的信息，不包含返回的结果集。
		因此，对于 SELECT 查询语句，你可以使用 db.Query 来执行，并通过遍历 Rows 对象获取多行结果。而对于非查询语句，如 INSERT、UPDATE、DELETE，你可以使用 db.Exec 来执行，并获取执行结果的信息。

		总结一下适用场景：
		当你需要获取多行结果，并将其映射到切片对象中时，使用 db.Select。
		当你只需要获取一行结果，并将其映射到单个对象中时，使用 db.Get。
		当你需要执行复杂的查询语句或非查询语句时，使用 db.Query。
		当你需要更灵活的结果处理和数据提取时，使用 db.Queryx。
		当你只需要获取查询结果的第一行时，使用 db.QueryRow。


		1. row.Scan 方法用于将查询结果的列值按顺序映射到传入的变量中。这意味着你需要按照查询结果的列顺序，将变量作为参数传递给 Scan 方法。 row.Scan(&var1, &var2, ...)
		2. row.StructScan 方法用于将查询结果的列值映射到结构体的字段上。它会根据结构体字段的标签（如 db 标签）来确定映射关系，而不依赖于列的顺序。 row.StructScan(&structVar)
		一般来说，如果你的查询结果是一个简单的单行结果，且你只关心少量的列值，可以使用 Row.Scan 来将结果直接扫描到变量中。而如果你的查询结果需要映射到结构体对象，并且需要处理复杂的字段映射和类型转换，建议使用 StructScan 来简化操作。


	*/)

// 用来存储连接实例（ DB实例不是连接，而是代表数据库的抽象 ）
// 如果想让一个包中的标识符（如变量、常量、类型、函数等）能被外部的包使用，那么标识符必须是对外可见的（public）。
// 在Go语言中是通过标识符的首字母大/小写来控制标识符的对外可见（public）/不可见（private）的。在一个包内部只有首字母大写的标识符才是对外可见的。
var Db *sqlx.DB

/*
在实际开发中，我们建议将 SQL 查询语句封装在一个单独的服务中，以达到更好的模块化和可维护性。您可以将 SQL 查询服务与其他服务一起组成一个或多个模块，以构建完整的应用程序。
*/
// 连接数据库(初始化数据库)
func ConnectDB(cfg map[string]interface{}) (*sqlx.DB, error) {
	// Sqlx是对GoLang标准database/sql的扩展。其特点是：
	// 1.把SQL执行的结果集转化成数据结构(Struct、Maps、Slices)。
	// 2.支持问号(?)或命名的Prepared Statements，避免SQL注入的安全问题
	// sqlx 简化了查询的方式，让你可以更容易的使用查询，sqlx 插入，更新，删除数据的方式与原生的一样。除此外，sqlx还有NamedQuery，NamedExec方法可以让sql语句带有名称的占位符

	/*
		采用传入配置的话是这样
		func ConnectDB(config *iris.Configuration) *sqlx.DB {
			dbtype := config.Other["DB_TYPE"].(string) //config *iris.Configuration
			dbUser := config.Other["DB_USER"].(string)
			dbPassword := config.Other["DB_PASSWORD"].(string)
			dbHost := config.Other["DB_HOST"].(string)
			dbCharset := config.Other["DB_CHARSET"].(string)
			dbName := config.Other["DB_NAME"].(string)
		}
	*/

	// 获取配置传参
	// fmt.Printf("type: %T \n", cfg) // type: map[string]interface {}
	// println("cfg:", cfg) //  cfg: map[DB_CHARSET:utf8mb4 DB_HOST:localhost:3306 ...]

	// 获取数据库配置项

	/*
		类型断言 由于接口是一般类型不知道具体类型，如果要转成具体类型就需要使用类型断言。
		不带检测的类型断言 t := i.(T)   表示看看能不能将 i 转换成 T 再赋给 t
		带检测的类型断言 t, ok:= i.(T)  表示 将 i 转换成 T 转换成功 ok 为true再把转化内容赋给 t
		如果需要区分多种类型，可以使用 type switch 断言 判断类型
		func findType(i interface{}) {
			switch x := i.(type) {
			case int:
				println(x, "is int")
			case string:
				println(x, "is string")
			case nil:
				println(x, "is nil")
			default:
				println(x, "not type matched")
			}
		}
	*/
	dbType, ok := cfg["DB_TYPE"].(string)
	if !ok {
		return nil, fmt.Errorf("数据库类型配置错误")
	}
	// 下划线“_”是特殊标识符意思是忽略这个变量，因为此时不需要知道返回的错误值忽略了error变量
	dbUser, _ := cfg["DB_USER"].(string)
	dbPassword, _ := cfg["DB_PASSWORD"].(string)
	dbHost := cfg["DB_HOST"].(string)
	dbCharset := cfg["DB_CHARSET"].(string)
	dbName := cfg["DB_NAME"].(string)
	dbCheckTabeName := models.UserTableName

	// 连接数据库配置 数据库DSN(Data Source Name)数据源链接
	conn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ")/" + dbName + "?charset=" + dbCharset + "&parseTime=true&loc=Local"
	// println(dbType, conn)

	// sqlx连接MySQL ，建立连接的方式与 golang自带的大差不差
	//database, err := sqlx.Open("数据库类型", "用户名:密码@tcp(地址:端口)/数据库名")
	//database, err := sqlx.Open(dbType, conn)

	// 强制连接并测试它是否有效
	//err=db.Ping()

	/*
		在某些情况下，您可能希望同时打开数据库和连接；
		例如，为了在初始化阶段捕获配置问题。
		您可以使用 一次性完成此操作Connect，它会打开一个新的数据库并尝试Ping.
		该MustConnect变体在遇到错误时会 panic，适合在您的包的模块级别使用。
	*/

	// 打开同时并测试连接 使用 Connect 连接，会验证是否连接成功，(sqlx库的相当于是把sql库的Open方法、Ping方法结合到了一个Connect方法中):
	database, err := sqlx.Connect(dbType, conn)
	// 连接失败
	if err != nil {
		// fmt.Printf("打开数据库失败，连接服务器异常：%v\n", err)
		return nil, fmt.Errorf("打开数据库失败，连接服务器异常：%v", err)
	}

	// defer用于资源的释放，会在函数返回之前进行调用。
	// 如果您忽略了其中一件事情，它们使用的连接可能会一直保留到垃圾回收为止，并且您的数据库最终会一次创建更多的连接以补偿其使用的连接。Rows.Close()可以安全地多次调用它，所以不要害怕在可能不需要的地方调用它。
	// 关闭sqlx.DB数据连接对象(资源释放)
	// defer database.Close() // 注意这行代码要写在上面连接失败err判断的下面

	// if err := database.Ping(); err != nil {
	// 	println(err)
	// }

	// // 打开同时并测试连接，出错时Go中可以抛出一个panic的异常 （使用 MustConnect 连接的话，验证失败不成功直接panic）
	// database := sqlx.MustConnect(dbtype, conn)

	// 默认情况下，同时打开的连接数 (使用中 + 空闲) 没有限制。
	// 当值为0或者小于0时表示无限制（0为默认配置）
	database.SetMaxOpenConns(20) // 设置最大连接数

	// 默认情况下 sql.DB 会在连接池中最多保留 2 个空闲连接
	// 小于或等于0时，不保留任何空闲连接
	database.SetMaxIdleConns(10) // 设置最大空闲数

	// 设置了可重用连接的最大时间长度
	// 设置为0，表示没有最大生存期，并且连接会被重用
	// database.SetConnMaxLifetime(time.Hour)
	// database.SetConnMaxLifetime(time.Minute * 3)

	Db = database

	/*
		但每次操作数据库是，都建立新的数据库链接。
		不仅浪费数据库连接池资源、还容易在并发高情况下造成无法链接数据库。
		于是我们需要想办法在链接数据库前 检查是否复用数据库链接 或 重新建立数据库链接？

		这里我们构建函数它返回一个数据库操作的结构体指针，专门用来执行数据库操作，
		需要注意的是，删除函数内之前的延后defer关闭链接函数，否则链接在函数体内就关闭了，调用方就无法使用数据库了。

	*/

	// 检查表是否存在
	if !isTableExists(dbName, dbCheckTabeName) {
		// 如果表不存在，从文件获取SQL语句并创建表
		if err := createTableFromFile(); err != nil {
			println("初始化数据库失败:", err)
		}
		initDB() // 初始化管理员数据
	}

	return database, nil
}

// 检查数据表是否存在
func isTableExists(dbName string, tabName string) bool {
	var exists bool
	// QueryRow用来从服务器获取一行数据。
	// 它首先从连接池获取一个数据库连接，然后使用Query执行sql查询，并返回一个Row对象；其中Row对象内部又有自己的Rows对象。
	err := Db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ?)", dbName, tabName).Scan(&exists)
	if err != nil {
		// 处理查询错误
		return false
	}
	return exists
}

// 从文件中获取SQL语句并执行创建表
func createTableFromFile() error {

	file, err := os.Open("./config/db.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	createTableQuerys := string(data)
	// println(createTableQuerys)
	// println("_________________________")

	// 拆分 SQL 语句为单个 CREATE TABLE 语句
	queries := strings.Split(createTableQuerys, ";") // 分割字符串

	// 执行每个 CREATE TABLE 语句
	for _, query := range queries {
		// println("_________________________")
		// println(query)
		// 遍历语句时最后一个语句可能为空字符串，因此你可能需要在遍历之前使用len函数检查statements数组的长度，然后在执行时跳过空字符串。
		if query == "" {
			continue
		}

		// 处理执行SQL语句错误
		_, err := Db.Exec(query + ";")
		if err != nil {
			return err
		}
	}
	return nil
}

// 初始化管理员数据
func initDB() {

	userModel := models.NewUserModel(Db)

	data := map[string]interface{}{
		"ciphers":  "admin",
		"username": "admin",
		"nickname": "管理员",
		"sex":      0,
	}

	// 调取创建用户模型 - 返回新插入数据的id
	uid, err := userModel.Create(data)
	if err != nil {
		println("Create Error: ", err.Error())
		return
	}

	// 初始化管理员权限
	admin := map[string]interface{}{
		"uid":    uid,
		"vip":    1,
		"manage": 1,
	}
	errM := userModel.CreateMaterial(admin)
	if errM != nil {
		println("Create Error: ", errM.Error())
		return
	}

	println(uid, "管理员数据初始化成功 \r\n用户名：", data["username"].(string), " 密码：", data["ciphers"].(string))

}

func InsertDB() {

	/*

		// 去掉SQL语句中的反引号 或 SET sql_mode='ANSI_QUOTES' SQL服务器模式包括ANSI_QUOTES模式选项，可以用双引号替代反引号将识别符括起来：

		sql := `
		CREATE TABLE IF NOT EXISTS user_center (
			uid bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID(自动)',
			username varchar(32) NOT NULL DEFAULT '' COMMENT '帐号(开放平台末绑定自动添加openid)',
			ciphers char(32) NOT NULL DEFAULT '' COMMENT '密码(为空则是开放台的没有绑定)',
			intime datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '注册时间',
			uptime datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
			PRIMARY KEY (uid),
			UNIQUE KEY username (username),
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户中心表';



		// Exec 和 MustExec 从连接池中获取连接并在服务器上执行提供的查询。
		// 对于不支持即席查询执行的驱动程序，可能会在幕后创建一个准备好的语句来执行。在返回结果之前连接被返回到池中。
		// result, err := Db.Exec(sql)

		// sqlx提供两个几个MustXXX方法。
		// Must方法是为了简化错误处理而出现的，当开发者确定SQL操作不会返回错误的时候就可以使用Must方法，
		// 但是如果真的出现了未知错误的时候，这个方法内部会触发panic，
		// 开发者需要有一个兜底的方案来处理这个panic，比如使用recover。
		// db.MustExec是一个便捷方法，它执行给定的SQL语句，并且在执行过程中发生错误时会立即抛出异常（panic）。这意味着如果执行过程中出现任何错误，程序将会终止并报告错误信息。通常情况下，db.MustExec用于执行不需要返回结果的SQL语句，如插入、更新或删除操作。
		// db.Exec也执行给定的SQL语句，但在执行过程中出现错误时，它会返回一个error对象，而不会终止程序。这允许您在代码中处理执行过程中可能出现的错误情况。db.Exec适用于需要处理执行结果或可能出现错误的查询语句。
		// 因此，db.MustExec和db.Exec的主要区别在于错误处理方式。db.MustExec在执行过程中遇到错误时会直接终止程序，而db.Exec会返回错误对象，让您有机会处理错误。选择使用哪种方法取决于您对错误处理的偏好和具体场景的需求。
		Db.MustExec(sql) // 您可以使用MustExec，它会在出错时死机 函数在链接出错时会panic。

	*/

	// // 获取时间，该时间带有时区等信息，获取的为当前地区所用时区的时间
	// now := time.Now() // 2022-12-01 15:52:47.310419 +0800 CST m=+5.643172070
	// uid := now.UnixMilli()
	// println(uid)
	// println(now.UnixMicro())
	// println("INSERT INTO sys_user (uid, username, nickname, uptime) VALUES (?, ?, ?, ?)", uid, "feikeq1", "肥客泉", now)
	// insertResult := Db.MustExec("INSERT INTO sys_user (uid, username, nickname, uptime) VALUES (?, ?, ?, ?)", uid, "feikeq1", "肥客泉", now)
	// lastInsertId, _ := insertResult.LastInsertId()
	// println("初始化数据库成功，管理员ID：", lastInsertId)

	// sqlx中的扩展 .Queryx 跟sql中Query的行为是一样的，但是返回的结果是sqlx.Rows，同样对scan做了增强：
	/*
		新增字段：如果您使用db.Select函数，您需要定义一个结构体来映射查询结果。如果您的数据库表添加了一个新字段，您需要更新结构体以包含该字段。否则，当您尝试使用db.Select函数时，将会出现错误。
		如果您不想更新结构体，您可以使用db.Queryx函数和map[string]interface{}来代替db.Select函数和结构体。这样，即使数据库表添加了新字段，您也不需要更新代码。以下是一个示例：

		```go
		rows, err := db.Queryx("SELECT * FROM mytable")
		//rows, err := m.DB.Queryx(sql, timeStartStr, timeEndStr) // 也可以使用问号的方式
		if err != nil {
			log.Fatalln(err)
		}
		// 使用 DB.Queryx 后需要在后面添加 defer rows.Close() 来关闭 sqlx.DB 数据连接对象做资源释放的操作。这样可以确保在执行完查询操作后，及时释放资源，避免资源泄漏和浪费。
		defer rows.Close()


		results := make([]map[string]interface{}, 0)
		// 遍历查询结果每一行
		for rows.Next() {
			row := make(map[string]interface{})
			err = rows.MapScan(row)
			if err != nil {
				log.Fatalln(err)
			}
			results = append(results, row)
		}

		```

		使用db.Unsafe函数可以将查询结果映射到map[string]interface{}中，而无需定义结构体。但是，这样做可能会导致一些安全问题，因为它允许SQL注入攻击。因此，我们建议您使用db.Queryx和map[string]interface{}来代替db.Select和结构体。
		如果您仍然想使用db.Unsafe函数，请确保您的代码已经过充分测试，并且您已经采取了必要的安全措施来防止SQL注入攻击。希望这可以帮到你。如果你还有其他问题，请随时问我。

	*/
	row, _ := Db.Queryx("SELECT CURRENT_TIMESTAMP()")
	var myTime time.Time

	for row.Next() {
		println("--------------1-------------")
		row.Scan(&myTime)
	}

	println("---------------2------------")

	fmt.Println(myTime)

	uid := utils.GenerateTimerID()
	_, intime := utils.FormatTimestamp(time.Now().Unix())
	pwd := utils.HashPassword("admin", intime)
	// println(uid)
	// println(intime)
	// println(pwd)

	// insert 新增
	// 数据库操作语句： INSERT INTO `sys_user` (`username`, `nickname`, `sex`, `intime`, `uptime`) VALUES ('test', '测试帐号',1,'2022-12-14','2022-12-14')

	// database, err := Db.Exec("INSERT INTO `sys_user` (`username`, `nickname`, `sex`, `intime`, `uptime`) VALUES (?, ?, ?, ?, ?)", fmt.Sprintf("test_%v", time.Now().UnixMicro()), "测试帐号", 1, time.Now(), time.Now())
	// if err != nil {
	// 	println("数据库Insert操作失败， ", err)
	// 	return
	// }

	// 有时候如果sql语句中的占位符过多，后面我们传参容易传错。因此有个准确度高的方法如下：
	database, err := Db.NamedExec("INSERT INTO `sys_user` (`uid`, `username`, `ciphers`, `nickname`, `sex`, `intime`, `uptime`) VALUES (:uid, :username, :ciphers, :nickname, :sex, :intime, :uptime)",
		map[string]interface{}{
			"uid":      uid,
			"ciphers":  pwd,
			"username": "admin",
			"nickname": "管理员",
			"sex":      0,
			"intime":   intime,
			"uptime":   intime,
		})
	if err != nil {
		println("数据库Insert操作失败， ", err)
		return
	}
	// 获取插入结果
	id, err := database.LastInsertId() // 新插入数据的id
	if err != nil {
		println("exec failed, ", err)
		return
	}

	println("添加数据成功：", id)
	/*
		sqlx 使用 LastInsertId() 获取不到主键ID怎么办?
		如果您在使用 sqlx 库的 LastInsertId() 方法时无法获取到主键ID，可能是由于以下几个原因：

		没有指定主键字段或主键字段未设置自增属性：LastInsertId() 方法依赖于数据库自动生成的主键值。如果您的表没有指定主键字段或主键字段没有设置自增属性，那么 LastInsertId() 方法将无法返回有效的主键ID。

		数据库不支持自动生成主键：某些数据库可能不支持自动生成主键，或者使用非标准的方式。在这种情况下，LastInsertId() 方法可能无法正常工作。您可以查阅您使用的数据库文档，了解其支持的主键生成方式，并采取相应的方法来获取主键ID。

		数据库驱动问题：有时，特定的数据库驱动可能存在问题，导致 LastInsertId() 方法无法正确返回主键ID。在这种情况下，您可以尝试更新或更换数据库驱动，或者查阅相关文档和社区支持来解决问题。

		如果 LastInsertId() 方法无法满足您的需求，您可以考虑使用其他方法来获取主键ID，具体取决于您使用的数据库和框架。以下是一些备选方案：

		使用 db.Get 或 db.QueryRow 方法查询插入记录后的主键值。例如，对于 MySQL 数据库，可以使用 SELECT LAST_INSERT_ID() 查询语句获取插入的最后一个自动生成的主键值。
		在插入记录后，根据特定条件再次查询数据库以获取主键ID。例如，根据插入的数据条件进行查询，或者使用其他唯一标识符来定位插入的记录，并获取主键ID。
		如果您使用的数据库支持返回插入后的完整记录，您可以使用 db.Get 或 db.QueryRow 方法获取完整记录，并从中提取主键ID。
		根据您使用的具体数据库和框架，可能还有其他适用的方法来获取主键ID。建议查阅相关文档、社区支持或特定数据库的文档，以获取更多针对您情况的指导和解决方案。
	*/
}

func SelectDB() {
	// select 查询

	// SQL注入问题：我们任何时候都不应该自己拼接SQL语句！

	type UserInfo struct {
		ID       int       `db:"uid"`
		UserName string    `db:"username"`
		NickName string    `db:"nickname"`
		Intime   time.Time `db:"intime"`
		Uptime   string    `db:"uptime"`
	}

	// var person UserInfo // 这种一般采用.get来装数据，.select 查询单行数据 采用结构的话会提示expected slice but got struct
	var person []UserInfo // 查询多行数据

	/*
		询占位符?在内部称为bindvars（查询占位符）,它非常重要。
		你应该始终使用它们向数据库发送值，因为它们可以防止SQL注入攻击。
		database/sql不尝试对查询文本进行任何验证；它与编码的参数一起按原样发送到服务器。
		除非驱动程序实现一个特殊的接口，否则在执行之前，查询是在服务器上准备的。
	*/
	// err := Db.Get(&person, "SELECT * FROM `sys_user` LIMIT 1") // 查询单行数据示例代码

	/*
		// they work with regular types as well
		var id int
		err = db.Get(&id, "SELECT count(*) FROM place")

			Get并Select用于rows.Scan可扫描类型和rows.StructScan不可扫描类型。
			它们大致类似于QueryRowand Query，其中 Get 可用于获取单个结果并对其进行扫描，
			而 Select 可用于获取结果片段：查询结果集里的数据要与定义的UserInfo结构一至才行

			在库中提供最常用的就是NamedQuery和NamedExec函数，一个是执行对查询参数命名并绑定，另一个则是对 CUD 操作的查询参数名的绑定
	*/
	// 通常用 .Get获取单条  或   NamedQuery 但 NamedQuery 真的不好用还得 for rows.Next() {} 来赋值
	err := Db.Select(&person, "SELECT `uid`,`username`,`nickname`,`intime`,`uptime` FROM `sys_user` LIMIT ?", 2) // 查询多行数据
	// 注意：在使用sqlx的db.Select时，不能将表名传递到“？”中，因为“？”只能用于插入值，而不能用于插入表名或列名。

	if err != nil {
		println("数据库select查询失败，", err.Error())
		return
	}

	println("数据查询成功：", person)

	// 拼接 NamedQuery 的 select 查询语句 (注意 0000-00-00 00:00:00 在sqlx库中的NamedQuery函数不支持冒号（:）作为命名参数的一部分，因为它会将其解释为命名参数的分隔符。解决方法是使用问号（?）作为命名参数的占位符，然后将命名参数作为第二个参数传递给函数。)
	// rows, err := Db.NamedQuery("SELECT `uid`,`username`,`nickname`FROM `sys_user` WHERE `uid`=:uid LIMIT ?", args)
	// if err != nil {
	// 	println("NamedQuery failed: ", err.Error())
	// }
	// // 程序结束后释放资源给连接池
	// defer rows.Close()
	// // 遍历查询结果每一行
	// for rows.Next() {
	// 	var user UserInfo
	// 	if err := rows.StructScan(&user); err != nil {
	// 		return nil, total, err
	// 	}
	// 	users = append(users, user)
	// }

}

// sqlx中的exec方法与原生sql中的exec使用基本一致：
func ExecDB() {
	// update 更新
	// UPDATE `categories` SET `title` = 'stu0003' WHERE `id` = 8
	// 也可以用 NamedExec
	resUpdate, errUpdate := Db.Exec("UPDATE `categories` SET `title` = ? WHERE `id` = ?", "更新来了哦", 8)
	if errUpdate != nil {
		println("数据库update更新失败，", errUpdate)
		return
	}
	rowUpdate, errUpdate := resUpdate.RowsAffected() // 操作影响的行数
	if errUpdate != nil {
		println("rows failed, ", errUpdate)
	}
	println("数据更新成功，更新了记录条数：", rowUpdate)

	// delete 删除
	// DELETE FROM `myapp`.`categories` WHERE `id` = 4
	resDelete, errDelete := Db.Exec("DELETE FROM `categories` WHERE `id` = ?", 15)
	if errDelete != nil {
		println("数据库delete删除失败，", errDelete)
		return
	}

	rowDelete, errDelete := resDelete.RowsAffected() // 操作影响的行数
	if errDelete != nil {
		println("rows failed, ", errDelete)
	}

	println("成功删除行数：", rowDelete)

}
