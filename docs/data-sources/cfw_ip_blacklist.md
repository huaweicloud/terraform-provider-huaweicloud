---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ip_blacklist"
description: |-
  Use this data source to get the imported IP blacklist information within HuaweiCloud.
---

# huaweicloud_cfw_ip_blacklist

Use this data source to get the imported IP blacklist information within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_ip_blacklist" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The IP blacklist query result.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `records` - The list of IP blacklist records.

  The [records](#records_struct) structure is documented below.

* `total` - The total number of IP blacklist in the firewall instance.

<a name="records_struct"></a>
The `records` block supports:

* `name` - The IP blacklist file name, corresponding to the export file name.

* `effect_scope` - The effect scope of the IP blacklist.  
  The valid values are as follows:
  + **1**: EIP.
  + **2**: NAT.

* `import_status` - The import status of the IP blacklist.  
  The valid values are as follows:
  + **2**: Generating.
  + **1**: Success.
  + **0**: Failed.

* `import_time` - The IP blacklist import time.

* `error_message` - The error message when import fails.
