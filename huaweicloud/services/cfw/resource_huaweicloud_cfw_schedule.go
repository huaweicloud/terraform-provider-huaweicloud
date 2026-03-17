package cfw

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// @API CFW POST /v1/{project_id}/schedule
// @API CFW PUT /v1/{project_id}/schedule/{schedule_id}
// @API CFW DELETE /v1/{project_id}/schedule/{schedule_id}
// @API CFW GET /v1/{project_id}/schedules
func ResourceSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScheduleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"object_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Field `description` can be edited to be empty.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"periodic": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     periodicSchema(),
			},
			"absolute": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem:     absoluteSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"ref_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func absoluteSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func periodicSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"week_mask": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"start_week": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"end_week": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildSchedulePeriodicWeekMaskParams(rawArray []interface{}) []int {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]int, 0, len(rawArray))
	for _, v := range rawArray {
		intVal, ok := v.(int)
		if !ok {
			continue
		}

		rst = append(rst, intVal)
	}

	return rst
}

func buildSchedulePeriodicParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"type":       rawMap["type"],
			"start_time": rawMap["start_time"],
			"end_time":   rawMap["end_time"],
			"week_mask":  buildSchedulePeriodicWeekMaskParams(rawMap["week_mask"].([]interface{})),
			"start_week": utils.ValueIgnoreEmpty(rawMap["start_week"]),
			"end_week":   utils.ValueIgnoreEmpty(rawMap["end_week"]),
		})
	}

	return rst
}

func buildScheduleAbsoluteParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rstMap := map[string]interface{}{
		"start_time": utils.ValueIgnoreEmpty(rawMap["start_time"]),
		"end_time":   utils.ValueIgnoreEmpty(rawMap["end_time"]),
	}

	return rstMap
}

func buildCreateScheduleBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"object_id":   d.Get("object_id"),
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"periodic":    buildSchedulePeriodicParams(d.Get("periodic").([]interface{})),
		"absolute":    buildScheduleAbsoluteParams(d.Get("absolute").([]interface{})),
	}
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/schedule"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateScheduleBodyParam(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CFW schedule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CFW schedule: ID is not found in API response")
	}
	d.SetId(id)

	return resourceScheduleRead(ctx, d, meta)
}

func GetScheduleById(client *golangsdk.ServiceClient, objectId, id string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/schedules"
		offset  = 0
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?object_id=%s&limit=1024", objectId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		target := utils.PathSearch(fmt.Sprintf("[?schedule_id=='%s']|[0]", id), records, nil)
		if target != nil {
			return target, nil
		}

		offset += len(records)
	}

	return nil, golangsdk.ErrDefault404{}
}

func flattenPeriodicAttr(respBody interface{}) []map[string]interface{} {
	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"type":       utils.PathSearch("type", v, nil),
			"start_time": utils.PathSearch("start_time", v, nil),
			"end_time":   utils.PathSearch("end_time", v, nil),
			"week_mask":  utils.PathSearch("week_mask", v, nil),
			"start_week": utils.PathSearch("start_week", v, nil),
			"end_week":   utils.PathSearch("end_week", v, nil),
		})
	}

	return rst
}

func convertAbsoluteStringToInt(respValue interface{}) int {
	stringValue, ok := respValue.(string)
	if !ok {
		return 0
	}

	r, err := strconv.Atoi(stringValue)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", stringValue)
	}

	return r
}

func flattenAbsoluteAttr(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"start_time": convertAbsoluteStringToInt(utils.PathSearch("start_time", respBody, nil)),
		"end_time":   convertAbsoluteStringToInt(utils.PathSearch("end_time", respBody, nil)),
	}

	return []map[string]interface{}{rstMap}
}

func resourceScheduleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "cfw"
		objectId = d.Get("object_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	respBody, err := GetScheduleById(client, objectId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW schedule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("periodic", flattenPeriodicAttr(utils.PathSearch("periodic", respBody, nil))),
		d.Set("absolute", flattenAbsoluteAttr(utils.PathSearch("absolute", respBody, nil))),
		d.Set("ref_count", utils.PathSearch("ref_count", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateScheduleBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"object_id":   d.Get("object_id"),
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"periodic":    buildSchedulePeriodicParams(d.Get("periodic").([]interface{})),
		"absolute":    buildScheduleAbsoluteParams(d.Get("absolute").([]interface{})),
	}
}

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/schedule/{schedule_id}"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{schedule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateScheduleBodyParam(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating CFW schedule: %s", err)
	}

	return resourceScheduleRead(ctx, d, meta)
}

func resourceScheduleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/schedule/{schedule_id}"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{schedule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CFW schedule: %s", err)
	}

	return nil
}

func resourceScheduleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <object_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("object_id", parts[0])
}
