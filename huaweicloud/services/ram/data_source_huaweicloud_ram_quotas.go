package ram

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM GET /v1/resource-shares/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,
		Schema: map[string]*schema.Schema{
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotasSchema(),
			},
		},
	}
}

func quotasSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listQuotasHttpUrl := "v1/resource-shares/quotas"
	listQuotasProduct := "ram"
	listQuotasClient, err := cfg.NewServiceClient(listQuotasProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	listQuotasPath := listQuotasClient.Endpoint + listQuotasHttpUrl

	listQuotasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var quotas []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listQuotasPath + buildListQuotasQueryParams(marker)
		listQuotasResp, err := listQuotasClient.Request("GET", queryPath, &listQuotasOpt)
		if err != nil {
			return diag.Errorf("error retrieving RAM quotas, %s", err)
		}

		listQuotasRespBody, err := utils.FlattenResponse(listQuotasResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageQuotas := flattenQuotasResp(listQuotasRespBody)
		quotas = append(quotas, onePageQuotas...)
		marker = utils.PathSearch("page_info.next_marker", listQuotasRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("quotas", quotas),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListQuotasQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenQuotasResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	quotas := make([]interface{}, 0)

	curJson := utils.PathSearch("quotas", resp, make([]interface{}, 0))
	curMap := curJson.(map[string]interface{})
	if resources, ok := curMap["resources"]; ok {
		childArray := resources.([]interface{})
		rst := make([]interface{}, 0, len(childArray))
		for _, v := range childArray {
			rst = append(rst, map[string]interface{}{
				"type":  utils.PathSearch("type", v, nil),
				"quota": utils.PathSearch("quota", v, nil),
				"min":   utils.PathSearch("min", v, nil),
				"max":   utils.PathSearch("max", v, nil),
				"used":  utils.PathSearch("used", v, nil),
			})
		}
		quotas = append(quotas, map[string]interface{}{
			"resources": rst,
		})
	} else {
		quotas = append(quotas, map[string]interface{}{
			"resources": []interface{}{},
		})
	}

	return quotas
}
