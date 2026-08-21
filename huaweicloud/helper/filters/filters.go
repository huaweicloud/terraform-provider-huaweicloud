package filters

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Filter func(item gjson.Result) bool

type QueryCond struct {
	Key      string
	Operator string
	Value    any
}

type JsonFilter struct {
	node     string
	jsonData any
	queries  []QueryCond
	filter   Filter
}

func New() *JsonFilter {
	return &JsonFilter{
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
	if dt != nil && dt.Kind() == reflect.Slice {
		return f.filterSlice()
	}
	return f.filterJson()
}

func (f *JsonFilter) filterSlice() (any, error) {
	if len(f.queries) > 0 || f.filter != nil {
		items, err := normalizeSliceViaJSON(f.jsonData)
		if err != nil {
			return f.jsonData, err
		}

		var filtered interface{} = items
		if len(f.queries) > 0 {
			filtered = filterArrayItems(items, f.queries)
		}
		return f.applyFilter(filtered), nil
	}

	rv := reflect.ValueOf(f.jsonData)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return f.jsonData, nil
	}

	items := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		items[i] = rv.Index(i).Interface()
	}
	return f.applyFilter(items), nil
}

func (f *JsonFilter) filterJson() (any, error) {
	if f.node == "" {
		return nil, fmt.Errorf("`From` cannot be empty")
	}

	mp, ok := f.jsonData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse object")
	}

	if len(f.queries) > 0 || f.filter != nil {
		normalized, err := normalizeViaJSON(mp)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize filter data: %w", err)
		}
		mp = normalized
	}

	nodeValue, err := getNestedValue(mp, f.node, ".")
	var filtered interface{}
	if err != nil {
		filtered = nil
	} else if arr, ok := normalizeToInterfaceSlice(nodeValue); ok {
		if len(f.queries) == 0 {
			filtered = arr
		} else {
			filtered = filterArrayItems(arr, f.queries)
		}
	} else {
		filtered = nodeValue
	}

	mp = putMap(f.node, mp, f.applyFilter(filtered))
	return mp, nil
}

func normalizeToInterfaceSlice(v interface{}) ([]interface{}, bool) {
	if arr, ok := v.([]interface{}); ok {
		return arr, true
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, false
	}

	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result, true
}

func filterArrayItems(items []interface{}, queries []QueryCond) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range items {
		vm, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if matchMapConditions(vm, queries) {
			result = append(result, vm)
		}
	}
	return result
}

func matchMapConditions(vm map[string]interface{}, queries []QueryCond) bool {
	for _, q := range queries {
		nv, err := getNestedValue(vm, q.Key, ".")
		if err != nil {
			return false
		}
		matched, err := matchCondition(nv, q.Operator, q.Value)
		if err != nil || !matched {
			return false
		}
	}
	return true
}

func matchCondition(x interface{}, operator string, y interface{}) (bool, error) {
	switch operator {
	case "=":
		return equalValues(x, y), nil
	case "contains":
		return strContainsCondition(x, y)
	case "has":
		return has(x, y)
	case "hasContains":
		return hasContain(x, y)
	default:
		return false, fmt.Errorf("invalid operator %s", operator)
	}
}

func equalValues(x, y interface{}) bool {
	fx, okX := toFloat64(x)
	fy, okY := toFloat64(y)
	if okX && okY {
		return fx == fy
	}
	return reflect.DeepEqual(x, y)
}

func normalizeViaJSON(data any) (map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func normalizeSliceViaJSON(data any) ([]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func strContainsCondition(x, y interface{}) (bool, error) {
	xv, okX := x.(string)
	if !okX {
		return false, fmt.Errorf("%v must be string", x)
	}
	yv, okY := y.(string)
	if !okY {
		return false, fmt.Errorf("%v must be string", y)
	}
	return strings.Contains(strings.ToLower(xv), strings.ToLower(yv)), nil
}

func toFloat64(v interface{}) (float64, bool) {
	switch u := v.(type) {
	case int:
		return float64(u), true
	case int8:
		return float64(u), true
	case int16:
		return float64(u), true
	case int32:
		return float64(u), true
	case int64:
		return float64(u), true
	case float32:
		return float64(u), true
	case float64:
		return u, true
	default:
		return 0, false
	}
}

func isIndex(in string) bool {
	return strings.HasPrefix(in, "[") && strings.HasSuffix(in, "]")
}

func getIndex(in string) (int, error) {
	if !isIndex(in) {
		return -1, fmt.Errorf("invalid index")
	}
	is := strings.TrimLeft(in, "[")
	is = strings.TrimRight(is, "]")
	return strconv.Atoi(is)
}

func getNestedValue(input interface{}, node, separator string) (interface{}, error) {
	pp := strings.Split(node, separator)
	for _, n := range pp {
		if isIndex(n) {
			arr, ok := input.([]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid node name %s", n)
			}
			indx, err := getIndex(n)
			if err != nil {
				return input, err
			}
			arrLen := len(arr)
			if arrLen == 0 || indx > arrLen-1 {
				return nil, errors.New("empty array")
			}
			input = arr[indx]
		} else {
			validNode := false
			if mp, ok := input.(map[string]interface{}); ok {
				input, ok = mp[n]
				validNode = ok
			}

			if mp, ok := input.(map[string][]interface{}); ok {
				input, ok = mp[n]
				validNode = ok
			}

			if !validNode {
				return nil, fmt.Errorf("invalid node name %s", n)
			}
		}
	}

	return input, nil
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

func valueIsNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
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
		if xVal.IsValid() && !valueIsNil(xVal) && isEqual(xVal, yVal) {
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
		if xVal.IsValid() && !valueIsNil(xVal) && isEqual(xVal, yVal) {
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
