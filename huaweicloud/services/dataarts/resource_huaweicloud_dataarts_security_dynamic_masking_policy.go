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
	dynamicMaskingPolicyResourceNotFoundCodes = []string{
		"DLS.6036", // Workspace does not exist, error code key is "errorCode".
		"DLS.2040", // The resource does not exist, error code key is "error_code".
	}
	dynamicMaskingPolicyResourceDeleteNotFoundCodes = []string{
		"DLS.6036", // Workspace does not exist, error code key is "errorCode".
		"DLS.2052", // The resource does not exist, error code key is "error_code".
	}
	dynamicMaskingPolicyNonUpdatableParams = []string{
		"workspace_id",
		"datasource_type",
		"cluster_id",
		"cluster_name",
		"database_name",
		"table_name",
		"schema_name",
	}
)

// @API DataArtsStudio POST /v1/{project_id}/security/masking/dynamic/policies
// @API DataArtsStudio GET /v1/{project_id}/security/masking/dynamic/policies/{id}
// @API DataArtsStudio PUT /v1/{project_id}/security/masking/dynamic/policies/{id}
// @API DataArtsStudio POST /v1/{project_id}/security/masking/dynamic/policies/batch-delete
func ResourceSecurityDynamicMaskingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityDynamicMaskingPolicyCreate,
		ReadContext:   resourceSecurityDynamicMaskingPolicyRead,
		UpdateContext: resourceSecurityDynamicMaskingPolicyUpdate,
		DeleteContext: resourceSecurityDynamicMaskingPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(dynamicMaskingPolicyNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityDynamicMaskingPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dynamic masking policy is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the dynamic masking policy belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the dynamic masking policy.`,
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data source type of the dynamic masking policy.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster corresponding to the data source.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the cluster corresponding to the data source.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the data table.`,
			},
			"conn_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the data connection.`,
			},
			"conn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the data connection.`,
			},
			"policy_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `The list of dynamic masking policy configurations.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The field name in the data table.`,
						},
						"column_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The field type in the data table.`,
						},
						"algorithm_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The algorithm type of dynamic masking.`,
						},
						"algorithm_detail_dto": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The algorithm detail object of dynamic masking, in JSON format.`,
						},
					},
				},
			},

			// Optional parameters.
			"table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the data table.`,
			},
			"user_groups": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The list of user groups, separated by commas (,).`,
			},
			"users": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The list of user names, separated by commas (,).`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The schema name corresponding to the DWS data source.`,
			},

			// Attributes.
			"sync_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current synchronization status of the policy.`,
			},
			"sync_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronization message of the policy.`,
			},
			"sync_log": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronization log of the policy.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the policy, in RFC3339 format.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the policy.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the policy, in RFC3339 format.`,
			},
			"update_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest updater of the policy.`,
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
func createSecurityDynamicMaskingPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) (string, error) {
	httpUrl := "v1/{project_id}/security/masking/dynamic/policies"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildSecurityDynamicMaskingPolicyBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opts)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	policyId := utils.PathSearch("id", respBody, "").(string)
	if policyId == "" {
		return "", errors.New("unable to find the ID of the dynamic masking policy from the API response")
	}

	return policyId, nil
}

func resourceSecurityDynamicMaskingPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	policyId, err := createSecurityDynamicMaskingPolicy(client, d, d.Get("workspace_id").(string))
	if err != nil {
		return diag.Errorf("error creating dynamic masking policy: %s", err)
	}

	d.SetId(policyId)

	return resourceSecurityDynamicMaskingPolicyRead(ctx, d, meta)
}

func buildSecurityDynamicMaskingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":            d.Get("name"),
		"datasource_type": d.Get("datasource_type"),
		"cluster_id":      d.Get("cluster_id"),
		"cluster_name":    d.Get("cluster_name"),
		"database_name":   d.Get("database_name"),
		"table_name":      d.Get("table_name"),
		"conn_name":       d.Get("conn_name"),
		"conn_id":         d.Get("conn_id"),
		"policy_list":     buildSecurityDynamicMaskingPolicies(d.Get("policy_list").(*schema.Set)),
		"table_id":        utils.ValueIgnoreEmpty(d.Get("table_id")),
		"user_groups":     utils.ValueIgnoreEmpty(d.Get("user_groups")),
		"users":           utils.ValueIgnoreEmpty(d.Get("users")),
		"schema_name":     utils.ValueIgnoreEmpty(d.Get("schema_name")),
	}
}

