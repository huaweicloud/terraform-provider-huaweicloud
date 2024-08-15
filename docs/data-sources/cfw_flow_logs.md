---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_flow_logs"
description: |-
  Use this data source to get the list of CFW flow logs.
---

# huaweicloud_cfw_flow_logs

Use this data source to get the list of CFW flow logs.

-> **NOTE:** Up to 1000 logs can be retrieved. Set filter criteria to narrow down the search scope.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_cfw_flow_logs" "test" {
  fw_instance_id = var.fw_instance_id
  start_time     = var.start_time
  end_time       = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `direction` - (Optional, String) Specifies the direction. The values can be **out2in** and **in2out**.

* `start_time` - (Required, String) Specifies the start time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `end_time` - (Required, String) Specifies the end time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `src_ip` - (Optional, String) Specifies the source IP address.

* `src_port` - (Optional, Int) Specifies the source port.

* `dst_ip` - (Optional, String) Specifies the destination IP address.

* `dst_port` - (Optional, Int) Specifies the destination port.

* `app` - (Optional, String) Specifies the application protocol.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project id.

* `src_region_name` - (Optional, String) Specifies the source region name.

* `dst_region_name` - (Optional, String) Specifies the destination region name.

* `src_province_name` - (Optional, String) Specifies the source province name.

* `dst_province_name` - (Optional, String) Specifies the destination province name.

* `src_city_name` - (Optional, String) Specifies the source city name.

* `dst_city_name` - (Optional, String) Specifies the destination city name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The flow log records.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `direction` - The direction, which can be inbound or outbound.

* `end_time` - The end time.

* `src_ip` - The source IP address.

* `dst_ip` - The destination IP address.

* `bytes` - The flow log bytes.

* `start_time` - The start time.

* `log_id` - The document ID.

* `src_port` - The source port.

* `app` - The application protocol.

* `dst_port` - The destination port.

* `protocol` - The protocol type.

* `packets` - The number of packets.

* `dst_host` - The destination host.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `src_province_id` - The source province ID.

* `src_province_name` - The source province name.

* `src_city_id` - The source city ID.

* `src_city_name` - The source city name.

* `dst_province_id` - The destination province ID.

* `dst_province_name` - The destination province name.

* `dst_city_id` - The destination city ID.

* `dst_city_name` - The destination city name.
