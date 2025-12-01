package hss

import (
	"context"
	"errors"
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

// @API HSS POST /v5/{project_id}/ransomware/protection/policy
// @API HSS DELETE /v5/{project_id}/ransomware/protection/policy
// @API HSS PUT /v5/{project_id}/ransomware/protection/policy
// @API HSS GET /v5/{project_id}/ransomware/protection/policy
func ResourceRansomwareProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRansomwareProtectionPolicyCreate,
		ReadContext:   resourceRansomwareProtectionPolicyRead,
		UpdateContext: resourceRansomwareProtectionPolicyUpdate,
		DeleteContext: resourceRansomwareProtectionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRansomwareProtectionPolicyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"enterprise_project_id",
			"operating_system",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protection_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protection_directory": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protection_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The field `enterprise_project_id` does not exist in API response body.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"exclude_directory": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The field `process_whitelist.hash` does not exist in API response body.
			// Instead of backfilling data for this field, provide an another attribute field.
			"process_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     processWhitelistSchema(),
			},
			// The field `agent_id_list` does not exist in API response body.
			// Field `agent_id_list` using only in update operation.
			"agent_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Field `runtime_detection_status` using only in update operation.
			"runtime_detection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ai_protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `bait_protection_status` using only in update operation.
			"bait_protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"process_whitelist_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     processWhitelistAttributeSchema(),
			},
			"count_associated_server": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func processWhitelistAttributeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func processWhitelistSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildRansomwareProtectionPolicyQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}

	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildProcessWhitelistRequestOpt(processWhitelist []interface{}) []map[string]interface{} {
	if len(processWhitelist) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(processWhitelist))
	for _, v := range processWhitelist {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"path": utils.ValueIgnoreEmpty(rawMap["path"]),
			"hash": utils.ValueIgnoreEmpty(rawMap["hash"]),
		})
	}

	return rst
}

func buildCreateRansomwareProtectionPolicyRequestOpt(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name":          d.Get("policy_name"),
		"protection_mode":      d.Get("protection_mode"),
		"protection_directory": d.Get("protection_directory"),
		"protection_type":      d.Get("protection_type"),
		"operating_system":     d.Get("operating_system"),
		"deploy_mode":          utils.ValueIgnoreEmpty(d.Get("deploy_mode")),
		"exclude_directory":    utils.ValueIgnoreEmpty(d.Get("exclude_directory")),
		"process_whitelist":    buildProcessWhitelistRequestOpt(d.Get("process_whitelist").([]interface{})),
		"ai_protection_status": utils.ValueIgnoreEmpty(d.Get("ai_protection_status")),
	}
}

func buildProtectionPolicyByPolicyNameQueryParams(policyName, epsId string) string {
	rst := fmt.Sprintf("?policy_name=%s", policyName)
	if epsId == "" {
		return rst
	}

	return fmt.Sprintf("%s&enterprise_project_id=%s", rst, epsId)
}

func buildProtectionPolicyByPolicyNameQueryParamsWithOffset(requestPath string, offset int) string {
	if offset == 0 {
		return requestPath
	}

	return fmt.Sprintf("%s&offset=%d", requestPath, offset)
}

// The search for policy_name is a fuzzy search, so pagination is required.
func QueryProtectionPolicyByPolicyName(client *golangsdk.ServiceClient, policyName, epsId, region string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/ransomware/protection/policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildProtectionPolicyByPolicyNameQueryParams(policyName, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	var (
		allResult []interface{}
		offset    int
	)

	for {
		requestPathWithOffset := buildProtectionPolicyByPolicyNameQueryParamsWithOffset(requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		allResult = append(allResult, dataList...)
		offset += len(dataList)
	}

	expression := fmt.Sprintf("[?policy_name == '%s']|[0]", policyName)
	targetPolicy := utils.PathSearch(expression, allResult, nil)
	if targetPolicy == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return targetPolicy, nil
}

func resourceRansomwareProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		epsId      = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product    = "hss"
		httpUrl    = "v5/{project_id}/ransomware/protection/policy"
		policyName = d.Get("policy_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRansomwareProtectionPolicyQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         utils.RemoveNil(buildCreateRansomwareProtectionPolicyRequestOpt(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating HSS ransomware protection policy: %s", err)
	}

	policyDetail, err := QueryProtectionPolicyByPolicyName(client, policyName, epsId, region)
	if err != nil {
		return diag.Errorf("error querying HSS ransomware protection policy by policy name: %s", err)
	}

	policyId := utils.PathSearch("policy_id", policyDetail, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating HSS ransomware protection policy: policy ID is empty")
	}

	d.SetId(policyId)

	// Fields `agent_id_list`, `runtime_detection_status`, `bait_protection_status` using only in update operation.
	updateFields := []string{"agent_id_list", "runtime_detection_status", "bait_protection_status"}
	if d.HasChanges(updateFields...) {
		if err := updateRansomwareProtectionPolicy(client, d, cfg); err != nil {
			return diag.Errorf("error updating HSS ransomware protection policy in create operation: %s", err)
		}
	}

	return resourceRansomwareProtectionPolicyRead(ctx, d, meta)
}

