package cdn

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

// @API CDN GET /v1.0/cdn/configuration/domains/{domain_name}/domain-verifies
func DataSourceDomainOwnerVerification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainOwnerVerificationRead,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the accelerated domain.`,
			},

			// Attributes.
			"dns_verify_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DNS resolution type.`,
			},
			"dns_verify_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DNS resolution host record name.`,
			},
			"file_verify_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file verification URL address.`,
			},
			"verify_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The verification domain name.`,
			},
			"file_verify_filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file verification filename.`,
			},
			"verify_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The verification value, which is the resolution value or file content.`,
			},
			"file_verify_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of file verification domain names.`,
			},
		},
	}
}

func dataSourceDomainOwnerVerificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v1.0/cdn/configuration/domains/{domain_name}/domain-verifies"
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return diag.Errorf("error querying domain owner verification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("domain_name", utils.PathSearch("domain_name", respBody, nil)),
		d.Set("dns_verify_type", utils.PathSearch("dns_verify_type", respBody, nil)),
		d.Set("dns_verify_name", utils.PathSearch("dns_verify_name", respBody, nil)),
		d.Set("file_verify_url", utils.PathSearch("file_verify_url", respBody, nil)),
		d.Set("verify_domain_name", utils.PathSearch("verify_domain_name", respBody, nil)),
		d.Set("file_verify_filename", utils.PathSearch("file_verify_filename", respBody, nil)),
		d.Set("verify_content", utils.PathSearch("verify_content", respBody, nil)),
		d.Set("file_verify_domains", utils.PathSearch("file_verify_domains", respBody, make([]interface{}, 0))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
