package dds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/errorlog-download
func DataSourceDDSErrorLogLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSErrorLogLinksRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the instance.`,
			},
			"file_name_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the names of the files.`,
			},
			"node_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the node IDs to which the files belong.`,
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the error logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the file name.`,
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the link status.`,
						},
						"file_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error file size.`,
						},
						"file_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the file link.`,
						},
						"updated_at": {
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

func dataSourceDDSErrorLogLinksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	resp, err := waitForErrorLogLinkCompleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error retrieving error log links: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("links", flattenErrorLogLinks(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func waitForErrorLogLinkCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"FINISH"},
		Refresh:      errorLogLinksRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutRead),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for error log links to be finished: %s", err)
	}

	return res, nil
}

func errorLogLinksRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getErrorLogLinks(client, d)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status", resp, "").(string)
		if status == "" {
			return nil, "ERROR", fmt.Errorf("unable to find status from API response")
		}

		return resp, status, nil
	}
}

func getErrorLogLinks(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/errorlog-download"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"file_name_list": utils.ValueIgnoreEmpty(d.Get("file_name_list")),
			"node_id_list":   utils.ValueIgnoreEmpty(d.Get("node_id_list")),
		}),
	}

	resp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting error log links: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error flattening response: %s", err)
	}

	return respBody, nil
}

func flattenErrorLogLinks(resp interface{}) []map[string]interface{} {
	links := utils.PathSearch("list", resp, make([]interface{}, 0)).([]interface{})
	if len(links) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		result = append(result, map[string]interface{}{
			"file_name": utils.PathSearch("file_name", link, nil),
			"node_name": utils.PathSearch("node_name", link, nil),
			"status":    utils.PathSearch("status", link, nil),
			"file_size": utils.PathSearch("file_size", link, nil),
			"file_link": utils.PathSearch("file_link", link, nil),
			"updated_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("updated_at", link, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
