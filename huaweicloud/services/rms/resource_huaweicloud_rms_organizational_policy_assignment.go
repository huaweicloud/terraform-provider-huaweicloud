package rms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Config PUT /v1/resource-manager/organizations/{organization_id}/policy-assignments
// @API Config GET /v1/resource-manager/organizations/{organization_id}/policy-assignments/{id}
// @API Config DELETE /v1/resource-manager/organizations/{organization_id}/policy-assignments/{id}
// @API Config GET /v1/resource-manager/organizations/{organization_id}/policy-assignment-statuses
func ResourceOrganizationalPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationalPolicyAssignmentCreateOrUpdate,
		ReadContext:   resourceOrganizationalPolicyAssignmentRead,
		UpdateContext: resourceOrganizationalPolicyAssignmentCreateOrUpdate,
		DeleteContext: resourceOrganizationalPolicyAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceOrganizationalPolicyAssignmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the organization.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the organizational policy assignment.",
			},
			"excluded_accounts": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The excluded accounts of the organizational policy assignment.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the organizational policy assignment.",
			},
			"policy_definition_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The ID of the policy definition.",
				ExactlyOneOf: []string{"function_urn"},
			},
			"function_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The function URN used to create the custom policy.",
			},
			"period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The period of the policy rule check.",
			},
			"policy_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the region to which the filtered resources belong.",
						},
						"resource_provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The service name to which the filtered resources belong.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The resource type of the filtered resources.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The resource ID used to filter a specified resources.",
						},
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The tag name used to filter resources.",
						},
						"tag_value": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							RequiredWith: []string{"policy_filter.0.tag_key"},
							Description:  "The tag value used to filter resources.",
						},
					},
				},
				Description: "The configuration used to filter resources.",
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsJSON,
				},
				Description: "The rule definition of the organizational policy assignment.",
			},
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation owner ID.",
			},
			"organization_policy_assignment_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation organization policy assignment URN.",
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

func buildOrganizationalPolicyFilter(filters []interface{}) map[string]interface{} {
	if len(filters) < 1 || filters[0] == nil {
		return nil
	}
	filter := filters[0].(map[string]interface{})
	return map[string]interface{}{
		"region_id":         utils.ValueIgnoreEmpty(filter["region"]),
		"resource_provider": utils.ValueIgnoreEmpty(filter["resource_provider"]),
		"resource_type":     utils.ValueIgnoreEmpty(filter["resource_type"]),
		"resource_id":       utils.ValueIgnoreEmpty(filter["resource_id"]),
		"tag_key":           utils.ValueIgnoreEmpty(filter["tag_key"]),
		"tag_value":         utils.ValueIgnoreEmpty(filter["tag_value"]),
	}
}

func buildOrganizationalRuleParameters(parameters map[string]interface{}) (map[string]interface{}, error) {
	if parameters == nil {
		return nil, nil
	}
	result := make(map[string]interface{})
	for k, jsonVal := range parameters {
		var value interface{}
		err := json.Unmarshal([]byte(jsonVal.(string)), &value)
		if err != nil {
			return result, fmt.Errorf("error analyzing parameter value: %s", err)
		}
		result[k] = map[string]interface{}{
			"value": value,
		}
	}
	return result, nil
}

