package cnad

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD PUT /v1/cnad/protected-ips/tags
func ResourceProtectedIpTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectedIpTagCreate,
		ReadContext:   resourceProtectedIpTagRead,
		UpdateContext: resourceProtectedIpTagUpdate,
		DeleteContext: resourceProtectedIpTagDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{"protected_ip_id", "tag"}),

		Schema: map[string]*schema.Schema{
			"protected_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected IP ID.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the tag to be set on the protected IP.`,
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

func buildProtectedIpTagBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":  d.Get("protected_ip_id"),
		"tag": d.Get("tag"),
	}
}

func resourceProtectedIpTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/cnad/protected-ips/tags"
		product = "aad"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildProtectedIpTagBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error setting CNAD protected IP tag: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceProtectedIpTagRead(ctx, d, meta)
}

func resourceProtectedIpTagRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedIpTagUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedIpTagDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to set tag on protected IP. Deleting this 
resource will not change the current CNAD protected IP tag, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
