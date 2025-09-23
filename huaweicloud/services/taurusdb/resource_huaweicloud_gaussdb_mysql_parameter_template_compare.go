package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/configurations/comparison
func ResourceGaussDBMysqlTemplateCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterTemplateCompareCreate,
		ReadContext:   resourceParameterTemplateCompareRead,
		DeleteContext: resourceParameterTemplateCompareDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"differences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBTemplateCompareDifferencesSchema(),
			},
		},
	}
}

func gaussDBTemplateCompareDifferencesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"parameter_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceParameterTemplateCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/comparison"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateParameterTemplateCompareBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL parameter template compare: %s", err)
	}

	sourceConfigurationId := d.Get("source_configuration_id").(string)
	targetConfigurationId := d.Get("target_configuration_id").(string)
	d.SetId(fmt.Sprintf("%s/%s", sourceConfigurationId, targetConfigurationId))

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("differences", flattenGaussDBParameterTemplateCompareResponseBody(createRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateParameterTemplateCompareBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_configuration_id": d.Get("source_configuration_id"),
		"target_configuration_id": d.Get("target_configuration_id"),
	}
	return bodyParams
}

func flattenGaussDBParameterTemplateCompareResponseBody(resp interface{}) []interface{} {
	differencesJson := utils.PathSearch("differences", resp, make([]interface{}, 0))
	differencesArray := differencesJson.([]interface{})
	if len(differencesArray) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(differencesArray))
	for _, v := range differencesArray {
		rst = append(rst, map[string]interface{}{
			"parameter_name": utils.PathSearch("parameter_name", v, nil),
			"source_value":   utils.PathSearch("source_value", v, nil),
			"target_value":   utils.PathSearch("target_value", v, nil),
		})
	}
	return rst
}

func resourceParameterTemplateCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceParameterTemplateCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameter template compare resource is not supported. The resource is only removed from the" +
		"state, the GaussDB MySQL instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
