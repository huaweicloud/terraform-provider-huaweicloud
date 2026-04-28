package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/subscriptions
func DataSourceRdsSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSubscriptionsRead,

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
			"publication_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_cloud": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"publication_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscription_db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionsSchema(),
			},
		},
	}
}

func subscriptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_cloud": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subscription_database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_subscription": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationSubscriptionSchema(),
			},
			"local_subscription": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     localSubscriptionSchema(),
			},
			"job_schedule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobScheduleSchema(),
			},
		},
	}
}

func publicationSubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"subscription_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_instance_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subscription_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func localSubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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

func jobScheduleSchema() *schema.Resource {
	return &schema.Resource{
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
				Elem:     oneTimeOccurrenceSchema(),
			},
			"frequency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     frequencySchema(),
			},
			"daily_frequency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dailyFrequencySchema(),
			},
			"duration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     durationSchema(),
			},
		},
	}
}

func oneTimeOccurrenceSchema() *schema.Resource {
	return &schema.Resource{
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
}

func frequencySchema() *schema.Resource {
	return &schema.Resource{
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
}

func dailyFrequencySchema() *schema.Resource {
	return &schema.Resource{
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
}

func durationSchema() *schema.Resource {
	return &schema.Resource{
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
}

func dataSourceRdsSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/replication/subscriptions"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildListSubscriptionsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving RDS subscriptions: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
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
		d.Set("subscriptions", flattenListSubscriptionsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListSubscriptionsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"
	if v, ok := d.GetOk("publication_id"); ok {
		queryParams += fmt.Sprintf("&publication_id=%v", v)
	}
	if v, ok := d.GetOk("is_cloud"); ok {
		isCloud, _ := strconv.ParseBool(v.(string))
		queryParams += fmt.Sprintf("&is_cloud=%v", isCloud)
	}
	if v, ok := d.GetOk("publication_name"); ok {
		queryParams += fmt.Sprintf("&publication_name=%v", v)
	}
	if v, ok := d.GetOk("subscription_db_name"); ok {
		queryParams += fmt.Sprintf("&subscription_db_name=%v", v)
	}
	return queryParams
}

func flattenListSubscriptionsBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	subscriptions := utils.PathSearch("subscriptions", resp, make([]interface{}, 0)).([]interface{})
	if len(subscriptions) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = map[string]interface{}{
			"id":                       utils.PathSearch("id", sub, nil),
			"status":                   utils.PathSearch("status", sub, nil),
			"publication_id":           utils.PathSearch("publication_id", sub, nil),
			"publication_name":         utils.PathSearch("publication_name", sub, nil),
			"is_cloud":                 utils.PathSearch("is_cloud", sub, nil),
			"subscription_database":    utils.PathSearch("subscription_database", sub, nil),
			"subscription_type":        utils.PathSearch("subscription_type", sub, nil),
			"publication_subscription": flattenPublicationSubscription(sub),
			"local_subscription":       flattenLocalSubscription(sub),
			"job_schedule":             flattenJobSchedule(sub),
		}
	}

	return result
}

func flattenPublicationSubscription(sub interface{}) []map[string]interface{} {
	ps := utils.PathSearch("publication_subscription", sub, nil)
	if ps == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"subscription_instance_name": utils.PathSearch("subscription_instance_name", ps, nil),
			"subscription_instance_ip":   utils.PathSearch("subscription_instance_ip", ps, nil),
			"subscription_instance_id":   utils.PathSearch("subscription_instance_id", ps, nil),
		},
	}
}

func flattenLocalSubscription(sub interface{}) []map[string]interface{} {
	ls := utils.PathSearch("local_subscription", sub, nil)
	if ls == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"publication_instance_id":   utils.PathSearch("publication_instance_id", ls, nil),
			"publication_instance_name": utils.PathSearch("publication_instance_name", ls, nil),
		},
	}
}

func flattenJobSchedule(sub interface{}) []map[string]interface{} {
	js := utils.PathSearch("job_schedule", sub, nil)
	if js == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                  utils.PathSearch("id", js, nil),
			"job_schedule_type":   utils.PathSearch("job_schedule_type", js, nil),
			"one_time_occurrence": flattenOneTimeOccurrence(js),
			"frequency":           flattenFrequency(js),
			"daily_frequency":     flattenDailyFrequency(js),
			"duration":            flattenDuration(js),
		},
	}
}

func flattenOneTimeOccurrence(js interface{}) []map[string]interface{} {
	oto := utils.PathSearch("one_time_occurrence", js, nil)
	if oto == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"active_start_date": utils.PathSearch("active_start_date", oto, nil),
			"active_start_time": utils.PathSearch("active_start_time", oto, nil),
		},
	}
}

func flattenFrequency(js interface{}) []map[string]interface{} {
	freq := utils.PathSearch("frequency", js, nil)
	if freq == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"freq_type":                      utils.PathSearch("freq_type", freq, nil),
			"freq_interval":                  utils.PathSearch("freq_interval", freq, nil),
			"freq_interval_weekly":           utils.PathSearch("freq_interval_weekly", freq, nil),
			"freq_interval_day_monthly":      utils.PathSearch("freq_interval_day_monthly", freq, nil),
			"freq_interval_monthly":          utils.PathSearch("freq_interval_monthly", freq, nil),
			"freq_relative_interval_monthly": utils.PathSearch("freq_relative_interval_monthly", freq, nil),
		},
	}
}

func flattenDailyFrequency(js interface{}) []map[string]interface{} {
	df := utils.PathSearch("daily_frequency", js, nil)
	if df == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"freq_subday_type":     utils.PathSearch("freq_subday_type", df, nil),
			"active_start_time":    utils.PathSearch("active_start_time", df, nil),
			"active_end_time":      utils.PathSearch("active_end_time", df, nil),
			"freq_subday_interval": utils.PathSearch("freq_subday_interval", df, nil),
			"freq_interval_unit":   utils.PathSearch("freq_interval_unit", df, nil),
		},
	}
}

func flattenDuration(js interface{}) []map[string]interface{} {
	dur := utils.PathSearch("duration", js, nil)
	if dur == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"active_start_date": utils.PathSearch("active_start_date", dur, nil),
			"active_end_date":   utils.PathSearch("active_end_date", dur, nil),
		},
	}
}
