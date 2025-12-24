package bms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}
// @API BMS PUT /v1/{project_id}/baremetalservers/{server_id}
// @API BMS POST /v1/{project_id}/baremetalservers
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/changeos
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/nics
// @API BMS POST /v1/{project_id}/baremetalservers/action
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/nics/delete
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/tags/action
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/metadata
// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}/tags
// @API BMS GET /v1/{project_id}/jobs/{job_id}
// @API VPC GET /v2.0/ports/{port_id}
// @API BSS GET /V2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceBmsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmsInstanceCreate,
		ReadContext:   resourceBmsInstanceRead,
		UpdateContext: resourceBmsInstanceUpdate,
		DeleteContext: resourceBmsInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// just stash the hash for state & diff comparisons
				StateFunc:        utils.HashAndHexEncode,
				DiffSuppressFunc: utils.SuppressUserData,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"eip_charge_mode", "sharetype", "bandwidth_size",
				},
			},
			"eip_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "sharetype", "bandwidth_size",
				},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PER", "WHOLE",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_size",
				},
			},
			"bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "sharetype",
				},
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"traffic", "bandwidth",
				}, true),
				ConflictsWith: []string{"eip_id"},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"system_disk_size",
				},
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"system_disk_type",
				},
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 59,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"power_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ON", "OFF", "REBOOT",
				}, false),
			},
			"charging_mode": common.SchemaChargingMode([]string{}),
			"period_unit":   common.SchemaPeriodUnit([]string{}),
			"period":        common.SchemaPeriod([]string{}),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// To avoid triggering changes metadata is not backfilled during read.
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBmsInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl    = "v1/{project_id}/baremetalservers"
		product    = "bms"
		bssProduct = "bssv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}
	bssClient, err := cfg.NewServiceClient(bssProduct, region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating BMS instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error creating BMS instance: order_id is not found in API response")
	}
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	if _, ok := d.GetOk("metadata"); ok {
		err = updateInstanceMetadata(d, client, buildUpdateInstanceMetadataBodyParams(d))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if action, ok := d.GetOk("power_action"); ok && action == "OFF" {
		err = updatePowerAction(ctx, d, client, action.(string), schema.TimeoutCreate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBmsInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              d.Get("name"),
		"imageRef":          d.Get("image_id"),
		"flavorRef":         d.Get("flavor_id"),
		"metadata":          buildCreateInstanceMetadataBodyParams(d),
		"user_data":         utils.ValueIgnoreEmpty(buildCreateInstanceUserDataBodyParams(d)),
		"adminPass":         utils.ValueIgnoreEmpty(d.Get("admin_pass")),
		"key_name":          utils.ValueIgnoreEmpty(d.Get("key_pair")),
		"vpcid":             d.Get("vpc_id"),
		"security_groups":   buildCreateInstanceSecurityGroupsBodyParams(d),
		"availability_zone": d.Get("availability_zone"),
		"nics":              buildCreateInstanceNicsBodyParams(d),
		"data_volumes":      utils.ValueIgnoreEmpty(buildCreateInstanceDataVolumesBodyParams(d)),
		"extendparam":       buildCreateInstanceExtendParamBodyParams(d, cfg),
		"publicip":          buildCreateInstancePublicIpBodyParams(d),
		"root_volume":       buildCreateInstanceRootVolumeBodyParams(d),
		"server_tags":       utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
	return map[string]interface{}{
		"server": bodyParams,
	}
}

func buildCreateInstanceMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"op_svc_userid": d.Get("user_id"),
		"agency_name":   utils.ValueIgnoreEmpty(d.Get("agency_name")),
	}
	return bodyParams
}

func buildCreateInstanceUserDataBodyParams(d *schema.ResourceData) string {
	userData := d.Get("user_data").(string)
	if userData == "" {
		return ""
	}
	if _, err := base64.StdEncoding.DecodeString(userData); err != nil {
		userData = base64.StdEncoding.EncodeToString([]byte(userData))
	}
	return userData
}

func buildCreateInstanceSecurityGroupsBodyParams(d *schema.ResourceData) []interface{} {
	rawSecGroups := d.Get("security_groups").(*schema.Set)
	if rawSecGroups.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, rawSecGroups.Len())
	for _, v := range rawSecGroups.List() {
		rst = append(rst, map[string]interface{}{
			"id": v,
		})
	}
	return rst
}

