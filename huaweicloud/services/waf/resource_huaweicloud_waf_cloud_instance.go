package waf

import (
	"context"
	"errors"
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
	ChargingModePrePaid  ChargingMode = "prePaid"
	ChargingModePostPaid ChargingMode = "postPaid"

	SpecCodeIntroduction SpecCode = "detection"    // Introduction edition.
	SpecCodeStandard     SpecCode = "professional" // Standard edition (The old is professional edition).
	SpecCodeProfessional SpecCode = "enterprise"   // Professional edition (The old is enterprise edition).
	SpecCodePlatinum     SpecCode = "ultimate"     // Platinum edition (The old is ultimate edition).

	ResourceTypeInstance        ResourceType = "hws.resource.type.waf"                 // prepaid resource type
	ResourceTypeBandwidth       ResourceType = "hws.resource.type.waf.bandwidth"       // prepaid resource type
	ResourceTypeDomain          ResourceType = "hws.resource.type.waf.domain"          // prepaid resource type
	ResourceTypeRule            ResourceType = "hws.resource.type.waf.rule"            // prepaid resource type
	ResourceTypePayPerUseDomain ResourceType = "hws.resource.type.waf.payperusedomain" // postpaid resource type
)

func expackProductSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of extended packages.",
			},
		},
	}
}

// @API WAF DELETE /v1/{project_id}/waf/postpaid
// @API WAF POST /v1/{project_id}/waf/postpaid
// @API WAF POST /v1/{project_id}/waf/subscription/batchalter/prepaid-cloud-waf
// @API WAF POST /v1/{project_id}/waf/subscription/purchase/prepaid-cloud-waf
// @API WAF GET /v1/{project_id}/waf/subscription
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceCloudInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudInstanceCreate,
		ReadContext:   resourceCloudInstanceRead,
		UpdateContext: resourceCloudInstanceUpdate,
		DeleteContext: resourceCloudInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
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
				Description: "Specifies the region where the cloud WAF is located.",
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, false),
				Description: "Specifies the charging mode of the cloud WAF.",
			},
			"website": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"resource_spec_code", "period_unit", "period", "auto_renew",
					"bandwidth_expack_product", "domain_expack_product", "rule_expack_product"},
				Description: "Specifies the website to which the account belongs.",
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SpecCodeIntroduction),
					string(SpecCodeStandard),
					string(SpecCodeProfessional),
					string(SpecCodePlatinum),
				}, false),
				Description: "Specifies the specification of the cloud WAF.",
			},
			"period_unit": common.SchemaPeriodUnit(nil),
			"period":      common.SchemaPeriod(nil),
			"auto_renew":  common.SchemaAutoRenewUpdatable(nil),
			"bandwidth_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "Specifies the configuration of the bandwidth extended packages.",
			},
			"domain_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "Specifies the configuration of the domain extended packages.",
			},
			"rule_expack_product": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        expackProductSchema(),
				MaxItems:    1,
				Description: "Specifies the configuration of the rule extended packages.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the cloud WAF belongs.",
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

	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		orderId, err := createPrePaidCloudInstance(wafClient, cfg, region, d)
		if err != nil {
			return diag.Errorf("error creating prepaid cloud WAF: %s", err)
		}
		if orderId == nil {
			return diag.Errorf("error creating prepaid cloud WAF, cause cannot find order id in response")
		}

		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("the order is not completed while creating prepaid cloud WAF: %s", err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		if err := validatePostPaidParameter(d); err != nil {
			return diag.FromErr(err)
		}

		instance, err := createPostPaidCloudInstance(wafClient, cfg, region, d)
		if err != nil {
			return diag.Errorf("error creating postpaid cloud WAF: %s", err)
		}
		if instance == nil {
			return diag.Errorf("error creating postpaid cloud WAF, cause cannot find instance in response")
		}

		resourceId, err := flattenResourceIdByType(instance, ResourceTypePayPerUseDomain)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	}

	return resourceCloudInstanceRead(ctx, d, meta)
}

func validatePostPaidParameter(d *schema.ResourceData) error {
	if _, ok := d.GetOk("website"); !ok {
		return fmt.Errorf("`website` must be specified in postpaid charging mode")
	}
	return nil
}

func createPrePaidCloudInstance(wafClient *golangsdk.ServiceClient, cfg *config.Config, region string,
	d *schema.ResourceData) (*string, error) {
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
		EnterpriseProjectId:        cfg.GetEnterpriseProjectID(d),
	}
	return clouds.Create(wafClient, opts)
}

