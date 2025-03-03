package rds

import (
	"context"
	"fmt"
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

var sqlServerDatabaseCopyNonUpdatableParams = []string{"instance_id", "procedure_name", "db_name_source", "db_name_target"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/database/procedure
func ResourceSQLServerDatabaseCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLServerDatabaseCopyCreate,
		ReadContext:   resourceSQLServerDatabaseCopyRead,
		UpdateContext: resourceSQLServerDatabaseCopyUpdate,
		DeleteContext: resourceSQLServerDatabaseCopyDelete,

		CustomizeDiff: config.FlexibleForceNew(sqlServerDatabaseCopyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"procedure_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name_target": {
				Type:     schema.TypeString,
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

func resourceSQLServerDatabaseCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/database/procedure"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateSQLServerDatabaseCopyBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS SQLServer database copy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("instance_id").(string), d.Get("procedure_name").(string)))

	return nil
}

func buildCreateSQLServerDatabaseCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"procedure_name": d.Get("procedure_name"),
		"params":         buildCreateSQLServerDatabaseCopyParamsBody(d),
	}
	return bodyParams
}

func buildCreateSQLServerDatabaseCopyParamsBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name_source": d.Get("db_name_source"),
		"db_name_target": d.Get("db_name_target"),
	}
	return bodyParams
}

func resourceSQLServerDatabaseCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSQLServerDatabaseCopyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSQLServerDatabaseCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS SQLServer database copy resource is not supported. The resource is only removed from the" +
		"state, the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
