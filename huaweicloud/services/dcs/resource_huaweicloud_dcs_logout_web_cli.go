package dcs

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

var dcsLogoutWebCliNonUpdatableParams = []string{"instance_id", "client_id"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/webcli/logout
func ResourceDcsLogoutWebCli() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsLogoutWebCliCreate,
		ReadContext:   resourceDcsLogoutWebCliRead,
		UpdateContext: resourceDcsLogoutWebCliUpdate,
		DeleteContext: resourceDcsLogoutWebCliDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsLogoutWebCliNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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

func resourceDcsLogoutWebCliCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/webcli/logout"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	reqBody := make(map[string]interface{})
	if v, ok := d.GetOk("client_id"); ok {
		reqBody["client_id"] = v.(string)
	}

	createOpt := golangsdk.RequestOpts{
		JSONBody: reqBody,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{
			200, 204,
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS logout web cli resource: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func resourceDcsLogoutWebCliRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsLogoutWebCliUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsLogoutWebCliDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS WebCli logout resource is not supported. The resource is only removed from the state"
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
