package taurusdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/nodes/name
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/priority
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
// @API GaussDBforMySQL GET /v3/{project_id}/instances/details
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
func ResourceGaussDBMysqlNodeConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlNodeConfigCreate,
		UpdateContext: resourceGaussDBMysqlNodeConfigUpdate,
		ReadContext:   resourceGaussDBMysqlNodeConfigRead,
		DeleteContext: resourceGaussDBMysqlNodeConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussDBMysqlNodeConfigImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBMysqlNodeConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	if _, ok := d.GetOk("name"); ok {
		err = updateNodeName(ctx, d, client, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if _, ok := d.GetOk("priority"); ok {
		err = updateNodePriority(ctx, d, client, schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(d.Get("node_id").(string))

	return resourceGaussDBMysqlNodeConfigRead(ctx, d, meta)
}

func resourceGaussDBMysqlNodeConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	if d.HasChange("name") {
		err = updateNodeName(ctx, d, client, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("priority") {
		err = updateNodePriority(ctx, d, client, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBMysqlNodeConfigRead(ctx, d, meta)
}

func updateNodeName(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/name"
	)

	instanceId := d.Get("instance_id").(string)
	nodeId := d.Get("node_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLInstanceNodeNameBodyParams(d))

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL instance(%s) node(%s) name: %s", instanceId, nodeId, err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL instance(%s) node(%s) name, job_id is not found in the"+
			" response", instanceId, nodeId)
	}

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId.(string), d.Timeout(timeout))
	if err != nil {
		return fmt.Errorf("error waiting for updating instance(%s) node(%s) name job to complete: %s", instanceId,
			nodeId, err)
	}

	return nil
}

func buildUpdateGaussDBMySQLInstanceNodeNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"node_id": d.Get("node_id"),
		"name":    d.Get("name"),
	}
	bodyParams := map[string]interface{}{
		"node_list": []map[string]interface{}{params},
	}
	return bodyParams
}

func updateNodePriority(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/priority"
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
	createOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLInstanceNodePriorityBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL instance(%s) node(%s) prority: %s", instanceId, nodeId, err)
	}

	// wait 10 seconds so that the job can be completed
	// lintignore:R018
	time.Sleep(10 * time.Second)

	return nil
}

func buildUpdateGaussDBMySQLInstanceNodePriorityBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"priority": strconv.Itoa(d.Get("priority").(int)),
	}
	return bodyParams
}

func resourceGaussDBMysqlNodeConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/details?instance_ids={instance_id}"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB MySQL instance")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	node := utils.PathSearch(fmt.Sprintf("instances[0].nodes[?id=='%s']|[0]", d.Id()), getRespBody, nil)
	if node == nil {
		log.Printf("[WARN] failed to get GaussDB MySQL instance node by ID(%s)", d.Id())
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instances[0].id", getRespBody, nil)),
		d.Set("node_id", utils.PathSearch("id", node, nil)),
		d.Set("name", utils.PathSearch("name", node, nil)),
		d.Set("priority", utils.PathSearch("priority", node, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBMysqlNodeConfigDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting instance node config resource is not supported. The instance node config resource is only " +
		"removed from the state, the GaussDB MySQL instance node config remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceGaussDBMysqlNodeConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<node_id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("node_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
