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

// @API DSC GET /v1/{project_id}/metadata/catalog/classification-tops
func DataSourceCatalogTopClassifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCatalogTopClassificationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the catalog top classifications are located.`,
			},

			// Optional parameters.
			"label_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"type_id"},
				Description:  `The ID of the group label.`,
			},
			"type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the data type.`,
			},

			// Attributes.
			"classifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"classification_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the classification.`,
						},
						"hit_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of matched records.`,
						},
						"column_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        catalogTopClassificationColumnDetails(),
							Description: `The column detail list of the classification.`,
						},
					},
				},
				Description: `The list of classifications.`,
			},
		},
	}
}

func catalogTopClassificationColumnDetails() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the asset corresponding to the database.`,
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the asset corresponding to the database.`,
			},
			"column_fqn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fully qualified name of the column.`,
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the database.`,
			},
			"match_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        catalogTopClassificationColumnDetailsMatchInfos(),
				Description: `The match information list.`,
			},
		},
	}
}

func catalogTopClassificationColumnDetailsMatchInfos() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"classification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the classification.`,
			},
			"classification_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the classification.`,
			},
			"match_content_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The matched content count.`,
			},
			"match_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The match rate (percentage).`,
			},
			"matched_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The match detail.`,
			},
			"matched_examples": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        catalogTopClassificationMatchedExamples(),
				Description: `The matched example list.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the rule.`,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the rule.`,
			},
			"security_level_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the security level.`,
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the security level.`,
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The color corresponding to the security level.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the template.`,
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the template.`,
			},
		},
	}
}

func catalogTopClassificationMatchedExamples() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"context": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The match context.`,
			},
			"line_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The line number of the match.`,
			},
			"matched_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The matched content.`,
			},
			"nlp_revised": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it has been NLP revised.`,
			},
		},
	}
}

func buildCatalogTopClassificationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("label_id"); ok {
		res = fmt.Sprintf("%s&label_id=%v", res, v)
	}

	if v, ok := d.GetOk("type_id"); ok {
		res = fmt.Sprintf("%s&type_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func listCatalogTopClassifications(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{project_id}/metadata/catalog/classification-tops"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildCatalogTopClassificationsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceCatalogTopClassificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := listCatalogTopClassifications(client, d)
	if err != nil {
		return diag.Errorf("error retrieving catalog top 5 classifications: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("classifications", flattenCatalogTopClassifications(utils.PathSearch("detection_classifications",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalogTopClassifications(classifications []interface{}) []interface{} {
	if len(classifications) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(classifications))
	for _, v := range classifications {
		rst = append(rst, map[string]interface{}{
			"classification_name": utils.PathSearch("classification_name", v, nil),
			"hit_number":          utils.PathSearch("hit_number", v, nil),
			"column_details": flattenCatalogTopClassificationColumnDetails(utils.PathSearch("column_details",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenCatalogTopClassificationColumnDetails(columns []interface{}) []interface{} {
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
			"match_infos": flattenCatalogTopClassificationMatchInfos(utils.PathSearch("match_infos",
				v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenCatalogTopClassificationMatchInfos(matchInfos []interface{}) []interface{} {
	if len(matchInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(matchInfos))
	for _, v := range matchInfos {
		rst = append(rst, map[string]interface{}{
			"classification_id":   utils.PathSearch("classification_id", v, nil),
			"classification_name": utils.PathSearch("classification_name", v, nil),
			"match_content_cnt":   utils.PathSearch("match_content_cnt", v, nil),
			"match_rate":          utils.PathSearch("match_rate", v, nil),
			"matched_detail":      utils.PathSearch("matched_detail", v, nil),
			"matched_examples": flattenCatalogTopClassificationMatchedExamples(utils.PathSearch("matched_examples",
				v, make([]interface{}, 0)).([]interface{})),
			"rule_id":              utils.PathSearch("rule_id", v, nil),
			"rule_name":            utils.PathSearch("rule_name", v, nil),
			"security_level_id":    utils.PathSearch("security_level_id", v, nil),
			"security_level_name":  utils.PathSearch("security_level_name", v, nil),
			"security_level_color": utils.PathSearch("security_level_color", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"template_name":        utils.PathSearch("template_name", v, nil),
		})
	}

	return rst
}

func flattenCatalogTopClassificationMatchedExamples(examples []interface{}) []interface{} {
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
