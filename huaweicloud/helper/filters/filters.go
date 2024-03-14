package filters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/thedevsaddam/gojsonq"
)

type QueryCond struct {
	Key      string
	Operator string
	Value    any
}

type JsonFilter struct {
	query    *gojsonq.JSONQ
	node     string
	jsonData any
	queries  []QueryCond
}

func (f *JsonFilter) GetQ() *gojsonq.JSONQ {
	return f.query
}

func New() *JsonFilter {
	return &JsonFilter{
		query: gojsonq.New().
			Macro("has", has).
			Macro("hasContains", hasContain),
		queries: make([]QueryCond, 0),
	}
}

func (f *JsonFilter) Data(data any) *JsonFilter {
	f.jsonData = data
	return f
}

func (f *JsonFilter) From(node string) *JsonFilter {
	f.node = node
	return f
}

func (f *JsonFilter) GetFrom() string {
	return f.node
}

func (f *JsonFilter) GetQueries() []QueryCond {
	return f.queries
}

func (f *JsonFilter) Where(key, operator string, val any) *JsonFilter {
	if val == nil {
		return f
	}
	f.queries = append(f.queries, QueryCond{
		Key:      key,
		Operator: operator,
		Value:    val,
	})
	return f
}

func (f *JsonFilter) Get() (any, error) {
	dt := reflect.TypeOf(f.jsonData)
	if dt.Kind() == reflect.Slice {
		return f.filterSlice()
	}
	return f.filterJson()
}

func (f *JsonFilter) filterSlice() (any, error) {
	b, err := json.Marshal(f.jsonData)
	if err != nil {
		return f.jsonData, err
	}

	query := f.query.JSONString(string(b))
	for _, q := range f.queries {
		query = query.Where(q.Key, q.Operator, q.Value)
	}
	return query.Get(), nil
}

func (f *JsonFilter) filterJson() (any, error) {
	if f.node == "" {
		return nil, fmt.Errorf("`From` cannot be empty")
	}

	b, err := json.Marshal(f.jsonData)
	if err != nil {
		return f.jsonData, err
	}

	query := f.query.JSONString(string(b)).From(f.node)

	for _, q := range f.queries {
		query = query.Where(q.Key, q.Operator, q.Value)
	}

	if f.node == "" {
		return query.Get(), nil
	}

	switch mp := f.jsonData.(type) {
	case map[string]interface{}:
		mp[f.node] = query.Get()
		return mp, nil
	default:
		return nil, fmt.Errorf("failed to parse object")
	}
}

func toStrSlice(in any) []string {
	rf := reflect.ValueOf(in)
	if rf.Kind() != reflect.Array && rf.Kind() != reflect.Slice {
		return nil
	}

	rst := make([]string, 0)
	for i := 0; i < rf.Len(); i++ {
		rst = append(rst, fmt.Sprintf("%v", rf.Index(i).Interface()))
	}

	return rst
}

// Contains reports whether v is present in s.
// Copied from the source code of go1.21.
func contains[S ~[]E, E comparable](s S, v E) bool {
	return index(s, v) >= 0
}

// Index returns the index of the first occurrence of v in s, or -1 if not present.
// Copied from the source code of go1.21.
func index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func has(x interface{}, y interface{}) (bool, error) {
	xVal := toStrSlice(x)
	rf := reflect.ValueOf(y)
	if rf.IsValid() && rf.Kind() != reflect.Slice {
		yy := fmt.Sprintf("%v", y)
		return contains(xVal, yy), nil
	}

	yy := toStrSlice(y)
	for _, v := range yy {
		if contains(xVal, v) {
			return true, nil
		}
	}

	return false, nil
}

func hasContain(x interface{}, y interface{}) (bool, error) {
	xArr := toStrSlice(x)
	yArr := make([]string, 0)
	rf := reflect.ValueOf(y)
	if rf.IsValid() && rf.Kind() != reflect.Slice {
		yy := fmt.Sprintf("%v", y)
		yArr = append(yArr, yy)
	} else {
		yArr = toStrSlice(y)
	}

	for _, yv := range yArr {
		for _, xv := range xArr {
			if strings.Contains(xv, yv) {
				return true, nil
			}
		}
	}
	return false, nil
}
