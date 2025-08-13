package deprecated

import (
	"context"
	"fmt"
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

// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/protection-enable
// @API HSS DELETE /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets
// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/configs/batch-query
// @API HSS PUT /v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets

// ResourceCCEProtection
// 1. Currently, the APIs used by this resource have not been released in the official environment and can only be called
// in the `cn-north-7` region.
// 2. Because the API has not been released to the official environment, no documentation is provided.
// 3. Due to the limitations of the test environment, the capabilities provided by this resource are not reliable. The
// current test cases cannot be executed 100% successfully. We need to re-verify after the open API is launched on the
// official website.
func ResourceCCEProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCEProtectionCreate,
		ReadContext:   resourceCCEProtectionRead,
		UpdateContext: resourceCCEProtectionUpdate,
		DeleteContext: resourceCCEProtectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCCEProtectionImportState,
		},

		DeprecationMessage: "Abandon the resource and use a disposable resource instead.",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(12 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CCE cluster type.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CCE cluster ID.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CCE cluster name.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the charging mode.`,
			},
			"cce_protection_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CCE protection type.`,
			},
			// Fields `enterprise_project_id`, `agent_version`, `runtime_info`, `schedule_info`, and `invoked_service`
			// have no response values and do not need to be computed.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"agent_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the agent version.`,
			},
			"runtime_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        runTimeInfoSchemas(),
				Description: `Specifies the container runtime configuration.`,
			},
			"schedule_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        scheduleInfoSchemas(),
				Description: `Specifies the node scheduling information.`,
			},
			"invoked_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the invoked service.`,
			},
			"prefer_packet_cycle": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to prefer the package period quota.`,
			},
			"auto_upgrade": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable automatic agent upgrade.`,
			},

			// Attributes
			"protect_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of nodes in the cluster that have protection enabled.`,
			},
			"protect_interrupt_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of nodes interrupted by cluster protection.`,
			},
			"unprotect_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of unprotected nodes in the cluster.`,
			},
			"node_total_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of cluster nodes.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID.`,
			},
			"protect_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protection type.`,
			},
			"protect_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protection status.`,
			},
			"fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reasons for protection failure.`,
			},
		},
	}
}

func runTimeInfoSchemas() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"runtime_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func scheduleInfoSchemas() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"pod_tolerances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildCreateCCEProtectionBodyParam(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"cluster_type":        d.Get("cluster_type"),
		"cluster_id":          d.Get("cluster_id"),
		"cluster_name":        d.Get("cluster_name"),
		"charging_mode":       d.Get("charging_mode"),
		"cce_protection_type": d.Get("cce_protection_type"),
		"prefer_packet_cycle": d.Get("prefer_packet_cycle"),
	}

	return bodyParam
}

