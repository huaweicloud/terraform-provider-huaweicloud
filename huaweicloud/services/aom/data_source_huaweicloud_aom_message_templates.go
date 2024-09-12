package aom

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

// @API AOM GET /v2/{project_id}/events/notification/templates
func DataSourceMessageTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMessageTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"message_templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"templates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"topic": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
					},
				},
			},
		},
	}
}

func dataSourceMessageTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listMessageTemplates(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("message_templates", flattenMessageTemplates(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listMessageTemplates(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/events/notification/templates"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards folder: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards folder: %s", err)
	}

	return listRespBody, nil
}

func flattenMessageTemplates(templates []interface{}) []interface{} {
	if len(templates) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(templates))
	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"name":                  utils.PathSearch("name", template, nil),
			"source":                utils.PathSearch("source", template, nil),
			"templates":             flattenMessageTemplateDetailTemplates(utils.PathSearch("templates", template, nil)),
			"locale":                utils.PathSearch("locale", template, nil),
			"description":           utils.PathSearch("desc", template, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", template, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", template, float64(0)).(float64))/1000, true),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("modify_time", template, float64(0)).(float64))/1000, true),
		})
	}
	return result
}
