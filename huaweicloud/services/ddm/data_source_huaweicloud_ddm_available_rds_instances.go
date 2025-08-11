package ddm

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

// @API DDM GET /v1/{project_id}/instances/{instance_id}/rds
func DataSourceDdmAvailableRdsInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmAvailableRdsInstancesRead,
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
			"instances": {
				Type:     schema.TypeList,
				Elem:     parameterAvailableRdsInstancesSchema(),
				Computed: true,
			},
		},
	}
}

func parameterAvailableRdsInstancesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_software_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDdmAvailableRdsInstancesRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmAvailableRdsInstances: Query the List of DDM instance available rds instances
	var (
		getDdmAvailableRdsInstancesHttpUrl = "v1/{project_id}/instances/{instance_id}/rds"
		getDdmAvailableRdsInstancesProduct = "ddm"
	)
	getDdmAvailableRdsInstancesClient, err := cfg.NewServiceClient(getDdmAvailableRdsInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	getDdmAvailableRdsInstancesPath := getDdmAvailableRdsInstancesClient.Endpoint + getDdmAvailableRdsInstancesHttpUrl
	getDdmAvailableRdsInstancesPath = strings.ReplaceAll(getDdmAvailableRdsInstancesPath, "{project_id}",
		getDdmAvailableRdsInstancesClient.ProjectID)
	getDdmAvailableRdsInstancesPath = strings.ReplaceAll(getDdmAvailableRdsInstancesPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	getDdmAvailableRdsInstancesResp, err := pagination.ListAllItems(
		getDdmAvailableRdsInstancesClient,
		"offset",
		getDdmAvailableRdsInstancesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		diag.Errorf("error retrieving DDM instance available rds instances: %s", err)
	}

	getDdmAvailableRdsInstancesRespJson, err := json.Marshal(getDdmAvailableRdsInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var getDdmAvailableRdsInstancesRespBody any
	err = json.Unmarshal(getDdmAvailableRdsInstancesRespJson, &getDdmAvailableRdsInstancesRespBody)
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
		d.Set("instances", flattenGetInstancesResponseBody(getDdmAvailableRdsInstancesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetInstancesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"project_id":              utils.PathSearch("project_id", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"engine_name":             utils.PathSearch("engine_name", v, nil),
			"engine_software_version": utils.PathSearch("engine_software_version", v, nil),
			"private_ip":              utils.PathSearch("private_ip", v, nil),
			"mode":                    utils.PathSearch("mode", v, nil),
			"port":                    utils.PathSearch("port", v, nil),
			"az_code":                 utils.PathSearch("az_code", v, nil),
			"time_zone":               utils.PathSearch("time_zone", v, nil),
		})
	}
	return rst
}
