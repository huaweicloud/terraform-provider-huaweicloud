package cnad

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

// @API AAD GET /v1/unblockservice/{domain_id}/block-list
func DataSourceAdvancedBlockList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedBlockListRead,
		Schema: map[string]*schema.Schema{
			"blocking_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of blocked IPs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The blocked IP address.`,
						},
						"blocking_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the IP was blocked.`,
						},
						"estimated_unblocking_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The estimated time when the IP will be unblocked.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the blocked IP.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceAdvancedBlockListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v1/unblockservice/{domain_id}/block-list"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CNAD advanced block list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	blockingList := utils.PathSearch("blocking_list", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("blocking_list", flattenAdvancedBlockList(blockingList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAdvancedBlockList(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, item := range respArray {
		rst = append(rst, map[string]interface{}{
			"ip":                        utils.PathSearch("ip", item, nil),
			"blocking_time":             utils.PathSearch("blocking_time", item, nil),
			"estimated_unblocking_time": utils.PathSearch("estimated_unblocking_time", item, nil),
			"status":                    utils.PathSearch("status", item, nil),
		})
	}
	return rst
}
