// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package orm

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

// Filter marks filtering the fields which does not exist in the fields of the operated table.
// Note that this function supports only single table operations.
func (m *Model) Filter() *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function Filter supports only single table operations")
	}
	model := m.getModel()
	model.filter = true
	return model
}

// Fields 指定需要操作的表字段，包括查询字段、写入字段、更新字段等, 多个字段之间用','连接。
// The parameter <fieldNamesOrMapStruct> can be type of string/map/*map/struct/*struct.
func (m *Model) Fields(fieldNamesOrMapStruct ...interface{}) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	switch {
	// String slice.
	case length >= 2:
		model := m.getModel()
		model.fields = gstr.Join(m.mappingAndFilterToTableFields(gconv.Strings(fieldNamesOrMapStruct)), ",")
		return model
	// It need type asserting.
	case length == 1:
		model := m.getModel()
		switch r := fieldNamesOrMapStruct[0].(type) {
		case string:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields([]string{r}), ",")
		case []string:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields(r), ",")
		default:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields(gutil.Keys(r)), ",")
		}
		return model
	}
	return m
}

// FieldsEx sets the excluded operation fields of the model, multiple fields joined using char ','.
// Note that this function supports only single table operations.
// The parameter <fieldNamesOrMapStruct> can be type of string/map/*map/struct/*struct.
func (m *Model) FieldsEx(fieldNamesOrMapStruct ...interface{}) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	model := m.getModel()
	switch {
	case length >= 2:
		model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(gconv.Strings(fieldNamesOrMapStruct)), ",")
		return model
	case length == 1:
		switch r := fieldNamesOrMapStruct[0].(type) {
		case string:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields([]string{r}), ",")
		case []string:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(r), ",")
		default:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(gutil.Keys(r)), ",")
		}
		return model
	}
	return m
}

// FieldsName (原: FieldsStr)检索并返回表中所有的字段名，各字段之间用','连接。
// 可选参数 <prefix> 为每个字段指定前缀, eg: FieldsStr("u.").
func (m *Model) FieldsName(prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		panic(err)
	}
	if len(tableFields) == 0 {
		panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
	}
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	newFields := ""
	for _, k := range fieldsArray {
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += prefixStr + k
	}
	newFields = m.db.QuoteString(newFields)
	return newFields
}

// FieldsExName (原FieldsExStr) 检索并返回表中例外的字段的名称。
//
// 参数<excludeFields>是被排除，不被检索的字段。各个字段之间用','连接。
// 可选参数 <prefix> 为每个字段指定前缀, eg: FieldsExStr("id", "u.").
func (m *Model) FieldsExName(excludeFields string, prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		panic(err)
	}
	if len(tableFields) == 0 {
		panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
	}
	fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(excludeFields, ","))
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	newFields := ""
	for _, k := range fieldsArray {
		if fieldsExSet.Contains(k) {
			continue
		}
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += prefixStr + k
	}
	newFields = m.db.QuoteString(newFields)
	return newFields
}

// FieldValue (原FindValue) retrieves and returns single field value by Model.WherePri and Model.Value.
// Also see Model.WherePri and Model.Value.
func (m *Model) FieldValue(fieldsAndWhere ...interface{}) (Value, error) {
	if len(fieldsAndWhere) >= 2 {
		return m.WherePri(fieldsAndWhere[1], fieldsAndWhere[2:]...).Fields(gconv.String(fieldsAndWhere[0])).value()
	}
	if len(fieldsAndWhere) == 1 {
		return m.Fields(gconv.String(fieldsAndWhere[0])).value()
	}
	return m.value()
}

// Array (原FindArray) queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
// Also see Model.WherePri and Model.Value.
func (m *Model) FieldList(fieldsAndWhere ...interface{}) ([]Value, error) {
	if len(fieldsAndWhere) >= 2 {
		return m.WherePri(fieldsAndWhere[1], fieldsAndWhere[2:]...).Fields(gconv.String(fieldsAndWhere[0])).array()
	}
	if len(fieldsAndWhere) == 1 {
		return m.Fields(gconv.String(fieldsAndWhere[0])).array()
	}
	return m.array()
}

// HasField determine whether the field exists in the table.
func (m *Model) HasField(field string) (bool, error) {
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		return false, err
	}
	if len(tableFields) == 0 {
		return false, fmt.Errorf(`empty table fields for table "%s"`, m.tables)
	}
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	for _, f := range fieldsArray {
		if f == field {
			return true, nil
		}
	}
	return false, nil
}
