package geminidb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

var geminiDbInstanceNonUpdatableParams = []string{"datastore", "datastore.*.type", "datastore.*.storage_engine",
	"datastore.*.version", "availability_zone", "vpc_id", "subnet_id", "mode", "flavor.*.storage", "product_type",
	"product_type", "dedicated_resource_id", "availability_zone_detail", "availability_zone_detail.*.primary_availability_zone",
	"availability_zone_detail.*.secondary_availability_zone", "charging_mode", "period_unit", "period",
}

// @API GaussDBforNoSQL POST /v3/{project_id}/instances
// @API GaussDBforNoSQL GET /v3/{project_id}/jobs
// @API GaussDBforNoSQL GET /v3/{project_id}/instances
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/name
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/password
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/ssl-option
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/port
// @API GaussDBforNoSQL PUT /v3.1/{project_id}/configurations/{config_id}/apply
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/volume
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/resize
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/enlarge-node
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/reduce-node
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/security-group
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/backups/policy
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/tags/action
// @API GaussDBforNoSQL DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
func ResourceGeminiDbInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDbInstanceCreate,
		ReadContext:   resourceGeminiDbInstanceRead,
		UpdateContext: resourceGeminiDbInstanceUpdate,
		DeleteContext: resourceGeminiDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(geminiDbInstanceNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     geminiDbInstanceDatastoreSchema(),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     geminiDbInstanceFlavorSchema(),
			},
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     geminiDbInstanceBackupStrategySchema(),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_option": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"availability_zone_detail": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     geminiDbInstanceAvailabilityZoneDetailSchema(),
			},
			"delete_node_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// charge info: charging_mode, period_unit, period, auto_renew
			// make ForceNew false here but do nothing in update method!
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"tags":       common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDbInstanceGroupsSchema(),
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lb_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dual_active_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDbInstanceDualActiveInfoSchema(),
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func geminiDbInstanceDatastoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"patch_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"whole_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceBackupStrategySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keep_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceAvailabilityZoneDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"primary_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secondary_availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceGroupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDbInstanceGroupVolumeSchema(),
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDbInstanceGroupNodesSchema(),
			},
		},
	}
	return &sc
}

func geminiDbInstanceGroupVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceGroupNodesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_reduce": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func geminiDbInstanceDualActiveInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceGeminiDbInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances"
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGeminiDbInstanceBodyParams(d, region, cfg.GetEnterpriseProjectID(d)))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GeminiDB instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating GeminiDB instance: ID is not found in API response")
	}
	d.SetId(id)

	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		orderId := utils.PathSearch("order_id", createRespBody, nil)
		if orderId == nil {
			return diag.Errorf("error creating GeminiDB instance: order_id is not found in API response")
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		jobId := utils.PathSearch("job_id", createRespBody, nil)
		if jobId == nil {
			return diag.Errorf("error creating GeminiDB instance: job_id is not found in API response")
		}
		err = checkGeminiDbInstanceJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = utils.CreateResourceTags(client, d, "instances", d.Id())
	if err != nil {
		return diag.Errorf("error creating GeminiDB instance(%s) tags: %s", d.Id(), err)
	}

	// This is a workaround to avoid db connection issue
	// lintignore:R018
	time.Sleep(360 * time.Second)

	return resourceGeminiDbInstanceRead(ctx, d, meta)
}

func buildCreateGeminiDbInstanceBodyParams(d *schema.ResourceData, region, epsId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                     d.Get("name"),
		"datastore":                buildCreateGeminiDbInstanceDatastoreBody(d),
		"region":                   region,
		"availability_zone":        d.Get("availability_zone"),
		"vpc_id":                   d.Get("vpc_id"),
		"subnet_id":                d.Get("subnet_id"),
		"security_group_id":        d.Get("security_group_id"),
		"password":                 d.Get("password"),
		"mode":                     d.Get("mode"),
		"flavor":                   buildCreateGeminiDbInstanceFlavorBody(d),
		"product_type":             utils.ValueIgnoreEmpty(d.Get("product_type")),
		"configuration_id":         utils.ValueIgnoreEmpty(d.Get("configuration_id")),
		"backup_strategy":          buildCreateGeminiDbInstanceBackupStrategyBody(d),
		"enterprise_project_id":    epsId,
		"ssl_option":               utils.ValueIgnoreEmpty(d.Get("ssl_option")),
		"dedicated_resource_id":    utils.ValueIgnoreEmpty(d.Get("dedicated_resource_id")),
		"port":                     utils.ValueIgnoreEmpty(d.Get("port")),
		"availability_zone_detail": buildCreateGeminiDbInstanceAvailabilityZoneDetailBody(d),
		"charge_info":              buildCreateGeminiDbInstanceChargeInfoBody(d),
	}
	if port, ok := d.GetOk("port"); ok {
		bodyParams["port"] = strconv.Itoa(port.(int))
	}
	if v, ok := d.GetOk("ssl_option"); ok {
		if v == "on" {
			bodyParams["ssl_option"] = "1"
		} else if v == "off" {
			bodyParams["ssl_option"] = "0"
		}
	}
	return bodyParams
}

