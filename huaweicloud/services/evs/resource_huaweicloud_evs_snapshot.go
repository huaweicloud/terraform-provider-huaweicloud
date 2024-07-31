package evs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/evs/v2/snapshots"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS DELETE /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS GET /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS PUT /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS POST /v2/{project_id}/cloudsnapshots
func ResourceEvsSnapshotV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsSnapshotV2Create,
		ReadContext:   resourceEvsSnapshotV2Read,
		UpdateContext: resourceEvsSnapshotV2Update,
		DeleteContext: resourceEvsSnapshotV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			// To avoid triggering changes metadata is not backfilled during read.
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceEvsSnapshotV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	evsClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating EVS storage client: %s", err)
	}

	createOpts := &snapshots.CreateOpts{
		VolumeID:    d.Get("volume_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Force:       d.Get("force").(bool),
		Metadata:    utils.ExpandToStringMap(d.Get("metadata").(map[string]interface{})),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := snapshots.Create(evsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating EVS snapshot: %s", err)
	}

	// Wait for the snapshot to become available.
	log.Printf("[DEBUG] Waiting for volume to become available")
	err = snapshots.WaitForStatus(evsClient, v.ID, "available", int(d.Timeout(schema.TimeoutCreate)/time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	// Store the ID now
	d.SetId(v.ID)
	return resourceEvsSnapshotV2Read(ctx, d, meta)
}

func resourceEvsSnapshotV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	evsClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating EVS storage client: %s", err)
	}

	v, err := snapshots.Get(evsClient, d.Id()).Extract()
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "snapshot")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	mErr := multierror.Append(nil,
		d.Set("volume_id", v.VolumeID),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("status", v.Status),
		d.Set("size", v.Size),
	)

	if mErr.ErrorOrNil() != nil {
		diag.FromErr(mErr)
	}
	return nil
}

func resourceEvsSnapshotV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	evsClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating EVS storage client: %s", err)
	}

	updateOpts := snapshots.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	_, err = snapshots.Update(evsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating EVS snapshot: %s", err)
	}

	return resourceEvsSnapshotV2Read(ctx, d, meta)
}

func resourceEvsSnapshotV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	evsClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating EVS storage client: %s", err)
	}

	if err := snapshots.Delete(evsClient, d.Id()).ExtractErr(); err != nil {
		// When the resource does not exist, calling the delete API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "snapshot")
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

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for snapshot (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// snapshotStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an snapshot.
func snapshotStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := snapshots.Get(client, id).Extract()
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" || v.Status == "error_deleting" {
			return v, v.Status, fmt.Errorf("there was an error creating or deleting the snapshot. " +
				"Please check with your cloud admin or check the API logs " +
				"to see why this error occurred")
		}

		return v, v.Status, nil
	}
}
