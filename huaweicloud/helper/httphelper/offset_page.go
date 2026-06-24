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

	log.Printf("[DEBUG] [IsEmpty] [%v] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

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
	offset, err := strconv.Atoi(q.Get(p.OffsetKey))
	if err != nil {
		log.Printf("[ERROR] [NextOffset] [%v] failed to parse offset key %q, value: %q, err: %s",
			p.uuid, p.OffsetKey, q.Get(p.OffsetKey), err)
	}
	if offset == 0 {
		offset = p.DefaultOffset
	}

	// get `limit` according to the following priority:
	limit, err := strconv.Atoi(q.Get(p.LimitKey))
	if err != nil {
		log.Printf("[ERROR] [NextOffset] [%v] failed to parse limit key %q, value: %q, err: %s",
			p.uuid, p.LimitKey, q.Get(p.LimitKey), err)
	}
	if limit == 0 {
		limit = p.DefaultLimit
		count, err := p.DataCount()
		if err != nil {
			log.Printf("[ERROR] [NextOffset] [%v] failed to get data count, err: %s", p.uuid, err)
		}
		if count > 0 {
			limit = count
		}
		// save `limit` in query path
		q.Set(p.LimitKey, strconv.Itoa(limit))
	}

	log.Printf("[DEBUG] [NextOffset] [%v] offset: %v, limit: %v; OffsetKey: %s, LimitKey: %s",
		p.uuid, offset, limit, p.OffsetKey, p.LimitKey)
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

	log.Printf("[DEBUG] [NextPageURL] [%v] NextPageURL: %v", p.uuid, currentURL.String())
	return currentURL.String(), nil
}
