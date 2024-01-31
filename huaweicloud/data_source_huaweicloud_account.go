package huaweicloud

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v3/auth/domains
func DataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating IAM client: %s", err)
	}

	// ResourceBase: https://iam.{CLOUD}/v3/auth/
	identityClient.ResourceBase += "auth/"
	allPages, err := domains.List(identityClient, nil).AllPages()
	if err != nil {
		return diag.Errorf("failed to query account: %s", err)
	}

	accounts, err := domains.ExtractDomains(allPages)
	if err != nil {
		return diag.Errorf("failed to extract details of account: %s", err)
	}

	if len(accounts) == 0 {
		return diag.Errorf("failed to query account: you are not authorized to perform the action")
	}

	result := accounts[0]
	d.SetId(result.ID)

	mErr := multierror.Append(nil,
		d.Set("name", result.Name),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
