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

// Due to a lack of testing conditions, this resource's documentation is not publicly available.

// @API DEW POST /v1.0/{project_id}/keystores/{keystore_id}/disable
// @API DEW POST /v1.0/{project_id}/keystores/{keystore_id}/enable
func ResourceKmsDedicatedKeystoreAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsDedicatedKeystoreActionCreate,
		ReadContext:   resourceKmsDedicatedKeystoreActionRead,
		UpdateContext: resourceKmsDedicatedKeystoreActionUpdate,
		DeleteContext: resourceKmsDedicatedKeystoreActionDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"keystore_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"keystore_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a dedicated keystore.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action to be performed on the dedicated keystore.`,
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

func updateDedicatedKeystoreStatus(client *golangsdk.ServiceClient, action string, keystoreId string) error {
	httpUrl := ""
	switch action {
	case "enable":
		httpUrl = "v1.0/{project_id}/keystores/{keystore_id}/enable"
	case "disable":
		httpUrl = "v1.0/{project_id}/keystores/{keystore_id}/disable"
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{keystore_id}", keystoreId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceKmsDedicatedKeystoreActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "kms"
		action     = d.Get("action").(string)
		keystoreId = d.Get("keystore_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	err = updateDedicatedKeystoreStatus(client, action, keystoreId)
	if err != nil {
		return diag.Errorf("error %s DEW dedicated keystore in creation operation: %s", action, err)
	}

	d.SetId(keystoreId)

	return resourceKmsDedicatedKeystoreActionRead(ctx, d, meta)
}

func resourceKmsDedicatedKeystoreActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsDedicatedKeystoreActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "kms"
		action     = d.Get("action").(string)
		keystoreId = d.Get("keystore_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	err = updateDedicatedKeystoreStatus(client, action, keystoreId)
	if err != nil {
		return diag.Errorf("error %s DEW dedicated keystore in update operation: %s", action, err)
	}

	return resourceKmsDedicatedKeystoreActionRead(ctx, d, meta)
}

func resourceKmsDedicatedKeystoreActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to perform actions on a dedicated keystore.
Deleting this resource will not recover the dedicated keystore, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
