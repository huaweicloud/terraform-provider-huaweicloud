package fgs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/events
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/events

func ResourceFunctionEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionEventCreate,
		ReadContext:   resourceFunctionEventRead,
		UpdateContext: resourceFunctionEventUpdate,
		DeleteContext: resourceFunctionEventDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFunctionEventImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function event is located.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The URN of the function to which the event blongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The function event name.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The function event content.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update (UTC) time of the function event, in RFC3339 format.`,
			},
		},
	}
}

func resourceFunctionEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/events"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", d.Get("function_urn").(string))
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildFunctionEventRequestBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph function event: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("id", respBody, "")
	d.SetId(resourceId.(string))

	return resourceFunctionEventRead(ctx, d, meta)
}

func buildFunctionEventRequestBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":    d.Get("name"),
		"content": d.Get("content"),
	}
}

func parseFunctionEventLastModified(timeNum interface{}) interface{} {
	var result string
	if timeNum == nil {
		return result
	}

	switch num := timeNum.(type) {
	case int:
		log.Printf("[DEBUG] The type of attribute 'last_modified' response value is 'int'")
		result = utils.FormatTimeStampRFC3339(int64(num), true)
	case float64:
		log.Printf("[DEBUG] The type of attribute 'last_modified' response value is 'float64'")
		// Ignore loss of precision.
		result = utils.FormatTimeStampRFC3339(int64(num), true)
	}

	return result
}

func resourceFunctionEventRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		resourceId = d.Id()
		httpUrl    = "v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", d.Get("function_urn").(string))
	getPath = strings.ReplaceAll(getPath, "{event_id}", resourceId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph function event")
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error extrieving function event (%s): %s", resourceId, err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("content", utils.PathSearch("content", respBody, nil)),
		d.Set("updated_at", parseFunctionEventLastModified(utils.PathSearch("last_modified", respBody, nil))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving function event fields: %s", err)
	}
	return nil
}

func resourceFunctionEventUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}"
		resourceId = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", d.Get("function_urn").(string))
	updatePath = strings.ReplaceAll(updatePath, "{event_id}", resourceId)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildFunctionEventRequestBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating FunctionGraph function event (%s): %s", resourceId, err)
	}
	return resourceFunctionEventRead(ctx, d, meta)
}

func resourceFunctionEventDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", d.Get("function_urn").(string))
	deletePath = strings.ReplaceAll(deletePath, "{event_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting FunctionGraph function event")
	}
	return nil
}

func resourceFunctionEventImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<function_urn>/<name>', but got '%s'", d.Id())
	}

	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/events"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{function_urn}", parts[0])

	queryOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	queryResp, err := client.Request("GET", queryPath, &queryOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving function event: %s", err)
	}

	respBody, err := utils.FlattenResponse(queryResp)
	if err != nil {
		return nil, err
	}

	eventId := utils.PathSearch(fmt.Sprintf("events[?name=='%s']|[0].id", parts[1]), respBody, "").(string)
	if eventId == "" {
		return nil, fmt.Errorf("unable to find the resource ID of the function event")
	}

	d.SetId(eventId)
	return []*schema.ResourceData{d}, d.Set("function_urn", parts[0])
}
