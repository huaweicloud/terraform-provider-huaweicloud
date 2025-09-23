package cce

import (
	"context"
	"net/url"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /openid/v1/jwks
func DataSourceCCEClusterOpenIDJWKS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterOpenIDJWKSRead,

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
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kty": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"n": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"e": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEClusterOpenIDJWKSRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)

	var mErr *multierror.Error

	var (
		getOpenIDJWKSHttpUrl = "openid/v1/jwks"
		getOpenIDJWKSProduct = "cce"
	)
	getOpenIDJWKSClient, err := cfg.NewServiceClient(getOpenIDJWKSProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	u, _ := url.Parse(getOpenIDJWKSClient.Endpoint)
	u.Host = clusterID + "." + u.Host
	rbUrl := u.String()

	getOpenIDJWKSPath := rbUrl + getOpenIDJWKSHttpUrl

	getOpenIDJWKSOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Accept": "application/jwk-set+json",
		},
	}

	getOpenIDJWKSResp, err := getOpenIDJWKSClient.Request("GET", getOpenIDJWKSPath, &getOpenIDJWKSOpt)

	if err != nil {
		return diag.Errorf("error retrieving CCE cluster public keys: %s", err)
	}

	getOpenIDJWKSRespBody, err := utils.FlattenResponse(getOpenIDJWKSResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("keys", flattenClusterOpenIDJWKS(getOpenIDJWKSRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterOpenIDJWKS(resp interface{}) []map[string]interface{} {
	keys := utils.PathSearch("keys", resp, []interface{}{}).([]interface{})
	if len(keys) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, len(keys))
	for i, key := range keys {
		res[i] = map[string]interface{}{
			"use": utils.PathSearch("use", key, nil),
			"kty": utils.PathSearch("kty", key, nil),
			"kid": utils.PathSearch("kid", key, nil),
			"alg": utils.PathSearch("alg", key, nil),
			"n":   utils.PathSearch("n", key, nil),
			"e":   utils.PathSearch("e", key, nil),
		}
	}
	return res
}
