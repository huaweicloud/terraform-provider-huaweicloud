// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DBSS
// ---------------------------------------------------------------

package dbss

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

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
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
				ForceNew:    true,
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
				ForceNew:    true,
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
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  `The charging period.`,
				ValidateFunc: validation.IntBetween(1, 9),
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
				ForceNew:    true,
				Description: `The description of the instance.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the IP address.`,
			},
			"tags": common.TagsForceNewSchema(),
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createInstance: Create a DBSS instance.
	var (
		createInstanceHttpUrl = "v2/{project_id}/dbss/audit/charge/period/order"
		createInstanceProduct = "dbss"
	)
	createInstanceClient, err := cfg.NewServiceClient(createInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	createInstancePath := createInstanceClient.Endpoint + createInstanceHttpUrl
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", createInstanceClient.ProjectID)

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	productId, err := getOrderProductId(d, cfg, region)
	if err != nil {
		return diag.Errorf("error getting product ID: %s", err)
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d, productId, cfg))
	createInstanceResp, err := createInstanceClient.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating Instance: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId, err := jmespath.Search("order_id", createInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating DBSS instance: ID is not found in API response")
	}

	// pay order
	resourceId, err := payOrder(ctx, d, cfg, orderId.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	err = createInstaceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the create of instance (%s) to complete: %s", d.Id(), err)
	}
	return resourceInstanceRead(ctx, d, meta)
}

func getOrderProductId(d *schema.ResourceData, cfg *config.Config, region string) (string, error) {
	var (
		getOrderProductIdHttpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		getOrderProductIdProduct = "bss"
	)
	getOrderProductIdClient, err := cfg.NewServiceClient(getOrderProductIdProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS Client: %s", err)
	}

	getOrderProductIdPath := getOrderProductIdClient.Endpoint + getOrderProductIdHttpUrl

	getOrderProductIdOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOrderProductIdOpt.JSONBody = utils.RemoveNil(buildGetFlavorsBodyParams(d,
		getOrderProductIdClient.ProjectID, region))
	getOrderProductIdResp, err := getOrderProductIdClient.Request("POST",
		getOrderProductIdPath, &getOrderProductIdOpt)

	if err != nil {
		return "", fmt.Errorf("error getting DBSS order product id: %s", err)
	}

	getCbhOrderProductIdRespBody, err := utils.FlattenResponse(getOrderProductIdResp)
	if err != nil {
		return "", err
	}
	curJson := utils.PathSearch("official_website_rating_result.product_rating_results",
		getCbhOrderProductIdRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return "", fmt.Errorf("fail to get DBSS order product id")
	}
	productId := utils.PathSearch("product_id", curArray[0], "")
	return productId.(string), nil
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
	region := cfg.GetRegion(d)
	var (
		payOrderHttpUrl = "v3/orders/customer-orders/pay"
		payOrderProduct = "bssv2"
	)
	bssClient, err := cfg.NewServiceClient(payOrderProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS Client: %s", err)
	}

	payOrderPath := bssClient.Endpoint + payOrderHttpUrl
	payOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	payOrderOpt.JSONBody = utils.RemoveNil(buildPayOrderBodyParams(orderId))
	_, err = bssClient.Request("POST", payOrderPath, &payOrderOpt)
	if err != nil {
		return "", fmt.Errorf("error pay DBSS order(%s): %s", orderId, err)
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
	return resourceId, err
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
		"enterprise_project_id": utils.ValueIgnoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
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
			"product_spec_desc":  utils.ValueIgnoreEmpty(d.Get("resource_spec_code")),
			"resource_spec_code": utils.ValueIgnoreEmpty(d.Get("resource_spec_code")),
			"resource_type":      "hws.resource.type.dbss",
		},
	}
}

func createInstaceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createInstaceWaitingHttpUrl = "v1/{project_id}/dbss/audit/jobs/{resource_id}"
				createInstaceWaitingProduct = "dbss"
			)
			createInstaceWaitingClient, err := cfg.NewServiceClient(createInstaceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Instance Client: %s", err)
			}

			createInstaceWaitingPath := createInstaceWaitingClient.Endpoint + createInstaceWaitingHttpUrl
			createInstaceWaitingPath = strings.ReplaceAll(createInstaceWaitingPath, "{project_id}", createInstaceWaitingClient.ProjectID)
			createInstaceWaitingPath = strings.ReplaceAll(createInstaceWaitingPath, "{resource_id}", fmt.Sprintf("%v", d.Id()))

			createInstaceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createInstaceWaitingResp, err := createInstaceWaitingClient.Request("GET", createInstaceWaitingPath, &createInstaceWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createInstaceWaitingRespBody, err := utils.FlattenResponse(createInstaceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			statusRaw, err := jmespath.Search(`jobs[0].status`, createInstaceWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `jobs[0].status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"SUCCESS",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createInstaceWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createInstaceWaitingRespBody, status, nil
			}

			return createInstaceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getInstance: Query the DBSS instance detail
	var (
		getInstanceHttpUrl = "v1/{project_id}/dbss/audit/instances"
		getInstanceProduct = "dbss"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving instance")
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instances, err := jmespath.Search("servers", getInstanceRespBody)
	if err != nil {
		diag.Errorf("error parsing servers from response= %#v", getInstanceRespBody)
	}

	instance, err := FilterInstances(instances.([]interface{}), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving instance")
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
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func FilterInstances(instances []interface{}, id string) (interface{}, error) {
	if len(instances) != 0 {
		for _, v := range instances {
			instance := v.(map[string]interface{})
			if instance["resource_id"] == id {
				return v, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Id()

	if d.HasChange("enterprise_project_id") {
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "auditInstance",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := common.MigrateEnterpriseProject(ctx, cfg, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
		return diag.Errorf("Error unsubscribing DBSS order = %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      waitForInstanceDelete(d, meta),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting DBSS instance: %s", err)
	}

	return nil
}

func waitForInstanceDelete(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	return func() (interface{}, string, error) {
		var (
			getInstanceHttpUrl = "v1/{project_id}/dbss/audit/instances"
			getInstanceProduct = "dbss"
		)
		getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
		if err != nil {
			return nil, "error", fmt.Errorf("error creating Instance Client: %s", err)
		}

		getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
		getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)

		getInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)

		if err != nil {
			return nil, "error", fmt.Errorf("error retrieving instance: %s", err)
		}

		getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
		if err != nil {
			return nil, "error", err
		}

		instances, err := jmespath.Search("servers", getInstanceRespBody)
		if err != nil {
			diag.Errorf("error parsing servers from response= %#v", getInstanceRespBody)
		}

		instance, err := FilterInstances(instances.([]interface{}), d.Id())
		if err != nil {
			// "instance" is useless, just return it to make the WaitForState break
			return "instance", "COMPLETE", nil
		}

		return instance, "PENDING", nil
	}
}
