package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

var htapLTSConfigNoneUpdatableParams = []string{
	"instance", "log_type",
}

// @API TaurusDB POST /v3/{project_id}/starrocks/instances/logs/lts-configs
// @API TaurusDB GET /v3/{project_id}/starrocks/instances/logs/lts-configs
// @API TaurusDB DELETE /v3/{project_id}/starrocks/instances/logs/lts-configs
func ResourceTaurusDBHtapStarrocksLtsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksLtsConfigCreateOrUpdate,
		ReadContext:   resourceTaurusDBHtapStarrocksLtsConfigRead,
		UpdateContext: resourceTaurusDBHtapStarrocksLtsConfigCreateOrUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksLtsConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBHtapStarrocksLtsConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(htapLTSConfigNoneUpdatableParams),

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
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_stream_id": {
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

func resourceTaurusDBHtapStarrocksLtsConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/starrocks/instances/logs/lts-configs"
		instanceID = d.Get("instance_id").(string)
		logType    = d.Get("log_type").(string)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildTaurusDBHtapStarrocksLtsConfigBodyParams(d),
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating TaurusDB HTAP StarRocks LTS config: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, logType))

	return resourceTaurusDBHtapStarrocksLtsConfigRead(ctx, d, meta)
}

func buildTaurusDBHtapStarrocksLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceTaurusDBHtapStarrocksLtsConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		logType    = d.Get("log_type").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	ltsConfig, err := GetTaurusDBHtapStarrocksLtsConfig(client, instanceId, logType)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving TaurusDB HTAP StarRocks LTS config")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("log_type", logType),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", ltsConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", ltsConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetTaurusDBHtapStarrocksLtsConfig(client *golangsdk.ServiceClient, instanceID, logType string) (interface{}, error) {
	listPath := client.Endpoint + "v3/{project_id}/starrocks/instances/logs/lts-configs"
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?instance_id=%s", instanceID)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error listing TaurusDB HTAP StarRocks LTS configs: %s", err)
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("instance_lts_configs|[0].lts_configs|[?log_type=='%s']|[0]", logType)
	ltsConfig := utils.PathSearch(searchPath, listRespBody, nil)
	if ltsConfig == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	enabled := utils.PathSearch("enabled", ltsConfig, false).(bool)
	if !enabled {
		return nil, golangsdk.ErrDefault404{}
	}
	return ltsConfig, nil
}

func resourceTaurusDBHtapStarrocksLtsConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/starrocks/instances/logs/lts-configs"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteTaurusDBHtapStarrocksLtsConfigBodyParams(d))

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// The instance not exists or the log type config does not exist.
		err := common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.280001")
		return common.CheckDeletedDiag(d, err, "error deleting TaurusDB HTAP StarRocks LTS config")
	}

	return resourceTaurusDBHtapStarrocksLtsConfigRead(ctx, d, meta)
}

func buildDeleteTaurusDBHtapStarrocksLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceTaurusDBHtapStarrocksLtsConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<log_type>, but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
