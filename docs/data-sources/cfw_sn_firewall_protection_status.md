---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_sn_firewall_protection_status"
description: |-
  Use this data source to get the SN firewall protection status information within HuaweiCloud.
---

# huaweicloud_cfw_sn_firewall_protection_status

Use this data source to get the SN firewall protection status information within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_sn_firewall_protection_status" "test" {
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

* `data` - The SN firewall protection status data.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `protection_status` - The firewall protection status.  
  The valid values are as follows:
  + **0**: Normal status.
  + **1**: Bypass in progress.
  + **2**: Bypass succeeded.
  + **3**: Bypass failed.
  + **4**: Recovery in progress.
  + **5**: Recovery failed.

* `id` - The firewall instance ID.

* `object_id` - The protected object ID.

* `failed_eip_list` - The list of EIPs that failed to bypass.

* `failed_eip_id_list` - The list of EIP IDs that failed to bypass.
