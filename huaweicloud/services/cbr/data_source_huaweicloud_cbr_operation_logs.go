package cbr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR GET /v3/{project_id}/operation-logs
func DataSourceOperationLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperationLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsSchema(),
			},
		},
	}
}

func dataSourceCbrOperationLogsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"checkpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ended_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsErrorInfoSchema(),
			},
			"extra_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsExtraInfoSchema(),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsErrorInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsExtraInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsBackupSchema(),
			},
			"common": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsCommonSchema(),
			},
			"delete": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsDeleteSchema(),
			},
			"sync": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsSyncSchema(),
			},
			"remove_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsRemoveResourcesSchema(),
			},
			"replication": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsReplicationSchema(),
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsResourceSchema(),
			},
			"restore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsRestoreSchema(),
			},
			"vault_delete": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsVaultDeleteSchema(),
			},
		},
	}
}

func dataSourceCbrOperationLogsBackupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_consistency_error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_consistency_error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app_consistency_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incremental": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsCommonSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsDeleteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsSyncSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sync_backup_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delete_backup_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"err_sync_backup_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsRemoveResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fail_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsResourceSchema(),
			},
		},
	}
}

func dataSourceCbrOperationLogsReplicationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination_backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_checkpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_checkpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"extra_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceCbrOperationLogsResourceExtraInfoSchema(),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsResourceExtraInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCbrOperationLogsRestoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCbrOperationLogsVaultDeleteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fail_delete_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_delete_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildOperationLogsQueryParams(d *schema.ResourceData, epsID string) string {
	queryParams := []string{"limit=1000"}
	if epsID != "" {
		queryParams = append(queryParams, fmt.Sprintf("enterprise_project_id=%v", epsID))
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = append(queryParams, fmt.Sprintf("end_time=%v", v))
	}
	if v, ok := d.GetOk("operation_type"); ok {
		queryParams = append(queryParams, fmt.Sprintf("operation_type=%v", v))
	}
	if v, ok := d.GetOk("provider_id"); ok {
		queryParams = append(queryParams, fmt.Sprintf("provider_id=%v", v))
	}
	if v, ok := d.GetOk("resource_id"); ok {
		queryParams = append(queryParams, fmt.Sprintf("resource_id=%v", v))
	}
	if v, ok := d.GetOk("resource_name"); ok {
		queryParams = append(queryParams, fmt.Sprintf("resource_name=%v", v))
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = append(queryParams, fmt.Sprintf("start_time=%v", v))
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = append(queryParams, fmt.Sprintf("status=%v", v))
	}
	if v, ok := d.GetOk("vault_id"); ok {
		queryParams = append(queryParams, fmt.Sprintf("vault_id=%v", v))
	}
	if v, ok := d.GetOk("vault_name"); ok {
		queryParams = append(queryParams, fmt.Sprintf("vault_name=%v", v))
	}

	return "?" + strings.Join(queryParams, "&")
}

func flattenOperationLogs(logs []interface{}) []map[string]interface{} {
	if len(logs) == 0 {
		return nil
	}

	var rst []map[string]interface{}
	for _, log := range logs {
		item := map[string]interface{}{
			"checkpoint_id":  utils.PathSearch("checkpoint_id", log, nil),
			"created_at":     utils.PathSearch("created_at", log, nil),
			"ended_at":       utils.PathSearch("ended_at", log, nil),
			"error_info":     flattenErrorInfo(utils.PathSearch("error_info", log, nil)),
			"extra_info":     flattenExtraInfo(utils.PathSearch("extra_info", log, nil)),
			"id":             utils.PathSearch("id", log, nil),
			"operation_type": utils.PathSearch("operation_type", log, nil),
			"policy_id":      utils.PathSearch("policy_id", log, nil),
			"provider_id":    utils.PathSearch("provider_id", log, nil),
			"started_at":     utils.PathSearch("started_at", log, nil),
			"status":         utils.PathSearch("status", log, nil),
			"updated_at":     utils.PathSearch("updated_at", log, nil),
			"vault_id":       utils.PathSearch("vault_id", log, nil),
			"vault_name":     utils.PathSearch("vault_name", log, nil),
		}
		rst = append(rst, item)
	}
	return rst
}

