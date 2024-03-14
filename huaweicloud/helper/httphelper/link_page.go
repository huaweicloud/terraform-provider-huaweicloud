package httphelper

import (
	"log"

	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk/pagination"
)

type LinkPager struct {
	pagination.LinkedPageBase
	uuid string

	DataPath string
	LinkExp  string
}

// IsEmpty returns true if an ImagePage contains no Image results.
func (p LinkPager) IsEmpty() (bool, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return false, err
	}

	count := len(rst.Get(p.DataPath).Array())

	log.Printf("[DEBUG] [LinkPager] [%v] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

	return count == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (p LinkPager) NextPageURL() (string, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return "", err
	}

	val, err := jmespath.Search(p.LinkExp, rst.Value())
	log.Printf("[DEBUG] [LinkPager] [%v] NextPageURL: %v, LinkExp: %s, error: %v", p.uuid, val, p.LinkExp, err)
	if err != nil {
		return "", err
	}

	link, _ := val.(string)
	return link, nil
}
