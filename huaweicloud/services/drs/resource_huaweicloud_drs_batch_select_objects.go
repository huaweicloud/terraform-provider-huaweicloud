package drs

import (
	"context"
	"fmt"
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

var batchSelectObjectsNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.selected",
	"jobs.*.sync_database",
	"jobs.*.job",
	"jobs.*.job.*.id",
	"jobs.*.job.*.parent_id",
	"jobs.*.job.*.object_type",
	"jobs.*.job.*.object_name",
	"jobs.*.job.*.select",
	"jobs.*.job.*.object_alias_name",
}

// @API DRS PUT /v3/{project_id}/jobs/batch-select-objects
func ResourceDrsBatchSelectObjects() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchSelectObjectsCreate,
		ReadContext:   resourceBatchSelectObjectsRead,
		UpdateContext: resourceBatchSelectObjectsUpdate,
		DeleteContext: resourceBatchSelectObjectsDelete,

		CustomizeDiff: config.FlexibleForceNew(batchSelectObjectsNonUpdatableParams),

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
				Elem:     batchSelectObjectsJobsSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"all_counts": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     batchSelectObjectsResultsSchema(),
			},
		},
	}
}

func batchSelectObjectsJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"selected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sync_database": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"job": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     batchSelectObjectsJobSchema(),
			},
		},
	}
}

func batchSelectObjectsJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_alias_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func batchSelectObjectsResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
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
	}
}

func buildBatchSelectObjectsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"jobs": buildBatchSelectObjectsJobsParams(d),
	}
}

func buildBatchSelectObjectsJobsParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("jobs").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawConfig := d.GetRawConfig()
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for i, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"job_id":   rawMap["job_id"],
			"selected": utils.GetNestedObjectFromRawConfig(rawConfig, fmt.Sprintf("jobs.%d.selected", i)),
			"sync_database": utils.GetNestedObjectFromRawConfig(
				rawConfig, fmt.Sprintf("jobs.%d.sync_database", i)),
			"job": buildBatchSelectObjectsJobParams(rawMap["job"].([]interface{})),
		})
	}

	return rst
}

func buildBatchSelectObjectsJobParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":                utils.ValueIgnoreEmpty(rawMap["id"]),
			"parent_id":         utils.ValueIgnoreEmpty(rawMap["parent_id"]),
			"object_type":       utils.ValueIgnoreEmpty(rawMap["object_type"]),
			"object_name":       utils.ValueIgnoreEmpty(rawMap["object_name"]),
			"select":            utils.ValueIgnoreEmpty(rawMap["select"]),
			"object_alias_name": utils.ValueIgnoreEmpty(rawMap["object_alias_name"]),
		})
	}

	return rst
}

func resourceBatchSelectObjectsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/batch-select-objects"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildBatchSelectObjectsBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch selecting DRS job objects: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("all_counts", utils.PathSearch("all_counts", respBody, nil)),
		d.Set("results", flattenBatchSelectObjectsResults(
			utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS batch select objects fields: %s", mErr)
	}

	return nil
}

func flattenBatchSelectObjectsResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"job_id":     utils.PathSearch("job_id", result, nil),
			"status":     utils.PathSearch("status", result, nil),
			"error_code": utils.PathSearch("error_code", result, nil),
			"error_msg":  utils.PathSearch("error_msg", result, nil),
		})
	}

	return rst
}

func resourceBatchSelectObjectsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBatchSelectObjectsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBatchSelectObjectsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch select DRS job objects. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information
    from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
