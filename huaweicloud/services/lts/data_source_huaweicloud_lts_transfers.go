package lts

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

// @API LTS GET /v2/{project_id}/transfers
func DataSourceLtsTransfers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsTransfersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query log transfers.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log group to which the log transfers and log streams belong.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log stream to be transferred in the log transfer.`,
			},
			"transfers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the transfer.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group to which the log transfer belongs.`,
						},
						"log_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group to which the log transfer belongs.`,
						},
						"log_streams": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataTransfersLogStreamsSchema(),
							Description: `The configuration of the log streams that to be transferred.`,
						},
						"log_transfer_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataTransfersLogTransferInfoSchema(),
							Description: `The configuration of the log transfer.`,
						},
					},
				},
				Description: `The list of log transfers.`,
			},
		},
	}
}

func dataTransfersLogStreamsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the log stream.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the log stream.`,
			},
		},
	}
	return &sc
}

func dataTransfersLogTransferInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"log_transfer_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the log transfer.`,
			},
			"log_transfer_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The mode of the log transfer.`,
			},
			"log_storage_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The format of the log transfer.`,
			},
			"log_transfer_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the log transfer.`,
			},
			"log_agency_transfer": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataTransfersLogAgencySchema(),
				Description: `The configuration of the agency transfer.`,
			},
			"log_transfer_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataTransfersLogDetailSchema(),
				Description: `The detail of the log transfer configuration.`,
			},
		},
	}
	return &sc
}

func dataTransfersLogAgencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"agency_domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the delegator account.`,
			},
			"agency_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the delegator account.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The agency name created by the delegator account.`,
			},
			"agency_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID of the delegator account.`,
			},
		},
	}
	return &sc
}

func dataTransfersLogDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"obs_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The length of the transfer interval for an OBS transfer task.`,
			},
			"obs_period_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit of the transfer interval for an OBS transfer task.`,
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the OBS bucket, which is the log transfer destination object.`,
			},
			"obs_transfer_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The storage path of the OBS bucket, which is the log transfer destination.`,
			},
			"obs_dir_prefix_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom prefix of the transfer path.`,
			},
			"obs_prefix_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The transfer file prefix of an OBS transfer task.`,
			},
			"obs_eps_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID of an OBS transfer task.`,
			},
			"obs_encrypted_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether OBS bucket encryption is enabled.`,
			},
			"obs_encrypted_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The KMS key ID for an OBS transfer task.`,
			},
			"obs_time_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time zone for an OBS transfer task.`,
			},
			"obs_time_zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `ID of the time zone for an OBS transfer task.`,
			},
			"dis_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the DIS stream.`,
			},
			"dis_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the DIS stream.`,
			},
			"kafka_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the kafka instance.`,
			},
			"kafka_topic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kafka topic.`,
			},
			"delivery_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of tag fields will be delivered when transferring.`,
			},
		},
	}
	return &sc
}

// The filter parameter 'log_transfer_type' is unavailable.
func buildQueryTransfersBodyParams(d *schema.ResourceData) string {
	res := ""
	if logGroupName, ok := d.GetOk("log_group_name"); ok {
		return fmt.Sprintf("%s&log_group_name=%v", res, logGroupName)
	}
	if logStreamName, ok := d.GetOk("log_stream_name"); ok {
		return fmt.Sprintf("%s&log_stream_name=%v", res, logStreamName)
	}
	return res
}

func dataSourceLtsTransfersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/transfers?limit=100"
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildQueryTransfersBodyParams(d)

	getTransferOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// The query parameter 'offset' is unavailable
	requestResp, err := client.Request("GET", listPath, &getTransferOpt)
	if err != nil {
		return diag.Errorf("error querying log transfers: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	transfers := utils.PathSearch("log_transfers", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("transfers", flattenTransfers(transfers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTransfers(transfers []interface{}) []interface{} {
	result := make([]interface{}, 0, len(transfers))
	for _, v := range transfers {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("log_transfer_id", v, nil),
			"log_group_id":      utils.PathSearch("log_group_id", v, nil),
			"log_group_name":    utils.PathSearch("log_group_name", v, nil),
			"log_streams":       flattenTransfersElemLogStreams(utils.PathSearch("log_streams", v, make([]interface{}, 0)).([]interface{})),
			"log_transfer_info": flattenTransfersElemLogTransferInfo(utils.PathSearch("log_transfer_info", v, nil)),
		})
	}
	return result
}

func flattenTransfersElemLogStreams(logStreams []interface{}) []interface{} {
	result := make([]interface{}, 0, len(logStreams))
	for _, v := range logStreams {
		result = append(result, map[string]interface{}{
			"log_stream_id":   utils.PathSearch("log_stream_id", v, nil),
			"log_stream_name": utils.PathSearch("log_stream_name", v, nil),
		})
	}
	return result
}

func flattenTransfersElemLogTransferInfo(transferInfo interface{}) []interface{} {
	if transferInfo == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"log_transfer_type":   utils.PathSearch("log_transfer_type", transferInfo, nil),
			"log_transfer_mode":   utils.PathSearch("log_transfer_mode", transferInfo, nil),
			"log_storage_format":  utils.PathSearch("log_storage_format", transferInfo, nil),
			"log_transfer_status": utils.PathSearch("log_transfer_status", transferInfo, nil),
			"log_agency_transfer": flattenLogTransferInfoLogAgency(transferInfo),
			"log_transfer_detail": flattenDataLogTransferInfoLogTransferDetail(transferInfo),
		},
	}
}

func flattenDataLogTransferInfoLogTransferDetail(resp interface{}) []interface{} {
	logTransferDetail := utils.PathSearch("log_transfer_detail", resp, nil)
	if logTransferDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"obs_period":           utils.PathSearch("obs_period", logTransferDetail, nil),
			"obs_period_unit":      utils.PathSearch("obs_period_unit", logTransferDetail, nil),
			"obs_bucket_name":      utils.PathSearch("obs_bucket_name", logTransferDetail, nil),
			"obs_transfer_path":    utils.PathSearch("obs_transfer_path", logTransferDetail, nil),
			"obs_dir_prefix_name":  utils.PathSearch("obs_dir_pre_fix_name", logTransferDetail, nil),
			"obs_prefix_name":      utils.PathSearch("obs_prefix_name", logTransferDetail, nil),
			"obs_eps_id":           utils.PathSearch("obs_eps_id", logTransferDetail, nil),
			"obs_encrypted_enable": utils.PathSearch("obs_encrypted_enable", logTransferDetail, nil),
			"obs_encrypted_id":     utils.PathSearch("obs_encrypted_id", logTransferDetail, nil),
			"obs_time_zone":        utils.PathSearch("obs_time_zone", logTransferDetail, nil),
			"obs_time_zone_id":     utils.PathSearch("obs_time_zone_id", logTransferDetail, nil),
			"dis_id":               utils.PathSearch("dis_id", logTransferDetail, nil),
			"dis_name":             utils.PathSearch("dis_name", logTransferDetail, nil),
			"kafka_id":             utils.PathSearch("kafka_id", logTransferDetail, nil),
			"kafka_topic":          utils.PathSearch("kafka_topic", logTransferDetail, nil),
			"delivery_tags":        utils.PathSearch("tags", logTransferDetail, nil),
		},
	}
}
