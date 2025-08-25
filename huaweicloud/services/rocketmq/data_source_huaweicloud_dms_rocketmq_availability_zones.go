package rocketmq

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

// @API RocketMQ GET /v2/available-zones
func DataSourceRocketMQAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRocketMQAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the availability zones are located.",
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the availability zone.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the availability zone.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The code of the availability zone.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port number of the availability zone.",
						},
						"sold_out": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the availability zone is sold out.",
						},
						"resource_availability": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether there are available resources in the availability zone.",
						},
						"default_az": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the availability zone is default.",
						},
						"remain_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The remaining time of the availability zone.",
						},
						"ipv6_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether IPv6 is supported in the availability zone.",
						},
					},
				},
				Description: "The list of availability zones.",
			},
		},
	}
}

func dataSourceRocketMQAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/available-zones"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error querying availability zones: %s", err)
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
		d.Set("availability_zones", flattenRocketMQAvailabilityZones(utils.PathSearch("available_zones",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRocketMQAvailabilityZones(azs []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(azs))
	for i, az := range azs {
		result[i] = map[string]interface{}{
			"id":                    utils.PathSearch("id", az, nil),
			"name":                  utils.PathSearch("name", az, nil),
			"code":                  utils.PathSearch("code", az, nil),
			"port":                  utils.PathSearch("port", az, nil),
			"sold_out":              utils.PathSearch("soldOut", az, nil),
			"resource_availability": utils.PathSearch("resource_availability", az, nil),
			"default_az":            utils.PathSearch("default_az", az, nil),
			"remain_time":           utils.PathSearch("remain_time", az, nil),
			"ipv6_enable":           utils.PathSearch("ipv6_enable", az, nil),
		}
	}
	return result
}
