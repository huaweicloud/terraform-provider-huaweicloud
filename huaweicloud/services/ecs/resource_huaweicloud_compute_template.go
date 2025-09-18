package ecs

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

var templateNonUpdatableParams = []string{"name", "description", "version_description", "template_data",
	"template_data.*.flavor_id", "template_data.*.name", "template_data.*.description", "template_data.*.availability_zone_id",
	"template_data.*.enterprise_project_id", "template_data.*.auto_recovery", "template_data.*.os_profile",
	"template_data.*.os_profile.*.key_name", "template_data.*.os_profile.*.user_data",
	"template_data.*.os_profile.*.iam_agency_name", "template_data.*.os_profile.*.enable_monitoring_service",
	"template_data.*.security_group_ids", "template_data.*.network_interfaces", "template_data.*.block_device_mappings",
	"template_data.*.market_options", "template_data.*.market_options.*.market_type", "template_data.*.market_options.*.spot_options",
	"template_data.*.market_options.*.spot_options.*.spot_price",
	"template_data.*.market_options.*.spot_options.*.block_duration_minutes",
	"template_data.*.market_options.*.spot_options.*.instance_interruption_behavior", "template_data.*.internet_access",
	"template_data.*.internet_access.*.publicip", "template_data.*.internet_access.*.publicip.*.publicip_type",
	"template_data.*.internet_access.*.publicip.*.charging_mode", "template_data.*.internet_access.*.publicip.*.bandwidth",
	"template_data.*.internet_access.*.publicip.*.bandwidth.*.share_type",
	"template_data.*.internet_access.*.publicip.*.bandwidth.*.size",
	"template_data.*.internet_access.*.publicip.*.bandwidth.*.charge_mode",
	"template_data.*.internet_access.*.publicip.*.bandwidth.*.id", "template_data.*.metadata",
	"template_data.*.tag_options"}

// @API ECS POST /v3/{project_id}/launch-templates
// @API ECS GET /v3/{project_id}/launch-template-versions
// @API ECS DELETE /v3/{project_id}/launch-templates/{launch_template_id}
func ResourceComputeTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeTemplateCreate,
		ReadContext:   resourceComputeTemplateRead,
		UpdateContext: resourceComputeTemplateUpdate,
		DeleteContext: resourceComputeTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(templateNonUpdatableParams),

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
			"template_data": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func templateTemplateDataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"os_profile": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataOsProfileSchema(),
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     templateTemplateDataNetworkInterfacesSchema(),
			},
			"block_device_mappings": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     templateTemplateDataBlockDeviceMappingsSchema(),
			},
			"market_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataMarketOptionsSchema(),
			},
			"internet_access": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataInternetAccessSchema(),
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     templateTemplateDataTagOptionsSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataOsProfileSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"iam_agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_monitoring_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func templateTemplateDataNetworkInterfacesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"virsubnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"attachment": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataNetworkInterfacesAttachmentSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataNetworkInterfacesAttachmentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"device_index": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func templateTemplateDataBlockDeviceMappingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cmk_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"attachment": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataBlockDeviceMappingsAttachmentSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataBlockDeviceMappingsAttachmentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"boot_index": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func templateTemplateDataMarketOptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"market_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"spot_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataMarketOptionsSpotOptionsSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataMarketOptionsSpotOptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"spot_price": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"block_duration_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_interruption_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func templateTemplateDataInternetAccessSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"publicip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataInternetAccessPublicIpSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataInternetAccessPublicIpSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"publicip_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     templateTemplateDataInternetAccessPublicIpBandwidthSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataInternetAccessPublicIpBandwidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func templateTemplateDataTagOptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     templateTemplateDataTagOptionsTagsSchema(),
			},
		},
	}
	return &sc
}

func templateTemplateDataTagOptionsTagsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceComputeTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/launch-templates"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateTemplateBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("launch_template_id", createRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("error creating ECS template: launch_template_id is not found in the response")
	}

	d.SetId(templateId)

	return resourceComputeTemplateRead(ctx, d, meta)
}

func buildCreateTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                d.Get("name"),
		"template_data":       buildCreateTemplateTemplateDataBodyParams(d.Get("template_data")),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
	}
	return map[string]interface{}{
		"launch_template": bodyParams,
	}
}

