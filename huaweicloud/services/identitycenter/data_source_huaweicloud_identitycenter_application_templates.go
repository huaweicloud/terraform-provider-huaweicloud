package identitycenter

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/application-templates
func DataSourceIdentityCenterApplicationTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterApplicationTemplatesRead,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sso_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_version": {
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
					},
				},
			},
		},
	}
}

func resourceIdentityCenterApplicationTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/application-templates"
		product = "identitycenter"
	)

	query := "?application_id=" + d.Get("application_id").(string)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	path := client.Endpoint + httpUrl + query

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET",
		path, &opt)

	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("application_templates", flattenApplicationTemplates(utils.PathSearch("application_templates",
			respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplicationTemplates(data interface{}) []interface{} {
	templates := data.([]interface{})
	if len(templates) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(templates))
	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"display_name":            utils.PathSearch("application.display.display_name", template, nil),
			"description":             utils.PathSearch("application.display.display_name", template, nil),
			"application_type":        utils.PathSearch("application.application_type", template, nil),
			"sso_protocol":            utils.PathSearch("sso_protocol", template, nil),
			"template_id":             utils.PathSearch("template_id", template, nil),
			"template_version":        utils.PathSearch("template_version", template, nil),
			"security_config":         flattenSecurityConfig(utils.PathSearch("security_config", template, nil)),
			"service_provider_config": flattenServiceProviderConfig(utils.PathSearch("service_provider_config", template, nil)),
			"response_config":         marshalJsonFormatParams("response config", utils.PathSearch("response_config", template, nil)),
			"response_schema_config":  marshalJsonFormatParams("response schema config", utils.PathSearch("response_schema_config", template, nil)),
		})
	}
	return result
}
