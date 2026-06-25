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

var instanceDatabaseRoleNonUpdatableParams = []string{"instance_id", "name", "password"}

// @API GaussDB POST /v3.1/{project_id}/instances/{instance_id}/db-role
// @API GaussDB GET /v3.1/{project_id}/instances/{instance_id}/db-role
// @API GaussDB GET /v3/{project_id}/instances
func ResourceGaussdbInstanceDatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussdbInstanceDatabaseRoleCreate,
		ReadContext:   resourceGaussdbInstanceDatabaseRoleRead,
		UpdateContext: resourceGaussdbInstanceDatabaseRoleUpdate,
		DeleteContext: resourceGaussdbInstanceDatabaseRoleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussdbInstanceDatabaseRoleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceDatabaseRoleNonUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceDatabaseRoleAttributeSchema(),
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

func instanceDatabaseRoleAttributeSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func resourceGaussdbInstanceDatabaseRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-role"
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
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussdbInstanceDatabaseRoleBodyParams(d))

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
		return diag.Errorf("error creating GaussDB instance database role: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	resourceId := fmt.Sprintf("%s/%s", instanceId, name)
	d.SetId(resourceId)

	return resourceGaussdbInstanceDatabaseRoleRead(ctx, d, meta)
}

func resourceGaussdbInstanceDatabaseRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/db-role"
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
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB instance database role")
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

	name := d.Get("name").(string)
	role := utils.PathSearch(fmt.Sprintf("roles[?name=='%s']|[0]", name), getRespBody, nil)
	if role == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB instance database role")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", d.Get("instance_id").(string)),
		d.Set("name", utils.PathSearch("name", role, nil)),
		d.Set("attribute", flattenGetInstanceDatabaseRoleAttributeBody(utils.PathSearch("attribute", role, nil))),
		d.Set("memberof", utils.PathSearch("memberof", role, nil)),
		d.Set("lock_status", utils.PathSearch("lock_status", role, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussdbInstanceDatabaseRoleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussdbInstanceDatabaseRoleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance database role resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceGaussdbInstanceDatabaseRoleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func buildCreateGaussdbInstanceDatabaseRoleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":     d.Get("name"),
		"password": d.Get("password"),
	}
}

func flattenGetInstanceDatabaseRoleAttributeBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"rolsuper":            utils.PathSearch("rolsuper", resp, false),
			"rolinherit":          utils.PathSearch("rolinherit", resp, false),
			"rolcreaterole":       utils.PathSearch("rolcreaterole", resp, false),
			"rolcreatedb":         utils.PathSearch("rolcreatedb", resp, false),
			"rolcanlogin":         utils.PathSearch("rolcanlogin", resp, false),
			"rolconnlimit":        utils.PathSearch("rolconnlimit", resp, 0),
			"rolreplication":      utils.PathSearch("rolreplication", resp, false),
			"rolbypassrls":        utils.PathSearch("rolbypassrls", resp, false),
			"rolpassworddeadline": utils.PathSearch("rolpassworddeadline", resp, nil),
		},
	}
}
