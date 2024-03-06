// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DBSS
// ---------------------------------------------------------------

package dbss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS GET /v1/{project_id}/dbss/audit/specification
func DataSourceDbssFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDbssFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the flavor.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the availability zone which the flavor belongs to.`,
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the level of the flavor.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `Specifies the memory size(GB) in the flavor.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the number of CPUs.`,
			},
			"proxy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the maximum supported database instances.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        dbssFlavorsFlavorSchema(),
				Computed:    true,
				Description: `Indicates the list of DBSS flavors.`,
			},
		},
	}
}

func dbssFlavorsFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the flavor.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the level of the flavor.`,
			},
			"proxy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum supported database instances.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of CPUs.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the memory size(GB) in the flavor.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the availability zones which the flavor belongs to`,
			},
		},
	}
	return &sc
}

func resourceDbssFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDbssFlavors: Query the List of DBSS flavors
	var (
		getDbssFlavorsHttpUrl = "v1/{project_id}/dbss/audit/specification"
		getDbssFlavorsProduct = "dbss"
	)
	getDbssFlavorsClient, err := cfg.NewServiceClient(getDbssFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	getDbssFlavorsPath := getDbssFlavorsClient.Endpoint + getDbssFlavorsHttpUrl
	getDbssFlavorsPath = strings.ReplaceAll(getDbssFlavorsPath, "{project_id}", getDbssFlavorsClient.ProjectID)

	getDbssFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDbssFlavorsResp, err := getDbssFlavorsClient.Request("GET", getDbssFlavorsPath, &getDbssFlavorsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DBSS flavors")
	}

	getDbssFlavorsRespBody, err := utils.FlattenResponse(getDbssFlavorsResp)
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
		d.Set("flavors", filterGetFlavorsResponseBodyFlavor(
			flattenGetFlavorsResponseBodyFlavor(getDbssFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetFlavorsResponseBodyFlavor(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("specification", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		memory := utils.PathSearch("ram", v, float64(0)).(float64) / 1024
		rst = append(rst, map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"level":              utils.PathSearch("level", v, nil),
			"proxy":              utils.PathSearch("proxy", v, nil),
			"vcpus":              utils.PathSearch("vcpus", v, nil),
			"memory":             memory,
			"availability_zones": utils.PathSearch("azs", v, nil),
		})
	}
	return rst
}

func filterGetFlavorsResponseBodyFlavor(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("flavor_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("availability_zone"); ok {
			azs := utils.PathSearch("availability_zones", v, make([]interface{}, 0))
			if !utils.IsStrContainsSliceElement(param.(string), utils.ExpandToStringList(azs.([]interface{})), false, true) {
				continue
			}
		}
		if param, ok := d.GetOk("level"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("level", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("memory"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("memory", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("vcpus"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("vcpus", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("proxy"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("proxy", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
