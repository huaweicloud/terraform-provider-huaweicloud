// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DCS
// ---------------------------------------------------------------

package dcs

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

// @API DCS GET /v2/{project_id}/instances
func DataSourceDcsInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDcsInstanceRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of an instance.`,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RUNNING", "ERROR", "RESTARTING", "FROZEN", "EXTENDING", "RESTORING", "FLUSHING",
				}, false),
				Description: `Specifies the cache instance status.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the subnet Network ID.`,
			},
			"capacity": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the cache capacity. Unit: GB.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        InstanceInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of DCS instances.`,
			},
		},
	}
}

func InstanceInstanceSchema() *schema.Resource {
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
				Description: `Indicates the name of an instance.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates a cache engine.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of a cache engine.`,
			},
			"capacity": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the cache capacity. Unit: GB.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the flavor of the cache instance.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Specifies the code of the AZ where the cache node resides.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of VPC which the instance belongs to.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of VPC which the instance belongs to.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of subnet which the instance belongs to.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of subnet which the instance belongs to.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the security group which the instance belongs to.`,
			},
			"security_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of security group which the instance belongs to.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the enterprise project id of the dcs instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of an instance.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP address of the DCS instance.`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance time window starts.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance time window ends.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the charging mode of the cache instance.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the port of the cache instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache instance status.`,
			},
			"used_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the size of the used memory. Unit: MB.`,
			},
			"max_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total memory size. Unit: MB.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the domain name of the instance.`,
			},
			"access_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the username used for accessing a DCS Memcached instance.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the order that created the instance.`,
			},
			"tags": common.TagsComputedSchema(),
		},
	}
	return &sc
}

func resourceDcsInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDCSInstances: Query the List of DCS instances.
	var (
		getDCSInstancesHttpUrl = "v2/{project_id}/instances"
		getDCSInstancesProduct = "dcs"
	)
	getDCSInstancesClient, err := cfg.NewServiceClient(getDCSInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	getDCSInstancesPath := getDCSInstancesClient.Endpoint + getDCSInstancesHttpUrl
	getDCSInstancesPath = strings.ReplaceAll(getDCSInstancesPath, "{project_id}", getDCSInstancesClient.ProjectID)

	getDCSInstancesQueryParams := buildGetDCSInstancesQueryParams(d)
	getDCSInstancesPath += getDCSInstancesQueryParams

	getDCSInstancesResp, err := pagination.ListAllItems(
		getDCSInstancesClient,
		"offset",
		getDCSInstancesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DcsInstance")
	}

	getDCSInstancesRespJson, err := json.Marshal(getDCSInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDCSInstancesRespBody interface{}
	err = json.Unmarshal(getDCSInstancesRespJson, &getDCSInstancesRespBody)
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
		d.Set("instances", flattenGetDCSInstancesResponseBodyInstance(getDCSInstancesRespBody, getDCSInstancesClient)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDCSInstancesResponseBodyInstance(resp interface{}, client *golangsdk.ServiceClient) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		// capacity
		capacity := utils.PathSearch("capacity", v, nil)
		capacityMinor := utils.PathSearch("capacity", v, nil)
		if capacity == 0 {
			capacity, _ = strconv.ParseFloat(capacityMinor.(string), floatBitSize)
		}

		securityGroupID := utils.PathSearch("security_group_id", v, nil)
		// If security_group_id is not set, the default value is returned: securityGroupId. Change it to empty.
		if securityGroupID == "securityGroupId" {
			securityGroupID = ""
		}
		id := utils.PathSearch("instance_id", v, nil)
		// save tags
		var tagMap interface{}
		if resourceTags, err := tags.Get(client, "instances", id.(string)).Extract(); err == nil {
			tagMap = utils.TagsToMap(resourceTags.Tags)
		} else {
			log.Printf("[WARN] Error fetching tags of DCS instance (%s): %s", id.(string), err)
		}
		rst = append(rst, map[string]interface{}{
			"id":                    id,
			"name":                  utils.PathSearch("name", v, nil),
			"engine":                utils.PathSearch("engine", v, nil),
			"engine_version":        utils.PathSearch("engine_version", v, nil),
			"capacity":              capacity,
			"flavor":                utils.PathSearch("spec_code", v, nil),
			"availability_zones":    utils.PathSearch("availability_zones", v, nil),
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"vpc_name":              utils.PathSearch("vpc_name", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"subnet_name":           utils.PathSearch("subnet_name", v, nil),
			"security_group_id":     securityGroupID,
			"security_group_name":   utils.PathSearch("security_group_name", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"private_ip":            utils.PathSearch("ip", v, nil),
			"maintain_begin":        utils.PathSearch("maintain_begin", v, nil),
			"maintain_end":          utils.PathSearch("maintain_end", v, nil),
			"charging_mode":         chargingMode[int(utils.PathSearch("charging_mode", v, nil).(float64))],
			"port":                  utils.PathSearch("port", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"used_memory":           utils.PathSearch("used_memory", v, nil),
			"max_memory":            utils.PathSearch("max_memory", v, nil),
			"domain_name":           utils.PathSearch("domain_name", v, nil),
			"access_user":           utils.PathSearch("access_user", v, nil),
			"order_id":              utils.PathSearch("order_id", v, nil),
			"tags":                  tagMap,
		})
	}
	return rst
}

func buildGetDCSInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("private_ip"); ok {
		res = fmt.Sprintf("%s&ip=%v", res, v)
	}

	if v, ok := d.GetOk("capacity"); ok {
		res = fmt.Sprintf("%s&capacity=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
