---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_access_control_logs"
description: |-
  Use this data source to get the list of CFW access control logs.
---

# huaweicloud_cfw_access_control_logs

Use this data source to get the list of CFW access control logs.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_cfw_access_control_logs" "test" {
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

* `start_time` - (Required, String) Specifies the start time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `end_time` - (Required, String) Specifies the end time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `src_ip` - (Optional, String) Specifies the source IP address.

* `src_port` - (Optional, Int) Specifies the source port.

* `dst_ip` - (Optional, String) Specifies the destination IP address.

* `dst_port` - (Optional, Int) Specifies the destination port.

* `app` - (Optional, String) Specifies the application protocol.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `rule_name` - (Optional, String) Specifies the rule name.

* `action` - (Optional, String) Specifies the action. The values can be **allow** and **deny**.

* `src_region_name` - (Optional, String) Specifies the source region name.

* `dst_region_name` - (Optional, String) Specifies the destination region name.

* `src_province_name` - (Optional, String) Specifies the source province name.

* `dst_province_name` - (Optional, String) Specifies the destination province name.

* `src_city_name` - (Optional, String) Specifies the source city name.

* `dst_city_name` - (Optional, String) Specifies the destination city name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The access control log records.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `src_ip` - The source IP address.

* `src_port` - The source port.

* `dst_port` - The destination port.

* `app` - The application protocol.

* `rule_name` - The rule name.

* `rule_id` - The rule ID.

* `hit_time` - The hit time.

* `log_id` - The document ID.

* `dst_ip` - The destination IP address.

* `protocol` - The protocol type.

* `action` - The action.

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

* `dst_host` - The destination host.
