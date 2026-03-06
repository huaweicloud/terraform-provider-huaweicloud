package rds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var publicationNonUpdatableParams = []string{"instance_id", "publication_name", "publication_database",
	"is_create_snapshot_immediately"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/replication/publications
// @API RDS GET /v3/{project_id}/jobs
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/publications
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/replication/publications/{publication_id}
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/replication/publications
func ResourcePublication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicationCreate,
		UpdateContext: resourcePublicationUpdate,
		ReadContext:   resourcePublicationRead,
		DeleteContext: resourcePublicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePublicationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(publicationNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publication_database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_create_snapshot_immediately": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"subscription_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationSubscriptionOptionsSchema(),
			},
			"job_schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     publicationJobScheduleSchema(),
			},
			"is_select_all_table": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"extend_tables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     publicationTablesSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func publicationSubscriptionOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"independent_agent": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"snapshot_always_available": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"replicate_ddl": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"allow_initialize_from_backup": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
		},
	}
}

func publicationJobScheduleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_schedule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"one_time_occurrence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationJobScheduleOneTimeOccurrenceSchema(),
			},
			"frequency": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationJobScheduleFrequencySchema(),
			},
			"daily_frequency": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationJobScheduleDailyFrequencySchema(),
			},
			"duration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationJobScheduleDurationSchema(),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func publicationJobScheduleOneTimeOccurrenceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_start_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func publicationJobScheduleFrequencySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"freq_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"freq_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"freq_interval_weekly": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"freq_interval_day_monthly": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"freq_interval_monthly": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"freq_relative_interval_monthly": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func publicationJobScheduleDailyFrequencySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"freq_subday_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"freq_subday_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"freq_interval_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func publicationJobScheduleDurationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"active_start_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active_end_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func publicationTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"columns": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_key": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter_statement": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationTablesFilterSchema(),
			},
			"article_properties": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     publicationTablesArticlePropertiesSchema(),
			},
		},
	}
}

func publicationTablesFilterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"relation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func publicationTablesArticlePropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"destination_object_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_object_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insert_delivery_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insert_stored_procedure": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_delivery_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_stored_procedure": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_delivery_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_stored_procedure": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePublicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	bodyParams, err := buildCreatePublicationBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	bodyParams = utils.RemoveNil(bodyParams)
	if _, ok := bodyParams["tables"]; !ok {
		// tables can not be empty
		bodyParams["tables"] = make([]interface{}, 0)
	}
	createOpt.JSONBody = bodyParams

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS publication: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating RDS publication: job_id is not found in API response")
	}

	getRespBody, err := getInstanceJob(client, jobId.(string))
	if err != nil {
		return diag.Errorf("error getting RDS job(%s) info: %s", jobId, err)
	}

	publicationId := utils.PathSearch("job.entities.publicationId", getRespBody, nil)
	if publicationId == nil {
		return diag.Errorf("error creating RDS publication: publicationId is not found in API response")
	}

	d.SetId(publicationId.(string))

	if err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error creating RDS publication (%s): %s", instanceId, err)
	}

	return resourcePublicationRead(ctx, d, meta)
}

func buildCreatePublicationBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	tables, err := buildPublicationTablesBodyParams(d.Get("tables"))
	if err != nil {
		return nil, err
	}
	isCreateSnapshotImmediately, _ := strconv.ParseBool(d.Get("is_create_snapshot_immediately").(string))
	bodyParams := map[string]interface{}{
		"publication_name":               d.Get("publication_name"),
		"publication_database":           d.Get("publication_database"),
		"is_create_snapshot_immediately": isCreateSnapshotImmediately,
		"subscription_options":           buildPublicationSubscriptionOptionsBodyParams(d.Get("subscription_options")),
		"job_schedule":                   buildPublicationJobScheduleBodyParams(d.Get("job_schedule")),
		"extend_tables":                  utils.ValueIgnoreEmpty(d.Get("extend_tables").(*schema.Set).List()),
		"tables":                         tables,
	}
	if v, ok := d.GetOk("is_select_all_table"); ok {
		isSelectAllTable, _ := strconv.ParseBool(v.(string))
		bodyParams["is_select_all_table"] = isSelectAllTable
	}
	return bodyParams, nil
}

