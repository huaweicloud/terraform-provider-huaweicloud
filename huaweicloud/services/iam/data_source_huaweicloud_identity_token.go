package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GET /v3/auth/tokens
func DataSourceIdentityToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityTokenRead,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expires_at": {

				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"methods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
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
					},
				},
			},
			"domain": {
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
					},
				},
			},
			"project": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
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
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"roles": {
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
					},
				},
			},
			"user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
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
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password_expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityTokenRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	product := "identity"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating identity client: %s", err)
	}

	token := d.Get("token").(string)
	if token == "" {
		return diag.Errorf("error getting identity token")
	}

	getTokenOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Subject-Token": token,
		},
	}

	resp, err := client.Request("GET", client.Endpoint+"v3/auth/tokens", &getTokenOpts)
	if err != nil {
		return diag.Errorf("failed to validate token: %s", err)
	}

	getTokenRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("expires_at", utils.PathSearch("token.expires_at", getTokenRespBody, nil)),
		d.Set("methods", utils.PathSearch("token.methods", getTokenRespBody, nil)),
		d.Set("catalog", flattenCatalog(getTokenRespBody)),
		d.Set("domain", flattenDomain(getTokenRespBody)),
		d.Set("project", flattenProject(getTokenRespBody)),
		d.Set("roles", flattenRoles(getTokenRespBody)),
		d.Set("issued_at", utils.PathSearch("token.issued_at", getTokenRespBody, nil)),
		d.Set("user", flattenUser(getTokenRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCatalog(getTokenRespBody interface{}) []map[string]interface{} {
	catalogRaw := utils.PathSearch("token.catalog", getTokenRespBody, nil)
	if catalogRaw == nil {
		return nil
	}
	catalog := catalogRaw.([]interface{})
	res := make([]map[string]interface{}, len(catalog))
	for i, v := range catalog {
		res[i] = map[string]interface{}{
			"endpoints": flattenEndpoints(v),
			"id":        utils.PathSearch("id", v, nil),
			"name":      utils.PathSearch("name", v, nil),
			"type":      utils.PathSearch("type", v, nil),
		}
	}
	return res
}

func flattenEndpoints(getEndpoints interface{}) []map[string]interface{} {
	endpointsRaw := utils.PathSearch("endpoints", getEndpoints, nil)
	if endpointsRaw == nil {
		return nil
	}
	endpoints := endpointsRaw.([]interface{})
	res := make([]map[string]interface{}, len(endpoints))
	for i, v := range endpoints {
		res[i] = map[string]interface{}{
			"id":        utils.PathSearch("id", v, nil),
			"interface": utils.PathSearch("interface", v, nil),
			"region":    utils.PathSearch("region", v, nil),
			"region_id": utils.PathSearch("region_id", v, nil),
			"url":       utils.PathSearch("url", v, nil),
		}
	}
	return res
}

func flattenDomain(getTokenRespBody interface{}) []map[string]interface{} {
	domainRaw := utils.PathSearch("token.domain", getTokenRespBody, nil)
	if domainRaw == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"id":   utils.PathSearch("id", domainRaw, nil),
			"name": utils.PathSearch("name", domainRaw, nil),
		},
	}
	return res
}

func flattenOtherDomain(getOtherDomain interface{}) []map[string]interface{} {
	otherDomainRaw := utils.PathSearch("domain", getOtherDomain, nil)
	if otherDomainRaw == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"id":   utils.PathSearch("id", otherDomainRaw, nil),
			"name": utils.PathSearch("name", otherDomainRaw, nil),
		},
	}
	return res
}

func flattenProject(getTokenRespBody interface{}) []map[string]interface{} {
	projectRaw := utils.PathSearch("token.project", getTokenRespBody, nil)
	if projectRaw == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain": flattenOtherDomain(projectRaw),
			"id":     utils.PathSearch("id", projectRaw, nil),
			"name":   utils.PathSearch("name", projectRaw, nil),
		},
	}
	return res
}

func flattenRoles(getTokenRespBody interface{}) []map[string]interface{} {
	rolesRaw := utils.PathSearch("token.roles", getTokenRespBody, nil)
	if rolesRaw == nil {
		return nil
	}
	roles := rolesRaw.([]interface{})
	res := make([]map[string]interface{}, len(roles))
	for i, v := range roles {
		res[i] = map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
		}
	}
	return res
}

func flattenUser(getTokenRespBody interface{}) []map[string]interface{} {
	userRaw := utils.PathSearch("token.user", getTokenRespBody, nil)
	if userRaw == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain":              flattenOtherDomain(userRaw),
			"id":                  utils.PathSearch("id", userRaw, nil),
			"name":                utils.PathSearch("name", userRaw, nil),
			"password_expires_at": utils.PathSearch("password_expires_at", userRaw, nil),
		},
	}
	return res
}
