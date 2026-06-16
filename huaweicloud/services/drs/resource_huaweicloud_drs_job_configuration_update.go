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

var jobConfigurationUpdateNonUpdatableParams = []string{
	"job_id",
	"values",
	"values.*.parameter_name",
	"values.*.parameter_value",
}

// @API DRS PUT /v5/{project_id}/jobs/{job_id}/modify-configuration
func ResourceJobConfigurationUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobConfigurationUpdateCreate,
		ReadContext:   resourceJobConfigurationUpdateRead,
		UpdateContext: resourceJobConfigurationUpdateUpdate,
		DeleteContext: resourceJobConfigurationUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(jobConfigurationUpdateNonUpdatableParams),

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
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func buildJobConfigurationUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("values").([]interface{})
	rstArray := make([]map[string]interface{}, 0, len(rawArray))

	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		rstArray = append(rstArray, map[string]interface{}{
			"parameter_name":  rawMap["parameter_name"],
			"parameter_value": rawMap["parameter_value"],
		})
	}

	return map[string]interface{}{
		"values": rstArray,
	}
}

func resourceJobConfigurationUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}/modify-configuration"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes:          []int{200, 201, 204},
		JSONBody:         buildJobConfigurationUpdateBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS job configuration: %s", err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	return nil
}

func resourceJobConfigurationUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobConfigurationUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobConfigurationUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update DRS job configuration. Deleting this resource
	will not rollback the configuration update operation, but will only remove the resource information from the tf
	state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
