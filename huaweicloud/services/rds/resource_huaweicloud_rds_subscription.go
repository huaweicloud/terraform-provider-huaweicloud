package rds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

var rdsSubscriptionNonUpdatableParams = []string{"instance_id", "subscription_database", "subscription_type",
	"initialize_at", "current_publication_id", "initialize_info", "initialize_info.*.file_source",
	"initialize_info.*.backup_id", "initialize_info.*.bucket_name", "initialize_info.*.file_path",
	"initialize_info.*.file_name", "initialize_info.*.overwrite_restore", "independent_agent", "bak_file_name",
	"bak_bucket_name", "local_subscription", "local_subscription.*.publication_id", "local_subscription.*.publication_name",
}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/replication/subscriptions
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/subscriptions
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/replication/subscriptions
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/replication/subscriptions
func ResourceRdsSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsSubscriptionCreate,
		UpdateContext: resourceRdsSubscriptionUpdate,
		ReadContext:   resourceRdsSubscriptionRead,
		DeleteContext: resourceRdsSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsSubscriptionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(rdsSubscriptionNonUpdatableParams),

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
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initialize_at": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionJobScheduleSchema(),
			},
			"local_subscription": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionLocalSubscriptionSchema(),
			},
			"initialize_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionInitializeInfoSchema(),
			},
			"independent_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"bak_file_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bak_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"is_cloud": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func rdsSubscriptionJobScheduleSchema() *schema.Resource {
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
				Elem:     rdsSubscriptionOneTimeOccurrenceSchema(),
			},
			"frequency": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionFrequencySchema(),
			},
			"daily_frequency": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionDailyFrequencySchema(),
			},
			"duration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     rdsSubscriptionDurationSchema(),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func rdsSubscriptionOneTimeOccurrenceSchema() *schema.Resource {
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

func rdsSubscriptionFrequencySchema() *schema.Resource {
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

func rdsSubscriptionDailyFrequencySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"freq_subday_type": {
				Type:     schema.TypeString,
				Optional: true,
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

func rdsSubscriptionDurationSchema() *schema.Resource {
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

func rdsSubscriptionInitializeInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overwrite_restore": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
		},
	}
}

func rdsSubscriptionLocalSubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"publication_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publication_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRdsSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
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
	createOpt.JSONBody = utils.RemoveNil(buildCreateRdsSubscriptionBodyParams(d))

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
		return diag.Errorf("error creating RDS subscription: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	subscriptions, err := getInstanceSubscriptions(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.Errorf("error creating RDS subscription: %s", err)
	}

	publicationId := d.Get("local_subscription").([]interface{})[0].(map[string]interface{})["publication_id"].(string)
	searchPath := fmt.Sprintf("subscriptions[?publication_id=='%s'&&subscription_database=='%s']|[0].id", publicationId,
		d.Get("subscription_database").(string))
	subscriptionId := utils.PathSearch(searchPath, subscriptions, "").(string)
	if subscriptionId == "" {
		return diag.Errorf("error creating RDS subscription: ID is not found in API response")
	}

	d.SetId(subscriptionId)

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating RDS subscription: job_id is not found in API response")
	}

	if err = checkRDSInstanceJobFinish(client, jobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error creating RDS(%s) subscription: %s", d.Id(), err)
	}

	return resourceRdsSubscriptionRead(ctx, d, meta)
}

func buildCreateRdsSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subscriptions": buildCreateRdsSubscriptionSubscriptionsBodyParams(d),
	}
	return bodyParams
}

