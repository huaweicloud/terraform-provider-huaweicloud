package rfs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var stackSetDeploymentNonUpdatableParams = []string{
	"stack_set_name",
	"deployment_targets",
	"deployment_targets.*.regions",
	"deployment_targets.*.domain_ids",
	"deployment_targets.*.domain_ids_uri",
	"deployment_targets.*.organizational_unit_ids",
	"deployment_targets.*.domain_id_filter_type",
	"stack_set_id",
	"template_body",
	"template_uri",
	"vars_uri",
	"vars_body",
	"operation_preferences",
	"operation_preferences.*.region_concurrency_type",
	"operation_preferences.*.region_order",
	"operation_preferences.*.failure_tolerance_count",
	"operation_preferences.*.failure_tolerance_percentage",
	"operation_preferences.*.max_concurrent_count",
	"operation_preferences.*.max_concurrent_percentage",
	"operation_preferences.*.failure_tolerance_mode",
	"call_identity",
}

// @API RFS POST /v1/stack-sets/{stack_set_name}/deployments
func ResourceStackSetDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStackSetDeploymentCreate,
		ReadContext:   resourceStackSetDeploymentRead,
		UpdateContext: resourceStackSetDeploymentUpdate,
		DeleteContext: resourceStackSetDeploymentDelete,

		CustomizeDiff: config.FlexibleForceNew(stackSetDeploymentNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_targets": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     deploymentTargetsSchema(),
			},
			"stack_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vars_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vars_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_preferences": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     operationPreferencesSchema(),
			},
			"call_identity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func operationPreferencesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_concurrency_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_order": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"failure_tolerance_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"failure_tolerance_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrent_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_concurrent_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"failure_tolerance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func deploymentTargetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_ids_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizational_unit_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_id_filter_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildStackSetDeploymentBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"deployment_targets":    buildDeploymentTargetsParams(d.Get("deployment_targets").([]interface{})),
		"stack_set_id":          utils.ValueIgnoreEmpty(d.Get("stack_set_id")),
		"template_body":         utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":          utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_uri":              utils.ValueIgnoreEmpty(d.Get("vars_uri")),
		"vars_body":             utils.ValueIgnoreEmpty(d.Get("vars_body")),
		"operation_preferences": buildOperationPreferencesParams(d.Get("operation_preferences").([]interface{})),
		"call_identity":         utils.ValueIgnoreEmpty(d.Get("call_identity")),
	}

	return body
}

func buildStringListParams(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	return rawArray
}

func buildDeploymentTargetsParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"regions":                 rawMap["regions"],
		"domain_ids":              buildStringListParams(rawMap["domain_ids"].([]interface{})),
		"domain_ids_uri":          utils.ValueIgnoreEmpty(rawMap["domain_ids_uri"]),
		"organizational_unit_ids": buildStringListParams(rawMap["organizational_unit_ids"].([]interface{})),
		"domain_id_filter_type":   utils.ValueIgnoreEmpty(rawMap["domain_id_filter_type"]),
	}
}

func buildOperationPreferencesParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"region_concurrency_type":      utils.ValueIgnoreEmpty(rawMap["region_concurrency_type"]),
		"region_order":                 buildStringListParams(rawMap["region_order"].([]interface{})),
		"failure_tolerance_count":      utils.ValueIgnoreEmpty(rawMap["failure_tolerance_count"]),
		"failure_tolerance_percentage": utils.ValueIgnoreEmpty(rawMap["failure_tolerance_percentage"]),
		"max_concurrent_count":         utils.ValueIgnoreEmpty(rawMap["max_concurrent_count"]),
		"max_concurrent_percentage":    utils.ValueIgnoreEmpty(rawMap["max_concurrent_percentage"]),
		"failure_tolerance_mode":       utils.ValueIgnoreEmpty(rawMap["failure_tolerance_mode"]),
	}
}

func resourceStackSetDeploymentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		stackSetName = d.Get("stack_set_name").(string)
		httpUrl      = "v1/stack-sets/{stack_set_name}/deployments"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", stackSetName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildStackSetDeploymentBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error triggering RFS stack set deployment: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationId := utils.PathSearch("stack_set_operation_id", respBody, "").(string)
	if operationId == "" {
		return diag.Errorf("error triggering RFS stack set deployment: Operation ID is not found in API response")
	}

	d.SetId(operationId)

	return nil
}

func resourceStackSetDeploymentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceStackSetDeploymentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceStackSetDeploymentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to trigger stack set deployment. Deleting this resource
    will not cancel the deployment operation, but will only remove resource information from
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
