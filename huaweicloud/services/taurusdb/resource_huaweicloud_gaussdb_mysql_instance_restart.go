package taurusdb

import (
	"context"
	"fmt"
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

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/restart
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/restart
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
func ResourceGaussDBMysqlRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlRestartCreate,
		ReadContext:   resourceGaussDBMysqlRestartRead,
		DeleteContext: resourceGaussDBMysqlRestartDelete,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delay": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGaussDBMysqlRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	nodeId := d.Get("node_id").(string)
	if len(nodeId) > 0 {
		err = restartNode(ctx, d, client)
	} else {
		err = restartInstance(ctx, d, client)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func restartInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/restart"
	)

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMySQLRestartBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting GaussDB MySQL instance(%s): %s", instanceId, err)
	}
	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error restarting GaussDB MySQL instance(%s), job_id is not found in the response", instanceId)
	}

	d.SetId(instanceId)

	// if delay is true, then the task is scheduled task, it is not needed to wait
	if !d.Get("delay").(bool) {
		err = checkGaussDBMySQLJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error waiting for restarting instance(%s) job to complete: %s", instanceId, err)
		}
	}

	return nil
}

func restartNode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/restart"
	)

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{node_id}", nodeId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMySQLRestartBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error restarting GaussDB MySQL instance(%s) node(%s): %s", instanceId, nodeId, err)
	}
	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error restarting GaussDB MySQL instance(%s) node(%s), job_id is not found in the "+
			"response", instanceId, nodeId)
	}

	d.SetId(nodeId)

	// if delay is true, then the task is scheduled task, it is not needed to wait
	if !d.Get("delay").(bool) {
		err = checkGaussDBMySQLJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error waiting for restarting instance(%s) node(%s) job to complete: %s", instanceId,
				nodeId, err)
		}
	}

	return nil
}

func buildCreateGaussDBMySQLRestartBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"delay": d.Get("delay"),
	}
	return bodyParams
}

func resourceGaussDBMysqlRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBMysqlRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the GaussDB MySQL instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
