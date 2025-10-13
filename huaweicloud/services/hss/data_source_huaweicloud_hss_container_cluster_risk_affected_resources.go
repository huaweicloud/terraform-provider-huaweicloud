package hss

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

// Due to the lack of the `risk_id` parameter, this resource cannot be tested.

// @API HSS GET /v5/{project_id}/container/cluster/risk/{risk_id}/affected-resources
func DataSourceContainerClusterRiskAffectedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerClusterRiskAffectedResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"risk_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hit_rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hit_path_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"first_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerClusterRiskAffectedResourcesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("resource_name"); ok {
		queryParams = fmt.Sprintf("%s&resource_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("resource_type"); ok {
		queryParams = fmt.Sprintf("%s&resource_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerClusterRiskAffectedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/cluster/risk/{risk_id}/affected-resources"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{risk_id}", d.Get("risk_id").(string))
	requestPath += buildContainerClusterRiskAffectedResourcesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container cluster risk affected resources: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenContainerClusterRiskAffectedResources(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerClusterRiskAffectedResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_type":   utils.PathSearch("resource_type", v, nil),
			"namespace":       utils.PathSearch("namespace", v, nil),
			"hit_rule":        utils.PathSearch("hit_rule", v, nil),
			"hit_path_list":   utils.ExpandToStringList(utils.PathSearch("hit_path_list", v, make([]interface{}, 0)).([]interface{})),
			"first_scan_time": utils.PathSearch("first_scan_time", v, nil),
			"last_scan_time":  utils.PathSearch("last_scan_time", v, nil),
		})
	}

	return rst
}
