package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/evs/v2/snapshots"
)

func ResourceEvsSnapshotV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceEvsSnapshotV2Create,
		Read:   resourceEvsSnapshotV2Read,
		Update: resourceEvsSnapshotV2Update,
		Delete: resourceEvsSnapshotV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceEvsSnapshotV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	evsClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	createOpts := &snapshots.CreateOpts{
		VolumeID:    d.Get("volume_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Force:       d.Get("force").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := snapshots.Create(evsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS snapshot: %s", err)
	}

	// Wait for the snapshot to become available.
	log.Printf("[DEBUG] Waiting for volume to become available")
	err = snapshots.WaitForStatus(evsClient, v.ID, "available", int(d.Timeout(schema.TimeoutCreate)/time.Second))
	if err != nil {
		return err
	}

	// Store the ID now
	d.SetId(v.ID)
	return resourceEvsSnapshotV2Read(d, meta)
}

func resourceEvsSnapshotV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	evsClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	v, err := snapshots.Get(evsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "snapshot")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	d.Set("volume_id", v.VolumeID)
	d.Set("name", v.Name)
	d.Set("description", v.Description)
	d.Set("status", v.Status)
	d.Set("size", v.Size)

	return nil
}

func resourceEvsSnapshotV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	evsClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	updateOpts := snapshots.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	_, err = snapshots.Update(evsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud EVS snapshot: %s", err)
	}

	return resourceEvsSnapshotV2Read(d, meta)
}

func resourceEvsSnapshotV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	evsClient, err := config.BlockStorageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud EVS storage client: %s", err)
	}

	if err := snapshots.Delete(evsClient, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "snapshot")
	}

	// Wait for the snapshot to delete before moving on.
	log.Printf("[DEBUG] Waiting for snapshot (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    snapshotStateRefreshFunc(evsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// snapshotStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an HuaweiCloud snapshot.
func snapshotStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := snapshots.Get(client, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" || v.Status == "error_deleting" {
			return v, v.Status, fmt.Errorf("There was an error creating or deleting the snapshot. " +
				"Please check with your cloud admin or check the API logs " +
				"to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}
