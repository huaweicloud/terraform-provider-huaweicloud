package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/public-ips
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceOpenGaussEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussEipAssociateCreate,
		ReadContext:   resourceOpenGaussEipAssociateRead,
		DeleteContext: resourceOpenGaussEipAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenGaussEipAssociateImportState,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOpenGaussEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	err = bindOrUnbindEip(ctx, d, client, "BIND", schema.TimeoutCreate)
	if err != nil {
		return diag.Errorf("error binding EIP to GaussDB OpenGauss(%s): %s", instanceId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceId, d.Get("node_id").(string)))

	return resourceOpenGaussEipAssociateRead(ctx, d, meta)
}

func resourceOpenGaussEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	res, err := getGaussDBOpenGaussInstancesById(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	instance := utils.PathSearch("instances[0]", res, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	node := utils.PathSearch(fmt.Sprintf("nodes[?id == '%s']|[0]", d.Get("node_id").(string)), instance, nil)
	if node == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	publicIp := utils.PathSearch("public_ip", node, nil)
	if publicIp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	boundEips, err := listInstanceEips(d, client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP associated with GaussDB OpenGauss")
	}

	boundEip := utils.PathSearch(fmt.Sprintf("public_ips[?public_ip_address=='%s']|[0]", publicIp.(string)),
		*boundEips, nil)
	if boundEip == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("id", instance, nil)),
		d.Set("node_id", utils.PathSearch("id", node, nil)),
		d.Set("public_ip", publicIp),
		d.Set("public_ip_id", utils.PathSearch("public_ip_id", boundEip, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listInstanceEips(d *schema.ResourceData, client *golangsdk.ServiceClient) (*interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/public-ips"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	return &listRespBody, nil
}

func resourceOpenGaussEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	err = bindOrUnbindEip(ctx, d, client, "UNBIND", schema.TimeoutDelete)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertUndefinedErrInto404Err(
				common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200712"),
				409, "error_code", "DBS.200011"),
			fmt.Sprintf("error unbinding EIP from GaussDB OpenGauss(%s)", instanceId))
	}

	return nil
}

func bindOrUnbindEip(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, action, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip"
	)

	instanceId := d.Get("instance_id").(string)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{node_id}", d.Get("node_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildBindOrUnbindEipBodyParams(d, action))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", updatePath, &updateOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	updateRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("job_id is not found in the response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(timeout))
	if err != nil {
		return err
	}

	return nil
}

func buildBindOrUnbindEipBodyParams(d *schema.ResourceData, action string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       action,
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
	}
	return bodyParams
}

func resourceOpenGaussEipAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<node_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("node_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
