package drs

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

// @API DRS GET /v5.1/{project_id}/jobs/{job_id}/db-object
func DataSourceDrsDbObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDbObjectRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_root_db": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dbObjectTargetRootDbSchema(),
			},
			"object_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"all": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"schemas": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sync_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"all": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tables": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     dbObjectTargetTablesSchema(),
									},
								},
							},
						},
						"tables": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dbObjectTargetTablesSchema(),
						},
						"total_table_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_synchronized": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"max_table_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dbObjectTargetRootDbSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dbObjectTargetTablesSchema() *schema.Resource {
	return &schema.Resource{

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"all": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"db_alias_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_alias_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filtered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"filter_conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"config_conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_synchronized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key_for_data_filtering": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_for_data_filtering": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key_for_column_filtering": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filtered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"additional": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDbObjectQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?type=%s", d.Get("type").(string))
}

func dataSourceDrsDbObjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5.1/{project_id}/jobs/{job_id}/db-object"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath+buildDbObjectQueryParams(d), &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS database object: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
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
		d.Set("target_root_db", flattenDbObjectTargetRootDb(
			utils.PathSearch("target_root_db", respBody, nil))),
		d.Set("object_info", flattenObjectInfo(
			utils.PathSearch("object_info", respBody, nil))),
		d.Set("max_table_num", utils.PathSearch("max_table_num", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("object_scope", utils.PathSearch("object_scope", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDbObjectTargetRootDb(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"db_name":     utils.PathSearch("db_name", resp, nil),
			"db_encoding": utils.PathSearch("db_encoding", resp, nil),
		},
	}
}

func flattenObjectInfo(resp interface{}) []interface{} {
	objectInfoMap, ok := resp.(map[string]interface{})
	if !ok || len(objectInfoMap) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(objectInfoMap))
	for key, value := range objectInfoMap {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"key":             key,
			"sync_type":       utils.PathSearch("sync_type", item, nil),
			"name":            utils.PathSearch("name", item, nil),
			"all":             utils.PathSearch("all", item, nil),
			"schemas":         flattenDbObjectSchemas(utils.PathSearch("schemas", item, nil)),
			"tables":          flattenDbObjectTables(utils.PathSearch("tables", item, nil)),
			"total_table_num": utils.PathSearch("total_table_num", item, nil),
			"is_synchronized": utils.PathSearch("is_synchronized", item, nil),
		})
	}

	return result
}

func flattenDbObjectSchemas(resp interface{}) []interface{} {
	schemasMap, ok := resp.(map[string]interface{})
	if !ok || len(schemasMap) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(schemasMap))
	for key, value := range schemasMap {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"key":       key,
			"sync_type": utils.PathSearch("sync_type", item, nil),
			"name":      utils.PathSearch("name", item, nil),
			"all":       utils.PathSearch("all", item, nil),
			"tables":    flattenDbObjectTables(utils.PathSearch("tables", item, nil)),
		})
	}

	return result
}

func flattenDbObjectTables(resp interface{}) []interface{} {
	tablesMap, ok := resp.(map[string]interface{})
	if !ok || len(tablesMap) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(tablesMap))
	for key, value := range tablesMap {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"key":               key,
			"sync_type":         utils.PathSearch("sync_type", item, nil),
			"type":              utils.PathSearch("type", item, nil),
			"name":              utils.PathSearch("name", item, nil),
			"all":               utils.PathSearch("all", item, nil),
			"db_alias_name":     utils.PathSearch("db_alias_name", item, nil),
			"schema_alias_name": utils.PathSearch("schema_alias_name", item, nil),
			"filtered":          utils.PathSearch("filtered", item, nil),
			"filter_conditions": utils.PathSearch("filter_conditions", item, nil),
			"config_conditions": utils.PathSearch("config_conditions", item, nil),
			"is_synchronized":   utils.PathSearch("is_synchronized", item, nil),
			"columns":           flattenDbObjectColumns(utils.PathSearch("columns", item, nil)),
		})
	}

	return result
}

func flattenDbObjectColumns(resp interface{}) []interface{} {
	columnsMap, ok := resp.(map[string]interface{})
	if !ok || len(columnsMap) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(columnsMap))
	for key, value := range columnsMap {
		item, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"key":       key,
			"sync_type": utils.PathSearch("sync_type", item, nil),
			"primary_key_for_data_filtering": utils.PathSearch(
				"primary_key_for_data_filtering", item, nil),
			"index_for_data_filtering": utils.PathSearch(
				"index_for_data_filtering", item, nil),
			"name": utils.PathSearch("name", item, nil),
			"type": utils.PathSearch("type", item, nil),
			"primary_key_for_column_filtering": utils.PathSearch(
				"primary_key_for_column_filtering", item, nil),
			"filtered":       utils.PathSearch("filtered", item, nil),
			"additional":     utils.PathSearch("additional", item, nil),
			"operation_type": utils.PathSearch("operation_type", item, nil),
			"value":          utils.PathSearch("value", item, nil),
		})
	}

	return result
}
