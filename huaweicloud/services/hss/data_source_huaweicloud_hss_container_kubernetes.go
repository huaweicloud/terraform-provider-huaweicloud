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

// @API HSS GET /v5/{project_id}/container/kubernetes
func DataSourceContainerKubernetes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_container": {
				Type:     schema.TypeBool,
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
			"last_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pod_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risky": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"low_risk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"medium_risk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"high_risk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fatal_risk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=100"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("container_name"); ok {
		queryParams = fmt.Sprintf("%s&container_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("pod_name"); ok {
		queryParams = fmt.Sprintf("%s&pod_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("image_name"); ok {
		queryParams = fmt.Sprintf("%s&image_name=%v", queryParams, v)
	}
	if d.Get("cluster_container").(bool) {
		queryParams = fmt.Sprintf("%s&cluster_container=%v", queryParams, true)
	}

	return queryParams
}

func dataSourceContainerKubernetesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		epsId          = cfg.GetEnterpriseProjectID(d)
		product        = "hss"
		httpUrl        = "v5/{project_id}/container/kubernetes"
		offset         = 0
		result         = make([]interface{}, 0)
		totalNum       float64
		lastUpdateTime float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildContainerKubernetesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS container kubernetes: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		lastUpdateTime = utils.PathSearch("last_update_time", getRespBody, float64(0)).(float64)
		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
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
		d.Set("total_num", totalNum),
		d.Set("last_update_time", lastUpdateTime),
		d.Set("data_list", flattenContainerKubernetesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesDataList(containerList []interface{}) []interface{} {
	if len(containerList) == 0 {
		return nil
	}

	result := make([]interface{}, len(containerList))
	for i, container := range containerList {
		result[i] = map[string]interface{}{
			"id":             utils.PathSearch("id", container, nil),
			"region_id":      utils.PathSearch("region_id", container, nil),
			"container_id":   utils.PathSearch("container_id", container, nil),
			"container_name": utils.PathSearch("container_name", container, nil),
			"image_name":     utils.PathSearch("image_name", container, nil),
			"status":         utils.PathSearch("status", container, nil),
			"create_time":    utils.PathSearch("create_time", container, nil),
			"cpu_limit":      utils.PathSearch("cpu_limit", container, nil),
			"memory_limit":   utils.PathSearch("memory_limit", container, nil),
			"restart_count":  utils.PathSearch("restart_count", container, nil),
			"pod_name":       utils.PathSearch("pod_name", container, nil),
			"cluster_name":   utils.PathSearch("cluster_name", container, nil),
			"cluster_id":     utils.PathSearch("cluster_id", container, nil),
			"cluster_type":   utils.PathSearch("cluster_type", container, nil),
			"risky":          utils.PathSearch("risky", container, nil),
			"low_risk":       utils.PathSearch("low_risk", container, nil),
			"medium_risk":    utils.PathSearch("medium_risk", container, nil),
			"high_risk":      utils.PathSearch("high_risk", container, nil),
			"fatal_risk":     utils.PathSearch("fatal_risk", container, nil),
		}
	}

	return result
}
