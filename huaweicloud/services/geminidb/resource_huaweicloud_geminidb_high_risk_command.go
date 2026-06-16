package geminidb

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

var highRiskCommandNonUpdatableParams = []string{
	"instance_id",
	"origin_name",
}

// @API GeminiDB PUT /v3/{project_id}/instances/{instance_id}/high-risk-commands
// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/high-risk-commands
func ResourceHighRiskCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHighRiskCommandCreate,
		ReadContext:   resourceHighRiskCommandRead,
		UpdateContext: resourceHighRiskCommandUpdate,
		DeleteContext: resourceHighRiskCommandDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceHighRiskCommandImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(highRiskCommandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"origin_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func buildHighRiskCommandBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"commands": []map[string]interface{}{
			{
				"origin_name": d.Get("origin_name"),
				"name":        d.Get("name"),
			},
		},
	}

	return bodyParams
}

func modifyHighRiskCommand(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/high-risk-commands"
	reqPath := client.Endpoint + httpUrl
	reqPath = strings.ReplaceAll(reqPath, "{project_id}", client.ProjectID)
	reqPath = strings.ReplaceAll(reqPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildHighRiskCommandBodyParams(d),
	}

	_, err := client.Request("PUT", reqPath, &createOpt)

	return err
}

func resourceHighRiskCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		originName = d.Get("origin_name").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	err = modifyHighRiskCommand(client, d, instanceId)
	if err != nil {
		return diag.Errorf("error modifying GeminiDB high risk command: %s", err)
	}

	resourceId := fmt.Sprintf("%s/%s", instanceId, originName)

	d.SetId(resourceId)

	return resourceHighRiskCommandRead(ctx, d, meta)
}

func resourceHighRiskCommandRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	resourceId := strings.Split(d.Id(), "/")
	commandInfo, err := GetHighRiskCommandInfo(client, resourceId[0], resourceId[1])
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB high risk commands")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", resourceId[0]),
		d.Set("origin_name", utils.PathSearch("origin_name", commandInfo, nil)),
		d.Set("name", utils.PathSearch("name", commandInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetHighRiskCommandInfo(client *golangsdk.ServiceClient, instanceId, originName string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/high-risk-commands"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	commandInfo := utils.PathSearch(fmt.Sprintf("commands|[?origin_name=='%s']|[0]", originName), respBody, nil)
	if commandInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return commandInfo, nil
}

func resourceHighRiskCommandUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	resourceId := strings.Split(d.Id(), "/")
	err = modifyHighRiskCommand(client, d, resourceId[0])
	if err != nil {
		return diag.Errorf("error updating GeminiDB high risk command: %s", err)
	}

	return resourceHighRiskCommandRead(ctx, d, meta)
}

func resourceHighRiskCommandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GeminiDB high risk command modify resource is not supported. The resource is only removed from the " +
		"state, the resource remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceHighRiskCommandImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<origin_name>', but got '%s'", importedId)
	}

	d.SetId(importedId)

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("origin_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
