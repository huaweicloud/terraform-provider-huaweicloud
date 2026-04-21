package dataarts

import (
	"context"
	"fmt"
	"strconv"
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

var securityDataRecognitionRuleGroupNonUpdatableParams = []string{
	"workspace_id",
}

// @API DataArtsStudio POST /v1/{project_id}/security/data-classification/rule/group
// @API DataArtsStudio GET /v1/{project_id}/security/data-classification/rule/group
// @API DataArtsStudio PUT /v1/{project_id}/security/data-classification/rule/group/{id}
// @API DataArtsStudio POST /v1/{project_id}/security/data-classification/rule/group/batch-delete
func ResourceSecurityDataRecognitionRuleGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityDataRecognitionRuleGroupCreate,
		ReadContext:   resourceSecurityDataRecognitionRuleGroupRead,
		UpdateContext: resourceSecurityDataRecognitionRuleGroupUpdate,
		DeleteContext: resourceSecurityDataRecognitionRuleGroupDelete,

		CustomizeDiff: config.FlexibleForceNew(securityDataRecognitionRuleGroupNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityDataRecognitionRuleGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the data recognition rule group is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the rule group belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the data recognition rule group.`,
			},
			"rule_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of data recognition rule IDs that the rule group contains.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the data recognition rule group.`,
			},

			// Attributes.
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the data recognition rule group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the data recognition rule group, in RFC3339 format.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the data recognition rule group.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the data recognition rule group, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildSecurityDataRecognitionRuleGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"rule_ids":    d.Get("rule_ids").(*schema.Set).List(),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func createSecurityDataRecognitionRuleGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/security/data-classification/rule/group"
		workspaceId = d.Get("workspace_id").(string)
	)

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildSecurityDataRecognitionRuleGroupBodyParams(d)),
	}

	resp, err := client.Request("POST", path, &opts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceSecurityDataRecognitionRuleGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := createSecurityDataRecognitionRuleGroup(client, d)
	if err != nil {
		return diag.Errorf("error creating DataArts Security data recognition rule group: %s", err)
	}

	groupId := utils.PathSearch("uuid", respBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable to find the ID of the DataArts Security data recognition rule group from the API response")
	}
	d.SetId(groupId)

	return resourceSecurityDataRecognitionRuleGroupRead(ctx, d, meta)
}

func listSecurityDataRecognitionRuleGroups(client *golangsdk.ServiceClient, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/data-classification/rule/group?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opts)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		ruleGroups := utils.PathSearch("rule_groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, ruleGroups...)
		if len(ruleGroups) < limit {
			break
		}
		offset += len(ruleGroups)
	}
	return result, nil
}

func GetSecurityDataRecognitionRuleGroupById(client *golangsdk.ServiceClient, workspaceId, groupId string) (interface{}, error) {
	groups, err := listSecurityDataRecognitionRuleGroups(client, workspaceId)
	if err != nil {
		return nil, err
	}

	group := utils.PathSearch(fmt.Sprintf("[?uuid=='%s']|[0]", groupId), groups, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/security/data-classification/rule/group",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the data recognition rule group (%s) does not exist", groupId)),
			},
		}
	}
	return group, nil
}

func resourceSecurityDataRecognitionRuleGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		groupId     = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetSecurityDataRecognitionRuleGroupById(client, workspaceId, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving DataArts Security data recognition rule group (%s)", groupId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("rule_ids", utils.PathSearch("rules[*].uuid", respBody, make([]interface{}, 0)).([]interface{})),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("created_by", utils.PathSearch("created_by", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("created_at", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_by", utils.PathSearch("updated_by", respBody, nil)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("updated_at", respBody, float64(0)).(float64))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func updateSecurityDataRecognitionRuleGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/security/data-classification/rule/group/{id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{id}", d.Id())

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody: map[string]interface{}{
			"name":        d.Get("name"),
			"rule_ids":    d.Get("rule_ids").(*schema.Set).List(),
			"description": d.Get("description"),
		},
	}
	_, err := client.Request("PUT", path, &opts)
	return err
}

func resourceSecurityDataRecognitionRuleGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if d.HasChanges("name", "rule_ids", "description") {
		err = updateSecurityDataRecognitionRuleGroup(client, d)
		if err != nil {
			return diag.Errorf("error updating DataArts Security data recognition rule group (%s): %s", d.Id(), err)
		}
	}

	return resourceSecurityDataRecognitionRuleGroupRead(ctx, d, meta)
}

func deleteSecurityDataRecognitionRuleGroup(client *golangsdk.ServiceClient, workspaceId, groupId string) error {
	httpUrl := "v1/{project_id}/security/data-classification/rule/group/batch-delete"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody: map[string]interface{}{
			"rule_group_ids": []string{groupId},
		},
		OkCodes: []int{204},
	}
	_, err := client.Request("POST", path, &opts)
	return err
}

func resourceSecurityDataRecognitionRuleGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		groupId     = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if err := deleteSecurityDataRecognitionRuleGroup(client, workspaceId, groupId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DataArts Security data recognition rule group (%s)", groupId))
	}

	return nil
}

func resourceSecurityDataRecognitionRuleGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
