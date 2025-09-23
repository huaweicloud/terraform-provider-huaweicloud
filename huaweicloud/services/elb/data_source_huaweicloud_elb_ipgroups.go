package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/ipgroups
func DataSourceElbIpGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbIpGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipgroups": {
				Type:     schema.TypeList,
				Elem:     ipgroupsSchema(),
				Computed: true,
			},
		},
	}
}

func ipgroupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
	}
	return &sc
}

func dataSourceElbIpGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listIpGroupsHttpUrl = "v3/{project_id}/elb/ipgroups"
		listIpGroupsProduct = "elb"
	)
	listIpGroupsClient, err := cfg.NewServiceClient(listIpGroupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listIpGroupsPath := listIpGroupsClient.Endpoint + listIpGroupsHttpUrl
	listIpGroupsPath = strings.ReplaceAll(listIpGroupsPath, "{project_id}", listIpGroupsClient.ProjectID)
	listIpGroupsQueryParams := buildListIpGroupsQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listIpGroupsPath += listIpGroupsQueryParams
	listIpGroupsResp, err := pagination.ListAllItems(
		listIpGroupsClient,
		"marker",
		listIpGroupsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB IP groups")
	}

	listIpGroupsRespJson, err := json.Marshal(listIpGroupsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listIpGroupsRespBody interface{}
	err = json.Unmarshal(listIpGroupsRespJson, &listIpGroupsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ipgroups", flattenListIpGroupsBody(listIpGroupsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListIpGroupsQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("ipgroup_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("ip_address"); ok {
		res = fmt.Sprintf("%s&ip_list=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListIpGroupsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("ipgroups", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"project_id":            utils.PathSearch("project_id", v, nil),
			"listeners":             utils.PathSearch("listeners", v, nil),
			"ip_list":               utils.PathSearch("ip_list", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"updated_at":            utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
