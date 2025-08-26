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

// @API AAD GET /v1/aad/instances/{instance_id}/{ip}/ddos-statistics
func DataSourceIpDdosStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpDdosStatisticsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attack_kbps_peak": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"in_kbps_peak": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ddos_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildIpDdosStatisticsQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?start_time=%v", d.Get("start_time"))
	queryParams += fmt.Sprintf("&end_time=%v", d.Get("end_time"))

	return queryParams
}

func dataSourceIpDdosStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v1/aad/instances/{instance_id}/{ip}/ddos-statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{ip}", d.Get("ip").(string))
	requestPath += buildIpDdosStatisticsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD IP ddos statistics: %s", err)
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
		d.Set("attack_kbps_peak", utils.PathSearch("attack_kbps_peak", respBody, nil)),
		d.Set("in_kbps_peak", utils.PathSearch("in_kbps_peak", respBody, nil)),
		d.Set("ddos_count", utils.PathSearch("ddos_count", respBody, nil)),
		d.Set("timestamp", utils.PathSearch("timestamp", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
