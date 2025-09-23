package dc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

var virtualInterfaceSwitchoverNonUpdatableParams = []string{"resource_id", "operation", "resource_type"}

// @API DC POST /v3/{project_id}/dcaas/switchover-test
// @API DC GET /v3/{project_id}/dcaas/switchover-test
func ResourceVirtualInterfaceSwitchover() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualInterfaceSwitchoverCreate,
		ReadContext:   resourceVirtualInterfaceSwitchoverRead,
		UpdateContext: resourceVirtualInterfaceSwitchoverUpdate,
		DeleteContext: resourceVirtualInterfaceSwitchoverDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(virtualInterfaceSwitchoverNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operate_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualInterfaceSwitchoverCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/switchover-test"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateVirtualInterfaceSwitchoverBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DC virtual interface switchover: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("switchover_test_record.id", createRespBody, "")
	if id == "" {
		return diag.Errorf("error creating DC virtual interface switchover: id is not found in API response")
	}
	d.SetId(id.(string))

	err = checkVirtualInterfaceSwitchoverFinish(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVirtualInterfaceSwitchoverRead(ctx, d, meta)
}

func buildCreateVirtualInterfaceSwitchoverBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_id":   d.Get("resource_id"),
		"operation":     d.Get("operation"),
		"resource_type": utils.ValueIgnoreEmpty(d.Get("resource_type")),
	}

	return map[string]interface{}{
		"switchover_test_record": bodyParams,
	}
}

func resourceVirtualInterfaceSwitchoverRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	switchover, err := getVirtualInterfaceSwitchoverRecord(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting DC virtual interface switchover record")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resource_id", utils.PathSearch("resource_id", switchover, nil)),
		d.Set("resource_type", utils.PathSearch("resource_type", switchover, nil)),
		d.Set("operation", utils.PathSearch("operation", switchover, nil)),
		d.Set("start_time", utils.PathSearch("start_time", switchover, nil)),
		d.Set("end_time", utils.PathSearch("end_time", switchover, nil)),
		d.Set("operate_status", utils.PathSearch("operate_status", switchover, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func checkVirtualInterfaceSwitchoverFinish(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      virtualInterfaceSwitchoverRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DC virtual interface switchover(%s) to be completed: %s ", d.Id(), err)
	}
	return nil
}

func virtualInterfaceSwitchoverRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		switchover, err := getVirtualInterfaceSwitchoverRecord(client, d)
		if err != nil {
			return nil, "ERROR", err
		}

		operateStatus := utils.PathSearch("operate_status", switchover, "").(string)
		if operateStatus == "COMPLETE" {
			return switchover, operateStatus, nil
		}
		if operateStatus == "ERROR" {
			return nil, operateStatus, errors.New("error test virtual interface switchover")
		}

		return switchover, "PENDING", nil
	}
}

func getVirtualInterfaceSwitchoverRecord(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/dcaas/switchover-test?resource_id={resource_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{resource_id}", d.Get("resource_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("switchover_test_records[?id=='%s']|[0]", d.Id())
	switchover := utils.PathSearch(searchPath, getRespBody, nil)
	if switchover == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return switchover, nil
}

func resourceVirtualInterfaceSwitchoverUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVirtualInterfaceSwitchoverDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DC virtual interface switchover is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
