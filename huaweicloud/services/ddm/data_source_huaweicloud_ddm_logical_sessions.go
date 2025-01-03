// Generated by PMS #430
package ddm

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDdmLogicalSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdmLogicalSessionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DDM instance ID.`,
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the keyword filtered by the session result.`,
			},
			"logical_processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the logical sessions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the session ID`,
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the current user.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address and port number.`,
						},
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the database name.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status of the SQL statement.`,
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the connection status.`,
						},
						"info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SQL statement that is being executed.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the duration of a connection, in seconds.`,
						},
					},
				},
			},
		},
	}
}

type LogicalSessionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLogicalSessionsDSWrapper(d *schema.ResourceData, meta interface{}) *LogicalSessionsDSWrapper {
	return &LogicalSessionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDdmLogicalSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLogicalSessionsDSWrapper(d, meta)
	shoLogProRst, err := wrapper.ShowLogicalProcesses()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showLogicalProcessesToSchema(shoLogProRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DDM GET /v3/{project_id}/instances/{instance_id}/logical-processes
func (w *LogicalSessionsDSWrapper) ShowLogicalProcesses() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ddm")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/logical-processes"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"keyword": w.Get("keyword"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("logical_processes", "offset", "limit", 0).
		Request().
		Result()
}

func (w *LogicalSessionsDSWrapper) showLogicalProcessesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("logical_processes", schemas.SliceToList(body.Get("logical_processes"),
			func(logPro gjson.Result) any {
				return map[string]any{
					"id":      logPro.Get("id").Value(),
					"user":    logPro.Get("user").Value(),
					"host":    logPro.Get("host").Value(),
					"db":      logPro.Get("db").Value(),
					"state":   logPro.Get("state").Value(),
					"command": logPro.Get("command").Value(),
					"info":    logPro.Get("info").Value(),
					"time":    logPro.Get("time").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
