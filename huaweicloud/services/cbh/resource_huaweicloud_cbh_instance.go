package cbh

import (
	"context"
	"encoding/json"
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
	"github.com/jmespath/go-jmespath"
)

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
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product ID of the CBH server.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the CBH instance.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID list of the security group.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the availability zone name.`,
			},
			"hx_password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the front end login password.`,
				Sensitive:   true,
			},
			"bastion_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the bastion.`,
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
				Required:    true,
				ForceNew:    true,
				Description: `Specifies whether auto renew is enabled. Valid values are "true" and "false". Defaults to **false**.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies a image ID.`,
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the inject user data.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the initial password.`,
				Sensitive:   true,
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the secret key of the admin.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP address of the subnet.`,
			},
			"public_ip": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     CBHInstancePublicIPSchema(),
				Optional: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     CBHInstanceRootVolumeSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_volume": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     CBHInstanceDataVolumeSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"slave_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Specifies the slave availability zone name. The slave machine will be created when
this field is not empty.`,
			},
			"metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the metadata of the service.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the IPv6 network is enabled.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the end time.`,
			},
			"relative_resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the new capacity expansion.`,
			},
			"product_info": {
				Type:     schema.TypeList,
				Elem:     CBHInstanceProductInfoSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"create", "renewals", "change",
				}, false),
				Description: `Specifies the type of the network operation.`,
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the elastic IP.`,
			},
			"exp_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expire time of the instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time of the instance.`,
			},
			"release_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the release time of the instance.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the server id of the instance.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private ip of the instance.`,
			},
			"task_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task status of the instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether the instance image can be upgraded.`,
			},
			"instance_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the instance.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the resource.`,
			},
			"alter_permit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether the front-end displays the capacity expansion button.`,
			},
			"bastion_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version of the instance image.`,
			},
			"new_bastion_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest version of the instance image.`,
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the bastion.`,
			},
		},
	}
}

func CBHInstancePublicIPSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the elastic IP.`,
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the elastic IP address.`,
			},
			"eip": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     CBHInstancePublicIPEipSchema(),
				Optional: true,
			},
		},
	}
	return &sc
}

func CBHInstancePublicIPEipSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of EIP.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the product ID of the IP associated with.`,
			},
			"bandwidth": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     CBHInstanceEipBandwidthSchema(),
				Optional: true,
			},
		},
	}
	return &sc
}

func CBHInstanceEipBandwidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the size of the bandwidth.`,
			},
			"share_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the share type. Only PER is supported.`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the charge type. The value can be traffic or empty.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the product ID of the bandwidth associated with.`,
			},
		},
	}
	return &sc
}

func CBHInstanceRootVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the size of the root volume, unit is GB.`,
			},
			"extend_param": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the info of the volume.`,
			},
		},
	}
	return &sc
}

func CBHInstanceDataVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of volume.`,
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the size of the data volume, unit is GB.`,
			},
			"extend_param": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the info of the volume.`,
			},
		},
	}
	return &sc
}

func CBHInstanceProductInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ID of the product.`,
			},
			"resource_size_measure_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resource capacity measurement ID.`,
			},
			"resource_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the size of the resource capacity.`,
			},
		},
	}
	return &sc
}

func resourceCBHInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createInstance: create CBH instance
	instanceKey, slaveInstanceKey, err := createInstance(d, config, region)
	if err != nil {
		return err
	}
	if instanceKey == nil {
		return diag.Errorf("error creating CbhInstance: instance_key is empty")
	}

	// createOrder: create instance order
	resourceId, err := createOrder(ctx, d, config, region, instanceKey.(string))
	if err != nil {
		return err
	}
	// createSlaveOrder: create slave instance order
	if slaveInstanceKey != nil {
		_, err = createOrder(ctx, d, config, region, slaveInstanceKey.(string))
		if err != nil {
			return err
		}
	}

	instanceId, err := getInstanceIdByResourceId(d, config, region, resourceId)
	if err != nil {
		return err
	}

	d.SetId(instanceId)

	return resourceCBHInstanceRead(ctx, d, meta)
}

