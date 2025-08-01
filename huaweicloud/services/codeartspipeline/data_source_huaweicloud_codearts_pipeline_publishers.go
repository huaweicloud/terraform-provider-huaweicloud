package codeartspipeline

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

// @API CodeArtsPipeline GET /v1/{domain_id}/publisher/query-all
func DataSourceCodeArtsPipelinePublishers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelinePublishersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the publisher name.`,
			},
			"publishers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the publisher list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the publisher ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the publisher name.`,
						},
						"en_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the publisher English name.`,
						},
						"support_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the support URL.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"logo_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the logo URL.`,
						},
						"website": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the website URL.`,
						},
						"source_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the source URL.`,
						},
						"auth_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the authorization status.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"last_update_user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the updater ID.`,
						},
						"last_update_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the updater name.`,
						},
						"last_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelinePublishersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	publishers, err := getPipelinePublishers(client, cfg.DomainID, d.Get("name").(string))
	if err != nil {
		return diag.Errorf("error getting CodeArts Pipeline publishers: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("publishers", publishers),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPipelinePublishers(client *golangsdk.ServiceClient, domainId, name string) ([]interface{}, error) {
	getHttpUrl := "v1/{domain_id}/publisher/query-all?limit=10"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if name != "" {
		getPath += fmt.Sprintf("&name=%s", name)
	}

	offset := 0
	rst := make([]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, fmt.Errorf("error flatten response: %s", err)
		}

		publishers := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(publishers) == 0 {
			break
		}

		for _, publisher := range publishers {
			rst = append(rst, map[string]interface{}{
				"id":                    utils.PathSearch("publisher_unique_id", publisher, nil),
				"name":                  utils.PathSearch("name", publisher, nil),
				"en_name":               utils.PathSearch("en_name", publisher, nil),
				"support_url":           utils.PathSearch("support_url", publisher, nil),
				"description":           utils.PathSearch("description", publisher, nil),
				"logo_url":              utils.PathSearch("logo_url", publisher, nil),
				"website":               utils.PathSearch("website", publisher, nil),
				"source_url":            utils.PathSearch("source_url", publisher, nil),
				"auth_status":           utils.PathSearch("auth_status", publisher, nil),
				"user_id":               utils.PathSearch("user_id", publisher, nil),
				"last_update_user_id":   utils.PathSearch("last_update_user_id", publisher, nil),
				"last_update_user_name": utils.PathSearch("last_update_user_name", publisher, nil),
				"last_update_time":      utils.PathSearch("last_update_time", publisher, nil),
			})
		}

		offset += 10
	}

	return rst, nil
}
