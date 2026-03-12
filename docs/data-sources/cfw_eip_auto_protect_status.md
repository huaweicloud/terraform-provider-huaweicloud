---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_auto_protect_status"
description: |-
  Use this data source to get the EIP auto protect status information within HuaweiCloud.
---

# huaweicloud_cfw_eip_auto_protect_status

Use this data source to get the EIP auto protect status information within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_eip_auto_protect_status" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The firewall status object.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `available_eip_count` - The number of EIPs that can be protected.

* `beyond_max_count` - Whether the EIP count limit is exceeded.

* `eip_protected_self` - The number of protected EIPs.

* `eip_total` - The total number of EIPs.

* `eip_un_protected` - The number of unprotected EIPs.

* `object_id` - The protected object ID.

* `status` - Whether the auto protection for new EIPs is enabled.  
  The valid values are as follows:
  + **1**: Yes.
  + **0**: No.
