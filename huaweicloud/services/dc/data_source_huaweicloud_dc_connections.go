// Generated by PMS #210
package dc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDcConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the connection.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the connection.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the connection.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the connection.`,
			},
			"hosting_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies operations connection ID by which hosted connections are filtered.`,
			},
			"port_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the port used by the connection.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the enterprise project to which the connections belong.`,
			},
			"direct_connects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All connections that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the connection.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the connection.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the connection.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the connection.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the connection.`,
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The connection bandwidth, in Mbit/s.`,
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access location information of the DC.`,
						},
						"port_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the port used by the connection.`,
						},
						"provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The line carrier of the connection.`,
						},
						"provider_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the carrier's leased line.`,
						},
						"support_feature": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Lists the features supported by the connection.`,
						},
						"vgw_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The gateway type of the DC.`,
						},
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The VLAN allocated to the hosted connection.`,
						},
						"hosting_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the operations connection on which the hosted connection is created.`,
						},
						"device_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the device connected to the connection.`,
						},
						"lag_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the LAG to which the connection belongs.`,
						},
						"ies_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of an IES edge site.`,
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The billing mode.`,
						},
						"peer_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The location of the on-premises facility at the other end of the connection.`,
						},
						"peer_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The carrier connected to the connection.`,
						},
						"peer_port_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The peer port type.`,
						},
						"public_border_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The public border group of the AZ, indicating whether the site is a HomeZones site.`,
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The customer email information.`,
						},
						"onestopdc_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of a full-service connection.`,
						},
						"modified_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The new bandwidth after the line bandwidth is changed.`,
						},
						"change_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a renewal change.`,
						},
						"ratio_95peak": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The percentage of the minimum bandwidth for 95th percentile billing.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the connection belongs.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the connection.`,
						},
						"apply_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application time of the connection, in RFC3339 format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the connection, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

type ConnectionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newConnectionsDSWrapper(d *schema.ResourceData, meta interface{}) *ConnectionsDSWrapper {
	return &ConnectionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDcConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newConnectionsDSWrapper(d, meta)
	lisDirConRst, err := wrapper.ListDirectConnects()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listDirectConnectsToSchema(lisDirConRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DC GET /v3/{project_id}/dcaas/direct-connects
func (w *ConnectionsDSWrapper) ListDirectConnects() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dc")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/dcaas/direct-connects"
	params := map[string]any{
		"hosting_id":            w.Get("hosting_id"),
		"enterprise_project_id": w.Get("enterprise_project_id"),
		"id":                    w.Get("connection_id"),
		"name":                  w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("direct_connects", "page_info.next_marker", "marker").
		Filter(
			filters.New().From("direct_connects").
				Where("type", "=", w.Get("type")).
				Where("status", "=", w.Get("status")).
				Where("port_type", "=", w.Get("port_type")),
		).
		Request().
		Result()
}

func (w *ConnectionsDSWrapper) listDirectConnectsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("direct_connects", schemas.SliceToList(body.Get("direct_connects"),
			func(directConnects gjson.Result) any {
				return map[string]any{
					"id":                    directConnects.Get("id").Value(),
					"name":                  directConnects.Get("name").Value(),
					"type":                  directConnects.Get("type").Value(),
					"status":                directConnects.Get("status").Value(),
					"description":           directConnects.Get("description").Value(),
					"bandwidth":             directConnects.Get("bandwidth").Value(),
					"location":              directConnects.Get("location").Value(),
					"port_type":             directConnects.Get("port_type").Value(),
					"provider":              directConnects.Get("provider").Value(),
					"provider_status":       directConnects.Get("provider_status").Value(),
					"support_feature":       schemas.SliceToStrList(directConnects.Get("support_feature")),
					"vgw_type":              directConnects.Get("vgw_type").Value(),
					"vlan":                  directConnects.Get("vlan").Value(),
					"hosting_id":            directConnects.Get("hosting_id").Value(),
					"device_id":             directConnects.Get("device_id").Value(),
					"lag_id":                directConnects.Get("lag_id").Value(),
					"ies_id":                directConnects.Get("ies_id").Value(),
					"charge_mode":           directConnects.Get("charge_mode").Value(),
					"peer_location":         directConnects.Get("peer_location").Value(),
					"peer_provider":         directConnects.Get("peer_provider").Value(),
					"peer_port_type":        directConnects.Get("peer_port_type").Value(),
					"public_border_group":   directConnects.Get("public_border_group").Value(),
					"email":                 directConnects.Get("email").Value(),
					"onestopdc_status":      directConnects.Get("onestopdc_status").Value(),
					"modified_bandwidth":    directConnects.Get("modified_bandwidth").Value(),
					"change_mode":           directConnects.Get("change_mode").Value(),
					"ratio_95peak":          directConnects.Get("ratio_95peak").Value(),
					"enterprise_project_id": directConnects.Get("enterprise_project_id").Value(),
					"tags":                  w.setDirectConnectsTags(directConnects),
					"apply_time":            w.setDirConAppTime(directConnects),
					"created_at":            w.setDirConCreTime(directConnects),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*ConnectionsDSWrapper) setDirectConnectsTags(data gjson.Result) map[string]string {
	tags := make(map[string]string)
	tagList := data.Get("tags").Array()
	for _, v := range tagList {
		tags[v.Get("key").String()] = v.Get("value").String()
	}
	return tags
}

func (*ConnectionsDSWrapper) setDirConAppTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("apply_time").String(), "2006-01-02T15:04:05.000Z")/1000, false)
}

func (*ConnectionsDSWrapper) setDirConCreTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("create_time").String(), "2006-01-02T15:04:05.000Z")/1000, false)
}