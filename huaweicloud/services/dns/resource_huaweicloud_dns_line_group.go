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

// @API DNS POST /v2.1/linegroups
// @API DNS GET /v2.1/linegroups/{linegroup_id}
// @API DNS PUT /v2.1/linegroups/{linegroup_id}
// @API DNS DELETE /v2.1/linegroups/{linegroup_id}
func ResourceLineGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLineGroupCreate,
		UpdateContext: resourceLineGroupUpdate,
		ReadContext:   resourceLineGroupRead,
		DeleteContext: resourceLineGroupDelete,
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
				Description: `The line group name.`,
			},
			"lines": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of the resolution line IDs.`,
				Set:         schema.HashString,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The line group description.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the line group.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the line group.`,
			},
		},
	}
}

func resourceLineGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dns", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if err := createLineGroup(client, d); err != nil {
		return diag.FromErr(err)
	}

	err = waitForLineGroupCreatedOrUpdated(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLineGroupRead(ctx, d, meta)
}

func createLineGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	createLineGroupHttpUrl := "v2.1/linegroups"
	createLineGroupPath := client.Endpoint + createLineGroupHttpUrl
	createLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateLineGroupBodyParams(d)),
	}
	createLineGroupResp, err := client.Request("POST", createLineGroupPath, &createLineGroupOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS line group: %s", err)
	}

	createLineGroupRespBody, err := utils.FlattenResponse(createLineGroupResp)
	if err != nil {
		return err
	}

	lineId := utils.PathSearch("line_id", createLineGroupRespBody, "").(string)
	if lineId == "" {
		return fmt.Errorf("unable to find the related line ID of the DNS line group from the API response")
	}
	d.SetId(lineId)
	return nil
}

func waitForLineGroupCreatedOrUpdated(ctx context.Context, client *golangsdk.ServiceClient, lineGroupId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"COMPLETED"},
		Pending:      []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Refresh:      lineGroupStatusRefreshFunc(client, lineGroupId, []string{"ACTIVE"}),
		Timeout:      timeOut,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS line group (%s) to be completed: %s", lineGroupId, err)
	}
	return nil
}

func buildCreateOrUpdateLineGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnsLines := d.Get("lines").(*schema.Set).List()

	return map[string]interface{}{
		// Required parameters.
		"name":  d.Get("name"),
		"lines": dnsLines,
		// Optional parameters.
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

// GetLineGroupById is a method used to get line group information by specified ID.
func GetLineGroupById(client *golangsdk.ServiceClient, lineGroupId string) (interface{}, error) {
	getLineGroupHttpUrl := "v2.1/linegroups/{linegroup_id}"
	getLineGroupPath := client.Endpoint + getLineGroupHttpUrl
	getLineGroupPath = strings.ReplaceAll(getLineGroupPath, "{linegroup_id}", lineGroupId)

	getLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLineGroupResp, err := client.Request("GET", getLineGroupPath, &getLineGroupOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getLineGroupResp)
}

func resourceLineGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	lineGroup, err := GetLineGroupById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS line group")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", lineGroup, nil)),
		d.Set("lines", utils.PathSearch("lines", lineGroup, nil)),
		d.Set("description", utils.PathSearch("description", lineGroup, nil)),
		// Attributes
		d.Set("status", utils.PathSearch("status", lineGroup, nil)),
		d.Set("created_at", utils.PathSearch("created_at", lineGroup, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", lineGroup, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLineGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if err := updateLineGroup(client, d); err != nil {
		return diag.FromErr(err)
	}

	err = waitForLineGroupCreatedOrUpdated(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLineGroupRead(ctx, d, meta)
}

func updateLineGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateLineGroupHttpUrl := "v2.1/linegroups/{linegroup_id}"
	updateLineGroupPath := client.Endpoint + updateLineGroupHttpUrl
	updateLineGroupPath = strings.ReplaceAll(updateLineGroupPath, "{linegroup_id}", d.Id())

	updateLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateLineGroupBodyParams(d)),
	}
	_, err := client.Request("PUT", updateLineGroupPath, &updateLineGroupOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS line group: %s", err)
	}
	return nil
}

func resourceLineGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		deleteLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"
		lineGroupId            = d.Id()
	)
	client, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	deleteLineGroupPath := client.Endpoint + deleteLineGroupHttpUrl
	deleteLineGroupPath = strings.ReplaceAll(deleteLineGroupPath, "{linegroup_id}", lineGroupId)

	deleteLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteLineGroupPath, &deleteLineGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DNS line group")
	}

	if err := waitForLineGroupDeleted(ctx, client, lineGroupId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForLineGroupDeleted(ctx context.Context, client *golangsdk.ServiceClient, lineGroupId string, timeOut time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING_CREATE", "PENDING_DELETE", "PENDING_UPDATE", "ERROR", "FREEZE", "DISABLE"},
		Refresh:      lineGroupStatusRefreshFunc(client, lineGroupId, nil),
		Timeout:      timeOut,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS line group (%s) to be deleted: %s", lineGroupId, err)
	}
	return nil
}

func lineGroupStatusRefreshFunc(client *golangsdk.ServiceClient, lineGroupId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lineGroup, err := GetLineGroupById(client, lineGroupId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "", err
		}

		status := utils.PathSearch("status", lineGroup, "").(string)
		if utils.StrSliceContains(targets, status) {
			return lineGroup, "COMPLETED", nil
		}

		return lineGroup, status, nil
	}
}
