// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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
		return diag.Errorf("error creating DWS Client: %s", err)
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

	id, err := jmespath.Search("snapshot.id", createDwsSnapshotRespBody)
	if err != nil {
		return diag.Errorf("error creating DWS snapshot: ID is not found in API response")
	}
	d.SetId(id.(string))

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
				return nil, "ERROR", fmt.Errorf("error creating DWS Client: %s", err)
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
			statusRaw, err := jmespath.Search(`snapshot.status`, createDwsSnapshotWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `snapshot.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

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

func resourceDwsSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsSnapshot: Query the DWS snapshot.
	var (
		getDwsSnapshotHttpUrl = "v1.0/{project_id}/snapshots/{snapshot_id}"
		getDwsSnapshotProduct = "dws"
	)
	getDwsSnapshotClient, err := cfg.NewServiceClient(getDwsSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsSnapshotPath := getDwsSnapshotClient.Endpoint + getDwsSnapshotHttpUrl
	getDwsSnapshotPath = strings.ReplaceAll(getDwsSnapshotPath, "{project_id}", getDwsSnapshotClient.ProjectID)
	getDwsSnapshotPath = strings.ReplaceAll(getDwsSnapshotPath, "{snapshot_id}", d.Id())

	getDwsSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	getDwsSnapshotResp, err := getDwsSnapshotClient.Request("GET", getDwsSnapshotPath, &getDwsSnapshotOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, parseSnapshotNotFoundError(err), "error retrieving DWS snapshot")
	}

	getDwsSnapshotRespBody, err := utils.FlattenResponse(getDwsSnapshotResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
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

func parseSnapshotNotFoundError(respErr error) error {
	var apiErr interface{}
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr != nil {
			return pErr
		}
		errCode, err := jmespath.Search(`error_code`, apiErr)
		if err != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", err.Error())
		}

		if errCode == `DWS.5149` {
			return golangsdk.ErrDefault404{}
		}
	}
	return respErr
}

func resourceDwsSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsSnapshot: delete DWS snapshot
	var (
		deleteDwsSnapshotHttpUrl = "v1.0/{project_id}/snapshots/{snapshot_id}"
		deleteDwsSnapshotProduct = "dws"
	)
	deleteDwsSnapshotClient, err := cfg.NewServiceClient(deleteDwsSnapshotProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	deleteDwsSnapshotPath := deleteDwsSnapshotClient.Endpoint + deleteDwsSnapshotHttpUrl
	deleteDwsSnapshotPath = strings.ReplaceAll(deleteDwsSnapshotPath, "{project_id}", deleteDwsSnapshotClient.ProjectID)
	deleteDwsSnapshotPath = strings.ReplaceAll(deleteDwsSnapshotPath, "{snapshot_id}", d.Id())

	deleteDwsSnapshotOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	_, err = deleteDwsSnapshotClient.Request("DELETE", deleteDwsSnapshotPath, &deleteDwsSnapshotOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS snapshot: %s", err)
	}

	return nil
}
