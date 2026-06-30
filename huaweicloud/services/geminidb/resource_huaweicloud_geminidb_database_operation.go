package geminidb

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var geminiDbDatabaseOperationNonUpdatableParams = []string{
	"instance_id",
	"action",
	"db_id",
}

// @API GeminiDB PUT /v3/{project_id}/instances/{instance_id}/databases
func ResourceGeminiDBDatabaseOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBDatabaseOperationCreate,
		ReadContext:   resourceGeminiDBDatabaseOperationRead,
		UpdateContext: resourceGeminiDBDatabaseOperationUpdate,
		DeleteContext: resourceGeminiDBDatabaseOperationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDbDatabaseOperationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_id": {
				Type:     schema.TypeInt,
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

func resourceGeminiDBDatabaseOperationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGeminiDBDatabaseOperationBodyParams(d))

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating GeminiDB database: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error operating GeminiDB database: job_id is not found in API response")
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error operating GeminiDB database: %s", err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId.String())

	return nil
}

func resourceGeminiDBDatabaseOperationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBDatabaseOperationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBDatabaseOperationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB database operation resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildCreateGeminiDBDatabaseOperationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action": d.Get("action"),
	}

	if dbID, ok := d.GetOk("db_id"); ok {
		bodyParams["db_id"] = dbID
	}

	return bodyParams
}
