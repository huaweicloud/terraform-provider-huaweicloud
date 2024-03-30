package httphelper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/stretchr/testify/assert"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
)

var testCaseMapper = make(map[string]*HttpTestCase)

func newHttpMockServer(t *testing.T) *httptest.Server {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		testCase, ok := testCaseMapper[req.Header.Get("input-key")]
		if !ok {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		// assert method
		assert.Equal(t, testCase.method, req.Method)
		// assert query
		assert.Equal(t, testCase.getQuery(), req.URL.Query().Encode())
		// assert body
		b, _ := io.ReadAll(req.Body)
		body := string(b)
		if body == "" || body == "null\n" {
			body = "{}"
		}
		assert.Equal(t, testCase.getBody(), strings.TrimSpace(body))

		// response
		if testCase.httpCode != 0 {
			w.WriteHeader(testCase.httpCode)
		}
		_, _ = fmt.Fprint(w, testCase.getBody())
	}))
	return server
}

func TestGet(t *testing.T) {
	server := newHttpMockServer(t)
	defer server.Close()

	testCases := []*HttpTestCase{
		{
			name:   "test GET 1",
			method: "GET",
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name:   "test GET 2",
			method: "GET",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name:   "test GET 3",
			method: "GET",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			httpCode: 201,
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rst, err := tc.newClient(server).
				OkCode(200).
				Request().
				Result()

			if tc.httpCode != 0 && tc.httpCode != 200 {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &golangsdk.ErrUnexpectedResponseCode{})
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.getBody(), strings.TrimSpace(rst.Raw))
			}
		})
	}
}

func TestPost(t *testing.T) {
	server := newHttpMockServer(t)
	defer server.Close()

	testCases := []*HttpTestCase{
		{
			name:   "test POST 1",
			method: "POST",
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name:   "test POST 2",
			method: "POST",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
				"arr1": []string{"1", "2"},
				"arr2": []float32{3.1455555555, 111.222222, 1},
				"arr3": []bool{true, false, true},
				"att1": 1,
				"att2": false,
				"att3": true,
				"att4": map[string][]string{
					"val": []string{"val1", "val2", "val3"},
				},
				"att5": map[string]any{
					"mv1": "str",
					"mv2": 12,
					"mv3": false,
					"mv4": []string{"1", "2"},
					"mv5": []int{1, 2, 3},
					"mv6": []bool{true, false, true},
				},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name:   "test POST 3",
			method: "POST",
			body: map[string]any{
				"users": map[string]any{
					"name": "zhangsan",
					"age":  12,
				},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name:   "test POST 4",
			method: "POST",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			body: map[string]any{
				"users": map[string]any{
					"name": "zhangsan",
					"age":  12,
				},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rst, err := tc.newClient(server).
				OkCode(200).
				Request().
				Result()

			if tc.httpCode != 0 && tc.httpCode != 200 {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &golangsdk.ErrUnexpectedResponseCode{})
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.getBody(), strings.TrimSpace(rst.Raw))
			}
		})
	}
}

