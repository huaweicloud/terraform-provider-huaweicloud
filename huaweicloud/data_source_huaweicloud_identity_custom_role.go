package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/policies"
)

func DataSourceIdentityCustomRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityCustomRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"references": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceIdentityCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	allPages, err := policies.List(identityClient).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query roles: %s", err)
	}

	roles, err := policies.ExtractPageRoles(allPages)

	conditions := map[string]interface{}{}

	if v, ok := d.GetOk("name"); ok {
		conditions["name"] = v.(string)
	}
	if v, ok := d.GetOk("id"); ok {
		conditions["id"] = v.(string)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		conditions["domain_id"] = v.(string)
	}
	if v, ok := d.GetOk("references"); ok {
		conditions["references"] = v.(int)
	}
	if v, ok := d.GetOk("description"); ok {
		conditions["description"] = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		conditions["type"] = v.(string)
	}

	var allRoles []policies.Role

	for _, role := range roles {
		if rolesFilter(role, conditions) {
			allRoles = append(allRoles, role)
		}
	}

	if len(allRoles) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allRoles) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allRoles)
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}
	role := allRoles[0]

	return dataSourceIdentityCustomRoleAttributes(d, config, &role)
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityCustomRoleAttributes(d *schema.ResourceData, config *Config, role *policies.Role) error {
	log.Printf("[DEBUG] huaweicloud_identity_role details: %#v", role)

	d.SetId(role.ID)
	d.Set("name", role.Name)
	d.Set("domain_id", role.DomainId)
	d.Set("references", role.References)
	d.Set("catalog", role.Catalog)
	d.Set("description", role.Description)
	d.Set("type", role.Type)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmt.Errorf("Error marshalling policy: %s", err)
	}

	d.Set("policy", string(policy))

	return nil
}

func rolesFilter(role policies.Role, conditions map[string]interface{}) bool {
	if v, ok := conditions["name"]; ok && v != role.Name {
		return false
	}
	if v, ok := conditions["id"]; ok && v != role.ID {
		return false
	}
	if v, ok := conditions["domain_id"]; ok && v != role.DomainId {
		return false
	}
	if v, ok := conditions["references"]; ok && v != role.References {
		return false
	}
	if v, ok := conditions["description"]; ok && v != role.Description {
		return false
	}
	if v, ok := conditions["type"]; ok && v != role.Type {
		return false
	}
	return true
}
