package filters

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/thedevsaddam/gojsonq"
	"github.com/tidwall/gjson"
)

type Filter func(item gjson.Result) bool

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
	filter   Filter
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
		filter:  nil,
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

func (f *JsonFilter) Filter(filter Filter) *JsonFilter {
	f.filter = filter
	return f
}

func (f *JsonFilter) GetFilter() Filter {
	return f.filter
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
	return f.applyFilter(query.Get()), nil
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
		mp = putMap(f.node, mp, f.applyFilter(query.Get()))
		return mp, nil
	default:
		return nil, fmt.Errorf("failed to parse object")
	}
}

func (f *JsonFilter) applyFilter(slice any) any {
	if f.filter == nil || slice == nil {
		return slice
	}

	rv := reflect.ValueOf(slice)
	if !rv.IsValid() || rv.IsNil() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return slice
	}

	resultData := make([]any, 0)
	for i := 0; i < rv.Len(); i++ {
		val := rv.Index(i).Interface()
		b, err := json.Marshal(val)
		if err != nil {
			log.Printf("[ERROR] failed to apply custom filters: %s", err)
			continue
		}

		if f.filter(gjson.ParseBytes(b)) {
			resultData = append(resultData, val)
		}
	}

	return resultData
}

func putMap(keyPath string, mp map[string]any, val any) map[string]any {
	keys := strings.Split(keyPath, ".")

	node := mp
	for i := 0; i < len(keys)-1; i++ {
		v, ok := node[keys[i]].(map[string]any)
		if !ok {
			return mp
		}
		node = v
	}

	key := keys[len(keys)-1]
	node[key] = val
	return mp
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
	if y == nil {
		return true, nil
	}

	rf := reflect.ValueOf(x)
	switch rf.Kind() {
	case reflect.Array, reflect.Slice:
		return sliceHas(x, y)
	case reflect.Map:
		return mapHas(x, y)
	default:
		return false, fmt.Errorf("[has] unsupported comparison type: %s", rf.Kind())
	}
}

func mapHas(x interface{}, y interface{}) (bool, error) {
	xRef := reflect.ValueOf(x)
	yRef := reflect.ValueOf(y)
	if xRef.Kind() != reflect.Map || yRef.Kind() != reflect.Map {
		return false, fmt.Errorf("[mapHas] types must all be map: %s %s", xRef.Kind(), xRef.Kind())
	}
	for _, k := range yRef.MapKeys() {
		yVal := yRef.MapIndex(k)
		xVal := xRef.MapIndex(k)
		if xVal.IsValid() && !xVal.IsNil() && isEqual(xVal, yVal) {
			continue
		}
		return false, nil
	}

	return true, nil
}

func sliceHas(x interface{}, y interface{}) (bool, error) {
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
	rf := reflect.ValueOf(x)
	switch rf.Kind() {
	case reflect.Array, reflect.Slice:
		return sliceHasContain(x, y)
	case reflect.Map:
		return mapHasContain(x, y)
	default:
		return false, fmt.Errorf("[hasContain] unsupported comparison type: %s", rf.Kind())
	}
}

func sliceHasContain(x interface{}, y interface{}) (bool, error) {
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

func mapHasContain(x interface{}, y interface{}) (bool, error) {
	if reflect.TypeOf(x).Kind() != reflect.Map || reflect.TypeOf(y).Kind() != reflect.Map {
		return false, fmt.Errorf("[mapHas] types must all be map: %s %s", reflect.TypeOf(x), reflect.TypeOf(y))
	}

	xRef := reflect.ValueOf(x)
	yRef := reflect.ValueOf(y)

	keys := yRef.MapKeys()
	for _, k := range keys {
		yVal := yRef.MapIndex(k)
		xVal := xRef.MapIndex(k)
		if xVal.IsValid() && !xVal.IsNil() && isEqual(xVal, yVal) {
			return true, nil
		}
	}

	return false, nil
}

// isEqual compare two basic type variables for equality
// x and y may be of int32 and interface type, and are converted to string before comparison
func isEqual(x, y reflect.Value) bool {
	if x.Kind() == reflect.Pointer {
		x = x.Elem()
	}
	if y.Kind() == reflect.Pointer {
		y = y.Elem()
	}

	vx := fmt.Sprintf("%v", x.Interface())
	vy := fmt.Sprintf("%v", y.Interface())
	return vx == vy
}
