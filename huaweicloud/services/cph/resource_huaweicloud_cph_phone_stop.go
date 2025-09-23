package cph

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var PhoneStopNonUpdatableParams = []string{"phone_id"}

// @API CPH POST /v1/{project_id}/cloud-phone/phones/batch-stop
// @API CPH GET /v1/{project_id}/cloud-phone/phones/{phone_id}
func ResourcePhoneStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePhoneStopCreate,
		UpdateContext: resourcePhoneStopUpdate,
		ReadContext:   resourcePhoneStopRead,
		DeleteContext: resourcePhoneStopDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(PhoneStopNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"phone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the CPH phone.`,
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

func resourcePhoneStopCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// createPhoneStop: create CPH phone stop
	createPhoneStopHttpUrl := "v1/{project_id}/cloud-phone/phones/batch-stop"
	createPhoneStopPath := client.Endpoint + createPhoneStopHttpUrl
	createPhoneStopPath = strings.ReplaceAll(createPhoneStopPath, "{project_id}", client.ProjectID)

	createPhoneStopOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"phone_ids": []string{d.Get("phone_id").(string)},
		},
	}

	createPhoneStopResp, err := client.Request("POST", createPhoneStopPath, &createPhoneStopOpt)
	if err != nil {
		return diag.Errorf("error creating CPH phone stop: %s", err)
	}

	resp, err := utils.FlattenResponse(createPhoneStopResp)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("jobs|[0].phone_id", resp, "").(string)
	if id == "" {
		return diag.Errorf("Unable to find the phone ID from the API response")
	}
	d.SetId(id)

	errorCode := utils.PathSearch("jobs|[0].error_code", resp, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("jobs|[0].error_msg", resp, "").(string)
		return diag.Errorf("failed to stop CPH phone (phone_id: %s), error_code: %s, error_msg: %s", id, errorCode, errorMsg)
	}

	err = checkPhoneStopStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePhoneStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhoneStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhoneStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CPH stop resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkPhoneStopStatus(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      phoneStateRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CPH phone stop to be completed: %s", err)
	}
	return nil
}

func phoneStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPhoneRespBody, err := getPhoneDetail(client, id)
		if err != nil {
			return nil, "ERROR", err
		}

		// Status is 8, indicates the phone is closed.
		phoneStatus := utils.PathSearch("status", getPhoneRespBody, float64(0)).(float64)
		if int(phoneStatus) == 8 {
			return getPhoneRespBody, "COMPLETED", nil
		}
		if int(phoneStatus) == -9 {
			return getPhoneRespBody, "ERROR", fmt.Errorf("failed to turn off cloud phone: %s", id)
		}
		return getPhoneRespBody, "PENDING", nil
	}
}

func getPhoneDetail(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getPhoneHttpUrl := "v1/{project_id}/cloud-phone/phones/{phone_id}"
	getPhonePath := client.Endpoint + getPhoneHttpUrl
	getPhonePath = strings.ReplaceAll(getPhonePath, "{project_id}", client.ProjectID)
	getPhonePath = strings.ReplaceAll(getPhonePath, "{phone_id}", id)

	getPhoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPhoneResp, err := client.Request("GET", getPhonePath, &getPhoneOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getPhoneResp)
}
