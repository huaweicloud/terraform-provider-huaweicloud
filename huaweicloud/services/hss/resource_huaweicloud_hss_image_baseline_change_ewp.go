package hss

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

// @API HSS POST /v5/{project_id}/image/baseline/extended-weak-password
func ResourceImageBaselineChangeEWP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageBaselineChangeEWPCreate,
		ReadContext:   resourceImageBaselineChangeEWPRead,
		UpdateContext: resourceImageBaselineChangeEWPUpdate,
		DeleteContext: resourceImageBaselineChangeEWPDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{"extended_weak_password"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// `extended_weak_password` can be set to an empty list.
			"extended_weak_password": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildImageBaselineChangeEWPBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"extended_weak_password": utils.ExpandToStringList(d.Get("extended_weak_password").([]interface{})),
	}
}

func resourceImageBaselineChangeEWPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		httpUrl = "v5/{project_id}/image/baseline/extended-weak-password"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildImageBaselineChangeEWPBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating HSS image baseline extended-weak-password: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceImageBaselineChangeEWPRead(ctx, d, meta)
}

func resourceImageBaselineChangeEWPRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageBaselineChangeEWPUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageBaselineChangeEWPDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This is a one-time action resource used to update the extended-weak-password configuration of the
    HSS image baseline. Deleting this resource will not clear the corresponding request record, but will only remove the
    resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
