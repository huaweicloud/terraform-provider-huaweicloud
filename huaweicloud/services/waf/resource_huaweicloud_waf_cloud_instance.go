package waf

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf/v1/clouds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type (
	ChargingMode string
	SpecCode     string
	ResourceType string
)

const (
	ChargingModePrePaid ChargingMode = "prePaid"

	SpecCodeIntroduction SpecCode = "detection"    // Introduction edition.
	SpecCodeStandard     SpecCode = "professional" // Standard edition (The old is professional edition).
	SpecCodeProfessional SpecCode = "enterprise"   // Professional edition (The old is enterprise edition).
	SpecCodePlatinum     SpecCode = "ultimate"     // Platinum edition (The old is ultimate edition).

	ResourceTypeInstance  ResourceType = "hws.resource.type.waf"
	ResourceTypeBandwidth ResourceType = "hws.resource.type.waf.bandwidth"
	ResourceTypeDomain    ResourceType = "hws.resource.type.waf.domain"
	ResourceTypeRule      ResourceType = "hws.resource.type.waf.rule"
)

func expackProductSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of extended packages.",
			},
		},
	}
}

func ResourceCloudInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudInstanceCreate,
		ReadContext:   resourceCloudInstanceRead,
		UpdateContext: resourceCloudInstanceUpdate,
		DeleteContext: resourceCloudInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the cloud WAF is located.",
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SpecCodeIntroduction),
					string(SpecCodeStandard),
					string(SpecCodeProfessional),
					string(SpecCodePlatinum),
				}, false),
				Description: "The specification of the cloud WAF.",
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
				}, false),
				Description: "The charging mode of the cloud WAF.",
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"bandwidth_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "The configuration of the bandwidth extended packages.",
			},
			"domain_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "The configuration of the domain extended packages.",
			},
			"rule_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "The configuration of the rule extended packages.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the cloud WAF belongs.",
			},
			// Attributes
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current status of the cloud WAF.",
			},
		},
	}
}

func buildExtendedPackages(packages []interface{}) *clouds.ExpackProductInfo {
	if len(packages) < 1 {
		return nil
	}

	expack := packages[0].(map[string]interface{})
	return &clouds.ExpackProductInfo{
		ResourceSize: expack["resource_size"].(int),
	}
}

func getAutoRenewValue(autoRenew string) bool {
	if autoRenew == "" {
		return false
	}
	result, _ := strconv.ParseBool(autoRenew)
	return result
}

func resourceCloudInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF v1 client: %s", err)
	}

	opts := clouds.CreateOpts{
		ProjectId:   wafClient.ProviderClient.ProjectID,
		IsAutoPay:   utils.Bool(true),
		IsAutoRenew: utils.Bool(getAutoRenewValue(d.Get("auto_renew").(string))),
		RegionId:    region,
		ProductInfo: &clouds.ProductInfo{
			ResourceSpecCode: d.Get("resource_spec_code").(string),
			PeriodType:       d.Get("period_unit").(string),
			PeriodNum:        d.Get("period").(int),
		},
		BandwidthExpackProductInfo: buildExtendedPackages(d.Get("bandwidth_expack_product").([]interface{})),
		DomainExpackProductInfo:    buildExtendedPackages(d.Get("domain_expack_product").([]interface{})),
		RuleExpackProductInfo:      buildExtendedPackages(d.Get("rule_expack_product").([]interface{})),
		EnterpriseProjectId:        common.GetEnterpriseProjectID(d, cfg),
	}
	orderId, err := clouds.Create(wafClient, opts)
	if err != nil {
		return diag.Errorf("error creating cloud WAF: %s", err)
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("the order is not completed while creating cloud WAF: %#v", err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	return resourceCloudInstanceRead(ctx, d, meta)
}

// Because the request value and the response value of the WAF API are different, we need to match the response value.
// For example, the response value is 'waf.detection', but the request value is 'detection'.
// The 'detection' is what we need, but the prefix 'waf.' is unnecessary.
func analysisSpecCode(specCode string) (string, error) {
	// All specification codes:
	// + waf.detection
	// + waf.instance.professional
	// + waf.enterprise
	// + waf.ultimate
	re := regexp.MustCompile(`^waf\.(?:[a-z]+\.)?(\w+)$`)
	result := re.FindStringSubmatch(specCode)
	// The right length of the string match is two, first is the match result of the regex string,
	// second is the match result of the regex group.
	if len(result) < 2 {
		return "", fmt.Errorf("invalid specification code, want 'waf.xxx' or 'waf.xxx.xxx', but '%s'.", specCode)
	}
	return result[1], nil
}

func QueryCloudInstance(client *golangsdk.ServiceClient, instanceId string) (*clouds.Instance, error) {
	default404Err := golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("the cloud WAF (%s) does not exist", instanceId)),
		},
	}
	resp, err := clouds.Get(client)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Resources) < 1 {
		return nil, default404Err
	}
	for _, val := range resp.Resources {
		if val.Type == string(ResourceTypeInstance) && val.ID == instanceId {
			return resp, nil
		}
	}
	return nil, default404Err
}

func flattenExtendedPackages(res clouds.Resource) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"resource_size": res.Size,
		},
	}
}

func resourceCloudInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF v1 client: %s", err)
	}

	instanceId := d.Id()
	resp, err := QueryCloudInstance(wafClient, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "cloud WAF")
	}
	log.Printf("[DEBUG] The resources list of the cloud WAF is: %#v", resp.Resources)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
	)
	for _, val := range resp.Resources {
		if val.Type == string(ResourceTypeInstance) {
			specCode, err := analysisSpecCode(val.SpecCode)
			if err != nil {
				return diag.FromErr(err)
			}
			mErr = multierror.Append(mErr,
				d.Set("resource_spec_code", specCode),
				d.Set("status", val.Status),
			)
		}
		if val.Type == string(ResourceTypeBandwidth) {
			mErr = multierror.Append(mErr, d.Set("bandwidth_expack_product", flattenExtendedPackages(val)))
		}
		if val.Type == string(ResourceTypeDomain) {
			mErr = multierror.Append(mErr, d.Set("domain_expack_product", flattenExtendedPackages(val)))
		}
		if val.Type == string(ResourceTypeRule) {
			mErr = multierror.Append(mErr, d.Set("rule_expack_product", flattenExtendedPackages(val)))
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving cluod WAF fields: %s", err)
	}
	return nil
}

func updateExtendedPackages(ctx context.Context, wafClient *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config) error {
	var (
		region            = cfg.GetRegion(d)
		isNewExpack       = false
		doesExpackChanged = false

		createOpts = clouds.CreateOpts{
			ProjectId:           wafClient.ProviderClient.ProjectID,
			IsAutoPay:           utils.Bool(true),
			IsAutoRenew:         utils.Bool(getAutoRenewValue(d.Get("auto_renew").(string))),
			RegionId:            region,
			EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
		}
		updateOpts = clouds.UpdateOpts{
			ProjectId:           wafClient.ProviderClient.ProjectID,
			IsAutoPay:           utils.Bool(true),
			EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
		}
	)
	if d.HasChange("bandwidth_expack_product.0.resource_size") {
		old, _ := d.GetChange("bandwidth_expack_product.0.resource_size")
		if old.(int) == 0 {
			createOpts.BandwidthExpackProductInfo = buildExtendedPackages(d.Get("bandwidth_expack_product").([]interface{}))
			isNewExpack = true
		} else {
			updateOpts.BandwidthExpackProductInfo = buildExtendedPackages(d.Get("bandwidth_expack_product").([]interface{}))
			doesExpackChanged = true
		}
	}
	if d.HasChange("domain_expack_product.0.resource_size") {
		old, _ := d.GetChange("domain_expack_product.0.resource_size")
		if old.(int) == 0 {
			createOpts.DomainExpackProductInfo = buildExtendedPackages(d.Get("domain_expack_product").([]interface{}))
			isNewExpack = true
		} else {
			updateOpts.DomainExpackProductInfo = buildExtendedPackages(d.Get("domain_expack_product").([]interface{}))
			doesExpackChanged = true
		}
	}
	if d.HasChange("rule_expack_product.0.resource_size") {
		old, _ := d.GetChange("rule_expack_product.0.resource_size")
		if old.(int) == 0 {
			createOpts.RuleExpackProductInfo = buildExtendedPackages(d.Get("rule_expack_product").([]interface{}))
			isNewExpack = true
		} else {
			updateOpts.RuleExpackProductInfo = buildExtendedPackages(d.Get("rule_expack_product").([]interface{}))
			doesExpackChanged = true
		}
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}

	if isNewExpack {
		orderId, err := clouds.Create(wafClient, createOpts)
		if err != nil {
			return fmt.Errorf("error creating extended packages: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("the order is not completed while creating extended packages: %#v", err)
		}
	}

	if doesExpackChanged {
		orderId, err := clouds.Update(wafClient, updateOpts)
		if err != nil {
			return fmt.Errorf("error updating extended packages: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("the order is not completed while updating extended packages: %#v", err)
		}
	}

	return nil
}

func resourceCloudInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF v1 client: %s", err)
	}

	instanceId := d.Id()
	if d.HasChange("resource_spec_code") {
		var (
			opts = clouds.UpdateOpts{
				ProjectId:           cfg.GetProjectID(region),
				IsAutoPay:           utils.Bool(true),
				EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
			}
		)

		if d.HasChange("resource_spec_code") {
			opts.ProductInfo = &clouds.UpdateProductInfo{
				ResourceSpecCode: d.Get("resource_spec_code").(string),
			}
		}
		orderId, err := clouds.Update(wafClient, opts)
		if err != nil {
			return diag.Errorf("error updating cloud WAF (%s): %s", instanceId, err)
		}

		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("the order is not completed while updating specification code: %#v", err)
		}
	}

	if d.HasChanges("bandwidth_expack_product", "domain_expack_product", "rule_expack_product") {
		err = updateExtendedPackages(ctx, wafClient, d, cfg)
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
			return diag.Errorf("error updating the auto-renew of the cloud WAF (%s): %s", instanceId, err)
		}
	}

	return resourceCloudInstanceRead(ctx, d, meta)
}

func cloudInstanceDeleteRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var unprotectedHostId string = ""
		resp, err := QueryCloudInstance(client, instanceId)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return unprotectedHostId, "DELETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceCloudInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	instanceId := d.Id()
	err := common.UnsubscribePrePaidResource(d, cfg, []string{instanceId})
	if err != nil {
		return diag.Errorf("error unsubscribing cloud WAF: %s", err)
	}

	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF v1 client: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      cloudInstanceDeleteRefreshFunc(wafClient, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting to delete the cloud WAF (%s): %s", instanceId, err)
	}
	return nil
}
