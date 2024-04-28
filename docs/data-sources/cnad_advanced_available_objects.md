---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_available_objects"
description: ""
---

# huaweicloud_cnad_advanced_available_objects

Use this data source to get the list of CNAD advanced available protected objects.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_cnad_advanced_available_objects" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the CNAD advanced instance ID. You can find it through data source
  `huaweicloud_cnad_advanced_instances`.

* `protected_object_id` - (Optional, String) Specifies the CNAD advanced protected object ID which you want to query.
  The protected object ID must be the ID of the Elastic IP, which in the same region with the CNAD advanced instance.

* `ip_address` - (Optional, String) Specifies the CNAD advanced protected object IP which you want to query.
  The IP must be the Elastic IP, which in the same region with the CNAD advanced instance.

* `type` - (Optional, String) Specifies the type of the protected object. This field means the type of resource which
  bound to Elastic IP. Valid values are: **VPN**, **NAT**, **VIP**,**CCI**, **EIP**, **ELB**, **REROUTING_IP**,
  **SubEni** and **NetInterFace**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `protected_objects` - Indicates the list of the CNAD advanced available protected objects.
  The [protected_objects](#advancedAvailableProtectedObjects) structure is documented below.

<a name="advancedAvailableProtectedObjects"></a>
The `protected_objects` block supports:

* `id` - Indicates the ID of the protected object.

* `ip_address` - Indicates the IP of the protected object.

* `type` - Indicates the type of the protected object.
