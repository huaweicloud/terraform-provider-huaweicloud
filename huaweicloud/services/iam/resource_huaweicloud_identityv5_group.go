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
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v5/groups"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildV5CreateGroupBodyParams(d)),
	}
	createGroupResp, err := client.Request("POST", createPath, &createOpt)
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

func buildV5CreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name":  d.Get("group_name").(string),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func GetV5GroupById(client *golangsdk.ServiceClient, groupId string) (interface{}, error) {
	httpUrl := "v5/groups/{group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("group", respBody, nil), nil
}

func resourceV5GroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	group, err := GetV5GroupById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting IAM group")
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
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	updateChanges := []string{
		"group_name",
		"description",
	}
	if d.HasChanges(updateChanges...) {
		updateGroupHttpUrl := "v5/groups/{group_id}"
		updateGroupPath := client.Endpoint + updateGroupHttpUrl
		updateGroupPath = strings.ReplaceAll(updateGroupPath, "{group_id}", userGroupId)
		updateGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildV5UpdateGroupBodyParams(d),
		}
		_, err := client.Request("PUT", updateGroupPath, &updateGroupOpt)
		if err != nil {
			return diag.Errorf("error updating user group (%s): %s", userGroupId, err)
		}
	}

	return resourceV5GroupRead(ctx, d, meta)
}

func buildV5UpdateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"new_group_name":        d.Get("group_name").(string),
		"new_group_description": d.Get("description").(string),
	}
}

func resourceV5GroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		userGroupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteGroupHttpUrl := "v5/groups/{group_id}"
	deleteGroupPath := client.Endpoint + deleteGroupHttpUrl
	deleteGroupPath = strings.ReplaceAll(deleteGroupPath, "{group_id}", userGroupId)
	deleteGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteGroupPath, &deleteGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting user group (%s): %s", userGroupId, err)
	}

	return nil
}
