package rds

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

// @API RDS GET /v3/{project_id}/instances
func DataSourceRdsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstancesRead,

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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_type": {
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ha_replication_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"param_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"time_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db": {
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
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"user_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"volume": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_encryption_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
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
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"private_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"public_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "rds"
		httpUrl = "v3/{project_id}/instances"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListInstancesQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving RDS instances: %s", err)
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", flattenListInstancesBody(listRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("datastore_type"); ok {
		res = fmt.Sprintf("%s&datastore_type=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&subnet_id=%v", res, v)
	}
	if v, ok := d.GetOk("group_type"); ok {
		res = fmt.Sprintf("%s&group_type=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListInstancesBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}

	enterpriseProjectId, enterpriseProjectIdOk := d.GetOk("enterprise_project_id")
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0)
	for _, v := range curArray {
		enterpriseProjectIdRaw := utils.PathSearch("enterprise_project_id", v, nil)
		if enterpriseProjectIdOk && enterpriseProjectId != enterpriseProjectIdRaw {
			continue
		}
		instance := map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"region":                utils.PathSearch("region", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"created":               utils.PathSearch("created", v, nil),
			"ha_replication_mode":   utils.PathSearch("ha.replication_mode", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"security_group_id":     utils.PathSearch("security_group_id", v, nil),
			"flavor":                utils.PathSearch("flavor_ref", v, nil),
			"time_zone":             utils.PathSearch("time_zone", v, nil),
			"ssl_enable":            utils.PathSearch("enable_ssl", v, nil),
			"enterprise_project_id": enterpriseProjectIdRaw,
			"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", v, make([]interface{}, 0))),
			"public_ips":            utils.PathSearch("public_ips", v, nil),
			"private_ips":           utils.PathSearch("private_ips", v, nil),
			"fixed_ip":              utils.PathSearch("private_ips[0]", v, nil),
			"volume":                flattenInstancesVolume(v),
			"db":                    flattenInstancesDb(v),
			"backup_strategy":       flattenInstancesBackupStrategy(v),
			"nodes":                 flattenInstancesNodes(v),
		}
		instanceType := utils.PathSearch("type", v, "").(string)
		var az []string
		if instanceType == "Ha" {
			masterNodeAz := utils.PathSearch("nodes[?role=='master']|[0].availability_zone", v, "").(string)
			slaveNodeAz := utils.PathSearch("nodes[?role=='slave']|[0].availability_zone", v, "").(string)
			az = []string{masterNodeAz, slaveNodeAz}
		} else if instanceType == "Single" || instanceType == "Replica" {
			az = []string{utils.PathSearch("nodes[0].availability_zone", v, "").(string)}
		}
		instance["availability_zone"] = az

		rst = append(rst, instance)
	}
	return rst
}

func flattenInstancesVolume(instance interface{}) []interface{} {
	volume := map[string]interface{}{
		"type":               utils.PathSearch("volume.type", instance, nil),
		"size":               utils.PathSearch("volume.size", instance, nil),
		"disk_encryption_id": utils.PathSearch("disk_encryption_id", instance, nil),
	}

	return []interface{}{volume}
}

func flattenInstancesDb(instance interface{}) []interface{} {
	database := map[string]interface{}{
		"type":      utils.PathSearch("datastore.type", instance, nil),
		"version":   utils.PathSearch("datastore.version", instance, nil),
		"port":      utils.PathSearch("port", instance, nil),
		"user_name": utils.PathSearch("db_user_name", instance, nil),
	}
	return []interface{}{database}
}

func flattenInstancesBackupStrategy(instance interface{}) []interface{} {
	backup := map[string]interface{}{
		"start_time": utils.PathSearch("backup_strategy.start_time", instance, nil),
		"keep_days":  utils.PathSearch("backup_strategy.keep_days", instance, nil),
	}
	return []interface{}{backup}
}

func flattenInstancesNodes(instance interface{}) []interface{} {
	nodesJson := utils.PathSearch("nodes", instance, make([]interface{}, 0))
	nodeArray := nodesJson.([]interface{})
	if len(nodeArray) < 1 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodeArray))
	for _, v := range nodeArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
		})
	}
	return rst
}
