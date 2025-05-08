package ga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/health-checks
func DataSourceHealthChecks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHealthChecksRead,
		Schema: map[string]*schema.Schema{
			"health_check_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the health check.",
			},
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the endpoint group to which the health check belongs.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the health check.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The front end protocol of the health check used.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether health check is enabled.",
			},
			"health_checks": {
				Type:        schema.TypeList,
				Elem:        healthChecksSchema(),
				Computed:    true,
				Description: "The list of the health checks.",
			},
		},
	}
}

func healthChecksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the health check.",
			},
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the endpoint group to which the health check belongs.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the health check.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The front end protocol of the health check used.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the health check.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time interval of the health check.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timeout of the health check.",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The max retries of the health check.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether health check is enabled.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the health check.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the health check.",
			},
			"frozen_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The frozen details of cloud services or resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a cloud service or resource.`,
						},
						"effect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the resource after being forzen.`,
						},
						"scene": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The service scenario.`,
						},
					},
				},
			},
		},
	}
	return &sc
}

func dataSourceHealthChecksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/health-checks"
		product = "ga"
		mErr    *multierror.Error
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildListHealthChecksQueryParams(d)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving GA health checks: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("health_checks", filterListHealthChecksResponseBody(flattenListHealthChecksResponseBody(respBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListHealthChecksResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	var (
		protocol = d.Get("protocol").(string)
		enabled  = d.Get("enabled").(string)
		rst      = make([]interface{}, 0, len(all))
	)

	for _, v := range all {
		if protocol != "" && protocol != utils.PathSearch("protocol", v, "").(string) {
			continue
		}

		if enabled != "" && enabled != fmt.Sprint(utils.PathSearch("enabled", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenListHealthChecksResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("health_checks", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"endpoint_group_id": utils.PathSearch("endpoint_group_id", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"protocol":          utils.PathSearch("protocol", v, nil),
			"port":              utils.PathSearch("port", v, nil),
			"interval":          utils.PathSearch("interval", v, nil),
			"timeout":           utils.PathSearch("timeout", v, nil),
			"max_retries":       utils.PathSearch("max_retries", v, nil),
			"enabled":           utils.PathSearch("enabled", v, nil),
			"created_at":        utils.PathSearch("created_at", v, nil),
			"updated_at":        utils.PathSearch("updated_at", v, nil),
			"frozen_info":       flattenHealthChecksFrozenInfo(utils.PathSearch("frozen_info", v, nil)),
		})
	}
	return rst
}

func flattenHealthChecksFrozenInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	frozenInfo := map[string]interface{}{
		"status": utils.PathSearch("status", resp, nil),
		"effect": utils.PathSearch("effect", resp, nil),
		"scene":  utils.PathSearch("scene", resp, []string{}),
	}

	return []map[string]interface{}{frozenInfo}
}

func buildListHealthChecksQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("health_check_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("endpoint_group_id"); ok {
		res = fmt.Sprintf("%s&endpoint_group_id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
