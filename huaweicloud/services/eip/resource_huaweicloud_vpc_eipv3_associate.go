package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{project_id}/eip/publicips/{publicip_id}/associate-instance
// @API EIP GET /v3/{project_id}/eip/publicips/{publicip_id}
// @API EIP POST /v3/{project_id}/eip/publicips/{publicip_id}/disassociate-instance
func ResourceEipv3Associate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipv3AssociateCreate,
		ReadContext:   resourceEipv3AssociateRead,
		DeleteContext: resourceEipv3AssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"publicip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"associate_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"associate_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEipv3AssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/eip/publicips/{publicip_id}/associate-instance"
		product = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	publicIpIid := d.Get("publicip_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{publicip_id}", publicIpIid)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateVpcEipv3BodyParams(d))
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating VPC EIP associate: %s", err)
	}

	d.SetId(publicIpIid)

	return resourceEipv3AssociateRead(ctx, d, meta)
}

func buildCreateVpcEipv3BodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"associate_instance_type": d.Get("associate_instance_type"),
		"associate_instance_id":   d.Get("associate_instance_id"),
	}
	bodyParams := map[string]interface{}{
		"publicip": params,
	}
	return bodyParams
}

func resourceEipv3AssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/eip/publicips/{publicip_id}"
		product = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{publicip_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC EIP associate")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	associateInstanceId := utils.PathSearch("publicip.associate_instance_id", getRespBody, "").(string)
	if associateInstanceId == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving VPC EIP associate")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("publicip_id", utils.PathSearch("publicip.id", getRespBody, nil)),
		d.Set("associate_instance_type", utils.PathSearch("publicip.associate_instance_type", getRespBody, nil)),
		d.Set("associate_instance_id", associateInstanceId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEipv3AssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/eip/publicips/{publicip_id}/disassociate-instance"
		product = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{publicip_id}", d.Get("publicip_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertUndefinedErrInto404Err(err, 409, "error_code", "EIP.7902"),
			"error deleting VPC EIP associate")
	}

	return nil
}
