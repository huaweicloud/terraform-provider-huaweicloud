package pagination

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/chnsz/golangsdk"
)

type OffsetPage interface {
	LastOffset() int
}

// OffsetPageBase is used for paging queries based on offset and limit.
type OffsetPageBase struct {
	Offset int
	Limit  int

	PageResult
}

// LastOffset returns index of the last element of the page.
func (current OffsetPageBase) LastOffset() int {
	q := current.URL.Query()
	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		offset = current.Offset
		q.Set("offset", strconv.Itoa(offset))
	}
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = current.Limit
		q.Set("limit", strconv.Itoa(limit))
	}
	return offset + limit
}

// NextPageURL generates the URL for the page of results after this one.
func (current OffsetPageBase) NextPageURL() (string, error) {
	currentURL := current.URL
	q := currentURL.Query()

	if q.Get("offset") == "" && q.Get("limit") == "" {
		// Without offset and limit, the page just a SinglePageBase.
		return "", nil
	}

	q.Set("offset", strconv.Itoa(current.LastOffset()))
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface.
func (current OffsetPageBase) IsEmpty() (bool, error) {
	if b, ok := current.Body.([]interface{}); ok {
		return len(b) == 0, nil
	}
	err := golangsdk.ErrUnexpectedType{}
	err.Expected = "[]interface{}"
	err.Actual = fmt.Sprintf("%v", reflect.TypeOf(current.Body))
	return true, err
}

// GetBody returns the page's body. This method is needed to satisfy the Page interface.
func (current OffsetPageBase) GetBody() interface{} {
	return current.Body
}
