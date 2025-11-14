package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/auto-volume-expand
func DataSourceVolumeAutoExpandConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVolumeAutoExpandConfigurationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the kafka instance is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"auto_volume_expand_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether disk auto-expansion is enabled.`,
			},
			"expand_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The threshold that triggers disk auto-expansion.`,
			},
			"max_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum volume size for disk auto-expansion.`,
			},
			"expand_increment": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The percentage of each disk auto-expansion.`,
			},
		},
	}
}

func dataSourceVolumeAutoExpandConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/auto-volume-expand"
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving volume auto-expansion configuration of Kafka instance (%s): %s", instanceId, err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("auto_volume_expand_enable", utils.PathSearch("auto_volume_expand_enable", respBody, nil)),
		d.Set("expand_threshold", utils.PathSearch("expand_threshold", respBody, nil)),
		d.Set("max_volume_size", utils.PathSearch("max_volume_size", respBody, nil)),
		d.Set("expand_increment", utils.PathSearch("expand_increment", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
