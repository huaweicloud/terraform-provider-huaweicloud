package rds

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var unlockNodeReadonlyStatusNonUpdatableParams = []string{"instance_id", "status_preservation_time"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/unlock-node-readonly-status
// @API RDS GET /v3/{project_id}/instances
func ResourceUnlockNodeReadonlyStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnlockNodeReadonlyStatusCreate,
		ReadContext:   resourceUnlockNodeReadonlyStatusRead,
		UpdateContext: resourceUnlockNodeReadonlyStatusUpdate,
		DeleteContext: resourceUnlockNodeReadonlyStatusDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(unlockNodeReadonlyStatusNonUpdatableParams),

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
			"status_preservation_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceUnlockNodeReadonlyStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/unlock-node-readonly-status"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateUnlockNodeReadonlyStatusBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS unlock node readonly status: %s", err)
	}

	d.SetId(d.Get("instance_id").(string))

	return nil
}

func buildCreateUnlockNodeReadonlyStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status_preservation_time": d.Get("status_preservation_time"),
	}
	return bodyParams
}

func resourceUnlockNodeReadonlyStatusRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUnlockNodeReadonlyStatusUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUnlockNodeReadonlyStatusDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS unlock node readonly status resource is not supported. The resource is only removed " +
		"from the state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
