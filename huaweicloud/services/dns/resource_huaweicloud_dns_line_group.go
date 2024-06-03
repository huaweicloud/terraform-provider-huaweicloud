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

// @API DNS POST /v2.1/linegroups
// @API DNS GET /v2.1/linegroups/{linegroup_id}
// @API DNS PUT /v2.1/linegroups/{linegroup_id}
// @API DNS DELETE /v2.1/linegroups/{linegroup_id}
func ResourceDNSLineGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSLineGroupCreate,
		UpdateContext: resourceDNSLineGroupUpdate,
		ReadContext:   resourceDNSLineGroupRead,
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
				Description: `Specifies the line group name.`,
			},
			"lines": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the lines.`,
				Set:         schema.HashString,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the line group description. A maximum of 255 characters are allowed.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Resource creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Resource update time.`,
			},
		},
	}
}

func resourceDNSLineGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	createDNSLineGroupClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	// createDNSLineGroup: Create the DNS line group
	if err := createDNSLineGroup(createDNSLineGroupClient, d); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := waitForDNSLineGroupCreateOrUpdate(ctx, createDNSLineGroupClient, d, timeout); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSLineGroupRead(ctx, d, meta)
}

func createDNSLineGroup(lineGroupClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var createDNSLineGroupHttpUrl = "v2.1/linegroups"

	createDNSLineGroupPath := lineGroupClient.Endpoint + createDNSLineGroupHttpUrl
	createDNSLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createDNSLineGroupOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateDNSLineGroupBodyParams(d))
	createDNSLineGroupResp, err := lineGroupClient.Request("POST", createDNSLineGroupPath,
		&createDNSLineGroupOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS line group: %s", err)
	}

	createDNSLineGroupRespBody, err := utils.FlattenResponse(createDNSLineGroupResp)
	if err != nil {
		return err
	}

	id, err := jmespath.Search("line_id", createDNSLineGroupRespBody)
	if err != nil {
		return fmt.Errorf("error creating DNS line group: ID is not found in API response")
	}
	d.SetId(id.(string))
	return nil
}

func waitForDNSLineGroupCreateOrUpdate(ctx context.Context, lineGroupClient *golangsdk.ServiceClient,
	d *schema.ResourceData, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING_CREATE", "PENDING_UPDATE"},
		Refresh:      dnsLineGroupStatusRefreshFunc(d, lineGroupClient),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS line group (%s) to be ACTIVE: %s", d.Id(), err)
	}
	return nil
}

func buildCreateOrUpdateDNSLineGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnsLines := d.Get("lines").(*schema.Set).List()

	return map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"lines":       utils.ValueIgnoreEmpty(dnsLines),
	}
}

func resourceDNSLineGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getDNSLineGroup: Query the DNS line group detail.
	var (
		getDNSLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"
		getDNSLineGroupProduct = "dns"
	)
	getDNSLineGroupClient, err := cfg.NewServiceClient(getDNSLineGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	getDNSLineGroupPath := getDNSLineGroupClient.Endpoint + getDNSLineGroupHttpUrl
	getDNSLineGroupPath = strings.ReplaceAll(getDNSLineGroupPath, "{linegroup_id}", d.Id())

	getDNSLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSLineGroupResp, err := getDNSLineGroupClient.Request("GET", getDNSLineGroupPath,
		&getDNSLineGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS line group")
	}

	getDNSLineGroupRespBody, err := utils.FlattenResponse(getDNSLineGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDNSLineGroupRespBody, nil)),
		d.Set("lines", utils.PathSearch("lines", getDNSLineGroupRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getDNSLineGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getDNSLineGroupRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getDNSLineGroupRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getDNSLineGroupRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSLineGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateDNSLineGroupClient, err := cfg.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	updateDNSLineGroupChanges := []string{
		"name",
		"description",
		"lines",
	}

	if d.HasChanges(updateDNSLineGroupChanges...) {
		// updateDNSLineGroup: Update the DNS line group
		if err := updateDNSLineGroup(updateDNSLineGroupClient, d); err != nil {
			return diag.FromErr(err)
		}

		timeout := d.Timeout(schema.TimeoutUpdate)
		if err := waitForDNSLineGroupCreateOrUpdate(ctx, updateDNSLineGroupClient, d, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDNSLineGroupRead(ctx, d, meta)
}

func updateDNSLineGroup(lineGroupClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var updateDNSLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"

	updateDNSLineGroupPath := lineGroupClient.Endpoint + updateDNSLineGroupHttpUrl
	updateDNSLineGroupPath = strings.ReplaceAll(updateDNSLineGroupPath, "{linegroup_id}", d.Id())

	updateDNSLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSLineGroupOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateDNSLineGroupBodyParams(d))
	_, err := lineGroupClient.Request("PUT", updateDNSLineGroupPath,
		&updateDNSLineGroupOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS line group: %s", err)
	}
	return nil
}

func resourceLineGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDNSLineGroup: Delete the DNS line group
	var (
		deleteDNSLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"
		deleteDNSLineGroupProduct = "dns"
	)
	deleteDNSLineGroupClient, err := cfg.NewServiceClient(deleteDNSLineGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	deleteDNSLineGroupPath := deleteDNSLineGroupClient.Endpoint + deleteDNSLineGroupHttpUrl
	deleteDNSLineGroupPath = strings.ReplaceAll(deleteDNSLineGroupPath, "{linegroup_id}", d.Id())

	deleteDNSLineGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = deleteDNSLineGroupClient.Request("DELETE", deleteDNSLineGroupPath, &deleteDNSLineGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting DNS line group: %s", err)
	}

	if err := waitForDNSLineGroupDeleted(ctx, deleteDNSLineGroupClient, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForDNSLineGroupDeleted(ctx context.Context, lineGroupClient *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING_CREATE", "PENDING_DELETE", "PENDING_UPDATE", "ERROR", "FREEZE", "DISABLE"},
		Refresh:      dnsLineGroupStatusRefreshFunc(d, lineGroupClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS line group (%s) to be DELETED: %s", d.Id(), err)
	}
	return nil
}

func dnsLineGroupStatusRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var getDNSLineGroupHttpUrl = "v2.1/linegroups/{linegroup_id}"

		getDNSLineGroupPath := client.Endpoint + getDNSLineGroupHttpUrl
		getDNSLineGroupPath = strings.ReplaceAll(getDNSLineGroupPath, "{linegroup_id}", d.Id())
		getDNSLineGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		getDNSLineGroupResp, err := client.Request("GET", getDNSLineGroupPath, &getDNSLineGroupOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "", err
		}

		getDNSLineGroupRespBody, err := utils.FlattenResponse(getDNSLineGroupResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getDNSLineGroupRespBody, "")
		return getDNSLineGroupRespBody, parseStatus(status.(string)), nil
	}
}
