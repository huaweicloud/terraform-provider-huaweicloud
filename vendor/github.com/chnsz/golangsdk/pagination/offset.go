package pagination

import (
	"strconv"
)

// OffsetPagebase is used for paging queries based on "limit" and "offset".
// Pagination parameters must be specified explicitly.
// You can change the parameter name of the pagination by overriding the GetPaginateParamsName method.
type OffsetPagebase struct {
	PageResult
}

// GetPaginateParamsName Overwrite this function for changing the parameter name of the pagination
// Default value is "limit", "offset"
func (current OffsetPagebase) GetPaginateParamsName() (string, string) {
	return "limit", "offset"
}

// NextPageURL generates the URL for the page of results after this one.
func (current OffsetPagebase) NextPageURL() (string, error) {
	currentURL := current.URL

	q := currentURL.Query()
	limitKey, offsetKey := current.GetPaginateParamsName()
	limit := q.Get(limitKey)
	if limit == "" {
		limit = "100"
	}

	offset := q.Get(offsetKey)
	if offset == "" {
		offset = "0"
	}

	sizeVal, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		return "", err
	}

	offsetVal, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		return "", err
	}

	offset = strconv.Itoa(int(offsetVal + sizeVal))
	q.Set(offsetKey, offset)
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

// GetBody returns the page's body. This method is needed to satisfy the
// Page interface.
func (current OffsetPagebase) GetBody() interface{} {
	return current.Body
}
