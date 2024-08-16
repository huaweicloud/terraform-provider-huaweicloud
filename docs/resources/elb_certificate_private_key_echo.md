---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_certificate_private_key_echo"
description: |-
  Manages an ELB certificate private key echo resource within HuaweiCloud.
---

# huaweicloud_elb_certificate_private_key_echo

Manages an ELB certificate private key echo resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_elb_certificate_private_key_echo" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to project ID.

## Import

The ELB certificate private key echo can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_certificate_private_key_echo.test <id>
```
