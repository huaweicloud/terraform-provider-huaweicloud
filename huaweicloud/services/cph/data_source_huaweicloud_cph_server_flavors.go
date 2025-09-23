// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CPH
// ---------------------------------------------------------------

package cph

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

// @API CPH GET /v1/{project_id}/cloud-phone/server-models
func DataSourceServerFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceServerFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  `The type of the CPH server flavor.`,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The vcpus of the CPH server.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The ram of the CPH server in GB.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        serverFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of flavor detail.`,
			},
		},
	}
}

func serverFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the flavor.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The vcpus of the CPH server.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ram of the CPH server in GB.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the CPH server flavor.`,
			},
			"extend_spec": {
				Type:     schema.TypeList,
				Elem:     serverFlavorsFlavorsExtendSpecSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func serverFlavorsFlavorsExtendSpecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vcpus": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the vcpus.`,
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the ram.`,
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the disk.`,
			},
			"network_interface": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the network interface.`,
			},
			"gpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the gpu.`,
			},
			"bms_flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended description of the bms flavor.`,
			},
			"gpu_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The gpu count.`,
			},
			"numa_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The numa count.`,
			},
		},
	}
	return &sc
}

func resourceServerFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of CPH server flavors
	var (
		listFlavorsHttpUrl = "v1/{project_id}/cloud-phone/server-models"
		listFlavorsProduct = "cph"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsqueryParams := buildListFlavorsQueryParams(d)
	listFlavorsPath += listFlavorsqueryParams

	listFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listFlavorsResp, err := listFlavorsClient.Request("GET", listFlavorsPath, &listFlavorsOpt)
	if err != nil {
		return diag.Errorf("error retrieving CPH server flavors: %s", err)
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
		d.Set("flavors", filterListServerModelsFlavors(
			flattenListServerModelsFlavors(listFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListServerModelsFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("server_models", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id":   utils.PathSearch("server_model_name", v, nil),
			"vcpus":       utils.PathSearch("cpu", v, nil),
			"memory":      utils.PathSearch("memory", v, nil),
			"type":        fmt.Sprint(utils.PathSearch("product_type", v, nil)),
			"extend_spec": flattenFlavorsExtendSpec(v),
		})
	}
	return rst
}

func flattenFlavorsExtendSpec(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("extend_spec", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"vcpus":             utils.PathSearch("cpu", curJson, nil),
			"memory":            utils.PathSearch("memory", curJson, nil),
			"disk":              utils.PathSearch("disk", curJson, nil),
			"network_interface": utils.PathSearch("network_interface", curJson, nil),
			"gpu":               utils.PathSearch("gpu", curJson, nil),
			"bms_flavor":        utils.PathSearch("bms_flavor", curJson, nil),
			"gpu_count":         utils.PathSearch("gpu_count", curJson, nil),
			"numa_count":        utils.PathSearch("numa_count", curJson, nil),
		},
	}
	return rst
}

func filterListServerModelsFlavors(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("vcpus"); ok && fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("vcpus", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("memory"); ok && fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("memory", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&product_type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
