package mrs

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

// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/files
func DataSourceClusterFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterFilesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the cluster files are located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The directory path of the file.`,
			},
			"files": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of files under the specified path.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path_suffix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The path suffix of the under the queried directory.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the file.`,
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The group to which the file belongs.`,
						},
						"permission": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The permission of the file.`,
						},
						"replication": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The replication factor of the file.`,
						},
						"block_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The block size of the file.`,
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The length of the file.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the file.`,
						},
						"children_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of entries under this directory.`,
						},
						"access_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the file was accessed, in RFC3339 format.`,
						},
						"modification_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the file was modified, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func listClusterFiles(client *golangsdk.ServiceClient, clusterId, directory string) ([]interface{}, error) {
	var (
		httpURL = "v2/{project_id}/clusters/{cluster_id}/files"
		result  = make([]interface{}, 0)
		limit   = 100
		// The offset indicates the number of pages, starts from 1.
		offset = 1
	)

	listPath := client.Endpoint + httpURL
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	listPath = fmt.Sprintf("%s?limit=%d&path=%s", listPath, limit, directory)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		files := utils.PathSearch("files", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, files...)
		if len(files) < limit {
			break
		}

		offset++
	}

	return result, nil
}

func dataSourceClusterFilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	directory := d.Get("path").(string)
	clusterId := d.Get("cluster_id").(string)
	files, err := listClusterFiles(client, clusterId, directory)
	if err != nil {
		return diag.Errorf("error retrieving cluster (%s) file list under the specified path (%s): %s",
			clusterId, directory, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("files", flattenClusterFiles(files)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterFiles(files []interface{}) []map[string]interface{} {
	if len(files) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(files))
	for _, v := range files {
		rst = append(rst, map[string]interface{}{
			"path_suffix":  utils.PathSearch("path_suffix", v, nil),
			"owner":        utils.PathSearch("owner", v, nil),
			"group":        utils.PathSearch("group", v, nil),
			"permission":   utils.PathSearch("permission", v, nil),
			"replication":  utils.PathSearch("replication", v, nil),
			"block_size":   utils.PathSearch("block_size", v, nil),
			"length":       utils.PathSearch("length", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"children_num": utils.PathSearch("children_num", v, nil),
			"access_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("access_time",
				v, float64(0)).(float64))/1000, false),
			"modification_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("modification_time",
				v, float64(0)).(float64))/1000, false),
		})
	}

	return rst
}
