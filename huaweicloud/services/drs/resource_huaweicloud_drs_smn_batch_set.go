package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsSmnBatchSetNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.status",
	"jobs.*.engine_type",
	"alarm_notify_info",
	"alarm_notify_info.*.topic_urn",
	"alarm_notify_info.*.delay_time",
	"alarm_notify_info.*.rto_delay",
	"alarm_notify_info.*.rpo_delay",
	"alarm_notify_info.*.alarm_to_user",
	"alarm_notify_info.*.subscriptions.*.protocol",
	"alarm_notify_info.*.subscriptions.*.endpoints",
}

// @API DRS POST /v3/{project_id}/jobs/batch-set-smn
func ResourceDrsSmnBatchSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsSmnBatchSetCreate,
		ReadContext:   resourceDrsSmnBatchSetRead,
		UpdateContext: resourceDrsSmnBatchSetUpdate,
		DeleteContext: resourceDrsSmnBatchSetDelete,

		CustomizeDiff: config.FlexibleForceNew(drsSmnBatchSetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"alarm_notify_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
									},
									"endpoints": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
									},
								},
							},
						},
						"topic_urn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delay_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rto_delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rpo_delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"alarm_to_user": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDrsSmnBatchSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"jobs":              buildDrsSmnBatchSetJobsBodyParams(d),
		"alarm_notify_info": buildDrsSmnBatchSetAlarmNotifyInfoBodyParams(d),
	}
	return bodyParams
}

func buildDrsSmnBatchSetJobsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	jobsRaw := d.Get("jobs").([]interface{})
	if len(jobsRaw) == 0 {
		return nil
	}

	jobs := make([]map[string]interface{}, 0, len(jobsRaw))
	for _, jobRaw := range jobsRaw {
		job, ok := jobRaw.(map[string]interface{})
		if !ok {
			return nil
		}
		jobs = append(jobs, map[string]interface{}{
			"job_id":      job["job_id"],
			"status":      job["status"],
			"engine_type": job["engine_type"],
		})
	}
	return jobs
}

func buildDrsSmnBatchSetAlarmNotifyInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	alarmNotifyInfoRaw := d.Get("alarm_notify_info").([]interface{})
	if len(alarmNotifyInfoRaw) == 0 {
		return nil
	}

	alarmNotifyInfo, ok := alarmNotifyInfoRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}

	subscriptionsRaw := alarmNotifyInfo["subscriptions"].([]interface{})
	subscriptions := make([]map[string]interface{}, 0, len(subscriptionsRaw))
	for _, subRaw := range subscriptionsRaw {
		sub, ok := subRaw.(map[string]interface{})
		if !ok {
			return nil
		}
		endpointsRaw := sub["endpoints"].([]interface{})
		endpoints := make([]string, 0, len(endpointsRaw))
		for _, ep := range endpointsRaw {
			endpoints = append(endpoints, ep.(string))
		}

		subscriptions = append(subscriptions, map[string]interface{}{
			"protocol":  sub["protocol"],
			"endpoints": endpoints,
		})
	}

	bodyParams := map[string]interface{}{
		"topic_urn":     alarmNotifyInfo["topic_urn"],
		"delay_time":    alarmNotifyInfo["delay_time"],
		"rto_delay":     alarmNotifyInfo["rto_delay"],
		"rpo_delay":     alarmNotifyInfo["rpo_delay"],
		"alarm_to_user": alarmNotifyInfo["alarm_to_user"],
	}

	if len(subscriptions) > 0 {
		bodyParams["subscriptions"] = subscriptions
	}

	return bodyParams
}

func resourceDrsSmnBatchSetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/batch-set-smn"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDrsSmnBatchSetBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch setting SMN for DRS jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(
		nil,
		d.Set("results", flattenDrsSmnBatchSetResults(respBody)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS batch set SMN fields: %s", mErr)
	}

	return nil
}

func flattenDrsSmnBatchSetResults(respBody interface{}) []interface{} {
	resultsRaw := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(resultsRaw) == 0 {
		return nil
	}

	results := make([]interface{}, 0, len(resultsRaw))
	for _, result := range resultsRaw {
		results = append(results, map[string]interface{}{
			"id":     utils.PathSearch("id", result, nil),
			"status": utils.PathSearch("status", result, nil),
		})
	}

	return results
}

func resourceDrsSmnBatchSetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDrsSmnBatchSetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDrsSmnBatchSetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch set SMN for DRS jobs. Deleting this resource
    will not undo the SMN configuration, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
