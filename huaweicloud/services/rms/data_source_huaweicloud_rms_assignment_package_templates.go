// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RMS
// ---------------------------------------------------------------

package rms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Config GET /v1/resource-manager/conformance-packs/templates
func DataSourceTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTemplatesRead,
		Schema: map[string]*schema.Schema{
			"template_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of a built-in assignment package template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description for a built-in assignment package template.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        templatesTemplateSchema(),
				Computed:    true,
				Description: `Indicates the list of RMS assignment package templates.`,
			},
		},
	}
}

func templatesTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a built-in assignment package template.`,
			},
			"template_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a built-in assignment package template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description for a built-in assignment package template.`,
			},
			"template_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the content of a built-in assignment package template.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Elem:        templatesTemplateParameterSchema(),
				Computed:    true,
				Description: `Indicates the parameters for a built-in assignment package template.`,
			},
		},
	}
	return &sc
}

func templatesTemplateParameterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a parameter for a built-in assignment package template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of a parameter for a built-in assignment package template.`,
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the default value of a parameter for a built-in assignment package.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of a parameter for a built-in assignment package template.`,
			},
		},
	}
	return &sc
}

func resourceTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getTemplates: Query the List of RMS assignment package templates
	var (
		getTemplatesHttpUrl = "v1/resource-manager/conformance-packs/templates"
		getTemplatesProduct = "rms"
	)
	getTemplatesClient, err := cfg.NewServiceClient(getTemplatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	getTemplatesBasePath := getTemplatesClient.Endpoint + getTemplatesHttpUrl

	getTemplatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var templates []interface{}
	var marker string
	var getTemplatesPath string
	for {
		getTemplatesPath = getTemplatesBasePath + buildGetTemplatesQueryParams(d, marker)
		getTemplatesResp, err := getTemplatesClient.Request("GET", getTemplatesPath, &getTemplatesOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving RMS assignment package templates")
		}

		getTemplatesRespBody, err := utils.FlattenResponse(getTemplatesResp)
		if err != nil {
			return diag.FromErr(err)
		}
		templates = append(templates, flattenGetTemplatesResponseBodyTemplate(d, getTemplatesRespBody)...)
		marker = utils.PathSearch("page_info.next_marker", getTemplatesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("templates", templates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetTemplatesResponseBodyTemplate(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("value", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	description := d.Get("description").(string)
	for _, v := range curArray {
		descriptionRaw := utils.PathSearch("description", v, "").(string)
		if description != "" && description != descriptionRaw {
			continue
		}
		parameters := utils.PathSearch("parameters", v, nil)
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"template_key":  utils.PathSearch("template_key", v, nil),
			"description":   descriptionRaw,
			"template_body": utils.PathSearch("template_body", v, nil),
			"parameters":    flattenGetTemplatesResponseBodyTemplateParam(parameters),
		})
	}
	return rst
}

func flattenGetTemplatesResponseBodyTemplateParam(parameters interface{}) interface{} {
	if parameters == nil {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	if parametersRaw, ok := parameters.(map[string]interface{}); ok {
		// the key of the parametersRaw is the name of the parameter, the name value which in the value of parametersRaw
		// is empty, so set the key of parametersRaw to the name
		for k, v := range parametersRaw {
			parameter := map[string]interface{}{
				"name":        k,
				"description": utils.PathSearch("description", v, nil),
				"type":        utils.PathSearch("type", v, nil),
			}
			// the default value has multiple types, convert it to json and then return
			defaultValue := utils.PathSearch("default_value", v, nil)
			if defaultValue != nil {
				defaultValueBytes, _ := json.Marshal(defaultValue)
				parameter["default_value"] = string(defaultValueBytes)
			}
			rst = append(rst, parameter)
		}
	}

	return rst
}

func buildGetTemplatesQueryParams(d *schema.ResourceData, marker string) string {
	res := "?limit=200"
	if v, ok := d.GetOk("template_key"); ok {
		res = fmt.Sprintf("%s&template_key=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
