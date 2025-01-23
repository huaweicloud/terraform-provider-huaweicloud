package cfw

import (
	"context"
	"fmt"
	"strconv"
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

var alarmConfigNonUpdatableParams = []string{
	"fw_instance_id", "alarm_type",
}

// @API CFW PUT /v1/{project_id}/cfw/alarm/config
// @API CFW GET /v1/{project_id}/cfw/alarm/config
func ResourceAlarmConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmConfigCreateOrUpdate,
		ReadContext:   resourceAlarmConfigRead,
		UpdateContext: resourceAlarmConfigCreateOrUpdate,
		DeleteContext: resourceAlarmConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAlarmConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(alarmConfigNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall ID.`,
			},
			"alarm_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the alarm type.`,
			},
			"alarm_time_period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the alarm period.`,
			},
			"frequency_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the alarm triggering frequency.`,
			},
			"frequency_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the alarm frequency time range.`,
			},
			"severity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm severity.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm URN.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The alarm language.`,
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The username.`,
			},
		},
	}
}

func resourceAlarmConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAlarmConfigHttpUrl = "v1/{project_id}/cfw/alarm/config?fw_instance_id={fw_instance_id}"
		createAlarmConfigProduct = "cfw"
	)
	createAlarmConfigClient, err := cfg.NewServiceClient(createAlarmConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	fwInstanceID := d.Get("fw_instance_id").(string)
	createAlarmConfigPath := createAlarmConfigClient.Endpoint + createAlarmConfigHttpUrl
	createAlarmConfigPath = strings.ReplaceAll(createAlarmConfigPath, "{project_id}",
		createAlarmConfigClient.ProjectID)
	createAlarmConfigPath = strings.ReplaceAll(createAlarmConfigPath, "{fw_instance_id}", fwInstanceID)

	createAlarmConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAlarmConfigOpt.JSONBody = buildCreateAlarmConfigBodyParams(d)
	_, err = createAlarmConfigClient.Request("PUT", createAlarmConfigPath, &createAlarmConfigOpt)
	if err != nil {
		return diag.Errorf("error creating CFW alarm configuration: %s", err)
	}

	if d.IsNewResource() {
		alarmType := d.Get("alarm_type").(int)
		d.SetId(fwInstanceID + "/" + strconv.Itoa(alarmType))
	}

	return resourceAlarmConfigRead(ctx, d, meta)
}

func buildCreateAlarmConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"enable_status":     1,
		"alarm_type":        d.Get("alarm_type").(int),
		"alarm_time_period": d.Get("alarm_time_period").(int),
		"frequency_count":   d.Get("frequency_count").(int),
		"frequency_time":    d.Get("frequency_time").(int),
		"severity":          d.Get("severity").(string),
		"topic_urn":         d.Get("topic_urn").(string),
	}
}

func resourceAlarmConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	readAlarmConfigProduct := "cfw"
	readAlarmConfigClient, err := cfg.NewServiceClient(readAlarmConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	fwInstanceID := d.Get("fw_instance_id").(string)
	alarmType := d.Get("alarm_type").(int)
	alarm, err := GetAlarmConfig(readAlarmConfigClient, fwInstanceID, alarmType)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW alarm configuration")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("alarm_type", utils.PathSearch("alarm_type", alarm, nil)),
		d.Set("alarm_time_period", utils.PathSearch("alarm_time_period", alarm, nil)),
		d.Set("frequency_count", utils.PathSearch("frequency_count", alarm, nil)),
		d.Set("frequency_time", utils.PathSearch("frequency_time", alarm, nil)),
		d.Set("severity", utils.PathSearch("severity", alarm, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", alarm, nil)),
		d.Set("language", utils.PathSearch("language", alarm, nil)),
		d.Set("username", utils.PathSearch("username", alarm, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAlarmConfig(client *golangsdk.ServiceClient, fwInstanceID string, alarmType int) (interface{}, error) {
	getAlarmConfigHttpUrl := "v1/{project_id}/cfw/alarm/config?fw_instance_id={fw_instance_id}"
	getAlarmConfigPath := client.Endpoint + getAlarmConfigHttpUrl
	getAlarmConfigPath = strings.ReplaceAll(getAlarmConfigPath, "{project_id}", client.ProjectID)
	getAlarmConfigPath = strings.ReplaceAll(getAlarmConfigPath, "{fw_instance_id}", fwInstanceID)

	getAlarmConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getAlarmConfigPath, &getAlarmConfigOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	findAlarmStr := fmt.Sprintf("alarm_configs[?alarm_type==`%d`]|[0]", alarmType)
	alarm := utils.PathSearch(findAlarmStr, respBody, nil)
	status := int(utils.PathSearch("enable_status", alarm, float64(2)).(float64))
	// If the status of the alarm is off or not found.
	if status == 0 || status == 2 {
		return nil, golangsdk.ErrDefault404{}
	}

	return alarm, nil
}

func resourceAlarmConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAlarmConfigHttpUrl = "v1/{project_id}/cfw/alarm/config?fw_instance_id={fw_instance_id}"
		deleteAlarmConfigProduct = "cfw"
	)
	deleteAlarmConfigClient, err := cfg.NewServiceClient(deleteAlarmConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	fwInstanceID := d.Get("fw_instance_id").(string)
	deleteAlarmConfigPath := deleteAlarmConfigClient.Endpoint + deleteAlarmConfigHttpUrl
	deleteAlarmConfigPath = strings.ReplaceAll(deleteAlarmConfigPath, "{project_id}", deleteAlarmConfigClient.ProjectID)
	deleteAlarmConfigPath = strings.ReplaceAll(deleteAlarmConfigPath, "{fw_instance_id}", fwInstanceID)

	alarmType := d.Get("alarm_type").(int)
	deleteAlarmConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteAlarmConfigOpt.JSONBody = buildDeleteAlarmConfigBodyParams(alarmType)
	_, err = deleteAlarmConfigClient.Request("PUT", deleteAlarmConfigPath, &deleteAlarmConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW alarm configuration",
		)
	}
	return nil
}

// Restore common parameters to default values.
func buildDeleteAlarmConfigBodyParams(alarmType int) map[string]interface{} {
	return map[string]interface{}{
		"enable_status":     0,
		"alarm_type":        alarmType,
		"alarm_time_period": 0,
		"topic_urn":         "",
	}
}

func resourceAlarmConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <fw_instance_id>/<alarm_type>")
	}

	d.Set("fw_instance_id", parts[0])

	alarmType, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid alarm_type: %s", parts[1])
	}
	d.Set("alarm_type", alarmType)

	return []*schema.ResourceData{d}, nil
}
