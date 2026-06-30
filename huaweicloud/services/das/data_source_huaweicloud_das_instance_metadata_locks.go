package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/instances/{instance_id}/metadata-locks
func DataSourceInstanceMetadataLocks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceMetadataLocksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the instance metadata locks are located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance ID.",
			},
			"db_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database user ID.",
			},

			// Optional parameters.
			"thread_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The session ID for filtering.",
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The database name for filtering.",
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The table name for filtering.",
			},

			// Attributes.
			"metadata_locks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of metadata locks that matched filter parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"thread_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session ID of the thread holding or waiting for the metadata lock.",
						},
						"lock_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lock status.",
						},
						"lock_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lock mode.",
						},
						"lock_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lock type.",
						},
						"lock_duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lock duration.",
						},
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database where the lock is held.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the table on which the lock is held.",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database user associated with the session.",
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The session duration.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host from which the session is connected.",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database name of the session.",
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The command being executed by the session.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the session.",
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SQL statement being executed by the session.",
						},
						"trx_exec_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution time of the current transaction.",
						},
						"block_process": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of processes that are blocking the current lock.",
							Elem:        metadataLockProcessSchema(),
						},
						"wait_process": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of processes that are waiting for the current lock.",
							Elem:        metadataLockProcessSchema(),
						},
					},
				},
			},
		},
	}
}

func metadataLockProcessSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session ID of the process.",
			},
			"user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database user associated with the process.",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The connecting host of the process.",
			},
			"database": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name of the process.",
			},
			"command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The command being executed by the process.",
			},
			"time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session duration of the process.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the process.",
			},
			"sql": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SQL statement being executed by the process.",
			},
			"trx_executed_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The transaction duration of the process.",
			},
		},
	}
}

func dataSourceInstanceMetadataLocksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	metadataLocks, err := listInstanceMetadataLocks(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS instance metadata locks: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metadata_locks", flattenInstanceMetadataLocks(metadataLocks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildInstanceMetadataLocksQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?db_user_id=%v", d.Get("db_user_id"))

	if v, ok := d.GetOk("thread_id"); ok {
		res = fmt.Sprintf("%s&thread_id=%v", res, v)
	}
	if v, ok := d.GetOk("database"); ok {
		res = fmt.Sprintf("%s&database=%v", res, v)
	}
	if v, ok := d.GetOk("table"); ok {
		res = fmt.Sprintf("%s&table=%v", res, v)
	}

	return res
}

func listInstanceMetadataLocks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/metadata-locks"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildInstanceMetadataLocksQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("metadata_locks", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenInstanceMetadataLocks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"thread_id":     utils.PathSearch("thread_id", item, nil),
			"lock_status":   utils.PathSearch("lock_status", item, nil),
			"lock_mode":     utils.PathSearch("lock_mode", item, nil),
			"lock_type":     utils.PathSearch("lock_type", item, nil),
			"lock_duration": utils.PathSearch("lock_duration", item, nil),
			"table_schema":  utils.PathSearch("table_schema", item, nil),
			"table_name":    utils.PathSearch("table_name", item, nil),
			"user":          utils.PathSearch("user", item, nil),
			"time":          utils.PathSearch("time", item, nil),
			"host":          utils.PathSearch("host", item, nil),
			"database":      utils.PathSearch("database", item, nil),
			"command":       utils.PathSearch("command", item, nil),
			"state":         utils.PathSearch("state", item, nil),
			"sql":           utils.PathSearch("sql", item, nil),
			"trx_exec_time": utils.PathSearch("trx_exec_time", item, nil),
			"block_process": flattenMetadataLockProcesses(
				utils.PathSearch("block_process", item, make([]interface{}, 0)).([]interface{})),
			"wait_process": flattenMetadataLockProcesses(
				utils.PathSearch("wait_process", item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenMetadataLockProcesses(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", item, nil),
			"user":              utils.PathSearch("user", item, nil),
			"host":              utils.PathSearch("host", item, nil),
			"database":          utils.PathSearch("database", item, nil),
			"command":           utils.PathSearch("command", item, nil),
			"time":              utils.PathSearch("time", item, nil),
			"state":             utils.PathSearch("state", item, nil),
			"sql":               utils.PathSearch("sql", item, nil),
			"trx_executed_time": utils.PathSearch("trx_executed_time", item, nil),
		})
	}

	return result
}
