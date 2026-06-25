package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/configuration/{database}
func DataSourceTaurusDBHtapStarrocksDbReplicationConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksDbReplicationConfigRead,

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
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_instance_level_sync": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_repl_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationDatabaseInfoSchema(),
			},
			"table_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationTableConfigCheckResultSchema(),
			},
			"table_repl_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationTableReplConfigComputedSchema(),
			},
			"new_table_repl_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationTableReplConfigComputedSchema(),
			},
			"target_database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_tables_change": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_error_of_alter_table": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_support_reg_exp": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksDbReplicationConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		database   = d.Get("database").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	respBody, err := getReplicationConfigByDatabase(client, instanceId, database)
	if err != nil {
		return diag.Errorf("error retrieving HTAP instance(%s) database (%s) replication config: %s", instanceId, database, err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("region", region),
		d.Set("source_instance_id", utils.PathSearch("source_instance_id", respBody, nil)),
		d.Set("source_node_id", utils.PathSearch("source_node_id", respBody, nil)),
		d.Set("is_instance_level_sync", utils.PathSearch("is_instance_level_sync", respBody, nil)),
		d.Set("database_repl_scope", utils.PathSearch("database_repl_scope", respBody, nil)),
		d.Set("database_info", flattenStarrocksReplicationDatabaseInfo(respBody)),
		d.Set("table_infos", flattenStarrocksReplicationTableInfos(respBody)),
		d.Set("table_repl_config", flattenStarrocksReplicationTableReplConfig(utils.PathSearch("table_repl_config", respBody, nil))),
		d.Set("new_table_repl_config", flattenStarrocksReplicationTableReplConfig(utils.PathSearch("new_table_repl_config", respBody, nil))),
		d.Set("target_database_name", utils.PathSearch("target_database_name", respBody, nil)),
		d.Set("is_support_reg_exp", utils.PathSearch("is_support_reg_exp", respBody, nil)),
		d.Set("is_tables_change", utils.PathSearch("is_tables_change", respBody, false).(bool)),
		d.Set("error_msg", utils.PathSearch("error_msg", respBody, nil)),
		d.Set("last_error_of_alter_table", utils.PathSearch("last_error_of_alter_table", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getReplicationConfigByDatabase(client *golangsdk.ServiceClient, instanceId, database string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/configuration/{database}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{database}", database)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
