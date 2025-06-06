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

// @API CBR GET /v3/{project_id}/backups
func DataSourceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region where the CBR backups are located.`,
			},
			"checkpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the restore point ID.`,
			},
			"dec": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies the dedicated cloud tag, which only takes effect in dedicated cloud scenarios.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time when the backup ends.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup type. The value can be backup or replication.`,
			},
			"incremental": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether incremental backup is used.`,
			},
			"member_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup sharing status.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup name.`,
			},
			"own_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the owning type of backup. private backups are queried by default.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parent backup ID.`,
			},
			"resource_az": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource availability zones.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource ID.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource name.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type.`,
			},
			"show_replication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to show replication records.`,
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort key.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time when the backup starts.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status.`,
			},
			"used_percent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the using percent of the occupied vault capacity.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the vault ID.`,
			},
			"backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The backup list.`,
				Elem:        dataBackupsSchema(),
			},
		},
	}
}

func dataBackupsSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"checkpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restore point ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup description.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time.`,
			},
			"extend_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The extended information.`,
				Elem:        dataExtendInfoSchema(),
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup ID.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup name.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent backup ID.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID.`,
			},
			"protected_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup time.`,
			},
			"resource_az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource availability zone.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource ID.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name.`,
			},
			"resource_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The resource size, in GB.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup status.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The vault ID.`,
			},
			"replication_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The replication record.`,
				Elem:        dataReplicationRecordsSchema(),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID.`,
			},
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup provider ID.`,
			},
			// Because of the circular dependency problem, the type of this field is defined as a json format string.
			"children": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The children backup list.`,
			},
			"incremental": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether incremental backup is used.`,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The backup snapshot type.`,
			},
		},
	}

	return sc
}

func dataReplicationRecordsSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the replication.`,
			},
			"destination_backup_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the destination backup used for replication.`,
			},
			"destination_checkpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The record ID of the destination backup used for replication.`,
			},
			"destination_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the replication destination project.`,
			},
			"destination_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication destination region.`,
			},
			"destination_vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination vault ID.`,
			},
			"extra_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The additional information of the replication.`,
				Elem:        dataExtraInfoSchema(),
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication record ID.`,
			},
			"source_backup_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the source backup used for replication.`,
			},
			"source_checkpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the source backup record used for replication.`,
			},
			"source_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the replication source project.`,
			},
			"source_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication source region.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The replication status.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the vault where the backup resides.`,
			},
		},
	}

	return sc
}

func dataExtraInfoSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The replication progress.`,
			},
			"fail_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error code. This field is empty if the operation is successful.`,
			},
			"fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error cause.`,
			},
			"auto_trigger": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether replication is automatically scheduled.`,
			},
			"destinatio_vault_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination vault ID.`,
			},
		},
	}

	return sc
}

func dataExtendInfoSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_trigger": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the backup is automatically generated.`,
			},
			"bootable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the backup is a system disk backup.`,
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Snapshot ID of the disk backup.`,
			},
			"support_lld": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to allow lazy loading for fast restoration.`,
			},
			"supported_restore_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restoration mode.`,
			},
			"os_images_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The ID list of images created using backups.`,
				Elem:        dataOsImagesDataSchema(),
			},
			"contain_system_disk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the VM backup data contains system disk data.`,
			},
			"encrypted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the backup is encrypted.`,
			},
			"system_disk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the disk is a system disk.`,
			},
			"is_multi_az": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether multi-AZ backup redundancy is used.`,
			},
		},
	}

	return sc
}

func dataOsImagesDataSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image ID.`,
			},
		},
	}

	return sc
}

