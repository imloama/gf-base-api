package db

import (
	"database/sql"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/encoding/gbinary"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"strings"
	"time"
)

func Query2Structs(pointer interface{}, sql string, args ...interface{})error{
	rows, err := g.DB().Query(sql, args)
	if err!=nil{
		return err
	}
	result, err := ConvertRowsToResult(rows)
	if err!=nil{
		return err
	}
	err = result.Structs(pointer)
	if err!=nil{
		return err
	}
	return nil
}

func ConvertRowsToResult(rows *sql.Rows) (gdb.Result, error) {
	if !rows.Next() {
		return nil, nil
	}
	// Column names and types.
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columnTypes := make([]string, len(columns))
	columnNames := make([]string, len(columns))
	for k, v := range columns {
		columnTypes[k] = v.DatabaseTypeName()
		columnNames[k] = v.Name()
	}
	var (
		values   = make([]interface{}, len(columnNames))
		records  = make(gdb.Result, 0)
		scanArgs = make([]interface{}, len(values))
	)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err := rows.Scan(scanArgs...); err != nil {
			return records, err
		}
		row := make(gdb.Record)
		for i, value := range values {
			if value == nil {
				row[columnNames[i]] = gvar.New(nil)
			} else {
				row[columnNames[i]] = gvar.New(ConvertFieldValueToLocalValue(value, columnTypes[i]))
			}
		}
		records = append(records, row)
		if !rows.Next() {
			break
		}
	}
	return records, nil
}

func ConvertFieldValueToLocalValue(fieldValue interface{}, fieldType string) interface{} {
	// If there's no type retrieved, it returns the `fieldValue` directly
	// to use its original data type, as `fieldValue` is type of interface{}.
	if fieldType == "" {
		return fieldValue
	}
	t, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	t = strings.ToLower(t)
	switch t {
	case
		"binary",
		"varbinary",
		"blob",
		"tinyblob",
		"mediumblob",
		"longblob":
		return gconv.Bytes(fieldValue)

	case
		"int",
		"tinyint",
		"small_int",
		"smallint",
		"medium_int",
		"mediumint",
		"serial":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint(gconv.String(fieldValue))
		}
		return gconv.Int(gconv.String(fieldValue))

	case
		"int8", // For pgsql, int8 = bigint.
		"big_int",
		"bigint",
		"bigserial":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint64(gconv.String(fieldValue))
		}
		return gconv.Int64(gconv.String(fieldValue))

	case "real":
		return gconv.Float32(gconv.String(fieldValue))

	case
		"float",
		"double",
		"decimal",
		"money",
		"numeric",
		"smallmoney":
		return gconv.Float64(gconv.String(fieldValue))

	case "bit":
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") {
			return 1
		}
		if strings.EqualFold(s, "false") {
			return 0
		}
		return gbinary.BeDecodeToInt64(gconv.Bytes(fieldValue))

	case "bool":
		return gconv.Bool(fieldValue)

	case "date":
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t).Format("Y-m-d")
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.Format("Y-m-d")

	case
		"datetime",
		"timestamp",
		"timestamptz":
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t)
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.String()

	default:
		// Auto detect field type, using key match.
		switch {
		case strings.Contains(t, "text") || strings.Contains(t, "char") || strings.Contains(t, "character"):
			return gconv.String(fieldValue)

		case strings.Contains(t, "float") || strings.Contains(t, "double") || strings.Contains(t, "numeric"):
			return gconv.Float64(gconv.String(fieldValue))

		case strings.Contains(t, "bool"):
			return gconv.Bool(gconv.String(fieldValue))

		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			return fieldValue

		case strings.Contains(t, "int"):
			return gconv.Int(gconv.String(fieldValue))

		case strings.Contains(t, "time"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.String()

		case strings.Contains(t, "date"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.Format("Y-m-d")

		default:
			return gconv.String(fieldValue)
		}
	}
}