func buildCreateTemplateTemplateDataBodyParams(templateDataRaw interface{}) map[string]interface{} {
	templateData := templateDataRaw.([]interface{})
	if len(templateData) == 0 {
		return nil
	}

	if v, ok := templateData[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"flavor_id":             utils.ValueIgnoreEmpty(v["flavor_id"]),
			"name":                  utils.ValueIgnoreEmpty(v["name"]),
			"description":           utils.ValueIgnoreEmpty(v["description"]),
			"availability_zone_id":  utils.ValueIgnoreEmpty(v["availability_zone_id"]),
			"enterprise_project_id": utils.ValueIgnoreEmpty(v["enterprise_project_id"]),
			"auto_recovery":         utils.ValueIgnoreEmpty(v["auto_recovery"]),
			"os_profile":            buildCreateTemplateTemplateDataOsProfileBodyParams(v["os_profile"]),
			"security_group_ids":    utils.ValueIgnoreEmpty(v["security_group_ids"].(*schema.Set).List()),
			"network_interfaces":    buildCreateTemplateTemplateDataNetworkInterfacesBodyParams(v["network_interfaces"]),
			"block_device_mappings": buildCreateTemplateTemplateDataBlockDeviceMappingsBodyParams(v["block_device_mappings"]),
			"market_options":        buildCreateTemplateTemplateDataMarketOptionsBodyParams(v["market_options"]),
			"internet_access":       buildCreateTemplateTemplateDataInternetAccessBodyParams(v["internet_access"]),
			"metadata":              utils.ValueIgnoreEmpty(v["metadata"]),
			"tag_options":           buildCreateTemplateTemplateDataTagOptionsBodyParams(v["tag_options"]),
		}

		return bodyParams
	}
	return nil
}

func buildCreateTemplateTemplateDataOsProfileBodyParams(osProfileRaw interface{}) map[string]interface{} {
	osProfile := osProfileRaw.([]interface{})
	if len(osProfile) == 0 {
		return nil
	}

	if v, ok := osProfile[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"key_name":                  utils.ValueIgnoreEmpty(v["key_name"]),
			"user_data":                 utils.ValueIgnoreEmpty(v["user_data"]),
			"iam_agency_name":           utils.ValueIgnoreEmpty(v["iam_agency_name"]),
			"enable_monitoring_service": utils.ValueIgnoreEmpty(v["enable_monitoring_service"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataNetworkInterfacesBodyParams(networkInterfacesRaw interface{}) []interface{} {
	networkInterfaces := networkInterfacesRaw.(*schema.Set)
	if networkInterfaces.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, networkInterfaces.Len())
	for _, networkInterface := range networkInterfaces.List() {
		if v, ok := networkInterface.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"virsubnet_id": utils.ValueIgnoreEmpty(v["virsubnet_id"]),
				"attachment":   buildCreateTemplateTemplateDataNetworkInterfacesAttachmentBodyParams(v["attachment"]),
			})
		}
	}

	return rst
}

func buildCreateTemplateTemplateDataNetworkInterfacesAttachmentBodyParams(attachmentRaw interface{}) map[string]interface{} {
	attachment := attachmentRaw.([]interface{})
	if len(attachment) == 0 {
		return nil
	}

	if v, ok := attachment[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"device_index": utils.ValueIgnoreEmpty(v["device_index"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataBlockDeviceMappingsBodyParams(blockDeviceMappingsRaw interface{}) []interface{} {
	blockDeviceMappings := blockDeviceMappingsRaw.(*schema.Set)
	if blockDeviceMappings.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, blockDeviceMappings.Len())
	for _, blockDeviceMapping := range blockDeviceMappings.List() {
		if v, ok := blockDeviceMapping.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"source_id":   utils.ValueIgnoreEmpty(v["source_id"]),
				"source_type": utils.ValueIgnoreEmpty(v["source_type"]),
				"encrypted":   utils.ValueIgnoreEmpty(v["encrypted"]),
				"cmk_id":      utils.ValueIgnoreEmpty(v["cmk_id"]),
				"volume_type": utils.ValueIgnoreEmpty(v["volume_type"]),
				"volume_size": utils.ValueIgnoreEmpty(v["volume_size"]),
				"attachment":  buildCreateTemplateTemplateDataBlockDeviceMappingsAttachmentBodyParams(v["attachment"]),
			})
		}
	}

	return rst
}

