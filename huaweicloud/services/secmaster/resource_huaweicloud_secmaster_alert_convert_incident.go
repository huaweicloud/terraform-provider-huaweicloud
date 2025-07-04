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

var alertConvertIncidentNonUpdatableParams = []string{
	"workspace_id", "ids", "title", "incident_type",
	"incident_type.*.id",
	"incident_type.*.category",
	"incident_type.*.incident_type",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/alerts/batch-order
func ResourceAlertConvertIncident() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertConvertIncidentCreate,
		ReadContext:   resourceAlertConvertIncidentRead,
		UpdateContext: resourceAlertConvertIncidentUpdate,
		DeleteContext: resourceAlertConvertIncidentDelete,

		CustomizeDiff: config.FlexibleForceNew(alertConvertIncidentNonUpdatableParams),

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
			"ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incident_type": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"incident_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"title": {
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

func buildAlertConvertIncidentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ids":              utils.ExpandToStringList(d.Get("ids").([]interface{})),
		"incident_content": buildIncidentContentBodyParams(d),
	}

	return bodyParams
}

func buildIncidentContentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"title":         utils.ValueIgnoreEmpty(d.Get("title")),
		"incident_type": buildIncidentTypeBodyParams(d.Get("incident_type").([]interface{})),
	}

	return bodyParams
}

func buildIncidentTypeBodyParams(incidentType []interface{}) map[string]interface{} {
	if len(incidentType) == 0 {
		return nil
	}

	incident, ok := incidentType[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"id":            utils.ValueIgnoreEmpty(incident["id"]),
		"category":      utils.ValueIgnoreEmpty(incident["category"]),
		"incident_type": utils.ValueIgnoreEmpty(incident["incident_type"]),
	}

	return bodyParams
}

func resourceAlertConvertIncidentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts/batch-order"
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
		JSONBody:         utils.RemoveNil(buildAlertConvertIncidentBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error converting alert into incident: %s", err)
	}

	d.SetId(workspaceId)

	return nil
}

func resourceAlertConvertIncidentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlertConvertIncidentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlertConvertIncidentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
