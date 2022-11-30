package dms

import (
	"context"
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/jmespath/go-jmespath"
)

func ResourceDmsRocketMQInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQInstanceCreate,
		UpdateContext: resourceDmsRocketMQInstanceUpdate,
		ReadContext:   resourceDmsRocketMQInstanceRead,
		DeleteContext: resourceDmsRocketMQInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
				Description: `Specifies the name of the DMS RocketMQ instance`,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9]*$`),
						"An instance name starts with a letter and can contain only letters,"+
							"digits, underscores (_), and hyphens (-)"),
					validation.StringLenBetween(4, 64),
				),
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the version of the RocketMQ engine.`,
			},
			"storage_space": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  `Specifies the message storage capacity, Unit: GB.`,
				ValidateFunc: validation.IntBetween(300, 3000),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a subnet`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a security group`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the list of availability zone names`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies a product ID`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the storage I/O specification`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  `Specifies the description of the DMS RocketMQ instance.`,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"ssl_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the RocketMQ SASL_SSL is enabled.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to support IPv6`,
			},
			"enable_publicip": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable public access.`,
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the EIP bound to the instance.`,
			},
			"broker_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the broker numbers.`,
			},
			"retention_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ACL access control.`,
			},
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lisenter_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advertised_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DMS RocketMQ instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DMS RocketMQ instance type. Value: cluster.`,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window starts. The format is HH:mm:ss.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window ends. The format is HH:mm:ss.`,
			},
			"used_storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the used message storage space. Unit: GB.`,
			},
			"publicip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public IP address.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the node quantity.`,
			},
			"new_spec_billing_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether billing based on new specifications is enabled.`,
			},
			"enable_acl": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether access control is enabled.`,
			},
			"namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the metadata address.`,
			},
			"broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the service data address.`,
			},
			"public_namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network metadata address.`,
			},
			"public_broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network service data address.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource specifications.`,
			},
		},
	}
}

func resourceDmsRocketMQInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createRocketmqInstance: create DMS rocketmq instance
	var (
		createRocketmqInstanceHttpUrl = "v2/{project_id}/instances"
		createRocketmqInstanceProduct = "dms"
	)
	createRocketmqInstanceClient, err := config.NewServiceClient(createRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	createRocketmqInstancePath := createRocketmqInstanceClient.Endpoint + createRocketmqInstanceHttpUrl
	createRocketmqInstancePath = strings.ReplaceAll(createRocketmqInstancePath, "{project_id}", createRocketmqInstanceClient.ProjectID)

	createRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	var availableZones []string
	zoneIDs, ok := d.GetOk("available_zones")
	if ok {
		availableZones = utils.ExpandToStringList(zoneIDs.([]interface{}))
	} else {
		// convert the codes of the availability zone into ids
		azCodes := d.Get("availability_zones").([]interface{})
		availableZones, err = getAvailableZoneIDByCode(config, region, azCodes)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	createRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildCreateRocketmqInstanceBodyParams(d, config, availableZones))
	createRocketmqInstanceResp, err := createRocketmqInstanceClient.Request("POST", createRocketmqInstancePath, &createRocketmqInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: %s", err)
	}

	createRocketmqInstanceRespBody, err := utils.FlattenResponse(createRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("instance_id", createRocketmqInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: ID is not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      DmsRocketmqInstanceStateRefreshFunc(createRocketmqInstanceClient, id.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to create: %s", id.(string), err)
	}

	d.SetId(id.(string))

	if _, ok = d.GetOk("cross_vpc_accesses"); ok {
		if err = updateCrossVpcAccesses(createRocketmqInstanceClient, d); err != nil {
			return diag.Errorf("Failed to update default advertised IP: %v", err)
		}
	}

	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func buildCreateRocketmqInstanceBodyParams(d *schema.ResourceData, config *config.Config,
	availableZones []string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              utils.ValueIngoreEmpty(d.Get("name")),
		"description":       utils.ValueIngoreEmpty(d.Get("description")),
		"engine":            "reliability",
		"engine_version":    utils.ValueIngoreEmpty(d.Get("engine_version")),
		"storage_space":     utils.ValueIngoreEmpty(d.Get("storage_space")),
		"vpc_id":            utils.ValueIngoreEmpty(d.Get("vpc_id")),
		"subnet_id":         utils.ValueIngoreEmpty(d.Get("subnet_id")),
		"security_group_id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
		"available_zones":   availableZones,
		"product_id":        utils.ValueIngoreEmpty(d.Get("flavor_id")),
		"ssl_enable":        utils.ValueIngoreEmpty(d.Get("ssl_enable")),
		"storage_spec_code": utils.ValueIngoreEmpty(d.Get("storage_spec_code")),
		"ipv6_enable":       utils.ValueIngoreEmpty(d.Get("ipv6_enable")),
		"enable_publicip":   utils.ValueIngoreEmpty(d.Get("enable_publicip")),
		"publicip_id":       utils.ValueIngoreEmpty(d.Get("publicip_id")),
		"broker_num":        utils.ValueIngoreEmpty(d.Get("broker_num")),
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	updateRocketmqInstancehasChanges := []string{
		"name",
		"description",
		"security_group_id",
		"retention_policy",
		"cross_vpc_accesses",
	}

	if d.HasChanges(updateRocketmqInstancehasChanges...) {
		// updateRocketmqInstance: update DMS rocketmq instance
		var (
			updateRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
			updateRocketmqInstanceProduct = "dms"
		)
		updateRocketmqInstanceClient, err := config.NewServiceClient(updateRocketmqInstanceProduct, region)
		if err != nil {
			return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
		}

		updateRocketmqInstancePath := updateRocketmqInstanceClient.Endpoint + updateRocketmqInstanceHttpUrl
		updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{project_id}", updateRocketmqInstanceClient.ProjectID)
		updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

		updateRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildUpdateRocketmqInstanceBodyParams(d, config))
		_, err = updateRocketmqInstanceClient.Request("PUT", updateRocketmqInstancePath, &updateRocketmqInstanceOpt)
		if err != nil {
			return diag.Errorf("error updating DmsRocketMQInstance: %s", err)
		}
		if d.HasChange("cross_vpc_accesses") {
			if err = updateCrossVpcAccesses(updateRocketmqInstanceClient, d); err != nil {
				return diag.Errorf("error updating DMS rocketMQ Cross-VPC access information: %s", err)
			}
		}
	}
	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func buildUpdateRocketmqInstanceBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":       utils.ValueIngoreEmpty(d.Get("description")),
		"security_group_id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
		"retention_policy":  utils.ValueIngoreEmpty(d.Get("retention_policy")),
	}
	if d.HasChange("name") {
		bodyParams["name"] = utils.ValueIngoreEmpty(d.Get("name"))
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqInstance: Query DMS rocketmq instance
	var (
		getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstanceProduct = "dms"
	)
	getRocketmqInstanceClient, err := config.NewServiceClient(getRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	getRocketmqInstancePath := getRocketmqInstanceClient.Endpoint + getRocketmqInstanceHttpUrl
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", getRocketmqInstanceClient.ProjectID)
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	getRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqInstanceResp, err := getRocketmqInstanceClient.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQInstance")
	}

	getRocketmqInstanceRespBody, err := utils.FlattenResponse(getRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// convert the ids of the availability zone into codes
	var availableZoneCodes []string
	availableZoneIDs := utils.PathSearch("available_zones", getRocketmqInstanceRespBody, nil)
	if availableZoneIDs != nil {
		azIDs := make([]string, 0)
		for _, v := range availableZoneIDs.([]interface{}) {
			azIDs = append(azIDs, v.(string))
		}
		availableZoneCodes, err = getAvailableZoneCodeByID(config, region, azIDs)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	crossVpcInfo := utils.PathSearch("cross_vpc_info", getRocketmqInstanceRespBody, nil)
	var crossVpcAccess []map[string]interface{}
	if crossVpcInfo != nil {
		crossVpcAccess, err = flattenConnectPorts(crossVpcInfo.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRocketmqInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRocketmqInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRocketmqInstanceRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRocketmqInstanceRespBody, nil)),
		d.Set("specification", utils.PathSearch("specification", getRocketmqInstanceRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getRocketmqInstanceRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getRocketmqInstanceRespBody, nil)),
		d.Set("flavor_id", utils.PathSearch("product_id", getRocketmqInstanceRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", getRocketmqInstanceRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getRocketmqInstanceRespBody, nil)),
		d.Set("availability_zones", availableZoneCodes),
		d.Set("maintain_begin", utils.PathSearch("maintain_begin", getRocketmqInstanceRespBody, nil)),
		d.Set("maintain_end", utils.PathSearch("maintain_end", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_space", utils.PathSearch("total_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("used_storage_space", utils.PathSearch("used_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_publicip", utils.PathSearch("enable_publicip", getRocketmqInstanceRespBody, nil)),
		d.Set("publicip_id", utils.PathSearch("publicip_id", getRocketmqInstanceRespBody, nil)),
		d.Set("publicip_address", utils.PathSearch("publicip_address", getRocketmqInstanceRespBody, nil)),
		d.Set("ssl_enable", utils.PathSearch("ssl_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_spec_code", utils.PathSearch("storage_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("ipv6_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("node_num", utils.PathSearch("node_num", getRocketmqInstanceRespBody, nil)),
		d.Set("new_spec_billing_enable", utils.PathSearch("new_spec_billing_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_acl", utils.PathSearch("enable_acl", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_num", utils.PathSearch("broker_num", getRocketmqInstanceRespBody, nil)),
		d.Set("namesrv_address", utils.PathSearch("namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_address", utils.PathSearch("broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_namesrv_address", utils.PathSearch("public_namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_broker_address", utils.PathSearch("public_broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("resource_spec_code", utils.PathSearch("resource_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("cross_vpc_accesses", crossVpcAccess),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// deleteRocketmqInstance: Delete DMS rocketmq instance
	var (
		deleteRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		deleteRocketmqInstanceProduct = "dms"
	)
	deleteRocketmqInstanceClient, err := config.NewServiceClient(deleteRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	deleteRocketmqInstancePath := deleteRocketmqInstanceClient.Endpoint + deleteRocketmqInstanceHttpUrl
	deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{project_id}", deleteRocketmqInstanceClient.ProjectID)
	deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	deleteRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteRocketmqInstanceClient.Request("DELETE", deleteRocketmqInstancePath, &deleteRocketmqInstanceOpt)
	if err != nil {
		return diag.Errorf("error deleting DmsRocketMQInstance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"},
		Target:       []string{"DELETED"},
		Refresh:      DmsRocketmqInstanceStateRefreshFunc(deleteRocketmqInstanceClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func DmsRocketmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRocketmqInstancePath := client.Endpoint + "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", client.ProjectID)
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceID))
		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		v, err := client.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)
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
		status := utils.PathSearch("status", respBody, "").(string)
		return respBody, status, nil
	}
}
