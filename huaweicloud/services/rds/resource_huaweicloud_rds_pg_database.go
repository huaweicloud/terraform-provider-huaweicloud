// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pgDatabaseNonUpdatableParams = []string{"instance_id", "name", "template", "character_set", "lc_collate",
	"lc_ctype", "is_revoke_public_privilege"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/database
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/detail
// @API RDS POST /v3/{project_id}/instances/{instance_id}/database/owner
// @API RDS POST /v3/{project_id}/instances/{instance_id}/database/update
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/database/{db_name}
func ResourcePgDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgDatabaseCreate,
		UpdateContext: resourcePgDatabaseUpdate,
		ReadContext:   resourcePgDatabaseRead,
		DeleteContext: resourcePgDatabaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgDatabaseNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database user.`,
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the database template.`,
			},
			"character_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database character set.`,
			},
			"lc_collate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database collocation.`,
			},
			"lc_ctype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database classification.`,
			},
			"is_revoke_public_privilege": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to revoke the PUBLIC CREATE permission of the public schema.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the database description.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the database size, in bytes.`,
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

func resourcePgDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPgDatabase: create RDS PostgreSQL database.
	var (
		createPgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database"
		createPgDatabaseProduct = "rds"
	)
	createPgDatabaseClient, err := cfg.NewServiceClient(createPgDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPgDatabasePath := createPgDatabaseClient.Endpoint + createPgDatabaseHttpUrl
	createPgDatabasePath = strings.ReplaceAll(createPgDatabasePath, "{project_id}", createPgDatabaseClient.ProjectID)
	createPgDatabasePath = strings.ReplaceAll(createPgDatabasePath, "{instance_id}", instanceId)

	createPgDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestBody := buildCreatePgDatabaseBodyParams(d)
	log.Printf("[DEBUG] Create RDS PostgreSQL database options: %#v", requestBody)
	createPgDatabaseOpt.JSONBody = utils.RemoveNil(requestBody)
	retryFunc := func() (interface{}, bool, error) {
		_, err = createPgDatabaseClient.Request("POST", createPgDatabasePath, &createPgDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createPgDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL database: %s", err)
	}

	dbName := d.Get("name").(string)
	d.SetId(instanceId + "/" + dbName)

	return resourcePgDatabaseRead(ctx, d, meta)
}

func buildCreatePgDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                       d.Get("name"),
		"owner":                      utils.ValueIgnoreEmpty(d.Get("owner")),
		"template":                   utils.ValueIgnoreEmpty(d.Get("template")),
		"character_set":              utils.ValueIgnoreEmpty(d.Get("character_set")),
		"lc_collate":                 utils.ValueIgnoreEmpty(d.Get("lc_collate")),
		"lc_ctype":                   utils.ValueIgnoreEmpty(d.Get("lc_ctype")),
		"is_revoke_public_privilege": utils.ValueIgnoreEmpty(d.Get("is_revoke_public_privilege")),
		"comment":                    utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourcePgDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPgDatabase: query RDS PostgreSQL database
	var (
		getPgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?db={db_name}&page=1&limit=100"
		getPgDatabaseProduct = "rds"
	)
	getPgDatabaseClient, err := cfg.NewServiceClient(getPgDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getPgDatabasePath := getPgDatabaseClient.Endpoint + getPgDatabaseHttpUrl
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{project_id}", getPgDatabaseClient.ProjectID)
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{instance_id}", instanceId)
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{db_name}", dbName)

	getPgDatabaseResp, err := pagination.ListAllItems(
		getPgDatabaseClient,
		"page",
		getPgDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL database")
	}

	getPgDatabaseRespJson, err := json.Marshal(getPgDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getPgDatabaseRespBody interface{}
	err = json.Unmarshal(getPgDatabaseRespJson, &getPgDatabaseRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getPgDatabaseRespBody, nil)
	if database == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", database, nil)),
		d.Set("owner", utils.PathSearch("owner", database, nil)),
		d.Set("character_set", utils.PathSearch("character_set", database, nil)),
		d.Set("lc_collate", utils.PathSearch("collate_set", database, nil)),
		d.Set("size", utils.PathSearch("size", database, nil)),
		d.Set("description", utils.PathSearch("comment", database, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePgDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updatePgDatabaseProduct = "rds"
	)
	updatePgDatabaseClient, err := cfg.NewServiceClient(updatePgDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.HasChange("owner") {
		// updatePgDatabaseOwner: update RDS PostgreSQL database owner
		if err = updatePgDatabaseOwner(ctx, d, updatePgDatabaseClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		// updatePgDatabaseDescription: update RDS PostgreSQL database description
		if err = updatePgDatabaseDescription(ctx, d, updatePgDatabaseClient); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourcePgDatabaseRead(ctx, d, meta)
}

func updatePgDatabaseOwner(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		updatePgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/owner"
	)

	instanceId := d.Get("instance_id").(string)
	updatePgDatabasePath := client.Endpoint + updatePgDatabaseHttpUrl
	updatePgDatabasePath = strings.ReplaceAll(updatePgDatabasePath, "{project_id}", client.ProjectID)
	updatePgDatabasePath = strings.ReplaceAll(updatePgDatabasePath, "{instance_id}", instanceId)

	updatePgDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updatePgDatabaseOpt.JSONBody = utils.RemoveNil(buildUpdatePgDatabaseOwnerBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", updatePgDatabasePath, &updatePgDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS PostgreSQL database owner: %s", err)
	}
	return nil
}

func buildUpdatePgDatabaseOwnerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"owner":    d.Get("owner"),
		"database": d.Get("name"),
	}
	return bodyParams
}

func updatePgDatabaseDescription(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		updatePgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/update"
	)

	instanceId := d.Get("instance_id").(string)
	updatePgDatabasePath := client.Endpoint + updatePgDatabaseHttpUrl
	updatePgDatabasePath = strings.ReplaceAll(updatePgDatabasePath, "{project_id}", client.ProjectID)
	updatePgDatabasePath = strings.ReplaceAll(updatePgDatabasePath, "{instance_id}", instanceId)

	updatePgDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updatePgDatabaseOpt.JSONBody = utils.RemoveNil(buildUpdatePgDatabaseDescriptionBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", updatePgDatabasePath, &updatePgDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS PostgreSQL database description: %s", err)
	}
	return nil
}

func buildUpdatePgDatabaseDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    d.Get("name"),
		"comment": d.Get("description"),
	}
	return bodyParams
}

func resourcePgDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePgDatabase: delete RDS PostgreSQL database
	var (
		deletePgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/{db_name}"
		deletePgDatabaseProduct = "rds"
	)
	deletePgDatabaseClient, err := cfg.NewServiceClient(deletePgDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deletePgDatabasePath := deletePgDatabaseClient.Endpoint + deletePgDatabaseHttpUrl
	deletePgDatabasePath = strings.ReplaceAll(deletePgDatabasePath, "{project_id}", deletePgDatabaseClient.ProjectID)
	deletePgDatabasePath = strings.ReplaceAll(deletePgDatabasePath, "{instance_id}", instanceId)
	deletePgDatabasePath = strings.ReplaceAll(deletePgDatabasePath, "{db_name}", d.Get("name").(string))

	deletePgDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deletePgDatabaseClient.Request("DELETE", deletePgDatabasePath, &deletePgDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deletePgDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS PostgreSQL database: %s", err)
	}

	return nil
}
