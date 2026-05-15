package dds

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var readonlyNodeNonUpdatableParams = []string{
	"instance_id",
	"delay",
}

// @API DDS POST /v3/{project_id}/instances/{instance_id}/readonly-node
// @API DDS GET /v3/{project_id}/instances
// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/readonly-node
// @API DDS POST /v3/{project_id}/instances/{instance_id}/resize
// @API DDS POST /v3/{project_id}/instances/{instance_id}/enlarge-volume
// @API DDS GET /v3/{project_id}/instances/{instance_id}/disk-usage
// @API DDS POST /v3/{project_id}/instances/{instance_id}/modify-internal-ip
// @API DDS GET /v3/{project_id}/jobs
// @API BSS GET /v2/orders/customer-orders/details
// @API BSS POST /v2/orders/suscriptions/resources/query
func ResourceReadonlyNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReadonlyNodeCreate,
		ReadContext:   resourceReadonlyNodeRead,
		UpdateContext: resourceReadonlyNodeUpdate,
		DeleteContext: resourceReadonlyNodeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceReadonlyNodeImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(readonlyNodeNonUpdatableParams),

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
			"spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateReadonlyNodeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"spec_code":   d.Get("spec_code"),
		"num":         1,
		"delay":       utils.ValueIgnoreEmpty(d.Get("delay")),
		"is_auto_pay": true,
	}

	return bodyParams
}

func resourceReadonlyNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		createHttpUrl = "v3/{project_id}/instances/{instance_id}/readonly-node"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateReadonlyNodeBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding readonly node to DDS instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", respBody, "").(string)
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if orderId == "" && jobId == "" {
		return diag.Errorf("error adding readonly node: unable to find `order_id` or `job_id` from the API response")
	}

	if orderId != "" {
		nodeId, err := waitForPrePaidOrderCompleted(ctx, cfg, d, d.Timeout(schema.TimeoutCreate), orderId)
		if err != nil {
			return diag.FromErr(err)
		}

		if nodeId == "" {
			return diag.Errorf("error adding readonly node to DDS instance: unable to find node ID")
		}

		d.SetId(nodeId)
	}

	if jobId != "" {
		nodeId, err := waitForJobCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), instanceId, jobId)
		if err != nil {
			return diag.FromErr(err)
		}

		if nodeId == "" {
			return diag.Errorf("error adding readonly node to DDS instance: unable to find node ID")
		}

		d.SetId(nodeId)
	}

	// The readonly node volume enlarge
	if err = modifyReadonlyNodeVolumeSize(ctx, cfg, client, d, instanceId, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	// The readonly node update private IP
	if err = modifyReadonlyNodePrivateIp(ctx, cfg, client, d, instanceId, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return resourceReadonlyNodeRead(ctx, d, meta)
}

func modifyReadonlyNodeVolumeSize(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceId, nodeId string) error {
	targetSize, ok := d.GetOk("size")
	if !ok {
		return nil
	}

	volume, err := getReadonlyNodeVolumeInfo(client, instanceId, nodeId)
	if err != nil {
		return fmt.Errorf("error retrieving readonly node volume information: %s", err)
	}

	volumeSize := utils.PathSearch("size", volume, "").(string)
	if convertStringToInt(targetSize.(string)).(int) < convertStringToInt(volumeSize).(int) {
		return errors.New("err updating the readonly node volume size: the volume does not support shrink")
	}

	if convertStringToInt(targetSize.(string)).(int) > convertStringToInt(volumeSize).(int) {
		err = updateReadonlyNodeVolumeSize(ctx, cfg, client, d, d.Timeout(schema.TimeoutCreate), instanceId)
		if err != nil {
			return err
		}
	}

	return nil
}

func modifyReadonlyNodePrivateIp(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceId, nodeId string) error {
	targePrivateIp, ok := d.GetOk("private_ip")
	if !ok {
		return nil
	}

	nodeInfo, err := GetReadonlyNodeInfo(client, instanceId, nodeId)
	if err != nil {
		return fmt.Errorf("error query the readonly node information: %s", err)
	}

	privateIp := utils.PathSearch("private_ip", nodeInfo, "").(string)
	if targePrivateIp.(string) == privateIp {
		return nil
	}

	err = updateReadonlyNodePrivateIp(ctx, cfg, client, d, d.Timeout(schema.TimeoutCreate), instanceId)
	return err
}

func waitForPrePaidOrderCompleted(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	timeout time.Duration, orderId string) (string, error) {
	var nodeId string
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating BSS v2 client: %s", err)
	}

	err = common.WaitOrderComplete(ctx, bssClient, orderId, timeout)
	if err != nil {
		return "", err
	}

	_, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderId, timeout)
	if err != nil {
		return "", err
	}

	resourceIds, err := cbc.GetResourceIDsByOrder(bssClient, orderId, 0)
	if err != nil {
		return "", fmt.Errorf("error retrieving resource IDs of DDS instance readonly node order (%s): %s", orderId, err)
	}

	for _, resourceId := range resourceIds {
		if strings.Contains(resourceId, ".vm") {
			readonlyId := strings.Split(resourceId, ".")
			nodeId = readonlyId[0]
			break
		}
	}

	return nodeId, nil
}

func waitForJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, jobId string) (string, error) {
	instanceBody, err := listInstanceGroups(client, instanceId)
	if err != nil {
		return "", err
	}

	resourceId := utils.PathSearch("groups[?type=='readonly'].nodes[]|[?status=='creating'].id|[0]", instanceBody, "").(string)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, jobId),
		Timeout:      timeout,
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for the job (%s) completed: %s ", jobId, err)
	}

	return resourceId, nil
}

func convertStringToInt(rawValue string) interface{} {
	if rawValue == "" {
		return nil
	}

	r, err := strconv.Atoi(rawValue)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", rawValue)
		return nil
	}

	return r
}

func listInstanceGroups(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances?id={instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("instances|[0]", respBody, nil), nil
}

func getReadonlyNodeVolumeInfo(client *golangsdk.ServiceClient, instanceId, nodeId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/disk-usage"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	volumeInfo := utils.PathSearch(fmt.Sprintf("volumes[?entity_id=='%s']|[0]", nodeId), respBody, nil)
	if volumeInfo == nil {
		return nil, errors.New("error retrieving volume from the API response")
	}

	return volumeInfo, nil
}

func resourceReadonlyNodeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	nodeInfo, err := GetReadonlyNodeInfo(client, instanceId, d.Id())
	if err != nil {
		// When the DDS instance does not exist, the response HTTP status code of the query API is 400
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.280110"),
			"error retrieving readonly node information")
	}

	volumeInfo, err := getReadonlyNodeVolumeInfo(client, instanceId, d.Id())
	if err != nil {
		log.Printf("[Warn] error retrieving the readonly node volume information")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("spec_code", utils.PathSearch("spec_code", nodeInfo, nil)),
		d.Set("name", utils.PathSearch("name", nodeInfo, nil)),
		d.Set("status", utils.PathSearch("status", nodeInfo, nil)),
		d.Set("role", utils.PathSearch("role", nodeInfo, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", nodeInfo, nil)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", nodeInfo, nil)),
		d.Set("size", utils.PathSearch("size", volumeInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetReadonlyNodeInfo(client *golangsdk.ServiceClient, instacneId, nodeId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances?id={instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instacneId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	nodeInfo := utils.PathSearch(fmt.Sprintf("instances[].groups[].nodes[]|[?id=='%s']|[0]", nodeId), respBody, nil)
	if nodeInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return nodeInfo, nil
}

func resourceReadonlyNodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	if d.HasChanges("spec_code") {
		err = updateReadonlyNodeSpecification(ctx, cfg, client, d, d.Timeout(schema.TimeoutUpdate), instanceId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("size") {
		err = updateReadonlyNodeVolumeSize(ctx, cfg, client, d, d.Timeout(schema.TimeoutUpdate), instanceId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("private_ip") {
		err = updateReadonlyNodePrivateIp(ctx, cfg, client, d, d.Timeout(schema.TimeoutUpdate), instanceId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceReadonlyNodeRead(ctx, d, meta)
}

func buildUpdateaReadonlyNodeSpecificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"resize": map[string]interface{}{
			"target_type":      "readonly",
			"target_id":        d.Id(),
			"target_spec_code": d.Get("spec_code"),
		},
		"is_auto_pay": true,
	}

	return params
}

func updateReadonlyNodeSpecification(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData, timeout time.Duration, instanceId string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/resize"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateaReadonlyNodeSpecificationBodyParams(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		requestResp, err := client.Request("POST", requestPath, &requestOpt)
		retry, err := handleMultiOperationsError(err)
		return requestResp, retry, err
	}

	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating the readonly node specification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	err = waitForReadonlyNodeOperationCompleted(ctx, cfg, client, d, timeout, respBody)
	if err != nil {
		return fmt.Errorf("error waiting for the readonly node specification update to complete: %s", err)
	}

	return nil
}

func buildUpdateaReadonlyNodeVolumeSize(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"volume": map[string]interface{}{
			"size":     d.Get("size"),
			"node_ids": []string{d.Id()},
		},
		"is_auto_pay": true,
	}

	return params
}

func updateReadonlyNodeVolumeSize(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData, timeout time.Duration, instanceId string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/enlarge-volume"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateaReadonlyNodeVolumeSize(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		requestResp, err := client.Request("POST", requestPath, &requestOpt)
		retry, err := handleMultiOperationsError(err)
		return requestResp, retry, err
	}

	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating the readonly node volume size: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	err = waitForReadonlyNodeOperationCompleted(ctx, cfg, client, d, timeout, respBody)
	if err != nil {
		return fmt.Errorf("error waiting for the readonly node volume size update to complete: %s", err)
	}

	return nil
}

func buildUpdateaReadonlyNodePrivateIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"node_id": d.Id(),
		"new_ip":  d.Get("private_ip"),
	}

	return params
}

func updateReadonlyNodePrivateIp(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient,
	d *schema.ResourceData, timeout time.Duration, instanceId string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/modify-internal-ip"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildUpdateaReadonlyNodePrivateIpBodyParams(d),
	}

	retryFunc := func() (interface{}, bool, error) {
		requestResp, err := client.Request("POST", requestPath, &requestOpt)
		retry, err := handleMultiOperationsError(err)
		return requestResp, retry, err
	}

	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      timeout,
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating the readonly node private IP: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return err
	}

	err = waitForReadonlyNodeOperationCompleted(ctx, cfg, client, d, timeout, respBody)
	if err != nil {
		return fmt.Errorf("error waiting for the readonly node private IP update to complete: %s", err)
	}

	return nil
}

func resourceReadonlyNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/readonly-node"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"node_list": []string{d.Id()},
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		requestResp, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return requestResp, retry, err
	}

	resp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.280407", "DBS.280110"),
			"error deleting DDS instance readonly node")
	}

	respBody, err := utils.FlattenResponse(resp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForReadonlyNodeOperationCompleted(ctx, cfg, client, d, d.Timeout(schema.TimeoutDelete), respBody)
	if err != nil {
		return diag.Errorf("error waiting for the readonly node delete to complete: %s", err)
	}

	return nil
}

func waitForReadonlyNodeOperationCompleted(ctx context.Context, cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, respBody interface{}) error {
	region := cfg.GetRegion(d)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	orderId := utils.PathSearch("order_id", respBody, "").(string)
	if orderId == "" && jobId == "" {
		return errors.New("error operating readonly node: unable to find `order_id` or `job_id` from the API response")
	}

	if orderId != "" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, orderId, timeout)
		if err != nil {
			return err
		}
	}

	if jobId != "" {
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, jobId),
			Timeout:      timeout,
			Delay:        60 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for the job (%s) completed: %s ", jobId, err)
		}
	}

	return nil
}

func resourceReadonlyNodeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
