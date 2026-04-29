package geminidb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

var geminiDbAccountNonUpdatableParams = []string{
	"instance_id",
	"name",
}

// @API GeminiDBforNoSQL POST /v3/{project_id}/instances
// @API GeminiDBforNoSQL GET /v3/{project_id}/jobs
// @API GeminiDBforNoSQL GET /v3/{project_id}/redis/instances/{instance_id}/db-users
// @API GeminiDBforNoSQL POST /v3/{project_id}/redis/instances/{instance_id}/db-users
// @API GeminiDBforNoSQL PUT /v3/{project_id}/redis/instances/{instance_id}/db-users/password
// @API GeminiDBforNoSQL PUT /v3/{project_id}/redis/instances/{instance_id}/db-users/privilege
// @API GeminiDBforNoSQL DELETE /v3/{project_id}/redis/instances/{instance_id}/db-users
func ResourceGeminidbAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDbAccountCreate,
		ReadContext:   resourceGeminiDbAccountRead,
		UpdateContext: resourceGeminiDbAccountUpdate,
		DeleteContext: resourceGeminiDbAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGeminiDBAccountImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDbAccountNonUpdatableParams),

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
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"privilege": {
				Type:     schema.TypeString,
				Required: true,
			},
			"databases": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGeminiDbAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/db-users"
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGeminiDbAccountBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}

	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GeminiDB account: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("instance_id").(string) + "/" + d.Get("name").(string))

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GeminiDB account: job_id is not found in API response")
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGeminiDbAccountRead(ctx, d, meta)
}

func buildCreateGeminiDbAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"users": buildCreateGeminiDbUsersBody(d),
	}
	return bodyParams
}

func buildCreateGeminiDbUsersBody(d *schema.ResourceData) []map[string]interface{} {
	rst := []map[string]interface{}{
		{
			"name":      d.Get("name"),
			"password":  d.Get("password"),
			"privilege": d.Get("privilege"),
			"databases": d.Get("databases").(*schema.Set).List(),
		},
	}
	return rst
}

func resourceGeminiDbAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/db-users?name={name}"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{name}", d.Get("name").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB account")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	user := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", d.Get("name").(string)), getRespBody, nil)
	if user == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", user, nil)),
		d.Set("type", utils.PathSearch("type", user, nil)),
		d.Set("privilege", utils.PathSearch("privilege", user, nil)),
		d.Set("databases", utils.PathSearch("databases", user, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDbAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	err = updateGeminiDbAccountPassword(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = updateGeminiDbAccountPrivilege(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGeminiDbAccountRead(ctx, d, meta)
}

func updateGeminiDbAccountPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChange("password") {
		return nil
	}

	var (
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/db-users/password"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes:          []int{204},
		JSONBody:         buildUpdateGeminiDbAccountPasswordBodyParams(d),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}

	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating GeminiDB account password: %s", err)
	}

	return nil
}

func buildUpdateGeminiDbAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}
	return bodyParams
}

func updateGeminiDbAccountPrivilege(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	if !d.HasChanges("privilege", "databases") {
		return nil
	}

	updatePath := client.Endpoint + "v3/{project_id}/redis/instances/{instance_id}/db-users/privilege"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildUpdateGeminiDbAccountPrivilegeBodyParams(d),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}

	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error updating GeminiDB account: %s", err)
	}

	updRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error updating GeminiDB account: job_id is not found in the response: %s", d.Id())
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}

func buildUpdateGeminiDbAccountPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"users": []map[string]interface{}{
			{
				"name":      d.Get("name"),
				"privilege": d.Get("privilege"),
				"databases": d.Get("databases").(*schema.Set).List(),
			},
		},
	}
	return bodyParams
}

func resourceGeminiDbAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/db-users"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGeminiDbAccountBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleDeletionError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.03280008"),
			"error retrieving Geminidb account")
	}
	deleteRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting GeminiDB account: job_id is not found in the response: %s", d.Id())
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildDeleteGeminiDbAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"names": []string{d.Get("name").(string)},
	}
	return bodyParams
}

func resourceGeminiDBAccountImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	instanceID := parts[0]
	name := parts[1]

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceID),
		d.Set("name", name),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
