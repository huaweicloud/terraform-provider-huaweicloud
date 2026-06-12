package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var socMappingStatusNonUpdatableParams = []string{"workspace_id", "mapping_id"}

// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}/status
func ResourceSocMappingStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSocMappingStatusCreate,
		UpdateContext: resourceSocMappingStatusUpdate,
		ReadContext:   resourceSocMappingStatusRead,
		DeleteContext: resourceSocMappingStatusDelete,

		CustomizeDiff: config.FlexibleForceNew(socMappingStatusNonUpdatableParams),

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
			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
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

func buildSocMappingStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"status": d.Get("status"),
	}
}

func resourceSocMappingStatusCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}/status"
		workspaceId   = d.Get("workspace_id").(string)
		mappingId     = d.Get("mapping_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{mapping_id}", mappingId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		JSONBody:         buildSocMappingStatusBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating mapping status: %s", err)
	}

	d.SetId(mappingId)

	return nil
}

func resourceSocMappingStatusRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSocMappingStatusUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		updateHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}/status"
		workspaceId   = d.Get("workspace_id").(string)
		mappingId     = d.Get("mapping_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{mapping_id}", mappingId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		JSONBody:         buildSocMappingStatusBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating mapping status: %s", err)
	}

	return nil
}

func resourceSocMappingStatusDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource to update mapping status. Deleting this resource will not change
		the status of the currently mapping, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
