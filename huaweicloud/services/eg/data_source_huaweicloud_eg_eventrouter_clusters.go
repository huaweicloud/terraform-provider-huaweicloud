package eg

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

// @API EG GET /v1/{project_id}/eventrouter/clusters
func DataSourceEventRouterClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventRouterClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the event router clusters are located.",
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the event router cluster to be queried for fuzzy matching.",
			},

			// Attributes.
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of the event router clusters that matched filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the event router cluster.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the event router cluster.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the event router cluster.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source type of the event router cluster.",
						},
						"sink_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The sink type of the event router cluster.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID to which the event router cluster belongs.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subnet ID to which the event router cluster belongs.",
						},
						"availability_zones": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone names of the event router cluster.",
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The flavor of the event router cluster.",
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charging mode of the event router cluster.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the event router cluster.",
						},
						"job_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of jobs running in the event router cluster.",
						},
						"err_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The error code of the event router cluster.",
						},
						"err_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The error message of the event router cluster.",
						},
						"public_access_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether public access is enabled for the event router cluster.",
						},
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The NAT gateway ID of the event router cluster.",
						},
						"eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP ID of the event router cluster.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the event router cluster, in RFC3339 format.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the event router cluster, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func buildEventRouterClustersQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&fuzzy_name=%v", res, v)
	}
	return res
}

func listEventRouterClusters(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/eventrouter/clusters?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildEventRouterClustersQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		clusters := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, clusters...)
		if len(clusters) < limit {
			break
		}
		offset += len(clusters)
	}

	return result, nil
}

func flattenEventRouterClusters(clusters []interface{}) []map[string]interface{} {
	if len(clusters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("cluster_id", cluster, nil),
			"name":                  utils.PathSearch("name", cluster, nil),
			"description":           utils.PathSearch("description", cluster, nil),
			"source_type":           utils.PathSearch("source_type", cluster, nil),
			"sink_type":             utils.PathSearch("sink_type", cluster, nil),
			"vpc_id":                utils.PathSearch("vpc_id", cluster, nil),
			"subnet_id":             utils.PathSearch("subnet_id", cluster, nil),
			"availability_zones":    utils.PathSearch("zone_names", cluster, nil),
			"flavor":                utils.PathSearch("flavor", cluster, nil),
			"charging_mode":         utils.PathSearch("charging_mode", cluster, nil),
			"status":                utils.PathSearch("status", cluster, nil),
			"job_count":             utils.PathSearch("job_count", cluster, nil),
			"err_code":              utils.PathSearch("err_code", cluster, nil),
			"err_message":           utils.PathSearch("err_message", cluster, nil),
			"public_access_enabled": utils.PathSearch("public_access_enabled", cluster, nil),
			"nat_id":                utils.PathSearch("nat_id", cluster, nil),
			"eip_id":                utils.PathSearch("eip_id", cluster, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_time",
				cluster, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_time",
				cluster, "").(string))/1000, false),
		})
	}

	return result
}

func dataSourceEventRouterClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	clusters, err := listEventRouterClusters(client, d)
	if err != nil {
		return diag.Errorf("error querying event router clusters: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("clusters", flattenEventRouterClusters(clusters)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
