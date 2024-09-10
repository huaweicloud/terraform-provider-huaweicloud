// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DNS
// ---------------------------------------------------------------

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
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS GET /v2.1/customlines
// @API DNS POST /v2.1/customlines
// @API DNS DELETE /v2.1/customlines/{line_id}
// @API DNS PUT /v2.1/customlines/{line_id}
func ResourceDNSCustomLine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSCustomLineCreate,
		UpdateContext: resourceDNSCustomLineUpdate,
		ReadContext:   resourceDNSCustomLineRead,
		DeleteContext: resourceDNSCustomLineDelete,
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
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the custom line name.`,
			},
			"ip_segments": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				MaxItems:    50,
				Description: `Specifies the IP address range.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom line description. A maximum of 255 characters are allowed.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource status.`,
			},
		},
	}
}

func resourceDNSCustomLineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	createDNSCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	// createDNSCustomLine: create DNS custom line.
	if err := createDNSCustomLine(createDNSCustomLineClient, d); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := waitForDNSCustomLineCreateOrUpdate(ctx, createDNSCustomLineClient, d, timeout); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSCustomLineRead(ctx, d, meta)
}

func createDNSCustomLine(customLineClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		createDNSCustomLineHttpUrl = "v2.1/customlines"
	)

	createDNSCustomLinePath := customLineClient.Endpoint + createDNSCustomLineHttpUrl
	createDNSCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createDNSCustomLineOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateDNSCustomLineBodyParams(d))
	createDNSCustomLineResp, err := customLineClient.Request("POST", createDNSCustomLinePath,
		&createDNSCustomLineOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS custom line: %s", err)
	}

	createDNSCustomLineRespBody, err := utils.FlattenResponse(createDNSCustomLineResp)
	if err != nil {
		return err
	}

	id, err := jmespath.Search("line_id", createDNSCustomLineRespBody)
	if err != nil {
		return fmt.Errorf("error creating DNS custom line: ID is not found in API response")
	}
	d.SetId(id.(string))
	return nil
}

func waitForDNSCustomLineCreateOrUpdate(ctx context.Context, customLineClient *golangsdk.ServiceClient,
	d *schema.ResourceData, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      dnsCustomLineStatusRefreshFunc(d, customLineClient),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS custom line (%s) to be ACTIVE : %s", d.Id(), err)
	}
	return nil
}

func buildCreateOrUpdateDNSCustomLineBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_segments": utils.ValueIgnoreEmpty(d.Get("ip_segments")),
	}
	return bodyParams
}

func resourceDNSCustomLineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDNSCustomLine: Query DNS custom line
	var (
		getDNSCustomLineHttpUrl = "v2.1/customlines"
		getDNSCustomLineProduct = "dns"
	)
	getDNSCustomLineClient, err := cfg.NewServiceClient(getDNSCustomLineProduct, region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	getDNSCustomLinePath := getDNSCustomLineClient.Endpoint + getDNSCustomLineHttpUrl
	getDNSCustomLinePath += buildGetDNSCustomLineQueryParams(d)

	getDNSCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSCustomLineResp, err := getDNSCustomLineClient.Request("GET", getDNSCustomLinePath,
		&getDNSCustomLineOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS custom line")
	}

	getDNSCustomLineRespBody, err := utils.FlattenResponse(getDNSCustomLineResp)
	if err != nil {
		return diag.FromErr(err)
	}

	customLineMap, err := flattenCustomLineResponseBody(getDNSCustomLineRespBody, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS custom line")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", customLineMap["name"]),
		d.Set("ip_segments", customLineMap["ip_segments"]),
		d.Set("status", customLineMap["status"]),
		d.Set("description", customLineMap["description"]),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCustomLineResponseBody(resp interface{}, d *schema.ResourceData) (map[string]interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("the custom line response is empty")
	}

	curJson := utils.PathSearch("lines", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		lineId := utils.PathSearch("line_id", v, "")
		if d.Id() == lineId.(string) {
			return map[string]interface{}{
				"name":        utils.PathSearch("name", v, nil),
				"ip_segments": utils.PathSearch("ip_segments", v, nil),
				"status":      utils.PathSearch("status", v, nil),
				"description": utils.PathSearch("description", v, nil),
			}, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func buildGetDNSCustomLineQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?line_id=%s", d.Id())
}

func resourceDNSCustomLineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateDNSCustomLineClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	updateDNSCustomLineChanges := []string{
		"name",
		"description",
		"ip_segments",
	}

	if d.HasChanges(updateDNSCustomLineChanges...) {
		// updateDNSCustomLine: Update DNS custom line
		if err := updateDNSCustomLine(updateDNSCustomLineClient, d); err != nil {
			return diag.FromErr(err)
		}

		timeout := d.Timeout(schema.TimeoutUpdate)
		if err := waitForDNSCustomLineCreateOrUpdate(ctx, updateDNSCustomLineClient, d, timeout); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceDNSCustomLineRead(ctx, d, meta)
}

func updateDNSCustomLine(customLineClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		updateDNSCustomLineHttpUrl = "v2.1/customlines/{line_id}"
	)

	updateDNSCustomLinePath := customLineClient.Endpoint + updateDNSCustomLineHttpUrl
	updateDNSCustomLinePath = strings.ReplaceAll(updateDNSCustomLinePath, "{line_id}", d.Id())

	updateDNSCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSCustomLineOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateDNSCustomLineBodyParams(d))
	_, err := customLineClient.Request("PUT", updateDNSCustomLinePath,
		&updateDNSCustomLineOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS custom line: %s", err)
	}
	return nil
}

func resourceDNSCustomLineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDNSCustomLine: Delete DNS custom line
	var (
		deleteDNSCustomLineHttpUrl = "v2.1/customlines/{line_id}"
		deleteDNSCustomLineProduct = "dns"
	)
	deleteDNSCustomLineClient, err := cfg.NewServiceClient(deleteDNSCustomLineProduct, region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	deleteDNSCustomLinePath := deleteDNSCustomLineClient.Endpoint + deleteDNSCustomLineHttpUrl
	deleteDNSCustomLinePath = strings.ReplaceAll(deleteDNSCustomLinePath, "{line_id}", d.Id())

	deleteDNSCustomLineOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = deleteDNSCustomLineClient.Request("DELETE", deleteDNSCustomLinePath, &deleteDNSCustomLineOpt)
	if err != nil {
		return diag.Errorf("error deleting DNS custom line: %s", err)
	}

	if err := waitForDNSCustomLineDeleted(ctx, deleteDNSCustomLineClient, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForDNSCustomLineDeleted(ctx context.Context, customLineClient *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      dnsCustomLineStatusRefreshFunc(d, customLineClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS custom line (%s) to be DELETED: %s", d.Id(), err)
	}
	return nil
}

func dnsCustomLineStatusRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getDNSCustomLineHttpUrl = "v2.1/customlines"
		)

		getDNSCustomLinePath := client.Endpoint + getDNSCustomLineHttpUrl
		getDNSCustomLinePath += buildGetDNSCustomLineQueryParams(d)
		getDNSCustomLineOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		getDNSCustomLineResp, err := client.Request("GET", getDNSCustomLinePath, &getDNSCustomLineOpt)
		if err != nil {
			return nil, "", err
		}

		getDNSCustomLineRespBody, err := utils.FlattenResponse(getDNSCustomLineResp)
		if err != nil {
			return nil, "", err
		}
		customLineMap, err := flattenCustomLineResponseBody(getDNSCustomLineRespBody, d)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return customLineMap, "", err
		}

		status := customLineMap["status"]
		return getDNSCustomLineRespBody, parseStatus(status.(string)), nil
	}
}
