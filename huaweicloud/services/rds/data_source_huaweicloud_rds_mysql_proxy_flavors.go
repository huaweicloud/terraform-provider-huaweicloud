package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/proxy/flavors
func DataSourceRdsMysqlProxyFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlProxyFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cpu": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mem": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"az_status": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
func dataSourceRdsMysqlProxyFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/flavors"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)
	listBasePath = strings.ReplaceAll(listBasePath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	limit := 100
	offset := 0
	flavorGroups := make(map[string][]map[string]interface{}, 0)
	for {
		listPath := listBasePath + buildListMysqlProxyFlavorsQueryParams(limit, offset)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS MySQL proxy flavors: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		flavorGroupsMap, maxCount := flattenMysqlProxyFlavorGroupsResp(listRespBody)
		for groupType, flavors := range flavorGroupsMap {
			flavorGroups[groupType] = append(flavorGroups[groupType], flavors...)
		}
		if maxCount < limit {
			break
		}
		offset += limit
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	res := make([]map[string]interface{}, 0, len(flavorGroups))
	for groupType, flavors := range flavorGroups {
		res = append(res, map[string]interface{}{
			"group_type": groupType,
			"flavors":    flavors,
		})
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("flavor_groups", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListMysqlProxyFlavorsQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func flattenMysqlProxyFlavorGroupsResp(listRespBody interface{}) (map[string][]map[string]interface{}, int) {
	flavorGroupsJson := utils.PathSearch("compute_flavor_groups", listRespBody, nil)
	if flavorGroupsJson == nil {
		return nil, 0
	}

	result := make(map[string][]map[string]interface{})
	maxCount := 0
	flavorGroupsArray := flavorGroupsJson.([]interface{})
	for _, flavorGroup := range flavorGroupsArray {
		groupType := utils.PathSearch("group_type", flavorGroup, "").(string)
		flavors := flattenMysqlProxyFlavorGroupFlavorsResp(flavorGroup)
		result[groupType] = flavors
		if len(flavors) > maxCount {
			maxCount = len(flavors)
		}
	}
	return result, maxCount
}

func flattenMysqlProxyFlavorGroupFlavorsResp(flavors interface{}) []map[string]interface{} {
	flavorsJson := utils.PathSearch("compute_flavors", flavors, nil)
	if flavorsJson == nil {
		return nil
	}

	flavorsArray := flavorsJson.([]interface{})
	result := make([]map[string]interface{}, 0, len(flavorsArray))
	for _, flavor := range flavorsArray {
		result = append(result, map[string]interface{}{
			"id":        utils.PathSearch("id", flavor, nil),
			"code":      utils.PathSearch("code", flavor, nil),
			"cpu":       utils.PathSearch("cpu", flavor, nil),
			"mem":       utils.PathSearch("mem", flavor, nil),
			"db_type":   utils.PathSearch("db_type", flavor, nil),
			"az_status": utils.PathSearch("az_status", flavor, nil),
		})
	}

	return result
}
