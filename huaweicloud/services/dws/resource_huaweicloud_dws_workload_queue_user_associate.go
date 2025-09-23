package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users/batch-create
// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users/batch-delete
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users
func ResourceWorkloadQueueUserAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkloadQueueUserAssociateCreate,
		UpdateContext: resourceWorkloadQueueUserAssociateUpdate,
		ReadContext:   resourceWorkloadQueueUserAssociateRead,
		DeleteContext: resourceWorkloadQueueUserAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkloadQueueUserAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the workload queue name to associate with the users.`,
			},
			"user_names": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the user names bound to the workload queue.",
			},
		},
	}
}

func resourceWorkloadQueueUserAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	queueName := d.Get("queue_name").(string)
	userNames := d.Get("user_names").(*schema.Set)
	err = bindUserNamesToQueue(client, clusterId, queueName, userNames)
	if err != nil {
		return diag.Errorf("error binding the users to the workload queue: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, queueName))

	return resourceWorkloadQueueUserAssociateRead(ctx, d, meta)
}

func buildWorkloadQueueUserAssociateBodyParams(queueName string, users *schema.Set) map[string]interface{} {
	result := make([]interface{}, users.Len())
	for i, v := range users.List() {
		result[i] = map[string]interface{}{
			"user_name": v.(string),
		}
	}

	return map[string]interface{}{
		"queue_name": queueName,
		"user_list":  result,
	}
}

func GetAssociatedUserNames(client *golangsdk.ServiceClient, clusterId, queueName string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getPath = strings.ReplaceAll(getPath, "{queue_name}", queueName)

	opt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	for {
		getPathWithOffset := fmt.Sprintf("%s?limit=100&offset=%d", getPath, offset)
		requestResp, err := client.Request("GET", getPathWithOffset, &opt)
		if err != nil {
			// "DWS.0047": The cluster ID does not exist, the status code is 404.
			// "DWS.9999": The workload queue does not exist, the status code is 500, error code is "DWS.9999".
			return nil, common.ConvertExpected500ErrInto404Err(err, "error_code", "DWS.9999")
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		userNames := utils.PathSearch("user_list[*].user_name", respBody, make([]interface{}, 0)).([]interface{})
		if len(userNames) < 1 {
			break
		}
		result = append(result, userNames...)
		offset += len(userNames)
	}

	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return result, nil
}

func resourceWorkloadQueueUserAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	associatedUsers, err := GetAssociatedUserNames(client, d.Get("cluster_id").(string), d.Get("queue_name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workload queue associated user")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_names", associatedUsers.([]interface{})),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the workload queue associated users: %s", err)
	}
	return nil
}

func bindUserNamesToQueue(client *golangsdk.ServiceClient, clusterId, queueName string, userNames *schema.Set) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users/batch-create"
	bindPath := client.Endpoint + httpUrl
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{cluster_id}", clusterId)
	bindPath = strings.ReplaceAll(bindPath, "{queue_name}", queueName)

	bindOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         buildWorkloadQueueUserAssociateBodyParams(queueName, userNames),
	}
	_, err := client.Request("POST", bindPath, &bindOpt)
	return err
}

func unbindUserNamesFromQueue(client *golangsdk.ServiceClient, clusterId, queueName string, userNames *schema.Set) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/queues/{queue_name}/users/batch-delete"
	unbindPath := client.Endpoint + httpUrl
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{cluster_id}", clusterId)
	unbindPath = strings.ReplaceAll(unbindPath, "{queue_name}", queueName)

	unbindOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         buildWorkloadQueueUserAssociateBodyParams(queueName, userNames),
	}
	_, err := client.Request("POST", unbindPath, &unbindOpt)
	return err
}

func resourceWorkloadQueueUserAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	if d.HasChanges("user_names") {
		var (
			clusterId      = d.Get("cluster_id").(string)
			queueName      = d.Get("queue_name").(string)
			oldRaw, newRaw = d.GetChange("user_names")
			addSet         = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
			rmSet          = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
		)

		if rmSet.Len() > 0 {
			err = unbindUserNamesFromQueue(client, clusterId, queueName, rmSet)
			if err != nil {
				return diag.Errorf("error unbinding users from the workload queue: %s", err)
			}
		}

		if addSet.Len() > 0 {
			err = bindUserNamesToQueue(client, clusterId, queueName, addSet)
			if err != nil {
				if err != nil {
					return diag.Errorf("error binding users to the workload queue: %s", err)
				}
			}
		}
	}
	return resourceWorkloadQueueUserAssociateRead(ctx, d, meta)
}

func resourceWorkloadQueueUserAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	err = unbindUserNamesFromQueue(client, d.Get("cluster_id").(string), d.Get("queue_name").(string), d.Get("user_names").(*schema.Set))
	if err != nil {
		// The API respondse includes these cases about resource not found.
		//   1. "DWS.0047": The cluster ID does not exist, the status code is 404.
		//   2. The unbound user has been deleted, the key name of the error code is "workload_res_code" and the value is 1 (int type).
		//   3. When the queue name does not exist, the status code is 200, so the CheckDeleted logic cannot be judged. (2024-08-26).
		return common.CheckDeletedDiag(d, common.ConvertExpected500ErrInto404Err(err, "workload_res_code"),
			"error unbinding associated users")
	}
	return nil
}

func resourceWorkloadQueueUserAssociateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<queue_name>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
		d.Set("queue_name", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
