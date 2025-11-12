package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/publications
func DataSourceRdsPublications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsPublicationsRead,
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
			"publication_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publication_db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriber_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsSchema(),
			},
		},
	}
}

func publicationsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsSubscriptionOptionsSchema(),
			},
			"job_schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsJobScheduleSchema(),
			},
			"is_select_all_table": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"extend_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsTablesSchema(),
			},
		},
	}
	return &sc
}

func publicationsSubscriptionOptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"independent_agent": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"snapshot_always_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replicate_ddl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_initialize_from_backup": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsJobScheduleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_schedule_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"one_time_occurrence": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsJobScheduleOneTimeOccurrenceSchema(),
			},
			"frequency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsJobScheduleFrequencySchema(),
			},
			"daily_frequency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsJobScheduleDailyFrequencySchema(),
			},
			"duration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsJobScheduleDurationSchema(),
			},
		},
	}
	return &sc
}

func publicationsJobScheduleOneTimeOccurrenceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsJobScheduleFrequencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"freq_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"freq_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"freq_interval_weekly": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"freq_interval_day_monthly": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"freq_interval_monthly": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"freq_relative_interval_monthly": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsJobScheduleDailyFrequencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"freq_subday_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"freq_subday_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"freq_interval_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsJobScheduleDurationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_end_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsTablesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter_statement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsTablesFilterSchema(),
			},
			"article_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationsTablesArticlePropertiesSchema(),
			},
		},
	}
	return &sc
}

func publicationsTablesFilterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"relation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filters": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func publicationsTablesArticlePropertiesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination_object_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_object_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insert_delivery_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insert_stored_procedure": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_delivery_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_stored_procedure": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_delivery_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_stored_procedure": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRdsPublicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	listPublicationsQueryParams := buildGetPublicationsQueryParams(d)
	listPath += listPublicationsQueryParams

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving RDS publications: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespJsonBody interface{}
	err = json.Unmarshal(listRespJson, &listRespJsonBody)
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
		d.Set("publications", flattenPublicationsBody(listRespJsonBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetPublicationsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("publication_name"); ok {
		res = fmt.Sprintf("%s&publication_name=%v", res, v)
	}
	if v, ok := d.GetOk("publication_db_name"); ok {
		res = fmt.Sprintf("%s&publication_db_name=%v", res, v)
	}
	if v, ok := d.GetOk("subscriber_instance_id"); ok {
		res = fmt.Sprintf("%s&subscriber_instance_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenPublicationsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("publications", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]any, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"publication_name":     utils.PathSearch("publication_name", v, nil),
			"publication_database": utils.PathSearch("publication_database", v, nil),
			"subscription_count":   utils.PathSearch("subscription_count", v, nil),
			"subscription_options": flattenPublicationsSubscriptionOptionsBody(v),
			"job_schedule":         flattenPublicationsJobScheduleBody(v),
			"is_select_all_table":  utils.PathSearch("is_select_all_table", v, nil),
			"extend_tables":        utils.PathSearch("extend_tables", v, nil),
			"tables":               flattenPublicationsTablesBody(v),
		})
	}
	return rst
}

func flattenPublicationsSubscriptionOptionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("subscription_options", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"independent_agent":            utils.PathSearch("independent_agent", curJson, nil),
			"snapshot_always_available":    utils.PathSearch("snapshot_always_available", curJson, nil),
			"replicate_ddl":                utils.PathSearch("replicate_ddl", curJson, nil),
			"allow_initialize_from_backup": utils.PathSearch("allow_initialize_from_backup", curJson, nil),
		},
	}
	return rst
}

func flattenPublicationsJobScheduleBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("job_schedule", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"id":                  utils.PathSearch("id", curJson, nil),
			"job_schedule_type":   utils.PathSearch("job_schedule_type", curJson, nil),
			"one_time_occurrence": flattenPublicationsJobScheduleOneTimeOccurrenceBody(curJson),
			"frequency":           flattenPublicationsJobScheduleFrequencyBody(curJson),
			"daily_frequency":     flattenPublicationsJobScheduleDailyFrequencyBody(curJson),
			"duration":            flattenPublicationsJobScheduleDurationBody(curJson),
		},
	}
	return rst
}

