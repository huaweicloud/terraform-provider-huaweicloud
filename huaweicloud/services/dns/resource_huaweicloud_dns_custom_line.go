package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS POST /v2.1/customlines
// @API DNS GET /v2.1/customlines
// @API DNS PUT /v2.1/customlines/{line_id}
// @API DNS DELETE /v2.1/customlines/{line_id}
func ResourceCustomLine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomLineCreate,
		UpdateContext: resourceCustomLineUpdate,
		ReadContext:   resourceCustomLineRead,
		DeleteContext: resourceCustomLineDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The region where the custom line is located`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The custom line name.`,
			},
			"ip_segments": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				MaxItems:    50,
				Description: `The IP address range.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom line description.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource status.`,
			},
		},
	}
}

func resourceCustomLineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var region string
	cfg := meta.(*config.Config)
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
		cfg.RegionClient = true
	}

	createCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	if err := createCustomLine(createCustomLineClient, d); err != nil {
		return diag.FromErr(err)
	}

	err = waitForCustomLineCreateOrUpdate(ctx, createCustomLineClient, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCustomLineRead(ctx, d, meta)
}

func createCustomLine(customLineClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	createCustomLineHttpUrl := "v2.1/customlines"
	createCustomLinePath := customLineClient.Endpoint + createCustomLineHttpUrl
	createCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateCustomLineBodyParams(d)),
	}
	createCustomLineResp, err := customLineClient.Request("POST", createCustomLinePath, &createCustomLineOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS custom line: %s", err)
	}

	createCustomLineRespBody, err := utils.FlattenResponse(createCustomLineResp)
	if err != nil {
		return err
	}

	customLineId := utils.PathSearch("line_id", createCustomLineRespBody, "").(string)
	if customLineId == "" {
		return fmt.Errorf("unable to find the DNS custom line ID from the API response")
	}
	d.SetId(customLineId)
	return nil
}

func waitForCustomLineCreateOrUpdate(ctx context.Context, customLineClient *golangsdk.ServiceClient,
	customLineId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      customLineStatusRefreshFunc(customLineClient, customLineId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS custom line (%s) to be ACTIVE : %s", customLineId, err)
	}
	return nil
}

func buildCreateOrUpdateCustomLineBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"ip_segments": utils.ExpandToStringList(d.Get("ip_segments").(*schema.Set).List()),
		// When description is updated to empty, the value of this field must be specified as an empty string.
		"description": d.Get("description"),
	}
	return bodyParams
}

// GetCustomLineById is a method used to query custom line detail by its ID.
func GetCustomLineById(client *golangsdk.ServiceClient, customLineId string) (interface{}, error) {
	getCustomLineHttpUrl := "v2.1/customlines"
	getPath := client.Endpoint + getCustomLineHttpUrl
	getPath += fmt.Sprintf("?line_id=%s", customLineId)

	getCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getCustomLineOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	customLine := utils.PathSearch("lines[0]", respBody, nil)
	if customLine == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return customLine, nil
}

func resourceCustomLineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var region string
	cfg := meta.(*config.Config)
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
		cfg.RegionClient = true
	}

	getCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	customLine, err := GetCustomLineById(getCustomLineClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS custom line")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", customLine, nil)),
		d.Set("ip_segments", utils.PathSearch("ip_segments", customLine, nil)),
		d.Set("status", utils.PathSearch("status", customLine, nil)),
		d.Set("description", utils.PathSearch("description", customLine, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCustomLineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var region string
	cfg := meta.(*config.Config)
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
		cfg.RegionClient = true
	}
	updateCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	if err := updateCustomLine(updateCustomLineClient, d); err != nil {
		return diag.FromErr(err)
	}

	err = waitForCustomLineCreateOrUpdate(ctx, updateCustomLineClient, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceCustomLineRead(ctx, d, meta)
}

func updateCustomLine(customLineClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateCustomLineHttpUrl := "v2.1/customlines/{line_id}"
	updateCustomLinePath := customLineClient.Endpoint + updateCustomLineHttpUrl
	updateCustomLinePath = strings.ReplaceAll(updateCustomLinePath, "{line_id}", d.Id())

	updateCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateCustomLineBodyParams(d)),
	}
	_, err := customLineClient.Request("PUT", updateCustomLinePath, &updateCustomLineOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS custom line: %s", err)
	}
	return nil
}

func resourceCustomLineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		deleteCustomLineHttpUrl = "v2.1/customlines/{line_id}"
		customLineId            = d.Id()
		region                  = ""
	)
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
		cfg.RegionClient = true
	}

	deleteCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	deleteCustomLinePath := deleteCustomLineClient.Endpoint + deleteCustomLineHttpUrl
	deleteCustomLinePath = strings.ReplaceAll(deleteCustomLinePath, "{line_id}", customLineId)

	deleteCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteCustomLineClient.Request("DELETE", deleteCustomLinePath, &deleteCustomLineOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DNS custom line")
	}

	if err := waitForCustomLineDeleted(ctx, deleteCustomLineClient, customLineId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForCustomLineDeleted(ctx context.Context, customLineClient *golangsdk.ServiceClient, customLineId string,
	timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      customLineStatusRefreshFunc(customLineClient, customLineId),
		Timeout:      timeOut,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS custom line (%s) to be DELETED: %s", customLineId, err)
	}
	return nil
}

func customLineStatusRefreshFunc(client *golangsdk.ServiceClient, customLineId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		customLine, err := GetCustomLineById(client, customLineId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return customLine, "", err
		}

		return customLine, parseStatus(utils.PathSearch("status", customLine, "").(string)), nil
	}
}
