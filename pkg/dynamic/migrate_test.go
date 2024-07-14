package dynamic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Neightly/juzhongmishi/pkg/dynamic"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 执行的sql随驱动包括版本不同而不同，因此该Example有失败的可能性，但不影响起到示例作用。
func Example() {
	// 准备DB
	dialector := sqlite.Open(`:memory:`)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &sqlLogger{logger.Discard},
	})
	assertError(err)

	// 表结构定义
	var tab = dynamic.NewTable("example_dynamic_0",
		dynamic.Field{Name: "ID", Kind: dynamic.Uint, Tag: `gorm:"primary_key;column:id" json:"id"`},
		dynamic.Field{Name: "Name", Kind: dynamic.String, Tag: `gorm:"column:name" json:"name"`},
		dynamic.Field{Name: "Age", Kind: dynamic.Int, Tag: `gorm:"column:age" json:"age"`},
	)

	// 创建表结构
	err = tab.CreateTable(db)
	assertError(err)
	defer tab.DropTable(db)

	// 插入2条记录
	record1, err1 := tab.Insert(db, []any{uint(77), "seven", int(7)})
	record2, err2 := tab.Insert(db, []any{uint(99), nil, int(9)}) // 省略name字段
	assertError(err1)
	assertError(err2)
	json.NewEncoder(os.Stderr).Encode(record1) // {"id":77,"name":"seven","age":7}
	json.NewEncoder(os.Stderr).Encode(record2) // {"id":99,"name":"","age":9}

	// 查询第2条记录
	destLast, err := tab.Last(db, []any{"id = ?", uint(99)})
	assertError(err)
	json.NewEncoder(os.Stderr).Encode(destLast) // {"id":99,"name":"","age":9}

	// 查询全部2条条记录
	destAll, err := tab.Slice(db, []any{"id > ? AND name IS NOT NULL", uint(0)})
	assertError(err)
	json.NewEncoder(os.Stderr).Encode(destAll) // [{"id":77,"name":"seven","age":7},{"id":99,"name":"","age":9}]

	// 更新第1条记录，如果包含主键则作为Where条件
	err = tab.Updates(db, []any{uint(77), "dummy", int(3)}, "name = ?", "other")
	assertError(err)

	// 删除第2条记录，如果条件是一个数字会按主键对待
	err = tab.Delete(db, []any{uint(99)})
	assertError(err)

	// Output:
	// SELECT count(*) FROM sqlite_master WHERE type='table' AND name="example_dynamic_0"
	// CREATE TABLE `example_dynamic_0` (`id` integer PRIMARY KEY AUTOINCREMENT,`name` text,`age` integer)
	// INSERT INTO `example_dynamic_0` (`name`,`age`,`id`) VALUES ("seven",7,77) RETURNING `id`
	// INSERT INTO `example_dynamic_0` (`name`,`age`,`id`) VALUES ("",9,99) RETURNING `id`
	// SELECT * FROM `example_dynamic_0` WHERE id = 99 ORDER BY `example_dynamic_0`.`id` DESC LIMIT 1
	// SELECT * FROM `example_dynamic_0` WHERE id > 0 AND name IS NOT NULL
	// UPDATE `example_dynamic_0` SET `name`="dummy",`age`=3 WHERE name = "other" AND `id` = 77
	// DELETE FROM `example_dynamic_0` WHERE `example_dynamic_0`.`id` = 99
	// PRAGMA foreign_keys
	// DROP TABLE IF EXISTS `example_dynamic_0`
}

func assertError(err error) {
	if err != nil {
		panic(err)
	}
}

type sqlLogger struct {
	logger.Interface
}

func (l *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	fmt.Fprintln(os.Stdout, sql)
}
