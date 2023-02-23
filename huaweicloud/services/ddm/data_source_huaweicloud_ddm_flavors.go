// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

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
				Optional:    true,
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
			"flavor_groups": {
				Type:        schema.TypeList,
				Elem:        FlavorsFlavorGroupSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM flavor.`,
			},
		},
	}
}

func FlavorsFlavorGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the compute resource architecture type.`,
			},
			"flavors": {
				Type:        schema.TypeList,
				Elem:        FlavorsFlavorGroupFlavorSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM compute flavors.`,
			},
		},
	}
	return &sc
}

func FlavorsFlavorGroupFlavorSchema() *schema.Resource {
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

	getDdmFlavorsQueryParams := buildGetDdmFlavorsQueryParams(d)
	getDdmFlavorsPath += getDdmFlavorsQueryParams

	getDdmFlavorsResp, err := pagination.ListAllItems(
		getDdmFlavorsClient,
		"offset",
		getDdmFlavorsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmFlavors")
	}

	getDdmFlavorsRespJson, err := json.Marshal(getDdmFlavorsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println(string(getDdmFlavorsRespJson))
	var getDdmFlavorsRespBody interface{}
	err = json.Unmarshal(getDdmFlavorsRespJson, &getDdmFlavorsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	groupType := d.Get("group_type").(string)
	typeCode := d.Get("type_code").(string)
	code := d.Get("code").(string)
	iaasCode := d.Get("iaas_code").(string)

	groups := flattenGetFlavorsResponseBodyFlavorGroup(getDdmFlavorsRespBody, groupType, typeCode, code, iaasCode)
	bytes, err := json.Marshal(groups)
	fmt.Println("groups info: ", string(bytes))

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flavor_groups", flattenGetFlavorsResponseBodyFlavorGroup(getDdmFlavorsRespBody, groupType,
			typeCode, code, iaasCode)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetFlavorsResponseBodyFlavorGroup(resp interface{}, groupType, typeCode, code, iaasCode string) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("computeFlavorGroups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		flavorGroupType := utils.PathSearch("groupType", v, nil)
		if groupType != "" && groupType != flavorGroupType {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"groupType": flavorGroupType,
			"flavors":   flattenFlavorGroupFlavors(v, typeCode, code, iaasCode),
		})
	}
	return rst
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

func buildGetDdmFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("engine_id"); ok {
		res = fmt.Sprintf("%s&engine_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