func buildCreateGeminiDbInstanceDatastoreBody(d *schema.ResourceData) map[string]interface{} {
	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 0 {
		return nil
	}
	datastore, ok := datastoreRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"type":           datastore["type"],
		"storage_engine": datastore["storage_engine"],
		"version":        datastore["version"],
	}
	return rst
}

func buildCreateGeminiDbInstanceFlavorBody(d *schema.ResourceData) []map[string]interface{} {
	flavorRaw := d.Get("flavor").([]interface{})
	flavor := flavorRaw[0].(map[string]interface{})
	rst := []map[string]interface{}{
		{
			"num":       strconv.Itoa(flavor["num"].(int)),
			"size":      strconv.Itoa(flavor["size"].(int)),
			"storage":   flavor["storage"],
			"spec_code": flavor["spec_code"],
		},
	}
	return rst
}

func buildCreateGeminiDbInstanceBackupStrategyBody(d *schema.ResourceData) map[string]interface{} {
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 0 {
		return nil
	}
	backupStrategy, ok := backupStrategyRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"start_time": backupStrategy["start_time"],
		"keep_days":  utils.ValueIgnoreEmpty(backupStrategy["keep_days"]),
	}
	return rst
}

func buildCreateGeminiDbInstanceAvailabilityZoneDetailBody(d *schema.ResourceData) map[string]interface{} {
	availabilityZoneDetailRaw := d.Get("availability_zone_detail").([]interface{})
	if len(availabilityZoneDetailRaw) == 0 {
		return nil
	}
	availabilityZoneDetail, ok := availabilityZoneDetailRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rst := map[string]interface{}{
		"primary_availability_zone":   availabilityZoneDetail["primary_availability_zone"],
		"secondary_availability_zone": availabilityZoneDetail["secondary_availability_zone"],
	}
	return rst
}

func buildCreateGeminiDbInstanceChargeInfoBody(d *schema.ResourceData) map[string]interface{} {
	chargingMode := d.Get("charging_mode").(string)
	if chargingMode != "prePaid" {
		return nil
	}

	rst := map[string]interface{}{
		"charge_mode":   d.Get("charging_mode"),
		"period_type":   d.Get("period_unit"),
		"period_num":    d.Get("period"),
		"is_auto_renew": d.Get("auto_renew"),
		"is_auto_pay":   "true",
	}
	return rst
}

func resourceGeminiDbInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}
	instance, err := getGeminiDbInstance(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB instance")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("datastore", flattenGeminiDbInstanceResponseBodyDatastore(instance)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("mode", utils.PathSearch("mode", instance, nil)),
		d.Set("flavor", flattenGeminiDbInstanceResponseBodyFlavor(d, instance)),
		d.Set("product_type", utils.PathSearch("product_type", instance, nil)),
		d.Set("backup_strategy", flattenGeminiDbInstanceResponseBodyBackupStrategy(instance)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("dedicated_resource_id", utils.PathSearch("dedicated_resource_id", instance, nil)),
		d.Set("availability_zone_detail", flattenGaussDBProxyResponseBodyAvailabilityZoneDetail(instance)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("db_user_name", utils.PathSearch("db_user_name", instance, nil)),
		d.Set("groups", flattenGaussDBProxyResponseBodyGroups(instance)),
		d.Set("time_zone", utils.PathSearch("time_zone", instance, nil)),
		d.Set("actions", utils.PathSearch("actions", instance, nil)),
		d.Set("lb_ip_address", utils.PathSearch("lb_ip_address", instance, nil)),
		d.Set("dual_active_info", flattenGaussDBProxyResponseBodyDualActiveInfo(instance)),
		d.Set("created", utils.PathSearch("created", instance, nil)),
		d.Set("updated", utils.PathSearch("updated", instance, nil)),
	)

	port, _ := strconv.Atoi(utils.PathSearch("port", instance, "0").(string))
	mErr = multierror.Append(mErr, d.Set("port", port))

	payMode := utils.PathSearch("pay_mode", instance, "").(string)
	if payMode == "1" {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	}

	if err = utils.SetResourceTagsToState(d, client, "instances", d.Id()); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGeminiDbInstanceResponseBodyDatastore(instance interface{}) []interface{} {
	datastore := utils.PathSearch("datastore", instance, nil)
	if datastore == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"type":            utils.PathSearch("type", datastore, nil),
			"storage_engine":  utils.PathSearch("engine", instance, nil),
			"version":         utils.PathSearch("version", datastore, nil),
			"patch_available": utils.PathSearch("patch_available", datastore, nil),
			"whole_version":   utils.PathSearch("whole_version", datastore, nil),
		},
	}
	return rst
}

