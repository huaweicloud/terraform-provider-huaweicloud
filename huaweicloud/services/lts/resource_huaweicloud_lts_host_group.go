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

// @API LTS POST /v3/{project_id}/lts/host-group-list
// @API LTS DELETE /v3/{project_id}/lts/host-group
// @API LTS POST /v3/{project_id}/lts/host-group
// @API LTS PUT /v3/{project_id}/lts/host-group
func ResourceHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostGroupCreate,
		UpdateContext: resourceHostGroupUpdate,
		ReadContext:   resourceHostGroupRead,
		DeleteContext: resourceHostGroupDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the host group.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the host.`,
			},
			"host_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the ID list of hosts to join the host group.`,
			},
			"agent_access_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of the host group.`,
			},
			"labels": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the custom label list of the host group.`,
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func resourceHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createHostGroup: Create an LTS HostGroup.
	var (
		createHostGroupHttpUrl = "v3/{project_id}/lts/host-group"
		createHostGroupProduct = "lts"
	)
	createHostGroupClient, err := cfg.NewServiceClient(createHostGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	createHostGroupPath := createHostGroupClient.Endpoint + createHostGroupHttpUrl
	createHostGroupPath = strings.ReplaceAll(createHostGroupPath, "{project_id}", createHostGroupClient.ProjectID)

	createHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	createHostGroupOpt.JSONBody = utils.RemoveNil(buildCreateHostGroupBodyParams(d))
	createHostGroupResp, err := createHostGroupClient.Request("POST", createHostGroupPath, &createHostGroupOpt)
	if err != nil {
		return diag.Errorf("error creating HostGroup: %s", err)
	}

	createHostGroupRespBody, err := utils.FlattenResponse(createHostGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("host_group_id", createHostGroupRespBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable to find the LTS host group ID from the API response")
	}
	d.SetId(groupId)

	return resourceHostGroupRead(ctx, d, meta)
}

func buildCreateHostGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_group_name":   utils.ValueIgnoreEmpty(d.Get("name")),
		"host_group_type":   utils.ValueIgnoreEmpty(d.Get("type")),
		"host_id_list":      utils.ValueIgnoreEmpty(d.Get("host_ids").(*schema.Set).List()),
		"host_group_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"agent_access_type": utils.ValueIgnoreEmpty(d.Get("agent_access_type")),
		"labels":            utils.ValueIgnoreEmpty(d.Get("labels").(*schema.Set).List()),
	}
	return bodyParams
}

func resourceHostGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getHostGroup: Query the LTS HostGroup detail
	var (
		getHostGroupHttpUrl = "v3/{project_id}/lts/host-group-list"
		getHostGroupProduct = "lts"
	)
	getHostGroupClient, err := cfg.NewServiceClient(getHostGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	getHostGroupPath := getHostGroupClient.Endpoint + getHostGroupHttpUrl
	getHostGroupPath = strings.ReplaceAll(getHostGroupPath, "{project_id}", getHostGroupClient.ProjectID)

	getHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getHostGroupOpt.JSONBody = utils.RemoveNil(BuildGetOrDeleteHostGroupBodyParams(d.Id()))
	getHostGroupResp, err := getHostGroupClient.Request("POST", getHostGroupPath, &getHostGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving HostGroup")
	}

	getHostGroupRespBody, err := utils.FlattenResponse(getHostGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("result[?host_group_id=='%s']|[0]", d.Id())
	getHostGroupRespBody = utils.PathSearch(jsonPath, getHostGroupRespBody, nil)
	if getHostGroupRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("host_group_name", getHostGroupRespBody, nil)),
		d.Set("type", utils.PathSearch("host_group_type", getHostGroupRespBody, nil)),
		d.Set("host_ids", utils.PathSearch("host_id_list", getHostGroupRespBody, nil)),
		d.Set("agent_access_type", utils.PathSearch("agent_access_type", getHostGroupRespBody, nil)),
		d.Set("labels", utils.PathSearch("labels", getHostGroupRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("host_group_tag", getHostGroupRespBody, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getHostGroupRespBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("update_time", getHostGroupRespBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateHostGroupChanges := []string{
		"name",
		"host_ids",
		"tags",
		"labels",
	}

	if d.HasChanges(updateHostGroupChanges...) {
		// updateHostGroup: Update an LTS HostGroup.
		var (
			updateHostGroupHttpUrl = "v3/{project_id}/lts/host-group"
			updateHostGroupProduct = "lts"
		)
		updateHostGroupClient, err := cfg.NewServiceClient(updateHostGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS Client: %s", err)
		}

		updateHostGroupPath := updateHostGroupClient.Endpoint + updateHostGroupHttpUrl
		updateHostGroupPath = strings.ReplaceAll(updateHostGroupPath, "{project_id}", updateHostGroupClient.ProjectID)

		updateHostGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateHostGroupOpt.JSONBody = buildUpdateHostGroupBodyParams(d)
		_, err = updateHostGroupClient.Request("PUT", updateHostGroupPath, &updateHostGroupOpt)
		if err != nil {
			return diag.Errorf("error updating HostGroup: %s", err)
		}
	}
	return resourceHostGroupRead(ctx, d, meta)
}

func buildUpdateHostGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_group_id": d.Id(),
		"host_id_list":  d.Get("host_ids").(*schema.Set).List(),
		"labels":        d.Get("labels").(*schema.Set).List(),
	}

	if d.HasChange("name") {
		bodyParams["host_group_name"] = d.Get("name")
	}

	// When deleting all tags, the value received by the interface must be an empty array.
	tags := utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))
	if tags == nil {
		bodyParams["host_group_tag"] = make([]interface{}, 0)
	} else {
		bodyParams["host_group_tag"] = tags
	}

	return bodyParams
}

func resourceHostGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteHostGroup: Delete an existing LTS HostGroup
	var (
		deleteHostGroupHttpUrl = "v3/{project_id}/lts/host-group"
		deleteHostGroupProduct = "lts"
	)
	deleteHostGroupClient, err := cfg.NewServiceClient(deleteHostGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	deleteHostGroupPath := deleteHostGroupClient.Endpoint + deleteHostGroupHttpUrl
	deleteHostGroupPath = strings.ReplaceAll(deleteHostGroupPath, "{project_id}", deleteHostGroupClient.ProjectID)

	deleteHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	deleteHostGroupOpt.JSONBody = utils.RemoveNil(BuildGetOrDeleteHostGroupBodyParams(d.Id()))
	_, err = deleteHostGroupClient.Request("DELETE", deleteHostGroupPath, &deleteHostGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting HostGroup: %s", err)
	}

	return nil
}

func BuildGetOrDeleteHostGroupBodyParams(id string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_group_id_list": []string{
			id,
		},
	}
	return bodyParams
}
