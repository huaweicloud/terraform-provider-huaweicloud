package dds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceNodeDeleteNonUpdatableParams = []string{
	"instance_id", "num", "node_list",
}

// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/nodes
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
func ResourceInstanceNodeDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceNodeDeleteCreate,
		ReadContext:   resourceInstanceNodeDeleteRead,
		UpdateContext: resourceInstanceNodeDeleteUpdate,
		DeleteContext: resourceInstanceNodeDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceNodeDeleteNonUpdatableParams),

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
			"num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildInstanceNodeDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"num":       utils.ValueIgnoreEmpty(d.Get("num")),
		"node_list": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("node_list").([]interface{}))),
	}
}

func resourceInstanceNodeDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deleteHttpUrl = "v3/{project_id}/instances/{instance_id}/nodes"
		instanceId    = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildInstanceNodeDeleteBodyParams(d)),
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}

	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error deleting instance nodes: %s", err)
	}

	d.SetId(instanceId)

	respBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	orderId := utils.PathSearch("order_id", respBody, "").(string)

	if jobId == "" && orderId == "" {
		return diag.Errorf("error deleting DDS instance node: job_id or order_id is not found in API response")
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

func resourceInstanceNodeDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceNodeDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceNodeDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting instance node delete resource is not supported. The delete resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
