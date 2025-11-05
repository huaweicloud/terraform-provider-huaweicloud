package dew

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceStatusNonUpdatableParams = []string{"instance_id"}

// @API DEW POST /v1/{project_id}/dew/cpcs/instances/{instance_id}/enable
// @API DEW POST /v1/{project_id}/dew/cpcs/instances/{instance_id}/disable
func ResourceCpcsInstanceStatusAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCpcsInstanceStatusActionCreate,
		ReadContext:   resourceCpcsInstanceStatusActionRead,
		UpdateContext: resourceCpcsInstanceStatusActionUpdate,
		DeleteContext: resourceCpcsInstanceStatusActionDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceStatusNonUpdatableParams),

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
			"action": {
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

func updateInstanceAction(client *golangsdk.ServiceClient, action, instanceId string) error {
	httpUrl := ""
	switch action {
	case "enable":
		httpUrl = "v1/{project_id}/dew/cpcs/instances/{instance_id}/enable"
	case "disable":
		httpUrl = "v1/{project_id}/dew/cpcs/instances/{instance_id}/disable"
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceCpcsInstanceStatusActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		product    = "kms"
		action     = d.Get("action").(string)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	err = updateInstanceAction(client, action, instanceId)
	if err != nil {
		return diag.Errorf("error %s DEW CPCS instance in creation operation: %s", action, err)
	}

	d.SetId(instanceId)

	return resourceCpcsInstanceStatusActionRead(ctx, d, meta)
}

func resourceCpcsInstanceStatusActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCpcsInstanceStatusActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		product    = "kms"
		action     = d.Get("action").(string)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	err = updateInstanceAction(client, action, instanceId)
	if err != nil {
		return diag.Errorf("error %s DEW CPCS instance in update operation: %s", action, err)
	}

	return resourceCpcsInstanceStatusActionRead(ctx, d, meta)
}

func resourceCpcsInstanceStatusActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not affect the actual instance status, but will only remove the resource
	information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
