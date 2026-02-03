package iam

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

// @API IAM POST /v3/groups
// @API IAM GET /v3/groups/{group_id}
// @API IAM PATCH /v3/groups/{group_id}
// @API IAM DELETE /v3/groups/{group_id}
func ResourceV3Group() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3GroupCreate,
		ReadContext:   resourceV3GroupRead,
		UpdateContext: resourceV3GroupUpdate,
		DeleteContext: resourceV3GroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group.",
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the group.",
			},
		},
	}
}

func buildV3CreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
}

func resourceV3GroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/groups"
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"group": buildV3CreateGroupBodyParams(d),
		}),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating group: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("group.id", respBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable to find the IAM group ID from the API response")
	}
	d.SetId(groupId)

	return resourceV3GroupRead(ctx, d, meta)
}

func GetV3GroupById(client *golangsdk.ServiceClient, groupId string) (interface{}, error) {
	httpUrl := "v3/groups/{group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceV3GroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	group, err := GetV3GroupById(client, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving group (%s)", groupId))
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("group.name", group, nil)),
		d.Set("description", utils.PathSearch("group.description", group, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting group fields: %s", err)
	}
	return nil
}

func buildV3UpdateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
}

func resourceV3GroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/groups/{group_id}"
		groupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", groupId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"group": buildV3UpdateGroupBodyParams(d),
		}),
	}

	_, err = client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating group (%s): %s", groupId, err)
	}

	return resourceV3GroupRead(ctx, d, meta)
}

func resourceV3GroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/groups/{group_id}"
		groupId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", groupId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting group (%s): %s", groupId, err)
	}
	return nil
}
