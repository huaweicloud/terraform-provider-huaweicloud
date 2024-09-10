package er

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/er/v3/associations"
	"github.com/chnsz/golangsdk/openstack/er/v3/propagations"
	"github.com/chnsz/golangsdk/openstack/er/v3/routes"
	"github.com/chnsz/golangsdk/openstack/er/v3/routetables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/associations
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/propagations
// @API ER GET /v3/{project_id}/enterprise-router/route-tables/{route_table_id}/static-routes
func DataSourceRouteTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouteTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the ER instance and route table are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the ER instance to which the route table belongs.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The route table ID used to query specified route table.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name used to filter the route tables.`,
			},
			"tags": common.TagsSchema(),
			// Attributes
			"route_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The route table ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the route table.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the route table.`,
						},
						"associations": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        routeTableRelationshipSchemaResource(),
							Description: `The association configuration of the route table.`,
						},
						"propagations": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        routeTableRelationshipSchemaResource(),
							Description: `The propagation configuration of the route table.`,
						},
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The route ID.`,
									},
									"destination": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The destination address (CIDR) of the route.`,
									},
									"is_blackhole": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether route is the black hole route.`,
									},
									"attachments": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attachment_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The ID of the nexthop attachment.`,
												},
												"attachment_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The type of the nexthop attachment.`,
												},
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The ID of the resource associated with the attachment.`,
												},
											},
										},
										Description: `The details of the attachment corresponding to the route.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The current status of the route.`,
									},
								},
							},
							Description: `The route details of the route table.`,
						},
						"is_default_association": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this route table is the default association route table.`,
						},
						"is_default_propagation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this route table is the default propagation route table.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the route table.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The tags configuration of the route table.`,
						},
					},
				},
			},
		},
	}
}

func routeTableRelationshipSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the association/propagation.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attachment ID corresponding to the routing association/propagation.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The attachment type corresponding to the routing association/propagation.`,
			},
		},
	}
}

func queryRouteTableAssociations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := associations.List(client, instanceId, routeTableId, associations.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp))
	for i, association := range resp {
		result[i] = map[string]interface{}{
			"attachment_id":   association.AttachmentId,
			"id":              association.ID,
			"attachment_type": association.ResourceType,
		}
	}
	return result, nil
}

func queryRouteTablePropagations(client *golangsdk.ServiceClient, instanceId, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := propagations.List(client, instanceId, routeTableId, propagations.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp))
	for i, propagation := range resp {
		result[i] = map[string]interface{}{
			"attachment_id":   propagation.AttachmentId,
			"id":              propagation.ID,
			"attachment_type": propagation.ResourceType,
		}
	}
	return result, nil
}

func queryRouteTableRoutes(client *golangsdk.ServiceClient, routeTableId string) ([]map[string]interface{},
	error) {
	resp, err := routes.List(client, routeTableId, routes.ListOpts{})
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(resp))
	for i, route := range resp {
		rr := map[string]interface{}{
			"destination":  route.Destination,
			"is_blackhole": route.IsBlackHole,
			"id":           route.ID,
			"status":       route.Status,
		}
		if len(route.Attachments) < 1 {
			result[i] = rr
			continue
		}

		attachments := make([]map[string]interface{}, len(route.Attachments))
		for i, attachment := range route.Attachments {
			attachments[i] = map[string]interface{}{
				"attachment_id":   attachment.AttachmentId,
				"attachment_type": attachment.ResourceType,
				"resource_id":     attachment.ResourceId,
			}
		}
		rr["attachments"] = attachments
		result[i] = rr
	}
	return result, nil
}

func filterRouteTablesByTags(d *schema.ResourceData, all []routetables.RouteTable) ([]routetables.RouteTable, error) {
	filter := map[string]interface{}{
		"ID":   d.Get("route_table_id"),
		"Name": d.Get("name"),
	}
	filterResult, err := utils.FilterSliceWithField(all, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting security groups list: %s", err)
	}

	tagFilter := d.Get("tags").(map[string]interface{})
	result := make([]routetables.RouteTable, 0, len(filterResult))
	for _, val := range filterResult {
		item := val.(routetables.RouteTable)
		tagmap := utils.TagsToMap(item.Tags)

		if !utils.HasMapContains(tagmap, tagFilter) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func flattenRouteTables(client *golangsdk.ServiceClient, instanceId string,
	all []routetables.RouteTable) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, routeTable := range all {
		routeTableId := routeTable.ID

		associationList, _ := queryRouteTableAssociations(client, instanceId, routeTableId)
		propagationList, _ := queryRouteTablePropagations(client, instanceId, routeTableId)
		routeList, _ := queryRouteTableRoutes(client, routeTableId)

		result[i] = map[string]interface{}{
			"id":                     routeTableId,
			"name":                   routeTable.Name,
			"description":            routeTable.Description,
			"associations":           associationList,
			"propagations":           propagationList,
			"routes":                 routeList,
			"is_default_association": routeTable.IsDefaultAssociation,
			"is_default_propagation": routeTable.IsDefaultPropagation,
			"status":                 routeTable.Status,
			// The time results are not the time in RF3339 format without milliseconds.
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(routeTable.CreatedAt)/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(routeTable.UpdatedAt)/1000, false),
			"tags":       utils.TagsToMap(routeTable.Tags),
		}
	}
	return result
}

func dataSourceRouteTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := routetables.List(client, instanceId, routetables.ListOpts{})
	if err != nil {
		return diag.Errorf("error retrieving route tables: %s", err)
	}
	filterResult, err := filterRouteTablesByTags(d, resp)
	if err != nil {
		return diag.Errorf("error retrieving route tables: %s", err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("route_tables", flattenRouteTables(client, instanceId, filterResult)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving route table list field: %s", mErr)
	}
	return nil
}
