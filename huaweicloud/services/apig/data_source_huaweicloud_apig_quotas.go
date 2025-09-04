package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instance/configs
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the quotas.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of quotas.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the quota.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the quota.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the quota.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the quota.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the quota, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func listQuotas(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instance/configs?limit={limit}"
		limit   = 500
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		quotas := utils.PathSearch("configs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, quotas...)
		if len(quotas) < limit {
			break
		}

		offset += len(quotas)
	}

	return result, nil
}

func flattenQuotas(quotas []interface{}) []map[string]interface{} {
	if len(quotas) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(quotas))
	for _, val := range quotas {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("config_id", val, nil),
			"name":        utils.PathSearch("config_name", val, nil),
			"value":       utils.PathSearch("config_value", val, nil),
			"description": utils.PathSearch("remark", val, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch(
				"config_time", val, "").(float64)/1000), false),
		})
	}

	return result
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	quotas, err := listQuotas(client)
	if err != nil {
		return diag.Errorf("error getting quotas: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(quotas)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
