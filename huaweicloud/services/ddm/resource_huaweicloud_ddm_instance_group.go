package ddm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceGroupNonUpdatableParams = []string{
	"instance_id", "name", "type", "flavor_id", "nodes", "nodes.*.available_zone", "nodes.*.subnet_id",
}

// @API DDM POST /v3/{project_id}/instances/{instance_id}/groups
// @API DDM GET /v3/{project_id}/instances/{instance_id}/groups
// @API DDM GET /v3/{project_id}/jobs/{job_id}
// @API DDM GET /v1/{project_id}/instances/{instance_id}
// @API DDM DELETE /v3/{project_id}/instances/{instance_id}/groups/{group_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS POST /v2/orders/suscriptions/resources/query
func ResourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceGroupCreate,
		UpdateContext: resourceInstanceGroupUpdate,
		ReadContext:   resourceInstanceGroupRead,
		DeleteContext: resourceInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceGroupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceGroupNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Elem:     instanceGroupNodesSchema(),
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_load_balance": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_num_per_node": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_num_per_node": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func instanceGroupNodesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/groups"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateInstanceGroupBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if jobId == "" && orderId == "" {
		return diag.Errorf("error creating DDM instance group: job_id and order_id are not found in API response")
	}

	if jobId != "" {
		groupId := utils.PathSearch("group_id", createRespBody, "").(string)
		if groupId == "" {
			return diag.Errorf("error creating DDM instance group: group_id is not found in API response")
		}

		d.SetId(groupId)
		err = checkJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err := cbc.PaySubscriptionOrder(bssClient, orderId); err != nil {
			return diag.Errorf("error paying for DDM instance group order (%s): %s", orderId, err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		nodeId, err := common.WaitOrderAllResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for DDM instance group order(%s) complete: %s", orderId, err)
		}
		group, err := getInstanceGroup(client, instanceId, "", strings.TrimSuffix(nodeId, ".vm"))
		if err != nil {
			return diag.FromErr(err)
		}
		groupId := utils.PathSearch("id", group, nil)
		if groupId == nil {
			return diag.Errorf("error creating DDM instance group: group ID is not found in API response")
		}
		d.SetId(groupId.(string))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(instanceId, client),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to running: %s", instanceId, err)
	}

	return resourceInstanceGroupRead(ctx, d, meta)
}

func buildCreateInstanceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"type":      d.Get("type"),
		"flavor_id": d.Get("flavor_id"),
		"nodes":     buildCreateInstanceGroupNodesBodyParams(d),
	}
	return bodyParams
}

func buildCreateInstanceGroupNodesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	nodes := d.Get("nodes").([]interface{})
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(nodes))
	for _, v := range nodes {
		raw := v.(map[string]interface{})
		node := map[string]interface{}{
			"available_zone": raw["available_zone"],
			"subnet_id":      raw["subnet_id"],
		}

		rst = append(rst, node)
	}
	return rst
}

func resourceInstanceGroupUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	group, err := getInstanceGroup(client, d.Get("instance_id").(string), d.Id(), "")
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DDM instance group")
	}
	if group == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", group, nil)),
		d.Set("type", utils.PathSearch("role", group, nil)),
		d.Set("endpoint", utils.PathSearch("endpoint", group, nil)),
		d.Set("ipv6_endpoint", utils.PathSearch("ipv6_endpoint", group, nil)),
		d.Set("is_load_balance", utils.PathSearch("is_load_balance", group, nil)),
		d.Set("is_default_group", utils.PathSearch("is_default_group", group, nil)),
		d.Set("cpu_num_per_node", utils.PathSearch("cpu_num_per_node", group, nil)),
		d.Set("mem_num_per_node", utils.PathSearch("mem_num_per_node", group, nil)),
		d.Set("architecture", utils.PathSearch("architecture", group, nil)),
		d.Set("nodes", flattenGetInstanceGroupBodyNodes(group)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstanceGroup(client *golangsdk.ServiceClient, instanceId, groupId, nodeId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/groups"
	)

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	var group interface{}
	var searchPath string
	if groupId != "" {
		searchPath = fmt.Sprintf("group_list|[?id=='%s']|[0]", groupId)
	} else {
		searchPath = fmt.Sprintf("group_list[?node_list[?id=='%s']]|[0]", nodeId)
	}
	for {
		getPath := getBasePath + buildPageQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		group = utils.PathSearch(searchPath, getRespBody, nil)
		if group != nil {
			break
		}
		groups := utils.PathSearch("group_list", getRespBody, make([]interface{}, 0)).([]interface{})
		offset += len(groups)
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if offset >= int(totalCount) {
			break
		}
	}
	return group, nil
}

func buildPageQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%d", offset)
}

func flattenGetInstanceGroupBodyNodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("node_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"available_zone": utils.PathSearch("az", v, nil),
		})
	}
	return rst
}

func resourceInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/groups/{group_id}"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddmInstanceStatusRefreshFunc(instanceId, client),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DDM instance group")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	orderId := utils.PathSearch("order_id", deleteRespBody, "").(string)
	if jobId == "" && orderId == "" {
		return diag.Errorf("error deleting DDM instance group: job_id and order_id are not found in API response")
	}

	if jobId != "" {
		err = checkJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(instanceId, client),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to running: %s", instanceId, err)
	}

	return nil
}

func resourceInstanceGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import id, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func checkJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      jobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DDM job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func jobRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v3/{project_id}/jobs/{job_id}"
		)

		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
		getResp, err := client.Request("GET", getPath, &getJobStatusOpt)
		if err != nil {
			return nil, "Failed", err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "Failed", err
		}

		status := utils.PathSearch("job.status", getRespBody, "").(string)
		if status == "" {
			return nil, "Failed", errors.New("job is not found")
		}
		if status == "Failed" {
			return getRespBody, "Failed", nil
		}
		if status == "Completed" {
			return getRespBody, "Completed", nil
		}

		return getRespBody, "Pending", nil
	}
}
