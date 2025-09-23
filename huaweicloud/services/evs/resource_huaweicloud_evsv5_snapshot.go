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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
func ResourceV5Snapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceV5SnapshotCreate,
		ReadContext:   ResourceV5SnapshotRead,
		UpdateContext: ResourceV5SnapshotUpdate,
		DeleteContext: ResourceV5SnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(evsv5SnapshotNonUpdatableParams),
			config.MergeDefaultTags(),
		),

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
			"tags": common.TagsSchema(),
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
			// The incremental default value is true, same as api doc.
			"incremental": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
				Elem:     evsv5ResourceSnapshotChainSchema(),
			},
			"snapshot_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func evsv5ResourceSnapshotChainSchema() *schema.Resource {
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
		"tags":                  utils.ValueIgnoreEmpty(d.Get("tags")),
	}
	if d.Get("instant_access").(bool) {
		body["instant_access"] = true
	}
	if !d.Get("incremental").(bool) {
		body["incremental"] = false
	}
	return map[string]interface{}{"snapshot": body}
}

func buildUpdateEvsv5SnapshotBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"name":           d.Get("name"),
		"description":    d.Get("description"),
		"instant_access": d.Get("instant_access"),
	}
	return map[string]interface{}{"snapshot": body}
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

func ResourceV5SnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForEvsv5SnapshotJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for EVS v5 snapshot (%s) job success: %s", d.Id(), err)
		}
	}

	return ResourceV5SnapshotRead(ctx, d, meta)
}

func handleEvsv5SnapshotMultiOperationsError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errResp, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errResp.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errResp.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}
		errorCode, errorCodeErr := jmespath.Search("error.code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		// If the response status code is 409 and the error code "EVS.2409" appears, it means that a snapshot is being
		// created under the volume. We need to wait and try again.
		if errorCode == "EVS.2409" {
			return true, err
		}
	}
	return false, err
}

func waitingForEvsv5SnapshotJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getEvsv5SnapshotJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in EVS job (%s) detail API response", jobID)
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			if status == "FAIL" {
				return respBody, status, fmt.Errorf("the EVS job (%s) status is FAIL, the fail reason is: %s",
					jobID, utils.PathSearch("fail_reason", respBody, "").(string))
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getEvsv5SnapshotJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v2/{project_id}/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying EVS job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func ResourceV5SnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	respBody, err := GetEvsv5SnapshotDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EVS v5 snapshot")
	}
	snapshot := utils.PathSearch("snapshot", respBody, nil)

	mErr := multierror.Append(
		d.Set("volume_id", utils.PathSearch("volume_id", snapshot, nil)),
		d.Set("name", utils.PathSearch("name", snapshot, nil)),
		d.Set("description", utils.PathSearch("description", snapshot, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", snapshot, nil)),
		d.Set("tags", utils.PathSearch("tags", snapshot, map[string]interface{}{})),
		d.Set("created_at", utils.PathSearch("created_at", snapshot, nil)),
		d.Set("instant_access", utils.PathSearch("instant_access", snapshot, nil)),
		d.Set("incremental", utils.PathSearch("incremental", snapshot, nil)),
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
		d.Set("snapshot_chains", flattenEvsv5ResourceSnapshotChains(utils.PathSearch("snapshot_chains", snapshot, []interface{}{}).([]interface{}))),
		d.Set("snapshot_group_id", utils.PathSearch("snapshot_group_id", snapshot, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceV5SnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChanges("name", "description", "instant_access") {
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
	return ResourceV5SnapshotRead(ctx, d, meta)
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
	requestBody := map[string]interface{}{
		"tags": expandEvsv5SnapshotTags(tags),
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err := client.Request("POST", apiPath, &requestOpt)
	return err
}

func ResourceV5SnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func flattenEvsv5ResourceSnapshotChains(chains []interface{}) []map[string]interface{} {
	if len(chains) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(chains))
	for _, c := range chains {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", c, nil),
			"availability_zone": utils.PathSearch("availability_zone", c, nil),
			"snapshot_count":    utils.PathSearch("snapshot_count", c, nil),
			"capacity":          utils.PathSearch("capacity", c, nil),
			"volume_id":         utils.PathSearch("volume_id", c, nil),
			"category":          utils.PathSearch("category", c, nil),
			"created_at":        utils.PathSearch("created_at", c, nil),
			"updated_at":        utils.PathSearch("updated_at", c, nil),
		})
	}
	return rst
}
