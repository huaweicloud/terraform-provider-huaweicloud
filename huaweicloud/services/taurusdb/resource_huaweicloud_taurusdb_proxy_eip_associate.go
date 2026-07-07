package taurusdb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var proxyEipAssociateNoneUpdatableParams = []string{
	"instance_id", "proxy_id", "public_ip", "public_ip_id",
}

// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/bind
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/proxies
// @API EIP GET /v1/{project_id}/publicips
func ResourceTaurusDBProxyEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBProxyEipAssociateCreate,
		ReadContext:   resourceTaurusDBProxyEipAssociateRead,
		UpdateContext: resourceTaurusDBProxyEipAssociateUpdate,
		DeleteContext: resourceTaurusDBProxyEipAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBProxyEipAssociateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(proxyEipAssociateNoneUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceTaurusDBProxyEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		proxyID    = d.Get("proxy_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/bind"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createPath = strings.ReplaceAll(createPath, "{proxy_id}", proxyID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildTaurusDBProxyEipAssociateBodyParams(d),
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		needRetry, err := handleMultiOperationsError(err)
		return res, needRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     taurusDBProxyInstanceStateRefreshFunc(client, instanceID, proxyID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error binding EIP to TaurusDB(%s) Proxy(%s): %s", instanceID, proxyID, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, proxyID))

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error associating EIP to TaurusDB(%s) Proxy(%s), job_id is not found in the response.", instanceID, proxyID)
	}
	expectedEip := d.Get("public_ip").(string)
	err = waitForProxyEipBindTaskCompleted(ctx, client, instanceID, proxyID, expectedEip, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for associating EIP to TaurusDB(%s) Proxy(%s) to complete: %s", instanceID, proxyID, err)
	}
	return resourceTaurusDBProxyEipAssociateRead(ctx, d, meta)
}

func buildTaurusDBProxyEipAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
		"bind":         "true",
	}
	return bodyParams
}

func taurusDBProxyInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, proxyID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		proxyStatus, err := getTaurusDBProxyStatus(client, instanceID, proxyID)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return proxyStatus, "DELETED", nil
			}
			return nil, "", err
		}
		return proxyStatus, proxyStatus.(string), nil
	}
}

func waitForProxyEipBindTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceID, proxyID,
	expectedEip string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"BINDING"},
		Target:       []string{"BOUND"},
		Refresh:      taurusDBProxyEipBindStateRefreshFunc(client, instanceID, proxyID, expectedEip),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func taurusDBProxyEipBindStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, proxyID, expectedEip string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		proxyEip, err := GetTaurusDBProxyEip(client, instanceID, proxyID)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return expectedEip, "BINDING", nil
			}
			return nil, "", err
		}
		if proxyEip != expectedEip {
			return proxyEip, "BINDING", nil
		}
		return proxyEip, "BOUND", nil
	}
}

func getTaurusDBProxyStatus(client *golangsdk.ServiceClient, instanceID, proxyID string) (interface{}, error) {
	searchExpression := fmt.Sprintf("proxy_list[?proxy.pool_id=='%s']|[0].proxy.status", proxyID)
	proxyStatus, err := getGaussDBProxy(client, instanceID, searchExpression)
	if err != nil {
		return nil, err
	}
	return proxyStatus, nil
}

func GetTaurusDBProxyEip(client *golangsdk.ServiceClient, instanceID, proxyID string) (interface{}, error) {
	searchExpression := fmt.Sprintf("proxy_list[?proxy.pool_id=='%s']|[0].proxy.eip", proxyID)
	proxyEip, err := getGaussDBProxy(client, instanceID, searchExpression)
	if err != nil {
		return nil, err
	}
	return proxyEip, nil
}

func resourceTaurusDBProxyEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		proxyID    = d.Get("proxy_id").(string)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	proxyEip, err := GetTaurusDBProxyEip(client, instanceID, proxyID)
	if err != nil {
		// Check if the proxy itself is gone (real 404 from API)
		var errDefault404 golangsdk.ErrDefault404
		if errors.As(err, &errDefault404) {
			return common.CheckDeletedDiag(d, err, "error retrieving EIP associated with TaurusDB Proxy")
		}
		return diag.Errorf("error retrieving EIP associated with TaurusDB Proxy: %s", err)
	}

	publicIP, ok := proxyEip.(string)
	if !ok {
		return diag.Errorf("error retrieving EIP associated with TaurusDB Proxy: %s", err)
	}

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	publicIpID, err := common.GetEipIDbyAddress(vpcClient, publicIP, "all_granted_eps")
	if err != nil {
		return diag.Errorf("unable to get ID of public IP(%s): %s", publicIP, err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("proxy_id", proxyID),
		d.Set("public_ip", publicIP),
		d.Set("public_ip_id", publicIpID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBProxyEipAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBProxyEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceID = d.Get("instance_id").(string)
		proxyID    = d.Get("proxy_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/bind"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deletePath = strings.ReplaceAll(deletePath, "{proxy_id}", proxyID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildTaurusDBProxyEipUnbindBodyParams(d),
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", deletePath, &deleteOpt)
		needRetry, err := handleMultiOperationsError(err)
		return res, needRetry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     taurusDBProxyInstanceStateRefreshFunc(client, instanceID, proxyID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error unbinding EIP from TaurusDB(%s) Proxy(%s): %s", instanceID, proxyID, err)
	}
	deleteRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error unbinding EIP from TaurusDB(%s) Proxy(%s), job_id is not found in the response.", instanceID, proxyID)
	}
	err = waitForProxyEipUnBindTaskCompleted(ctx, client, instanceID, proxyID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for unbinding EIP from TaurusDB(%s) Proxy(%s) to complete: %s", instanceID, proxyID, err)
	}
	return nil
}

func buildTaurusDBProxyEipUnbindBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"public_ip": d.Get("public_ip"),
		"bind":      "false",
	}
	return bodyParams
}

func waitForProxyEipUnBindTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	instanceID, proxyID string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"UNBINDING"},
		Target:       []string{"UNBOUND"},
		Refresh:      taurusDBProxyEipUnbindStateRefreshFunc(client, instanceID, proxyID),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return err
	}
	return nil
}

func taurusDBProxyEipUnbindStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, proxyID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		proxyEip, err := GetTaurusDBProxyEip(client, instanceID, proxyID)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "null", "UNBOUND", nil
			}
			return nil, "", err
		}
		return proxyEip, "UNBINDING", nil
	}
}

func resourceTaurusDBProxyEipAssociateImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<proxy_id>, but got `%s`", d.Id())
	}

	instanceId := parts[0]
	proxyId := parts[1]

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("proxy_id", proxyId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
