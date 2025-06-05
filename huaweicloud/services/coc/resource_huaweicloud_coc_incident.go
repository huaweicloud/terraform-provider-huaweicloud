package coc

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

var incidentNonUpdatableParams = []string{"regions", "enterprise_project", "start_time", "current_cloud_service",
	"incident_level", "is_service_interrupt", "incident_type", "incident_ownership", "incident_title",
	"incident_description", "incident_source", "incident_assignee", "assignee_scene", "assignee_role", "creator"}

// @API COC POST /v1/external/incident/create
// @API COC GET /v1/external/incident/{incident_num}
func ResourceIncident() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentCreate,
		ReadContext:   resourceIncidentRead,
		UpdateContext: resourceIncidentUpdate,
		DeleteContext: resourceIncidentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(incidentNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"incident_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_service_interrupt": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"incident_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"current_cloud_service": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incident_ownership": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incident_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"incident_assignee": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"assignee_scene": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignee_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"warroom_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"handle_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enum_data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filed_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enum_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildIncidentCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"incident_level":        d.Get("incident_level"),
		"is_service_interrupt":  d.Get("is_service_interrupt"),
		"incident_type":         d.Get("incident_type"),
		"incident_title":        d.Get("incident_title"),
		"incident_source":       d.Get("incident_source"),
		"creator":               d.Get("creator"),
		"region":                utils.ExpandToStringList(d.Get("regions").([]interface{})),
		"enterprise_project":    utils.ExpandToStringList(d.Get("enterprise_project").([]interface{})),
		"start_time":            utils.ValueIgnoreEmpty(d.Get("start_time")),
		"current_cloud_service": utils.ExpandToStringList(d.Get("current_cloud_service").([]interface{})),
		"incident_ownership":    utils.ValueIgnoreEmpty(d.Get("incident_ownership")),
		"incident_description":  utils.ValueIgnoreEmpty(d.Get("incident_description")),
		"incident_assignee":     utils.ExpandToStringList(d.Get("incident_assignee").([]interface{})),
		"assignee_scene":        utils.ValueIgnoreEmpty(d.Get("assignee_scene")),
		"assignee_role":         utils.ValueIgnoreEmpty(d.Get("assignee_role")),
	}

	return bodyParams
}

func resourceIncidentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/external/incident/create"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildIncidentCreateOpts(d)),
	}

	createIncidentResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC incident: %s", err)
	}

	createIncidentRespBody, err := utils.FlattenResponse(createIncidentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	incidentNum := utils.PathSearch("data.incident_num", createIncidentRespBody, "").(string)
	if incidentNum == "" {
		return diag.Errorf("unable to find the COC incident num from the API response")
	}
	d.SetId(incidentNum)

	return resourceIncidentRead(ctx, d, meta)
}

func resourceIncidentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	readIncidentRespBody, err := GetIncident(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00067003"),
			fmt.Sprintf("error retrieving incident (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("regions", utils.PathSearch("data.region", readIncidentRespBody, nil)),
		d.Set("enterprise_project", utils.PathSearch("data.enterprise_project", readIncidentRespBody, nil)),
		d.Set("current_cloud_service", utils.PathSearch("data.current_cloud_service", readIncidentRespBody, nil)),
		d.Set("incident_level", utils.PathSearch("data.incident_level", readIncidentRespBody, nil)),
		d.Set("is_service_interrupt", utils.PathSearch("data.is_service_interrupt", readIncidentRespBody, nil)),
		d.Set("start_time", utils.PathSearch("data.start_time", readIncidentRespBody, nil)),
		d.Set("incident_type", utils.PathSearch("data.incident_type", readIncidentRespBody, nil)),
		d.Set("incident_title", utils.PathSearch("data.incident_title", readIncidentRespBody, nil)),
		d.Set("incident_description", utils.PathSearch("data.incident_description", readIncidentRespBody, nil)),
		d.Set("incident_source", utils.PathSearch("data.incident_source", readIncidentRespBody, nil)),
		d.Set("incident_ownership", utils.PathSearch("data.incident_ownership", readIncidentRespBody, nil)),
		d.Set("incident_assignee", utils.PathSearch("data.incident_assignee", readIncidentRespBody, nil)),
		d.Set("assignee_scene", utils.PathSearch("data.assignee_scene", readIncidentRespBody, nil)),
		d.Set("assignee_role", utils.PathSearch("data.assignee_role", readIncidentRespBody, nil)),
		d.Set("warroom_id", utils.PathSearch("data.warroom_id", readIncidentRespBody, nil)),
		d.Set("handle_time", utils.PathSearch("data.handle_time", readIncidentRespBody, nil)),
		d.Set("status", utils.PathSearch("data.status", readIncidentRespBody, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", readIncidentRespBody, nil)),
		d.Set("creator", utils.PathSearch("data.creator", readIncidentRespBody, nil)),
		d.Set("enum_data_list", flattenIncidentEnumDataListParams(readIncidentRespBody)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC incident fields: %s", err)
	}

	return nil
}

func GetIncident(client *golangsdk.ServiceClient, incidentNum string) (interface{}, error) {
	httpUrl := "v1/external/incident/{incident_num}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{incident_num}", incidentNum)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	readIncidentResp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}
	readIncidentRespBody, err := utils.FlattenResponse(readIncidentResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening incident: %s", err)
	}
	return readIncidentRespBody, nil
}

func flattenIncidentEnumDataListParams(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("data.enum_data_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"filed_key": utils.PathSearch("filed_key", v, nil),
			"enum_key":  utils.PathSearch("enum_key", v, nil),
			"name_zh":   utils.PathSearch("name_zh", v, nil),
			"name_en":   utils.PathSearch("name_en", v, nil),
		}
	}
	return rst
}

func resourceIncidentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIncidentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting incident resource is not supported. The incident resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
