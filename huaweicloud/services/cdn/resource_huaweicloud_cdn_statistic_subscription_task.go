package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN POST /v1/cdn/statistics/subscription-tasks
// @API CDN GET /v1/cdn/statistics/subscription-tasks
// @API CDN PUT /v1/cdn/statistics/subscription-tasks/{id}
// @API CDN DELETE /v1/cdn/statistics/subscription-tasks/{id}
func ResourceStatisticSubscriptionTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatisticSubscriptionTaskCreate,
		ReadContext:   resourceStatisticSubscriptionTaskRead,
		UpdateContext: resourceStatisticSubscriptionTaskUpdate,
		DeleteContext: resourceStatisticSubscriptionTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceStatisticSubscriptionTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the subscription task.`,
			},
			"period_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The type of the subscription task.`,
			},
			"emails": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The email addresses to receive the operation reports.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of domain names to subscribe.`,
			},
			"report_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the operation report.`,
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the subscription task, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the subscription task, in RFC3339 format.`,
			},
		},
	}
}

func buildStatisticSubscriptionTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"period_type": d.Get("period_type"),
		"emails":      d.Get("emails"),
		"domain_name": d.Get("domain_name"),
		"report_type": d.Get("report_type"),
	}
}

func createStatisticSubscriptionTask(client *golangsdk.ServiceClient, bodyParams map[string]interface{}) (interface{}, error) {
	httpUrl := "v1/cdn/statistics/subscription-tasks"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: bodyParams,
		OkCodes:  []int{200, 204},
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func listStatisticSubscriptionTasks(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/cdn/statistics/subscription-tasks?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		tasks := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)
		if len(tasks) < limit {
			break
		}
		offset += len(tasks)
	}

	return result, nil
}

func getStatisticSubscriptionTaskByName(client *golangsdk.ServiceClient, taskName string) (interface{}, error) {
	tasks, err := listStatisticSubscriptionTasks(client)
	if err != nil {
		return nil, err
	}

	task := utils.PathSearch(fmt.Sprintf("[?name == '%s']|[0]", taskName), tasks, nil)
	if task == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/cdn/statistics/subscription-tasks",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the subscription task with name '%s' has been removed", taskName)),
			},
		}
	}
	return task, nil
}

func GetStatisticSubscriptionTaskById(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	tasks, err := listStatisticSubscriptionTasks(client)
	if err != nil {
		return nil, err
	}

	task := utils.PathSearch(fmt.Sprintf("[?id == `%s`]|[0]", taskId), tasks, nil)
	if task == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/cdn/statistics/subscription-tasks",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the subscription task with ID '%s' has been removed", taskId)),
			},
		}
	}
	return task, nil
}

func resourceStatisticSubscriptionTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	bodyParams := buildStatisticSubscriptionTaskBodyParams(d)
	task, err := createStatisticSubscriptionTask(client, bodyParams)
	if err != nil {
		return diag.Errorf("error creating CDN statistic subscription task: %s", err)
	}

	taskId := utils.PathSearch("id", task, float64(0)).(float64)
	d.SetId(strconv.Itoa(int(taskId)))

	return resourceStatisticSubscriptionTaskRead(ctx, d, meta)
}

func resourceStatisticSubscriptionTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		taskId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	task, err := GetStatisticSubscriptionTaskById(client, taskId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting subscription task (%s)", taskId))
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", task, nil)),
		d.Set("period_type", utils.PathSearch("period_type", task, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", task, nil)),
		d.Set("report_type", utils.PathSearch("report_type", task, nil)),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", task, float64(0)).(float64))/1000, false)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", task, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateStatisticSubscriptionTask(client *golangsdk.ServiceClient, taskId string, bodyParams map[string]interface{}) error {
	httpUrl := "v1/cdn/statistics/subscription-tasks/{id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{id}", taskId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: bodyParams,
		OkCodes:  []int{200},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceStatisticSubscriptionTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		taskId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	bodyParams := buildStatisticSubscriptionTaskBodyParams(d)
	if err := updateStatisticSubscriptionTask(client, taskId, bodyParams); err != nil {
		return diag.Errorf("error updating subscription task (%s): %s", taskId, err)
	}

	return resourceStatisticSubscriptionTaskRead(ctx, d, meta)
}

func deleteStatisticSubscriptionTask(client *golangsdk.ServiceClient, taskId string) error {
	httpUrl := "v1/cdn/statistics/subscription-tasks/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", taskId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceStatisticSubscriptionTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                               = meta.(*config.Config)
		taskId                            = d.Id()
		statisticSubscriptionTaskErrCodes = []string{
			"CDN.0001",
		}
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = deleteStatisticSubscriptionTask(client, taskId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", statisticSubscriptionTaskErrCodes...),
			fmt.Sprintf("error deleting subscription task (%s)", taskId))
	}

	return nil
}

func resourceStatisticSubscriptionTaskImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()

	// Check if it's a numeric ID (task ID) or a name
	if _, err := strconv.ParseInt(importedId, 10, 64); err == nil {
		// It's a numeric ID
		d.SetId(importedId)
		return []*schema.ResourceData{d}, nil
	}

	// It's a name, need to query by name
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	task, err := getStatisticSubscriptionTaskByName(client, importedId)
	if err != nil {
		return nil, err
	}
	taskId := utils.PathSearch("id", task, float64(0)).(float64)
	d.SetId(strconv.Itoa(int(taskId)))

	return []*schema.ResourceData{d}, nil
}
