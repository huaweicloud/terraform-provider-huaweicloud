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

var v5PolicyUserAttachNonUpdatableParams = []string{"policy_id", "user_id"}

// @API IAM POST /v5/policies/{policy_id}/attach-user
// @API IAM POST /v5/policies/{policy_id}/detach-user
// @API IAM GET /v5/users/{user_id}/attached-policies
func ResourceV5PolicyUserAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5PolicyUserAttachCreate,
		ReadContext:   resourceV5PolicyUserAttachRead,
		UpdateContext: resourceV5PolicyUserAttachUpdate,
		DeleteContext: resourceV5PolicyUserAttachDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5PolicyUserAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v5PolicyUserAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the identity policy to be attached.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM user associated with the identity policy.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the identity policy associated with the user.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the attached identity policy.`,
			},
			"attached_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the identity policy was attached to the user, in RFC3339 format.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5PolicyUserAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		userId   = d.Get("user_id").(string)
		policyId = d.Get("policy_id").(string)
		httpUrl  = "v5/policies/{policy_id}/attach-user"
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"user_id": userId,
		},
	}
	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error attaching policy(%s) to the specified user(%s): %s", policyId, userId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", policyId, userId))

	return resourceV5PolicyUserAttachRead(ctx, d, meta)
}

func GetV5UserAttachedPolicy(client *golangsdk.ServiceClient, userId, policyId string) (interface{}, error) {
	httpUrl := "v5/users/{user_id}/attached-policies"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{user_id}", userId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	marker := ""
	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		policy := utils.PathSearch(fmt.Sprintf("attached_policies[?policy_id=='%s']|[0]", policyId), respBody, nil)
		if policy != nil {
			return policy, nil
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v5/users/{user_id}/attached-policies",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the policy (%s) associated with the user (%s) does not exist", policyId, userId)),
		},
	}
}

func resourceV5PolicyUserAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		userId   = d.Get("user_id").(string)
		policyId = d.Get("policy_id").(string)
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := GetV5UserAttachedPolicy(client, userId, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving policy associated with user")
	}

	mErr := multierror.Append(
		d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("urn", utils.PathSearch("urn", policy, nil)),
		d.Set("attached_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("attached_at",
			policy, "").(string))/1000, false)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceV5PolicyUserAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV5PolicyUserAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		userId   = d.Get("user_id").(string)
		policyId = d.Get("policy_id").(string)
		httpUrl  = "v5/policies/{policy_id}/detach-user"
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"user_id": userId,
		},
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error detaching policy(%s) from user(%s)", policyId, userId))
	}

	return nil
}

func resourceV5PolicyUserAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<user_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("policy_id", parts[0]),
		d.Set("user_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
