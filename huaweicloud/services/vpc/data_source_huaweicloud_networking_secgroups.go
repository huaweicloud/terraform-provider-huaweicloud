package vpc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceNetworkingSecGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
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

func getSecGroupDetail(secGroup groups.SecurityGroup) map[string]interface{} {
	return map[string]interface{}{
		"id":                    secGroup.ID,
		"name":                  secGroup.Name,
		"enterprise_project_id": secGroup.EnterpriseProjectId,
		"description":           secGroup.Description,
		"created_at":            secGroup.CreatedAt,
		"updated_at":            secGroup.UpdatedAt,
	}
}

func filterAvailableSecGroups(d *schema.ResourceData, secGroups []groups.SecurityGroup) ([]map[string]interface{},
	[]string) {
	secGroupCopy := secGroups

	// build filter by description content.
	if desc, ok := d.GetOk("description"); ok {
		tmp := make([]groups.SecurityGroup, 0, len(secGroupCopy))
		for _, secgroup := range secGroupCopy {
			if strings.Contains(secgroup.Description, desc.(string)) {
				// Filter all security groups with keywords in their description.
				tmp = append(tmp, secgroup)
			}
		}
		secGroupCopy = tmp
	}

	result := make([]map[string]interface{}, len(secGroupCopy))
	ids := make([]string, len(secGroupCopy))
	for i, secGroup := range secGroupCopy {
		ids[i] = secGroup.ID
		result[i] = getSecGroupDetail(secGroup)
	}

	return result, ids
}

func dataSourceNetworkingSecGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking client: %s", err)
	}

	// The List method currently does not support filtering by keyword in the description. Therefore, keyword filtering
	// is implemented by manually filtering the description value of the List method return.
	listOpts := groups.ListOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: config.DataGetEnterpriseProjectID(d),
	}
	allSecGroups, err := groups.List(networkingClient, listOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error getting security groups: %s", err)
	}
	logp.Printf("[DEBUG] Retrieved Security Groups: %+v", allSecGroups)

	resp, ids := filterAvailableSecGroups(d, allSecGroups)
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("security_groups", resp),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}
