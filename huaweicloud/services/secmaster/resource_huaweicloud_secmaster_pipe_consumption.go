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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption
func ResourcePipeConsumption() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipeConsumptionCreate,
		ReadContext:   resourcePipeConsumptionRead,
		UpdateContext: resourcePipeConsumptionUpdate,
		DeleteContext: resourcePipeConsumptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePipeConsumptionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"workspace_id",
			"pipe_id",
		}),

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
			"pipe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_id": {
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

func resourcePipeConsumptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		pipeId      = d.Get("pipe_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{pipe_id}", pipeId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"pipe_id": pipeId,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster pipe consumption: %s", err)
	}

	d.SetId(pipeId)

	return resourcePipeConsumptionRead(ctx, d, meta)
}

func GetPipeConsumptionByName(client *golangsdk.ServiceClient, workspaceId, pipeId string) (interface{}, error) {
	createPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{pipe_id}", pipeId)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", createPath, &createOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", respBody, "disable").(string)
	if status == "disable" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourcePipeConsumptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetPipeConsumptionByName(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster pipe consumption")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("access_point", utils.PathSearch("access_point", respBody, nil)),
		d.Set("pipe_name", utils.PathSearch("pipe_name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("subscription_name", utils.PathSearch("subscription_name", respBody, nil)),
		d.Set("table_id", utils.PathSearch("table_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePipeConsumptionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePipeConsumptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/pipes/{pipe_id}/consumption"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{pipe_id}", d.Id())
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster pipe consumption: %s", err)
	}

	return nil
}

func resourcePipeConsumptionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<pipe_id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
		d.Set("pipe_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
