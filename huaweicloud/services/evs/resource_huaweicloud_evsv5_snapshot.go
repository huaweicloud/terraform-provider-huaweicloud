package evs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var evsv5SnapshotNonUpdatableParams = []string{
	"volume_id",
	"enterprise_project_id",
	"incremental",
}

// @API EVS POST /v5/{project_id}/snapshots
// @API EVS PUT /v5/{project_id}/snapshots/{snapshot_id}
// @API EVS GET /v5/{project_id}/snapshots/{snapshot_id}
// @API EVS DELETE /v5/{project_id}/snapshots/{snapshot_id}
// @API EVS POST /v5/{project_id}/snapshots/{snapshot_id}/tags/create
// @API EVS POST /v5/{project_id}/snapshots/{snapshot_id}/tags/delete
func ResourceEvsv5Snapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsv5SnapshotCreate,
		ReadContext:   resourceEvsv5SnapshotRead,
		UpdateContext: resourceEvsv5SnapshotUpdate,
		DeleteContext: resourceEvsv5SnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(evsv5SnapshotNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instant_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cmk_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instant_access_retention_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_chains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     evsv5SnapshotChainSchemaForResource(),
			},
			"snapshot_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func evsv5SnapshotChainSchemaForResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateEvsv5SnapshotBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	body := map[string]interface{}{
		"volume_id":             utils.ValueIgnoreEmpty(d.Get("volume_id")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		"instant_access":        utils.ValueIgnoreEmpty(d.Get("instant_access")),
		"incremental":           utils.ValueIgnoreEmpty(d.Get("incremental")),
	}
	if tags, ok := d.Get("tags").(map[string]interface{}); ok && len(tags) > 0 {
		tagsMap := make(map[string]string, len(tags))
		for k, v := range tags {
			if strVal, ok := v.(string); ok {
				tagsMap[k] = strVal
			}
		}
		body["tags"] = tagsMap
	}
	return map[string]interface{}{"snapshot": utils.RemoveNil(body)}
}

func buildUpdateEvsv5SnapshotBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"instant_access": utils.ValueIgnoreEmpty(d.Get("instant_access")),
	}
	return map[string]interface{}{"snapshot": utils.RemoveNil(body)}
}

func GetEvsv5SnapshotDetail(client *golangsdk.ServiceClient, snapshotID string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/snapshots/{snapshot_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", snapshotID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceEvsv5SnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshots"
		product = "evs"
		epsID   = cfg.GetEnterpriseProjectID(d)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEvsv5SnapshotBodyParams(d, epsID)),
	}

	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", requestPath, &requestOpt)
		retry, err := handleEvsv5SnapshotMultiOperationsError(err)
		return resp, retry, err
	}

	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshEvsv5SnapshotStatus(client, d.Id()),
		WaitTarget:   []string{"available"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating EVS v5 snapshot: %s", err)
	}

	respBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotId := utils.PathSearch("snapshot_id", respBody, "").(string)
	if snapshotId == "" {
		return diag.Errorf("error creating EVS v5 snapshot: ID is not found in API response")
	}

	d.SetId(snapshotId)

	return resourceEvsv5SnapshotRead(ctx, d, meta)
}

func refreshEvsv5SnapshotStatus(client *golangsdk.ServiceClient, snapshotId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetEvsv5SnapshotDetail(client, snapshotId)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("snapshot.status", respBody, "").(string)
		if status == "" {
			return respBody, "ERROR", errors.New("status is not found in API response")
		}
		return respBody, status, nil
	}
}

func handleEvsv5SnapshotMultiOperationsError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}
		errorCode, errorCodeErr := jmespath.Search("error_code||errCode", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		if errorCode == "EVS.2409" {
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		return true, err
	}
	return false, err
}

