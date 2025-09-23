package deh

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dehInstanceNonUpdatableParams = []string{
	"availability_zone",
	"host_type",
	"metadata",
	"charging_mode",
	"period_unit",
	"period",
	"enterprise_project_id",
}

// @API DEH POST /v1.0/{project_id}/dedicated-hosts
// @API DEH GET /v1/{project_id}/jobs/{job_id}
// @API DEH GET /v1.0/{project_id}/dedicated-hosts/{dedicated_host_id}
// @API DEH PUT /v1.0/{project_id}/dedicated-hosts/{dedicated_host_id}
// @API DEH POST /v1.0/{project_id}/dedicated-host-tags/{dedicated_host_id}/tags/action
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDehInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDehInstanceCreate,
		ReadContext:   resourceDehInstanceRead,
		UpdateContext: resourceDehInstanceUpdate,
		DeleteContext: resourceDehInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(dehInstanceNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_placement": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"host_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dehInstanceHostPropertiesSchema(),
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allocated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_uuids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sys_tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dehInstanceHostPropertiesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_instance_capacities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dehInstanceHostPropertiesAvailableInstanceCapacitiesSchema(),
			},
		},
	}
	return &sc
}

func dehInstanceHostPropertiesAvailableInstanceCapacitiesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDehInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl    = "v1.0/{project_id}/dedicated-hosts"
		product    = "deh"
		ecsProduct = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEH client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDehInstanceBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DEH instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", createRespBody, nil)
	if orderId != nil {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	ecsClient, err := cfg.NewServiceClient(ecsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	jobId := utils.PathSearch("jobId || job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating DEH instance: job_id is not found in API response")
	}
	res, err := checkDehJobFinish(ctx, ecsClient, jobId.(string), 0, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId := utils.PathSearch("entities.entities.entities[0].resourceId", res, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}
	d.SetId(resourceId)

	return resourceDehInstanceRead(ctx, d, meta)
}

func buildCreateDehInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"quantity":          1,
		"availability_zone": d.Get("availability_zone"),
		"name":              d.Get("name"),
		"host_type":         d.Get("host_type"),
		"auto_placement":    utils.ValueIgnoreEmpty(d.Get("auto_placement")),
		"metadata":          utils.ValueIgnoreEmpty(d.Get("metadata")),
		"tags":              utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"extend_param":      buildCreateDehInstanceExtendParamsBodyParams(d, cfg),
	}
	return bodyParams
}

func buildCreateDehInstanceExtendParamsBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := make(map[string]interface{})
	if v, ok := d.GetOk("charging_mode"); ok {
		bodyParams["charging_mode"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		bodyParams["period_type"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		bodyParams["period_num"] = v
	}
	if v, ok := d.GetOk("auto_renew"); ok {
		bodyParams["is_auto_renew"] = v
	}
	if len(bodyParams) > 0 {
		bodyParams["is_auto_pay"] = "true"
	}
	if v := cfg.GetEnterpriseProjectID(d); v != "" {
		bodyParams["enterprise_project_id"] = v
	}
	if len(bodyParams) > 0 {
		return bodyParams
	}
	return nil
}

func checkDehJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string, delay int,
	timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      dehJobStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        time.Duration(delay) * time.Second,
		PollInterval: 10 * time.Second,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for DEH job (%s) to be completed: %s ", jobID, err)
	}
	return res, nil
}

func dehJobStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v1/{project_id}/jobs/{job_id}"
		)

		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "FAIL" {
			return nil, "FAIL", errors.New("the DEH instance create job execution has failed")
		}
		if status == "SUCCESS" {
			return getRespBody, "SUCCESS", nil
		}
		return getRespBody, "PENDING", nil
	}
}

func resourceDehInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "deh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEH client: %s", err)
	}

	getRespBody, err := getDehInstance(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DEH instance")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("availability_zone", utils.PathSearch("dedicated_host.availability_zone", getRespBody, nil)),
		d.Set("name", utils.PathSearch("dedicated_host.name", getRespBody, nil)),
		d.Set("host_type", utils.PathSearch("dedicated_host.host_properties.host_type", getRespBody, nil)),
		d.Set("auto_placement", utils.PathSearch("dedicated_host.auto_placement", getRespBody, nil)),
		d.Set("metadata", flattenMetadata(d, getRespBody)),
		d.Set("state", utils.PathSearch("dedicated_host.state", getRespBody, nil)),
		d.Set("available_vcpus", utils.PathSearch("dedicated_host.available_vcpus", getRespBody, nil)),
		d.Set("available_memory", utils.PathSearch("dedicated_host.available_memory", getRespBody, nil)),
		d.Set("allocated_at", utils.PathSearch("dedicated_host.allocated_at", getRespBody, nil)),
		d.Set("instance_total", utils.PathSearch("dedicated_host.instance_total", getRespBody, nil)),
		d.Set("instance_uuids", utils.PathSearch("dedicated_host.instance_uuids", getRespBody, nil)),
		d.Set("tags", utils.PathSearch("dedicated_host.tags", getRespBody, nil)),
		d.Set("sys_tags", utils.PathSearch("dedicated_host.sys_tags", getRespBody, nil)),
		d.Set("host_properties", flattenDehHostProperties(getRespBody)),
	)
	mErr = multierror.Append(mErr, d.Set("metadata", flattenMetadata(d, getRespBody)))

	chargingMode := utils.PathSearch("dedicated_host.metadata.charging_mode", getRespBody, nil)
	if chargingMode != nil && chargingMode == "1" {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMetadata(d *schema.ResourceData, resp interface{}) map[string]interface{} {
	rawMetadata := d.Get("metadata").(map[string]interface{})
	metadata := utils.PathSearch("dedicated_host.metadata", resp, nil)
	rst := make(map[string]interface{})
	if metadata != nil {
		for k, v := range metadata.(map[string]interface{}) {
			if _, ok := rawMetadata[k]; ok {
				rst[k] = v
			}
		}
	}
	return rst
}

func flattenDehHostProperties(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dedicated_host.host_properties", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"host_type":                     utils.PathSearch("host_type", curJson, nil),
			"host_type_name":                utils.PathSearch("host_type_name", curJson, nil),
			"vcpus":                         utils.PathSearch("vcpus", curJson, nil),
			"cores":                         utils.PathSearch("cores", curJson, nil),
			"sockets":                       utils.PathSearch("sockets", curJson, nil),
			"memory":                        utils.PathSearch("memory", curJson, nil),
			"available_instance_capacities": flattenDehHostPropertyAvailableInstanceCapacities(curJson),
		},
	}

	return rst
}

func flattenDehHostPropertyAvailableInstanceCapacities(resp interface{}) []interface{} {
	curJson := utils.PathSearch("available_instance_capacities", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor": utils.PathSearch("flavor", v, nil),
		})
	}

	return rst
}

func resourceDehInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "deh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEH client: %s", err)
	}

	if d.HasChanges("name", "auto_placement") {
		err = updateDehInstance(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = updateDehInstanceTags(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the DEH instance (%s): %s", d.Id(), err)
		}
	}

	return resourceDehInstanceRead(ctx, d, meta)
}

func updateDehInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl = "v1.0/{project_id}/dedicated-hosts/{dedicated_host_id}"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{dedicated_host_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateDehInstanceBodyParams(d))

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating DEH instance: %s", err)
	}
	return nil
}

func buildUpdateDehInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"auto_placement": d.Get("auto_placement"),
		"name":           d.Get("name"),
	}
	bodyParams := map[string]interface{}{
		"dedicated_host": params,
	}
	return bodyParams
}

func updateDehInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	if len(oMap) > 0 {
		err := doUpdateDehInstanceTags(d, client, oMap, "delete")
		if err != nil {
			return fmt.Errorf("error deleting tags from DEH instance (%s): %s", d.Id(), err)
		}
	}

	if len(nMap) > 0 {
		err := doUpdateDehInstanceTags(d, client, nMap, "create")
		if err != nil {
			return fmt.Errorf("error adding tags to DEH instance (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func doUpdateDehInstanceTags(d *schema.ResourceData, client *golangsdk.ServiceClient, updateTags map[string]interface{},
	action string) error {
	var (
		httpUrl = "v1.0/{project_id}/dedicated-host-tags/{dedicated_host_id}/tags/action"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{dedicated_host_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateDehInstanceTagsBodyParams(updateTags, action))

	_, err := client.Request("POST", updatePath, &updateOpt)
	return err
}

func buildUpdateDehInstanceTagsBodyParams(updateTags map[string]interface{}, action string) map[string]interface{} {
	tags := make([]interface{}, 0, len(updateTags))
	for key, value := range updateTags {
		tags = append(tags, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
	bodyParams := map[string]interface{}{
		"action": action,
		"tags":   tags,
	}
	return bodyParams
}

func resourceDehInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.Get("charging_mode").(string) != "prePaid" {
		errorMsg := "Deleting postPaid DEH instance is not supported. The restoration record is only removed from the " +
			"state, but it remains in the cloud."
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorMsg,
			},
		}
	}

	var (
		product = "deh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEH client: %s", err)
	}

	err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
	if err != nil {
		return diag.Errorf("error deleting DEH instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"deleted"},
		Refresh:    dehInstanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      20 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DEH instance (%s) to be deleted: %s ", d.Id(), err)
	}

	return nil
}

func dehInstanceStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getDehInstance(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "deleted", nil
			}
			return nil, "error", err
		}
		return res, "pending", nil
	}
}

func getDehInstance(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	var (
		httpUrl = "v1.0/{project_id}/dedicated-hosts/{dedicated_host_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dedicated_host_id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}
