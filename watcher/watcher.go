package watcher

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
)

func doToPropertiesMap(source *map[string]any) *map[string]any {
	dataMap := map[string]any{}
	json, err := gjson.DecodeToJson(source)
	if err != nil {
		panic(err)
	}
	propertiesString := json.MustToPropertiesString()
	propertiesArray := gstr.SplitAndTrim(propertiesString, "\n")
	for _, item := range propertiesArray {
		trim := gstr.SplitAndTrim(item, "=")
		dataMap[trim[0]] = trim[1]
	}
	return &dataMap
}

func Compare(source *map[string]any, target *map[string]any) *map[string]any {
	m := map[string]any{}
	sourceMap := *source
	targetMap := *target
	for key, value := range sourceMap {
		targetValue := targetMap[key]
		if targetValue == nil || value != targetValue {
			m[key] = value
			continue
		}
	}
	return &m
}
