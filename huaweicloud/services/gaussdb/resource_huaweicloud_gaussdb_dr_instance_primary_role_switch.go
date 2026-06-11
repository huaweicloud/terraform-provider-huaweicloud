package gaussdb

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

var gaussDbDrInstancePrimaryRoleSwitchNonUpdatableParams = []string{
	"instance_id",
	"disaster_type",
	"post_process_config",
}

// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/switchover
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbDrInstancePrimaryRoleSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbDrInstancePrimaryRoleSwitchCreate,
		ReadContext:   resourceGaussDbDrInstancePrimaryRoleSwitchRead,
		UpdateContext: resourceGaussDbDrInstancePrimaryRoleSwitchUpdate,
		DeleteContext: resourceGaussDbDrInstancePrimaryRoleSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbDrInstancePrimaryRoleSwitchNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
			"disaster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"post_process_config": {
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

func resourceGaussDbDrInstancePrimaryRoleSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/switchover"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildDrInstancePrimaryRoleSwitchBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error switching GaussDB DR instance primary role: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job_id from the API response")
	}

	d.SetId(d.Get("instance_id").(string))

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating GaussDB DR instance primary role: %s", err)
	}

	return nil
}

func buildDrInstancePrimaryRoleSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disaster_type":       d.Get("disaster_type"),
		"post_process_config": utils.ValueIgnoreEmpty(d.Get("post_process_config")),
	}
	return bodyParams
}

func resourceGaussDbDrInstancePrimaryRoleSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrInstancePrimaryRoleSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrInstancePrimaryRoleSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB DR instance primary role switch resource is not supported. The resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
