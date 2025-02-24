package css

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/chnsz/golangsdk/openstack/css/v1/snapshots"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CSS DELETE /v1.0/{project_id}/clusters/{clusterId}/index_snapshot/{snapId}
// @API CSS POST /v1.0/{project_id}/clusters/{clusterId}/index_snapshot
// @API CSS GET /v1.0/{project_id}/clusters/{clusterId}/index_snapshots
// @API CSS GET /v1.0/{project_id}/clusters/{clusterId}
func ResourceCssSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssSnapshotCreate,
		ReadContext:   resourceCssSnapshotRead,
		DeleteContext: resourceCssSnapshotDelete,

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

func resourceCssSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	createOpts := &snapshots.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Indices:     d.Get("index").(string),
	}

	log.Printf("[DEBUG] create options: %#v", createOpts)
	snap, err := snapshots.Create(cssClient, createOpts, clusterID).Extract()
	if err != nil {
		return diag.Errorf("error creating CSS snapshot: %s", err)
	}

	// Store the snapshot ID
	d.SetId(snap.ID)

	log.Printf("[DEBUG] waiting for snapshot (%s) to complete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    cssSnapshotStateRefreshFunc(cssClient, clusterID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for snapshot (%s) to complete: %s", d.Id(), err)
	}

	return resourceCssSnapshotRead(ctx, d, meta)
}

func resourceCssSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.CssV1Client(region)

	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	snapList, err := snapshots.List(cssClient, clusterID).Extract()
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error getting CSS cluster snapshot")
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
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	log.Printf("[DEBUG] retrieved the sanpshot %s: %+v", d.Id(), snap)

	mErr := multierror.Append(nil,
		d.Set("name", snap.Name),
		d.Set("description", snap.Description),
		d.Set("status", snap.Status),
		d.Set("index", snap.Indices),
		d.Set("cluster_id", snap.ClusterID),
		d.Set("cluster_name", snap.ClusterName),
		// Method is more suitable for backup_type
		d.Set("backup_type", snap.Method),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCssSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	if err := snapshots.Delete(cssClient, clusterID, d.Id()).ExtractErr(); err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"), "error deleting CSS cluster snapshot")
	}

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
			return nil, "NOTEXIST", fmt.Errorf("the specified snapshot %s not exist", id)
		}

		return snap, snap.Status, nil
	}
}

func resourceCssSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for CSS snapshot. Format must be <cluster id>/<snapshot id>")
		return nil, err
	}
	clusterID := parts[0]
	snapshotID := parts[1]

	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssClient, err := conf.CssV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS client: %s", err)
	}

	// check the css cluster whether exists
	if _, err := cluster.Get(cssClient, clusterID); err != nil {
		return nil, err
	}

	d.Set("cluster_id", clusterID)
	d.SetId(snapshotID)

	return []*schema.ResourceData{d}, nil
}
