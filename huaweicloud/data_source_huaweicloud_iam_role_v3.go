package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	sdkroles "github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
)

func dataSourceIAMRoleV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMRoleV3Read,

		Schema: map[string]*schema.Schema{
			"projects": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"domains": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"project_domains": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"others": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceIAMRoleV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud client: %s", err)
	}

	domainID, err := getDomainID(config, client)
	if err != nil {
		return fmt.Errorf("Error getting the domain id, err=%s", err)

	}

	roles, err := dsAllRolesOfDomain(domainID, client)
	if err != nil {
		return err
	}
	if roles != nil {
		d.Set("projects", roles["XA"])
		d.Set("domains", roles["AX"])
		d.Set("project_domains", roles["AA"])
		d.Set("others", roles["XX"])
	}
	d.SetId("roles")
	return nil
}

func dsListRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]map[string]string, error) {
	old := client.Endpoint
	defer func() { client.Endpoint = old }()
	client.Endpoint = "https://iam.myhwclouds.com:443/v3/"

	opts := sdkroles.ListOpts{
		DomainID: domainID,
	}
	allPages, err := sdkroles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("List roles failed, err=%s", err)
	}

	all, err := sdkroles.ExtractRoles(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract roles failed, err=%s", err)
	}
	if len(all) == 0 {
		return nil, nil
	}

	r := map[string]map[string]string{
		"AX": make(map[string]string, 0),
		"XA": make(map[string]string, 0),
		"AA": make(map[string]string, 0),
		"XX": make(map[string]string, 0),
	}
	for _, item := range all {
		rtype, ok := item.Extra["type"].(string)
		if !ok {
			log.Printf("[DEBUG] Can not retrieve type of role:%#v", item)
			continue
		}

		dn, ok := item.Extra["display_name"].(string)
		if !ok {
			log.Printf("[DEBUG] Can not retrieve name ofrole:%#v", item)
			continue
		}

		desc, ok := item.Extra["description"].(string)
		if !ok {
			log.Printf("[DEBUG] Can not retrieve description of role:%#v", item)
			continue
		}

		r[rtype][dn] = desc
	}
	return r, nil
}

func dsAllRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]map[string]string, error) {
	roles, err := dsListRolesOfDomain("", client)
	if err != nil {
		return nil, fmt.Errorf("Error listing global roles, err=%s", err)
	}

	customRoles, err := dsListRolesOfDomain(domainID, client)
	if err != nil {
		return nil, fmt.Errorf("Error listing domain's custom roles, err=%s", err)
	}

	if customRoles == nil {
		return roles, nil
	}
	if roles == nil {
		return customRoles, nil
	}
	for k, v := range customRoles {
		for k1, v1 := range v {
			roles[k][k1] = v1
		}
	}
	return roles, nil
}
