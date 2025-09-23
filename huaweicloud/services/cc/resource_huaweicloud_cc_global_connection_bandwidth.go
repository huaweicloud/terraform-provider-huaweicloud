package cc

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

// @API CC POST /v3/{domain_id}/gcb/gcbandwidths
// @API CC GET /v3/{domain_id}/gcb/gcbandwidths/{id}
// @API CC PUT /v3/{domain_id}/gcb/gcbandwidths/{id}
// @API CC DELETE /v3/{domain_id}/gcb/gcbandwidths/{id}
// @API CC POST /v3/gcb/{resource_id}/tags/create
// @API CC POST /v3/gcb/{resource_id}/tags/delete
func ResourceGlobalConnectionBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalConnectionBandwidthCreate,
		UpdateContext: resourceGlobalConnectionBandwidthUpdate,
		ReadContext:   resourceGlobalConnectionBandwidthRead,
		DeleteContext: resourceGlobalConnectionBandwidthDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bordercross": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sla_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remote_area": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"spec_code_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"binding_service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_share": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"frozen": {
				Type:     schema.TypeBool,
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

func resourceGlobalConnectionBandwidthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	createGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths"
	createGCBPath := client.Endpoint + createGCBHttpUrl
	createGCBPath = strings.ReplaceAll(createGCBPath, "{domain_id}", cfg.DomainID)
	createGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"globalconnection_bandwidth": utils.RemoveNil(buildCreateGCBParams(d, cfg.GetEnterpriseProjectID(d))),
		},
	}

	createGCBResp, err := client.Request("POST", createGCBPath, &createGCBOpt)
	if err != nil {
		return diag.Errorf("error creating global connection bandwidth: %s", err)
	}
	createGCBRespBody, err := utils.FlattenResponse(createGCBResp)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("globalconnection_bandwidth.id", createGCBRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating global connection bandwidth: ID is not found in API response")
	}

	d.SetId(id)

	if v, ok := d.GetOk("binding_service"); ok && v.(string) != "ALL" {
		err = updateGCB(client, d, cfg.DomainID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGlobalConnectionBandwidthRead(ctx, d, meta)
}

func buildCreateGCBParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  d.Get("name"),
		"type":                  d.Get("type"),
		"size":                  d.Get("size"),
		"charge_mode":           d.Get("charge_mode"),
		"bordercross":           d.Get("bordercross"),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		"sla_level":             utils.ValueIgnoreEmpty(d.Get("sla_level")),
		"local_area":            utils.ValueIgnoreEmpty(d.Get("local_area")),
		"remote_area":           utils.ValueIgnoreEmpty(d.Get("remote_area")),
		"spec_code_id":          utils.ValueIgnoreEmpty(d.Get("spec_code_id")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return params
}

func resourceGlobalConnectionBandwidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	getGCBPath := client.Endpoint + getGCBHttpUrl
	getGCBPath = strings.ReplaceAll(getGCBPath, "{domain_id}", cfg.DomainID)
	getGCBPath = strings.ReplaceAll(getGCBPath, "{id}", d.Id())
	getGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGCBResp, err := client.Request("GET", getGCBPath, &getGCBOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global connection bandwidth")
	}
	getGCBRespBody, err := utils.FlattenResponse(getGCBResp)
	if err != nil {
		return diag.FromErr(err)
	}
	getGCBRespBody = utils.PathSearch("globalconnection_bandwidth", getGCBRespBody, nil)
	if getGCBRespBody == nil {
		return diag.Errorf("error getting global connection bandwidth: %s is not found in API response",
			"globalconnection_bandwidth")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", getGCBRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getGCBRespBody, nil)),
		d.Set("size", int(utils.PathSearch("size", getGCBRespBody, float64(0)).(float64))),
		d.Set("charge_mode", utils.PathSearch("charge_mode", getGCBRespBody, nil)),
		d.Set("bordercross", utils.PathSearch("bordercross", getGCBRespBody, false)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getGCBRespBody, nil)),
		d.Set("sla_level", utils.PathSearch("sla_level", getGCBRespBody, nil)),
		d.Set("local_area", utils.PathSearch("local_site_code", getGCBRespBody, nil)),
		d.Set("remote_area", utils.PathSearch("remote_site_code", getGCBRespBody, nil)),
		d.Set("spec_code_id", utils.PathSearch("spec_code_id", getGCBRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getGCBRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", getGCBRespBody, nil))),
		d.Set("enable_share", utils.PathSearch("enable_share", getGCBRespBody, false)),
		d.Set("binding_service", utils.PathSearch("binding_service", getGCBRespBody, nil)),
		d.Set("instances", flattenInstances(utils.PathSearch("instances", getGCBRespBody, make([]interface{}, 0)))),
		d.Set("frozen", utils.PathSearch("frozen", getGCBRespBody, false)),
		d.Set("created_at", utils.PathSearch("created_at", getGCBRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getGCBRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstances(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, val := range rawArray {
		params := map[string]interface{}{
			"id":     utils.PathSearch("id", val, nil),
			"type":   utils.PathSearch("type", val, nil),
			"region": utils.PathSearch("region_id", val, nil),
		}
		rst = append(rst, params)
	}
	return rst
}

func resourceGlobalConnectionBandwidthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	updateChanges := []string{
		"name",
		"size",
		"charge_mode",
		"sla_level",
		"binding_service",
		"spec_code_id",
		"description",
	}

	if d.HasChanges(updateChanges...) {
		err = updateGCB(client, d, cfg.DomainID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("tags") {
		tagErr := updateTags(client, d)
		if tagErr != nil {
			return diag.Errorf("error updating tags of global connection bandwidth (%s): %s", d.Id(), tagErr)
		}
	}

	return resourceGlobalConnectionBandwidthRead(ctx, d, meta)
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	manageTagsHttpUrl := "v3/gcb/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", d.Id())
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

func updateGCB(client *golangsdk.ServiceClient, d *schema.ResourceData, domainID string) error {
	updateGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	updateGCBPath := client.Endpoint + updateGCBHttpUrl
	updateGCBPath = strings.ReplaceAll(updateGCBPath, "{domain_id}", domainID)
	updateGCBPath = strings.ReplaceAll(updateGCBPath, "{id}", d.Id())
	updateGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"globalconnection_bandwidth": utils.RemoveNil(buildUpdateGCBParams(d)),
		},
	}

	_, err := client.Request("PUT", updateGCBPath, &updateGCBOpt)
	if err != nil {
		return fmt.Errorf("error updating global connection bandwidth: %s", err)
	}

	return nil
}

func buildUpdateGCBParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":            d.Get("name"),
		"size":            d.Get("size"),
		"charge_mode":     d.Get("charge_mode"),
		"sla_level":       utils.ValueIgnoreEmpty(d.Get("sla_level")),
		"binding_service": utils.ValueIgnoreEmpty(d.Get("binding_service")),
		"spec_code_id":    utils.ValueIgnoreEmpty(d.Get("spec_code_id")),
		"description":     d.Get("description"),
	}
	return params
}

func resourceGlobalConnectionBandwidthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	deleteGCBPath := client.Endpoint + deleteGCBHttpUrl
	deleteGCBPath = strings.ReplaceAll(deleteGCBPath, "{domain_id}", cfg.DomainID)
	deleteGCBPath = strings.ReplaceAll(deleteGCBPath, "{id}", d.Id())
	deleteGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteGCBPath, &deleteGCBOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting global connection bandwidth")
	}

	return nil
}
