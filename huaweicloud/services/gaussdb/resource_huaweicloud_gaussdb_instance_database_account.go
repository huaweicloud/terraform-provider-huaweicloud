package gaussdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

var gaussdbInstanceDatabaseAccountNonUpdatableParams = []string{"instance_id", "name", "is_login_only"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/db-user
// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/db-user/password
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/db-users
// @API GaussDB GET /v3/{project_id}/instances
func ResourceGaussDBInstanceDatabaseAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBInstanceDatabaseAccountCreate,
		UpdateContext: resourceGaussDBInstanceDatabaseAccountUpdate,
		ReadContext:   resourceGaussDBInstanceDatabaseAccountRead,
		DeleteContext: resourceGaussDBInstanceDatabaseAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussDBInstanceDatabaseAccountImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussdbInstanceDatabaseAccountNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"is_login_only": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rolsuper": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolinherit": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolcreaterole": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolcreatedb": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolcanlogin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolconnlimit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rolreplication": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolbypassrls": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rolpassworddeadline": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"memberof": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBInstanceDatabaseAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-user"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBInstanceDatabaseAccountBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB instance database account: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("instance_id").(string), d.Get("name").(string)))

	return resourceGaussDBInstanceDatabaseAccountRead(ctx, d, meta)
}

func buildCreateGaussDBInstanceDatabaseAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}

	if v, ok := d.GetOk("is_login_only"); ok {
		bodyParams["is_login_only"] = v.(string) == "true"
	}

	return bodyParams
}

func resourceGaussDBInstanceDatabaseAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB instance database account")
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	accountName := d.Get("name").(string)
	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getRespBody, nil)
	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB instance database account")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", d.Get("instance_id").(string)),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("attribute", flattenGaussDBInstanceDatabaseAccountAttribute(account)),
		d.Set("memberof", utils.PathSearch("memberof", account, nil)),
		d.Set("lock_status", utils.PathSearch("lock_status", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDBInstanceDatabaseAccountAttribute(account interface{}) []map[string]interface{} {
	attribute := utils.PathSearch("attribute", account, nil)
	if attribute == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"rolsuper":            utils.PathSearch("attribute.rolsuper", account, false),
			"rolinherit":          utils.PathSearch("attribute.rolinherit", account, false),
			"rolcreaterole":       utils.PathSearch("attribute.rolcreaterole", account, false),
			"rolcreatedb":         utils.PathSearch("attribute.rolcreatedb", account, false),
			"rolcanlogin":         utils.PathSearch("attribute.rolcanlogin", account, false),
			"rolconnlimit":        utils.PathSearch("attribute.rolconnlimit", account, 0),
			"rolreplication":      utils.PathSearch("attribute.rolreplication", account, false),
			"rolbypassrls":        utils.PathSearch("attribute.rolbypassrls", account, false),
			"rolpassworddeadline": utils.PathSearch("attribute.rolpassworddeadline", account, ""),
		},
	}
}

func resourceGaussDBInstanceDatabaseAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	if d.HasChange("password") {
		if err = updateGaussDBInstanceDatabaseAccountPassword(d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBInstanceDatabaseAccountRead(ctx, d, meta)
}

func updateGaussDBInstanceDatabaseAccountPassword(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/db-user/password"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = buildUpdateGaussDBInstanceDatabaseAccountPasswordBodyParams(d)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB instance database account password: %s", err)
	}

	return nil
}

func buildUpdateGaussDBInstanceDatabaseAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}
}

func resourceGaussDBInstanceDatabaseAccountDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance database account resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceGaussDBInstanceDatabaseAccountImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	d.SetId(fmt.Sprintf("%s/%s", parts[0], parts[1]))
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