func flattenGeminiDbInstanceResponseBodyFlavor(d *schema.ResourceData, instance interface{}) []interface{} {
	specCode := utils.PathSearch("groups[0].nodes[0].spec_code", instance, nil)
	nodes := utils.PathSearch("groups[0].nodes", instance, make([]interface{}, 0)).([]interface{})
	sizeRaw := utils.PathSearch("groups[0].volume.size", instance, "0").(string)
	size, _ := strconv.Atoi(sizeRaw)

	flavorRaw := d.Get("flavor").([]interface{})
	var storage string
	if len(flavorRaw) > 0 {
		if flavor, ok := flavorRaw[0].(map[string]interface{}); ok {
			storage = flavor["storage"].(string)
		}
	}

	rst := []interface{}{
		map[string]interface{}{
			"num":       len(nodes),
			"size":      size,
			"spec_code": specCode,
			"storage":   storage,
		},
	}
	return rst
}

func flattenGeminiDbInstanceResponseBodyBackupStrategy(instance interface{}) []interface{} {
	backupStrategy := utils.PathSearch("backup_strategy", instance, nil)
	if backupStrategy == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"start_time": utils.PathSearch("start_time", backupStrategy, nil),
			"keep_days":  utils.PathSearch("keep_days", backupStrategy, nil),
		},
	}
	return rst
}

func flattenGaussDBProxyResponseBodyAvailabilityZoneDetail(instance interface{}) []interface{} {
	availabilityZoneDetail := utils.PathSearch("availability_zone_detail", instance, nil)
	if availabilityZoneDetail == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"primary_availability_zone":   utils.PathSearch("primary_availability_zone", availabilityZoneDetail, nil),
			"secondary_availability_zone": utils.PathSearch("secondary_availability_zone", availabilityZoneDetail, nil),
		},
	}
	return rst
}

func flattenGaussDBProxyResponseBodyGroups(instance interface{}) []interface{} {
	curJson := utils.PathSearch("groups", instance, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("id", v, nil),
			"status": utils.PathSearch("status", v, nil),
			"volume": flattenGaussDBProxyResponseBodyGroupVolume(v),
			"nodes":  flattenGaussDBProxyResponseBodyGroupNodes(v),
		})
	}
	return rst
}

func flattenGaussDBProxyResponseBodyGroupVolume(group interface{}) []interface{} {
	volume := utils.PathSearch("volume", group, nil)
	if volume == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"size": utils.PathSearch("size", volume, nil),
			"used": utils.PathSearch("used", volume, nil),
		},
	}
	return rst
}

func flattenGaussDBProxyResponseBodyGroupNodes(group interface{}) []interface{} {
	curJson := utils.PathSearch("nodes", group, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"subnet_id":         utils.PathSearch("subnet_id", v, nil),
			"private_ip":        utils.PathSearch("private_ip", v, nil),
			"public_ip":         utils.PathSearch("public_ip", v, nil),
			"spec_code":         utils.PathSearch("spec_code", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"support_reduce":    utils.PathSearch("support_reduce", v, nil),
		})
	}
	return rst
}

func flattenGaussDBProxyResponseBodyDualActiveInfo(instance interface{}) []interface{} {
	dualActiveInfo := utils.PathSearch("dual_active_info", instance, nil)
	if dualActiveInfo == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"role":                    utils.PathSearch("role", dualActiveInfo, nil),
			"status":                  utils.PathSearch("status", dualActiveInfo, nil),
			"destination_instance_id": utils.PathSearch("destination_instance_id", dualActiveInfo, nil),
			"destination_region":      utils.PathSearch("destination_region", dualActiveInfo, nil),
		},
	}
	return rst
}

func resourceGeminiDbInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	err = updateGeminiDbInstanceName(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstancePassword(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceSslOption(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstancePort(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceConfigurationId(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceVolumeSize(ctx, d, client, bssClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceFlavor(ctx, d, client, bssClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceNodeNum(ctx, d, client, bssClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceSecurityGroupId(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbInstanceBackupStrategy(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "instances", d.Id())
		if err != nil {
			return diag.Errorf("error updating GeminiDB instance(%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the GeminiDB instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "nosql",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err = cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGeminiDbInstanceRead(ctx, d, meta)
}

func updateGeminiDbInstanceName(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("name") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/name",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateGeminiDbInstanceNameBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) name: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}
	return bodyParams
}

func updateGeminiDbInstancePassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("password") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/password",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateGeminiDbInstancePasswordBodyParams(d),
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) password: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstancePasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": d.Get("password"),
	}
	return bodyParams
}

func updateGeminiDbInstanceSslOption(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("ssl_option") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/ssl-option",
		httpMethod:          "POST",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateGeminiDbInstanceSslOptionBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		checkJobExpression:  "job_id",
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) SSL option: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceSslOptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ssl_option": d.Get("ssl_option"),
	}
	return bodyParams
}

func updateGeminiDbInstancePort(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("port") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/port",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateGeminiDbInstancePortBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		checkJobExpression:  "job_id",
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) port: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstancePortBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"port": d.Get("port"),
	}
	return bodyParams
}

func updateGeminiDbInstanceConfigurationId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("configuration_id") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3.1/{project_id}/configurations/{config_id}/apply",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"config_id": d.Get("configuration_id").(string)},
		updateBodyParams:    buildUpdateGeminiDbInstanceConfigurationIdBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		checkJobExpression:  "job_id",
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) configuration ID: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceConfigurationIdBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids": []string{d.Id()},
	}
	return bodyParams
}

func updateGeminiDbInstanceVolumeSize(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	if !d.HasChange("flavor.0.size") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/volume",
		httpMethod:           "PUT",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     buildUpdateGeminiDbInstanceVolumeSizeBodyParams(d),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) volume size: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceVolumeSizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"size": d.Get("flavor.0.size").(int),
	}
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		bodyParams["is_auto_pay"] = true
	}
	return bodyParams
}

func updateGeminiDbInstanceFlavor(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	if !d.HasChange("flavor.0.spec_code") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/resize",
		httpMethod:           "PUT",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     buildUpdateGeminiDbInstanceFlavorBodyParams(d),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) flavor: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceFlavorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resize": map[string]interface{}{
			"target_spec_code": d.Get("flavor.0.spec_code").(string),
		},
	}
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		bodyParams["is_auto_pay"] = "true"
	}
	return bodyParams
}

func updateGeminiDbInstanceNodeNum(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient) error {
	if !d.HasChanges("flavor.0.num", "delete_node_list") {
		return nil
	}

	oldNodeNumRaw, newNodeNumRaw := d.GetChange("flavor.0.num")
	oldNodeNum := oldNodeNumRaw.(int)
	newNodeNum := newNodeNumRaw.(int)
	deleteNodeList := d.Get("delete_node_list").(*schema.Set)
	if newNodeNum > oldNodeNum {
		if deleteNodeList.Len() > 0 {
			return errors.New("the delete_node_list cannot be set when the number of nodes is increased")
		}
		err := enlargeGeminiDbInstanceNodeNUm(ctx, d, client, bssClient, newNodeNum-oldNodeNum)
		if err != nil {
			return err
		}
	}
	if deleteNodeList.Len() > 0 {
		for _, deleteNode := range deleteNodeList.List() {
			bodyParams := buildReduceGeminiDbInstanceNodeBodyParams(deleteNode.(string))
			err := reduceGeminiDbInstanceNodeNum(ctx, d, client, bssClient, bodyParams)
			if err != nil {
				return err
			}
		}
	} else if oldNodeNum > newNodeNum {
		bodyParams := buildReduceGeminiDbInstanceNodeNumBodyParams()
		for i := 0; i < oldNodeNum-newNodeNum; i++ {
			err := reduceGeminiDbInstanceNodeNum(ctx, d, client, bssClient, bodyParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildReduceGeminiDbInstanceNodeBodyParams(nodeId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_list": []string{nodeId},
	}
	return bodyParams
}

func buildReduceGeminiDbInstanceNodeNumBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"num": 1,
	}
	return bodyParams
}

func enlargeGeminiDbInstanceNodeNUm(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient,
	expandNum int) error {
	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/enlarge-node",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     buildEnlargeGeminiDbInstanceNodeNUmBodyParams(d, expandNum),
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error enlarging GeminiDB instance(%s) node num: %s", d.Id(), err)
	}
	return nil
}

