package smn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	messageDetectionNonUpdatableParams = []string{"topic_urn", "protocol", "endpoint", "extension"}

	detectionResultMap = map[string]string{
		"0": "available",
		"1": "unexecuted",
		"2": "unavailable",
	}

	detectionNotFoundCodes = []string{"SMN.00013030"}
)

// @API SMN POST /v2/{project_id}/notifications/topics/{topic_urn}/detection
// @API SMN GET /v2/{project_id}/notifications/topics/{topic_urn}/detection/{task_id}
func ResourceMessageDetection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMessageDetectionCreate,
		UpdateContext: resourceMessageDetectionUpdate,
		ReadContext:   resourceMessageDetectionRead,
		DeleteContext: resourceMessageDetectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(messageDetectionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource identifier of a topic.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protocol type.`,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the endpoint address to be detected.`,
			},
			"extension": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the extended key/value for subscriptions over HTTP or HTTPS.`,
			},
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The message detection result.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceMessageDetectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	topicUrn := d.Get("topic_urn").(string)

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	// createMessageDetection: create SMN message detection
	createMessageDetectionHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/detection"
	createMessageDetectionPath := client.Endpoint + createMessageDetectionHttpUrl
	createMessageDetectionPath = strings.ReplaceAll(createMessageDetectionPath, "{project_id}", client.ProjectID)
	createMessageDetectionPath = strings.ReplaceAll(createMessageDetectionPath, "{topic_urn}", topicUrn)

	createMessageDetectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createMessageDetectionOpt.JSONBody = utils.RemoveNil(buildCreateMessageDetectionBodyParams(d))
	createMessageDetectionResp, err := client.Request("POST", createMessageDetectionPath, &createMessageDetectionOpt)
	if err != nil {
		return diag.Errorf("error creating SMN message detection: %s", err)
	}

	createMessageDetectionRespBody, err := utils.FlattenResponse(createMessageDetectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("task_id", createMessageDetectionRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMN message detection task: task ID is not found in API response")
	}
	d.SetId(id)

	// Check whether detection task has completed
	err = checkDetectionTaskCompleted(ctx, client, topicUrn, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMessageDetectionRead(ctx, d, meta)
}

func buildCreateMessageDetectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"protocol": d.Get("protocol"),
		"endpoint": d.Get("endpoint"),
	}

	if v, ok := d.GetOk("extension"); ok {
		bodyParams["extension"] = map[string]interface{}{
			"header": v,
		}
	}

	return bodyParams
}

func resourceMessageDetectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	topicUrn := d.Get("topic_urn").(string)

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}
	resp, err := getMessageDetectionDetail(client, topicUrn, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", detectionNotFoundCodes...),
			"error querying SMN message detection detail")
	}
	status := fmt.Sprintf("%v", utils.PathSearch("status", resp, float64(0)).(float64))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("result", detectionResultMap[status]),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMessageDetectionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMessageDetectionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting message detection resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkDetectionTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, topicUrn, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"COMPLETED"},
		Refresh:      detectionTaskStateRefreshFunc(client, topicUrn, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for SMN (%s) message detection task completed: %s", topicUrn, err)
	}
	return nil
}

func detectionTaskStateRefreshFunc(client *golangsdk.ServiceClient, topicUrn, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getMessageDetectionDetail(client, topicUrn, id)
		if err != nil {
			return resp, "ERROR", err
		}

		status := fmt.Sprintf("%v", utils.PathSearch("status", resp, float64(0)).(float64))
		if status == "1" {
			return resp, "RUNNING", nil
		}
		return resp, "COMPLETED", nil
	}
}

func getMessageDetectionDetail(client *golangsdk.ServiceClient, topicUrn, id string) (interface{}, error) {
	getMessageDetectionHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/detection/{task_id}"
	getMessageDetectionPath := client.Endpoint + getMessageDetectionHttpUrl
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{project_id}", client.ProjectID)
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{topic_urn}", topicUrn)
	getMessageDetectionPath = strings.ReplaceAll(getMessageDetectionPath, "{task_id}", id)

	getMessageDetectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getMessageDetectionResp, err := client.Request("GET", getMessageDetectionPath, &getMessageDetectionOpt)
	if err != nil {
		return nil, err
	}

	getMessageDetectionRespBody, err := utils.FlattenResponse(getMessageDetectionResp)
	if err != nil {
		return nil, err
	}

	return getMessageDetectionRespBody, nil
}
