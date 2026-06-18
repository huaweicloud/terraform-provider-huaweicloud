package gaussdb

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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/parse/schema-table
func DataSourceSqlTextSchemaTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlTextSchemaTableRead,

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
			"sql_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTextSchemaTableDatabaseTablesSchema(),
			},
		},
	}
}

func sqlTextSchemaTableDatabaseTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlTextSchemaTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/parse/schema-table"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	getOpt.JSONBody = utils.RemoveNil(buildGetSqlTextSchemaTableBodyParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB SQL text schema table: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("database_tables", flattenGetSqlTextSchemaTableDatabaseTablesBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSqlTextSchemaTableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_text": d.Get("sql_text"),
	}
	return bodyParams
}

func flattenGetSqlTextSchemaTableDatabaseTablesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("database_tables", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"table_name":  utils.PathSearch("table_name", v, nil),
			"schema_name": utils.PathSearch("schema_name", v, nil),
		})
	}
	return rst
}
