package taurusdb

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/public-ips/bind
// @API GaussDBforNoSQL GET /v3.1/{project_id}/instances/{instance_id}
// @API GaussDBforNoSQL PUT /v3/{project_id}/instances/{instance_id}/public-ips/unbind
// @API EIP GET /v1/{project_id}/publicips
func ResourceGaussMysqlEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussMysqlEipAssociateCreate,
		ReadContext:   resourceGaussMysqlEipAssociateRead,
		DeleteContext: resourceGaussMysqlEipAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a GaussDB MySQL instance.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the EIP address to be bound.`,
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the EIP.`,
			},
		},
	}
}

func resourceGaussMysqlEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/public-ips/bind"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildGaussMySQLEipAssociateBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error binding EIP to GaussDB MySQL(%s): %s", instanceID, err)
	}

	d.SetId(instanceID)

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error associating EIP to GaussDB MySQL(%s), job_id is not found in the response.", instanceID)
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGaussMysqlEipAssociateRead(ctx, d, meta)
}

func buildGaussMySQLEipAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
	}
	return bodyParams
}

func resourceGaussMysqlEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.280238"),
			"error retrieving EIP associated with GaussDB MySQL")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	publicIP := utils.PathSearch("instance.public_ips", getRespBody, "").(string)
	if publicIP == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving EIP associated with GaussDB MySQL")
	}

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	publicID, err := common.GetEipIDbyAddress(vpcClient, publicIP, "all_granted_eps")
	if err != nil {
		return diag.Errorf("unable to get ID of public IP(%s): %s", publicIP, err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance.id", getRespBody, nil)),
		d.Set("public_ip", publicIP),
		d.Set("public_ip_id", publicID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussMysqlEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/public-ips/unbind"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error unbinding EIP from GaussDB MySQL(%s): %s", instanceID, err)
	}

	deleteRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error unbinding EIP from GaussDB MySQL(%s), job_id is not found in the response.", instanceID)
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
