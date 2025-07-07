package rms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/rms/v1/policyassignments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	AssignmentTypeBuiltin = "builtin"
	AssignmentTypeCustom  = "custom"

	AssignmentStatusDisabled   = "Disabled"
	AssignmentStatusEnabled    = "Enabled"
	AssignmentStatusEvaluating = "Evaluating"
)

// @API Config PUT /v1/resource-manager/domains/{domain_id}/policy-assignments
// @API Config GET /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}
// @API Config POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/disable
// @API Config POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/enable
// @API Config PUT /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}
// @API Config POST /v1/resource-manager/{resource_type}/{resource_id}/tags/create
// @API Config POST /v1/resource-manager/{resource_type}/{resource_id}/tags/delete
// @API Config GET /v1/resource-manager/{resource_type}/{resource_id}/tags
func ResourcePolicyAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceePolicyAssignmentCreate,
		ReadContext:   resourceePolicyAssignmentRead,
		UpdateContext: resourceePolicyAssignmentUpdate,
		DeleteContext: resourceePolicyAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the policy assignment.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy assignment.",
			},
			"policy_definition_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the policy definition.",
			},
			"period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The period of the policy rule check.",
			},
			"policy_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the region to which the filtered resources belong.",
						},
						"resource_provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The service name to which the filtered resources belong.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource type of the filtered resources.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The resource ID used to filter a specified resources.",
						},
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The tag name used to filter resources.",
						},
						"tag_value": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"policy_filter.0.tag_key"},
							Description:  "The tag value used to filter resources.",
						},
					},
				},
				Description: "The configuration used to filter resources.",
			},
			"custom_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_urn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The function URN used to create the custom policy.",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The authorization type of the custom policy.",
						},
						"auth_value": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsJSON,
							},
							Description: "The authorization value of the custom policy.",
						},
					},
				},
				Description: "The configuration of the custom policy.",
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
				},
				Description: "The rule definition of the policy assignment.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The expect status of the policy.",
			},
			"tags": common.TagsSchema(),
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the policy assignment.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time.",
			},
		},
	}
}

func buildPolicyFilter(filters []interface{}) policyassignments.PolicyFilter {
	if len(filters) < 1 {
		return policyassignments.PolicyFilter{}
	}
	filter := filters[0].(map[string]interface{})
	return policyassignments.PolicyFilter{
		RegionId:         filter["region"].(string),
		ResourceProvider: filter["resource_provider"].(string),
		ResourceType:     filter["resource_type"].(string),
		ResourceId:       filter["resource_id"].(string),
		TagKey:           filter["tag_key"].(string),
		TagValue:         filter["tag_value"].(string),
	}
}

func buildCustomPolicy(policies []interface{}) (*policyassignments.CustomPolicy, error) {
	if len(policies) < 1 {
		return nil, nil
	}
	policy := policies[0].(map[string]interface{})
	result := policyassignments.CustomPolicy{
		FunctionUrn: policy["function_urn"].(string),
		AuthType:    policy["auth_type"].(string),
	}
	authValues := make(map[string]interface{})
	for k, jsonVal := range policy["auth_value"].(map[string]interface{}) {
		var value interface{}
		err := json.Unmarshal([]byte(jsonVal.(string)), &value)
		if err != nil {
			return &result, fmt.Errorf("error analyzing authorization value: %s", err)
		}
		authValues[k] = value
	}
	result.AuthValue = authValues

	return &result, nil
}

func buildRuleParameters(parameters map[string]interface{}) (map[string]policyassignments.PolicyParameterValue, error) {
	if len(parameters) < 1 {
		return nil, nil
	}
	result := make(map[string]policyassignments.PolicyParameterValue)
	for k, jsonVal := range parameters {
		var value interface{}
		err := json.Unmarshal([]byte(jsonVal.(string)), &value)
		if err != nil {
			return result, fmt.Errorf("error analyzing parameter value: %s", err)
		}
		result[k] = policyassignments.PolicyParameterValue{
			Value: value,
		}
	}
	return result, nil
}

