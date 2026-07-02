package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dataClassTypeLayoutNonUpdatableParams = []string{
	"workspace_id",
	"dataclass_id",
	"type_id",
	"layout_id",
	"layout_name",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/types/layout
func ResourceDataClassTypeLayout() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataClassTypeLayoutCreate,
		ReadContext:   resourceDataClassTypeLayoutRead,
		UpdateContext: resourceDataClassTypeLayoutUpdate,
		DeleteContext: resourceDataClassTypeLayoutDelete,

		CustomizeDiff: config.FlexibleForceNew(dataClassTypeLayoutNonUpdatableParams),

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
			"dataclass_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"layout_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"layout_name": {
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

func buildDataClassTypeLayoutBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":          d.Get("type_id"),
		"layout_id":   utils.ValueIgnoreEmpty(d.Get("layout_id")),
		"layout_name": utils.ValueIgnoreEmpty(d.Get("layout_name")),
	}
}

func resourceDataClassTypeLayoutCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/dataclasses/{dataclass_id}/types/layout"
		workspaceId   = d.Get("workspace_id").(string)
		dataclassId   = d.Get("dataclass_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{dataclass_id}", dataclassId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildDataClassTypeLayoutBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error binding data class type with layout: %s", err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	return nil
}

func resourceDataClassTypeLayoutRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDataClassTypeLayoutUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDataClassTypeLayoutDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for binding data class type with layout. Deleting this
		resource will not unbind the layout from the data class type, but will only remove the resource information from
		the tfstate file.`

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
