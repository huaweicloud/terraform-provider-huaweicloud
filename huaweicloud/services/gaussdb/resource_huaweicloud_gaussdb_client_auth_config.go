package gaussdb

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussDBClientAuthConfigNonUpdatableParams = []string{"instance_id", "type", "database", "user", "address"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/hba-info
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}/hba-info
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/hba-info
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/hba-info
func ResourceGaussDbClientAuthConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBClientAuthConfigCreate,
		ReadContext:   resourceGaussDBClientAuthConfigRead,
		DeleteContext: resourceGaussDBClientAuthConfigDelete,
		UpdateContext: resourceGaussDBClientAuthConfigUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussClientAuthConfigImport,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussDBClientAuthConfigNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
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

func resourceGaussDBClientAuthConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBClientAuthConfigBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB client auth config: %s", err)
	}

	id := fmt.Sprintf("%s:%s:%s:%s:%s", instanceID, d.Get("type").(string), d.Get("database").(string),
		d.Get("user").(string), d.Get("address").(string))
	d.SetId(id)
	return resourceGaussDBClientAuthConfigRead(ctx, d, meta)
}

func buildCreateGaussDBClientAuthConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"hba_confs": []map[string]interface{}{
			buildCreateGaussDBClientAuthConfigChildBody(d),
		},
	}
	return bodyParams
}

func buildCreateGaussDBClientAuthConfigChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"type":     d.Get("type").(string),
		"database": d.Get("database").(string),
		"user":     d.Get("user").(string),
		"address":  d.Get("address").(string),
		"method":   d.Get("method").(string),
	}

	return params
}

func resourceGaussDBClientAuthConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var offset int

	typeRaw := d.Get("type").(string)
	databaseRaw := d.Get("database").(string)
	userRaw := d.Get("user").(string)
	addressRaw := d.Get("address").(string)
	methodRaw := d.Get("method").(string)

	var hbaConf interface{}

	for {
		getPath := getBasePath + buildGaussDBClientAuthConfigQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving GaussDB client auth config")
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		hbaConf = utils.PathSearch(
			fmt.Sprintf("hba_confs|[?type=='%s' && database=='%s' && user=='%s' && address=='%s']|[0]",
				typeRaw, databaseRaw, userRaw, addressRaw),
			getRespBody,
			nil,
		)

		if hbaConf != nil {
			break
		}

		hbaConfs := utils.PathSearch("hba_confs", getRespBody, []interface{}{}).([]interface{})

		if len(hbaConfs) < 100 {
			break
		}

		offset += 100
	}

	if hbaConf == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB client auth config")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("type", typeRaw),
		d.Set("database", databaseRaw),
		d.Set("user", userRaw),
		d.Set("address", addressRaw),
		d.Set("method", methodRaw),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDBClientAuthConfigQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func resourceGaussDBClientAuthConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	oldMethodRaw, newMethodRaw := d.GetChange("method")

	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBClientAuthConfigBodyParams(
		oldMethodRaw.(string),
		newMethodRaw.(string),
		d,
	))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error updating GaussDB client auth config: %s", err)
	}

	return resourceGaussDBClientAuthConfigRead(ctx, d, meta)
}

func buildUpdateGaussDBClientAuthConfigBodyParams(oldMethod, newMethod string, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"before_conf": map[string]interface{}{
			"type":     d.Get("type").(string),
			"database": d.Get("database").(string),
			"user":     d.Get("user").(string),
			"address":  d.Get("address").(string),
			"method":   oldMethod,
		},
		"after_conf": map[string]interface{}{
			"type":     d.Get("type").(string),
			"database": d.Get("database").(string),
			"user":     d.Get("user").(string),
			"address":  d.Get("address").(string),
			"method":   newMethod,
		},
	}
}

func resourceGaussDBClientAuthConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	deleteOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBClientAuthConfigBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200823"),
			"error deleting GaussDB client auth config")
	}

	return nil
}

func resourceGaussClientAuthConfigImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>:<type>:<database>:<user>:<address>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("type", parts[1]),
		d.Set("database", parts[2]),
		d.Set("user", parts[3]),
		d.Set("address", parts[4]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
