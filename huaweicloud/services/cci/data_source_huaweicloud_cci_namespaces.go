package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"
	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API CCI GET /api/v1/namespaces/{name}
// @API CCI GET /apis/networking.cci.io/v1beta1/namespaces/{ns}/networks
// @API CCI GET /api/v1/namespaces
func DataSourceCciNamespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCciNamespacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"general-computing", "gpu-accelerated",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_expend_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"warmup_pool_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"recycling_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"container_network_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rbac_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_cidr": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"network_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func isNamespaceParamMatch(d *schema.ResourceData, ns namespaces.Namespace) bool {
	if val, ok := d.GetOk("type"); ok {
		if val.(string) != ns.Metadata.Annotations.Flavor {
			return false
		}
	}
	if val, ok := d.GetOk("enterprise_project_id"); ok {
		if val.(string) != ns.Metadata.Labels.EnterpriseProjectID {
			return false
		}
	}
	return true
}

func flattenVpcNetwork(networkList []networks.Network) []map[string]interface{} {
	if len(networkList) < 1 {
		return nil
	}

	network := networkList[0]
	return []map[string]interface{}{
		{
			"name":              network.Metadata.Name,
			"security_group_id": network.Metadata.Annotations["network.alpha.kubernetes.io/default-security-group"],
			"vpc": []map[string]interface{}{
				{
					"id":          network.Spec.AttachedVPC,
					"subnet_id":   network.Spec.SubnetID,
					"subnet_cidr": network.Spec.Cidr,
					"network_id":  network.Spec.NetworkID,
				},
			},
		},
	}
}

func filterDataCciNamespaces(d *schema.ResourceData, client *golangsdk.ServiceClient,
	nsList []namespaces.Namespace) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, 0, len(nsList))
	ids := make([]string, 0)
	for _, ns := range nsList {
		if isNamespaceParamMatch(d, ns) {
			netList, err := networks.List(client, ns.Metadata.Name)
			if err != nil {
				return result, ids, err
			}
			nsParams := map[string]interface{}{
				"id":                        ns.Metadata.UID,
				"name":                      ns.Metadata.Name,
				"type":                      ns.Metadata.Annotations.Flavor,
				"enterprise_project_id":     ns.Metadata.Labels.EnterpriseProjectID,
				"auto_expend_enabled":       ns.Metadata.Annotations.AutoExpend,
				"warmup_pool_size":          ns.Metadata.Annotations.PoolSize,
				"recycling_interval":        ns.Metadata.Annotations.RecyclingInterval,
				"container_network_enabled": isContainNetworkEnabled(ns.Metadata.Annotations.NetworkEnable),
				"rbac_enabled":              ns.Metadata.Labels.RbacEnable,
				"created_at":                ns.Metadata.CreationTimestamp,
				"status":                    ns.Status.Phase,
				"network":                   flattenVpcNetwork(netList),
			}

			result = append(result, nsParams)
			ids = append(ids, ns.Metadata.UID)
		}
	}

	return result, ids, nil
}

func dataSourceCciNamespacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 client: %s", err)
	}
	betaClient, err := conf.CciV1BetaClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 beta1 client: %s", err)
	}

	var nsList []namespaces.Namespace
	if ns, ok := d.GetOk("name"); ok {
		nsResp, err := namespaces.Get(client, ns.(string)).Extract()
		if err != nil {
			return diag.Errorf("Error getting the namespace (%s) from the server: %s", ns.(string), err)
		}
		nsList = append(nsList, *nsResp)
	} else {
		pages, err := namespaces.List(client, namespaces.ListOpts{}).AllPages()
		if err != nil {
			return diag.Errorf("Error finding the namespace list from the server: %s", err)
		}
		nsList, err = namespaces.ExtractNamespaces(pages)
		if err != nil {
			return diag.Errorf("Error extracting CCI namespaces: %s", err)
		}
	}

	resp, ids, err := filterDataCciNamespaces(d, betaClient, nsList)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", resp),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("Error saving the namespace's (%v) fields to state: %s", ids, mErr)
	}
	return nil
}
