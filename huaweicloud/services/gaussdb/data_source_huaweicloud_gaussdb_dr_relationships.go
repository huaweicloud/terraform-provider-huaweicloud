package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3.5/{project_id}/disaster-recovery/relations
func DataSourceGaussDbDrRelationships() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbDrRelationshipsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dr_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dr_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbDrRelationshipsRelationsSchema(),
			},
		},
	}
}

func gaussDbDrRelationshipsRelationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disaster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disaster_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"synchronization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"precheck_failed_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"slave_region_instance_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbDrRelationshipsRegionInstanceInfoSchema(),
			},
			"master_region_instance_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbDrRelationshipsRegionInstanceInfoSchema(),
			},
		},
	}
}

func gaussDbDrRelationshipsRegionInstanceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbDrRelationshipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3.5/{project_id}/disaster-recovery/relations"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildGetGaussDbDrRelationshipsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB DR relationships: %s", err)
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("relations", flattenGetDrRelationsRelationshipsBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussDbDrRelationshipsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("instance_name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("dr_role"); ok {
		res = fmt.Sprintf("%s&dr_role=%v", res, v)
	}
	if v, ok := d.GetOk("dr_type"); ok {
		res = fmt.Sprintf("%s&dr_type=%v", res, v)
	}
	if v, ok := d.GetOk("dr_status"); ok {
		res = fmt.Sprintf("%s&dr_status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetDrRelationsRelationshipsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("relations", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"disaster_type":               utils.PathSearch("disaster_type", v, nil),
			"name":                        utils.PathSearch("name", v, nil),
			"disaster_role":               utils.PathSearch("disaster_role", v, nil),
			"created":                     utils.PathSearch("created", v, nil),
			"updated":                     utils.PathSearch("updated", v, nil),
			"id":                          utils.PathSearch("id", v, nil),
			"synchronization_id":          utils.PathSearch("synchronization_id", v, nil),
			"status":                      utils.PathSearch("status", v, nil),
			"precheck_failed_reason":      utils.PathSearch("precheck_failed_reason", v, nil),
			"instance_id":                 utils.PathSearch("instance_id", v, nil),
			"instance_name":               utils.PathSearch("instance_name", v, nil),
			"instance_status":             utils.PathSearch("instance_status", v, nil),
			"actions":                     utils.PathSearch("actions", v, nil),
			"slave_region_instance_info":  flattenGetDrRelationsRegionInstanceInfoBody(utils.PathSearch("slave_region_instance_info", v, nil)),
			"master_region_instance_info": flattenGetDrRelationsRegionInstanceInfoBody(utils.PathSearch("master_region_instance_info", v, nil)),
		})
	}
	return res
}

func flattenGetDrRelationsRegionInstanceInfoBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"region_code":  utils.PathSearch("region_code", resp, nil),
			"instance_id":  utils.PathSearch("instance_id", resp, nil),
			"project_id":   utils.PathSearch("project_id", resp, nil),
			"project_name": utils.PathSearch("project_name", resp, nil),
			"ip_address":   utils.PathSearch("ip_address", resp, nil),
		},
	}
}
