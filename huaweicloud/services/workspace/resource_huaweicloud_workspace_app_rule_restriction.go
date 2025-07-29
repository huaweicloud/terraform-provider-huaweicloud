package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WORKSPACE POST /v1/{project_id}/app-center/app-restricted-rules
// @API WORKSPACE POST /v1/{project_id}/app-center/app-restricted-rules/actions/batch-delete
// @API WORKSPACE GET /v1/{project_id}/app-center/app-restricted-rules
func ResourceAppRuleRestriction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppRuleRestrictionCreate,
		ReadContext:   resourceAppRuleRestrictionRead,
		UpdateContext: resourceAppRuleRestrictionUpdate,
		DeleteContext: resourceAppRuleRestrictionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application restricted rule is located.`,
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application rule IDs to be restricted.`,
			},
		},
	}
}

func buildAppRestrictedRuleBodyParams(ruleIds []string) map[string]interface{} {
	return map[string]interface{}{
		"items": ruleIds,
	}
}

func createAppRestrictedRule(client *golangsdk.ServiceClient, ruleIds []string) error {
	createPath := "v1/{project_id}/app-center/app-restricted-rules"
	createPath = client.Endpoint + createPath
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppRestrictedRuleBodyParams(ruleIds),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceAppRuleRestrictionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = createAppRestrictedRule(client, d.Get("rule_ids").([]string))
	if err != nil {
		return diag.Errorf("error creating Workspace application restricted rule: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAppRuleRestrictionRead(ctx, d, meta)
}

func convertRulesToIdList(rules []interface{}) []string {
	ruleIds := make([]string, 0, len(rules))
	for _, rule := range rules {
		if ruleMap, ok := rule.(map[string]interface{}); ok {
			if id, ok := ruleMap["id"].(string); ok {
				ruleIds = append(ruleIds, id)
			}
		}
	}
	return ruleIds
}

func ListAppRestrictedRuleIds(client *golangsdk.ServiceClient) ([]string, error) {
	var (
		listPathWithLimit = "v1/{project_id}/app-center/app-restricted-rules?limit={limit}"
		offset            = 0
		limit             = 100
		results           = make([]interface{}, 0)
	)

	listPathWithLimit = client.Endpoint + listPathWithLimit
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%s", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		rules := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, rules...)
		if len(rules) < limit {
			break
		}

		offset += len(rules)
	}

	if len(results) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/app-center/app-restricted-rules",
				RequestId: "NONE",
				Body:      []byte("querying application restricted rule failed."),
			},
		}
	}

	return convertRulesToIdList(results), nil
}

func resourceAppRuleRestrictionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	ruleIds, err := ListAppRestrictedRuleIds(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying Workspace application restricted rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rule_ids", ruleIds),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func compareRuleIds(currentRules, newRules []string) (rulesToAdd, rulesToDelete []string) {
	// Find rules to add (in new but not in current)
	rulesToAdd = make([]string, 0)
	for _, newRule := range newRules {
		found := false
		for _, currentRule := range currentRules {
			if newRule == currentRule {
				found = true
				break
			}
		}
		if !found {
			rulesToAdd = append(rulesToAdd, newRule)
		}
	}

	// Find rules to delete (in current but not in new)
	rulesToDelete = make([]string, 0)
	for _, currentRule := range currentRules {
		found := false
		for _, newRule := range newRules {
			if currentRule == newRule {
				found = true
				break
			}
		}
		if !found {
			rulesToDelete = append(rulesToDelete, currentRule)
		}
	}

	return rulesToAdd, rulesToDelete
}

func resourceAppRuleRestrictionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	// Update app_rule_ids if changed
	if d.HasChange("rule_ids") {
		// Get current restricted rules
		currentRuleIds, err := ListAppRestrictedRuleIds(client)
		if err != nil {
			return diag.Errorf("error querying current application restricted rules: %s", err)
		}

		// Get schema update rules
		newRuleIds := d.Get("rule_ids").([]string)

		// Compare rules to find what needs to be added or deleted
		rulesToAdd, rulesToDelete := compareRuleIds(
			currentRuleIds,
			newRuleIds)

		// Add new rules if any
		if len(rulesToAdd) > 0 {
			err = createAppRestrictedRule(client, rulesToAdd)
			if err != nil {
				return diag.Errorf("error adding new application restricted rules: %s", err)
			}
		}

		// Delete old rules if any
		if len(rulesToDelete) > 0 {
			err = deleteAppRestrictedRule(client, rulesToDelete)
			if err != nil {
				return diag.Errorf("error deleting old application restricted rules: %s", err)
			}
		}
	}

	return resourceAppRuleRestrictionRead(ctx, d, meta)
}

func deleteAppRestrictedRule(client *golangsdk.ServiceClient, ruleIds []string) error {
	deletePath := "v1/{project_id}/app-center/app-restricted-rules/actions/batch-delete"
	deletePath = client.Endpoint + deletePath
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"items": ruleIds,
		},
	}

	_, err := client.Request("POST", deletePath, &deleteOpt)
	return err
}

func resourceAppRuleRestrictionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = deleteAppRestrictedRule(client, d.Get("rule_ids").([]string))
	if err != nil {
		return diag.Errorf("error deleting Workspace application restricted rule: %s", err)
	}

	return resourceAppRuleRestrictionRead(ctx, d, meta)
}
