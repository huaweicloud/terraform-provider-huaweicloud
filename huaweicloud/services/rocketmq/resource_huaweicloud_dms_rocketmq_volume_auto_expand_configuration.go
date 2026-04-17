package rocketmq

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

var volumeAutoExpandConfigurationNonUpdatableParams = []string{
	"instance_id",
	"auto_volume_expand_enable",
	"expand_threshold",
	"max_volume_size",
	"expand_increment",
}

// @API RocketMQ PUT /v2/{project_id}/instances/{instance_id}/auto-volume-expand
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/auto-volume-expand
func ResourceVolumeAutoExpandConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeAutoExpandConfigurationCreate,
		ReadContext:   resourceVolumeAutoExpandConfigurationRead,
		UpdateContext: resourceVolumeAutoExpandConfigurationUpdate,
		DeleteContext: resourceVolumeAutoExpandConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(volumeAutoExpandConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instance is located to be configured volume auto-expansion.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance.`,
			},

			// Optional parameters.
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
			"max_volume_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum volume size for disk auto-expansion, in GB.`,
			},
			"expand_increment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The percentage of storage space to be expanded out of the total instance storage space, in percentage (%).`,
			},

			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildVolumeAutoExpandConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"auto_volume_expand_enable": d.Get("auto_volume_expand_enable"),
		"expand_threshold":          utils.ValueIgnoreEmpty(d.Get("expand_threshold")),
		"max_volume_size":           utils.ValueIgnoreEmpty(d.Get("max_volume_size")),
		"expand_increment":          utils.ValueIgnoreEmpty(d.Get("expand_increment")),
	}
}

func updateVolumeAutoExpandConfiguration(client *golangsdk.ServiceClient, instanceId string, params map[string]interface{}) error {
	httpUrl := "v2/{project_id}/instances/{instance_id}/auto-volume-expand"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes:          []int{204},
		JSONBody:         utils.RemoveNil(params),
	}

	_, err := client.Request("PUT", updatePath, &opts)
	return err
}

func GetVolumeAutoExpandConfiguration(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/auto-volume-expand"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &opts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceVolumeAutoExpandConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	if err = updateVolumeAutoExpandConfiguration(client, instanceId, buildVolumeAutoExpandConfigurationBodyParams(d)); err != nil {
		return diag.Errorf("error configuring volume auto-expansion for the instance (%s): %s", instanceId, err)
	}

	d.SetId(instanceId)

	return resourceVolumeAutoExpandConfigurationRead(ctx, d, meta)
}

func resourceVolumeAutoExpandConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)

	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := GetVolumeAutoExpandConfiguration(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting volume auto-expansion configuration of the instance (%s): %s",
			instanceId, err))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("auto_volume_expand_enable", utils.PathSearch("auto_volume_expand_enable", respBody, nil)),
		d.Set("expand_threshold", utils.PathSearch("expand_threshold", respBody, nil)),
		d.Set("max_volume_size", utils.PathSearch("max_volume_size", respBody, nil)),
		d.Set("expand_increment", utils.PathSearch("expand_increment", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVolumeAutoExpandConfigurationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVolumeAutoExpandConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to configure volume auto-expansion. Deleting this resource will not
clear the configuration, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