func buildProtectionPolicyByPolicyIdQueryParams(policyId, epsId string) string {
	rst := fmt.Sprintf("?protect_policy_id=%s", policyId)
	if epsId == "" {
		return rst
	}

	return fmt.Sprintf("%s&enterprise_project_id=%s", rst, epsId)
}

func QueryProtectionPolicyByPolicyId(client *golangsdk.ServiceClient, policyId, epsId, region string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/ransomware/protection/policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildProtectionPolicyByPolicyIdQueryParams(policyId, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	policyDetail := utils.PathSearch("data_list|[0]", respBody, nil)
	if policyDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return policyDetail, nil
}

func resourceRansomwareProtectionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	policyDetail, err := QueryProtectionPolicyByPolicyId(client, d.Id(), epsId, region)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving HSS ransomware protection policy")
	}

	processWhitelistAttribute := utils.PathSearch("process_whitelist", policyDetail, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policy_name", utils.PathSearch("policy_name", policyDetail, nil)),
		d.Set("protection_mode", utils.PathSearch("protection_mode", policyDetail, nil)),
		d.Set("bait_protection_status", utils.PathSearch("bait_protection_status", policyDetail, nil)),
		d.Set("deploy_mode", utils.PathSearch("deploy_mode", policyDetail, nil)),
		d.Set("protection_directory", utils.PathSearch("protection_directory", policyDetail, nil)),
		d.Set("protection_type", utils.PathSearch("protection_type", policyDetail, nil)),
		d.Set("exclude_directory", utils.PathSearch("exclude_directory", policyDetail, nil)),
		d.Set("runtime_detection_status", utils.PathSearch("runtime_detection_status", policyDetail, nil)),
		d.Set("count_associated_server", utils.PathSearch("count_associated_server", policyDetail, nil)),
		d.Set("operating_system", utils.PathSearch("operating_system", policyDetail, nil)),
		d.Set("process_whitelist_attribute", flattenProcessWhitelistAttribute(processWhitelistAttribute)),
		d.Set("default_policy", utils.PathSearch("default_policy", policyDetail, nil)),
		d.Set("ai_protection_status", utils.PathSearch("ai_protection_status", policyDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenProcessWhitelistAttribute(processWhitelistAttribute []interface{}) []interface{} {
	if len(processWhitelistAttribute) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(processWhitelistAttribute))
	for _, v := range processWhitelistAttribute {
		rst = append(rst, map[string]interface{}{
			"path": utils.PathSearch("path", v, nil),
			"hash": utils.PathSearch("hash", v, nil),
		})
	}

	return rst
}

func buildUpdateRansomwareProtectionPolicyRequestOpt(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_id":                d.Id(),
		"policy_name":              d.Get("policy_name"),
		"protection_mode":          d.Get("protection_mode"),
		"protection_directory":     d.Get("protection_directory"),
		"protection_type":          d.Get("protection_type"),
		"operating_system":         d.Get("operating_system"),
		"bait_protection_status":   utils.ValueIgnoreEmpty(d.Get("bait_protection_status")),
		"deploy_mode":              utils.ValueIgnoreEmpty(d.Get("deploy_mode")),
		"exclude_directory":        utils.ValueIgnoreEmpty(d.Get("exclude_directory")),
		"agent_id_list":            utils.ExpandToStringList(d.Get("agent_id_list").([]interface{})),
		"runtime_detection_status": utils.ValueIgnoreEmpty(d.Get("runtime_detection_status")),
		"process_whitelist":        buildProcessWhitelistRequestOpt(d.Get("process_whitelist").([]interface{})),
		"ai_protection_status":     utils.ValueIgnoreEmpty(d.Get("ai_protection_status")),
	}
}

func updateRansomwareProtectionPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	requestPath := client.Endpoint + "v5/{project_id}/ransomware/protection/policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildRansomwareProtectionPolicyQueryParams(cfg.GetEnterpriseProjectID(d, QueryAllEpsValue))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": cfg.GetRegion(d)},
		JSONBody:         utils.RemoveNil(buildUpdateRansomwareProtectionPolicyRequestOpt(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceRansomwareProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err = updateRansomwareProtectionPolicy(client, d, cfg); err != nil {
		return diag.Errorf("error updating HSS ransomware protection policy in update operation: %s", err)
	}

	return resourceRansomwareProtectionPolicyRead(ctx, d, meta)
}

func buildDeleteRansomwareProtectionPolicyQueryParams(policyId, epsId string) string {
	rst := fmt.Sprintf("?policy_id=%s", policyId)
	if epsId == "" {
		return rst
	}

	return fmt.Sprintf("%s&enterprise_project_id=%s", rst, epsId)
}

func resourceRansomwareProtectionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
		httpUrl = "v5/{project_id}/ransomware/protection/policy"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteRansomwareProtectionPolicyQueryParams(d.Id(), epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS ransomware protection policy: %s", err)
	}

	return nil
}

func resourceRansomwareProtectionPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <enterprise_project_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[0])
}