func flattenErrorInfo(errorInfo interface{}) []map[string]interface{} {
	if errorInfo == nil {
		return nil
	}
	item := map[string]interface{}{
		"code":    utils.PathSearch("code", errorInfo, nil),
		"message": utils.PathSearch("message", errorInfo, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfo(extraInfo interface{}) []map[string]interface{} {
	if extraInfo == nil {
		return nil
	}
	item := map[string]interface{}{
		"backup":           flattenExtraInfoBackup(utils.PathSearch("backup", extraInfo, nil)),
		"common":           flattenExtraInfoCommon(utils.PathSearch("common", extraInfo, nil)),
		"delete":           flattenExtraInfoDelete(utils.PathSearch("delete", extraInfo, nil)),
		"sync":             flattenExtraInfoSync(utils.PathSearch("sync", extraInfo, nil)),
		"remove_resources": flattenExtraInfoRemoveResources(utils.PathSearch("remove_resources", extraInfo, nil)),
		"replication":      flattenExtraInfoReplication(utils.PathSearch("replication", extraInfo, nil)),
		"resource":         flattenExtraInfoResource(utils.PathSearch("resource", extraInfo, nil)),
		"restore":          flattenExtraInfoRestore(utils.PathSearch("restore", extraInfo, nil)),
		"vault_delete":     flattenExtraInfoVaultDelete(utils.PathSearch("vault_delete", extraInfo, nil)),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoBackup(backup interface{}) []map[string]interface{} {
	if backup == nil {
		return nil
	}
	item := map[string]interface{}{
		"app_consistency_error_code":    utils.PathSearch("app_consistency_error_code", backup, nil),
		"app_consistency_error_message": utils.PathSearch("app_consistency_error_message", backup, nil),
		"app_consistency_status":        utils.PathSearch("app_consistency_status", backup, nil),
		"backup_id":                     utils.PathSearch("backup_id", backup, nil),
		"backup_name":                   utils.PathSearch("backup_name", backup, nil),
		"incremental":                   utils.PathSearch("incremental", backup, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoCommon(common interface{}) []map[string]interface{} {
	if common == nil {
		return nil
	}
	item := map[string]interface{}{
		"progress":   utils.PathSearch("progress", common, nil),
		"request_id": utils.PathSearch("request_id", common, nil),
		"task_id":    utils.PathSearch("task_id", common, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoDelete(extraInfoDelete interface{}) []map[string]interface{} {
	if extraInfoDelete == nil {
		return nil
	}
	item := map[string]interface{}{
		"backup_id":   utils.PathSearch("backup_id", extraInfoDelete, nil),
		"backup_name": utils.PathSearch("backup_name", extraInfoDelete, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoSync(sync interface{}) []map[string]interface{} {
	if sync == nil {
		return nil
	}
	item := map[string]interface{}{
		"sync_backup_num":     utils.PathSearch("sync_backup_num", sync, nil),
		"delete_backup_num":   utils.PathSearch("delete_backup_num", sync, nil),
		"err_sync_backup_num": utils.PathSearch("err_sync_backup_num", sync, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoRemoveResources(removeResources interface{}) []map[string]interface{} {
	if removeResources == nil {
		return nil
	}
	item := map[string]interface{}{
		"fail_count":  utils.PathSearch("fail_count", removeResources, nil),
		"total_count": utils.PathSearch("total_count", removeResources, nil),
		"resources":   flattenExtraInfoResource(utils.PathSearch("resources", removeResources, nil)),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoReplication(replication interface{}) []map[string]interface{} {
	if replication == nil {
		return nil
	}
	item := map[string]interface{}{
		"destination_backup_id":     utils.PathSearch("destination_backup_id", replication, nil),
		"destination_checkpoint_id": utils.PathSearch("destination_checkpoint_id", replication, nil),
		"destination_project_id":    utils.PathSearch("destination_project_id", replication, nil),
		"destination_region":        utils.PathSearch("destination_region", replication, nil),
		"source_backup_id":          utils.PathSearch("source_backup_id", replication, nil),
		"source_checkpoint_id":      utils.PathSearch("source_checkpoint_id", replication, nil),
		"source_project_id":         utils.PathSearch("source_project_id", replication, nil),
		"source_region":             utils.PathSearch("source_region", replication, nil),
		"source_backup_name":        utils.PathSearch("source_backup_name", replication, nil),
		"destination_backup_name":   utils.PathSearch("destination_backup_name", replication, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoResource(resource interface{}) []map[string]interface{} {
	if resource == nil {
		return nil
	}
	item := map[string]interface{}{
		"extra_info": utils.PathSearch("extra_info", resource, nil),
		"id":         utils.PathSearch("id", resource, nil),
		"name":       utils.PathSearch("name", resource, nil),
		"type":       utils.PathSearch("type", resource, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoRestore(restore interface{}) []map[string]interface{} {
	if restore == nil {
		return nil
	}
	item := map[string]interface{}{
		"backup_id":            utils.PathSearch("backup_id", restore, nil),
		"backup_name":          utils.PathSearch("backup_name", restore, nil),
		"target_resource_id":   utils.PathSearch("target_resource_id", restore, nil),
		"target_resource_name": utils.PathSearch("target_resource_name", restore, nil),
	}
	return []map[string]interface{}{item}
}

func flattenExtraInfoVaultDelete(vaultDelete interface{}) []map[string]interface{} {
	if vaultDelete == nil {
		return nil
	}
	item := map[string]interface{}{
		"fail_delete_count":  utils.PathSearch("fail_delete_count", vaultDelete, nil),
		"total_delete_count": utils.PathSearch("total_delete_count", vaultDelete, nil),
	}
	return []map[string]interface{}{item}
}

func dataSourceOperationLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsID   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v3/{project_id}/operation-logs"
		product = "cbr"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOperationLogsQueryParams(d, epsID)

	var allLogs []interface{}
	offset := 0

	for {
		pagedPath := requestPath
		pagedPath += fmt.Sprintf("&offset=%d", offset)

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", pagedPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying CBR operation logs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		logs := utils.PathSearch("operation_logs", respBody, []interface{}{}).([]interface{})
		if len(logs) == 0 {
			break
		}
		allLogs = append(allLogs, logs...)
		offset += len(logs)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("operation_logs", flattenOperationLogs(allLogs)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the CBR operation logs: %s", mErr)
	}
	return nil
}