func TestPostPutPatchDelete(t *testing.T) {
	server := newHttpMockServer(t)
	defer server.Close()

	basicTc := []*HttpTestCase{
		{
			name: "no_param",
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name: "no_body_param",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name: "no_query_param",
			body: map[string]any{
				"users": map[string]any{
					"name": "zhangsan",
					"age":  12,
				},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name: "all_param",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			body: map[string]any{
				"users": map[string]any{
					"name": "zhangsan",
					"age":  12,
				},
			},
			response: map[string]any{
				"users": []map[string]any{
					{
						"name": "zhangsan",
						"age":  12,
					},
					{
						"name": "lisi",
						"age":  13,
					},
					{
						"name": "wangwu",
						"age":  16,
					},
				},
			},
		}, {
			name: "no_response",
			query: map[string]any{
				"name": "zhangsan",
				"age":  []int{12, 13, 14},
			},
			body: map[string]any{
				"users": map[string]any{
					"name": "zhangsan",
					"age":  12,
				},
			},
		},
	}

	testCases := make([]*HttpTestCase, 0)
	for i, v := range basicTc {
		tc := v.copy()
		tc.method = "PUT"
		tc.name = fmt.Sprintf("test %s %s %v", tc.method, tc.name, i+1)
		testCases = append(testCases, tc)
	}

	for i, v := range basicTc {
		tc := v.copy()
		tc.method = "PATCH"
		tc.name = fmt.Sprintf("test %s %s %v", tc.method, tc.name, i+1)
		testCases = append(testCases, tc)
	}

	for i, v := range basicTc {
		tc := v.copy()
		tc.method = "POST"
		tc.name = fmt.Sprintf("test %s %s %v", tc.method, tc.name, i+1)
		testCases = append(testCases, tc)
	}

	for i, v := range basicTc {
		tc := v.copy()
		tc.method = "DELETE"
		tc.name = fmt.Sprintf("test %s %s %v", tc.method, tc.name, i+1)
		testCases = append(testCases, tc)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hh := tc.newClient(server).
				OkCode(200).
				Request()

			rst, err := hh.Result()
			assert.NoError(t, err)

			if tc.httpCode != 0 && tc.httpCode != 200 {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &golangsdk.ErrUnexpectedResponseCode{})
			} else {
				assert.Equal(t, tc.getBody(), strings.TrimSpace(rst.Raw))
			}

			data, err := hh.Data()
			assert.NoError(t, err)
			if tc.httpCode != 0 && tc.httpCode != 200 {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &golangsdk.ErrUnexpectedResponseCode{})
			} else {
				b, err := json.Marshal(data)
				assert.NoError(t, err)
				body := string(b)
				assert.Equal(t, body, strings.TrimSpace(rst.Raw))
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	map1 := map[string]interface{}{
		"a": 1,
		"b": "hello",
		"c": []interface{}{1, 2, 3},
		"d": map[string]interface{}{
			"x": 10,
			"y": 20,
			"w": []interface{}{1, 2, 3},
		},
	}

	map2 := map[string]interface{}{
		"b": "world",
		"c": []interface{}{4, 5},
		"d": map[string]interface{}{
			"y": 30,
			"z": 40,
			"w": []interface{}{4, 5},
		},
		"e": "new",
	}

	expected := map[string]interface{}{
		"a": 1,
		"b": "world",
		"c": []interface{}{1, 2, 3, 4, 5},
		"d": map[string]interface{}{
			"x": 10,
			"y": 30,
			"z": 40,
			"w": []interface{}{1, 2, 3, 4, 5},
		},
		"e": "new",
	}

	rst := mergeMaps(map1, map2)
	assert.Equal(t, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", rst))
}

func TestMarshalQueryParams(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		expected string
	}{
		{
			name: "int",
			params: map[string]any{
				"zero":   0,
				"val":    12,
				"int_32": 2147483647,
				"int_64": 9223372036854775807,
			},
			expected: "?int_32=2147483647&int_64=9223372036854775807&val=12&zero=0",
		},
		{
			name: "string",
			params: map[string]any{
				"zero": "",
				"foo":  "bar",
				"arr1": []string{"a", "b", "", "c"},
				"arr2": []any{"a", "b", "", "c"},
			},
			expected: "?arr1=a&arr1=b&arr1=c&arr2=a&arr2=b&arr2=&arr2=c&foo=bar",
		},
		{
			name: "float32",
			params: map[string]any{
				"zero": float32(0),
				"val":  float32(123),
				"max":  math.MaxFloat32,
				"arr1": []float32{1, 2, 0, math.MaxFloat32},
				"arr2": []any{1, 2, 3, 0, math.MaxFloat32},
			},
			expected: "?arr1=1&arr1=2&arr1=0&arr1=3.4028235e%2B38&arr2=1&arr2=2&arr2=3&arr2=0&" +
				"arr2=3.4028234663852886e%2B38&max=3.4028234663852886e%2B38&val=123&zero=0",
		},
		{
			name: "float64",
			params: map[string]any{
				"zero": float64(0),
				"val":  float64(123),
				"max":  math.MaxFloat64,
				"arr1": []float64{1, 2, 3, 0, math.MaxFloat64},
				"arr2": []any{1, 2, 3, 0, math.MaxFloat64},
			},
			expected: "?arr1=1&arr1=2&arr1=3&arr1=0&arr1=1.7976931348623157e%2B308&arr2=1&arr2=2&arr2=3&arr2=0&" +
				"arr2=1.7976931348623157e%2B308&max=1.7976931348623157e%2B308&val=123&zero=0",
		},
		{
			name: "bool",
			params: map[string]any{
				"zero": false,
				"val":  true,
				"arr1": []bool{true, true, false, false},
				"arr2": []any{true, true, false, false},
			},
			expected: "?arr1=true&arr1=true&arr1=false&arr1=false&arr2=true&arr2=true&arr2=false&arr2=false&val=true&zero=false",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rst := marshalQueryParams(tc.params)
			fmt.Println(tc.name, rst)
			assert.Equal(t, tc.expected, rst)
		})
	}
}

