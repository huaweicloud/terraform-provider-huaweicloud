package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/groups
// @API IAM GET /v5/groups/{group_id}
// @API IAM PUT /v5/groups/{group_id}
// @API IAM DELETE /v5/groups/{group_id}
func ResourceV5Group() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5GroupCreate,
		ReadContext:   resourceV5GroupRead,
		UpdateContext: resourceV5GroupUpdate,
		DeleteContext: resourceV5GroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the user group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the user group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the user group.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the user group.`,
			},
		},
	}
}

func resourceV5GroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createGroupHttpUrl := "v5/groups"
	createGroupPath := iamClient.Endpoint + createGroupHttpUrl
	createGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateGroupBodyParams(d)),
	}
	createGroupResp, err := iamClient.Request("POST", createGroupPath, &createGroupOpt)
	if err != nil {
		return diag.Errorf("error creating IAM group: %s", err)
	}

	createGroupBody, err := utils.FlattenResponse(createGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	userGroupId := utils.PathSearch("group.group_id", createGroupBody, "").(string)
	if userGroupId == "" {
		return diag.Errorf("unable to find the user group ID from the API response")
	}

	d.SetId(userGroupId)
	return resourceV5GroupRead(ctx, d, meta)
}

func buildCreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name":  d.Get("group_name").(string),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceV5GroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getGroupHttpUrl := "v5/groups/{group_id}"
	getGroupPath := iamClient.Endpoint + getGroupHttpUrl
	getGroupPath = strings.ReplaceAll(getGroupPath, "{group_id}", d.Id())
	getGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getGroupResp, err := iamClient.Request("GET", getGroupPath, &getGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM group")
	}

	getGroupRespBody, err := utils.FlattenResponse(getGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	group := utils.PathSearch("group", getGroupRespBody, nil)
	if group == nil {
		return common.CheckDeletedDiag(d, err, "unable to find the user group in the API response")
	}

	mErr := multierror.Append(
		d.Set("description", utils.PathSearch("description", group, nil)),
		d.Set("group_name", utils.PathSearch("group_name", group, nil)),
		d.Set("created_at", utils.PathSearch("created_at", group, nil)),
		d.Set("urn", utils.PathSearch("urn", group, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5GroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		userGroupId = d.Id()
	)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	updateChanges := []string{
		"group_name",
		"description",
	}
	if d.HasChanges(updateChanges...) {
		updateGroupHttpUrl := "v5/groups/{group_id}"
		updateGroupPath := iamClient.Endpoint + updateGroupHttpUrl
		updateGroupPath = strings.ReplaceAll(updateGroupPath, "{group_id}", userGroupId)
		updateGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateGroupBodyParams(d),
		}
		_, err := iamClient.Request("PUT", updateGroupPath, &updateGroupOpt)
		if err != nil {
			return diag.Errorf("error updating user group (%s): %s", userGroupId, err)
		}
	}

	return resourceV5GroupRead(ctx, d, meta)
}

func buildUpdateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_group_name":        d.Get("group_name").(string),
		"new_group_description": d.Get("description").(string),
	}
	return bodyParams
}

func resourceV5GroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		userGroupId = d.Id()
	)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteGroupHttpUrl := "v5/groups/{group_id}"
	deleteGroupPath := iamClient.Endpoint + deleteGroupHttpUrl
	deleteGroupPath = strings.ReplaceAll(deleteGroupPath, "{group_id}", userGroupId)
	deleteGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = iamClient.Request("DELETE", deleteGroupPath, &deleteGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting user group (%s): %s", userGroupId, err)
	}

	return nil
}
