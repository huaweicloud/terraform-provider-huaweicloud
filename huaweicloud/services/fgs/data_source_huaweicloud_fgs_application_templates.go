package fgs

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

// @API FunctionGraph GET /v2/{project_id}/fgs/application/templates
func DataSourceApplicationTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application templates are located.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The runtime name used to query the application templates.`,
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The category used to query the application templates.`,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The template ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The template name.`,
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The template runtime.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The template category.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of template.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the function application.`,
						},
					},
				},
				Description: `All application templates that match the filter parameters.`,
			},
		},
	}
}

func buildApplicationTemplatesQueryParams(d *schema.ResourceData) string {
	if runtime, ok := d.GetOk("runtime"); ok {
		return fmt.Sprintf("&runtime=%v", runtime)
	}
	return ""
}

func getApplicationTemplates(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/application/templates?maxitems=100"
		marker  float64
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildApplicationTemplatesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%v", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		templates := utils.PathSearch("templates", respBody, make([]interface{}, 0)).([]interface{})
		if len(templates) < 1 {
			break
		}
		result = append(result, templates...)
		// In this API, marker has the same meaning as offset.
		nextMarker := utils.PathSearch("next_marker", respBody, float64(0)).(float64)
		if nextMarker == marker || nextMarker == 0 {
			// Make sure the next marker value is correct, not the previous marker or zero (in the last page).
			break
		}
		marker = nextMarker
	}

	return result, nil
}

func filterListTemplates(d *schema.ResourceData, templates []interface{}) []interface{} {
	result := templates

	if category, ok := d.GetOk("category"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?category=='%v']", category), result, make([]interface{}, 0)).([]interface{})
	}

	return result
}

func flattenListTemplates(templates []interface{}) []map[string]interface{} {
	if len(templates) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(templates))
	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", template, nil),
			"name":        utils.PathSearch("name", template, nil),
			"runtime":     utils.PathSearch("runtime", template, nil),
			"category":    utils.PathSearch("category", template, nil),
			"type":        utils.PathSearch("type", template, nil),
			"description": utils.PathSearch("description", template, nil),
		})
	}
	return result
}

func dataSourceApplicationTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	templateList, err := getApplicationTemplates(client, d)
	if err != nil {
		return diag.Errorf("error retrieving application templates: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("templates", flattenListTemplates(filterListTemplates(d, templateList))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
