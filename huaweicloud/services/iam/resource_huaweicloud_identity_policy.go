package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

var policyNonUpdatableParams = []string{"name", "path", "description"}

// @API IAM POST /v5/policies
// @API IAM GET /v5/policies/{policy_id}
// @API IAM DELETE /v5/policies/{policy_id}
// @API IAM POST /v5/policies/{policy_id}/versions
// @API IAM GET /v5/policies/{policy_id}/versions
// @API IAM GET /v5/policies/{policy_id}/versions/{version_id}
// @API IAM DELETE /v5/policies/{policy_id}/versions/{version_id}
func ResourceIdentityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityPolicyCreate,
		ReadContext:   resourceIdentityPolicyRead,
		UpdateContext: resourceIdentityPolicyUpdate,
		DeleteContext: resourceIdentityPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(policyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_to_delete": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attachment_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPolicyHttpUrl := "v5/policies"
	createPolicyPath := client.Endpoint + createPolicyHttpUrl
	createPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePolicyBodyParams(d)),
	}
	createPolicyResp, err := client.Request("POST", createPolicyPath, &createPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM identity policy: %s", err)
	}
	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("policy.policy_id", createPolicyRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating IAM identity policy: policy_id is not found in API response")
	}
	d.SetId(id.(string))

	return resourceIdentityPolicyRead(ctx, d, meta)
}

func buildCreatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name":     d.Get("name"),
		"policy_document": d.Get("policy_document").(string),
		"path":            utils.ValueIgnoreEmpty(d.Get("path")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceIdentityPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getPolicyHttpUrl := "v5/policies/{policy_id}"
	getPolicyPath := client.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", d.Id())
	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var getPolicyResp *http.Response

	getPolicyResp, err = client.Request("GET", getPolicyPath, &getPolicyOpt)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); !ok || !d.IsNewResource() {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM identity policy")
		}

		// if got 404 error in new resource, wait 10 seconds and try again
		// lintignore:R018
		time.Sleep(10 * time.Second)
		getPolicyResp, err = client.Request("GET", getPolicyPath, &getPolicyOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM identity policy")
		}
	}
	getPolicyRespBody, err := utils.FlattenResponse(getPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := utils.PathSearch("policy", getPolicyRespBody, nil)
	if policy == nil {
		return diag.Errorf("error getting IAM identity policy: policy is not found in API response")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("path", utils.PathSearch("path", policy, nil)),
		d.Set("description", utils.PathSearch("description", policy, nil)),
		d.Set("urn", utils.PathSearch("urn", policy, nil)),
		d.Set("policy_type", utils.PathSearch("policy_type", policy, nil)),
		d.Set("default_version_id", utils.PathSearch("default_version_id", policy, nil)),
		d.Set("attachment_count", utils.PathSearch("attachment_count", policy, nil)),
		d.Set("created_at", utils.PathSearch("created_at", policy, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", policy, nil)),
	)

	defaultVersion, err := getIdentityPolicyVersionDocument(client, d, d.Get("default_version_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(mErr,
		d.Set("policy_document", utils.PathSearch("document", defaultVersion, nil)),
	)

	versions, err := getIdentityPolicyVersions(client, d)
	if err != nil {
		return diag.Errorf("error retrieving the versions of IAM identity policy(%s): %s", d.Id(), err)
	}

	mErr = multierror.Append(mErr,
		d.Set("version_ids", utils.PathSearch("[*].version_id", versions, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM identity policy fields: %s", err)
	}

	return nil
}

func getIdentityPolicyVersionDocument(client *golangsdk.ServiceClient, d *schema.ResourceData, versionID string) (interface{}, error) {
	getPolicyVersionHttpUrl := "v5/policies/{policy_id}/versions/{version_id}"
	getPolicyVersionPath := client.Endpoint + getPolicyVersionHttpUrl
	getPolicyVersionPath = strings.ReplaceAll(getPolicyVersionPath, "{policy_id}", d.Id())
	getPolicyVersionPath = strings.ReplaceAll(getPolicyVersionPath, "{version_id}", versionID)
	getPolicyVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyVersionResp, err := client.Request("GET", getPolicyVersionPath, &getPolicyVersionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM identity policy versions: %s", err)
	}
	getPolicyVersionRespBody, err := utils.FlattenResponse(getPolicyVersionResp)
	if err != nil {
		return nil, err
	}

	policyVersion := utils.PathSearch("policy_version", getPolicyVersionRespBody, nil)
	if policyVersion == nil {
		return nil, fmt.Errorf("error getting IAM identity policy versions: policy_version is not found in API response")
	}

	return policyVersion, nil
}

func getIdentityPolicyVersions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	getPolicyVersionsHttpUrl := "v5/policies/{policy_id}/versions"
	getPolicyVersionsPath := client.Endpoint + getPolicyVersionsHttpUrl
	getPolicyVersionsPath = strings.ReplaceAll(getPolicyVersionsPath, "{policy_id}", d.Id())
	getPolicyVersionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var policyVersions []interface{}
	for {
		getPolicyVersionsResp, err := client.Request("GET", getPolicyVersionsPath, &getPolicyVersionsOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving IAM identity policy versions: %s", err)
		}
		getPolicyVersionsRespBody, err := utils.FlattenResponse(getPolicyVersionsResp)
		if err != nil {
			return nil, err
		}

		policyVersionsTemp := utils.PathSearch("versions", getPolicyVersionsRespBody, nil)
		if policyVersionsTemp == nil {
			return nil, fmt.Errorf("error getting IAM identity policy versions: versions is not found in API response")
		}

		policyVersions = append(policyVersions, policyVersionsTemp.([]interface{})...)

		marker := utils.PathSearch("page_info.next_marker", getPolicyVersionsRespBody, nil)
		if marker == nil {
			break
		}
		getPolicyVersionsPath += fmt.Sprintf("?marker=%s", marker)
	}

	return policyVersions, nil
}

func removeIdentityPolicyVersion(client *golangsdk.ServiceClient, d *schema.ResourceData, versionID string) error {
	deletePolicyVersionHttpUrl := "v5/policies/{policy_id}/versions/{version_id}"
	deletePolicyVersionPath := client.Endpoint + deletePolicyVersionHttpUrl
	deletePolicyVersionPath = strings.ReplaceAll(deletePolicyVersionPath, "{policy_id}", d.Id())
	deletePolicyVersionPath = strings.ReplaceAll(deletePolicyVersionPath, "{version_id}", versionID)
	deletePolicyVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePolicyVersionPath, &deletePolicyVersionOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildAddPolicyVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_document": d.Get("policy_document").(string),
		"set_as_default":  true,
	}
	return bodyParams
}

func addIdentityPolicyVersion(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	addPolicyVersionHttpUrl := "v5/policies/{policy_id}/versions"
	addPolicyVersionPath := client.Endpoint + addPolicyVersionHttpUrl
	addPolicyVersionPath = strings.ReplaceAll(addPolicyVersionPath, "{policy_id}", d.Id())
	addPolicyVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAddPolicyVersionBodyParams(d)),
	}

	_, err := client.Request("POST", addPolicyVersionPath, &addPolicyVersionOpt)
	if err != nil {
		return err
	}

	return nil
}

func handleAddVersionError409(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}

		errCode := utils.PathSearch("error_code", apiError, nil)
		if apiError == nil {
			return false, err
		}

		// PAP5.0028: versions per policy limit exceeded
		if errCode == "PAP5.0028" {
			return true, err
		}
	}
	return false, err
}

