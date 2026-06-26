package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/tables
func DataSourceGaussDBInstanceDatabaseTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInstanceDatabaseTablesRead,

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
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_name_keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     databaseTablesSchema(),
			},
		},
	}
}

func databaseTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBInstanceDatabaseTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/tables"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildGetDatabaseTablesQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance database tables: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tables", flattenGetDatabaseTablesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDatabaseTablesQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&db_name=%v", res, d.Get("db_name"))
	res = fmt.Sprintf("%s&schema_name=%v", res, d.Get("schema_name"))
	if v, ok := d.GetOk("table_name_keyword"); ok {
		res = fmt.Sprintf("%s&table_name_keyword=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetDatabaseTablesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("database_tables", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"table_name": utils.PathSearch("table_name", v, nil),
		})
	}
	return res
}
