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

var pgPluginNonUpdatableParams = []string{"instance_id", "name", "database_name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/extensions
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/extensions
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/extensions
func ResourceRdsPgPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPgInstancePluginCreate,
		ReadContext:   resourceRdsPgInstancePluginRead,
		UpdateContext: resourceRdsPgInstancePluginUpdate,
		DeleteContext: resourceRdsPgInstancePluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgPluginNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_preload_libraries": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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

func buildCreatePgPluginBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name":  d.Get("database_name"),
		"extension_name": d.Get("name"),
	}
	return bodyParams
}

func queryPluginDetail(client *golangsdk.ServiceClient, instanceId, databaseName, name string) (interface{}, error) {
	listPgPluginHttpUrl := "v3/{project_id}/instances/{instance_id}/extensions?database_name={database_name}"
	listPgPluginPath := client.Endpoint + listPgPluginHttpUrl
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{project_id}", client.ProjectID)
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{instance_id}", instanceId)
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{database_name}", databaseName)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		listPgPluginPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		if errCode := parseErrCode(resp); errCode == "DBS.280238" || errCode == "DBS.200823" {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v3/{project_id}/instances/{instance_id}/extensions",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the PostgreSQL plugin (%s) does not exist", name)),
				},
			}
		}
		return nil, err
	}
	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var bodyJson interface{}
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Reading RDS PostgreSQL plugin response: %#v", bodyJson)

	pluginDetail := utils.PathSearch(fmt.Sprintf("extensions[?name=='%s']|[?created]|[0]", name), bodyJson, nil)
	if pluginDetail == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/instances/{instance_id}/extensions",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the PostgreSQL plugin (%s) does not exist or its creation status is not created", name)),
			},
		}
	}

	return pluginDetail, nil
}

func parseErrCode(resp interface{}) string {
	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] Error marshaling PostgreSQL plugin: %s", err)
		return ""
	}

	var bodyJson interface{}
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		log.Printf("[ERROR] Error unmarshal PostgreSQL plugin: %s", err)
		return ""
	}
	errCode := utils.PathSearch("error.code", bodyJson, "")
	return errCode.(string)
}

func resourceRdsPgInstancePluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var createPgPluginProduct = "rds"
	createPgPluginClient, err := cfg.NewServiceClient(createPgPluginProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	pluginDetail, err := queryPluginDetail(createPgPluginClient, instanceId, d.Get("database_name").(string), d.Get("name").(string))
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return diag.FromErr(err)
		}
	}

	if pluginDetail != nil {
		return diag.Errorf("The RDS PostgreSQL plugin already created: %#v", pluginDetail)
	}

	createPgPluginHttpUrl := "v3/{project_id}/instances/{instance_id}/extensions"
	createPgPluginPath := createPgPluginClient.Endpoint + createPgPluginHttpUrl
	createPgPluginPath = strings.ReplaceAll(createPgPluginPath, "{project_id}", createPgPluginClient.ProjectID)
	createPgPluginPath = strings.ReplaceAll(createPgPluginPath, "{instance_id}", instanceId)

	jsonBody := utils.RemoveNil(buildCreatePgPluginBody(d))
	log.Printf("[DEBUG] Create RDS PostgreSQL plugin: %#v", jsonBody)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = createPgPluginClient.Request("POST", createPgPluginPath, &reqOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createPgPluginClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL plugin: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", d.Get("instance_id").(string), d.Get("database_name").(string), d.Get("name").(string))
	d.SetId(id)

	return resourceRdsPgInstancePluginRead(ctx, d, cfg)
}

func resourceRdsPgInstancePluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getPgPluginClient, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return diag.Errorf("invalid ID format in ReadContext: %s", d.Id())
	}
	pluginDetail, err := queryPluginDetail(getPgPluginClient, parts[0], parts[1], parts[2])
	log.Printf("[DEBUG] RDS PostgreSQL plugin detail: %#v", pluginDetail)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving PostgreSQL plugin")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("instance_id", parts[0]),
		d.Set("name", utils.PathSearch("name", pluginDetail, "")),
		d.Set("database_name", utils.PathSearch("database_name", pluginDetail, "")),
		d.Set("version", utils.PathSearch("version", pluginDetail, "")),
		d.Set("shared_preload_libraries", utils.PathSearch("shared_preload_libraries", pluginDetail, "")),
		d.Set("description", utils.PathSearch("description", pluginDetail, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeletePgPluginBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name":  d.Get("database_name"),
		"extension_name": d.Get("name"),
	}
	return bodyParams
}

func resourceRdsPgInstancePluginUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsPgInstancePluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	deletePgPluginClient, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deletePgPluginHttpUrl := "v3/{project_id}/instances/{instance_id}/extensions"
	deletePgPluginPath := deletePgPluginClient.Endpoint + deletePgPluginHttpUrl
	deletePgPluginPath = strings.ReplaceAll(deletePgPluginPath, "{project_id}", deletePgPluginClient.ProjectID)
	deletePgPluginPath = strings.ReplaceAll(deletePgPluginPath, "{instance_id}", d.Get("instance_id").(string))

	jsonBody := utils.RemoveNil(buildDeletePgPluginBody(d))
	log.Printf("[DEBUG] Delete RDS PostgreSQL plugin: %#v", jsonBody)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err = deletePgPluginClient.Request("DELETE", deletePgPluginPath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deletePgPluginClient, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error deleting RDS PostgreSQL plugin: %s", err)
	}

	return nil
}
