package dds

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS GET /v3/{project_id}/instances/{instance_id}/kill-op-rule
func DataSourceKillOpRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKillOpRulesRead,

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
			"operation_types": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespaces": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_summary": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_types": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespaces": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ips": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"secs_running": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildKillOpRulesQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("operation_types"); ok {
		queryParams = fmt.Sprintf("%s&operation_types=%v", queryParams, v)
	}

	if v, ok := d.GetOk("namespaces"); ok {
		queryParams = fmt.Sprintf("%s&namespaces=%v", queryParams, v)
	}

	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}

	if v, ok := d.GetOk("plan_summary"); ok {
		queryParams = fmt.Sprintf("%s&plan_summary=%v", queryParams, v)
	}

	return queryParams
}

func listKillOpRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/kill-op-rule?limit={limit}"
		instanceId = d.Get("instance_id").(string)
		limit      = 100
		offset     = 0
		result     = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildKillOpRulesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		dataResp := utils.PathSearch("rules", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, dataResp...)
		if len(dataResp) < limit {
			break
		}

		offset += len(dataResp)
	}

	return result, nil
}

func dataSourceKillOpRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	resp, err := listKillOpRules(client, d)
	if err != nil {
		return diag.Errorf("error retrieving killOp rules: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenKillOpRules(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKillOpRules(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"operation_types": utils.PathSearch("operation_types", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"namespaces":      utils.PathSearch("namespaces", v, nil),
			"client_ips":      utils.PathSearch("client_ips", v, nil),
			"plan_summary":    utils.PathSearch("plan_summary", v, nil),
			"node_type":       utils.PathSearch("node_type", v, nil),
			"max_concurrency": utils.PathSearch("max_concurrency", v, nil),
			"secs_running":    utils.PathSearch("secs_running", v, nil),
		})
	}

	return result
}
