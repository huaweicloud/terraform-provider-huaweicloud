package geminidb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceRestartNonUpdatableParams = []string{
	"instance_id",
	"node_id",
}

// @API GeminiDB POST /v3/{project_id}/instances/{instance_id}/restart
// @API GeminiDB POST /v3/{project_id}/instances
// @API GeminiDB GET /v3/{project_id}/jobs
func ResourceInstanceRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceRestartCreate,
		ReadContext:   resourceInstanceRestartRead,
		UpdateContext: resourceInstanceRestartUpdate,
		DeleteContext: resourceInstanceRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(instanceRestartNonUpdatableParams),

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildInstanceRestartBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id": utils.ValueIgnoreEmpty(d.Get("node_id")),
	}

	return bodyParams
}

func resourceInstanceRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	respBody, err := restartInstanceOrNode(ctx, client, d, instanceId)
	if err != nil {
		return diag.Errorf("error restarting GeminiDb instance or node: %s", err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error restarting instance or node: unable to find job ID from API response")
	}

	d.SetId(jobId)

	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for instance or node to restart to complete: %s", err)
	}

	return nil
}

func restartInstanceOrNode(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/restart"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildInstanceRestartBodyParams(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		r, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return r, retry, err
	}
	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp.(*http.Response))
}

func resourceInstanceRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceInstanceRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceInstanceRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB restart instance or node resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
