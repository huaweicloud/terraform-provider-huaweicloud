package dds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/configurations/comparison
func ResourceDDSParameterTemplateCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceParameterTemplateCompareCreate,
		ReadContext:   resourceParameterTemplateCompareRead,
		DeleteContext: resourceParameterTemplateCompareDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func resourceParameterTemplateCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/comparison"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateParameterTemplateCompareBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error comparing DDS parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("differences", flattenCompareResponseBodyDifferences(createRespBody)),
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

func flattenCompareResponseBodyDifferences(resp interface{}) []interface{} {
	differences := utils.PathSearch("differences", resp, make([]interface{}, 0)).([]interface{})
	if len(differences) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(differences))
	for _, difference := range differences {
		rst = append(rst, map[string]interface{}{
			"parameter_name": utils.PathSearch("parameter_name", difference, nil),
			"source_value":   utils.PathSearch("source_value", difference, nil),
			"target_value":   utils.PathSearch("target_value", difference, nil),
		})
	}
	return rst
}

func resourceParameterTemplateCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceParameterTemplateCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameter template compare resource is not supported. The resource is only removed from the" +
		"state, the DDS parameter template remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
