package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentitySecurityCompliance
// @API IAM GET /v3/domains/{domain_id}/config/security_compliance
// @API IAM GET /v3/domains/{domain_id}/config/security_compliance/{option}
func DataSourceIdentitySecurityCompliance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentitySecurityComplianceRead,

		Schema: map[string]*schema.Schema{
			"option": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"password_regex", "password_regex_description"}, false),
				Description:  "Specifies the query type, which supports password_regex or password_regex_description.",
			},

			"password_regex": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password_regex_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentitySecurityComplianceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	optionStr := d.Get("option").(string)
	if optionStr == "" {
		return showSecurityCompliance(iamClient, d, cfg)
	}
	return showSecurityComplianceByOption(iamClient, optionStr, d, cfg)
}

func showSecurityCompliance(iamClient *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) diag.Diagnostics {
	showSecurityCompliancePath := iamClient.Endpoint + "v3/domains/{domain_id}/config/security_compliance"
	showSecurityCompliancePath = strings.ReplaceAll(showSecurityCompliancePath, "{domain_id}", cfg.DomainID)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", showSecurityCompliancePath, &options)
	if err != nil {
		return diag.Errorf("showSecurityCompliance error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("password_regex", utils.PathSearch("config.security_compliance.password_regex", respBody, "")),
		d.Set("password_regex_description", utils.PathSearch("config.security_compliance.password_regex_description", respBody, "")),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SecurityCompliance fields: %s", err)
	}
	return nil
}

func showSecurityComplianceByOption(iamClient *golangsdk.ServiceClient, optionStr string, d *schema.ResourceData,
	cfg *config.Config) diag.Diagnostics {
	showSecurityComplianceByOptionPath := iamClient.Endpoint + "v3/domains/{domain_id}/config/security_compliance/" + optionStr
	showSecurityComplianceByOptionPath = strings.ReplaceAll(showSecurityComplianceByOptionPath, "{domain_id}", cfg.DomainID)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", showSecurityComplianceByOptionPath, &options)

	if err != nil {
		return diag.Errorf("showSecurityComplianceByOption error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if optionStr == "password_regex" {
		err = d.Set("password_regex", utils.PathSearch("config.password_regex", respBody, ""))
	} else if optionStr == "password_regex_description" {
		err = d.Set("password_regex_description", utils.PathSearch("config.password_regex_description", respBody, ""))
	}
	if err != nil {
		return diag.Errorf("error setting showSecurityComplianceByOption fields: %s", err)
	}
	return nil
}
