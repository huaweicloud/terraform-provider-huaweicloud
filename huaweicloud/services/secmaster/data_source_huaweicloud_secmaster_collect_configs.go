package secmaster

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

// @API SecMaster GET /v1/{project_id}/collector/cloudlogs/config
func DataSourceCollectConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCollectConfigRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"csvc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_statistics": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"all_vendors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     allVendorsSchema(),
			},
			"config_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     configStatisticsSchema(),
			},
			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     datasetInfoSchema(),
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataListSchema(),
			},
			"lts_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ltsSetsSchema(),
			},
		},
	}
}

func allVendorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csvc_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     csvcListSchema(),
			},
		},
	}
}

func csvcListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sourceListSchema(),
			},
		},
	}
}

func sourceListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc_display": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csvc_help": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_display": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_help": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func configStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"daily_traffic": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_all_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_in_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vendor_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func datasetInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_region": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     referenceSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     targetInfoSchema(),
			},
		},
	}
}

func referenceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc_display": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"csvc_help": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_display": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_help": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func targetInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pipe": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shards": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_all_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_successful_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_all_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_in_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_in_num_last_one_hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataListDatasetSchema(),
			},
		},
	}
}

func dataListDatasetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"all_accounts": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"new_account_auto_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_all_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_successful_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sink_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     referenceSchema(),
			},
			"target": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     targetInfoSchema(),
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     accountSchema(),
			},
		},
	}
}

func accountSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sink_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_log_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"log_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func ltsSetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildListCollectConfigQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d&domain_id=%s", limit, offset, d.Get("domain_id").(string))

	if v, ok := d.GetOk("region_id"); ok {
		queryParams = fmt.Sprintf("%s&region_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("csvc"); ok {
		queryParams = fmt.Sprintf("%s&csvc=%v", queryParams, v)
	}
	if v, ok := d.GetOk("query_statistics"); ok {
		queryParams = fmt.Sprintf("%s&query_statistics=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceCollectConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/collector/cloudlogs/config"
		limit   = 500
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath+buildListCollectConfigQueryParams(d, limit, offset), &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster collect config: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("workspace_id", utils.PathSearch("workspace_id", respBody, nil)),
		d.Set("dataspace_id", utils.PathSearch("dataspace_id", respBody, nil)),
		d.Set("dataspace_name", utils.PathSearch("dataspace_name", respBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", respBody, nil)),
		d.Set("all_vendors", flattenAllVendors(utils.PathSearch("all_vendors", respBody, make([]interface{}, 0)))),
		d.Set("config_statistics", flattenConfigStatistics(utils.PathSearch("config_statistics", respBody, nil))),
		d.Set("datasets", flattenDatasetInfos(utils.PathSearch("datasets", respBody, make([]interface{}, 0)))),
		d.Set("data_list", flattenDataList(utils.PathSearch("data_list", respBody, make([]interface{}, 0)))),
		d.Set("lts_sets", flattenLtsSets(utils.PathSearch("lts_sets", respBody, make([]interface{}, 0)))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAllVendors(allVendors interface{}) []interface{} {
	if allVendors == nil {
		return nil
	}

	vendorList, ok := allVendors.([]interface{})
	if !ok || len(vendorList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(vendorList))
	for _, item := range vendorList {
		result = append(result, map[string]interface{}{
			"cloud_vendor": utils.PathSearch("cloud_vendor", item, nil),
			"csvc_list":    flattenCsvcList(utils.PathSearch("csvc_list", item, make([]interface{}, 0))),
		})
	}

	return result
}

func flattenCsvcList(csvcList interface{}) []interface{} {
	if csvcList == nil {
		return nil
	}

	list, ok := csvcList.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		result = append(result, map[string]interface{}{
			"csvc":        utils.PathSearch("csvc", item, nil),
			"source_list": flattenSourceList(utils.PathSearch("source_list", item, make([]interface{}, 0))),
		})
	}

	return result
}

func flattenSourceList(sourceList interface{}) []interface{} {
	if sourceList == nil {
		return nil
	}

	list, ok := sourceList.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		result = append(result, map[string]interface{}{
			"csvc_display":   utils.PathSearch("csvc_display", item, nil),
			"csvc_help":      utils.PathSearch("csvc_help", item, nil),
			"source_display": utils.PathSearch("source_display", item, nil),
			"source_help":    utils.PathSearch("source_help", item, nil),
			"link":           utils.PathSearch("link", item, nil),
		})
	}

	return result
}

func flattenConfigStatistics(configStatistics interface{}) []interface{} {
	if configStatistics == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"account_num":     utils.PathSearch("account_num", configStatistics, nil),
			"daily_traffic":   utils.PathSearch("daily_traffic", configStatistics, nil),
			"log_num":         utils.PathSearch("log_num", configStatistics, nil),
			"product_all_num": utils.PathSearch("product_all_num", configStatistics, nil),
			"product_in_num":  utils.PathSearch("product_in_num", configStatistics, nil),
			"vendor_num":      utils.PathSearch("vendor_num", configStatistics, nil),
		},
	}
}

func flattenDatasetInfos(datasets interface{}) []interface{} {
	if datasets == nil {
		return nil
	}

	datasetList, ok := datasets.([]interface{})
	if !ok || len(datasetList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(datasetList))
	for _, item := range datasetList {
		result = append(result, map[string]interface{}{
			"csvc":        utils.PathSearch("csvc", item, nil),
			"enable":      utils.PathSearch("enable", item, nil),
			"is_region":   utils.PathSearch("is_region", item, nil),
			"source_id":   utils.PathSearch("source_id", item, nil),
			"source_name": utils.PathSearch("source_name", item, nil),
			"type":        utils.PathSearch("type", item, nil),
			"reference":   flattenReference(utils.PathSearch("reference", item, nil)),
			"target":      flattenTargetInfo(utils.PathSearch("target", item, nil)),
		})
	}

	return result
}

func flattenReference(reference interface{}) []interface{} {
	if reference == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"csvc_display":   utils.PathSearch("csvc_display", reference, nil),
			"csvc_help":      utils.PathSearch("csvc_help", reference, nil),
			"source_display": utils.PathSearch("source_display", reference, nil),
			"source_help":    utils.PathSearch("source_help", reference, nil),
			"link":           utils.PathSearch("link", reference, nil),
		},
	}
}

