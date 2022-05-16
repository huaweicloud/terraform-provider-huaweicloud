package cpts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	clusterTypeShare   = "shared-cluster-internet"
	clusterTypePrivate = "private-cluster"

	runStatusRunning   = 0
	runStatusFinished  = 2
	runStatusToRunning = 9

	operationEnable = "enable"
	operationStop   = "stop"
)

func ResourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,
		ReadContext:   resourceTaskRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 42),
			},

			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"benchmark_concurrency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      100,
				ValidateFunc: validation.IntBetween(0, 2000000),
			},

			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"operation": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{operationEnable, operationStop}, false),
			},

			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	projectId := int32(d.Get("project_id").(int))
	benchConcurrent := int32(d.Get("benchmark_concurrency").(int))
	createOpts := &model.CreateTaskRequest{
		Body: &model.CreateTaskRequestBody{
			Name:            d.Get("name").(string),
			ProjectId:       projectId,
			BenchConcurrent: &benchConcurrent,
		},
	}

	response, err := client.CreateTask(createOpts)
	if err != nil {
		return diag.Errorf("error creating CPTS task: %s", err)
	}

	if response.TaskId == nil {
		return diag.Errorf("error creating CPTS task: id not found in api response")
	}

	d.SetId(strconv.Itoa(int(*response.TaskId)))
	return resourceTaskRead(ctx, d, meta)
}

func resourceTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	response, err := client.ShowTask(&model.ShowTaskRequest{
		TaskId: int32(id),
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CPTS task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", response.Taskinfo.Name),
		d.Set("project_id", response.Taskinfo.ProjectId),
		d.Set("benchmark_concurrency", response.Taskinfo.BenchConcurrent),
		d.Set("status", response.Taskinfo.RunStatus),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	projectId := int32(d.Get("project_id").(int))

	if d.HasChanges("benchmark_concurrency", "name") {
		benchConcurrent := int32(d.Get("benchmark_concurrency").(int))
		_, err = client.UpdateTask(&model.UpdateTaskRequest{
			TaskId: int32(id),
			Body: &model.UpdateTaskRequestBody{
				Id:              int32(id),
				Name:            d.Get("name").(string),
				ProjectId:       projectId,
				BenchConcurrent: &benchConcurrent,
			},
		})

		if err != nil {
			return diag.Errorf("error updating the task %q: %s", id, err)
		}
	}

	// Enable or stop task
	if d.HasChange("operation") {
		op := d.Get("operation").(string)
		//Enable task
		if op == operationEnable {
			updateStatusRequest := model.UpdateTaskStatusRequest{
				TestSuiteId: projectId,
				TaskId:      int32(id),
			}

			clusterId := int32(d.Get("cluster_id").(int))
			if clusterId > 0 {
				updateStatusRequest.Body = &model.UpdateTaskStatusRequestBody{
					ClusterType: clusterTypePrivate,
					ClusterId:   clusterId,
					Status:      runStatusToRunning,
					NetworkInfo: &model.NetworkInfo{},
				}
			} else {
				updateStatusRequest.Body = &model.UpdateTaskStatusRequestBody{
					ClusterType: clusterTypeShare,
					ClusterId:   0,
					Status:      runStatusToRunning,
					NetworkInfo: &model.NetworkInfo{
						NetworkType: "internet",
					},
				}
			}

			_, err := client.UpdateTaskStatus(&updateStatusRequest)
			if err != nil {
				return diag.Errorf("error starting the task %q: %s", id, err)
			}

		}

		// stop task
		if op == operationStop {
			_, err := client.UpdateTaskStatus(&model.UpdateTaskStatusRequest{
				TestSuiteId: projectId,
				TaskId:      int32(id),
				Body: &model.UpdateTaskStatusRequestBody{
					ClusterType: clusterTypeShare,
					ClusterId:   -1,
					Status:      runStatusFinished,
					NetworkInfo: &model.NetworkInfo{},
				},
			})
			if err != nil {
				return diag.Errorf("error stopping the task %q: %s", id, err)
			}
		}

		err = waitingforTaskFinished(ctx, client, id, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaskRead(ctx, d, meta)
}

func resourceTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the task ID must be integer: %s", err)
	}

	deleteOpts := &model.DeleteTaskRequest{
		TaskId: int32(id),
	}

	_, err = client.DeleteTask(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting CPTS task %q: %s", id, err)
	}

	return nil
}
func waitingforTaskFinished(ctx context.Context, client *v1.CptsClient, id int64, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{fmt.Sprint(runStatusRunning), fmt.Sprint(runStatusToRunning)},
		Target:  []string{fmt.Sprint(runStatusFinished)},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.ShowTask(&model.ShowTaskRequest{TaskId: int32(id)})
			if err != nil {
				return nil, "", err
			}
			status := resp.Taskinfo.RunStatus
			return resp, fmt.Sprintf("%d", status), nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CPTS task (%d) to be finished: %s", id, err)
	}
	return nil
}