func buildCreateInstanceNicsBodyParams(d *schema.ResourceData) []interface{} {
	rawNics := d.Get("nics").([]interface{})
	if len(rawNics) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawNics))
	for _, v := range rawNics {
		if nicRaw, ok := v.(map[string]interface{}); ok {
			nic := map[string]interface{}{
				"subnet_id": nicRaw["subnet_id"],
			}
			if len(nicRaw["ip_address"].(string)) > 0 {
				nic["ip_address"] = nicRaw["ip_address"]
			}
			rst = append(rst, nic)
		}
	}
	return rst
}

func buildCreateInstanceDataVolumesBodyParams(d *schema.ResourceData) []interface{} {
	rawDataDisks := d.Get("data_disks").([]interface{})
	if len(rawDataDisks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawDataDisks))
	for _, v := range rawDataDisks {
		if dataDisk, ok := v.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"volumetype": dataDisk["type"],
				"size":       dataDisk["size"],
			})
		}
	}
	return rst
}

func buildCreateInstanceExtendParamBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chargingMode":          d.Get("charging_mode"),
		"periodType":            d.Get("period_unit"),
		"periodNum":             d.Get("period"),
		"isAutoPay":             "true",
		"isAutoRenew":           utils.ValueIgnoreEmpty(d.Get("auto_renew")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	return bodyParams
}

func buildCreateInstancePublicIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("eip_id"); ok {
		bodyParams := map[string]interface{}{
			"id": v,
		}
		return bodyParams
	}
	if v, ok := d.GetOk("iptype"); ok {
		bodyParams := map[string]interface{}{
			"iptype":      v,
			"bandwidth":   buildCreateInstancePublicIpBandwidthBodyParams(d),
			"extendparam": buildCreateInstancePublicIpExtendParamBodyParams(d),
		}
		return map[string]interface{}{
			"eip": bodyParams,
		}
	}

	return nil
}

func buildCreateInstancePublicIpBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sharetype":  d.Get("sharetype"),
		"size":       d.Get("bandwidth_size"),
		"chargemode": utils.ValueIgnoreEmpty(d.Get("bandwidth_charge_mode")),
	}
	return bodyParams
}

func buildCreateInstancePublicIpExtendParamBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chargingMode": d.Get("eip_charge_mode"),
	}
	return bodyParams
}

func buildCreateInstanceRootVolumeBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("system_disk_type")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"volumetype": v,
		"size":       d.Get("system_disk_size"),
	}
	return bodyParams
}

func resourceBmsInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	getRespBody, err := getInstance(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving BMS instance")
	}
	status := utils.PathSearch("server.status", getRespBody, "").(string)
	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving BMS instance")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("server.name", getRespBody, nil)),
		d.Set("image_id", utils.PathSearch(`server.metadata."metering.image_id"`, getRespBody, nil)),
		d.Set("flavor_id", utils.PathSearch("server.flavor.id", getRespBody, nil)),
		d.Set("host_id", utils.PathSearch("server.hostId", getRespBody, nil)),
		d.Set("nics", flattenInstanceNics(cfg, region, getRespBody)),
		d.Set("key_pair", utils.PathSearch("server.key_name", getRespBody, nil)),
		d.Set("security_groups", flattenInstanceSecurityGroups(getRespBody)),
		d.Set("status", status),
		d.Set("user_id", utils.PathSearch("server.metadata.op_svc_userid", getRespBody, nil)),
		d.Set("image_name", utils.PathSearch("server.metadata.image_name", getRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("server.metadata.vpc_id", getRespBody, nil)),
		d.Set("agency_name", utils.PathSearch("server.metadata.agency_name", getRespBody, nil)),
		d.Set("availability_zone", utils.PathSearch(`server."OS-EXT-AZ:availability_zone"`, getRespBody, nil)),
		d.Set("description", utils.PathSearch("server.description", getRespBody, nil)),
		d.Set("user_data", utils.PathSearch(`server."OS-EXT-SRV-ATTR:user_data"`, getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("server.enterprise_project_id", getRespBody, nil)),
		d.Set("disk_ids", flattenInstanceDiskIds(getRespBody)),
		d.Set("public_ip", flattenInstancePublicIp(getRespBody)),
	)

	if resourceTags, err := tags.Get(client, "baremetalservers", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] error fetching tags of BMS instance (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstance(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenInstanceNics(cfg *config.Config, region string, respBody interface{}) []interface{} {
	curJson := utils.PathSearch("server.addresses", respBody, nil)
	if curJson == nil {
		return nil
	}

	networkingClient, err := cfg.NewServiceClient("networkv2", region)
	if err != nil {
		log.Printf("error creating networking client: %s", err)
	}
	rst := make([]interface{}, 0)
	for _, addrs := range curJson.(map[string]interface{}) {
		for _, addr := range addrs.([]interface{}) {
			addType := utils.PathSearch(`"OS-EXT-IPS:type"`, addr, "").(string)
			if addType != "fixed" {
				continue
			}
			portId := utils.PathSearch(`"OS-EXT-IPS:port_id"`, addr, "").(string)
			nic := map[string]interface{}{
				"ip_address":  utils.PathSearch("addr", addr, nil),
				"mac_address": utils.PathSearch(`"OS-EXT-IPS-MAC:mac_addr"`, addr, nil),
				"port_id":     portId,
			}
			port, err := getNicPort(networkingClient, portId)
			if err != nil {
				log.Printf("[WARN] error retrieving instance nics: failed to fetch port %s", portId)
			} else {
				nic["subnet_id"] = utils.PathSearch("port.network_id", port, nil)
			}
			rst = append(rst, nic)
		}
	}
	return rst
}

func getNicPort(client *golangsdk.ServiceClient, portId string) (interface{}, error) {
	var (
		httpUrl = "v2.0/ports/{port_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{port_id}", portId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving port: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func flattenInstanceSecurityGroups(respBody interface{}) []interface{} {
	curJson := utils.PathSearch("server.security_groups", respBody, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, utils.PathSearch("id", v, nil))
	}
	return rst
}

func flattenInstanceDiskIds(respBody interface{}) []interface{} {
	curJson := utils.PathSearch(`server."os-extended-volumes:volumes_attached"`, respBody, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, utils.PathSearch("id", v, nil))
	}
	return rst
}

func flattenInstancePublicIp(respBody interface{}) interface{} {
	curJson := utils.PathSearch("server.addresses", respBody, nil)
	if curJson == nil {
		return nil
	}

	for _, addrs := range curJson.(map[string]interface{}) {
		for _, addr := range addrs.([]interface{}) {
			addType := utils.PathSearch(`"OS-EXT-IPS:type"`, addr, "").(string)
			if addType == "floating" {
				return utils.PathSearch("addr", addr, nil)
			}
		}
	}
	return nil
}

func resourceBmsInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	instanceId := d.Id()

	err = checkImageIdUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// if power_action is changed from OFF to ON, the instance should be start first
	powerAction := d.Get("power_action").(string)
	if d.HasChanges("power_action") && powerAction == "ON" {
		if err = updatePowerAction(ctx, d, client, powerAction, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("image_id") {
		err = updateInstanceImage(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("name") {
		err = updateInstanceName(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("agency_name") {
		err = updateInstanceMetadata(d, client, buildUpdateInstanceAgencyNameBodyParams(d))
		if err != nil {
			return diag.Errorf("error updating BMS instance agency name: %s", err)
		}
	}

	if d.HasChanges("metadata") {
		err = updateInstanceMetadata(d, client, buildUpdateInstanceMetadataBodyParams(d))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "baremetalservers", instanceId)
		if err != nil {
			return diag.Errorf("error updating tags of bms server: %s", err)
		}
	}

	// Security group parammeters are missing in the network card, this is a legacy feature.
	// The first network card can not be deleted, api error message:{"error": {"message": "primary port can not be deleted.", "code":"BMS.0222"}}
	if d.HasChange("nics") {
		err = updateInstanceNics(ctx, d, client, schema.TimeoutUpdate)
		if err != nil {
			return diag.Errorf("error updating NICs: %s", err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "bms_server",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	// if power_action is changed from ON to OFF/REBOOT, the instance should be close/restart at the end
	if d.HasChanges("power_action") && (powerAction == "OFF" || powerAction == "REBOOT") {
		if err = updatePowerAction(ctx, d, client, powerAction, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBmsInstanceRead(ctx, d, meta)
}

func checkImageIdUpdate(d *schema.ResourceData) error {
	if d.HasChange("image_id") {
		return nil
	}
	for _, v := range []string{"admin_pass", "key_pair", "user_id", "user_data"} {
		if d.HasChange(v) {
			return fmt.Errorf("%s can only be modified when image_id is modified", v)
		}
	}
	return nil
}

func updateInstanceImage(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	powerAction := d.Get("power_action").(string)
	// Only a stopped BMS or a BMS on which changing the OS failed supports changing OS
	if powerAction != "OFF" {
		if err := updatePowerAction(ctx, d, client, "OFF", schema.TimeoutUpdate); err != nil {
			return err
		}
	}
	err := updateInstanceImageId(ctx, d, client)
	if err != nil {
		return err
	}
	// the instance will be started when changing OS, so it should be stop if power_action is OFF
	if powerAction == "OFF" {
		if err = updatePowerAction(ctx, d, client, "OFF", schema.TimeoutUpdate); err != nil {
			return err
		}
	}
	return nil
}

func updateInstanceImageId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/changeos"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceImageIdBodyParams(d))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating BMS instance image ID: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("jobId", updateRespBody, "").(string)
	if jobId == "" {
		return errors.New("error updating BMS instance image ID: jobId is not found in API response")
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
}

func buildUpdateInstanceImageIdBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"imageid":   d.Get("image_id"),
		"adminpass": utils.ValueIgnoreEmpty(d.Get("admin_pass")),
		"keyname":   utils.ValueIgnoreEmpty(d.Get("key_pair")),
		"userid":    utils.ValueIgnoreEmpty(d.Get("user_id")),
		"metadata":  buildUpdateInstanceImageIdMetadataBodyParams(d),
	}

	return map[string]interface{}{
		"os-change": bodyParams,
	}
}

func buildUpdateInstanceImageIdMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("user_data")
	if !ok {
		return nil
	}
	userData := v.(string)
	if _, err := base64.StdEncoding.DecodeString(userData); err != nil {
		userData = base64.StdEncoding.EncodeToString([]byte(userData))
	}
	bodyParams := map[string]interface{}{
		"user_data": userData,
	}

	return bodyParams
}

func updateInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	addNics, deleteNics := getDiffNics(d)
	if len(deleteNics) > 0 {
		err := deleteInstanceNics(ctx, d, client, timeout, deleteNics)
		if err != nil {
			return err
		}
	}
	if len(addNics) > 0 {
		err := addInstanceNics(ctx, d, client, timeout, addNics)
		if err != nil {
			return err
		}
	}
	return nil
}

func addInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	addNics []interface{}) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/nics"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildAddInstanceNicsBodyParams(d, addNics))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error adding nics to BMS instance: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return errors.New("error adding nics to BMS instance: job_id is not found in API response")
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(timeout))
}

func buildAddInstanceNicsBodyParams(d *schema.ResourceData, addNics []interface{}) map[string]interface{} {
	bodyParams := make([]map[string]interface{}, 0, len(addNics))
	securityGroups := buildAddInstanceNicSecGroupsBodyParams(d)
	for _, nic := range addNics {
		if v, ok := nic.(map[string]interface{}); ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"subnet_id":       v["subnet_id"].(string),
				"ip_address":      utils.ValueIgnoreEmpty(v["ip_address"]),
				"security_groups": utils.ValueIgnoreEmpty(securityGroups),
			})
		}
	}
	return map[string]interface{}{
		"nics": bodyParams,
	}
}

func buildAddInstanceNicSecGroupsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	segGroups := d.Get("security_groups").(*schema.Set)
	rst := make([]map[string]interface{}, 0, segGroups.Len())
	for _, segGroup := range segGroups.List() {
		rst = append(rst, map[string]interface{}{
			"id": segGroup,
		})
	}
	return rst
}

func deleteInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	deleteNics []interface{}) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/nics/delete"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildDeleteInstanceNicsBodyParams(deleteNics)

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error deleting nics to BMS instance: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return errors.New("error deleting nics from BMS instance: job_id is not found in API response")
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(timeout))
}

