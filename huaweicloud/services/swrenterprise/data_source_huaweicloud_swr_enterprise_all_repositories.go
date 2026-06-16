package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/{project_id}/repositories
func DataSourceSwrEnterpriseAllRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseAllRepositoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repositories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates the repository list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the repository ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the repository name.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the namespace name.",
						},
						"namespace_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the namespace ID.",
						},
						"tag_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the number of artifact tags in a repository.",
						},
						"pull_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the total number of downloads.",
						},
						"artifact_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the total number of artifact packages.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the description.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the creation time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the update time.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the ID of an SWR Enterprise Edition instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the name of an SWR Enterprise Edition instance.",
						},
						"resource_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the resource URN.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseAllRepositoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	repositories, err := listAllRepositories(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("repositories", flattenRepositories(repositories)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func listAllRepositories(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/repositories"
		result  = make([]interface{}, 0)
		limit   = 100
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	if val, ok := d.GetOk("name"); ok {
		listPath = fmt.Sprintf("%s&name=%s", listPath, val)
	}

	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		repos := utils.PathSearch("repositories", respBody, make([]interface{}, 0)).([]interface{})

		result = append(result, repos...)

		if len(repos) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func flattenRepositories(repositories []interface{}) []map[string]interface{} {
	if len(repositories) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0)
	for _, repo := range repositories {
		result = append(result, map[string]interface{}{
			"id":             int(utils.PathSearch("id", repo, float64(0)).(float64)),
			"name":           utils.PathSearch("name", repo, nil),
			"namespace_name": utils.PathSearch("namespace_name", repo, nil),
			"namespace_id":   int(utils.PathSearch("namespace_id", repo, float64(0)).(float64)),
			"tag_count":      int(utils.PathSearch("tag_count", repo, float64(0)).(float64)),
			"pull_count":     int(utils.PathSearch("pull_count", repo, float64(0)).(float64)),
			"artifact_count": int(utils.PathSearch("artifact_count", repo, float64(0)).(float64)),
			"description":    utils.PathSearch("description", repo, nil),
			"created_at":     utils.PathSearch("created_at", repo, nil),
			"updated_at":     utils.PathSearch("updated_at", repo, nil),
			"instance_id":    utils.PathSearch("instance_id", repo, nil),
			"instance_name":  utils.PathSearch("instance_name", repo, nil),
			"resource_urn":   utils.PathSearch("resource_urn", repo, nil),
		})
	}
	return result
}
