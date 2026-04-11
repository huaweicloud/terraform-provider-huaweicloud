package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var configurationResetNonUpdatableParams = []string{
	"config_id",
}

// @API RDS PUT /v3/{project_id}/configurations/{config_id}/reset
func ResourceRdsConfigurationReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsConfigurationResetCreate,
		ReadContext:   resourceRdsConfigurationResetRead,
		UpdateContext: resourceRdsConfigurationResetUpdate,
		DeleteContext: resourceRdsConfigurationResetDelete,

		CustomizeDiff: config.FlexibleForceNew(configurationResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"config_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"need_restart": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceRdsConfigurationResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}/reset"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{config_id}", d.Get("config_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS configuration reset: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("config_id").(string))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("config_name", utils.PathSearch("config_name", createRespBody, nil)),
		d.Set("need_restart", utils.PathSearch("need_restart", createRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRdsConfigurationResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsConfigurationResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsConfigurationResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS parameter template reset resource is not supported. The resource is only removed from the" +
		"state, the parameter template still remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