func buildSecurityDynamicMaskingPolicies(policies *schema.Set) []map[string]interface{} {
	if policies.Len() == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, policies.Len())
	for _, policy := range policies.List() {
		result = append(result, map[string]interface{}{
			"column_name":          utils.PathSearch("column_name", policy, nil),
			"column_type":          utils.PathSearch("column_type", policy, nil),
			"algorithm_type":       utils.PathSearch("algorithm_type", policy, nil),
			"algorithm_detail_dto": utils.StringToJson(utils.PathSearch("algorithm_detail_dto", policy, "").(string)),
		})
	}

	return result
}

func GetSecurityDynamicMaskingPolicyById(client *golangsdk.ServiceClient, workspaceId, policyId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/security/masking/dynamic/policies/{id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", policyId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	resp, err := client.Request("GET", getPath, &opts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceSecurityDynamicMaskingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetSecurityDynamicMaskingPolicyById(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(common.ConvertExpected400ErrInto404Err(err, "errorCode",
				dynamicMaskingPolicyResourceNotFoundCodes...), "error_code", dynamicMaskingPolicyResourceNotFoundCodes...),
			fmt.Sprintf("error retrieving dynamic masking policy (%s)", d.Id()),
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("datasource_type", utils.PathSearch("datasource_type", respBody, nil)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", respBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", respBody, nil)),
		d.Set("database_name", utils.PathSearch("database_name", respBody, nil)),
		d.Set("table_name", utils.PathSearch("table_name", respBody, nil)),
		d.Set("conn_name", utils.PathSearch("conn_name", respBody, nil)),
		d.Set("conn_id", utils.PathSearch("conn_id", respBody, nil)),
		d.Set("policy_list", flattenSecurityDynamicMaskingPolicyList(utils.PathSearch("policy_list", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("table_id", utils.PathSearch("table_id", respBody, nil)),
		d.Set("user_groups", utils.PathSearch("user_groups", respBody, nil)),
		d.Set("users", utils.PathSearch("users", respBody, nil)),
		d.Set("schema_name", utils.PathSearch("schema_name", respBody, nil)),
		d.Set("sync_status", utils.PathSearch("sync_status", respBody, nil)),
		d.Set("sync_msg", utils.PathSearch("sync_msg", respBody, nil)),
		d.Set("sync_log", utils.PathSearch("sync_log", respBody, nil)),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
			respBody, float64(0)).(float64))/1000, false)),
		d.Set("create_user", utils.PathSearch("create_user", respBody, nil)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
			respBody, float64(0)).(float64))/1000, false)),
		d.Set("update_user", utils.PathSearch("update_user", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityDynamicMaskingPolicyList(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"column_name":    utils.PathSearch("column_name", policy, nil),
			"column_type":    utils.PathSearch("column_type", policy, nil),
			"algorithm_type": utils.PathSearch("algorithm_type", policy, nil),
			// 'algorithm_detail' and 'algorithm_detail_dto' only differ in type, but they represent the same content.
			// 'algorithm_detail' is a JSON string without null values.
			// 'algorithm_detail_dto' is a struct object that may contain null values.
			"algorithm_detail_dto": utils.PathSearch("algorithm_detail", policy, nil),
		})
	}

	return result
}

func updateSecurityDynamicMaskingPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) error {
	httpUrl := "v1/{project_id}/security/masking/dynamic/policies/{id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildSecurityDynamicMaskingPolicyBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opts)
	return err
}

func resourceSecurityDynamicMaskingPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = updateSecurityDynamicMaskingPolicy(client, d, d.Get("workspace_id").(string))
	if err != nil {
		return diag.Errorf("error updating dynamic masking policy (%s): %s", d.Id(), err)
	}

	return resourceSecurityDynamicMaskingPolicyRead(ctx, d, meta)
}

func deleteSecurityDynamicMaskingPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, workspaceId string) error {
	httpUrl := "v1/{project_id}/security/masking/dynamic/policies/batch-delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody: map[string]interface{}{
			"ids": []string{d.Id()},
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", deletePath, &opts)
	return err
}

func resourceSecurityDynamicMaskingPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = deleteSecurityDynamicMaskingPolicy(client, d, workspaceId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(common.ConvertExpected400ErrInto404Err(err, "errorCode",
				dynamicMaskingPolicyResourceDeleteNotFoundCodes...), "error_code", dynamicMaskingPolicyResourceDeleteNotFoundCodes...),
			fmt.Sprintf("error deleting dynamic masking policy (%s)", d.Id()),
		)
	}

	return nil
}
func resourceSecurityDynamicMaskingPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
