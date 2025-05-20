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

func DataSourceRdsConfigurationHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsConfigurationHistoriesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"histories": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"new_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"applied": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceRdsConfigurationHistoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/configuration-histories"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var histories []interface{}
	offset := 0
	for {
		path := basePath + buildFullQueryParams(d, offset)
		getResp, err := client.Request("GET", path, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS configuration histories: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		parsed := flattenConfigurationHistoriesBody(getRespBody)
		if len(parsed) == 0 {
			break
		}

		histories = append(histories, parsed...)
		offset += 10
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("histories", histories),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildFullQueryParams(d *schema.ResourceData, offset int) string {
	params := fmt.Sprintf("?limit=10&offset=%d", offset)
	if v, ok := d.GetOk("start_time"); ok {
		params += fmt.Sprintf("&start_time=%v", v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		params += fmt.Sprintf("&end_time=%v", v)
	}
	if v, ok := d.GetOk("param_name"); ok {
		params += fmt.Sprintf("&param_name=%v", v)
	}
	return params
}

func flattenConfigurationHistoriesBody(resp interface{}) []interface{} {
	historiesJSON := utils.PathSearch("histories", resp, make([]interface{}, 0))
	historiesArr := historiesJSON.([]interface{})

	if len(historiesArr) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(historiesArr))
	for _, h := range historiesArr {
		result = append(result, map[string]interface{}{
			"parameter_name": utils.PathSearch("parameter_name", h, nil),
			"old_value":      utils.PathSearch("old_value", h, nil),
			"new_value":      utils.PathSearch("new_value", h, nil),
			"update_result":  utils.PathSearch("update_result", h, nil),
			"applied":        utils.PathSearch("applied", h, nil),
			"update_time":    utils.PathSearch("update_time", h, nil),
			"apply_time":     utils.PathSearch("apply_time", h, nil),
		})
	}
	return result
}
