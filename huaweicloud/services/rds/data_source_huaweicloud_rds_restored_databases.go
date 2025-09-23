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

// @API RDS POST /v3/{project_id}/{engine}/instances/history/databases
func DataSourceRdsRestoredDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsRestoredDatabasesRead,
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
			"database_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"table_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Elem:     restoredDatabasesInstancesSchema(),
				Computed: true,
			},
		},
	}
}

func restoredDatabasesInstancesSchema() *schema.Resource {
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
				Elem:     restoredDatabasesInstancesDatabasesSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func restoredDatabasesInstancesDatabasesSchema() *schema.Resource {
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
		},
	}
	return &sc
}

func dataSourceRdsRestoredDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/{engine}/instances/history/databases"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{engine}", d.Get("engine").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetRestoredDatabasesParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS restored databases: %s", err)
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
		d.Set("database_limit", utils.PathSearch("database_limit", getRespBody, nil)),
		d.Set("table_limit", utils.PathSearch("table_limit", getRespBody, nil)),
		d.Set("instances", flattenRestoredDatabasesInstances(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetRestoredDatabasesParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids":       d.Get("instance_ids").([]interface{}),
		"restore_time":       d.Get("restore_time").(string),
		"instance_name_like": utils.ValueIgnoreEmpty(d.Get("instance_name_like").(string)),
		"database_name_like": utils.ValueIgnoreEmpty(d.Get("database_name_like").(string)),
	}
	return bodyParams
}

func flattenRestoredDatabasesInstances(resp interface{}) []map[string]interface{} {
	instancesJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	instancesArray := instancesJson.([]interface{})
	if len(instancesArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(instancesArray))
	for _, instance := range instancesArray {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", instance, nil),
			"name":         utils.PathSearch("name", instance, nil),
			"total_tables": utils.PathSearch("total_tables", instance, nil),
			"databases":    flattenRestoredDatabasesDatabases(instance),
		})
	}
	return result
}

func flattenRestoredDatabasesDatabases(resp interface{}) []map[string]interface{} {
	databasesJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	databasesArray := databasesJson.([]interface{})
	if len(databasesArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(databasesArray))
	for _, database := range databasesArray {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", database, nil),
			"total_tables": utils.PathSearch("total_tables", database, nil),
		})
	}
	return result
}
