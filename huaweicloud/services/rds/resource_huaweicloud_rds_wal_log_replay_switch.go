package rds

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

var walLogReplaySwitchNonUpdatableParams = []string{"instance_id", "pause_log_replay", "enable_force_new"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/log-replay/update
// @API RDS GET /v3/{project_id}/instances
func ResourceRdsWalLogReplaySwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWalLogReplaySwitchCreate,
		UpdateContext: resourceWalLogReplaySwitchUpdate,
		ReadContext:   resourceWalLogReplaySwitchRead,
		DeleteContext: resourceWalLogReplaySwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(walLogReplaySwitchNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
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
			"pause_log_replay": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceWalLogReplaySwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v3/{project_id}/instances/{instance_id}/log-replay/update"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	requestBody := map[string]interface{}{
		"pause_log_replay": d.Get("pause_log_replay").(string),
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
			"x-language":   "en-us",
		},
		JSONBody: requestBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating WAL replay status for RDS instance(%s): %s", instanceID, err)
	}

	res, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	message := utils.PathSearch("message", res, "")
	if message.(string) != "operate successfully" {
		return diag.Errorf("error pausing/resuming WAL log replay for RDS instance(%s), response: %v", instanceID, message)
	}

	d.SetId(instanceID)
	return nil
}

func resourceWalLogReplaySwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWalLogReplaySwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWalLogReplaySwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting WAL replay resource is not supported. The resource is only removed from state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
