package swr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v3/manage/shared-repositories
func DataSourceSwrv3SharedRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrv3SharedRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"shared_by": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the sharing mode.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the organization name.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the image repository name.`,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether the sharing has expired.`,
			},
			"repos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the repositories.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the repository ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the repository name.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the repository size.`,
						},
						"organization": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the organization that a repository belongs to.`,
						},
						"num_download": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the repository downloads.`,
						},
						"status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the image shared by others has expired.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the repository type.`,
						},
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether a repository is a public repository.`,
						},
						"num_images": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of images in a repository.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the repository description.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when a repository was created. It is the UTC standard time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when a repository was updated. It is the UTC standard time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrv3SharedRepositoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listRepositoriesHttpUrl := "v3/manage/shared-repositories"
	listRepositoriesPath := client.Endpoint + listRepositoriesHttpUrl
	listRepositoriesPath += buildV3ListSharedRepositoriesQueryParams(d)
	listRepositoriesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	marker := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listRepositoriesPath
		if marker != 0 {
			currentPath += fmt.Sprintf("&marker=%v", marker)
		}
		listRepositoriesResp, err := client.Request("GET", currentPath, &listRepositoriesOpt)
		if err != nil {
			return diag.Errorf("error querying SWR repositories: %s", err)
		}
		listRepositoriesRespBody, err := utils.FlattenResponse(listRepositoriesResp)
		if err != nil {
			return diag.Errorf("error flattening SWR repositories response : %s", err)
		}

		repositories := utils.PathSearch("repos", listRepositoriesRespBody, make([]interface{}, 0)).([]interface{})
		for _, repository := range repositories {
			results = append(results, map[string]interface{}{
				"id":           utils.PathSearch("id", repository, nil),
				"name":         utils.PathSearch("name", repository, nil),
				"size":         utils.PathSearch("size", repository, nil),
				"organization": utils.PathSearch("namespace_name", repository, nil),
				"num_download": utils.PathSearch("num_download", repository, nil),
				"status":       utils.PathSearch("status", repository, nil),
				"category":     utils.PathSearch("category", repository, nil),
				"is_public":    utils.PathSearch("is_public", repository, nil),
				"description":  utils.PathSearch("description", repository, nil),
				"num_images":   utils.PathSearch("num_images", repository, nil),
				"created_at":   utils.PathSearch("created_at", repository, nil),
				"updated_at":   utils.PathSearch("updated_at", repository, nil),
			})
		}

		marker = int(utils.PathSearch("nextMarker", listRepositoriesRespBody, float64(-1)).(float64))
		if marker == -1 {
			break
		}
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("repos", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV3ListSharedRepositoriesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"
	res = fmt.Sprintf("%s&shared_by=%v", res, d.Get("shared_by"))
	if v, ok := d.GetOk("organization"); ok {
		res = fmt.Sprintf("%s&namespace=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	return res
}
