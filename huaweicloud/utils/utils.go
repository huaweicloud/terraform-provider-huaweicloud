package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"
)

// ConvertStructToMap converts an instance of struct to a map object, and
// changes each key of fileds to the value of 'nameMap' if the key in it
// or to its corresponding lowercase.
func ConvertStructToMap(obj interface{}, nameMap map[string]string) (map[string]interface{}, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, marshal failed:%v", err)
	}

	m, err := regexp.Compile(`"[a-z0-9A-Z_]+":`)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, compile regular express failed")
	}
	nb := m.ReplaceAllFunc(
		b,
		func(src []byte) []byte {
			k := string(src[1 : len(src)-2])
			v, ok := nameMap[k]
			if !ok {
				v = strings.ToLower(k)
			}
			return []byte(fmt.Sprintf("\"%s\":", v))
		},
	)
	log.Printf("[DEBUG]convertStructToMap:: before change b =%s", b)
	log.Printf("[DEBUG]convertStructToMap:: after change nb=%s", nb)

	p := make(map[string]interface{})
	err = json.Unmarshal(nb, &p)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]convertStructToMap:: map= %#v\n", p)
	return p, nil
}

// ExpandToStringList takes the result for an array of strings and returns a []string
func ExpandToStringList(v []interface{}) []string {
	s := make([]string, 0, len(v))
	for _, val := range v {
		if strVal, ok := val.(string); ok && strVal != "" {
			s = append(s, strVal)
		}
	}

	return s
}

// ExpandToStringListPointer takes the result for an array of strings and returns a pointer of the array
func ExpandToStringListPointer(v []interface{}) *[]string {
	s := ExpandToStringList(v)

	return &s
}

// ExpandToIntList takes the result for an array of intgers and returns a []int
func ExpandToIntList(v []interface{}) []int {
	s := make([]int, 0, len(v))
	for _, val := range v {
		if intVal, ok := val.(int); ok {
			s = append(s, intVal)
		}
	}
	return s
}

// ExpandToInt32List takes the result for an array of intgers and returns a []int32
func ExpandToInt32List(v []interface{}) []int32 {
	s := make([]int32, 0, len(v))
	for _, val := range v {
		if intVal, ok := val.(int); ok {
			s = append(s, int32(intVal))
		}
	}
	return s
}

// ExpandToInt32ListPointer takes the result for an array of in32 and returns a pointer of the array
func ExpandToInt32ListPointer(v []interface{}) *[]int32 {
	s := ExpandToInt32List(v)

	return &s
}

// ExpandToStringListBySet takes the result for a set of strings and returns a []string
func ExpandToStringListBySet(v *schema.Set) []string {
	s := make([]string, 0, v.Len())
	for _, val := range v.List() {
		if strVal, ok := val.(string); ok && strVal != "" {
			s = append(s, strVal)
		}
	}

	return s
}

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
func flattenToStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func pointersMapToStringList(pointers map[string]*string) map[string]interface{} {
	list := make(map[string]interface{}, len(pointers))
	for i, v := range pointers {
		list[i] = *v
	}
	return list
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func NormalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

// StrSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func StrSliceContains(haystack []string, needle string) bool {
	return IsStrContainsSliceElement(needle, haystack, false, true)
}

// StrSliceContainsAnother checks whether a string slice (b) contains another string slice (s).
func StrSliceContainsAnother(b []string, s []string) bool {
	// The empty set is the subset of any set.
	if len(s) < 1 {
		return true
	}
	for _, v := range s {
		if !StrSliceContains(b, v) {
			return false
		}
	}
	return true
}

// IsStrContainsSliceElement returns true if the string exists in given slice or contains in one of slice elements when
// open exact flag. Also you can ignore case for this check.
func IsStrContainsSliceElement(str string, sl []string, ignoreCase, isExcat bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for _, s := range sl {
		if ignoreCase {
			s = strings.ToLower(s)
		}
		if isExcat && s == str {
			return true
		}
		if !isExcat && strings.Contains(str, s) {
			return true
		}
	}
	return false
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)
	return buffer.Bytes(), err
}

