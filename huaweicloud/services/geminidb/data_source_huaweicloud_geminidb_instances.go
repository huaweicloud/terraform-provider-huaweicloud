package geminidb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

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
			"datastore_type": {
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
						"product_type": {
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
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
									"whole_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated": {
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
							Elem: &schema.Resource{
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
							},
						},
						"pay_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintenance_window": {
							Type:     schema.TypeString,
							Computed: true,
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
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
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
										},
									},
									"nodes": {
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
												"role": {
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
										},
									},
								},
							},
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
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dedicated_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_encryption_id": {
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
						"dr_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dual_active_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_instance_node_num": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_instance_spec_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ccm_cert_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cert_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildInstancesQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("datastore_type"); ok {
		queryParams = fmt.Sprintf("%s&datastore_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("mode"); ok {
		queryParams = fmt.Sprintf("%s&mode=%v", queryParams, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		queryParams = fmt.Sprintf("%s&vpc_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		queryParams = fmt.Sprintf("%s&subnet_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildInstancesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the GeminiDB instances: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		instances := utils.PathSearch("instances", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		if len(instances) < limit {
			break
		}

		offset += len(instances)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instances", flattenInstances(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstances(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"port":               utils.PathSearch("port", v, nil),
			"mode":               utils.PathSearch("mode", v, nil),
			"product_type":       utils.PathSearch("product_type", v, nil),
			"region":             utils.PathSearch("region", v, nil),
			"datastore":          flattenInstanceDatastoreInfo(utils.PathSearch("datastore", v, nil)),
			"engine":             utils.PathSearch("engine", v, nil),
			"created":            utils.PathSearch("created", v, nil),
			"updated":            utils.PathSearch("updated", v, nil),
			"db_user_name":       utils.PathSearch("db_user_name", v, nil),
			"vpc_id":             utils.PathSearch("vpc_id", v, nil),
			"subnet_id":          utils.PathSearch("subnet_id", v, nil),
			"security_group_id":  utils.PathSearch("security_group_id", v, nil),
			"backup_strategy":    flattenInstanceBackupStrategyInfo(utils.PathSearch("backup_strategy", v, nil)),
			"pay_mode":           utils.PathSearch("pay_mode", v, nil),
			"maintenance_window": utils.PathSearch("maintenance_window", v, nil),
			"groups": flattenInstanceGroupsInfo(
				utils.PathSearch("groups", v, make([]interface{}, 0)).([]interface{})),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"time_zone":             utils.PathSearch("time_zone", v, nil),
			"actions":               utils.PathSearch("actions", v, nil),
			"dedicated_resource_id": utils.PathSearch("dedicated_resource_id", v, nil),
			"disk_encryption_id":    utils.PathSearch("disk_encryption_id", v, nil),
			"lb_ip_address":         utils.PathSearch("lb_ip_address", v, nil),
			"lb_port":               utils.PathSearch("lb_port", v, nil),
			"availability_zone":     utils.PathSearch("availability_zone", v, nil),
			"dr_instance_id":        utils.PathSearch("dr_instance_id", v, nil),
			"dual_active_info":      flattenInstanceDaulActiveInfo(utils.PathSearch("dual_active_info", v, nil)),
			"ccm_cert_info":         flattenInstanceCertInfo(utils.PathSearch("ccm_cert_info", v, nil)),
		})
	}

	return result
}

func flattenInstanceDatastoreInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"type":            utils.PathSearch("type", resp, nil),
		"version":         utils.PathSearch("version", resp, nil),
		"patch_available": utils.PathSearch("patch_available", resp, nil),
		"whole_version":   utils.PathSearch("whole_version", resp, nil),
	}

	return []interface{}{result}
}

func flattenInstanceBackupStrategyInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"start_time": utils.PathSearch("start_time", resp, nil),
		"keep_days":  utils.PathSearch("keep_days", resp, nil),
	}

	return []interface{}{result}
}

func flattenInstanceGroupsInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":     utils.PathSearch("id", v, nil),
			"status": utils.PathSearch("status", v, nil),
			"volume": flattenInstanceVolumeInfo(utils.PathSearch("volume", v, nil)),
			"nodes": flattenInstanceNodesInfo(
				utils.PathSearch("nodes", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenInstanceVolumeInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"size": utils.PathSearch("size", resp, nil),
		"used": utils.PathSearch("used", resp, nil),
	}

	return []interface{}{result}
}

func flattenInstanceNodesInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"subnet_id":         utils.PathSearch("subnet_id", v, nil),
			"private_ip":        utils.PathSearch("private_ip", v, nil),
			"public_ip":         utils.PathSearch("public_ip", v, nil),
			"spec_code":         utils.PathSearch("spec_code", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
			"support_reduce":    utils.PathSearch("support_reduce", v, nil),
		})
	}

	return result
}

func flattenInstanceDaulActiveInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"role":                           utils.PathSearch("role", resp, nil),
		"status":                         utils.PathSearch("status", resp, nil),
		"destination_instance_id":        utils.PathSearch("destination_instance_id", resp, nil),
		"destination_region":             utils.PathSearch("destination_region", resp, nil),
		"destination_instance_name":      utils.PathSearch("destination_instance_name", resp, nil),
		"destination_instance_node_num":  utils.PathSearch("destination_instance_node_num", resp, nil),
		"destination_instance_spec_code": utils.PathSearch("destination_instance_spec_code", resp, nil),
	}

	return []interface{}{result}
}

func flattenInstanceCertInfo(resp interface{}) []interface{} {
	if resp == nil || len(resp.(map[string]interface{})) == 0 {
		return nil
	}

	result := map[string]interface{}{
		"cert_id":   utils.PathSearch("cert_id", resp, nil),
		"cert_type": utils.PathSearch("cert_type", resp, nil),
	}

	return []interface{}{result}
}
