// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type queryParma struct {
	resp    interface{}
	cpuArch string
	code    string
	vcpus   string
	memory  string
}

type queryRes struct {
	x86Flavors []interface{}
	armFlavors []interface{}
	offset     int
	limit      int
	x86Total   int
	armTotal   int
}

// @API DDM GET /v2/{project_id}/flavors
func DataSourceDdmFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmFlavorsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of an engine.`,
			},
			"cpu_arch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the compute resource architecture type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"X86", "ARM",
				}, false),
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VM flavor types recorded in DDM.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the number of CPUs.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the memory size. Unit GB.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        FlavorsFlavorSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM compute flavors.`,
			},
		},
	}
}

func FlavorsFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a flavor.`,
			},
			"cpu_arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the compute resource architecture type.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VM flavor types recorded in DDM.`,
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of CPUs.`,
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the memory size.`,
			},
		},
	}
	return &sc
}

func resourceDdmFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmFlavors: Query the List of DDM flavors
	var (
		getDdmFlavorsHttpUrl = "v2/{project_id}/flavors"
		getDdmFlavorsProduct = "ddm"
	)
	getDdmFlavorsClient, err := cfg.NewServiceClient(getDdmFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getDdmFlavorsPath := getDdmFlavorsClient.Endpoint + getDdmFlavorsHttpUrl
	getDdmFlavorsPath = strings.ReplaceAll(getDdmFlavorsPath, "{project_id}", getDdmFlavorsClient.ProjectID)

	getDdmFlavorsQueryParams := buildGetDdmFlavorsQueryParams(d, 0)
	getDdmFlavorsPath += getDdmFlavorsQueryParams
	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	cpuArch := d.Get("cpu_arch").(string)
	code := d.Get("code").(string)
	var vcpus string
	if v, ok := d.GetOk("vcpus"); ok {
		vcpus = strconv.Itoa(v.(int))
	}
	var memory string
	if v, ok := d.GetOk("memory"); ok {
		memory = strconv.Itoa(v.(int))
	}

	x86Flavors := make([]interface{}, 0)
	armFlavors := make([]interface{}, 0)
	for {
		getDdmFlavorsResp, err := getDdmFlavorsClient.Request("GET", getDdmFlavorsPath, &getInstanceOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving DdmFlavors")
		}
		getDdmFlavorsRespBody, err := utils.FlattenResponse(getDdmFlavorsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res := flattenGetFlavorsResponseBodyFlavorGroup(&queryParma{
			resp:    getDdmFlavorsRespBody,
			cpuArch: cpuArch,
			code:    code,
			vcpus:   vcpus,
			memory:  memory,
		})
		x86Flavors = append(x86Flavors, res.x86Flavors...)
		armFlavors = append(armFlavors, res.armFlavors...)
		if res.offset+res.limit >= res.x86Total && res.offset+res.limit >= res.armTotal {
			break
		}
		getDdmFlavorsPath = updatePathOffset(getDdmFlavorsPath, res.offset+res.limit)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	var flavors []interface{}
	flavors = append(flavors, x86Flavors...)
	flavors = append(flavors, armFlavors...)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetFlavorsResponseBodyFlavorGroup(parma *queryParma) *queryRes {
	if parma.resp == nil {
		return &queryRes{}
	}
	curJson := utils.PathSearch("computeFlavorGroups", parma.resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	var x86Flavors []interface{}
	var armFlavors []interface{}
	var offset, limit, x86Total, armTotal float64

	for _, v := range curArray {
		flavorCPUArch := utils.PathSearch("groupType", v, "")
		if parma.cpuArch != "" && parma.cpuArch != flavorCPUArch {
			continue
		}
		offset = utils.PathSearch("offset", v, float64(0)).(float64)
		limit = utils.PathSearch("limit", v, float64(0)).(float64)
		if flavorCPUArch == "X86" {
			x86Flavors = flattenFlavorGroupFlavors(v, flavorCPUArch.(string), parma)
			x86Total = utils.PathSearch("total", v, float64(0)).(float64)
		} else {
			armFlavors = flattenFlavorGroupFlavors(v, flavorCPUArch.(string), parma)
			armTotal = utils.PathSearch("total", v, float64(0)).(float64)
		}
	}
	return &queryRes{
		x86Flavors: x86Flavors,
		armFlavors: armFlavors,
		offset:     int(offset),
		limit:      int(limit),
		x86Total:   int(x86Total),
		armTotal:   int(armTotal),
	}
}

func flattenFlavorGroupFlavors(resp interface{}, flavorGroupType string, parma *queryParma) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("computeFlavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		code := utils.PathSearch("code", v, "")
		vcpus, _ := utils.PathSearch("cpu", v, "").(string)
		memory, _ := utils.PathSearch("mem", v, "").(string)
		if parma.code != "" && parma.code != code {
			continue
		}
		if parma.vcpus != "" && parma.vcpus != vcpus {
			continue
		}
		if parma.memory != "" && parma.memory != memory {
			continue
		}
		vcpusNum, _ := strconv.Atoi(vcpus)
		memoryNum, _ := strconv.Atoi(memory)
		rst = append(rst, map[string]interface{}{
			"id":       utils.PathSearch("id", v, nil),
			"cpu_arch": flavorGroupType,
			"code":     code,
			"vcpus":    vcpusNum,
			"memory":   memoryNum,
		})
	}
	return rst
}

func buildGetDdmFlavorsQueryParams(d *schema.ResourceData, offset int) string {
	res := ""
	if v, ok := d.GetOk("engine_id"); ok {
		res = fmt.Sprintf("%s&engine_id=%v", res, v)
	}
	res = fmt.Sprintf("%s&offset=%v", res, offset)
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func updatePathOffset(path string, offset int) string {
	index := strings.Index(path, "offset")
	return fmt.Sprintf("%soffset=%v", path[:index], offset)
}
