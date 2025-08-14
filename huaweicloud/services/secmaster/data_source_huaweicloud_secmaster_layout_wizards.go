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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}/wizards
func DataSourceLayoutWizards() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLayoutWizardsRead,

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
			"layout_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
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
						"en_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"en_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wizard_json": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_binding": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"binding_button": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"button_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"button_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"is_built_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"boa_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLayoutWizardsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}/wizards"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", d.Get("workspace_id").(string))
	listPath = strings.ReplaceAll(listPath, "{layout_id}", d.Get("layout_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving layout wizards: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data", listRespBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenLayoutWizards(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLayoutWizards(wizards []interface{}) []interface{} {
	if len(wizards) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(wizards))
	for _, v := range wizards {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"en_name":        utils.PathSearch("en_name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"en_description": utils.PathSearch("en_description", v, nil),
			"wizard_json":    utils.PathSearch("wizard_json", v, nil),
			"creator_id":     utils.PathSearch("creator_id", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"workspace_id":   utils.PathSearch("workspace_id", v, nil),
			"is_binding":     utils.PathSearch("is_binding", v, nil),
			"binding_button": flattenBindingButtons(utils.PathSearch("binding_buttons", v, make([]interface{}, 0)).([]interface{})),
			"is_built_in":    utils.PathSearch("is_built_in", v, nil),
			"boa_version":    utils.PathSearch("boa_version", v, nil),
			"version":        utils.PathSearch("version", v, nil),
		})
	}

	return result
}

func flattenBindingButtons(buttons []interface{}) []interface{} {
	if len(buttons) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(buttons))
	for _, v := range buttons {
		result = append(result, map[string]interface{}{
			"button_id":   utils.PathSearch("button_id", v, nil),
			"button_name": utils.PathSearch("button_name", v, nil),
		})
	}

	return result
}
