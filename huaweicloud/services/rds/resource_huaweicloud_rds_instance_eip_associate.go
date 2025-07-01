package rds

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipAssociateNonUpdatableParams = []string{"instance_id", "public_ip", "public_ip_id"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/public-ip
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
func ResourceRdsInstanceEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsInstanceEipAssociateCreate,
		ReadContext:   resourceRdsInstanceEipAssociateRead,
		UpdateContext: resourceRdsInstanceEipAssociateUpdate,
		DeleteContext: resourceRdsInstanceEipAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(eipAssociateNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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
			"instance_id": {
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

func resourceRdsInstanceEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	bodyParams := buildBindRdsInstanceEipBodyParams(d)
	err = bindOrUnbindEip(ctx, d, client, schema.TimeoutCreate, bodyParams)
	if err != nil {
		return diag.Errorf("error binding EIP to RDS instance(%s): %s", d.Id(), err)
	}

	return resourceRdsInstanceEipAssociateRead(ctx, d, meta)
}

func buildBindRdsInstanceEipBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
		"is_bind":      true,
	}
	return bodyParams
}

func resourceRdsInstanceEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
		product = "rds"
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
		return common.CheckDeletedDiag(d, err, "error retrieving EIP associated with GaussDB MySQL")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	publicIP := utils.PathSearch("instances|[0].public_ips[0]", getRespBody, nil)
	if publicIP == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving EIP associated with GaussDB MySQL")
	}

	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	publicID, err := common.GetEipIDbyAddress(vpcClient, publicIP.(string), "all_granted_eps")
	if err != nil {
		return diag.Errorf("unable to get ID of public IP(%s): %s", publicIP, err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instances|[0].id", getRespBody, nil)),
		d.Set("public_ip", publicIP),
		d.Set("public_ip_id", publicID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRdsInstanceEipAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsInstanceEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	bodyParams := buildUnbindRdsInstanceEipBodyParams()
	err = bindOrUnbindEip(ctx, d, client, schema.TimeoutDelete, bodyParams)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200001"),
			fmt.Sprintf("error unbinding EIP from RDS instance(%s)", d.Id()))
	}

	return nil
}

func bindOrUnbindEip(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	bodyParam map[string]interface{}) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/public-ip"
	)

	instanceID := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(bodyParam)

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		d.SetId(instanceID)
	}

	deleteRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("job_id is not found in the response")
	}

	return checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(timeout))
}

func buildUnbindRdsInstanceEipBodyParams() map[string]interface{} {
	bodyParams := map[string]interface{}{
		"is_bind": false,
	}
	return bodyParams
}
