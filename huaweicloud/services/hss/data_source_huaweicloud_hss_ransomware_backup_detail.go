package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/ransomware/backup/{backup_id}/detail
func DataSourceRansomwareBackupDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRansomwareBackupDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// `created_at` in API documentation as a Long type, but it actually returns a string type.
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRansomwareBackupDetailQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceRansomwareBackupDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		backupId = d.Get("backup_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v5/{project_id}/ransomware/backup/{backup_id}/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{backup_id}", backupId)
	requestPath += buildRansomwareBackupDetailQueryParams(epsId)

	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS ransomware backup detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error retrieving HSS ransomware backup detail: ID is not found in API response")
	}

	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("image_type", utils.PathSearch("image_type", respBody, nil)),
		d.Set("vault_id", utils.PathSearch("vault_id", respBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("resource_size", utils.PathSearch("resource_size", respBody, nil)),
		d.Set("resource_id", utils.PathSearch("resource_id", respBody, nil)),
		d.Set("resource_type", utils.PathSearch("resource_type", respBody, nil)),
		d.Set("resource_name", utils.PathSearch("resource_name", respBody, nil)),
		d.Set("children", flattenRansomwareBackupDetailChildren(
			utils.PathSearch("children", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRansomwareBackupDetailChildren(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"image_type":    utils.PathSearch("image_type", v, nil),
			"vault_id":      utils.PathSearch("vault_id", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"resource_size": utils.PathSearch("resource_size", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
		})
	}

	return rst
}
