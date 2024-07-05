package swr

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/shared-repositories
func DataSourceSharedRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"center": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_name": {
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
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeBool,
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

func dataSourceSharedRepositoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	// listSharedRepositories: Query the list of SWR shared repositories.
	resp, err := getSharedRepositories(client, d)
	if err != nil {
		return diag.Errorf("error retrieving SWR shared repositories: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("repositories", flattenSharedRepositoriesResponse(resp, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSharedRepositories(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	listSharedRepositoriesHttpUrl := "v2/manage/shared-repositories"
	listSharedRepositoriesPath := client.Endpoint + listSharedRepositoriesHttpUrl
	listSharedRepositoriesQueryParams := buildListSharedRepositoriesQueryParams(d)
	listSharedRepositoriesPath += listSharedRepositoriesQueryParams

	listSharedRepositoriesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	results := make([]interface{}, 0)
	offset := 0
	for {
		currentPath := fmt.Sprintf("%slimit=10&offset=%d", listSharedRepositoriesPath, offset)
		listSharedRepositoriesResp, err := client.Request("GET", currentPath, &listSharedRepositoriesOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying SWR shared repositories: %s", err)
		}
		listSharedRepositoriesRespBody, err := utils.FlattenResponse(listSharedRepositoriesResp)
		if err != nil {
			return nil, fmt.Errorf("error retrieving SWR shared repositories: %s", err)
		}
		sharedRepositories := listSharedRepositoriesRespBody.([]interface{})
		total := 0
		if len(sharedRepositories) > 0 {
			total = int(utils.PathSearch("total_range", sharedRepositories[0], float64(0)).(float64))
		}
		results = append(results, sharedRepositories...)
		offset += len(sharedRepositories)
		if offset == total {
			break
		}
	}

	return results, nil
}

func buildListSharedRepositoriesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("organization"); ok {
		res = fmt.Sprintf("%s&namespace=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("center"); ok {
		res = fmt.Sprintf("%s&center=%v", res, v)
	}
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}

	if res != "" {
		return fmt.Sprintf("?%s&", res[1:])
	}
	return "?"
}

func flattenSharedRepositoriesResponse(rawParams []interface{}, d *schema.ResourceData) []interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	results := make([]interface{}, 0)
	for _, v := range rawParams {
		name := utils.PathSearch("name", v, "").(string)
		if val, ok := d.GetOk("name"); ok && name != val {
			continue
		}
		results = append(results, map[string]interface{}{
			"organization":  utils.PathSearch("namespace", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"domain_name":   utils.PathSearch("domain_name", v, nil),
			"status":        utils.PathSearch("status", v, false),
			"category":      utils.PathSearch("category", v, nil),
			"is_public":     utils.PathSearch("is_public", v, false),
			"description":   utils.PathSearch("description", v, nil),
			"size":          int(utils.PathSearch("size", v, float64(0)).(float64)),
			"num_images":    int(utils.PathSearch("num_images", v, float64(0)).(float64)),
			"num_download":  int(utils.PathSearch("num_download", v, float64(0)).(float64)),
			"path":          utils.PathSearch("path", v, nil),
			"internal_path": utils.PathSearch("internal_path", v, nil),
			"tags":          utils.PathSearch("tags", v, nil),
			"created_at":    utils.PathSearch("created_at", v, nil),
			"updated_at":    utils.PathSearch("updated_at", v, nil),
		})
	}
	return results
}