func resourceEvsv5SnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshots/{snapshot_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EVS v5 snapshot")
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	snapshot := utils.PathSearch("snapshot", respBody, nil)

	mErr := multierror.Append(
		d.Set("volume_id", utils.PathSearch("volume_id", snapshot, nil)),
		d.Set("name", utils.PathSearch("name", snapshot, nil)),
		d.Set("description", utils.PathSearch("description", snapshot, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", snapshot, nil)),
		d.Set("tags", utils.PathSearch("tags", snapshot, map[string]interface{}{})),
		d.Set("instant_access", utils.PathSearch("instant_access", snapshot, nil)),
		d.Set("incremental", utils.PathSearch("incremental", snapshot, nil)),
		d.Set("created_at", utils.PathSearch("created_at", snapshot, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", snapshot, nil)),
		d.Set("size", utils.PathSearch("size", snapshot, nil)),
		d.Set("status", utils.PathSearch("status", snapshot, nil)),
		d.Set("encrypted", utils.PathSearch("encrypted", snapshot, nil)),
		d.Set("cmk_id", utils.PathSearch("cmk_id", snapshot, nil)),
		d.Set("category", utils.PathSearch("category", snapshot, nil)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", snapshot, nil)),
		d.Set("retention_at", utils.PathSearch("retention_at", snapshot, nil)),
		d.Set("instant_access_retention_at", utils.PathSearch("instant_access_retention_at", snapshot, nil)),
		d.Set("snapshot_type", utils.PathSearch("snapshot_type", snapshot, nil)),
		d.Set("progress", utils.PathSearch("progress", snapshot, nil)),
		d.Set("encrypt_algorithm", utils.PathSearch("encrypt_algorithm", snapshot, nil)),
		d.Set("snapshot_chains", utils.PathSearch("snapshot_chains", snapshot, nil)),
		d.Set("snapshot_group_id", utils.PathSearch("snapshot_group_id", snapshot, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvsv5SnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshots/{snapshot_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateEvsv5SnapshotBodyParams(d)),
	}
	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating EVS v5 snapshot: %s", err)
	}
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		if err := updateEvsv5SnapshotTags(client, d.Id(), oldTags.(map[string]interface{}), false); err != nil {
			return diag.Errorf("error deleting old tags for EVS v5 snapshot (%s): %s", d.Id(), err)
		}
		if err := updateEvsv5SnapshotTags(client, d.Id(), newTags.(map[string]interface{}), true); err != nil {
			return diag.Errorf("error creating new tags for EVS v5 snapshot (%s): %s", d.Id(), err)
		}
	}
	return resourceEvsv5SnapshotRead(ctx, d, meta)
}

func updateEvsv5SnapshotTags(client *golangsdk.ServiceClient, snapshotID string, tags map[string]interface{}, isCreate bool) error {
	if len(tags) == 0 {
		return nil
	}
	var (
		apiPath string
	)
	if isCreate {
		apiPath = client.Endpoint + "v5/{project_id}/snapshots/{snapshot_id}/tags/create"
	} else {
		apiPath = client.Endpoint + "v5/{project_id}/snapshots/{snapshot_id}/tags/delete"
	}
	apiPath = strings.ReplaceAll(apiPath, "{project_id}", client.ProjectID)
	apiPath = strings.ReplaceAll(apiPath, "{snapshot_id}", snapshotID)
	body := map[string]interface{}{}
	body["tags"] = expandEvsv5SnapshotTags(tags)
	requestBody := body
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err := client.Request("POST", apiPath, &requestOpt)
	return err
}

func resourceEvsv5SnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshots/{snapshot_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting EVS v5 snapshot")
	}

	if err := waitingForEvsv5SnapshotDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for EVS v5 snapshot (%s) deleted: %s", d.Id(), err)
	}

	return nil
}

func waitingForEvsv5SnapshotDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetEvsv5SnapshotDetail(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			status := utils.PathSearch("snapshot.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", errors.New("status is not found in API response")
			}

			if status == "error_deleting" {
				return respBody, status, errors.New("an error occurred while deleting the EVS v5 snapshot. " +
					"Please check with your cloud admin or check the API logs " +
					"to see why this error occurred")
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func expandEvsv5SnapshotTags(rawTags map[string]interface{}) []interface{} {
	tags := make([]interface{}, 0, len(rawTags))
	for k, v := range rawTags {
		tags = append(tags, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}
	return tags
}
