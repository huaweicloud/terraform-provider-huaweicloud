package backups

import (
	"github.com/huaweicloud/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToBackupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a Backup.
type UpdateOpts struct {
	//Keep Days
	KeepDays *int `json:"keep_days" required:"true"`
	//Start Time
	StartTime string `json:"start_time" required:"true"`
	//Period
	Period string `json:"period" required:"true"`
	//DifferentialPeriod
	DifferentialPeriod string `json:"differential_period" required:"true"`
}

// ToBackupUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToBackupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "backup_policy")
}

// Update accepts a UpdateOpts struct and uses the values to update a Backup.The response code from api is 200
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToBackupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200},
		MoreHeaders: RequestOpts.MoreHeaders}
	_, r.Err = c.Put(resourceURL(c, id), b, nil, reqOpt)
	return
}
