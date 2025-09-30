package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var twoFactorLoginConfigNonUpdatableParams = []string{
	"enabled",
	"auth_type",
	"host_id_list",
	"topic_display_name",
	"topic_urn",
	"enterprise_project_id"}

// @API HSS POST /v5/{project_id}/setting/two-factor-login/config
func ResourceSettingTwoFactorLoginConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSettingTwoFactorLoginConfigCreate,
		ReadContext:   resourceSettingTwoFactorLoginConfigRead,
		UpdateContext: resourceSettingTwoFactorLoginConfigUpdate,
		DeleteContext: resourceSettingTwoFactorLoginConfigDelete,

		CustomizeDiff: config.FlexibleForceNew(twoFactorLoginConfigNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topic_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
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

func buildTwoFactorLoginConfigQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildTwoFactorLoginConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enabled":            d.Get("enabled"),
		"auth_type":          d.Get("auth_type"),
		"host_id_list":       utils.ExpandToStringList(d.Get("host_id_list").([]interface{})),
		"topic_display_name": utils.ValueIgnoreEmpty(d.Get("topic_display_name")),
		"topic_urn":          utils.ValueIgnoreEmpty(d.Get("topic_urn")),
	}

	return bodyParams
}

func resourceSettingTwoFactorLoginConfigCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/setting/two-factor-login/config"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildTwoFactorLoginConfigQueryParams(epsId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTwoFactorLoginConfigBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error setting the two-factor login configuration: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func resourceSettingTwoFactorLoginConfigRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceSettingTwoFactorLoginConfigUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceSettingTwoFactorLoginConfigDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to set two-factor login configuration. Deleting this resource
    will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
