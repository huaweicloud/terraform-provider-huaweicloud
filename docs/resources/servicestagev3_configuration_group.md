---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_configuration_group"
description: |-
  Manages a configuration group resource within HuaweiCloud.
---

# huaweicloud_servicestagev3_configuration_group

Manages a configuration group resource within HuaweiCloud.

## Example Usage

```hcl
variables "group_name" {}

resource "huaweicloud_servicestagev3_configuration_group" "test" {
  name = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the configuration group.  
  The valid length is limited from `2` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.  
  The name must start with a letter and end with a letter or a digit.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the configuration group.  
  The maximum length is `256` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also configuration group ID.

* `creator` - The creator of the configuration group.

* `created_at` - The creation time of the configuration group, in RFC3339 format.

## Import

The resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_servicestagev3_configuration_group.test <id>
```
