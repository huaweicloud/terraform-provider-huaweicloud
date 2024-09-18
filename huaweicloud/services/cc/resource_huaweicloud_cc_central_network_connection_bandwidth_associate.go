package cc

import (
	"context"
	"fmt"
	"net/http"
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

// @API CC GET /v3/{domain_id}/gcn/central-network/{central_network_id}/connections
// @API CC PUT /v3/{domain_id}/gcn/central-network/{central_network_id}/connections/{connection_id}
func ResourceCentralNetworkConnectionBandwidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCentralNetworkConnectionBandwidthAssociateCreateOrUpdate,
		UpdateContext: resourceCentralNetworkConnectionBandwidthAssociateCreateOrUpdate,
		ReadContext:   resourceCentralNetworkConnectionBandwidthAssociateRead,
		DeleteContext: resourceCentralNetworkConnectionBandwidthAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCentralNetworkConnectionBandwidthAssociateImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"central_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the central network to which the connection belongs.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the connection.`,
			},
			"global_connection_bandwidth_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the global connection bandwidth.`,
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The bandwidth size of the connection.`,
			},
		},
	}
}

func resourceCentralNetworkConnectionBandwidthAssociateCreateOrUpdate(ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/connections/{connection_id}"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}
	connID := d.Get("connection_id").(string)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)
	path = strings.ReplaceAll(path, "{central_network_id}", d.Get("central_network_id").(string))
	path = strings.ReplaceAll(path, "{connection_id}", connID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	opt.JSONBody = buildCreateCentralNetworkConnectionBandwidthAssociateBodyParams(d)
	resp, err := client.Request("PUT", path, &opt)
	if err != nil {
		return diag.Errorf("error creating central network connection bandwidth associate: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	timeOut := schema.TimeoutUpdate
	if d.IsNewResource() {
		d.SetId(connID)
		timeOut = schema.TimeoutCreate
	}

	err = centralNetworkConnectionBandwidthAssociateWaitingForStateCompleted(ctx, d, meta, d.Timeout(timeOut))
	if err != nil {
		return diag.Errorf("error waiting for the central network connection bandwidth associate creation to complete: %s", err)
	}

	return resourceCentralNetworkConnectionBandwidthAssociateRead(ctx, d, meta)
}

func buildCreateCentralNetworkConnectionBandwidthAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"central_network_connection": map[string]interface{}{
			"bandwidth_type":                 "BandwidthPackage",
			"global_connection_bandwidth_id": d.Get("global_connection_bandwidth_id").(string),
			"bandwidth_size":                 d.Get("bandwidth_size").(int),
		},
	}
}
func centralNetworkConnectionBandwidthAssociateWaitingForStateCompleted(ctx context.Context,
	d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := readCentralNetworkConnection(meta, d)
			if err != nil {
				return nil, "ERROR", err
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`central_network_connections[0].state`, respBody, nil)
			if status == nil {
				return nil, "ERROR", fmt.Errorf("error parsing state from response body")
			}

			if utils.StrSliceContains([]string{"FAILED", "DELETED"}, status.(string)) {
				return respBody, "", fmt.Errorf("unexpected status '%s'", status.(string))
			}

			if status.(string) == "AVAILABLE" {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCentralNetworkConnectionBandwidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	resp, err := readCentralNetworkConnection(meta, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving central network connection: %s")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	connection := utils.PathSearch("central_network_connections", respBody, make([]interface{}, 0))
	if v, ok := connection.([]interface{}); ok && len(v) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	bandwidthType := utils.PathSearch("[0].bandwidth_type", connection, "TestBandwidth").(string)
	if bandwidthType == "TestBandwidth" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no bandwidth found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("central_network_id", d.Get("central_network_id").(string)),
		d.Set("connection_id", d.Id()),
		d.Set("global_connection_bandwidth_id", utils.PathSearch("[0].global_connection_bandwidth_id", connection, "").(string)),
		d.Set("bandwidth_size", int(utils.PathSearch("[0].bandwidth_size", connection, float64(0)).(float64))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func readCentralNetworkConnection(meta interface{}, d *schema.ResourceData) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/connections"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)
	path = strings.ReplaceAll(path, "{central_network_id}", d.Get("central_network_id").(string))
	path += "?id=" + d.Id()

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	return client.Request("GET", path, &opt)
}

func resourceCentralNetworkConnectionBandwidthAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/connections/{connection_id}"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)
	path = strings.ReplaceAll(path, "{central_network_id}", d.Get("central_network_id").(string))
	path = strings.ReplaceAll(path, "{connection_id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	opt.JSONBody = buildDeletedCentralNetworkConnectionBandwidthAssociate()

	_, err = client.Request("PUT", path, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "GCN.101505"),
			"error deleting central network connection bandwidth associate")
	}

	err = centralNetworkConnectionBandwidthAssociateWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the central network connection bandwidth associate deletion to complete: %s", err)
	}

	return nil
}

func buildDeletedCentralNetworkConnectionBandwidthAssociate() map[string]interface{} {
	return map[string]interface{}{
		"central_network_connection": map[string]interface{}{
			"bandwidth_type": "TestBandwidth",
		},
	}
}

func resourceCentralNetworkConnectionBandwidthAssociateImportState(_ context.Context,
	d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <central_network_id>/<connection_id>")
	}

	d.Set("central_network_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