func flattenPublicationsJobScheduleOneTimeOccurrenceBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("one_time_occurrence", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", curJson, nil),
			"active_start_time": utils.PathSearch("active_start_time", curJson, nil),
		},
	}
	return rst
}

func flattenPublicationsJobScheduleFrequencyBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("frequency", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"freq_type":                      utils.PathSearch("freq_type", curJson, nil),
			"freq_interval":                  utils.PathSearch("freq_interval", curJson, nil),
			"freq_interval_weekly":           utils.PathSearch("freq_interval_weekly", curJson, nil),
			"freq_interval_day_monthly":      utils.PathSearch("freq_interval_day_monthly", curJson, nil),
			"freq_interval_monthly":          utils.PathSearch("freq_interval_monthly", curJson, nil),
			"freq_relative_interval_monthly": utils.PathSearch("freq_relative_interval_monthly", curJson, nil),
		},
	}
	return rst
}

func flattenPublicationsJobScheduleDailyFrequencyBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("daily_frequency", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"freq_subday_type":     utils.PathSearch("freq_subday_type", curJson, nil),
			"active_start_time":    utils.PathSearch("active_start_time", curJson, nil),
			"active_end_time":      utils.PathSearch("active_end_time", curJson, nil),
			"freq_subday_interval": utils.PathSearch("freq_subday_interval", curJson, nil),
			"freq_interval_unit":   utils.PathSearch("freq_interval_unit", curJson, nil),
		},
	}
	return rst
}

func flattenPublicationsJobScheduleDurationBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("duration", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", curJson, nil),
			"active_end_date":   utils.PathSearch("active_end_date", curJson, nil),
		},
	}
	return rst
}

func flattenPublicationsTablesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tables", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]any, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"table_name":         utils.PathSearch("table_name", v, nil),
			"schema":             utils.PathSearch("schema", v, nil),
			"columns":            utils.PathSearch("columns", v, nil),
			"primary_key":        utils.PathSearch("primary_key", v, nil),
			"filter_statement":   utils.PathSearch("filter_statement", v, nil),
			"filter":             flattenPublicationsTablesFilterBody(v),
			"article_properties": flattenPublicationsTablesArticlePropertiesBody(v),
		})
	}
	return rst
}

func flattenPublicationsTablesFilterBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("filter", resp, nil)
	if curJson == nil {
		return nil
	}

	filtersJson, _ := json.Marshal(utils.PathSearch("filters", curJson, nil))
	rst := []interface{}{
		map[string]interface{}{
			"relation":  utils.PathSearch("relation", curJson, nil),
			"column":    utils.PathSearch("column", curJson, nil),
			"condition": utils.PathSearch("condition", curJson, nil),
			"value":     utils.PathSearch("value", curJson, nil),
			"filters":   string(filtersJson),
		},
	}
	return rst
}

func flattenPublicationsTablesArticlePropertiesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("article_properties", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"destination_object_name":  utils.PathSearch("destination_object_name", curJson, nil),
			"destination_object_owner": utils.PathSearch("destination_object_owner", curJson, nil),
			"insert_delivery_format":   utils.PathSearch("insert_delivery_format", curJson, nil),
			"insert_stored_procedure":  utils.PathSearch("insert_stored_procedure", curJson, nil),
			"update_delivery_format":   utils.PathSearch("update_delivery_format", curJson, nil),
			"update_stored_procedure":  utils.PathSearch("update_stored_procedure", curJson, nil),
			"delete_delivery_format":   utils.PathSearch("delete_delivery_format", curJson, nil),
			"delete_stored_procedure":  utils.PathSearch("delete_stored_procedure", curJson, nil),
		},
	}
	return rst
}
