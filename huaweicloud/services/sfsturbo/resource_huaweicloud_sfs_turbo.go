package sfsturbo

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	shareTypeStandard    = "STANDARD"
	shareTypePerformance = "PERFORMANCE"
	shareTypeHpc         = "HPC"
	shareTypeHpcCache    = "HPC_CACHE"
)

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{id}/action
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{id}
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{id}
// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/{id}/tags/{key}
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/{id}/tags
// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/{id}/tags/action
func ResourceSFSTurbo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSTurboCreate,
		ReadContext:   resourceSFSTurboRead,
		UpdateContext: resourceSFSTurboUpdate,
		DeleteContext: resourceSFSTurboDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"share_proto": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "NFS",
				ValidateFunc: validation.StringInSlice([]string{"NFS"}, false),
			},
			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					shareTypeStandard, shareTypePerformance, shareTypeHpc, shareTypeHpcCache,
				}, false),
				Default: shareTypeStandard,
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
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"crypt_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enhanced": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"hpc_bandwidth": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"enhanced", "hpc_cache_bandwidth"},
			},
			"hpc_cache_bandwidth": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enhanced", "hpc_bandwidth"},
			},
			// Editing this field will only take effect when the `security_group_id` field is changed.
			"auto_create_security_group_rules": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"dedicated_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags":          common.TagsSchema(),
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateParameter(d *schema.ResourceData) error {
	_, isHpcBandwidthSet := d.GetOk("hpc_bandwidth")
	_, isHpcCacheBandwidthSet := d.GetOk("hpc_cache_bandwidth")
	switch d.Get("share_type").(string) {
	case shareTypeHpc:
		if !isHpcBandwidthSet {
			return errors.New("`hpc_bandwidth` is required when share type is HPC")
		}
	case shareTypeHpcCache:
		if !isHpcCacheBandwidthSet {
			return errors.New("`hpc_cache_bandwidth` is required when share type is HPC_CACHE")
		}
		if d.Get("charging_mode").(string) == "prePaid" {
			return errors.New("HPC_CACHE share type only support in postpaid charging mode")
		}
	default:
		if isHpcBandwidthSet || isHpcCacheBandwidthSet {
			return errors.New("`hpc_bandwidth` and `hpc_cache_bandwidth` cannot be set when share type is" +
				" STANDARD or PERFORMANCE")
		}
	}
	return nil
}

func buildTurboShareTypeParam(d *schema.ResourceData) string {
	shareType := d.Get("share_type").(string)
	if shareType == shareTypeStandard || shareType == shareTypePerformance {
		return shareType
	}
	// For `HPC` and `HPC_CACHE` types, create API does not validate the share type.
	// We can fill in `STANDARD` or `PERFORMANCE`.
	return shareTypeStandard
}

func buildTurboMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	rstMap := map[string]interface{}{
		"crypt_key_id":                     utils.ValueIgnoreEmpty(d.Get("crypt_key_id")),
		"dedicated_flavor":                 utils.ValueIgnoreEmpty(d.Get("dedicated_flavor")),
		"dedicated_storage_id":             utils.ValueIgnoreEmpty(d.Get("dedicated_storage_id")),
		"auto_create_security_group_rules": utils.ValueIgnoreEmpty(d.Get("auto_create_security_group_rules")),
	}

	switch d.Get("share_type").(string) {
	case shareTypeHpc:
		rstMap["expand_type"] = "hpc"
		rstMap["hpc_bw"] = d.Get("hpc_bandwidth")
	case shareTypeHpcCache:
		rstMap["expand_type"] = "hpc_cache"
		rstMap["hpc_bw"] = d.Get("hpc_cache_bandwidth")
	default:
		if _, ok := d.GetOk("enhanced"); ok {
			rstMap["expand_type"] = "bandwidth"
		}
	}

	return rstMap
}

func buildTurboShareCreateBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                  d.Get("name"),
		"size":                  d.Get("size"),
		"share_proto":           utils.ValueIgnoreEmpty(d.Get("share_proto")),
		"vpc_id":                d.Get("vpc_id"),
		"subnet_id":             d.Get("subnet_id"),
		"security_group_id":     d.Get("security_group_id"),
		"availability_zone":     d.Get("availability_zone"),
		"backup_id":             utils.ValueIgnoreEmpty(d.Get("backup_id")),
		"share_type":            buildTurboShareTypeParam(d),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"metadata":              buildTurboMetadataBodyParams(d),
	}
}

func buildTurboBillingPeriodTypeParam(d *schema.ResourceData) int {
	if d.Get("period_unit").(string) == "month" {
		return 2
	}
	return 3
}

func buildTurboBillingAutoRenewParam(d *schema.ResourceData) int {
	if d.Get("auto_renew").(string) == "true" {
		return 1
	}
	return 0
}

func buildTurboBillingCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	if d.Get("charging_mode") != "prePaid" {
		return nil
	}

	return map[string]interface{}{
		"period_num":    d.Get("period"),
		"is_auto_pay":   1,
		"period_type":   buildTurboBillingPeriodTypeParam(d),
		"is_auto_renew": buildTurboBillingAutoRenewParam(d),
	}
}

func buildSFSTurboCreateBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"share":     buildTurboShareCreateBodyParams(cfg, d),
		"bss_param": buildTurboBillingCreateBodyParams(d),
	}
}

func GetTurboDetail(client *golangsdk.ServiceClient, shareId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", shareId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForTurboStatusReady(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetTurboDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in SFS Turbo detail API response")
			}

			if utils.StrSliceContains([]string{"303", "800"}, status) {
				return respBody, "ERROR", fmt.Errorf("unexpected status: '%s'", status)
			}

			if utils.StrSliceContains([]string{"200"}, status) {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceSFSTurboCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "sfs-turbo"
		httpUrl = "v1/{project_id}/sfs-turbo/shares"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	if err := validateParameter(d); err != nil {
		return diag.FromErr(err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSFSTurboCreateBodyParams(cfg, d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId := utils.PathSearch("orderId", respBody, "").(string)
		if orderId == "" {
			return diag.Errorf(`error creating SFS Turbo: unable to find the order ID,
			 this is a COM (Cloud Order Management) error, please contact service for help and check your order
			 status on the console.`)
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
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
	} else {
		id := utils.PathSearch("id", respBody, "").(string)
		if id == "" {
			return diag.Errorf("error creating SFS Turbo: unable to find the share ID")
		}
		d.SetId(id)
	}

	if err := waitingForTurboStatusReady(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo (%s) to become ready: %s", d.Id(), err)
	}

	// add tags
	if err := utils.CreateResourceTags(client, d, "sfs-turbo", d.Id()); err != nil {
		return diag.Errorf("error setting tags of SFS Turbo %s: %s", d.Id(), err)
	}

	return resourceSFSTurboRead(ctx, d, meta)
}

func flattenSizeAttribute(size string) interface{} {
	// size is a string of float64, should convert it to int
	if fsize, err := strconv.ParseFloat(size, 64); err == nil {
		return int(fsize)
	}

	return nil
}

func flattenStatusAttribute(respBody interface{}) interface{} {
	if subStatus := utils.PathSearch("sub_status", respBody, nil); subStatus != nil {
		return subStatus
	}

	return utils.PathSearch("status", respBody, nil)
}

func resourceSFSTurboRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	respBody, err := GetTurboDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "SFS.TURBO.9000"),
			"error retrieving SFS Turbo")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("share_proto", utils.PathSearch("share_proto", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", respBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", respBody, nil)),
		d.Set("available_capacity", utils.PathSearch("avail_capacity", respBody, nil)),
		d.Set("export_location", utils.PathSearch("export_location", respBody, nil)),
		d.Set("crypt_key_id", utils.PathSearch("crypt_key_id", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("backup_id", utils.PathSearch("backup_id", respBody, nil)),
		d.Set("size", flattenSizeAttribute(utils.PathSearch("size", respBody, "").(string))),
		d.Set("status", flattenStatusAttribute(respBody)),
	)

	expandType := utils.PathSearch("expand_type", respBody, "").(string)
	// `HPC` and `HPC_CACHE` are custom types. `STANDARD` and `PERFORMANCE` are system types.
	switch expandType {
	case "hpc":
		mErr = multierror.Append(
			mErr,
			d.Set("share_type", shareTypeHpc),
			d.Set("hpc_bandwidth", utils.PathSearch("hpc_bw", respBody, nil)),
		)
	case "hpc_cache":
		mErr = multierror.Append(
			mErr,
			d.Set("share_type", shareTypeHpcCache),
			d.Set("hpc_cache_bandwidth", utils.PathSearch("hpc_bw", respBody, nil)),
		)
	default:
		mErr = multierror.Append(mErr, d.Set("share_type", utils.PathSearch("share_type", respBody, nil)))
		if expandType == "bandwidth" {
			mErr = multierror.Append(mErr, d.Set("enhanced", true))
		} else {
			mErr = multierror.Append(mErr, d.Set("enhanced", false))
		}
	}

	// set tags
	err = utils.SetResourceTagsToState(d, client, "sfs-turbo", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func convertHpcCacheBandwidth(d *schema.ResourceData) (int, error) {
	if d.Get("share_type") == shareTypeHpcCache {
		hpcCacheBandwidth := d.Get("hpc_cache_bandwidth").(string)
		if !strings.HasSuffix(hpcCacheBandwidth, "G") {
			return 0, fmt.Errorf("the unit of HPC cache bandwidth is missing or invalid, want 'xG',"+
				" but got '%s'", hpcCacheBandwidth)
		}
		bandwidth, err := strconv.Atoi(strings.TrimRight(hpcCacheBandwidth, "G"))
		if err != nil {
			return 0, fmt.Errorf("failed to convert HPC cache bandwidth (%s) to int value: %s",
				hpcCacheBandwidth, err)
		}
		return bandwidth, nil
	}
	return 0, nil
}

func buildTurboExtendOpts(newSize, hpcCacheBandwidth int, isPrePaid bool) map[string]interface{} {
	extendMap := map[string]interface{}{
		"new_size":      newSize,
		"new_bandwidth": hpcCacheBandwidth,
	}

	if isPrePaid {
		extendMap["bss_param"] = map[string]interface{}{
			"is_auto_pay": 1,
		}
	}

	return map[string]interface{}{
		"extend": extendMap,
	}
}

func waitingForTurboSubStatusReady(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetTurboDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			subStatus := utils.PathSearch("sub_status", respBody, "").(string)
			if subStatus == "" {
				return nil, "ERROR", errors.New("sub_status is not found in SFS Turbo detail API response")
			}

			// '321' indicate expansion failed
			// '332' indicate changing security group failed
			if utils.StrSliceContains([]string{"321", "332"}, subStatus) {
				return respBody, "ERROR", fmt.Errorf("unexpected sub_status: '%s'", subStatus)
			}

			// '221' indicate expansion succeeded
			// '232' indicate changing security group succeeded
			if utils.StrSliceContains([]string{"221", "232"}, subStatus) {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceSFSTurboUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	if d.HasChanges("size", "hpc_cache_bandwidth") {
		old, newSize := d.GetChange("size")
		if old.(int) > newSize.(int) {
			return diag.Errorf("shrinking SFS Turbo size is not supported")
		}

		if d.Get("share_type") != shareTypeHpcCache && d.HasChange("hpc_cache_bandwidth") {
			return diag.Errorf("only `HPC_CACHE` share type support updating HPC cache bandwidth")
		}

		hpcCacheBandwidth, err := convertHpcCacheBandwidth(d)
		if err != nil {
			return diag.FromErr(err)
		}

		isPrePaid := d.Get("charging_mode").(string) == "prePaid"
		requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{id}/action"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildTurboExtendOpts(newSize.(int), hpcCacheBandwidth, isPrePaid),
		}
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error expanding SFS Turbo size: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		if isPrePaid {
			orderId := utils.PathSearch("orderId", respBody, "").(string)
			if orderId == "" {
				return diag.Errorf("unable to find the order ID, this is a COM (Cloud Order Management) error, " +
					"please contact service for help and check your order status on the console.")
			}
			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if err := waitingForTurboSubStatusReady(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for SFS Turbo sub_status to be ready: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err := updateSFSTurboTags(client, d); err != nil {
			return diag.Errorf("error updating tags of SFS Turbo %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the SFS Turbo (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("name") {
		requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{id}/action"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200, 201, 202, 204},
			JSONBody: map[string]interface{}{
				"change_name": map[string]interface{}{
					"name": d.Get("name"),
				},
			},
		}
		_, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating name of SFS Turbo: %s", err)
		}
	}

	if d.HasChange("security_group_id") {
		requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{id}/action"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestBody := map[string]interface{}{
			"change_security_group": map[string]interface{}{
				"security_group_id": d.Get("security_group_id"),
			},
			"auto_create_security_group_rules": utils.ValueIgnoreEmpty(d.Get("auto_create_security_group_rules")),
		}

		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(requestBody),
		}
		_, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating security group of SFS Turbo: %s", err)
		}

		if err := waitingForTurboSubStatusReady(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for SFS Turbo sub_status to be ready: %s", err)
		}
	}

	return resourceSFSTurboRead(ctx, d, meta)
}

func updateSFSTurboTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	// remove old tags
	oldKeys := getOldTagKeys(d)
	if err := utils.DeleteResourceTagsWithKeys(client, oldKeys, "sfs-turbo", d.Id()); err != nil {
		return err
	}

	// set new tags
	return utils.CreateResourceTags(client, d, "sfs-turbo", d.Id())
}

func getOldTagKeys(d *schema.ResourceData) []string {
	oRaw, _ := d.GetChange("tags")
	var tagKeys []string
	if oMap := oRaw.(map[string]interface{}); len(oMap) > 0 {
		for k := range oMap {
			tagKeys = append(tagKeys, k)
		}
	}
	return tagKeys
}

func waitingForTurboSubStatusDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetTurboDetail(client, d.Id())
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "Resource Not Found", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceSFSTurboDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" {
		err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
		if err != nil {
			return diag.Errorf("error unsubscribing SFS Turbo: %s", err)
		}
	} else {
		requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{id}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err := client.Request("DELETE", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error deleting SFS Turbo: %s", err)
		}
	}

	if err := waitingForTurboSubStatusDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo to be deleted: %s", err)
	}
	return nil
}
