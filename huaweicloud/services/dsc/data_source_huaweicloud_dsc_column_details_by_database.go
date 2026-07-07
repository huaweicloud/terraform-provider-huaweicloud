package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/metadata/catalog/column-details/database-dim
func DataSourceDscColumnDetailsByDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscColumnDetailsByDatabaseRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"label_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the group label ID for filtering. Either `label_id` or `type_id` must be specified.",
			},
			"type_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type ID for filtering. Either `label_id` or `type_id` must be specified.",
			},
			"results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column details list by database dimension.",
				Elem:        dscColumnDetailsDatabaseDimSchema(),
			},
		},
	}
}

func dscColumnDetailsDatabaseDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The match count.",
			},
			"db_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name.",
			},
			"db_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database type.",
			},
			"tables": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The table information list.",
				Elem:        dscColumnDetailsTableDimSchema(),
			},
		},
	}
}

func dscColumnDetailsTableDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The match count.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscColumnDetailsColumnDimSchema(),
			},
		},
	}
}

func dscColumnDetailsColumnDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The column name.",
			},
			"classification_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The classification tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"level_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The level tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildDscColumnDetailsByDatabaseQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("label_id"); ok {
		queryParams = fmt.Sprintf("%s&label_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type_id"); ok {
		queryParams = fmt.Sprintf("%s&type_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscColumnDetailsByDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/metadata/catalog/column-details/database-dim"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	currentPath := requestPath + buildDscColumnDetailsByDatabaseQueryParams(d)

	resp, err := client.Request("GET", currentPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC column details by database: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("results", flattenDscColumnDetailsDatabaseDim(
			utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscColumnDetailsDatabaseDim(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, v := range results {
		rst = append(rst, map[string]interface{}{
			"asset_name": utils.PathSearch("asset_name", v, nil),
			"count":      utils.PathSearch("count", v, nil),
			"db_name":    utils.PathSearch("db_name", v, nil),
			"db_type":    utils.PathSearch("db_type", v, nil),
			"tables": flattenDscColumnDetailsTableDim(
				utils.PathSearch("tables", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscColumnDetailsTableDim(tables []interface{}) []interface{} {
	if len(tables) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tables))
	for _, v := range tables {
		rst = append(rst, map[string]interface{}{
			"count":      utils.PathSearch("count", v, nil),
			"table_name": utils.PathSearch("table_name", v, nil),
			"columns": flattenDscColumnDetailsColumnDim(
				utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscColumnDetailsColumnDim(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, map[string]interface{}{
			"column_name":         utils.PathSearch("column_name", v, nil),
			"classification_tags": utils.PathSearch("classification_tags", v, nil),
			"level_tags":          utils.PathSearch("level_tags", v, nil),
		})
	}

	return rst
}
