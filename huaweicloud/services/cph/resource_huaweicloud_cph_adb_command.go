package cph

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var commandNonUpdatableParams = []string{
	"command",
	"content",
	"phone_ids",
	"server_ids",
}

// @API CPH POST /v1/{project_id}/cloud-phone/phones/commands
// @API CPH GET /v1/{project_id}/cloud-phone/jobs
func ResourceAdbCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdbCommandCreate,
		UpdateContext: resourceAdbCommandUpdate,
		ReadContext:   resourceAdbCommandRead,
		DeleteContext: resourceAdbCommandDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(commandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ADB command.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the content.`,
			},
			"phone_ids": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"phone_ids", "server_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `Specifies the IDs of the CPH phone.`,
			},
			"server_ids": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"phone_ids", "server_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `Specifies the IDs of CPH server.`,
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

func resourceAdbCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// createAdbCommand: create CPH adb command
	createAdbCommandHttpUrl := "v1/{project_id}/cloud-phone/phones/commands"
	createAdbCommandPath := client.Endpoint + createAdbCommandHttpUrl
	createAdbCommandPath = strings.ReplaceAll(createAdbCommandPath, "{project_id}", client.ProjectID)

	createAdbCommandOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createAdbCommandOpt.JSONBody = utils.RemoveNil(buildcreateAdbCommandBodyParams(d))
	createAdbCommandResp, err := client.Request("POST", createAdbCommandPath, &createAdbCommandOpt)
	if err != nil {
		return diag.Errorf("error creating CPH adb command: %s", err)
	}

	createAdbCommandRespBody, err := utils.FlattenResponse(createAdbCommandResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("request_id", createAdbCommandRespBody, "").(string)
	if id == "" {
		return diag.Errorf("Unable to find the request ID from the API response")
	}
	d.SetId(id)

	phoneIds := utils.PathSearch("jobs[?error_code=='CPS.0005'].phone_id", createAdbCommandRespBody, make([]interface{}, 0)).([]interface{})
	serverIds := utils.PathSearch("jobs[?error_code=='CPS.0097'].server_id", createAdbCommandRespBody, make([]interface{}, 0)).([]interface{})
	if len(phoneIds) > 0 {
		log.Printf("[WARN] Phone not found: %v", phoneIds)
	}
	if len(serverIds) > 0 {
		log.Printf("[WARN] Parameter: server_id is invalid: %v", serverIds)
	}

	err = checkJobStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("error waiting for CPH adb job completed: %s", err),
			},
		}
	}

	return nil
}

func buildcreateAdbCommandBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"command":    d.Get("command"),
		"content":    d.Get("content"),
		"phone_ids":  utils.ValueIgnoreEmpty(d.Get("phone_ids")),
		"server_ids": utils.ValueIgnoreEmpty(d.Get("server_ids")),
	}

	return bodyParams
}

func resourceAdbCommandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAdbCommandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAdbCommandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting adb command resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkJobStatus(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStateRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func jobStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobsHttpUrl := "v1/{project_id}/cloud-phone/jobs?request_id={request_id}"
		getJobsPath := client.Endpoint + getJobsHttpUrl
		getJobsPath = strings.ReplaceAll(getJobsPath, "{project_id}", client.ProjectID)
		getJobsPath = strings.ReplaceAll(getJobsPath, "{request_id}", id)

		getJobsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getJobsResp, err := client.Request("GET", getJobsPath, &getJobsOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		getJobsRespBody, err := utils.FlattenResponse(getJobsResp)
		if err != nil {
			return nil, "ERROR", err
		}

		// status is 1, indicates the job is Running
		runningStatus := utils.PathSearch("jobs[?status==`1`]", getJobsRespBody, make([]interface{}, 0)).([]interface{})
		// status is -1, indicates the job is failed
		failedPhoneIds := utils.PathSearch("jobs[?status==`-1`].phone_id", getJobsRespBody, make([]interface{}, 0)).([]interface{})
		if len(runningStatus) > 0 {
			return getJobsRespBody, "PENDING", nil
		}
		if len(failedPhoneIds) > 0 {
			return getJobsRespBody, "ERROR", fmt.Errorf("failed phone ids: %v", failedPhoneIds)
		}
		return getJobsRespBody, "COMPLETED", nil
	}
}
