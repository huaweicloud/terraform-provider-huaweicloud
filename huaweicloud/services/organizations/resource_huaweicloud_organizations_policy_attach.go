package organizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations/policies/{policy_id}/detach
// @API Organizations POST /v1/organizations/policies/{policy_id}/attach
// @API Organizations GET /v1/organizations/policies/{policy_id}/attached-entities
func ResourcePolicyAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAttachCreate,
		ReadContext:   resourcePolicyAttachRead,
		DeleteContext: resourcePolicyAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the policy.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The unique ID of the root, OU, or account.`,
			},
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
		},
	}
}

func resourcePolicyAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/policies/{policy_id}/attach"
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
		JSONBody:         buildPolicyAttachBodyParams(entityId),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error attaching entity (%s) to policy (%s): %s", entityId, policyId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", policyId, entityId))

	return resourcePolicyAttachRead(ctx, d, meta)
}

func buildPolicyAttachBodyParams(entityId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"entity_id": entityId,
	}
	return bodyParams
}

func listPolicyAttachedEntities(client *golangsdk.ServiceClient, policyId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/policies/{policy_id}/attached-entities"
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

		entities := utils.PathSearch("attached_entities", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, entities...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func GetPolicyAttachedEntity(client *golangsdk.ServiceClient, policyId string, entityId string) (interface{}, error) {
	entities, err := listPolicyAttachedEntities(client, policyId)
	if err != nil {
		return nil, err
	}

	entity := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", entityId), entities, nil)
	if entity == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/policies/{policy_id}/attached-entities",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("unable to find the entity (%s) attached to the policy (%s)", entityId, policyId)),
			},
		}
	}

	return entity, nil
}

func resourcePolicyAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg  = meta.(*config.Config)
		mErr *multierror.Error
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<policy_id>/<entity_id>', but got '%s'", d.Id())
	}

	policyId := parts[0]
	entityId := parts[1]
	attachedEntity, err := GetPolicyAttachedEntity(client, policyId, entityId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("the entity (%s) is not attached to the policy (%s)", entityId, policyId))
	}

	mErr = multierror.Append(
		mErr,
		d.Set("policy_id", policyId),
		d.Set("entity_name", utils.PathSearch("name", attachedEntity, nil)),
		d.Set("entity_id", utils.PathSearch("id", attachedEntity, nil)),
		d.Set("entity_type", utils.PathSearch("type", attachedEntity, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePolicyAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/policies/{policy_id}/detach"
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
		JSONBody:         buildPolicyAttachBodyParams(entityId),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error detaching entity (%s) from policy (%s)", entityId, policyId),
		)
	}

	return nil
}
