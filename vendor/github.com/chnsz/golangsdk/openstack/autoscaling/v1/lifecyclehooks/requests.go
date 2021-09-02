package lifecyclehooks

import "github.com/chnsz/golangsdk"

// CreateOpts is a struct which will be used to create a lifecycle hook.
type CreateOpts struct {
	// The lifecycle hook name.
	// The name contains only letters, digits, underscores (_), and hyphens (-), and cannot exceed 32 characters.
	Name string `json:"lifecycle_hook_name" required:"true"`
	// The lifecycle hook type.
	// INSTANCE_TERMINATING: The hook suspends the instance when the instance is terminated.
	// INSTANCE_LAUNCHING: The hook suspends the instance when the instance is started.
	Type string `json:"lifecycle_hook_type" required:"true"`
	// The default lifecycle hook callback operation: ABANDON and CONTINUE.
	// By default, this operation is performed when the timeout duration expires.
	DefaultResult string `json:"default_result,omitempty"`
	// The lifecycle hook timeout duration, which ranges from 300 to 86400 in the unit of second.
	// The default value is 3600.
	Timeout int `json:"default_timeout,omitempty"`
	// The unique topic in SMN.
	// This parameter specifies a notification object for a lifecycle hook.
	NotificationTopicURN string `json:"notification_topic_urn" required:"true"`
	// The customized notification, which contains no more than 256 characters in length.
	// The message cannot contain the following characters: <>&'().
	NotificationMetadata string `json:"notification_metadata,omitempty"`
}

// CreateOptsBuilder is an interface by which can serialize the create parameters.
type CreateOptsBuilder interface {
	ToLifecycleHookCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToLifecycleHookCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method which can be able to access to create the hook of autoscaling service.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder, groupID string) (r CreateResult) {
	b, err := opts.ToLifecycleHookCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client, groupID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Get is a method to obtains the hook detail of autoscaling service.
func Get(client *golangsdk.ServiceClient, groupID, hookName string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, groupID, hookName), &r.Body, nil)
	return
}

// List is a method to obtains a hook array of the autoscaling service.
func List(client *golangsdk.ServiceClient, groupID string) (r ListResult) {
	_, r.Err = client.Get(listURL(client, groupID), &r.Body, nil)
	return
}

// UpdateOpts is a struct which will be used to update the existing hook.
type UpdateOpts struct {
	// The lifecycle hook type
	Type string `json:"lifecycle_hook_type,omitempty"`
	// The default lifecycle hook callback operation: ABANDON and CONTINUE.
	// By default, this operation is performed when the timeout duration expires.
	DefaultResult string `json:"default_result,omitempty"`
	// The lifecycle hook timeout duration, which ranges from 300 to 86400 in the unit of second.
	// The default value is 3600.
	Timeout int `json:"default_timeout,omitempty"`
	// The unique topic in SMN.
	// This parameter specifies a notification object for a lifecycle hook.
	NotificationTopicURN string `json:"notification_topic_urn,omitempty"`
	// The customized notification, which contains no more than 256 characters in length.
	// The message cannot contain the following characters: <>&'().
	NotificationMetadata string `json:"notification_metadata,omitempty"`
}

// CreateOptsBuilder is an interface by which can serialize the create parameters
type UpdateOptsBuilder interface {
	ToLifecycleHookUpdateMap() (map[string]interface{}, error)
}

func (opts UpdateOpts) ToLifecycleHookUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method which can be able to access to udpate the existing hook of the autoscaling service.
func Update(client *golangsdk.ServiceClient, opts UpdateOptsBuilder, groupID, hookName string) (r UpdateResult) {
	b, err := opts.ToLifecycleHookUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(resourceURL(client, groupID, hookName), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

//Delete is a method which can be able to access to delete the existing hook of the autoscaling service.
func Delete(client *golangsdk.ServiceClient, groupID, hookName string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, groupID, hookName), nil)
	return
}
