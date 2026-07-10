package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/protect-measure-baseline
func DataSourceProtectMeasureBaseline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProtectMeasureBaselineRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protect_measure_baseline": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     protectMeasureBaselineSchema(),
			},
		},
	}
}

func protectMeasureBaselineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     measureRequirementDetailSchema(),
			},
		},
	}
}

func measureRequirementDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_type_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     protectMeasureDataTypeDetailSchema(),
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"measure_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     protectMeasureDetailSchema(),
			},
			"protect_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func protectMeasureDataTypeDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"life_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func protectMeasureDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"life_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"measure_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceProtectMeasureBaselineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/protect-measure-baseline"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC protect measure baseline: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("protect_measure_baseline", flattenProtectMeasureBaseline(
			utils.PathSearch("protect_measure_baseline", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProtectMeasureBaseline(resp interface{}) []interface{} {
	baselineMap, ok := resp.(map[string]interface{})
	if !ok || len(baselineMap) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(baselineMap))
	for key, value := range baselineMap {
		details, ok := value.([]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    key,
			"detail": flattenMeasureRequirementDetailList(details),
		})
	}

	return rst
}

func flattenMeasureRequirementDetailList(details []interface{}) []interface{} {
	if len(details) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(details))
	for _, v := range details {
		rst = append(rst, map[string]interface{}{
			"create_time": utils.PathSearch("create_time", v, nil),
			"data_type_info": flattenProtectMeasureDataTypeDetail(
				utils.PathSearch("data_type_info", v, nil)),
			"id": utils.PathSearch("id", v, nil),
			"measure_info": flattenProtectMeasureDetail(
				utils.PathSearch("measure_info", v, nil)),
			"protect_level": utils.PathSearch("protect_level", v, nil),
			"update_time":   utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}

func flattenProtectMeasureDataTypeDetail(obj interface{}) []interface{} {
	if obj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category":     utils.PathSearch("category", obj, nil),
			"create_time":  utils.PathSearch("create_time", obj, nil),
			"data_type":    utils.PathSearch("data_type", obj, nil),
			"data_type_id": utils.PathSearch("data_type_id", obj, nil),
			"id":           utils.PathSearch("id", obj, nil),
			"is_deleted":   utils.PathSearch("is_deleted", obj, nil),
			"life_cycle":   utils.PathSearch("life_cycle", obj, nil),
			"update_time":  utils.PathSearch("update_time", obj, nil),
		},
	}
}

func flattenProtectMeasureDetail(obj interface{}) []interface{} {
	if obj == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"create_time":  utils.PathSearch("create_time", obj, nil),
			"description":  utils.PathSearch("description", obj, nil),
			"id":           utils.PathSearch("id", obj, nil),
			"is_deleted":   utils.PathSearch("is_deleted", obj, nil),
			"life_cycle":   utils.PathSearch("life_cycle", obj, nil),
			"measure_type": utils.PathSearch("measure_type", obj, nil),
			"name":         utils.PathSearch("name", obj, nil),
			"update_time":  utils.PathSearch("update_time", obj, nil),
		},
	}
}
