package dataarts

import (
	"context"
	"errors"
	"fmt"
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

var (
	securityResourcePermissionPolicyNonUpdatableParams = []string{
		"workspace_id",
	}
	securityResourcePermissionPolicyResourceNotFoundCodes = []string{
		"DLS.3080", // The resource permission policy is not found.
		"DLS.6036", // The Workspace ID is not found, error code key is "errorCode".
	}
)

// @API DataArtsStudio POST /v1/{project_id}/security/permission-resource
// @API DataArtsStudio GET /v1/{project_id}/security/permission-resource/{policy_id}
// @API DataArtsStudio PUT /v1/{project_id}/security/permission-resource/{policy_id}
// @API DataArtsStudio POST /v1/{project_id}/security/permission-resource/batch-delete
func ResourceSecurityResourcePermissionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityResourcePermissionPolicyCreate,
		ReadContext:   resourceSecurityResourcePermissionPolicyRead,
		UpdateContext: resourceSecurityResourcePermissionPolicyUpdate,
		DeleteContext: resourceSecurityResourcePermissionPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(securityResourcePermissionPolicyNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityResourcePermissionPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource permission policy is located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the resource permission policy belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the resource permission policy.`,
			},
			"resources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `The list of resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the resource.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the resource.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the resource.`,
						},
					},
				},
			},
			"members": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `The list of members.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the member.`,
						},
						"member_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the member.`,
						},
						"member_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the member.`,
						},
					},
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the resource permission policy, in RFC3339 format.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the resource permission policy.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the resource permission policy, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildSecurityResourcePermissionPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name": d.Get("name"),
		"resources":   buildSecurityResourcePermissionPolicyResources(d.Get("resources").(*schema.Set)),
		"members":     buildSecurityResourcePermissionPolicyMembers(d.Get("members").(*schema.Set)),
	}
}

func buildSecurityResourcePermissionPolicyResources(resources *schema.Set) []map[string]interface{} {
	if resources.Len() == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, resources.Len())
	for _, v := range resources.List() {
		result = append(result, map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}

	return result
}

func buildSecurityResourcePermissionPolicyMembers(members *schema.Set) []map[string]interface{} {
	if members.Len() == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, members.Len())
	for _, v := range members.List() {
		result = append(result, map[string]interface{}{
			"member_id":   utils.PathSearch("member_id", v, nil),
			"member_name": utils.PathSearch("member_name", v, nil),
			"member_type": utils.PathSearch("member_type", v, nil),
		})
	}

	return result
}

func createSecurityResourcePermissionPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) (string, error) {
	httpUrl := "v1/{project_id}/security/permission-resource"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         buildSecurityResourcePermissionPolicyBodyParams(d),
	}

	resp, err := client.Request("POST", path, &opts)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	policyId := utils.PathSearch("policy_id", respBody, "").(string)
	if policyId == "" {
		return "", errors.New("unable to find the policy ID of the resource permission policy from the API response")
	}

	return policyId, nil
}

func resourceSecurityResourcePermissionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	policyId, err := createSecurityResourcePermissionPolicy(client, d, workspaceId)
	if err != nil {
		return diag.Errorf("error creating resource permission policy: %s", err)
	}

	d.SetId(policyId)

	return resourceSecurityResourcePermissionPolicyRead(ctx, d, meta)
}

// GetSecurityResourcePermissionPolicyById is a method used to query a resource permission policy by its ID.
func GetSecurityResourcePermissionPolicyById(client *golangsdk.ServiceClient, workspaceId, policyId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/security/permission-resource/{policy_id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{policy_id}", policyId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	resp, err := client.Request("GET", path, &opts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceSecurityResourcePermissionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		policyId    = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetSecurityResourcePermissionPolicyById(client, workspaceId, policyId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(common.ConvertExpected400ErrInto404Err(err, "errorCode",
				securityResourcePermissionPolicyResourceNotFoundCodes...), "error_code",
				securityResourcePermissionPolicyResourceNotFoundCodes...),
			fmt.Sprintf("error retrieving resource permission policy (%s)", policyId),
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("policy_name", respBody, nil)),
		d.Set("resources", flattenSecurityResourcePermissionPolicyResources(utils.PathSearch("resources", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("members", flattenSecurityResourcePermissionPolicyMembers(utils.PathSearch("members",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
			respBody, float64(0)).(float64))/1000, false)),
		d.Set("create_user", utils.PathSearch("create_user", respBody, nil)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
			respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityResourcePermissionPolicyResources(resources []interface{}) []interface{} {
	if len(resources) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		result = append(result, map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
		})
	}

	return result
}

func flattenSecurityResourcePermissionPolicyMembers(members []interface{}) []interface{} {
	if len(members) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(members))
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"member_id":   utils.PathSearch("member_id", member, nil),
			"member_name": utils.PathSearch("member_name", member, nil),
			"member_type": utils.PathSearch("member_type", member, nil),
		})
	}

	return result
}

func updateSecurityResourcePermissionPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) error {
	httpUrl := "v1/{project_id}/security/permission-resource/{policy_id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{policy_id}", d.Id())

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         buildSecurityResourcePermissionPolicyBodyParams(d),
	}

	_, err := client.Request("PUT", path, &opts)
	return err
}

func resourceSecurityResourcePermissionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err = updateSecurityResourcePermissionPolicy(client, d, workspaceId); err != nil {
		return diag.Errorf("error updating resource permission policy (%s): %s", d.Id(), err)
	}

	return resourceSecurityResourcePermissionPolicyRead(ctx, d, meta)
}

func deleteSecurityResourcePermissionPolicy(client *golangsdk.ServiceClient, workspaceId, policyId string) error {
	httpUrl := "v1/{project_id}/security/permission-resource/batch-delete"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody: map[string]interface{}{
			"ids": []string{policyId},
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", path, &opts)
	return err
}

func resourceSecurityResourcePermissionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		policyId    = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err = deleteSecurityResourcePermissionPolicy(client, workspaceId, policyId); err != nil {
		// DLS.6036: The Workspace ID is not found.
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", "DLS.6036"),
			fmt.Sprintf("error deleting resource permission policy (%s)", policyId),
		)
	}

	return nil
}

func resourceSecurityResourcePermissionPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
