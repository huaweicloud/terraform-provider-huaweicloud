package hss

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/host-management/groups
func DataSourceHostGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"risk_host_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"unprotect_host_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
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
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"risk_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unprotect_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceHostGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		epsId     = cfg.GetEnterpriseProjectID(d)
		groupName = d.Get("name").(string)
		product   = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	allHostGroups, err := queryHostGroupsByName(client, region, epsId, groupName)
	if err != nil {
		return diag.Errorf("error querying HSS host groups: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	targetGroups := filterHostGroups(allHostGroups, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenHostGroups(targetGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func convertNumValueToString(num interface{}) string {
	if num == nil {
		return ""
	}

	if v, ok := num.(float64); ok {
		return fmt.Sprintf("%v", v)
	}

	return ""
}

func filterHostGroups(hostGroups []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(hostGroups))
	for _, v := range hostGroups {
		if groupID, ok := d.GetOk("group_id"); ok &&
			fmt.Sprint(groupID) != utils.PathSearch("group_id", v, "").(string) {
			continue
		}

		if hostNum, ok := d.GetOk("host_num"); ok &&
			hostNum.(string) != convertNumValueToString(utils.PathSearch("host_num", v, nil)) {
			continue
		}

		if riskHostNum, ok := d.GetOk("risk_host_num"); ok &&
			riskHostNum.(string) != convertNumValueToString(utils.PathSearch("risk_host_num", v, nil)) {
			continue
		}

		if unprotectHostNum, ok := d.GetOk("unprotect_host_num"); ok &&
			unprotectHostNum.(string) != convertNumValueToString(utils.PathSearch("unprotect_host_num", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenHostGroups(hostGroups []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(hostGroups))
	for _, v := range hostGroups {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("group_id", v, nil),
			"name":               utils.PathSearch("group_name", v, nil),
			"host_num":           utils.PathSearch("host_num", v, nil),
			"risk_host_num":      utils.PathSearch("risk_host_num", v, nil),
			"unprotect_host_num": utils.PathSearch("unprotect_host_num", v, nil),
			"host_ids":           utils.ExpandToStringList(utils.PathSearch("host_id_list", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}
