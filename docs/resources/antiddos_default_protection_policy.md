---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_default_protection_policy"
description: |-
  Manages Cloud Native Anti-DDos default protection policy resource within HuaweiCloud.
---

# huaweicloud_antiddos_default_protection_policy

Manages Cloud Native Anti-DDos default protection policy resource within HuaweiCloud.

-> Destroying the resource actually sets the field `traffic_threshold` to 120 Mbps.

## Example Usage

```hcl
resource "huaweicloud_antiddos_default_protection_policy" "test" {
  traffic_threshold = 150
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `traffic_threshold` - (Required, Int) Specifies the traffic cleaning threshold in Mbps.
  The value can be `10`, `30`, `50`, `70`, `100`, `120`, `150`, `200`, `250`, `300`, `1,000` Mbps.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Cloud Native Anti-DDos default protection policy resource can be imported using `id`. e.g.

```bash
$ terraform import huaweicloud_antiddos_default_protection_policy.test <id>
```
