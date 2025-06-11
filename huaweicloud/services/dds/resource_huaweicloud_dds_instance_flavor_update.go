package dds

import (
	"context"
	"encoding/json"
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

var ddsInstanceFlavorUpdateNonUpdatableParams = []string{
	"instance_id", "target_spec_code", "target_type", "target_id", "target_ids",
}

// @API DDS POST /v3/{project_id}/instances/{instance_id}/resize
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
func ResourceDdsInstanceFlavorUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsInstanceFlavorUpdateCreate,
		ReadContext:   resourceDdsInstanceFlavorUpdateRead,
		UpdateContext: resourceDdsInstanceFlavorUpdateUpdate,
		DeleteContext: resourceDdsInstanceFlavorUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(ddsInstanceFlavorUpdateNonUpdatableParams),

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
			"target_spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceDdsInstanceFlavorUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/resize"
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
	createBody := buildCreateDdsInstanceFlavorUpdateBodyParams(d)
	createOpt.JSONBody = utils.RemoveNil(createBody)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		if handleDowngradeError(err) {
			delete(createBody, "is_auto_pay")
			createOpt.JSONBody = utils.RemoveNil(createBody)
			return res, true, err
		}
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
		return diag.Errorf("error creating DDS instance flavor update: %s", err)
	}

	d.SetId(instanceId)

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if jobId == "" && orderId == "" {
		return diag.Errorf("error creating DDS instance flavor update: job_id or order_id is not found in API response")
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

func buildCreateDdsInstanceFlavorUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"target_type":      utils.ValueIgnoreEmpty(d.Get("target_type")),
		"target_id":        utils.ValueIgnoreEmpty(d.Get("target_id")),
		"target_ids":       utils.ValueIgnoreEmpty(d.Get("target_ids").(*schema.Set).List()),
		"target_spec_code": d.Get("target_spec_code"),
	}
	bodyParams := map[string]interface{}{
		"resize":      params,
		"is_auto_pay": true,
	}
	return bodyParams
}

func resourceDdsInstanceFlavorUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdsInstanceFlavorUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdsInstanceFlavorUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DDS instance flavor update resource is not supported. The resource is only removed from the " +
		"state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func handleDowngradeError(err error) bool {
	if err == nil {
		return false
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false
		}

		errorCode := utils.PathSearch("error_code", apiError, "").(string)
		if errorCode == "" {
			return false
		}

		// if error code is DBS.239037, then indicates it is downgrade scene, is_auto_pay must be null
		return errorCode == "DBS.239037"
	}
	return false
}
