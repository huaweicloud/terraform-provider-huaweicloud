// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RMS
// ---------------------------------------------------------------

package rms

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Config PUT /v1/resource-manager/domains/{domain_id}/tracker-config
// @API Config GET /v1/resource-manager/domains/{domain_id}/tracker-config
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/tracker-config
func ResourceRecorder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecorderUpdate,
		UpdateContext: resourceRecorderUpdate,
		ReadContext:   resourceRecorderRead,
		DeleteContext: resourceRecorderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IAM agency name.`,
			},
			"selector": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     recorderSelectorConfigSchema(),
				Required: true,
			},
			"obs_channel": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Elem:         recorderOBSChannelConfigSchema(),
				Optional:     true,
				AtLeastOneOf: []string{"smn_channel"},
			},
			"smn_channel": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     recorderSMNChannelConfigSchema(),
				Optional: true,
			},
		},
	}
}

func recorderSelectorConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_supported": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to select all supported resources.`,
			},
			"resource_types": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the resource type list.`,
			},
		},
	}
	return &sc
}

func recorderOBSChannelConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the OBS bucket name.`,
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the region where this bucket is located.`,
			},
		},
	}
	return &sc
}

func recorderSMNChannelConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SMN topic URN.`,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region where this SMN topic is located.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the project ID where this SMN topic is located.`,
			},
		},
	}
	return &sc
}

func buildRecorderBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	channelOpts, err := buildRecorderChannelConfigRequest(d)
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"agency_name": d.Get("agency_name"),
		"selector":    buildRecorderSelectorConfigRequest(d.Get("selector")),
		"channel":     channelOpts,
	}
	return bodyParams, nil
}

func buildRecorderSelectorConfigRequest(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	resources := make([]interface{}, 0)
	isSupportAll := raw["all_supported"].(bool)
	// resource_types is only valid when all_supported is false
	if !isSupportAll {
		resources = raw["resource_types"].(*schema.Set).List()
	}

	return map[string]interface{}{
		"all_supported":  isSupportAll,
		"resource_types": resources,
	}
}

func buildRecorderChannelConfigRequest(d *schema.ResourceData) (map[string]interface{}, error) {
	obsOpts := buildRecorderOBSChannelConfigRequest(d.Get("obs_channel"))
	smnOpts, err := buildRecorderSMNChannelConfigRequest(d.Get("smn_channel"))
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"obs": utils.ValueIgnoreEmpty(obsOpts),
		"smn": utils.ValueIgnoreEmpty(smnOpts),
	}
	return bodyParams, nil
}

func buildRecorderOBSChannelConfigRequest(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"bucket_name": raw["bucket"].(string),
		"region_id":   raw["region"].(string),
	}
}

func buildRecorderSMNChannelConfigRequest(rawParams interface{}) (map[string]interface{}, error) {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil, nil
	}

	raw := rawArray[0].(map[string]interface{})
	var (
		err       error
		topicURN  = raw["topic_urn"].(string)
		region    = raw["region"].(string)
		projectID = raw["project_id"].(string)
	)
	if region == "" || projectID == "" {
		region, projectID, err = parseTopicURN(topicURN)
		if err != nil {
			return nil, err
		}
	}

	params := map[string]interface{}{
		"topic_urn":  topicURN,
		"region_id":  region,
		"project_id": projectID,
	}
	return params, nil
}

func resourceRecorderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		updateRecorderHttpUrl = "v1/resource-manager/domains/{domain_id}/tracker-config"
		updateRecorderProduct = "rms"
	)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateRecorderClient, err := cfg.NewServiceClient(updateRecorderProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	updateRecorderPath := updateRecorderClient.Endpoint + updateRecorderHttpUrl
	updateRecorderPath = strings.ReplaceAll(updateRecorderPath, "{domain_id}", cfg.DomainID)

	updateRecorderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	updateOpts, err := buildRecorderBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] the RMS recorder request options: %#v", updateOpts)
	updateRecorderOpt.JSONBody = utils.RemoveNil(updateOpts)
	_, err = updateRecorderClient.Request("PUT", updateRecorderPath, &updateRecorderOpt)
	if err != nil {
		return diag.Errorf("error creating or updating RMS recorder: %s", err)
	}

	d.SetId(cfg.DomainID)
	return resourceRecorderRead(ctx, d, meta)
}

func resourceRecorderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		getRecorderHttpUrl = "v1/resource-manager/domains/{domain_id}/tracker-config"
		getRecorderProduct = "rms"
	)

	cfg := meta.(*config.Config)
	getRecorderClient, err := cfg.NewServiceClient(getRecorderProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	getRecorderPath := getRecorderClient.Endpoint + getRecorderHttpUrl
	getRecorderPath = strings.ReplaceAll(getRecorderPath, "{domain_id}", cfg.DomainID)

	getRecorderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRecorderResp, err := getRecorderClient.Request("GET", getRecorderPath, &getRecorderOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS recorder")
	}

	getRecorderRespBody, err := utils.FlattenResponse(getRecorderResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("agency_name", utils.PathSearch("agency_name", getRecorderRespBody, nil)),
		d.Set("selector", flattenRecorderSelectorConfig(getRecorderRespBody)),
		d.Set("obs_channel", flattenRecorderOBSChannelConfig(getRecorderRespBody)),
		d.Set("smn_channel", flattenRecorderSMNChannelConfig(getRecorderRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRecorderSelectorConfig(resp interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"all_supported":  utils.PathSearch("selector.all_supported", resp, nil),
			"resource_types": utils.PathSearch("selector.resource_types", resp, nil),
		},
	}
}

func flattenRecorderOBSChannelConfig(resp interface{}) []interface{} {
	curJson := utils.PathSearch("channel.obs", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing channel.obs from response")
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"bucket": utils.PathSearch("bucket_name", curJson, nil),
			"region": utils.PathSearch("region_id", curJson, nil),
		},
	}
}

func flattenRecorderSMNChannelConfig(resp interface{}) []interface{} {
	curJson := utils.PathSearch("channel.smn", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing channel.smn from response")
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"topic_urn":  utils.PathSearch("topic_urn", curJson, nil),
			"region":     utils.PathSearch("region_id", curJson, nil),
			"project_id": utils.PathSearch("project_id", curJson, nil),
		},
	}
}

func resourceRecorderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		deleteRecorderHttpUrl = "v1/resource-manager/domains/{domain_id}/tracker-config"
		deleteRecorderProduct = "rms"
	)

	cfg := meta.(*config.Config)
	deleteRecorderClient, err := cfg.NewServiceClient(deleteRecorderProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	deleteRecorderPath := deleteRecorderClient.Endpoint + deleteRecorderHttpUrl
	deleteRecorderPath = strings.ReplaceAll(deleteRecorderPath, "{domain_id}", cfg.DomainID)

	deleteRecorderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteRecorderClient.Request("DELETE", deleteRecorderPath, &deleteRecorderOpt)
	if err != nil {
		return diag.Errorf("error deleting RMS recorder: %s", err)
	}

	return nil
}

// parseTopicURN is used to parse the region and project_ID, the format of topic_urn is as follows:
// urn:smn:{region}:{project ID}:{topic name}
func parseTopicURN(urn string) (region, projectID string, err error) {
	parts := strings.Split(urn, ":")
	if len(parts) != 5 {
		err = fmt.Errorf("cannot get region and project_id from topic_urn, " +
			"please check the format of topic_urn or specify them manually.")
		return
	}

	region = parts[2]
	projectID = parts[3]
	return
}
