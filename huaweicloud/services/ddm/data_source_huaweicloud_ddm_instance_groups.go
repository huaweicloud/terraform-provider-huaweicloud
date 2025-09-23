package ddm

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

// @API DDM GET /v3/{project_id}/instances/{instance_id}/groups
func DataSourceDdmInstanceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmInstanceGroupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_list": {
				Type:     schema.TypeList,
				Elem:     InstanceGroupsGroupSchema(),
				Computed: true,
			},
		},
	}
}

func InstanceGroupsGroupSchema() *schema.Resource {
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
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_load_balance": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cpu_num_per_node": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_num_per_node": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Elem:     InstanceGroupsGroupNodeSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func InstanceGroupsGroupNodeSchema() *schema.Resource {
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
			"az": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDdmInstanceGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/groups"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	var res []interface{}
	offset := 0
	for {
		getPath := getBasePath + buildGetDdmInstanceGroupsQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving DDM instance groups")
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		groups := flattenGetInstanceGroupsResponseBodyGroup(getRespBody)
		res = append(res, groups...)
		totalAccount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if len(res) >= int(totalAccount) {
			break
		}
		offset++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_list", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
func buildGetDdmInstanceGroupsQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%d", offset)
}

func flattenGetInstanceGroupsResponseBodyGroup(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("group_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"role":             utils.PathSearch("role", v, nil),
			"endpoint":         utils.PathSearch("endpoint", v, nil),
			"ipv6_endpoint":    utils.PathSearch("ipv6_endpoint", v, nil),
			"is_load_balance":  utils.PathSearch("is_load_balance", v, nil),
			"is_default_group": utils.PathSearch("is_default_group", v, nil),
			"cpu_num_per_node": utils.PathSearch("cpu_num_per_node", v, nil),
			"mem_num_per_node": utils.PathSearch("mem_num_per_node", v, nil),
			"architecture":     utils.PathSearch("architecture", v, nil),
			"node_list":        flattenGetInstanceGroupsResponseBodyGroupNode(v),
		})
	}
	return rst
}

func flattenGetInstanceGroupsResponseBodyGroupNode(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("node_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
			"az":   utils.PathSearch("az", v, nil),
		})
	}
	return rst
}
