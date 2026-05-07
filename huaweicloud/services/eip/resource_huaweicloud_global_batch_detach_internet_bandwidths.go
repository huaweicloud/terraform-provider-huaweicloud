package eip

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/global-eips/batch-detach-internet-bandwidths
// @API EIP GET /v3/{domain_id}/geip/jobs/{job_id}
func ResourceBatchDetachInternetBandwidths() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDetachInternetBandwidthsCreate,
		ReadContext:   resourceBatchDetachInternetBandwidthsRead,
		UpdateContext: resourceBatchDetachInternetBandwidthsUpdate,
		DeleteContext: resourceBatchDetachInternetBandwidthsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"global_eips": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     globalEipInternetBandwidthSchema(),
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func globalEipInternetBandwidthSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"global_eip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func buildDetachInternetBandwidthsBodyParams(d *schema.ResourceData) map[string]interface{} {
	globalEipsRaw := d.Get("global_eips").([]interface{})
	globalEips := make([]map[string]interface{}, 0, len(globalEipsRaw))

	for _, item := range globalEipsRaw {
		v := item.(map[string]interface{})
		globalEips = append(globalEips, map[string]interface{}{
			"global_eip_id": v["global_eip_id"],
		})
	}

	bodyParams := map[string]interface{}{
		"global_eips": globalEips,
	}

	return bodyParams
}

func resourceBatchDetachInternetBandwidthsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{domain_id}/global-eips/batch-detach-internet-bandwidths"
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDetachInternetBandwidthsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch detaching internet bandwidths from global EIPs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job_id from the API response")
	}

	err = waitForJobCompleted(ctx, d.Timeout(schema.TimeoutCreate), jobId, cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for batch detach internet bandwidths job completed: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate EIP request ID: %s", err)
	}

	d.SetId(requestId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("job_id", jobId))
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBatchDetachInternetBandwidthsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDetachInternetBandwidthsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDetachInternetBandwidthsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch detach internet bandwidths. Deleting this
    resource will not clear of corresponding request record, but will only remove resource information from the tf state
    file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func waitForJobCompleted(ctx context.Context, timeout time.Duration, jobId string, domainID string,
	client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshJobStatusFunc(client, jobId, domainID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshJobStatusFunc(client *golangsdk.ServiceClient, jobId, domainID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobHttpUrl := "v3/{domain_id}/geip/jobs/{job_id}"
		getJobPath := client.Endpoint + getJobHttpUrl
		getJobPath = strings.ReplaceAll(getJobPath, "{domain_id}", domainID)
		getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", jobId)
		getJobOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getJobResp, err := client.Request("GET", getJobPath, &getJobOpt)
		if err != nil {
			return nil, "ERROR", err
		}
		getJobRespBody, err := utils.FlattenResponse(getJobResp)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("job.status", getJobRespBody, "").(string)
		if status == "" {
			return nil, "ERROR", errors.New("unable to find job status from the API response")
		}

		if status == "FINISH_ROLLBACK_SUCC" {
			return nil, "FAILURE", fmt.Errorf("job fail: %s", utils.PathSearch("job.error_message", getJobRespBody, nil))
		} else if status == "FINISH_SUCC" {
			return getJobRespBody, "SUCCESS", nil
		}
		return getJobRespBody, "PENDING", nil
	}
}
