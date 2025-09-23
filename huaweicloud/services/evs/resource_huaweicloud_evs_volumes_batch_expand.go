package evs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var volumesBatchExpandNonUpdatableParams = []string{
	"volumes",
	"volumes.*.id",
	"volumes.*.new_size",
	"is_auto_pay",
}

// @API EVS POST /v5/{project_id}/volumes/batch-extend
// @API EVS GET /v1/{project_id}/jobs/{job_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
func ResourceVolumesBatchExpand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumesBatchExpandCreate,
		ReadContext:   resourceVolumesBatchExpandRead,
		UpdateContext: resourceVolumesBatchExpandUpdate,
		DeleteContext: resourceVolumesBatchExpandDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(volumesBatchExpandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"new_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"is_auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
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

func buildVolumesBatchExpandBodyParams(d *schema.ResourceData) map[string]interface{} {
	volumes := d.Get("volumes").([]interface{})
	volumesList := make([]map[string]interface{}, 0, len(volumes))

	for _, v := range volumes {
		volume, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		volumesList = append(volumesList, map[string]interface{}{
			"id":       volume["id"],
			"new_size": volume["new_size"],
		})
	}

	bodyParams := map[string]interface{}{
		"bss_param": map[string]interface{}{
			"is_auto_pay": d.Get("is_auto_pay"),
		},
		"volumes": volumesList,
	}

	return bodyParams
}

func getExpandJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying EVS job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitingForExpandVolumeJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getExpandJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in EVS job (%s) detail API response", jobID)
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			if status == "FAIL" {
				return respBody, status, fmt.Errorf("the EVS job (%s) status is FAIL, the fail reason is: %s",
					jobID, utils.PathSearch("fail_reason", respBody, "").(string))
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceVolumesBatchExpandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/volumes/batch-extend"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVolumesBatchExpandBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch expanding EVS volumes: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	if jobID := utils.PathSearch("job_id", respBody, "").(string); jobID != "" {
		if err := waitingForExpandVolumeJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for the job (%s) to succeed: %s", jobID, err)
		}
	}

	if orderID := utils.PathSearch("order_id", respBody, "").(string); orderID != "" {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the order (%s) is not completed while expanding EVS volume: %v", orderID, err)
		}
		if _, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return nil
}

func resourceVolumesBatchExpandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumesBatchExpandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVolumesBatchExpandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to expand volumes.
Deleting this resource will not reset the expanded volumes, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
