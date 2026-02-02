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
func DataSourceV5AuthorizationSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5AuthorizationSchemaRead,

		Schema: map[string]*schema.Schema{
			"service_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The service name abbreviation to be queried.`,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authorization item name.`,
						},
						"permission_only": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: `Whether the authorization item is only used as a permission point and
does not correspond to any operation.`,
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
										Description: `The service custom conditional attribute list and some global attribute list supported
  by the authorization item and resource.`,
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the resource type is mandatory for this authorization item.`,
									},
									"urn_template": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The uniform resource name template for the resource.`,
									},
								},
							},
						},
						"access_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The access level granted when using this authorization item in a policy.`,
						},
						"aliases": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of authorization item aliases.`,
						},
						"condition_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: `The service custom conditional attribute list and some global attribute list supported by the
  authorization items and are independent of the resources.`,
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
				Description: `The list of authorization items supported by the cloud service.`,
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
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The English description.`,
									},
									"zh_cn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The Chinese description.`,
									},
								},
							},
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The condition key name.`,
						},
						"multi_valued": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the condition value is multi-valued.`,
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type of the condition value.`,
						},
					},
				},
				Description: `The list of condition keys supported by the cloud service.`,
			},
			"operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dependent_actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The other authorization item list that this operation may require.`,
						},
						"operation_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The action of the operation.`,
						},
						"operation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OpenAPI operation identifier.`,
						},
					},
				},
				Description: `The list of operations supported by the cloud service.`,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"urn_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name template for the resource.`,
						},
						"type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type name of the resource.`,
						},
					},
				},
				Description: `The list of resources supported by the cloud service.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version number of the service authorization summary.`,
			},
		},
	}
}

func dataSourceV5AuthorizationSchemaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v5/authorization-schemas/services/{service_code}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{service_code}", d.Get("service_code").(string))

	createLoginProfileOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, createLoginProfileOpt)
	if err != nil {
		return diag.Errorf("error getting IAM authorization schema: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("actions", flattenV5AuthorizationSchemaActions(utils.PathSearch("actions", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("conditions", flattenV5AuthorizationSchemaConditions(utils.PathSearch("conditions", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("operations", flattenV5AuthorizationSchemaOperations(utils.PathSearch("operations", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("resources", flattenV5AuthorizationSchemaResources(utils.PathSearch("resources", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("version", utils.PathSearch("version", respBody, "").(string)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM authorization schema fields: %s", err)
	}

	return nil
}

func flattenV5AuthorizationSchemaActions(actions []interface{}) []interface{} {
	if len(actions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(actions))
	for _, v := range actions {
		result = append(result, map[string]interface{}{
			"name":            utils.PathSearch("name", v, nil),
			"permission_only": utils.PathSearch("permission_only", v, nil),
			"resources": flattenV5AuthorizationSchemaActionResources(utils.PathSearch("resources", v,
				make([]interface{}, 0)).([]interface{})),
			"access_level":   utils.PathSearch("access_level", v, nil),
			"aliases":        utils.PathSearch("aliases", v, nil),
			"condition_keys": utils.PathSearch("condition_keys", v, nil),
			"description":    flattenV5AuthorizationSchemaDescription(utils.PathSearch("description", v, nil)),
		})
	}

	return result
}

func flattenV5AuthorizationSchemaActionResources(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"condition_keys": utils.PathSearch("condition_keys", v, nil),
			"required":       utils.PathSearch("required", v, nil),
			"urn_template":   utils.PathSearch("urn_template", v, nil),
		})
	}

	return result
}

func flattenV5AuthorizationSchemaDescription(description interface{}) []interface{} {
	if description == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"en_us": utils.PathSearch("en_US", description, nil),
			"zh_cn": utils.PathSearch("zh_CN", description, nil),
		},
	}
}

func flattenV5AuthorizationSchemaConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(conditions))
	for _, v := range conditions {
		result = append(result, map[string]interface{}{
			"description":  flattenV5AuthorizationSchemaDescription(utils.PathSearch("description", v, nil)),
			"key":          utils.PathSearch("key", v, nil),
			"multi_valued": utils.PathSearch("multi_valued", v, nil),
			"value_type":   utils.PathSearch("value_type", v, nil),
		})
	}

	return result
}

func flattenV5AuthorizationSchemaOperations(operations []interface{}) []interface{} {
	if len(operations) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(operations))
	for _, v := range operations {
		result = append(result, map[string]interface{}{
			"dependent_actions": utils.PathSearch("dependent_actions", v, nil),
			"operation_action":  utils.PathSearch("operation_action", v, nil),
			"operation_id":      utils.PathSearch("operation_id", v, nil),
		})
	}

	return result
}

func flattenV5AuthorizationSchemaResources(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"urn_template": utils.PathSearch("urn_template", v, nil),
			"type_name":    utils.PathSearch("type_name", v, nil),
		})
	}

	return result
}
