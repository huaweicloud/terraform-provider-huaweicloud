package flavors

import (
	"net/http"

	"github.com/huaweicloud/golangsdk"
)

//ListOptsBuilder list builder
type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

// ListOpts to list site
type ListOpts struct {
	//SiteIDS query by site ids
	SiteIDS string `q:"site_ids"`

	//Name query by name
	Name string `q:"name"`

	//Limit query limit
	Limit string `q:"limit"`

	//Offset query begin index
	Offset string `q:"offset"`

	//ID query by id
	ID string `q:"id"`

	//Area query by area
	Area string `q:"area"`

	//Province query by province
	Province string `q:"province"`

	//City query by city
	City string `q:"city"`

	//Operator query by operator
	Operator string `q:"operator"`
}

// ToListQuery converts ListOpts structures to query string
func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func List(client *golangsdk.ServiceClient, listOpts ListOptsBuilder) (r GetResult) {
	url := ListURL(client)
	if listOpts != nil {
		query, err := listOpts.ToListQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
