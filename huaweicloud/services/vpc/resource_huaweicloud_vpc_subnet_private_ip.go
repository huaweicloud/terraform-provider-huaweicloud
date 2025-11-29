package vpc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v1/{project_id}/privateips
// @API VPC GET /v1/{project_id}/privateips/{privateip_id}
// @API VPC DELETE /v1/{project_id}/privateips/{privateip_id}

var privateIPSyncNonUpdatableParams = []string{"subnet_id", "ip_address", "device_owner"}

func ResourceSubnetPrivateIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetPrivateIPCreate,
		ReadContext:   resourceSubnetPrivateIPRead,
		UpdateContext: resourceSubnetPrivateIPUpdate,
		DeleteContext: resourceSubnetPrivateIPDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(privateIPSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSubnetPrivateIPBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"privateips": []map[string]interface{}{
			{
				"subnet_id":    d.Get("subnet_id"),
				"device_owner": utils.ValueIgnoreEmpty(d.Get("device_owner")),
				"ip_address":   utils.ValueIgnoreEmpty(d.Get("ip_address")),
			},
		},
	}
	return bodyParams
}

func resourceSubnetPrivateIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	ctreateSubnetPrivateIPHttpUrl := "v1/{project_id}/privateips"
	ctreateSubnetPrivateIPPath := client.Endpoint + ctreateSubnetPrivateIPHttpUrl
	ctreateSubnetPrivateIPPath = strings.ReplaceAll(ctreateSubnetPrivateIPPath, "{project_id}", client.ProjectID)
	createSubnetPrivateIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createSubnetPrivateIPOpt.JSONBody = utils.RemoveNil(buildSubnetPrivateIPBodyParams(d))
	createSubnetPrivateIPResp, err := client.Request("POST", ctreateSubnetPrivateIPPath, &createSubnetPrivateIPOpt)
	if err != nil {
		return diag.Errorf("error creating VPC subnet private IP: %s", err)
	}

	createSubnetPrivateIPRespBody, err := utils.FlattenResponse(createSubnetPrivateIPResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("privateips[0].id", createSubnetPrivateIPRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPC subnet private IP: ID is not found in API response")
	}
	d.SetId(id)

	return resourceSubnetPrivateIPRead(ctx, d, meta)
}

func resourceSubnetPrivateIPRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	getSubnetPrivateIPHttpUrl := "v1/{project_id}/privateips/{privateip_id}"
	getSubnetPrivateIPPath := client.Endpoint + getSubnetPrivateIPHttpUrl
	getSubnetPrivateIPPath = strings.ReplaceAll(getSubnetPrivateIPPath, "{project_id}", client.ProjectID)
	getSubnetPrivateIPPath = strings.ReplaceAll(getSubnetPrivateIPPath, "{privateip_id}", d.Id())
	getSubnetPrivateIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getSubnetPrivateIPResp, err := client.Request("GET", getSubnetPrivateIPPath, &getSubnetPrivateIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC subnet private IP")
	}

	getSubnetPrivateIPRespBody, err := utils.FlattenResponse(getSubnetPrivateIPResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("privateip.subnet_id", getSubnetPrivateIPRespBody, nil)),
		d.Set("ip_address", utils.PathSearch("privateip.ip_address", getSubnetPrivateIPRespBody, nil)),
		d.Set("status", utils.PathSearch("privateip.status", getSubnetPrivateIPRespBody, nil)),
		d.Set("device_owner", utils.PathSearch("privateip.device_owner", getSubnetPrivateIPRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSubnetPrivateIPUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSubnetPrivateIPDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	deleteSubnetPrivateIPHttpUrl := "v1/{project_id}/privateips/{privateip_id}"
	deleteSubnetPrivateIPPath := client.Endpoint + deleteSubnetPrivateIPHttpUrl
	deleteSubnetPrivateIPPath = strings.ReplaceAll(deleteSubnetPrivateIPPath, "{project_id}", client.ProjectID)
	deleteSubnetPrivateIPPath = strings.ReplaceAll(deleteSubnetPrivateIPPath, "{privateip_id}", d.Id())
	deleteSubnetPrivateIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteSubnetPrivateIPPath, &deleteSubnetPrivateIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC subnet private IP")
	}
	return nil
}
