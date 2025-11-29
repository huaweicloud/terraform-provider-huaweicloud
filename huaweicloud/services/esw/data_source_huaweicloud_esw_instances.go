package esw

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ESW GET /v3/{project_id}/l2cg/instances
func DataSourceEswInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEswInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
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
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     availabilityZonesSchema(),
						},
						"tunnel_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     tunnelInfoSchema(),
						},
						"charge_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     chargeInfosSchema(),
						},
					},
				},
			},
		},
	}
}

func availabilityZonesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"primary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func tunnelInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func chargeInfosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceEswInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/l2cg/instances"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listQueryParams := buildListQueryParams(d)
	listPath += listQueryParams

	listResp, err := pagination.ListAllItems(
		client,
		"marker",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving ESW instances: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenEswInstancesBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenEswInstancesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"project_id":         utils.PathSearch("project_id", v, nil),
			"region":             utils.PathSearch("region", v, nil),
			"flavor_ref":         utils.PathSearch("flavor_ref", v, nil),
			"ha_mode":            utils.PathSearch("ha_mode", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"created_at":         utils.PathSearch("created_at", v, nil),
			"updated_at":         utils.PathSearch("updated_at", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"availability_zones": flattenInstancesAvailabilityZones(v),
			"tunnel_info":        flattenInstancesTunnelInfo(v),
			"charge_infos":       flattenInstancesChargeInfos(v),
		})
	}
	return rst
}

func flattenInstancesAvailabilityZones(resp interface{}) []interface{} {
	curJson := utils.PathSearch("availability_zones", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"primary": utils.PathSearch("primary", curJson, nil),
			"standby": utils.PathSearch("standby", curJson, nil),
		},
	}
	return rst
}

func flattenInstancesTunnelInfo(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tunnel_info", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"vpc_id":       utils.PathSearch("vpc_id", curJson, nil),
			"virsubnet_id": utils.PathSearch("virsubnet_id", curJson, nil),
			"tunnel_ip":    utils.PathSearch("tunnel_ip", curJson, nil),
			"tunnel_port":  int(utils.PathSearch("tunnel_port", curJson, float64(0)).(float64)),
			"tunnel_type":  utils.PathSearch("tunnel_type", curJson, nil),
		},
	}
	return rst
}

func flattenInstancesChargeInfos(resp interface{}) []interface{} {
	curJson := utils.PathSearch("charge_infos", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"charge_mode": utils.PathSearch("charge_mode", curJson, nil),
		},
	}
	return rst
}
