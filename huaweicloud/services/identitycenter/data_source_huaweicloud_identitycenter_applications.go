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

// @API IDENTITYCENTER GET /v1/instances/{instance_id}/applications
func DataSourceIdentityCenterApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterApplicationsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_provider_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignment_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"assignment_required": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_urn": {
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
						"application_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portal_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"visible": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"visibility": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sign_in_options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"origin": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"application_url": {
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

func buildListApplicationsParams(marker string) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func listApplications(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/instances/{instance_id}/applications"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	queryParams := buildListApplicationsParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center applications: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		applications := utils.PathSearch("applications", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, applications...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListApplicationsParams(marker)
	}
	return result, nil
}

func dataSourceIdentityCenterApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	applications, err := listApplications(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("applications", flattenApplications(applications)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplications(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"application_urn":          utils.PathSearch("application_urn", application, nil),
			"application_provider_urn": utils.PathSearch("application_provider_urn", application, nil),
			"created_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("created_date", application, float64(0)).(float64))/1000, false),
			"description":         utils.PathSearch("description", application, nil),
			"instance_urn":        utils.PathSearch("instance_urn", application, nil),
			"name":                utils.PathSearch("name", application, nil),
			"status":              utils.PathSearch("status", application, nil),
			"application_account": utils.PathSearch("application_account", application, nil),
			"assignment_config":   flattenAssignmentConfig(utils.PathSearch("assignment_config", application, nil)),
			"portal_options":      flattenPortalOptions(utils.PathSearch("portal_options", application, nil)),
		})
	}
	return result
}

func flattenAssignmentConfig(cfg interface{}) []map[string]interface{} {
	if cfg == nil || len(cfg.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"assignment_required": utils.PathSearch("assignment_required", cfg, nil),
		},
	}
}

func flattenPortalOptions(options interface{}) []map[string]interface{} {
	if options == nil || len(options.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"visible":         utils.PathSearch("visible", options, nil),
			"visibility":      utils.PathSearch("visibility", options, nil),
			"sign_in_options": flattenSignInOptions(utils.PathSearch("sign_in_options", options, nil)),
		},
	}
}

func flattenSignInOptions(options interface{}) []map[string]interface{} {
	if options == nil || len(options.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"origin":          utils.PathSearch("origin", options, nil),
			"application_url": utils.PathSearch("application_url", options, nil),
		},
	}
}
