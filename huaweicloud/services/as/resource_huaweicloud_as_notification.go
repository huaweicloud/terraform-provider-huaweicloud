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
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
		product        = "autoscaling"
		scalingGroupID = d.Get("scaling_group_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	putPath := client.Endpoint + httpUrl
	putPath = strings.ReplaceAll(putPath, "{project_id}", client.ProjectID)
	putPath = strings.ReplaceAll(putPath, "{scaling_group_id}", scalingGroupID)
	putOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateOrUpdateASNotificationBodyParams(d),
	}

	_, err = client.Request("PUT", putPath, &putOpt)
	if err != nil {
		return diag.Errorf("error creating or updating AS notification: %s", err)
	}

	d.SetId(d.Get("topic_urn").(string))
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
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		mErr           *multierror.Error
		httpUrl        = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}"
		product        = "autoscaling"
		scalingGroupID = d.Get("scaling_group_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_group_id}", scalingGroupID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// When the group does not exist, the response error information of the detailed API is as follows:
		// {"error": {"code": "AS.2007","message": "The AS group does not exist."}}.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007"), "error retrieving AS notification")
	}

	getASNotificationRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationMap := filterTargetASNotificationByTopicUrn(getASNotificationRespBody, d.Id())
	if len(notificationMap) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

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
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		deleteUrl      = "autoscaling-api/v1/{project_id}/scaling_notification/{scaling_group_id}/{topic_urn}"
		product        = "autoscaling"
		scalingGroupID = d.Get("scaling_group_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	deletePath := client.Endpoint + deleteUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{scaling_group_id}", scalingGroupID)
	deletePath = strings.ReplaceAll(deletePath, "{topic_urn}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// When the group does not exist, the response error message of the delete API is:
		// {"error": {"code": "AS.2007","message": "The AS group does not exist."}}.
		// When AS notification does not exist, the response HTTP status code of the delete API is 404
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007"), "error deleting AS notification")
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
