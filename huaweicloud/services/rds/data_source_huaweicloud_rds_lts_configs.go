package rds

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

// @API RDS GET /v3/{project_id}/{engine}/instances/logs/lts-configs
func DataSourceRdsLtsConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsLtsConfigsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_lts_configs": {
				Type:     schema.TypeList,
				Elem:     instanceLtsConfigsSchema(),
				Computed: true,
			},
		},
	}
}

func instanceLtsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlAccountsRead,
		Schema: map[string]*schema.Schema{
			"lts_configs": {
				Type:     schema.TypeList,
				Elem:     lstConfigsSchema(),
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeList,
				Elem:     instanceInfoSchema(),
				Computed: true,
			},
		},
	}
}

func lstConfigsSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlAccountsRead,
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

func instanceInfoSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlAccountsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
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
			"engine_category": {
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
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceRdsLtsConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	basePath := client.Endpoint + httpUrl + buildLtsConfigsQueryParams(d)
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{engine}", d.Get("engine").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", basePath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS LTS configs: %s", err)
	}

	body, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	flatted := flattenInstanceLtsConfigs(body)

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_lts_configs", flatted),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildLtsConfigsQueryParams(d *schema.ResourceData) string {
	query := ""

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		query += fmt.Sprintf("&enterprise_project_id=%v", v)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		query += fmt.Sprintf("&instance_id=%v", v)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		query += fmt.Sprintf("&instance_name=%v", v)
	}

	if v, ok := d.GetOk("sort"); ok {
		query += fmt.Sprintf("&sort=%v", v)
	}

	if v, ok := d.GetOk("instance_status"); ok {
		query += fmt.Sprintf("&instance_status=%v", v)
	}

	if query == "" {
		return ""
	}

	return "?" + query[1:]
}

func flattenInstanceLtsConfigs(resp interface{}) []interface{} {
	rawList := utils.PathSearch("instance_lts_configs", resp, nil)
	if rawList == nil {
		return nil
	}

	cur, ok := rawList.([]interface{})
	if !ok || len(cur) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(cur))
	for _, v := range cur {
		out = append(out, map[string]interface{}{
			"lts_configs": flattenLtsConfigs(utils.PathSearch("lts_configs", v, nil)),
			"instance":    flattenInstanceBasic(utils.PathSearch("instance", v, nil)),
		})
	}
	return out
}

func flattenLtsConfigs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	arr, ok := resp.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		out = append(out, map[string]interface{}{
			"log_type":      utils.PathSearch("log_type", v, nil),
			"lts_group_id":  utils.PathSearch("lts_group_id", v, nil),
			"lts_stream_id": utils.PathSearch("lts_stream_id", v, nil),
			"enabled":       utils.PathSearch("enabled", v, nil),
		})
	}
	return out
}

func flattenInstanceBasic(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":                    utils.PathSearch("id", resp, nil),
			"name":                  utils.PathSearch("name", resp, nil),
			"engine_name":           utils.PathSearch("engine_name", resp, nil),
			"engine_version":        utils.PathSearch("engine_version", resp, nil),
			"engine_category":       utils.PathSearch("engine_category", resp, nil),
			"status":                utils.PathSearch("status", resp, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", resp, nil),
			"actions":               utils.PathSearch("actions", resp, nil),
		},
	}
}
