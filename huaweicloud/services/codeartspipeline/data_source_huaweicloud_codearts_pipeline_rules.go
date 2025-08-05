package codeartspipeline

import (
	"context"
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

// @API CodeArtsPipeline GET /v2/{domain_id}/rules/query
func DataSourceCodeArtsPipelineRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the rule name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the rule type.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the rule list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the rule ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the rule name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the rule type.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the rule version.`,
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operator.`,
						},
						"operate_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the operate time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{domain_id}/rules/query?limit=10"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{domain_id}", cfg.DomainID)
	listPath += buildCodeArtsPipelineRulesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error getting rules: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(listRespBody, projectNotFoundError); err != nil {
			return diag.Errorf("error getting rules: %s", err)
		}

		rules := utils.PathSearch("data", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}

		for _, rule := range rules {
			rst = append(rst, map[string]interface{}{
				"id":           utils.PathSearch("id", rule, nil),
				"name":         utils.PathSearch("name", rule, nil),
				"type":         utils.PathSearch("type", rule, nil),
				"version":      utils.PathSearch("version", rule, nil),
				"operator":     utils.PathSearch("operator", rule, nil),
				"operate_time": utils.PathSearch("operate_time", rule, nil),
			})
		}

		offset += len(rules)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeArtsPipelineRulesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("project_id"); ok {
		res += fmt.Sprintf("&cloud_project_id=%v", v)
	}
	if v, ok := d.GetOk("name"); ok {
		res += fmt.Sprintf("&name=%v", v)
	}
	if v, ok := d.GetOk("type"); ok {
		res += fmt.Sprintf("&type=%v", v)
	}

	return res
}
