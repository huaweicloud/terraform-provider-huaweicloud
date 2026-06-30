package vpc

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

// @API VPC GET /v2.0/vpc/peerings
func DataSourcePeeringConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePeeringConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the peering connections are located.`,
			},

			// Optional parameters.
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the VPC peering connection to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the VPC peering connection to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the VPC peering connection to be queried.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The project ID of the VPC to be queried.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the requester's VPC to be queried.`,
			},

			// Attributes.
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the peering connection.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the peering connection.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the peering connection.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the peering connection.`,
						},
						"request_vpc_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the requester's VPC.`,
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The project ID of the requester to which the VPC of peering connection belongs.`,
									},
								},
							},
							Description: `The information of the requester's VPC.`,
						},
						"accept_vpc_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the accepter's VPC.`,
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The project ID of the accepter to which the VPC of peering connection belongs.`,
									},
								},
							},
							Description: `The information of the accepter's VPC.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the peering connection, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the peering connection, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of peering connections that matched the filter parameters.`,
			},
		},
	}
}

func buildListPeeringConnectionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("connection_id"); ok {
		queryParams = fmt.Sprintf("%s&id=%s", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%s", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%s", queryParams, v)
	}
	if v, ok := d.GetOk("project_id"); ok {
		queryParams = fmt.Sprintf("%s&tenant_id=%s", queryParams, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		queryParams = fmt.Sprintf("%s&vpc_id=%s", queryParams, v)
	}
	return queryParams
}

func listPeeringConnections(client *golangsdk.ServiceClient, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v2.0/vpc/peerings?limit={limit}"
		limit   = 2000
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += queryParams

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// Unfortunately, the API does not return the next marker, so just the first page can be queried.
	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("peerings", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenPeeringConnectionRequestVpcInfo(requestVpcInfo map[string]interface{}) []map[string]interface{} {
	if len(requestVpcInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc_id":     utils.PathSearch("vpc_id", requestVpcInfo, nil),
			"project_id": utils.PathSearch("tenant_id", requestVpcInfo, nil),
		},
	}
}

func flattenPeeringConnectionAcceptVpcInfo(acceptVpcInfo map[string]interface{}) []map[string]interface{} {
	if len(acceptVpcInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"vpc_id":     utils.PathSearch("vpc_id", acceptVpcInfo, nil),
			"project_id": utils.PathSearch("tenant_id", acceptVpcInfo, nil),
		},
	}
}

func flattenPeeringConnections(peerings []interface{}) []map[string]interface{} {
	if len(peerings) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(peerings))
	for _, peering := range peerings {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("id", peering, nil),
			"name":   utils.PathSearch("name", peering, nil),
			"status": utils.PathSearch("status", peering, nil),
			"request_vpc_info": flattenPeeringConnectionRequestVpcInfo(utils.PathSearch("request_vpc_info",
				peering, make(map[string]interface{})).(map[string]interface{})),
			"accept_vpc_info": flattenPeeringConnectionAcceptVpcInfo(utils.PathSearch("accept_vpc_info",
				peering, make(map[string]interface{})).(map[string]interface{})),
			"description": utils.PathSearch("description", peering, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
				peering, "").(string), "2006-01-02T15:04:05")/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
				peering, "").(string), "2006-01-02T15:04:05")/1000, false),
		})
	}
	return result
}

func dataSourcePeeringConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC peering connection client: %s", err)
	}

	perringConnections, err := listPeeringConnections(client, buildListPeeringConnectionsQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying VPC peering connections: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connections", flattenPeeringConnections(perringConnections)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