func buildCreateTemplateTemplateDataBlockDeviceMappingsAttachmentBodyParams(attachmentRaw interface{}) map[string]interface{} {
	attachment := attachmentRaw.([]interface{})
	if len(attachment) == 0 {
		return nil
	}

	if v, ok := attachment[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"boot_index": utils.ValueIgnoreEmpty(v["boot_index"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataMarketOptionsBodyParams(marketOptionsRaw interface{}) map[string]interface{} {
	marketOptions := marketOptionsRaw.([]interface{})
	if len(marketOptions) == 0 {
		return nil
	}

	if v, ok := marketOptions[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"market_type":  utils.ValueIgnoreEmpty(v["market_type"]),
			"spot_options": buildCreateTemplateTemplateDataMarketOptionsSpotOptionsBodyParams(v["spot_options"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataMarketOptionsSpotOptionsBodyParams(spotOptionsRaw interface{}) map[string]interface{} {
	spotOptions := spotOptionsRaw.([]interface{})
	if len(spotOptions) == 0 {
		return nil
	}

	if v, ok := spotOptions[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"spot_price":                     utils.ValueIgnoreEmpty(v["spot_price"]),
			"block_duration_minutes":         utils.ValueIgnoreEmpty(v["block_duration_minutes"]),
			"instance_interruption_behavior": utils.ValueIgnoreEmpty(v["instance_interruption_behavior"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataInternetAccessBodyParams(internetAccessRaw interface{}) map[string]interface{} {
	internetAccess := internetAccessRaw.([]interface{})
	if len(internetAccess) == 0 {
		return nil
	}

	if v, ok := internetAccess[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"publicip": buildCreateTemplateTemplateDataInternetAccessPublicIpBodyParams(v["publicip"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataInternetAccessPublicIpBodyParams(publicIpRaw interface{}) map[string]interface{} {
	publicIp := publicIpRaw.([]interface{})
	if len(publicIp) == 0 {
		return nil
	}

	if v, ok := publicIp[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"publicip_type": utils.ValueIgnoreEmpty(v["publicip_type"]),
			"charging_mode": utils.ValueIgnoreEmpty(v["charging_mode"]),
			"bandwidth":     buildCreateTemplateTemplateDataInternetAccessPublicIpBandwidthBodyParams(v["bandwidth"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataInternetAccessPublicIpBandwidthBodyParams(bandwidthRaw interface{}) map[string]interface{} {
	bandwidth := bandwidthRaw.([]interface{})
	if len(bandwidth) == 0 {
		return nil
	}

	if v, ok := bandwidth[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"share_type":  utils.ValueIgnoreEmpty(v["share_type"]),
			"size":        utils.ValueIgnoreEmpty(v["size"]),
			"charge_mode": utils.ValueIgnoreEmpty(v["charge_mode"]),
			"id":          utils.ValueIgnoreEmpty(v["id"]),
		}
		return rst
	}

	return nil
}

func buildCreateTemplateTemplateDataTagOptionsBodyParams(tagOptionsRaw interface{}) []interface{} {
	tagOptions := tagOptionsRaw.(*schema.Set)
	if tagOptions.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, tagOptions.Len())
	for _, tagOption := range tagOptions.List() {
		if v, ok := tagOption.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"tags": buildCreateTemplateTemplateDataTagOptionsTagsBodyParams(v["tags"]),
			})
		}
	}

	return rst
}

func buildCreateTemplateTemplateDataTagOptionsTagsBodyParams(tagsRaw interface{}) []interface{} {
	tags := tagsRaw.(*schema.Set)
	if tags.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, tags.Len())
	for _, tag := range tags.List() {
		if v, ok := tag.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"key":   utils.ValueIgnoreEmpty(v["key"]),
				"value": utils.ValueIgnoreEmpty(v["value"]),
			})
		}
	}

	return rst
}

func resourceComputeTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/launch-template-versions?launch_template_id={launch_template_id}"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{launch_template_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ECS template")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	template := utils.PathSearch("launch_template_versions[0]", getRespBody, nil)
	if template == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving ECS template")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("template_data", flattenTemplateTemplateData(template)),
		d.Set("version_description", utils.PathSearch("version_description", template, nil)),
		d.Set("version_number", utils.PathSearch("version_number", template, nil)),
		d.Set("version_id", utils.PathSearch("version_id", template, nil)),
		d.Set("created_at", utils.PathSearch("created_at", template, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTemplateTemplateData(template interface{}) []interface{} {
	curJson := utils.PathSearch("template_data", template, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"flavor_id":             utils.PathSearch("flavor_id", curJson, nil),
			"name":                  utils.PathSearch("name", curJson, nil),
			"description":           utils.PathSearch("description", curJson, nil),
			"availability_zone_id":  utils.PathSearch("availability_zone_id", curJson, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", curJson, nil),
			"auto_recovery":         utils.PathSearch("auto_recovery", curJson, nil),
			"os_profile":            flattenTemplateTemplateDataOsProfile(curJson),
			"security_group_ids":    utils.PathSearch("security_group_ids", curJson, nil),
			"network_interfaces":    flattenTemplateTemplateDataNetworkInterfaces(curJson),
			"block_device_mappings": flattenTemplateTemplateDataBlockDeviceMappings(curJson),
			"market_options":        flattenTemplateTemplateDataMarketOptions(curJson),
			"internet_access":       flattenTemplateTemplateDataInternetAccess(curJson),
			"metadata":              utils.PathSearch("metadata", curJson, nil),
			"tag_options":           flattenTemplateTemplateTagOptionsData(curJson),
		},
	}
	return rst
}

func flattenTemplateTemplateDataOsProfile(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("os_profile", templateData, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"key_name":                  utils.PathSearch("key_name", curJson, nil),
			"user_data":                 utils.PathSearch("user_data", curJson, nil),
			"iam_agency_name":           utils.PathSearch("iam_agency_name", curJson, nil),
			"enable_monitoring_service": utils.PathSearch("enable_monitoring_service", curJson, nil),
		},
	}
	return rst
}

func flattenTemplateTemplateDataNetworkInterfaces(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("network_interfaces", templateData, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"virsubnet_id": utils.PathSearch("virsubnet_id", v, nil),
			"attachment":   flattenTemplateTemplateDataNetworkInterfacesAttachment(v),
		})
	}
	return rst
}

func flattenTemplateTemplateDataNetworkInterfacesAttachment(networkInterfaces interface{}) []interface{} {
	curJson := utils.PathSearch("attachment", networkInterfaces, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"device_index": utils.PathSearch("device_index", curJson, nil),
		},
	}
	return rst
}

func flattenTemplateTemplateDataBlockDeviceMappings(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("block_device_mappings", templateData, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"source_id":   utils.PathSearch("source_id", v, nil),
			"source_type": utils.PathSearch("source_type", v, nil),
			"encrypted":   utils.PathSearch("encrypted", v, nil),
			"cmk_id":      utils.PathSearch("cmk_id", v, nil),
			"volume_type": utils.PathSearch("volume_type", v, nil),
			"volume_size": utils.PathSearch("volume_size", v, nil),
			"attachment":  flattenTemplateTemplateDataBlockDeviceMappingsAttachment(v),
		})
	}
	return rst
}

func flattenTemplateTemplateDataBlockDeviceMappingsAttachment(blockDeviceMappings interface{}) []interface{} {
	curJson := utils.PathSearch("attachment", blockDeviceMappings, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"boot_index": utils.PathSearch("boot_index", curJson, nil),
		},
	}
	return rst
}

func flattenTemplateTemplateDataMarketOptions(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("market_options", templateData, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"market_type":  utils.PathSearch("market_type", curJson, nil),
			"spot_options": flattenTemplateTemplateDataMarketOptionsSpotOptions(curJson),
		},
	}
	return rst
}

func flattenTemplateTemplateDataMarketOptionsSpotOptions(marketOptions interface{}) []interface{} {
	curJson := utils.PathSearch("spot_options", marketOptions, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"spot_price":                     utils.PathSearch("spot_price", curJson, nil),
			"block_duration_minutes":         utils.PathSearch("block_duration_minutes", curJson, nil),
			"instance_interruption_behavior": utils.PathSearch("instance_interruption_behavior", curJson, nil),
		},
	}
	return rst
}

func flattenTemplateTemplateDataInternetAccess(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("internet_access", templateData, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"publicip": flattenTemplateTemplateDataInternetAccessPublicIp(curJson),
		},
	}
	return rst
}

func flattenTemplateTemplateDataInternetAccessPublicIp(internetAccess interface{}) []interface{} {
	curJson := utils.PathSearch("publicip", internetAccess, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"publicip_type": utils.PathSearch("publicip_type", curJson, nil),
			"charging_mode": utils.PathSearch("charging_mode", curJson, nil),
			"bandwidth":     flattenTemplateTemplateDataInternetAccessPublicIpBandwidth(curJson),
		},
	}
	return rst
}

func flattenTemplateTemplateDataInternetAccessPublicIpBandwidth(publicIp interface{}) []interface{} {
	curJson := utils.PathSearch("bandwidth", publicIp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"share_type":  utils.PathSearch("share_type", curJson, nil),
			"size":        utils.PathSearch("size", curJson, nil),
			"charge_mode": utils.PathSearch("charge_mode", curJson, nil),
			"id":          utils.PathSearch("id", curJson, nil),
		},
	}
	return rst
}

func flattenTemplateTemplateTagOptionsData(templateData interface{}) []interface{} {
	curJson := utils.PathSearch("tag_options", templateData, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"tags": flattenTemplateTemplateDataTagOptionsTags(v),
		})
	}
	return rst
}

func flattenTemplateTemplateDataTagOptionsTags(tagOptions interface{}) []interface{} {
	curJson := utils.PathSearch("tags", tagOptions, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func resourceComputeTemplateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/launch-templates/{launch_template_id}"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{launch_template_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ECS template")
	}

	return nil
}
