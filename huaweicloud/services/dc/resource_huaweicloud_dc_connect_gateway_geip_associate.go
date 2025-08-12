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

var connectGatewayGeipAssociateNonUpdatableParams = []string{"connect_gateway_id", "global_eip_id", "type"}

// @API DC POST /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/binding-global-eips
// @API DC GET /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/binding-global-eips
// @API DC POST /v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/unbinding-global-eips
func ResourceDcConnectGatewayGeipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcConnectGatewayGeipAssociateCreate,
		ReadContext:   resourceDcConnectGatewayGeipAssociateRead,
		UpdateContext: resourceDcConnectGatewayGeipAssociateUpdate,
		DeleteContext: resourceDcConnectGatewayGeipAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: globalGatewayGeipAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(connectGatewayGeipAssociateNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"connect_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_eip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"global_eip_segment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ie_vtep_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDcConnectGatewayGeipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/binding-global-eips"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connect_gateway_id}", d.Get("connect_gateway_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDcConnectGatewayGeipAssociateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DC connect gateway global EIP associate: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("connect_gateway_id").(string), d.Get("global_eip_id").(string)))

	err = waitForConnectGatewayGeipAvailable(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcConnectGatewayGeipAssociateRead(ctx, d, meta)
}

func buildCreateDcConnectGatewayGeipAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"global_eips": []map[string]interface{}{
			{
				"global_eip_id": d.Get("global_eip_id"),
				"type":          utils.ValueIgnoreEmpty(d.Get("type")),
			},
		},
	}

	return bodyParams
}

func resourceDcConnectGatewayGeipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	globalEip, err := getConnectGatewayGeip(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DC connect gateway global EIP associate")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("connect_gateway_id", d.Get("connect_gateway_id")),
		d.Set("global_eip_id", utils.PathSearch("global_eip_id", globalEip, nil)),
		d.Set("type", utils.PathSearch("type", globalEip, nil)),
		d.Set("global_eip_segment_id", utils.PathSearch("global_eip_segment_id", globalEip, nil)),
		d.Set("status", utils.PathSearch("status", globalEip, nil)),
		d.Set("error_message", utils.PathSearch("error_message", globalEip, nil)),
		d.Set("cidr", utils.PathSearch("cidr", globalEip, nil)),
		d.Set("address_family", utils.PathSearch("address_family", globalEip, nil)),
		d.Set("ie_vtep_ip", utils.PathSearch("ie_vtep_ip", globalEip, nil)),
		d.Set("created_time", utils.PathSearch("created_time", globalEip, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDcConnectGatewayGeipAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcConnectGatewayGeipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/unbinding-global-eips"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{connect_gateway_id}", d.Get("connect_gateway_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteDcConnectGatewayGeipAssociateBodyParams(d),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DC connect gateway global EIP associate")
	}

	err = waitForConnectGatewayGeipDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteDcConnectGatewayGeipAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"global_eips": []map[string]interface{}{
			{
				"global_eip_id": d.Get("global_eip_id"),
			},
		},
	}

	return bodyParams
}

func waitForConnectGatewayGeipDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      connectGatewayGeipRefreshFunc(client, d),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DC connect gateway global EIP to be deleted: %s ", err)
	}
	return nil
}

func waitForConnectGatewayGeipAvailable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"BIND_SUCCESSFULLY"},
		Refresh:      connectGatewayGeipRefreshFunc(client, d),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DC connect gateway global EIP to available: %s ", err)
	}
	return nil
}

func connectGatewayGeipRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		globalEip, err := getConnectGatewayGeip(client, d)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", globalEip, "").(string)
		if status == "BIND_SUCCESSFULLY" {
			return globalEip, status, nil
		}
		if status == "ERROR" {
			return nil, status, errors.New("error retrieving connect gateway global EIP")
		}

		return globalEip, "PENDING", nil
	}
}

func getConnectGatewayGeip(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/dcaas/connect-gateways/{connect_gateway_id}/binding-global-eips?global_eip_id={global_eip_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connect_gateway_id}", d.Get("connect_gateway_id").(string))
	getPath = strings.ReplaceAll(getPath, "{global_eip_id}", d.Get("global_eip_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	globalEip := utils.PathSearch("global_eips|[0]", getRespBody, nil)
	if globalEip == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return globalEip, nil
}

func globalGatewayGeipAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <connect_gateway_id>/<global_eip_id>")
	}
	mErr := multierror.Append(nil,
		d.Set("connect_gateway_id", parts[0]),
		d.Set("global_eip_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
