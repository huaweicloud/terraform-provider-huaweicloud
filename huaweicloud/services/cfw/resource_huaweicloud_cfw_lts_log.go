package cfw

import (
	"context"
	"encoding/json"
	"errors"
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

// @API CFW POST /v1/{project_id}/cfw/logs/configuration
// @API CFW GET /v1/{project_id}/cfw/logs/configuration
// @API CFW PUT /v1/{project_id}/cfw/logs/configuration
func ResourceLtsLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsLogCreate,
		ReadContext:   resourceLtsLogRead,
		UpdateContext: resourceLtsLogUpdate,
		DeleteContext: resourceLtsLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the firewall.`,
			},
			"lts_log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `LTS log group ID.`,
			},
			"lts_attack_log_stream_enable": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `LTS attack log stream switch.`,
			},
			"lts_access_log_stream_enable": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `LTS access log stream switch.`,
			},
			"lts_flow_log_stream_enable": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `LTS flow log stream switch.`,
			},
			"lts_attack_log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `LTS attack log stream ID.`,
			},
			"lts_access_log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `LTS access log stream ID.`,
			},
			"lts_flow_log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `LTS flow log stream ID.`,
			},
		},
	}
}

func resourceLtsLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/cfw/logs/configuration"
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id"))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = buildLtsLogConfigurationBodyParams(d)
	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error creating CFW lts log configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CFW lts log configuration: ID is not found in API response")
	}
	d.SetId(id)

	return resourceLtsLogRead(ctx, d, meta)
}

func parseError(err error) error {
	var errCode golangsdk.ErrDefault400
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}
		errorCode := utils.PathSearch("error_code", apiError, nil)
		if errorCode == nil {
			return fmt.Errorf("error parsing error_code from response")
		}
		if errorCode == "CFW.00200005" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func buildLtsLogConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	// A value of 1 for lts_enable indicates connection to the LTS service.
	return map[string]interface{}{
		"fw_instance_id":               d.Get("fw_instance_id"),
		"lts_enable":                   1,
		"lts_log_group_id":             d.Get("lts_log_group_id"),
		"lts_attack_log_stream_enable": d.Get("lts_attack_log_stream_enable"),
		"lts_access_log_stream_enable": d.Get("lts_access_log_stream_enable"),
		"lts_flow_log_stream_enable":   d.Get("lts_flow_log_stream_enable"),
		"lts_attack_log_stream_id":     d.Get("lts_attack_log_stream_id"),
		"lts_access_log_stream_id":     d.Get("lts_access_log_stream_id"),
		"lts_flow_log_stream_id":       d.Get("lts_flow_log_stream_id"),
	}
}

func resourceLtsLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	product := "cfw"
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	id := d.Id()
	configuration, err := getLtsLog(client, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW lts log configuration")
	}

	ltsEnable := utils.PathSearch("lts_enable", configuration, float64(0)).(float64)
	if int(ltsEnable) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CFW lts log configuration")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("fw_instance_id", utils.PathSearch("fw_instance_id", configuration, nil)),
		d.Set("lts_log_group_id", utils.PathSearch("lts_log_group_id", configuration, nil)),
		d.Set("lts_attack_log_stream_id", utils.PathSearch("lts_attack_log_stream_id", configuration, nil)),
		d.Set("lts_access_log_stream_id", utils.PathSearch("lts_access_log_stream_id", configuration, nil)),
		d.Set("lts_flow_log_stream_id", utils.PathSearch("lts_flow_log_stream_id", configuration, nil)),
		d.Set("lts_attack_log_stream_enable", utils.PathSearch("lts_attack_log_stream_enable", configuration, nil)),
		d.Set("lts_access_log_stream_enable", utils.PathSearch("lts_access_log_stream_enable", configuration, nil)),
		d.Set("lts_flow_log_stream_enable", utils.PathSearch("lts_flow_log_stream_enable", configuration, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getLtsLog(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cfw/logs/configuration"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", id)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	configuration := utils.PathSearch("data", respBody, nil)
	if configuration == nil {
		return nil, fmt.Errorf("error parsing data from response= %#v", respBody)
	}

	return configuration, nil
}

func resourceLtsLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/cfw/logs/configuration"
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id"))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = buildLtsLogConfigurationBodyParams(d)
	_, err = client.Request("PUT", path, &opt)
	if err != nil {
		return diag.Errorf("error updating CFW lts log configuration: %s", err)
	}

	return resourceLtsLogRead(ctx, d, meta)
}

func resourceLtsLogDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/cfw/logs/configuration"
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	id := d.Id()
	configuration, err := getLtsLog(client, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW lts log configuration")
	}

	ltsEnable := utils.PathSearch("lts_enable", configuration, float64(0)).(float64)
	if int(ltsEnable) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CFW lts log configuration")
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", d.Get("fw_instance_id"))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = buildDeleteLtsLogConfigurationBodyParams(d)
	_, err = client.Request("PUT", path, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW lts log configuration",
		)
	}

	return nil
}

func buildDeleteLtsLogConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	// A value of 0 for lts_enable indicates disconnection from the LTS service.
	// A value of 0 for lts_attack_log_stream_enable, lts_access_log_stream_enable, and lts_flow_log_stream_enable
	// indicates that the corresponding log stream is disabled.
	return map[string]interface{}{
		"fw_instance_id":               d.Get("fw_instance_id"),
		"lts_enable":                   0,
		"lts_log_group_id":             d.Get("lts_log_group_id"),
		"lts_attack_log_stream_enable": 0,
		"lts_access_log_stream_enable": 0,
		"lts_flow_log_stream_enable":   0,
		"lts_attack_log_stream_id":     "",
		"lts_access_log_stream_id":     "",
		"lts_flow_log_stream_id":       "",
	}
}
