package httphelper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chnsz/golangsdk"
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