// RemoveDuplicateElem removes duplicate elements from slice
func RemoveDuplicateElem(s []string) []string {
	result := []string{}
	tmpMap := map[string]byte{}
	for _, e := range s {
		l := len(tmpMap)
		tmpMap[e] = 0
		if len(tmpMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func RemoveNil(data map[string]interface{}) map[string]interface{} {
	withoutNil := make(map[string]interface{})

	for k, v := range data {
		if v == nil {
			continue
		}

		switch v := v.(type) {
		case map[string]interface{}:
			withoutNil[k] = RemoveNil(v)
		default:
			withoutNil[k] = v
		}
	}

	return withoutNil
}

func IsResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(golangsdk.ErrDefault404)
	return ok
}

// FormatTimeStampRFC3339 is used to unify the time format to RFC-3339 and return a time string.
func FormatTimeStampRFC3339(timestamp int64) string {
	createTime := time.Unix(timestamp, 0)
	return createTime.Format(time.RFC3339)
}

// FormatTimeStampUTC is used to unify the unix second time to UTC time string, format: YYYY-MM-DD HH:MM:SS.
func FormatTimeStampUTC(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02 15:04:05")
}

// FormatTimeStampUTC is used to unify the unix second time to UTC time string, format: YYYY-MM-DD HH:MM:SS.
func FormatUTCTimeStamp(utcTime string) (int64, error) {
	timestamp, err := time.Parse("2006-01-02 15:04:05", utcTime)
	if err != nil {
		return 0, fmt.Errorf("unable to prase the time: %s", utcTime)
	}
	return timestamp.Unix(), nil
}

// EncodeBase64String is used to encode a string by base64.
func EncodeBase64String(str string) string {
	strByte := []byte(str)
	return base64.StdEncoding.EncodeToString(strByte)
}

// EncodeBase64IfNot is used to encode a string by base64 if it not a base64 string.
func EncodeBase64IfNot(str string) string {
	if _, err := base64.StdEncoding.DecodeString(str); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return str
}

// IsIPv4Address is used to check whether the addr string is IPv4 format
func IsIPv4Address(addr string) bool {
	pattern := "^((25[0-5]|2[0-4]\\d|(1\\d{2}|[1-9]?\\d))\\.){3}(25[0-5]|2[0-4]\\d|(1\\d{2}|[1-9]?\\d))$"
	matched, _ := regexp.MatchString(pattern, addr)
	return matched
}

// This function compares whether there is a containment relationship between two maps, that is,
// whether map A (rawMap) contains map B (filterMap).
//
//	Map A is {'foo': 'bar'} and filter map B is {'foo': 'bar'} or {'foo': 'bar,dor'} will return true.
//	Map A is {'foo': 'bar'} and filter map B is {'foo': 'dor'} or {'foo1': 'bar'} will return false.
//	Map A is {'foo': 'bar'} and filter map B is {'foo': ''} will return true.
//	Map A is {'foo': 'bar'} and filter map B is {'': 'bar'} or {'': ''} will return false.
//
// The value of filter map 'bar,for' means that the object value can be either 'bar' or 'dor'.
// Note: There is no spaces before and after the delimiter (,).
func HasMapContains(rawMap map[string]string, filterMap map[string]interface{}) bool {
	if len(filterMap) < 1 {
		return true
	}

	hasContain := true
	for key, value := range filterMap {
		hasContain = hasContain && hasMapContain(rawMap, key, value.(string))
	}

	return hasContain
}

func hasMapContain(rawMap map[string]string, filterKey, filterValue string) bool {
	if rawTag, ok := rawMap[filterKey]; ok {
		if filterValue != "" {
			filterTagValues := strings.Split(filterValue, ",")
			return StrSliceContains(filterTagValues, rawTag)
		} else {
			return true
		}
	} else {
		return false
	}
}

// WriteToPemFile is used to write the keypair to Pem file.
func WriteToPemFile(path, privateKey string) (err error) {
	// If the private key exists, give it write permission for editing (-rw-------) for root user.
	if _, err = ioutil.ReadFile(path); err == nil {
		err = os.Chmod(path, 0600)
		if err != nil {
			return
		}

		defer func() {
			// read-only permission (-r--------).
			mErr := multierror.Append(err, os.Chmod(path, 0400))
			err = mErr.ErrorOrNil()
		}()
	}
	if err = ioutil.WriteFile(path, []byte(privateKey), 0600); err != nil {
		return err
	}
	return nil
}

/*
MarshalValue is used to marshal the value of struct in huaweicloud-sdk-go-v3, like this:
type Xxxx struct { value string }
*/
func MarshalValue(i interface{}) string {
	if i == nil {
		return ""
	}

	jsonRaw, err := json.Marshal(i)
	if err != nil {
		log.Printf("[WARN] failed to marshal %#v: %s", i, err)
		return ""
	}

	return strings.Trim(string(jsonRaw), `"`)
}

// RandomString returns a random string with a fixed length. You can also define a custom random character set.
// Note: make sure the number is not a negative integer or a big integer.
func RandomString(n int, allowedChars ...[]rune) (result string) {
	var letters []rune

	if len(allowedChars) == 0 {
		// Using default seed.
		letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	} else {
		letters = allowedChars[0]
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] The number (input n) cannot be a negative integer or a large integer: %#v", r)
		}
	}()
	b := make([]rune, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}

	result = string(b)
	return
}

// IsDebugOrHigher returns a bool type parameter, which specifies whether to print log
var validLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}

func IsDebugOrHigher() bool {
	logLevel := os.Getenv("TF_LOG_PROVIDER")
	if logLevel == "" {
		logLevel = os.Getenv("TF_LOG")
	}

	if logLevel != "" {
		if isValidLogLevel(logLevel) {
			logLevel = strings.ToUpper(logLevel)
			return logLevel == "DEBUG" || logLevel == "TRACE"
		} else {
			log.Printf("[WARN] Invalid log level: %q. Valid levels are: %+v", logLevel, validLevels)
		}
	}
	return false
}

func isValidLogLevel(level string) bool {
	for _, l := range validLevels {
		if strings.ToUpper(level) == string(l) {
			return true
		}
	}

	return false
}

// PathSearch evaluates a JMESPath expression against input data and returns the result.
func PathSearch(expression string, obj interface{}, defaultValue interface{}) interface{} {
	v, err := jmespath.Search(expression, obj)
	if err != nil || v == nil {
		return defaultValue
	}
	return v
}

// FlattenResponse returns the api response body if it's not empty
func FlattenResponse(resp *http.Response) (interface{}, error) {
	var respBody interface{}
	defer resp.Body.Close()
	// Don't decode JSON when there is no content
	if resp.StatusCode == http.StatusNoContent {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		return resp, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return respBody, nil
}
