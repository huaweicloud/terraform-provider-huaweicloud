package iam

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

var v5PolicyAgencyAttachNonUpdatableParams = []string{
	"policy_id",
	"agency_id",
}

// @API IAM POST /v5/policies/{policy_id}/attach-agency
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
// @API IAM POST /v5/policies/{policy_id}/detach-agency
func ResourceV5PolicyAgencyAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5PolicyAgencyAttachCreate,
		ReadContext:   resourceV5PolicyAgencyAttachRead,
		UpdateContext: resourceV5PolicyAgencyAttachUpdate,
		DeleteContext: resourceV5PolicyAgencyAttachDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5PolicyAgencyAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v5PolicyAgencyAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the IAM V5 policy.",
			},
			"agency_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the IAM agency.",
			},

			// Attributes.
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the IAM V5 policy.",
			},
			"policy_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URN of the IAM V5 policy.",
			},
			"attached_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the IAM V5 policy was attached to the agency.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func attachV5PolicyToAgency(client *golangsdk.ServiceClient, agencyId, policyId string) error {
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
	return err
}

func resourceV5PolicyAgencyAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Get("agency_id").(string)
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("identity", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = attachV5PolicyToAgency(client, agencyId, policyId)
	if err != nil {
		return diag.Errorf("error attaching IAM V5 policy (%s) to agency (%s): %s", policyId, agencyId, err)
	}

	d.SetId(policyId + "/" + agencyId)

	return resourceV5PolicyAgencyAttachRead(ctx, d, meta)
}

func GetV5AgencyAttachedPolicy(client *golangsdk.ServiceClient, agencyId, policyId string) (interface{}, error) {
	attachedPolicies, err := listV5AgencyAttachedPolicies(client, agencyId)
	if err != nil {
		return nil, err
	}

	attachedPolicy := utils.PathSearch(fmt.Sprintf("[?policy_id=='%s']|[0]", policyId), attachedPolicies, nil)
	if attachedPolicy == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/agencies/{agency_id}/attached-policies",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The attached policy (%s) of agency (%s) not found", policyId, agencyId)),
			},
		}
	}
	return attachedPolicy, nil
}

func resourceV5PolicyAgencyAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Get("agency_id").(string)
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("identity", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	attachedPolicy, err := GetV5AgencyAttachedPolicy(client, agencyId, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving attached policy of agency")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_name", utils.PathSearch("policy_name", attachedPolicy, nil)),
		d.Set("policy_urn", utils.PathSearch("urn", attachedPolicy, nil)),
		d.Set("attached_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("attached_at",
			attachedPolicy, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5PolicyAgencyAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func detachV5PolicyFromAgency(client *golangsdk.ServiceClient, agencyId, policyId string) error {
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
	return err
}

func resourceV5PolicyAgencyAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Get("agency_id").(string)
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("identity", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = detachV5PolicyFromAgency(client, agencyId, policyId)
	if err != nil {
		return diag.Errorf("error detaching IAM V5 policy (%s) from agency (%s): %s", policyId, agencyId, err)
	}

	return nil
}

func resourceV5PolicyAgencyAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, want '<policy_id>/<agency_id>', but got '%s'", importedId)
	}

	d.Set("policy_id", parts[0])
	d.Set("agency_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
