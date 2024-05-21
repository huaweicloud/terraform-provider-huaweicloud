package dds

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

const pageLimit = 10

// @API DDS GET /v3/{project_id}/backups
func DataSourceDDSBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsBackupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"end_time"},
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"begin_time"},
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"datastore": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsBackupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getBackupHttpUrl := "v3/{project_id}/backups"
	getBackupPath := client.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", client.ProjectID)
	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pagelimit is `10`
	getBackupPath += fmt.Sprintf("?limit=%v", pageLimit)
	getBackupPath = buildQueryBackupListPath(d, getBackupPath)

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getBackupPath + fmt.Sprintf("&offset=%d", currentTotal)
		getBackupResp, err := client.Request("GET", currentPath, &getBackupOpt)
		if err != nil {
			return diag.Errorf("error retrieving backups: %s", err)
		}
		getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		backups := utils.PathSearch("backups", getBackupRespBody, make([]interface{}, 0)).([]interface{})
		for _, backup := range backups {
			// filter result
			instanceName := utils.PathSearch("instance_name", backup, "").(string)
			backupName := utils.PathSearch("name", backup, "").(string)
			status := utils.PathSearch("status", backup, "").(string)
			description := utils.PathSearch("description", backup, "").(string)
			if val, ok := d.GetOk("instance_name"); ok && instanceName != val {
				continue
			}
			if val, ok := d.GetOk("backup_name"); ok && backupName != val {
				continue
			}
			if val, ok := d.GetOk("status"); ok && status != val {
				continue
			}
			if val, ok := d.GetOk("description"); ok && description != val {
				continue
			}
			results = append(results, map[string]interface{}{
				"id":            utils.PathSearch("id", backup, nil),
				"name":          utils.PathSearch("name", backup, nil),
				"instance_id":   utils.PathSearch("instance_id", backup, nil),
				"instance_name": utils.PathSearch("instance_name", backup, nil),
				"datastore":     flattenGetBackupResponseDatastore(backup),
				"type":          utils.PathSearch("type", backup, nil),
				"begin_time":    utils.PathSearch("begin_time", backup, nil),
				"end_time":      utils.PathSearch("end_time", backup, nil),
				"status":        utils.PathSearch("status", backup, nil),
				"size":          utils.PathSearch("size", backup, 0),
				"description":   utils.PathSearch("description", backup, nil),
			})
		}

		// `totalCount` means the number of all `backups`, and type is float64.
		currentTotal += len(backups)
		totalCount := utils.PathSearch("total_count", getBackupRespBody, float64(0))
		if int(totalCount.(float64)) == currentTotal {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("backups", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryBackupListPath(d *schema.ResourceData, getBackupPath string) string {
	if instId, ok := d.GetOk("instance_id"); ok {
		getBackupPath += fmt.Sprintf("&instance_id=%s", instId)
	}
	if backupId, ok := d.GetOk("backup_id"); ok {
		getBackupPath += fmt.Sprintf("&backup_id=%s", backupId)
	}
	if backupType, ok := d.GetOk("backup_type"); ok {
		getBackupPath += fmt.Sprintf("&backup_type=%s", backupType)
	}
	if mode, ok := d.GetOk("mode"); ok {
		getBackupPath += fmt.Sprintf("&mode=%s", mode)
	}
	if beginTime, ok := d.GetOk("begin_time"); ok {
		getBackupPath += fmt.Sprintf("&begin_time=%s", beginTime)
		getBackupPath += fmt.Sprintf("&end_time=%s", d.Get("end_time"))
	}
	return getBackupPath
}
