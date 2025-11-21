package rgc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/managed-accounts/{managed_account_id}/controls/{control_id}
func DataSourceAccountEnabledControl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountEnabledControlRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     accountEnabledControlsSchema(),
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func accountEnabledControlsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"behavior": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regional_preference": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceAccountEnabledControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getAccountEnabledControlProduct = "rgc"
	getAccountEnabledControlClient, err := cfg.NewServiceClient(getAccountEnabledControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getAccountEnabledControlRespBody, err := getAccountEnabledControl(getAccountEnabledControlClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC account enabled control: %s", err)
	}

	controlDetail := parseControlDetail(getAccountEnabledControlRespBody)

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("control_detail", controlDetail),
		d.Set("regions", utils.PathSearch("regions", getAccountEnabledControlRespBody, nil)),
		d.Set("state", utils.PathSearch("state", getAccountEnabledControlRespBody, nil)),
		d.Set("message", utils.PathSearch("message", getAccountEnabledControlRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getAccountEnabledControlRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAccountEnabledControl(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	managedAccountId := d.Get("managed_account_id").(string)
	controlId := d.Get("control_id").(string)

	var (
		getAccountEnabledControlHttpUrl = "v1/governance/managed-accounts/{managed_account_id}/controls/{control_id}"
	)
	getAccountEnabledControlHttpPath := client.Endpoint + getAccountEnabledControlHttpUrl
	getAccountEnabledControlHttpPath = strings.ReplaceAll(getAccountEnabledControlHttpPath, "{managed_account_id}", managedAccountId)
	getAccountEnabledControlHttpPath = strings.ReplaceAll(getAccountEnabledControlHttpPath, "{control_id}", controlId)

	getAccountEnabledControlHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountEnabledControlHttpResp, err := client.Request("GET", getAccountEnabledControlHttpPath, &getAccountEnabledControlHttpOpt)
	if err != nil {
		return nil, err
	}
	getAccountEnabledControlRespBody, err := utils.FlattenResponse(getAccountEnabledControlHttpResp)
	if err != nil {
		return nil, err
	}
	return getAccountEnabledControlRespBody, nil
}

func parseControlDetail(respBody interface{}) []interface{} {
	controlDetailList := make([]interface{}, 0)

	controlDetail := utils.PathSearch("control", respBody, nil)
	if controlDetail != nil {
		controlDetailList = append(controlDetailList, controlDetail)
	}

	return controlDetailList
}
