// 动态创建表结构，动态增删改查。
// FIXME: 为了简化操作，只支持原生基本类型。不包括自定义类型，例如：type Long int64
// 参考示例参见Example。
// 如果需要更复杂的条件处理，可以参考Example自行使用StructZero、StructValue、SliceVariable进行组合。
package dynamic

import (
	"reflect"

	"gorm.io/gorm"
)

// Kind对应于reflect.Kind，但只支持基本类型
type Kind reflect.Kind

const (
	_ Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128

	String = Kind(reflect.String)
)

var kindType = [...]reflect.Type{
	Bool:       reflect.TypeOf(bool(false)),
	Int:        reflect.TypeOf(int(0)),
	Int8:       reflect.TypeOf(int8(0)),
	Int16:      reflect.TypeOf(int16(0)),
	Int32:      reflect.TypeOf(int32(0)),
	Int64:      reflect.TypeOf(int64(0)),
	Uint:       reflect.TypeOf(uint(0)),
	Uint8:      reflect.TypeOf(uint8(0)),
	Uint16:     reflect.TypeOf(uint16(0)),
	Uint32:     reflect.TypeOf(uint32(0)),
	Uint64:     reflect.TypeOf(uint64(0)),
	Uintptr:    reflect.TypeOf(uintptr(0)),
	Float32:    reflect.TypeOf(float32(0)),
	Float64:    reflect.TypeOf(float64(0)),
	Complex64:  reflect.TypeOf(complex64(0)),
	Complex128: reflect.TypeOf(complex128(0)),
	String:     reflect.TypeOf(string("")),
}

// Field代表struct的一个字段，也代表数据库表的一个列。
type Field struct {
	Name string // ID
	Kind Kind   // Uint
	Tag  string // `gorm:"primary_key;column:id" json:"id"`
}

// Table代表一个动态表：表名可以指定，表结构也可以动态指定。
type Table struct {
	tableName  string
	structType reflect.Type
}

// 通过指定表名和各个列定义创建一个动态表
func NewTable(table string, fields ...Field) *Table {
	fs := &Table{tableName: table}

	structFields := make([]reflect.StructField, len(fields))
	for i, field := range fields {
		structFields[i].Name = field.Name
		structFields[i].Type = kindType[field.Kind]
		structFields[i].Tag = reflect.StructTag(field.Tag)
	}

	fs.structType = reflect.StructOf(structFields)
	return fs
}

// helper方法，在链式操作之前明确指定要操作的表名
func (tab *Table) DBWithTable(db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{NewDB: true}).
		Table(tab.tableName)
}

// CreateTable将动态表迁移至数据库
func (tab *Table) CreateTable(db *gorm.DB) error {
	schema := tab.StructZero()
	return tab.DBWithTable(db).
		Migrator().
		AutoMigrate(schema)
}

// DropTable将动态表从数据库种删除
func (tab *Table) DropTable(db *gorm.DB) error {
	schema := tab.StructZero()
	return tab.DBWithTable(db).
		Migrator().
		DropTable(schema)
}

// NOTE: 允许values的数量和字段数量不一样，但是对应位置的类型必须一致，使用nil可以跳过某个位置。
// 例如：[]any{uint(77), nil, int(7)}
func (tab *Table) Insert(db *gorm.DB, values []any) (record any, err error) {
	record = tab.StructValue(values)
	err = tab.DBWithTable(db).
		Create(record).Error
	return
}

// First根据条件查询第一条记录，主键增序。
// conditions故意不设计成...interface{}的形式，意在让调用者输入nil来明确表明没有条件。
func (tab *Table) First(db *gorm.DB, conditions []any) (record any, err error) {
	record = tab.StructZero()
	err = tab.DBWithTable(db).
		First(record, conditions...).Error
	return
}

