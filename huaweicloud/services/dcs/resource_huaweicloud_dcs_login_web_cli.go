package dcs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dcsLoginWebCliNonUpdatableParams = []string{"instance_id", "password"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/webcli/auth
func ResourceDcsLoginWebCli() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsLoginWebCliCreate,
		ReadContext:   resourceDcsLoginWebCliRead,
		UpdateContext: resourceDcsLoginWebCliUpdate,
		DeleteContext: resourceDcsLoginWebCliDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsLoginWebCliNonUpdatableParams),

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
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceDcsLoginWebCliCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/webcli/auth"
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
	if v, ok := d.GetOk("password"); ok {
		reqBody["password"] = v.(string)
	}

	createOpt := golangsdk.RequestOpts{
		JSONBody: reqBody,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS login web cli resource: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing DCS login response: %s", err)
	}

	clientID := utils.PathSearch("client_id", respBody, "").(string)
	if clientID == "" {
		return diag.Errorf("client_id not found in login response")
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	if err := d.Set("client_id", clientID); err != nil {
		return diag.Errorf("failed to set client_id: %s", err)
	}

	return nil
}

func resourceDcsLoginWebCliRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsLoginWebCliUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsLoginWebCliDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS WebCli login resource is not supported. The resource is only removed from the state"
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
