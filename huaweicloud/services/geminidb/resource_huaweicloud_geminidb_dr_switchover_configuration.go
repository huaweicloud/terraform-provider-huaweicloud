package geminidb

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

var geminiDbDRSwitchoverConfigurationParams = []string{
	"instance_id",
}

// @API GeminiDB PUT /v3/{project_id}/instances/disaster-recovery/settings
// @API GeminiDB GET /v3/{project_id}/instances/disaster-recovery/settings
func ResourceGeminiDBDRSwitchoverConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBDRSwitchoverConfigurationCreateOrUpdate,
		ReadContext:   resourceGeminiDBDRSwitchoverConfigurationRead,
		UpdateContext: resourceGeminiDBDRSwitchoverConfigurationCreateOrUpdate,
		DeleteContext: resourceGeminiDBDRSwitchoverConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDbDRSwitchoverConfigurationParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switchover_ratio": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sync_delay": {
				Type:     schema.TypeInt,
				Optional: true,
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

func resourceGeminiDBDRSwitchoverConfigurationCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/disaster-recovery/settings"
	createOrUpdatePath := client.Endpoint + httpUrl
	createOrUpdatePath = strings.ReplaceAll(createOrUpdatePath, "{project_id}", client.ProjectID)

	createOrUpdateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDisasterRecoverySettingsBody(d),
	}

	_, err = client.Request("PUT", createOrUpdatePath, &createOrUpdateOpt)
	if err != nil {
		return diag.Errorf("error setting disaster recovery settings: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	d.SetId(instanceID)

	return resourceGeminiDBDRSwitchoverConfigurationRead(ctx, d, meta)
}

func buildDisasterRecoverySettingsBody(d *schema.ResourceData) map[string]interface{} {
	instanceID := d.Get("instance_id").(string)

	settings := map[string]interface{}{
		"instance_id":      instanceID,
		"switchover_ratio": utils.ValueIgnoreEmpty(d.Get("switchover_ratio")),
		"sync_delay":       utils.ValueIgnoreEmpty(d.Get("sync_delay")),
	}

	body := map[string]interface{}{
		"disaster_recovery_settings": []map[string]interface{}{settings},
	}

	return body
}

func resourceGeminiDBDRSwitchoverConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Id()

	httpUrl := "v3/{project_id}/instances/disaster-recovery/settings?instance_id={instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving disaster recovery settings")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	settings := utils.PathSearch(fmt.Sprintf("disaster_recovery_settings[?instance_id=='%s']|[0]", instanceID), respBody, nil)
	if settings == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving disaster recovery settings")
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("switchover_ratio", utils.PathSearch("switchover_ratio", settings, nil)),
		d.Set("sync_delay", utils.PathSearch("sync_delay", settings, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBDRSwitchoverConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB DR Configuring DR Switchover for an Instance is not supported." +
		"The Configuring DR Switchover for an Instance resource is only removed from the state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