func buildPolicyAssignmentCreateOpts(d *schema.ResourceData) (policyassignments.CreateOpts, error) {
	result := policyassignments.CreateOpts{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               AssignmentTypeBuiltin,
		PolicyFilter:       buildPolicyFilter(d.Get("policy_filter").([]interface{})),
		PolicyDefinitionId: d.Get("policy_definition_id").(string),
		Period:             d.Get("period").(string),
	}
	customPolicy, err := buildCustomPolicy(d.Get("custom_policy").([]interface{}))
	if err != nil {
		return result, err
	}
	result.CustomPolicy = customPolicy
	if customPolicy != nil {
		result.Type = AssignmentTypeCustom
	}

	parameters, err := buildRuleParameters(d.Get("parameters").(map[string]interface{}))
	if err != nil {
		return result, err
	}
	result.Parameters = parameters

	return result, nil
}

func updatePolicyAssignmentStatus(client *golangsdk.ServiceClient, domainId, assignmentId,
	statusConfig string) (err error) {
	switch statusConfig {
	case AssignmentStatusDisabled:
		err = policyassignments.Disable(client, domainId, assignmentId)
	case AssignmentStatusEnabled:
		err = policyassignments.Enable(client, domainId, assignmentId)
	}
	return
}

func policyAssignmentRefreshFunc(client *golangsdk.ServiceClient, domainId,
	assignmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := policyassignments.Get(client, domainId, assignmentId)
		if err != nil {
			return resp, "ERROR", err
		}
		return resp, resp.Status, nil
	}
}

func createResourceTags(client *golangsdk.ServiceClient, id string, tags map[string]interface{}) error {
	createTagsHttpUrl := "v1/resource-manager/{resource_type}/{resource_id}/tags/create"
	createTagsPath := client.Endpoint + createTagsHttpUrl
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", "config:policyAssignments")
	createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", id)

	tagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	tagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}
	_, err := client.Request("POST", createTagsPath, &tagsOpt)
	return err
}

func resourceePolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.RmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}
	opts, err := buildPolicyAssignmentCreateOpts(d)
	if err != nil {
		return diag.Errorf("error creating the create option structure of the RMS policy assignment: %s", err)
	}
	domainId := cfg.DomainID
	resp, err := policyassignments.Create(client, domainId, opts)
	if err != nil {
		return diag.Errorf("error creating policy assignment resource: %s", err)
	}

	assignmentId := resp.ID
	d.SetId(assignmentId)

	// it will take too long time to become enabled when the resources are very huge.
	// so we wait for the enabled status only when user want to disable it during creating.
	if statusConfig := d.Get("status").(string); statusConfig == AssignmentStatusDisabled {
		log.Printf("[DEBUG] Waiting for the policy assignment (%s) status to become enabled, then disable it", assignmentId)
		stateConf := &resource.StateChangeConf{
			Pending:                   []string{AssignmentStatusDisabled, AssignmentStatusEvaluating},
			Target:                    []string{AssignmentStatusEnabled},
			Refresh:                   policyAssignmentRefreshFunc(client, domainId, assignmentId),
			Timeout:                   d.Timeout(schema.TimeoutCreate),
			Delay:                     10 * time.Second,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 2,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the policy assignment (%s) status to become enabled: %s",
				assignmentId, err)
		}

		err = updatePolicyAssignmentStatus(client, domainId, assignmentId, statusConfig)
		if err != nil {
			return diag.Errorf("error disabling the status of the policy assignment: %s", err)
		}
	}
	if rawTags := d.Get("tags").(map[string]interface{}); len(rawTags) > 0 {
		err = createResourceTags(client, assignmentId, rawTags)
		if err != nil {
			return diag.Errorf("error creating the policy assignment tags: %s", err)
		}
	}

	return resourceePolicyAssignmentRead(ctx, d, meta)
}

