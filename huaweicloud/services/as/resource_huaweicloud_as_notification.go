// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AS
// ---------------------------------------------------------------

package as

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS PUT /autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}
// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}/{topic_urn}
func ResourceAsNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsNotificationPut,
		UpdateContext: resourceAsNotificationPut,
		ReadContext:   resourceAsNotificationRead,
		DeleteContext: resourceAsNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAsNotificationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the AS group ID.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the unique topic URN of the SMN.`,
			},
			"events": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"SCALING_UP", "SCALING_UP_FAIL", "SCALING_DOWN", "SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL",
					}, false),
				},
				Description: `Specifies the topic scene of AS group.`,
			},
			"topic_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The topic name in SMN.`,
			},
		},
	}
}

func resourceAsNotificationPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// putASNotification: put an AS notification.
	var (
		putASNotificationHttpUrl = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
		putASNotificationProduct = "autoscaling"
	)
	putASNotificationClient, err := cfg.NewServiceClient(putASNotificationProduct, region)
	if err != nil {
		return diag.Errorf("error creating AutoScaling Client: %s", err)
	}

	putASNotificationPath := putASNotificationClient.Endpoint + putASNotificationHttpUrl
	putASNotificationPath = strings.ReplaceAll(putASNotificationPath, "{project_id}",
		putASNotificationClient.ProjectID)
	putASNotificationPath = strings.ReplaceAll(putASNotificationPath, "{scaling_group_id}",
		fmt.Sprintf("%v", d.Get("scaling_group_id")))

	putASNotificationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	putASNotificationOpt.JSONBody = buildCreateOrUpdateASNotificationBodyParams(d)
	_, err = putASNotificationClient.Request("PUT", putASNotificationPath, &putASNotificationOpt)
	if err != nil {
		return diag.Errorf("error creating or updating AS notification: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	d.SetId(topicUrn)
	return resourceAsNotificationRead(ctx, d, meta)
}

func buildCreateOrUpdateASNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"topic_urn":   d.Get("topic_urn"),
		"topic_scene": utils.ValueIgnoreEmpty(d.Get("events")),
	}
	return bodyParams
}

func resourceAsNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getASNotification: Query the AS notification.
	var (
		getASNotificationHttpUrl = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
		getASNotificationProduct = "autoscaling"
	)
	getASNotificationClient, err := cfg.NewServiceClient(getASNotificationProduct, region)
	if err != nil {
		return diag.Errorf("error creating AutoScaling Client: %s", err)
	}

	getASNotificationPath := getASNotificationClient.Endpoint + getASNotificationHttpUrl
	getASNotificationPath = strings.ReplaceAll(getASNotificationPath, "{project_id}",
		getASNotificationClient.ProjectID)
	getASNotificationPath = strings.ReplaceAll(getASNotificationPath, "{scaling_group_id}",
		fmt.Sprintf("%v", d.Get("scaling_group_id")))

	getASNotificationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getASNotificationResp, err := getASNotificationClient.Request("GET", getASNotificationPath,
		&getASNotificationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AS notification")
	}

	getASNotificationRespBody, err := utils.FlattenResponse(getASNotificationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationMap := filterTargetASNotificationByTopicUrn(getASNotificationRespBody, d.Id())
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("events", notificationMap["events"]),
		d.Set("topic_name", notificationMap["topic_name"]),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterTargetASNotificationByTopicUrn(resp interface{}, topicUrn string) map[string]interface{} {
	notificationMap := make(map[string]interface{})
	if resp == nil {
		return notificationMap
	}

	curJson := utils.PathSearch("topics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		urn := utils.PathSearch("topic_urn", v, "")
		if topicUrn == urn.(string) {
			notificationMap["events"] = utils.PathSearch("topic_scene", v, nil)
			notificationMap["topic_name"] = utils.PathSearch("topic_name", v, "")
			break
		}
	}
	return notificationMap
}

func resourceAsNotificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteASNotification: Delete the AS notification.
	var (
		deleteASNotificationHttpUrl = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}" +
			"/{topic_urn}"
		deleteASNotificationProduct = "autoscaling"
	)
	deleteASNotificationClient, err := cfg.NewServiceClient(deleteASNotificationProduct, region)
	if err != nil {
		return diag.Errorf("error creating AutoScaling Client: %s", err)
	}

	deleteASNotificationPath := deleteASNotificationClient.Endpoint + deleteASNotificationHttpUrl
	deleteASNotificationPath = strings.ReplaceAll(deleteASNotificationPath, "{project_id}",
		deleteASNotificationClient.ProjectID)
	deleteASNotificationPath = strings.ReplaceAll(deleteASNotificationPath, "{scaling_group_id}",
		fmt.Sprintf("%v", d.Get("scaling_group_id")))
	deleteASNotificationPath = strings.ReplaceAll(deleteASNotificationPath, "{topic_urn}",
		fmt.Sprintf("%v", d.Id()))

	deleteASNotificationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteASNotificationClient.Request("DELETE", deleteASNotificationPath, &deleteASNotificationOpt)
	if err != nil {
		return diag.Errorf("error deleting AS notification: %s", err)
	}
	return nil
}

func resourceAsNotificationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <scaling_group_id>/<topic_urn>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("scaling_group_id", parts[0]),
		d.Set("topic_urn", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
