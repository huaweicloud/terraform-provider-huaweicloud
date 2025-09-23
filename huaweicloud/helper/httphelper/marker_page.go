package httphelper

import (
	"log"
	"net/url"
	"strings"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type MarkerPager struct {
	pagination.MarkerPageBase
	uuid string

	DataPath  string
	MarkerKey string
	NextExp   string
}

// IsEmpty returns true if a ListResult no association.
func (p MarkerPager) IsEmpty() (bool, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return false, err
	}
	count := len(rst.Get(p.DataPath).Array())

	log.Printf("[DEBUG] [MarkerPager] [%v] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

	return count == 0, nil
}

// LastMarker returns the last marker index in a ListResult.
func (p MarkerPager) LastMarker() (string, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return "", err
	}

	next := utils.PathSearch(p.NextExp, rst.Value(), "").(string)
	log.Printf("[DEBUG] [MarkerPager] [%v] last marker: %s, nextPath: %s, error: %s", p.uuid, next, p.NextExp, err)

	if !strings.Contains(next, "?") {
		return next, nil
	}

	u, err := url.Parse(next)
	if err != nil {
		return "", err
	}

	return u.Query().Get(p.MarkerKey), nil
}

func (p MarkerPager) NextPageURL() (string, error) {
	mark, err := p.Owner.LastMarker()
	log.Printf("[DEBUG] [MarkerPager] [%v] next mark: %s", p.uuid, mark)
	if err != nil {
		log.Printf("[ERROR] [MarkerPager] [%v] failed to get last marker: %s", p.uuid, err)
		return "", err
	}

	if mark == "" {
		log.Printf("[DEBUG] [MarkerPager] [%v] not found next mark, stop query", p.uuid)
		return "", nil
	}

	currentURL := p.URL
	q := currentURL.Query()
	q.Set(p.MarkerKey, mark)
	currentURL.RawQuery = q.Encode()

	log.Printf("[DEBUG] [MarkerPager] [%v] NextPageURL: %v", p.uuid, currentURL.String())
	return currentURL.String(), nil
}