func flattenTargetInfo(target interface{}) []interface{} {
	if target == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"pipe":   utils.PathSearch("pipe", target, nil),
			"shards": utils.PathSearch("shards", target, nil),
			"ttl":    utils.PathSearch("ttl", target, nil),
		},
	}
}

func flattenDataList(dataList interface{}) []interface{} {
	if dataList == nil {
		return nil
	}

	dataListRaw, ok := dataList.([]interface{})
	if !ok || len(dataListRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataListRaw))
	for _, item := range dataListRaw {
		result = append(result, map[string]interface{}{
			"csvc":                     utils.PathSearch("csvc", item, nil),
			"vendor":                   utils.PathSearch("vendor", item, nil),
			"process_status":           utils.PathSearch("process_status", item, nil),
			"account_all_num":          utils.PathSearch("account_all_num", item, nil),
			"account_successful_num":   utils.PathSearch("account_successful_num", item, nil),
			"log_all_num":              utils.PathSearch("log_all_num", item, nil),
			"log_in_num":               utils.PathSearch("log_in_num", item, nil),
			"log_in_num_last_one_hour": utils.PathSearch("log_in_num_last_one_hour", item, nil),
			"last_modified_time":       utils.PathSearch("last_modified_time", item, nil),
			"datasets":                 flattenDataListDatasets(utils.PathSearch("datasets", item, make([]interface{}, 0))),
		})
	}

	return result
}

func flattenDataListDatasets(datasets interface{}) []interface{} {
	if datasets == nil {
		return nil
	}

	datasetList, ok := datasets.([]interface{})
	if !ok || len(datasetList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(datasetList))
	for _, item := range datasetList {
		result = append(result, map[string]interface{}{
			"source_id":               utils.PathSearch("source_id", item, nil),
			"source_name":             utils.PathSearch("source_name", item, nil),
			"enable":                  utils.PathSearch("enable", item, nil),
			"process_status":          utils.PathSearch("process_status", item, nil),
			"alert":                   utils.PathSearch("alert", item, nil),
			"all_accounts":            utils.PathSearch("all_accounts", item, nil),
			"new_account_auto_access": utils.PathSearch("new_account_auto_access", item, nil),
			"region_id":               utils.PathSearch("region_id", item, nil),
			"workspace_id":            utils.PathSearch("workspace_id", item, nil),
			"account_all_num":         utils.PathSearch("account_all_num", item, nil),
			"account_successful_num":  utils.PathSearch("account_successful_num", item, nil),
			"sink_msg":                utils.PathSearch("sink_msg", item, nil),
			"reference":               flattenReference(utils.PathSearch("reference", item, nil)),
			"target":                  flattenTargetInfo(utils.PathSearch("target", item, nil)),
			"accounts":                flattenAccounts(utils.PathSearch("accounts", item, make([]interface{}, 0))),
		})
	}

	return result
}

func flattenAccounts(accounts interface{}) []interface{} {
	if accounts == nil {
		return nil
	}

	accountList, ok := accounts.([]interface{})
	if !ok || len(accountList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(accountList))
	for _, item := range accountList {
		result = append(result, map[string]interface{}{
			"account_id":     utils.PathSearch("account_id", item, nil),
			"name":           utils.PathSearch("name", item, nil),
			"process_status": utils.PathSearch("process_status", item, nil),
			"sink_msg":       utils.PathSearch("sink_msg", item, nil),
			"last_log_date":  utils.PathSearch("last_log_date", item, nil),
			"log_count":      utils.PathSearch("log_count", item, nil),
		})
	}

	return result
}

func flattenLtsSets(ltsSets interface{}) []interface{} {
	if ltsSets == nil {
		return nil
	}

	ltsSetsList, ok := ltsSets.([]interface{})
	if !ok || len(ltsSetsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(ltsSetsList))
	for _, item := range ltsSetsList {
		result = append(result, map[string]interface{}{
			"config_name":   utils.PathSearch("config_name", item, nil),
			"enable":        utils.PathSearch("enable", item, nil),
			"log_group_id":  utils.PathSearch("log_group_id", item, nil),
			"log_stream_id": utils.PathSearch("log_stream_id", item, nil),
			"log_type":      utils.PathSearch("log_type", item, nil),
			"pipe_alias":    utils.PathSearch("pipe_alias", item, nil),
			"type_prefix":   utils.PathSearch("type_prefix", item, nil),
		})
	}

	return result
}
