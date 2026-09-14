package cbr

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var vaultPoliciesAssociateNonUpdatableParams = []string{
	"vault_id",
}

// @API CBR POST /v3/{project_id}/vaults/{vault_id}/associatepolicy
// @API CBR POST /v3/{project_id}/vaults/{vault_id}/dissociatepolicy
// @API CBR GET /v3/{project_id}/policies
func ResourceVaultPoliciesAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultPoliciesAssociateCreate,
		ReadContext:   resourceVaultPoliciesAssociateRead,
		UpdateContext: resourceVaultPoliciesAssociateUpdate,
		DeleteContext: resourceVaultPoliciesAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(vaultPoliciesAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to associate the policies to the vault.`,
			},

			// Required parameters.
			"vault_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the CBR vault to which the policies will be associated.`,
			},
			"policies": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the ID of the policy to associate with the vault.`,
						},
						"destination_vault_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the destination vault ID when associating a replication policy.`,
						},
					},
				},
				Description: `Specifies the policy details to associate with the CBR vault.`,
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

func resourceVaultPoliciesAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Get("vault_id").(string)
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	err = updatePoliciesBinding(client, vaultId, schema.NewSet(schema.HashString, nil), d.Get("policies"))
	if err != nil {
		return diag.Errorf("error associating policies to CBR vault (%s): %s", vaultId, err)
	}

	d.SetId(vaultId)

	// Avoid hitting cache nodes belonging to the same process to make sure the API response data is updated.
	// lintignore:R018
	time.Sleep(30 * time.Second)

	return resourceVaultPoliciesAssociateRead(ctx, d, meta)
}

func listPoliciesByVaultId(client *golangsdk.ServiceClient, vaultId string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/policies?vault_id={vault_id}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{vault_id}", vaultId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{}), nil
}

// GetVaultAssociatedPolicies queries the policies associated with the specified vault.
// Returns ErrDefault404 when the vault has no associated policies.
func GetVaultAssociatedPolicies(client *golangsdk.ServiceClient, vaultId string) ([]interface{}, error) {
	associatedPolicies, err := listPoliciesByVaultId(client, vaultId)
	if err != nil {
		return nil, err
	}

	if len(associatedPolicies) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/policies",
				RequestId: "NONE",
				Body:      fmt.Appendf(nil, "the vault (%s) has no associated policies", vaultId),
			},
		}
	}
	return associatedPolicies, nil
}

func flattenVaultAssociatedPolicies(policyList []interface{}, vaultId string) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(policyList))
	for _, val := range policyList {
		policy := map[string]interface{}{
			"id": utils.PathSearch("id", val, ""),
		}

		if destVaultId := utils.PathSearch(fmt.Sprintf("associated_vaults[?vault_id=='%s'].destination_vault_id|[0]",
			vaultId), val, "").(string); destVaultId != "" {
			policy["destination_vault_id"] = destVaultId
		}
		results = append(results, policy)
	}
	return results
}

func resourceVaultPoliciesAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Id()
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	policyList, err := GetVaultAssociatedPolicies(client, vaultId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying policies associated with CBR vault (%s)", vaultId))
	}
	if len(policyList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			fmt.Sprintf("error querying policies associated with CBR vault (%s)", vaultId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vault_id", vaultId),
		d.Set("policies", flattenVaultAssociatedPolicies(policyList, vaultId)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVaultPoliciesAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Id()
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	if d.HasChange("policies") {
		oldRaw, newRaw := d.GetChange("policies")
		err = updatePoliciesBinding(client, vaultId, oldRaw, newRaw)
		if err != nil {
			return diag.Errorf("error updating policies associated with CBR vault (%s): %s", vaultId, err)
		}
		// Avoid hitting cache nodes belonging to the same process to make sure the API response data is updated.
		// lintignore:R018
		time.Sleep(30 * time.Second)
	}

	return resourceVaultPoliciesAssociateRead(ctx, d, meta)
}

func resourceVaultPoliciesAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		vaultId = d.Id()
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	err = updatePoliciesBinding(client, vaultId, d.Get("policies"), schema.NewSet(schema.HashString, nil))
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error disassociating policies from CBR vault (%s)", vaultId))
	}

	return nil
}
