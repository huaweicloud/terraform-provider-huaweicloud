// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
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

// @API DDM GET /v1/{project_id}/instances
func DataSourceDdmInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project id.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the engine version.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        InstancesInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM instance.`,
			},
		},
	}
}

func InstancesInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the DDM instance.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of availability zones.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a security group.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of nodes.`,
			},
			"access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address for accessing the DDM instance.`,
			},
			"access_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the port for accessing the DDM instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the enterprise project id.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the engine version.`,
			},
		},
	}
	return &sc
}

func resourceDdmInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmInstances: Query the List of DDM instances
	var (
		getDdmInstancesHttpUrl = "v1/{project_id}/instances"
		getDdmInstancesProduct = "ddm"
	)
	getDdmInstancesClient, err := cfg.NewServiceClient(getDdmInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getDdmInstancesPath := getDdmInstancesClient.Endpoint + getDdmInstancesHttpUrl
	getDdmInstancesPath = strings.ReplaceAll(getDdmInstancesPath, "{project_id}", getDdmInstancesClient.ProjectID)

	getDdmInstancesResp, err := pagination.ListAllItems(
		getDdmInstancesClient,
		"offset",
		getDdmInstancesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmInstances")
	}

	getDdmInstancesRespJson, err := json.Marshal(getDdmInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDdmInstancesRespBody interface{}
	err = json.Unmarshal(getDdmInstancesRespJson, &getDdmInstancesRespBody)
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
		d.Set("instances", flattenGetInstancesResponseBodyInstance(d, getDdmInstancesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetInstancesResponseBodyInstance(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	name := d.Get("name").(string)
	status := d.Get("status").(string)
	enterpriseProjectId := d.Get("enterprise_project_id").(string)
	engineVersion := d.Get("engine_version").(string)
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		azCodes := utils.PathSearch("available_zone", v, "")
		availabilityZones := strings.Split(azCodes.(string), ",")
		instanceName := utils.PathSearch("name", v, nil)
		instanceStatus := utils.PathSearch("status", v, nil)
		instanceEnterpriseProjectId := utils.PathSearch("enterprise_project_id", v, nil)
		instanceEngineVersion := utils.PathSearch("engine_version", v, nil)
		if name != "" && name != instanceName {
			continue
		}
		if status != "" && status != instanceStatus {
			continue
		}
		if enterpriseProjectId != "" && enterpriseProjectId != instanceEnterpriseProjectId {
			continue
		}
		if engineVersion != "" && engineVersion != instanceEngineVersion {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"status":                instanceStatus,
			"name":                  instanceName,
			"availability_zones":    availabilityZones,
			"vpc_id":                utils.PathSearch("vpc_id", v, nil),
			"subnet_id":             utils.PathSearch("subnet_id", v, nil),
			"security_group_id":     utils.PathSearch("security_group_id", v, nil),
			"node_num":              utils.PathSearch("node_count", v, nil),
			"access_ip":             utils.PathSearch("access_ip", v, nil),
			"access_port":           utils.PathSearch("access_port", v, nil),
			"enterprise_project_id": instanceEnterpriseProjectId,
			"engine_version":        instanceEngineVersion,
		})
	}
	return rst
}
