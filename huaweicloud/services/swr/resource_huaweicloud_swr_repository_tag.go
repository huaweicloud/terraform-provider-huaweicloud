package swr

import (
	"context"
	"errors"
	"fmt"
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

var repositoryTagNonUpdatableParams = []string{
	"organization", "repository", "name", "source_tag", "destination_tag", "override",
}

// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/tags
// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}
// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}
func ResourceSwrRepositoryTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrRepositoryTagCreate,
		UpdateContext: resourceSwrRepositoryTagUpdate,
		ReadContext:   resourceSwrRepositoryTagRead,
		DeleteContext: resourceSwrRepositoryTagDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrRepositoryTagImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(repositoryTagNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the organization.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the repository.`,
			},
			"source_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the source tag.`,
			},
			"destination_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the destination tag.`,
			},
			"override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to override.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"tag_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the tag ID.`,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the repository ID.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image tag.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image ID.`,
			},
			"manifest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image manifest.`,
			},
			"digest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image digest.`,
			},
			"schema": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the docker schema.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image pull path.`,
			},
			"internal_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the image internal pull path.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the image size.`,
			},
			"is_trusted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the image is trusted.`,
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"tag_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the tag type.`,
			},
		},
	}
}

func resourceSwrRepositoryTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)

	createHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/tags"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{namespace}", organization)
	createPath = strings.ReplaceAll(createPath, "{repository}", strings.ReplaceAll(repository, "/", "$"))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrRepositoryTagBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR repository tag: %s", err)
	}

	tag := d.Get("destination_tag").(string)
	d.SetId(organization + "/" + repository + "/" + tag)
	d.Set("tag", tag)

	return resourceSwrRepositoryTagRead(ctx, d, meta)
}

func buildCreateSwrRepositoryTagBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_tag":      d.Get("source_tag"),
		"destination_tag": d.Get("destination_tag"),
		"override":        d.Get("override"),
	}

	return bodyParams
}

func resourceSwrRepositoryTagRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	tag := d.Get("tag").(string)

	getHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{namespace}", organization)
	getPath = strings.ReplaceAll(getPath, "{repository}", strings.ReplaceAll(repository, "/", "$"))
	getPath = strings.ReplaceAll(getPath, "{tag}", tag)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR repository tag")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tag_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("repository_id", utils.PathSearch("repo_id", getRespBody, nil)),
		d.Set("image_id", utils.PathSearch("image_id", getRespBody, nil)),
		d.Set("manifest", utils.PathSearch("manifest", getRespBody, nil)),
		d.Set("digest", utils.PathSearch("digest", getRespBody, nil)),
		d.Set("schema", utils.PathSearch("schema", getRespBody, nil)),
		d.Set("path", utils.PathSearch("path", getRespBody, nil)),
		d.Set("internal_path", utils.PathSearch("internal_path", getRespBody, nil)),
		d.Set("size", utils.PathSearch("size", getRespBody, nil)),
		d.Set("is_trusted", utils.PathSearch("is_trusted", getRespBody, nil)),
		d.Set("created", utils.PathSearch("created", getRespBody, nil)),
		d.Set("updated", utils.PathSearch("updated", getRespBody, nil)),
		d.Set("tag_type", utils.PathSearch("tag_type", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrRepositoryTagUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrRepositoryTagDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	tag := d.Get("tag").(string)

	deleteHttpUrl := "v2/manage/namespaces/{namespace}/repos/{repository}/tags/{tag}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{namespace}", organization)
	deletePath = strings.ReplaceAll(deletePath, "{repository}", strings.ReplaceAll(repository, "/", "$"))
	deletePath = strings.ReplaceAll(deletePath, "{tag}", tag)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR repository tag")
	}

	return nil
}

func resourceSwrRepositoryTagImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 3 {
			return nil, errors.New("invalid id format, must be <organization_name>/<repository_name>/<tag> or " +
				"<organization_name>,<repository_name>,<tag>")
		}
	} else {
		// reform ID to be separated by slashes
		id := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		d.SetId(id)
	}

	d.Set("organization", parts[0])
	d.Set("repository", parts[1])
	d.Set("tag", parts[2])

	return []*schema.ResourceData{d}, nil
}
