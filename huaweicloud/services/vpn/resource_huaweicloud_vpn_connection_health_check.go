// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN POST /v5/{project_id}/connection-monitors
// @API VPN DELETE /v5/{project_id}/connection-monitors/{id}
// @API VPN GET /v5/{project_id}/connection-monitors/{id}
func ResourceConnectionHealthCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionHealthCheckCreate,
		ReadContext:   resourceConnectionHealthCheckRead,
		DeleteContext: resourceConnectionHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the VPN connection to monitor.`,
			},
			"destination_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The destination IP address of the VPN connection.`,
			},
			"source_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source IP address of the VPN connection.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the connection health check.`,
			},
		},
	}
}

func resourceConnectionHealthCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createConnectionHealthCheck: Create a VPN ConnectionHealthCheck.
	var (
		createConnectionHealthCheckHttpUrl = "v5/{project_id}/connection-monitors"
		createConnectionHealthCheckProduct = "vpn"
	)
	createConnectionHealthCheckClient, err := cfg.NewServiceClient(createConnectionHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createConnectionHealthCheckPath := createConnectionHealthCheckClient.Endpoint + createConnectionHealthCheckHttpUrl
	createConnectionHealthCheckPath = strings.ReplaceAll(createConnectionHealthCheckPath, "{project_id}",
		createConnectionHealthCheckClient.ProjectID)

	createConnectionHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createConnectionHealthCheckOpt.JSONBody = utils.RemoveNil(buildCreateConnectionHealthCheckBodyParams(d))
	createConnectionHealthCheckResp, err := createConnectionHealthCheckClient.Request("POST",
		createConnectionHealthCheckPath, &createConnectionHealthCheckOpt)
	if err != nil {
		return diag.Errorf("error creating VPN connection health check: %s", err)
	}

	createConnectionHealthCheckRespBody, err := utils.FlattenResponse(createConnectionHealthCheckResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("connection_monitor.id", createConnectionHealthCheckRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN connection health check: ID is not found in API response")
	}
	d.SetId(id)

	return resourceConnectionHealthCheckRead(ctx, d, meta)
}

func buildCreateConnectionHealthCheckBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"connection_monitor": map[string]interface{}{
			"vpn_connection_id": utils.ValueIgnoreEmpty(d.Get("connection_id")),
		},
	}
	return bodyParams
}

func resourceConnectionHealthCheckRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getConnectionHealthCheck: Query the VPN ConnectionHealthCheck detail
	var (
		getConnectionHealthCheckHttpUrl = "v5/{project_id}/connection-monitors/{id}"
		getConnectionHealthCheckProduct = "vpn"
	)
	getConnectionHealthCheckClient, err := cfg.NewServiceClient(getConnectionHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getConnectionHealthCheckPath := getConnectionHealthCheckClient.Endpoint + getConnectionHealthCheckHttpUrl
	getConnectionHealthCheckPath = strings.ReplaceAll(getConnectionHealthCheckPath, "{project_id}",
		getConnectionHealthCheckClient.ProjectID)
	getConnectionHealthCheckPath = strings.ReplaceAll(getConnectionHealthCheckPath, "{id}", d.Id())

	getConnectionHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getConnectionHealthCheckResp, err := getConnectionHealthCheckClient.Request("GET",
		getConnectionHealthCheckPath, &getConnectionHealthCheckOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN connection health check")
	}

	getConnectionHealthCheckRespBody, err := utils.FlattenResponse(getConnectionHealthCheckResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("destination_ip", utils.PathSearch("connection_monitor.destination_ip", getConnectionHealthCheckRespBody, nil)),
		d.Set("source_ip", utils.PathSearch("connection_monitor.source_ip", getConnectionHealthCheckRespBody, nil)),
		d.Set("status", utils.PathSearch("connection_monitor.status", getConnectionHealthCheckRespBody, nil)),
		d.Set("connection_id", utils.PathSearch("connection_monitor.vpn_connection_id", getConnectionHealthCheckRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceConnectionHealthCheckDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteConnectionHealthCheck: Delete an existing VPN ConnectionHealthCheck
	var (
		deleteConnectionHealthCheckHttpUrl = "v5/{project_id}/connection-monitors/{id}"
		deleteConnectionHealthCheckProduct = "vpn"
	)
	deleteConnectionHealthCheckClient, err := cfg.NewServiceClient(deleteConnectionHealthCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	deleteConnectionHealthCheckPath := deleteConnectionHealthCheckClient.Endpoint + deleteConnectionHealthCheckHttpUrl
	deleteConnectionHealthCheckPath = strings.ReplaceAll(deleteConnectionHealthCheckPath, "{project_id}",
		deleteConnectionHealthCheckClient.ProjectID)
	deleteConnectionHealthCheckPath = strings.ReplaceAll(deleteConnectionHealthCheckPath, "{id}", d.Id())

	deleteConnectionHealthCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteConnectionHealthCheckClient.Request("DELETE", deleteConnectionHealthCheckPath,
		&deleteConnectionHealthCheckOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "VPN.0001"),
			"error deleting VPN connection health check",
		)
	}

	return nil
}
