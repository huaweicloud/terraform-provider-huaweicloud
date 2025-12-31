package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/module-templates/restriction
func DataSourceCollectorModuleRestrictions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCollectorModuleRestrictionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"module_restrictions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"example": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"required": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"restrictions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"logic": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"title": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"template_field_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"title": {
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
					},
				},
			},
		},
	}
}

func buildCollectorModuleRestrictionsParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"template_ids": utils.ExpandToStringList(d.Get("template_ids").([]interface{})),
	}
}

func dataSourceCollectorModuleRestrictionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/module-templates/restriction"
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCollectorModuleRestrictionsParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving collector module restrictions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	respData, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("module_restrictions", flattenCollectorModuleRestrictions(respData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorModuleRestrictions(templateFields []interface{}) []interface{} {
	if len(templateFields) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(templateFields))
	for _, v := range templateFields {
		rst = append(rst, map[string]interface{}{
			"template_id": utils.PathSearch("template_id", v, nil),
			"fields":      flattenTemplateFileds(utils.PathSearch("fields", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenTemplateFileds(fieldsInfo []interface{}) []interface{} {
	if len(fieldsInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(fieldsInfo))
	for _, v := range fieldsInfo {
		rst = append(rst, map[string]interface{}{
			"default_value":     utils.PathSearch("default_value", v, nil),
			"description":       utils.PathSearch("description", v, nil),
			"example":           utils.PathSearch("example", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"required":          utils.PathSearch("required", v, nil),
			"restrictions":      flattenRestrictionsInfo(utils.PathSearch("restrictions", v, make([]interface{}, 0)).([]interface{})),
			"template_field_id": utils.PathSearch("template_field_id", v, nil),
			"title":             utils.PathSearch("title", v, nil),
			"type":              utils.PathSearch("type", v, nil),
		})
	}

	return rst
}

func flattenRestrictionsInfo(restrictionsInfo []interface{}) []interface{} {
	if len(restrictionsInfo) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(restrictionsInfo))
	for _, v := range restrictionsInfo {
		rst = append(rst, map[string]interface{}{
			"logic": utils.PathSearch("logic", v, nil),
			"title": utils.PathSearch("title", v, nil),
			"type":  utils.PathSearch("type", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
