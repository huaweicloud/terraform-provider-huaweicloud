---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_count"
description: |-
  Use this data source to get the CFW EIP count.
---

# huaweicloud_cfw_eip_count

Use this data source to get the CFW EIP count.

## Example Usage

```hcl
variable object_id {}

data "huaweicloud_cfw_eip_count" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.
  You can get this value from data source `huaweicloud_cfw_firewalls`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  For enterprise users, if omitted, all enterprise project will be used.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `eip_protected` - The total number of EIPs protected by all firewalls in the account.

* `eip_protected_self` - The number of EIPs protected by the current firewall.

* `eip_total` - The total number of EIPs.
