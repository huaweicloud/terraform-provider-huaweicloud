// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

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
				Description: `Specifies the ID of the policy.`,
			},
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the unique ID of the root, OU, or account.`,
			},
			"entity_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the entity.`,
			},
			"entity_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the entity.`,
			},
		},
	}
}

func resourcePolicyAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPolicyAttach: create Organizations policy attach
	var (
		createPolicyAttachHttpUrl = "v1/organizations/policies/{policy_id}/attach"
		createPolicyAttachProduct = "organizations"
	)
	createPolicyAttachClient, err := cfg.NewServiceClient(createPolicyAttachProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	policyId := d.Get("policy_id").(string)
	createPolicyAttachPath := createPolicyAttachClient.Endpoint + createPolicyAttachHttpUrl
	createPolicyAttachPath = strings.ReplaceAll(createPolicyAttachPath, "{policy_id}", policyId)

	createPolicyAttachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPolicyAttachOpt.JSONBody = utils.RemoveNil(buildPolicyAttachBodyParams(d))
	_, err = createPolicyAttachClient.Request("POST", createPolicyAttachPath, &createPolicyAttachOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations policy attach: %s", err)
	}

	entityId := d.Get("entity_id").(string)
	d.SetId(fmt.Sprintf("%s/%s", policyId, entityId))

	return resourcePolicyAttachRead(ctx, d, meta)
}

func buildPolicyAttachBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"entity_id": d.Get("entity_id"),
	}
	return bodyParams
}

func resourcePolicyAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPolicyAttach: Query Organizations policy attach
	var (
		getPolicyAttachHttpUrl = "v1/organizations/policies/{policy_id}/attached-entities"
		getPolicyAttachProduct = "organizations"
	)
	getPolicyAttachClient, err := cfg.NewServiceClient(getPolicyAttachProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	// Split policy_id and entity_id from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <policy_id>/<entity_id>")
	}
	policyId := parts[0]
	entityId := parts[1]

	getPolicyAttachPath := getPolicyAttachClient.Endpoint + getPolicyAttachHttpUrl
	getPolicyAttachPath = strings.ReplaceAll(getPolicyAttachPath, "{policy_id}", policyId)

	getPolicyAttachResp, err := pagination.ListAllItems(
		getPolicyAttachClient,
		"marker",
		getPolicyAttachPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations policy attach")
	}

	getPolicyAttachRespJson, err := json.Marshal(getPolicyAttachResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getPolicyAttachRespBody interface{}
	err = json.Unmarshal(getPolicyAttachRespJson, &getPolicyAttachRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	attachedEntity := utils.PathSearch(fmt.Sprintf("attached_entities[?id=='%s']|[0]", entityId),
		getPolicyAttachRespBody, nil)

	if attachedEntity == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePolicyAttach: Delete Organizations policy attach
	var (
		deletePolicyAttachHttpUrl = "v1/organizations/policies/{policy_id}/detach"
		deletePolicyAttachProduct = "organizations"
	)
	deletePolicyAttachClient, err := cfg.NewServiceClient(deletePolicyAttachProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePolicyAttachPath := deletePolicyAttachClient.Endpoint + deletePolicyAttachHttpUrl
	deletePolicyAttachPath = strings.ReplaceAll(deletePolicyAttachPath, "{policy_id}",
		d.Get("policy_id").(string))

	deletePolicyAttachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	deletePolicyAttachOpt.JSONBody = utils.RemoveNil(buildPolicyAttachBodyParams(d))
	_, err = deletePolicyAttachClient.Request("POST", deletePolicyAttachPath, &deletePolicyAttachOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations policy attach: %s", err)
	}

	return nil
}
