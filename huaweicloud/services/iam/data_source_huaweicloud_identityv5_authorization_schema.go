package iam

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

// @API IAM GET /v5/authorization-schemas/services/{service_code}
func DataSourceIdentityV5AuthorizationSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5AuthorizationSchemaRead,

		Schema: map[string]*schema.Schema{
			"service_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_keys": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"required": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"urn_template": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"condition_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en_us": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zh_cn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en_us": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zh_cn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_valued": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"value_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dependent_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"operation_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"urn_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityV5AuthorizationSchemaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	serviceCode := d.Get("service_code").(string)

	createLoginProfileHttpUrl := "v5/authorization-schemas/services/{service_code}"
	createLoginProfilePath := iamClient.Endpoint + createLoginProfileHttpUrl
	createLoginProfilePath = strings.ReplaceAll(createLoginProfilePath, "{service_code}", serviceCode)

	createLoginProfileOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createLoginProfileResp, err := iamClient.Request("GET", createLoginProfilePath, createLoginProfileOpt)
	if err != nil {
		return diag.Errorf("error get IAM authorization schema: %s", err)
	}

	createLoginProfileBody, err := utils.FlattenResponse(createLoginProfileResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate uuid: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("actions", flattenActions(utils.PathSearch("actions", createLoginProfileBody, nil))),
		d.Set("conditions", flattenConditions(utils.PathSearch("conditions", createLoginProfileBody, nil))),
		d.Set("operations", flattenOperations(utils.PathSearch("operations", createLoginProfileBody, nil))),
		d.Set("resources", flattenResources(utils.PathSearch("resources", createLoginProfileBody, nil))),
		d.Set("version", utils.PathSearch("version", createLoginProfileBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM authorization schema fields: %s", err)
	}

	return nil
}

func flattenActions(actions interface{}) []interface{} {
	if actions == nil {
		return nil
	}

	actionList := actions.([]interface{})
	result := make([]interface{}, len(actionList))

	for i, v := range actionList {
		action := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"name":            action["name"],
			"permission_only": action["permission_only"],
			"resources":       flattenActionResources(action["resources"]),
			"access_level":    action["access_level"],
			"aliases":         action["aliases"],
			"condition_keys":  action["condition_keys"],
			"description":     flattenActionDescription(action["description"]),
		}
	}

	return result
}

func flattenActionResources(resources interface{}) []interface{} {
	if resources == nil {
		return nil
	}

	resourceList := resources.([]interface{})
	result := make([]interface{}, len(resourceList))

	for i, v := range resourceList {
		resource := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"condition_keys": resource["condition_keys"],
			"required":       resource["required"],
			"urn_template":   resource["urn_template"],
		}
	}

	return result
}

func flattenActionDescription(description interface{}) []interface{} {
	if description == nil {
		return nil
	}

	desc := description.(map[string]interface{})
	result := make([]interface{}, 1)
	result[0] = map[string]interface{}{
		"en_us": desc["en_US"],
		"zh_cn": desc["zh_CN"],
	}

	return result
}

func flattenConditions(conditions interface{}) []interface{} {
	if conditions == nil {
		return nil
	}

	conditionList := conditions.([]interface{})
	result := make([]interface{}, len(conditionList))

	for i, v := range conditionList {
		condition := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"description":  flattenConditionDescription(condition["description"]),
			"key":          condition["key"],
			"multi_valued": condition["multi_valued"],
			"value_type":   condition["value_type"],
		}
	}

	return result
}

func flattenConditionDescription(description interface{}) []interface{} {
	if description == nil {
		return nil
	}

	desc := description.(map[string]interface{})
	result := make([]interface{}, 1)
	result[0] = map[string]interface{}{
		"en_us": desc["en_US"],
		"zh_cn": desc["zh_CN"],
	}

	return result
}

func flattenOperations(operations interface{}) []interface{} {
	if operations == nil {
		return nil
	}

	operationList := operations.([]interface{})
	result := make([]interface{}, len(operationList))

	for i, v := range operationList {
		operation := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"dependent_actions": operation["dependent_actions"],
			"operation_action":  operation["operation_action"],
			"operation_id":      operation["operation_id"],
		}
	}

	return result
}

func flattenResources(resources interface{}) []interface{} {
	if resources == nil {
		return nil
	}

	resourceList := resources.([]interface{})
	result := make([]interface{}, len(resourceList))

	for i, v := range resourceList {
		resource := v.(map[string]interface{})
		result[i] = map[string]interface{}{
			"urn_template": resource["urn_template"],
			"type_name":    resource["type_name"],
		}
	}

	return result
}