func flattenPolicyFilter(filter policyassignments.PolicyFilter) []map[string]interface{} {
	if reflect.DeepEqual(filter, policyassignments.PolicyFilter{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"region":            filter.RegionId,
			"resource_provider": filter.ResourceProvider,
			"resource_type":     filter.ResourceType,
			"resource_id":       filter.ResourceId,
			"tag_key":           filter.TagKey,
			"tag_value":         filter.TagValue,
		},
	}
}

func flattenCustomPolicy(customPolicy policyassignments.CustomPolicy) ([]map[string]interface{}, error) {
	if reflect.DeepEqual(customPolicy, policyassignments.CustomPolicy{}) {
		return nil, nil
	}

	authValues := make(map[string]interface{})
	for k, v := range customPolicy.AuthValue {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("generate json string failed: %s", err)
		}
		authValues[k] = string(jsonBytes)
	}
	return []map[string]interface{}{
		{
			"function_urn": customPolicy.FunctionUrn,
			"auth_type":    customPolicy.AuthType,
			"auth_value":   authValues,
		},
	}, nil
}

func flattenPolicyParameters(parameters map[string]policyassignments.PolicyParameterValue) (map[string]interface{},
	error) {
	if len(parameters) < 1 {
		return nil, nil
	}

	result := make(map[string]interface{})
	for k, v := range parameters {
		jsonBytes, err := json.Marshal(v.Value)
		if err != nil {
			return nil, fmt.Errorf("generate json string failed: %s", err)
		}
		result[k] = string(jsonBytes)
	}
	return result, nil
}

func getTags(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getTagsHttpUrl := "v1/resource-manager/{resource_type}/{resource_id}/tags"
	getTagsPath := client.Endpoint + getTagsHttpUrl
	getTagsPath = strings.ReplaceAll(getTagsPath, "{resource_type}", "config:policyAssignments")
	getTagsPath = strings.ReplaceAll(getTagsPath, "{resource_id}", id)

	tagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	tagsResp, err := client.Request("GET", getTagsPath, &tagsOpt)
	if err != nil {
		return nil, err
	}

	tagsRespBody, err := utils.FlattenResponse(tagsResp)
	if err != nil {
		return nil, err
	}
	return tagsRespBody, nil
}

func resourceePolicyAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.RmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	assignmentId := d.Id()
	resp, err := policyassignments.Get(client, cfg.DomainID, assignmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "RMS policy assignment")
	}

	customPolicy, err := flattenCustomPolicy(resp.CustomPolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	parameters, err := flattenPolicyParameters(resp.Parameters)
	if err != nil {
		return diag.FromErr(err)
	}
	tagsRespBody, err := getTags(client, assignmentId)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("type", resp.Type),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("policy_definition_id", resp.PolicyDefinitionId),
		d.Set("period", resp.Period),
		d.Set("policy_filter", flattenPolicyFilter(resp.PolicyFilter)),
		d.Set("custom_policy", customPolicy),
		d.Set("parameters", parameters),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", tagsRespBody, nil))),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving policy assignment resource (%s) fields: %s", assignmentId, mErr)
	}
	return nil
}

func buildPolicyAssignmentUpdateOpts(d *schema.ResourceData) (policyassignments.UpdateOpts, error) {
	result := policyassignments.UpdateOpts{
		Name:               d.Get("name").(string),
		Description:        utils.String(d.Get("description").(string)),
		Type:               AssignmentTypeBuiltin,
		PolicyFilter:       buildPolicyFilter(d.Get("policy_filter").([]interface{})),
		PolicyDefinitionId: d.Get("policy_definition_id").(string),
		Period:             d.Get("period").(string),
	}
	customPolicy, err := buildCustomPolicy(d.Get("custom_policy").([]interface{}))
	if err != nil {
		return result, err
	}
	result.CustomPolicy = customPolicy
	if customPolicy != nil {
		result.Type = AssignmentTypeCustom
	}

	parameters, err := buildRuleParameters(d.Get("parameters").(map[string]interface{}))
	if err != nil {
		return result, err
	}
	result.Parameters = parameters

	return result, nil
}