func buildDeleteInstanceNicsBodyParams(addNics []interface{}) interface{} {
	bodyParams := make([]interface{}, 0, len(addNics))
	for _, nic := range addNics {
		if v, ok := nic.(map[string]interface{}); ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"id": v["port_id"].(string),
			})
		}
	}
	return map[string]interface{}{
		"nics": bodyParams,
	}
}

func updatePowerAction(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, powerAction, timeout string) error {
	var bodyParams interface{}
	var action string
	switch powerAction {
	case "ON":
		bodyParams = buildStartupInstanceBodyParams(d)
		action = "starting"
	case "OFF":
		bodyParams = buildShutdownInstanceBodyParams(d)
		action = "stopping"
	case "REBOOT":
		bodyParams = buildRebootInstanceBodyParams(d)
		action = "rebooting"
	default:
		return fmt.Errorf("the value of power_action(%s) is error, it should be in [ON, OFF, BEBOOT]", powerAction)
	}

	var (
		httpUrl = "v1/{project_id}/baremetalservers/action"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = bodyParams

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error %s BMS instance: %s", action, err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error %s BMS instance: job_id is not found in API response", action)
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(timeout))
}

func buildStartupInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"servers": []map[string]interface{}{
			{
				"id": d.Id(),
			},
		},
	}
	return map[string]interface{}{
		"os-start": bodyParams,
	}
}

func buildShutdownInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": "HARD",
		"servers": []map[string]interface{}{
			{
				"id": d.Id(),
			},
		},
	}
	return map[string]interface{}{
		"os-stop": bodyParams,
	}
}

func buildRebootInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": "HARD",
		"servers": []map[string]interface{}{
			{
				"id": d.Id(),
			},
		},
	}
	return map[string]interface{}{
		"reboot": bodyParams,
	}
}

func updateInstanceName(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceNameBodyParams(d)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating BMS instance name: %s", err)
	}

	return nil
}

func buildUpdateInstanceNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}
	return map[string]interface{}{
		"server": bodyParams,
	}
}

func updateInstanceMetadata(d *schema.ResourceData, client *golangsdk.ServiceClient, bodyParams interface{}) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/metadata"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = bodyParams

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateInstanceAgencyNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agency_name": d.Get("agency_name"),
	}
	return map[string]interface{}{
		"metadata": bodyParams,
	}
}

func buildUpdateInstanceMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": d.Get("metadata"),
	}
	return bodyParams
}

// Get the list of new and to-be-deleted network cards.
func getDiffNics(d *schema.ResourceData) (addList []interface{}, removeList []interface{}) {
	oldNics, newNics := d.GetChange("nics")
	oldList := oldNics.([]interface{})
	newList := newNics.([]interface{})
	for _, ov := range oldList {
		om := ov.(map[string]interface{})
		oSubnetId := om["subnet_id"].(string)
		oIpAddress := om["ip_address"].(string)
		needRemove := true
		for _, nv := range newList {
			nm := nv.(map[string]interface{})
			nSubnetId := nm["subnet_id"].(string)
			nIpAddress := nm["ip_address"].(string)

			if oSubnetId == nSubnetId && oIpAddress == nIpAddress {
				needRemove = false
				break
			}
		}
		if needRemove {
			removeList = append(removeList, ov)
		}
	}

	for _, nv := range newList {
		nm := nv.(map[string]interface{})
		nSubnetId := nm["subnet_id"].(string)
		nIpAddress := nm["ip_address"].(string)
		needAdd := true
		for _, ov := range oldList {
			om := ov.(map[string]interface{})
			oSubnetId := om["subnet_id"].(string)
			oIpAddress := om["ip_address"].(string)
			if nSubnetId == oSubnetId && nIpAddress == oIpAddress {
				needAdd = false
				break
			}
		}
		if needAdd {
			addList = append(addList, nv)
		}
	}
	return
}

func resourceBmsInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product    = "bms"
		eipProduct = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	serverID := d.Id()
	publicIP := d.Get("public_ip").(string)
	diskIds := d.Get("disk_ids").([]interface{})

	resourceIDs := make([]string, 0, 2+len(diskIds))

	resourceIDs = append(resourceIDs, serverID)

	if len(diskIds) > 0 {
		for _, diskId := range diskIds {
			resourceIDs = append(resourceIDs, diskId.(string))
		}
	}

	// unsubscribe the eip if necessary
	if _, ok := d.GetOk("iptype"); ok && publicIP != "" && d.Get("eip_charge_mode").(string) == "prePaid" {
		eipClient, err := cfg.NewServiceClient(eipProduct, region)
		if err != nil {
			return diag.Errorf("error creating networking client: %s", err)
		}

		eipId, err := getEipIdByAddress(eipClient, publicIP)
		if err != nil {
			return diag.FromErr(err)
		}

		resourceIDs = append(resourceIDs, eipId)
	}

	if err := common.UnsubscribePrePaidResource(d, cfg, resourceIDs); err != nil {
		return diag.Errorf("error unsubscribing BMS server: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "ACTIVE", "SHUTOFF"},
		Target:       []string{"DELETED"},
		Refresh:      waitForBmsInstanceDelete(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting BMS instance: %s", err)
	}

	return nil
}

func getEipIdByAddress(client *golangsdk.ServiceClient, address string) (string, error) {
	var (
		httpUrl = "v3/{project_id}/eip/publicips"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetEipByAddressQueryParams(address)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving EIP by address(%s): %s", address, err)
	}
	getRestBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return "", err
	}

	eipId := utils.PathSearch("publicips[0].id", getRestBody, nil)
	if eipId == nil {
		return "", fmt.Errorf("error retrieving EIP by address(%s)", address)
	}

	return eipId.(string), nil
}

func buildGetEipByAddressQueryParams(address string) string {
	return fmt.Sprintf("?enterprise_project_id=all_granted_eps&public_ip_address=%s", address)
}

func waitForBmsInstanceDelete(client *golangsdk.ServiceClient, serverId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getInstance(client, serverId)
		if err != nil {
			return getRespBody, "Deleting", err
		}
		status := utils.PathSearch("server.status", getRespBody, "").(string)

		return getRespBody, status, nil
	}
}
