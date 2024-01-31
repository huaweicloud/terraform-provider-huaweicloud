// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/models/template
func DataSourceModelTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceModelTemplatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of model. The valid values are **Classification** and **Common**.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AI engine.`,
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model runtime environment.`,
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Keywords to search in name or description. Fuzzy match is supported.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        modelTemplateTemplatesSchema(),
				Computed:    true,
				Description: `The list of model templates.`,
			},
		},
	}
}

func modelTemplateTemplatesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Template name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Template description.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Architecture type. The valid values are **X86_64** and **AARCH64**.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of model. The valid values are **Classification** and **Common**.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The AI engine.`,
			},
			"environment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model runtime environment.`,
			},
			"template_docs": {
				Type:     schema.TypeList,
				Elem:     modelTemplatesDocsSchema(),
				Computed: true,
			},
			"template_inputs": {
				Type:     schema.TypeList,
				Elem:     modelTemplatesInputsSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func modelTemplatesDocsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"doc_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `HTTP(S) link of the document.`,
			},
			"doc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Document name.`,
			},
		},
	}
	return &sc
}

func modelTemplatesInputsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the input parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the input parameter.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the input parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the input parameter.`,
			},
		},
	}
	return &sc
}

func resourceModelTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listModelTemplateHttpUrl = "v1/{project_id}/models/template"
		listModelTemplateProduct = "modelarts"
	)
	listModelTemplateClient, err := cfg.NewServiceClient(listModelTemplateProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	listModelTemplatePath := listModelTemplateClient.Endpoint + listModelTemplateHttpUrl
	listModelTemplatePath = strings.ReplaceAll(listModelTemplatePath, "{project_id}", listModelTemplateClient.ProjectID)

	listModelTemplatequeryParams := buildListModelTemplateQueryParams(d)
	listModelTemplatePath += listModelTemplatequeryParams

	listModelTemplateResp, err := pagination.ListAllItems(
		listModelTemplateClient,
		"offset",
		listModelTemplatePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ModelArts model template")
	}

	listModelTemplateRespJson, err := json.Marshal(listModelTemplateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listModelTemplateRespBody interface{}
	err = json.Unmarshal(listModelTemplateRespJson, &listModelTemplateRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("templates", flattenListModelTemplatetemplates(listModelTemplateRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListModelTemplatetemplates(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("templates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("template_id", v, nil),
			"name":            utils.PathSearch("template_name_en", v, nil),
			"description":     utils.PathSearch("description_en", v, nil),
			"arch":            utils.PathSearch("arch", v, nil),
			"type":            utils.PathSearch("template_labels.theme", v, nil),
			"engine":          utils.PathSearch("template_labels.engine", v, nil),
			"environment":     utils.PathSearch("template_labels.environment", v, nil),
			"template_docs":   flattenModelTemplatesDocs(v),
			"template_inputs": flattenModelTemplatesInputs(v),
		})
	}
	return rst
}

func flattenModelTemplatesDocs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("template_docs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"doc_url":  utils.PathSearch("doc_url", v, nil),
			"doc_name": utils.PathSearch("doc_name", v, nil),
		})
	}
	return rst
}

func flattenModelTemplatesInputs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("template_inputs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("input_id", v, nil),
			"type":        utils.PathSearch("input_type", v, nil),
			"name":        utils.PathSearch("name_en", v, nil),
			"description": utils.PathSearch("description_en", v, nil),
		})
	}
	return rst
}

func buildListModelTemplateQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&theme=%v", res, v)
	}

	if v, ok := d.GetOk("engine"); ok {
		res = fmt.Sprintf("%s&engine=%v", res, v)
	}

	if v, ok := d.GetOk("environment"); ok {
		res = fmt.Sprintf("%s&env=%v", res, v)
	}

	if v, ok := d.GetOk("keyword"); ok {
		res = fmt.Sprintf("%s&keyword=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
