package rds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dbLogsShrinkingNonUpdatableParams = []string{"instance_id", "db_name"}

// @API RDS POST /v3.1/{project_id}/instances/{instance_id}/db-shrink
// @API RDS GET /v3/{project_id}/instances
func ResourceRdsDbLogsShrinking() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsDbLogsShrinkingCreate,
		ReadContext:   resourceRdsDbLogsShrinkingRead,
		UpdateContext: resourceRdsDbLogsShrinkingUpdate,
		DeleteContext: resourceRdsDbLogsShrinkingDelete,

		CustomizeDiff: config.FlexibleForceNew(dbLogsShrinkingNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the database.`,
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

func resourceRdsDbLogsShrinkingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-shrink"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDbLogsShrinkingBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error shrinking database logs of RDS instance(%s): %s", instanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	resp := utils.PathSearch("resp", createRespBody, nil)
	if resp == nil {
		return diag.Errorf("unable to find the resp from the response: %s", err)
	}
	if resp.(string) != "successful" {
		return diag.Errorf("error shrinking database logs of RDS instance(%s), the response is: %s", instanceId,
			resp.(string))
	}

	d.SetId(instanceId)

	return nil
}

func buildCreateDbLogsShrinkingBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
	}
	return bodyParams
}

func resourceRdsDbLogsShrinkingRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsDbLogsShrinkingUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsDbLogsShrinkingDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting database logs shrinking is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before database logs shrinking."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