func getEarliestVersionID(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("version_to_delete"); ok {
		return v.(string), nil
	}
	versions, err := getIdentityPolicyVersions(client, d)
	if err != nil {
		return "", fmt.Errorf("error retrieving the versions of IAM identity policy(%s): %s", d.Id(), err)
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("the list of policy versions is empty")
	}

	return utils.PathSearch("[-1].version_id", versions, "").(string), nil
}

func resourceIdentityPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	policyID := d.Id()

	err = addIdentityPolicyVersion(client, d)
	if err != nil {
		removeVersion, err := handleAddVersionError409(err)
		// if get a "versions per policy limit exceeded" error
		// remove a version and try again
		if !removeVersion {
			return diag.Errorf("error adding a new version of IAM identity policy(%s): %s", policyID, err)
		}
		versionID, err := getEarliestVersionID(client, d)
		if err != nil {
			return diag.FromErr(err)
		}

		err = removeIdentityPolicyVersion(client, d, versionID)
		if err != nil {
			return diag.Errorf("error removing the earliest version of IAM identity policy(%s): %s", policyID, err)
		}

		err = addIdentityPolicyVersion(client, d)
		if err != nil {
			return diag.Errorf("error adding a new version of IAM identity policy(%s): %s", policyID, err)
		}
	}

	return resourceIdentityPolicyRead(ctx, d, meta)
}

func resourceIdentityPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePolicyHttpUrl := "v5/policies/{policy_id}"
	deletePolicyPath := client.Endpoint + deletePolicyHttpUrl
	deletePolicyPath = strings.ReplaceAll(deletePolicyPath, "{policy_id}", d.Id())
	deletePolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePolicyPath, &deletePolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IAM identity policy")
	}

	return nil
}
