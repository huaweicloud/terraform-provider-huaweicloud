package iam

import (
	"context"
	"encoding/json"
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

const pageLimit = 100

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
func ResourceIAMServiceAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMServiceAgencyCreate,
		ReadContext:   resourceIAMServiceAgencyRead,
		UpdateContext: resourceIAMServiceAgencyUpdate,
		DeleteContext: resourceIAMServiceAgencyDelete,

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
			"delegated_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"trust_policy": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceIAMServiceAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		JSONBody:         utils.RemoveNil(buildCreateAgencyBodyParams(d)),
	}
	createAgencyResp, err := client.Request("POST", createAgencyPath, &createAgencyOpt)
	if err != nil {
		return diag.Errorf("error creating IAM service agency: %s", err)
	}
	createAgencyRespBody, err := utils.FlattenResponse(createAgencyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("agency.agency_id", createAgencyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating IAM service agency: agency_id is not found in API response")
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

	return resourceIAMServiceAgencyRead(ctx, d, meta)
}

func buildCreateAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agency_name":          d.Get("name"),
		"trust_policy":         buildTrustPolicy(d.Get("delegated_service_name").(string)),
		"path":                 utils.ValueIgnoreEmpty(d.Get("path")),
		"max_session_duration": utils.ValueIgnoreEmpty(d.Get("duration")),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func buildTrustPolicy(delegatedServiceName string) string {
	// only the delegated_service_name is variable
	statement := make([]map[string]interface{}, 1)
	v := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Service": []string{delegatedServiceName},
		},
		"Action": []string{"sts:agencies:assume"},
	}
	statement[0] = v
	trustPolicy := map[string]interface{}{
		"Version":   "5.0",
		"Statement": statement,
	}
	s, _ := json.Marshal(trustPolicy)

	return string(s)
}

func resourceIAMServiceAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	getAgencyResp, err := client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM service")
	}
	getAgencyRespBody, err := utils.FlattenResponse(getAgencyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	agency := utils.PathSearch("agency", getAgencyRespBody, nil)
	if agency == nil {
		return diag.Errorf("error getting IAM service agency: agency is not found in API response")
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
		d.Set("delegated_service_name", utils.PathSearch("Statement[0].Principal.Service[0]", policy, nil)),
		d.Set("policy_names", policyNames),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM service agency fields: %s", err)
	}

	return nil
}

func flattenTagsToMap(tags interface{}) map[string]interface{} {
	if tagArray, ok := tags.([]interface{}); ok {
		result := make(map[string]interface{})
		for _, val := range tagArray {
			if t, ok := val.(map[string]interface{}); ok {
				result[t["tag_key"].(string)] = t["tag_value"]
			}
		}
		return result
	}
	return nil
}

func listAttachedPolicyInfo(client *golangsdk.ServiceClient, agencyID, infoType string) ([]string, error) {
	listPolicyHttpUrl := "v5/agencies/{agency_id}/attached-policies"
	listPolicyPath := client.Endpoint + listPolicyHttpUrl
	listPolicyPath = strings.ReplaceAll(listPolicyPath, "{agency_id}", agencyID)
	listPolicyPath += fmt.Sprintf("?limit=%v", pageLimit)
	listPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	path := listPolicyPath
	rst := make([]string, 0)
	for {
		listPolicyResp, err := client.Request("GET", path, &listPolicyOpt)
		if err != nil {
			return nil, fmt.Errorf("error getting IAM policies: %s", err)
		}
		listPolicyRespBody, err := utils.FlattenResponse(listPolicyResp)
		if err != nil {
			return nil, err
		}

		// only get policy names or policy IDs
		jsonPath := fmt.Sprintf("attached_policies[*].%s", infoType)
		policyInfos := utils.PathSearch(jsonPath, listPolicyRespBody, make([]interface{}, 0))
		rst = append(rst, utils.ExpandToStringList(policyInfos.([]interface{}))...)

		marker := utils.PathSearch("page_info.next_marker", listPolicyRespBody, "")
		if marker == "" {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", listPolicyPath, marker)
	}
	return rst, nil
}

func resourceIAMServiceAgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			return diag.Errorf("error updating IAM service agency: %s", err)
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

	return resourceIAMServiceAgencyRead(ctx, d, meta)
}

