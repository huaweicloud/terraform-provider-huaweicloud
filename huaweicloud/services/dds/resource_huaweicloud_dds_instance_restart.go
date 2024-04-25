package dds

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/restart
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSInstanceRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSInstanceRestartCreate,
		ReadContext:   resourceDDSInstanceRestartRead,
		DeleteContext: resourceDDSInstanceRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"target_id"},
			},
			"target_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"target_type"},
			},
		},
	}
}

func resourceDDSInstanceRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}
	instId := d.Get("instance_id").(string)

	restartOpts := instances.RestartOpts{
		TargetId: instId,
	}
	if v, ok := d.GetOk("target_id"); ok {
		restartOpts.TargetType = d.Get("target_type").(string)
		restartOpts.TargetId = v.(string)
	}

	// restart instance
	err = restartInstance(ctx, client, d.Timeout(schema.TimeoutCreate), instId, restartOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instId)

	return nil
}

func resourceDDSInstanceRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDSInstanceRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func restartInstance(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, instId string,
	restartOpts instances.RestartOpts) error {
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.RestartInstance(client, instId, restartOpts)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting instance: %s", err)
	}
	resp := r.(*instances.CommonResp)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, resp.JobId),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s ", resp.JobId, err)
	}

	return nil
}
