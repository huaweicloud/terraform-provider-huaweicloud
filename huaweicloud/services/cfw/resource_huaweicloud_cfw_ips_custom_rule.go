package cfw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW POST /v1/{project_id}/ips/custom-rule
// @API CFW PUT /v1/{project_id}/ips/custom-rule/{ips_cfw_id}
// @API CFW GET /v1/{project_id}/ips/custom-rule
// @API CFW GET /v1/{project_id}/ips/custom-rule/{ips_cfw_id}
// @API CFW POST /v1/{project_id}/ips/custom-rule/batch-delete
func ResourceIpsCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpsCustomRuleCreate,
		ReadContext:   resourceIpsCustomRuleRead,
		UpdateContext: resourceIpsCustomRuleUpdate,
		DeleteContext: resourceIpsCustomRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpsCustomRuleImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"fw_instance_id",
			"ips_name",
			"protocol",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			// Unable to retrieve the value of field `fw_instance_id` from the details API.
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ips_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"action_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"affected_os": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"attack_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"contents": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     contentsSchema(),
			},
			"direction": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"dst_port": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     portSchema(),
			},
			"severity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"software": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"src_port": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     portSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"config_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func contentsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Required: true}),
			},
			// Between 1 ~ 65535  Required
			"depth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Required: true}),
			},
			// Defaults to false.
			"is_hex": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Defaults to false.
			"is_ignore": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Defaults to false.
			"is_uri": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// Between 0 ~ 65535
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// Between 0 ~ 1
			"relative_position": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func portSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"port_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Required: true}),
			},
			"ports": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildIpsCustomRuleContentsParam(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"content":           rawMap["content"],
			"depth":             rawMap["depth"],
			"is_hex":            rawMap["is_hex"],
			"is_ignore":         rawMap["is_ignore"],
			"is_uri":            rawMap["is_uri"],
			"offset":            rawMap["offset"],
			"relative_position": rawMap["relative_position"],
		})
	}

	return rst
}

func buildIpsCustomRulePortParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"port_type": rawMap["port_type"],
		"ports":     utils.ValueIgnoreEmpty(rawMap["ports"]),
	}
}

func buildCreateIpsCustomRuleBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"ips_name":       d.Get("ips_name"),
		"protocol":       d.Get("protocol"),
		"action_type":    d.Get("action_type"),
		"affected_os":    d.Get("affected_os"),
		"attack_type":    d.Get("attack_type"),
		"contents":       buildIpsCustomRuleContentsParam(d.Get("contents").([]interface{})),
		"direction":      d.Get("direction"),
		"dst_port":       buildIpsCustomRulePortParam(d.Get("dst_port").([]interface{})),
		"severity":       d.Get("severity"),
		"software":       d.Get("software"),
		"src_port":       buildIpsCustomRulePortParam(d.Get("src_port").([]interface{})),
	}
}

func buildListIpsCustomRuleQueryParams(d *schema.ResourceData) string {
	var (
		ipsName      = d.Get("ips_name").(string)
		fwInstanceId = d.Get("fw_instance_id").(string)
		limit        = 1024
	)

	return fmt.Sprintf("?ips_name=%s&fw_instance_id=%s&limit=%d", ipsName, fwInstanceId, limit)
}

func getIpsCustomRuleByList(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/ips/custom-rule"
		ipsName = d.Get("ips_name").(string)
		offset  = 0
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildListIpsCustomRuleQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := requestPath + fmt.Sprintf("&offset=%d", offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		// The `ips_name` filter parameter in the API is a fuzzy search; we need to perform a precise match afterwards.
		targetRule := utils.PathSearch(fmt.Sprintf("[?ips_name=='%s']|[0]", ipsName), records, nil)
		if targetRule != nil {
			return targetRule, nil
		}

		// If the target data is not found on the current page, continue the search from the next offset.
		offset += len(records)
	}

	return nil, golangsdk.ErrDefault404{}
}

func buildIpsCustomRuleDetailQueryParams(fwInstanceId string) string {
	return fmt.Sprintf("?fw_instance_id=%s", fwInstanceId)
}

func GetIpsCustomRuleDetail(client *golangsdk.ServiceClient, fwInstanceId, ipsCfwId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/ips/custom-rule/{ips_cfw_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{ips_cfw_id}", ipsCfwId)
	requestPath += buildIpsCustomRuleDetailQueryParams(fwInstanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	contents := utils.PathSearch("data.contents", respBody, make([]interface{}, 0)).([]interface{})
	protocol := int(utils.PathSearch("data.protocol", respBody, float64(0)).(float64))
	severity := int(utils.PathSearch("data.severity", respBody, float64(0)).(float64))

	if len(contents) == 00 || protocol == -99 || severity == -99 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func refreshIpsCustomRuleConfigStatusFunc(client *golangsdk.ServiceClient, fwInstanceId, ipsCfwId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetIpsCustomRuleDetail(client, fwInstanceId, ipsCfwId)
		if err != nil {
			return nil, "ERROR", err
		}

		configStatus := int(utils.PathSearch("data.config_status", respBody, float64(0)).(float64))
		if configStatus == 2 {
			return respBody, "COMPLETED", nil
		}

		if configStatus == 3 {
			return respBody, "ERROR", errors.New("unexpected config status `3` detected")
		}

		return respBody, "PENDING", nil
	}
}

func waitingForIpsCustomRuleConfigSuccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	fwInstanceId := d.Get("fw_instance_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshIpsCustomRuleConfigStatusFunc(client, fwInstanceId, d.Id()),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceIpsCustomRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ips/custom-rule"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateIpsCustomRuleBodyParam(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CFW IPS custom rule: %s", err)
	}

	customRule, err := getIpsCustomRuleByList(client, d)
	if err != nil {
		return diag.Errorf("error querying CFW IPS custom rule after creation: %s", err)
	}

	id := utils.PathSearch("ips_cfw_id", customRule, "").(string)
	if id == "" {
		return diag.Errorf("error creating CFW IPS custom rule: ID is not found in API response")
	}

	d.SetId(id)

	if err := waitingForIpsCustomRuleConfigSuccess(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CFW IPS custom rule (%s) creation to config success: %s", d.Id(), err)
	}

	return resourceIpsCustomRuleRead(ctx, d, meta)
}

func flattenIpsCustomRuleContentsAttr(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"content":           utils.PathSearch("content", v, nil),
			"depth":             utils.PathSearch("depth", v, nil),
			"is_hex":            utils.PathSearch("is_hex", v, nil),
			"is_ignore":         utils.PathSearch("is_ignore", v, nil),
			"is_uri":            utils.PathSearch("is_uri", v, nil),
			"offset":            utils.PathSearch("offset", v, nil),
			"relative_position": utils.PathSearch("relative_position", v, nil),
		})
	}
	return rst
}

