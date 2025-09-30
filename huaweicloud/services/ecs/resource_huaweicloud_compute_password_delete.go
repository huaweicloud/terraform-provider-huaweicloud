package ecs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var passwordDeleteNonUpdatableParams = []string{"server_id"}

// @API ECS DELETE /v1/{project_id}/cloudservers/{server_id}/os-server-password
func ResourceComputePasswordDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputePasswordDeleteCreate,
		ReadContext:   resourceComputePasswordDeleteRead,
		UpdateContext: resourceComputePasswordDeleteUpdate,
		DeleteContext: resourceComputePasswordDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(passwordDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceComputePasswordDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/cloudservers/{server_id}/os-server-password"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", d.Get("server_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS password delete: %s", err)
	}

	d.SetId(d.Get("server_id").(string))

	return nil
}

func resourceComputePasswordDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputePasswordDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputePasswordDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ECS password delete resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
