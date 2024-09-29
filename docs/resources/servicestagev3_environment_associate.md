---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_environment_associate"
description: |-
  Use this resource to bind resources to the environment within HuaweiCloud.
---

# huaweicloud_servicestagev3_environment_associate

Use this resource to bind resources to the environment within HuaweiCloud.

-> An environment can only have one resource.

~> Do not use this resource at the same time as resource `huaweicloud_servicestage_environment`.

## Example Usage

```hcl
variable "environment_id" {}
variable "eip_id" {}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = var.environment_id

  resources {
    id   = var.eip_id
    type = "eip"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the environment and resources are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the environment ID associated with the resources.  
  Changing this will create a new resource.

* `resources` - (Required, List) Specifies the information about the associated resources.
  The [resources](#servicestage_v3_env_associated_resources) structure is documented below.

<a name="servicestage_v3_env_associated_resources"></a>
The `resources` block supports:

* `id` - (Required, String) Specifies the ID of the resource to be associated.

* `type` - (Required, String) Specifies the type of the resource to be associated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the environment ID), in UUID format.

## Import

Associate resources can be imported using the `id` (`environment_id`), e.g.

```bash
$ terraform import huaweicloud_servicestagev3_environment_associate.test <id>
```
