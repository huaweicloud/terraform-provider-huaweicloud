package dds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ddsInstanceNodeNumUpdateNonUpdatableParams = []string{
	"instance_id", "type", "spec_code", "num", "volume", "volume.*.size",
}

// @API DDS POST /v3/{project_id}/instances/{instance_id}/enlarge
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
func ResourceDdsInstanceNodeNumUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsInstanceNodeNumUpdateCreate,
		ReadContext:   resourceDdsInstanceNodeNumUpdateRead,
		UpdateContext: resourceDdsInstanceNodeNumUpdateUpdate,
		DeleteContext: resourceDdsInstanceNodeNumUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(ddsInstanceNodeNumUpdateNonUpdatableParams),

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"num": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Required: true,
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
		},
	}
}

func resourceDdsInstanceNodeNumUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/enlarge"
		product = "dds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createBody := buildCreateDdsInstanceNodeNumUpdateBodyParams(d)
	createOpt.JSONBody = utils.RemoveNil(createBody)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating DDS instance node num update: %s", err)
	}

	d.SetId(instanceId)

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if jobId == "" && orderId == "" {
		return diag.Errorf("error creating DDS instance node num update: job_id or order_id is not found in API response")
	}

	if orderId != "" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if jobId != "" {
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, jobId),
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        60 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", jobId, err)
		}
	}

	err = waitForInstanceReady(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildCreateDdsInstanceNodeNumUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"spec_code":   d.Get("spec_code"),
		"num":         d.Get("num"),
		"volume":      buildCreateDdsInstanceNodeVolumeUpdateBodyParams(d),
		"is_auto_pay": true,
	}
	return bodyParams
}

func buildCreateDdsInstanceNodeVolumeUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	if volumeRaw, ok := d.GetOk("volume"); ok {
		if v, ok := volumeRaw.([]interface{})[0].(map[string]interface{}); ok {
			params := map[string]interface{}{
				"size": v["size"].(string),
			}
			return params
		}
	}
	return nil
}

func resourceDdsInstanceNodeNumUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdsInstanceNodeNumUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdsInstanceNodeNumUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DDS instance node num update resource is not supported. The resource is only removed from the " +
		"state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
