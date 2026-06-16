package das

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var searchPathSwitchNonUpdatableParams = []string{
	"connection_id",
}

// @API DAS POST /v3/{project_id}/connections/{connection_id}/clouddba-edit-search-path-flag
// @API DAS GET /v3/{project_id}/connections/{connection_id}/clouddba-get-search-path-flag
func ResourceSearchPathSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSearchPathSwitchCreate,
		ReadContext:   resourceSearchPathSwitchRead,
		UpdateContext: resourceSearchPathSwitchUpdate,
		DeleteContext: resourceSearchPathSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(searchPathSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the search path switch is located.",
			},

			// Required parameters.
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database connection (DB user ID).",
			},
			"switch_on": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the search path switch.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildSearchPathSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"search_path_flag": d.Get("switch_on").(bool),
	}

	return bodyParams
}

func buildSearchPathSwitchRequest(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, *golangsdk.RequestOpts) {
	var (
		httpUrl    = "v3/{project_id}/connections/{connection_id}/clouddba-edit-search-path-flag"
		createPath = client.Endpoint + httpUrl
	)
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", d.Get("connection_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildSearchPathSwitchBodyParams(d),
	}

	return createPath, &createOpt
}

func buildSearchPathSwitchGetRequest(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, *golangsdk.RequestOpts) {
	var (
		httpUrl = "v3/{project_id}/connections/{connection_id}/clouddba-get-search-path-flag"
		getPath = client.Endpoint + httpUrl
	)
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", d.Get("connection_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	return getPath, &getOpt
}

func searchPathSwitchRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath, getOpt := buildSearchPathSwitchGetRequest(client, d)
		requestResp, err := client.Request("GET", getPath, getOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, "ERROR", err
		}

		searchPathFlag := utils.PathSearch("search_path_flag", respBody, nil)
		if searchPathFlag == nil {
			return nil, "ERROR", errors.New("'search_path_flag' field is not found in the response of API")
		}

		currentFlag := searchPathFlag.(bool)
		targetFlag := d.Get("switch_on").(bool)

		if currentFlag == targetFlag {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitForSearchPathSwitchComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      searchPathSwitchRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the DAS search path switch to become %v: %s", d.Get("switch_on").(bool), err)
	}
	return nil
}

func resourceSearchPathSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	createPath, createOpt := buildSearchPathSwitchRequest(client, d)
	_, err = client.Request("POST", createPath, createOpt)
	if err != nil {
		return diag.Errorf("error switching DAS search path switch: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	if err = waitForSearchPathSwitchComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for the DAS search path switch to complete: %s", err)
	}

	return resourceSearchPathSwitchRead(ctx, d, meta)
}

func resourceSearchPathSwitchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	getPath, getOpt := buildSearchPathSwitchGetRequest(client, d)
	requestResp, err := client.Request("GET", getPath, getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errorCodeStr", "DAS.5010"),
			"error querying DAS search path switch")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("switch_on", utils.PathSearch("search_path_flag", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSearchPathSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	updatePath, updateOpt := buildSearchPathSwitchRequest(client, d)
	_, err = client.Request("POST", updatePath, updateOpt)
	if err != nil {
		return diag.Errorf("error switching DAS search path switch: %s", err)
	}

	if err = waitForSearchPathSwitchComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for the DAS search path switch to complete: %s", err)
	}

	return resourceSearchPathSwitchRead(ctx, d, meta)
}

func resourceSearchPathSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for switching the search path switch. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
