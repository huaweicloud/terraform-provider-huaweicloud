package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/cluster
func DataSourceCpcsClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCpcsClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the cluster name.`,
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the service type of the cluster.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sort direction.`,
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the clusters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task ID.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The account ID.`,
						},
						"ccsp_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The CCSP cluster ID.`,
						},
						"distributed_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The distribution type.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster ID.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster name.`,
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service type of the cluster.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the cluster.`,
						},
						"instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of service instances in the cluster.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the cluster.`,
						},
						"progress_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The progress information.`,
						},
						"vsm_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of VSM instances used by the cluster.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation time of the cluster.`,
						},
						"shared_ccsp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the CCSP is shared.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID.`,
						},
						"az": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone.`,
						},
					},
				},
			},
		},
	}
}

func buildDataSourceCpcsClustersQueryParams(d *schema.ResourceData, pageNum int) string {
	rst := fmt.Sprintf("?page_num=%d", pageNum)

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("service_type"); ok {
		rst += fmt.Sprintf("&service_type=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	return rst
}

func dataSourceCpcsClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/dew/cpcs/cluster"
		product     = "kms"
		pageNum     = 1
		allClusters = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithPageNum := requestPath + buildDataSourceCpcsClustersQueryParams(d, pageNum)
		resp, err := client.Request("GET", requestPathWithPageNum, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DEW CPCS clusters: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.Errorf("error flattening DEW CPCS clusters response: %s", err)
		}

		results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(results) == 0 {
			break
		}
		allClusters = append(allClusters, results...)

		totalNum := int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		if len(allClusters) >= totalNum {
			break
		}

		pageNum++
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("clusters", flattenCpcsClustersResponseBody(allClusters)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCpcsClustersResponseBody(clusters []interface{}) []interface{} {
	if len(clusters) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		result = append(result, map[string]interface{}{
			"task_id":               utils.PathSearch("task_id", cluster, nil),
			"project_id":            utils.PathSearch("project_id", cluster, nil),
			"domain_id":             utils.PathSearch("domain_id", cluster, nil),
			"ccsp_id":               utils.PathSearch("ccsp_id", cluster, nil),
			"distributed_type":      utils.PathSearch("distributed_type", cluster, nil),
			"cluster_id":            utils.PathSearch("cluster_id", cluster, nil),
			"cluster_name":          utils.PathSearch("cluster_name", cluster, nil),
			"service_type":          utils.PathSearch("service_type", cluster, nil),
			"type":                  utils.PathSearch("type", cluster, nil),
			"instance_num":          utils.PathSearch("instance_num", cluster, nil),
			"status":                utils.PathSearch("status", cluster, nil),
			"progress_info":         utils.PathSearch("progress_info", cluster, nil),
			"vsm_num":               utils.PathSearch("vsm_num", cluster, nil),
			"create_time":           utils.PathSearch("create_time", cluster, nil),
			"shared_ccsp":           utils.PathSearch("shared_ccsp", cluster, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", cluster, nil),
			"az":                    utils.PathSearch("az", cluster, nil),
		})
	}

	return result
}
