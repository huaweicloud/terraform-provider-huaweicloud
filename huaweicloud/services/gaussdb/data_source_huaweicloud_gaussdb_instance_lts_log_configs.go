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

// @API GaussDB GET /v3/{project_id}/instances/logs/lts-config
func DataSourceGaussDBInstanceLtsLogConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInstanceLtsLogConfigsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_lts_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceLtsConfigsSchema(),
			},
		},
	}
}

func instanceLtsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceLtsDetailSchema(),
			},
			"lts_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ltsConfigSchema(),
			},
		},
	}
}

func instanceLtsDetailSchema() *schema.Resource {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     datastoreSchema(),
			},
			"frozen_flag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datastoreSchema() *schema.Resource {
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

func ltsConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lts_stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBInstanceLtsLogConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/logs/lts-config"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildGetInstanceLtsLogConfigsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance LTS log configs: %s", err)
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
		d.Set("instance_lts_configs", flattenGetInstanceLtsLogConfigsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceLtsLogConfigsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_mode"); ok {
		res = fmt.Sprintf("%s&instance_mode=%v", res, v)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetInstanceLtsLogConfigsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance_lts_configs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance":    flattenGetInstanceLtsDetailBody(v),
			"lts_configs": flattenGetLtsConfigsBody(v),
		})
	}
	return res
}

func flattenGetInstanceLtsDetailBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":                    utils.PathSearch("id", curJson, nil),
			"name":                  utils.PathSearch("name", curJson, nil),
			"mode":                  utils.PathSearch("mode", curJson, nil),
			"status":                utils.PathSearch("status", curJson, nil),
			"datastore":             flattenGetDatastoreBody(curJson),
			"frozen_flag":           utils.PathSearch("frozen_flag", curJson, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", curJson, nil),
			"actions":               utils.PathSearch("actions", curJson, []interface{}{}),
		},
	}
}

func flattenGetDatastoreBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("datastore", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("type", curJson, nil),
			"version": utils.PathSearch("version", curJson, nil),
		},
	}
}

func flattenGetLtsConfigsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("lts_configs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"log_type":      utils.PathSearch("log_type", v, nil),
			"lts_group_id":  utils.PathSearch("lts_group_id", v, nil),
			"lts_stream_id": utils.PathSearch("lts_stream_id", v, nil),
			"enabled":       utils.PathSearch("enabled", v, nil),
		})
	}
	return res
}
