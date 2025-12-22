package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var socMappingDeleteNonUpdatableParams = []string{"workspace_id", "mapping_id"}

// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}
func ResourceSocMappingDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSocMappingDeleteCreate,
		UpdateContext: resourceSocMappingDeleteUpdate,
		ReadContext:   resourceSocMappingDeleteRead,
		DeleteContext: resourceSocMappingDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(socMappingDeleteNonUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceSocMappingDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		mappingId     = d.Get("mapping_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/mappings/{mapping_id}"
	)

	createClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{mapping_id}", mappingId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err = createClient.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting mapping: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return nil
}

func resourceSocMappingDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSocMappingDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSocMappingDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource to delete mapping. Deleting this resource will not change
	  the status of the currently mapping resource, but will only remove the resource information from the tfstate file.`

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
