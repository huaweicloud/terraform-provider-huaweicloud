---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_image_trigger"
description: ""
---

# huaweicloud_swr_image_trigger

Manages a SWR image trigger within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "repository_name" {}
variable "cluster_id" {}
variable "namespace" {}
variable "name" {}

resource "huaweicloud_swr_image_trigger" "test"{
  organization    = var.organization_name
  repository      = var.repository_name
  workload_type   = "deployments"
  workload_name   = "test_name"
  cluster_id      = var.cluster_id
  namespace       = var.namespace
  condition_value = ".*"
  enabled         = "true"
  name            = var.name
  condition_type  = "all"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization.

  Changing this parameter will create a new resource.

* `repository` - (Required, String, ForceNew) Specifies the name of the repository.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the trigger name.

  Changing this parameter will create a new resource.

* `workload_type` - (Required, String, ForceNew) Specifies the type of the application.
  Value options: **deployments**, **statefulsets**.

  Changing this parameter will create a new resource.

* `workload_name` - (Required, String, ForceNew) Specifies the name of the application.

  Changing this parameter will create a new resource.

* `namespace` - (Required, String, ForceNew) Specifies the namespace where the application is located.

  Changing this parameter will create a new resource.

* `condition_type` - (Required, String, ForceNew) Specifies the trigger condition type.
  Value options **all**, **tag**, **regular**.

  Changing this parameter will create a new resource.

* `condition_value` - (Required, String, ForceNew) Specifies the trigger condition value. Value options:
  + When condition_type is set to `all`, set this parameter to `.*`.
  + When condition_type is set to `tag`, set this parameter to specific image tags separated by semicolons (;).
  + When condition_type is set to `regular`, set this parameter to a regular expression.

  Changing this parameter will create a new resource.

* `cluster_id` - (Optional, String, ForceNew) Specifies the ID of the cluster.
  It is required when type is set to `cce`.

  Changing this parameter will create a new resource.

* `cluster_name` - (Optional, String, ForceNew) Specifies the name of the cluster.

  Changing this parameter will create a new resource.

* `container` - (Optional, String, ForceNew) Specifies the name of the container to be updated.
  By default, all containers using this image are updated.

  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the trigger type.
  Value options **cce**, **cci**. Default to **cce**.

  Changing this parameter will create a new resource.

* `enabled` - (Optional, String) Specifies whether to enable the trigger.
  Value options **true**, **false**. Default to **true**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time.

* `creator_name` - Indicates the creator name of the trigger.

## Import

The swr image trigger can be imported using the organization name, repository name
and trigger name separated by a slash, e.g.:

Only when repository name is with no slashes, can use slashes to separate.

```bash
$ terraform import huaweicloud_swr_image_trigger.test <organization_name>/<repository_name>/<trigger_name>
```

Using comma to separate is available for repository name with slashes or not.

```bash
$ terraform import huaweicloud_swr_image_trigger.test <organization_name>,<repository_name>,<trigger_name>
```