func flattenIpsCustomRuledPortAttr(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	portMap := map[string]interface{}{
		"port_type": utils.PathSearch("port_type", respBody, nil),
		"ports":     utils.PathSearch("ports", respBody, nil),
	}

	return []map[string]interface{}{portMap}
}

func resourceIpsCustomRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "cfw"
		fwInstanceId = d.Get("fw_instance_id").(string)
		ipsCfwId     = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	respBody, err := GetIpsCustomRuleDetail(client, fwInstanceId, ipsCfwId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW IPS custom rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ips_name", utils.PathSearch("data.ips_name", respBody, nil)),
		d.Set("protocol", utils.PathSearch("data.protocol", respBody, nil)),
		d.Set("action_type", utils.PathSearch("data.action", respBody, nil)),
		d.Set("affected_os", utils.PathSearch("data.affected_os", respBody, nil)),
		d.Set("attack_type", utils.PathSearch("data.attack_type", respBody, nil)),
		d.Set("contents", flattenIpsCustomRuleContentsAttr(utils.PathSearch("data.contents", respBody, nil))),
		d.Set("direction", utils.PathSearch("data.direction", respBody, nil)),
		d.Set("dst_port", flattenIpsCustomRuledPortAttr(utils.PathSearch("data.dst_port", respBody, nil))),
		d.Set("severity", utils.PathSearch("data.severity", respBody, nil)),
		d.Set("software", utils.PathSearch("data.software", respBody, nil)),
		d.Set("src_port", flattenIpsCustomRuledPortAttr(utils.PathSearch("data.src_port", respBody, nil))),
		d.Set("config_status", utils.PathSearch("data.config_status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateIpsCustomRuleBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"ips_name":       d.Get("ips_name"),
		"protocol":       d.Get("protocol"),
		"action_type":    d.Get("action_type"),
		"affected_os":    d.Get("affected_os"),
		"attack_type":    d.Get("attack_type"),
		"contents":       buildIpsCustomRuleContentsParam(d.Get("contents").([]interface{})),
		"direction":      d.Get("direction"),
		"dst_port":       buildIpsCustomRulePortParam(d.Get("dst_port").([]interface{})),
		"severity":       d.Get("severity"),
		"software":       d.Get("software"),
		"src_port":       buildIpsCustomRulePortParam(d.Get("src_port").([]interface{})),
	}
}

func handleRetryUpdateOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode := utils.PathSearch("error_code", apiError, "").(string)
		if errorCode == "" {
			return false, errors.New("unable to find error code from API response")
		}

		errorMsg := utils.PathSearch("error_msg", apiError, "").(string)
		if errorMsg == "" {
			return false, errors.New("unable to find error message from API response")
		}

		if errorCode == "CFW.00109003" {
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}

func resourceIpsCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ips/custom-rule/{ips_cfw_id}"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{ips_cfw_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateIpsCustomRuleBodyParam(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = client.Request("PUT", requestPath, &requestOpt)
		retry, err := handleRetryUpdateOperationsError(err)
		return nil, retry, err
	}

	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating CFW IPS custom rule: %s", err)
	}

	if err := waitingForIpsCustomRuleConfigSuccess(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("error waiting for CFW IPS custom rule (%s) update to config success: %s", d.Id(), err)
	}

	return resourceIpsCustomRuleRead(ctx, d, meta)
}

func buildDeleteIpsCustomRuleBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"ips_ids":        []string{d.Id()},
	}
}

func waitingForIpsCustomRuleDeleteSuccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	fwInstanceId := d.Get("fw_instance_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetIpsCustomRuleDetail(client, fwInstanceId, d.Id())
			if err != nil {
				var errCode404 golangsdk.ErrDefault404
				if errors.As(err, &errCode404) {
					return "deleted", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceIpsCustomRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ips/custom-rule/batch-delete"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDeleteIpsCustomRuleBodyParam(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CFW IPS custom rule: %s", err)
	}

	if err := waitingForIpsCustomRuleDeleteSuccess(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for CFW IPS custom rule (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func resourceIpsCustomRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <fw_instance_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("fw_instance_id", parts[0])
}
