package geminidb

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

// @API GaussDBforNoSQL GET /v3/{project_id}/instances
func DataSourceGaussDBInfluxInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInfluxInstancesRead,

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
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     gaussDBInfluxInstanceDatastoreSchema(),
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_strategy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     gaussDBInfluxInstanceBackupStrategySchema(),
						},
						"pay_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_begin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintain_end": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     gaussDBInfluxInstanceGroupSchema(),
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dedicated_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lb_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lb_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
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

func gaussDBInfluxInstanceDatastoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"patch_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func gaussDBInfluxInstanceBackupStrategySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keep_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func gaussDBInfluxInstanceGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBInfluxInstanceGroupVolumeSchema(),
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBInfluxInstanceGroupNodesSchema(),
			},
		},
	}
	return &sc
}

func gaussDBInfluxInstanceGroupVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func gaussDBInfluxInstanceGroupNodesSchema() *schema.Resource {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_reduce": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceGaussDBInfluxInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	queryParams := buildGetGaussDBInfluxInstancesQueryParams(d, "influxdb")

	return dataSourceGaussDBNoSQLInstancesRead(d, meta, queryParams)
}

func buildGetGaussDBInfluxInstancesQueryParams(d *schema.ResourceData, datastoreType string) string {
	res := fmt.Sprintf("?datastore_type=%v", datastoreType)
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("mode"); ok {
		res = fmt.Sprintf("%s&mode=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&subnet_id=%v", res, v)
	}
	return res
}

func dataSourceGaussDBNoSQLInstancesRead(d *schema.ResourceData, meta interface{}, queryParams string) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances"
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += queryParams

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB influx instances")
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instances", flattenGetGaussDBNoSQLInstances(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussDBNoSQLInstances(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		instance := map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"port":                  utils.PathSearch("port", v, nil),
			"mode":                  utils.PathSearch("mode", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"datastore":             flattenInstanceDatastore(v),
			"engine":                utils.PathSearch("engine", v, nil),
			"db_user_name":          utils.PathSearch("db_user_name", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"security_group_id":     utils.PathSearch("security_group_id", v, nil),
			"backup_strategy":       flattenInstanceBackupStrategy(v),
			"pay_mode":              utils.PathSearch("pay_mode", v, nil),
			"groups":                flattenInstanceGroups(v),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"time_zone":             utils.PathSearch("time_zone", v, nil),
			"actions":               utils.PathSearch("actions", v, nil),
			"dedicated_resource_id": utils.PathSearch("dedicated_resource_id", v, nil),
			"lb_ip_address":         utils.PathSearch("lb_ip_address", v, nil),
			"lb_port":               utils.PathSearch("lb_port", v, nil),
			"availability_zone":     utils.PathSearch("availability_zone", v, nil),
			"created_at":            utils.PathSearch("created", v, nil),
			"updated_at":            utils.PathSearch("updated", v, nil),
		}
		maintainWindow := strings.Split(utils.PathSearch("maintenance_window", v, "").(string), "-")
		if len(maintainWindow) == 2 {
			instance["maintain_begin"] = maintainWindow[0]
			instance["maintain_end"] = maintainWindow[1]
		}
		rst = append(rst, instance)
	}
	return rst
}

func flattenInstanceDatastore(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"engine":          utils.PathSearch("datastore.type", resp, nil),
		"version":         utils.PathSearch("datastore.version", resp, nil),
		"patch_available": utils.PathSearch("datastore.patch_available", resp, nil),
	})
	return rst
}

func flattenInstanceBackupStrategy(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"start_time": utils.PathSearch("backup_strategy.start_time", resp, nil),
		"keep_days":  utils.PathSearch("backup_strategy.keep_days", resp, nil),
	})
	return rst
}

func flattenInstanceGroups(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("id", v, nil),
			"status": utils.PathSearch("status", v, nil),
			"volume": flattenInstanceGroupVolume(v),
			"nodes":  flattenInstanceGroupNodes(v),
		})
	}
	return rst
}

func flattenInstanceGroupVolume(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"size": utils.PathSearch("volume.size", resp, nil),
		"used": utils.PathSearch("volume.used", resp, nil),
	})
	return rst
}

func flattenInstanceGroupNodes(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"subnet_id":         utils.PathSearch("subnet_id", v, nil),
			"private_ip":        utils.PathSearch("private_ip", v, nil),
			"public_ip":         utils.PathSearch("public_ip", v, nil),
			"spec_code":         utils.PathSearch("spec_code", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"support_reduce":    utils.PathSearch("support_reduce", v, nil),
		})
	}
	return rst
}
