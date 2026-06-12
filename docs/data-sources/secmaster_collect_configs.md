---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collect_configs"
description: |-
  Use this data source to get the list of SecMaster collect configs within HuaweiCloud.
---

# huaweicloud_secmaster_collect_configs

Use this data source to get the list of SecMaster collect configs within HuaweiCloud.

## Example Usage

```hcl
variable "region_id" {}
variable "domain_id" {}

data "huaweicloud_secmaster_collect_configs" "test" {
  region_id        = var.region_id
  domain_id        = var.domain_id
  query_statistics = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `region_id` - (Required, String) Specifies the region ID of the collect config.

* `domain_id` - (Required, String) Specifies the domain ID of the collect config.

* `csvc` - (Optional, String) Specifies the cloud product to filter.

* `query_statistics` - (Optional, Bool) Specifies whether to query statistics.

* `sort_key` - (Optional, String) Specifies the sort key.

* `sort_dir` - (Optional, String) Specifies the sort direction.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workspace_id` - The workspace ID.

* `dataspace_id` - The dataspace ID.

* `dataspace_name` - The dataspace name.

* `project_id` - The project ID.

* `all_vendors` - The list of all cloud vendors.

  The [all_vendors](#all_vendors_struct) structure is documented below.

* `config_statistics` - The statistics of the collect config.

  The [config_statistics](#config_statistics_struct) structure is documented below.

* `datasets` - The list of dataset information.

  The [datasets](#datasets_struct) structure is documented below.

* `data_list` - The list of data details.

  The [data_list](#data_list_struct) structure is documented below.

* `lts_sets` - The list of LTS configurations.

  The [lts_sets](#lts_sets_struct) structure is documented below.

<a name="all_vendors_struct"></a>
The `all_vendors` block supports:

* `cloud_vendor` - The cloud vendor name.

* `csvc_list` - The list of cloud products.

  The [csvc_list](#csvc_list_struct) structure is documented below.

<a name="csvc_list_struct"></a>
The `csvc_list` block supports:

* `csvc` - The cloud product name.

* `source_list` - The list of log sources.

  The [source_list](#source_list_struct) structure is documented below.

<a name="source_list_struct"></a>
The `source_list` block supports:

* `csvc_display` - The cloud product display name.

* `csvc_help` - The cloud product help information.

* `source_display` - The log source display name.

* `source_help` - The log source help information.

* `link` - The help link.

<a name="config_statistics_struct"></a>
The `config_statistics` block supports:

* `account_num` - The number of accounts.

* `daily_traffic` - The daily traffic in bytes.

* `log_num` - The number of logs.

* `product_all_num` - The total number of cloud products.

* `product_in_num` - The number of enabled cloud products.

* `vendor_num` - The number of cloud vendors.

<a name="datasets_struct"></a>
The `datasets` block supports:

* `csvc` - The cloud product name.

* `enable` - The enable status.

* `is_region` - Whether it is a region resource.

* `source_id` - The data source ID.

* `source_name` - The data source name.

* `type` - The type of the dataset.

* `reference` - The reference information.

  The [reference](#reference_struct) structure is documented below.

* `target` - The target information.

  The [target](#target_struct) structure is documented below.

<a name="reference_struct"></a>
The `reference` block supports:

* `csvc_display` - The cloud product display name.

* `csvc_help` - The cloud product help information.

* `source_display` - The log source display name.

* `source_help` - The log source help information.

* `link` - The help link.

<a name="target_struct"></a>
The `target` block supports:

* `pipe` - The data pipe name.

* `shards` - The number of shards.

* `ttl` - The time to live in days.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `csvc` - The cloud product name.

* `vendor` - The cloud vendor name.

* `process_status` - The process status.

* `account_all_num` - The total number of accounts.

* `account_successful_num` - The number of successful accounts.

* `log_all_num` - The total number of logs.

* `log_in_num` - The number of ingested logs.

* `log_in_num_last_one_hour` - The number of logs ingested in the last hour.

* `last_modified_time` - The last modified time.

* `datasets` - The list of datasets.

  The [data_list_datasets](#data_list_datasets_struct) structure is documented below.

<a name="data_list_datasets_struct"></a>
The `data_list_datasets` block supports:

* `source_id` - The data source ID.

* `source_name` - The data source name.

* `enable` - The enable status.

* `process_status` - The process status.

* `alert` - Whether automatic alert conversion is enabled.

* `all_accounts` - Whether all managed accounts are accessed.

* `new_account_auto_access` - Whether automatic synchronization of new accounts is enabled.

* `region_id` - The region ID.

* `workspace_id` - The workspace ID.

* `account_all_num` - The total number of accounts.

* `account_successful_num` - The number of successful accounts.

* `sink_msg` - The sink message.

* `reference` - The reference information.

  The [reference](#reference_struct) structure is documented above.

* `target` - The target information.

  The [target](#target_struct) structure is documented above.

* `accounts` - The list of managed accounts.

  The [accounts](#accounts_struct) structure is documented below.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `account_id` - The account ID.

* `name` - The account name.

* `process_status` - The process status.

* `sink_msg` - The sink message.

* `last_log_date` - The last log date.

* `log_count` - The log count.

<a name="lts_sets_struct"></a>
The `lts_sets` block supports:

* `config_name` - The configuration name.

* `enable` - The enable status.

* `log_group_id` - The log group ID.

* `log_stream_id` - The log stream ID.

* `log_type` - The log type.

* `pipe_alias` - The pipe alias.

* `type_prefix` - The log type prefix.
