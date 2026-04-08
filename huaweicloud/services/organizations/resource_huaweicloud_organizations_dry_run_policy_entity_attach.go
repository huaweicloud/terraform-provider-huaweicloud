package organizations

import (
	"context"
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

var dryRunPolicyEntityAttachNonUpdatableParams = []string{"policy_id", "entity_id"}

// @API Organizations POST /v1/organizations/dry-run-policies/{policy_id}/attach
// @API Organizations GET /v1/organizations/dry-run-policies/{policy_id}/attached-entities
// @API Organizations POST /v1/organizations/dry-run-policies/{policy_id}/detach
func ResourceDryRunPolicyEntityAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDryRunPolicyEntityAttachCreate,
		ReadContext:   resourceDryRunPolicyEntityAttachRead,
		UpdateContext: resourceDryRunPolicyEntityAttachUpdate,
		DeleteContext: resourceDryRunPolicyEntityAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDryRunPolicyEntityAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(dryRunPolicyEntityAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dry-run policy.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the entity (root, OU, or account).`,
			},
			// Attributes.
			"entity_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the entity.`,
			},
			"entity_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the entity.`,
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

func resourceDryRunPolicyEntityAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/dry-run-policies/{policy_id}/attach"
		policyId = d.Get("policy_id").(string)
		entityId = d.Get("entity_id").(string)
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"entity_id": entityId,
		},
	}
	if _, err = client.Request("POST", createPath, &createOpt); err != nil {
		return diag.Errorf("error attaching entity (%s) to dry-run policy (%s): %s", entityId, policyId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", policyId, entityId))

	return resourceDryRunPolicyEntityAttachRead(ctx, d, meta)
}

func listAttachedEntitiesForDryRunPolicy(client *golangsdk.ServiceClient, policyId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/dry-run-policies/{policy_id}/attached-entities"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{policy_id}", policyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s?marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		attachedEntities := utils.PathSearch("attached_entities", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, attachedEntities...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

// GetAttachedEntityForDryRunPolicy is a method used to query the attached entity for a dry-run policy.
func GetAttachedEntityForDryRunPolicy(client *golangsdk.ServiceClient, policyId string, entityId string) (interface{}, error) {
	attachedEntities, err := listAttachedEntitiesForDryRunPolicy(client, policyId)
	if err != nil {
		return nil, err
	}

	attachedEntity := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", entityId), attachedEntities, nil)
	if attachedEntity != nil {
		return attachedEntity, nil
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/organizations/dry-run-policies/{policy_id}/attached-entities",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the attached entity (%s) for dry-run policy (%s) does not exist", entityId, policyId)),
		},
	}
}

func resourceDryRunPolicyEntityAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	policyId := d.Get("policy_id").(string)
	entityId := d.Get("entity_id").(string)
	attachedEntity, err := GetAttachedEntityForDryRunPolicy(client, policyId, entityId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error retrieving attached entity (%s) for dry-run policy (%s)", entityId, policyId),
		)
	}

	mErr := multierror.Append(
		d.Set("policy_id", policyId),
		d.Set("entity_id", utils.PathSearch("id", attachedEntity, nil)),
		d.Set("entity_name", utils.PathSearch("name", attachedEntity, nil)),
		d.Set("entity_type", utils.PathSearch("type", attachedEntity, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDryRunPolicyEntityAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDryRunPolicyEntityAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/dry-run-policies/{policy_id}/detach"
		policyId = d.Get("policy_id").(string)
		entityId = d.Get("entity_id").(string)
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"entity_id": entityId,
		},
	}
	if _, err = client.Request("POST", deletePath, &deleteOpt); err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error detaching entity (%s) from dry-run policy (%s)", entityId, policyId),
		)
	}

	return nil
}

func resourceDryRunPolicyEntityAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<entity_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("policy_id", parts[0]),
		d.Set("entity_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