// Last根据条件查询第一条记录，主键降序。
// conditions故意不设计成...interface{}的形式，意在让调用者输入nil来明确表明没有条件。
func (tab *Table) Last(db *gorm.DB, conditions []any) (record any, err error) {
	record = tab.StructZero()
	err = tab.DBWithTable(db).
		Last(record, conditions...).Error
	return
}

// Last根据条件查询第一条记录，无序除非另外指定。
// conditions故意不设计成...interface{}的形式，意在让调用者输入nil来明确表明没有条件。
func (tab *Table) Take(db *gorm.DB, conditions []any) (record any, err error) {
	record = tab.StructZero()
	err = tab.DBWithTable(db).
		Take(record, conditions...).Error
	return
}

// Find根据条件查询所有记录。
// conditions故意不设计成...interface{}的形式，意在让调用者输入nil来明确表明没有条件。
func (tab *Table) Slice(db *gorm.DB, conditions []any) (records any, err error) {
	records = tab.SliceVariable(0, 20)
	err = tab.DBWithTable(db).
		Find(records, conditions...).Error
	return
}

// Updates使用query和args组成的查询条件，将记录值更新为values。
// NOTE: 允许values的数量和字段数量不一样，但是对应位置的类型必须一致，使用nil可以跳过某个位置。
// 例如：[]any{uint(77), nil, int(7)}
func (tab *Table) Updates(db *gorm.DB, values []any, query any, args ...any) (err error) {
	dest := tab.StructValue(values)
	err = tab.DBWithTable(db).
		Where(query, args...).
		Updates(dest).Error
	return
}

// Delete根据条件删除记录。
// conditions故意不设计成...interface{}的形式，意在让调用者输入nil来明确表明没有条件。
func (tab *Table) Delete(db *gorm.DB, conditions []any) (err error) {
	dest := tab.StructZero()
	err = tab.DBWithTable(db).Delete(dest, conditions...).Error
	return
}

// 不包含用户指定的值，相当于var ptr = new(struct{Xxx})
// 主要用于读取：First/Last/Take，也用于创建表结构
func (tab *Table) StructZero() any {
	ptrValue := reflect.New(tab.structType)
	return ptrValue.Interface()
}

// 包含了用户指定的数值，相当于var ptr = new(struct{Xxx}); (*ptr).Xxx = value;
// 主要用于写入：Insert/Update，也可以作为Where条件。
// NOTE: 允许values的数量和字段数量不一样，但是对应位置的类型必须一致，使用nil可以跳过某个位置。
// 例如：[]any{uint(77), nil, int(7)}
func (tab *Table) StructValue(values []any) any {
	ptrValue := reflect.New(tab.structType) // var ptr = new(struct{Xxx})
	elem := ptrValue.Elem()                 // var elem  = *ptr

	// 允许values的数量和字段数量不一样，但是对应位置的类型必须一致，使用nil可以跳过某个位置。
	var N = elem.NumField()
	if len(values) < N {
		N = len(values)
	}
	for i := 0; i < N; i++ {
		field := elem.Field(i)
		if !(field.IsValid() && field.CanSet()) || values[i] == nil {
			continue
		}
		vt := reflect.TypeOf(values[i])
		if vt == field.Type() { // vt.AssignableTo(field.Type())更宽松，但没必要。
			field.Set(reflect.ValueOf(values[i])) // elem.Xxx = value
		}
	}

	return ptrValue.Interface()
}

// 不包含用户指定的值，相当于var ptr = new([]struct{Xxx}); *ptr = make([]struct{Xxx}, 0, 20)
// 主要用于读取：Find
func (tab *Table) SliceVariable(len, cap int) any {
	sliceType := reflect.SliceOf(tab.structType)                     // []struct{Xxx}
	ptrSliceValue := reflect.New(sliceType)                          // var ptr = new([]struct{Xxx})
	ptrSliceValue.Elem().Set(reflect.MakeSlice(sliceType, len, cap)) // *ptr = make([]struct{Xxx}, 0, 20)
	return ptrSliceValue.Interface()
}
