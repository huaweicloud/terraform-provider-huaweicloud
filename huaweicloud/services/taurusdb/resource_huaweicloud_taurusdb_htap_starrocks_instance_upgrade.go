package taurusdb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var starrocksInstanceUpgradeNoneUpdatableParams = []string{
	"instance_id", "starrocks_instance_id", "delay", "is_skip_validate",
}

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/db-upgrade
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBHtapStarrocksInstanceUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksInstanceUpgradeCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksInstanceUpgradeRead,
		UpdateContext: resourceTaurusDBHtapStarrocksInstanceUpgradeUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksInstanceUpgradeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(starrocksInstanceUpgradeNoneUpdatableParams),

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
			"starrocks_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delay": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"is_skip_validate": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
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

func resourceTaurusDBHtapStarrocksInstanceUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	jobId, err := upgradeStarrocksInstance(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for upgrading StarRocks instance job (%s) to complete: %s", jobId, err)
	}
	return nil
}

func upgradeStarrocksInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/db-upgrade"
	)

	instanceId := d.Get("instance_id").(string)
	starrocksInstanceId := d.Get("starrocks_instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildUpgradeStarrocksInstanceBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     htapInstanceStateRefreshFunc(client, instanceId, starrocksInstanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return "", fmt.Errorf("error upgrading TaurusDB Htap StarRocks instance(%s): %s", starrocksInstanceId, err)
	}
	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return "", err
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return "", fmt.Errorf("error upgrading StarRocks instance(%s), job_id is not found in the response",
			starrocksInstanceId)
	}

	return jobId, nil
}

func buildUpgradeStarrocksInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"delay":            d.Get("delay"),
		"is_skip_validate": d.Get("is_skip_validate"),
	}
	return bodyParams
}

func resourceTaurusDBHtapStarrocksInstanceUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksInstanceUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksInstanceUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting upgrade resource is not supported. The upgrade resource is only removed from the state," +
		" the StarRocks instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
