package cce

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCCEClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClustersV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clusters": {
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
						"flavor_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"billing_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_network_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_subnet_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authentication_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_network_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"masters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"availability_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"kube_config_raw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_clusters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate_authority_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"certificate_users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_certificate_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_key_data": {
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
	}
}

func dataSourceCCEClustersV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := clusters.ListOpts{
		ID:                  d.Get("cluster_id").(string),
		Name:                d.Get("name").(string),
		Type:                d.Get("cluster_type").(string),
		Phase:               d.Get("status").(string),
		VpcID:               d.Get("vpc_id").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
	}

	refinedClusters, err := clusters.List(cceClient, listOpts)
	logp.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve clusters: %s", err)
	}

	ids := make([]string, 0, len(refinedClusters))
	clustersToSet := make([]map[string]interface{}, 0, len(refinedClusters))

	for _, v := range refinedClusters {

		ids = append(ids, v.Metadata.Id)

		cluster := map[string]interface{}{
			"name":                   v.Metadata.Name,
			"id":                     v.Metadata.Id,
			"status":                 v.Status.Phase,
			"flavor_id":              v.Spec.Flavor,
			"cluster_version":        v.Spec.Version,
			"cluster_type":           v.Spec.Type,
			"description":            v.Spec.Description,
			"billing_mode":           v.Spec.BillingMode,
			"vpc_id":                 v.Spec.HostNetwork.VpcId,
			"subnet_id":              v.Spec.HostNetwork.SubnetId,
			"container_network_cidr": v.Spec.ContainerNetwork.Cidr,
			"container_network_type": v.Spec.ContainerNetwork.Mode,
			"eni_subnet_id":          v.Spec.EniNetwork.SubnetId,
			"eni_subnet_cidr":        v.Spec.EniNetwork.Cidr,
			"authentication_mode":    v.Spec.Authentication.Mode,
			"security_group_id":      v.Spec.HostNetwork.SecurityGroup,
			"enterprise_project_id":  v.Spec.ExtendParam["enterpriseProjectId"],
			"service_network_cidr":   v.Spec.KubernetesSvcIPRange,
		}

		var endpoints []map[string]interface{}
		for _, endpoint := range v.Status.Endpoints {

			mapping := map[string]interface{}{
				"url":  endpoint.Url,
				"type": endpoint.Type,
			}
			endpoints = append(endpoints, mapping)
		}
		cluster["endpoints"] = endpoints

		r := clusters.GetCert(cceClient, v.Metadata.Id)

		kubeConfigRaw, err := utils.JsonMarshal(r.Body)

		if err != nil {
			logp.Printf("Error marshaling r.Body: %s", err)
		}

		cluster["kube_config_raw"] = string(kubeConfigRaw)

		cert, err := r.Extract()

		if err != nil {
			logp.Printf("Error retrieving HuaweiCloud CCE cluster cert: %s", err)
		}

		//Set Certificate Clusters
		var clusterList []map[string]interface{}
		for _, clusterObj := range cert.Clusters {
			clusterCert := make(map[string]interface{})
			clusterCert["name"] = clusterObj.Name
			clusterCert["server"] = clusterObj.Cluster.Server
			clusterCert["certificate_authority_data"] = clusterObj.Cluster.CertAuthorityData
			clusterList = append(clusterList, clusterCert)
		}
		cluster["certificate_clusters"] = clusterList

		//Set Certificate Users
		var userList []map[string]interface{}
		for _, userObj := range cert.Users {
			userCert := make(map[string]interface{})
			userCert["name"] = userObj.Name
			userCert["client_certificate_data"] = userObj.User.ClientCertData
			userCert["client_key_data"] = userObj.User.ClientKeyData
			userList = append(userList, userCert)
		}
		cluster["certificate_users"] = userList

		// Set masters
		var masterList []map[string]interface{}
		for _, masterObj := range v.Spec.Masters {
			master := make(map[string]interface{})
			master["availability_zone"] = masterObj.MasterAZ
			masterList = append(masterList, master)
		}
		cluster["masters"] = masterList

		clustersToSet = append(clustersToSet, cluster)
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("ids", ids),
		d.Set("clusters", clustersToSet),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting cce clusters fields: %s", err)
	}
	return nil
}
