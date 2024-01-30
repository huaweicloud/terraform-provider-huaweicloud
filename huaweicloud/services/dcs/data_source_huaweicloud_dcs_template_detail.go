// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DCS
// ---------------------------------------------------------------

package dcs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/config-templates/{template_id}
func DataSourceTemplateDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTemplateDetailRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the template.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the template.`,
			},
			"params": {
				Type:        schema.TypeList,
				Elem:        templateDetailParamSchema(),
				Optional:    true,
				Description: `Specifies the ID of the template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the template.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine version.`,
			},
			"cache_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DCS instance type.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the product edition.`,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the template.`,
			},
		},
	}
}

func templateDetailParamSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the param.`,
			},
			"param_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the param.`,
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the default of the param.`,
			},
			"value_range": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the value range of the param.`,
			},
			"value_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the value type of the param.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the param.`,
			},
			"need_restart": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the DCS instance need restart.`,
			},
		},
	}
	return &sc
}

func resourceTemplateDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDcsTemplateDetail: Query the DCS template detail.
	var (
		getDcsTemplateDetailHttpUrl = "v2/{project_id}/config-templates/{template_id}"
		getDcsTemplateDetailProduct = "dcs"
	)
	getDcsTemplateDetailClient, err := cfg.NewServiceClient(getDcsTemplateDetailProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getDcsTemplateDetailPath := getDcsTemplateDetailClient.Endpoint + getDcsTemplateDetailHttpUrl
	getDcsTemplateDetailPath = strings.ReplaceAll(getDcsTemplateDetailPath, "{project_id}",
		getDcsTemplateDetailClient.ProjectID)
	getDcsTemplateDetailPath = strings.ReplaceAll(getDcsTemplateDetailPath, "{template_id}",
		d.Get("template_id").(string))

	getDcsTemplateDetailQueryParams := buildGetDcsTemplateDetailQueryParams(d)
	getDcsTemplateDetailPath += getDcsTemplateDetailQueryParams

	getDcsTemplateDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDcsTemplateDetailResp, err := getDcsTemplateDetailClient.Request("GET", getDcsTemplateDetailPath,
		&getDcsTemplateDetailOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DCS template detail")
	}

	getDcsTemplateDetailRespBody, err := utils.FlattenResponse(getDcsTemplateDetailResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDcsTemplateDetailRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getDcsTemplateDetailRespBody, nil)),
		d.Set("engine", utils.PathSearch("engine", getDcsTemplateDetailRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getDcsTemplateDetailRespBody, nil)),
		d.Set("cache_mode", utils.PathSearch("cache_mode", getDcsTemplateDetailRespBody, nil)),
		d.Set("product_type", utils.PathSearch("product_type", getDcsTemplateDetailRespBody, nil)),
		d.Set("storage_type", utils.PathSearch("storage_type", getDcsTemplateDetailRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getDcsTemplateDetailRespBody, nil)),
		d.Set("params", filterGetDcsTemplateDetailResponseBodyParam(
			flattenGetDcsTemplateDetailResponseBodyParam(getDcsTemplateDetailRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDcsTemplateDetailResponseBodyParam(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("params", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"param_id":      utils.PathSearch("param_id", v, nil),
			"param_name":    utils.PathSearch("param_name", v, nil),
			"default_value": utils.PathSearch("default_value", v, nil),
			"value_range":   utils.PathSearch("value_range", v, nil),
			"value_type":    utils.PathSearch("value_type", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"need_restart":  utils.PathSearch("need_restart", v, nil),
		})
	}
	return rst
}

func filterGetDcsTemplateDetailResponseBodyParam(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	paramsMap := buildRawParamsMap(d)
	log.Printf("paramsMap info: %#v", paramsMap)
	for _, v := range all {
		paramName := utils.PathSearch("param_name", v, "").(string)
		if len(paramsMap) > 0 && !paramsMap[paramName] {
			continue
		}
		rst = append(rst, v)
	}
	return rst
}

func buildRawParamsMap(d *schema.ResourceData) map[string]bool {
	params := d.Get("params").([]interface{})
	paramsMap := make(map[string]bool)
	for _, param := range params {
		if v, ok := param.(map[string]interface{}); ok && v["param_name"].(string) != "" {
			paramsMap[v["param_name"].(string)] = true
		}
	}
	return paramsMap
}

func buildGetDcsTemplateDetailQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
