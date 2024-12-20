package live

import (
	"context"
	"fmt"
	"strings"
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

// @API Live PUT /v1/{project_id}/domain/hls
// @API Live GET /v1/{project_id}/domain/hls
func ResourceHlsConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHlsConfigurationCreate,
		ReadContext:   resourceHlsConfigurationRead,
		UpdateContext: resourceHlsConfigurationUpdate,
		DeleteContext: resourceHlsConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceHlsConfigImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"hls_fragment": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"hls_ts_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"hls_min_frags": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceHlsConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = updateHlsConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating Live HLS configuration: %s", err)
	}

	// After successfully calling the update HLS configuration API,
	// When calling the query API immediately, there may be a situation where the configuration has not been updated,
	// So we need to add a waiting logic.
	err = waitingForHlsConfigurationComplete(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for Live HLS configuration creation to complete: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceHlsConfigurationRead(ctx, d, meta)
}

func waitingForHlsConfigurationComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getRespBody, err := ReadHlsConfiguration(client, d.Get("domain_name").(string))
			if err != nil {
				return nil, "ERROR", err
			}

			curJson := utils.PathSearch("application", getRespBody, make([]interface{}, 0))
			curArray := curJson.([]interface{})
			if len(curArray) == 0 {
				return nil, "ERROR", fmt.Errorf("the application is not found in query API response")
			}

			// Check if the response of the query API matches the target value.
			hlsFragmentResp := utils.PathSearch("hls_fragment", curArray[0], float64(0)).(float64)
			hlsTsCountResp := utils.PathSearch("hls_ts_count", curArray[0], float64(0)).(float64)
			hlsMinFragsResp := utils.PathSearch("hls_min_frags", curArray[0], float64(0)).(float64)
			applicationRepMap := d.Get("application").([]interface{})[0].(map[string]interface{})
			if int(hlsFragmentResp) == applicationRepMap["hls_fragment"].(int) &&
				int(hlsTsCountResp) == applicationRepMap["hls_ts_count"].(int) &&
				int(hlsMinFragsResp) == applicationRepMap["hls_min_frags"].(int) {
				return getRespBody, "COMPLETED", nil
			}

			return getRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func updateHlsConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	hlsConfigurationHttpUrl := "v1/{project_id}/domain/hls"
	hlsConfigurationPath := client.Endpoint + hlsConfigurationHttpUrl
	hlsConfigurationPath = strings.ReplaceAll(hlsConfigurationPath, "{project_id}", client.ProjectID)

	hlsConfigurationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildHlsConfigurationBodyParams(d),
	}

	_, err := client.Request("PUT", hlsConfigurationPath, &hlsConfigurationOpt)
	return err
}

func buildHlsConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	hlsConfigurationBodyParams := map[string]interface{}{
		"push_domain": d.Get("domain_name"),
		"application": buildHlsConfigurationApplicationBodyParams(d),
	}

	return hlsConfigurationBodyParams
}

func buildHlsConfigurationApplicationBodyParams(d *schema.ResourceData) []map[string]interface{} {
	application := d.Get("application").([]interface{})
	applicationBodyParams := make([]map[string]interface{}, 0)
	for _, v := range application {
		applicationMap := v.(map[string]interface{})
		applicationBodyParam := map[string]interface{}{
			"name":          applicationMap["name"],
			"hls_fragment":  applicationMap["hls_fragment"],
			"hls_ts_count":  applicationMap["hls_ts_count"],
			"hls_min_frags": applicationMap["hls_min_frags"],
		}
		applicationBodyParams = append(applicationBodyParams, applicationBodyParam)
	}

	return applicationBodyParams
}

func ReadHlsConfiguration(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/domain/hls"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?push_domain=%s", domainName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceHlsConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := ReadHlsConfiguration(client, d.Get("domain_name").(string))
	// When the `domain_name` does not exist, calling the query API will return a `404` status code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live HLS configuration")
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("push_domain", getRespBody, nil)),
		d.Set("application", flattenApplication(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApplication(getRespBody interface{}) []interface{} {
	if getRespBody == nil {
		return nil
	}

	curJson := utils.PathSearch("application", getRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0)
	for _, v := range curArray {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"hls_fragment":  utils.PathSearch("hls_fragment", v, nil),
			"hls_ts_count":  utils.PathSearch("hls_ts_count", v, nil),
			"hls_min_frags": utils.PathSearch("hls_min_frags", v, nil),
		})
	}

	return result
}

func resourceHlsConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = updateHlsConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating Live HLS configuration: %s", err)
	}

	err = waitingForHlsConfigurationComplete(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error waiting for Live HLS configuration update to complete: %s", err)
	}

	return resourceHlsConfigurationRead(ctx, d, meta)
}

func resourceHlsConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceHlsConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	if importedId == "" {
		return nil, fmt.Errorf("invalid format specified for import ID, `domain_name` is empty")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("domain_name", importedId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
