package httphelper

import (
	"log"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/chnsz/golangsdk/pagination"
)

// URLFunc the URL generation function for the next query.
//
// return nil, nil // will terminate the query.
type URLFunc func(body *gjson.Result, curUrl url.URL) (*url.URL, error)

type CustomPager struct {
	pagination.PageResult
	uuid string

	NextURLFunc URLFunc
	DataPath    string
}

func (p CustomPager) IsEmpty() (bool, error) {
	rst, err := bodyToGJson(p.Body)
	if err != nil {
		return true, err
	}

	count := len(rst.Get(p.DataPath).Array())
	log.Printf("[DEBUG] [CustomPager] [%s] response count: %v, dataPath: %s", p.uuid, count, p.DataPath)

	return count == 0, nil
}

func (p CustomPager) NextPageURL() (string, error) {
	body, err := bodyToGJson(p.Body)
	if err != nil {
		return "", err
	}

	u, err := p.NextURLFunc(body, p.URL)
	if err != nil {
		return "", err
	}

	if u == nil {
		log.Printf("[ERROR] [CustomPager] [%v] next URL is empty, stop query", p.uuid)
		return "", nil
	}

	log.Printf("[DEBUG] [CustomPager] [%s] NextPageURL: %s", p.uuid, u.String())
	return u.String(), nil
}

func (p CustomPager) GetBody() any {
	return p.Body
}
