package coc

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

// @API COC GET /v1/atomics/{atomic_unique_key}
func DataSourceCocDocumentAtomicDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocDocumentAtomicDetailRead,

		Schema: map[string]*schema.Schema{
			"atomic_unique_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the unique identifier of an atomic capability.`,
			},
			"atomic_name_zh": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the Chinese name.`,
			},
			"atomic_name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the English name.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the tag information.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the atomic capability input.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter variable name.`,
						},
						"param_name_zh": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the Chinese name of the parameter.`,
						},
						"param_name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the English name of the parameter.`,
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the field is required.`,
						},
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter type: constant/dictionary.`,
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the minimum value.`,
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the maximum value.`,
						},
						"min_len": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the minimum length.`,
						},
						"max_len": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the maximum length.`,
						},
						"pattern": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the regular restriction expression.`,
						},
					},
				},
			},
			"outputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the atomic capability output.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"supported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether output is supported.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the output type.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCocDocumentAtomicDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/atomics/{atomic_unique_key}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{atomic_unique_key}", d.Get("atomic_unique_key").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error querying document atomic detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error querying document atomic detail: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("atomic_name_zh", utils.PathSearch("data.atomic_name_zh", respBody, nil)),
		d.Set("atomic_name_en", utils.PathSearch("data.atomic_name_en", respBody, nil)),
		d.Set("tags", utils.PathSearch("data.tags", respBody, nil)),
		d.Set("inputs", flattenCocDocumentAtomicDetailInputs(
			utils.PathSearch("data.inputs", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("outputs", flattenCocDocumentAtomicDetailOutputs(
			utils.PathSearch("data.outputs", respBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCocDocumentAtomicDetailInputs(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"param_key":     utils.PathSearch("param_key", params, nil),
			"param_name_zh": utils.PathSearch("param_name_zh", params, nil),
			"param_name_en": utils.PathSearch("param_name_en", params, nil),
			"required":      utils.PathSearch("required", params, nil),
			"param_type":    utils.PathSearch("param_type", params, nil),
			"min":           utils.PathSearch("min", params, nil),
			"max":           utils.PathSearch("max", params, nil),
			"min_len":       utils.PathSearch("min_len", params, nil),
			"max_len":       utils.PathSearch("max_len", params, nil),
			"pattern":       utils.PathSearch("pattern", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenCocDocumentAtomicDetailOutputs(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"supported": utils.PathSearch("supported", param, nil),
			"type":      utils.PathSearch("type", param, nil),
		},
	}

	return rst
}
