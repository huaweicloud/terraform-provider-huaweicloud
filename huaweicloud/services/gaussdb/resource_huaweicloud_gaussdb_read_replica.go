package gaussdb

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

var gaussDbReadReplicaNonUpdatableParams = []string{
	"instance_id",
	"availability_zone",
	"flavor_ref",
	"configuration_id",
}

func ResourceGaussDbReadReplica() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbReadReplicaCreate,
		UpdateContext: resourceGaussDbReadReplicaUpdate,
		ReadContext:   resourceGaussDbReadReplicaRead,
		DeleteContext: resourceGaussDbReadReplicaDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussDbReadReplicaImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(gaussDbReadReplicaNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_ref": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_names": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDbReadReplicaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/readonly-nodes"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDbReadReplicaBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB read replica: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.Errorf("error creating GaussDB read replica: %s", err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GaussDB read replica: job_id is not found in API response")
	}

	readReplica, err := getGaussDbReadReplica(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.Errorf("error creating GaussDB read replica: %s", err)
	}
	readReplicaId := utils.PathSearch("nodes[?status=='creating']|[0].id", readReplica, nil)
	if readReplicaId == nil {
		return diag.Errorf("error creating GaussDB read replica: ID is not found")
	}

	d.SetId(readReplicaId.(string))

	if err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 60, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for GaussDB read replica job (%s) to be completed: %s", jobId, err)
	}

	return resourceGaussDbReadReplicaRead(ctx, d, meta)
}

func buildCreateGaussDbReadReplicaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_distribution": []map[string]interface{}{
			{
				"num":               1,
				"availability_zone": d.Get("availability_zone"),
				"flavor_ref":        d.Get("flavor_ref"),
				"configuration_id":  d.Get("configuration_id"),
			},
		},
	}
	return bodyParams
}

func resourceGaussDbReadReplicaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceRes, err := getGaussDBOpenGaussInstancesById(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	instance := utils.PathSearch(fmt.Sprintf("instances[0].nodes[?id=='%s']|[0]", d.Id()), instanceRes, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB instance")
	}

	readReplicaRes, err := getGaussDbReadReplica(client, d.Get("instance_id").(string))
	if err != nil {
		return diag.Errorf("error retrieving GaussDB read replica: %s", err)
	}
	readReplica := utils.PathSearch(fmt.Sprintf("nodes[?id=='%s']|[0]", d.Id()), readReplicaRes, nil)
	if readReplica == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB read replica")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", readReplica, nil)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", readReplica, nil)),
		d.Set("status", utils.PathSearch("status", readReplica, nil)),
		d.Set("flavor_ref", utils.PathSearch("flavor_ref", instance, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", instance, nil)),
		d.Set("data_ip", utils.PathSearch("data_ip", instance, nil)),
		d.Set("component_names", utils.PathSearch("component_names", instance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getGaussDbReadReplica(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/readonly-nodes"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceGaussDbReadReplicaUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbReadReplicaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/readonly-nodes"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDbReadReplicaBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.216030"),
			"error deleting GaussDB read replica")
	}

	deleteRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting GaussDB read replica: job_id is not found in API response")
	}
	if err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 60, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for GaussDB read replica deletion job (%s) to be completed: %s", jobId, err)
	}

	return nil
}

func buildDeleteGaussDbReadReplicaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_ids": []string{d.Id()},
	}
	return bodyParams
}

func resourceGaussDbReadReplicaImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
