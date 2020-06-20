package tags

import (
	"github.com/huaweicloud/golangsdk"
)

//ActionOptsBuilder is an interface from which can build the request of creating/deleting tags
type ActionOptsBuilder interface {
	ToTagsActionMap() (map[string]interface{}, error)
}

//ActionOpts is a struct contains the parameters of creating/deleting tags
type ActionOpts struct {
	Action string        `json:"action" required:"ture"`
	Tags   []ResourceTag `json:"tags" required:"true"`
}

//ToTagsActionMap build the action request in json format
func (opts ActionOpts) ToTagsActionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func doAction(client *golangsdk.ServiceClient, srvType, id string, opts ActionOptsBuilder) (r ActionResult) {
	b, err := opts.ToTagsActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, srvType, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 204},
	})
	return
}

//Create is a method of creating tags by id
func Create(client *golangsdk.ServiceClient, srvType, id string, tags []ResourceTag) (r ActionResult) {
	opts := ActionOpts{
		Tags:   tags,
		Action: "create",
	}
	return doAction(client, srvType, id, opts)
}

//Delete is a method of deleting tags by id
func Delete(client *golangsdk.ServiceClient, srvType, id string, tags []ResourceTag) (r ActionResult) {
	opts := ActionOpts{
		Tags:   tags,
		Action: "delete",
	}
	return doAction(client, srvType, id, opts)
}

//Get is a method of getting the tags by id
func Get(client *golangsdk.ServiceClient, srvType, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, srvType, id), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202, 200},
		MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
	})
	return
}

//List is a method of getting the tags of all service
func List(client *golangsdk.ServiceClient, srvType string) (r ListResult) {
	_, r.Err = client.Get(listURL(client, srvType), &r.Body, nil)
	return
}
