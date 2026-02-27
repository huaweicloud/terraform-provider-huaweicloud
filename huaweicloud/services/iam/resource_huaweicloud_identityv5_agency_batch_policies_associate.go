package iam

import (
	"context"
	"fmt"
	"log"
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

var (
	v5AgencyBatchPoliciesAssociateNonUpdatableParams = []string{
		"agency_id",
	}
	objSliceParamKeysForV5AgencyBatchPoliciesAssociate = []string{
		"policies",
	}
)

// @API IAM POST /v5/policies/{policy_id}/attach-agency
// @API IAM POST /v5/policies/{policy_id}/detach-agency
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
func ResourceV5AgencyBatchPoliciesAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5AgencyBatchPoliciesAssociateCreateOrUpdate,
		ReadContext:   resourceV5AgencyBatchPoliciesAssociateRead,
		UpdateContext: resourceV5AgencyBatchPoliciesAssociateCreateOrUpdate,
		DeleteContext: resourceV5AgencyBatchPoliciesAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(v5AgencyBatchPoliciesAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"agency_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM agency to which the policies will be attached.`,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the identity policy.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the identity policy.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URN of the identity policy.`,
						},
						"attached_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attached time of the identity policy, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of policies to be attached to the agency.`,
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

			// Internal attributes.
			"policies_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the identity policy.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'policies'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func orderAgencyAssociatedPoliciesByPoliciesOrigin(policies, policiesOrigin []interface{}) []interface{} {
	if len(policiesOrigin) < 1 {
		return policies
	}

	sortedPolicies := make([]interface{}, 0, len(policies))
	policiesCopy := policies
	for _, policyOrigin := range policiesOrigin {
		policyIdOrigin := utils.PathSearch("id", policyOrigin, "").(string)
		for index, policy := range policiesCopy {
			if utils.PathSearch("policy_id", policy, "").(string) == policyIdOrigin {
				// Add the found policy to the sorted policies list.
				sortedPolicies = append(sortedPolicies, policiesCopy[index])
				// Remove the processed policy from the original array.
				policiesCopy = append(policiesCopy[:index], policiesCopy[index+1:]...)
				break
			}
		}
	}

	return sortedPolicies
}

func ListV5AgencyAssociatedPolicies(client *golangsdk.ServiceClient, agencyId string, policiesOrigin []interface{},
	ignoreNotFound ...bool) ([]interface{}, error) {
	var (
		httpUrl = "v5/agencies/{agency_id}/attached-policies?limit={limit}"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{agency_id}", agencyId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		if marker != "" {
			listPath = fmt.Sprintf("%s&marker=%s", listPath, marker)
		}
		resp, err := client.Request("GET", listPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}
		policies := utils.PathSearch("attached_policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	parsedPolicies := orderAgencyAssociatedPoliciesByPoliciesOrigin(result, policiesOrigin)
	if len(parsedPolicies) < 1 && (len(ignoreNotFound) < 1 || !ignoreNotFound[0]) {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/agencies/{agency_id}/attached-policies",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("all policies were detached from the agency (%s)", agencyId)),
			},
		}
	}

	return parsedPolicies, nil
}

func findDeletePoliciesFromAgency(originPolicies, rawConfigPolicies []interface{}) []interface{} {
	if len(originPolicies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(originPolicies))
	for _, policy := range originPolicies {
		if utils.PathSearch(fmt.Sprintf("length([?id == '%v'])",
			utils.PathSearch("id", policy, "")), rawConfigPolicies, float64(0)).(float64) < 1 {
			// If the new policy list does not contain this policy, it is considered that this policy is no longer associated with the agency.
			result = append(result, policy)
		}
	}
	return result
}

func findAddPoliciesToAgency(rawConfigPolicies, remoteStatePolicies []interface{}) []interface{} {
	if len(rawConfigPolicies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(rawConfigPolicies))
	for _, policy := range rawConfigPolicies {
		if utils.PathSearch(fmt.Sprintf("length([?policy_id == '%v'])",
			utils.PathSearch("id", policy, "")), remoteStatePolicies, float64(0)).(float64) < 1 {
			// If the remote state policy list does not contain this policy, it is considered that this policy is a newly associated policy.
			result = append(result, policy)
		}
	}
	return result
}

