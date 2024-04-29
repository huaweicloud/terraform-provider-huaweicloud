---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_policy"
description: ""
---

# huaweicloud_cnad_advanced_policy

Manages a CNAD advanced policy resource within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cnad_advanced_instances" "test" {}

resource "huaweicloud_cnad_advanced_policy" "test" {
  instance_id = huaweicloud_cnad_advanced_instances.test.instances[0].instance_id
  name        = "test-policy"
  threshold   = 100
  udp         = "block"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the CNAD advanced instance ID.
  You can find it through data source `huaweicloud_cnad_advanced_instances`.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the policy name, the maximum length is 255 characters.

* `udp` - (Optional, String) Specifies whether to block the UDP protocol. Valid values are **block** and **unblock**.

* `threshold` - (Optional, Int) Specifies the cleaning threshold, the value ranges from 100 Mbps to 1000 Mbps.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `block_location` - The location block list.

* `block_protocol` - The protocol block list.

* `connection_protection` - Whether enable connection protection.

* `connection_protection_list` - The connection protection list.

* `fingerprint_count` - The fingerprint count.

* `port_block_count` - The number of port blockages.

* `watermark_count` - The number of watermarks.

## Import

The CNAD advanced policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cnad_advanced_policy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `udp`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_cnad_advanced_policy" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      udp,
    ]
  }
}
```