func buildPublicationSubscriptionOptionsBodyParams(subscriptionOptionsRaw interface{}) map[string]interface{} {
	rawParams := subscriptionOptionsRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{}
		if v["independent_agent"] != nil {
			independentAgent, _ := strconv.ParseBool(v["independent_agent"].(string))
			bodyParams["independent_agent"] = independentAgent
		}
		if v["snapshot_always_available"] != nil {
			snapshotAlwaysAvailable, _ := strconv.ParseBool(v["snapshot_always_available"].(string))
			bodyParams["snapshot_always_available"] = snapshotAlwaysAvailable
		}
		if v["replicate_ddl"] != nil {
			replicateDdl, _ := strconv.ParseBool(v["replicate_ddl"].(string))
			bodyParams["replicate_ddl"] = replicateDdl
		}
		if v["allow_initialize_from_backup"] != nil {
			allowInitializeFromBackup, _ := strconv.ParseBool(v["allow_initialize_from_backup"].(string))
			bodyParams["allow_initialize_from_backup"] = allowInitializeFromBackup
		}
		return bodyParams
	}

	return nil
}

func buildPublicationJobScheduleBodyParams(jobScheduleRaw interface{}) map[string]interface{} {
	rawParams := jobScheduleRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"job_schedule_type":   utils.ValueIgnoreEmpty(v["job_schedule_type"].(string)),
			"one_time_occurrence": buildPublicationJobScheduleOneTimeOccurrenceBodyParams(v["one_time_occurrence"]),
			"frequency":           buildPublicationJobScheduleFrequencyBodyParams(v["frequency"]),
			"daily_frequency":     buildPublicationJobScheduleDailyFrequencyBodyParams(v["daily_frequency"]),
			"duration":            buildPublicationJobScheduleDurationBodyParams(v["duration"]),
		}
		return bodyParams
	}

	return nil
}

func buildPublicationJobScheduleOneTimeOccurrenceBodyParams(oneTimeOccurrenceRaw interface{}) map[string]interface{} {
	rawParams := oneTimeOccurrenceRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"active_start_date": utils.ValueIgnoreEmpty(v["active_start_date"]),
			"active_start_time": utils.ValueIgnoreEmpty(v["active_start_time"]),
		}
		return bodyParams
	}

	return nil
}

func buildPublicationJobScheduleFrequencyBodyParams(frequencyRaw interface{}) map[string]interface{} {
	rawParams := frequencyRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"freq_type":                      utils.ValueIgnoreEmpty(v["freq_type"]),
			"freq_interval":                  utils.ValueIgnoreEmpty(v["freq_interval"]),
			"freq_interval_weekly":           utils.ValueIgnoreEmpty(v["freq_interval_weekly"].(*schema.Set).List()),
			"freq_interval_day_monthly":      utils.ValueIgnoreEmpty(v["freq_interval_day_monthly"]),
			"freq_interval_monthly":          utils.ValueIgnoreEmpty(v["freq_interval_monthly"]),
			"freq_relative_interval_monthly": utils.ValueIgnoreEmpty(v["freq_relative_interval_monthly"]),
		}
		return bodyParams
	}

	return nil
}

func buildPublicationJobScheduleDailyFrequencyBodyParams(dailyFrequencyRaw interface{}) map[string]interface{} {
	rawParams := dailyFrequencyRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"freq_subday_type":     utils.ValueIgnoreEmpty(v["freq_subday_type"]),
			"active_start_time":    utils.ValueIgnoreEmpty(v["active_start_time"]),
			"active_end_time":      utils.ValueIgnoreEmpty(v["active_end_time"]),
			"freq_subday_interval": utils.ValueIgnoreEmpty(v["freq_subday_interval"]),
			"freq_interval_unit":   utils.ValueIgnoreEmpty(v["freq_interval_unit"]),
		}
		return bodyParams
	}

	return nil
}

func buildPublicationJobScheduleDurationBodyParams(durationRaw interface{}) map[string]interface{} {
	rawParams := durationRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"active_start_date": utils.ValueIgnoreEmpty(v["active_start_date"]),
			"active_end_date":   utils.ValueIgnoreEmpty(v["active_end_date"]),
		}
		return bodyParams
	}

	return nil
}

