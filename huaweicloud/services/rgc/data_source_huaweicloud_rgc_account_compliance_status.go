package rgc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/managed-accounts/{managed_account_id}/compliance-status
func DataSourceAccountComplianceStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountComplianceStatusRead,
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
				Optional: true,
			},
			"compliance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccountComplianceStatusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getAccountComplianceStatusProduct = "rgc"
	getAccountComplianceStatusClient, err := cfg.NewServiceClient(getAccountComplianceStatusProduct, region)
	if err != nil {
		return diag.Errorf("Error creating RGC client: %s", err)
	}

	getAccountComplianceStatusRespBody, err := getAccountComplianceStatus(getAccountComplianceStatusClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC account compliance: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("compliance_status", utils.PathSearch("compliance_status", getAccountComplianceStatusRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAccountComplianceStatus(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	managedAccountId := d.Get("managed_account_id").(string)
	var (
		httpUrl = "v1/governance/managed-accounts/{managed_account_id}/compliance-status"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{managed_account_id}", managedAccountId)
	getPath += buildAccountComplianceStatusQueryParams(d)

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

func buildAccountComplianceStatusQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("control_id"); ok {
		res = fmt.Sprintf("%s&control_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
