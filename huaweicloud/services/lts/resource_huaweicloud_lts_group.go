package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/groups
// @API LTS GET /v2/{project_id}/groups
// @API LTS POST /v2/{project_id}/groups/{log_group_id}
// @API LTS DELETE /v2/{project_id}/groups/{log_group_id}
// @API LTS POST /v1/{project_id}/{resource_type}/{resource_id}/tags/action
func ResourceLTSGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID to which the log group belongs.",
			},
			// Attributes
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/groups"
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		createPath = fmt.Sprintf("%s?enterprise_project_id=%s", createPath, epsId)
	}

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateGroupBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating log group: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	logGroupId := utils.PathSearch("log_group_id", respBody, "").(string)
	if logGroupId == "" {
		return diag.Errorf("unable to find the LTS log group ID from the API response")
	}

	d.SetId(logGroupId)

	if _, ok := d.GetOk("tags"); ok {
		groupId := d.Id()
		if err := updateTags(client, "groups", groupId, d); err != nil {
			return diag.Errorf("error creating tags of log group %s: %s", groupId, err)
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func buildCreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_group_name": d.Get("group_name"),
		"ttl_in_days":    utils.ValueIgnoreEmpty(d.Get("ttl_in_days")),
	}
	return bodyParams
}

func ignoreSysEpsTag(tags map[string]interface{}) map[string]interface{} {
	delete(tags, "_sys_enterprise_project_id")
	return tags
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/groups"
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving log group")
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error parsing the log group: %s", err)
	}

	groupResult := utils.PathSearch(fmt.Sprintf("log_groups|[?log_group_id=='%s']|[0]", groupId), respBody, nil)
	if groupResult == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, fmt.Sprintf("unable to find log group by its ID (%s)", groupId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_name", utils.PathSearch("log_group_name", groupResult, nil)),
		// Using the `delete` method in the `ignoreSysEpsTag` method will change the original value,
		// so assign a value to `enterprise_project_id` parmater before assigning a value to `tags` parmater.
		d.Set("enterprise_project_id", utils.PathSearch("tag._sys_enterprise_project_id", groupResult, "")),
		d.Set("tags", ignoreSysEpsTag(utils.PathSearch("tag", groupResult, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("ttl_in_days", utils.PathSearch("ttl_in_days", groupResult, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("creation_time", groupResult, float64(0)).(float64))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/groups/{log_group_id}"
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	if d.HasChange("ttl_in_days") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{log_group_id}", groupId)

		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         utils.RemoveNil(buildUpdateGroupBodyParams(d)),
		}

		_, err = client.Request("POST", updatePath, &updateOpts)
		if err != nil {
			return diag.Errorf("error updating log group (%s): %s", groupId, err)
		}
	}

	if d.HasChange("tags") {
		if err := updateTags(client, "groups", groupId, d); err != nil {
			return diag.Errorf("error updating tags of log group %s: %s", groupId, err)
		}
	}

	return resourceGroupRead(ctx, d, meta)
}

func buildUpdateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ttl_in_days": d.Get("ttl_in_days"),
	}
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/groups/{log_group_id}"
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{log_group_id}", groupId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting log group")
	}
	return nil
}
