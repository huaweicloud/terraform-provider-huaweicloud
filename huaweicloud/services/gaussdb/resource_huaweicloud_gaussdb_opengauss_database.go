package gaussdb

import (
	"context"
	"fmt"
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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/database
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/databases
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}/database
func ResourceOpenGaussDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussDatabaseCreate,
		ReadContext:   resourceOpenGaussDatabaseRead,
		DeleteContext: resourceOpenGaussDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenGaussDatabaseImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
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
				ForceNew:    true,
				Description: `Specifies the ID of the GaussDB instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database name.`,
			},
			"character_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the database character set.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the database user.`,
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the name of the database template.`,
			},
			"lc_collate": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the database collation.`,
			},
			"lc_ctype": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the database classification.`,
			},
			"size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the database size.`,
			},
			"compatibility_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the database compatibility type.`,
			},
		},
	}
}

func resourceOpenGaussDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/database"
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
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussDatabaseBodyParams(d))

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
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB openGauss database: %s", err)
	}

	databaseName := d.Get("name").(string)
	d.SetId(instanceID + "/" + databaseName)

	return resourceOpenGaussDatabaseRead(ctx, d, meta)
}

func buildCreateOpenGaussDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"character_set": utils.ValueIgnoreEmpty(d.Get("character_set")),
		"owner":         utils.ValueIgnoreEmpty(d.Get("owner")),
		"template":      utils.ValueIgnoreEmpty(d.Get("template")),
		"lc_collate":    utils.ValueIgnoreEmpty(d.Get("lc_collate")),
		"lc_ctype":      utils.ValueIgnoreEmpty(d.Get("lc_ctype")),
	}
	return bodyParams
}

func resourceOpenGaussDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/databases"
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
	}

	// actually, offset is pageNo in API, not offset
	var offset int
	var database interface{}
	databaseName := d.Get("name").(string)

	for {
		getPath := getBasePath + buildOpenGaussDatabaseQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving GaussDB openGauss database")
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		database = utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", databaseName), getRespBody, nil)
		if database != nil {
			break
		}
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if int(totalCount) <= (offset+1)*100 {
			break
		}
		offset++
	}
	if database == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB openGauss database")
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
		d.Set("compatibility_type", utils.PathSearch("compatibility_type", database, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOpenGaussDatabaseQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func resourceOpenGaussDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/database?database_name={database_name}"
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
	deletePath = strings.ReplaceAll(deletePath, "{database_name}", d.Get("name").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

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
			"error deleting GaussDB openGauss database")
	}

	return nil
}

func resourceOpenGaussDatabaseImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
