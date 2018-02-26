package flavors

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

//list the flavors informations about a specified id of database
func List(client *golangsdk.ServiceClient, dataStoreID string, region string) (r ListResult) {

	_, r.Err = client.Get(listURL(client, dataStoreID, region), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
