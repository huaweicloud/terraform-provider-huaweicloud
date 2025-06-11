package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventUnblockIpNonUpdatableParams = []string{"data_list", "data_list.*.host_id", "data_list.*.src_ip",
	"data_list.*.login_type", "enterprise_project_id"}

// @API HSS PUT /v5/{project_id}/event/blocked-ip
func ResourceEventUnblockIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventUnblockIpCreate,
		ReadContext:   resourceEventUnblockIpRead,
		UpdateContext: resourceEventUnblockIpUpdate,
		DeleteContext: resourceEventUnblockIpDelete,

		CustomizeDiff: config.FlexibleForceNew(eventUnblockIpNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"login_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateEventUnblockIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	dataList := d.Get("data_list").([]interface{})
	dataListRequestBody := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		dataListRequestBody = append(dataListRequestBody, map[string]interface{}{
			"host_id":    utils.PathSearch("host_id", v, nil),
			"src_ip":     utils.PathSearch("src_ip", v, nil),
			"login_type": utils.PathSearch("login_type", v, nil),
		})
	}

	return map[string]interface{}{
		"data_list": dataListRequestBody,
	}
}

func resourceEventUnblockIpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/event/blocked-ip"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	enterpriseProjectID := cfg.GetEnterpriseProjectID(d)
	if enterpriseProjectID != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", enterpriseProjectID)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateEventUnblockIpBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error unblocking HSS IP: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceEventUnblockIpRead(ctx, d, meta)
}

func resourceEventUnblockIpRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventUnblockIpUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventUnblockIpDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to HSS unblock IP. Deleting this resource
    will not change the current HSS unblock IP, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
