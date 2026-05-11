package geminidb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var memoryMappingNonUpdatableParams = []string{"source_instance_id", "target_instance_id"}

// @API GaussDBforNoSQL POST /v3/{project_id}/dbcache/mapping
// @API GaussDBforNoSQL GET /v3/{project_id}/dbcache/mappings
// @API GaussDBforNoSQL DELETE /v3/{project_id}/dbcache/mapping
// @API GaussDBforNoSQL GET /v3/{project_id}/jobs
func ResourceMemoryMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemoryMappingCreate,
		ReadContext:   resourceMemoryMappingRead,
		UpdateContext: resourceMemoryMappingUpdate,
		DeleteContext: resourceMemoryMappingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(memoryMappingNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildMemoryMappingBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_instance_id": d.Get("source_instance_id"),
		"target_instance_id": d.Get("target_instance_id"),
	}

	return bodyParams
}

func resourceMemoryMappingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v3/{project_id}/dbcache/mapping"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildMemoryMappingBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating memory acceleration mapping: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating memory acceleration mapping: unable to find job ID")
	}

	queryParams := map[string]string{
		"source_instance_id": d.Get("source_instance_id").(string),
		"target_instance_id": d.Get("target_instance_id").(string),
	}

	memoryMapping, err := GetMemoryMappingInfo(client, queryParams)
	if err != nil {
		return diag.Errorf("error retrieving memory acceleration mapping information: %s", err)
	}

	mappingId := utils.PathSearch("id", memoryMapping, "").(string)
	if mappingId == "" {
		return diag.Errorf("error creating memory acceleration mapping: unable to find mapping ID")
	}

	d.SetId(mappingId)

	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating memory acceleration mapping: %s", err)
	}

	return resourceMemoryMappingRead(ctx, d, meta)
}

func resourceMemoryMappingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	queryParams := map[string]string{
		"id": d.Id(),
	}
	memoryMappingInfo, err := GetMemoryMappingInfo(client, queryParams)
	if err != nil {
		// When the memory mapping does not exist, the response HTTP status code of the query API is 200
		// and return nil
		return common.CheckDeletedDiag(d, err, "error retrieving memory acceleration mapping")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", memoryMappingInfo, nil)),
		d.Set("source_instance_id", utils.PathSearch("source_instance_id", memoryMappingInfo, nil)),
		d.Set("source_instance_name", utils.PathSearch("source_instance_name", memoryMappingInfo, nil)),
		d.Set("target_instance_id", utils.PathSearch("target_instance_id", memoryMappingInfo, nil)),
		d.Set("target_instance_name", utils.PathSearch("target_instance_name", memoryMappingInfo, nil)),
		d.Set("status", utils.PathSearch("status", memoryMappingInfo, nil)),
		d.Set("created", utils.PathSearch("created", memoryMappingInfo, nil)),
		d.Set("updated", utils.PathSearch("updated", memoryMappingInfo, nil)),
		d.Set("rule_count", utils.PathSearch("rule_count", memoryMappingInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetMemoryMappingInfo(client *golangsdk.ServiceClient, queryParams map[string]string) (interface{}, error) {
	httpUrl := "v3/{project_id}/dbcache/mappings"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if len(queryParams) > 0 {
		var queryParts []string
		for key, value := range queryParams {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
		}

		getPath += "?" + strings.Join(queryParts, "&")
	}

	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	mappingInfo := utils.PathSearch("dbcache_mappings|[0]", respBody, nil)
	if mappingInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return mappingInfo, nil
}

func resourceMemoryMappingUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMemoryMappingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dbcache/mapping"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"id": d.Id(),
		},
	}

	resp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the memory mapping does not exist, the response HTTP status code of the deletion API is 400.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.03050010"),
			fmt.Sprintf("error deleting memory acceleration mapping, the error message: %s", err))
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting memory acceleration mapping: unable to find job ID")
	}

	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error deleting memory acceleration mapping: %s", err)
	}

	return nil
}
