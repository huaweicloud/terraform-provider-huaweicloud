---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_addon"
description: |-
  Manages a CCE Autopilot add-on resource within huaweicloud.
---

# huaweicloud_cce_autopilot_addon

Manages a CCE Autopilot add-on resource within huaweicloud.

## Example Usage

### Basic Usage

variable "cluster_id" {}
variable "basic_json_string" {}
variable "flavor_json_string" {}
variable "custom_json_string" {}

```hcl
resource "huaweicloud_cce_autopilot_addon" "test" {
  cluster_id          = var.cluster_id
  version             = "1.4.3"
  addon_template_name = "log-agent"
  values = {
    "basic"  = jsonencode(var.basic_json_string)
    "flavor" = jsonencode(var.flavor_json_string)
    "custom" = jsonencode(var.custom_json_string)
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE autopilot add-on resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new add-on resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `addon_template_name` - (Required, String, NonUpdatable) Specifies the name of the add-on template to be installed,
  for example, **coredns**.

* `values` - (Required, Map) Specifies the add-on template installation parameters, varying depending on the add-on.
  The values of this map should be json strings.

* `version` - (Optional, String) Specifies the version of the add-on.

* `name` - (Optional, String) Specifies the add-on name.

* `alias` - (Optional, String) Specifies the add-on alias.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the add-on resource.
  
* `created_at` - The time when the add-on was created.

* `updated_at` - The time when the add-on was updated.

* `status` - The status of the add-on.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The autopilot add-on can be imported using the add-on ID, e.g.

```bash
 $ terraform import huaweicloud_cce_autopilot_addon.myaddon <addon_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `values`. It is generally
recommended running `terraform plan` after importing a add-on. You can then decide if changes should be applied to
the add-on, or the resource definition should be updated to align with the add-on. Also you can ignore changes as
below.

```hcl
resource "huaweicloud_cce_autopilot_addon" "myaddon" {
    ...

  lifecycle {
    ignore_changes = [
      values
    ]
  }
}
```