type HttpTestCase struct {
	name     string
	method   string
	query    map[string]any
	body     map[string]any
	httpCode int
	response any
}

func (tc *HttpTestCase) copy() *HttpTestCase {
	return &HttpTestCase{
		name:     tc.name,
		method:   tc.method,
		query:    tc.query,
		body:     tc.body,
		httpCode: tc.httpCode,
		response: tc.response,
	}
}

func (tc *HttpTestCase) saveMapper() string {
	key := fmt.Sprintf("%p", tc)
	testCaseMapper[key] = tc
	return key
}

func (tc *HttpTestCase) getBody() string {
	if len(tc.body) == 0 {
		return "{}"
	}

	b, err := json.Marshal(tc.body)
	if err != nil {
		fmt.Println("failed to Marshal body:", err)
		return ""
	}
	return string(b)
}

func (tc *HttpTestCase) getQuery() string {
	params := marshalQueryParams(tc.query)
	return strings.Replace(params, "?", "", 1)
}

//nolint:gosec
func (tc *HttpTestCase) newClient(server *httptest.Server) *HttpHelper {
	client := server.Client()
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	svsClient := &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{
			HTTPClient: *client,
		},
		Endpoint: server.URL,
	}

	return New(svsClient).
		Method(tc.method).
		URI("test").
		Body(tc.body).
		Query(tc.query).
		Headers(map[string]string{
			"input-key": tc.saveMapper(),
		})
}

func TestHttpHelperURI(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	url1 := "https://github.com/"
	url2 := "http://fadsdfsaefs132.com"
	url3 := "github.com"
	url4 := "fads#!$_@dfsaefs132asdkml"
	url5 := ""

	helper.URI(url1)
	assert.Equal(t, url1, helper.url)

	helper.URI(url2)
	assert.Equal(t, url2, helper.url)

	helper.URI(url3)
	assert.Equal(t, url3, helper.url)

	helper.URI(url4)
	assert.Equal(t, url4, helper.url)

	helper.URI(url5)
	assert.Equal(t, url5, helper.url)
}

func TestHttpHelperMethod(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	method1 := "GET"
	method2 := "POST"
	method3 := "UPDATE"
	method4 := "DELETE"
	method5 := "get"
	method6 := "foo"
	method7 := ""

	helper.Method(method1)
	assert.Equal(t, method1, helper.method)

	helper.Method(method2)
	assert.Equal(t, method2, helper.method)

	helper.Method(method3)
	assert.Equal(t, method3, helper.method)

	helper.Method(method4)
	assert.Equal(t, method4, helper.method)

	helper.Method(method5)
	assert.Equal(t, method5, helper.method)

	helper.Method(method6)
	assert.Equal(t, method6, helper.method)

	helper.Method(method7)
	assert.Equal(t, method7, helper.method)
}

