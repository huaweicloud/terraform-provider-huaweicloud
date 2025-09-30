package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var topicQuotaNonUpdatableParams = []string{"instance_id", "topic"}

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota
// @API Kafka GET /v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota
// @API Kafka DELETE /v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
func ResourceTopicQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicQuotaCreate,
		ReadContext:   resourceTopicQuotaRead,
		UpdateContext: resourceTopicQuotaUpdate,
		DeleteContext: resourceTopicQuotaDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTopicQuotaImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(topicQuotaNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region where the topic quota is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the topic quota belongs.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic.`,
			},
			"producer_byte_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The producer rate limit. The unit is byte/s.`,
			},
			"consumer_byte_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The consumer rate limit. The unit is byte/s.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildTopicQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topic":              d.Get("topic"),
		"producer-byte-rate": utils.ValueIgnoreEmpty(d.Get("producer_byte_rate")),
		"consumer-byte-rate": utils.ValueIgnoreEmpty(d.Get("consumer_byte_rate")),
	}
}

func resourceTopicQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		topic      = d.Get("topic").(string)
		httpUrl    = "v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota"
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildTopicQuotaBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating topic quota of the instance (%s): %s", instanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the topic quota of the instance (%s) to be created: %s", instanceId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceId, topic))

	return resourceTopicQuotaRead(ctx, d, meta)
}

func getTopicQuotas(client *golangsdk.ServiceClient, instanceId, topicName string) (interface{}, error) {
	var (
		httpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota"
		result  = make([]interface{}, 0)
		offset  = 0
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	// The `limit` default is `10` and `keyword` is fuzzy match.
	listPath = fmt.Sprintf("%s?type=topic&limit=100", listPath)
	if topicName != "" {
		listPath = fmt.Sprintf("%s&keyword=%s", listPath, topicName)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		quotas := utils.PathSearch("quotas", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, quotas...)
		offset += len(quotas)
		// offset cannot be greater than or equal to the total number of quotas.
		if offset >= int(utils.PathSearch("count", respBody, float64(0)).(float64)) {
			break
		}
	}

	return result, nil
}

func GetTopicQuotaByTopicName(client *golangsdk.ServiceClient, instanceId, topicName string) (interface{}, error) {
	quotas, err := getTopicQuotas(client, instanceId, topicName)
	if err != nil {
		return nil, err
	}

	topicQuota := utils.PathSearch(fmt.Sprintf("[?topic=='%s']|[0]", topicName), quotas, nil)
	if topicQuota == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The topic (%s) quota of the instance (%s) was not found.", topicName, instanceId)),
			},
		}
	}
	return topicQuota, nil
}

func resourceTopicQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		topicName  = d.Get("topic").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := GetTopicQuotaByTopicName(client, instanceId, topicName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting topic (%s) quota of the instance (%s)", topicName, instanceId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topic", utils.PathSearch("topic", respBody, nil)),
		d.Set("producer_byte_rate", utils.PathSearch(`"producer-byte-rate"`, respBody, nil)),
		d.Set("consumer_byte_rate", utils.PathSearch(`"consumer-byte-rate"`, respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTopicQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota"
		topic      = d.Get("topic").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	updatePath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildTopicQuotaBodyParams(d)),
	}

	resp, err := client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating topic (%s) quota of the instance (%s): %s", topic, instanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error waiting for the topic (%s) quota of the instance (%s) to be updated: %s", topic, instanceId, err)
	}

	return resourceTopicQuotaRead(ctx, d, meta)
}

func resourceTopicQuotaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/kafka/{project_id}/instances/{instance_id}/kafka-topic-quota"
		topicName  = d.Get("topic").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	// When the topic quota does not exist, the delete API will return 200, but the status in the query task API response body is "FAILED".
	// So we need to check if the topic quota exists before deleting it.
	_, err = GetTopicQuotaByTopicName(client, instanceId, topicName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("The topic (%s) quota of the instance (%s) was not found",
			topicName, instanceId))
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	deletePath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: map[string]interface{}{
			"topic": topicName,
		},
	}

	resp, err := client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting topic (%s) quota of the instance (%s)",
			topicName, instanceId))
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the topic (%s) quota of the instance (%s) to be deleted: %s",
			topicName, instanceId, err)
	}

	return nil
}

func resourceTopicQuotaImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<topic>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("topic", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