func createInstance(d *schema.ResourceData, config *config.Config, region string) (interface{}, interface{}, diag.Diagnostics) {
	// createInstance: create CBH instance
	var (
		createInstanceHttpUrl = "v1/{project_id}/cbs/instance/create"
		createInstanceProduct = "cbh"
	)
	createInstanceClient, err := config.NewServiceClient(createInstanceProduct, region)
	if err != nil {
		return "", "", diag.Errorf("error creating CBHInstance Client: %s", err)
	}

	createInstancePath := createInstanceClient.Endpoint + createInstanceHttpUrl
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", createInstanceClient.ProjectID)

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d, config))
	createInstanceResp, err := createInstanceClient.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", "", diag.Errorf("error creating CBHInstance: err: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", "", diag.FromErr(err)
	}

	instanceKey, err := jmespath.Search("instance_key", createInstanceRespBody)
	if err != nil {
		return "", "", diag.Errorf("error creating CbhInstance: instance_key is not found in API response")
	}

	slaveInstanceKey, err := jmespath.Search("slaveInstanceKey", createInstanceRespBody)
	if err != nil {
		return "", "", diag.Errorf("error creating CbhInstance: slaveInstanceKey is not found in API response")
	}
	return instanceKey, slaveInstanceKey, nil
}

func createOrder(ctx context.Context, d *schema.ResourceData, config *config.Config, region,
	instanceKey string) (string, diag.Diagnostics) {
	var (
		createInstanceHttpUrl = "v1/{project_id}/cbs/period/order"
		payOrderHttpUrl       = "v3/orders/customer-orders/pay"
		createInstanceProduct = "cbh"
		payOrderProduct       = "bssv2"
	)
	createInstanceClient, err := config.NewServiceClient(createInstanceProduct, region)
	if err != nil {
		return "", diag.Errorf("error creating CBHInstance Client: %s", err)
	}

	createInstancePath := createInstanceClient.Endpoint + createInstanceHttpUrl
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", createInstanceClient.ProjectID)

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(buildCreateOrderParams(d, config, region, instanceKey))
	createInstanceResp, err := createInstanceClient.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", diag.Errorf("error creating CBHOrder: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", diag.FromErr(err)
	}

	orderId, err := jmespath.Search("order_id", createInstanceRespBody)
	if err != nil {
		return "", diag.Errorf("error creating CBHOrder: order ID is not found in API response")
	}

	bssClient, err := config.NewServiceClient(payOrderProduct, region)
	if err != nil {
		return "", diag.Errorf("error creating BSS v2 Client: %s", err)
	}

	payOrderPath := bssClient.Endpoint + payOrderHttpUrl

	// pay order
	err = payOrder(bssClient, orderId.(string), payOrderPath)
	if err != nil {
		return "", diag.Errorf("error pay CBHOrder: %s", err)
	}

	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId.(string),
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return "", diag.Errorf("error waiting for replica order resource %s complete: %s", orderId.(string), err)
	}

	return resourceId, nil
}

func getInstanceIdByResourceId(d *schema.ResourceData, config *config.Config, region, resourceId string) (string,
	diag.Diagnostics) {
	instances, err := getInstanceList(d, config, region)
	if err != nil {
		return "", diag.Errorf("%s", err)
	}
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["resource_id"].(string) == resourceId {
			return instance["instance_id"].(string), nil
		}
	}
	return "", diag.Errorf("error get instance by resource_id: %s", resourceId)
}

func getInstanceList(d *schema.ResourceData, config *config.Config, region string) ([]interface{}, error) {
	// getCbhInstances: Query the List of CBH instances
	var (
		getCbhInstancesHttpUrl = "v1/{project_id}/cbs/instance/list"
		getCbhInstancesProduct = "cbh"
	)
	getCbhInstancesClient, err := config.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CbhInstances Client: %s", err)
	}

	getCbhInstancesPath := getCbhInstancesClient.Endpoint + getCbhInstancesHttpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", getCbhInstancesClient.ProjectID)

	getCbhInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCbhInstancesResp, err := getCbhInstancesClient.Request("GET", getCbhInstancesPath, &getCbhInstancesOpt)

	if err != nil {
		return nil, err
	}

	getCbhInstancesRespBody, err := utils.FlattenResponse(getCbhInstancesResp)
	if err != nil {
		return nil, err
	}
	return flattenGetInstancesResponseBodyInstance(getCbhInstancesRespBody), nil
}

