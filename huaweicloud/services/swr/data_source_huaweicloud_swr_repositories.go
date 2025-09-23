// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SWR
// ---------------------------------------------------------------

package swr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/repos
func DataSourceRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_public": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_images": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_download": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"total_range": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listRepositories: Query the list of SWR repositories.
	var (
		listRepositoriesHttpUrl = "v2/manage/repos"
		listRepositoriesProduct = "swr"
	)

	listRepositoriesClient, err := cfg.NewServiceClient(listRepositoriesProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	listRepositoriesPath := listRepositoriesClient.Endpoint + listRepositoriesHttpUrl

	offset := 0
	listRepositoriesQueryParams := buildListPublicRepositoriesQueryParams(d, offset)
	listRepositoriesPath += listRepositoriesQueryParams

	listRepositoriesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	results := make([]map[string]interface{}, 0)
	for {
		listRepositoriesResp, err := listRepositoriesClient.Request("GET", listRepositoriesPath, &listRepositoriesOpt)
		if err != nil {
			return diag.Errorf("error querying SWR repositories: %s", err)
		}
		listRepositoriesRespBody, err := utils.FlattenResponse(listRepositoriesResp)
		if err != nil {
			return diag.Errorf("error retrieving SWR repositories: %s", err)
		}
		repositories := listRepositoriesRespBody.([]interface{})
		total := 0
		if len(repositories) > 0 {
			total = int(utils.PathSearch("total_range", repositories[0], float64(0)).(float64))
		}
		for _, repository := range repositories {
			// filter result by name
			name := utils.PathSearch("name", repository, "").(string)
			if val, ok := d.GetOk("name"); ok && name != val {
				continue
			}
			results = append(results, map[string]interface{}{
				"organization":  utils.PathSearch("namespace", repository, nil),
				"name":          utils.PathSearch("name", repository, nil),
				"category":      utils.PathSearch("category", repository, nil),
				"is_public":     utils.PathSearch("is_public", repository, false),
				"description":   utils.PathSearch("description", repository, nil),
				"size":          int(utils.PathSearch("size", repository, float64(0)).(float64)),
				"num_images":    int(utils.PathSearch("num_images", repository, float64(0)).(float64)),
				"num_download":  int(utils.PathSearch("num_download", repository, float64(0)).(float64)),
				"path":          utils.PathSearch("path", repository, nil),
				"internal_path": utils.PathSearch("internal_path", repository, nil),
				"tags":          utils.PathSearch("tags", repository, nil),
				"status":        utils.PathSearch("status", repository, false),
				"total_range":   int(utils.PathSearch("total_range", repository, float64(0)).(float64)),
				"created_at":    utils.PathSearch("created_at", repository, nil),
				"updated_at":    utils.PathSearch("updated_at", repository, nil),
			})
		}
		offset += len(repositories)
		if offset == total {
			break
		}
		index := strings.Index(listRepositoriesPath, "offset")
		listRepositoriesPath = fmt.Sprintf("%soffset=%v", listRepositoriesPath[:index], offset)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("repositories", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListPublicRepositoriesQueryParams(d *schema.ResourceData, offset int) string {
	res := ""

	if v, ok := d.GetOk("organization"); ok {
		res = fmt.Sprintf("%s&namespace=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("category"); ok {
		res = fmt.Sprintf("%s&category=%v", res, v)
	}
	if v, ok := d.GetOk("is_public"); ok {
		res = fmt.Sprintf("%s&is_public=%v", res, v)
	}

	if res != "" {
		return "?" + res[1:] + fmt.Sprintf("&limit=10&offset=%v", offset)
	}
	return fmt.Sprintf("?limit=10&offset=%v", offset)
}
