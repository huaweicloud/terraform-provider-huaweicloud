package taurusdb

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

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB POST /v3/{project_id}/logs/lts-configs
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}
// @API TaurusDB GET /v3/{project_id}/logs/lts-configs
// @API TaurusDB DELETE /v3/{project_id}/logs/lts-configs
func ResourceTaurusDBLtsLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBLtsLogCreateOrUpdate,
		ReadContext:   resourceTaurusDBLtsLogRead,
		UpdateContext: resourceTaurusDBLtsLogCreateOrUpdate,
		DeleteContext: resourceTaurusDBLtsLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBLtsLogImportState,
		},

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the TaurusDB instance.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the LTS log.`,
			},
			"lts_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the LTS log group.`,
			},
			"lts_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the LTS log stream.`,
			},
		},
	}
}

func resourceTaurusDBLtsLogCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/logs/lts-configs"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildGaussDBMysqlLtsLogBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating TaurusDB LTS log: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, logType))

	return resourceTaurusDBLtsLogRead(ctx, d, meta)
}

func buildGaussDBMysqlLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"log_configs": []map[string]interface{}{
			{
				"instance_id":   d.Get("instance_id"),
				"log_type":      d.Get("log_type"),
				"lts_group_id":  d.Get("lts_group_id"),
				"lts_stream_id": d.Get("lts_stream_id"),
			},
		},
	}
	return ltsConfigs
}

func resourceTaurusDBLtsLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	var (
		httpUrl = "v3/{project_id}/logs/lts-configs"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?instance_id=%s", d.Get("instance_id").(string))

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB LTS configs: %s", err)
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.Errorf("error marshaling TaurusDB LTS configs: %s", err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.Errorf("error unmarshalling TaurusDB LTS configs: %s", err)
	}
	searchPath := fmt.Sprintf("instance_lts_configs|[0].lts_configs|[?log_type=='%s' && enabled]|[0]", d.Get("log_type").(string))
	ltsConfig := utils.PathSearch(searchPath, listRespBody, nil)
	if ltsConfig == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_lts_configs[0].instance.id", listRespBody, nil)),
		d.Set("log_type", utils.PathSearch("log_type", ltsConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", ltsConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", ltsConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBLtsLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/logs/lts-configs"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBMysqlLtsLogBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.280001"),
			"error deleting TaurusDB LTS log")
	}

	return resourceTaurusDBLtsLogRead(ctx, d, meta)
}

func buildDeleteGaussDBMysqlLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"log_configs": []map[string]interface{}{
			{
				"instance_id": d.Get("instance_id"),
				"log_type":    d.Get("log_type"),
			},
		},
	}
	return ltsConfigs
}

func resourceTaurusDBLtsLogImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<log_type>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
