package sfsturbo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/targets
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}/policy
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}/attributes
func ResourceOBSTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOBSTargetCreate,
		ReadContext:   resourceOBSTargetRead,
		UpdateContext: resourceOBSTargetUpdate,
		DeleteContext: resourceOBSTargetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOBSTargetImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_system_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"obs": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     obsSchema(),
			},
			"delete_data_in_file_system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
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

func obsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_export_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"events": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"prefix": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"suffix": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dir_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"gid": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func buildCreateOBSTargetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_system_path": d.Get("file_system_path"),
		"obs":              buildOBSBody(d.Get("obs.0").(map[string]interface{})),
	}
	return bodyParams
}

func buildOBSBody(obsData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"bucket":     obsData["bucket"],
		"endpoint":   obsData["endpoint"],
		"policy":     buildPolicyBodyParams(obsData["policy"].([]interface{})),
		"attributes": buildAttributesBodyParams(obsData["attributes"].([]interface{})),
	}
}

func buildPolicyBodyParams(rawPolicy []interface{}) map[string]interface{} {
	if len(rawPolicy) == 0 {
		return nil
	}

	policy, ok := rawPolicy[0].(map[string]interface{})
	if !ok {
		return nil
	}

	policyParams := map[string]interface{}{
		"auto_export_policy": buildExportPolicyBodyParams(policy["auto_export_policy"].([]interface{})),
	}

	return policyParams
}

func buildExportPolicyBodyParams(rawExportPolicy []interface{}) map[string]interface{} {
	if len(rawExportPolicy) == 0 {
		return nil
	}

	exportPolicy, ok := rawExportPolicy[0].(map[string]interface{})
	if !ok {
		return nil
	}

	events := utils.ExpandToStringList(exportPolicy["events"].([]interface{}))

	exportPolicyParams := map[string]interface{}{
		"prefix": utils.ValueIgnoreEmpty(exportPolicy["prefix"]),
		"suffix": utils.ValueIgnoreEmpty(exportPolicy["suffix"]),
	}

	if len(events) > 0 {
		exportPolicyParams["events"] = events
	}

	return exportPolicyParams
}

func buildAttributesBodyParams(rawAttributes []interface{}) map[string]interface{} {
	if len(rawAttributes) == 0 {
		return nil
	}

	attributes, ok := rawAttributes[0].(map[string]interface{})
	if !ok {
		return nil
	}

	attributesParams := map[string]interface{}{
		"file_mode": bulidAttributeMode(attributes["file_mode"].(string)),
		"dir_mode":  bulidAttributeMode(attributes["dir_mode"].(string)),
		"uid":       attributes["uid"],
		"gid":       attributes["gid"],
	}

	return attributesParams
}

func bulidAttributeMode(str string) interface{} {
	resp, err := strconv.Atoi(str)
	if err != nil {
		return nil
	}

	return resp
}

func resourceOBSTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
	)

	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets"
	createObsTargetPath := client.ResourceBaseURL() + createObsTargetHttpUrl
	createObsTargetPath = strings.ReplaceAll(createObsTargetPath, "{share_id}", shareId)

	createObsTargetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createObsTargetOpt.JSONBody = utils.RemoveNil(buildCreateOBSTargetBodyParams(d))
	createObsTargetResp, err := client.Request("POST", createObsTargetPath, &createObsTargetOpt)
	if err != nil {
		return diag.Errorf("error creating OBS target to the SFS Turbo: %s", err)
	}

	createObsTargetRespBody, err := utils.FlattenResponse(createObsTargetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	targetId := utils.PathSearch("target_id", createObsTargetRespBody, "").(string)
	if targetId == "" {
		return diag.Errorf("unable to find the OBS target ID (source for the SFS Turbo) from the API response")
	}

	d.SetId(targetId)

	err = obsTargetWaitingForStateCompleted(ctx, client, shareId, targetId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of OBS target (%s) to complete: %s", d.Id(), err)
	}

	return resourceOBSTargetRead(ctx, d, meta)
}