func payOrder(bssClient *golangsdk.ServiceClient, orderId, payOrderPath string) error {
	payOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	payOrderOpt.JSONBody = utils.RemoveNil(buildPayOrderBodyParams(orderId))
	_, err := bssClient.Request("POST", payOrderPath, &payOrderOpt)
	return err
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server": buildCreateInstanceParams(d, config),
	}
	return bodyParams
}

func buildCreateInstanceParams(d *schema.ResourceData, config *config.Config) interface{} {
	bodyParams := map[string]interface{}{
		"image_ref":               utils.ValueIngoreEmpty(d.Get("image_id")),
		"flavor_ref":              utils.ValueIngoreEmpty(d.Get("flavor_id")),
		"instance_name":           utils.ValueIngoreEmpty(d.Get("name")),
		"user_data":               utils.ValueIngoreEmpty(d.Get("user_data")),
		"password":                utils.ValueIngoreEmpty(d.Get("password")),
		"key_name":                utils.ValueIngoreEmpty(d.Get("key_name")),
		"vpc_id":                  utils.ValueIngoreEmpty(d.Get("vpc_id")),
		"nics":                    buildCreateInstanceNicsChildBody(d),
		"public_ip":               buildCreateInstancePublicIpChildBody(d),
		"root_volume":             buildCreateInstanceRootVolumeChildBody(d),
		"data_volume":             buildCreateInstanceDataVolumeChildBody(d),
		"security_groups":         buildCreateInstanceSecurityGroupsChildBody(d),
		"availability_zone":       utils.ValueIngoreEmpty(d.Get("availability_zone")),
		"region":                  config.GetRegion(d),
		"slave_availability_zone": utils.ValueIngoreEmpty(d.Get("slave_availability_zone")),
		"metadata":                utils.ValueIngoreEmpty(d.Get("metadata")),
		"hx_password":             utils.ValueIngoreEmpty(d.Get("hx_password")),
		"bastion_type":            utils.ValueIngoreEmpty(d.Get("bastion_type")),
		"ipv6_enable":             utils.ValueIngoreEmpty(d.Get("ipv6_enable")),
		"end_time":                utils.ValueIngoreEmpty(d.Get("end_time")),
	}
	return bodyParams
}

func buildCreateOrderParams(d *schema.ResourceData, config *config.Config, region,
	instanceKey string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_key":         instanceKey,
		"region_id":            region,
		"end_time":             utils.ValueIngoreEmpty(d.Get("end_time")),
		"cloud_service_type":   "hws.service.type.cbh",
		"period_num":           utils.ValueIngoreEmpty(d.Get("period")),
		"subscription_num":     1,
		"relative_resource_id": utils.ValueIngoreEmpty(d.Get("relative_resource_id")),
		"product_infos":        buildCreateInstanceProductInfoChildBody(d),
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParams["charging_mode"] = 0
	}
	if d.Get("period_unit").(string) == "year" {
		bodyParams["period_type"] = 1
	} else {
		bodyParams["period_type"] = 2
	}
	if d.Get("auto_renew").(string) == "true" {
		bodyParams["is_auto_renew"] = 0
	} else {
		bodyParams["is_auto_renew"] = 1
	}
	return bodyParams
}

func buildCreateInstanceNicsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"subnet_id":  utils.ValueIngoreEmpty(d.Get("subnet_id")),
			"ip_address": utils.ValueIngoreEmpty(d.Get("ip_address")),
		},
	}
}

func buildCreateInstancePublicIpChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("public_ip").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"id":         utils.ValueIngoreEmpty(raw["id"]),
		"public_eip": utils.ValueIngoreEmpty(raw["address"]),
		"eip":        buildCreateInstancePublicIpEipChildBody(d),
	}

	return params
}

func buildCreateInstancePublicIpEipChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("public_ip.0.eip").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"type":      utils.ValueIngoreEmpty(raw["type"]),
		"flavor_id": utils.ValueIngoreEmpty(raw["flavor_id"]),
		"bandwidth": buildCreateInstancePublicIpEipBandwidthChildBody(d),
	}

	return params
}

func buildCreateInstancePublicIpEipBandwidthChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("public_ip.0.eip.0.bandwidth").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"size":        utils.ValueIngoreEmpty(raw["size"]),
		"share_type":  utils.ValueIngoreEmpty(raw["share_type"]),
		"charge_mode": utils.ValueIngoreEmpty(raw["charge_mode"]),
		"flavor_id":   utils.ValueIngoreEmpty(raw["flavor_id"]),
	}

	return params
}

func buildCreateInstanceRootVolumeChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("root_volume").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"volume_type":  utils.ValueIngoreEmpty(raw["type"]),
		"size":         utils.ValueIngoreEmpty(raw["size"]),
		"extend_param": utils.ValueIngoreEmpty(raw["extend_param"]),
	}

	return params
}

func buildCreateInstanceDataVolumeChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("data_volume").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"volume_type":  utils.ValueIngoreEmpty(raw["type"]),
		"size":         utils.ValueIngoreEmpty(raw["size"]),
		"extend_param": utils.ValueIngoreEmpty(raw["extend_param"]),
	}

	return params
}

func buildCreateInstanceSecurityGroupsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
		},
	}
}

func buildCreateInstanceProductInfoChildBody(d *schema.ResourceData) interface{} {
	rawParams := d.Get("product_info").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0)
	for _, rawParam := range rawParams {
		raw := rawParam.(map[string]interface{})
		param := map[string]interface{}{
			"product_id":               utils.ValueIngoreEmpty(raw["product_id"]),
			"cloud_service_type":       "hws.service.type.cbh",
			"resource_type":            "hws.resource.type.cbh.ins",
			"resource_spec_code":       utils.ValueIngoreEmpty(d.Get("flavor_id")),
			"available_zone_id":        utils.ValueIngoreEmpty(d.Get("availability_zone")),
			"resource_size_measure_id": utils.ValueIngoreEmpty(raw["resource_size_measure_id"]),
			"resource_size":            utils.ValueIngoreEmpty(raw["resource_size"]),
		}
		params = append(params, param)
	}

	return params
}

func buildPayOrderBodyParams(orderId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"order_id":     orderId,
		"use_coupon":   "NO",
		"use_discount": "NO",
	}
	return bodyParams
}

func resourceCBHInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// updateInstance: update the CBH instance
	var (
		bindEipHttpUrl        = "v1/{project_id}/cbs/instance/{server_id}/eip/bind"
		unbindEipHttpUrl      = "v1/{project_id}/cbs/instance/{server_id}/eip/unbind"
		updateAdminPassword   = "v1/{project_id}/cbs/instance/password"
		updateNetworkHttpUrl  = "v1/{project_id}/cbs/{server_id}/network/change"
		updateInstanceProduct = "cbh"
	)

	if d.HasChanges("public_ip") {
		oPublicIpIdRaw, nPublicIpIdRaw := d.GetChange("public_ip.0.id")
		oPublicIpEipRaw, nPublicIpEipRaw := d.GetChange("public_ip.0.address")
		oPublicIpId := oPublicIpIdRaw.(string)
		nPublicIpId := nPublicIpIdRaw.(string)
		oPublicIpEip := oPublicIpEipRaw.(string)
		nPublicIpEip := nPublicIpEipRaw.(string)
		operateType := 0
		if len(nPublicIpId) == 0 || len(nPublicIpEip) == 0 {
			operateType = 1
		}
		if len(oPublicIpId) == 0 || len(oPublicIpEip) == 0 {
			operateType = 2
		}
		switch operateType {
		case 1:
			err := unbindEip(d, config, region, updateInstanceProduct, oPublicIpId, oPublicIpEip, unbindEipHttpUrl)
			if err != nil {
				return diag.Errorf("error unbind eip from CBH instance: %s", err)
			}
		case 2:
			err := bindEip(d, config, region, updateInstanceProduct, nPublicIpId, nPublicIpEip, bindEipHttpUrl)
			if err != nil {
				return diag.Errorf("error bind eip to CBH instance: %s", err)
			}
		default:
			err := unbindEip(d, config, region, updateInstanceProduct, oPublicIpId, oPublicIpEip, unbindEipHttpUrl)
			if err != nil {
				return diag.Errorf("error unbind eip from CBH instance: %s", err)
			}
			err = bindEip(d, config, region, updateInstanceProduct, nPublicIpId, nPublicIpEip, bindEipHttpUrl)
			if err != nil {
				// if bind new eip fail, then bind the old eip to CBH instance
				_ = bindEip(d, config, region, updateInstanceProduct, oPublicIpId, oPublicIpEip, bindEipHttpUrl)
				return diag.Errorf("error bind eip to CBH instance: %s", err)
			}
		}
	}

	if d.HasChanges("network_type", "subnet_id", "ip_address", "security_group_id") {
		err := updateNetwork(d, config, region, updateInstanceProduct, updateNetworkHttpUrl)
		if err != nil {
			return diag.Errorf("error update CBH network: %s", err)
		}
	}

	if d.HasChanges("password") {
		err := updatePassword(d, config, region, updateInstanceProduct, updateAdminPassword)
		if err != nil {
			return diag.Errorf("error update CBH admin password: %s", err)
		}
	}
	return resourceCBHInstanceRead(ctx, d, meta)
}

func bindEip(d *schema.ResourceData, config *config.Config, region, product, publicipId, publicEip,
	httpUrl string) error {
	updateInstanceClient, err := config.NewServiceClient(product, region)
	if err != nil {
		return err
	}

	getCbhInstancesPath := updateInstanceClient.Endpoint + httpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", updateInstanceClient.ProjectID)
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{server_id}", d.Id())

	bindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	bindEipOpt.JSONBody = utils.RemoveNil(buildUpdateEipBodyParams(publicipId, publicEip))
	bindEipResp, err := updateInstanceClient.Request("POST", getCbhInstancesPath, &bindEipOpt)
	if err != nil {
		return err
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

func unbindEip(d *schema.ResourceData, config *config.Config, region, product, publicipId, publicEip,
	httpUrl string) error {
	updateInstanceClient, err := config.NewServiceClient(product, region)
	if err != nil {
		return err
	}

	getCbhInstancesPath := updateInstanceClient.Endpoint + httpUrl
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", updateInstanceClient.ProjectID)
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{server_id}", d.Id())

	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	unbindEipOpt.JSONBody = utils.RemoveNil(buildUpdateEipBodyParams(publicipId, publicEip))
	unbindEipResp, err := updateInstanceClient.Request("POST", getCbhInstancesPath, &unbindEipOpt)
	if err != nil {
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

func updateNetwork(d *schema.ResourceData, config *config.Config, region, product, httpUrl string) error {
	updateInstanceClient, err := config.NewServiceClient(product, region)
	if err != nil {
		return err
	}

	updateNetworkPath := updateInstanceClient.Endpoint + httpUrl
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{project_id}", updateInstanceClient.ProjectID)
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{server_id}", d.Id())

	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildUpdateNetworkBodyParams(d, config))
	_, err = updateInstanceClient.Request("POST", updateNetworkPath, &updateNetworkOpt)

	return err
}

func updatePassword(d *schema.ResourceData, config *config.Config, region, product, httpUrl string) error {
	updateInstanceClient, err := config.NewServiceClient(product, region)
	if err != nil {
		return err
	}

	updatePasswordPath := updateInstanceClient.Endpoint + httpUrl
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", updateInstanceClient.ProjectID)

	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updatePasswordOpt.JSONBody = utils.RemoveNil(buildUpdateAdminPasswordParams(d, config))
	updatePasswordResp, err := updateInstanceClient.Request("PUT", updatePasswordPath, &updatePasswordOpt)
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

func buildUpdateEipBodyParams(publicipId, publicEip string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"publicip_id": publicipId,
		"public_eip":  publicEip,
	}
	return bodyParams
}

func buildUpdateAdminPasswordParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_password": utils.ValueIngoreEmpty(d.Get("password")),
		"server_id":    d.Id(),
	}
	return bodyParams
}

func buildUpdateNetworkBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":            utils.ValueIngoreEmpty(d.Get("network_type")),
		"nics":            buildUpdateInstanceNicsChildBody(d),
		"security_groups": buildUpdateInstanceSecurityGroupsChildBody(d),
	}
	return bodyParams
}

func buildUpdateInstanceNicsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"subnet_id":  utils.ValueIngoreEmpty(d.Get("subnet_id").(string)),
			"ip_address": utils.ValueIngoreEmpty(d.Get("ip_address").(string)),
		},
	}
}

func buildUpdateInstanceSecurityGroupsChildBody(d *schema.ResourceData) interface{} {
	return []map[string]interface{}{
		{
			"security_group_id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
		},
	}
}

func resourceCBHInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	instances, err := getInstanceList(d, config, region)
	if err != nil {
		return diag.Errorf("%s", err)
	}
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["instance_id"].(string) == d.Id() {
			mErr = multierror.Append(
				mErr,
				d.Set("region", region),
				d.Set("publicip_id", instance["publicip_id"]),
				d.Set("exp_time", instance["exp_time"]),
				d.Set("start_time", instance["start_time"]),
				d.Set("end_time", instance["end_time"]),
				d.Set("release_time", instance["release_time"]),
				d.Set("name", instance["name"]),
				d.Set("instance_id", instance["instance_id"]),
				d.Set("private_ip", instance["private_ip"]),
				d.Set("task_status", instance["task_status"]),
				d.Set("status", instance["status"]),
				d.Set("vpc_id", instance["vpc_id"]),
				d.Set("subnet_id", instance["subnet_id"]),
				d.Set("security_group_id", instance["security_group_id"]),
				d.Set("flavor_id", instance["flavor_id"]),
				d.Set("update", instance["update"]),
				d.Set("instance_key", instance["instance_key"]),
				d.Set("resource_id", instance["resource_id"]),
				d.Set("bastion_type", instance["bastion_type"]),
				d.Set("alter_permit", instance["alter_permit"]),
				d.Set("bastion_version", instance["bastion_version"]),
				d.Set("new_bastion_version", instance["new_bastion_version"]),
				d.Set("instance_status", instance["instance_status"]),
				d.Set("description", instance["description"]),
				d.Set("auto_renew", instance["auto_renew"]),
			)
			break
		}
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCBHInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var (
		getCbhInstancesProduct = "cbh"
	)

	getCbhInstancesClient, err := config.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH Client: %s", err)
	}

	id := d.Id()
	resourceId, err := getResourceIdByInstanceId(d, config, region, id)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err = common.UnsubscribePrePaidResource(d, config, []string{resourceId}); err != nil {
			return diag.Errorf("error unsubscribe CBH instance: %s", err)
		}
	} else {
		return diag.Errorf("only the charging_mode of prePaid is support")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    CbhInstanceStateRefreshFunc(getCbhInstancesClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("timeout waiting for vault deletion to complete: %s", err)
	}
	d.SetId("")

	return nil
}

func getResourceIdByInstanceId(d *schema.ResourceData, config *config.Config, region, instanceId string) (string, error) {
	instances, err := getInstanceList(d, config, region)
	if err != nil {
		return "", err
	}
	for _, v := range instances {
		instance := v.(map[string]interface{})
		if instance["instance_id"].(string) == instanceId {
			return instance["resource_id"].(string), nil
		}
	}
	return "", fmt.Errorf("error get resource_id by instance_id: %s", instanceId)
}

func CbhInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
		instances := flattenGetInstancesResponseBodyInstance(respBody)
		for _, value := range instances {
			instance := value.(map[string]interface{})
			if instance["instance_id"].(string) == instanceID {
				return instance, instance["status"].(string), nil
			}
		}
		return respBody, "DELETED", nil
	}
}
