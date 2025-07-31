// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DBSS
// ---------------------------------------------------------------

package dbss

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS GET /v1/{project_id}/dbss/audit/instances
// @API DBSS GET /v1/{project_id}/dbss/audit/jobs/{resource_id}
// @API DBSS POST /v2/{project_id}/dbss/audit/charge/period/order
// @API BSS POST /v2/bills/ratings/period-resources/subscribe-rate
// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API DBSS POST /v1/{project_id}/{resource_type}/{resource_id}/tags/create
// @API DBSS DELETE /v1/{project_id}/{resource_type}/{resource_id}/tags/delete
// @API DBSS PUT /v1/{project_id}/dbss/audit/instances/{instance_id}
// @API DBSS POST /v1/{project_id}/dbss/audit/security-group
// @API DBSS POST /v1/{project_id}/dbss/audit/instance/start
// @API DBSS POST /v1/{project_id}/dbss/audit/instance/stop
// @API DBSS POST /v1/{project_id}/dbss/audit/instance/reboot
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceInstanceDelete,
		UpdateContext: resourceInstanceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance name.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The availability zone to which the instnce belongs.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the flavor.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The resource specifications.`,
			},
			// The splicing specification of this field lacks documentation. Please refer to the test case before using it.
			"product_spec_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `schema: Required; The product specification description in json string format.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID of the NIC.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The security group to which the instance belongs.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: `Enterprise project ID.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Billing mode.`,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
				}, false),
			},
			"period_unit": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period unit.`,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The charging period.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether auto renew is enabled. Valid values are "true" and "false".`,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the instance.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the IP address.`,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop", "reboot",
				}, false),
			},
			"tags": common.TagsSchema(),
			"connect_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection address.`,
			},
			"connect_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IPv6 address.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expired time`,
			},
			"port_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the port that the EIP is bound to.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance status.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the audit instance.`,
			},

			// Deprecated
			"product_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: `product_id is deprecated.`,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/dbss/audit/charge/period/order"
		product = "dbss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	productId, err := getOrderProductId(d, cfg, region)
	if err != nil {
		return diag.Errorf("error getting product ID: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateInstanceBodyParams(d, productId, cfg)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DBSS instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderID := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderID == "" {
		return diag.Errorf("error creating DBSS instance: order_id is not found in API response")
	}

	// pay order
	resourceId, err := payOrder(ctx, d, cfg, orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	err = waitingForInstanceStateCompleted(ctx, d, client, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the instance (%s) creation to complete: %s", d.Id(), err)
	}

	action := d.Get("action").(string)
	if action == "stop" || action == "reboot" {
		resp, err := QueryTargetDBSSInstance(client, d.Id())
		if err != nil {
			return diag.Errorf("error retrieving DBSS instance")
		}

		instanceId := utils.PathSearch("id", resp, "").(string)
		if instanceId == "" {
			return diag.Errorf("failed to operation the DBSS instance: 'instance_id' is not found in API response")
		}

		err = operateDbssInstance(ctx, d, client, action, instanceId)
		if err != nil {
			return diag.Errorf("error '%s' the DBSS instance in creation: %s", action, err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func getOrderProductId(d *schema.ResourceData, cfg *config.Config, region string) (string, error) {
	var (
		httpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		product = "bss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildGetFlavorsBodyParams(d, client.ProjectID, region)),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error getting the price of DBSS prepaid product: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return "", err
	}

	expression := "official_website_rating_result.product_rating_results|[0].product_id"
	productID := utils.PathSearch(expression, getRespBody, "").(string)
	if productID == "" {
		return "", fmt.Errorf("error getting the price of DBSS prepaid product: product_id is not found in API response")
	}

	return productID, nil
}

func buildGetFlavorsBodyParams(d *schema.ResourceData, projectId, region string) map[string]interface{} {
	periodUnit := d.Get("period_unit").(string)
	var periodType string
	if periodUnit == "month" {
		periodType = "2"
	} else {
		periodType = "3"
	}

	params := make(map[string]interface{})
	params["id"] = "1"
	params["cloud_service_type"] = "hws.service.type.dbss"
	params["resource_type"] = "hws.resource.type.dbss"
	params["resource_spec"] = d.Get("resource_spec_code")
	params["region"] = region
	params["period_type"] = periodType
	params["period_num"] = utils.ValueIgnoreEmpty(d.Get("period"))
	params["subscription_num"] = "1"

	bodyParams := map[string]interface{}{
		"project_id":    projectId,
		"product_infos": []map[string]interface{}{params},
	}
	return bodyParams
}

func payOrder(ctx context.Context, d *schema.ResourceData, cfg *config.Config, orderId string) (string, error) {
	var (
		region          = cfg.GetRegion(d)
		payOrderHttpUrl = "v3/orders/customer-orders/pay"
		payOrderProduct = "bssv2"
	)
	bssClient, err := cfg.NewServiceClient(payOrderProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS client: %s", err)
	}

	payOrderPath := bssClient.Endpoint + payOrderHttpUrl
	payOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildPayOrderBodyParams(orderId)),
	}

	_, err = bssClient.Request("POST", payOrderPath, &payOrderOpt)
	if err != nil {
		return "", fmt.Errorf("error paying DBSS order(%s): %s", orderId, err)
	}

	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", err
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", fmt.Errorf("error waiting for DBSS instance order %s complete: %s", orderId, err)
	}
	return resourceId, nil
}

func buildPayOrderBodyParams(orderId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"order_id":     orderId,
		"use_coupon":   "NO",
		"use_discount": "NO",
	}
	return bodyParams
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, productId string, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"region":                cfg.GetRegion(d),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"comment":               utils.ValueIgnoreEmpty(d.Get("description")),
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"cloud_service_type":    "hws.service.type.dbss",
		"flavor_ref":            utils.ValueIgnoreEmpty(d.Get("flavor")),
		"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"nics":                  buildCreateInstanceNicsRequestBody(d),
		"security_groups":       buildCreateInstanceSecurityGroupsRequestBody(d),
		"product_infos":         buildCreateInstanceProductInfosRequestBody(d, productId),
		"subscription_num":      1,
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"period_num":            utils.ValueIgnoreEmpty(d.Get("period")),
	}

	chargingMode := d.Get("charging_mode").(string)
	if chargingMode == "prePaid" {
		bodyParams["charging_mode"] = 0
	}

	periodUnit := d.Get("period_unit").(string)
	if periodUnit == "month" {
		bodyParams["period_type"] = 2
	} else {
		bodyParams["period_type"] = 3
	}

	autoRenew := d.Get("auto_renew").(string)
	if autoRenew == "true" {
		bodyParams["is_auto_renew"] = 1
	} else {
		bodyParams["is_auto_renew"] = 0
	}
	return bodyParams
}

