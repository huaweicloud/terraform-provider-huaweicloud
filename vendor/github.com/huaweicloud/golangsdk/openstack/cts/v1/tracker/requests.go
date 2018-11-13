package tracker

import (
	"reflect"

	"github.com/huaweicloud/golangsdk"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the attributes you want to see returned.
type ListOpts struct {
	TrackerName    string `q:"tracker_name"`
	BucketName     string
	FilePrefixName string
	Status         string
}

// List returns collection of Tracker. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(client *golangsdk.ServiceClient, opts ListOpts) ([]Tracker, error) {
	var r ListResult
	_, r.Err = client.Get(rootURL(client), &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	allTracker, err := r.ExtractTracker()
	if err != nil {
		return nil, err
	}

	return FilterTracker(allTracker, opts)
}

func FilterTracker(tracker []Tracker, opts ListOpts) ([]Tracker, error) {

	var refinedTracker []Tracker
	var matched bool
	m := map[string]interface{}{}

	if opts.BucketName != "" {
		m["BucketName"] = opts.BucketName
	}
	if opts.FilePrefixName != "" {
		m["FilePrefixName"] = opts.FilePrefixName
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}

	if len(m) > 0 && len(tracker) > 0 {
		for _, trackers := range tracker {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&trackers, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedTracker = append(refinedTracker, trackers)
			}
		}
	} else {
		refinedTracker = tracker
	}
	return refinedTracker, nil
}

func getStructField(v *Tracker, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTrackerCreateMap() (map[string]interface{}, error)
}

// CreateOptsWithSMN contains the options for create a Tracker. This object is
// passed to tracker.Create().
type CreateOptsWithSMN struct {
	BucketName                string                    `json:"bucket_name" required:"true"`
	FilePrefixName            string                    `json:"file_prefix_name,omitempty"`
	SimpleMessageNotification SimpleMessageNotification `json:"smn,omitempty"`
}

// CreateOpts contains the options for create a Tracker. This object is
// passed to tracker.Create().
type CreateOpts struct {
	BucketName     string `json:"bucket_name" required:"true"`
	FilePrefixName string `json:"file_prefix_name,omitempty"`
}

type SimpleMessageNotification struct {
	IsSupportSMN          bool     `json:"is_support_smn"`
	TopicID               string   `json:"topic_id"`
	Operations            []string `json:"operations" required:"true"`
	IsSendAllKeyOperation bool     `json:"is_send_all_key_operation"`
	NeedNotifyUserList    []string `json:"need_notify_user_list,omitempty"`
}

// ToTrackerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToTrackerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// ToTrackerCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOptsWithSMN) ToTrackerCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new tracker based on the values in CreateOpts. To extract
// the tracker name  from the response, call the Extract method on the
// CreateResult.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToTrackerCreateMap()

	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// UpdateOptsWithSMN contains all the values needed to update a  tracker
type UpdateOptsWithSMN struct {
	Status                    string                    `json:"status,omitempty"`
	BucketName                string                    `json:"bucket_name" required:"true"`
	FilePrefixName            string                    `json:"file_prefix_name,omitempty"`
	SimpleMessageNotification SimpleMessageNotification `json:"smn,omitempty"`
}

// UpdateOpts contains all the values needed to update a  tracker
type UpdateOpts struct {
	Status         string `json:"status,omitempty"`
	BucketName     string `json:"bucket_name" required:"true"`
	FilePrefixName string `json:"file_prefix_name,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToTrackerUpdateMap() (map[string]interface{}, error)
}

// ToTrackerUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOptsWithSMN) ToTrackerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func (opts UpdateOpts) ToTrackerUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Update(client *golangsdk.ServiceClient, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToTrackerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular tracker.
func Delete(client *golangsdk.ServiceClient) (r DeleteResult) {
	_, r.Err = client.Delete(rootURL(client), &golangsdk.RequestOpts{
		OkCodes:  []int{204},
		JSONBody: nil,
	})
	return
}
