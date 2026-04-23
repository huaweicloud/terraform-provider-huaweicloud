package rfs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/{project_id}/stacks/{stack_name}/execution-plans
// @API RFS GET /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/metadata
// @API RFS DELETE /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}
func ResourceExecutionPlanV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExecutionPlanV2Create,
		ReadContext:   resourceExecutionPlanV2Read,
		UpdateContext: resourceExecutionPlanV2Update,
		DeleteContext: resourceExecutionPlanV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceExecutionPlanV2ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"stack_name",
			"execution_plan_name",
			"stack_id",
			"template_body",
			"template_uri",
			"description",
			"vars_structure",
			"vars_structure.*.var_key",
			"vars_structure.*.var_value",
			"vars_structure.*.encryption",
			"vars_structure.*.encryption.*.kms",
			"vars_structure.*.encryption.*.kms.*.id",
			"vars_structure.*.encryption.*.kms.*.cipher_text",
			"vars_body",
			"vars_uri",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `template_body` has no response value.
			"template_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Field `template_uri` has no response value.
			"template_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vars_structure": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     varsStructureSchema(),
			},
			"vars_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `vars_uri` has no response value.
			"vars_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"execution_plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vars_uri_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apply_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"summary": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     varsStructureSummarySchema(),
			},
		},
	}
}

func varsStructureSummarySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_add": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_update": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_delete": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_import": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func varsStructureSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"var_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"var_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encryption": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     varsStructureEncryptionSchema(),
			},
		},
	}
}

func varsStructureEncryptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kms": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     varsStructureEncryptionKmsSchema(),
			},
		},
	}
}

func varsStructureEncryptionKmsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cipher_text": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildCreateExecutionPlanV2BodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"execution_plan_name": d.Get("execution_plan_name"),
		"stack_id":            utils.ValueIgnoreEmpty(d.Get("stack_id")),
		"template_body":       utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":        utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"vars_structure":      buildVarsStructureParams(d.Get("vars_structure").([]interface{})),
		"vars_body":           utils.ValueIgnoreEmpty(d.Get("vars_body")),
		"vars_uri":            utils.ValueIgnoreEmpty(d.Get("vars_uri")),
	}
}

func buildVarsStructureParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"var_key":    rawMap["var_key"],
			"var_value":  rawMap["var_value"],
			"encryption": buildVarsStructureEncryptionParams(rawMap["encryption"].([]interface{})),
		})
	}

	return result
}

func buildVarsStructureEncryptionParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"kms": buildVarsStructureEncryptionKmsParams(rawMap["kms"].([]interface{})),
	}
}

func buildVarsStructureEncryptionKmsParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":          rawMap["id"],
		"cipher_text": rawMap["cipher_text"],
	}
}

func ReadExecutionPlanV2Detail(client *golangsdk.ServiceClient, reqUUID, stackName, planName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", planName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForExecutionPlanV2Available(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, reqUUID string) error {
	stackName := d.Get("stack_name").(string)
	executionPlanName := d.Get("execution_plan_name").(string)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadExecutionPlanV2Detail(client, reqUUID, stackName, executionPlanName)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("unable to find status in API response")
			}

			if status == "CREATION_FAILED" {
				return respBody, status, nil
			}

			if status == "AVAILABLE" {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceExecutionPlanV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		product           = "rfs"
		httpUrl           = "v1/{project_id}/stacks/{stack_name}/execution-plans"
		stackName         = d.Get("stack_name").(string)
		executionPlanName = d.Get("execution_plan_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateExecutionPlanV2BodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating RFS execution plan: %s", err)
	}

	d.SetId(executionPlanName)

	if err := waitingForExecutionPlanV2Available(ctx, client, d, d.Timeout(schema.TimeoutCreate), reqUUID); err != nil {
		return diag.Errorf("error waiting for RFS execution plan to be available: %s", err)
	}

	return resourceExecutionPlanV2Read(ctx, d, meta)
}

func flattenEncryptionKms(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":          utils.PathSearch("id", respBody, nil),
			"cipher_text": utils.PathSearch("cipher_text", respBody, nil),
		},
	}
}

func flattenEncryptionStructure(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"kms": flattenEncryptionKms(utils.PathSearch("kms", respBody, nil)),
		},
	}
}

func flattenVarsStructure(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok || len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		result = append(result, map[string]interface{}{
			"var_key":    utils.PathSearch("var_key", v, nil),
			"var_value":  utils.PathSearch("var_value", v, nil),
			"encryption": flattenEncryptionStructure(utils.PathSearch("encryption", v, nil)),
		})
	}

	return result
}

func flattenExecutionPlanSummary(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"resource_add":    utils.PathSearch("resource_add", respBody, nil),
			"resource_update": utils.PathSearch("resource_update", respBody, nil),
			"resource_delete": utils.PathSearch("resource_delete", respBody, nil),
			"resource_import": utils.PathSearch("resource_import", respBody, nil),
		},
	}
}

func resourceExecutionPlanV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		product           = "rfs"
		stackName         = d.Get("stack_name").(string)
		executionPlanName = d.Get("execution_plan_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	respBody, err := ReadExecutionPlanV2Detail(client, reqUUID, stackName, executionPlanName)
	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving RFS execution plan")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_name", utils.PathSearch("stack_name", respBody, nil)),
		d.Set("execution_plan_name", utils.PathSearch("execution_plan_name", respBody, nil)),
		d.Set("stack_id", utils.PathSearch("stack_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("vars_structure", flattenVarsStructure(utils.PathSearch("vars_structure", respBody, nil))),
		d.Set("vars_body", utils.PathSearch("vars_body", respBody, nil)),
		d.Set("execution_plan_id", utils.PathSearch("execution_plan_id", respBody, nil)),
		d.Set("vars_uri_content", utils.PathSearch("vars_uri_content", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("status_message", utils.PathSearch("status_message", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("apply_time", utils.PathSearch("apply_time", respBody, nil)),
		d.Set("summary", flattenExecutionPlanSummary(utils.PathSearch("summary", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceExecutionPlanV2Update(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceExecutionPlanV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		product           = "rfs"
		stackName         = d.Get("stack_name").(string)
		executionPlanName = d.Get("execution_plan_name").(string)
		httpUrl           = "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", executionPlanName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting RFS execution plan: %s", err)
	}

	return nil
}

func resourceExecutionPlanV2ImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID,"+
			" want '<stack_name>/<execution_plan_name>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(
		d.Set("stack_name", parts[0]),
		d.Set("execution_plan_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
