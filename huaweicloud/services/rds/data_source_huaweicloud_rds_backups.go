// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/backups
func DataSourceRdsBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: rdsBackupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Instance ID.`,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Backup ID.`,
			},
			"backup_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Backup type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"auto", "manual", "fragment", "incremental",
				}, false),
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Start time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `End time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Backup name.`,
			},
			"backups": {
				Type:        schema.TypeList,
				Elem:        BackupsBackupSchema(),
				Computed:    true,
				Description: `Backup list. For details, see Data structure of the Backup field.`,
			},
		},
	}
}

func BackupsBackupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup ID.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `RDS instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup type.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Backup size in KB.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup status.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"associated_with_ddm": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a DDM instance has been associated.`,
			},
			"datastore": {
				Type:     schema.TypeList,
				Elem:     BackupsBackupDatastoreSchema(),
				Computed: true,
			},
			"databases": {
				Type:        schema.TypeList,
				Elem:        BackupsBackupDatabasesSchema(),
				Computed:    true,
				Description: `Database been backed up.`,
			},
		},
	}
	return &sc
}

func BackupsBackupDatastoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `DB engine.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `DB engine version.`,
			},
		},
	}
	return &sc
}

func BackupsBackupDatabasesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Database to be backed up for Microsoft SQL Server.`,
			},
		},
	}
	return &sc
}

func rdsBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listBackups: Query the List of RDS Backups.
	var (
		listBackupsHttpUrl = "v3/{project_id}/backups"
		listBackupsProduct = "rds"
	)
	listBackupsClient, err := cfg.NewServiceClient(listBackupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listBackupsPath := listBackupsClient.Endpoint + listBackupsHttpUrl
	listBackupsPath = strings.ReplaceAll(listBackupsPath, "{project_id}", listBackupsClient.ProjectID)

	listBackupsQueryParams := buildListBackupsQueryParams(d)
	listBackupsPath += listBackupsQueryParams

	listBackupsResp, err := pagination.ListAllItems(
		listBackupsClient,
		"offset",
		listBackupsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Backup")
	}

	listBackupsRespJson, err := json.Marshal(listBackupsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listBackupsRespBody interface{}
	err = json.Unmarshal(listBackupsRespJson, &listBackupsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("backups", filterListBackupsBodyBackup(flattenListBackupsBodyBackup(listBackupsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListBackupsBodyBackup(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("backups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"instance_id":         utils.PathSearch("instance_id", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"type":                utils.PathSearch("type", v, nil),
			"size":                utils.PathSearch("size", v, nil),
			"status":              utils.PathSearch("status", v, nil),
			"begin_time":          utils.PathSearch("begin_time", v, nil),
			"end_time":            utils.PathSearch("end_time", v, nil),
			"associated_with_ddm": utils.PathSearch("associated_with_ddm", v, nil),
			"datastore":           flattenBackupDatastore(v),
			"databases":           flattenBackupDatabases(v),
		})
	}
	return rst
}

func flattenBackupDatastore(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("datastore", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("type", curJson, nil),
			"version": utils.PathSearch("version", curJson, nil),
		},
	}
	return rst
}

func flattenBackupDatabases(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func filterListBackupsBodyBackup(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok && param != utils.PathSearch("name", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListBackupsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("backup_id"); ok {
		res = fmt.Sprintf("%s&backup_id=%v", res, v)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if v, ok := d.GetOk("backup_type"); ok {
		res = fmt.Sprintf("%s&backup_type=%v", res, v)
	}

	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%v", res, v)
	}

	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
