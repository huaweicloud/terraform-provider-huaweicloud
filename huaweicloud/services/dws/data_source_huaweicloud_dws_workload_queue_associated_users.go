// Generated by PMS #327
package dws

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDwsWorkloadQueueAssociatedUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDwsWorkloadQueueAssociatedUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workload queue name bound to the users.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All users that associated with the specified workload queue.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user.`,
						},
						"occupy_resource_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of the resources used by the user to run jobs.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource name.`,
									},
									"resource_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The resource value.`,
									},
									"value_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource attribute unit.`,
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

type WorkloadQueueAssociatedUsersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newWorkloadQueueAssociatedUsersDSWrapper(d *schema.ResourceData, meta interface{}) *WorkloadQueueAssociatedUsersDSWrapper {
	return &WorkloadQueueAssociatedUsersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDwsWorkloadQueueAssociatedUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newWorkloadQueueAssociatedUsersDSWrapper(d, meta)
	lisWorQueUseRst, err := wrapper.ListWorkloadQueueUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listWorkloadQueueUsersToSchema(lisWorQueUseRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users
func (w *WorkloadQueueAssociatedUsersDSWrapper) ListWorkloadQueueUsers() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dws")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	uri = strings.ReplaceAll(uri, "{queue_name}", w.Get("queue_name").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("user_list", "offset", "limit", 100).
		Request().
		Result()
}

func (w *WorkloadQueueAssociatedUsersDSWrapper) listWorkloadQueueUsersToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("users", schemas.SliceToList(body.Get("user_list"),
			func(users gjson.Result) any {
				return map[string]any{
					"name": users.Get("user_name").Value(),
					"occupy_resource_list": schemas.SliceToList(users.Get("occupy_resource_list"),
						func(occupyResourceList gjson.Result) any {
							return map[string]any{
								"resource_name":  occupyResourceList.Get("resource_name").Value(),
								"resource_value": occupyResourceList.Get("resource_value").Value(),
								"value_unit":     occupyResourceList.Get("value_unit").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
