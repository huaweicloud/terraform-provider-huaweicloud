package rc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RC POST /v1/resource-center/run-query
func DataSourceRcRunQuery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRcRunQueryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRcRunQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resource-center/run-query"
	)

	client, err := cfg.NewServiceClient("rc", region)
	if err != nil {
		return diag.Errorf("error creating RC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildRcRunQueryBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error running RC advanced query: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("query_info", flattenRunQueryInfo(respBody)),
		d.Set("results", flattenRunQueryResults(
			utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildRcRunQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expression": d.Get("expression"),
	}
	return bodyParams
}

func flattenRunQueryInfo(respBody interface{}) []interface{} {
	queryInfo := utils.PathSearch("query_info", respBody, nil)
	if queryInfo == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"select_fields": utils.PathSearch("select_fields", queryInfo, nil),
		},
	}

	return rst
}

func flattenRunQueryResults(results []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, utils.JsonToString(v))
	}

	return rst
}