func TestHttpHelperBody(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	body1 := map[string]any{"foo1": "bar1"}
	body2 := map[string]any{"foo1": "bar1", "foo2": "bar2"}
	body3 := map[string]any{}
	body4 := map[string]any{"foo1": 123}
	body5 := map[string]any{"foo1": 123, "foo2": 456}
	body6 := map[string]any{"foo1": map[string]int{"foo111": 123}}

	helper.Body(body1)
	assert.Equal(t, body1, helper.body)

	helper.Body(body2)
	assert.Equal(t, body2, helper.body)

	helper.Body(body3)
	assert.Equal(t, body3, helper.body)

	helper.Body(body4)
	assert.Equal(t, body4, helper.body)

	helper.Body(body5)
	assert.Equal(t, body5, helper.body)

	helper.Body(body6)
	assert.Equal(t, body6, helper.body)
}

func TestHttpHelperQuery(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	query1 := map[string]any{"foo1": "bar1"}
	query2 := map[string]any{"foo1": "bar1", "foo2": "bar2"}
	query3 := map[string]any{}
	query4 := map[string]any{"foo1": 123}
	query5 := map[string]any{"foo1": 123, "foo2": 456}
	query6 := map[string]any{"foo1": map[string]int{"foo111": 123}}

	helper.Query(query1)
	assert.Equal(t, query1, helper.query)

	helper.Query(query2)
	assert.Equal(t, query2, helper.query)

	helper.Query(query3)
	assert.Equal(t, query3, helper.query)

	helper.Query(query4)
	assert.Equal(t, query4, helper.query)

	helper.Query(query5)
	assert.Equal(t, query5, helper.query)

	helper.Query(query6)
	assert.Equal(t, query6, helper.query)
}

func TestHttpHelperHeaders(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	header1 := map[string]string{"foo1": "bar1"}
	header2 := map[string]string{"foo1": "bar1", "foo2": "bar2"}
	header3 := map[string]string{}

	helper.Headers(header1)
	assert.Equal(t, "bar1", helper.requestOpts.MoreHeaders["foo1"])

	helper.Headers(header2)
	assert.Equal(t, "bar1", helper.requestOpts.MoreHeaders["foo1"])
	assert.Equal(t, "bar2", helper.requestOpts.MoreHeaders["foo2"])

	helper.Headers(header3)
	assert.Equal(t, "", helper.requestOpts.MoreHeaders["foo3"])
}

func TestHttpHelperOkCode(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	code1 := 200
	code2 := 404
	code3 := 502
	code4 := 10086123
	code5 := 0
	code6 := -10

	helper.OkCode(code1)
	assert.Equal(t, []int{code1}, helper.requestOpts.OkCodes)

	helper.OkCode(code2)
	assert.Equal(t, []int{code2}, helper.requestOpts.OkCodes)

	helper.OkCode(code3)
	assert.Equal(t, []int{code3}, helper.requestOpts.OkCodes)

	helper.OkCode(code4)
	assert.Equal(t, []int{code4}, helper.requestOpts.OkCodes)

	helper.OkCode(code5)
	assert.Equal(t, []int{code5}, helper.requestOpts.OkCodes)

	helper.OkCode(code6)
	assert.Equal(t, []int{code6}, helper.requestOpts.OkCodes)
}

func TestMarkerPager(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	// test empty
	dataPath1 := ""
	markerKey1 := ""
	nexExp1 := ""

	helper.MarkerPager(dataPath1, nexExp1, markerKey1)
	page := helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath1, page.(MarkerPager).DataPath)
	assert.Equal(t, markerKey1, page.(MarkerPager).MarkerKey)
	assert.Equal(t, nexExp1, page.(MarkerPager).NextExp)

	// test normal
	dataPath2 := "vpcs"
	markerKey2 := "marker"
	nexExp2 := "page_info.next_marker"

	helper.MarkerPager(dataPath2, nexExp2, markerKey2)
	page = helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath2, page.(MarkerPager).DataPath)
	assert.Equal(t, markerKey2, page.(MarkerPager).MarkerKey)
	assert.Equal(t, nexExp2, page.(MarkerPager).NextExp)
}

