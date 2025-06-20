package processor

import (
	"fmt"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/largehong/lkv/memkv"
)

func Regexp(expr string, item string) (map[string]string, error) {
	r, err := regexp2.Compile(expr, regexp2.RE2)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	match, err := r.FindStringMatch(item)
	if err != nil {
		return nil, err
	}
	for _, group := range match.Groups() {
		if group.Name == "0" {
			continue
		}
		data[group.Name] = group.Capture.String()
	}
	return data, nil
}

func Regexps(expr string, items []string) (map[string][]string, error) {
	r, err := regexp2.Compile(expr, regexp2.RE2)
	if err != nil {
		return nil, err
	}

	data := make(map[string][]string)

	for _, s := range items {
		match, err := r.FindStringMatch(s)
		if err != nil {
			return nil, err
		}
		for _, group := range match.Groups() {
			if group.Name == "0" {
				continue
			}
			values, ok := data[group.Name]
			if !ok {
				values = make([]string, 0)
			}
			values = append(values, group.Capture.String())
			data[group.Name] = values
		}
	}
	return data, nil
}

func Unique(items []string) []string {
	m := make(map[string]struct{})
	for _, item := range items {
		m[item] = struct{}{}
	}

	d := make([]string, 0)
	for k := range m {
		d = append(d, k)
	}

	return d
}

func GetMemKVKeys(kvs []memkv.KV) []string {
	keys := make([]string, 0)
	for _, kv := range kvs {
		keys = append(keys, kv.Key)
	}
	return keys
}

func GetMemKVValues(kvs []memkv.KV) []any {
	values := make([]any, 0)
	for _, kv := range kvs {
		values = append(values, kv.Value)
	}
	return values
}

func FuncMaps() map[string]any {
	return map[string]any{
		"regexps":   Regexps,
		"regexp":    Regexp,
		"unique":    Unique,
		"keys":      GetMemKVKeys,
		"values":    GetMemKVValues,
		"join":      strings.Join,
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"title":     strings.ToTitle,
		"now":       time.Now,
		"timestamp": time.Now().Unix,
		"sprintf":   fmt.Sprintf,
	}
}
