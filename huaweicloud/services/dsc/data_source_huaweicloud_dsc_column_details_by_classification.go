package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/metadata/catalog/column-details/classification-dim
func DataSourceDscColumnDetailsByClassification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscColumnDetailsByClassificationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the group label ID for filtering. Either `label_id` or `type_id` must be specified.",
			},
			"type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type ID for filtering. Either `label_id` or `type_id` must be specified.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column details list by classification dimension.",
				Elem:        dscTemplateDimSchema(),
			},
		},
	}
}

func dscTemplateDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template name.",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The match count.",
			},
			"classifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The classification information list.",
				Elem:        dscClassificationDimSchema(),
			},
		},
	}
}

func dscClassificationDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"classification_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The classification name.",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The match count.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information and match information list.",
				Elem:        dscColumnInfoAndMatchInfoSchema(),
			},
		},
	}
}

func dscColumnInfoAndMatchInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset ID.",
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"column_fqn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The column fully qualified name.",
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database type.",
			},
			"match_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The match information list.",
				Elem:        dscClassificationMatchInfoSchema(),
			},
		},
	}
}

func dscClassificationMatchInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"classification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The classification ID.",
			},
			"classification_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The classification name.",
			},
			"match_content_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The matched content count.",
			},
			"match_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The match rate (percentage).",
			},
			"matched_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The match detail.",
			},
			"matched_examples": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The matched example list.",
				Elem:        dscMatchedExamplesSchema(),
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule name.",
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The security level color.",
			},
			"security_level_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level ID.",
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level name.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template ID.",
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template name.",
			},
		},
	}
}

func dscMatchedExamplesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"context": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The match context.",
			},
			"line_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The line number of the match.",
			},
			"matched_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The matched content.",
			},
			"nlp_revised": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether it has been NLP revised.",
			},
		},
	}
}

func buildDscColumnDetailsByClassificationQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("label_id"); ok {
		queryParams = fmt.Sprintf("%s&label_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type_id"); ok {
		queryParams = fmt.Sprintf("%s&type_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscColumnDetailsByClassificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/metadata/catalog/column-details/classification-dim"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildDscColumnDetailsByClassificationQueryParams(d)

	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC column details by classification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("results", flattenDscTemplateDim(
			utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscTemplateDim(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"template_name": utils.PathSearch("template_name", v, nil),
			"count":         utils.PathSearch("count", v, nil),
			"classifications": flattenDscClassificationDim(
				utils.PathSearch("classifications", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscClassificationDim(classifications []interface{}) []interface{} {
	if len(classifications) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(classifications))
	for _, v := range classifications {
		rst = append(rst, map[string]interface{}{
			"classification_name": utils.PathSearch("classification_name", v, nil),
			"count":               utils.PathSearch("count", v, nil),
			"columns": flattenDscColumnInfoAndMatchInfo(
				utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscColumnInfoAndMatchInfo(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, map[string]interface{}{
			"asset_id":   utils.PathSearch("asset_id", v, nil),
			"asset_name": utils.PathSearch("asset_name", v, nil),
			"column_fqn": utils.PathSearch("column_fqn", v, nil),
			"db_type":    utils.PathSearch("db_type", v, nil),
			"match_infos": flattenDscClassificationMatchInfo(
				utils.PathSearch("match_infos", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscClassificationMatchInfo(matchInfos []interface{}) []interface{} {
	if len(matchInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(matchInfos))
	for _, v := range matchInfos {
		rst = append(rst, map[string]interface{}{
			"classification_id":    utils.PathSearch("classification_id", v, nil),
			"classification_name":  utils.PathSearch("classification_name", v, nil),
			"match_content_cnt":    utils.PathSearch("match_content_cnt", v, nil),
			"match_rate":           utils.PathSearch("match_rate", v, nil),
			"matched_detail":       utils.PathSearch("matched_detail", v, nil),
			"matched_examples":     flattenDscMatchedExamples(utils.PathSearch("matched_examples", v, make([]interface{}, 0)).([]interface{})),
			"rule_id":              utils.PathSearch("rule_id", v, nil),
			"rule_name":            utils.PathSearch("rule_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"template_name":        utils.PathSearch("template_name", v, nil),
		})
	}

	return rst
}

func flattenDscMatchedExamples(examples []interface{}) []interface{} {
	if len(examples) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(examples))
	for _, v := range examples {
		rst = append(rst, map[string]interface{}{
			"context":         utils.PathSearch("context", v, nil),
			"line_number":     utils.PathSearch("line_number", v, nil),
			"matched_content": utils.PathSearch("matched_content", v, nil),
			"nlp_revised":     utils.PathSearch("nlp_revised", v, nil),
		})
	}

	return rst
}
