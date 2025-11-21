package rgc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/rgc/templates/{template_name}/deploy-params
func DataSourceTemplateDeployParams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTemplateDeployParamsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     TemplateDeployParamsSchema(),
			},
		},
	}
}

func TemplateDeployParamsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"default": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nullable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func dataSourceTemplateDeployParamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getTemplateDeployParamsProduct = "rgc"
	getTemplateDeployParamsClient, err := cfg.NewServiceClient(getTemplateDeployParamsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getTemplateDeployParamsRespBody, err := getTemplateDeployParams(getTemplateDeployParamsClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC template deploy params: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("variables", parseTemplateVariables(getTemplateDeployParamsRespBody)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseTemplateVariables(respBody interface{}) []interface{} {
	variablesList := make([]interface{}, 0)

	variables := utils.PathSearch("variables", respBody, nil)
	if variables != nil {
		variablesSlice := variables.([]interface{})
		for _, variable := range variablesSlice {
			variablesMap := variable.(map[string]interface{})
			if v, ok := variablesMap["default"]; ok {
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					variablesMap["default"] = ""
				}
				jsonStr := string(jsonBytes)
				variablesMap["default"] = jsonStr
			}

			variablesList = append(variablesList, variablesMap)
		}
	}
	return variablesList
}

func getTemplateDeployParams(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	templateName := d.Get("template_name").(string)
	var (
		httpUrl = "v1/rgc/templates/{template_name}/deploy-params"
	)
	getPath := client.Endpoint + httpUrl
	getPath += buildTemplateDeployParamsQueryParams(d)
	getPath = strings.ReplaceAll(getPath, "{template_name}", templateName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func buildTemplateDeployParamsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("template_version"); ok {
		res = fmt.Sprintf("%s&version=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
