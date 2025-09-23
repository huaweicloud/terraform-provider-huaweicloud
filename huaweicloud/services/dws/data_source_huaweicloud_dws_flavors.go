// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/node-types
func DataSourceDwsFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDwsFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The availability zone name.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The vcpus of the dws node flavor.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The ram of the dws node flavor in GB.`,
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of datastore.`,
				ValidateFunc: validation.StringInSlice([]string{
					"dws", "hybrid", "stream",
				}, false),
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        dwsFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of flavor detail.`,
			},
		},
	}
}

func dwsFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the dws node flavor.`,
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of datastore.`,
			},
			"datastore_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of datastore.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The vcpus of the dws node flavor.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ram of the dws node flavor in GB.`,
			},
			"volumetype": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Disk type.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The default disk size in GB.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The list of availability zones.`,
			},
			"elastic_volume_specs": {
				Type:        schema.TypeList,
				Elem:        dwsFlavorsFlavorsElasticVolumeSpecSchema(),
				Computed:    true,
				Description: `The typical specification, If the volume specification is elastic.`,
			},
		},
	}
	return &sc
}

func dwsFlavorsFlavorsElasticVolumeSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"step": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Disk size increment step.`,
			},
			"min_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Minimum disk size.`,
			},
			"max_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Maximum disk size.`,
			},
		},
	}
	return &sc
}

func resourceDwsFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of DWS cluster flavors
	var (
		listFlavorsHttpUrl = "v2/{project_id}/node-types"
		listFlavorsProduct = "dws"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DwsFlavors Client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}
	listFlavorsResp, err := listFlavorsClient.Request("GET", listFlavorsPath, &listFlavorsOpt)
	if err != nil {
		return diag.Errorf("error retrieving DWS flavors: %s", err)
	}

	listFlavorsRespBody, err := utils.FlattenResponse(listFlavorsResp)
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
		d.Set("flavors", filterListNodeTypesFlavors(
			flattenListNodeTypesFlavors(listFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListNodeTypesFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("node_types", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id":            utils.PathSearch("spec_name", v, nil),
			"datastore_type":       utils.PathSearch("datastore_type", v, nil),
			"datastore_version":    utils.PathSearch("datastores[0].version", v, nil),
			"vcpus":                utils.PathSearch("vcpus", v, nil),
			"memory":               utils.PathSearch("ram", v, nil),
			"volumetype":           utils.PathSearch("detail[?type=='LOCAL_DISK' || type=='SSD' ].type|[0]", v, nil),
			"size":                 utils.PathSearch("detail[?type=='LOCAL_DISK' || type=='SSD' ].value|[0]|to_number(@)", v, nil),
			"availability_zones":   utils.PathSearch("availability_zones[?status=='normal'].code", v, nil),
			"elastic_volume_specs": flattenFlavorsElasticVolumeSpecs(v),
		})
	}
	return rst
}

func flattenFlavorsElasticVolumeSpecs(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("elastic_volume_specs[0]", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"step":     utils.PathSearch("step", curJson, nil),
			"min_size": utils.PathSearch("min_size", curJson, nil),
			"max_size": utils.PathSearch("max_size", curJson, nil),
		},
	}
	return rst
}

func filterListNodeTypesFlavors(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("availability_zone"); ok {
			availabilityZones := utils.ExpandToStringList(utils.PathSearch("availability_zones", v, []string{}).([]interface{}))
			if !utils.StrSliceContains(availabilityZones, param.(string)) {
				continue
			}
		}
		if param, ok := d.GetOk("vcpus"); ok && fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("vcpus", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("memory"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("memory", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("datastore_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("datastore_type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
