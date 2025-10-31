package rgc

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC POST /v1/best-practice/detect
// @API RGC GET /v1/best-practice/detection-overview
// @API RGC GET /v1/best-practice/status
func ResourceBestPractice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBestPracticeCreate,
		UpdateContext: resourceBestPracticeUpdate,
		ReadContext:   resourceBestPracticeRead,
		DeleteContext: resourceBestPracticeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"total_score": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"detect_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBestPracticeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createBestPracticeHttpUrl = "v1/best-practice/detect"
		createBestPracticeProduct = "rgc"
	)

	createBestPracticeClient, err := cfg.NewServiceClient(createBestPracticeProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	createBestPracticePath := createBestPracticeClient.Endpoint + createBestPracticeHttpUrl
	createBestPracticeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = createBestPracticeClient.Request("POST", createBestPracticePath, &createBestPracticeOpt)
	if err != nil {
		return diag.Errorf("error creating best practice: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"in_progress"},
		Target:       []string{"succeeded"},
		Refresh:      bestPracticeStateRefreshFunc(createBestPracticeClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for best practice to create: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return resourceBestPracticeRead(ctx, d, meta)
}

func resourceBestPracticeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getBestPracticeOverviewHttpUrl = "v1/best-practice/detection-overview"
		getBestPracticeOverviewProduct = "rgc"
	)

	getBestPracticeOverviewClient, err := cfg.NewServiceClient(getBestPracticeOverviewProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getBestPracticeOverviewPath := getBestPracticeOverviewClient.Endpoint + getBestPracticeOverviewHttpUrl

	getBestPracticeOverviewOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBestPracticeOverviewResp, err := getBestPracticeOverviewClient.Request("GET", getBestPracticeOverviewPath, &getBestPracticeOverviewOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving best practice")
	}

	getBestPracticeOverviewRespBody, err := utils.FlattenResponse(getBestPracticeOverviewResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("total_score", utils.PathSearch("total_score", getBestPracticeOverviewRespBody, nil)),
		d.Set("detect_time", utils.PathSearch("detect_time", getBestPracticeOverviewRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())

}

func resourceBestPracticeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBestPracticeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RGC best-practice resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func bestPracticeStateRefreshFunc(client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getBestPracticeStatusHttpUrl := "v1/best-practice/status"
		getBestPracticeStatusPath := client.Endpoint + getBestPracticeStatusHttpUrl

		getBestPracticeStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getBestPracticeStatusResp, err := client.Request("GET", getBestPracticeStatusPath, &getBestPracticeStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getBestPracticeStatusRespBody, err := utils.FlattenResponse(getBestPracticeStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getBestPracticeStatusRespBody, "").(string)
		if status == "" {
			return nil, "", errors.New("error getting best practice status")
		}

		return getBestPracticeStatusRespBody, status, nil
	}
}