func buildOrganizationalPolicyAssignmentCreateOpts(d *schema.ResourceData) (map[string]interface{}, error) {
	metadata := map[string]interface{}{
		"description":          d.Get("description").(string),
		"policy_filter":        buildOrganizationalPolicyFilter(d.Get("policy_filter").([]interface{})),
		"policy_definition_id": utils.ValueIgnoreEmpty(d.Get("policy_definition_id").(string)),
		"period":               utils.ValueIgnoreEmpty(d.Get("period").(string)),
		"function_urn":         utils.ValueIgnoreEmpty(d.Get("function_urn").(string)),
	}

	parameters, err := buildOrganizationalRuleParameters(d.Get("parameters").(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	metadata["parameters"] = parameters

	res := map[string]interface{}{
		"organization_policy_assignment_name": d.Get("name"),
		"excluded_accounts":                   d.Get("excluded_accounts").(*schema.Set).List(),
	}

	if _, ok := d.GetOk("policy_definition_id"); ok {
		res["managed_policy_assignment_metadata"] = metadata
	} else {
		res["custom_policy_assignment_metadata"] = metadata
	}

	return res, nil
}

func resourceOrganizationalPolicyAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrganizationalPolicyAssignment: Create a RMS organizational policy assignment.
	var (
		createOrgPolicyAssignmentHttpUrl = "v1/resource-manager/organizations/{organization_id}/policy-assignments"
		createOrgPolicyAssignmentProduct = "rms"
	)
	createOrgPolicyAssignmentClient, err := cfg.NewServiceClient(createOrgPolicyAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	createOrgPolicyAssignmentPath := createOrgPolicyAssignmentClient.Endpoint + createOrgPolicyAssignmentHttpUrl
	createOrgPolicyAssignmentPath = strings.ReplaceAll(createOrgPolicyAssignmentPath, "{organization_id}", d.Get("organization_id").(string))

	createOrgPolicyAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildOrganizationalPolicyAssignmentCreateOpts(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOrgPolicyAssignmentOpt.JSONBody = utils.RemoveNil(createOpts)
	createOrgPolicyAssignmentResp, err := createOrgPolicyAssignmentClient.Request("PUT",
		createOrgPolicyAssignmentPath, &createOrgPolicyAssignmentOpt)
	if err != nil {
		return diag.Errorf("error creating RMS organizational policy assignment: %s", err)
	}

	createOrgPolicyAssignmentRespBody, err := utils.FlattenResponse(createOrgPolicyAssignmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	waitTimeout := d.Timeout(schema.TimeoutUpdate)

	if d.IsNewResource() {
		id := utils.PathSearch("organization_policy_assignment_id", createOrgPolicyAssignmentRespBody, "").(string)
		if id == "" {
			return diag.Errorf("error creating RMS assignment package: ID is not found in API response")
		}
		d.SetId(id)

		waitTimeout = d.Timeout(schema.TimeoutCreate)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{"CREATE_SUCCESSFUL", "UPDATE_SUCCESSFUL"},
		Pending:      []string{"CREATE_IN_PROGRESS", "UPDATE_IN_PROGRESS"},
		Refresh:      refreshOrgPolicyAssignmentDeployStatus(d, createOrgPolicyAssignmentClient),
		Timeout:      waitTimeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS organizational assignment Package (%s) to be created: %s",
			d.Id(), err)
	}

	return resourceOrganizationalPolicyAssignmentRead(ctx, d, meta)
}

func flattenOrganizationalPolicyFilter(filter interface{}) []map[string]interface{} {
	if filter == nil {
		return nil
	}

	filterMap := filter.(map[string]interface{})
	return []map[string]interface{}{
		{
			"region":            filterMap["region"],
			"resource_provider": filterMap["resource_provider"],
			"resource_type":     filterMap["resource_type"],
			"resource_id":       filterMap["resource_id"],
			"tag_key":           filterMap["tag_key"],
			"tag_value":         filterMap["tag_value"],
		},
	}
}

func flattenOrganizationalPolicyParameters(parameters interface{}) (map[string]interface{}, error) {
	if parameters == nil {
		return nil, nil
	}

	result := make(map[string]interface{})
	for k, v := range parameters.(map[string]interface{}) {
		val := v.(map[string]interface{})
		jsonBytes, err := json.Marshal(val["value"])
		if err != nil {
			return nil, fmt.Errorf("generate json string of %s failed: %s", k, err)
		}
		result[k] = string(jsonBytes)
	}
	return result, nil
}

func resourceOrganizationalPolicyAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganizationalPolicyAssignment: Query the RMS organizational policy assignment
	var (
		getOrgPolicyAssignmentHttpUrl = "v1/resource-manager/organizations/{organization_id}/policy-assignments/{id}"
		getOrgPolicyAssignmentProduct = "rms"
	)
	getOrgPolicyAssignmentClient, err := cfg.NewServiceClient(getOrgPolicyAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getOrgPolicyAssignmentPath := getOrgPolicyAssignmentClient.Endpoint + getOrgPolicyAssignmentHttpUrl
	getOrgPolicyAssignmentPath = strings.ReplaceAll(getOrgPolicyAssignmentPath, "{organization_id}",
		d.Get("organization_id").(string))
	getOrgPolicyAssignmentPath = strings.ReplaceAll(getOrgPolicyAssignmentPath, "{id}", d.Id())

	getOrgPolicyAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrgPolicyAssignmentResp, err := getOrgPolicyAssignmentClient.Request("GET", getOrgPolicyAssignmentPath,
		&getOrgPolicyAssignmentOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS organizational assignment package")
	}

	getOrgPolicyAssignmentRespBody, err := utils.FlattenResponse(getOrgPolicyAssignmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	parameters, err := flattenOrganizationalPolicyParameters(utils.PathSearch("parameters", getOrgPolicyAssignmentRespBody, nil))
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("organization_id", utils.PathSearch("organization_id", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("name", utils.PathSearch("organization_policy_assignment_name", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("excluded_accounts", utils.PathSearch("excluded_accounts", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("policy_definition_id", utils.PathSearch("policy_definition_id", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("function_urn", utils.PathSearch("function_urn", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("period", utils.PathSearch("period", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("policy_filter", flattenOrganizationalPolicyFilter(utils.PathSearch("policy_filter", getOrgPolicyAssignmentRespBody, nil))),
		d.Set("parameters", parameters),
		d.Set("owner_id", utils.PathSearch("owner_id", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("organization_policy_assignment_urn", utils.PathSearch("organization_policy_assignment_urn",
			getOrgPolicyAssignmentRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getOrgPolicyAssignmentRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getOrgPolicyAssignmentRespBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving organizational policy assignment resource (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceOrganizationalPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrganizationalPolicyAssignment: Delete an existing RMS organizational assignment package
	var (
		deleteOrgPolicyAssignmentHttpUrl = "v1/resource-manager/organizations/{organization_id}/policy-assignments/{id}"
		deleteOrgPolicyAssignmentProduct = "rms"
	)
	deleteOrgPolicyAssignmentClient, err := cfg.NewServiceClient(deleteOrgPolicyAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	deleteOrgPolicyAssignmentPath := deleteOrgPolicyAssignmentClient.Endpoint + deleteOrgPolicyAssignmentHttpUrl
	deleteOrgPolicyAssignmentPath = strings.ReplaceAll(deleteOrgPolicyAssignmentPath, "{organization_id}",
		d.Get("organization_id").(string))
	deleteOrgPolicyAssignmentPath = strings.ReplaceAll(deleteOrgPolicyAssignmentPath, "{id}", d.Id())

	deleteOrgPolicyAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteOrgPolicyAssignmentClient.Request("DELETE", deleteOrgPolicyAssignmentPath,
		&deleteOrgPolicyAssignmentOpt)
	if err != nil {
		return diag.Errorf("error deleting RMS organizational assignment package: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{""},
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Refresh:      refreshOrgPolicyAssignmentDeployStatus(d, deleteOrgPolicyAssignmentClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS organizational assignment Package (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func refreshOrgPolicyAssignmentDeployStatus(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		// getDeployStatus: Query the RMS organizational assignment
		var (
			getDeployStatusHttpUrl = "v1/resource-manager/organizations/{organization_id}/policy-assignment-statuses"
		)

		getDeployStatusPath := client.Endpoint + getDeployStatusHttpUrl
		getDeployStatusPath = strings.ReplaceAll(getDeployStatusPath, "{organization_id}", d.Get("organization_id").(string))

		getDeployStatusPath += fmt.Sprintf("?organization_policy_assignment_id=%s", d.Id())

		getDeployStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getDeployStatusResp, err := client.Request("GET", getDeployStatusPath, &getDeployStatusOpt)
		if err != nil {
			return nil, "", err
		}
		getDeployStatusRespBody, err := utils.FlattenResponse(getDeployStatusResp)
		if err != nil {
			return nil, "", err
		}
		return getDeployStatusRespBody, utils.PathSearch("value[0].organization_policy_assignment_status",
			getDeployStatusRespBody, "").(string), nil
	}
}

// resourceOrganizationalPolicyAssignmentImportState use to import an id with format <organization_id>/<id>
func resourceOrganizationalPolicyAssignmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <organization_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("organization_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
