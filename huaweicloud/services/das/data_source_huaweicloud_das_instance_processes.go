package das

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

// @API DAS GET /v3/{project_id}/instances/{instance_id}/processes
func DataSourceInstanceProcesses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceProcessesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the instance processes are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the processes belong.`,
			},
			"db_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the database user.`,
			},

			// Optional parameters.
			"db_user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database user.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the database.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the instance node.`,
			},

			// Attributes.
			"processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of processes that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the process.`,
						},
						"db_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database user.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The host of the process.`,
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database.`,
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The command being executed.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The duration of the process, in seconds.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The state of the process.`,
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL statement being executed.`,
						},
						"trx_executed_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The duration of the transaction, in seconds.`,
						},
					},
				},
			},
		},
	}
}

func buildInstanceProcessesQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("&db_user_id=%v", d.Get("db_user_id"))

	if v, ok := d.GetOk("db_user_name"); ok {
		res = fmt.Sprintf("%s&user=%v", res, v)
	}
	if v, ok := d.GetOk("db_name"); ok {
		res = fmt.Sprintf("%s&database=%v", res, v)
	}
	if v, ok := d.GetOk("node_id"); ok {
		res = fmt.Sprintf("%s&node_id=%v", res, v)
	}

	return res
}

func listInstanceProcesses(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/processes?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildInstanceProcessesQueryParams(d)

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
		processes := utils.PathSearch("processes", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, processes...)
		if len(processes) < limit {
			break
		}
		offset += len(processes)
	}

	return result, nil
}

func flattenInstanceProcesses(processes []interface{}) []map[string]interface{} {
	if len(processes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(processes))
	for _, process := range processes {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", process, nil),
			"db_user_name":      utils.PathSearch("user", process, nil),
			"host":              utils.PathSearch("host", process, nil),
			"db_name":           utils.PathSearch("database", process, nil),
			"command":           utils.PathSearch("command", process, nil),
			"time":              utils.PathSearch("time", process, nil),
			"state":             utils.PathSearch("state", process, nil),
			"sql":               utils.PathSearch("sql", process, nil),
			"trx_executed_time": utils.PathSearch("trx_executed_time", process, nil),
		})
	}
	return result
}

func dataSourceInstanceProcessesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	processes, err := listInstanceProcesses(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS instance processes: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("processes", flattenInstanceProcesses(processes)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