func buildCreateRdsSubscriptionSubscriptionsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	subscription := map[string]interface{}{
		"subscription_database": d.Get("subscription_database"),
		"subscription_type":     d.Get("subscription_type"),
		"initialize_at":         d.Get("initialize_at"),
		"job_schedule":          buildRdsSubscriptionJobScheduleBodyParams(d.Get("job_schedule")),
		"initialize_info":       buildRdsSubscriptionInitializeInfoBodyParams(d.Get("initialize_info")),
		"bak_file_name":         utils.ValueIgnoreEmpty(d.Get("bak_file_name")),
		"bak_bucket_name":       utils.ValueIgnoreEmpty(d.Get("bak_bucket_name")),
		"local_subscription":    buildRdsSubscriptionLocalSubscriptionBodyParams(d.Get("local_subscription")),
	}

	if v, ok := d.GetOk("independent_agent"); ok {
		subscription["independent_agent"] = v == "true"
	}

	return []map[string]interface{}{subscription}
}

func buildRdsSubscriptionJobScheduleBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"job_schedule_type":   utils.ValueIgnoreEmpty(v["job_schedule_type"]),
			"one_time_occurrence": buildRdsSubscriptionOneTimeOccurrenceBodyParams(v["one_time_occurrence"]),
			"frequency":           buildRdsSubscriptionFrequencyBodyParams(v["frequency"]),
			"daily_frequency":     buildRdsSubscriptionDailyFrequencyBodyParams(v["daily_frequency"]),
			"duration":            buildRdsSubscriptionDurationBodyParams(v["duration"]),
		}
		return bodyParams
	}

	return nil
}

func buildRdsSubscriptionOneTimeOccurrenceBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
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

func buildRdsSubscriptionFrequencyBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
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

func buildRdsSubscriptionDailyFrequencyBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
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

func buildRdsSubscriptionDurationBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
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

func buildRdsSubscriptionInitializeInfoBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"file_source": utils.ValueIgnoreEmpty(v["file_source"]),
			"backup_id":   utils.ValueIgnoreEmpty(v["backup_id"]),
			"bucket_name": utils.ValueIgnoreEmpty(v["bucket_name"]),
			"file_path":   utils.ValueIgnoreEmpty(v["file_path"]),
			"file_name":   utils.ValueIgnoreEmpty(v["file_name"]),
		}
		if overwriteRestore, ok := v["overwrite_restore"]; ok {
			bodyParams["overwrite_restore"] = overwriteRestore == "true"
		}
		return bodyParams
	}

	return nil
}

func buildRdsSubscriptionLocalSubscriptionBodyParams(raw interface{}) map[string]interface{} {
	if raw == nil {
		return nil
	}
	rawParams := raw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if v, ok := rawParams[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"publication_id":   v["publication_id"],
			"publication_name": utils.ValueIgnoreEmpty(v["publication_name"]),
		}
		return bodyParams
	}
	return nil
}

func resourceRdsSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	subscriptions, err := getInstanceSubscriptions(client, d.Get("instance_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS subscription")
	}

	subscription := utils.PathSearch(fmt.Sprintf("subscriptions[?id=='%s']|[0]", d.Id()), subscriptions, nil)
	if subscription == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving RDS subscription")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subscription_database", utils.PathSearch("subscription_database", subscription, nil)),
		d.Set("subscription_type", utils.PathSearch("subscription_type", subscription, nil)),
		d.Set("job_schedule", flattenRdsSubscriptionJobSchedule(subscription)),
		d.Set("local_subscription", flattenRdsSubscriptionLocalSubscription(subscription)),
		d.Set("status", utils.PathSearch("status", subscription, nil)),
		d.Set("is_cloud", strconv.FormatBool(utils.PathSearch("is_cloud", subscription, false).(bool))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstanceSubscriptions(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}

	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func flattenRdsSubscriptionJobSchedule(resp interface{}) []interface{} {
	curJson := utils.PathSearch("job_schedule", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                  utils.PathSearch("id", curJson, nil),
			"job_schedule_type":   utils.PathSearch("job_schedule_type", curJson, nil),
			"one_time_occurrence": flattenRdsSubscriptionOneTimeOccurrence(curJson),
			"frequency":           flattenRdsSubscriptionFrequency(curJson),
			"daily_frequency":     flattenRdsSubscriptionDailyFrequency(curJson),
			"duration":            flattenRdsSubscriptionDuration(curJson),
		},
	}
}

