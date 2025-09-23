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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/log-replay/database
func DataSourceRdsReadReplicaRestorableDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsReadReplicaRestorableDatabasesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_tables": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"table_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     replicaRestorableDatabasesSchema(),
			},
		},
	}
}

func replicaRestorableDatabasesSchema() *schema.Resource {
	return &schema.Resource{
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
				Elem:     replicaRestorableSchemasSchema(),
				Computed: true,
			},
		},
	}
}

func replicaRestorableSchemasSchema() *schema.Resource {
	return &schema.Resource{
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
				Elem:     replicaRestorableTablesSchema(),
				Computed: true,
			},
		},
	}
}
func replicaRestorableTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsReadReplicaRestorableDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/log-replay/database"
	getUrl := client.Endpoint + httpUrl
	getUrl = strings.ReplaceAll(getUrl, "{project_id}", client.ProjectID)
	getUrl = strings.ReplaceAll(getUrl, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getUrl, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS read replica restorable databases: %s", err)
	}

	body, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("database_limit", utils.PathSearch("database_limit", body, nil)),
		d.Set("total_tables", utils.PathSearch("total_tables", body, nil)),
		d.Set("table_limit", utils.PathSearch("table_limit", body, nil)),
		d.Set("databases", flattenRestorableDatabases(body)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRestorableDatabases(body interface{}) []interface{} {
	raw := utils.PathSearch("databases", body, nil)
	if raw == nil {
		return nil
	}

	databases, ok := raw.([]interface{})
	if !ok || len(databases) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(databases))
	for _, db := range databases {
		out = append(out, map[string]interface{}{
			"name":         utils.PathSearch("name", db, nil),
			"total_tables": utils.PathSearch("total_tables", db, nil),
			"schemas":      flattenRestorableSchemas(db),
		})
	}
	return out
}

func flattenRestorableSchemas(db interface{}) []interface{} {
	raw := utils.PathSearch("schemas", db, nil)
	if raw == nil {
		return nil
	}

	schemas, ok := raw.([]interface{})
	if !ok || len(schemas) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(schemas))
	for _, sc := range schemas {
		out = append(out, map[string]interface{}{
			"name":         utils.PathSearch("name", sc, nil),
			"total_tables": utils.PathSearch("total_tables", sc, nil),
			"tables":       flattenRestorableTables(sc),
		})
	}
	return out
}

func flattenRestorableTables(sc interface{}) []interface{} {
	raw := utils.PathSearch("tables", sc, nil)
	if raw == nil {
		return nil
	}

	tables, ok := raw.([]interface{})
	if !ok || len(tables) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(tables))
	for _, t := range tables {
		out = append(out, map[string]interface{}{
			"name": utils.PathSearch("name", t, nil),
		})
	}
	return out
}
