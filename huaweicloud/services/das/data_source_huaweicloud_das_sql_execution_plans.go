package das

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

// @API DAS GET /v3/{project_id}/instances/{instance_id}/sql/explain
func DataSourceSqlExecutionPlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlExecutionPlansRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the SQL execution plans are located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database instance.",
			},
			"db_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database user ID.",
			},
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database name.",
			},
			"sql": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SQL statement.",
			},

			// Attributes.
			"plans": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of SQL execution plans.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the execution plan step.",
						},
						"select_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The select type of the query.",
						},
						"table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table name.",
						},
						"partitions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The partitions that the query will match.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access type.",
						},
						"possible_keys": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The possible keys that could be used.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The actual key used.",
						},
						"key_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The length of the key used.",
						},
						"ref": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The column or constant used with the key to select rows.",
						},
						"rows": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of rows MySQL estimates it must examine.",
						},
						"filtered": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The percentage of rows filtered by the table condition.",
						},
						"extra": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Additional information.",
						},
					},
				},
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The error message if the SQL execution failed.",
			},
		},
	}
}

func dataSourceSqlExecutionPlansRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	executionPlans, errorMessage, err := getSqlExecutionPlans(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS SQL execution plans: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plans", flattenSqlExecutionPlans(executionPlans)),
		d.Set("error_message", errorMessage),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSqlExecutionPlans(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, string, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql/explain"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildSqlExecutionPlansQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, "", err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, "", err
	}

	executionPlans := utils.PathSearch("execution_plans", respBody, make([]interface{}, 0)).([]interface{})
	errorMessage := utils.PathSearch("error_message", respBody, "").(string)

	return executionPlans, errorMessage, nil
}

func buildSqlExecutionPlansQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?db_user_id=%v&database=%v&sql=%v",
		d.Get("db_user_id").(string),
		d.Get("database").(string),
		d.Get("sql").(string))

	return res
}

func flattenSqlExecutionPlans(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", item, nil),
			"select_type":   utils.PathSearch("select_type", item, nil),
			"table":         utils.PathSearch("table", item, nil),
			"partitions":    utils.PathSearch("partitions", item, nil),
			"type":          utils.PathSearch("type", item, nil),
			"possible_keys": utils.PathSearch("possible_keys", item, nil),
			"key":           utils.PathSearch("key", item, nil),
			"key_len":       utils.PathSearch("key_len", item, nil),
			"ref":           utils.PathSearch("ref", item, nil),
			"rows":          utils.PathSearch("rows", item, nil),
			"filtered":      utils.PathSearch("filtered", item, nil),
			"extra":         utils.PathSearch("extra", item, nil),
		})
	}

	return result
}
