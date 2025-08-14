package hss

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/protection-enable
// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/configs/batch-query
func ResourceContainerKubernetesClusterProtectionEnable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterProtectionEnableCreate,
		ReadContext:   resourceClusterProtectionEnableRead,
		UpdateContext: resourceClusterProtectionEnableUpdate,
		DeleteContext: resourceClusterProtectionEnableDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"cluster_type",
			"cluster_id",
			"cluster_name",
			"charging_mode",
			"cce_protection_type",
			"prefer_packet_cycle",
			"enterprise_project_id",
			"monitor_protection_statuses",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located.",
			},
			// Body Params
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of CCE cluster.",
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the charging mode.",
			},
			"cce_protection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the CCE protection type.",
			},
			"prefer_packet_cycle": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to prioritize the use of packet cycle quota.",
			},
			// Query Params
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the resource belongs.",
			},
			// The `monitor_protection_statuses` field is a custom field used to define which protection statuses to monitor.
			// The user determines the monitoring status of the CCE integrated protection configuration.
			"monitor_protection_statuses": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the monitor protection status.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"protect_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protection status.",
			},
		},
	}
}

func buildCreateClusterProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"cluster_name":        d.Get("cluster_name"),
		"cluster_id":          d.Get("cluster_id"),
		"cluster_type":        utils.ValueIgnoreEmpty(d.Get("cluster_type")),
		"charging_mode":       utils.ValueIgnoreEmpty(d.Get("charging_mode")),
		"cce_protection_type": utils.ValueIgnoreEmpty(d.Get("cce_protection_type")),
	}

	if d.Get("prefer_packet_cycle").(bool) {
		rst["prefer_packet_cycle"] = true
	}

	return rst
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

func ReadCCEProtection(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) (interface{}, error) {
	var (
		clusterID   = d.Get("cluster_id").(string)
		clusterName = d.Get("cluster_name").(string)
		epsID       = cfg.GetEnterpriseProjectID(d)
		region      = cfg.GetRegion(d)
		httpUrl     = "v5/{project_id}/container/kubernetes/clusters/configs/batch-query"
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
		JSONBody: buildQueryCCEProtectionBodyParam(clusterID, clusterName),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func getMonitorProtectionStatuses(d *schema.ResourceData) []string {
	monitorProtectionStatuses := utils.ExpandToStringList(d.Get("monitor_protection_statuses").([]interface{}))
	if len(monitorProtectionStatuses) == 0 {
		monitorProtectionStatuses = append(monitorProtectionStatuses, "protecting")
	}

	return monitorProtectionStatuses
}

func waitingForProtectStatusProtecting(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	cfg *config.Config, timeout time.Duration) (interface{}, error) {
	monitorProtectionStatuses := getMonitorProtectionStatuses(d)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadCCEProtection(client, cfg, d)
			if err != nil {
				return nil, "ERROR", err
			}

			clusterRespBody := utils.PathSearch("data_list|[0]", respBody, nil)
			if clusterRespBody == nil {
				return nil, "ERROR", errors.New("configuration is empty in API response")
			}

			protectStatus := utils.PathSearch("protect_status", clusterRespBody, "").(string)
			if protectStatus == "" {
				return nil, "ERROR", errors.New("protect_status is not found in API response")
			}

			if utils.StrSliceContains(monitorProtectionStatuses, protectStatus) {
				return clusterRespBody, "COMPLETED", nil
			}

			if protectStatus == "error_protect" {
				failReason := utils.PathSearch("fail_reason", clusterRespBody, "").(string)
				return clusterRespBody, "ERROR", fmt.Errorf("the protection status is `error_protect`: %s", failReason)
			}

			// Due to the lack of test conditions to verify the accuracy of the logic, other states are considered
			// to be states that need to be waited for.
			return clusterRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func resourceClusterProtectionEnableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		epsId     = cfg.GetEnterpriseProjectID(d)
		clusterID = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container/kubernetes/clusters/protection-enable"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateClusterProtectionBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error enabling HSS container kubernetes cluster protection: %s", err)
	}
	d.SetId(clusterID)

	configResp, err := waitingForProtectStatusProtecting(ctx, client, d, cfg, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for HSS container kubernetes cluster (%s) protection enable to complete: %s",
			clusterID, err)
	}

	if err := d.Set("protect_status", utils.PathSearch("protect_status", configResp, "").(string)); err != nil {
		return diag.Errorf("error setting protect status: %s", err)
	}

	return resourceClusterProtectionEnableRead(ctx, d, meta)
}

func resourceClusterProtectionEnableRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceClusterProtectionEnableUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceClusterProtectionEnableDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time action resource, so deletion is not supported
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary: `This resource is a one-time action resource used to enable HSS container cluster protection. Deleting this
			resource will not disable the protection, but will only remove the resource information from the tf state file.`,
		},
	}
}
