package drs

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

var objectCompareNonUpdatableParams = []string{
	"job_id",
	"compare_task_num",
}

// @API DRS POST /v3/{project_id}/jobs/{job_id}/object/compare
func ResourceObjectCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectCompareCreate,
		ReadContext:   resourceObjectCompareRead,
		UpdateContext: resourceObjectCompareUpdate,
		DeleteContext: resourceObjectCompareDelete,

		CustomizeDiff: config.FlexibleForceNew(objectCompareNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"compare_task_num": {
				Type:     schema.TypeInt,
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

func buildObjectCompareBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"compare_task_num": d.Get("compare_task_num"),
	}
}

func resourceObjectCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/{job_id}/object/compare"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildObjectCompareBodyParams(d),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error starting DRS object compare task: %s", err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	return nil
}

func resourceObjectCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceObjectCompareUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceObjectCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to start a DRS object compare task. Deleting this 
resource will not undo the compare operation, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