func createPostPaidCloudInstance(wafClient *golangsdk.ServiceClient, cfg *config.Config, region string,
	d *schema.ResourceData) (*clouds.Instance, error) {
	opts := clouds.CreatePostPaidOpts{
		Region:              region,
		ConsoleArea:         d.Get("website").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
	return clouds.CreatePostPaid(wafClient, opts)
}

func flattenResourceIdByType(instance *clouds.Instance, resourceType ResourceType) (string, error) {
	for _, v := range instance.Resources {
		if v.Type == string(resourceType) {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("cannot find target resource type (%s) from response", resourceType)
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
		return "", fmt.Errorf("invalid specification code, want 'waf.xxx' or 'waf.xxx.xxx', but '%s'", specCode)
	}
	return result[1], nil
}

func QueryCloudInstance(client *golangsdk.ServiceClient, instanceId, epsId string) (*clouds.Instance,
	ChargingMode, error) {
	default404Err := golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("the cloud WAF (%s) does not exist", instanceId)),
		},
	}
	resp, err := clouds.GetWithEpsID(client, epsId)
	if err != nil {
		return nil, "", err
	}
	if resp == nil || len(resp.Resources) < 1 {
		return nil, "", default404Err
	}
	for _, val := range resp.Resources {
		if val.Type == string(ResourceTypeInstance) && val.ID == instanceId {
			return resp, ChargingModePrePaid, nil
		}
		if val.Type == string(ResourceTypePayPerUseDomain) && val.ID == instanceId {
			return resp, ChargingModePostPaid, nil
		}
	}
	return nil, "", default404Err
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
	epsId := d.Get("enterprise_project_id").(string)
	resp, chargingMode, err := QueryCloudInstance(wafClient, instanceId, epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "cloud WAF")
	}
	log.Printf("[DEBUG] The resources list of the cloud WAF is: %#v", resp.Resources)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("charging_mode", chargingMode),
	)
	for _, val := range resp.Resources {
		switch val.Type {
		case string(ResourceTypeInstance):
			specCode, err := analysisSpecCode(val.SpecCode)
			if err != nil {
				return diag.FromErr(err)
			}
			mErr = multierror.Append(mErr,
				d.Set("resource_spec_code", specCode),
				d.Set("status", val.Status),
			)
		case string(ResourceTypeBandwidth):
			mErr = multierror.Append(mErr, d.Set("bandwidth_expack_product", flattenExtendedPackages(val)))
		case string(ResourceTypeDomain):
			mErr = multierror.Append(mErr, d.Set("domain_expack_product", flattenExtendedPackages(val)))
		case string(ResourceTypeRule):
			mErr = multierror.Append(mErr, d.Set("rule_expack_product", flattenExtendedPackages(val)))
		case string(ResourceTypePayPerUseDomain):
			// fill the domain status for postpaid resource type
			mErr = multierror.Append(mErr, d.Set("status", val.Status))
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving cloud WAF fields: %s", err)
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
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
		updateOpts = clouds.UpdateOpts{
			ProjectId:           wafClient.ProviderClient.ProjectID,
			IsAutoPay:           utils.Bool(true),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
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
			return fmt.Errorf("the order is not completed while creating extended packages: %v", err)
		}
	}

	if doesExpackChanged {
		orderId, err := clouds.Update(wafClient, updateOpts)
		if err != nil {
			return fmt.Errorf("error updating extended packages: %s", err)
		}

		err = common.WaitOrderComplete(ctx, bssClient, *orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("the order is not completed while updating extended packages: %v", err)
		}
	}

	return nil
}

func resourceCloudInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("charging_mode").(string) == string(ChargingModePostPaid) {
		return diag.Errorf("the postpaid charging mode cloud WAF instances cannot be updated.")
	}

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
				EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
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
			return diag.Errorf("the order is not completed while updating specification code: %v", err)
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

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "waf",
			RegionId:     region,
			ProjectId:    wafClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCloudInstanceRead(ctx, d, meta)
}

func cloudInstanceDeleteRefreshFunc(client *golangsdk.ServiceClient, instanceId, epsId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := QueryCloudInstance(client, instanceId, epsId)
		var errDefault404 golangsdk.ErrDefault404
		if errors.As(err, &errDefault404) {
			return "success_deleted", "DELETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceCloudInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF v1 client: %s", err)
	}

	chargingMode := d.Get("charging_mode").(string)
	epsId := cfg.GetEnterpriseProjectID(d)
	if chargingMode == string(ChargingModePostPaid) {
		opts := clouds.DeletePostPaidOpts{
			Region:              region,
			EnterpriseProjectId: epsId,
		}
		if err = clouds.DeletePostPaid(wafClient, opts); err != nil {
			return diag.Errorf("error deleting the postpaid cloud WAF: %s", err)
		}
		return nil
	}

	instanceId := d.Id()
	err = common.UnsubscribePrePaidResource(d, cfg, []string{instanceId})
	if err != nil {
		// When the resource does not exist, the API for unsubscribing prePaid resource will return a `400` status code,
		// and the response body is as follows:
		// {"error_code": "CBC.30000067",
		// "error_msg": "Unsubscription not supported. This resource has been deleted or the subscription to this resource has
		// not been synchronized to ..."}
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "CBC.30000067"),
			"error unsubscribing WAF cloud instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      cloudInstanceDeleteRefreshFunc(wafClient, instanceId, epsId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting to delete the prepaid cloud WAF (%s): %s", instanceId, err)
	}
	return nil
}
