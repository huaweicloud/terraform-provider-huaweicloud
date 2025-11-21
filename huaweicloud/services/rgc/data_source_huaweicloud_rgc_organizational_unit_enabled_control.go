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

// @API RGC GET /v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls/{control_id}
func DataSourceOrganizationalUnitEnabledControl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationalUnitEnabledControlRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     organizationalUnitEnabledControlSchema(),
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

func organizationalUnitEnabledControlSchema() *schema.Resource {
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

func dataSourceOrganizationalUnitEnabledControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getOrganizationalUnitControlProduct := "rgc"
	getOrganizationalUnitControlClient, err := cfg.NewServiceClient(getOrganizationalUnitControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating rgc client: %s", err)
	}

	getOrganizationalUnitEnabledControlRespBody, err := getOrganizationalUnitEnabledControl(getOrganizationalUnitControlClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC organizational unit enabled control: %s", err)
	}

	control := parseOrganizationalUnitControlDetail(getOrganizationalUnitEnabledControlRespBody)

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("control", control),
		d.Set("regions", utils.PathSearch("regions", getOrganizationalUnitEnabledControlRespBody, nil)),
		d.Set("state", utils.PathSearch("state", getOrganizationalUnitEnabledControlRespBody, nil)),
		d.Set("message", utils.PathSearch("message", getOrganizationalUnitEnabledControlRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getOrganizationalUnitEnabledControlRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getOrganizationalUnitEnabledControl(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	managedOrganizationalUnitId := d.Get("managed_organizational_unit_id").(string)
	controlId := d.Get("control_id").(string)

	var (
		httpUrl = "v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls/{control_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{managed_organizational_unit_id}", managedOrganizationalUnitId)
	getPath = strings.ReplaceAll(getPath, "{control_id}", controlId)

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
	return getRespBody, nil
}

func parseOrganizationalUnitControlDetail(respBody interface{}) []interface{} {
	controlDetailList := make([]interface{}, 0)

	controlDetail := utils.PathSearch("control", respBody, nil)
	if controlDetail != nil {
		controlDetailList = append(controlDetailList, controlDetail)
	}

	return controlDetailList
}
