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

// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/risks/query
func DataSourceContainerKubernetesClustersRisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesClustersRisksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"detect_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"images_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"baseline_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vul_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"event_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protect_node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charging_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesClustersRisksQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildContainerKubernetesClustersRisksBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"cluster_id_list": utils.ExpandToStringList(d.Get("cluster_id_list").([]interface{})),
		"detect_type":     utils.ValueIgnoreEmpty(d.Get("detect_type")),
	}
}

func dataSourceContainerKubernetesClustersRisksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container/kubernetes/clusters/risks/query"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerKubernetesClustersRisksQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildContainerKubernetesClustersRisksBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS container kubernetes clusters risks: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("data_list", flattenContainerKubernetesClustersRisksDataList(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesClustersRisksDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"images_num":        utils.PathSearch("images_num", v, nil),
			"baseline_risk_num": utils.PathSearch("baseline_risk_num", v, nil),
			"vul_num":           utils.PathSearch("vul_num", v, nil),
			"event_num":         utils.PathSearch("event_num", v, nil),
			"protect_node_num":  utils.PathSearch("protect_node_num", v, nil),
			"node_total_num":    utils.PathSearch("node_total_num", v, nil),
			"cluster_id":        utils.PathSearch("cluster_id", v, nil),
			"charging_mode":     utils.PathSearch("charging_mode", v, nil),
		})
	}

	return rst
}
