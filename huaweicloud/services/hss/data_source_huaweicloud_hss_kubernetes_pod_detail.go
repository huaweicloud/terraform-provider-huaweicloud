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

// @API HSS GET /v5/{project_id}/kubernetes/{pod_name}/pod/detail
func DataSourceKubernetesPodDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesPodDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name_attr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
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
			"node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pod_ip": {
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
			"containers": {
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

func buildKubernetesPodDetailQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceKubernetesPodDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		podName = d.Get("pod_name").(string)
		httpUrl = "v5/{project_id}/kubernetes/{pod_name}/pod/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{pod_name}", podName)
	requestPath += buildKubernetesPodDetailQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS kubernetes pod detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pod_name_attr", utils.PathSearch("pod_name", respBody, nil)),
		d.Set("namespace_name", utils.PathSearch("namespace_name", respBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", respBody, nil)),
		d.Set("node_name", utils.PathSearch("node_name", respBody, nil)),
		d.Set("label", utils.PathSearch("label", respBody, nil)),
		d.Set("cpu", utils.PathSearch("cpu", respBody, nil)),
		d.Set("memory", utils.PathSearch("memory", respBody, nil)),
		d.Set("cpu_limit", utils.PathSearch("cpu_limit", respBody, nil)),
		d.Set("memory_limit", utils.PathSearch("memory_limit", respBody, nil)),
		d.Set("node_ip", utils.PathSearch("node_ip", respBody, nil)),
		d.Set("pod_ip", utils.PathSearch("pod_ip", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("containers", flattenKubernetesPodDetailContainers(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKubernetesPodDetailContainers(resp interface{}) []map[string]interface{} {
	containersResp := utils.PathSearch("containers", resp, make([]interface{}, 0)).([]interface{})
	if len(containersResp) == 0 {
		return nil
	}

	containersResult := make([]map[string]interface{}, 0, len(containersResp))
	for _, v := range containersResp {
		containersResult = append(containersResult, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"region_id":      utils.PathSearch("region_id", v, nil),
			"container_id":   utils.PathSearch("container_id", v, nil),
			"container_name": utils.PathSearch("container_name", v, nil),
			"image_name":     utils.PathSearch("image_name", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"cpu_limit":      utils.PathSearch("cpu_limit", v, nil),
			"memory_limit":   utils.PathSearch("memory_limit", v, nil),
			"restart_count":  utils.PathSearch("restart_count", v, nil),
			"pod_name":       utils.PathSearch("pod_name", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"cluster_type":   utils.PathSearch("cluster_type", v, nil),
			"risky":          utils.PathSearch("risky", v, nil),
			"low_risk":       utils.PathSearch("low_risk", v, nil),
			"medium_risk":    utils.PathSearch("medium_risk", v, nil),
			"high_risk":      utils.PathSearch("high_risk", v, nil),
			"fatal_risk":     utils.PathSearch("fatal_risk", v, nil),
		})
	}

	return containersResult
}
