package cbh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH POST /v1/{project_id}/cbs/{server_id}/network/change
// @API CBH POST /v1/{project_id}/cbs/instance/{server_id}/eip/bind
// @API CBH POST /v1/{project_id}/cbs/instance/{server_id}/eip/unbind
// @API CBH POST /v1/{project_id}/cbs/instance/create
// @API CBH GET /v1/{project_id}/cbs/instance/list
// @API CBH PUT /v1/{project_id}/cbs/instance/password
// @API CBH POST /v1/{project_id}/cbs/period/order
// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceCBHInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBHInstanceCreate,
		UpdateContext: resourceCBHInstanceUpdate,
		ReadContext:   resourceCBHInstanceRead,
		DeleteContext: resourceCBHInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the name of the CBH instance.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product ID of the CBH server.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the security group.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the availability zone name.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password for logging in to the management console.`,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
				}, false),
				Description: `Specifies the charging mode of the read replica instance.`,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
				Description: `Specifies the charging period unit of the instance.`,
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 9),
				Description:  `Specifies the charging period of the read replica instance.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether auto renew is enabled.`,
			},
			"subnet_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP address of the subnet.`,
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the elastic IP.`,
			},
			"public_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"public_ip_id"},
				Description:  `Specifies the elastic IP address.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the IPv6 network is enabled.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private ip of the instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version of the instance image.`,
			},
		},
	}
}

func resourceCBHInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createCbhInstanceProduct = "cbh"
	)

	createCbhInstanceClient, err := cfg.NewServiceClient(createCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH Client: %s", err)
	}

	// createInstance: create CBH instance
	instanceKey, err := createInstance(d, cfg, createCbhInstanceClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if instanceKey == nil {
		return diag.Errorf("error creating CbhInstance: instance_key is empty")
	}

	// createOrder: create instance order
	orderId, err := createOrder(d, cfg, createCbhInstanceClient, region, instanceKey.(string))
	if err != nil {
		return diag.FromErr(err)
	}
	// pay order
	resourceId, err := payOrder(ctx, d, cfg, orderId)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId, err := getInstanceIdByResourceId(createCbhInstanceClient, resourceId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    cbhInstanceStateRefreshFunc(createCbhInstanceClient, instanceId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("timeout waiting for instance to active: %s", err)
	}

	return resourceCBHInstanceRead(ctx, d, meta)
}

func createInstance(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient) (interface{}, error) {
	// createInstance: create CBH instance
	var (
		createInstanceHttpUrl = "v1/{project_id}/cbs/instance/create"
	)

	createInstancePath := client.Endpoint + createInstanceHttpUrl
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", client.ProjectID)

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	params, err := buildCreateInstanceBodyParams(d, cfg)
	if err != nil {
		return "", err
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(params)
	createInstanceResp, err := client.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", fmt.Errorf("error creating CBHInstance: err: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", err
	}

	instanceKey, err := jmespath.Search("instance_key", createInstanceRespBody)
	if err != nil {
		return "", fmt.Errorf("error creating CbhInstance: instance_key is not found in API response")
	}
	return instanceKey, nil
}

func createOrder(d *schema.ResourceData, cfg *config.Config, client *golangsdk.ServiceClient,
	region, instanceKey string) (string, error) {
	var (
		createOrderHttpUrl = "v1/{project_id}/cbs/period/order"
	)

	createOrderPath := client.Endpoint + createOrderHttpUrl
	createOrderPath = strings.ReplaceAll(createOrderPath, "{project_id}", client.ProjectID)

	createOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	productId, err := getOrderProductId(d, cfg, region)
	if err != nil {
		return "", err
	}
	createOrderOpt.JSONBody = utils.RemoveNil(buildCreateOrderParams(d, strings.TrimSpace(productId), region, instanceKey))
	createOrderResp, err := client.Request("POST", createOrderPath, &createOrderOpt)
	if err != nil {
		return "", fmt.Errorf("error creating CBHOrder: %s", err)
	}

	createOrderRespBody, err := utils.FlattenResponse(createOrderResp)
	if err != nil {
		return "", err
	}

	orderId, err := jmespath.Search("order_id", createOrderRespBody)
	if err != nil {
		return "", fmt.Errorf("error creating CBHOrder: order ID is not found in API response")
	}
	return orderId.(string), nil
}

func getInstanceIdByResourceId(client *golangsdk.ServiceClient, resourceId string) (string, error) {
	instances, err := getInstanceList(client)
	if err != nil {
		return "", err
	}
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["resourceId"] != nil && instance["resourceId"].(string) == resourceId {
			return instance["instanceId"].(string), nil
		}
	}
	return "", fmt.Errorf("error getting instance_id by resource_id: %s", resourceId)
}

func getInstanceList(client *golangsdk.ServiceClient) ([]interface{}, error) {
	// getCbhInstances: Query the List of CBH instances
	var (
		getCbhInstancesHttpUrl = "v1/{project_id}/cbs/instance/list"
	)

	getCbhInstancesPath := client.Endpoint + getCbhInstancesHttpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)

	getCbhInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCbhInstancesResp, err := client.Request("GET", getCbhInstancesPath, &getCbhInstancesOpt)

	if err != nil {
		return nil, err
	}

	getCbhInstancesRespBody, err := utils.FlattenResponse(getCbhInstancesResp)
	if err != nil {
		return nil, err
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})
	return instances, nil
}

func payOrder(ctx context.Context, d *schema.ResourceData, cfg *config.Config, orderId string) (string, error) {
	region := cfg.GetRegion(d)
	var (
		payOrderHttpUrl = "v3/orders/customer-orders/pay"
		payOrderProduct = "bssv2"
	)
	bssClient, err := cfg.NewServiceClient(payOrderProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS v2 Client: %s", err)
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
		return "", fmt.Errorf("error pay CBH order(%s): %s", orderId, err)
	}
	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", err
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", fmt.Errorf("error waiting for CBH instance order %s complete: %s", orderId, err)
	}
	return resourceId, err
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	params, err := buildCreateInstanceParams(d, cfg)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"server": params,
	}
	return bodyParams, nil
}

func buildCreateInstanceParams(d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	publicIp, err := buildCreateInstancePublicIpChildBody(d, cfg)
	if err != nil {
		return nil, err
	}
	bodyParams := map[string]interface{}{
		"flavor_ref":        utils.ValueIngoreEmpty(d.Get("flavor_id")),
		"instance_name":     utils.ValueIngoreEmpty(d.Get("name")),
		"vpc_id":            utils.ValueIngoreEmpty(d.Get("vpc_id")),
		"nics":              buildCreateInstanceNicsChildBody(d),
		"public_ip":         publicIp,
		"security_groups":   buildCreateInstanceSecurityGroupsChildBody(d),
		"availability_zone": utils.ValueIngoreEmpty(d.Get("availability_zone")),
		"region":            cfg.GetRegion(d),
		"hx_password":       utils.ValueIngoreEmpty(d.Get("password")),
		"bastion_type":      "OEM",
		"ipv6_enable":       utils.ValueIngoreEmpty(d.Get("ipv6_enable")),
	}
	return bodyParams, nil
}

func buildCreateOrderParams(d *schema.ResourceData, productId, region, instanceKey string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_key":       instanceKey,
		"region_id":          region,
		"cloud_service_type": "hws.service.type.cbh",
		"period_num":         utils.ValueIngoreEmpty(d.Get("period")),
		"subscription_num":   1,
		"product_infos":      buildCreateInstanceProductInfoChildBody(d, productId),
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParams["charging_mode"] = 0
	}
	if d.Get("period_unit").(string) == "year" {
		bodyParams["period_type"] = 3
	} else {
		bodyParams["period_type"] = 2
	}
	if d.Get("auto_renew").(string) == "true" {
		bodyParams["is_auto_renew"] = 1
	} else {
		bodyParams["is_auto_renew"] = 0
	}
	return bodyParams
}

func buildCreateInstanceNicsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"subnet_id":  utils.ValueIngoreEmpty(d.Get("subnet_id")),
			"ip_address": utils.ValueIngoreEmpty(d.Get("subnet_address")),
		},
	}
}

func buildCreateInstancePublicIpChildBody(d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	publicIpId := d.Get("public_ip_id").(string)
	if publicIpId == "" {
		return nil, nil
	}

	publicIp := d.Get("public_ip").(string)
	if publicIp == "" {
		address, err := getPublicAddressById(d, cfg, publicIpId)
		if err != nil {
			return nil, err
		}
		publicIp = address
	}
	params := map[string]interface{}{
		"id":         publicIpId,
		"public_eip": publicIp,
	}
	return params, nil
}

func buildCreateInstanceSecurityGroupsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
		},
	}
}

func buildCreateInstanceProductInfoChildBody(d *schema.ResourceData, productId string) interface{} {
	param := map[string]interface{}{
		"product_id":               productId,
		"cloud_service_type":       "hws.service.type.cbh",
		"resource_type":            "hws.resource.type.cbh.ins",
		"resource_spec_code":       utils.ValueIngoreEmpty(d.Get("flavor_id")),
		"resource_size_measure_id": "14",
		"resource_size":            "1",
	}

	return []map[string]interface{}{param}
}

func buildPayOrderBodyParams(orderId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"order_id":     orderId,
		"use_coupon":   "NO",
		"use_discount": "NO",
	}
	return bodyParams
}

func getOrderProductId(d *schema.ResourceData, cfg *config.Config, region string) (string, error) {
	var (
		getCbhOrderProductIdHttpUrl = "v2/bills/ratings/period-resources/subscribe-rate"
		getCbhOrderProductIdProduct = "bss"
	)
	getCbhOrderProductIdClient, err := cfg.NewServiceClient(getCbhOrderProductIdProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS Client: %s", err)
	}

	getCbhOrderProductIdPath := getCbhOrderProductIdClient.Endpoint + getCbhOrderProductIdHttpUrl

	getCbhOrderProductIdOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCbhOrderProductIdOpt.JSONBody = utils.RemoveNil(buildGetCbhFlavorsBodyParams(d,
		getCbhOrderProductIdClient.ProjectID, region))
	getCbhOrderProductIdResp, err := getCbhOrderProductIdClient.Request("POST",
		getCbhOrderProductIdPath, &getCbhOrderProductIdOpt)

	if err != nil {
		return "", fmt.Errorf("error getting CBH order product id: %s", err)
	}

	getCbhOrderProductIdRespBody, err := utils.FlattenResponse(getCbhOrderProductIdResp)
	if err != nil {
		return "", err
	}
	curJson := utils.PathSearch("official_website_rating_result.product_rating_results",
		getCbhOrderProductIdRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return "", fmt.Errorf("fail to get CBH order product id")
	}
	productId := utils.PathSearch("product_id", curArray[0], "")
	return productId.(string), nil
}

func buildGetCbhFlavorsBodyParams(d *schema.ResourceData, projectId, region string) map[string]interface{} {
	periodUnit := d.Get("period_unit").(string)
	var periodType string
	if periodUnit == "month" {
		periodType = "2"
	} else {
		periodType = "3"
	}

	params := make(map[string]interface{})
	params["id"] = "1"
	params["cloud_service_type"] = "hws.service.type.cbh"
	params["resource_type"] = "hws.resource.type.cbh.ins"
	params["resource_spec"] = utils.ValueIngoreEmpty(d.Get("flavor_id"))
	params["region"] = region
	params["period_type"] = periodType
	params["period_num"] = utils.ValueIngoreEmpty(d.Get("period"))
	params["subscription_num"] = "1"

	bodyParams := map[string]interface{}{
		"project_id":    projectId,
		"product_infos": []map[string]interface{}{params},
	}
	return bodyParams
}

func resourceCBHInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateInstance: update the CBH instance
	var (
		bindEipHttpUrl           = "v1/{project_id}/cbs/instance/{server_id}/eip/bind"
		unbindEipHttpUrl         = "v1/{project_id}/cbs/instance/{server_id}/eip/unbind"
		updateAdminPassword      = "v1/{project_id}/cbs/instance/password"
		updateCBHInstanceProduct = "cbh"
	)

	updateCbhInstanceClient, err := cfg.NewServiceClient(updateCBHInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH Client: %s", err)
	}

	if d.HasChanges("public_ip_id", "public_ip") {
		oPublicIpIdRaw, nPublicIpIdRaw := d.GetChange("public_ip_id")
		oPublicIpRaw, nPublicIpRaw := d.GetChange("public_ip")
		oPublicIpId := strings.TrimSpace(oPublicIpIdRaw.(string))
		nPublicIpId := strings.TrimSpace(nPublicIpIdRaw.(string))
		oPublicIp := oPublicIpRaw.(string)
		nPublicIp := nPublicIpRaw.(string)

		if oPublicIpId == nPublicIpId && oPublicIp != nPublicIp {
			return diag.Errorf("the public ip is not match the public ip id")
		}
		if len(oPublicIpId) > 0 {
			err = unbindEip(d, updateCbhInstanceClient, oPublicIpId, oPublicIp, unbindEipHttpUrl)
			if err != nil {
				return diag.Errorf("error unbind eip from CBH instance: %s", err)
			}
		}
		if len(nPublicIpId) > 0 {
			err = bindEip(d, updateCbhInstanceClient, cfg, nPublicIpId, nPublicIp, bindEipHttpUrl)
			if err != nil {
				// if bind new eip fail, then bind the old eip to CBH instance
				if len(oPublicIpId) > 0 {
					_ = bindEip(d, updateCbhInstanceClient, cfg, oPublicIpId, oPublicIp, bindEipHttpUrl)
				}
				return diag.Errorf("error bind eip to CBH instance: %s", err)
			}
		}
	}

	if d.HasChanges("password") {
		err = updatePassword(d, updateCbhInstanceClient, updateAdminPassword)
		if err != nil {
			return diag.Errorf("error update CBH admin password: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		resourceId, err := getResourceIdByInstanceId(updateCbhInstanceClient, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CBH instance (%s): %s", d.Id(), err)
		}
	}
	return resourceCBHInstanceRead(ctx, d, meta)
}

func bindEip(d *schema.ResourceData, client *golangsdk.ServiceClient, cfg *config.Config, publicIpId,
	publicIp, httpUrl string) error {
	getCbhInstancesPath := client.Endpoint + httpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{server_id}", d.Id())

	bindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	if publicIp == "" {
		address, err := getPublicAddressById(d, cfg, publicIpId)
		if err != nil {
			return err
		}
		publicIp = address
	}
	bindEipOpt.JSONBody = utils.RemoveNil(buildUpdateEipBodyParams(publicIpId, publicIp))
	bindEipResp, err := client.Request("POST", getCbhInstancesPath, &bindEipOpt)
	if err != nil {
		return fmt.Errorf("error bind EIP to CBH instance: %s", err)
	}
	bindEipRespBody, err := utils.FlattenResponse(bindEipResp)
	if err != nil {
		return err
	}
	bindEipRespBodyBytes, _ := json.Marshal(bindEipRespBody)
	if !strings.Contains(string(bindEipRespBodyBytes), "success") {
		return fmt.Errorf("fail bind eip to CBH instance: %s", string(bindEipRespBodyBytes))
	}
	return nil
}

func unbindEip(d *schema.ResourceData, client *golangsdk.ServiceClient, publicIpId, publicIp, httpUrl string) error {
	getCbhInstancesPath := client.Endpoint + httpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{server_id}", d.Id())

	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	unbindEipOpt.JSONBody = utils.RemoveNil(buildUpdateEipBodyParams(publicIpId, publicIp))
	unbindEipResp, err := client.Request("POST", getCbhInstancesPath, &unbindEipOpt)
	if err != nil {
		if apiErr, ok := err.(golangsdk.ErrDefault400); ok {
			var respBody interface{}
			if jsonErr := json.Unmarshal(apiErr.Body, &respBody); jsonErr != nil {
				return jsonErr
			}
			errCode := utils.PathSearch("error_code", respBody, "")
			// these two error code indicate the eip has been unbound or deleted
			if errCode == "CBH.10020010" || errCode == "CBH.10020009" {
				log.Printf("[WARN] Failed to unbind EIP (ID: %s, IP: %s) from CBH instance(%s): %s", publicIpId,
					publicIp, d.Id(), apiErr)
				return nil
			}
		}
		return fmt.Errorf("error unbind EIP from CBH instance: %s", err)
	}
	unbindEipRespBody, err := utils.FlattenResponse(unbindEipResp)
	if err != nil {
		return err
	}
	unbindEipRespBodyBytes, _ := json.Marshal(unbindEipRespBody)
	if !strings.Contains(string(unbindEipRespBodyBytes), "success") {
		return fmt.Errorf("fail unbind eip from CBH instance: %s", string(unbindEipRespBodyBytes))
	}
	return nil
}

func updatePassword(d *schema.ResourceData, client *golangsdk.ServiceClient, httpUrl string) error {
	updatePasswordPath := client.Endpoint + httpUrl
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", client.ProjectID)

	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updatePasswordOpt.JSONBody = utils.RemoveNil(buildUpdateAdminPasswordParams(d))
	updatePasswordResp, err := client.Request("PUT", updatePasswordPath, &updatePasswordOpt)
	if err != nil {
		return fmt.Errorf("error update CBH instance password: %s", err)
	}

	updatePasswordRespBody, err := utils.FlattenResponse(updatePasswordResp)
	if err != nil {
		return err
	}
	updatePasswordRespBodyBytes, _ := json.Marshal(updatePasswordRespBody)
	if !strings.Contains(string(updatePasswordRespBodyBytes), "success") {
		return fmt.Errorf("fail update CBH instance password: %s", string(updatePasswordRespBodyBytes))
	}
	return nil
}

func buildUpdateEipBodyParams(publicIpId, publicEip string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publicip_id": publicIpId,
		"public_eip":  publicEip,
	}
	return bodyParams
}

func buildUpdateAdminPasswordParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_password": utils.ValueIngoreEmpty(d.Get("password")),
		"server_id":    d.Id(),
	}
	return bodyParams
}

func resourceCBHInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getCbhInstanceProduct = "cbh"
	)

	getCbhInstanceClient, err := cfg.NewServiceClient(getCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH Client: %s", err)
	}

	var mErr *multierror.Error
	instances, err := getInstanceList(getCbhInstanceClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(instances) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	for index, v := range instances {
		instance := v.(map[string]interface{})
		if instance["instanceId"].(string) != d.Id() {
			if index == len(instances)-1 {
				return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
			}
			continue
		}
		publicIpId := instance["publicId"]
		var publicIp string
		if publicIpId != nil && strings.TrimSpace(publicIpId.(string)) != "" {
			publicIp, err = getPublicAddressById(d, cfg, strings.TrimSpace(publicIpId.(string)))
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); !ok {
					return diag.FromErr(err)
				}
				publicIpId = ""
			}
		}
		mErr = multierror.Append(
			mErr,
			d.Set("region", region),
			d.Set("public_ip_id", publicIpId),
			d.Set("public_ip", publicIp),
			d.Set("name", instance["name"]),
			d.Set("private_ip", instance["privateIp"]),
			d.Set("status", instance["status"]),
			d.Set("vpc_id", instance["vpcId"]),
			d.Set("subnet_id", instance["subnetId"]),
			d.Set("security_group_id", instance["securityGroupId"]),
			d.Set("flavor_id", instance["specification"]),
			d.Set("availability_zone", instance["zone"]),
			d.Set("version", instance["bastionVersion"]),
		)
		break
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPublicAddressById(d *schema.ResourceData, cfg *config.Config, publicIpId string) (string, error) {
	region := cfg.GetRegion(d)
	networkingClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return "", fmt.Errorf("error creating VPC v1 client: %s", err)
	}
	publicIp, err := eips.Get(networkingClient, publicIpId).Extract()
	if err != nil {
		return "", err
	}
	return publicIp.PublicAddress, nil
}

func resourceCBHInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getCbhInstancesProduct = "cbh"
	)

	getCbhInstanceClient, err := cfg.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH Client: %s", err)
	}

	id := d.Id()
	resourceId, err := getResourceIdByInstanceId(getCbhInstanceClient, id)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	if v, ok := d.GetOk("charging_mode"); !ok || v.(string) != "prePaid" {
		return diag.Errorf("only the charging_mode of prePaid is support")
	}
	if err = common.UnsubscribePrePaidResource(d, cfg, []string{resourceId}); err != nil {
		return diag.Errorf("error unsubscribe CBH instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    cbhInstanceStateRefreshFunc(getCbhInstanceClient, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("timeout waiting for vault deletion to complete: %s", err)
	}
	d.SetId("")

	return nil
}

func getResourceIdByInstanceId(client *golangsdk.ServiceClient, instanceId string) (string, error) {
	instances, err := getInstanceList(client)
	if err != nil {
		return "", err
	}
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["instanceId"].(string) == instanceId {
			return instance["resourceId"].(string), nil
		}
	}
	return "", fmt.Errorf("error get resource_id by instance_id: %s", instanceId)
}

func cbhInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getCbhInstancesPath := client.Endpoint + "v1/{project_id}/cbs/instance/list"
		getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)
		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		v, err := client.Request("GET", getCbhInstancesPath, &getRocketmqInstanceOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}
		respBody, err := utils.FlattenResponse(v)
		if err != nil {
			return nil, "", err
		}
		instances := utils.PathSearch("instance", respBody, make([]interface{}, 0)).([]interface{})
		for _, value := range instances {
			instance := value.(map[string]interface{})
			if instance["instanceId"].(string) == instanceID {
				return instance, instance["status"].(string), nil
			}
		}
		return respBody, "DELETED", nil
	}
}
