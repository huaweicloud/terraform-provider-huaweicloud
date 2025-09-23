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

// @API VPC POST /v3/{project_id}/vpc/traffic-mirror-filters
// @API VPC GET /v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}
// @API VPC PUT /v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}
// @API VPC DELETE /v3/{project_id}/vpc/traffic-mirror-filters/{traffic_mirror_filter_id}
func ResourceTrafficMirrorFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficMirrorFilterCreate,
		ReadContext:   resourceTrafficMirrorFilterRead,
		UpdateContext: resourceTrafficMirrorFilterUpdate,
		DeleteContext: resourceTrafficMirrorFilterDelete,

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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func getTrafficMirrorFilterHttpUrl(d *schema.ResourceData, client *golangsdk.ServiceClient) string {
	trafficMirrorFilterPath := client.ResourceBaseURL() + "vpc/traffic-mirror-filters"
	if d.Id() != "" {
		trafficMirrorFilterPath += "/" + d.Id()
	}
	return trafficMirrorFilterPath
}

func buildTrafficMirrorFilterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"traffic_mirror_filter": map[string]interface{}{
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}
	return bodyParams
}

func resourceTrafficMirrorFilterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	ctreateTrafficMirrorFilterPath := getTrafficMirrorFilterHttpUrl(d, client)
	createTrafficMirrorFilterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createTrafficMirrorFilterOpt.JSONBody = utils.RemoveNil(buildTrafficMirrorFilterBodyParams(d))
	createTrafficMirrorFilterResp, err := client.Request("POST", ctreateTrafficMirrorFilterPath, &createTrafficMirrorFilterOpt)
	if err != nil {
		return diag.Errorf("error creating traffic mirror filter: %s", err)
	}

	createTrafficMirrorFilterRespBody, err := utils.FlattenResponse(createTrafficMirrorFilterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("traffic_mirror_filter.id", createTrafficMirrorFilterRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating traffic mirror filter: ID is not found in API response")
	}
	d.SetId(id)

	return resourceTrafficMirrorFilterRead(ctx, d, meta)
}

func resourceTrafficMirrorFilterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorFilterPath := getTrafficMirrorFilterHttpUrl(d, client)
	getTrafficMirrorFilterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorFilterResp, err := client.Request("GET", getTrafficMirrorFilterPath, &getTrafficMirrorFilterOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC traffic mirror filter")
	}

	getTrafficMirrorFilterRespBody, err := utils.FlattenResponse(getTrafficMirrorFilterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("traffic_mirror_filter.name", getTrafficMirrorFilterRespBody, nil)),
		d.Set("description", utils.PathSearch("traffic_mirror_filter.description", getTrafficMirrorFilterRespBody, nil)),
		d.Set("created_at", utils.PathSearch("traffic_mirror_filter.created_at", getTrafficMirrorFilterRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("traffic_mirror_filter.updated_at", getTrafficMirrorFilterRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTrafficMirrorFilterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateTrafficMirrorFilterOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateTrafficMirrorFilterOpts.JSONBody = utils.RemoveNil(buildTrafficMirrorFilterBodyParams(d))
		updateTrafficMirrorFilterPath := getTrafficMirrorFilterHttpUrl(d, client)
		_, err = client.Request("PUT", updateTrafficMirrorFilterPath, &updateTrafficMirrorFilterOpts)
		if err != nil {
			return diag.Errorf("error updating traffic mirror filter: %s", err)
		}
	}

	return resourceTrafficMirrorFilterRead(ctx, d, meta)
}

func resourceTrafficMirrorFilterDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	deleteTrafficMirrorFilterPath := getTrafficMirrorFilterHttpUrl(d, client)
	deleteTrafficMirrorFilterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteTrafficMirrorFilterPath, &deleteTrafficMirrorFilterOpt)
	if err != nil {
		return diag.Errorf("error deleting traffic mirror filter: %s", err)
	}
	return nil
}
