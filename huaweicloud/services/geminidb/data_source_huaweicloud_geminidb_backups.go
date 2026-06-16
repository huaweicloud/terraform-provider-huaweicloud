package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v4/{project_id}/backups
func DataSourceGeminiDBBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBBackupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBBackupsBackupSchema(),
			},
		},
	}
}

func geminiDBBackupsBackupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBBackupsDatastoreSchema(),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBBackupsDatabaseTablesSchema(),
			},
		},
	}
}

func geminiDBBackupsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func geminiDBBackupsDatabaseTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGeminiDBBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v4/{project_id}/backups"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildListGeminiDBBackupsQueryParams(d)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving GeminiDB backups: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB backups: %s", err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB backups: %s", err)
	}

	backups := flattenListGeminiDBBackups(respBody)
	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("backups", backups),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListGeminiDBBackupsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("datastore_type"); ok {
		queryParams = fmt.Sprintf("%s&datastore_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("backup_id"); ok {
		queryParams = fmt.Sprintf("%s&backup_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("backup_type"); ok {
		queryParams = fmt.Sprintf("%s&backup_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenListGeminiDBBackups(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("backups", resp, make([]interface{}, 0))

	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"instance_id":     utils.PathSearch("instance_id", v, nil),
			"instance_name":   utils.PathSearch("instance_name", v, nil),
			"datastore":       flattenGeminiDBBackupsDatastore(v),
			"type":            utils.PathSearch("type", v, nil),
			"size":            utils.PathSearch("size", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"begin_time":      utils.PathSearch("begin_time", v, nil),
			"end_time":        utils.PathSearch("end_time", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"database_tables": flattenGeminiDBBackupsDatabaseTables(v),
		})
	}

	return rst
}

func flattenGeminiDBBackupsDatastore(backup interface{}) []map[string]interface{} {
	datastoreRaw := utils.PathSearch("datastore", backup, nil)
	if datastoreRaw == nil {
		return nil
	}

	datastore := map[string]interface{}{
		"type":    utils.PathSearch("type", datastoreRaw, nil),
		"version": utils.PathSearch("version", datastoreRaw, nil),
	}

	return []map[string]interface{}{datastore}
}

func flattenGeminiDBBackupsDatabaseTables(backup interface{}) []map[string]interface{} {
	databaseTablesRaw := utils.PathSearch("database_tables", backup, nil)
	if databaseTablesRaw == nil {
		return nil
	}

	databaseTablesSlice, ok := databaseTablesRaw.([]interface{})
	if !ok {
		return nil
	}

	databaseTables := make([]map[string]interface{}, 0, len(databaseTablesSlice))
	for _, tableRaw := range databaseTablesSlice {
		table, ok := tableRaw.(map[string]interface{})
		if !ok {
			continue
		}

		tableMap := map[string]interface{}{
			"database_name": utils.PathSearch("database_name", table, nil),
			"table_names":   utils.PathSearch("table_names", table, nil),
		}
		databaseTables = append(databaseTables, tableMap)
	}

	return databaseTables
}
