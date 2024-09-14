package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/traffic-mirror-sessions
// @API VPC GET /v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}
// @API VPC PUT /v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}
// @API VPC PUT /v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}/remove-sources
// @API VPC PUT /v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}/add-sources
// @API VPC DELETE /v3/{project_id}/vpc/traffic-mirror-sessions/{traffic_mirror_session_id}
func ResourceTrafficMirrorSession() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficMirrorSessionCreate,
		ReadContext:   resourceTrafficMirrorSessionRead,
		UpdateContext: resourceTrafficMirrorSessionUpdate,
		DeleteContext: resourceTrafficMirrorSessionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_sources": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"traffic_mirror_target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_mirror_target_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_network_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"packet_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getTrafficMirrorSessionHttpUrl(d *schema.ResourceData, client *golangsdk.ServiceClient) string {
	trafficMirrorSessionPath := client.ResourceBaseURL() + "vpc/traffic-mirror-sessions"
	if d.Id() != "" {
		trafficMirrorSessionPath += "/" + d.Id()
	}
	return trafficMirrorSessionPath
}

func removeTrafficMirrorSourcesHttpUrl(d *schema.ResourceData, client *golangsdk.ServiceClient) string {
	trafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
	return trafficMirrorSessionPath + "/remove-sources"
}

func addTrafficMirrorSourcesHttpUrl(d *schema.ResourceData, client *golangsdk.ServiceClient) string {
	trafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
	return trafficMirrorSessionPath + "/add-sources"
}

func buildCreateTrafficMirrorSessionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"traffic_mirror_session": map[string]interface{}{
			"name":                       d.Get("name"),
			"description":                d.Get("description"),
			"traffic_mirror_filter_id":   d.Get("traffic_mirror_filter_id"),
			"traffic_mirror_sources":     d.Get("traffic_mirror_sources"),
			"traffic_mirror_target_id":   d.Get("traffic_mirror_target_id"),
			"traffic_mirror_target_type": d.Get("traffic_mirror_target_type"),
			"virtual_network_id":         utils.ValueIgnoreEmpty(d.Get("virtual_network_id")),
			"packet_length":              utils.ValueIgnoreEmpty(d.Get("packet_length")),
			"priority":                   utils.ValueIgnoreEmpty(d.Get("priority")),
			"enabled":                    utils.ValueIgnoreEmpty(d.Get("enabled")),
			"type":                       utils.ValueIgnoreEmpty(d.Get("type")),
		},
	}
	return bodyParams
}

func resourceTrafficMirrorSessionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	ctreateTrafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
	createTrafficMirrorSessionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createTrafficMirrorSessionOpt.JSONBody = utils.RemoveNil(buildCreateTrafficMirrorSessionBodyParams(d))
	createTrafficMirrorSessionResp, err := client.Request("POST", ctreateTrafficMirrorSessionPath, &createTrafficMirrorSessionOpt)
	if err != nil {
		return diag.Errorf("error creating traffic mirror filter: %s", err)
	}

	createTrafficMirrorSessionRespBody, err := utils.FlattenResponse(createTrafficMirrorSessionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("traffic_mirror_session.id", createTrafficMirrorSessionRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating traffic mirror filter: ID is not found in API response")
	}
	d.SetId(id)

	return resourceTrafficMirrorSessionRead(ctx, d, meta)
}

func resourceTrafficMirrorSessionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
	getTrafficMirrorSessionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorSessionResp, err := client.Request("GET", getTrafficMirrorSessionPath, &getTrafficMirrorSessionOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC traffic mirror filter")
	}

	resp, err := utils.FlattenResponse(getTrafficMirrorSessionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("traffic_mirror_session.name", resp, nil)),
		d.Set("description", utils.PathSearch("traffic_mirror_session.description", resp, nil)),
		d.Set("traffic_mirror_filter_id", utils.PathSearch("traffic_mirror_session.traffic_mirror_filter_id", resp, nil)),
		d.Set("traffic_mirror_sources", utils.PathSearch("traffic_mirror_session.traffic_mirror_sources", resp, nil)),
		d.Set("traffic_mirror_target_id", utils.PathSearch("traffic_mirror_session.traffic_mirror_target_id", resp, nil)),
		d.Set("traffic_mirror_target_type", utils.PathSearch("traffic_mirror_session.traffic_mirror_target_type", resp, nil)),
		d.Set("virtual_network_id", utils.PathSearch("traffic_mirror_session.virtual_network_id", resp, nil)),
		d.Set("packet_length", utils.PathSearch("traffic_mirror_session.packet_length", resp, nil)),
		d.Set("priority", utils.PathSearch("traffic_mirror_session.priority", resp, nil)),
		d.Set("enabled", utils.PathSearch("traffic_mirror_session.enabled", resp, nil)),
		d.Set("type", utils.PathSearch("traffic_mirror_session.type", resp, nil)),
		d.Set("created_at", utils.PathSearch("traffic_mirror_session.created_at", resp, nil)),
		d.Set("updated_at", utils.PathSearch("traffic_mirror_session.updated_at", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateTrafficMirrorSessionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"traffic_mirror_session": map[string]interface{}{
			"name":                       d.Get("name"),
			"description":                d.Get("description"),
			"traffic_mirror_filter_id":   d.Get("traffic_mirror_filter_id"),
			"traffic_mirror_target_id":   d.Get("traffic_mirror_target_id"),
			"traffic_mirror_target_type": d.Get("traffic_mirror_target_type"),
			"virtual_network_id":         d.Get("virtual_network_id"),
			"packet_length":              d.Get("packet_length"),
			"priority":                   d.Get("priority"),
			"enabled":                    d.Get("enabled"),
			"type":                       d.Get("type"),
		},
	}
	return bodyParams
}

func buildTrafficMirrorSourcesBodyParams(sources []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"traffic_mirror_session": map[string]interface{}{
			"traffic_mirror_sources": sources,
		},
	}
	return bodyParams
}

func resourceTrafficMirrorSessionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	if d.HasChangeExcept("traffic_mirror_sources") {
		updateTrafficMirrorSessionOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateTrafficMirrorSessionOpts.JSONBody = utils.RemoveNil(buildUpdateTrafficMirrorSessionBodyParams(d))
		updateTrafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
		_, err = client.Request("PUT", updateTrafficMirrorSessionPath, &updateTrafficMirrorSessionOpts)
		if err != nil {
			return diag.Errorf("error updating traffic mirror filter: %s", err)
		}
	}

	if d.HasChange("traffic_mirror_sources") {
		oldSources, newSources := d.GetChange("traffic_mirror_sources")
		if len(oldSources.([]interface{})) > 0 {
			removeTrafficMirrorSourcesOpts := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			removeTrafficMirrorSourcesOpts.JSONBody = utils.RemoveNil(buildTrafficMirrorSourcesBodyParams(oldSources.([]interface{})))
			removeTrafficMirrorSourcesPath := removeTrafficMirrorSourcesHttpUrl(d, client)
			_, err = client.Request("PUT", removeTrafficMirrorSourcesPath, &removeTrafficMirrorSourcesOpts)
			if err != nil {
				return diag.Errorf("error updating traffic mirror filter: %s", err)
			}
		}

		if len(newSources.([]interface{})) > 0 {
			addTrafficMirrorSourcesOpts := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			addTrafficMirrorSourcesOpts.JSONBody = utils.RemoveNil(buildTrafficMirrorSourcesBodyParams(newSources.([]interface{})))
			addTrafficMirrorSourcesPath := addTrafficMirrorSourcesHttpUrl(d, client)
			_, err = client.Request("PUT", addTrafficMirrorSourcesPath, &addTrafficMirrorSourcesOpts)
			if err != nil {
				return diag.Errorf("error updating traffic mirror filter: %s", err)
			}
		}
	}

	return resourceTrafficMirrorSessionRead(ctx, d, meta)
}

func resourceTrafficMirrorSessionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	deleteTrafficMirrorSessionPath := getTrafficMirrorSessionHttpUrl(d, client)
	deleteTrafficMirrorSessionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteTrafficMirrorSessionPath, &deleteTrafficMirrorSessionOpt)
	if err != nil {
		return diag.Errorf("error deleting traffic mirror filter: %s", err)
	}
	return nil
}
