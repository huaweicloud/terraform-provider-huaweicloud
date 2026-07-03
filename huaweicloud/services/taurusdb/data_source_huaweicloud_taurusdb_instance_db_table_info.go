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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/table-info
func DataSourceTaurusDBInstanceDbTableInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBInstanceDbTableInfoRead,

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
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"database_name"},
			},
			"database_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"table_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"table_meta_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"column_default": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_nullable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extra": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBInstanceDbTableInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/table-info"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	// Append query parameters
	queryParams := buildDbTableInfoQueryParams(d)
	if queryParams != "" {
		getPath += "?" + queryParams
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying TaurusDB instance database and table info: %s", err)
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

	databaseNames := flattenStringList(utils.PathSearch("database_names", getRespBody, []interface{}{}))
	tableNames := flattenStringList(utils.PathSearch("table_names", getRespBody, []interface{}{}))
	tableMetaInfos := flattenTableMetaInfos(getRespBody)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("database_names", databaseNames),
		d.Set("table_names", tableNames),
		d.Set("table_meta_infos", tableMetaInfos),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStringList(raw interface{}) []string {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func flattenTableMetaInfos(resp interface{}) []map[string]interface{} {
	tableMetaInfosRaw := utils.PathSearch("table_meta_infos", resp, []interface{}{}).([]interface{})
	if len(tableMetaInfosRaw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tableMetaInfosRaw))
	for _, item := range tableMetaInfosRaw {
		result = append(result, map[string]interface{}{
			"column_name":    utils.PathSearch("columnName", item, nil),
			"column_type":    utils.PathSearch("columnType", item, nil),
			"column_key":     utils.PathSearch("columnKey", item, nil),
			"column_default": utils.PathSearch("columnDefault", item, nil),
			"is_nullable":    utils.PathSearch("isNullable", item, nil),
			"extra":          utils.PathSearch("extra", item, nil),
		})
	}
	return result
}

func buildDbTableInfoQueryParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("database_name"); ok {
		params = fmt.Sprintf("database_name=%s", v.(string))
	}
	if v, ok := d.GetOk("table_name"); ok {
		if params != "" {
			params += "&"
		}
		params += fmt.Sprintf("table_name=%s", v.(string))
	}
	return params
}
