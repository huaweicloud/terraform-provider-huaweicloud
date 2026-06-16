package drs

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

var batchRetryTaskNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.is_sync_re_edit",
}

// @API DRS POST /v3/{project_id}/jobs/batch-retry-task
func ResourceBatchRetryTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchRetryTaskCreate,
		ReadContext:   resourceBatchRetryTaskRead,
		UpdateContext: resourceBatchRetryTaskUpdate,
		DeleteContext: resourceBatchRetryTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(batchRetryTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						// Define the field type as string to handle default values ​​for boolean types.
						"is_sync_re_edit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
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
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func convertStringToBool(stringValue string) interface{} {
	if stringValue == "" {
		return nil
	}

	if stringValue == "true" {
		return true
	}

	return false
}

func buildRetryTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("jobs").([]interface{})
	rst := make([]map[string]interface{}, 0, len(rawArray))

	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"job_id":          rawMap["job_id"],
			"is_sync_re_edit": convertStringToBool(rawMap["is_sync_re_edit"].(string)),
		})
	}

	return map[string]interface{}{
		"jobs": rst,
	}
}

func resourceBatchRetryTaskCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/batch-retry-task"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildRetryTaskBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch retrying DRS tasks: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) == 0 {
		return diag.Errorf("unable to find the results from the API response")
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("results", flattenRetryTaskResults(results)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting retry DRS task fields: %s", mErr)
	}

	return nil
}

func flattenRetryTaskResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", result, nil),
			"status":     utils.PathSearch("status", result, nil),
			"error_code": utils.PathSearch("error_code", result, nil),
			"error_msg":  utils.PathSearch("error_msg", result, nil),
		})
	}
	return rst
}

func resourceBatchRetryTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchRetryTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchRetryTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to retry DRS tasks. Deleting this resource will not
undo the retry operation, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