func buildPublicationTablesBodyParams(tablesRaw interface{}) ([]map[string]interface{}, error) {
	rawParams := tablesRaw.(*schema.Set)
	if rawParams.Len() == 0 {
		return nil, nil
	}

	rst := make([]map[string]interface{}, 0, rawParams.Len())
	for _, v := range rawParams.List() {
		if table, ok := v.(map[string]interface{}); ok {
			filter, err := buildPublicationTablesFilterBodyParams(table["filter"])
			if err != nil {
				return nil, err
			}
			rst = append(rst, map[string]interface{}{
				"table_name":         table["table_name"],
				"schema":             utils.ValueIgnoreEmpty(table["schema"]),
				"columns":            utils.ValueIgnoreEmpty(table["columns"].(*schema.Set).List()),
				"primary_key":        utils.ValueIgnoreEmpty(table["primary_key"].(*schema.Set).List()),
				"filter_statement":   utils.ValueIgnoreEmpty(table["filter_statement"]),
				"filter":             filter,
				"article_properties": buildPublicationTablesArticlePropertiesBodyParams(table["article_properties"]),
			})
		}
	}

	return rst, nil
}

func buildPublicationTablesFilterBodyParams(filterRaw interface{}) (map[string]interface{}, error) {
	rawParams := filterRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil, nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		filters, err := buildPublicationTablesFilterFiltersBodyParams(v["filters"])
		if err != nil {
			return nil, err
		}
		bodyParams := map[string]interface{}{
			"relation":  utils.ValueIgnoreEmpty(v["relation"]),
			"column":    utils.ValueIgnoreEmpty(v["column"]),
			"condition": utils.ValueIgnoreEmpty(v["condition"]),
			"value":     utils.ValueIgnoreEmpty(v["value"]),
			"filters":   filters,
		}
		return bodyParams, nil
	}

	return nil, nil
}

func buildPublicationTablesFilterFiltersBodyParams(filtersRaw interface{}) ([]interface{}, error) {
	rawParams := filtersRaw.(*schema.Set)
	if rawParams.Len() == 0 {
		return nil, nil
	}

	rst := make([]interface{}, 0, rawParams.Len())
	for _, v := range rawParams.List() {
		var filter interface{}
		err := json.Unmarshal([]byte(v.(string)), &filter)
		if err != nil {
			return nil, fmt.Errorf("unable to parse JSON of filter: %s", v.(string))
		}
		rst = append(rst, filter)
	}

	return rst, nil
}

func buildPublicationTablesArticlePropertiesBodyParams(filterRaw interface{}) map[string]interface{} {
	rawParams := filterRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"destination_object_name":  utils.ValueIgnoreEmpty(v["destination_object_name"]),
			"destination_object_owner": utils.ValueIgnoreEmpty(v["destination_object_owner"]),
			"insert_delivery_format":   utils.ValueIgnoreEmpty(v["insert_delivery_format"]),
			"insert_stored_procedure":  utils.ValueIgnoreEmpty(v["insert_stored_procedure"]),
			"update_delivery_format":   utils.ValueIgnoreEmpty(v["update_delivery_format"]),
			"update_stored_procedure":  utils.ValueIgnoreEmpty(v["update_stored_procedure"]),
			"delete_delivery_format":   utils.ValueIgnoreEmpty(v["delete_delivery_format"]),
			"delete_stored_procedure":  utils.ValueIgnoreEmpty(v["delete_stored_procedure"]),
		}
		return bodyParams
	}

	return nil
}

func resourcePublicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS publication")
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	publication := utils.PathSearch(fmt.Sprintf("publications[?id=='%s']|[0]", d.Id()), getRespBody, nil)
	if publication == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/instances/{instance_id}/replication/publications",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the RDS publication (%s) does not exist", d.Id())),
			},
		}, "error retrieving RDS publication")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("publication_name", utils.PathSearch("publication_name", publication, nil)),
		d.Set("publication_database", utils.PathSearch("publication_database", publication, nil)),
		d.Set("subscription_options", flattenPublicationSubscriptionOptions(publication)),
		d.Set("job_schedule", flattenPublicationJobSchedule(publication)),
		d.Set("is_select_all_table", strconv.FormatBool(
			utils.PathSearch("is_select_all_table", publication, false).(bool))),
		d.Set("extend_tables", utils.PathSearch("extend_tables", publication, nil)),
		d.Set("tables", flattenPublicationTables(publication)),
		d.Set("status", utils.PathSearch("status", publication, nil)),
		d.Set("subscription_count", utils.PathSearch("subscription_count", publication, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPublicationSubscriptionOptions(publication interface{}) []interface{} {
	subscriptionOptions := utils.PathSearch("subscription_options", publication, nil)
	if subscriptionOptions == nil || len(subscriptionOptions.(map[string]interface{})) == 0 {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"independent_agent": strconv.FormatBool(utils.PathSearch(
				"independent_agent", subscriptionOptions, false).(bool)),
			"snapshot_always_available": strconv.FormatBool(utils.PathSearch(
				"snapshot_always_available", subscriptionOptions, false).(bool)),
			"replicate_ddl": strconv.FormatBool(utils.PathSearch(
				"replicate_ddl", subscriptionOptions, false).(bool)),
			"allow_initialize_from_backup": strconv.FormatBool(utils.PathSearch(
				"allow_initialize_from_backup", subscriptionOptions, false).(bool)),
		},
	}
	return rst
}

func flattenPublicationJobSchedule(publication interface{}) []interface{} {
	jobSchedule := utils.PathSearch("job_schedule", publication, nil)
	if jobSchedule == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"job_schedule_type":   utils.PathSearch("job_schedule_type", jobSchedule, nil),
			"one_time_occurrence": flattenPublicationJobScheduleOneTimeOccurrence(jobSchedule),
			"frequency":           flattenPublicationJobScheduleFrequency(jobSchedule),
			"daily_frequency":     flattenPublicationJobScheduleDailyFrequency(jobSchedule),
			"duration":            flattenPublicationJobScheduleDuration(jobSchedule),
			"id":                  utils.PathSearch("id", jobSchedule, nil),
		},
	}
	return rst
}

func flattenPublicationJobScheduleOneTimeOccurrence(jobSchedule interface{}) []interface{} {
	oneTimeOccurrence := utils.PathSearch("one_time_occurrence", jobSchedule, nil)
	if oneTimeOccurrence == nil || len(oneTimeOccurrence.(map[string]interface{})) == 0 {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", oneTimeOccurrence, nil),
			"active_start_time": utils.PathSearch("active_start_time", oneTimeOccurrence, nil),
		},
	}
	return rst
}

func flattenPublicationJobScheduleFrequency(jobSchedule interface{}) []interface{} {
	frequency := utils.PathSearch("frequency", jobSchedule, nil)
	if frequency == nil || len(frequency.(map[string]interface{})) == 0 {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"freq_type":                      utils.PathSearch("freq_type", frequency, nil),
			"freq_interval":                  utils.PathSearch("freq_interval", frequency, nil),
			"freq_interval_weekly":           utils.PathSearch("freq_interval_weekly", frequency, nil),
			"freq_interval_day_monthly":      utils.PathSearch("freq_interval_day_monthly", frequency, nil),
			"freq_interval_monthly":          utils.PathSearch("freq_interval_monthly", frequency, nil),
			"freq_relative_interval_monthly": utils.PathSearch("freq_relative_interval_monthly", frequency, nil),
		},
	}
	return rst
}

func flattenPublicationJobScheduleDailyFrequency(jobSchedule interface{}) []interface{} {
	dailyFrequency := utils.PathSearch("daily_frequency", jobSchedule, nil)
	if dailyFrequency == nil || len(dailyFrequency.(map[string]interface{})) == 0 {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"freq_subday_type":     utils.PathSearch("freq_subday_type", dailyFrequency, nil),
			"active_start_time":    utils.PathSearch("active_start_time", dailyFrequency, nil),
			"active_end_time":      utils.PathSearch("active_end_time", dailyFrequency, nil),
			"freq_subday_interval": utils.PathSearch("freq_subday_interval", dailyFrequency, nil),
			"freq_interval_unit":   utils.PathSearch("freq_interval_unit", dailyFrequency, nil),
		},
	}
	return rst
}

func flattenPublicationJobScheduleDuration(jobSchedule interface{}) []interface{} {
	duration := utils.PathSearch("duration", jobSchedule, nil)
	if duration == nil || len(duration.(map[string]interface{})) == 0 {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", duration, nil),
			"active_end_date":   utils.PathSearch("active_end_date", duration, nil),
		},
	}
	return rst
}

func flattenPublicationTables(publication interface{}) []interface{} {
	tables := utils.PathSearch("tables", publication, nil)
	if tables == nil || len(tables.([]interface{})) == 0 {
		return nil
	}

	tableArray := tables.([]interface{})
	rst := make([]interface{}, 0, len(tableArray))
	for _, v := range tableArray {
		rst = append(rst, map[string]interface{}{
			"table_name":         utils.PathSearch("table_name", v, nil),
			"schema":             utils.PathSearch("schema", v, nil),
			"columns":            utils.PathSearch("columns", v, nil),
			"primary_key":        utils.PathSearch("primary_key", v, nil),
			"filter_statement":   utils.PathSearch("filter_statement", v, nil),
			"filter":             flattenPublicationTablesFilter(v),
			"article_properties": flattenPublicationTablesArticleProperties(v),
		})
	}
	return rst
}

