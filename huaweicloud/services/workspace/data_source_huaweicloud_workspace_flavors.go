package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

// @API Workspace GET /v2/{project_id}/products
func DataSourceWorkspaceFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkspaceFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_gpu": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"system_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charging_mode": {
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
		},
	}
}

func dataSourceWorkspaceFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of Workspace flavors
	var (
		listFlavorsHttpUrl = "v2/{project_id}/products"
		listFlavorsProduct = "workspace"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsqueryParams := buildListFlavorsQueryParams(d)
	listFlavorsPath += listFlavorsqueryParams

	listFlavorsResp, err := pagination.ListAllItems(
		listFlavorsClient,
		"offset",
		listFlavorsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Flavor")
	}

	listFlavorsRespJson, err := json.Marshal(listFlavorsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listFlavorsRespBody interface{}
	err = json.Unmarshal(listFlavorsRespJson, &listFlavorsRespBody)
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
		d.Set("flavors", filterListFlavorsBodyFlavor(flattenListFlavorsBodyFlavor(listFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListFlavorsBodyFlavor(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("products", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("product_id", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"architecture":     utils.PathSearch("architecture", v, nil),
			"vcpus":            convertStrToInt(utils.PathSearch("cpu", v, nil)),
			"memory":           utils.ConvertMemoryUnit(utils.PathSearch("memory", v, nil), 1),
			"is_gpu":           utils.PathSearch("is_gpu", v, nil),
			"system_disk_type": utils.PathSearch("system_disk_type", v, nil),
			"system_disk_size": utils.PathSearch("system_disk_size", v, nil),
			"description":      utils.PathSearch("descriptions", v, nil),
			"charging_mode":    normalizeChargingMode(utils.PathSearch("charge_mode", v, nil)),
			"status":           utils.PathSearch("status", v, nil),
		})
	}
	return rst
}

func normalizeChargingMode(chargeModeCode interface{}) string {
	if chargeModeCode == "0" {
		return "prePaid"
	}
	return "postPaid"
}

func convertStrToInt(str interface{}) int {
	num, err := strconv.Atoi(str.(string))
	if err != nil {
		fmt.Printf("convert string to int fail")
	}
	return num
}

func filterListFlavorsBodyFlavor(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("vcpus"); ok && param != utils.PathSearch("vcpus", v, nil) {
			continue
		}

		if param, ok := d.GetOk("memory"); ok && param != utils.PathSearch("memory", v, nil) {
			continue
		}
		rst = append(rst, v)
	}
	return rst
}

func buildListFlavorsQueryParams(d *schema.ResourceData) string {
	res := "charge_mode=1&status=normal"
	if v, ok := d.GetOk("availability_zone"); ok {
		res += fmt.Sprintf("&availability_zone=%v", v)
	}

	if v, ok := d.GetOk("os_type"); ok {
		res += fmt.Sprintf("&os_type=%v", v)
	}

	if res != "" {
		res = "?" + res[0:]
	}
	return res
}
