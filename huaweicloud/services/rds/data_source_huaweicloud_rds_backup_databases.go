package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/db-table-name
func DataSourceRdsBackupDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsBackupDatabasesRead,
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
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     backupDatabaseInfoSchema(),
			},
		},
	}
}

func backupDatabaseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsBackupDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/database/db-table-name"
	getUrl := client.Endpoint + httpUrl
	getUrl = strings.ReplaceAll(getUrl, "{project_id}", client.ProjectID)
	getUrl = strings.ReplaceAll(getUrl, "{instance_id}", d.Get("instance_id").(string))

	backupID := d.Get("backup_id").(string)
	getUrl = fmt.Sprintf("%s?backup_id=%s", getUrl, backupID)

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getUrl, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS backup databases: %s", err)
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
		d.Set("bucket_name", utils.PathSearch("bucket_name", body, nil)),
		d.Set("database_limit", utils.PathSearch("database_limit", body, nil)),
		d.Set("databases", flattenBackupDatabasesList(body)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackupDatabasesList(body interface{}) []interface{} {
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
			"database_name":    utils.PathSearch("database_name", db, nil),
			"backup_file_name": utils.PathSearch("backup_file_name", db, nil),
			"backup_file_size": utils.PathSearch("backup_file_size", db, nil),
		})
	}
	return out
}
