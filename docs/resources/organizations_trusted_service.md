---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_trusted_service"
description: ""
---

# huaweicloud_organizations_trusted_service

Manages an Organizations trusted service resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_organizations_trusted_service" "test"{
  service = "service.AOM"
}
```

## Argument Reference

The following arguments are supported:

* `service` - (Required, String, ForceNew) Specifies the name of the trusted service principal.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `enabled_at` - Indicates the date when the trusted service was integrated with Organizations.

## Import

The organizations trusted service can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_trusted_service.test <id>
```
