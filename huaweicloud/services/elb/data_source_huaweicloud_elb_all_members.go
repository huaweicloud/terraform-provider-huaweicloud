// Generated by PMS #570
package elb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceElbAllMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbAllMembersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the backend server name.`,
			},
			"weight": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `Specifies the weight of the backend server.`,
			},
			"subnet_cidr_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the ID of the subnet where the backend server works.`,
			},
			"address": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IP address of the backend server.`,
			},
			"protocol_port": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `Specifies the port used by the backend servers.`,
			},
			"member_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the backend server ID.`,
			},
			"operating_status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the operating status of the backend server.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the ID of the enterprise project.`,
			},
			"ip_version": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IP address version supported by the backend server group.`,
			},
			"pool_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the ID of the backend server group to which the backend server belongs.`,
			},
			"loadbalancer_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the ID of the load balancer with which the load balancer is associated.`,
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of backend servers.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the backend server ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the backend server name.`,
						},
						"member_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the type of the backend server.`,
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the private IP address bound to the backend server.`,
						},
						"subnet_cidr_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the IPv4 or IPv6 subnet where the backend server resides.`,
						},
						"protocol_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the port used by the backend server to receive requests.`,
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the weight of the backend server.`,
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the health status of the backend server.`,
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP version supported by the backend server.`,
						},
						"pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the backend server group to which the backend server belongs.`,
						},
						"loadbalancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the load balancer with which the backend server is associated.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the project where the backend server is used.`,
						},
						"status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the health status of the backend server.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the listener ID.`,
									},
									"operating_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the health status of the backend server.`,
									},
									"reason": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates why health check fails.`,
										Elem:        membersStatusReasonElem(),
									},
								},
							},
						},
						"reason": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates why health check fails.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reason_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the code of the health check failures.`,
									},
									"expected_response": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the expected HTTP status code.`,
									},
									"healthcheck_response": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the returned HTTP status code in the response.`,
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when a backend server was added.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when a backend server was updated.`,
						},
					},
				},
			},
		},
	}
}

// membersStatusReasonElem
// The Elem of "members.status.reason"
func membersStatusReasonElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"reason_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code of the health check failures.`,
			},
			"expected_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expected HTTP status code.`,
			},
			"healthcheck_response": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the returned HTTP status code in the response.`,
			},
		},
	}
}

type AllMembersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAllMembersDSWrapper(d *schema.ResourceData, meta interface{}) *AllMembersDSWrapper {
	return &AllMembersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceElbAllMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAllMembersDSWrapper(d, meta)
	listAllMembersRst, err := wrapper.ListAllMembers()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAllMembersToSchema(listAllMembersRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API ELB GET /v3/{project_id}/elb/members
func (w *AllMembersDSWrapper) ListAllMembers() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "elb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/elb/members"
	params := map[string]any{
		"name":                  w.ListToArray("name"),
		"weight":                w.ListToArray("weight"),
		"subnet_cidr_id":        w.ListToArray("subnet_cidr_id"),
		"address":               w.ListToArray("address"),
		"protocol_port":         w.ListToArray("protocol_port"),
		"id":                    w.ListToArray("member_id"),
		"operating_status":      w.ListToArray("operating_status"),
		"enterprise_project_id": w.ListToArray("enterprise_project_id"),
		"ip_version":            w.ListToArray("ip_version"),
		"pool_id":               w.ListToArray("pool_id"),
		"loadbalancer_id":       w.ListToArray("loadbalancer_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("members", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *AllMembersDSWrapper) listAllMembersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("members", schemas.SliceToList(body.Get("members"),
			func(members gjson.Result) any {
				return map[string]any{
					"id":               members.Get("id").Value(),
					"name":             members.Get("name").Value(),
					"member_type":      members.Get("member_type").Value(),
					"address":          members.Get("address").Value(),
					"subnet_cidr_id":   members.Get("subnet_cidr_id").Value(),
					"protocol_port":    members.Get("protocol_port").Value(),
					"weight":           members.Get("weight").Value(),
					"operating_status": members.Get("operating_status").Value(),
					"ip_version":       members.Get("ip_version").Value(),
					"pool_id":          members.Get("pool_id").Value(),
					"loadbalancer_id":  members.Get("loadbalancer_id").Value(),
					"project_id":       members.Get("project_id").Value(),
					"status": schemas.SliceToList(members.Get("status"),
						func(status gjson.Result) any {
							return map[string]any{
								"listener_id":      status.Get("listener_id").Value(),
								"operating_status": status.Get("operating_status").Value(),
								"reason":           w.setMembersStatusReason(status),
							}
						},
					),
					"reason": schemas.SliceToList(members.Get("reason"),
						func(reason gjson.Result) any {
							return map[string]any{
								"reason_code":          reason.Get("reason_code").Value(),
								"expected_response":    reason.Get("expected_response").Value(),
								"healthcheck_response": reason.Get("healthcheck_response").Value(),
							}
						},
					),
					"created_at": members.Get("created_at").Value(),
					"updated_at": members.Get("updated_at").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*AllMembersDSWrapper) setMembersStatusReason(status gjson.Result) any {
	return schemas.SliceToList(status.Get("reason"), func(reason gjson.Result) any {
		return map[string]any{
			"reason_code":          reason.Get("reason_code").Value(),
			"expected_response":    reason.Get("expected_response").Value(),
			"healthcheck_response": reason.Get("healthcheck_response").Value(),
		}
	})
}
