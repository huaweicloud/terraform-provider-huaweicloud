// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

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

// @API ModelArts GET /v1/{project_id}/notebooks/flavors
func DataSourceNotebookFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceNotebookFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Processor type. The valid values are: **CPU**, **GPU**, **ASCEND**.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cluster type.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        notebookFlavorsFlavorsSchema(),
				Computed:    true,
				Description: `The list of flavors.`,
			},
		},
	}
}

func notebookFlavorsFlavorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the flavor.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the flavor.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Architecture type. The valid values are **X86_64** and **AARCH64**.`,
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Processor type. The valid values are: **CPU**, **GPU**, **ASCEND**.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specification description.`,
			},
			"feature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Flavor type.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Memory size.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of vCPUs.`,
			},
			"free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Free flavor or not.`,
			},
			"sold_out": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether resources are sufficient.`,
			},
			"billing": {
				Type:     schema.TypeList,
				Elem:     notebookFlavorsFlavorsBillingInfoSchema(),
				Computed: true,
			},
			"gpu": {
				Type:     schema.TypeList,
				Elem:     notebookFlavorsFlavorsGpuInfoSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func notebookFlavorsFlavorsBillingInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Billing code.`,
			},
			"unit_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Billing unit.`,
			},
		},
	}
	return &sc
}

func notebookFlavorsFlavorsGpuInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"gpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of GPUs.`,
			},
			"gpu_memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `GPU memory.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `GPU type.`,
			},
		},
	}
	return &sc
}

func resourceNotebookFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listFlavors: Query the list of ModelArts notebook flavors
	var (
		listFlavorsHttpUrl = "v1/{project_id}/notebooks/flavors"
		listFlavorsProduct = "modelarts"
	)
	listFlavorsClient, err := cfg.NewServiceClient(listFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	listFlavorsPath := listFlavorsClient.Endpoint + listFlavorsHttpUrl
	listFlavorsPath = strings.ReplaceAll(listFlavorsPath, "{project_id}", listFlavorsClient.ProjectID)

	listFlavorsqueryParams := buildListFlavorsQueryParams(d)
	listFlavorsPath += listFlavorsqueryParams

	listFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	listFlavorsResp, err := listFlavorsClient.Request("GET", listFlavorsPath, &listFlavorsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving NotebookFlavors")
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
		d.Set("flavors", flattenListNotebookFlavorsFlavors(listFlavorsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListNotebookFlavorsFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("data[?!sold_out]", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"arch":        utils.PathSearch("arch", v, nil),
			"category":    utils.PathSearch("category", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"feature":     utils.PathSearch("feature", v, nil),
			"memory":      utils.PathSearch("memory", v, nil),
			"vcpus":       utils.PathSearch("vcpus", v, nil),
			"free":        utils.PathSearch("free", v, nil),
			"sold_out":    utils.PathSearch("sold_out", v, nil),
			"billing":     flattenFlavorsBilling(v),
			"gpu":         flattenFlavorsGpu(v),
		})
	}
	return rst
}

func flattenFlavorsBilling(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("billing", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"code":     utils.PathSearch("code", curJson, nil),
			"unit_num": utils.PathSearch("unit_num", curJson, nil),
		},
	}
	return rst
}

func flattenFlavorsGpu(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("gpu", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"gpu":        utils.PathSearch("gpu", curJson, nil),
			"gpu_memory": utils.PathSearch("gpu_memory", curJson, nil),
			"type":       utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func buildListFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("category"); ok {
		res = fmt.Sprintf("%s&category=%v", res, v)
	}

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
