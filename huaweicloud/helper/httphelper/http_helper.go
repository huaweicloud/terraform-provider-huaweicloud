package httphelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/tidwall/gjson"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type HttpHelper struct {
	client      *golangsdk.ServiceClient
	requestOpts *golangsdk.RequestOpts
	url         string
	method      string
	body        map[string]any
	query       map[string]any
	queryExt    map[string]any

	pager func(r pagination.PageResult) pagination.Page

	result golangsdk.Result
}

func New(client *golangsdk.ServiceClient) *HttpHelper {
	return &HttpHelper{
		client:   client,
		queryExt: make(map[string]any),
		requestOpts: &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
		},
	}
}

//nolint:revive
func (c *HttpHelper) URI(url string) *HttpHelper {
	c.url = url
	return c
}

func (c *HttpHelper) Method(method string) *HttpHelper {
	c.method = method
	return c
}

func (c *HttpHelper) Body(body map[string]any) *HttpHelper {
	c.body = body
	return c
}

func (c *HttpHelper) Query(query map[string]any) *HttpHelper {
	c.query = query
	return c
}

func (c *HttpHelper) Headers(headers map[string]string) *HttpHelper {
	for key, val := range headers {
		c.requestOpts.MoreHeaders[key] = val
	}
	return c
}

func (c *HttpHelper) OkCode(okCodes ...int) *HttpHelper {
	c.requestOpts.OkCodes = okCodes
	return c
}

func (c *HttpHelper) MarkerPager(dataPath, nextExp, markerKey string) *HttpHelper {
	timestamp, _ := uuid.GenerateUUID()
	c.pager = func(r pagination.PageResult) pagination.Page {
		p := MarkerPager{
			MarkerPageBase: pagination.MarkerPageBase{PageResult: r},
			uuid:           timestamp,
			DataPath:       dataPath,
			NextExp:        nextExp,
			MarkerKey:      markerKey,
		}
		p.Owner = p
		return p
	}

	return c
}

func (c *HttpHelper) PageSizePager(dataPath, pageNumKey, perPageKey string, perPage int) *HttpHelper {
	if perPage > 0 {
		c.queryExt[perPageKey] = perPage
	}

	timestamp, _ := uuid.GenerateUUID()
	c.pager = func(r pagination.PageResult) pagination.Page {
		return PageSizePager{
			OffsetPageBase: pagination.OffsetPageBase{PageResult: r},
			uuid:           timestamp,
			DataPath:       dataPath,
			PageNumKey:     pageNumKey,
			PerPageKey:     perPageKey,
		}
	}

	return c
}

func (c *HttpHelper) LinkPager(dataPath, linkExp string) *HttpHelper {
	timestamp, _ := uuid.GenerateUUID()
	c.pager = func(r pagination.PageResult) pagination.Page {
		return LinkPager{
			LinkedPageBase: pagination.LinkedPageBase{PageResult: r},
			uuid:           timestamp,
			DataPath:       dataPath,
			LinkExp:        linkExp,
		}
	}

	return c
}

func (c *HttpHelper) OffsetPager(dataPath, offsetKey, limitKey string, defaultLimit int) *HttpHelper {
	if defaultLimit > 0 {
		c.queryExt[limitKey] = defaultLimit
	}
	timestamp, _ := uuid.GenerateUUID()

	c.pager = func(r pagination.PageResult) pagination.Page {
		return OffsetPager{
			OffsetPageBase: pagination.OffsetPageBase{PageResult: r},
			uuid:           timestamp,
			DataPath:       dataPath,
			DefaultLimit:   defaultLimit,
			OffsetKey:      offsetKey,
			LimitKey:       limitKey,
		}
	}

	return c
}

func (c *HttpHelper) Request() *HttpHelper {
	c.buildURL()
	c.appendQueryParams()

	if c.pager != nil {
		c.requestWithPage()
		return c
	}
	c.requestNoPage()

	return c
}

func (c *HttpHelper) buildURL() *HttpHelper {
	endpoint := strings.TrimRight(c.client.Endpoint, "/")
	c.url = fmt.Sprintf("%s/%s", endpoint, strings.TrimLeft(c.url, "/"))
	c.url = strings.ReplaceAll(c.url, "{project_id}", c.client.ProjectID)
	return c
}