func TestSizePager(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	// test empty
	dataPath1 := ""
	pageNumKey1 := ""
	perPageKey1 := ""
	perPage1 := 0

	helper.PageSizePager(dataPath1, pageNumKey1, perPageKey1, perPage1)
	page := helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath1, page.(PageSizePager).DataPath)
	assert.Equal(t, pageNumKey1, page.(PageSizePager).PageNumKey)
	assert.Equal(t, perPageKey1, page.(PageSizePager).PerPageKey)

	// test normal
	dataPath2 := "vpcs"
	pageNumKey2 := "page"
	perPageKey2 := "page_info"
	perPage2 := 5

	helper.PageSizePager(dataPath2, pageNumKey2, perPageKey2, perPage2)
	page = helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath2, page.(PageSizePager).DataPath)
	assert.Equal(t, pageNumKey2, page.(PageSizePager).PageNumKey)
	assert.Equal(t, perPageKey2, page.(PageSizePager).PerPageKey)
}

func TestLinkPager(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	// test empty
	dataPath1 := ""
	linkExp1 := ""

	helper.LinkPager(dataPath1, linkExp1)
	page := helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath1, page.(LinkPager).DataPath)
	assert.Equal(t, linkExp1, page.(LinkPager).LinkExp)

	// test normal
	dataPath2 := "vpcs"
	linkExp2 := "cdssd"

	helper.LinkPager(dataPath2, linkExp2)
	page = helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath2, page.(LinkPager).DataPath)
	assert.Equal(t, linkExp2, page.(LinkPager).LinkExp)
}

func TestOffsetPager(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	// test empty
	dataPath1 := ""
	offsetKey1 := ""
	limitKey1 := ""
	defaultLimit1 := 0

	helper.OffsetPager(dataPath1, offsetKey1, limitKey1, defaultLimit1)
	page := helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath1, page.(OffsetPager).DataPath)
	assert.Equal(t, offsetKey1, page.(OffsetPager).OffsetKey)
	assert.Equal(t, limitKey1, page.(OffsetPager).LimitKey)
	assert.Equal(t, defaultLimit1, page.(OffsetPager).DefaultLimit)

	// test normal
	dataPath2 := "vpcs"
	offsetKey2 := "limit"
	limitKey2 := "limit_info"
	defaultLimit2 := 5

	helper.OffsetPager(dataPath2, offsetKey2, limitKey2, defaultLimit2)
	page = helper.pager(pagination.PageResult{})

	assert.NotNil(t, helper.pager)
	assert.Equal(t, dataPath2, page.(OffsetPager).DataPath)
	assert.Equal(t, offsetKey2, page.(OffsetPager).OffsetKey)
	assert.Equal(t, limitKey2, page.(OffsetPager).LimitKey)
	assert.Equal(t, defaultLimit2, page.(OffsetPager).DefaultLimit)
}

func TestFilter(t *testing.T) {
	// create helper
	client := &golangsdk.ServiceClient{}
	helper := New(client)

	// test empty
	filter1 := &filters.JsonFilter{}
	helper.Filter(filter1)
	assert.NotNil(t, helper.filters)
	assert.Equal(t, filter1, helper.filters[0])

	// test normal
	filter2 := filters.New()
	helper.Filter(filter2)
	assert.NotNil(t, helper.filters)
	assert.Equal(t, filter2, helper.filters[1])
}

func TestBodyToGJson(t *testing.T) {
	body0 := ""
	body1 := "foo111"
	body2 := 1234567
	body3 := []string{"aaa", "bbb"}
	body4 := map[string]string{"foo1": "bar1", "foo2": "bar2"}

	result0, err := bodyToGJson(body0)
	assert.Equal(t, body0, result0.Str)
	assert.Nil(t, err)

	result1, err := bodyToGJson(body1)
	assert.Equal(t, body1, result1.Str)
	assert.Nil(t, err)

	result2, err := bodyToGJson(body2)
	assert.Equal(t, strconv.Itoa(body2), result2.Raw)
	assert.Nil(t, err)

	result3, err := bodyToGJson(body3)
	assert.Equal(t, "[\"aaa\",\"bbb\"]\n", result3.String())
	assert.Nil(t, err)

	result4, err := bodyToGJson(body4)
	assert.Equal(t, "{\"foo1\":\"bar1\",\"foo2\":\"bar2\"}\n", result4.String())
	assert.Nil(t, err)
}
