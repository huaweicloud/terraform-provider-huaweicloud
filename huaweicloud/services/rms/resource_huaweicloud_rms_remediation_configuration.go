package rms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var remediationConfigurationNonUpdatableParams = []string{"policy_assignment_id"}

// @API Config PUT /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration
// @API Config GET /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration
func ResourceRemediationConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRemediationConfigurationCreateOrUpdate,
		ReadContext:   resourceRemediationConfigurationRead,
		UpdateContext: resourceRemediationConfigurationCreateOrUpdate,
		DeleteContext: resourceRemediationConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(remediationConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"policy_assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The policy assignment ID.`,
			},
			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The execution method of remediation.`,
			},
			"target_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of a remediation object.`,
			},
			"resource_parameter": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `The dynamic parameter of remediation.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The parameter name for passing the resource ID.`,
						},
					},
				},
			},
			"automatic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether remediation is automatic.`,
			},
			"static_parameter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The static parameters for the remediation execution.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"var_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The static parameter name.`,
						},
						"var_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The static parameter value.`,
						},
					},
				},
			},
			"auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The authorization type for remediation configurations.`,
			},
			"auth_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The information of dependent service authorization.`,
			},
			"maximum_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of retries allowed within a specified period.`,
			},
			"retry_attempt_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The time period during which the number of attempts specified in the maximum_attempts can be tried.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the remediation configuration was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the remediation configuration was updated.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who created the remediation configuration.`,
			},
		},
	}
}
func resourceRemediationConfigurationCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration"
		product = "rms"
	)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	policyAssignmentID := d.Get("policy_assignment_id").(string)

	createRemediationConfigurationPath := client.Endpoint + httpUrl
	createRemediationConfigurationPath = strings.ReplaceAll(createRemediationConfigurationPath, "{domain_id}", cfg.DomainID)
	createRemediationConfigurationPath = strings.ReplaceAll(createRemediationConfigurationPath, "{policy_assignment_id}", policyAssignmentID)

	createRemediationConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams, err := buildRemediationConfigurationBodyParams(d)
	if err != nil {
		return diag.Errorf("error building RMS remediation configuration request parameters: %s", err)
	}
	createRemediationConfigurationOpt.JSONBody = utils.RemoveNil(bodyParams)
	_, err = client.Request("PUT", createRemediationConfigurationPath, &createRemediationConfigurationOpt)
	if err != nil {
		return diag.Errorf("error creating RMS remediation configuration: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(policyAssignmentID)
	}
	return resourceRemediationConfigurationRead(ctx, d, meta)
}

func buildRemediationConfigurationBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"target_type":           d.Get("target_type"),
		"target_id":             d.Get("target_id"),
		"resource_parameter":    buildResourceParameterBodyParams(d.Get("resource_parameter")),
		"automatic":             utils.ValueIgnoreEmpty(d.Get("automatic")),
		"auth_type":             utils.ValueIgnoreEmpty(d.Get("auth_type")),
		"auth_value":            utils.ValueIgnoreEmpty(d.Get("auth_value")),
		"maximum_attempts":      utils.ValueIgnoreEmpty(d.Get("maximum_attempts")),
		"retry_attempt_seconds": utils.ValueIgnoreEmpty(d.Get("retry_attempt_seconds")),
	}
	staticParameter, err := buildStaticParameterBodyParams(d.Get("static_parameter"))
	if err != nil {
		return nil, err
	}
	bodyParams["static_parameter"] = staticParameter
	return bodyParams, nil
}

func buildResourceParameterBodyParams(resourceParameter interface{}) map[string]interface{} {
	if resourceParameter == nil {
		return nil
	}

	resourceParameterMap := resourceParameter.([]interface{})[0].(map[string]interface{})
	return map[string]interface{}{
		"resource_id": resourceParameterMap["resource_id"].(string),
	}
}

