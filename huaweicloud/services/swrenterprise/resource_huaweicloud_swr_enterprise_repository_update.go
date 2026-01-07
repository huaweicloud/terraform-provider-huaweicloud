package swrenterprise

import (
	"context"
	"net/url"
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

var enterpriseRepositoryUpdateNonUpdatableParams = []string{
	"instance_id", "namespace_name", "repository_name",
}

// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}
func ResourceSwrEnterpriseRepositoryUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseRepositoryUpdateCreateOrUpdate,
		UpdateContext: resourceSwrEnterpriseRepositoryUpdateCreateOrUpdate,
		ReadContext:   resourceSwrEnterpriseRepositoryUpdateRead,
		DeleteContext: resourceSwrEnterpriseRepositoryUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseRepositoryUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"repository_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repository name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the repository.`,
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

func resourceSwrEnterpriseRepositoryUpdateCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	repositoryName := d.Get("repository_name").(string)

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", namespaceName)
	updatePath = strings.ReplaceAll(updatePath, "{repository_name}", url.PathEscape(strings.ReplaceAll(repositoryName, "/", "%2F")))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"description": d.Get("description"),
		},
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR enterprise repository: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(instanceId + "/" + namespaceName + "/" + repositoryName)
	}

	return resourceSwrEnterpriseRepositoryUpdateRead(ctx, d, meta)
}

func resourceSwrEnterpriseRepositoryUpdateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/repositories/{repository_name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", d.Get("namespace_name").(string))
	getPath = strings.ReplaceAll(getPath, "{repository_name}", url.PathEscape(strings.ReplaceAll(d.Get("repository_name").(string), "/", "%2F")))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR repository")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseRepositoryUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise repository update resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
