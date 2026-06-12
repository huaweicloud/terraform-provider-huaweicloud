package secmaster

import (
	"context"
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

var nodeExpansionNonUpdatableParams = []string{"workspace_id", "node_id"}

// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/nodes/{node_id}
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/nodes
func ResourceNodeExpansion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeExpansionCreate,
		ReadContext:   resourceNodeExpansionRead,
		UpdateContext: resourceNodeExpansionUpdate,
		DeleteContext: resourceNodeExpansionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNodeExpansionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nodeExpansionNonUpdatableParams),

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Fields `custom_label`, `data_center`, `description`, `maintainer`, `network_plane` can be updated to empty,
			// so Computed is not added.
			"custom_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_center": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintainer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_plane": {
				Type:     schema.TypeString,
				Optional: true,
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

// To prevent the program from making unexpected check delete judgments, it is necessary to restrict the input parameters.
func precheckExpansionParamsAllEmpty(d *schema.ResourceData) error {
	fields := []string{"custom_label", "data_center", "description", "maintainer", "network_plane"}
	for _, field := range fields {
		if d.Get(field).(string) != "" {
			return nil
		}
	}

	return fmt.Errorf("the fields %v cannot all be empty, at least one must be configured", fields)
}

func buildNodeExpansionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"custom_label":  d.Get("custom_label"),
		"data_center":   d.Get("data_center"),
		"description":   d.Get("description"),
		"maintainer":    d.Get("maintainer"),
		"network_plane": d.Get("network_plane"),
	}

	return bodyParams
}

func configNodeExpansion(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/nodes/{node_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{node_id}", d.Get("node_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildNodeExpansionBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceNodeExpansionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		nodeId = d.Get("node_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := precheckExpansionParamsAllEmpty(d); err != nil {
		return diag.FromErr(err)
	}

	if err := configNodeExpansion(client, d); err != nil {
		return diag.Errorf("error configuring SecMaster node expansion in create: %s", err)
	}

	d.SetId(nodeId)

	return resourceNodeExpansionRead(ctx, d, meta)
}

func isNodeExpansionEmpty(nodeExpansion interface{}) bool {
	if nodeExpansion == nil {
		return true
	}

	fields := []string{"custom_label", "data_center", "description", "maintainer", "network_plane"}
	for _, field := range fields {
		if utils.PathSearch(field, nodeExpansion, "").(string) != "" {
			return false
		}
	}

	return true
}

func GetNodeExpansionByNodeId(client *golangsdk.ServiceClient, workspaceId, nodeId string) (interface{}, error) {
	listPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/nodes"
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	listPath = fmt.Sprintf("%s?node_id=%s", listPath, nodeId)
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

	nodeExpansion := utils.PathSearch("records|[0].node_expansion", respBody, nil)
	if isNodeExpansionEmpty(nodeExpansion) {
		return nil, golangsdk.ErrDefault404{}
	}

	return nodeExpansion, nil
}

func resourceNodeExpansionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		nodeId      = d.Get("node_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	nodeExpansion, err := GetNodeExpansionByNodeId(client, workspaceId, nodeId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster node expansion")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("custom_label", utils.PathSearch("custom_label", nodeExpansion, nil)),
		d.Set("data_center", utils.PathSearch("data_center", nodeExpansion, nil)),
		d.Set("description", utils.PathSearch("description", nodeExpansion, nil)),
		d.Set("maintainer", utils.PathSearch("maintainer", nodeExpansion, nil)),
		d.Set("network_plane", utils.PathSearch("network_plane", nodeExpansion, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNodeExpansionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := precheckExpansionParamsAllEmpty(d); err != nil {
		return diag.FromErr(err)
	}

	if err := configNodeExpansion(client, d); err != nil {
		return diag.Errorf("error configuring SecMaster node expansion in update: %s", err)
	}

	return resourceNodeExpansionRead(ctx, d, meta)
}

func resourceNodeExpansionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		nodeId      = d.Get("node_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/nodes/{node_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{node_id}", nodeId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"custom_label":  "",
			"data_center":   "",
			"description":   "",
			"maintainer":    "",
			"network_plane": "",
		},
	}

	_, err = client.Request("PUT", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster node expansion: %s", err)
	}

	return nil
}

func resourceNodeExpansionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
		d.Set("node_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