func buildDatasourceBackupQueryParams(d *schema.ResourceData, offset int) string {
	var (
		res              string
		normalParameters = []string{
			"checkpoint_id",
			"end_time",
			"enterprise_project_id",
			"image_type",
			"member_status",
			"name",
			"own_type",
			"parent_id",
			"resource_az",
			"resource_id",
			"resource_name",
			"resource_type",
			"sort",
			"start_time",
			"used_percent",
			"vault_id",
		}
	)

	for _, param := range normalParameters {
		if v, ok := d.GetOk(param); ok {
			res = fmt.Sprintf("%s&%s=%v", res, param, v)
		}
	}

	if d.Get("dec").(bool) {
		res = fmt.Sprintf("%s&dec=true", res)
	}

	if d.Get("incremental").(bool) {
		res = fmt.Sprintf("%s&incremental=true", res)
	}

	if d.Get("show_replication").(bool) {
		res = fmt.Sprintf("%s&show_replication=true", res)
	}

	if v, ok := d.GetOk("status"); ok {
		arr := strings.Split(v.(string), ",")
		for _, statusValue := range arr {
			res = fmt.Sprintf("%s&status=%v", res, statusValue)
		}
	}

	if offset != 0 {
		res = fmt.Sprintf("%s&offset=%v", res, offset)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func dataSourceBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/{project_id}/backups"
		product      = "cbr"
		totalBackups []interface{}
		offset       = 0
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		eachRequestPath := requestPath + buildDatasourceBackupQueryParams(d, offset)
		resp, err := client.Request("GET", eachRequestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CBR backups: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		backups := utils.PathSearch("backups", respBody, make([]interface{}, 0)).([]interface{})
		if len(backups) == 0 {
			break
		}

		totalBackups = append(totalBackups, backups...)
		offset += len(backups)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("backups", flattenDataBackups(totalBackups)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataBackups(totalBackups []interface{}) []interface{} {
	if len(totalBackups) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(totalBackups))
	for _, v := range totalBackups {
		result = append(result, map[string]interface{}{
			"checkpoint_id":         utils.PathSearch("checkpoint_id", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"expired_at":            utils.PathSearch("expired_at", v, nil),
			"extend_info":           flattenExtendInfoAttribute(v),
			"id":                    utils.PathSearch("id", v, nil),
			"image_type":            utils.PathSearch("image_type", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"parent_id":             utils.PathSearch("parent_id", v, nil),
			"project_id":            utils.PathSearch("project_id", v, nil),
			"protected_at":          utils.PathSearch("protected_at", v, nil),
			"resource_az":           utils.PathSearch("resource_az", v, nil),
			"resource_id":           utils.PathSearch("resource_id", v, nil),
			"resource_name":         utils.PathSearch("resource_name", v, nil),
			"resource_size":         utils.PathSearch("resource_size", v, nil),
			"resource_type":         utils.PathSearch("resource_type", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
			"vault_id":              utils.PathSearch("vault_id", v, nil),
			"replication_records":   flattenReplicationRecordsAttribute(v),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"provider_id":           utils.PathSearch("provider_id", v, nil),
			"children":              flattenChildrenAttribute(v),
			"incremental":           utils.PathSearch("incremental", v, nil),
			"version":               utils.PathSearch("version", v, nil),
		})
	}
	return result
}

func flattenChildrenAttribute(respBody interface{}) string {
	rawArray := utils.PathSearch("children", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return ""
	}

	return utils.MarshalValue(rawArray)
}

func flattenReplicationRecordsAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("replication_records", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		result = append(result, map[string]interface{}{
			"created_at":                utils.PathSearch("created_at", v, nil),
			"destination_backup_id":     utils.PathSearch("destination_backup_id", v, nil),
			"destination_checkpoint_id": utils.PathSearch("destination_checkpoint_id", v, nil),
			"destination_project_id":    utils.PathSearch("destination_project_id", v, nil),
			"destination_region":        utils.PathSearch("destination_region", v, nil),
			"destination_vault_id":      utils.PathSearch("destination_vault_id", v, nil),
			"extra_info":                flattenExtraInfoAttribute(v),
			"id":                        utils.PathSearch("id", v, nil),
			"source_backup_id":          utils.PathSearch("source_backup_id", v, nil),
			"source_checkpoint_id":      utils.PathSearch("source_checkpoint_id", v, nil),
			"source_project_id":         utils.PathSearch("source_project_id", v, nil),
			"source_region":             utils.PathSearch("source_region", v, nil),
			"status":                    utils.PathSearch("status", v, nil),
			"vault_id":                  utils.PathSearch("vault_id", v, nil),
		})
	}

	return result
}

func flattenExtraInfoAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("extra_info", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"progress":            utils.PathSearch("progress", rawMap, nil),
		"fail_code":           utils.PathSearch("fail_code", rawMap, nil),
		"fail_reason":         utils.PathSearch("fail_reason", rawMap, nil),
		"auto_trigger":        utils.PathSearch("auto_trigger", rawMap, nil),
		"destinatio_vault_id": utils.PathSearch("destinatio_vault_id", rawMap, nil),
	}

	return []interface{}{rstMap}
}

func flattenExtendInfoAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("extend_info", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"auto_trigger":           utils.PathSearch("auto_trigger", rawMap, nil),
		"bootable":               utils.PathSearch("bootable", rawMap, nil),
		"snapshot_id":            utils.PathSearch("snapshot_id", rawMap, nil),
		"support_lld":            utils.PathSearch("support_lld", rawMap, nil),
		"supported_restore_mode": utils.PathSearch("supported_restore_mode", rawMap, nil),
		"os_images_data":         flattenOsImagesDataAttribute(rawMap),
		"contain_system_disk":    utils.PathSearch("contain_system_disk", rawMap, nil),
		"encrypted":              utils.PathSearch("encrypted", rawMap, nil),
		"system_disk":            utils.PathSearch("system_disk", rawMap, nil),
		"is_multi_az":            utils.PathSearch("is_multi_az", rawMap, nil),
	}

	return []interface{}{rstMap}
}

func flattenOsImagesDataAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("os_images_data", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		result = append(result, map[string]interface{}{
			"image_id": utils.PathSearch("image_id", v, nil),
		})
	}

	return result
}
