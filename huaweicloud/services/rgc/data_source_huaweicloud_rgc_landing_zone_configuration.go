package rgc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/landing-zone/configuration
func DataSourceLandingZoneConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLandingZoneConfigurationRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"common_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"home_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_trail_type": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"identity_center_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_structure_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"logging_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logging_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_logging_bucket": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retention_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable_multi_az": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"logging_bucket": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retention_days": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable_multi_az": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"organization_structure": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organizational_unit_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organizational_unit_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_configuration_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLandingZoneConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getLandingZoneConfigurationProduct = "rgc"
	getLandingZoneConfigurationClient, err := cfg.NewServiceClient(getLandingZoneConfigurationProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getLandingZoneConfigurationRespBody, err := getLandingZoneConfiguration(getLandingZoneConfigurationClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC landing zone configuration: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("common_configuration", parseCommonConfiguration(getLandingZoneConfigurationRespBody)),
		d.Set("logging_configuration", getLoggingConfiguration(getLandingZoneConfigurationRespBody)),
		d.Set("organization_structure", utils.PathSearch("organization_structure", getLandingZoneConfigurationRespBody, nil)),
		d.Set("regions", utils.PathSearch("regions", getLandingZoneConfigurationRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseCommonConfiguration(respBody interface{}) []interface{} {
	commonConfigurationList := make([]interface{}, 0)

	commonConfiguration := utils.PathSearch("common_configuration", respBody, nil)
	if commonConfiguration != nil {
		value := commonConfiguration.(map[string]interface{})
		commonConfigurationList = append(commonConfigurationList, value)
	}

	return commonConfigurationList
}

func getLoggingConfiguration(respBody interface{}) []interface{} {
	loggingConfigurationList := make([]interface{}, 0)

	loggingConfiguration := utils.PathSearch("logging_configuration", respBody, nil)
	if loggingConfiguration != nil {
		loggingConfigurationMap := utils.RemoveNil(loggingConfiguration.(map[string]interface{}))
		value := make(map[string]interface{})
		if v, ok := loggingConfigurationMap["logging_buck_name"]; ok {
			value["logging_buck_name"] = v
		}

		if v, ok := loggingConfigurationMap["access_logging_bucket"]; ok {
			value["access_logging_bucket"] = []interface{}{v}
		}

		if v, ok := loggingConfigurationMap["logging_bucket"]; ok {
			value["logging_bucket"] = []interface{}{v}
		}
		loggingConfigurationList = append(loggingConfigurationList, value)
	}

	return loggingConfigurationList
}

func getLandingZoneConfiguration(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getLandingZoneConfigurationHttpUrl = "v1/landing-zone/configuration"
	)
	getLandingZoneConfigurationHttpPath := client.Endpoint + getLandingZoneConfigurationHttpUrl

	getLandingZoneConfigurationHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLandingZoneConfigurationHttpResp, err := client.Request("GET", getLandingZoneConfigurationHttpPath, &getLandingZoneConfigurationHttpOpt)
	if err != nil {
		return nil, err
	}
	getLandingZoneConfigurationRespBody, err := utils.FlattenResponse(getLandingZoneConfigurationHttpResp)
	if err != nil {
		return nil, err
	}
	return getLandingZoneConfigurationRespBody, nil
}