func deleteResourceTags(client *golangsdk.ServiceClient, id string, tags map[string]interface{}) error {
	deleteTagsHttpUrl := "v1/resource-manager/{resource_type}/{resource_id}/tags/delete"
	deleteTagsPath := client.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", "config:policyAssignments")
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", id)

	tagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	tagsOpt.JSONBody = map[string]interface{}{
		"tags": utils.ExpandResourceTags(tags),
	}
	_, err := client.Request("POST", deleteTagsPath, &tagsOpt)
	return err
}

func updateResourceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldTags, newTags := d.GetChange("tags")

	// remove old tags
	if oMap := oldTags.(map[string]interface{}); len(oMap) > 0 {
		err := deleteResourceTags(client, d.Id(), oMap)
		if err != nil {
			return err
		}
	}

	// set new tags
	if nMap := newTags.(map[string]interface{}); len(nMap) > 0 {
		err := createResourceTags(client, d.Id(), nMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceePolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.RmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	assignmentId := d.Id()
	domainId := cfg.DomainID

	if d.HasChange("status") {
		oldVal, newVal := d.GetChange("status")
		err = updatePolicyAssignmentStatus(client, domainId, d.Id(), d.Get("status").(string))
		if err != nil {
			return diag.Errorf("error updating the status of the policy assignment (%s): %s", assignmentId, err)
		}

		if newVal.(string) == AssignmentStatusEnabled {
			log.Printf("[DEBUG] Waiting for the policy assignment (%s) status to become %s.", assignmentId,
				strings.ToLower(newVal.(string)))
			stateConf := &resource.StateChangeConf{
				Pending:                   []string{oldVal.(string)},
				Target:                    []string{AssignmentStatusEvaluating, AssignmentStatusEnabled},
				Refresh:                   policyAssignmentRefreshFunc(client, domainId, assignmentId),
				Timeout:                   d.Timeout(schema.TimeoutUpdate),
				Delay:                     10 * time.Second,
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 2,
			}
			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("error waiting for the policy assignment (%s) status to become %s: %s",
					assignmentId, strings.ToLower(newVal.(string)), err)
			}
		}
	}
	if d.HasChangesExcept("status", "tags") {
		opts, err := buildPolicyAssignmentUpdateOpts(d)
		if err != nil {
			return diag.Errorf("error creating the update option structure of the RMS policy assignment: %s", err)
		}

		_, err = policyassignments.Update(client, domainId, assignmentId, opts)
		if err != nil {
			return diag.Errorf("error updating policy assignment resource (%s): %s", assignmentId, err)
		}
		currentStatus := d.Get("status").(string)
		log.Printf("[DEBUG] Waiting for the policy assignment (%s) status to become %s.", assignmentId,
			strings.ToLower(currentStatus))
		stateConf := &resource.StateChangeConf{
			Target:                    []string{currentStatus},
			Refresh:                   policyAssignmentRefreshFunc(client, domainId, assignmentId),
			Timeout:                   d.Timeout(schema.TimeoutUpdate),
			Delay:                     10 * time.Second,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 2,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the policy assignment (%s) status to become %s: %s",
				assignmentId, strings.ToLower(currentStatus), err)
		}
	}
	if d.HasChange("tags") {
		err := updateResourceTags(client, d)
		if err != nil {
			return diag.Errorf("error updating the policy assignment tags: %s", err)
		}
	}

	return resourceePolicyAssignmentRead(ctx, d, meta)
}

func resourceePolicyAssignmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.RmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}
	var (
		assignmentId = d.Id()
		domainId     = cfg.DomainID
	)
	if d.Get("status").(string) == AssignmentStatusEnabled {
		// Before delete policy assignment, we need to disable it.
		err = policyassignments.Disable(client, domainId, assignmentId)
		if err != nil {
			return diag.Errorf("failed to disable the policy assignment (%s): %s", assignmentId, err)
		}
	}

	err = policyassignments.Delete(client, domainId, assignmentId)
	if err != nil {
		return diag.Errorf("error deleting the policy assignment (%s): %s", assignmentId, err)
	}
	return nil
}
