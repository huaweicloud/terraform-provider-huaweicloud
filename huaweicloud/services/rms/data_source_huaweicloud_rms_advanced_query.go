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

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/run-query
func DataSourceAdvancedQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedQueryRead,

		Schema: map[string]*schema.Schema{
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

func dataSourceAdvancedQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getAdvancedQueryHttpUrl := "v1/resource-manager/domains/{domain_id}/run-query"
	getAdvancedQueryPath := client.Endpoint + getAdvancedQueryHttpUrl
	getAdvancedQueryPath = strings.ReplaceAll(getAdvancedQueryPath, "{domain_id}", cfg.DomainID)

	getAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAdvancedQueryBodyParams(d)),
	}

	requestPath := getAdvancedQueryPath
	resp, err := client.Request("POST", requestPath, &getAdvancedQueryOpt)
	if err != nil {
		return diag.Errorf("error doing RMS advanced query: %s", err)
	}

	getAdvancedQueryRespBody, err := utils.FlattenResponse(resp)
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
		d.Set("results", utils.PathSearch("results", getAdvancedQueryRespBody, nil)),
		d.Set("query_info", flattenQueryInfo(getAdvancedQueryRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAdvancedQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expression": d.Get("expression"),
	}
	return bodyParams
}

func flattenQueryInfo(getAdvancedQueryRespBody interface{}) []map[string]interface{} {
	queryInfo := utils.PathSearch("query_info", getAdvancedQueryRespBody, nil)
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
