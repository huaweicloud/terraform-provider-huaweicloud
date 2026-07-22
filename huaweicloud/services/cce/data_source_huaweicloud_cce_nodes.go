package cce

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{clusterid}/nodes
// @API ECS GET /v1/{project_id}/cloudservers/{id}/tags
func DataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ignore_details": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, _ string) ([]string, []error) {
					validValues := []string{"tags"}
					params := strings.Split(v.(string), ",")
					if !utils.StrSliceContainsAnother(validValues, params) {
						return nil, []error{fmt.Errorf("the value must within %s", validValues)}
					}
					return nil, nil
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_pair": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_volume": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volumetype": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"extend_params": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								}},
						},
						"data_volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volumetype": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"extend_params": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								}},
						},
						"billing_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CceV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	items, err := getNodesOfCCECluster(client, clusterId)
	if err != nil {
		return diag.Errorf("error retrieving nodes of CCE Cluster (%s): %s", clusterId, err)
	}

	filteredItems := filterNodesWithParameters(items, d)

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	ids := make([]string, 0)
	for _, item := range filteredItems {
		uid := utils.PathSearch("metadata.uid", item, "").(string)
		if uid != "" {
			ids = append(ids, uid)
		}
	}
	nodesToSet := flattenListNodesResponseBody(filteredItems, d, cfg, region)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nodes", nodesToSet),
		d.Set("ids", ids),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getNodesOfCCECluster(client *golangsdk.ServiceClient, clusterID string) ([]interface{}, error) {
	var (
		httpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/nodes"
		limit   = 2000
		marker  = ""
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := listPath + buildListNodesQueryParams(limit, marker)

		requestResp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)
		if len(items) < limit {
			break
		}
		marker = utils.PathSearch("pageInfo.nextMarker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func buildListNodesQueryParams(limit int, marker string) string {
	res := fmt.Sprintf("?limit=%d", limit)
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func filterNodesWithParameters(items []interface{}, d *schema.ResourceData) []interface{} {
	if len(items) == 0 {
		return items
	}

	var targetID, targetName, targetStatus string
	if val, ok := d.GetOk("node_id"); ok {
		targetID, _ = val.(string)
	}
	if val, ok := d.GetOk("name"); ok {
		targetName, _ = val.(string)
	}
	if val, ok := d.GetOk("status"); ok {
		targetStatus, _ = val.(string)
	}

	if targetID == "" && targetName == "" && targetStatus == "" {
		return items
	}

	var filteredItems []interface{}
	for _, item := range items {
		itemId := utils.PathSearch("metadata.uid", item, "").(string)
		itemName := utils.PathSearch("metadata.name", item, "").(string)
		itemStatus := utils.PathSearch("status.phase", item, "").(string)

		if targetID != "" && itemId != targetID {
			continue
		}
		if targetName != "" && itemName != targetName {
			continue
		}
		if targetStatus != "" && itemStatus != targetStatus {
			continue
		}
		filteredItems = append(filteredItems, item)
	}

	return filteredItems
}

func flattenListNodesResponseBody(items []interface{}, d *schema.ResourceData, cfg *config.Config, region string) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		node := map[string]interface{}{
			"id":                    utils.PathSearch("metadata.uid", item, nil),
			"name":                  utils.PathSearch("metadata.name", item, nil),
			"flavor_id":             utils.PathSearch("spec.flavor", item, nil),
			"availability_zone":     utils.PathSearch("spec.az", item, nil),
			"os":                    utils.PathSearch("spec.os", item, nil),
			"billing_mode":          utils.PathSearch("spec.billingMode", item, nil),
			"key_pair":              utils.PathSearch("spec.login.sshKey", item, nil),
			"subnet_id":             utils.PathSearch("spec.nodeNicSpec.primaryNic.subnetId", item, nil),
			"ecs_group_id":          utils.PathSearch("spec.ecsGroupId", item, nil),
			"server_id":             utils.PathSearch("status.serverId", item, nil),
			"public_ip":             utils.PathSearch("status.publicIP", item, nil),
			"private_ip":            utils.PathSearch("status.privateIP", item, nil),
			"status":                utils.PathSearch("status.phase", item, nil),
			"enterprise_project_id": utils.PathSearch("spec.serverEnterpriseProjectID", item, nil),
		}

		// Flatten hostname_config
		hostnameConfig := flattenNodeHostnameConfigFromSearch(utils.PathSearch("spec.hostnameConfig", item, nil))
		node["hostname_config"] = hostnameConfig

		// Flatten root_volume
		rootVolume := flattenNodeRootVolumeFromSearch(utils.PathSearch("spec.rootVolume", item, nil))
		node["root_volume"] = rootVolume

		// Flatten data_volumes
		dataVolumes := flattenNodeDataVolumesFromSearch(utils.PathSearch("spec.dataVolumes", item, make([]interface{}, 0)))
		node["data_volumes"] = dataVolumes

		// Fetch tags from ECS instance
		if !strings.Contains(d.Get("ignore_details").(string), "tags") {
			computeClient, err := cfg.ComputeV1Client(region)
			if err == nil {
				serverId := utils.PathSearch("status.serverId", item, "").(string)
				if serverId != "" {
					if resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract(); err == nil {
						tagsMap := utils.TagsToMap(resourceTags.Tags)
						node["tags"] = tagsMap
					} else {
						log.Printf("[WARN] Error fetching tags of CCE Node (%s): %s", serverId, err)
					}
				}
			}
		}

		res = append(res, node)
	}

	return res
}

func flattenNodeHostnameConfigFromSearch(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"type": utils.PathSearch("type", raw, nil),
		},
	}
}

func flattenNodeRootVolumeFromSearch(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"size":          utils.PathSearch("size", raw, nil),
			"volumetype":    utils.PathSearch("volumetype", raw, nil),
			"extend_params": utils.PathSearch("extendParam", raw, nil),
		},
	}
}

func flattenNodeDataVolumesFromSearch(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}
	curArray, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"size":          utils.PathSearch("size", v, nil),
			"volumetype":    utils.PathSearch("volumetype", v, nil),
			"extend_params": utils.PathSearch("extendParam", v, nil),
		})
	}
	return res
}
