package iam

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-ROLE/roles
func DataSourceIdentityCustomRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCustomRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"references": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityCustomRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	allPages, err := policies.List(identityClient).AllPages()
	if err != nil {
		return diag.Errorf("unable to query IAM custom policies: %s", err)
	}

	allPolicies, err := policies.ExtractPageRoles(allPages)
	if err != nil {
		return diag.Errorf("unable to extract IAM custom policies: %s", err)
	}

	conditions := map[string]interface{}{
		"ID":          d.Get("id").(string),
		"Name":        d.Get("name").(string),
		"Type":        d.Get("type").(string),
		"Description": d.Get("description").(string),
		"DomainId":    d.Get("domain_id").(string),
		"References":  d.Get("references").(int),
	}

	filterPolicies, err := utils.FilterSliceWithField(allPolicies, conditions)
	if err != nil {
		return diag.Errorf("filter IAM custom policies failed: %s", err)
	}

	if len(filterPolicies) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(filterPolicies) > 1 {
		return diag.Errorf("your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	role := filterPolicies[0].(policies.Role)
	log.Printf("[DEBUG] retrieve IAM custom policy: %#v", role)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return diag.Errorf("error marshaling the policy of IAM custom policy: %s", err)
	}

	d.SetId(role.ID)
	mErr := multierror.Append(nil,
		d.Set("name", role.Name),
		d.Set("domain_id", role.DomainId),
		d.Set("references", role.References),
		d.Set("catalog", role.Catalog),
		d.Set("description", role.Description),
		d.Set("type", role.Type),
		d.Set("policy", string(policy)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM custom policy fields: %s", err)
	}
	return nil
}
