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

var cloudLogResourceNonUpdatableParams = []string{
	"workspace_id",
	"domain_id",
	"resources",
	"resources.*.enable",
	"resources.*.region_id",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/resource
func ResourceCloudLogResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudLogResourceCreate,
		UpdateContext: resourceCloudLogResourceUpdate,
		ReadContext:   resourceCloudLogResourceRead,
		DeleteContext: resourceCloudLogResourceDelete,

		CustomizeDiff: config.FlexibleForceNew(cloudLogResourceNonUpdatableParams),

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
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func buildCloudLogResourceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_id": d.Get("domain_id"),
		"resources": buildResourceDataBodyParams(d.Get("resources").([]interface{})),
	}

	return bodyParams
}

func buildResourceDataBodyParams(resourceInfo []interface{}) []map[string]interface{} {
	if len(resourceInfo) == 0 {
		return nil
	}

	paramsInfo := make([]map[string]interface{}, 0, len(resourceInfo))
	for _, v := range resourceInfo {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"enable":    raw["enable"],
			"region_id": raw["region_id"],
		}

		paramsInfo = append(paramsInfo, params)
	}

	return paramsInfo
}

func resourceCloudLogResourceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/cloud-logs/resource"
		workspaceId   = d.Get("workspace_id").(string)
	)

	createClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCloudLogResourceBodyParams(d),
	}

	_, err = createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating cloud log resource: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return nil
}

func resourceCloudLogResourceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCloudLogResourceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCloudLogResourceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource to create cloud log resource. Deleting this resource will
		not change the status of the currently resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
