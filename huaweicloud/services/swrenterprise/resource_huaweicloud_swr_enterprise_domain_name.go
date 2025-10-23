package swrenterprise

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

var enterpriseDomainNameNonUpdatableParams = []string{
	"instance_id", "domain_name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/domainname
// @API SWR GET /v2/{project_id}/instances/{instance_id}/domainname
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/domainname/{domainname_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/domainname/{domainname_id}
func ResourceSwrEnterpriseDomainName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseDomainNameCreate,
		UpdateContext: resourceSwrEnterpriseDomainNameUpdate,
		ReadContext:   resourceSwrEnterpriseDomainNameRead,
		DeleteContext: resourceSwrEnterpriseDomainNameDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseDomainNameNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name.`,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SCM certificate ID.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the domain name type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"domain_name_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the domain name ID.`,
			},
		},
	}
}

func resourceSwrEnterpriseDomainNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/domainname"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrEnterpriseDomainNameBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR domain name: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("domain_name_info.uid", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find SWR instance domain name ID from the API response")
	}

	d.SetId(instanceId + "/" + id)

	return resourceSwrEnterpriseDomainNameRead(ctx, d, meta)
}

func buildCreateSwrEnterpriseDomainNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_name":    d.Get("domain_name"),
		"certificate_id": d.Get("certificate_id"),
	}

	return bodyParams
}

func resourceSwrEnterpriseDomainNameRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<domain_name_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	id := parts[1]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/domainname?uid={domain_name_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{domain_name_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting SWR domain name")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening SWR instance domain name response: %s", err)
	}

	domainName := utils.PathSearch("domain_name_infos[0]", getRespBody, nil)
	if domainName == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting SWR domain name")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("domain_name_id", id),
		d.Set("domain_name", utils.PathSearch("domain_name", domainName, nil)),
		d.Set("type", utils.PathSearch("type", domainName, nil)),
		d.Set("certificate_id", utils.PathSearch("certificate_id", domainName, nil)),
		d.Set("created_at", utils.PathSearch("created_at", domainName, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", domainName, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseDomainNameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	if d.HasChange("certificate_id") {
		updateHttpUrl := "v2/{project_id}/instances/{instance_id}/domainname/{domainname_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{domainname_id}", d.Get("domain_name_id").(string))
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateSwrEnterpriseDomainNameBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating SWR instance domain name: %s", err)
		}
	}

	return resourceSwrEnterpriseDomainNameRead(ctx, d, meta)
}

func buildUpdateSwrEnterpriseDomainNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"certificate_id": d.Get("certificate_id"),
	}

	return bodyParams
}

func resourceSwrEnterpriseDomainNameDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/domainname/{domainname_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{domainname_id}", d.Get("domain_name_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR domain name")
	}

	return nil
}
