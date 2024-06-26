---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_transfers"
description: |-
  Using this data source to query the list of log transfers within HuaweiCloud.
---

# huaweicloud_lts_transfers

Using this data source to query the list of log transfers within HuaweiCloud.

## Example Usage

```hcl
variable "log_group_name" {}

data "huaweicloud_lts_transfers" "test" {
  log_group_name = var.log_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query log transfers.
  If omitted, the provider-level region will be used.

* `log_group_name` - (Optional, String) Specifies the name of the log group to which the log transfers and log streams
  belong.

* `log_stream_name` - (Optional, String) Specifies the name of the log stream to be transferred in the log transfer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `transfers` - The list of log transfers.  
  The [transfers](#data_transfers_attr) structure is documented below.

<a name="data_transfers_attr"></a>
The `transfers` block supports:

* `id` - The ID of the transfer.

* `log_group_id` - The ID of the log group to which the log transfer belongs.

* `log_group_name` - The name of the log group to which the log transfer belongs.

* `log_streams` - The configuration of the log streams that to be transferred.  
  The [log_streams](#data_transfers_elem_log_streams) structure is documented below.

* `log_transfer_info` - The configuration of the log transfer.  
  The [log_transfer_info](#data_transfers_elem_log_transfer_info) structure is documented below.

<a name="data_transfers_elem_log_streams"></a>
The `log_streams` block supports:

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.

<a name="data_transfers_elem_log_transfer_info"></a>
The `log_transfer_info` block supports:

* `log_transfer_type` - The type of the log transfer.

* `log_transfer_mode` - The mode of the log transfer.

* `log_storage_format` - The format of the log transfer.

* `log_transfer_status` - The status of the log transfer.

* `log_agency_transfer` - The configuration of the agency transfer.  
  The [log_agency_transfer](#data_log_transfer_info_elem_log_agency_transfer) structure is documented below.

* `log_transfer_detail` - The detail of the log transfer configuration.  
  The [log_transfer_detail](#data_log_transfer_info_elem_log_transfer_detail) structure is documented below.

<a name="data_log_transfer_info_elem_log_agency_transfer"></a>
The `log_agency_transfer` block supports:

* `agency_domain_id` - The ID of the delegator account.

* `agency_domain_name` - The name of the delegator account.

* `agency_name` - The agency name created by the delegator account.

* `agency_project_id` - The project ID of the delegator account.

<a name="data_log_transfer_info_elem_log_transfer_detail"></a>
The `log_transfer_detail` block supports:

* `obs_period` - The length of the transfer interval for an OBS transfer task.

* `obs_period_unit` - The unit of the transfer interval for an OBS transfer task.

* `obs_bucket_name` - The name of the OBS bucket, which is the log transfer destination object.

* `obs_transfer_path` - The storage path of the OBS bucket, which is the log transfer destination.

* `obs_dir_prefix_name` - The custom prefix of the transfer path.

* `obs_prefix_name` - The transfer file prefix of an OBS transfer task.

* `obs_eps_id` - The enterprise project ID of an OBS transfer task.

* `obs_encrypted_enable` - Whether OBS bucket encryption is enabled.

* `obs_encrypted_id` - The KMS key ID for an OBS transfer task.

* `obs_time_zone` - The time zone for an OBS transfer task.

* `obs_time_zone_id` - ID of the time zone for an OBS transfer task.

* `dis_id` - The ID of the DIS stream.

* `dis_name` - The name of the DIS stream.

* `kafka_id` - The ID of the kafka instance.

* `kafka_topic` - The kafka topic.

* `delivery_tags` - The list of tag fields will be delivered when transferring.
