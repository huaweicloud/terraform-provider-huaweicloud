package bms

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

var instancePasswordDeleteNonUpdatableParams = []string{"server_id"}

// @API BMS DELETE /v1/{project_id}/baremetalservers/{server_id}/os-server-password
func ResourceInstancePasswordDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstancePasswordDeleteCreate,
		ReadContext:   resourceInstancePasswordDeleteRead,
		UpdateContext: resourceInstancePasswordDeleteUpdate,
		DeleteContext: resourceInstancePasswordDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(instancePasswordDeleteNonUpdatableParams),

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

func resourceInstancePasswordDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/os-server-password"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", d.Get("server_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating BMS instance password delete: %s", err)
	}

	d.SetId(d.Get("server_id").(string))

	return nil
}

func resourceInstancePasswordDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstancePasswordDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstancePasswordDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting BMS instance password delete resource is not supported. The resource is only removed from state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
