package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GET /v3/auth/tokens
func DataSourceIdentityUserTokenInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityUserTokenInfoRead,

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_catalog": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
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
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

func dataSourceIdentityUserTokenInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		getIdentityTokenHttpUrl = "v3/auth/tokens"
		product                 = "iam"
	)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Identity Token client: %s", err)
	}
	getIdentityTokenBasePath := client.Endpoint + getIdentityTokenHttpUrl
	getIdentityTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Subject-Token": d.Get("token").(string),
		},
	}
	noCatalog := d.Get("no_catalog").(string)
	if noCatalog == "true" {
		getIdentityTokenBasePath = getIdentityTokenBasePath + "?nocatalog=" + noCatalog
	}
	response, err := client.Request("GET", getIdentityTokenBasePath, &getIdentityTokenOpt)
	if err != nil {
		return diag.Errorf("error get identity token: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	tokenModel := utils.PathSearch("token", respBody, nil)
	return flattenResponse(tokenModel, noCatalog, d)
}

func flattenResponse(tokenModel interface{}, noCatalog string, d *schema.ResourceData) diag.Diagnostics {
	if tokenModel == nil {
		return nil
	}
	mErr := multierror.Append(nil,
		d.Set("expires_at", utils.PathSearch("expires_at", tokenModel, nil)),
		d.Set("issued_at", utils.PathSearch("issued_at", tokenModel, nil)),
		d.Set("methods", utils.PathSearch("methods", tokenModel, nil)),
		d.Set("domain", flattenDomain(utils.PathSearch("domain", tokenModel, nil))),
		d.Set("project", flattenProject(utils.PathSearch("project", tokenModel, nil))),
		d.Set("roles", flattenRoles(utils.PathSearch("roles", tokenModel, nil))),
		d.Set("user", flattenUser(utils.PathSearch("user", tokenModel, nil))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	if noCatalog != "true" {
		if e := d.Set("catalog", flattenCatalog(utils.PathSearch("catalog", tokenModel, nil))); e != nil {
			return diag.Errorf("error setting token fields catalog: %s", e)
		}
	}
	return nil
}

func flattenCatalog(catalogs interface{}) []map[string]interface{} {
	catalogsArr := catalogs.([]interface{})
	res := make([]map[string]interface{}, len(catalogsArr))
	for i, catalog := range catalogsArr {
		res[i] = map[string]interface{}{
			"endpoints": flattenEndpoints(utils.PathSearch("endpoints", catalog, nil)),
			"id":        utils.PathSearch("id", catalog, nil).(string),
			"name":      utils.PathSearch("name", catalog, nil).(string),
			"type":      utils.PathSearch("type", catalog, nil).(string),
		}
	}
	return res
}

func flattenEndpoints(endpoints interface{}) []map[string]interface{} {
	endpointsArr := endpoints.([]interface{})
	res := make([]map[string]interface{}, len(endpointsArr))
	for i, endpoint := range endpointsArr {
		res[i] = map[string]interface{}{
			"id":        utils.PathSearch("id", endpoint, nil).(string),
			"interface": utils.PathSearch("interface", endpoint, nil).(string),
			"region":    utils.PathSearch("region", endpoint, nil).(string),
			"region_id": utils.PathSearch("region_id", endpoint, nil).(string),
			"url":       utils.PathSearch("url", endpoint, nil).(string),
		}
	}
	return res
}

func flattenDomain(domain interface{}) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   utils.PathSearch("id", domain, nil).(string),
		"name": utils.PathSearch("name", domain, nil).(string),
	}
	return res
}

func flattenProject(projects interface{}) []map[string]interface{} {
	if projects == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain": flattenProjectDomain(utils.PathSearch("domain", projects, nil)),
			"id":     utils.PathSearch("id", projects, nil).(string),
			"name":   utils.PathSearch("name", projects, nil).(string),
		},
	}
	return res
}

func flattenProjectDomain(domain interface{}) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   utils.PathSearch("id", domain, nil).(string),
		"name": utils.PathSearch("name", domain, nil).(string),
	}
	return res
}

func flattenRoles(roles interface{}) []map[string]interface{} {
	rolesArr := roles.([]interface{})
	res := make([]map[string]interface{}, len(rolesArr))
	for i, role := range rolesArr {
		res[i] = map[string]interface{}{
			"id":   utils.PathSearch("id", role, nil).(string),
			"name": utils.PathSearch("name", role, nil).(string),
		}
	}
	return res
}

func flattenUser(user interface{}) []map[string]interface{} {
	if user == nil {
		return nil
	}
	res := []map[string]interface{}{
		{
			"domain":              flattenUserDomain(utils.PathSearch("domain", user, nil)),
			"id":                  utils.PathSearch("id", user, nil).(string),
			"name":                utils.PathSearch("name", user, nil).(string),
			"password_expires_at": utils.PathSearch("password_expires_at", user, nil).(string),
		},
	}
	return res
}

func flattenUserDomain(domain interface{}) map[string]interface{} {
	if domain == nil {
		return nil
	}
	res := map[string]interface{}{
		"id":   utils.PathSearch("id", domain, nil).(string),
		"name": utils.PathSearch("name", domain, nil).(string),
	}
	return res
}
