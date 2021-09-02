package topics

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpsBuilder is an interface which is used for creating a kafka topic
type CreateOpsBuilder interface {
	ToTopicCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters of create function
type CreateOps struct {
	// the name/ID of a topic
	Name string `json:"id" required:"true"`
	// topic partitions, value range: 1-100
	Partition int `json:"partition,omitempty"`
	// topic replications, value range: 1-3
	Replication int `json:"replication,omitempty"`
	// aging time in hours, value range: 1-168, defaults to 72
	RetentionTime int `json:"retention_time,omitempty"`

	SyncMessageFlush bool `json:"sync_message_flush,omitempty"`
	SyncReplication  bool `json:"sync_replication,omitempty"`
}

// ToTopicCreateMap is used for type convert
func (ops CreateOps) ToTopicCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a kafka topic with given parameters
func Create(client *golangsdk.ServiceClient, instanceID string, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToTopicCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// UpdateOptsBuilder is an interface which can build the map paramter of update function
type UpdateOptsBuilder interface {
	ToTopicUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the parameters of update function
type UpdateOpts struct {
	Topics []UpdateItem `json:"topics" required:"true"`
}

// UpdateItem represents the object of one topic in update function
type UpdateItem struct {
	// Name can not be updated
	Name             string `json:"id" required:"true"`
	Partition        *int   `json:"new_partition_numbers,omitempty"`
	RetentionTime    *int   `json:"retention_time,omitempty"`
	SyncMessageFlush *bool  `json:"sync_message_flush,omitempty"`
	SyncReplication  *bool  `json:"sync_replication,omitempty"`
}

// ToTopicUpdateMap is used for type convert
func (opts UpdateOpts) ToTopicUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method which can be able to update topics
func Update(client *golangsdk.ServiceClient, instanceID string, opts UpdateOptsBuilder) (r UpdateResult) {
	body, err := opts.ToTopicUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(rootURL(client, instanceID), body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Get an topic with detailed information by instance id and topic name
func Get(client *golangsdk.ServiceClient, instanceID, topicName string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, instanceID, topicName), &r.Body, nil)
	return
}

// List all topics belong to the instance id
func List(client *golangsdk.ServiceClient, instanceID string) (r ListResult) {
	_, r.Err = client.Get(rootURL(client, instanceID), &r.Body, nil)
	return
}

// Delete given topics belong to the instance id
func Delete(client *golangsdk.ServiceClient, instanceID string, topics []string) (r DeleteResult) {
	var delOpts = struct {
		Topics []string `json:"topics" required:"true"`
	}{Topics: topics}

	b, err := golangsdk.BuildRequestBody(delOpts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(deleteURL(client, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}
