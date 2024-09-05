package rms

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

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/{aggregator_id}/run-query
func DataSourceAggregatorAdvancedQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorAdvancedQueryRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource aggregator ID.`,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"select_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceAggregatorAdvancedQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getAggregatorAdvancedQueryHttpUrl := "v1/resource-manager/domains/{domain_id}/aggregators/{aggregator_id}/run-query"
	getAggregatorAdvancedQueryPath := client.Endpoint + getAggregatorAdvancedQueryHttpUrl
	getAggregatorAdvancedQueryPath = strings.ReplaceAll(getAggregatorAdvancedQueryPath, "{domain_id}", cfg.DomainID)
	getAggregatorAdvancedQueryPath = strings.ReplaceAll(getAggregatorAdvancedQueryPath, "{aggregator_id}", d.Get("aggregator_id").(string))

	getAggregatorAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAggregatorAdvancedQueryBodyParams(d)),
	}

	requestPath := getAggregatorAdvancedQueryPath
	resp, err := client.Request("POST", requestPath, &getAggregatorAdvancedQueryOpt)
	if err != nil {
		return diag.Errorf("error doing RMS aggregator advanced query: %s", err)
	}

	getAggregatorAdvancedQueryRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("results", utils.PathSearch("results", getAggregatorAdvancedQueryRespBody, nil)),
		d.Set("query_info", flattenAggregatorqueryInfo(getAggregatorAdvancedQueryRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorAdvancedQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expression": d.Get("expression"),
	}
	return bodyParams
}

func flattenAggregatorqueryInfo(getAggregatorAdvancedQueryRespBody interface{}) []map[string]interface{} {
	queryInfo := utils.PathSearch("query_info", getAggregatorAdvancedQueryRespBody, nil)
	if queryInfo == nil {
		return nil
	}

	res := []map[string]interface{}{
		{
			"select_fields": utils.PathSearch("select_fields", queryInfo, nil),
		},
	}

	return res
}
