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

var v5PolicyGroupAttachNonUpdatableParams = []string{"policy_id", "group_id"}

// @API IAM POST /v5/policies/{policy_id}/attach-group
// @API IAM POST /v5/policies/{policy_id}/detach-group
// @API IAM GET /v5/groups/{group_id}/attached-policies
func ResourceV5PolicyGroupAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5PolicyGroupAttachCreate,
		ReadContext:   resourceV5PolicyGroupAttachRead,
		UpdateContext: resourceV5PolicyGroupAttachUpdate,
		DeleteContext: resourceV5PolicyGroupAttachDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5PolicyGroupAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v5PolicyGroupAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the identity policy to be attached.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM user group associated with the identity policy.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the identity policy associated with the user group.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the attached identity policy.`,
			},
			"attached_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the identity policy was attached to the user group, in RFC3339 format.`,
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

func resourceV5PolicyGroupAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		groupId  = d.Get("group_id").(string)
		policyId = d.Get("policy_id").(string)
		httpUrl  = "v5/policies/{policy_id}/attach-group"
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
			"group_id": groupId,
		},
	}
	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error attaching policy(%s) to the specified user group(%s): %s", policyId, groupId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", policyId, groupId))

	return resourceV5PolicyGroupAttachRead(ctx, d, meta)
}

func GetV5GroupAttachedPolicy(client *golangsdk.ServiceClient, groupId, policyId string) (interface{}, error) {
	httpUrl := "v5/groups/{group_id}/attached-policies"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{group_id}", groupId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	marker := ""
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
			URL:       "/v5/groups/{group_id}/attached-policies",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the policy (%s) associated with the group (%s) does not exist", policyId, groupId)),
		},
	}
}

func resourceV5PolicyGroupAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		groupId  = d.Get("group_id").(string)
		policyId = d.Get("policy_id").(string)
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := GetV5GroupAttachedPolicy(client, groupId, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving policy associated with group")
	}

	mErr := multierror.Append(
		d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("urn", utils.PathSearch("urn", policy, nil)),
		d.Set("attached_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("attached_at",
			policy, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV5PolicyGroupAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV5PolicyGroupAttachDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		groupId  = d.Get("group_id").(string)
		policyId = d.Get("policy_id").(string)
		httpUrl  = "v5/policies/{policy_id}/detach-group"
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
			"group_id": groupId,
		},
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error detaching policy(%s) from group(%s)", policyId, groupId))
	}

	return nil
}

func resourceV5PolicyGroupAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<group_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("policy_id", parts[0]),
		d.Set("group_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
