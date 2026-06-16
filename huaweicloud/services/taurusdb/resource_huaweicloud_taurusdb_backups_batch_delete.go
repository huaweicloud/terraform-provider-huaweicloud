package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var backupsBatchDeleteNoneUpdatableParams = []string{"backup_ids"}

// @API TaurusDB DELETE /v3/{project_id}/backups
func ResourceTaurusDBBackupsBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBBackupsBatchDeleteCreate,
		ReadContext:   resourceTaurusDBBackupsBatchDeleteRead,
		UpdateContext: resourceTaurusDBBackupsBatchDeleteUpdate,
		DeleteContext: resourceTaurusDBBackupsBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(backupsBatchDeleteNoneUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"backup_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"success_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceTaurusDBBackupsBatchDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/backups"
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBBackupsBatchDeleteBodyParams(d))

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error batch deleting TaurusDB backups: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("success_count", int(utils.PathSearch("success_count", deleteRespBody, float64(0)).(float64))),
		d.Set("failed_count", int(utils.PathSearch("failed_count", deleteRespBody, float64(0)).(float64))),
		d.Set("failed_results", flattenBatchDeleteBackupFailedResults(deleteRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteGaussDBBackupsBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"backup_ids": utils.ExpandToStringList(d.Get("backup_ids").([]interface{})),
	}
	return bodyParams
}

func flattenBatchDeleteBackupFailedResults(respBody interface{}) []interface{} {
	failedResults := utils.PathSearch("failed_results", respBody, make([]interface{}, 0)).([]interface{})
	if len(failedResults) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(failedResults))
	for _, item := range failedResults {
		rst = append(rst, map[string]interface{}{
			"backup_id":  utils.PathSearch("backup_id", item, nil),
			"error_code": utils.PathSearch("error_code", item, nil),
			"error_msg":  utils.PathSearch("error_msg", item, nil),
		})
	}
	return rst
}

func resourceTaurusDBBackupsBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBBackupsBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
func resourceTaurusDBBackupsBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting TaurusDB backups batch delete resource is not supported. The TaurusDB backups " +
		"batch delete resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
