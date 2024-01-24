package dcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/backups
func DataSourceBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsBackupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the DCS instance.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time (UTC) of DCS backups.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time (UTC) of DCS backups.`,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the DCS instance backup.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the backup status.`,
			},
			"is_support_restore": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether restoration is supported.`,
			},
			"backup_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the format of the DCS instance backup.`,
			},
			"backups": {
				Type:        schema.TypeList,
				Elem:        backupRecordResponseSchema(),
				Computed:    true,
				Description: `Indicates the list of backup records.`,
			},
		},
	}
}

func backupRecordResponseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the DCS instance backup.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup name.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the ID of the DCS instance.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the size of the backup file (byte).`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup type.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the backup task is created.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which DCS instance backup is completed.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup status.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the DCS instance backup.`,
			},
			"is_support_restore": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether restoration is supported.`,
			},
			"backup_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the format of the DCS instance backup.`,
			},
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error code displayed for a backup failure.`,
			},
			"progress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup progress.`,
			},
		},
	}
	return &sc
}

func dataSourceDcsBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listBackupsHttpUrl = "v2/{project_id}/instances/{instance_id}/backups"
		listBackupsProduct = "dcs"
	)
	listBackupsClient, err := cfg.NewServiceClient(listBackupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	listBackupsBasePath := listBackupsClient.Endpoint + listBackupsHttpUrl
	listBackupsBasePath = strings.ReplaceAll(listBackupsBasePath, "{project_id}", listBackupsClient.ProjectID)
	listBackupsBasePath = strings.ReplaceAll(listBackupsBasePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getBackupsSchemasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 204,
		},
	}

	var currentTotal int
	var listBackupsPath string
	resultArr := make([]interface{}, 0)

	for {
		listBackupsPath = listBackupsBasePath + buildListBackupsQueryParams(d, currentTotal)
		listBackupsResp, err := listBackupsClient.Request("GET", listBackupsPath, &getBackupsSchemasOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving DCS backups")
		}

		listBackupsRespBody, err := utils.FlattenResponse(listBackupsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		backups := utils.PathSearch("backup_record_response", listBackupsRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("total_num", listBackupsRespBody, 0)
		resultArr = append(resultArr, backups...)
		currentTotal += len(backups)
		if float64(currentTotal) == total || total == 0 {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("backups", flattenListBackupsBody(filterListBackups(d, resultArr))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListBackupsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("backup_id", v, nil),
			"name":               utils.PathSearch("backup_name", v, nil),
			"instance_id":        utils.PathSearch("instance_id", v, nil),
			"size":               utils.PathSearch("size", v, nil),
			"type":               utils.PathSearch("backup_type", v, nil),
			"begin_time":         utils.PathSearch("created_at", v, nil),
			"end_time":           utils.PathSearch("updated_at", v, nil),
			"description":        utils.PathSearch("remark", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"is_support_restore": utils.PathSearch("is_support_restore", v, nil),
			"backup_format":      utils.PathSearch("backup_format", v, nil),
			"error_code":         utils.PathSearch("error_code", v, nil),
			"progress":           utils.PathSearch("progress", v, nil),
		})
	}
	return rst
}

func filterListBackups(d *schema.ResourceData, backupRecordArray []interface{}) []interface{} {
	if len(backupRecordArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(backupRecordArray))

	rawId, rawIdOK := d.GetOk("backup_id")
	rawName, rawNameOK := d.GetOk("name")
	rawType, rawTypeOk := d.GetOk("type")
	rawStatus, rawStatusOK := d.GetOk("status")
	rawIsSupportRestore, rawIsSupportRestoreOK := d.GetOk("is_support_restore")
	rawBackupFormat, backupFormatOK := d.GetOk("backup_format")

	for _, backupRecord := range backupRecordArray {
		id := utils.PathSearch("backup_id", backupRecord, nil)
		status := utils.PathSearch("status", backupRecord, nil)
		name := utils.PathSearch("backup_name", backupRecord, nil)
		backupType := utils.PathSearch("backup_type", backupRecord, nil)
		isSupportRestore := utils.PathSearch("is_support_restore", backupRecord, nil)
		backupFormat := utils.PathSearch("backup_format", backupRecord, nil)
		if rawIdOK && rawId != id {
			continue
		}
		if rawNameOK && rawName != name {
			continue
		}
		if rawTypeOk && rawType != backupType {
			continue
		}
		if rawStatusOK && rawStatus != status {
			continue
		}
		if rawIsSupportRestoreOK && rawIsSupportRestore != isSupportRestore {
			continue
		}
		if backupFormatOK && rawBackupFormat != backupFormat {
			continue
		}
		result = append(result, backupRecord)
	}

	return result
}

func buildListBackupsQueryParams(d *schema.ResourceData, offset int) string {
	res := fmt.Sprintf("?limit=10&offset=%v", offset)
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%v", res, v)
	}

	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	return res
}
