package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/agencies
// @API IAM GET /v5/agencies/{agency_id}
// @API IAM PUT /v5/agencies/{agency_id}
// @API IAM DELETE /v5/agencies/{agency_id}
// @API IAM GET /v5/policies
// @API IAM POST /v5/policies/{policy_id}/attach-agency
// @API IAM POST /v5/policies/{policy_id}/detach-agency
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
// @API IAM POST /v5/{resource_type}/{resource_id}/tags/create
// @API IAM POST /v5/{resource_type}/{resource_id}/tags/delete
func ResourceIAMTrustAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMTrustAgencyCreate,
		ReadContext:   resourceIAMTrustAgencyRead,
		UpdateContext: resourceIAMTrustAgencyUpdate,
		DeleteContext: resourceIAMTrustAgencyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trust_policy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"policy_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIAMTrustAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// get policy IDs by names
	policyIDs, err := getPolicyIDsByNames(client, d.Get("policy_names").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	createAgencyHttpUrl := "v5/agencies"
	createAgencyPath := client.Endpoint + createAgencyHttpUrl
	createAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateTrustAgencyBodyParams(d)),
	}
	createAgencyResp, err := client.Request("POST", createAgencyPath, &createAgencyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM trust agency: %s", err)
	}
	createAgencyRespBody, err := utils.FlattenResponse(createAgencyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("agency.agency_id", createAgencyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating IAM trust agency: agency_id is not found in API response")
	}
	d.SetId(id)

	// attach policies by ID
	for _, policyID := range policyIDs {
		if err = attachPolicyByID(client, d.Id(), policyID); err != nil {
			return diag.FromErr(err)
		}
	}

	// create tags
	if err := createTags(client, d.Get("tags").(map[string]interface{}), "agency", d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return resourceIAMTrustAgencyRead(ctx, d, meta)
}

func buildCreateTrustAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agency_name":          d.Get("name"),
		"trust_policy":         d.Get("trust_policy").(string),
		"path":                 utils.ValueIgnoreEmpty(d.Get("path")),
		"max_session_duration": utils.ValueIgnoreEmpty(d.Get("duration")),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceIAMTrustAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getAgencyHttpUrl := "v5/agencies/{agency_id}"
	getAgencyPath := client.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{agency_id}", d.Id())
	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var getAgencyResp *http.Response

	getAgencyResp, err = client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); !ok || !d.IsNewResource() {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM trust agency")
		}

		// if got 404 error in new resource, wait 10 seconds and try again
		// lintignore:R018
		time.Sleep(10 * time.Second)
		getAgencyResp, err = client.Request("GET", getAgencyPath, &getAgencyOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM trust agency")
		}
	}
	getAgencyRespBody, err := utils.FlattenResponse(getAgencyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	agency := utils.PathSearch("agency", getAgencyRespBody, nil)
	if agency == nil {
		return diag.Errorf("error getting IAM trust agency: agency is not found in API response")
	}

	// extract delegated service name
	var policy interface{}
	trustPolicy := utils.PathSearch("trust_policy", agency, nil)
	err = json.Unmarshal([]byte(trustPolicy.(string)), &policy)
	if err != nil {
		return diag.Errorf("error unmarshaling trust policy: %s", err)
	}

	// get attached policy names
	policyNames, err := listAttachedPolicyInfo(client, d.Id(), "policy_name")
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("agency_name", agency, nil)),
		d.Set("path", utils.PathSearch("path", agency, nil)),
		d.Set("duration", utils.PathSearch("max_session_duration", agency, nil)),
		d.Set("description", utils.PathSearch("description", agency, nil)),
		d.Set("tags", flattenTagsToMap(utils.PathSearch("tags", agency, make([]interface{}, 0)))),
		d.Set("urn", utils.PathSearch("urn", agency, nil)),
		d.Set("created_at", utils.PathSearch("created_at", agency, nil)),
		d.Set("trust_policy", utils.PathSearch("trust_policy", agency, nil)),
		d.Set("policy_names", policyNames),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM trust agency fields: %s", err)
	}

	return nil
}

func resourceIAMTrustAgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChanges("duration", "description") {
		updateAgencyHttpUrl := "v5/agencies/{agency_id}"
		updateAgencyPath := client.Endpoint + updateAgencyHttpUrl
		updateAgencyPath = strings.ReplaceAll(updateAgencyPath, "{agency_id}", d.Id())
		updateAgencyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"max_session_duration": d.Get("duration"),
				"description":          d.Get("description"),
			},
		}
		_, err := client.Request("PUT", updateAgencyPath, &updateAgencyOpt)
		if err != nil {
			return diag.Errorf("error updating IAM trust agency: %s", err)
		}
	}

	// update attached policy
	if d.HasChange("policy_names") {
		oRaw, nRaw := d.GetChange("policy_names")
		oMap := oRaw.(*schema.Set)
		nMap := nRaw.(*schema.Set)

		// list policy IDs which will be detached
		detachPolicyIDs, err := getPolicyIDsByNames(client, oMap.Difference(nMap).List())
		if err != nil {
			return diag.FromErr(err)
		}

		// detach policy by ID
		for _, detachPolicyID := range detachPolicyIDs {
			if err := detachPolicyByID(client, d.Id(), detachPolicyID); err != nil {
				return diag.FromErr(err)
			}
		}

		// list policy IDs which will be attached
		attachPolicyIDs, err := getPolicyIDsByNames(client, nMap.Difference(oMap).List())
		if err != nil {
			return diag.FromErr(err)
		}

		// attach policy by ID
		for _, attachPolicyID := range attachPolicyIDs {
			if err := attachPolicyByID(client, d.Id(), attachPolicyID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// update tags
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		// remove old tags
		if len(oMap) > 0 {
			if err = deleteTags(client, oMap, "agency", d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}

		// set new tags
		if len(nMap) > 0 {
			if err := createTags(client, nMap, "agency", d.Id()); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIAMTrustAgencyRead(ctx, d, meta)
}

func resourceIAMTrustAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// before deleting trust agency, require to detach policies first
	// get attached policy IDs
	policyIDs, err := listAttachedPolicyInfo(client, d.Id(), "policy_id")
	if err != nil {
		return diag.FromErr(err)
	}

	// detach policy by ID
	for _, policyID := range policyIDs {
		if err = detachPolicyByID(client, d.Id(), policyID); err != nil {
			return diag.FromErr(err)
		}
	}

	// delete trust agency
	deleteAgencyHttpUrl := "v5/agencies/{agency_id}"
	deleteAgencyPath := client.Endpoint + deleteAgencyHttpUrl
	deleteAgencyPath = strings.ReplaceAll(deleteAgencyPath, "{agency_id}", d.Id())
	deleteAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteAgencyPath, &deleteAgencyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM trust agency: %s", err)
	}

	return nil
}
