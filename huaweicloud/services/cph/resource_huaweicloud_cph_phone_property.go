package cph

import (
	"context"
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

var PhonePropertyNonUpdatableParams = []string{"phones", "phones.*.phone_id", "phones.*.property"}

// @API CPH POST /v1/{project_id}/cloud-phone/phones/batch-update-property
// @API CPH GET /v1/{project_id}/cloud-phone/jobs/{job_id}
func ResourcePhoneProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePhonePropertyCreate,
		UpdateContext: resourcePhonePropertyUpdate,
		ReadContext:   resourcePhonePropertyRead,
		DeleteContext: resourcePhonePropertyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(PhoneResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"phones": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the CPH phones.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the phone ID.`,
						},
						"property": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the phone property, the format is json string.`,
						},
					},
				},
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

func resourcePhonePropertyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// updatePhoneProperty: update CPH phone property
	updatePhonePropertyHttpUrl := "v1/{project_id}/cloud-phone/phones/batch-update-property"
	updatePhonePropertyPath := client.Endpoint + updatePhonePropertyHttpUrl
	updatePhonePropertyPath = strings.ReplaceAll(updatePhonePropertyPath, "{project_id}", client.ProjectID)

	updatePhonePropertyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updatePhonePropertyOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"phones": d.Get("phones"),
	})
	updatePhonePropertyResp, err := client.Request("POST", updatePhonePropertyPath, &updatePhonePropertyOpt)
	if err != nil {
		return diag.Errorf("error updating CPH phone property: %s", err)
	}

	resp, err := utils.FlattenResponse(updatePhonePropertyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("jobs|[0].job_id", resp, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}
	d.SetId(id)

	errorCode := utils.PathSearch("jobs|[0].error_code", resp, "").(string)
	if errorCode != "" {
		phoneId := utils.PathSearch("jobs|[0].phone_id", resp, "").(string)
		errorMsg := utils.PathSearch("jobs|[0].error_msg", resp, "").(string)
		return diag.Errorf("failed to updating CPH phone property (phone_id: %s), error_code: %s, error_msg: %s", phoneId, errorCode, errorMsg)
	}

	err = checkPhonePropertyJobStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePhonePropertyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhonePropertyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhonePropertyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CPH phone property resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkPhonePropertyJobStatus(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStatusRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
