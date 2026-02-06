package oms

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS GET /v2/{project_id}/sync-tasks
func DataSourceMigrationSyncTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrationSyncTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync_task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_cloud_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_kms": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_metadata_migration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_restore": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cdn": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"authentication_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"object_overwrite_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dst_storage_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistency_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildMigrationSyncTasksQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	return queryParams
}

func listMigrationSyncTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/sync-tasks?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildMigrationSyncTasksQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		tasks := utils.PathSearch("tasks", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)
		if len(tasks) < limit {
			break
		}

		offset += len(tasks)
	}

	return result, nil
}

func dataSourceMigrationSyncTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	resp, err := listMigrationSyncTasks(client, d)
	if err != nil {
		return diag.Errorf("error retrieving synchronization tasks: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenMigrationSyncTasks(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMigrationSyncTasks(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"sync_task_id":              utils.PathSearch("sync_task_id", v, nil),
			"src_cloud_type":            utils.PathSearch("src_cloud_type", v, nil),
			"src_region":                utils.PathSearch("src_region", v, nil),
			"src_bucket":                utils.PathSearch("src_bucket", v, nil),
			"dst_bucket":                utils.PathSearch("dst_bucket", v, nil),
			"dst_region":                utils.PathSearch("dst_region", v, nil),
			"description":               utils.PathSearch("description", v, nil),
			"status":                    utils.PathSearch("status", v, nil),
			"enable_kms":                utils.PathSearch("enable_kms", v, nil),
			"enable_metadata_migration": utils.PathSearch("enable_metadata_migration", v, nil),
			"enable_restore":            utils.PathSearch("enable_restore", v, nil),
			"app_id":                    utils.PathSearch("app_id", v, nil),
			"source_cdn":                flattenSourceCdnResp(utils.PathSearch("source_cdn", v, nil)),
			"object_overwrite_mode":     utils.PathSearch("object_overwrite_mode", v, nil),
			"dst_storage_policy":        utils.PathSearch("dst_storage_policy", v, nil),
			"consistency_check":         utils.PathSearch("consistency_check", v, nil),
			"create_time":               utils.PathSearch("create_time", v, nil),
			"last_start_time":           utils.PathSearch("last_start_time", v, nil),
		})
	}

	return result
}

func flattenSourceCdnResp(sourceCdnResp interface{}) []map[string]interface{} {
	if sourceCdnResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"domain":              utils.PathSearch("domain", sourceCdnResp, nil),
		"protocol":            utils.PathSearch("protocol", sourceCdnResp, nil),
		"authentication_type": utils.PathSearch("authentication_type", sourceCdnResp, nil),
	}

	return []map[string]interface{}{result}
}
