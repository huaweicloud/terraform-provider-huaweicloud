package dsc

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

var nonUpdatableParamsMeasureInfo = []string{
	"life_cycle",
	"measure_type",
}

// @API DSC POST /v1/{project_id}/protect-measures/batch-create
// @API DSC GET /v1/{project_id}/protect-measures
// @API DSC POST /v1/{project_id}/protect-measures/batch-update
// @API DSC POST /v1/{project_id}/protect-measures/batch-delete
func ResourceMeasureInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMeasureInfoCreate,
		ReadContext:   resourceMeasureInfoRead,
		UpdateContext: resourceMeasureInfoUpdate,
		DeleteContext: resourceMeasureInfoDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMeasureInfoImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsMeasureInfo),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the measure name.",
			},
			"life_cycle": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the data lifecycle stage.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the measure description.",
			},
			"measure_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the measure type.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"is_deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the measure is deleted.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time.",
			},
		},
	}
}

func buildCreateMeasureInfoBodyParams(d *schema.ResourceData) []interface{} {
	measureItem := map[string]interface{}{
		"name": d.Get("name"),
	}

	if v, ok := d.GetOk("description"); ok {
		measureItem["description"] = v
	}
	if v, ok := d.GetOk("life_cycle"); ok {
		measureItem["life_cycle"] = v
	}
	if v, ok := d.GetOk("measure_type"); ok {
		measureItem["measure_type"] = v
	}

	return []interface{}{measureItem}
}

func resourceMeasureInfoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		name      = d.Get("name").(string)
		lifeCycle = d.Get("life_cycle").(string)
		httpUrl   = "v1/{project_id}/protect-measures/batch-create"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildCreateMeasureInfoBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DSC measure info: %s", err)
	}

	// The create API does not return the specific resource ID, search by name to get the ID
	respBody, err := GetMeasureInfoByName(client, name, lifeCycle)
	if err != nil {
		return diag.Errorf("error retrieving DSC measure info after creation: %s", err)
	}

	measureId := fmt.Sprintf("%v", utils.PathSearch("id", respBody, ""))
	if measureId == "" || measureId == "<nil>" {
		return diag.Errorf("error creating DSC measure info: ID is not found in API response")
	}

	d.SetId(measureId)

	return resourceMeasureInfoRead(ctx, d, meta)
}

func GetMeasureInfo(client *golangsdk.ServiceClient, filterField, filterValue, lifeCycle string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/protect-measures"
		offset  = 0
		limit   = 1000
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + fmt.Sprintf("?life_cycle=%s&limit=%d&offset=%d", lifeCycle, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		measureList := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
		for _, measure := range measureList {
			fieldValue := fmt.Sprintf("%v", utils.PathSearch(filterField, measure, nil))
			if fieldValue == filterValue {
				return measure, nil
			}
		}

		if len(measureList) < limit {
			break
		}

		offset += len(measureList)
	}

	return nil, golangsdk.ErrDefault404{}
}

func GetMeasureInfoByName(client *golangsdk.ServiceClient, name, lifeCycle string) (interface{}, error) {
	return GetMeasureInfo(client, "name", name, lifeCycle)
}

func GetMeasureInfoById(client *golangsdk.ServiceClient, id, lifeCycle string) (interface{}, error) {
	return GetMeasureInfo(client, "id", id, lifeCycle)
}

func resourceMeasureInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		lifeCycle = d.Get("life_cycle").(string)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := GetMeasureInfoById(client, d.Id(), lifeCycle)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC measure info")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("life_cycle", utils.PathSearch("life_cycle", respBody, nil)),
		d.Set("measure_type", utils.PathSearch("measure_type", respBody, nil)),
		d.Set("is_deleted", utils.PathSearch("is_deleted", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateMeasureInfoBodyParams(d *schema.ResourceData) []interface{} {
	measureItem := map[string]interface{}{
		"id":   d.Id(),
		"name": d.Get("name"),
	}

	if v, ok := d.GetOk("description"); ok {
		measureItem["description"] = v
	}
	if v, ok := d.GetOk("life_cycle"); ok {
		measureItem["life_cycle"] = v
	}
	if v, ok := d.GetOk("measure_type"); ok {
		measureItem["measure_type"] = v
	}
	if v, ok := d.GetOk("is_deleted"); ok {
		measureItem["is_deleted"] = v
	}
	if v, ok := d.GetOk("create_time"); ok {
		measureItem["create_time"] = v
	}
	if v, ok := d.GetOk("update_time"); ok {
		measureItem["update_time"] = v
	}

	return []interface{}{measureItem}
}

func resourceMeasureInfoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/protect-measures/batch-update"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

		updateOpt := golangsdk.RequestOpts{
			MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
			KeepResponseBody: true,
			JSONBody:         buildUpdateMeasureInfoBodyParams(d),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating DSC measure info: %s", err)
		}
	}

	return resourceMeasureInfoRead(ctx, d, meta)
}

func resourceMeasureInfoDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/protect-measures/batch-delete"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	// The delete API accepts an array of IDs (Long type)
	measureIdInt, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("error converting measure ID to integer: %s", err)
	}
	deleteBody := []interface{}{measureIdInt}

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         deleteBody,
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "error_code", "dsc.40000070"), "error deleting DSC measure info")
	}

	return nil
}

func resourceMeasureInfoImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<id>/<life_cycle>', but got '%s'",
			importedId)
	}

	d.SetId(parts[0])

	mErr := multierror.Append(nil,
		d.Set("life_cycle", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
