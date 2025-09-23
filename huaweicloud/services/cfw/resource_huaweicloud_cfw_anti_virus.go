package cfw

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

var antiVirusNonUpdatableParams = []string{
	"object_id",
}

// @API CFW PUT /v1/{project_id}/anti-virus/switch
// @API CFW GET /v1/{project_id}/anti-virus/switch
// @API CFW PUT /v1/{project_id}/anti-virus/rule
// @API CFW GET /v1/{project_id}/anti-virus/rule
func ResourceAntiVirus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntiVirusCreate,
		ReadContext:   resourceAntiVirusRead,
		UpdateContext: resourceAntiVirusUpdate,
		DeleteContext: resourceAntiVirusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(antiVirusNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"scan_protocol_configs": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the scan protocol configurations.`,
				Elem:        scanProtocolConfigSchema(),
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

func scanProtocolConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `The antivirus action.`,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"protocol_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The protocol type.`,
			},
		},
	}
	return &sc
}

func resourceAntiVirusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAntiVirusSwitchHttpUrl = "v1/{project_id}/anti-virus/switch"
		createAntiVirusProduct       = "cfw"
	)
	createAntiVirusClient, err := cfg.NewServiceClient(createAntiVirusProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	createAntiVirusSwitchPath := createAntiVirusClient.Endpoint + createAntiVirusSwitchHttpUrl
	createAntiVirusSwitchPath = strings.ReplaceAll(createAntiVirusSwitchPath, "{project_id}",
		createAntiVirusClient.ProjectID)

	createAntiVirusSwitchOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	objectId := d.Get("object_id").(string)
	createAntiVirusSwitchOpt.JSONBody = buildAntiVirusSwitchBodyParams(objectId, 1)
	_, err = createAntiVirusClient.Request("PUT", createAntiVirusSwitchPath, &createAntiVirusSwitchOpt)
	if err != nil {
		return diag.Errorf("error creating CFW anti virus: %s", err)
	}

	err = enableAntiVirusConfig(createAntiVirusClient, d.Get("scan_protocol_configs").(*schema.Set).List(), d.Get("object_id").(string))
	if err != nil {
		return diag.Errorf("error creating CFW anti virus: %s", err)
	}

	d.SetId(objectId)

	return resourceAntiVirusRead(ctx, d, meta)
}

func buildAntiVirusSwitchBodyParams(objectId string, status int) map[string]interface{} {
	return map[string]interface{}{
		"anti_virus_status": status,
		"object_id":         objectId,
	}
}

func enableAntiVirusConfig(client *golangsdk.ServiceClient, configs []interface{}, objectID string) error {
	enableAntiVirusConfigHttpUrl := "v1/{project_id}/anti-virus/rule"
	enableAntiVirusConfigPath := client.Endpoint + enableAntiVirusConfigHttpUrl
	enableAntiVirusConfigPath = strings.ReplaceAll(enableAntiVirusConfigPath, "{project_id}", client.ProjectID)

	enableAntiVirusConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	enableAntiVirusConfigOpt.JSONBody = buildEnableAntiVirusConfigBodyParams(configs, objectID)
	_, err := client.Request("PUT", enableAntiVirusConfigPath, &enableAntiVirusConfigOpt)
	return err
}

func buildEnableAntiVirusConfigBodyParams(params []interface{}, objectID string) map[string]interface{} {
	antiVirusConfigs := make([]map[string]interface{}, len(params))
	for i, param := range params {
		configMap := param.(map[string]interface{})
		antiVirusConfigs[i] = map[string]interface{}{
			"action":        configMap["action"].(int),
			"protocol_type": configMap["protocol_type"].(int),
		}
	}
	return map[string]interface{}{
		"object_id":             objectID,
		"scan_protocol_configs": antiVirusConfigs,
	}
}

func resourceAntiVirusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	getAntiVirusProduct := "cfw"
	getAntiVirusClient, err := cfg.NewServiceClient(getAntiVirusProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	objectID := d.Id()
	antiVirusConfigs, err := GetAntiVirusConfigs(getAntiVirusClient, objectID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW anti virus")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("object_id", objectID),
		d.Set("scan_protocol_configs", flattenAntiVirusConfigs(antiVirusConfigs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAntiVirusConfigs(param interface{}) []interface{} {
	rawConfigs := param.([]interface{})
	configs := make([]interface{}, len(rawConfigs))
	for i, rawConfig := range rawConfigs {
		configs[i] = map[string]interface{}{
			"action":        utils.PathSearch("action", rawConfig, nil),
			"protocol_type": utils.PathSearch("protocol_type", rawConfig, nil),
		}
	}
	return configs
}

func GetAntiVirusConfigs(client *golangsdk.ServiceClient, objectID string) (interface{}, error) {
	getAntiVirusSwitchHttpUrl := "v1/{project_id}/anti-virus/switch"
	getAntiVirusSwitchPath := client.Endpoint + getAntiVirusSwitchHttpUrl
	getAntiVirusSwitchPath = strings.ReplaceAll(getAntiVirusSwitchPath, "{project_id}", client.ProjectID)
	getAntiVirusSwitchPath += "?object_id=" + objectID

	getAntiVirusSwitchOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getAntiVirusSwitchPath, &getAntiVirusSwitchOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	status := int(utils.PathSearch("data.anti_virus_status", respBody, float64(2)).(float64))
	// If the status is off or not found, consider the resource as non-existent.
	if status == 0 || status == 2 {
		return nil, golangsdk.ErrDefault404{}
	}

	getAntiVirusConfigHttpUrl := "v1/{project_id}/anti-virus/rule"
	getAntiVirusConfigPath := client.Endpoint + getAntiVirusConfigHttpUrl
	getAntiVirusConfigPath = strings.ReplaceAll(getAntiVirusConfigPath, "{project_id}", client.ProjectID)
	// The value of engine_type does not affect the query results.
	getAntiVirusConfigPath += "?engine_type=1"

	getAntiVirusConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	antiVirusConfigs := make([]interface{}, 0)
	offset := 0
	for {
		path := fmt.Sprintf("%s&object_id=%s&limit=50&offset=%d", getAntiVirusConfigPath, objectID, offset)
		resp, err := client.Request("GET", path, &getAntiVirusConfigOpt)

		if err != nil {
			return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		expression := "data.scan_protocol_configs[?action == `0` || action == `1`]"
		configs := utils.PathSearch(expression, respBody, make([]interface{}, 0))
		antiVirusConfigs = append(antiVirusConfigs, configs.([]interface{})...)

		offset += 50
		total := utils.PathSearch("data.total", respBody, float64(0))
		if int(total.(float64)) <= offset {
			break
		}
	}

	// If no enabled antivirus configurations are found, consider the resource as non-existent.
	if len(antiVirusConfigs) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return antiVirusConfigs, nil
}

func resourceAntiVirusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAntiVirusProduct := "cfw"
	updateAntiVirusClient, err := cfg.NewServiceClient(updateAntiVirusProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	objectId := d.Get("object_id").(string)
	rawOldConfigs, rawNewConfigs := d.GetChange("scan_protocol_configs")
	enableConfigSet := rawNewConfigs.(*schema.Set).Difference(rawOldConfigs.(*schema.Set))
	disableConfigSet := rawOldConfigs.(*schema.Set).Difference(rawNewConfigs.(*schema.Set))

	if disableConfigSet.Len() > 0 {
		err = disableAntiVirusConfig(updateAntiVirusClient, disableConfigSet.List(), objectId)
		if err != nil {
			return diag.Errorf("error updating CFW anti virus: %s", err)
		}
	}

	if enableConfigSet.Len() > 0 {
		err = enableAntiVirusConfig(updateAntiVirusClient, enableConfigSet.List(), objectId)
		if err != nil {
			return diag.Errorf("error updating CFW anti virus: %s", err)
		}
	}
	return resourceAntiVirusRead(ctx, d, meta)
}

func disableAntiVirusConfig(client *golangsdk.ServiceClient, configs []interface{}, objectID string) error {
	disableAntiVirusConfigHttpUrl := "v1/{project_id}/anti-virus/rule"
	disableAntiVirusConfigPath := client.Endpoint + disableAntiVirusConfigHttpUrl
	disableAntiVirusConfigPath = strings.ReplaceAll(disableAntiVirusConfigPath, "{project_id}", client.ProjectID)

	disableAntiVirusConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	disableAntiVirusConfigOpt.JSONBody = buildDisableAntiVirusConfigBodyParams(configs, objectID)
	_, err := client.Request("PUT", disableAntiVirusConfigPath, &disableAntiVirusConfigOpt)
	return err
}

func buildDisableAntiVirusConfigBodyParams(params []interface{}, objectID string) map[string]interface{} {
	antiVirusConfigs := make([]map[string]interface{}, len(params))
	for i, param := range params {
		configMap := param.(map[string]interface{})
		antiVirusConfigs[i] = map[string]interface{}{
			"action":        2,
			"protocol_type": configMap["protocol_type"].(int),
		}
	}
	return map[string]interface{}{
		"object_id":             objectID,
		"scan_protocol_configs": antiVirusConfigs,
	}
}

func resourceAntiVirusDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAntiVirusSwitchHttpUrl = "v1/{project_id}/anti-virus/switch"
		deleteAntiVirusProduct       = "cfw"
	)

	deleteAntiVirusClient, err := cfg.NewServiceClient(deleteAntiVirusProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	configs := d.Get("scan_protocol_configs").(*schema.Set).List()
	objectId := d.Get("object_id").(string)
	err = disableAntiVirusConfig(deleteAntiVirusClient, configs, objectId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW anti virus",
		)
	}

	deleteAntiVirusSwitchPath := deleteAntiVirusClient.Endpoint + deleteAntiVirusSwitchHttpUrl
	deleteAntiVirusSwitchPath = strings.ReplaceAll(deleteAntiVirusSwitchPath, "{project_id}",
		deleteAntiVirusClient.ProjectID)

	deleteAntiVirusSwitchOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteAntiVirusSwitchOpt.JSONBody = buildAntiVirusSwitchBodyParams(objectId, 0)
	_, err = deleteAntiVirusClient.Request("PUT", deleteAntiVirusSwitchPath, &deleteAntiVirusSwitchOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW alarm configuration",
		)
	}

	return nil
}
