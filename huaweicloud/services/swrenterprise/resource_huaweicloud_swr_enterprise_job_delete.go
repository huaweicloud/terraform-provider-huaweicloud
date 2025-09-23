package swrenterprise

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

var enterpriseJobDeleteNonUpdatableParams = []string{
	"job_id",
}

// @API SWR DELETE /v2/{project_id}/jobs/{job_id}
func ResourceSwrEnterpriseJobDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseJobDeleteCreate,
		UpdateContext: resourceSwrEnterpriseJobDeleteUpdate,
		ReadContext:   resourceSwrEnterpriseJobDeleteRead,
		DeleteContext: resourceSwrEnterpriseJobDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseJobDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise job ID.`,
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

func resourceSwrEnterpriseJobDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/jobs/{job_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{job_id}", d.Get("job_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SWR enterprise job: %s", err)
	}

	d.SetId(d.Get("job_id").(string))

	return nil
}

func resourceSwrEnterpriseJobDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseJobDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseJobDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise job delete resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