func buildEnlargeGeminiDbInstanceNodeNUmBodyParams(d *schema.ResourceData, expandNum int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"num":       expandNum,
		"subnet_id": d.Get("subnet_id"),
	}
	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		bodyParams["is_auto_pay"] = "true"
	}
	return bodyParams
}

func reduceGeminiDbInstanceNodeNum(ctx context.Context, d *schema.ResourceData, client, bssClient *golangsdk.ServiceClient,
	bodyParams interface{}) error {
	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:              "v3/{project_id}/instances/{instance_id}/reduce-node",
		httpMethod:           "POST",
		pathParams:           map[string]string{"instance_id": d.Id()},
		updateBodyParams:     bodyParams,
		isRetry:              true,
		timeout:              schema.TimeoutUpdate,
		checkJobExpression:   "job_id",
		checkOrderExpression: "order_id",
		bssClient:            bssClient,
		isWaitInstanceReady:  true,
	})
	if err != nil {
		return fmt.Errorf("error reducing GeminiDB instance(%s) node num: %s", d.Id(), err)
	}
	return nil
}

func updateGeminiDbInstanceSecurityGroupId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("security_group_id") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/security-group",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateGeminiDbInstanceSecurityGroupIdBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		checkJobExpression:  "job_id",
		isWaitInstanceReady: true,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) security group ID: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceSecurityGroupIdBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"security_group_id": d.Get("security_group_id"),
	}
	return bodyParams
}

func updateGeminiDbInstanceBackupStrategy(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("backup_strategy") {
		return nil
	}

	_, err := updateGeminiDbInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:          "v3/{project_id}/instances/{instance_id}/backups/policy",
		httpMethod:       "PUT",
		pathParams:       map[string]string{"instance_id": d.Id()},
		updateBodyParams: buildUpdateGeminiDbInstanceBackupStrategyBodyParams(d),
		isRetry:          true,
		timeout:          schema.TimeoutUpdate,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB instance(%s) backup strategy: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateGeminiDbInstanceBackupStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 0 {
		return nil
	}
	backupStrategy, ok := backupStrategyRaw[0].(map[string]interface{})
	if !ok {
		return nil
	}
	bodyParams := map[string]interface{}{
		"backup_policy": map[string]interface{}{
			"keep_days":  backupStrategy["keep_days"],
			"start_time": backupStrategy["start_time"],
			"period":     "1,2,3,4,5,6,7",
		},
	}
	return bodyParams
}

func resourceGeminiDbInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	if v, ok := d.GetOk("charging_mode"); ok && v == "prePaid" {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
			retry, err := handleDeletionError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe GeminiDB instance: %s", err)
		}
	} else {
		err = deleteRdsInstance(ctx, d, client)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error deleting GeminiDB instance")
		}
	}

	return nil
}

func deleteRdsInstance(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleDeletionError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	deleteRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return errors.New("error deleting GeminiDB instance: job_id is not found in the response")
	}
	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	return nil
}

func checkGeminiDbInstanceJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      geminiDbInstanceJobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GeminiDB job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func geminiDbInstanceJobRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v3/{project_id}/jobs?id={job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", jobId)

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return nil, "Failed", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "Failed", err
		}

		status := utils.PathSearch("jobs[0].status", getJobStatusRespBody, "").(string)
		if status == "" {
			return nil, "Failed", errors.New("job is not found")
		}
		if status == "Failed" {
			return getJobStatusRespBody, "Failed", nil
		}
		if status == "Completed" {
			return getJobStatusRespBody, "Completed", nil
		}

		return getJobStatusRespBody, "Pending", nil
	}
}

func geminiDbInstanceStatusRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getGeminiDbInstance(client, instanceID)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", instance, "").(string)
		if utils.StrSliceContains([]string{"abnormal", "createfail", "enlargefail"}, status) {
			return instance, "ERROR", fmt.Errorf("the instance status is: %s", status)
		}
		if status == "normal" {
			return instance, "ACTIVE", nil
		}
		return instance, status, nil
	}
}

func getGeminiDbInstance(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

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
	instance := utils.PathSearch("instances[0]", getRespBody, nil)
	if instance == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return instance, nil
}
