package block_devices

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/jobs"
)

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Get is a method to obtain the detail of volume attachment for specified cloud server.
func Get(c *golangsdk.ServiceClient, server_id string, volume_id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, server_id, volume_id), &r.Body, nil)
	return
}

// List is a method to obtain the list of volume attachment details for specified cloud server.
func List(c *golangsdk.ServiceClient, serverId string) ([]VolumeAttachment, error) {
	var rst golangsdk.Result
	_, err := c.Get(rootURL(c, serverId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	if err == nil {
		var r []VolumeAttachment
		rst.ExtractIntoSlicePtr(&r, "volumeAttachments")
		return r, nil
	}
	return nil, err
}

// AttachOpts is the structure required by the Attach method to mount the disk.
type AttachOpts struct {
	// Indicates the disk device name.
	// NOTE:
	//   The new disk device name cannot be the same as an existing one.
	//   This parameter is mandatory for Xen ECSs. Set the parameter value to /dev/sda for the system disks of such ECSs
	//   and to /dev/sdx for data disks, where x is a letter in alphabetical order. For example, if there are two data
	//   disks, set the device names of the two data disks to /dev/sdb and /dev/sdc, respectively. If you set a device
	//   name starting with /dev/vd, the system uses /dev/sd by default.
	//   For KVM ECSs, set the parameter value to /dev/vda for system disks. The device names for data disks of KVM ECSs
	//   are optional. If the device names of data disks are required, set them in alphabetical order. For example, if
	//   there are two data disks, set the device names of the two data disks to /dev/vdb and /dev/vdc, respectively.
	//   If you set a device name starting with /dev/sd, the system uses /dev/vd by default.
	Device string `json:"device,omitempty"`
	// Specifies the ID of the disk to be attached. The value is in UUID format.
	VolumeId string `json:"volumeId" required:"true"`
	// Specifies the ID of the ECS to which the disk will be attached.
	ServerId string `json:"-" required:"true"`
}

// Attach is a method to mount a disk to a specified cloud server.
func Attach(c *golangsdk.ServiceClient, opts AttachOpts) (*jobs.Job, error) {
	b, err := golangsdk.BuildRequestBody(opts, "volumeAttachment")
	if err != nil {
		return nil, err
	}

	var r jobs.Job
	_, err = c.Post(attachURL(c, opts.ServerId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// DetachOpts is the structure required by the Detach method to unmount a disk from the cloud server.
type DetachOpts struct {
	// Indicates whether to forcibly detach a data disk.
	//   If yes, set it to 1.
	//   If no, set it to 0.
	// Defaults to 0.
	DeleteFlag int `q:"delete_flag"`
	// Specifies the ID of the ECS to which the disk will be detached.
	ServerId string `json:"-" required:"true"`
}

// Detach is a method to unmount a disk from a specified cloud server.
func Detach(c *golangsdk.ServiceClient, volumeId string, opts DetachOpts) (*jobs.Job, error) {
	url := detachURL(c, opts.ServerId, volumeId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r jobs.Job
	_, err = c.DeleteWithBody(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}
