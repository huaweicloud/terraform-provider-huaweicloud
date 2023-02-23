// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				Description: `Specifies the ID of a engine.`,
			},
			"group_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the compute resource architecture type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"X86", "ARM",
				}, false),
			},
			"type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type code.`,
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VM flavor types recorded in DDM.`,
			},
			"iaas_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VM flavor types recorded by the IaaS layer.`,
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
				Description: `Indicates the compute resource architecture type.`,
			},
			"type_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource type code.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VM flavor types recorded in DDM.`,
			},
			"iaas_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the VM flavor types recorded by the IaaS layer.`,
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
		getDdmFlavorsProduct = "ddmv2"
	)
	getDdmFlavorsClient, err := cfg.NewServiceClient(getDdmFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DdmFlavors Client: %s", err)
	}

	getDdmFlavorsPath := getDdmFlavorsClient.Endpoint + getDdmFlavorsHttpUrl
	getDdmFlavorsPath = strings.ReplaceAll(getDdmFlavorsPath, "{project_id}", getDdmFlavorsClient.ProjectID)

	offset, limit, total := 0, 0, 0

	getDdmFlavorsQueryParams := buildGetDdmFlavorsQueryParams(d, offset)
	getDdmFlavorsPath += getDdmFlavorsQueryParams
	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	groupType := d.Get("group_type").(string)
	typeCode := d.Get("type_code").(string)
	code := d.Get("code").(string)
	iaasCode := d.Get("iaas_code").(string)

	allFlavors := make([]interface{}, 0)
	var flavors []interface{}
	for {
		getDdmFlavorsResp, err := getDdmFlavorsClient.Request("GET", getDdmFlavorsPath, &getInstanceOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving DdmFlavors")
		}
		getDdmFlavorsRespBody, err := utils.FlattenResponse(getDdmFlavorsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		flavors, offset, limit, total = flattenGetFlavorsResponseBodyFlavorGroup(getDdmFlavorsRespBody, groupType,
			typeCode, code, iaasCode)
		allFlavors = append(allFlavors, flavors...)
		if offset+limit >= total {
			break
		}
		getDdmFlavorsPath = updatePathOffset(getDdmFlavorsPath, offset+limit)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavors", allFlavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetFlavorsResponseBodyFlavorGroup(resp interface{}, groupType, typeCode, code, iaasCode string) ([]interface{}, int, int, int) {
	if resp == nil {
		return nil, 0, 0, 0
	}
	curJson := utils.PathSearch("computeFlavorGroups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		flavorGroupType := utils.PathSearch("groupType", v, nil)
		if groupType != "" && groupType != flavorGroupType {
			continue
		}
		offset := utils.PathSearch("offset", v, nil).(float64)
		limit := utils.PathSearch("limit", v, nil).(float64)
		total := utils.PathSearch("total", v, nil).(float64)
		return flattenFlavorGroupFlavors(v, typeCode, code, iaasCode), int(offset), int(limit), int(total)
	}
	return nil, 0, 0, 0
}

func flattenFlavorGroupFlavors(resp interface{}, typeCode, code, iaasCode string) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("computeFlavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		flavorTypeCode := utils.PathSearch("typeCode", v, nil)
		flavorCode := utils.PathSearch("code", v, nil)
		flavorIaasCode := utils.PathSearch("iaasCode", v, nil)
		if typeCode != "" && typeCode != flavorTypeCode {
			continue
		}
		if code != "" && code != flavorCode {
			continue
		}
		if iaasCode != "" && iaasCode != flavorIaasCode {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":        utils.PathSearch("id", v, nil),
			"type_code": flavorTypeCode,
			"code":      flavorCode,
			"iaas_code": flavorIaasCode,
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
	return fmt.Sprintf("%s&offset=%v", path[:index], offset)
}