func flattenRdsSubscriptionOneTimeOccurrence(resp interface{}) []interface{} {
	curJson := utils.PathSearch("one_time_occurrence", resp, nil)
	if curJson == nil || len(curJson.(map[string]interface{})) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", curJson, nil),
			"active_start_time": utils.PathSearch("active_start_time", curJson, nil),
		},
	}
}

func flattenRdsSubscriptionFrequency(resp interface{}) []interface{} {
	curJson := utils.PathSearch("frequency", resp, nil)
	if curJson == nil || len(curJson.(map[string]interface{})) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"freq_type":                      utils.PathSearch("freq_type", curJson, nil),
			"freq_interval":                  utils.PathSearch("freq_interval", curJson, nil),
			"freq_interval_weekly":           utils.PathSearch("freq_interval_weekly", curJson, nil),
			"freq_interval_day_monthly":      utils.PathSearch("freq_interval_day_monthly", curJson, nil),
			"freq_interval_monthly":          utils.PathSearch("freq_interval_monthly", curJson, nil),
			"freq_relative_interval_monthly": utils.PathSearch("freq_relative_interval_monthly", curJson, nil),
		},
	}
}

func flattenRdsSubscriptionDailyFrequency(resp interface{}) []interface{} {
	curJson := utils.PathSearch("daily_frequency", resp, nil)
	if curJson == nil || len(curJson.(map[string]interface{})) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"freq_subday_type":     utils.PathSearch("freq_subday_type", curJson, nil),
			"active_start_time":    utils.PathSearch("active_start_time", curJson, nil),
			"active_end_time":      utils.PathSearch("active_end_time", curJson, nil),
			"freq_subday_interval": utils.PathSearch("freq_subday_interval", curJson, nil),
			"freq_interval_unit":   utils.PathSearch("freq_interval_unit", curJson, nil),
		},
	}
}

func flattenRdsSubscriptionDuration(resp interface{}) []interface{} {
	curJson := utils.PathSearch("duration", resp, nil)
	if curJson == nil || len(curJson.(map[string]interface{})) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"active_start_date": utils.PathSearch("active_start_date", curJson, nil),
			"active_end_date":   utils.PathSearch("active_end_date", curJson, nil),
		},
	}
}

func flattenRdsSubscriptionLocalSubscription(resp interface{}) []interface{} {
	curJson := utils.PathSearch("local_subscription", resp, nil)
	if curJson == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"publication_id":            utils.PathSearch("publication_id", resp, nil),
			"publication_name":          utils.PathSearch("publication_name", resp, nil),
			"publication_instance_id":   utils.PathSearch("publication_instance_id", curJson, nil),
			"publication_instance_name": utils.PathSearch("publication_instance_name", curJson, nil),
		},
	}
}

func resourceRdsSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateRdsSubscriptionJobScheduleBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating RDS subscription: %s", err)
	}

	return resourceRdsSubscriptionRead(ctx, d, meta)
}

func buildUpdateRdsSubscriptionJobScheduleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subscription_ids": []string{d.Id()},
		"job_schedule":     buildRdsSubscriptionJobScheduleBodyParams(d.Get("job_schedule")),
	}
	return bodyParams
}

func resourceRdsSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
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
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteRdsSubscriptionBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		// if an instance is deleted and then delete a subscription, an operation conflict will be reported.
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
		return common.CheckDeletedDiag(d, err, "error deleting RDS subscription")
	}

	deleteRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting RDS subscription: job_id is not found in API response")
	}

	if err = checkRDSInstanceJobFinish(client, jobId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error deleting RDS subscription(%s): %s", d.Id(), err)
	}

	return nil
}

func buildDeleteRdsSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subscription_ids": []string{d.Id()},
	}
	return bodyParams
}

func resourceRdsSubscriptionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
