package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/expansion-parameters
func DataSourceGaussdbParamSetTemplateExpansion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbParamSetTemplateExpansionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"param_set_template_expansion": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussdbParamSetTemplateExpansionSchema(),
			},
		},
	}
}

func gaussdbParamSetTemplateExpansionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restart_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussdbParamSetTemplateExpansionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + "v3/{project_id}/expansion-parameters"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB parameter setting template for expansion: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	paramSetTemplateExpansion := flattenGetGaussdbParamSetTemplateExpansionBody(getRespBody)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("param_set_template_expansion", paramSetTemplateExpansion),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussdbParamSetTemplateExpansionBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("expansion_parameters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"value":            utils.PathSearch("value", v, nil),
			"restart_required": utils.PathSearch("restart_required", v, false),
			"value_range":      utils.PathSearch("value_range", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"description":      utils.PathSearch("description", v, nil),
		})
	}
	return res
}
