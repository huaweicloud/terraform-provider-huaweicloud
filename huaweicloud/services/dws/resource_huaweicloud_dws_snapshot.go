// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v1.0/{project_id}/snapshots
// @API DWS GET /v1.0/{project_id}/snapshots/{snapshot_id}
// @API DWS DELETE /v1.0/{project_id}/snapshots/{snapshot_id}
func ResourceDwsSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsSnapshotCreate,
		ReadContext:   resourceDwsSnapshotRead,
		DeleteContext: resourceDwsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Snapshot name.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `ID of the cluster for which you want to create a snapshot.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Snapshot description.`,
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when a snapshot starts to be created.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when a snapshot is complete.`,
			},
			"size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Snapshot size, in GB.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Snapshot status.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Snapshot type.`,
			},
		},
	}
}

func resourceDwsSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsSnapshot: create a DWS snapshot.
	var (
		createDwsSnapshotHttpUrl = "v1.0/{project_id}/snapshots"
		createDwsSnapshotProduct = "dws"
	)
	createDwsSnapshotClient, err := cfg.NewServiceClient(createDwsSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createDwsSnapshotPath := createDwsSnapshotClient.Endpoint + createDwsSnapshotHttpUrl
	createDwsSnapshotPath = strings.ReplaceAll(createDwsSnapshotPath, "{project_id}", createDwsSnapshotClient.ProjectID)

	createDwsSnapshotOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	createDwsSnapshotOpt.JSONBody = utils.RemoveNil(buildCreateDwsSnapshotBodyParams(d))
	createDwsSnapshotResp, err := createDwsSnapshotClient.Request("POST", createDwsSnapshotPath, &createDwsSnapshotOpt)
	if err != nil {
		return diag.Errorf("error creating DWS snapshot: %s", err)
	}

	createDwsSnapshotRespBody, err := utils.FlattenResponse(createDwsSnapshotResp)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotId := utils.PathSearch("snapshot.id", createDwsSnapshotRespBody, "").(string)
	if snapshotId == "" {
		return diag.Errorf("unable to find the DWS snapshot ID from the API response")
	}
	d.SetId(snapshotId)

	err = createDwsSnapshotWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of DWS Snapshot (%s) to complete: %s", d.Id(), err)
	}
	return resourceDwsSnapshotRead(ctx, d, meta)
}

func buildCreateDwsSnapshotBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"snapshot": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"cluster_id":  utils.ValueIgnoreEmpty(d.Get("cluster_id")),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func createDwsSnapshotWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createDwsSnapshotWaitingHttpUrl = "v1.0/{project_id}/snapshots/{snapshot_id}"
				createDwsSnapshotWaitingProduct = "dws"
			)
			createDwsSnapshotWaitingClient, err := cfg.NewServiceClient(createDwsSnapshotWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating DWS client: %s", err)
			}

			createDwsSnapshotWaitingPath := createDwsSnapshotWaitingClient.Endpoint + createDwsSnapshotWaitingHttpUrl
			createDwsSnapshotWaitingPath = strings.ReplaceAll(createDwsSnapshotWaitingPath, "{project_id}",
				createDwsSnapshotWaitingClient.ProjectID)
			createDwsSnapshotWaitingPath = strings.ReplaceAll(createDwsSnapshotWaitingPath, "{snapshot_id}", d.Id())

			createDwsSnapshotWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders:      requestOpts.MoreHeaders,
			}
			createDwsSnapshotWaitingResp, err := createDwsSnapshotWaitingClient.Request("GET",
				createDwsSnapshotWaitingPath, &createDwsSnapshotWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createDwsSnapshotWaitingRespBody, err := utils.FlattenResponse(createDwsSnapshotWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`snapshot.status`, createDwsSnapshotWaitingRespBody, "").(string)

			targetStatus := []string{
				"AVAILABLE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createDwsSnapshotWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"UNAVAILABLE",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createDwsSnapshotWaitingRespBody, status, nil
			}

			return createDwsSnapshotWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// GetSnapshotById is a method used to query snapshot detail.
func GetSnapshotById(client *golangsdk.ServiceClient, snapshotId string) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/snapshots/{snapshot_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{snapshot_id}", snapshotId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		// "DWS.5149": The copied snapshot does not exist.
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.5149")
	}

	return utils.FlattenResponse(resp)
}

func resourceDwsSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getDwsSnapshotClient, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getDwsSnapshotRespBody, err := GetSnapshotById(getDwsSnapshotClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS snapshot")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("snapshot.name", getDwsSnapshotRespBody, nil)),
		d.Set("cluster_id", utils.PathSearch("snapshot.cluster_id", getDwsSnapshotRespBody, nil)),
		d.Set("description", utils.PathSearch("snapshot.description", getDwsSnapshotRespBody, nil)),
		d.Set("started_at", utils.PathSearch("snapshot.started", getDwsSnapshotRespBody, nil)),
		d.Set("finished_at", utils.PathSearch("snapshot.finished", getDwsSnapshotRespBody, nil)),
		d.Set("size", utils.PathSearch("snapshot.size", getDwsSnapshotRespBody, nil)),
		d.Set("status", utils.PathSearch("snapshot.status", getDwsSnapshotRespBody, nil)),
		d.Set("type", utils.PathSearch("snapshot.type", getDwsSnapshotRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// deleteSnapshotById is a method used to delete snapshot.
func deleteSnapshotById(client *golangsdk.ServiceClient, snapshotId string) error {
	deleteHttpUrl := "v1.0/{project_id}/snapshots/{snapshot_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{snapshot_id}", snapshotId)

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// "DWS.0001": The snapshot has been deleted.
		return common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.0001")
	}

	return nil
}
func resourceDwsSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	deleteDwsSnapshotClient, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err = deleteSnapshotById(deleteDwsSnapshotClient, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DWS snapshot")
	}

	return nil
}
