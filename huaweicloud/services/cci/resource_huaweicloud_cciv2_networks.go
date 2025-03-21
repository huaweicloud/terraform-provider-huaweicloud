package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
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
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the namespace.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the namespace.`,
						},
						"annotations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The annotations of the namespace.`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The labels of the namespace.`,
						},
						"creation_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation timestamp of the namespace.`,
						},
						"resource_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource version of the namespace.`,
						},
						"self_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The self link of the namespace.`,
						},
						"uid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uid of the namespace.`,
						},
						"ip_families": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the IP families of the CCI network.`,
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the security group IDs of the CCI network.`,
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Specifies the subnet ID of the CCI network.`,
									},
								},
							},
							Description: `Specifies the subnets of the CCI network.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the namespace.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2NetworksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	// client, err := conf.CciV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCI v2 client: %s", err)
	}

	listNetworkHttpUrl := "apis/cci/v2/namespaces"
	listNetworkPath := client.Endpoint + listNetworkHttpUrl
	listNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listNetworksResp, err := client.Request("GET", listNetworkPath, &listNetworkOpt)
	if err != nil {
		return diag.Errorf("error querying CCI networks: %s", err)
	}

	listNetworksRespBody, err := utils.FlattenResponse(listNetworksResp)
	if err != nil {
		return diag.Errorf("error finding the networks list from the server: %s", err)
	}
	networks := utils.PathSearch("items", listNetworksRespBody, make([]interface{}, 0)).([]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", flattenNetworks(networks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNetworks(networks []interface{}) []interface{} {
	if len(networks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(networks))
	for _, v := range networks {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"self_link":          utils.PathSearch("metadata.selfLink", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"ip_families":        utils.PathSearch("spec.ipFamilies", v, nil),
			"security_group_ids": utils.PathSearch("spec.securityGroups", v, nil),
			"subnets":            utils.PathSearch("spec.subnets", v, nil),
			"status":             utils.PathSearch("status.state", v, nil),
		})
	}
	return rst
}
