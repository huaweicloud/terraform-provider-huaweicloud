// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"log"
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

// @API ModelArts GET /v1/{project_id}/resourceflavors
func DataSourceResourceFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceResourceFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tag key.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of resource flavor.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        resourceFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of resource flavors.`,
			},
		},
	}
}

func resourceFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Flavor ID.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The key/value pairs to associate with the flavor.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of resource flavor.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Computer architecture.`,
			},
			"cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Number of CPU cores.`,
			},
			"memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Memory size in GiB.`,
			},
			"gpu": {
				Type:     schema.TypeList,
				Elem:     resourceFlavorsFlavorsGpuSchema(),
				Computed: true,
			},
			"npu": {
				Type:     schema.TypeList,
				Elem:     resourceFlavorsFlavorsNpuSchema(),
				Computed: true,
			},
			"volume": {
				Type:        schema.TypeList,
				Elem:        resourceFlavorsFlavorsVolumeSchema(),
				Computed:    true,
				Description: `Data disks information.`,
			},
			"billing_modes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: `Billing mode supported by the flavor.`,
			},
			"job_flavors": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Training job types supported by the resource flavor.`,
			},
			"az_status": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Sales status of a resource specification in each AZ. The value is (AZ, Status).`,
			},
		},
	}
	return &sc
}

func resourceFlavorsFlavorsGpuSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `GPU type.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Number of GPUs.`,
			},
		},
	}
	return &sc
}

func resourceFlavorsFlavorsNpuSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `NPU type.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Number of NPUs.`,
			},
		},
	}
	return &sc
}

func resourceFlavorsFlavorsVolumeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Disk type.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Disk size, in GiB.`,
			},
		},
	}
	return &sc
}

func resourceResourceFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listResourceFlavors: Query the list of ModelArts resource flavors
	var (
		listResourceFlavorsHttpUrl = "v1/{project_id}/resourceflavors"
		listResourceFlavorsProduct = "modelarts"
	)
	listResourceFlavorsClient, err := cfg.NewServiceClient(listResourceFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	listResourceFlavorsPath := listResourceFlavorsClient.Endpoint + listResourceFlavorsHttpUrl
	listResourceFlavorsPath = strings.ReplaceAll(listResourceFlavorsPath, "{project_id}", listResourceFlavorsClient.ProjectID)

	listResourceFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	allItems := make([]interface{}, 0)
	nextMarker := ""
	for {
		queryPath := listResourceFlavorsPath + buildListResourceFlavorsQueryParams(d, nextMarker)
		listResourceFlavorsResp, err := listResourceFlavorsClient.Request("GET", queryPath, &listResourceFlavorsOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving resource flavors")
		}

		listResourceFlavorsRespBody, err := utils.FlattenResponse(listResourceFlavorsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		items := utils.PathSearch("items", listResourceFlavorsRespBody, make([]interface{}, 0)).([]interface{})

		if len(items) > 0 {
			allItems = append(allItems, items...)
		}

		nextMarker = utils.PathSearch("metadata.continue", listResourceFlavorsRespBody, "").(string)
		if len(nextMarker) == 0 {
			break
		}
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavors", filterListResourceFlavorsFlavors(flattenListResourceFlavorsFlavors(allItems), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListResourceFlavorsFlavors(curArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("metadata.name", v, nil),
			"tags":          utils.PathSearch("metadata.labels", v, nil),
			"type":          utils.PathSearch("spec.type", v, nil),
			"arch":          utils.PathSearch("spec.cpuArch", v, nil),
			"cpu":           utils.PathSearch("spec.cpu", v, nil),
			"memory":        utils.PathSearch("spec.memory", v, nil),
			"gpu":           flattenResourceFlavorsGpu(v),
			"npu":           flattenResourceFlavorsNpu(v),
			"volume":        flattenResourceFlavorsVolume(v),
			"billing_modes": utils.PathSearch("spec.billingModes", v, nil),
			"job_flavors":   utils.PathSearch("spec.jobFlavors", v, nil),
			"az_status":     utils.PathSearch("status.phase", v, nil),
		})
	}
	return rst
}

func flattenResourceFlavorsGpu(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("spec.gpu", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing spec.gpu from response")
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("type", curJson, nil),
			"size": utils.PathSearch("size", curJson, nil),
		},
	}
	return rst
}

func flattenResourceFlavorsNpu(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("spec.npu", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing spec.npu from response")
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("type", curJson, nil),
			"size": utils.PathSearch("size", curJson, nil),
		},
	}
	return rst
}

func flattenResourceFlavorsVolume(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.dataVolumes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"type": utils.PathSearch("volumeType", v, nil),
			"size": utils.PathSearch("size", v, nil),
		})
	}
	return rst
}

func filterListResourceFlavorsFlavors(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListResourceFlavorsQueryParams(d *schema.ResourceData, next string) string {
	res := ""
	if v, ok := d.GetOk("tag"); ok {
		res = fmt.Sprintf("%s&labelSelector=%v", res, v)
	}

	if len(next) > 0 {
		res = fmt.Sprintf("%s&continue=%v", res, next)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
