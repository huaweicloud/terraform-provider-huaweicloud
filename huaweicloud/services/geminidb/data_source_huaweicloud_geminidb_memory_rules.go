package geminidb

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

// @API GeminiDB GET /v3/{project_id}/dbcache/rules
func DataSourceMemoryRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemoryRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dbcache_mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_db_schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_db_table": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_db_schema": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_db_table": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value_columns": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ttl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_separator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value_separator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildMemoryRulesQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?dbcache_mapping_id=%v&limit=200", d.Get("dbcache_mapping_id"))

	if v, ok := d.GetOk("rule_id"); ok {
		queryParams = fmt.Sprintf("%s&rule_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		queryParams = fmt.Sprintf("%s&rule_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_db_schema"); ok {
		queryParams = fmt.Sprintf("%s&source_db_schema=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_db_table"); ok {
		queryParams = fmt.Sprintf("%s&source_db_table=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceMemoryRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dbcache/rules"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildMemoryRulesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the memory acceleration rules: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryRules := utils.PathSearch("rules", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(memoryRules) == 0 {
			break
		}

		result = append(result, memoryRules...)
		offset += len(memoryRules)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rules", flattenMemoryRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMemoryRules(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"source_db_schema": utils.PathSearch("source_db_schema", v, nil),
			"source_db_table":  utils.PathSearch("source_db_table", v, nil),
			"storage_type":     utils.PathSearch("storage_type", v, nil),
			"target_database":  utils.PathSearch("target_database", v, nil),
			"key_columns":      utils.PathSearch("key_columns", v, nil),
			"value_columns":    utils.PathSearch("value_columns", v, nil),
			"ttl":              utils.PathSearch("ttl", v, nil),
			"key_separator":    utils.PathSearch("key_separator", v, nil),
			"value_separator":  utils.PathSearch("value_separator", v, nil),
			"key_prefix":       utils.PathSearch("key_prefix", v, nil),
		})
	}

	return result
}
