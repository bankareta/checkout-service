package helpers

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// New membuat instance LogrusLogger sebagai middleware
func MinifyJson(r interface{}) map[string]interface{} {
	js, _ := jsoniter.Marshal(r)
	var m map[string]interface{}
	_ = jsoniter.Unmarshal(js, &m)
	for k, v := range m {
		s := fmt.Sprintf("%v", v)
		if len(s) > 10000 {
			delete(m, k)
		}
	}
	return m
}

func StructToMap(v interface{}) map[string]interface{} {
	js, err := jsoniter.Marshal(v)
	if err != nil {
		return nil
	}

	var maps map[string]interface{}

	err = jsoniter.Unmarshal(js, &maps)
	if err != nil {
		return nil
	}

	return maps
}
