package taurusdb

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

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/htap/tables
func DataSourceTaurusDBHtapPrimaryInstanceTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapPrimaryInstanceTablesRead,

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
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_tables": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     databaseTablesInfoSchema(),
			},
			"selected_tables": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     databaseTablesInfoSchema(),
			},
			"filter_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func databaseTablesInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTaurusDBHtapPrimaryInstanceTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		offset = 0
		result = make([]interface{}, 0)
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/htap/tables"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildListTaurusDBHtapPrimaryInstanceTablesBody(d),
	}

	for {
		currentPath := fmt.Sprintf("%s?limit=100&offset=%d", listPath, offset)
		resp, err := client.Request("POST", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP primary instance tables: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tables := utils.PathSearch("tables", respBody, make([]interface{}, 0)).([]interface{})
		if len(tables) == 0 {
			break
		}

		result = append(result, tables...)

		totalCount := utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if int(totalCount) == len(result) {
			break
		}

		offset += len(tables)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("tables", utils.ExpandToStringList(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListTaurusDBHtapPrimaryInstanceTablesBody(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"source_instance_id": d.Get("source_instance_id").(string),
		"database_tables":    buildDatabaseTablesInfo(d.Get("database_tables").([]interface{})),
		"selected_tables":    buildDatabaseTablesInfo(d.Get("selected_tables").([]interface{})),
		"filter_type":        d.Get("filter_type").(string),
	}
}

func buildDatabaseTablesInfo(raw []interface{}) []interface{} {
	res := make([]interface{}, 0, len(raw))
	for _, v := range raw {
		item := v.(map[string]interface{})
		tableInfo := map[string]interface{}{
			"database": item["database"].(string),
		}
		if tables, ok := item["tables"]; ok && tables != nil {
			tableInfo["tables"] = utils.ExpandToStringList(tables.([]interface{}))
		}
		res = append(res, tableInfo)
	}
	return res
}
