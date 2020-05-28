package tags

import (
	"github.com/huaweicloud/golangsdk"
)

//ActionOptsBuilder is an interface from which can build the request of creating/deleting group tags
type ActionOptsBuilder interface {
	ToGroupTagsActionMap() (map[string]interface{}, error)
}

//ActionOpts is a struct contains the parameters of creating/deleting group tags
type ActionOpts struct {
	Tags   []ResourceTag `json:"tags" required:"true"`
	Action string        `json:"action" required:"ture"`
}

//ToGroupTagsActionMap build the action request in json format
func (opts ActionOpts) ToGroupTagsActionMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func doAction(client *golangsdk.ServiceClient, id string, opts ActionOptsBuilder) (r ActionResult) {
	b, err := opts.ToGroupTagsActionMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(actionURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200, 204},
	})
	return
}

//Create is a method of creating group tags by id
func Create(client *golangsdk.ServiceClient, id string, tags []ResourceTag) (r ActionResult) {
	opts := ActionOpts{
		Tags:   tags,
		Action: "create",
	}
	return doAction(client, id, opts)
}

//Update is a method of updating group tags by id, used by hcs
func Update(client *golangsdk.ServiceClient, id string, tags []ResourceTag) (r ActionResult) {
	opts := ActionOpts{
		Tags:   tags,
		Action: "update",
	}
	return doAction(client, id, opts)
}

//Delete is a method of deleting group tags by id
func Delete(client *golangsdk.ServiceClient, id string, tags []ResourceTag) (r ActionResult) {
	opts := ActionOpts{
		Tags:   tags,
		Action: "delete",
	}
	return doAction(client, id, opts)
}

//Get is a method of getting the tags of the group by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

//List is a method of getting the tags of all groups
func List(client *golangsdk.ServiceClient) (r ListResult) {
	_, r.Err = client.Get(listURL(client), &r.Body, nil)
	return
}
