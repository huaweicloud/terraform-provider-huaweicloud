package httphelper

import (
	"log"
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

type PageSizePager struct {
	pagination.OffsetPageBase
	uuid string

	DataPath   string
	PageNumKey string
	PerPageKey string
}

// IsEmpty determines whether or not a page of Roles contains any results.
func (p PageSizePager) IsEmpty() (bool, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return true, err
	}

	count := len(rst.Get(p.DataPath).Array())
	log.Printf("[DEBUG] [PageSizePager] [%v] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

	return count == 0, nil
}

// NextOffset returns offset of the next element of the page.
func (p PageSizePager) CurrentPageNum() int {
	q := p.URL.Query()
	page, _ := strconv.Atoi(q.Get(p.PageNumKey))
	if page == 0 {
		return 1
	}
	return page
}

// NextPageURL generates the URL for the page of results after this one.
func (p PageSizePager) NextPageURL() (string, error) {
	currentPageNum := p.CurrentPageNum()
	currentURL := p.URL
	q := currentURL.Query()
	q.Set(p.PageNumKey, strconv.Itoa(currentPageNum+1))
	currentURL.RawQuery = q.Encode()

	log.Printf("[DEBUG] [PageSizePager] [%v] NextPageURL: %v", p.uuid, currentURL.String())
	return currentURL.String(), nil
}
