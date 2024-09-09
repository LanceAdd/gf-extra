package ops

import (
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func JsonExtract(jsonOrRawSql string, key string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_EXTRACT(%s, '$%s') AS %s", jsonOrRawSql, key, alias[0])
	}
	return fmt.Sprintf("JSON_EXTRACT(%s, '$%s')", jsonOrRawSql, key)
}

func JsonArrayExtract(jsonOrRawSql string, index int, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_EXTRACT(%s, '$[%d]') AS %s", jsonOrRawSql, index, alias[0])
	}
	return fmt.Sprintf("JSON_EXTRACT(%s, '$[%d]')", jsonOrRawSql, index)
}

func JsonSet(jsonOrRawSql string, key string, data any, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	encode, _ := gjson.EncodeString(data)
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_SET(%s, '$%s', %s) AS %s", jsonOrRawSql, key, encode, alias[0])
	}
	return fmt.Sprintf("JSON_SET(%s, '$%s', %s)", jsonOrRawSql, key, encode)
}

func JsonSetColumn(rawString string, key string, column string, alias ...string) string {
	if IsBlank(rawString) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_SET(%s, '$%s', %s) AS %s", rawString, key, column, alias[0])
	}
	return fmt.Sprintf("JSON_SET(%s, '$%s', %s)", rawString, key, column)
}

func JsonArraySet(jsonOrRawSql string, index int, data any, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	encode, _ := gjson.EncodeString(data)
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_SET(%s, '$[%d]', %s) AS %s", jsonOrRawSql, index, encode, alias[0])
	}
	return fmt.Sprintf("JSON_SET(%s, '$[%d]', %s)", jsonOrRawSql, index, encode)
}

func JsonContains(jsonOrRawSql string, data any, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	encode, _ := gjson.EncodeString(data)
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_CONTAINS(%s, '%s', '$') AS `%s`", jsonOrRawSql, encode, alias[0])
	}
	return fmt.Sprintf("JSON_CONTAINS(%s, '%s')", jsonOrRawSql, encode)
}

func JsonLength(jsonOrRawSql string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_LENGTH(%s) AS `%s`", jsonOrRawSql, alias[0])
	}
	return fmt.Sprintf("JSON_LENGTH(%s)", jsonOrRawSql)
}

func JsonDepth(jsonOrRawSql string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_DEPTH(%s) AS `%s`", jsonOrRawSql, alias[0])
	}
	return fmt.Sprintf("JSON_DEPTH(%s)", jsonOrRawSql)
}

func JsonType(jsonOrRawSql string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_TYPE(%s) AS `%s`", jsonOrRawSql, alias[0])
	}
	return fmt.Sprintf("JSON_TYPE(%s)", jsonOrRawSql)
}

func JsonValid(jsonOrRawSql string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_VALID(%s) AS `%s`", jsonOrRawSql, alias[0])
	}
	return fmt.Sprintf("JSON_VALID(%s)", jsonOrRawSql)
}

func JsonStorageSize(jsonOrRawSql string, alias ...string) string {
	if IsBlank(jsonOrRawSql) {
		return ""
	}
	if len(alias) > 0 {
		return fmt.Sprintf("JSON_STORAGE_SIZE(%s) AS `%s`", jsonOrRawSql, alias[0])
	}
	return fmt.Sprintf("JSON_STORAGE_SIZE(%s)", jsonOrRawSql)
}
