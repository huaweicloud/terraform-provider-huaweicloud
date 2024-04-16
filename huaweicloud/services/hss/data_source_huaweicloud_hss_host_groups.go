package hss

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

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
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		name   = d.Get("name").(string)
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	allHostGroups, err := queryHostGroups(client, region, epsId, name)
	if err != nil {
		return diag.Errorf("error querying host groups: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)
	targetGroups := filterHostGroups(allHostGroups, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenHostGroups(targetGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func int32ToString(i *int32) string {
	if i == nil {
		return ""
	}

	return fmt.Sprintf("%d", *i)
}

func filterHostGroups(groups []hssv5model.HostGroupItem, d *schema.ResourceData) []hssv5model.HostGroupItem {
	if len(groups) == 0 {
		return nil
	}

	rst := make([]hssv5model.HostGroupItem, 0, len(groups))
	for _, v := range groups {
		if groupID, ok := d.GetOk("group_id"); ok &&
			fmt.Sprint(groupID) != utils.StringValue(v.GroupId) {
			continue
		}

		if hostNum, ok := d.GetOk("host_num"); ok &&
			hostNum.(string) != int32ToString(v.HostNum) {
			continue
		}

		if riskHostNum, ok := d.GetOk("risk_host_num"); ok &&
			riskHostNum.(string) != int32ToString(v.RiskHostNum) {
			continue
		}

		if unprotectHostNum, ok := d.GetOk("unprotect_host_num"); ok &&
			unprotectHostNum.(string) != int32ToString(v.UnprotectHostNum) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenHostGroups(groups []hssv5model.HostGroupItem) []interface{} {
	if len(groups) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(groups))
	for _, v := range groups {
		rst = append(rst, map[string]interface{}{
			"id":                 v.GroupId,
			"name":               v.GroupName,
			"host_num":           v.HostNum,
			"risk_host_num":      v.RiskHostNum,
			"unprotect_host_num": v.UnprotectHostNum,
			"host_ids":           v.HostIdList,
		})
	}

	return rst
}
