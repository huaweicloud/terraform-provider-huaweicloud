---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_enterprise_project"
description: ""
---

# huaweicloud_enterprise_project

Use this resource to manage an enterprise project within HuaweiCloud.

-> **NOTE:** Deleting enterprise projects is not support. If you destroy a resource of enterprise project,
  the project is only disabled and removed from the state, but it remains in the cloud

## Example Usage

```hcl
resource "huaweicloud_enterprise_project" "test" {
  name        = "test"
  description = "example project"
}
```

## Argument Reference

* `name` - (Required, String) Specifies the name of the enterprise project.
  This parameter can contain `1` to `64` characters. Only English letters, Chinese characters, digits, underscores (_),
  and hyphens (-) are allowed.  
  The name must be unique in the domain and cannot include any form of the word "default" ("deFaulT", for instance).

* `description` - (Optional, String) Specifies the description of the enterprise project.

* `type` - (Optional, String) Specifies the type of the enterprise project.
  The valid values are **poc** and **prod**, defaults to **prod**.

* `enable` - (Optional, Bool) Specifies whether to enable the enterprise project. Defaults to **true**.

* `skip_disable_on_destroy` - (Optional, Bool) Specifies whether to skip disable the enterprise project on destroy.
  Defaults to **false**.

* `delete_flag` - (Optional, Bool) Specifies whether to delete enterprise project.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the enterprise project.

* `status` - Indicates the status of an enterprise project.
  + **1**: Indicates enabled.
  + **2**: Indicates disabled.

* `created_at` - Indicates the time (UTC) when the enterprise project was created. Example: **2018-05-18T06:49:06Z**.

* `updated_at` - Indicates the time (UTC) when the enterprise project was modified. Example: **2018-05-28T02:21:36Z**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Enterprise projects can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_enterprise_project.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `delete_flag`. It is generally
recommended running **terraform plan** after importing an enterprise project. You can then decide if changes should be
applied to the enterprise project, or the resource definition should be updated to align with the enterprise project.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_enterprise_project" "test" {
    ...

  lifecycle {
    ignore_changes = [
      delete_flag
    ]
  }
}
```