func flattenPublicationTablesFilter(table interface{}) []interface{} {
	filter := utils.PathSearch("filter", table, nil)
	if filter == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"relation":  utils.PathSearch("relation", filter, nil),
			"column":    utils.PathSearch("column", filter, nil),
			"condition": utils.PathSearch("condition", filter, nil),
			"value":     utils.PathSearch("value", filter, nil),
			"filters":   flattenPublicationTablesFilterFilters(filter),
		},
	}
	return rst
}

func flattenPublicationTablesFilterFilters(filter interface{}) []interface{} {
	filters := utils.PathSearch("filters", filter, nil)
	if filters == nil {
		return nil
	}

	filterArray := filters.([]interface{})
	rst := make([]interface{}, 0, len(filterArray))
	for _, v := range filterArray {
		jsonRaw, err := json.Marshal(v)
		if err != nil {
			log.Printf("[ERROR] unable to convert the filter to json")
		} else {
			rst = append(rst, string(jsonRaw))
		}
	}
	return rst
}

func flattenPublicationTablesArticleProperties(table interface{}) []interface{} {
	articleProperties := utils.PathSearch("article_properties", table, nil)
	if articleProperties == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"destination_object_name":  utils.PathSearch("destination_object_name", articleProperties, nil),
			"destination_object_owner": utils.PathSearch("destination_object_owner", articleProperties, nil),
			"insert_delivery_format":   utils.PathSearch("insert_delivery_format", articleProperties, nil),
			"insert_stored_procedure":  utils.PathSearch("insert_stored_procedure", articleProperties, nil),
			"update_delivery_format":   utils.PathSearch("update_delivery_format", articleProperties, nil),
			"update_stored_procedure":  utils.PathSearch("update_stored_procedure", articleProperties, nil),
			"delete_delivery_format":   utils.PathSearch("delete_delivery_format", articleProperties, nil),
			"delete_stored_procedure":  utils.PathSearch("delete_stored_procedure", articleProperties, nil),
		},
	}
	return rst
}

func resourcePublicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications/{publication_id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{publication_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	bodyParams, err := buildUpdatePublicationBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	bodyParams = utils.RemoveNil(bodyParams)
	if _, ok := bodyParams["tables"]; !ok {
		// tables can not be empty
		bodyParams["tables"] = make([]interface{}, 0)
	}
	updateOpt.JSONBody = bodyParams

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating RDS publication: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error updating RDS publication: job_id is not found in API response")
	}

	if err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("error updating RDS publication (%s): %s", d.Id(), err)
	}

	return resourcePublicationRead(ctx, d, meta)
}

func buildUpdatePublicationBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	tables, err := buildPublicationTablesBodyParams(d.Get("tables"))
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"subscription_options": buildPublicationSubscriptionOptionsBodyParams(d.Get("subscription_options")),
		"job_schedule":         buildPublicationJobScheduleBodyParams(d.Get("job_schedule")),
		"extend_tables":        utils.ValueIgnoreEmpty(d.Get("extend_tables").(*schema.Set).List()),
		"tables":               tables,
	}
	if v, ok := d.GetOk("is_select_all_table"); ok {
		isSelectAllTable, _ := strconv.ParseBool(v.(string))
		bodyParams["is_select_all_table"] = isSelectAllTable
	}

	return bodyParams, nil
}

func resourcePublicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publications"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = buildDeletePublicationBody(d.Id())

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		// if an instance is deleted and then delete a publication, an operation conflict will be reported.
		// Therefore, it is necessary to check whether the instance still exists.
		if retry {
			instance, err := GetRdsInstanceByID(client, d.Id())
			if err != nil {
				return res, false, err
			}
			if instance == nil {
				return res, false, golangsdk.ErrDefault404{}
			}
		}
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting RDS publication")
	}

	deleteRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting RDS publication: job_id is not found in API response")
	}

	if err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error deleting RDS publication (%s): %s", d.Id(), err)
	}

	return nil
}

func buildDeletePublicationBody(publicationId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publication_ids": []string{publicationId},
	}
	return bodyParams
}

func resourcePublicationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
