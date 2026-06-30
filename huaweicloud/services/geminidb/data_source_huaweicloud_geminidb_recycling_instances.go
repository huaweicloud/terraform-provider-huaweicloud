package geminidb

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/instances/recycle-instances
func DataSourceGeminiDBRecyclingInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBRecyclingInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBRecyclingInstancesInstanceSchema(),
			},
		},
	}
}

func geminiDBRecyclingInstancesInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_store": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBRecyclingInstancesDataStoreSchema(),
			},
			"charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retained_until": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func geminiDBRecyclingInstancesDataStoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGeminiDBRecyclingInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/recycle-instances"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB recycling instances: %s", err)
	}

	listRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	instances := flattenListGeminiDBRecyclingInstances(listRespBody)
	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instances", instances),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGeminiDBRecyclingInstances(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"mode":                  utils.PathSearch("mode", v, nil),
			"product_type":          utils.PathSearch("product_type", v, nil),
			"data_store":            flattenGeminiDBRecyclingInstancesDataStore(v),
			"charge_type":           utils.PathSearch("charge_type", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"backup_id":             utils.PathSearch("backup_id", v, nil),
			"created_at":            utils.PathSearch("created_at", v, nil),
			"deleted_at":            utils.PathSearch("deleted_at", v, nil),
			"retained_until":        utils.PathSearch("retained_until", v, nil),
		})
	}
	return rst
}

func flattenGeminiDBRecyclingInstancesDataStore(instance interface{}) []map[string]interface{} {
	dataStoreRaw := utils.PathSearch("data_store", instance, nil)
	if dataStoreRaw == nil {
		return nil
	}

	dataStore := map[string]interface{}{
		"type":    utils.PathSearch("type", dataStoreRaw, nil),
		"version": utils.PathSearch("version", dataStoreRaw, nil),
	}

	return []map[string]interface{}{dataStore}
}
