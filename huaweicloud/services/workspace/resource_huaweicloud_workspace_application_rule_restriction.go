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

// @API Workspace POST /v1/{project_id}/app-center/app-restricted-rules
// @API Workspace POST /v1/{project_id}/app-center/app-restricted-rules/actions/batch-delete
// @API Workspace GET /v1/{project_id}/app-center/app-restricted-rules
func ResourceApplicationRuleRestriction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationRuleRestrictionCreate,
		ReadContext:   resourceApplicationRuleRestrictionRead,
		UpdateContext: resourceApplicationRuleRestrictionUpdate,
		DeleteContext: resourceApplicationRuleRestrictionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the restricted application rules are located.`,
			},
			"rule_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application rule IDs to be restricted.`,
			},
		},
	}
}

func buildRestrictedApplicationRuleBodyParams(ruleIds []string) map[string]interface{} {
	return map[string]interface{}{
		"items": ruleIds,
	}
}

func createRestrictedApplicationRule(client *golangsdk.ServiceClient, ruleIds []string) error {
	createPath := "v1/{project_id}/app-center/app-restricted-rules"
	createPath = client.Endpoint + createPath
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRestrictedApplicationRuleBodyParams(ruleIds),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceApplicationRuleRestrictionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		ruleIds = d.Get("rule_ids").(*schema.Set)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = createRestrictedApplicationRule(client, utils.ExpandToStringListBySet(ruleIds))
	if err != nil {
		return diag.Errorf("error creating Workspace restricted application rules: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceApplicationRuleRestrictionRead(ctx, d, meta)
}

func listRestrictedApplicationRuleIds(client *golangsdk.ServiceClient) ([]string, error) {
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

		ruleIds := utils.PathSearch("items[*].id", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, ruleIds...)
		if len(ruleIds) < limit {
			break
		}

		offset += len(ruleIds)
	}

	return utils.ExpandToStringList(results), nil
}

func FilterRestrictedApplicationRuleIds(client *golangsdk.ServiceClient, managedRuleIds []string) ([]string, error) {
	filterRuleIds := make([]string, 0)

	ruleIds, err := listRestrictedApplicationRuleIds(client)
	if err != nil {
		return nil, err
	}

	for _, ruleId := range ruleIds {
		if utils.StrSliceContains(managedRuleIds, ruleId) {
			filterRuleIds = append(filterRuleIds, ruleId)
		}
	}

	if len(filterRuleIds) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/app-center/app-restricted-rules",
				RequestId: "NONE",
				Body:      []byte("querying restricted application rules failed."),
			},
		}
	}
	return filterRuleIds, nil
}

func resourceApplicationRuleRestrictionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		managedRuleIds = d.Get("rule_ids").(*schema.Set)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	ruleIds, err := FilterRestrictedApplicationRuleIds(client, utils.ExpandToStringListBySet(managedRuleIds))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying Workspace restricted application rules")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rule_ids", ruleIds),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceApplicationRuleRestrictionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	if d.HasChange("rule_ids") {
		oldRaws, newRaws := d.GetChange("rule_ids")
		rmRuleIds := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
		addRuleIds := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))

		if addRuleIds.Len() > 0 {
			err = createRestrictedApplicationRule(client, utils.ExpandToStringListBySet(addRuleIds))
			if err != nil {
				return diag.Errorf("error adding new restricted application rules: %s", err)
			}
		}

		if rmRuleIds.Len() > 0 {
			err = deleteRestrictedApplicationRule(client, utils.ExpandToStringListBySet(rmRuleIds))
			if err != nil {
				return diag.Errorf("error deleting old restricted application rules: %s", err)
			}
		}
	}

	return resourceApplicationRuleRestrictionRead(ctx, d, meta)
}

func deleteRestrictedApplicationRule(client *golangsdk.ServiceClient, ruleIds []string) error {
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

func resourceApplicationRuleRestrictionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		ruleIds = d.Get("rule_ids").(*schema.Set)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	err = deleteRestrictedApplicationRule(client, utils.ExpandToStringListBySet(ruleIds))
	if err != nil {
		return diag.Errorf("error deleting Workspace restricted application rules: %s", err)
	}

	return nil
}
