package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceDatabaseBatchUpgradeNoneUpdatableParams = []string{
	"databases_instance_infos", "delay",
	"databases_instance_infos.*.instance_id",
	"databases_instance_infos.*.current_version",
}

// @API TaurusDB POST /v3/{project_id}/instances/database-version/upgrade
func ResourceTaurusDBInstanceDatabaseBatchUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBInstanceDatabaseBatchUpgradeCreate,
		ReadContext:   resourceTaurusDBInstanceDatabaseBatchUpgradeRead,
		UpdateContext: resourceTaurusDBInstanceDatabaseBatchUpgradeUpdate,
		DeleteContext: resourceTaurusDBInstanceDatabaseBatchUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceDatabaseBatchUpgradeNoneUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"databases_instance_infos": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     instanceDatabaseBatchUpgradeInstanceInfosSchema(),
			},
			"delay": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
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

func instanceDatabaseBatchUpgradeInstanceInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTaurusDBInstanceDatabaseBatchUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/database-version/upgrade"
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
		JSONBody:         utils.RemoveNil(buildCreateInstanceDatabaseBatchUpgradeBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error upgrading TaurusDB instance database version in batches: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resp := utils.PathSearch("resp", createRespBody, "").(string)
	if resp != "success" {
		return diag.Errorf("error upgrading TaurusDB instance database version in batches: the response is not success, got: %s", resp)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id.String())

	return nil
}

func buildCreateInstanceDatabaseBatchUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	databasesInstanceInfos := buildInstanceDatabaseBatchUpgradeInstanceInfosBodyParams(d.Get("databases_instance_infos"))
	bodyParams := map[string]interface{}{
		"databases_instance_infos": databasesInstanceInfos,
		"delay":                    d.Get("delay"),
	}
	return bodyParams
}

func buildInstanceDatabaseBatchUpgradeInstanceInfosBodyParams(rawParams interface{}) []interface{} {
	if rawParams == nil {
		return nil
	}
	rawArray := rawParams.([]interface{})
	result := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if item, ok := v.(map[string]interface{}); ok {
			result = append(result, map[string]interface{}{
				"instance_id":     item["instance_id"],
				"current_version": item["current_version"],
			})
		}
	}
	return result
}

func resourceTaurusDBInstanceDatabaseBatchUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBInstanceDatabaseBatchUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBInstanceDatabaseBatchUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting instance database batch upgrade resource is not supported. The resource is only removed from the state," +
		" the TaurusDB instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
