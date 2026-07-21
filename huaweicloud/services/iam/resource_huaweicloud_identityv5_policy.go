package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v5PolicyNonUpdatableParams = []string{
	"name",
	"path",
	"description",
}

// @API IAM POST /v5/policies
// @API IAM GET /v5/policies/{policy_id}
// @API IAM GET /v5/policies/{policy_id}/versions
// @API IAM POST /v5/policies/{policy_id}/versions
// @API IAM DELETE /v5/policies/{policy_id}/versions/{version_id}
// @API IAM DELETE /v5/policies/{policy_id}
func ResourceV5Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5PolicyCreate,
		ReadContext:   resourceV5PolicyRead,
		UpdateContext: resourceV5PolicyUpdate,
		DeleteContext: resourceV5PolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(v5PolicyNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the identity policy.`,
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(o, n)
					return equal
				},
				Description: `The policy document of the identity policy, in JSON format.`,
			},

			// Optional parameters.
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource path of the identity policy.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the identity policy.`,
			},
			"version_to_delete": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the policy version to be deleted.`,
			},

			// Attributes.
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the identity policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the identity policy.`,
			},
			"default_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default version ID of the identity policy.`,
			},
			"version_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The version ID list of the identity policy.`,
			},
			"attachment_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The attachment count of the identity policy.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the identity policy.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the identity policy.`,
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

			// Deprecated attributes.
			"policy_type": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: `Use attribute 'type' instead.`,
			},
		},
	}
}

func buildCreateV5PolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name":     d.Get("name"),
		"policy_document": d.Get("policy_document").(string),
		"path":            utils.ValueIgnoreEmpty(d.Get("path")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceV5PolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/policies"
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateV5PolicyBodyParams(d)),
	}

	createPolicyResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating identity policy: %s", err)
	}
	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy.policy_id", createPolicyRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the policy ID from the API response")
	}
	d.SetId(policyId)

	return resourceV5PolicyRead(ctx, d, meta)
}

func handlePolicyQueryError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		return true, err
	}
	if _, ok := err.(golangsdk.ErrDefault409); ok {
		return true, err
	}
	return false, err
}

// isRetry is a optional parameter (defaults to false), if it is true, the function will retry the request when the
// error is retryable (404 or 409).
func GetV5PolicyById(ctx context.Context, client *golangsdk.ServiceClient, policyId string, timeout time.Duration,
	isRetry ...bool) (interface{}, error) {
	httpUrl := "v5/policies/{policy_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{policy_id}", policyId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var (
		requestResp *http.Response
		requestErr  error
	)

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		requestResp, requestErr = client.Request("GET", getPath, &getOpt)
		retryable, err := handlePolicyQueryError(requestErr)
		if retryable && (len(isRetry) > 0 && isRetry[0]) {
			// lintignore:R018
			time.Sleep(15 * time.Second)
			return retry.RetryableError(err)
		}
		if err != nil {
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	requestRespBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("policy", requestRespBody, nil), nil
}

func listV5PolicyVersions(client *golangsdk.ServiceClient, policyId string) ([]interface{}, error) {
	var (
		httpUrl = "v5/policies/{policy_id}/versions?limit={limit}"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{policy_id}", policyId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		versions := utils.PathSearch("versions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, versions...)
		if len(versions) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}
	return result, nil
}

func resourceV5PolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy, err := GetV5PolicyById(ctx, client, policyId, d.Timeout(schema.TimeoutRead), d.IsNewResource())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving identity policy (%s)", policyId))
	}

	versions, err := listV5PolicyVersions(client, policyId)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		// Required parameters.
		d.Set("name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("policy_document", utils.PathSearch("[?is_default]|[0].document", versions, nil)),
		// Optional parameters.
		d.Set("path", utils.PathSearch("path", policy, nil)),
		d.Set("description", utils.PathSearch("description", policy, nil)),
		// Attributes.
		d.Set("urn", utils.PathSearch("urn", policy, nil)),
		d.Set("type", utils.PathSearch("policy_type", policy, nil)),
		d.Set("default_version_id", utils.PathSearch("default_version_id", policy, nil)),
		d.Set("version_ids", utils.PathSearch("[*].version_id", versions, nil)),
		d.Set("attachment_count", utils.PathSearch("attachment_count", policy, nil)),
		d.Set("created_at", utils.PathSearch("created_at", policy, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", policy, nil)),
		// Deprecated attributes.
		d.Set("policy_type", utils.PathSearch("policy_type", policy, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving identity policy (%s) fields: %s", policyId, err)
	}
	return nil
}

func deleteV5PolicyVersion(client *golangsdk.ServiceClient, policyId string, deleteVersionId ...string) error {
	versions, err := listV5PolicyVersions(client, policyId)
	if err != nil {
		return err
	}

	var deleteVersion string
	if len(deleteVersionId) > 0 && deleteVersionId[0] != "" {
		deleteVersion = deleteVersionId[0]
	} else {
		// Delete earliest version by default if deleteVersionId is not provided
		deleteVersion = utils.PathSearch("[-1].version_id", versions, "").(string)
	}

	httpUrl := "v5/policies/{policy_id}/versions/{version_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deletePath = strings.ReplaceAll(deletePath, "{version_id}", deleteVersion)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func addV5PolicyVersion(client *golangsdk.ServiceClient, policyId, document string) error {
	httpUrl := "v5/policies/{policy_id}/versions"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{policy_id}", policyId)
	addOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"policy_document": document,
			"set_as_default":  true,
		}),
	}

	_, err := client.Request("POST", addPath, &addOpt)
	return err
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

func resourceV5PolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		err = addV5PolicyVersion(client, policyId, d.Get("policy_document").(string))
		if err != nil {
			isVersionUpperLimit, err := handleAddVersionError409(err)
			// if get a "versions per policy limit exceeded" error remove a earliest version and try again
			if !isVersionUpperLimit {
				return diag.Errorf("error adding a new version of identity policy(%s): %s", policyId, err)
			}

			err = deleteV5PolicyVersion(client, policyId, d.Get("version_to_delete").(string))
			if err != nil {
				return diag.Errorf("error removing the earliest version of identity policy(%s): %s", policyId, err)
			}

			err = addV5PolicyVersion(client, policyId, d.Get("policy_document").(string))
			if err != nil {
				return diag.Errorf("error adding a new version of identity policy(%s): %s", policyId, err)
			}
		}
	}

	return resourceV5PolicyRead(ctx, d, meta)
}

func resourceV5PolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteHttpUrl := "v5/policies/{policy_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting identity policy (%s)", policyId))
	}
	return nil
}
