package modelarts

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/notebooks
func DataSourceNotebooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotebooksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the notebooks are located.`,
			},

			// Optional parameters.
			"feature": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The feature type of the notebooks to be queried.`,
			},
			"notebook_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The notebook instance ID to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the notebooks to be queried. Fuzzy match is supported.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the notebooks to be queried.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID of the notebooks to be queried.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The flavor of the notebooks to be queried.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The image ID of the notebooks to be queried.`,
			},
			"billing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The billing type of the notebooks to be queried.`,
			},

			// Attributes.
			"notebooks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceNotebooksNotebookElemSchema(),
				Description: `The list of notebooks that match the filter parameters.`,
			},
		},
	}
}

func dataSourceNotebooksNotebookElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the notebook instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the notebook instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the notebook instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the notebook instance.`,
			},
			"feature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The feature type of the notebook instance.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor of the notebook instance.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace ID to which the notebook belongs.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dedicated resource pool ID of the notebook instance.`,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The dedicated resource pool name of the notebook instance.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image ID of the notebook instance.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image name of the notebook instance.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The image type of the notebook instance.`,
			},
			"image_swr_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SWR path of the image of the notebook instance.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access URL of the notebook instance.`,
			},
			"fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The failure reason of the notebook instance.`,
			},
			"auto_stop_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the auto stop feature is enabled for the notebook instance.`,
			},
			"lease_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The lease duration of the notebook instance, in milliseconds.`,
			},
			"lease_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The auto stop type of the notebook instance.`,
			},
			"key_pair": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH key pair name when SSH access is configured.`,
			},
			"ssh_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH access URI when SSH access is configured.`,
			},
			"allowed_access_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The allowed access IP list for SSH when SSH access is configured.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID of the notebook instance.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of the node where the notebook instance is located.`,
			},
			"jupyter_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The JupyterLab version of the notebook instance.`,
			},
			"billing_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The billing resource types of the notebook instance.`,
			},
			"custom_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom specification of the notebook instance, in JSON format.`,
			},
			"user_vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user VPC configuration of the notebook instance, in JSON format.`,
			},
			"volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage category of the notebook instance.`,
						},
						"ownership": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ownership type of the storage.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The storage capacity.`,
						},
						"uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage URI.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage ID.`,
						},
						"mount_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mount path of the storage.`,
						},
					},
				},
				Description: `The system volume configuration of the notebook instance.`,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The extended storage category.`,
						},
						"mount_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mount path of the extended storage.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source path of the extended storage.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the extended storage.`,
						},
						"mount_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mount type of the extended storage.`,
						},
					},
				},
				Description: `The extended storage list of the notebook instance.`,
			},
			"tags": common.TagsComputedSchema(`The tags of the notebook instance.`),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the notebook instance.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the notebook instance.`,
			},
		},
	}
}

func buildNotebooksQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("billing"); ok {
		res = fmt.Sprintf("%s&billing=%v", res, v)
	}
	if v, ok := d.GetOk("feature"); ok {
		res = fmt.Sprintf("%s&feature=%v", res, v)
	}
	if v, ok := d.GetOk("flavor_id"); ok {
		res = fmt.Sprintf("%s&flavor=%v", res, v)
	}
	if v, ok := d.GetOk("image_id"); ok {
		res = fmt.Sprintf("%s&image_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("notebook_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	return res
}

func listNotebooks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	// The ListNotebooks API allows 10, 20 or 50 for limit; use the maximum to reduce requests.
	var (
		httpUrl = "v1/{project_id}/notebooks?limit={limit}"
		limit   = 50
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildNotebooksQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		notebooks := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, notebooks...)

		if len(notebooks) < limit {
			break
		}
		offset += len(notebooks)
	}

	return result, nil
}

func flattenNotebookEndpointsAccess(endpoints []interface{}) (keyPair, sshURI string, ips []interface{}) {
	for _, ep := range endpoints {
		if fmt.Sprint(utils.PathSearch("service", ep, "")) != "SSH" {
			continue
		}
		keyPair = utils.PathSearch("key_pair_names[0]", ep, "").(string)
		sshURI = fmt.Sprint(utils.PathSearch("uri", ep, ""))
		ips = utils.PathSearch("allowed_access_ips", ep, make([]interface{}, 0)).([]interface{})
		break
	}
	return
}

func flattenNotebookVolumeForList(volume map[string]interface{}) []map[string]interface{} {
	if len(volume) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":       utils.PathSearch("category", volume, nil),
			"ownership":  utils.PathSearch("ownership", volume, nil),
			"size":       utils.PathSearch("capacity", volume, nil),
			"uri":        utils.PathSearch("uri", volume, nil),
			"id":         utils.PathSearch("id", volume, nil),
			"mount_path": utils.PathSearch("mount_path", volume, nil),
		},
	}
}

func flattenNotebookDataVolumes(dataVolumes []interface{}) []map[string]interface{} {
	if len(dataVolumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataVolumes))
	for _, dataVolume := range dataVolumes {
		result = append(result, map[string]interface{}{
			"type":       utils.PathSearch("category", dataVolume, nil),
			"mount_path": utils.PathSearch("mount_path", dataVolume, nil),
			"path":       utils.PathSearch("url", dataVolume, nil),
			"status":     utils.PathSearch("status", dataVolume, nil),
			"mount_type": utils.PathSearch("mount_type", dataVolume, nil),
		})
	}
	return result
}

func flattenNotebooks(notebooks []interface{}) []map[string]interface{} {
	if len(notebooks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(notebooks))
	for _, notebook := range notebooks {
		keyPair, sshURI, ips := flattenNotebookEndpointsAccess(utils.PathSearch("endpoints", notebook, make([]interface{}, 0)).([]interface{}))

		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", notebook, nil),
			"name":               utils.PathSearch("name", notebook, nil),
			"description":        utils.PathSearch("description", notebook, nil),
			"status":             utils.PathSearch("status", notebook, nil),
			"feature":            utils.PathSearch("feature", notebook, nil),
			"flavor_id":          utils.PathSearch("flavor", notebook, nil),
			"workspace_id":       utils.PathSearch("workspace_id", notebook, nil),
			"pool_id":            utils.PathSearch("pool.id", notebook, nil),
			"pool_name":          utils.PathSearch("pool.name", notebook, nil),
			"image_id":           utils.PathSearch("image.id", notebook, nil),
			"image_name":         utils.PathSearch("image.name", notebook, nil),
			"image_type":         utils.PathSearch("image.type", notebook, nil),
			"image_swr_path":     utils.PathSearch("image.swr_path", notebook, nil),
			"url":                utils.PathSearch("url", notebook, nil),
			"fail_reason":        utils.PathSearch("fail_reason", notebook, nil),
			"auto_stop_enabled":  utils.PathSearch("lease.enable", notebook, nil),
			"lease_duration":     utils.PathSearch("lease.duration", notebook, nil),
			"lease_type":         utils.PathSearch("lease.type", notebook, nil),
			"key_pair":           keyPair,
			"ssh_uri":            sshURI,
			"allowed_access_ips": ips,
			"user_id":            utils.PathSearch("user_id", notebook, nil),
			"ip":                 utils.PathSearch("ip", notebook, nil),
			"jupyter_version":    utils.PathSearch("jupyter_version", notebook, nil),
			"billing_items":      utils.PathSearch("billing_items", notebook, make([]interface{}, 0)).([]interface{}),
			"custom_spec":        utils.JsonToString(utils.PathSearch("custom_spec", notebook, nil)),
			"user_vpc":           utils.JsonToString(utils.PathSearch("user_vpc", notebook, nil)),
			"volume": flattenNotebookVolumeForList(utils.PathSearch("volume",
				notebook, make(map[string]interface{})).(map[string]interface{})),
			"data_volumes": flattenNotebookDataVolumes(utils.PathSearch("data_volumes", notebook, make([]interface{}, 0)).([]interface{})),
			"tags":         utils.FlattenTagsToMap(utils.PathSearch("tags", notebook, nil)),
			"created_at":   utils.FormatTimeStampRFC3339(int64(utils.PathSearch("lease.create_at", notebook, float64(0)).(float64))/1000, false),
			"updated_at":   utils.FormatTimeStampRFC3339(int64(utils.PathSearch("lease.update_at", notebook, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func dataSourceNotebooksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	notebooks, err := listNotebooks(client, d)
	if err != nil {
		return diag.Errorf("error querying ModelArts notebooks: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("notebooks", flattenNotebooks(notebooks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
