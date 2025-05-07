package kafka

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

// @API Kafka POST /v2/{project_id}/instances/{instance_id}/connector
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
// @API Kafka POST /v2/{project_id}/kafka/instances/{instance_id}/delete-connector
func ResourceDmsKafkaSmartConnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaSmartConnectCreate,
		ReadContext:   resourceDmsKafkaSmartConnectRead,
		DeleteContext: resourceDmsKafkaSmartConnectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaSmartConnectImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
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
				Description: `Specifies the ID of the Kafka instance.`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the specification code of the smart connect.`,
			},
			"bandwidth": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the bandwidth of the smart connect.`,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the node count of the smart connect.`,
			},
		},
	}
}

func resourceDmsKafkaSmartConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createKafkaSmartConnect: create DMS kafka smart connect
	var (
		createKafkaSmartConnectHttpUrl = "v2/{project_id}/instances/{instance_id}/connector"
		createKafkaSmartConnectProduct = "dmsv2"
	)
	createKafkaSmartConnectClient, err := cfg.NewServiceClient(createKafkaSmartConnectProduct, region)

	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createKafkaSmartConnectPath := createKafkaSmartConnectClient.Endpoint + createKafkaSmartConnectHttpUrl
	createKafkaSmartConnectPath = strings.ReplaceAll(createKafkaSmartConnectPath, "{project_id}",
		createKafkaSmartConnectClient.ProjectID)
	createKafkaSmartConnectPath = strings.ReplaceAll(createKafkaSmartConnectPath, "{instance_id}", instanceID)

	createKafkaSmartConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createKafkaSmartConnectOpt.JSONBody = utils.RemoveNil(buildCreateKafkaSmartConnectBodyParams(d))

	// create kafka smart connect is allowd only when the instance status is RUNNING.
	retryFunc := func() (interface{}, bool, error) {
		createKafkaSmartConnectResp, createErr := createKafkaSmartConnectClient.Request("POST",
			createKafkaSmartConnectPath, &createKafkaSmartConnectOpt)
		retry, err := handleMultiOperationsError(createErr)
		return createKafkaSmartConnectResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(createKafkaSmartConnectClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating DMS kafka smart connect: %v", err)
	}

	kafkaSmartConnectRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	connectorId := utils.PathSearch("connector_id", kafkaSmartConnectRespBody, nil)
	if connectorId == nil {
		return diag.Errorf("error creating DMS kafka smart connect: connector ID is not found in API response")
	}
	d.SetId(connectorId.(string))

	// enable smart connect, need to wait for the instance status to be RUNNING so that the job could be completed.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXTENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(createKafkaSmartConnectClient, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instanceID, err)
	}

	return resourceDmsKafkaSmartConnectRead(ctx, d, meta)
}

func buildCreateKafkaSmartConnectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"spec_code":     utils.ValueIgnoreEmpty(d.Get("storage_spec_code")),
		"specification": utils.ValueIgnoreEmpty(d.Get("bandwidth")),
		"node_cnt":      utils.ValueIgnoreEmpty(d.Get("node_count")),
	}
	return bodyParams
}

func resourceDmsKafkaSmartConnectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getKafkaSmartConnect: query DMS kafka smart connect
	var (
		getKafkaSmartConnectHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getKafkaSmartConnectProduct = "dms"
	)
	getKafkaSmartConnectClient, err := cfg.NewServiceClient(getKafkaSmartConnectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	getKafkaSmartConnectPath := getKafkaSmartConnectClient.Endpoint + getKafkaSmartConnectHttpUrl
	getKafkaSmartConnectPath = strings.ReplaceAll(getKafkaSmartConnectPath, "{project_id}",
		getKafkaSmartConnectClient.ProjectID)
	getKafkaSmartConnectPath = strings.ReplaceAll(getKafkaSmartConnectPath, "{instance_id}", instanceID)

	getKafkaSmartConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectResp, err := getKafkaSmartConnectClient.Request("GET", getKafkaSmartConnectPath,
		&getKafkaSmartConnectOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka smart connect")
	}

	getKafkaSmartConnectRespBody, respBodyerr := utils.FlattenResponse(getKafkaSmartConnectResp)
	if respBodyerr != nil {
		return diag.FromErr(respBodyerr)
	}

	connectorId := utils.PathSearch("connector_id", getKafkaSmartConnectRespBody, nil)
	if connectorId != d.Id() {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsKafkaSmartConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteKafkaSmartConnect: delete DMS kafka smart connect
	var (
		deleteKafkaSmartConnectHttpUrl = "v2/{project_id}/kafka/instances/{instance_id}/delete-connector"
		deleteKafkaSmartConnectProduct = "dmsv2"
	)
	deleteKafkaSmartConnectClient, err := cfg.NewServiceClient(deleteKafkaSmartConnectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteKafkaSmartConnectPath := deleteKafkaSmartConnectClient.Endpoint + deleteKafkaSmartConnectHttpUrl
	deleteKafkaSmartConnectPath = strings.ReplaceAll(deleteKafkaSmartConnectPath, "{project_id}",
		deleteKafkaSmartConnectClient.ProjectID)
	deleteKafkaSmartConnectPath = strings.ReplaceAll(deleteKafkaSmartConnectPath, "{instance_id}", instanceID)

	deleteKafkaSmartConnectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteKafkaSmartConnectOpt.JSONBody = utils.RemoveNil(buildDeleteKafkaSmartConnectBodyParams(d))

	// delete kafka smart connect is allowd only when the instance status is RUNNING.
	retryFunc := func() (interface{}, bool, error) {
		deleteKafkaSmartConnectResp, deleteErr := deleteKafkaSmartConnectClient.Request("POST",
			deleteKafkaSmartConnectPath, &deleteKafkaSmartConnectOpt)
		retry, err := handleMultiOperationsError(deleteErr)
		return deleteKafkaSmartConnectResp, retry, err
	}
	_, retryErr := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(deleteKafkaSmartConnectClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if retryErr != nil {
		return diag.Errorf("error deleting DMS kafka smart connect: %v", err)
	}

	// delete smart connect, need to wait for the instance status to be RUNNING so that the job could be completed.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CONNECTOR_DELETING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(deleteKafkaSmartConnectClient, instanceID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", instanceID, err)
	}

	return nil
}

func buildDeleteKafkaSmartConnectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_id": d.Get("instance_id"),
	}
	return bodyParams
}

// resourceDmsKafkaSmartConnectImportState is used to import an id with format <instance_id>/<id>
func resourceDmsKafkaSmartConnectImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <instance_id>/<id>")
	}
	d.Set("instance_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
