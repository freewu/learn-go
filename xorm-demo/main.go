package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"math/rand"
	"strconv"
	"time"
	"xorm.io/core"
)

var engine *xorm.Engine

type User struct {
	Id   		int64		`xorm:"pk autoincr comment('自增ID')"` // 自增ID
	Age 		uint		`xorm:"default(0) comment('年龄')"` // 年龄
	Name 		string		`xorm:"varchar(100) notnull unique 'user_name' comment('姓名')"` // 姓名
	CreatedAt 	time.Time 	`xorm:"created comment('创建时间')"` // 创建时间
	UpdatedAt 	time.Time 	`xorm:"updated comment('更新时间')"` // 更新时间
	Version 	uint 		`xorm:"version comment('乐观锁')"` // 要使用乐观锁，需要使用version标记
	DeletedAt 	time.Time	`xorm:"deleted comment('删除时间')"` // 软删除标记
}

func main() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@tcp(192.168.150.129:3306)/test?charset=utf8")
	if err != nil {
		fmt.Printf("New mysql engine error: %v\n",err)
		panic(err)
	}
	engine.ShowSQL(true) // 在控制台打印出生成的SQL语句
	engine.Logger().SetLevel(core.LOG_DEBUG) // 会在控制台打印调试及以上的信息
	// 将日志输出到指定文件
	//f, err := os.Create(""xorm-sql.log")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//engine.SetLogger(xorm.NewSimpleLogger(f))

	// 将日志记录到syslog 需要
	//logWriter, err := syslog.New(syslog.LOG_DEBUG, "xorm-sql")
	//if err != nil {
	//	log.Fatalf("Fail to create xorm system logger: %v\n", err)
	//}
	//logger := xorm.NewSimpleLogger(logWriter)
	//engine.SetLogger(logger)

	engine.SetMaxIdleConns(10) // 设置连接池的空闲数大小
	engine.SetMaxOpenConns(10) // 设置最大打开连接数

	// 名称映射规则
	// xorm内置了三种IMapper实现：core.SnakeMapper ， core.SameMapper和core.GonicMapper
	//SnakeMapper 支持struct为驼峰式命名，表结构为下划线命名之间的转换，这个是默认的Maper；
	//SameMapper 支持结构体名称和对应的表名称以及结构体field名称与对应的表字段名称相同的命名；
	//GonicMapper 和SnakeMapper很类似，但是对于特定词支持更好，比如ID会翻译成id而不是i_d
	engine.SetMapper(core.GonicMapper{})
	//engine.SetTableMapper(core.SameMapper{})
	//engine.SetColumnMapper(core.SnakeMapper{})

	/**
	// 字段类型映射
	go type's kind												value method								xorm type
	implemented Conversion										Conversion.ToDB / Conversion.FromDB	        Text
	int, int8, int16, int32, uint, uint8, uint16, uint32													Int
	int64, uint64																							BigInt
	float32																									Float
	float64																									Double
	complex64, complex128										json.Marshal / json.UnMarshal				Varchar(64)
	[]uint8																									Blob
	array, slice, map except []uint8							json.Marshal / json.UnMarshal				Text
	bool														1 or 0										Bool
	string																									Varchar(255)
	time.Time																								DateTime
	cascade struct												primary key field value						BigInt
	struct														json.Marshal / json.UnMarshal				Text
	 */

	// 添加表前缀
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "prefix_")
	engine.SetTableMapper(tbMapper)

	/**
	使用Table和Tag改变名称映射
		engine.Table() 指定的临时表名优先级最高
		TableName() string 其次
		Mapper 自动映射的表名优先级最后
	 */

	/**
	Column属性定义
	name	当前field对应的字段的名称，可选，如不写，则自动根据field名字和转换规则命名，如与其它关键字冲突，请使用单引号括起来
	pk	是否是Primary Key，如果在一个struct中有多个字段都使用了此标记，则这多个字段构成了复合主键，单主键当前支持int32,int,int64,uint32,uint,uint64,string这7种Go的数据类型，复合主键支持这7种Go的数据类型的组合。
	autoincr	是否是自增
	[not ]null 或 notnull	是否可以为空
	unique或unique(uniquename)	是否是唯一，如不加括号则该字段不允许重复；如加上括号，则括号中为联合唯一索引的名字，此时如果有另外一个或多个字段和本unique的uniquename相同，则这些uniquename相同的字段组成联合唯一索引
	index或index(indexname)	是否是索引，如不加括号则该字段自身为索引，如加上括号，则括号中为联合索引的名字，此时如果有另外一个或多个字段和本index的indexname相同，则这些indexname相同的字段组成联合索引
	extends	应用于一个匿名成员结构体或者非匿名成员结构体之上，表示此结构体的所有成员也映射到数据库中，extends可加载无限级
	-	这个Field将不进行字段映射
	->	这个Field将只写入到数据库而不从数据库读取
	<-	这个Field将只从数据库读取，而不写入到数据库
	created	这个Field将在Insert时自动赋值为当前时间
	updated	这个Field将在Insert或Update时自动赋值为当前时间
	deleted	这个Field将在Delete时设置为当前时间，并且当前记录不删除
	version	这个Field将会在insert时默认为1，每次更新自动加1
	default 0或default(0)	设置默认值，紧跟的内容如果是Varchar等需要加上单引号
	json	表示内容将先转成Json格式，然后存储到数据库中，数据库中的字段类型可以为Text或者二进制
	comment	设置字段的注释（当前仅支持mysql）
	 */

	/**
	xorm类型和各个数据库类型的对应表：
	xorm		mysql		sqlite3		postgres	remark
	BIT			BIT			INTEGER		BIT
	TINYINT		TINYINT		INTEGER		SMALLINT
	SMALLINT	SMALLINT	INTEGER		SMALLINT
	MEDIUMINT	MEDIUMINT	INTEGER		INTEGER
	INT			INT			INTEGER		INTEGER
	INTEGER		INTEGER		INTEGER		INTEGER
	BIGINT		BIGINT		INTEGER		BIGINT
	CHAR		CHAR		TEXT		CHAR
	VARCHAR		VARCHAR		TEXT		VARCHAR
	TINYTEXT	TINYTEXT	TEXT		TEXT
	TEXT		TEXT		TEXT		TEXT
	MEDIUMTEXT	MEDIUMTEXT	TEXT		TEXT
	LONGTEXT	LONGTEXT	TEXT		TEXT
	BINARY		BINARY		BLOB		BYTEA
	VARBINARY	VARBINARY	BLOB		BYTEA
	DATE		DATE		NUMERIC		DATE
	DATETIME	DATETIME	NUMERIC		TIMESTAMP
	TIME		TIME		NUMERIC		TIME
	TIMESTAMP	TIMESTAMP	NUMERIC		TIMESTAMP
	TIMESTAMPZ	TEXT		TEXT		TIMESTAMP with zone	timestamp with zone info
	REAL		REAL		REAL		REAL
	FLOAT		FLOAT		REAL		REAL
	DOUBLE		DOUBLE		REAL		DOUBLE PRECISION
	DECIMAL		DECIMAL		NUMERIC		DECIMAL
	NUMERIC		NUMERIC		NUMERIC		NUMERIC
	TINYBLOB	TINYBLOB	BLOB		BYTEA
	BLOB		BLOB		BLOB		BYTEA
	MEDIUMBLOB	MEDIUMBLOB	BLOB		BYTEA
	LONGBLOB	LONGBLOB	BLOB		BYTEA
	BYTEA		BLOB		BLOB		BYTEA
	BOOL		TINYINT		INTEGER		BOOLEAN
	SERIAL		INT			INTEGER		SERIAL		auto increment
	BIGSERIAL	BIGINT		INTEGER		BIGSERIAL	auto increment
	 */

	// 获取到数据库中所有的表，字段，索引的信息
	tables, err := engine.DBMetas()
	if err != nil {
		fmt.Printf("get database meta error: %v\n",err)
	}
	for e := range tables {
		t := tables[e]
		fmt.Printf("Type:%v Name:%v AutoIncrement:%v Charset:%v Comment:%v Updated:%v\n",t.Type,t.Name,t.AutoIncrement,t.Comment,t.Charset,t.Updated)
	}
	fmt.Printf("%v\n",tables)

	//Dump数据库结构和数据
	//engine.DumpAll(w io.Writer)
	err = engine.DumpAllToFile("dump.sql")
	if err != nil {
		fmt.Printf("dump database error: %v\n",err)
	}

	// Import 执行数据库SQL脚本
	//engine.Import(r io.Reader)
	//engine.ImportFile(fpath string)

	// 同步数据库结构
	err = engine.Sync2(&User{})
	if err != nil {
		fmt.Printf("sync struct to database error: %v\n",err)
	}
	/*
	* 自动检测和创建表，这个检测是根据表的名字
	* 自动检测和新增表中的字段，这个检测是根据字段名，同时对表中多余的字段给出警告信息
	* 自动检测，创建和删除索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称。因此这里需要注意，如果在一个有大量数据的表中引入新的索引，数据库可能需要一定的时间来建立索引。
	* 自动转换varchar字段类型到text字段类型，自动警告其它字段类型在模型和数据库之间不一致的情况。
	* 自动警告字段的默认值，是否为空信息在模型和数据库之间不匹配的情况
	 */

	// 改变xorm的时区 默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	// 插入数据
	var user User
	user.Name = "bluefrog"
	row, err := engine.Insert(&user) // INSERT INTO `prefix_user` (`user_name`,`created_at`,`updated_at`) VALUES (?, ?, ?) []interface {}{"bluefrog", "2021-09-18 16:15:40", "2021-09-18 16:15:40"}
	if err != nil {
		fmt.Printf("insert data error: %v\n",err)
	}
	fmt.Printf("insert effect row: %v\n",row)

	// 查询
	// 使用GET 来判读是否存在
	// SELECT `id`, `user_name`, `created_at`, `updated_at` FROM `prefix_user` AS `u` WHERE (u.user_name = ?) ORDER BY `name` ASC, `id` DESC LIMIT 1
	b, err := engine.Alias("u").Asc("user_name").Desc("id").Where("u.user_name = ?", "bluefrog").Get(&User{})
	if err != nil {
		fmt.Printf("select data error: %v\n",err)
	}
	fmt.Printf("get flag: %v\n",b) // true

	var u User
	// SELECT `id`, `user_name`, `created_at`, `updated_at` FROM `prefix_user` WHERE `id`=? ORDER BY id DESC LIMIT 1
	b, err = engine.ID(1).OrderBy("id DESC").Get(&u)
	if err != nil {
		fmt.Printf("select data error: %v\n",err)
	}
	fmt.Printf("get flag: %v\n",b) // true
	fmt.Printf("get data: %v\n",u) // {1 bluefrog 2021-09-18 16:15:40 +0800 CST 2021-09-18 16:15:40 +0800 CST}

	// 获取列表

	// 更新
	var user1 User
	b, err = engine.Id(1).Get(&user1)
	if err != nil {
		fmt.Printf("get data error: %v\n",err)
	}
	fmt.Printf("get flag: %v\n",b) // true
	user1.Name = "freewu"
	affected, err := engine.Id(1).Update(&user1)
	if err != nil {
		fmt.Printf("update data error: %v\n",err)
	}
	fmt.Printf("update affected: %v\n",affected) // update affected: 1

	// 通过添加Cols函数指定需要更新结构体中的哪些值，未指定的将不更新
	user1.Age = 18
	user1.Name = "xxxx"
	// UPDATE `prefix_user` SET `age` = ?, `updated_at` = ?, `version` = `version` + 1 WHERE `id`=? AND `version`=? []interface {}{18, "2021-09-23 10:03:12", 1, 0x2}
	affected, err = engine.Id(1).Cols("age").Update(&user1)
	if err != nil {
		fmt.Printf("update cols data error: %v\n",err)
	}
	fmt.Printf("update cols affected: %v\n",affected) // update affected: 1

	// 通过传入map[string]interface{}来进行更新，需要指定更新到哪个表
	//affected, err = engine.Table(new(User)).Where("id = ?",1).Update(map[string]interface{}{ "age": 6 })
	affected, err = engine.Table("prefix_user").ID(1).Update(map[string]interface{}{"age":6})
	//affected, err = engine.Table(new(User)).ID(1).Update(map[string]interface{}{"age":6,"version":2})
	// 如果有设置乐观锁 使用 Table(new(User)) + Update(map[string]interface{}) 会出现 call of reflect.Value.Interface on zero Value
	if err != nil {
		fmt.Printf("update map data error: %v\n",err)
	}
	fmt.Printf("update map affected: %v\n",affected) // update affected: 0

	// 删除
	// 软删除 需结构体定义 deleted
	affected, err = engine.Where("user_name = ?", "bluefrog").Delete(&User{})
	if err != nil {
		fmt.Printf("delete data error: %v\n",err)
	}
	fmt.Printf("delete affected: %v\n",affected) // delete affected: 1
	// 物理删除
	affected, err = engine.Where("user_name = ?", "bluefrog").Unscoped().Delete(&User{})
	if err != nil {
		fmt.Printf("unscoped delete data error: %v\n",err)
	}
	fmt.Printf("unscoped delete affected: %v\n",affected) // delete affected: 1

	// 执行原生sql
	sql :=`INSERT INTO prefix_user(user_name,age) VALUES (?, ?)`
	res, err := engine.Exec(sql, "bluefrog" + strconv.Itoa(rand.Int()), 20)
	if err != nil {
		fmt.Printf("insert sql error: %v\n",err)
	} else {
		lastInsertId, _ := res.LastInsertId()
		rowsAffected,_ := res.RowsAffected()
		fmt.Printf("sql: %v result.LastInsertId: %v RowsAffected: %v \n",sql,lastInsertId,rowsAffected)
	}

	sql = `UPDATE prefix_user SET age = ? WHERE user_name = ?`
	res, err = engine.Exec(sql, 91, "freewu")
	if err != nil {
		fmt.Printf("update sql error: %v\n",err)
	} else {
		lastInsertId, _ := res.LastInsertId()
		rowsAffected,_ := res.RowsAffected()
		fmt.Printf("sql: %v result.LastInsertId: %v RowsAffected: %v \n",sql,lastInsertId,rowsAffected)
	}



	// 事务

}
