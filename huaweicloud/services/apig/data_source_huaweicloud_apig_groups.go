package apig

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/api-groups
func DataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
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
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sl_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"on_sell_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url_domains": {
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
									"cname_status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ssl_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"min_ssl_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"verified_client_certificate_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"is_has_trusted_root_ca": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"sl_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"environment": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"variable": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The variable name.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The variable value.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the variable that the group has.",
												},
											},
										},
										Description: "The array of one or more environment variables.",
									},
									"environment_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the environment to which the variables belongs.",
									},
								},
							},
							Description: "The array of one or more environments of the associated group.",
						},
					},
				},
			},
		},
	}
}

func flattenGroupsUrlDomain(urlDomains []apigroups.UrlDomian) []map[string]interface{} {
	if len(urlDomains) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(urlDomains))
	for i, v := range urlDomains {
		result[i] = map[string]interface{}{
			"id":                                  v.Id,
			"name":                                v.DomainName,
			"cname_status":                        v.ResolutionStatus,
			"ssl_id":                              v.SSLId,
			"ssl_name":                            v.SSLName,
			"min_ssl_version":                     v.MinSSLVersion,
			"verified_client_certificate_enabled": v.VerifiedClientCertificateEnabled,
			"is_has_trusted_root_ca":              v.IsHasTrustedRootCA,
		}
	}

	return result
}

func flattenGroups(group apigroups.Group, variables []map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":             group.Id,
		"name":           group.Name,
		"status":         group.Status,
		"sl_domain":      group.Subdomain,
		"created_at":     group.RegistraionTime,
		"updated_at":     group.UpdateTime,
		"on_sell_status": group.OnSellStatus,
		"url_domains":    flattenGroupsUrlDomain(group.UrlDomians),
		"sl_domains":     group.Subdomains,
		"description":    group.Description,
		"is_default":     group.IsDefault,
		"environment":    variables,
	}

	return result
}

func dataSourceGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt := apigroups.ListOpts{
		Id:   d.Get("group_id").(string),
		Name: d.Get("name").(string),
	}
	pages, err := apigroups.List(client, d.Get("instance_id").(string), opt).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := apigroups.ExtractGroups(pages)
	if err != nil {
		return diag.Errorf("unable to get the group list form server: %v", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := make([]map[string]interface{}, 0)
	for _, group := range resp {
		variables, err := queryEnvironmentVariables(client, d.Get("instance_id").(string), group.Id)
		if err != nil {
			return diag.FromErr(err)
		}
		result = append(result, flattenGroups(group, flattenEnvironmentVariables(variables)))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", result),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving dedicated groups data source fields: %s", mErr)
	}
	return nil
}
