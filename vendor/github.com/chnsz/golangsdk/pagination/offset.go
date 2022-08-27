package pagination

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chnsz/golangsdk"
)

// OffsetPage is used for paging queries based on offset and limit.
type OffsetPage interface {
	Page

	NextOffset() int
}

// OffsetPageBase is used for paging queries based on offset and limit.
// DefaultOffset and DefaultLimit are used to calculate the next offset
type OffsetPageBase struct {
	PageResult

	DefaultOffset int
	DefaultLimit  int
}

// NextOffset returns offset of the next element of the page.
func (current OffsetPageBase) NextOffset() int {
	q := current.URL.Query()
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset == 0 {
		offset = current.DefaultOffset
	}

	// get `limit` according to the following priority:
	// query path > current page size > DefaultLimit
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit == 0 {
		count, err := current.getPageSize()
		if err != nil {
			limit = current.DefaultLimit
		}
		limit = count
		// save `limit` in query path
		q.Set("limit", strconv.Itoa(limit))
	}

	return offset + limit
}

// NextPageURL generates the URL for the page of results after this one.
func (current OffsetPageBase) NextPageURL() (string, error) {
	next := current.NextOffset()
	if next == 0 {
		return "", nil
	}

	currentURL := current.URL
	q := currentURL.Query()
	q.Set("offset", strconv.Itoa(next))
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface.
func (current OffsetPageBase) IsEmpty() (bool, error) {
	if pb, ok := current.Body.(map[string]interface{}); ok {
		for k, v := range pb {
			// ignore xxx_links
			if !strings.HasSuffix(k, "links") {
				// check the field's type. we only want []interface{} (which is really []map[string]interface{})
				switch vt := v.(type) {
				case []interface{}:
					return len(vt) == 0, nil
				}
			}
		}
	}
	if pb, ok := current.Body.([]interface{}); ok {
		return len(pb) == 0, nil
	}

	err := golangsdk.ErrUnexpectedType{}
	err.Expected = "map[string]interface{}/[]interface{}"
	err.Actual = fmt.Sprintf("%T", current.Body)
	return true, err
}

// GetBody returns the page's body. This method is needed to satisfy the Page interface.
func (current OffsetPageBase) GetBody() interface{} {
	return current.Body
}

// getPageSize calculates the count of items in an offset page
func (current OffsetPageBase) getPageSize() (int, error) {
	var pageSize int

	switch pb := current.Body.(type) {
	case map[string]interface{}:
		for k, v := range pb {
			// ignore xxx_links
			if !strings.HasSuffix(k, "links") {
				// check the field's type. we only want []interface{} (which is really []map[string]interface{})
				switch vt := v.(type) {
				case []interface{}:
					pageSize = len(vt)
					break
				}
			}
		}
	case []interface{}:
		pageSize = len(pb)
	default:
		err := golangsdk.ErrUnexpectedType{}
		err.Expected = "map[string]interface{}/[]interface{}"
		err.Actual = fmt.Sprintf("%T", pb)
		return 0, err
	}

	return pageSize, nil
}
