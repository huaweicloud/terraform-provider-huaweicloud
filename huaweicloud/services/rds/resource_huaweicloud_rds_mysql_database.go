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

var mysqlDatabaseNonUpdatableParams = []string{"instance_id", "name", "character_set"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/database
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/detail
// @API RDS POST /v3/{project_id}/instances/{instance_id}/database/update
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/database/{db_name}
func ResourceMysqlDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlDatabaseCreate,
		UpdateContext: resourceMysqlDatabaseUpdate,
		ReadContext:   resourceMysqlDatabaseRead,
		DeleteContext: resourceMysqlDatabaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlDatabaseNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS Mysql instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"character_set": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the character set used by the database.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the database description.`,
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

func resourceMysqlDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createMysqlDatabase: create RDS Mysql database.
	var (
		createMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database"
		createMysqlDatabaseProduct = "rds"
	)
	createMysqlDatabaseClient, err := cfg.NewServiceClient(createMysqlDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createMysqlDatabasePath := createMysqlDatabaseClient.Endpoint + createMysqlDatabaseHttpUrl
	createMysqlDatabasePath = strings.ReplaceAll(createMysqlDatabasePath, "{project_id}",
		createMysqlDatabaseClient.ProjectID)
	createMysqlDatabasePath = strings.ReplaceAll(createMysqlDatabasePath, "{instance_id}", instanceId)

	createMysqlDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createMysqlDatabaseOpt.JSONBody = utils.RemoveNil(buildCreateMysqlDatabaseBodyParams(d))
	log.Printf("[DEBUG] Create RDS Mysql database options: %#v", createMysqlDatabaseOpt)

	retryFunc := func() (interface{}, bool, error) {
		_, err = createMysqlDatabaseClient.Request("POST", createMysqlDatabasePath, &createMysqlDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createMysqlDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS Mysql database: %s", err)
	}

	dbName := d.Get("name").(string)
	d.SetId(instanceId + "/" + dbName)

	return resourceMysqlDatabaseRead(ctx, d, meta)
}

func buildCreateMysqlDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"character_set": utils.ValueIgnoreEmpty(d.Get("character_set")),
		"comment":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceMysqlDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMysqlDatabase: query RDS Mysql database
	var (
		getMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		getMysqlDatabaseProduct = "rds"
	)
	getMysqlDatabaseClient, err := cfg.NewServiceClient(getMysqlDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePath := getMysqlDatabaseClient.Endpoint + getMysqlDatabaseHttpUrl
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{project_id}",
		getMysqlDatabaseClient.ProjectID)
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{instance_id}", instanceId)

	getMysqlDatabaseResp, err := pagination.ListAllItems(
		getMysqlDatabaseClient,
		"page",
		getMysqlDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS Mysql database")
	}

	getMysqlDatabaseRespJson, err := json.Marshal(getMysqlDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getMysqlDatabaseRespBody interface{}
	err = json.Unmarshal(getMysqlDatabaseRespJson, &getMysqlDatabaseRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getMysqlDatabaseRespBody, nil)
	if database == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	characterSet := utils.PathSearch("character_set", database, "").(string)
	characterSet = strings.TrimSuffix(characterSet, "mb3")
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", database, nil)),
		d.Set("character_set", characterSet),
		d.Set("description", utils.PathSearch("comment", database, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMysqlDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.HasChange("description") {
		// updateMysqlDatabase: update RDS Mysql database
		var (
			updateMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/update"
			updateMysqlDatabaseProduct = "rds"
		)
		updateMysqlDatabaseClient, err := cfg.NewServiceClient(updateMysqlDatabaseProduct, region)
		if err != nil {
			return diag.Errorf("error creating RDS client: %s", err)
		}

		instanceId := d.Get("instance_id").(string)
		updateMysqlDatabasePath := updateMysqlDatabaseClient.Endpoint + updateMysqlDatabaseHttpUrl
		updateMysqlDatabasePath = strings.ReplaceAll(updateMysqlDatabasePath, "{project_id}",
			updateMysqlDatabaseClient.ProjectID)
		updateMysqlDatabasePath = strings.ReplaceAll(updateMysqlDatabasePath, "{instance_id}", instanceId)

		updateMysqlDatabaseOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateMysqlDatabaseOpt.JSONBody = utils.RemoveNil(buildUpdateMysqlDatabaseBodyParams(d))
		log.Printf("[DEBUG] Update RDS Mysql database options: %#v", updateMysqlDatabaseOpt)
		retryFunc := func() (interface{}, bool, error) {
			_, err = updateMysqlDatabaseClient.Request("POST", updateMysqlDatabasePath, &updateMysqlDatabaseOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rdsInstanceStateRefreshFunc(updateMysqlDatabaseClient, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating RDS Mysql database: %s", err)
		}
	}

	return resourceMysqlDatabaseRead(ctx, d, meta)
}

func buildUpdateMysqlDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceMysqlDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteMysqlDatabase: delete RDS Mysql database
	var (
		deleteMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/{db_name}"
		deleteMysqlDatabaseProduct = "rds"
	)
	deleteMysqlDatabaseClient, err := cfg.NewServiceClient(deleteMysqlDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteMysqlDatabasePath := deleteMysqlDatabaseClient.Endpoint + deleteMysqlDatabaseHttpUrl
	deleteMysqlDatabasePath = strings.ReplaceAll(deleteMysqlDatabasePath, "{project_id}",
		deleteMysqlDatabaseClient.ProjectID)
	deleteMysqlDatabasePath = strings.ReplaceAll(deleteMysqlDatabasePath, "{instance_id}", instanceId)
	deleteMysqlDatabasePath = strings.ReplaceAll(deleteMysqlDatabasePath, "{db_name}", d.Get("name").(string))

	deleteMysqlDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	log.Printf("[DEBUG] Delete RDS Mysql database options: %#v", deleteMysqlDatabaseOpt)
	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteMysqlDatabaseClient.Request("DELETE", deleteMysqlDatabasePath, &deleteMysqlDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteMysqlDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS Mysql database: %s", err)
	}

	return nil
}
