package pagination

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chnsz/golangsdk"
)

// PageSizeBase is used for paging queries based on "page" and "pageSize".
// You can change the parameter name of the current page number by overriding the GetPageName method.
type PageSizeBase struct {
	PageResult
}

// GetPageName Overwrite this function for changing the parameter name of the current page number
// Default value is "page"
func (current PageSizeBase) GetPageName() string {
	return "page"
}

// NextPageURL generates the URL for the page of results after this one.
func (current PageSizeBase) NextPageURL() (string, error) {
	currentURL := current.URL

	q := currentURL.Query()
	pageNum := q.Get(current.GetPageName())
	if pageNum == "" {
		pageNum = "1"
	}

	sizeVal, err := strconv.ParseInt(pageNum, 10, 32)
	if err != nil {
		return "", err
	}

	pageNum = strconv.Itoa(int(sizeVal + 1))
	q.Set(current.GetPageName(), pageNum)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface
func (current PageSizeBase) IsEmpty() (bool, error) {
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

// GetBody returns the linked page's body. This method is needed to satisfy the
// Page interface.
func (current PageSizeBase) GetBody() interface{} {
	return current.Body
}
