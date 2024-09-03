package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	v1groups "github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	v3groups "github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	v3rules "github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v3/{project_id}/vpc/security-group-rules
// @API VPC GET /v3/{project_id}/vpc/security-groups
// @API VPC GET /v1/{project_id}/security-groups
func DataSourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": securityGroupRuleSchema,
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getRuleListByGroupId(client *golangsdk.ServiceClient, groupId string) ([]map[string]interface{}, error) {
	listOpts := v3rules.ListOpts{
		SecurityGroupId: groupId,
	}
	resp, err := v3rules.List(client, listOpts)
	if err != nil {
		return nil, err
	}
	return flattenSecurityGroupRulesV3(resp)
}

func dataSourceNetworkingSecGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v3Client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	listOpts := v3groups.ListOpts{
		ID:                  d.Get("secgroup_id").(string),
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	allSecGroups, err := v3groups.List(v3Client, listOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// If the v3 API does not exist or has not been published in the specified region, set again using v1 API.
			return dataSourceNetworkingSecGroupReadV1(ctx, d, meta)
		}
		return diag.Errorf("unable to get security groups list: %s", err)
	}

	if len(allSecGroups) < 1 {
		return diag.Errorf("no Security Group found.")
	}

	if len(allSecGroups) > 1 {
		return diag.Errorf("more than one Security Groups found.")
	}

	secGroup := allSecGroups[0]
	d.SetId(secGroup.ID)
	log.Printf("[DEBUG] Retrieved Security Group (%s) by v3 client: %v", d.Id(), secGroup)

	rules, err := getRuleListByGroupId(v3Client, secGroup.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] The rules list of security group (%s) is: %v", d.Id(), rules)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", secGroup.Name),
		d.Set("description", secGroup.Description),
		d.Set("rules", rules),
		d.Set("created_at", secGroup.CreatedAt),
		d.Set("updated_at", secGroup.UpdatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func dataSourceNetworkingSecGroupReadV1(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v1 client: %s", err)
	}

	listOpts := v1groups.ListOpts{
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	pages, err := v1groups.List(v1Client, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}
	allGroups, err := v1groups.ExtractSecurityGroups(pages)
	if err != nil {
		return diag.Errorf("error retrieving security groups list: %s", err)
	}
	if len(allGroups) == 0 {
		return diag.Errorf("no sucurity group found, please change your search criteria and try again.")
	}
	log.Printf("[DEBUG] The retrieved group list is: %v", allGroups)

	filter := map[string]interface{}{
		"ID":   d.Get("secgroup_id"),
		"Name": d.Get("name"),
	}
	filterGroups, err := utils.FilterSliceWithField(allGroups, filter)
	if err != nil {
		return diag.Errorf("error filting security groups list: %s", err)
	}
	if len(filterGroups) < 1 {
		return diag.Errorf("No Security Group found.")
	}
	if len(filterGroups) > 1 {
		return diag.Errorf("More than one Security Groups found.")
	}

	resp := filterGroups[0].(v1groups.SecurityGroup)
	d.SetId(resp.ID)

	rules := flattenSecurityGroupRulesV1(&resp)
	log.Printf("[DEBUG] The retrieved rules list is: %v", rules)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("rules", rules),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
