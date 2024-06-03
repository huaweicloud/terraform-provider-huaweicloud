package dcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/restores
// @API DCS GET /v2/{project_id}/instances/{instance_id}/restores
func ResourceDcsRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsRestoreCreate,
		ReadContext:   resourceDcsRestoreRead,
		DeleteContext: resourceDcsRestoreDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
				Description: "Specifies the ID of the DCS instance to be restored.",
			},
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the backup ID used to restore the DCS instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the description of the DCS instance restoration.`,
			},
			"restore_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the restoration record.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the restoration record created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the restoration record completed.`,
			},
		},
	}
}

func resourceDcsRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		restoreDcsCreateHttpUrl = "v2/{project_id}/instances/{instance_id}/restores"
		restoreDcsCreateProduct = "dcs"
	)
	restoreDcsCreateClient, err := cfg.NewServiceClient(restoreDcsCreateProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	restoreDcsCreatePath := restoreDcsCreateClient.Endpoint + restoreDcsCreateHttpUrl
	restoreDcsCreatePath = strings.ReplaceAll(restoreDcsCreatePath, "{project_id}", restoreDcsCreateClient.ProjectID)
	restoreDcsCreatePath = strings.ReplaceAll(restoreDcsCreatePath, "{instance_id}", instanceID)

	restoreDcsCreateOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	restoreDcsCreateOpt.JSONBody = utils.RemoveNil(buildCreateRestoreBodyParams(d))

	restoreDcsCreateResp, err := restoreDcsCreateClient.Request("POST", restoreDcsCreatePath, &restoreDcsCreateOpt)
	if err != nil {
		return diag.Errorf("error restoring the instance (%s): %s", instanceID, err)
	}

	restoreDcsCreateRespBody, err := utils.FlattenResponse(restoreDcsCreateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	restoreId := utils.PathSearch("restore_id", restoreDcsCreateRespBody, "")
	if restoreId == nil {
		return diag.Errorf("unable to find the restore_id of the instance (%s): %s", instanceID, err)
	}

	d.SetId(restoreId.(string))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"waiting", "restoring"},
		Target:       []string{"succeed"},
		Refresh:      restoreRecordRefreshFunc(instanceID, d.Id(), restoreDcsCreateClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for restoring the instance (%s) to complete: %s", instanceID, err)
	}

	return resourceDcsRestoreRead(ctx, d, meta)
}

func buildCreateRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"backup_id": d.Get("backup_id"),
		"remark":    utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceDcsRestoreRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceID := d.Get("instance_id").(string)

	var (
		getRestoreProduct = "dcs"
		mErr              *multierror.Error
	)

	getRestoreClient, err := cfg.NewServiceClient(getRestoreProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	restoreRecord, err := GetRestoreRecord(instanceID, d.Id(), getRestoreClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if restoreRecord == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("backup_id", utils.PathSearch("backup_id", restoreRecord, nil)),
		d.Set("description", utils.PathSearch("restore_remark", restoreRecord, nil)),
		d.Set("restore_name", utils.PathSearch("restore_name", restoreRecord, nil)),
		d.Set("created_at", utils.PathSearch("created_at", restoreRecord, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", restoreRecord, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcsRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration record is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func restoreRecordRefreshFunc(instanceID, restoreID string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		restoreRecord, err := GetRestoreRecord(instanceID, restoreID, client)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", restoreRecord, "")
		return restoreRecord, status.(string), nil
	}
}

func GetRestoreRecord(instanceID, restoreID string, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getRestoreHttpUrl = "v2/{project_id}/instances/{instance_id}/restores"
		getRestorePath    string
		offset            int
	)

	getRestoreBasePath := client.Endpoint + getRestoreHttpUrl
	getRestoreBasePath = strings.ReplaceAll(getRestoreBasePath, "{project_id}", client.ProjectID)
	getRestoreBasePath = strings.ReplaceAll(getRestoreBasePath, "{instance_id}", instanceID)
	getRestoreOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	for {
		getRestoreQueryParams := fmt.Sprintf("?limit=1000&offset=%v", offset)
		getRestorePath = getRestoreBasePath + getRestoreQueryParams
		getRestoreResp, err := client.Request("GET", getRestorePath, &getRestoreOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving restoration records: %s", err)
		}

		getRestoreRespBody, err := utils.FlattenResponse(getRestoreResp)
		if err != nil {
			return nil, err
		}

		restoreRecords := utils.PathSearch("restore_record_response", getRestoreRespBody, make([]interface{}, 0)).([]interface{})
		totalNum := utils.PathSearch("total_num", getRestoreRespBody, float64(0))

		item := utils.PathSearch(fmt.Sprintf("restore_record_response|[?restore_id =='%s']|[0]", restoreID), getRestoreRespBody, nil)
		if item != nil {
			return item, nil
		}

		offset += len(restoreRecords)
		if offset >= int(totalNum.(float64)) {
			break
		}
	}
	return nil, fmt.Errorf("error getting restoration record by restore_id (%s)", restoreID)
}
