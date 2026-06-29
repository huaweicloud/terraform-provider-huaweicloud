package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/starrocks/instances/logs/lts-configs
func DataSourceTaurusDBHtapStarrocksLtsConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksLtsConfigsRead,

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
				Elem:     htapStarrocksLtsConfigsSchema(),
			},
		},
	}
}

func htapStarrocksLtsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapStarrocksLtsInstanceSchema(),
			},
			"lts_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapStarrocksLtsConfigSchema(),
			},
		},
	}
}

func htapStarrocksLtsInstanceSchema() *schema.Resource {
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
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func htapStarrocksLtsConfigSchema() *schema.Resource {
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

func dataSourceTaurusDBHtapStarrocksLtsConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/starrocks/instances/logs/lts-configs"
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildGetTaurusDBHtapStarrocksLtsConfigsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB HTAP StarRocks LTS configs: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.Errorf("error marshaling TaurusDB HTAP StarRocks LTS configs: %s", err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.Errorf("error unmarshalling TaurusDB HTAP StarRocks LTS configs: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_lts_configs", flattenTaurusDBHtapStarrocksLtsConfigs(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTaurusDBHtapStarrocksLtsConfigsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
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

func flattenTaurusDBHtapStarrocksLtsConfigs(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance_lts_configs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance":    flattenTaurusDBHtapStarrocksInstance(v),
			"lts_configs": flattenTaurusDBHtapStarrocksLtsConfigList(v),
		})
	}
	return res
}

func flattenTaurusDBHtapStarrocksInstance(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":                      utils.PathSearch("id", curJson, nil),
			"name":                    utils.PathSearch("name", curJson, nil),
			"mode":                    utils.PathSearch("mode", curJson, nil),
			"engine_name":             utils.PathSearch("engine_name", curJson, nil),
			"engine_version":          utils.PathSearch("engine_version", curJson, nil),
			"status":                  utils.PathSearch("status", curJson, nil),
			"enterprise_project_id":   utils.PathSearch("enterprise_project_id", curJson, nil),
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", curJson, nil),
		},
	}
}

func flattenTaurusDBHtapStarrocksLtsConfigList(resp interface{}) []interface{} {
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
