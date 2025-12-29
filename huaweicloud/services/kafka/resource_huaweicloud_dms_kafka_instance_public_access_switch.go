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

var publicAccessSwitchNonUpdatableParams = []string{
	"instance_id",
	"eip_address",
	"public_boundwidth",
	"publicip_id",
}

// @API Kafka POST /v1/{project_id}/instances/{instance_id}/public-boundwidth
func ResourcePublicAccessSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicAccessSwitchCreate,
		ReadContext:   resourcePublicAccessSwitchRead,
		UpdateContext: resourcePublicAccessSwitchUpdate,
		DeleteContext: resourcePublicAccessSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(publicAccessSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the public access switch are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},

			// Optional parameters.
			"eip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The elastic IP address of the Kafka instance.`,
			},
			"public_boundwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  `The public bandwidth of the Kafka instance.`,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The public IP ID of the Kafka instance.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildPublicAccessSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"eip_address":       d.Get("eip_address"),
		"public_boundwidth": d.Get("public_boundwidth"),
		"publicip_id":       d.Get("publicip_id"),
	}
}

func resourcePublicAccessSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	var (
		httpUrl    = "v1/{project_id}/instances/{instance_id}/public-boundwidth"
		instanceId = d.Get("instance_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildPublicAccessSwitchBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error modifying Kafka public IP access switch: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourcePublicAccessSwitchRead(ctx, d, meta)
}

func resourcePublicAccessSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePublicAccessSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePublicAccessSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for switching public IP access of Kafka instance. 
Deleting this resource will not clear the corresponding request record, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
