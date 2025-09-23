package codeartsbuild

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var taskActionNonUpdatableParams = []string{
	"job_id", "action", "build_no", "parameter", "scm", "scm.0.build_tag", "scm.0.build_commit_id",
}

// @API CodeArtsBuild POST /v1/job/execute
// @API CodeArtsBuild POST /v1/job/{job_id}/stop
func ResourceCodeArtsBuildTaskAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBuildTaskActionCreate,
		ReadContext:   resourceBuildTaskActionRead,
		UpdateContext: resourceBuildTaskActionUpdate,
		DeleteContext: resourceBuildTaskActionDelete,

		CustomizeDiff: config.FlexibleForceNew(taskActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the build task ID.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action.`,
				ValidateFunc: validation.StringInSlice([]string{
					"execute", "stop",
				}, false),
			},
			"build_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the build task number, start from 1.`,
			},
			"parameter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the parameter list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the parameter name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the parameter value.`,
						},
					},
				},
			},
			"scm": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the build execution SCM.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"build_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the build tag.`,
						},
						"build_commit_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the build commit ID.`,
						},
					},
				},
			},
			"daily_build_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the daily build number.`,
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

var actionFunc = map[string]func(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics{
	"execute": executeTask,
	"stop":    stopTask,
}

func resourceBuildTaskActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_build", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	action := d.Get("action").(string)
	return actionFunc[action](client, d)
}

func executeTask(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	httpUrl := "v1/job/execute"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildExecuteBuildTaskBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error executing CodeArts Build task: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody); err != nil {
		return diag.Errorf("error executing CodeArts Build task: %s", err)
	}

	buildNo := utils.PathSearch("result.actual_build_number", createRespBody, "").(string)
	if buildNo == "" {
		return diag.Errorf("unable to find the CodeArts Build task actual_build_number from the API response")
	}

	id := fmt.Sprintf("%s/%s", d.Get("job_id"), buildNo)
	d.SetId(id)
	d.Set("build_no", buildNo)
	d.Set("daily_build_number", utils.PathSearch("result.daily_build_number", createRespBody, nil))

	return nil
}

func buildExecuteBuildTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_id":    d.Get("job_id"),
		"parameter": buildExecuteBuildTaskParameters(d),
		"scm":       buildExecuteBuildTaskScms(d),
	}

	return bodyParams
}

func buildExecuteBuildTaskParameters(d *schema.ResourceData) interface{} {
	rawParameters := d.Get("parameter").(*schema.Set).List()
	if len(rawParameters) == 0 {
		return nil
	}

	parameters := make([]map[string]interface{}, 0, len(rawParameters))
	for _, p := range rawParameters {
		if parameter, ok := p.(map[string]interface{}); ok {
			parameterMap := map[string]interface{}{
				"name":  parameter["name"],
				"value": parameter["value"],
			}
			parameters = append(parameters, parameterMap)
		}
	}

	return parameters
}

func buildExecuteBuildTaskScms(d *schema.ResourceData) interface{} {
	rawScms := d.Get("scm").([]interface{})
	if len(rawScms) == 0 {
		return nil
	}

	if scm, ok := rawScms[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"build_tag":       utils.ValueIgnoreEmpty(scm["build_tag"]),
			"build_commit_id": utils.ValueIgnoreEmpty(scm["build_commit_id"]),
		}
	}

	return nil
}

func stopTask(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	jobId := d.Get("job_id").(string)
	httpUrl := "v1/job/{job_id}/stop"
	stopPath := client.Endpoint + httpUrl
	stopPath = strings.ReplaceAll(stopPath, "{job_id}", jobId)
	stopOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"build_no": utils.ValueIgnoreEmpty(d.Get("build_no")),
		}),
	}

	stopResp, err := client.Request("POST", stopPath, &stopOpt)
	if err != nil {
		return diag.Errorf("error stopping CodeArts Build task: %s", err)
	}
	stopRespBody, err := utils.FlattenResponse(stopResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(stopRespBody); err != nil {
		return diag.Errorf("error stopping CodeArts Build task: %s", err)
	}

	id := jobId
	if v, ok := d.GetOk("build_no"); ok {
		id += "/" + v.(string)
	}
	d.SetId(id)

	return nil
}

func resourceBuildTaskActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBuildTaskActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBuildTaskActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting build task action resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
