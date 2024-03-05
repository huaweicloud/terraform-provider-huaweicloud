package cbh

import (
	"context"
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

// @API CBH GET /v1/{project_id}/cbs/instance/list
// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/bind
// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/unbind
// @API CBH POST /v2/{project_id}/cbs/instance
// @API CBH GET /v2/{project_id}/cbs/instance/list
// @API CBH PUT /v2/{project_id}/cbs/instance/password
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
				Description: `Specifies the charging mode of the CBH instance.`,
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
				Description:  `Specifies the charging period of the CBH instance.`,
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
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether the IPv6 network is enabled.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Computed; The elastic IP address.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private IP of the instance.`,
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
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		createCbhInstanceProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(createCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	orderId, err := createCBHInstance(client, d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for CBH instance order %s complete: %s", orderId, err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}
	expression := fmt.Sprintf("[?resource_info.resource_id == '%s']|[0].server_id", resourceId)
	serverId, err := jmespath.Search(expression, instances)
	if err != nil || serverId == nil {
		return diag.Errorf("error creating CBH instance: ID is not found in API response")
	}

	d.SetId(serverId.(string))
	if err := waitingForCBHInstanceActive(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CBH instance (%s) creation to active: %s", d.Id(), err)
	}

	return resourceCBHInstanceRead(ctx, d, meta)
}

func createCBHInstance(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (string, error) {
	publicIp, err := buildCreateNetworkPublicIpBodyParam(d, cfg)
	if err != nil {
		return "", fmt.Errorf("error creating CBH instance: error building network public IP body: %s", err)
	}

	createInstancePath := client.Endpoint + "v2/{project_id}/cbs/instance"
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", client.ProjectID)
	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCBHInstanceBodyParam(d, cfg.GetRegion(d), publicIp)),
	}

	createInstanceResp, err := client.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", fmt.Errorf("error creating CBH instance: %s", err)
	}
	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", err
	}

	orderId, err := jmespath.Search("order_id", createInstanceRespBody)
	if err != nil || orderId == nil {
		return "", fmt.Errorf("error creating CBH instance: order_id is not found in API response")
	}
	return orderId.(string), nil
}

func buildCreateCBHInstanceBodyParam(d *schema.ResourceData, region string, publicIp interface{}) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"specification":     d.Get("flavor_id"),
		"instance_name":     d.Get("name"),
		"password":          d.Get("password"),
		"region":            region,
		"availability_zone": d.Get("availability_zone"),
		"charging_mode":     buildCreateChargingModeParam(d),
		"period_type":       buildCreatePeriodTypeParam(d),
		"period_num":        utils.ValueIngoreEmpty(d.Get("period")),
		"is_auto_renew":     buildCreateIsAutoRenewParam(d),
		"is_auto_pay":       1,
		"network":           buildCreateNetworkBodyParam(d, publicIp),
	}
	// The default value of the field `ipv6_enable` is false
	if d.Get("ipv6_enable").(bool) {
		bodyParam["ipv6_enable"] = true
	}
	return bodyParam
}

// Currently, the CBH instance only supports prePaid charging mode
func buildCreateChargingModeParam(d *schema.ResourceData) interface{} {
	if d.Get("charging_mode").(string) == "prePaid" {
		return 0
	}
	return nil
}

// Currently, the CBH instance only supports `year` and `month` period unit in prePaid charging mode.
func buildCreatePeriodTypeParam(d *schema.ResourceData) interface{} {
	if d.Get("period_unit").(string) == "year" {
		return 3
	}

	if d.Get("period_unit").(string) == "month" {
		return 2
	}
	return nil
}

// `1` indicates automatic renewal. `0` indicates non-automatic renewal.
func buildCreateIsAutoRenewParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_renew").(string) == "true" {
		return 1
	}

	if d.Get("auto_renew").(string) == "false" {
		return 0
	}
	return nil
}

func buildCreateNetworkBodyParam(d *schema.ResourceData, publicIp interface{}) interface{} {
	return map[string]interface{}{
		"vpc_id":          d.Get("vpc_id"),
		"subnet_id":       d.Get("subnet_id"),
		"public_ip":       publicIp,
		"security_groups": buildCreateNetworkSecurityGroupsBodyParam(d),
		"private_ip":      buildCreateNetworkPrivateIpBodyParam(d),
	}
}

func buildCreateNetworkPublicIpBodyParam(d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	publicIpId := d.Get("public_ip_id").(string)
	if publicIpId == "" {
		return nil, nil
	}

	publicIp, err := getPublicAddressById(d, cfg, publicIpId)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"id":         publicIpId,
		"public_eip": publicIp,
	}
	return params, nil
}

func buildCreateNetworkSecurityGroupsBodyParam(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"id": d.Get("security_group_id"),
		},
	}
}

func buildCreateNetworkPrivateIpBodyParam(d *schema.ResourceData) interface{} {
	if v, ok := d.GetOk("subnet_address"); ok {
		return map[string]interface{}{
			"ip": v,
		}
	}
	return nil
}

func getCBHInstanceList(client *golangsdk.ServiceClient) ([]interface{}, error) {
	getCbhInstancesPath := client.Endpoint + "v2/{project_id}/cbs/instance/list"
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)
	getCbhInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCbhInstancesResp, err := client.Request("GET", getCbhInstancesPath, &getCbhInstancesOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH instance list: %s", err)
	}

	getCbhInstancesRespBody, err := utils.FlattenResponse(getCbhInstancesResp)
	if err != nil {
		return nil, err
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})
	return instances, nil
}

func waitingForCBHInstanceActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	unexpectedStatus := []string{"SHUTOFF", "DELETING", "DELETED", "ERROR", "FROZEN"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "ACTIVE" {
				return instance, "COMPLETED", nil
			}

			if status == "" {
				return instance, "ERROR", fmt.Errorf("status is not found in list API response")
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return instance, status, nil
			}
			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// This method will be deleted when the datasource API is upgraded.
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

func resourceCBHInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		updateCBHInstanceProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(updateCBHInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	if d.HasChanges("public_ip_id") {
		if err := updatePublicIpId(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("password") {
		if err = updatePassword(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		instances, err := getCBHInstanceList(client)
		if err != nil {
			return diag.FromErr(err)
		}
		expression := fmt.Sprintf("[?server_id == '%s']|[0].resource_info.resource_id", d.Id())
		resourceId := utils.PathSearch(expression, instances, "").(string)
		if resourceId == "" {
			return diag.Errorf("error updating the auto-renew of the CBH instance (%s): "+
				"resource ID is not found in list API response", d.Id())
		}

		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CBH instance (%s): %s", d.Id(), err)
		}
	}
	return resourceCBHInstanceRead(ctx, d, meta)
}

func updatePublicIpId(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oPublicIpIdRaw, nPublicIpIdRaw := d.GetChange("public_ip_id")
	oPublicIpId := strings.TrimSpace(oPublicIpIdRaw.(string))
	nPublicIpId := strings.TrimSpace(nPublicIpIdRaw.(string))

	if len(oPublicIpId) > 0 {
		if err := unbindEip(client, d, oPublicIpId); err != nil {
			return err
		}
	}
	if len(nPublicIpId) > 0 {
		if err := bindEip(client, d, nPublicIpId); err != nil {
			// if bind new eip fail, then bind the old eip to CBH instance
			if len(oPublicIpId) > 0 {
				if err := bindEip(client, d, oPublicIpId); err != nil {
					log.Printf("[WARN] error bind old EIP: %s", err)
				}
			}
			return err
		}
	}
	return nil
}

func unbindEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId string) error {
	unbindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/unbind"
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{server_id}", d.Id())
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", unbindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error unbind EIP (%s) from CBH instance (%s): %s", publicIpId, d.Id(), err)
	}
	return nil
}

func bindEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId string) error {
	bindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/bind"
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{server_id}", d.Id())
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", bindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error bind EIP (%s) to CBH instance (%s): %s", publicIpId, d.Id(), err)
	}
	return nil
}

func updatePassword(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePasswordPath := client.Endpoint + "v2/{project_id}/cbs/instance/password"
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", client.ProjectID)
	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateAdminPasswordParams(d),
	}
	_, err := client.Request("PUT", updatePasswordPath, &updatePasswordOpt)
	if err != nil {
		return fmt.Errorf("error update CBH instance (%s) password: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateAdminPasswordParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"new_password": d.Get("password"),
		"server_id":    d.Id(),
	}
}

func resourceCBHInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr                  *multierror.Error
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		getCbhInstanceProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(getCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	// When EIP is not configured, the query interface field will return a space string.
	publicIpId := strings.TrimSpace(utils.PathSearch("network.public_id", instance, "").(string))

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("public_ip_id", publicIpId),
		d.Set("public_ip", utils.PathSearch("network.public_ip", instance, nil)),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("private_ip", utils.PathSearch("network.private_ip", instance, nil)),
		d.Set("subnet_address", utils.PathSearch("network.private_ip", instance, nil)),
		d.Set("status", utils.PathSearch("status_info.status", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("network.vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("network.subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("network.security_group_id", instance, nil)),
		d.Set("flavor_id", utils.PathSearch("resource_info.specification", instance, nil)),
		d.Set("availability_zone", utils.PathSearch("az_info.zone", instance, nil)),
		d.Set("version", utils.PathSearch("bastion_version", instance, nil)),
	)
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
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		getCbhInstancesProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	resourceId := utils.PathSearch("resource_info.resource_id", instance, "").(string)
	if resourceId == "" {
		return diag.Errorf("error deleting the CBH instance (%s): resource ID is not found in list API response", d.Id())
	}

	if err = common.UnsubscribePrePaidResource(d, cfg, []string{resourceId}); err != nil {
		return diag.Errorf("error unsubscribe CBH instance (%s): %s", d.Id(), err)
	}

	if err := waitingForCBHInstanceDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for CBH instance (%s) deleted: %s", d.Id(), err)
	}
	return nil
}

func waitingForCBHInstanceDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				m := map[string]string{"code": "COMPLETED"}
				return m, "COMPLETED", nil
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "DELETED" {
				return instance, "COMPLETED", nil
			}
			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
