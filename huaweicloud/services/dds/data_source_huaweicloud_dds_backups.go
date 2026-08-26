package dds

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
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
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/backups?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildBackupsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the DDS backups: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		backups := utils.PathSearch("backups", getRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, backups...)
		if len(backups) < limit {
			break
		}

		offset += len(backups)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("backups", flattenBackups(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackups(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
			"datastore":     flattenGetBackupResponseDatastore(v),
			"type":          utils.PathSearch("type", v, nil),
			"begin_time":    utils.PathSearch("begin_time", v, nil),
			"end_time":      utils.PathSearch("end_time", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"size":          utils.PathSearch("size", v, nil),
			"description":   utils.PathSearch("description", v, nil),
		})
	}

	return result
}

func buildBackupsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if instId, ok := d.GetOk("instance_id"); ok {
		queryParams += fmt.Sprintf("&instance_id=%s", instId)
	}
	if backupId, ok := d.GetOk("backup_id"); ok {
		queryParams += fmt.Sprintf("&backup_id=%s", backupId)
	}
	if backupType, ok := d.GetOk("backup_type"); ok {
		queryParams += fmt.Sprintf("&backup_type=%s", backupType)
	}
	if mode, ok := d.GetOk("mode"); ok {
		queryParams += fmt.Sprintf("&mode=%s", mode)
	}
	if beginTime, ok := d.GetOk("begin_time"); ok {
		queryParams += fmt.Sprintf("&begin_time=%s", beginTime)
	}
	if endTime, ok := d.GetOk("end_time"); ok {
		queryParams += fmt.Sprintf("&end_time=%s", endTime)
	}
	if instanceName, ok := d.GetOk("instance_name"); ok {
		queryParams += fmt.Sprintf("&instance_name=%s", instanceName)
	}
	if backupName, ok := d.GetOk("backup_name"); ok {
		queryParams += fmt.Sprintf("&backup_name=%s", backupName)
	}
	if backupStatus, ok := d.GetOk("status"); ok {
		queryParams += fmt.Sprintf("&backup_status=%s", backupStatus)
	}
	if description, ok := d.GetOk("description"); ok {
		queryParams += fmt.Sprintf("&backup_description=%s", description)
	}

	return queryParams
}
