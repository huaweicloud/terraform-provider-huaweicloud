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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/no-index-tables
func DataSourceRdsInstanceNoIndexTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstanceNoIndexTablesRead,
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
			"table_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"newest": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_diagnose_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsInstanceNoIndexTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/no-index-tables"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getInstanceNoIndexTablesQueryParams := buildGetInstanceNoIndexTablesQueryParams(d)
	getPath += getInstanceNoIndexTablesQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS instance no index tables: %s", err)
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
		d.Set("region", region),
		d.Set("tables", flattenInstanceNoIndexTablesBody(getRespBody)),
		d.Set("last_diagnose_timestamp", utils.PathSearch("last_diagnose_timestamp", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceNoIndexTablesQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?table_type=%v&newest=%v", d.Get("table_type"), d.Get("newest"))
}

func flattenInstanceNoIndexTablesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tables", resp, make([]any, 0))
	curArray := curJson.([]any)
	rst := make([]any, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]any{
			"db_name":     utils.PathSearch("db_name", v, nil),
			"schema_name": utils.PathSearch("schema_name", v, nil),
			"table_name":  utils.PathSearch("table_name", v, nil),
		})
	}
	return rst
}
