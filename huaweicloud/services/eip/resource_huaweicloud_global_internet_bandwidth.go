package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/geip/internet-bandwidths
// @API EIP DELETE /v3/{domain_id}/geip/internet-bandwidths/{id}
// @API EIP GET /v3/{domain_id}/geip/internet-bandwidths/{id}
// @API EIP PUT /v3/{domain_id}/geip/internet-bandwidths/{id}
// @API EIP POST /v3/internet-bandwidth/{resource_id}/tags/delete
// @API EIP POST /v3/internet-bandwidth/{resource_id}/tags/create
func ResourceGlobalInternetBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInternetBandwidthCreate,
		ReadContext:   resourceInternetBandwidthRead,
		UpdateContext: resourceInternetBandwidthUpdate,
		DeleteContext: resourceInternetBandwidthDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"access_site": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ingress_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ratio_95peak": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"frozen_info": {
				Type:     schema.TypeString,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInternetBandwidthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	createInternetBandwidthHttpUrl := "v3/{domain_id}/geip/internet-bandwidths"
	createInternetBandwidthPath := client.Endpoint + createInternetBandwidthHttpUrl
	createInternetBandwidthPath = strings.ReplaceAll(createInternetBandwidthPath, "{domain_id}", cfg.DomainID)

	createInternetBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"internet_bandwidth": utils.RemoveNil(buildCreateInternetBandwidthBodyParams(d, cfg.GetEnterpriseProjectID(d))),
		},
	}

	createInternetBandwidthResp, err := client.Request("POST", createInternetBandwidthPath, &createInternetBandwidthOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createInternetBandwidthRespBody, err := utils.FlattenResponse(createInternetBandwidthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("internet_bandwidth.id", createInternetBandwidthRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find internet bandwidth ID from the API response")
	}
	d.SetId(id)

	return resourceInternetBandwidthRead(ctx, d, meta)
}

func buildCreateInternetBandwidthBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"isp":                   d.Get("isp"),
		"access_site":           d.Get("access_site"),
		"charge_mode":           d.Get("charge_mode"),
		"size":                  d.Get("size"),
		"ingress_size":          utils.ValueIgnoreEmpty(d.Get("ingress_size")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"type":                  utils.ValueIgnoreEmpty(d.Get("type")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return bodyParams
}

func resourceInternetBandwidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getInternetBandwidthHttpUrl := "v3/{domain_id}/geip/internet-bandwidths/{id}"
	getInternetBandwidthPath := client.Endpoint + getInternetBandwidthHttpUrl
	getInternetBandwidthPath = strings.ReplaceAll(getInternetBandwidthPath, "{domain_id}", cfg.DomainID)
	getInternetBandwidthPath = strings.ReplaceAll(getInternetBandwidthPath, "{id}", d.Id())

	getInternetBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getInternetBandwidthResp, err := client.Request("GET", getInternetBandwidthPath, &getInternetBandwidthOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global internet bandwidth")
	}
	getInternetBandwidthRespBody, err := utils.FlattenResponse(getInternetBandwidthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	bandwidth := utils.PathSearch("internet_bandwidth", getInternetBandwidthRespBody, nil)
	if bandwidth == nil {
		return diag.Errorf("unable to find internet bandwidth from the API response")
	}

	mErr := multierror.Append(nil,
		d.Set("isp", utils.PathSearch("isp", bandwidth, nil)),
		d.Set("access_site", utils.PathSearch("access_site", bandwidth, nil)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", bandwidth, nil)),
		d.Set("size", utils.PathSearch("size", bandwidth, 0)),
		d.Set("ingress_size", utils.PathSearch("ingress_size", bandwidth, 0)),
		d.Set("description", utils.PathSearch("description", bandwidth, nil)),
		d.Set("name", utils.PathSearch("name", bandwidth, nil)),
		d.Set("type", utils.PathSearch("type", bandwidth, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", bandwidth, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", bandwidth, nil))),
		d.Set("ratio_95peak", utils.PathSearch("ratio_95peak", bandwidth, nil)),
		d.Set("frozen_info", utils.PathSearch("freezen_info", bandwidth, nil)),
		d.Set("created_at", utils.PathSearch("created_at", bandwidth, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", bandwidth, nil)),
		d.Set("status", utils.PathSearch("status", bandwidth, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting global internet bandwidth fields: %s", err)
	}
	return nil
}

func resourceInternetBandwidthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	updateChanges := []string{
		"charge_mode",
		"size",
		"description",
		"name",
		"ingress_size",
	}

	if d.HasChanges(updateChanges...) {
		updateInternetBandwidthHttpUrl := "v3/{domain_id}/geip/internet-bandwidths/{id}"
		updateInternetBandwidthPath := client.Endpoint + updateInternetBandwidthHttpUrl
		updateInternetBandwidthPath = strings.ReplaceAll(updateInternetBandwidthPath, "{domain_id}", cfg.DomainID)
		updateInternetBandwidthPath = strings.ReplaceAll(updateInternetBandwidthPath, "{id}", d.Id())

		updateInternetBandwidthOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"internet_bandwidth": utils.RemoveNil(buildUpdateInternetBandwidthBodyParams(d)),
			},
		}

		_, err = client.Request("PUT", updateInternetBandwidthPath, &updateInternetBandwidthOpt)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(client, d, "internet-bandwidth", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of global internet bandwidth (%s): %s", d.Id(), tagErr)
		}
	}

	return resourceInternetBandwidthRead(ctx, d, meta)
}

func buildUpdateInternetBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"charge_mode": d.Get("charge_mode"),
		"size":        d.Get("size"),
		"description": d.Get("description"),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
	}

	if d.Get("charge_mode").(string) != "95peak_guar" {
		bodyParams["ingress_size"] = utils.ValueIgnoreEmpty(d.Get("ingress_size"))
	}

	return bodyParams
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tagsType string, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	manageTagsHttpUrl := "v3/{tags_type}/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{tags_type}", tagsType)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", id)
	manageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	// remove old tags
	if len(oMap) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(oMap),
		}
		deleteTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "delete")
		_, err := client.Request("POST", deleteTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(nMap),
		}
		createTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "create")
		_, err := client.Request("POST", createTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceInternetBandwidthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deleteInternetBandwidthHttpUrl := "v3/{domain_id}/geip/internet-bandwidths/{id}"
	deleteInternetBandwidthPath := client.Endpoint + deleteInternetBandwidthHttpUrl
	deleteInternetBandwidthPath = strings.ReplaceAll(deleteInternetBandwidthPath, "{domain_id}", cfg.DomainID)
	deleteInternetBandwidthPath = strings.ReplaceAll(deleteInternetBandwidthPath, "{id}", d.Id())

	deleteInternetBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteInternetBandwidthPath, &deleteInternetBandwidthOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
