package cce

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /api/v3/projects/{project_id}/clusters/{id}/clustercert
func DataSourceCCEClusterCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterCertificateRead,

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
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"kube_config_raw": {
				Type:     schema.TypeString,
				Computed: true,
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
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_authority_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"insecure_skip_tls_verify": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"users": {
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
			"contexts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"current_context": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCCEClusterCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create CCE client : %s", err)
	}

	clusterID := d.Get("cluster_id").(string)

	opts := clusters.GetCertOpts{
		Duration: d.Get("duration").(int),
	}
	r := clusters.GetCert(cceClient, clusterID, opts)
	kubeConfigRaw, err := utils.JsonMarshal(r.Body)

	if err != nil {
		log.Printf("error marshaling r.Body: %s", err)
	}

	mErr := multierror.Append(nil, d.Set("kube_config_raw", string(kubeConfigRaw)))

	cert, err := r.Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve CCE cluster cert: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(mErr,
		d.Set("current_context", cert.CurrentContext),
		d.Set("clusters", flattenClusterCertClusters(cert)),
		d.Set("users", flattenClusterCertUsers(cert)),
		d.Set("contexts", flattenClusterCertContexts(cert)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting CCE clusters certificate: %s", err)
	}

	return nil
}

func flattenClusterCertClusters(cert *clusters.Certificate) []map[string]interface{} {
	certClusters := cert.Clusters
	res := make([]map[string]interface{}, len(certClusters))
	for i, cluster := range certClusters {
		res[i] = map[string]interface{}{
			"name":                       cluster.Name,
			"server":                     cluster.Cluster.Server,
			"certificate_authority_data": cluster.Cluster.CertAuthorityData,
			"insecure_skip_tls_verify":   cluster.Cluster.InsecureSkipTLSVerify,
		}
	}
	return res
}

func flattenClusterCertUsers(cert *clusters.Certificate) []map[string]interface{} {
	users := cert.Users
	res := make([]map[string]interface{}, len(users))
	for i, user := range users {
		res[i] = map[string]interface{}{
			"name":                    user.Name,
			"client_certificate_data": user.User.ClientCertData,
			"client_key_data":         user.User.ClientKeyData,
		}
	}
	return res
}

func flattenClusterCertContexts(cert *clusters.Certificate) []map[string]interface{} {
	contexts := cert.Contexts
	res := make([]map[string]interface{}, len(contexts))
	for i, c := range contexts {
		res[i] = map[string]interface{}{
			"name":    c.Name,
			"cluster": c.Context.Cluster,
			"user":    c.Context.User,
		}
	}
	return res
}
