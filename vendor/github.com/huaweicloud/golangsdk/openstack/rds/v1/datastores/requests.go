package datastores

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

//list the version informations about a specified type of database
func List(client *golangsdk.ServiceClient, dataStoreName string) (r ListResult) {

	_, r.Err = client.Get(listURL(client, dataStoreName), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}
