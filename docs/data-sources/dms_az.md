---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_az"
description: ""
---

# huaweicloud_dms_az

Use this data source to get the ID of an available HuaweiCloud dms az.

!> **WARNING:** It has been deprecated. This data source is used for the `available_zones` of the
`huaweicloud_dms_kafka_instance` and `huaweicloud_dms_rabbitmq_instance` resource.
Now argument `available_zones` has been deprecated, instead `availability_zones`,
this data source will no longer be used.

## Example Usage

```hcl
data "huaweicloud_dms_az" "az1" {
  code = "cn-north-4a"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms az. If omitted, the provider-level region will be
  used.

* `code` - (Optional, String) Specifies the code of an AZ.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.

* `name` - Indicates the name of an AZ.

* `port` - Indicates the port number of an AZ.

* `ipv6_enabled` - Whether the IPv6 network is enabled.