func createCCEProtection(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	var (
		httpUrl = "v5/{project_id}/container/kubernetes/clusters/protection-enable"
		epsID   = cfg.GetEnterpriseProjectID(d)
		region  = cfg.GetRegion(d)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsID != "" {
		requestPath += "?enterprise_project_id=" + epsID
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       region,
		},
		JSONBody: buildCreateCCEProtectionBodyParam(d),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func buildQueryCCEProtectionBodyParam(clusterID, clusterName string) interface{} {
	clusterBodyParam := map[string]interface{}{
		"cluster_id":   clusterID,
		"cluster_name": clusterName,
	}

	return map[string]interface{}{
		"cluster_info_list": []interface{}{clusterBodyParam},
	}
}

func ReadCCEProtection(client *golangsdk.ServiceClient, clusterID, clusterName, epsID, region string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/configs/batch-query"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsID != "" {
		requestPath += "?enterprise_project_id=" + epsID
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       region,
		},
		JSONBody: buildQueryCCEProtectionBodyParam(clusterID, clusterName),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForProtectStatusProtecting(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config, timeout time.Duration) error {
	var (
		clusterID   = d.Get("cluster_id").(string)
		clusterName = d.Get("cluster_name").(string)
		epsID       = cfg.GetEnterpriseProjectID(d)
		region      = cfg.GetRegion(d)
	)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadCCEProtection(client, clusterID, clusterName, epsID, region)
			if err != nil {
				return nil, "ERROR", err
			}

			clusterRespBody := utils.PathSearch("data_list|[0]", respBody, nil)
			if clusterRespBody == nil {
				return nil, "ERROR", fmt.Errorf("configuration is empty in API response")
			}

			protectStatus := utils.PathSearch("protect_status", clusterRespBody, "").(string)
			if protectStatus == "" {
				return nil, "ERROR", fmt.Errorf("protect_status is not found in API response")
			}

			if protectStatus == "protecting" {
				return clusterRespBody, "COMPLETED", nil
			}

			if protectStatus == "error_protect" {
				failReason := utils.PathSearch("fail_reason", clusterRespBody, "").(string)
				return clusterRespBody, "ERROR", fmt.Errorf("the protection status is error: %s", failReason)
			}

			// Due to the lack of test conditions to verify the accuracy of the logic, other states are considered
			// to be states that need to be waited for.
			return clusterRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func buildUpdateRuntimeInfoBodyParam(d *schema.ResourceData) []interface{} {
	rawArray, ok := d.Get("runtime_info").([]interface{})
	if !ok {
		return nil
	}
	rst := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		params := map[string]interface{}{
			"runtime_name": utils.PathSearch("runtime_name", v, nil),
			"runtime_path": utils.PathSearch("runtime_path", v, nil),
		}
		rst = append(rst, params)
	}
	return rst
}

func buildUpdateScheduleInfoBodyParam(d *schema.ResourceData) interface{} {
	rawArray, ok := d.Get("schedule_info").([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	return map[string]interface{}{
		"node_selector":  utils.PathSearch("schedule_info.0.node_selector", rawArray, nil),
		"pod_tolerances": utils.PathSearch("schedule_info.0.pod_tolerances", rawArray, nil),
	}
}

func buildUpdateCCEProtectionBodyParam(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"cluster_name":        d.Get("cluster_name"),
		"charging_mode":       d.Get("charging_mode"),
		"cce_protection_type": d.Get("cce_protection_type"),
		"agent_version":       utils.ValueIgnoreEmpty(d.Get("agent_version")),
		"auto_upgrade":        utils.ValueIgnoreEmpty(d.Get("auto_upgrade")),
		"runtime_info":        utils.ValueIgnoreEmpty(buildUpdateRuntimeInfoBodyParam(d)),
		"schedule_info":       utils.ValueIgnoreEmpty(buildUpdateScheduleInfoBodyParam(d)),
		"invoked_service":     utils.ValueIgnoreEmpty(d.Get("invoked_service")),
		"prefer_packet_cycle": utils.ValueIgnoreEmpty(d.Get("prefer_packet_cycle")),
	}

	return bodyParam
}

func updateCCEProtection(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	var (
		httpUrl   = "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
		epsID     = cfg.GetEnterpriseProjectID(d)
		clusterID = d.Get("cluster_id").(string)
		region    = cfg.GetRegion(d)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)
	if epsID != "" {
		requestPath += "?enterprise_project_id=" + epsID
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       region,
		},
		JSONBody: utils.RemoveNil(buildUpdateCCEProtectionBodyParam(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceCCEProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		clusterID = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := createCCEProtection(client, d, cfg); err != nil {
		return diag.Errorf("error opening HSS protection for CCE (%s): %s", clusterID, err)
	}
	d.SetId(clusterID)

	if err := waitingForProtectStatusProtecting(ctx, client, d, cfg, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CCE (%s) protection configuration open to complete: %s", clusterID, err)
	}

	if err := updateCCEProtection(client, d, cfg); err != nil {
		return diag.Errorf("error configuring HSS protection configuration for CCE (%s) in"+
			" creation operation: %s", clusterID, err)
	}

	return resourceCCEProtectionRead(ctx, d, meta)
}

func flattenPreferPacketCycle(clusterRespBody interface{}) bool {
	preferPacketCycle := int(utils.PathSearch("prefer_packet_cycle", clusterRespBody, float64(0)).(float64))
	return preferPacketCycle != 0
}

func resourceCCEProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "hss"
		clusterID   = d.Get("cluster_id").(string)
		clusterName = d.Get("cluster_name").(string)
		epsID       = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	respBody, err := ReadCCEProtection(client, clusterID, clusterName, epsID, region)
	if err != nil {
		return diag.Errorf("error retrieving HSS protection configuration for CCE: %s", err)
	}

	clusterRespBody := utils.PathSearch("data_list|[0]", respBody, nil)
	if clusterRespBody == nil {
		return diag.Errorf("error retrieving HSS protection configuration for CCE: configuration is empty in API" +
			" response")
	}

	protectStatus := utils.PathSearch("protect_status", clusterRespBody, "").(string)
	if protectStatus == "" {
		return diag.Errorf("error retrieving HSS protection configuration for CCE: protect_status is not found" +
			" in API response")
	}

	// After the deletion is complete, the data can still be found.
	// This `404` judgment can only cover some scenarios where resources do not exist.
	if protectStatus == "unprotect" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", utils.PathSearch("cluster_id", clusterRespBody, nil)),
		d.Set("protect_node_num", utils.PathSearch("protect_node_num", clusterRespBody, nil)),
		d.Set("protect_interrupt_node_num", utils.PathSearch("protect_interrupt_node_num", clusterRespBody, nil)),
		d.Set("unprotect_node_num", utils.PathSearch("unprotect_node_num", clusterRespBody, nil)),
		d.Set("node_total_num", utils.PathSearch("node_total_num", clusterRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", clusterRespBody, nil)),
		d.Set("charging_mode", utils.PathSearch("charging_mode", clusterRespBody, nil)),
		d.Set("auto_upgrade", utils.PathSearch("auto_upgrade", clusterRespBody, nil)),
		d.Set("prefer_packet_cycle", flattenPreferPacketCycle(clusterRespBody)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", clusterRespBody, nil)),
		d.Set("protect_type", utils.PathSearch("protect_type", clusterRespBody, nil)),
		d.Set("protect_status", utils.PathSearch("protect_status", clusterRespBody, nil)),
		d.Set("cluster_type", utils.PathSearch("cluster_type", clusterRespBody, nil)),
		d.Set("fail_reason", utils.PathSearch("fail_reason", clusterRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCCEProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		clusterID = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updateCCEProtection(client, d, cfg); err != nil {
		return diag.Errorf("error configuring HSS protection configuration for CCE (%s) in"+
			" update operation: %s", clusterID, err)
	}

	return resourceCCEProtectionRead(ctx, d, meta)
}

// The deletion operation only deletes some configurations and does not turn off protection.
// After the deletion is complete, the data can still be found.
func resourceCCEProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		epsID     = cfg.GetEnterpriseProjectID(d)
		clusterID = d.Get("cluster_id").(string)
		httpUrl   = "v5/{project_id}/container/kubernetes/clusters/{cluster_id}/daemonsets"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)
	if epsID != "" {
		requestPath += "?enterprise_project_id=" + epsID
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"region":       region,
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS protection configuration for CCE (%s): %s", clusterID, err)
	}

	return nil
}

func resourceCCEProtectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<cluster_name>',"+
			" but got '%s'", importedId)
	}
	d.SetId(parts[0])
	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
		d.Set("cluster_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
