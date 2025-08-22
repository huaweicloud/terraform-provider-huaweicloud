package apig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/features
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/features
func ResourceInstanceFeature() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceFeatureCreate,
		ReadContext:   resourceInstanceFeatureRead,
		UpdateContext: resourceInstanceFeatureUpdate,
		DeleteContext: resourceInstanceFeatureDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceFeatureImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
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
				Description: "Specified the ID of the dedicated instance to which the feature belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specified the name of the feature.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specified whether to enable the feature.",
			},
			// The format is `config="off"` or `config={\"max_timeout\":80000}`
			"config": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specified the detailed configuration of the feature.",
			},
		},
	}
}

func updateFeatureConfiguration(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceId, name string) (*http.Response, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/features"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildConfigInstanceFeatureParams(name, d.Get("enabled"), d.Get("config"))),
	}
	// The same configuration feature can only be modified once per minute.
	var resp *http.Response
	var reqErr error
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		resp, reqErr = client.Request("POST", path, &opts)
		isRetry, err := handleOperationError409(reqErr)
		if isRetry {
			// lintignore:R018
			time.Sleep(30 * time.Second)
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return resp, err
}

func resourceInstanceFeatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		name       = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}

	createResp, err := updateFeatureConfiguration(ctx, client, d, instanceId, name)
	if err != nil {
		return diag.Errorf("error creating instance feature: %s", err)
	}

	creatRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	featureName := utils.PathSearch("name", creatRespBody, "").(string)
	if featureName == "" {
		return diag.Errorf("unable to find the feature name from the API response")
	}

	d.SetId(featureName)
	return resourceInstanceFeatureRead(ctx, d, meta)
}

func handleOperationError409(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}

		errCode, searchErr := jmespath.Search("error_code", apiError)
		if searchErr != nil {
			return false, err
		}

		// APIG.3711: A configuration parameter can be modified only once per minute.
		if errCode == "APIG.3711" {
			return true, err
		}
	}
	return false, err
}

func buildConfigInstanceFeatureParams(featrueName string, enabled, cfg interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":   featrueName,
		"enable": enabled.(bool),
		"config": utils.ValueIgnoreEmpty(cfg.(string)),
	}
}

// GetInstanceFeature is a method that used to query the list of the features under specified APIG instance.
func GetInstanceFeature(client *golangsdk.ServiceClient, instanceId, featureName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/features?limit=500"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	offset := 0
	var feature interface{}
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		features := utils.PathSearch("features", respBody, make([]interface{}, 0)).([]interface{})
		if len(features) < 1 {
			break
		}

		if feature = utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", featureName), features, nil); feature != nil {
			return feature, nil
		}
		offset += len(features)
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceInstanceFeatureRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}

	feature, err := GetInstanceFeature(client, instanceId, d.Id())
	if err != nil {
		// When instance ID not exist, status code is 404, error code id APIG.3030
		return common.CheckDeletedDiag(d, err, "Instance feature configuration")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", feature, nil)),
		d.Set("enabled", utils.PathSearch("enable", feature, false)),
		d.Set("config", utils.PathSearch("config", feature, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceFeatureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		instanceId  = d.Get("instance_id").(string)
		featureName = d.Get("name").(string)
	)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	_, err = updateFeatureConfiguration(ctx, client, d, instanceId, featureName)
	if err != nil {
		return diag.Errorf("error creating instance feature: %s", err)
	}
	if err != nil {
		return diag.Errorf("error updating feature (%s) under specified instance (%s): %s", featureName, instanceId, err)
	}

	return resourceInstanceFeatureRead(ctx, d, meta)
}

func resourceInstanceFeatureDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceFeatureImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want <instance_id>/<name>, but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
	)
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
