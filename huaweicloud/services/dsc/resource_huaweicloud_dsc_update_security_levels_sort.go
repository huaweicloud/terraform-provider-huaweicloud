package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var updateSecurityLevelsSortNonUpdatableParams = []string{
	"level_id",
	"target_level_id",
}

// @API DSC PUT /v1/{project_id}/scan-security-levels/{level_id}/sort
func ResourceDscUpdateSecurityLevelsSort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscUpdateSecurityLevelsSortCreate,
		ReadContext:   resourceDscUpdateSecurityLevelsSortRead,
		UpdateContext: resourceDscUpdateSecurityLevelsSortUpdate,
		DeleteContext: resourceDscUpdateSecurityLevelsSortDelete,

		CustomizeDiff: config.FlexibleForceNew(updateSecurityLevelsSortNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"level_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_level_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDscUpdateSecurityLevelsSortCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-security-levels/{level_id}/sort"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{level_id}", d.Get("level_id").(string))

	body := map[string]interface{}{
		"target_level_id": d.Get("target_level_id").(string),
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         body,
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DSC security levels sort: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDscUpdateSecurityLevelsSortRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDscUpdateSecurityLevelsSortUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDscUpdateSecurityLevelsSortDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to adjust the display order of DSC security
levels. Deleting this resource will not revert the sort adjustment, but will only remove the resource
information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
