package aad

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

// Due to limited testing conditions, this data source cannot be tested and the API was not successfully called.

// @API AAD GET /v2/aad/instances/{instance_id}/ddos-info/flow/connection-numbers
func DataSourceConnectionNumbers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionNumbersRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the start time.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the end time.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the IP address.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The connection number data list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connection number name.",
						},
						"list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The connection number data items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The timestamp in milliseconds.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The connection number value.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildConnectionNumbersQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?start_time=%v", d.Get("start_time"))
	queryParams += fmt.Sprintf("&end_time=%v", d.Get("end_time"))
	queryParams += fmt.Sprintf("&ip=%v", d.Get("ip"))

	return queryParams
}

func dataSourceConnectionNumbersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "aad"
		httpUrl    = "v2/aad/instances/{instance_id}/ddos-info/flow/connection-numbers"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestPath += buildConnectionNumbersQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD connection numbers: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("data", flattenConnectionNumbersData(utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConnectionNumbersData(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"list": flattenConnectionNumbersList(utils.PathSearch("list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenConnectionNumbersList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"time":  utils.PathSearch("time", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
