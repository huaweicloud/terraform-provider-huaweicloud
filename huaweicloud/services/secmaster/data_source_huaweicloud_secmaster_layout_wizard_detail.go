package secmaster

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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards
func DataSourceLayoutWizardDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLayoutWizardDetailRead,

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
			"field_id": {
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

func dataSourceLayoutWizardDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards"
		layoutId = d.Get("layout_id").(string)
		fieldId  = d.Get("field_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", d.Get("workspace_id").(string))
	listPath = fmt.Sprintf("%s?layout_id=%s&field_id=%s", listPath, layoutId, fieldId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving layout wizard detail: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data", listRespBody, nil)

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenLayoutWizardDetail(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLayoutWizardDetail(layoutWizard interface{}) []interface{} {
	if layoutWizard == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":             utils.PathSearch("id", layoutWizard, nil),
		"name":           utils.PathSearch("name", layoutWizard, nil),
		"en_name":        utils.PathSearch("en_name", layoutWizard, nil),
		"description":    utils.PathSearch("description", layoutWizard, nil),
		"en_description": utils.PathSearch("en_description", layoutWizard, nil),
		"wizard_json":    utils.PathSearch("wizard_json", layoutWizard, nil),
		"creator_id":     utils.PathSearch("creator_id", layoutWizard, nil),
		"create_time":    utils.PathSearch("create_time", layoutWizard, nil),
		"update_time":    utils.PathSearch("update_time", layoutWizard, nil),
		"project_id":     utils.PathSearch("project_id", layoutWizard, nil),
		"workspace_id":   utils.PathSearch("workspace_id", layoutWizard, nil),
		"is_binding":     utils.PathSearch("is_binding", layoutWizard, nil),
		"binding_button": flattenWizardDetailBindingButtons(
			utils.PathSearch("binding_buttons", layoutWizard, make([]interface{}, 0)).([]interface{})),
		"is_built_in": utils.PathSearch("is_built_in", layoutWizard, nil),
		"boa_version": utils.PathSearch("boa_version", layoutWizard, nil),
		"version":     utils.PathSearch("version", layoutWizard, nil),
	}

	return []interface{}{result}
}

func flattenWizardDetailBindingButtons(buttons []interface{}) []interface{} {
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
