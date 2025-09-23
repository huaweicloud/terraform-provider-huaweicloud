package gaussdb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB PUT /v3.1/{project_id}/instances/{instance_id}/db-upgrade
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceOpenGaussInstanceUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussInstanceUpgradeCreate,
		ReadContext:   resourceOpenGaussInstanceUpgradeRead,
		DeleteContext: resourceOpenGaussInstanceUpgradeDelete,

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
				ForceNew: true,
			},
			"upgrade_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"upgrade_action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"upgrade_shard_num": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"upgrade_az": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOpenGaussInstanceUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-upgrade"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussInstanceUpgradeBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error upgrading GaussDB OpenGauss instance(%s): %s", instanceId, err)
	}

	d.SetId(instanceId)

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", createRespBody, nil)
	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if orderId == nil && jobId == nil {
		return diag.Errorf("error upgrading GaussDB OpenGauss instance flavor: order_id and job_id is not found" +
			" in API response")
	}
	if orderId != nil {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		// wait for order success
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if jobId != nil {
		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 10, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func buildCreateOpenGaussInstanceUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"upgrade_type":      d.Get("upgrade_type"),
		"upgrade_action":    utils.ValueIgnoreEmpty(d.Get("upgrade_action")),
		"target_version":    utils.ValueIgnoreEmpty(d.Get("target_version")),
		"upgrade_shard_num": utils.ValueIgnoreEmpty(d.Get("upgrade_shard_num")),
		"upgrade_az":        utils.ValueIgnoreEmpty(d.Get("upgrade_az")),
	}
	return bodyParams
}

func resourceOpenGaussInstanceUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussInstanceUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting upgrade resource is not supported. The upgrade resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
