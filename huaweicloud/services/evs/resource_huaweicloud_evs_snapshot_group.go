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

var snapshotGroupNonUpdatableParams = []string{
	"server_id",
	"volume_ids",
	"instant_access",
	"enterprise_project_id",
	"incremental",
}

// @API EVS POST /v5/{project_id}/snapshot-groups
// @API EVS PUT /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS GET /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS DELETE /v5/{project_id}/snapshot-groups/{snapshot_group_id}
// @API EVS POST /v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/create
// @API EVS POST /v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/delete
func ResourceEvsSnapshotGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsSnapshotGroupCreate,
		ReadContext:   resourceEvsSnapshotGroupRead,
		UpdateContext: resourceEvsSnapshotGroupUpdate,
		DeleteContext: resourceEvsSnapshotGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(snapshotGroupNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Fields `volume_ids`, `instant_access`, and `incremental` are not returned in the response body,
			// so Computed is not added.
			"volume_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instant_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Because this field can be edited to be empty, no Computed attribute is added.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func buildCreateSnapshotGroupBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	body := map[string]interface{}{
		"server_id":             utils.ValueIgnoreEmpty(d.Get("server_id")),
		"volume_ids":            utils.ValueIgnoreEmpty(d.Get("volume_ids")),
		"instant_access":        utils.ValueIgnoreEmpty(d.Get("instant_access")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		"tags":                  utils.ValueIgnoreEmpty(d.Get("tags")),
	}
	if d.Get("instant_access").(bool) {
		body["instant_access"] = true
	}
	return map[string]interface{}{"snapshot_group": body}
}

func buildUpdateSnapshotGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return map[string]interface{}{"snapshot_group": body}
}

func GetSnapshotGroupDetail(client *golangsdk.ServiceClient, snapshotGroupID string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", snapshotGroupID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceEvsSnapshotGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups"
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
		JSONBody:         utils.RemoveNil(buildCreateSnapshotGroupBodyParams(d, epsID)),
	}

	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", requestPath, &requestOpt)
		retry, err := handleMultiOperationsError(err)
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
		return diag.Errorf("error creating EVS snapshot group: %s", err)
	}

	respBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotGroupId := utils.PathSearch("snapshot_group_id", respBody, "").(string)
	if snapshotGroupId == "" {
		return diag.Errorf("error creating EVS snapshot group: ID is not found in API response")
	}

	d.SetId(snapshotGroupId)

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForSnapshotGroupJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for EVS snapshot group (%s) job success: %s", d.Id(), err)
		}
	}

	return resourceEvsSnapshotGroupRead(ctx, d, meta)
}

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}
		errorCode, errorCodeErr := jmespath.Search("error.code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		// If the response status code is 400 and the error code "EVS.2409" appears, it means that a snapshot is being
		// created under the volume. We need to wait and try again.
		if errorCode == "EVS.2409" {
			return true, err
		}
	}
	return false, err
}

func waitingForSnapshotGroupJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getSnapshotGroupJobDetail(client, jobID)
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

func getSnapshotGroupJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
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

func resourceEvsSnapshotGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	respBody, err := GetSnapshotGroupDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EVS snapshot group")
	}
	group := utils.PathSearch("snapshot_group", respBody, nil)

	mErr := multierror.Append(
		d.Set("server_id", utils.PathSearch("server_id", group, nil)),
		d.Set("name", utils.PathSearch("name", group, nil)),
		d.Set("description", utils.PathSearch("description", group, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", group, nil)),
		d.Set("tags", utils.PathSearch("tags", group, map[string]interface{}{})),
		d.Set("created_at", utils.PathSearch("created_at", group, nil)),
		d.Set("status", utils.PathSearch("status", group, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", group, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvsSnapshotGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	if d.HasChanges("name", "description") {
		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateSnapshotGroupBodyParams(d)),
		}
		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating EVS snapshot group: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		if err := updateSnapshotGroupTags(client, d.Id(), oldTags.(map[string]interface{}), false); err != nil {
			return diag.Errorf("error deleting old tags for EVS snapshot group (%s): %s", d.Id(), err)
		}
		if err := updateSnapshotGroupTags(client, d.Id(), newTags.(map[string]interface{}), true); err != nil {
			return diag.Errorf("error creating new tags for EVS snapshot group (%s): %s", d.Id(), err)
		}
	}
	return resourceEvsSnapshotGroupRead(ctx, d, meta)
}

func updateSnapshotGroupTags(client *golangsdk.ServiceClient, groupID string, tags map[string]interface{}, isCreate bool) error {
	if len(tags) == 0 {
		return nil
	}
	var (
		apiPath string
	)
	if isCreate {
		apiPath = client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/create"
	} else {
		apiPath = client.Endpoint + "v5/{project_id}/snapshot-groups/{snapshot_group_id}/tags/delete"
	}
	apiPath = strings.ReplaceAll(apiPath, "{project_id}", client.ProjectID)
	apiPath = strings.ReplaceAll(apiPath, "{snapshot_group_id}", groupID)
	requestBody := map[string]interface{}{
		"tags": expandSnapshotGroupTags(tags),
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err := client.Request("POST", apiPath, &requestOpt)
	return err
}

func resourceEvsSnapshotGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-groups/{snapshot_group_id}"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_group_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting EVS snapshot group")
	}

	if err := waitingForSnapshotGroupDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for EVS snapshot group (%s) deleted: %s", d.Id(), err)
	}

	return nil
}

func waitingForSnapshotGroupDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetSnapshotGroupDetail(client, d.Id())
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

func expandSnapshotGroupTags(rawTags map[string]interface{}) []interface{} {
	tags := make([]interface{}, 0, len(rawTags))
	for k, v := range rawTags {
		tags = append(tags, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}
	return tags
}
