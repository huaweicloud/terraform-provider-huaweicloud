package eip

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

// @API EIP GET /v3/{domain_id}/geip/internet-bandwidth-limits
func DataSourceGlobalEipInternetBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipInternetBandwidthLimitsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// The `sort_key` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `sort_dir` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `charge_mode` field in the API documentation is of type list, which is unnecessary.
			// Here, it is modified to type string.
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_bandwidth_limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ext_limit": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_ingress_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_ingress_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ratio_95peak": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipInternetBandwidthLimitsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if raw, ok := d.GetOk("fields"); ok {
		for _, f := range raw.([]interface{}) {
			queryParams += fmt.Sprintf("&fields=%s", f)
		}
	}
	if raw, ok := d.GetOk("sort_key"); ok {
		queryParams += fmt.Sprintf("&sort_key=%s", raw)
	}
	if raw, ok := d.GetOk("sort_dir"); ok {
		queryParams += fmt.Sprintf("&sort_dir=%s", raw)
	}
	if raw, ok := d.GetOk("charge_mode"); ok {
		queryParams += fmt.Sprintf("&charge_mode=%s", raw)
	}
	if raw, ok := d.GetOk("type"); ok {
		queryParams += fmt.Sprintf("&type=%s", raw)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceGlobalEipInternetBandwidthLimitsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/geip/internet-bandwidth-limits"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath += buildGlobalEipInternetBandwidthLimitsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving global EIP internet bandwidth limits: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("internet_bandwidth_limits", flattenGlobalEipInternetBandwidthLimits(
			utils.PathSearch("internet_bandwidth_limits", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipInternetBandwidthLimits(limits []interface{}) []interface{} {
	if len(limits) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(limits))
	for _, limit := range limits {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", limit, nil),
			"charge_mode": utils.PathSearch("charge_mode", limit, nil),
			"min_size":    utils.PathSearch("min_size", limit, nil),
			"ext_limit": flattenGlobalEipInternetBandwidthLimitsExtLimit(
				utils.PathSearch("ext_limit", limit, nil)),
			"max_size": utils.PathSearch("max_size", limit, nil),
			"type":     utils.PathSearch("type", limit, nil),
		})
	}

	return rst
}

func flattenGlobalEipInternetBandwidthLimitsExtLimit(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"min_ingress_size": utils.PathSearch("min_ingress_size", raw, nil),
			"max_ingress_size": utils.PathSearch("max_ingress_size", raw, nil),
			"ratio_95peak":     utils.PathSearch("ratio_95peak", raw, nil),
		},
	}
}
