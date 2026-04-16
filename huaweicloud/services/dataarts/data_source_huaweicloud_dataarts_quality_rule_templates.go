package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/quality/rule-templates
func DataSourceQualityRuleTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: qualityRuleTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the quality rule templates are located.`,
			},

			// Required parameter
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the quality rule templates belong.`,
			},

			// Optional parameters
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the quality rule template.`,
			},
			"category_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The category ID of the quality rule template.`,
			},
			"system_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query only system templates.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the quality rule template.`,
			},

			// Attributes
			"templates": {
				Type:        schema.TypeList,
				Elem:        qualityRuleTemplateSchema(),
				Computed:    true,
				Description: `The list of quality rule templates that matched filter parameters.`,
			},
		},
	}
}

func qualityRuleTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ID of the quality rule template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the quality rule template.`,
			},
			"category_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The category ID of the quality rule template.`,
			},
			"dimension": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dimension of the quality rule template.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the quality rule template.`,
			},
			"system_template": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the quality rule template is a system template.`,
			},
			"sql_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The definition relationship of the quality rule template.`,
			},
			"abnormal_table_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The abnormal table template of the quality rule template.`,
			},
			"result_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The result description of the quality rule template.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the quality rule template, in RFC3339 format.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the quality rule template.`,
			},
		},
	}
	return &sc
}

func buildQualityRuleTemplatesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("category_id"); ok {
		res = fmt.Sprintf("%s&category_id=%v", res, v)
	}
	if v, ok := d.GetOk("system_template"); ok {
		res = fmt.Sprintf("%s&system_template=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	return res
}

func queryQualityRuleTemplates(client *golangsdk.ServiceClient, workspaceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/quality/rule-templates?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		templates := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, templates...)
		if len(templates) < limit {
			break
		}

		offset += len(templates)
	}

	return result, nil
}

func flattenQualityRuleTemplates(templates []interface{}) []interface{} {
	result := make([]interface{}, 0, len(templates))

	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"id":                      utils.PathSearch("id", template, nil),
			"name":                    utils.PathSearch("name", template, nil),
			"category_id":             utils.PathSearch("category_id", template, nil),
			"dimension":               utils.PathSearch("dimension", template, nil),
			"type":                    utils.PathSearch("type", template, nil),
			"system_template":         utils.PathSearch("system_template", template, nil),
			"sql_info":                utils.PathSearch("sql_info", template, nil),
			"abnormal_table_template": utils.PathSearch("abnormal_table_template", template, nil),
			"result_description":      utils.PathSearch("result_description", template, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", template, float64(0)).(float64))/1000, false),
			"creator": utils.PathSearch("creator", template, nil),
		})
	}

	return result
}

func qualityRuleTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	templates, err := queryQualityRuleTemplates(client, d.Get("workspace_id").(string),
		buildQualityRuleTemplatesQueryParams(d))
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("templates", flattenQualityRuleTemplates(templates)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
