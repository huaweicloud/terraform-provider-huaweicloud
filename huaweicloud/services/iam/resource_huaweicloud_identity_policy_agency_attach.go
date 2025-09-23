package iam

import (
	"context"
	"errors"
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

var policyAgencyAttachNonUpdatableParams = []string{"policy_id", "agency_id"}

// @API IAM POST /v5/policies/{policy_id}/attach-agency
// @API IAM POST /v5/policies/{policy_id}/detach-agency
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
func ResourceIdentityPolicyAgencyAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityPolicyAgencyAttachCreate,
		ReadContext:   resourceIdentityPolicyAgencyAttachRead,
		UpdateContext: resourceIdentityPolicyAgencyAttachUpdate,
		DeleteContext: resourceIdentityPolicyAgencyAttachDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityPolicyAgencyAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(policyAgencyAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attached_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityPolicyAgencyAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	agencyID := d.Get("agency_id").(string)
	policyID := d.Get("policy_id").(string)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	attachPolicyAgencyHttpUrl := "v5/policies/{policy_id}/attach-agency"
	attachPolicyAgencyPath := client.Endpoint + attachPolicyAgencyHttpUrl
	attachPolicyAgencyPath = strings.ReplaceAll(attachPolicyAgencyPath, "{policy_id}", policyID)
	attachPolicyAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPolicyAgencyAttachBodyParams(d)),
	}
	_, err = client.Request("POST", attachPolicyAgencyPath, &attachPolicyAgencyOpt)
	if err != nil {
		return diag.Errorf("error attaching IAM identity agency(%s) with policy(%s): %s", agencyID, policyID, err)
	}

	d.SetId(policyID + "/" + agencyID)

	return resourceIdentityPolicyAgencyAttachRead(ctx, d, meta)
}

func buildPolicyAgencyAttachBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agency_id": d.Get("agency_id"),
	}
	return bodyParams
}

func resourceIdentityPolicyAgencyAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	agencyID := d.Get("agency_id").(string)
	policyID := d.Get("policy_id").(string)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getPolicyAgencyHttpUrl := "v5/agencies/{agency_id}/attached-policies"
	getPolicyAgencyPath := client.Endpoint + getPolicyAgencyHttpUrl
	getPolicyAgencyPath = strings.ReplaceAll(getPolicyAgencyPath, "{agency_id}", agencyID)

	queryParams := buildListPolicyAgencyAttachParams("")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", getPolicyAgencyPath+queryParams, &opt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving policies attached to agency")
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jsonPath := fmt.Sprintf("attached_policies[?policy_id=='%s'] | [0]", policyID)
		policy := utils.PathSearch(jsonPath, respBody, nil)
		if policy != nil {
			mErr := multierror.Append(nil,
				d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
				d.Set("policy_urn", utils.PathSearch("urn", policy, nil)),
				d.Set("attached_at", utils.PathSearch("attached_at", policy, nil)),
			)

			if err := mErr.ErrorOrNil(); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListPolicyAgencyAttachParams(marker)
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving policies attached to agency")
}

func buildListPolicyAgencyAttachParams(marker string) string {
	res := ""

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func resourceIdentityPolicyAgencyAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIdentityPolicyAgencyAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	agencyID := d.Get("agency_id").(string)
	policyID := d.Get("policy_id").(string)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	detachPolicyAgencyHttpUrl := "v5/policies/{policy_id}/detach-agency"
	detachPolicyAgencyPath := client.Endpoint + detachPolicyAgencyHttpUrl
	detachPolicyAgencyPath = strings.ReplaceAll(detachPolicyAgencyPath, "{policy_id}", policyID)
	detachPolicyAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPolicyAgencyAttachBodyParams(d)),
	}
	_, err = client.Request("POST", detachPolicyAgencyPath, &detachPolicyAgencyOpt)
	if err != nil {
		return diag.Errorf("error detaching IAM identity agency(%s) with policy(%s): %s", agencyID, policyID, err)
	}

	return nil
}

func resourceIdentityPolicyAgencyAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import id," +
			" must be <policy_id>/<agency_id>")
	}

	d.Set("policy_id", parts[0])
	d.Set("agency_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
