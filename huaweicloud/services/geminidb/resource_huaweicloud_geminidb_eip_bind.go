package geminidb

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var eipBindNonUpdatableParams = []string{
	"instance_id",
	"node_id",
	"public_ip",
	"public_ip_id",
}

// @API GeminiDB POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip
// @API GeminiDB POST /v3/{project_id}/instances
// @API GeminiDB GET /v3/{project_id}/jobs
// @API EIP GET /v1/{project_id}/publicips
func ResourceEipBind() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipBindCreate,
		ReadContext:   resourceEipBindRead,
		UpdateContext: resourceEipBindUpdate,
		DeleteContext: resourceEipBindDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEipBindImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(eipBindNonUpdatableParams),

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
			"node_id": {
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

func buildCreateEipBindBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       "BIND",
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
	}

	return bodyParams
}

func resourceEipBindCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		nodeId        = d.Get("node_id").(string)
		createHttpUrl = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{node_id}", nodeId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildCreateEipBindBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error binding EIP to GeminiDB instance node: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(nodeId)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error binding EIP to GeminiDB instance node: unable to find job ID from API response")
	}

	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the GeminiDB instance node bind EIP to complete: %s", err)
	}

	return resourceEipBindRead(ctx, d, meta)
}

func resourceEipBindRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	eipInfo, err := GetEipBindInfo(client, instanceId, d.Id())
	if err != nil {
		// When the GeminiDB instance does not exist, the response HTTP status code of the query API is 200
		// and return empty list
		return common.CheckDeletedDiag(d, err, "error retrieving the GeminiDB instance node EIP information")
	}

	ipAddress := utils.PathSearch("public_ip", eipInfo, "").(string)
	if ipAddress == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving the GeminiDB instance node EIP information")
	}

	var publicIpInfo interface{}
	if ipAddress != "" {
		vpcClient, err := cfg.NewServiceClient("vpc", region)
		if err != nil {
			return diag.Errorf("error creating VPC client: %s", err)
		}

		publicIpInfo, err = getAssociateEipInfo(vpcClient, ipAddress)
		if err != nil {
			log.Printf("[Warn] error retrieving the EIP information")
		}
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("node_id", utils.PathSearch("id", eipInfo, nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", eipInfo, nil)),
		d.Set("public_ip_id", utils.PathSearch("id", publicIpInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetEipBindInfo(client *golangsdk.ServiceClient, instanceId, nodeId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances?id={instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	eipInfo := utils.PathSearch(fmt.Sprintf("instances[].groups[].nodes[]|[?id=='%s']|[0]", nodeId), respBody, nil)
	if eipInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return eipInfo, nil
}

func getAssociateEipInfo(client *golangsdk.ServiceClient, ipAddress string) (interface{}, error) {
	httpUrl := "v1/{project_id}/publicips?public_ip_address={public_ip_address}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{public_ip_address}", ipAddress)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ipAddressInfo := utils.PathSearch("publicips|[0]", respBody, nil)
	if ipAddressInfo == nil {
		return nil, errors.New("error retrieving EIP information")
	}

	return ipAddressInfo, nil
}

func resourceEipBindUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func buildDeleteEipBindBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":       "UNBIND",
		"public_ip":    d.Get("public_ip"),
		"public_ip_id": d.Get("public_ip_id"),
	}

	return bodyParams
}

func resourceEipBindDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		nodeId     = d.Get("node_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/public-ip"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{node_id}", nodeId)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildDeleteEipBindBodyParams(d),
	}

	resp, err := client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.238009"),
			"error unbinding EIP from GeminiDB instance node")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error unbinding EIP from GeminiDB instance node: unable to find job ID from API response")
	}

	err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the GeminiDB instance node unbind EIP to complete: %s", err)
	}

	return nil
}

func resourceEipBindImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