func deletePoliciesFromAgency(client *golangsdk.ServiceClient, agencyId string, policies []interface{}) error {
	for _, policy := range policies {
		policyId := utils.PathSearch("id", policy, "").(string)
		httpUrl := "v5/policies/{policy_id}/detach-agency"
		detachPath := client.Endpoint + httpUrl
		detachPath = strings.ReplaceAll(detachPath, "{policy_id}", policyId)
		detachOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"agency_id": agencyId,
			},
		}
		_, err := client.Request("POST", detachPath, &detachOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the policy (%s) was already detached from the agency (%s)", policyId, agencyId)
				continue
			}
			return fmt.Errorf("error detaching IAM identity agency (%s) from policy (%s): %s", agencyId, policyId, err)
		}
	}
	return nil
}

func addPoliciesToAgency(client *golangsdk.ServiceClient, agencyId string, policies []interface{}) error {
	for _, policy := range policies {
		policyId := utils.PathSearch("id", policy, "").(string)

		httpUrl := "v5/policies/{policy_id}/attach-agency"
		attachPath := client.Endpoint + httpUrl
		attachPath = strings.ReplaceAll(attachPath, "{policy_id}", policyId)
		attachOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"agency_id": agencyId,
			},
		}
		_, err := client.Request("POST", attachPath, &attachOpt)
		if err != nil {
			return fmt.Errorf("error attaching IAM identity agency (%s) to policy(%s): %s", agencyId, policyId, err)
		}
	}
	return nil
}

func resourceV5AgencyBatchPoliciesAssociateCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		agencyId          = d.Get("agency_id").(string)
		rawConfigPolicies = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "policies").([]interface{})
		originPolicies    = d.Get("policies_origin").([]interface{})
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(agencyId)
	}

	remoteStatePolicies, err := ListV5AgencyAssociatedPolicies(client, agencyId, originPolicies, true)
	if err != nil {
		return diag.Errorf("error getting remote IAM agency associated policies: %s", err)
	}

	deletePolicies := findDeletePoliciesFromAgency(originPolicies, rawConfigPolicies)
	addPolicies := findAddPoliciesToAgency(rawConfigPolicies, remoteStatePolicies)

	if len(deletePolicies) > 0 {
		err = deletePoliciesFromAgency(client, agencyId, deletePolicies)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addPolicies) > 0 {
		err = addPoliciesToAgency(client, agencyId, addPolicies)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	// Only preserve the 'id' field in policies_origin, ignore other fields if they cause panic (no address).
	err = utils.RefreshObjectParamOriginValues(d, objSliceParamKeysForV5AgencyBatchPoliciesAssociate, utils.RefreshObjectParamOriginValuesOptions{
		PreservedFields: map[string][]string{
			"policies": {"id"},
		},
	})
	if err != nil {
		// Don't report an error if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceV5AgencyBatchPoliciesAssociateRead(ctx, d, meta)
}

func flattenV5AgencyAssociatedPolicies(policies []interface{}) []interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("policy_id", policy, nil),
			"name": utils.PathSearch("policy_name", policy, nil),
			"urn":  utils.PathSearch("urn", policy, nil),
			"attached_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("attached_at",
				policy, "").(string))/1000, false),
		})
	}
	return result
}

func resourceV5AgencyBatchPoliciesAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		agencyId       = d.Id()
		originPolicies = d.Get("policies_origin").([]interface{})
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	remoteStatePolicies, err := ListV5AgencyAssociatedPolicies(client, agencyId, originPolicies, false)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting remote IAM agency associated policies")
	}

	mErr := multierror.Append(nil,
		d.Set("agency_id", agencyId),
		d.Set("policies", flattenV5AgencyAssociatedPolicies(remoteStatePolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5AgencyBatchPoliciesAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Get("agency_id").(string)
		policies = d.Get("policies").([]interface{})
	)

	// The value of policies_origin is empty only when the resource is imported and the terraform apply command is not executed.
	// In this case, all information obtained from the remote service is used to remove policies from the agency.
	if originPolicies, ok := d.GetOk("policies_origin"); ok && len(originPolicies.([]interface{})) > 0 {
		log.Printf("[DEBUG] Find the custom policies configuration, according to it to remove policies from the agency (%v)", agencyId)
		policies = originPolicies.([]interface{})
	}

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = deletePoliciesFromAgency(client, agencyId, policies)
	return diag.FromErr(err)
}
