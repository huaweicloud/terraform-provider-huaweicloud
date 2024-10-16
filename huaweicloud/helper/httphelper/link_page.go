package httphelper

import (
	"log"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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

	link := utils.PathSearch(p.LinkExp, rst.Value(), "").(string)
	log.Printf("[DEBUG] [LinkPager] [%v] NextPageURL: %v, LinkExp: %s, error: %v", p.uuid, link, p.LinkExp, err)
	return link, nil
}
