package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var volumeAutoExpandConfigurationNonUpdatableParams = []string{
	"instance_id",
	"auto_volume_expand_enable",
	"expand_threshold",
	"expand_increment",
	"max_volume_size",
}

// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/auto-volume-expand
func ResourceVolumeAutoExpandConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeAutoExpandConfigurationCreate,
		ReadContext:   resourceVolumeAutoExpandConfigurationRead,
		UpdateContext: resourceVolumeAutoExpandConfigurationUpdate,
		DeleteContext: resourceVolumeAutoExpandConfigurationDelete,

		CustomizeDiff: config.FlexibleForceNew(volumeAutoExpandConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Kafka instance is located to be configured volume auto-expansion.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"auto_volume_expand_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable disk auto-expansion.`,
			},
			"expand_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The threshold for triggering disk auto-expansion, in percentage (%).`,
			},
			"expand_increment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The percentage of storage space to be expanded out of the total instance storage space, in percentage (%).`,
			},
			"max_volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum volume size for disk auto-expansion, in GB.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildVolumeAutoExpandConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"auto_volume_expand_enable": d.Get("auto_volume_expand_enable"),
		"expand_threshold":          utils.ValueIgnoreEmpty(d.Get("expand_threshold")),
		"expand_increment":          utils.ValueIgnoreEmpty(d.Get("expand_increment")),
		"max_volume_size":           utils.ValueIgnoreEmpty(d.Get("max_volume_size")),
	}
}

func resourceVolumeAutoExpandConfigurationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/auto-volume-expand"
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVolumeAutoExpandConfigurationBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error configuring volume auto-expansion of the instance(%s): %s", instanceId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceVolumeAutoExpandConfigurationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVolumeAutoExpandConfigurationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVolumeAutoExpandConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to configure volume auto-expansion. 
Deleting this resource will not clear the configuration, but will only remove the resource information 
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
