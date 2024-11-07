package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/{database_name}/instances/history/tables
func DataSourceRdsRestoredTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsRestoredTablesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"restore_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Elem:     restoredTablesInstancesSchema(),
				Computed: true,
			},
		},
	}
}

func restoredTablesInstancesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_tables": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Elem:     restoredTablesInstancesDatabasesSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func restoredTablesInstancesDatabasesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_tables": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"schemas": {
				Type:     schema.TypeList,
				Elem:     restoredTablesInstancesDatabasesSchemasSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func restoredTablesInstancesDatabasesSchemasSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_tables": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Elem:     restoredTablesInstancesDatabasesSchemasTablesSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func restoredTablesInstancesDatabasesSchemasTablesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRdsRestoredTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/{database_name}/instances/history/tables"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{database_name}", d.Get("engine").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetRestoredTablesParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS restored tables: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("table_limit", utils.PathSearch("table_limit", getRespBody, nil)),
		d.Set("instances", flattenRestoredTablesInstances(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetRestoredTablesParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids":       d.Get("instance_ids").([]interface{}),
		"restore_time":       d.Get("restore_time").(string),
		"instance_name_like": utils.ValueIgnoreEmpty(d.Get("instance_name_like").(string)),
		"database_name_like": utils.ValueIgnoreEmpty(d.Get("database_name_like").(string)),
		"table_name_like":    utils.ValueIgnoreEmpty(d.Get("table_name_like").(string)),
	}
	return bodyParams
}

func flattenRestoredTablesInstances(resp interface{}) []map[string]interface{} {
	instancesJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	instancesArray := instancesJson.([]interface{})
	if len(instancesArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(instancesArray))
	for _, v := range instancesArray {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"name":         utils.PathSearch("name", v, nil),
			"total_tables": utils.PathSearch("total_tables", v, nil),
			"databases":    flattenRestoredTablesDatabases(v),
		})
	}
	return result
}

func flattenRestoredTablesDatabases(resp interface{}) []map[string]interface{} {
	databasesJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databasesArray := databasesJson.([]interface{})
	if len(databasesArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(databasesArray))
	for _, v := range databasesArray {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"total_tables": utils.PathSearch("total_tables", v, nil),
			"schemas":      flattenRestoredTablesSchemas(v),
		})
	}
	return result
}

func flattenRestoredTablesSchemas(resp interface{}) []map[string]interface{} {
	schemasJson := utils.PathSearch("schemas", resp, make([]interface{}, 0))
	schemasArray := schemasJson.([]interface{})
	if len(schemasArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(schemasArray))
	for _, v := range schemasArray {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"total_tables": utils.PathSearch("total_tables", v, nil),
			"tables":       flattenRestoredTablesTables(v),
		})
	}
	return result
}

func flattenRestoredTablesTables(resp interface{}) []map[string]interface{} {
	tablesJson := utils.PathSearch("tables", resp, make([]interface{}, 0))
	tablesArray := tablesJson.([]interface{})
	if len(tablesArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tablesArray))
	for _, v := range tablesArray {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
		})
	}
	return result
}
