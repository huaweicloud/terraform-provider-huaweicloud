package lts

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eps"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const EPSTagKey string = "_sys_enterprise_project_id"

// @API LTS POST /v2/{project_id}/groups/{log_group_id}/streams
// @API LTS POST /v1/{project_id}/{resource_type}/{resource_id}/tags/action
// @API LTS POST /v1.0/{project_id}/lts/favorite
// @API LTS GET /v2/{project_id}/groups/{log_group_id}/streams
// @API LTS PUT /v2/{project_id}/groups/{log_group_id}/streams-ttl/{log_stream_id}
// @API LTS DELETE /v1.0/{project_id}/lts/favorite/{fav_res_id}
// @API LTS DELETE /v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}
func ResourceLTSStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamCreate,
		ReadContext:   resourceStreamRead,
		UpdateContext: resourceStreamUpdate,
		DeleteContext: resourceStreamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStreamImportState,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"is_favorite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to favorite the log stream.`,
			},
			// Attributes
			"filter_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/groups/{log_group_id}/streams"
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{log_group_id}", d.Get("group_id").(string))

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateStreamBodyParams(cfg, d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating LTS stream: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	streamId := utils.PathSearch("log_stream_id", respBody, "").(string)
	if streamId == "" {
		return diag.Errorf("unable to find the LTS stream ID from the API response")
	}

	d.SetId(streamId)

	if _, ok := d.GetOk("tags"); ok {
		streamId := d.Id()
		err = updateTags(client, "topics", streamId, d)
		if err != nil {
			return diag.Errorf("error creating tags of log stream %s: %s", streamId, err)
		}
	}

	if d.Get("is_favorite").(bool) {
		if err = favoriteLogStream(client, cfg, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceStreamRead(ctx, d, meta)
}

func buildCreateStreamBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_stream_name": d.Get("stream_name"),
		"ttl_in_days":     utils.ValueIgnoreEmpty(d.Get("ttl_in_days")),
	}

	userNoPermission := []string{"EPS.0004"}
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId == "" {
		return bodyParams
	}
	epsInfo, err := getEnterpriseProjectById(cfg, cfg.GetRegion(d), epsId)
	// If we catch error 403, it means that the user does not have EPS permissions, return immediately.
	if parsedErr := eps.ParseQueryError403(err, userNoPermission, "No permission, skip the enterprise project query"); parsedErr == nil {
		// Unable to set enterprise project ID for log stream via parameter 'enterprise_project_id' and
		// 'tags._sys_enterprise_project_id'. Currently, only parameter 'enterprise_project_name' is available.
		bodyParams["enterprise_project_name"] = utils.PathSearch("enterprise_project.name", epsInfo, nil)
	}
	// If not, insert enterprise project name into bodyParams.
	return bodyParams
}

func getEnterpriseProjectById(cfg *config.Config, region, epsId string) (interface{}, error) {
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return nil, fmt.Errorf("error creating EPS client: %s", err)
	}

	httpUrl := "v1.0/enterprise-projects/{enterprise_project_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{enterprise_project_id}", epsId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func GetLogStreams(client *golangsdk.ServiceClient, logGroupId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups/{log_group_id}/streams"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_group_id}", logGroupId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func GetLogStreamById(client *golangsdk.ServiceClient, logGroupId, streamId string) (interface{}, error) {
	streams, err := GetLogStreams(client, logGroupId)
	if err != nil {
		return nil, err
	}

	streamResult := utils.PathSearch(fmt.Sprintf("log_streams|[?log_stream_id=='%s']|[0]", streamId), streams, nil)
	if streamResult == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "v2/{project_id}/groups/{log_group_id}/streams",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the log stream (%s) has been deleted", streamId)),
			},
		}
	}

	return streamResult, nil
}

func resourceStreamRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		streamId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	streamResult, err := GetLogStreamById(client, d.Get("group_id").(string), streamId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find log stream by its ID (%s)", streamId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("stream_name", utils.PathSearch("log_stream_name", streamResult, nil)),
		d.Set("ttl_in_days", utils.PathSearch("ttl_in_days", streamResult, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("tag._sys_enterprise_project_id", streamResult, nil)),
		d.Set("tags", ignoreSysEpsTag(utils.PathSearch("tag", streamResult, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("is_favorite", utils.PathSearch("is_favorite", streamResult, nil)),
		d.Set("filter_count", utils.PathSearch("filter_count", streamResult, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("creation_time", streamResult, float64(0)).(float64))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		streamId = d.Id()
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	if d.HasChange("tags") {
		err = updateTags(client, "topics", streamId, d)
		if err != nil {
			return diag.Errorf("error updating tags of log stream (%s): %s", streamId, err)
		}
	}

	if d.HasChange("ttl_in_days") {
		if err = updateStreamTTL(client, d.Get("group_id").(string), streamId, d.Get("ttl_in_days")); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("is_favorite") {
		isFavorite := d.Get("is_favorite").(bool)
		if isFavorite {
			err = favoriteLogStream(client, cfg, d)
		} else {
			err = removeFavoriteLogStream(client, streamId)
		}

		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceStreamRead(ctx, d, meta)
}

func updateStreamTTL(client *golangsdk.ServiceClient, logGroupId, logStreamId string, ttlInDays interface{}) error {
	httpUrl := "v2/{project_id}/groups/{log_group_id}/streams-ttl/{log_stream_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{log_group_id}", logGroupId)
	updatePath = strings.ReplaceAll(updatePath, "{log_stream_id}", logStreamId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"ttl_in_days": ttlInDays,
		},
	}
	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating ttl_in_days of the log stream (%s): %s", logStreamId, err)
	}
	return nil
}

func favoriteLogStream(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1.0/{project_id}/lts/favorite"
		logStreamId = d.Id()
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFavoriteLogStreamBodyParams(cfg, d, logStreamId)),
	}
	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("unable to favorite log stream (%s): %s", logStreamId, err)
	}
	return nil
}

func buildFavoriteLogStreamBodyParams(cfg *config.Config, d *schema.ResourceData, logStreamId string) map[string]interface{} {
	return map[string]interface{}{
		"eps_id":                 utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"favorite_resource_id":   logStreamId,
		"favorite_resource_type": "log_stream",
		"log_group_id":           d.Get("group_id"),
		"log_stream_id":          logStreamId,
		"log_stream_name":        d.Get("stream_name"),
		// This parameter must be set to `true`, otherwise the favoriting will not take effect.
		"is_global": true,
	}
}

func removeFavoriteLogStream(client *golangsdk.ServiceClient, logStreamId string) error {
	httpUrl := "v1.0/{project_id}/lts/favorite/{fav_res_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{fav_res_id}", logStreamId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err := client.Request("DELETE", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error removing favorite log stream (%s): %s", logStreamId, err)
	}
	return nil
}

func resourceStreamDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}"
		streamId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{log_group_id}", d.Get("group_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{log_stream_id}", streamId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting log stream")
	}
	return nil
}

func resourceStreamImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<group_id>/<stream_id>', but '%s'", d.Id())
	}

	groupID := parts[0]
	streamID := parts[1]

	d.SetId(streamID)
	mErr := multierror.Append(nil,
		d.Set("group_id", groupID),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
