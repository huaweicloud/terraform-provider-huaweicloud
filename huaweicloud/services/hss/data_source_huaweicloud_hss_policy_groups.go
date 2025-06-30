package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/policy/groups
func DataSourcePolicyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyGroupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the policy group ID.",
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the policy group name.",
			},
			// The filter parameter does not take effect.​
			"container_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to query the container edition policy.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the policy group.",
						},
						"deletable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether a policy group can be deleted.",
						},
						"host_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of associated servers.",
						},
						"default_group": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether a policy group is the default policy group.",
						},
						"support_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The supported OS.",
						},
						"support_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The supported versions.",
						},
					},
				},
				Description: "The policy group list.",
			},
		},
	}
}

func dataSourcePolicyGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/policy/groups"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildPolicyGroupsQueryParams(d, cfg)
	allPolicyGroups := make([]interface{}, 0)
	offset := 0

	listPolicyGroupsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		// When using this API on the European site, not adding the `region` parameter in the headers will result in an
		// error. After consulting the HSS service, it is confirmed that the `region` needs to be added here.
		MoreHeaders: map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listPolicyGroupsOpt)

		if err != nil {
			return diag.Errorf("error retrieving HSS policy groups: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		policyGroupsResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})

		if len(policyGroupsResp) == 0 {
			break
		}
		allPolicyGroups = append(allPolicyGroups, policyGroupsResp...)
		offset += len(policyGroupsResp)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenPolicyGroups(allPolicyGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPolicyGroupsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=10"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("group_name"); ok {
		queryParams = fmt.Sprintf("%s&group_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_mode"); ok {
		queryParams = fmt.Sprintf("%s&container_mode=%v", queryParams, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}
	return queryParams
}

func flattenPolicyGroups(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"group_name":      utils.PathSearch("group_name", v, nil),
			"group_id":        utils.PathSearch("group_id", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"deletable":       utils.PathSearch("deletable", v, nil),
			"host_num":        utils.PathSearch("host_num", v, nil),
			"default_group":   utils.PathSearch("default_group", v, nil),
			"support_os":      utils.PathSearch("support_os", v, nil),
			"support_version": utils.PathSearch("support_version", v, nil),
		})
	}
	return rst
}
