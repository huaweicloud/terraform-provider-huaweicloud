package coc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var incidentHandleNonUpdatableParams = []string{"incident_num", "operator", "operate_key", "parameter"}

// @API COC POST /v1/external/incident/handle
func ResourceIncidentHandle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentHandleCreate,
		ReadContext:   resourceIncidentHandleRead,
		UpdateContext: resourceIncidentHandleUpdate,
		DeleteContext: resourceIncidentHandleDelete,

		CustomizeDiff: config.FlexibleForceNew(incidentHandleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"incident_num": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operate_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameter": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildIncidentHandleCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"incident_num": d.Get("incident_num"),
		"operator":     d.Get("operator"),
		"operate_key":  d.Get("operate_key"),
		"parameter":    buildIncidentHandleCreateOptsParamsBody(d.Get("parameter")),
	}

	return bodyParams
}

func buildIncidentHandleCreateOptsParamsBody(rawParams interface{}) map[string]interface{} {
	// have to send empty structure
	if rawParams == nil {
		return make(map[string]interface{}, 0)
	}

	return utils.RemoveNil(rawParams.(map[string]interface{}))
}

func resourceIncidentHandleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/external/incident/handle"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildIncidentHandleCreateOpts(d),
	}

	incidentNum := d.Get("incident_num").(string)

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC incident handle: %s", err)
	}

	d.SetId(incidentNum)

	return resourceIncidentHandleRead(ctx, d, meta)
}

func resourceIncidentHandleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIncidentHandleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIncidentHandleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting incident handle resource is not supported. The incident handle resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
