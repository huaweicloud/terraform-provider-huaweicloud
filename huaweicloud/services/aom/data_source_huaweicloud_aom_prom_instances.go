package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v1/{project_id}/aom/prometheus
func DataSourceAomPromInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAomPromInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prom_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prom_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cce_cluster_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prom_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_write_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_read_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_http_api_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prom_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAomPromInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listInstancesHttpUrl := "v1/{project_id}/aom/prometheus"
	listInstancesHttpUrl = strings.ReplaceAll(listInstancesHttpUrl, "{project_id}", client.ProjectID)
	listInstancesPath := client.Endpoint + listInstancesHttpUrl + buildListInstancesQueryParams(d)
	listInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}
	listInstancesResp, err := client.Request("GET", listInstancesPath, &listInstancesOpt)
	if err != nil {
		return diag.Errorf("error retrieving AOM prometheus instances: %s", err)
	}

	listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
	if err != nil {
		return diag.Errorf("error flattening AOM prometheus instances response: %s", err)
	}

	instances := utils.PathSearch("prometheus", listInstancesRespBody, make([]interface{}, 0)).([]interface{})
	results := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		results = append(results, map[string]interface{}{
			"id":                     utils.PathSearch("prom_id", instance, nil),
			"prom_name":              utils.PathSearch("prom_name", instance, nil),
			"prom_type":              utils.PathSearch("prom_type", instance, nil),
			"prom_version":           utils.PathSearch("prom_version", instance, nil),
			"enterprise_project_id":  utils.PathSearch("enterprise_project_id", instance, nil),
			"remote_write_url":       utils.PathSearch("prom_spec_config.remote_write_url", instance, nil),
			"remote_read_url":        utils.PathSearch("prom_spec_config.remote_read_url", instance, nil),
			"prom_http_api_endpoint": utils.PathSearch("prom_spec_config.prom_http_api_endpoint", instance, nil),
			"prom_status":            utils.PathSearch("prom_status", instance, nil),
			"is_deleted_tag":         utils.PathSearch("is_deleted_tag", instance, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("prom_create_timestamp", instance, float64(0)).(float64))/1000, true),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("prom_update_timestamp", instance, float64(0)).(float64))/1000, true),
			"deleted_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("deleted_time", instance, float64(0)).(float64))/1000, true),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListInstancesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("prom_id"); ok {
		res = fmt.Sprintf("%s&prom_id=%v", res, v)
	}
	if v, ok := d.GetOk("prom_type"); ok {
		res = fmt.Sprintf("%s&prom_type=%v", res, v)
	}
	if v, ok := d.GetOk("cce_cluster_enable"); ok {
		res = fmt.Sprintf("%s&cce_cluster_enable=%v", res, v)
	}
	if v, ok := d.GetOk("prom_status"); ok {
		res = fmt.Sprintf("%s&prom_status=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
