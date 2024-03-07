package httphelper

import (
	"log"
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

type OffsetPager struct {
	pagination.OffsetPageBase
	uuid string

	DataPath     string
	DefaultLimit int
	OffsetKey    string
	LimitKey     string
}

// IsEmpty checks whether a PolicyPage struct is empty.
func (p OffsetPager) IsEmpty() (bool, error) {
	count, err := p.DataCount()
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] [OffsetPager] [%v] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

	return count == 0, nil
}

func (p OffsetPager) DataCount() (int, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return 0, err
	}

	return len(rst.Get(p.DataPath).Array()), nil
}

func (p OffsetPager) NextOffset() int {
	q := p.URL.Query()
	offset, _ := strconv.Atoi(q.Get(p.OffsetKey))
	if offset == 0 {
		offset = p.DefaultOffset
	}

	// get `limit` according to the following priority:
	limit, _ := strconv.Atoi(q.Get(p.LimitKey))
	if limit == 0 {
		limit = p.DefaultLimit
		count, _ := p.DataCount()
		if count > 0 {
			limit = count
		}
		// save `limit` in query path
		q.Set(p.LimitKey, strconv.Itoa(limit))
	}

	log.Printf("[DEBUG] [OffsetPager] [%v] offset: %v, limit: %v; OffsetKey: %s, LimitKey: %s", p.uuid, offset, limit, p.OffsetKey, p.LimitKey)
	return offset + limit
}

// NextPageURL generates the URL for the page of results after this one.
func (p OffsetPager) NextPageURL() (string, error) {
	next := p.NextOffset()
	if next == 0 {
		return "", nil
	}

	currentURL := p.URL
	q := currentURL.Query()
	q.Set(p.OffsetKey, strconv.Itoa(next))
	currentURL.RawQuery = q.Encode()

	log.Printf("[DEBUG] [OffsetPager] [%v] NextPageURL: %v", p.uuid, currentURL.String())
	return currentURL.String(), nil
}
