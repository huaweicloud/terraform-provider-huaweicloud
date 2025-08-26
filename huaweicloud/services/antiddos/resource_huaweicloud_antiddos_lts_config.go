package antiddos

import (
	"context"
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

// @API ANTI-DDOS PUT /v1/{project_id}/antiddos/lts-config
// @API ANTI-DDOS GET /v1/{project_id}/antiddos/lts-config
func ResourceLtsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLtsConfigCreate,
		ReadContext:   resourceLtsConfigRead,
		UpdateContext: resourceLtsConfigUpdate,
		DeleteContext: resourceLtsConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLtsConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"enterprise_project_id"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			"lts_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the LTS group ID.`,
			},
			"lts_attack_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the LTS attack stream ID.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildUpdateLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"enabled": true,
		"lts_id_info": map[string]interface{}{
			"lts_group_id":         d.Get("lts_group_id"),
			"lts_attack_stream_id": d.Get("lts_attack_stream_id"),
		},
	}
}

func buildLtcConfigQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	return fmt.Sprintf("?enterprise_project_id=%s", cfg.GetEnterpriseProjectID(d))
}

func updateLtsConfig(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	requestPath := client.Endpoint + "v1/{project_id}/antiddos/lts-config"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLtcConfigQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         buildUpdateLtsConfigBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceLtsConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
		epsID   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	if err := updateLtsConfig(client, d, cfg); err != nil {
		return diag.Errorf("error updating Anti-DDoS LTS config in create operation: %s", err)
	}

	d.SetId(epsID)

	return resourceLtsConfigRead(ctx, d, meta)
}

func resourceLtsConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
		httpUrl = "v1/{project_id}/antiddos/lts-config"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLtcConfigQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS LTS config: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := utils.PathSearch("enabled", respBody, false).(bool)
	if !enabled {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("lts_group_id", utils.PathSearch("lts_id_info.lts_group_id", respBody, nil)),
		d.Set("lts_attack_stream_id", utils.PathSearch("lts_id_info.lts_attack_stream_id", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLtsConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	if err := updateLtsConfig(client, d, cfg); err != nil {
		return diag.Errorf("error updating Anti-DDoS LTS config in update operation: %s", err)
	}

	return resourceLtsConfigRead(ctx, d, meta)
}

func deleteLtsConfig(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	requestPath := client.Endpoint + "v1/{project_id}/antiddos/lts-config"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLtcConfigQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"enabled": false,
		},
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceLtsConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	if err := deleteLtsConfig(client, d, cfg); err != nil {
		return diag.Errorf("error disabling Anti-DDoS LTS config: %s", err)
	}
	return nil
}

func resourceLtsConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", d.Id())
}
