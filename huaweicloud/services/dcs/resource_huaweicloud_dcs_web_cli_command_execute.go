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

var dcsWebCliCommandExecuteNonUpdatableParams = []string{"instance_id", "client_id", "command", "database"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/webcli/command
func ResourceDcsWebCliCommandExecute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsWebCliCommandExecuteCreate,
		ReadContext:   resourceDcsWebCliCommandExecuteRead,
		UpdateContext: resourceDcsWebCliCommandExecuteUpdate,
		DeleteContext: resourceDcsWebCliCommandExecuteDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsWebCliCommandExecuteNonUpdatableParams),

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
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
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

func resourceDcsWebCliCommandExecuteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/webcli/command"
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

	createOpt := golangsdk.RequestOpts{
		JSONBody:         utils.RemoveNil(buildCreateDcsWebcliCommandBodyParams(d)),
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS web cli command execute resource: %s", err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	return nil
}

func buildCreateDcsWebcliCommandBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_id": d.Get("client_id"),
		"command":   d.Get("command"),
		"database":  d.Get("database"),
	}

	return bodyParams
}

func resourceDcsWebCliCommandExecuteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsWebCliCommandExecuteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsWebCliCommandExecuteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS web cli command execute resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
