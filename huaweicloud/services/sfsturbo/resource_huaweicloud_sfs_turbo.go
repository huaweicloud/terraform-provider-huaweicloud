package sfsturbo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sfs_turbo/v1/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	prepaidUnitMonth int = 2
	prepaidUnitYear  int = 3

	autoRenewDisabled int = 0
	autoRenewEnabled  int = 1

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

func buildTurboMetadataOpts(d *schema.ResourceData) shares.Metadata {
	metaOpts := shares.Metadata{}
	if v, ok := d.GetOk("crypt_key_id"); ok {
		metaOpts.CryptKeyID = v.(string)
	}
	if v, ok := d.GetOk("dedicated_flavor"); ok {
		metaOpts.DedicatedFlavor = v.(string)
	}
	if v, ok := d.GetOk("dedicated_storage_id"); ok {
		metaOpts.DedicatedStorageID = v.(string)
	}

	switch d.Get("share_type").(string) {
	case shareTypeHpc:
		metaOpts.ExpandType = "hpc"
		metaOpts.HpcBw = d.Get("hpc_bandwidth").(string)
	case shareTypeHpcCache:
		metaOpts.ExpandType = "hpc_cache"
		metaOpts.HpcBw = d.Get("hpc_cache_bandwidth").(string)
	default:
		if _, ok := d.GetOk("enhanced"); ok {
			metaOpts.ExpandType = "bandwidth"
		}
	}
	return metaOpts
}

func convertShareType(d *schema.ResourceData) string {
	shareType := d.Get("share_type").(string)
	if shareType == shareTypeStandard || shareType == shareTypePerformance {
		return shareType
	}
	// For `HPC` and `HPC_CACHE` types, create API does not validate the share type.
	// We can fill in `STANDARD` or `PERFORMANCE`.
	return shareTypeStandard
}

func buildTurboCreateOpts(cfg *config.Config, d *schema.ResourceData) shares.CreateOpts {
	result := shares.CreateOpts{
		Share: shares.Share{
			Name:                d.Get("name").(string),
			Size:                d.Get("size").(int),
			ShareProto:          d.Get("share_proto").(string),
			VpcID:               d.Get("vpc_id").(string),
			SubnetID:            d.Get("subnet_id").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			AvailabilityZone:    d.Get("availability_zone").(string),
			BackupID:            d.Get("backup_id").(string),
			ShareType:           convertShareType(d),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
			Metadata:            buildTurboMetadataOpts(d),
		},
	}
	if d.Get("charging_mode") == "prePaid" {
		billing := shares.BssParam{
			PeriodNum: d.Get("period").(int),
			IsAutoPay: utils.Int(1), // Always enable auto-pay.
		}
		if d.Get("period_unit").(string) == "month" {
			billing.PeriodType = prepaidUnitMonth
		} else {
			billing.PeriodType = prepaidUnitYear
		}
		if d.Get("auto_renew").(string) == "true" {
			billing.IsAutoRenew = utils.Int(autoRenewEnabled)
		} else {
			billing.IsAutoRenew = utils.Int(autoRenewDisabled)
		}
		result.BssParam = &billing
	}
	return result
}

func validateParameter(d *schema.ResourceData) error {
	_, isHpcBandwidthSet := d.GetOk("hpc_bandwidth")
	_, isHpcCacheBandwidthSet := d.GetOk("hpc_cache_bandwidth")
	switch d.Get("share_type").(string) {
	case shareTypeHpc:
		if !isHpcBandwidthSet {
			return fmt.Errorf("`hpc_bandwidth` is required when share type is HPC")
		}
	case shareTypeHpcCache:
		if !isHpcCacheBandwidthSet {
			return fmt.Errorf("`hpc_cache_bandwidth` is required when share type is HPC_CACHE")
		}
		if d.Get("charging_mode").(string) == "prePaid" {
			return fmt.Errorf("HPC_CACHE share type only support in postpaid charging mode")
		}
	default:
		if isHpcBandwidthSet || isHpcCacheBandwidthSet {
			return fmt.Errorf("`hpc_bandwidth` and `hpc_cache_bandwidth` cannot be set when share type is" +
				" STANDARD or PERFORMANCE")
		}
	}
	return nil
}

func resourceSFSTurboCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	if err := validateParameter(d); err != nil {
		return diag.FromErr(err)
	}

	createOpts := buildTurboCreateOpts(cfg, d)
	log.Printf("[DEBUG] create sfs turbo with option: %+v", createOpts)
	resp, err := shares.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating SFS Turbo: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId := resp.OrderId
		if orderId == "" {
			return diag.Errorf("unable to find the order ID, this is a COM (Cloud Order Management) error, " +
				"please contact service for help and check your order status on the console.")
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
		d.SetId(resp.ID)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      createWaitForSFSTurboStatus(sfsClient, resp.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for SFS Turbo (%s) to become ready: %s ", d.Id(), stateErr)
	}

	// add tags
	if err := utils.CreateResourceTags(sfsClient, d, "sfs-turbo", d.Id()); err != nil {
		return diag.Errorf("error setting tags of SFS Turbo %s: %s", d.Id(), err)
	}

	return resourceSFSTurboRead(ctx, d, meta)
}

func flattenSize(n *shares.Turbo) interface{} {
	// n.Size is a string of float64, should convert it to int
	if fsize, err := strconv.ParseFloat(n.Size, 64); err == nil {
		return int(fsize)
	}

	return nil
}

func flattenStatus(n *shares.Turbo) interface{} {
	if n.SubStatus != "" {
		return n.SubStatus
	}

	return n.Status
}

func resourceSFSTurboRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		if hasSpecifyErrorCode403(err, "SFS.TURBO.9000") {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "SFS Turbo")
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", n.Name),
		d.Set("share_proto", n.ShareProto),
		d.Set("vpc_id", n.VpcID),
		d.Set("subnet_id", n.SubnetID),
		d.Set("security_group_id", n.SecurityGroupID),
		d.Set("version", n.Version),
		d.Set("region", region),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("available_capacity", n.AvailCapacity),
		d.Set("export_location", n.ExportLocation),
		d.Set("crypt_key_id", n.CryptKeyID),
		d.Set("enterprise_project_id", n.EnterpriseProjectId),
		d.Set("backup_id", n.BackupId),
		d.Set("size", flattenSize(n)),
		d.Set("status", flattenStatus(n)),
	)

	// Cannot obtain the billing parameters for pre-paid.

	// `HPC` and `HPC_CACHE` are custom types. `STANDARD` and `PERFORMANCE` are system types.
	switch n.ExpandType {
	case "hpc":
		mErr = multierror.Append(
			mErr,
			d.Set("share_type", shareTypeHpc),
			d.Set("hpc_bandwidth", n.HpcBw),
		)
	case "hpc_cache":
		mErr = multierror.Append(
			mErr,
			d.Set("share_type", shareTypeHpcCache),
			d.Set("hpc_cache_bandwidth", n.HpcBw),
		)
	default:
		mErr = multierror.Append(mErr, d.Set("share_type", n.ShareType))
		if n.ExpandType == "bandwidth" {
			mErr = multierror.Append(mErr, d.Set("enhanced", true))
		} else {
			mErr = multierror.Append(mErr, d.Set("enhanced", false))
		}
	}

	// set tags
	err = utils.SetResourceTagsToState(d, sfsClient, "sfs-turbo", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTurboUpdateOpts(newSize, hpcCacheBandwidth int, isPrePaid bool) shares.ExpandOpts {
	expandOpts := shares.ExtendOpts{
		NewSize:      newSize,
		NewBandwidth: hpcCacheBandwidth,
	}
	if isPrePaid {
		expandOpts.BssParam = &shares.BssParamExtend{
			IsAutoPay: utils.Int(1),
		}
	}
	return shares.ExpandOpts{
		Extend: expandOpts,
	}
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

func resourceSFSTurboUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	resourceId := d.Id()
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
		updateOpts := buildTurboUpdateOpts(newSize.(int), hpcCacheBandwidth, isPrePaid)
		resp, err := shares.Expand(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error expanding SFS Turbo size: %s", err)
		}

		if isPrePaid {
			orderId := resp.OrderId
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
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      waitForSFSTurboSubStatus(sfsClient, resourceId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error updating SFS Turbo: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err := updateSFSTurboTags(sfsClient, d); err != nil {
			return diag.Errorf("error updating tags of SFS Turbo %s: %s", resourceId, err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the SFS Turbo (%s): %s", resourceId, err)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := shares.UpdateNameOpts{
			Name: d.Get("name").(string),
		}
		err = shares.UpdateName(sfsClient, d.Id(), updateNameOpts).Err
		if err != nil {
			return diag.Errorf("error updating name of SFS Turbo: %s", err)
		}
	}

	if d.HasChange("security_group_id") {
		updateSecurityGroupIdOpts := shares.UpdateSecurityGroupIdOpts{
			SecurityGroupId: d.Get("security_group_id").(string),
		}
		err = shares.UpdateSecurityGroupId(sfsClient, d.Id(), updateSecurityGroupIdOpts).Err
		if err != nil {
			return diag.Errorf("error updating security group ID of SFS Turbo: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      waitForSFSTurboSubStatus(sfsClient, resourceId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error updating SFS Turbo: %s", err)
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

func resourceSFSTurboDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	resourceId := d.Id()
	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" {
		err := common.UnsubscribePrePaidResource(d, cfg, []string{resourceId})
		if err != nil {
			return diag.Errorf("error unsubscribing SFS Turbo: %s", err)
		}
	} else {
		err = shares.Delete(sfsClient, resourceId).ExtractErr()
		if err != nil {
			if hasSpecifyErrorCode403(err, "SFS.TURBO.9000") || hasSpecifyErrorCode400(err, "SFS.TURBO.0002") {
				err = golangsdk.ErrDefault404{}
			}
			return common.CheckDeletedDiag(d, err, "SFS Turbo")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      deleteWaitForSFSTurboStatus(sfsClient, resourceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting SFS Turbo: %s", err)
	}
	return nil
}

func createWaitForSFSTurboStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			return nil, "ERROR", err
		}

		if utils.StrSliceContains([]string{"303", "800"}, resp.Status) {
			return resp, "ERROR", fmt.Errorf("unexpected status: '%s'", resp.Status)
		}

		if utils.StrSliceContains([]string{"200"}, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func deleteWaitForSFSTurboStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		return resp, "PENDING", nil
	}
}

func waitForSFSTurboSubStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		// '321' indicate expansion failed
		// '332' indicate changing security group failed
		if utils.StrSliceContains([]string{"321", "332"}, resp.SubStatus) {
			return resp, "ERROR", fmt.Errorf("unexpected status: '%s'", resp.SubStatus)
		}

		// '221' indicate expansion succeeded
		// '232' indicate changing security group succeeded
		if utils.StrSliceContains([]string{"221", "232"}, resp.SubStatus) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

// When the SFS Turbo does not exist, the response body example of the details interface is as follows:
// {"errCode":"SFS.TURBO.0002","errMsg":"cluster not found"}
func hasSpecifyErrorCode400(err error, specCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("errCode", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse errCode from response body: %s", parseErr)
			}

			if errorCode == specCode {
				return true
			}
		}
	}

	return false
}

// When the SFS Turbo does not exist, the response body example of the details interface is as follows:
// {"errCode":"SFS.TURBO.9000","errMsg":"no privileges to operate"}
func hasSpecifyErrorCode403(err error, specCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("errCode", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse errCode from response body: %s", parseErr)
			}

			if errorCode == specCode {
				return true
			}
		}
	}

	return false
}
