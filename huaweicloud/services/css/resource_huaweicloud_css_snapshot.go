package css

import (
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/chnsz/golangsdk/openstack/css/v1/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCssSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceCssSnapshotCreate,
		Read:   resourceCssSnapshotRead,
		Delete: resourceCssSnapshotDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCssSnapshotImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"index": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCssSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssClient, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	createOpts := &snapshots.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Indices:     d.Get("index").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	snap, err := snapshots.Create(cssClient, createOpts, clusterID).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CSS snapshot: %s", err)
	}

	// Store the snapshot ID
	d.SetId(snap.ID)

	logp.Printf("[DEBUG] Waiting for snapshot (%s) to complete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    cssSnapshotStateRefreshFunc(cssClient, clusterID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for snapshot (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceCssSnapshotRead(d, meta)
}

func resourceCssSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssClient, err := config.CssV1Client(region)

	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	snapList, err := snapshots.List(cssClient, clusterID).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "snapshot")
	}

	// find the snapshot by ID
	var snap snapshots.Snapshot
	for _, v := range snapList {
		if v.ID == d.Id() {
			snap = v
			break
		}
	}
	if snap.ID == "" {
		logp.Printf("[INFO] the snapshot %s does not exist", d.Id())
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] Retrieved the sanpshot %s: %+v", d.Id(), snap)

	d.Set("name", snap.Name)
	d.Set("description", snap.Description)
	d.Set("status", snap.Status)
	d.Set("index", snap.Indices)
	d.Set("cluster_id", snap.ClusterID)
	d.Set("cluster_name", snap.ClusterName)
	// Method is more suitable for backup_type
	d.Set("backup_type", snap.Method)

	return nil
}

func resourceCssSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssClient, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	if err := snapshots.Delete(cssClient, clusterID, d.Id()).ExtractErr(); err != nil {
		return common.CheckDeleted(d, err, "snapshot")
	}

	d.SetId("")
	return nil
}

// cssSnapshotStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an CSS cluster snapshot.
func cssSnapshotStateRefreshFunc(client *golangsdk.ServiceClient, clusterID, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapList, err := snapshots.List(client, clusterID).Extract()
		if err != nil {
			return nil, "FAILED", err
		}

		// find the snapshot by ID
		var snap snapshots.Snapshot
		for _, v := range snapList {
			if v.ID == id {
				snap = v
				break
			}
		}

		if snap.ID == "" {
			return nil, "NOTEXIST", fmtp.Errorf("The specified snapshot %s not exist", id)
		}

		return snap, snap.Status, nil
	}
}

func resourceCssSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for CSS snapshot. Format must be <cluster id>/<snapshot id>")
		return nil, err
	}
	clusterID := parts[0]
	snapshotID := parts[1]

	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssClient, err := config.CssV1Client(region)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud CSS client, err=%s", err)
	}

	// check the css cluster whether exists
	if _, err := cluster.Get(cssClient, clusterID); err != nil {
		return nil, err
	}

	d.Set("cluster_id", clusterID)
	d.SetId(snapshotID)

	return []*schema.ResourceData{d}, nil
}
