package antiddos

import (
	"context"
	"math"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ANTI-DDOS POST /v1/{project_id}/antiddos/default-config
// @API ANTI-DDOS DELETE /v1/{project_id}/antiddos/default-config
// @API ANTI-DDOS GET /v1/{project_id}/antiddos/default-config
func ResourceDefaultProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDefaultProtectionPolicyCreate,
		ReadContext:   resourceDefaultProtectionPolicyRead,
		UpdateContext: resourceDefaultProtectionPolicyUpdate,
		DeleteContext: resourceDefaultProtectionPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"traffic_threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The traffic cleaning threshold in Mbps.`,
			},
		},
	}
}

// ReadDefaultProtectionPolicy Test case will use this method, so the first letter is capitalized.
func ReadDefaultProtectionPolicy(client *golangsdk.ServiceClient) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/antiddos/default-config"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

// configDefaultProtectionPolicy Only the field traffic_pos_id is meaningful, other fields are meaningless.
func configDefaultProtectionPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, policyResp interface{}) error {
	createPath := client.Endpoint + "v1/{project_id}/antiddos/default-config"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	cleaningAccessPosID := utils.PathSearch("cleaning_access_pos_id", policyResp, float64(1)).(float64)
	createOpt.JSONBody = map[string]interface{}{
		"enable_L7":           utils.PathSearch("enable_L7", policyResp, nil),
		"traffic_pos_id":      getTrafficThresholdID(d.Get("traffic_threshold").(int)),
		"http_request_pos_id": utils.PathSearch("http_request_pos_id", policyResp, nil),
		// Make sure the `cleaning_access_pos_id` not larger than `8`.
		// Field `cleaning_access_pos_id` has no practical meaning in the request.
		// This will avoid error in partners cloud.
		"cleaning_access_pos_id": int(math.Min(cleaningAccessPosID, 8)),
		"app_type_id":            utils.PathSearch("app_type_id", policyResp, nil),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceDefaultProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS v1 client: %s", err)
	}

	policyResp, err := ReadDefaultProtectionPolicy(client)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS default protection policy in creation operation: %s", err)
	}

	if err := configDefaultProtectionPolicy(client, d, policyResp); err != nil {
		return diag.Errorf("error configuring Anti-DDoS default protection policy in creation operation: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return resourceDefaultProtectionPolicyRead(ctx, d, meta)
}

// The default protection policy always has a value and does not require checkDeleted logic.
func resourceDefaultProtectionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS v1 client: %s", err)
	}

	policyResp, err := ReadDefaultProtectionPolicy(client)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS default protection policy: %s", err)
	}

	trafficPosID := utils.PathSearch("traffic_pos_id", policyResp, nil)
	if trafficPosID == nil {
		return diag.Errorf("error retrieving Anti-DDoS default protection policy: traffic_pos_id is not found in" +
			" read API response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("traffic_threshold", getTrafficThresholdBandwidth(int(trafficPosID.(float64)))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDefaultProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS v1 client: %s", err)
	}

	policyResp, err := ReadDefaultProtectionPolicy(client)
	if err != nil {
		return diag.Errorf("error retrieving Anti-DDoS default protection policy in update operation: %s", err)
	}

	if err := configDefaultProtectionPolicy(client, d, policyResp); err != nil {
		return diag.Errorf("error configuring Anti-DDoS default protection policy in update operation: %s", err)
	}

	return resourceDefaultProtectionPolicyRead(ctx, d, meta)
}

// The actual effect of deleting the default protection policy is to set the Traffic Cleaning Threshold to 120.
func resourceDefaultProtectionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS v1 client: %s", err)
	}

	deletePath := client.Endpoint + "v1/{project_id}/antiddos/default-config"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Anti-DDoS default protection policy: %s", err)
	}
	return nil
}