func getPolicyIDsByNames(client *golangsdk.ServiceClient, policyNames []interface{}) ([]string, error) {
	listPolicyHttpUrl := "v5/policies"
	listPolicyPath := client.Endpoint + listPolicyHttpUrl
	listPolicyPath += fmt.Sprintf("?limit=%v", pageLimit)
	listPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	path := listPolicyPath
	rst := make([]string, 0)
	for {
		listPolicyResp, err := client.Request("GET", path, &listPolicyOpt)
		if err != nil {
			return nil, fmt.Errorf("error getting IAM policies: %s", err)
		}
		listPolicyRespBody, err := utils.FlattenResponse(listPolicyResp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("policies", listPolicyRespBody, make([]interface{}, 0))
		for _, policy := range policies.([]interface{}) {
			policyName := utils.PathSearch("policy_name", policy, "")
			if policyName == "" {
				return nil, fmt.Errorf("error getting policy name")
			}
			policyID := utils.PathSearch("policy_id", policy, "")
			if policyID == "" {
				return nil, fmt.Errorf("error getting policy ID for name(%s)", policyName)
			}

			// use policy name to filter result
			if utils.StrSliceContains(utils.ExpandToStringList(policyNames), policyName.(string)) {
				rst = append(rst, policyID.(string))
			}
		}

		// break when all names were found
		if len(rst) == len(policyNames) {
			break
		}

		marker := utils.PathSearch("page_info.next_marker", listPolicyRespBody, "")
		if marker == "" {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", listPolicyPath, marker)
	}
	return rst, nil
}

func attachPolicyByID(client *golangsdk.ServiceClient, agencyID, policyID string) error {
	attachHttpUrl := "v5/policies/{policy_id}/attach-agency"
	attachPath := client.Endpoint + attachHttpUrl
	attachPath = strings.ReplaceAll(attachPath, "{policy_id}", policyID)
	attachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"agency_id": agencyID,
		},
	}
	_, err := client.Request("POST", attachPath, &attachOpt)
	if err != nil {
		return fmt.Errorf("error attaching IAM service agency to policy(%s): %s", policyID, err)
	}
	return nil
}

func detachPolicyByID(client *golangsdk.ServiceClient, agencyID, policyID string) error {
	detachHttpUrl := "v5/policies/{policy_id}/detach-agency"
	detachPath := client.Endpoint + detachHttpUrl
	detachPath = strings.ReplaceAll(detachPath, "{policy_id}", policyID)
	detachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"agency_id": agencyID,
		},
	}
	_, err := client.Request("POST", detachPath, &detachOpt)
	if err != nil {
		return fmt.Errorf("error detaching IAM service agency to policy(%s): %s", policyID, err)
	}
	return nil
}

func createTags(createTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, id string) error {
	if len(tags) > 0 {
		createTagsHttpUrl := "v5/{resource_type}/{resource_id}/tags/create"
		createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", resourceType)
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", id)
		createTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"tags": expandResourceTags(tags),
			},
		}

		_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
		if err != nil {
			return fmt.Errorf("error creating tags: %s", err)
		}
	}
	return nil
}

func expandResourceTags(tagmap map[string]interface{}) []map[string]interface{} {
	var taglist []map[string]interface{}
	for k, v := range tagmap {
		tag := map[string]interface{}{
			"tag_key":   k,
			"tag_value": v,
		}
		taglist = append(taglist, tag)
	}
	return taglist
}

func deleteTags(deleteTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, id string) error {
	if len(tags) > 0 {
		deleteTagsHttpUrl := "v5/{resource_type}/{resource_id}/tags/delete"
		deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
		deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", resourceType)
		deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", id)
		deleteTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         expandTagsKeyToStringList(tags),
		}

		_, err := deleteTagsClient.Request("DELETE", deleteTagsPath, &deleteTagsOpt)
		if err != nil {
			return fmt.Errorf("error deleting tags: %s", err)
		}
	}
	return nil
}

func expandTagsKeyToStringList(tagmap map[string]interface{}) []string {
	var tagKeyList []string
	for k := range tagmap {
		tagKeyList = append(tagKeyList, k)
	}
	return tagKeyList
}

func resourceIAMServiceAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam_no_version", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// before deleting service agency, require to detach policies first
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

	// delete service agency
	deleteAgencyHttpUrl := "v5/agencies/{agency_id}"
	deleteAgencyPath := client.Endpoint + deleteAgencyHttpUrl
	deleteAgencyPath = strings.ReplaceAll(deleteAgencyPath, "{agency_id}", d.Id())
	deleteAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteAgencyPath, &deleteAgencyOpt)
	if err != nil {
		return diag.Errorf("error deleting IAM service agency: %s", err)
	}

	return nil
}