func (c *HttpHelper) appendQueryParams() {
	query := make(map[string]any)
	for k, v := range c.query {
		query[k] = v
	}
	for k, v := range c.queryExt {
		if _, ok := query[k]; ok {
			continue
		}
		query[k] = v
	}
	if len(query) == 0 {
		return
	}

	params := marshalQueryParams(query)
	if strings.Contains(c.url, "?") {
		c.url = c.url + "&" + strings.TrimLeft(params, "?")
	} else {
		c.url += params
	}
}

func (c *HttpHelper) Result() (*gjson.Result, error) {
	if c.result.Err != nil {
		return nil, c.result.Err
	}

	result, err := BodyToGJson(c.result.Body)
	if err != nil {
		return nil, err
	}
	if result == nil || !result.Exists() {
		return nil, golangsdk.ErrDefault404{}
	}

	return result, nil
}

func (c *HttpHelper) Data() (map[string]any, error) {
	data := make(map[string]any)
	err := c.ExtractInto(&data)
	return data, err
}

func (c *HttpHelper) ExtractInto(to any) error {
	if c.result.Err != nil {
		return c.result.Err
	}

	if reader, ok := c.result.Body.(io.Reader); ok {
		if readCloser, ok := reader.(io.Closer); ok {
			defer readCloser.Close()
		}
		return json.NewDecoder(reader).Decode(to)
	}

	b, err := jsonMarshal(c.result.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)

	return err
}

func (c *HttpHelper) requestWithPage() {
	allPage, err := pagination.NewPager(c.client, c.url, c.pager).AllPages()
	if err != nil {
		c.result.Err = err
		return
	}
	c.result.Body = allPage.GetBody()
}

func (c *HttpHelper) requestNoPage() {
	var err error

	switch c.method {
	case "HEAD":
		_, err = c.client.Head(c.url, c.requestOpts)
	case "GET":
		_, err = c.client.Get(c.url, &c.result.Body, c.requestOpts)
	case "POST":
		_, err = c.client.Post(c.url, c.body, &c.result.Body, c.requestOpts)
	case "PUT":
		_, err = c.client.Put(c.url, c.body, &c.result.Body, c.requestOpts)
	case "PATCH":
		_, err = c.client.Patch(c.url, c.body, &c.result.Body, c.requestOpts)
	case "DELETE":
		_, err = c.client.DeleteWithBodyResp(c.url, c.body, &c.result.Body, c.requestOpts)
	}

	c.result.Err = err
}

func marshalQueryParams(params map[string]any) string {
	query := url.Values{}

	for key, val := range params {
		v := reflect.ValueOf(val)
		if !v.IsValid() || v.IsZero() {
			continue
		}

		switch v.Kind() {
		case reflect.String:
			query.Add(key, v.String())
		case reflect.Int:
			query.Add(key, strconv.FormatInt(v.Int(), 10))
		case reflect.Bool:
			query.Add(key, strconv.FormatBool(v.Bool()))
		case reflect.Slice:
			switch v.Type().Elem() {
			case reflect.TypeOf(0):
				for i := 0; i < v.Len(); i++ {
					query.Add(key, strconv.FormatInt(v.Index(i).Int(), 10))
				}
			default:
				for i := 0; i < v.Len(); i++ {
					if v.Index(i).IsZero() {
						continue
					}
					query.Add(key, v.Index(i).String())
				}
			}
		case reflect.Map:
			if v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.String {
				var s []string
				for _, k := range v.MapKeys() {
					value := v.MapIndex(k).String()
					s = append(s, fmt.Sprintf("'%s':'%s'", k.String(), value))
				}
				query.Add(key, fmt.Sprintf("{%s}", strings.Join(s, ", ")))
			}
		}
	}

	u := &url.URL{RawQuery: query.Encode()}
	return u.String()
}

func BodyToGJson(body any) (*gjson.Result, error) {
	b, err := BodyToBytes(body)
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(b)
	return &result, nil
}

func BodyToBytes(body any) ([]byte, error) {
	if reader, ok := body.(io.Reader); ok {
		if readCloser, ok := reader.(io.Closer); ok {
			defer readCloser.Close()
		}

		return io.ReadAll(reader)
	}

	return jsonMarshal(body)
}

func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)
	return buffer.Bytes(), err
}
