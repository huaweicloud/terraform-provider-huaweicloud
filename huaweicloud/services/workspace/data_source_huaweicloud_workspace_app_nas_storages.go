package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/persistent-storages
func DataSourceAppNasStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppNasStoragesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the NAS storages are located.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the NAS storage to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the NAS storage to be queried.",
			},
			"storages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the NAS storage.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the NAS storage.",
						},
						"storage_metadata": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_handle": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The storage name.",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The storage type.",
									},
									"export_location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The storage access URL.",
									},
								},
							},
							Description: `The metadata of the corresponding storage.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the NAS storage, in RFC3339 format.",
						},
						"personal_folder_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of the personal folders under this NAS storage.",
						},
						"shared_folder_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of the shared folders under this NAS storage.",
						},
					},
				},
				Description: "All NAS storages that match the filter parameters.",
			},
		},
	}
}

func flattenAppNasStorages(policies []interface{}) []interface{} {
	result := make([]interface{}, 0, len(policies))

	for _, val := range policies {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", val, nil),
			"name":             utils.PathSearch("name", val, nil),
			"storage_metadata": flattenStorageMetadata(utils.PathSearch("storage_metadata", val, make([]interface{}, 0))),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				val, "").(string))/1000, false),
			"personal_folder_count": utils.PathSearch("user_claim_count", val, nil),
			"shared_folder_count":   utils.PathSearch("share_claim_count", val, nil),
		})
	}

	return result
}

func buildAppNasStoragesQueryParams(d *schema.ResourceData) string {
	res := ""
	if storageId, ok := d.GetOk("storage_id"); ok {
		res = fmt.Sprintf("%s&storage_id=%v", res, storageId)
	}
	if storageName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, storageName)
	}
	return res
}

func dataSourceAppNasStoragesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	storages, err := listAppNasStorages(client, buildAppNasStoragesQueryParams(d))
	if err != nil {
		// API error already formated in the list method.
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("storages", flattenAppNasStorages(storages)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting data source fields of the NAS storages: %s", err)
	}
	return nil
}