func getOBSTargetInfo(client *golangsdk.ServiceClient, shareId, targetId string) (interface{}, error) {
	getObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets/{target_id}"
	getObsTargetPath := client.ResourceBaseURL() + getObsTargetHttpUrl
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{share_id}", shareId)
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{target_id}", targetId)
	getObsTargetOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getObsTargetPath, &getObsTargetOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceOBSTargetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	respBody, err := getOBSTargetInfo(client, d.Get("share_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "SFS Turbo OBS target")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("file_system_path", utils.PathSearch("file_system_path", respBody, nil)),
		d.Set("obs", flattenGetOBSDataResponseBody(utils.PathSearch("obs", respBody, nil))),
		d.Set("status", utils.PathSearch("lifecycle", respBody, nil)),
		d.Set("created_at", utils.PathSearch("creation_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetOBSDataResponseBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"bucket":     utils.PathSearch("bucket", resp, nil),
			"endpoint":   utils.PathSearch("endpoint", resp, nil),
			"policy":     flattenPolicy(utils.PathSearch("policy", resp, nil)),
			"attributes": flattenAttributes(utils.PathSearch("attributes", resp, nil)),
		},
	}
}

func flattenPolicy(policy interface{}) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	policyResult := map[string]interface{}{
		"auto_export_policy": flattenExportPolicy(utils.PathSearch("auto_export_policy", policy, nil)),
	}

	return []map[string]interface{}{policyResult}
}

func flattenExportPolicy(exportPolicy interface{}) []map[string]interface{} {
	if exportPolicy == nil {
		return nil
	}

	exportPolicyResult := map[string]interface{}{
		"events": utils.PathSearch("events", exportPolicy, nil),
		"prefix": utils.PathSearch("prefix", exportPolicy, nil),
		"suffix": utils.PathSearch("suffix", exportPolicy, nil),
	}

	return []map[string]interface{}{exportPolicyResult}
}

func flattenAttributes(attributes interface{}) []map[string]interface{} {
	if attributes == nil {
		return nil
	}

	attributesResult := map[string]interface{}{
		"file_mode": flattenAttributeMode(utils.PathSearch("file_mode", attributes, nil)),
		"dir_mode":  flattenAttributeMode(utils.PathSearch("dir_mode", attributes, nil)),
		"uid":       utils.PathSearch("uid", attributes, nil),
		"gid":       utils.PathSearch("gid", attributes, nil),
	}

	return []map[string]interface{}{attributesResult}
}

func flattenAttributeMode(param interface{}) string {
	if param == nil {
		return ""
	}

	return fmt.Sprintf("%v", param)
}

func resourceOBSTargetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 Client: %s", err)
	}

	if d.HasChange("obs.0.policy.0.auto_export_policy") {
		if err := updateTargetPolicy(client, d, shareId, d.Id()); err != nil {
			return diag.Errorf("error updating the policy of the OBS target: %s", err)
		}
	}

	if d.HasChange("obs.0.attributes") {
		if err := updateTargetAttributes(client, d, shareId, d.Id()); err != nil {
			return diag.Errorf("error updating the attributes of the OBS target: %s", err)
		}
	}

	return resourceOBSTargetRead(ctx, d, meta)
}

func updateTargetPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, shareId, targetId string) error {
	httpUrl := "v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}/policy"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{share_id}", shareId)
	updatePath = strings.ReplaceAll(updatePath, "{target_id}", targetId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(updateExportPolicyBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func updateExportPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawExportPolicy := d.Get("obs.0.policy.0.auto_export_policy.0")

	exportPolicyParams := map[string]interface{}{
		"prefix": utils.ValueIgnoreEmpty(utils.PathSearch("prefix", rawExportPolicy, nil)),
		"suffix": utils.ValueIgnoreEmpty(utils.PathSearch("suffix", rawExportPolicy, nil)),
		"events": utils.ExpandToStringList(utils.PathSearch("events", rawExportPolicy, make([]interface{}, 0)).([]interface{})),
	}

	policyParams := map[string]interface{}{
		"policy": map[string]interface{}{
			"auto_export_policy": exportPolicyParams,
		},
	}

	return policyParams
}

func updateTargetAttributes(client *golangsdk.ServiceClient, d *schema.ResourceData, shareId, targetId string) error {
	httpUrl := "v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}/attributes"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{share_id}", shareId)
	updatePath = strings.ReplaceAll(updatePath, "{target_id}", targetId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(updateAttributesBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func updateAttributesBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawAttributes := d.Get("obs.0.attributes.0")

	filMode := utils.PathSearch("file_mode", rawAttributes, "").(string)
	dirMode := utils.PathSearch("dir_mode", rawAttributes, "").(string)

	attributesParams := map[string]interface{}{
		"file_mode": bulidAttributeMode(filMode),
		"dir_mode":  bulidAttributeMode(dirMode),
		"uid":       utils.PathSearch("uid", rawAttributes, nil),
		"gid":       utils.PathSearch("gid", rawAttributes, nil),
	}

	attributesBodyParams := map[string]interface{}{
		"attributes": attributesParams,
	}

	return attributesBodyParams
}

func resourceOBSTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
	)

	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 Client: %s", err)
	}

	deleteObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets/{target_id}"
	deleteObsTargetPath := client.ResourceBaseURL() + deleteObsTargetHttpUrl
	deleteObsTargetPath = strings.ReplaceAll(deleteObsTargetPath, "{share_id}", d.Get("share_id").(string))
	deleteObsTargetPath = strings.ReplaceAll(deleteObsTargetPath, "{target_id}", d.Id())

	if v, ok := d.GetOk("delete_data_in_file_system"); ok {
		deleteObsTargetPath += fmt.Sprintf("?delete_data_in_file_system=%v", v)
	}

	deleteObsTargetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteObsTargetPath, &deleteObsTargetOpt)
	if err != nil {
		return diag.Errorf("error deleting OBS target from SFS Turbo: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      obsTargetStatusRefreshFunc(client, shareId, d.Id(), true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func obsTargetWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, shareId, targetId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      obsTargetStatusRefreshFunc(client, shareId, targetId, false),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func obsTargetStatusRefreshFunc(client *golangsdk.ServiceClient, shareId, targetId string, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getOBSTargetInfo(client, shareId, targetId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}

			return nil, "ERROR", err
		}

		status := utils.PathSearch("lifecycle", respBody, "").(string)

		if utils.StrSliceContains([]string{"MISCONFIGURED", "FAILED"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status: '%s'", status)
		}

		if utils.StrSliceContains([]string{"AVAILABLE"}, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceOBSTargetImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<share_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("share_id", parts[0])
}
