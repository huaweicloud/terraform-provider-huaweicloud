// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDS
// ---------------------------------------------------------------

package dds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS GET /v3/{project_id}/instances
func DataSourceDdsInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdsInstanceRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the DB instance name.`,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Sharding", "ReplicaSet", "Single",
				}, true),
				Description: `Specifies the mode of the database instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the subnet Network ID.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        ddsInstanceInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of DDS instances.`,
			},
		},
	}
}

func ddsInstanceInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB instance name.`,
			},
			"ssl": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether to enable or disable SSL.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the database port number. The port range is 2100 to 9500.`,
			},
			"datastore": {
				Type:     schema.TypeList,
				Elem:     ddsInstanceInstanceDatastoreSchema(),
				Computed: true,
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Elem:     ddsInstanceInstanceBackupStrategySchema(),
				Computed: true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VPC ID`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the subnet Network ID.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the security group ID of the DDS instance.`,
			},
			"disk_encryption_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the disk encryption ID of the instance.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the mode of the database instance.`,
			},
			"db_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB Administator name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the the DB instance status.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the enterprise project id of the dds instance.`,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:        schema.TypeList,
							Elem:        ddsInstanceInstanceNodeSchema(),
							Computed:    true,
							Description: `Indicates the instance nodes information.`,
						},
					},
				},
			},
			"tags": common.TagsComputedSchema(),

			// deprecated
			"nodes": {
				Type:        schema.TypeList,
				Elem:        ddsInstanceInstanceNodeSchema(),
				Computed:    true,
				Description: `This field is deprecated.`,
			},
		},
	}
	return &sc
}

func ddsInstanceInstanceDatastoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB engine.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB instance version.`,
			},
			"storage_engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage engine of the DB instance.`,
			},
		},
	}
	return &sc
}

func ddsInstanceInstanceBackupStrategySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup time window.`,
			},
			"keep_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of days to retain the generated backup files.`,
			},
		},
	}
	return &sc
}

func ddsInstanceInstanceNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node name.`,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node role.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private IP address of a node.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the EIP that has been bound on a node.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node status.`,
			},
			"spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the node spec code.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the availability zone.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "schema: Deprecated",
			},
		},
	}
	return &sc
}

func resourceDdsInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getDDSInstances: Query the List of DDS instances.
	var (
		getDDSInstancesHttpUrl = "v3/{project_id}/instances"
		getDDSInstancesProduct = "dds"
	)
	getDDSInstancesClient, err := conf.NewServiceClient(getDDSInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DdsInstance Client: %s", err)
	}

	getDDSInstancesPath := getDDSInstancesClient.Endpoint + getDDSInstancesHttpUrl
	getDDSInstancesPath = strings.ReplaceAll(getDDSInstancesPath, "{project_id}", getDDSInstancesClient.ProjectID)

	getDDSInstancesQueryParams := buildGetDDSInstancesQueryParams(d)
	getDDSInstancesPath += getDDSInstancesQueryParams

	getDDSInstancesResp, err := pagination.ListAllItems(
		getDDSInstancesClient,
		"offset",
		getDDSInstancesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdsInstance")
	}

	getDDSInstancesRespJson, err := json.Marshal(getDDSInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDDSInstancesRespBody interface{}
	err = json.Unmarshal(getDDSInstancesRespJson, &getDDSInstancesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", flattenGetDDSInstancesResponseBodyInstance(getDDSInstancesRespBody, getDDSInstancesClient)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDDSInstancesResponseBodyInstance(resp interface{}, client *golangsdk.ServiceClient) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		id := utils.PathSearch("id", v, nil)

		sslEnable := true
		if utils.PathSearch("ssl", v, nil) == 0 {
			sslEnable = false
		}
		portStr := utils.PathSearch("port", v, nil)
		port, err := strconv.Atoi(portStr.(string))
		if err != nil {
			log.Printf("[WARNING] Port %s invalid, Type conversion error: %s", portStr, err)
		}

		// save tags
		var tagMap interface{}
		if resourceTags, err := tags.Get(client, "instances", id.(string)).Extract(); err == nil {
			tagMap = utils.TagsToMap(resourceTags.Tags)
		} else {
			log.Printf("[WARN] Error fetching tags of DDS instance (%s): %s", id.(string), err)
		}

		rst = append(rst, map[string]interface{}{
			"id":                    id,
			"name":                  utils.PathSearch("name", v, nil),
			"ssl":                   sslEnable,
			"port":                  port,
			"datastore":             flattenInstanceDatastore(v),
			"backup_strategy":       flattenInstanceBackupStrategy(v),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"security_group_id":     utils.PathSearch("security_group_id", v, nil),
			"disk_encryption_id":    utils.PathSearch("disk_encryption_id", v, nil),
			"mode":                  utils.PathSearch("mode", v, nil),
			"db_username":           utils.PathSearch("db_username", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"nodes":                 flattenInstanceNodes(v),
			"groups":                flattenInstanceGroups(v),
			"tags":                  tagMap,
		})
	}
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
			"type":   utils.PathSearch("type", v, nil),
			"id":     utils.PathSearch("id", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"status": utils.PathSearch("status", v, nil),
			"size":   utils.PathSearch("volume.size", v, nil),
			"used":   utils.PathSearch("volume.used", v, nil),
			"nodes":  flattenInstanceNodes(v),
		})
	}
	return rst
}

func flattenInstanceNodes(resp interface{}) []interface{} {
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
			"role":              utils.PathSearch("role", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"private_ip":        utils.PathSearch("private_ip", v, nil),
			"public_ip":         utils.PathSearch("public_ip", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"spec_code":         utils.PathSearch("spec_code", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
		})
	}
	return rst
}

func flattenInstanceDatastore(resp interface{}) interface{} {
	var rst []map[string]interface{}
	curJson := utils.PathSearch("datastore", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []map[string]interface{}{
		{
			"type":           utils.PathSearch("type", curJson, nil),
			"version":        utils.PathSearch("version", curJson, nil),
			"storage_engine": utils.PathSearch("storage_engine", curJson, nil),
		},
	}
	return rst
}

func flattenInstanceBackupStrategy(resp interface{}) interface{} {
	var rst []map[string]interface{}
	curJson := utils.PathSearch("backup_strategy", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []map[string]interface{}{
		{
			"start_time": utils.PathSearch("start_time", curJson, nil),
			"keep_days":  utils.PathSearch("keep_days", curJson, nil),
		},
	}
	return rst
}

func buildGetDDSInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
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

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