func buildStaticParameterBodyParams(staticParameter interface{}) ([]map[string]interface{}, error) {
	if rawArray, ok := staticParameter.([]interface{}); ok && len(rawArray) > 0 {
		staticParams := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			staticParamMap := v.(map[string]interface{})
			rawVarValue := staticParamMap["var_value"].(string)

			var varValue interface{}
			err := json.Unmarshal([]byte(rawVarValue), &varValue)
			if err != nil {
				return nil, err
			}
			staticParams[i] = map[string]interface{}{
				"var_key":   staticParamMap["var_key"].(string),
				"var_value": varValue,
			}
		}
		return staticParams, nil
	}
	return nil, nil
}

func resourceRemediationConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration"
		product = "rms"
	)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	getRemediationConfigurationPath := client.Endpoint + httpUrl
	getRemediationConfigurationPath = strings.ReplaceAll(getRemediationConfigurationPath, "{domain_id}", cfg.DomainID)
	getRemediationConfigurationPath = strings.ReplaceAll(getRemediationConfigurationPath, "{policy_assignment_id}", d.Id())

	getRemediationConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRemediationConfigurationResp, err := client.Request("GET", getRemediationConfigurationPath, &getRemediationConfigurationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS remediation configuration")
	}

	getRemediationConfigurationRespBody, err := utils.FlattenResponse(getRemediationConfigurationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("policy_assignment_id", d.Id()),
		d.Set("automatic", utils.PathSearch("automatic", getRemediationConfigurationRespBody, nil)),
		d.Set("target_type", utils.PathSearch("target_type", getRemediationConfigurationRespBody, nil)),
		d.Set("target_id", utils.PathSearch("target_id", getRemediationConfigurationRespBody, nil)),
		d.Set("static_parameter", flattenStaticParameter(getRemediationConfigurationRespBody)),
		d.Set("resource_parameter", flattentResourceParameter(getRemediationConfigurationRespBody)),
		d.Set("maximum_attempts", utils.PathSearch("maximum_attempts", getRemediationConfigurationRespBody, nil)),
		d.Set("retry_attempt_seconds", utils.PathSearch("retry_attempt_seconds", getRemediationConfigurationRespBody, nil)),
		d.Set("auth_type", utils.PathSearch("auth_type", getRemediationConfigurationRespBody, nil)),
		d.Set("auth_value", utils.PathSearch("auth_value", getRemediationConfigurationRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRemediationConfigurationRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRemediationConfigurationRespBody, nil)),
		d.Set("created_by", utils.PathSearch("created_by", getRemediationConfigurationRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStaticParameter(param interface{}) []map[string]interface{} {
	rawStaticParam := utils.PathSearch("static_parameter", param, nil)
	if rawStaticParam == nil {
		return nil
	}

	if rawArray, ok := rawStaticParam.([]interface{}); ok && len(rawArray) > 0 {
		staticParams := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			staticParamMap := v.(map[string]interface{})
			rawVarValue := utils.PathSearch("var_value", staticParamMap, nil)
			staticParams = append(staticParams, map[string]interface{}{
				"var_key":   utils.PathSearch("var_key", staticParamMap, nil),
				"var_value": utils.JsonToString(rawVarValue),
			})
		}
		return staticParams
	}
	return nil
}

func flattentResourceParameter(param interface{}) []interface{} {
	rawResourceParam := utils.PathSearch("resource_parameter", param, nil)
	if rawResourceParam == nil {
		return nil
	}

	resourceParamMap := rawResourceParam.(map[string]interface{})
	return []interface{}{
		map[string]interface{}{
			"resource_id": utils.PathSearch("resource_id", resourceParamMap, nil),
		},
	}
}

func resourceRemediationConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-configuration"
		product = "rms"
	)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	deleteRemediationConfigurationPath := client.Endpoint + httpUrl
	deleteRemediationConfigurationPath = strings.ReplaceAll(deleteRemediationConfigurationPath, "{domain_id}", cfg.DomainID)
	deleteRemediationConfigurationPath = strings.ReplaceAll(deleteRemediationConfigurationPath, "{policy_assignment_id}", d.Id())

	deleteRemediationConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteRemediationConfigurationPath, &deleteRemediationConfigurationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting RMS remediation configuration")
	}

	return nil
}
