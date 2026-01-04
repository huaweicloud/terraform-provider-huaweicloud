package secmaster

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var collectorChannelGroupNonUpdatableParams = []string{"workspace_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups/{group_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups/{group_id}
func ResourceCollectorChannelGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorChannelGroupCreate,
		ReadContext:   resourceCollectorChannelGroupRead,
		UpdateContext: resourceCollectorChannelGroupUpdate,
		DeleteContext: resourceCollectorChannelGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCollectorChannelGroupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(collectorChannelGroupNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func buildCollectorChannelGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}

	return bodyParams
}

func resourceCollectorChannelGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		groupName   = d.Get("name").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCollectorChannelGroupBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating colloector channel group: %s", err)
	}

	group, err := GetCollectorChannelGroupByName(client, workspaceId, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("group_id", group, "").(string)
	if groupId == "" {
		return diag.Errorf("error creating collector channel group: unable to find group ID")
	}

	d.SetId(groupId)

	return resourceCollectorChannelGroupRead(ctx, d, meta)
}

func GetCollectorChannelGroupByName(client *golangsdk.ServiceClient, workspaceId, groupName string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	listPath = fmt.Sprintf("%s?name=%s", listPath, groupName)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if groups, ok := respBody.([]interface{}); ok {
		if len(groups) > 0 {
			return groups[0], nil
		}
	}

	// When the collector channel group does not exist, the status code is `200`.
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the collector channel group (%s) does not exist", groupName)),
		},
	}
}

func resourceCollectorChannelGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		groupName   = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	group, err := GetCollectorChannelGroupByName(client, workspaceId, groupName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving collector channel group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", group, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCollectorChannelGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups/{group_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCollectorChannelGroupBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating collector channel group: %s", err)
	}

	return resourceCollectorChannelGroupRead(ctx, d, meta)
}

func resourceCollectorChannelGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/groups/{group_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the collector channel group does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting the colloector channel group")
	}

	return nil
}

func resourceCollectorChannelGroupImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
		d.Set("name", parts[1]),
	)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	group, err := GetCollectorChannelGroupByName(client, parts[0], parts[1])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving collector channel group from API response: %s", err)
	}

	groupId := utils.PathSearch("group_id", group, "").(string)
	if groupId == "" {
		return []*schema.ResourceData{d}, errors.New("unable to find collector channel group ID from API response")
	}

	d.SetId(groupId)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
