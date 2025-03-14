package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/yangtse/v2/namespaces/{namespace}/networks
func DataSourceV2Networks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2NetworksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_families": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     networksStatusSchema(),
						},
					},
				},
			},
		},
	}
}

func networksStatusSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_transition_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnet_attrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_v4_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_v6_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func dataSourceV2NetworksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	listNetworksHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks"
	listNetworksPath := client.Endpoint + listNetworksHttpUrl
	listNetworksPath = strings.ReplaceAll(listNetworksPath, "{namespace}", d.Get("namespace").(string))
	listNetworksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listNetworksResp, err := client.Request("GET", listNetworksPath, &listNetworksOpt)
	if err != nil {
		return diag.Errorf("error querying CCI networks: %s", err)
	}

	listNetworksRespBody, err := utils.FlattenResponse(listNetworksResp)
	if err != nil {
		return diag.Errorf("error getting the networks list from the server: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	networks := utils.PathSearch("items", listNetworksRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("networks", flattenNetworks(networks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworks(networks []interface{}) []interface{} {
	if len(networks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(networks))
	for _, v := range networks {
		subnets := utils.PathSearch("spec.subnets", v, make([]interface{}, 0)).([]interface{})
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"finalizers":         utils.PathSearch("metadata.finalizers", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"ip_families":        utils.PathSearch("spec.ipFamilies", v, nil),
			"security_group_ids": utils.PathSearch("spec.securityGroups", v, nil),
			"subnets":            flattenNetworkSubnets(subnets),
			"status":             flattenNetworkStatus(utils.PathSearch("status", v, nil)),
		})
	}
	return rst
}
