package rgc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var landingZoneNonUpdatableParams = []string{"home_region", "organization_structure_type", "deny_ungoverned_regions", "kms_key_id", "organization_structure_type"}

// @API RGC POST /v1/landing-zone/setup
// @API RGC GET /v1/landing-zone/status
// @API RGC POST /v1/landing-zone/delete
func ResourceLandingZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLandingZoneCreate,
		UpdateContext: resourceLandingZoneUpdate,
		ReadContext:   resourceLandingZoneRead,
		DeleteContext: resourceLandingZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(landingZoneNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"home_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_configuration_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_configuration_status": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"organization_structure": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organizational_unit_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organizational_unit_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"accounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"account_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"phone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"account_email": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"logging_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logging_bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"logging_bucket": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     loggingBaselineConfiguration(),
						},
						"access_logging_bucket": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     loggingBaselineConfiguration(),
						},
					},
				},
			},
			"organization_structure_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
			},
			"identity_store_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identity_center_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ENABLE",
			},
			"deny_ungoverned_regions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cloud_trail_type": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"baseline_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"landing_zone_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployed_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func loggingBaselineConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"retention_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_multi_az": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceLandingZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		setupLandingZoneHttpUrl = "v1/landing-zone/setup"
		setupLandingZoneProduct = "rgc"
	)

	setupLandingZoneClient, err := cfg.NewServiceClient(setupLandingZoneProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	setupLandingZonePath := setupLandingZoneClient.Endpoint + setupLandingZoneHttpUrl
	setupLandingZoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSetupLandingZoneBodyParams(d, "CREATE")),
		OkCodes:          []int{200},
	}
	_, err = setupLandingZoneClient.Request("POST", setupLandingZonePath, &setupLandingZoneOpt)
	if err != nil {
		return diag.Errorf("error setup landing zone: %s", err)
	}
	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"in_progress"},
		Target:       []string{"succeeded"},
		Refresh:      landingZoneStateRefreshFunc(setupLandingZoneClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for setup landing zone: %s", err)
	}

	return resourceLandingZoneRead(ctx, d, meta)
}

func resourceLandingZoneRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getLandingZoneStatusHttpUrl = "v1/landing-zone/status"
		getOrganizationUnitProduct  = "rgc"
	)

	getLandingZoneStatusClient, err := cfg.NewServiceClient(getOrganizationUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getLandingZoneStatusPath := getLandingZoneStatusClient.Endpoint + getLandingZoneStatusHttpUrl

	getLandingZoneStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLandingZoneStatusResp, err := getLandingZoneStatusClient.Request("GET", getLandingZoneStatusPath, &getLandingZoneStatusOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving landing zone")
	}

	getLandingZoneStatusRespBody, err := utils.FlattenResponse(getLandingZoneStatusResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("landing_zone_status", utils.PathSearch("landing_zone_status", getLandingZoneStatusRespBody, nil)),
		d.Set("deployed_version", utils.PathSearch("deployed_version", getLandingZoneStatusRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLandingZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateLandingZoneHttpUrl = "v1/landing-zone/setup"
		updateLandingZoneProduct = "rgc"
	)

	updateLandingZoneClient, err := cfg.NewServiceClient(updateLandingZoneProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	updateLandingZonePath := updateLandingZoneClient.Endpoint + updateLandingZoneHttpUrl

	updateLandingZoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSetupLandingZoneBodyParams(d, "UPDATE"),
		OkCodes:          []int{200},
	}

	_, err = updateLandingZoneClient.Request("POST", updateLandingZonePath, &updateLandingZoneOpt)
	if err != nil {
		return diag.Errorf("error update landing zone: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"in_progress"},
		Target:       []string{"succeeded"},
		Refresh:      landingZoneStateRefreshFunc(updateLandingZoneClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for update landing zone: %s", err)
	}

	return resourceLandingZoneRead(ctx, d, meta)
}

func resourceLandingZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteLandingZoneHttpUrl = "v1/landing-zone/delete"
		deleteLandingZoneProduct = "rgc"
	)

	deleteLandingZoneClient, err := cfg.NewServiceClient(deleteLandingZoneProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	deleteLandingZonePath := deleteLandingZoneClient.Endpoint + deleteLandingZoneHttpUrl

	deleteLandingZoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200},
	}

	_, err = deleteLandingZoneClient.Request("POST", deleteLandingZonePath, &deleteLandingZoneOpt)
	if err != nil {
		return diag.Errorf("error delete landing zone: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"in_progress"},
		Target:       []string{"succeeded"},
		Refresh:      landingZoneStateRefreshFunc(deleteLandingZoneClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for delete landing zone: %s", err)
	}

	return nil
}

func filterOrganizationStructureEmptyValue(bodyParams map[string]interface{}) {
	ouStructures := bodyParams["organization_structure"].([]interface{})

	for _, ouStructure := range ouStructures {
		ouStructureMap := ouStructure.(map[string]interface{})

		if ouStructureMap["organizational_unit_name"] == "" {
			delete(ouStructureMap, "organizational_unit_name")
		}

		accounts := ouStructureMap["accounts"].([]interface{})

		for _, account := range accounts {
			accountMap := account.(map[string]interface{})

			if accountMap["account_email"] == "" {
				delete(accountMap, "account_email")
			}

			if accountMap["phone"] == "" {
				delete(accountMap, "phone")
			}
		}
	}
}

func buildSetupLandingZoneBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"home_region":                    d.Get("home_region").(string),
		"setup_landing_zone_action_type": action,
		"region_configuration_list":      d.Get("region_configuration_list").([]interface{}),
		"organization_structure":         d.Get("organization_structure").([]interface{}),
	}

	filterOrganizationStructureEmptyValue(bodyParams)

	if v, ok := d.GetOk("organization_structure_type"); ok {
		bodyParams["organization_structure_type"] = v
	}

	if v, ok := d.GetOk("identity_store_email"); ok {
		bodyParams["identity_store_email"] = v
	}

	if v, ok := d.GetOk("identity_center_status"); ok {
		bodyParams["identity_center_status"] = v
	}

	if v, ok := d.GetOk("deny_ungoverned_regions"); ok {
		bodyParams["deny_ungoverned_regions"] = v
	}

	if v, ok := d.GetOk("cloud_trail_type"); ok {
		bodyParams["cloud_trail_type"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		bodyParams["kms_key_id"] = v
	}

	if v, ok := d.GetOk("baseline_version"); ok {
		bodyParams["baseline_version"] = v
	}

	bodyParams["logging_configuration"] = parseLoggingConfiguration(d)

	return bodyParams
}

func parseLoggingConfiguration(d *schema.ResourceData) map[string]interface{} {
	configuration := d.Get("logging_configuration")
	if configuration == nil || len(configuration.([]interface{})) == 0 {
		return nil
	}
	loggingConfigurationMap := make(map[string]interface{})
	loggingBucketName := utils.PathSearch("[0].logging_bucket_name", configuration, nil)
	if loggingBucketName != "" {
		loggingConfigurationMap["logging_bucket_name"] = loggingBucketName
	}
	loggingConfigurationMap["logging_bucket"] = buildLoggingBucketBodyParams(utils.PathSearch("[0].logging_bucket", configuration, nil))
	loggingConfigurationMap["access_logging_bucket"] = buildLoggingBucketBodyParams(utils.PathSearch("[0].access_logging_bucket", configuration, nil))
	return loggingConfigurationMap
}

func buildLoggingBucketBodyParams(loggingBucket interface{}) map[string]interface{} {
	if loggingBucket == nil || len(loggingBucket.([]interface{})) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"retention_days":  utils.PathSearch("[0].retention_days", loggingBucket, nil),
		"enable_multi_az": utils.PathSearch("[0].enable_multi_az", loggingBucket, nil),
	}
	return bodyParams
}

func landingZoneStateRefreshFunc(client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getLandingZoneStatusHttpUrl := "v1/landing-zone/status"
		getLandingZoneStatusPath := client.Endpoint + getLandingZoneStatusHttpUrl

		getLandingZoneStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getLandingZoneStatusResp, err := client.Request("GET", getLandingZoneStatusPath, &getLandingZoneStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getLandingZoneStatusRespBody, err := utils.FlattenResponse(getLandingZoneStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("landing_zone_status", getLandingZoneStatusRespBody, "").(string)
		if status == "failed" || status == "" {
			message := utils.PathSearch("message", getLandingZoneStatusRespBody, "")
			return nil, "", fmt.Errorf("status: %s; message: %s", status, message)
		}

		return getLandingZoneStatusRespBody, status, nil
	}
}
