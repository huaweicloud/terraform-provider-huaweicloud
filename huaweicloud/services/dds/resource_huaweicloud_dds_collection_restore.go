package dds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/restore/collections
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSCollectionRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSCollectionRestoreCreate,
		ReadContext:   resourceDDSCollectionRestoreRead,
		DeleteContext: resourceDDSCollectionRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
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
			"restore_collections": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"restore_database_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"collections": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"old_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"restore_collection_time": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"new_name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceDDSCollectionRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}

	instId := d.Get("instance_id").(string)

	// restore instance
	restoreHttpUrl := "v3/{project_id}/instances/{instance_id}/restore/collections"
	restorePath := client.Endpoint + restoreHttpUrl
	restorePath = strings.ReplaceAll(restorePath, "{project_id}", client.ProjectID)
	restorePath = strings.ReplaceAll(restorePath, "{instance_id}", instId)

	restoreOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectionRestoreBodyParams(d)),
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", restorePath, &restoreOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	restoreResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error restoring collections to instance(%s): %s", instId, err)
	}

	// get job ID
	restoreRespBody, err := utils.FlattenResponse(restoreResp.(*http.Response))
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}
	jobID := utils.PathSearch("job_id", restoreRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	d.SetId(jobID)

	// wait for job complete
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, jobID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s", jobID, err)
	}

	return nil
}

func buildCollectionRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	restoreCollections := d.Get("restore_collections").([]interface{})

	opts := make([]map[string]interface{}, 0, len(restoreCollections))
	for _, v := range restoreCollections {
		restoreCollection := v.(map[string]interface{})
		opts = append(opts, map[string]interface{}{
			"database":              restoreCollection["database"],
			"restore_database_time": utils.ValueIgnoreEmpty(restoreCollection["restore_database_time"]),
			"collections":           buildCollectionRestoreBodyParamsCollections(restoreCollection["collections"].([]interface{})),
		})
	}

	return map[string]interface{}{
		"restore_collections": opts,
	}
}

func buildCollectionRestoreBodyParamsCollections(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	opts := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		params := v.(map[string]interface{})
		opts = append(opts, map[string]interface{}{
			"old_name":                params["old_name"],
			"restore_collection_time": params["restore_collection_time"],
			"new_name":                utils.ValueIgnoreEmpty(params["new_name"]),
		})
	}

	return opts
}

func resourceDDSCollectionRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDSCollectionRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restore resource is not supported. The restore resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
