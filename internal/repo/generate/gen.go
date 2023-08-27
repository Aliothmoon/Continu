package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"strings"
)

const MysqlDsn = "root:admin@(localhost:3306)/ci?charset=utf8mb4&parseTime=True&loc=Local"

type Querier interface {
}

func main() {

	// 连接数据库
	db, err := gorm.Open(mysql.Open(MysqlDsn))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:        "internal/repo/query",
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:  true,  // generate pointer when field is nullable
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: true, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: true, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: false, // generate with gorm column type tag
		WithUnitTest:     false,
	})
	// 设置目标 db
	g.UseDB(db)
	g.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
		log.Println(columnName)
		return Case2Camel(columnName)
	})
	g.ApplyBasic(g.GenerateAllTable()...)

	//  后续有需求 可生成动态SQL
	//g.ApplyInterface(func(querier Querier) {})
	g.Execute()

	////配置MySQL连接参数
	//username := "root"     //账号
	//password := "root"     //密码
	//host := "localhost"    //数据库地址，可以是Ip或者域名
	//port := 3306           //数据库端口
	//Dbname := "myDatabase" //数据库名
	//timeout := "10s"       //连接超时，10秒
	//
	////拼接Dsn 参数
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
	//	username, password,
	//	host, port,
	//	Dbname, timeout)
	////连接Mysql, 获得DB类型实例，用于后面的数据库读写操作。
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}
}

func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	name = strings.ToLower(name[:1]) + name[1:]
	return strings.Replace(name, " ", "", -1)
}
