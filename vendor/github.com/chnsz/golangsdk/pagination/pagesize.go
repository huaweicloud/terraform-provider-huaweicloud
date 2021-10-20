package pagination

import (
	"fmt"
	"reflect"
	"strconv"

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
	if b, ok := current.Body.([]interface{}); ok {
		return len(b) == 0, nil
	}
	err := golangsdk.ErrUnexpectedType{}
	err.Expected = "[]interface{}"
	err.Actual = fmt.Sprintf("%v", reflect.TypeOf(current.Body))
	return true, err
}

// GetBody returns the linked page's body. This method is needed to satisfy the
// Page interface.
func (current PageSizeBase) GetBody() interface{} {
	return current.Body
}