func buildCreateInstanceNicsRequestBody(d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"ip_address": utils.ValueIgnoreEmpty(d.Get("ip_address")),
			"subnet_id":  utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		},
	}
}

func buildCreateInstanceSecurityGroupsRequestBody(d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id": utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		},
	}
}

func buildCreateInstanceProductInfosRequestBody(d *schema.ResourceData, productId string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"cloud_service_type": "hws.service.type.dbss",
			"product_id":         productId,
			"product_spec_desc":  utils.ValueIgnoreEmpty(d.Get("product_spec_desc")),
			"resource_spec_code": utils.ValueIgnoreEmpty(d.Get("resource_spec_code")),
			"resource_type":      "hws.resource.type.dbss",
		},
	}
}

func queryJobDetail(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	jobPath := client.Endpoint + "v1/{project_id}/dbss/audit/jobs/{resource_id}"
	jobPath = strings.ReplaceAll(jobPath, "{project_id}", client.ProjectID)
	jobPath = strings.ReplaceAll(jobPath, "{resource_id}", d.Id())
	jobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	jobResp, err := client.Request("GET", jobPath, &jobOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(jobResp)
}

func waitingForInstanceStateCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	t time.Duration) error {
	unexpectedStatus := []string{"ERROR"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobRespBody, err := queryJobDetail(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("jobs[0].status", jobRespBody, "").(string)
			if status == "" {
				return nil, "ERROR", fmt.Errorf("status is not found in DBSS job API response")
			}

			if status == "SUCCESS" {
				return jobRespBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return jobRespBody, status, nil
			}

			return jobRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func QueryTargetDBSSInstance(client *golangsdk.ServiceClient, resourceID string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/dbss/audit/instances"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("servers[?resource_id == '%s']|[0]", resourceID)
	instance := utils.PathSearch(expression, getRespBody, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return instance, nil
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "dbss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	instance, err := QueryTargetDBSSInstance(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DBSS instance")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("description", utils.PathSearch("comment", instance, nil)),
		d.Set("connect_ip", utils.PathSearch("connect_ip", instance, nil)),
		d.Set("connect_ipv6", utils.PathSearch("connect_ipv6", instance, nil)),
		d.Set("created_at", utils.PathSearch("created", instance, nil)),
		d.Set("expired_at", utils.PathSearch("expired", instance, nil)),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("port_id", utils.PathSearch("port_id", instance, nil)),
		d.Set("resource_spec_code", utils.PathSearch("resource_spec_code", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnetId", instance, nil)),
		d.Set("availability_zone", utils.PathSearch("zone", instance, nil)),
		d.Set("instance_id", utils.PathSearch("id", instance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dbss", region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err := updateNameAndDescription(client, d); err != nil {
			return diag.Errorf("error updating the DBSS instance: %s", err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "auditInstance",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		oldRaw, newRaw := d.GetChange("tags")
		oldMap := oldRaw.(map[string]interface{})
		newMap := newRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			if err = deleteTags(client, "auditInstance", d.Id(), oldMap); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(newMap) > 0 {
			if err := createTags(client, "auditInstance", d.Id(), newMap); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("security_group_id") {
		if err := updateSecurityGroupID(client, d); err != nil {
			return diag.Errorf("error updating DBSS security group: %s", err)
		}
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		instanceId := d.Get("instance_id").(string)
		if instanceId == "" {
			return diag.Errorf("editing action is currently not supported because of a failure in getting instance_id")
		}
		if action != "" {
			err = operateDbssInstance(ctx, d, client, action, instanceId)
			if err != nil {
				return diag.Errorf("error '%s' the DBSS instance in update: %s", action, err)
			}
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func updateNameAndDescription(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v1/{project_id}/dbss/audit/instances/{instance_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name":    utils.ValueIgnoreEmpty(d.Get("name")),
			"comment": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}

	resp, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", respBody, "").(string)
	if result != "success" {
		return fmt.Errorf("the update response value is not success")
	}

	return nil
}

func updateSecurityGroupID(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v1/{project_id}/dbss/audit/security-group"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resource_id":       d.Id(),
			"securitygroup_ids": []string{d.Get("security_group_id").(string)},
		},
	}
	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", updateRespBody, "").(string)
	if result != "success" {
		return fmt.Errorf("the security group's modification response value is not success")
	}

	return nil
}

func createTags(createTagsClient *golangsdk.ServiceClient, resourceType, resourceId string, tags map[string]interface{}) error {
	createTagsHttpUrl := "v1/{project_id}/{resource_type}/{resource_id}/tags/create"
	createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
	createTagsPath = strings.ReplaceAll(createTagsPath, "{project_id}", createTagsClient.ProjectID)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", resourceType)
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", resourceId)
	createTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	createTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
	if err != nil {
		return fmt.Errorf("error creating tags: %s", err)
	}
	return nil
}

func deleteTags(deleteTagsClient *golangsdk.ServiceClient, resourceType, resourceId string, tags map[string]interface{}) error {
	deleteTagsHttpUrl := "v1/{project_id}/{resource_type}/{resource_id}/tags/delete"
	deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{project_id}", deleteTagsClient.ProjectID)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", resourceType)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", resourceId)
	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}

	_, err := deleteTagsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags: %s", err)
	}
	return nil
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dbss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		return diag.Errorf("error unsubscribing DBSS instance (%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETE"},
		Refresh: func() (interface{}, string, error) {
			instance, err := QueryTargetDBSSInstance(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "instance", "COMPLETE", nil
				}
				return nil, "ERROR", err
			}

			return instance, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DBSS instance (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func operateDbssInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	action, instanceId string) error {
	httpUrl := "v1/{project_id}/dbss/audit/instance/{action}"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{action}", action)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"instance_id": instanceId,
		},
	}

	_, err := client.Request("POST", actionPath, &actionOpts)
	if err != nil {
		return err
	}

	err = waitingForOperationInstanceCompleted(ctx, d, client, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the instance '%s' to complete: %s", action, err)
	}

	return nil
}

func waitingForOperationInstanceCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := QueryTargetDBSSInstance(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			task := utils.PathSearch("task", respBody, "").(string)
			if task == "" {
				return nil, "ERROR", fmt.Errorf("the 'task' is not found in API response")
			}

			if task == "NO_TASK" {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}
