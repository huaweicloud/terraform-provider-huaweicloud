package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IDENTITYCENTER GET /v1/instances/{instance_id}/application-instances
func DataSourceIdentityCenterApplicationInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterApplicationInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_schema_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ttl": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"service_provider_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audience": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"require_request_signature": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"consumers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"binding": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"default_value": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"start_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"identity_provider_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issuer_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metadata_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_login_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_logout_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"active_certificate": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"algorithm": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"certificate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiry_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"issue_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"visible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_user_visible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"managed_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListApplicationInstancesParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func listApplicationInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/instances/{instance_id}/application-instances"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	queryParams := buildListApplicationInstancesParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center application instances: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		applications := utils.PathSearch("application_instances", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListApplicationInstancesParams(marker)
	}
	return result, nil
}

func dataSourceIdentityCenterApplicationInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	applicationInstances, err := listApplicationInstances(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("application_instances", flattenApplicationInstances(applicationInstances)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplicationInstances(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"name":                     utils.PathSearch("name", application, nil),
			"display_name":             utils.PathSearch("display.display_name", application, nil),
			"description":              utils.PathSearch("display.description", application, nil),
			"active_certificate":       flattenActiveCertificate(utils.PathSearch("active_certificate", application, nil)),
			"identity_provider_config": flattenIdentityProviderConfig(utils.PathSearch("identity_provider_config", application, nil)),
			"visible":                  utils.PathSearch("visible", application, nil),
			"status":                   utils.PathSearch("status", application, nil),
			"client_id":                utils.PathSearch("client_id", application, nil),
			"end_user_visible":         utils.PathSearch("end_user_visible", application, nil),
			"managed_account":          utils.PathSearch("managed_account", application, nil),
			"security_config":          flattenSecurityConfig(utils.PathSearch("security_config", application, nil)),
			"service_provider_config":  flattenServiceProviderConfig(utils.PathSearch("service_provider_config", application, nil)),
			"response_config":          marshalJsonFormatParams("response config", utils.PathSearch("response_config", application, nil)),
			"response_schema_config": marshalJsonFormatParams("response schema config",
				utils.PathSearch("response_schema_config", application, nil)),
		})
	}
	return result
}
