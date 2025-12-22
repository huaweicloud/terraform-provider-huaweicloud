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

var dataSpaceNonUpdatableParams = []string{"workspace_id", "dataspace_name"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}
func ResourceDataspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataspaceCreate,
		ReadContext:   resourceDataspaceRead,
		UpdateContext: resourceDataspaceUpdate,
		DeleteContext: resourceDataspaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataspaceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(dataSpaceNonUpdatableParams),

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
			"dataspace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"dataspace_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildDataspaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dataspace_name": d.Get("dataspace_name"),
		"description":    d.Get("description"),
	}

	return bodyParams
}

func resourceDataspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces"
		workspaceId   = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDataspaceBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating dataspace: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataspaceId := utils.PathSearch("dataspace_id", respBody, "").(string)
	if dataspaceId == "" {
		return diag.Errorf("unable to find the dataspace ID from the API response")
	}

	d.SetId(dataspaceId)

	return resourceDataspaceRead(ctx, d, meta)
}

func resourceDataspaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	dataspace, err := GetDataspaceInfo(client, workspaceId, d.Id())
	if err != nil {
		// When the dataspace does not exist, the response HTTP status code of the query API is 404
		return common.CheckDeletedDiag(d, err, "error retrieving dataspace")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("dataspace_name", utils.PathSearch("dataspace_name", dataspace, nil)),
		d.Set("description", utils.PathSearch("description", dataspace, nil)),
		d.Set("dataspace_type", utils.PathSearch("dataspace_type", dataspace, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", dataspace, nil)),
		d.Set("project_id", utils.PathSearch("project_id", dataspace, nil)),
		d.Set("create_by", utils.PathSearch("create_by", dataspace, nil)),
		d.Set("update_by", utils.PathSearch("update_by", dataspace, nil)),
		d.Set("create_time", utils.PathSearch("create_time", dataspace, nil)),
		d.Set("update_time", utils.PathSearch("update_time", dataspace, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDataspaceInfo(client *golangsdk.ServiceClient, workspaceId, dataspaceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{dataspace_id}", dataspaceId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateDataspaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
	}

	return bodyParams
}

func resourceDataspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{dataspace_id}", d.Id())

	if d.HasChange("description") {
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateDataspaceBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating dataspace: %s", err)
		}
	}

	return resourceDataspaceRead(ctx, d, meta)
}

func resourceDataspaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces/{dataspace_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{dataspace_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the dataspace does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting dataspace")
	}

	return nil
}

func resourceDataspaceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
