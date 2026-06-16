package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/instances
func DataSourceDataArtsStudioInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataArtsStudioInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DataArts Studio instances are located.`,
			},

			// Attributes.
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the DataArts Studio instance.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the DataArts Studio instance.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version (spec code) of the DataArts Studio instance.`,
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The order ID of the DataArts Studio instance.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID of the DataArts Studio instance.`,
						},
						"auto_renew": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Whether auto renew is enabled for the DataArts Studio instance.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the DataArts Studio instance.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID of the DataArts Studio instance.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet ID of the DataArts Studio instance.`,
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone of the DataArts Studio instance.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the DataArts Studio instance belongs.`,
						},
						"effective_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The effective time of the DataArts Studio instance, in RFC3339 format.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the DataArts Studio instance.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the DataArts Studio instance, in RFC3339 format.`,
						},
						"workspace_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace mode of the DataArts Studio instance.`,
						},
					},
				},
				Description: `The list of the DataArts Studio instances.`,
			},
		},
	}
}

func listDataArtsStudioInstances(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/instances?limit={limit}"
		limit   = 100
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
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		instances := utils.PathSearch("commodity_orders", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		if len(instances) < limit {
			break
		}
		offset += len(instances)
	}

	return result, nil
}

func parseStudioInstanceAutoRenew(autoRenew int) string {
	if autoRenew == 1 {
		return "true"
	}
	return "false"
}

func flattenDataArtsStudioInstances(instances []interface{}) []map[string]interface{} {
	if len(instances) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("resource_id", instance, nil),
			"name":       utils.PathSearch("resource_name", instance, nil),
			"version":    utils.PathSearch("resource_spec_code", instance, nil),
			"order_id":   utils.PathSearch("order_id", instance, nil),
			"product_id": utils.PathSearch("product_id", instance, nil),
			"auto_renew": parseStudioInstanceAutoRenew(int(utils.PathSearch("is_auto_renew",
				instance, float64(0)).(float64))),
			"status":                utils.PathSearch("status", instance, nil),
			"vpc_id":                utils.PathSearch("vpc_id", instance, nil),
			"subnet_id":             utils.PathSearch("net_id", instance, nil),
			"availability_zone":     utils.PathSearch("availability_zone", instance, nil),
			"enterprise_project_id": utils.PathSearch("eps_id", instance, nil),
			"effective_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("effective_time",
				instance, float64(0)).(float64))/1000, false),
			"created_by": utils.PathSearch("create_user", instance, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				instance, float64(0)).(float64))/1000, false),
			"workspace_mode": utils.PathSearch("work_space_mode", instance, nil),
		})
	}

	return result
}

func dataSourceDataArtsStudioInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	instances, err := listDataArtsStudioInstances(client)
	if err != nil {
		return diag.Errorf("error querying DataArts Studio instances: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenDataArtsStudioInstances(instances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
