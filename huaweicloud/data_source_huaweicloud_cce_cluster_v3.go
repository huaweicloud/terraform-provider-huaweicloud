package huaweicloud

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"highway_subnet_id": {
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
				Optional: true,
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
	}
}

func dataSourceCCEClusterV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cceClient, err := config.CceV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud CCE client : %s", err)
	}

	listOpts := clusters.ListOpts{
		ID:    d.Get("id").(string),
		Name:  d.Get("name").(string),
		Type:  d.Get("cluster_type").(string),
		Phase: d.Get("status").(string),
		VpcID: d.Get("vpc_id").(string),
	}

	refinedClusters, err := clusters.List(cceClient, listOpts)
	logp.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Cluster := refinedClusters[0]

	logp.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", Cluster.Metadata.Id, Cluster)

	d.SetId(Cluster.Metadata.Id)
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", Cluster.Metadata.Name),
		d.Set("status", Cluster.Status.Phase),
		d.Set("flavor_id", Cluster.Spec.Flavor),
		d.Set("cluster_version", Cluster.Spec.Version),
		d.Set("cluster_type", Cluster.Spec.Type),
		d.Set("description", Cluster.Spec.Description),
		d.Set("billing_mode", Cluster.Spec.BillingMode),
		d.Set("vpc_id", Cluster.Spec.HostNetwork.VpcId),
		d.Set("subnet_id", Cluster.Spec.HostNetwork.SubnetId),
		d.Set("highway_subnet_id", Cluster.Spec.HostNetwork.HighwaySubnet),
		d.Set("container_network_cidr", Cluster.Spec.ContainerNetwork.Cidr),
		d.Set("container_network_type", Cluster.Spec.ContainerNetwork.Mode),
		d.Set("eni_subnet_id", Cluster.Spec.EniNetwork.SubnetId),
		d.Set("eni_subnet_cidr", Cluster.Spec.EniNetwork.Cidr),
		d.Set("authentication_mode", Cluster.Spec.Authentication.Mode),
		d.Set("security_group_id", Cluster.Spec.HostNetwork.SecurityGroup),
		d.Set("enterprise_project_id", Cluster.Spec.ExtendParam["enterpriseProjectId"]),
		d.Set("service_network_cidr", Cluster.Spec.KubernetesSvcIPRange),
	)

	// set endpoints
	var v []map[string]interface{}
	for _, endpoint := range Cluster.Status.Endpoints {

		mapping := map[string]interface{}{
			"url":  endpoint.Url,
			"type": endpoint.Type,
		}
		v = append(v, mapping)
	}
	mErr = multierror.Append(mErr, d.Set("endpoints", v))

	// set kube_config_raw
	r := clusters.GetCert(cceClient, d.Id())
	kubeConfigRaw, err := utils.JsonMarshal(r.Body)
	if err != nil {
		logp.Printf("Error marshaling r.Body: %s", err)
	}
	mErr = multierror.Append(mErr, d.Set("kube_config_raw", string(kubeConfigRaw)))

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
	mErr = multierror.Append(mErr, d.Set("certificate_clusters", clusterList))

	//Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := make(map[string]interface{})
		userCert["name"] = userObj.Name
		userCert["client_certificate_data"] = userObj.User.ClientCertData
		userCert["client_key_data"] = userObj.User.ClientKeyData
		userList = append(userList, userCert)
	}
	mErr = multierror.Append(mErr, d.Set("certificate_users", userList))

	// Set masters
	var masterList []map[string]interface{}
	for _, masterObj := range Cluster.Spec.Masters {
		master := make(map[string]interface{})
		master["availability_zone"] = masterObj.MasterAZ
		masterList = append(masterList, master)
	}
	mErr = multierror.Append(mErr, d.Set("masters", masterList))

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting cluster fields: %s", err)
	}

	return nil
}
