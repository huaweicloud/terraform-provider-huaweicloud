package vpc

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	v1groups "github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	v3groups "github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

type v1Group = v1groups.SecurityGroup

type v3Group = v3groups.SecurityGroup

// @API VPC GET /v3/{project_id}/vpc/security-groups
// @API VPC GET /v1/{project_id}/security-groups
func DataSourceNetworkingSecGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenSecGroupDetail(secGroup interface{}) map[string]interface{} {
	var result map[string]interface{}
	switch entity := secGroup.(type) {
	case v1Group:
		result = map[string]interface{}{
			"id":                    entity.ID,
			"name":                  entity.Name,
			"enterprise_project_id": entity.EnterpriseProjectId,
			"description":           entity.Description,
		}
	case v3Group:
		result = map[string]interface{}{
			"id":                    entity.ID,
			"name":                  entity.Name,
			"enterprise_project_id": entity.EnterpriseProjectId,
			"description":           entity.Description,
			"created_at":            entity.CreatedAt,
			"updated_at":            entity.UpdatedAt,
		}
	}
	return result
}

func filterAvailableSecGroupsV1(secGroups []v1groups.SecurityGroup, name, descKey string) ([]map[string]interface{},
	[]string) {
	var secGroupsCopy []v1groups.SecurityGroup = secGroups

	// Build filter by name and description content.
	if name != "" {
		tmp := make([]v1groups.SecurityGroup, 0, len(secGroupsCopy))
		for _, secgroup := range secGroups {
			if name != secgroup.Name {
				continue
			}
			tmp = append(tmp, secgroup)
		}
		secGroupsCopy = tmp
	}
	if descKey != "" {
		tmp := make([]v1groups.SecurityGroup, 0, len(secGroupsCopy))
		for _, secgroup := range secGroups {
			if !strings.Contains(secgroup.Description, descKey) {
				continue
			}
			tmp = append(tmp, secgroup)
		}
		secGroupsCopy = tmp
	}

	result := make([]map[string]interface{}, len(secGroupsCopy))
	ids := make([]string, len(secGroupsCopy))
	for i, secGroup := range secGroupsCopy {
		ids[i] = secGroup.ID
		result[i] = flattenSecGroupDetail(secGroup)
	}

	return result, ids
}

func filterAvailableSecGroupsV3(secGroups []v3groups.SecurityGroup, descKey string) ([]map[string]interface{},
	[]string) {
	var secGroupsCopy []v3groups.SecurityGroup = secGroups

	// Filter all security groups with keywords in their description.
	if descKey != "" {
		tmp := make([]v3groups.SecurityGroup, 0, len(secGroupsCopy))
		for _, secgroup := range secGroups {
			if !strings.Contains(secgroup.Description, descKey) {
				continue
			}
			tmp = append(tmp, secgroup)
		}
		secGroupsCopy = tmp
	}

	result := make([]map[string]interface{}, len(secGroupsCopy))
	ids := make([]string, len(secGroupsCopy))
	for i, secGroup := range secGroupsCopy {
		ids[i] = secGroup.ID
		result[i] = flattenSecGroupDetail(secGroup)
	}

	return result, ids
}

func dataSourceNetworkingSecGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v3Client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return diag.Errorf("error creating networking v3 client: %s", err)
	}

	// The List method currently does not support filtering by keyword in the description. Therefore, keyword filtering
	// is implemented by manually filtering the description value of the List method return.
	listOpts := v3groups.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}

	var groupList []map[string]interface{}
	var ids []string
	allSecGroups, err := v3groups.List(v3Client, listOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// If the v3 API does not exist or has not been published in the specified region, set again using v1 API.
			return dataSourceNetworkingSecGroupsReadV1(ctx, d, meta)
		} else {
			return diag.Errorf("error getting security groups list: %s", err)
		}
	} else {
		groupList, ids = filterAvailableSecGroupsV3(allSecGroups, d.Get("description").(string))
	}
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("security_groups", groupList),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func dataSourceNetworkingSecGroupsReadV1(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error getting security groups: %s", err)
	}
	allSecGroups, err := v1groups.ExtractSecurityGroups(pages)
	if err != nil {
		return diag.Errorf("error retrieving security groups list: %s", err)
	}
	log.Printf("[DEBUG] Retrieved Security Groups: %+v", allSecGroups)

	groupList, ids := filterAvailableSecGroupsV1(allSecGroups, d.Get("name").(string), d.Get("description").(string))
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("security_groups", groupList),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
