package antiddos

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

// @API Anti-DDOS GET /v2/{project_id}/antiddos
func DataSourceEipDefenseStatusesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEipDefenseStatusesV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protection status of the EIP.`,
			},
			"ips": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IP address for filtering, supports partial matching.`,
			},
			"ddos_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of EIP defense statuses.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"floating_ip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the EIP.`,
						},
						"floating_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of the EIP.`,
						},
						"product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the EIP protection service.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The protection status of the EIP.`,
						},
						"clean_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The cleaning threshold.`,
						},
						"block_threshold": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The blackhole threshold.`,
						},
					},
				},
			},
		},
	}
}

func buildEipDefenseStatusesV2QueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"

	if v, ok := d.GetOk("status"); ok {
		queryParams += fmt.Sprintf("&status=%s", v.(string))
	}

	if v, ok := d.GetOk("ips"); ok {
		queryParams += fmt.Sprintf("&ips=%s", v.(string))
	}

	return queryParams
}

func dataSourceEipDefenseStatusesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/antiddos"
		product = "anti-ddos"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEipDefenseStatusesV2QueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := requestPath + fmt.Sprintf("&offset=%d", offset)

		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving Anti-DDoS V2 EIP defense statuses: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("ddosStatus", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		allData = append(allData, dataList...)
		offset += len(dataList)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ddos_status", flattenV2EipDefenseStatuses(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenV2EipDefenseStatuses(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(respArray))
	for _, respBody := range respArray {
		result = append(result, map[string]interface{}{
			"floating_ip_id":      utils.PathSearch("floating_ip_id", respBody, nil),
			"floating_ip_address": utils.PathSearch("floating_ip_address", respBody, nil),
			"product_type":        utils.PathSearch("product_type", respBody, nil),
			"status":              utils.PathSearch("status", respBody, nil),
			"clean_threshold":     utils.PathSearch("clean_threshold", respBody, nil),
			"block_threshold":     utils.PathSearch("block_threshold", respBody, nil),
		})
	}

	return result
}
