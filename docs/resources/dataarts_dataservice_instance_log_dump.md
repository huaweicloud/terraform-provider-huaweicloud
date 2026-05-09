---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_instance_log_dump"
description: |-
  Use this resource to manage OBS or LTS log dump for DataArts DataService exclusive cluster within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_instance_log_dump

Use this resource to manage OBS or LTS log dump for DataArts DataService exclusive cluster within HuaweiCloud.

~> 1. Destroying this resource will disable log dump.
   </br>2. The interval for operating log dump on the same cluster must be at least `30` seconds.

-> Only one resource can be created under an exclusive cluster.

## Example Usage

### Enable OBS log dump

```hcl
variable "workspace_id" {}
variable "instance_id" {}

resource "huaweicloud_dataarts_dataservice_instance_log_dump" "test" {
  workspace_id = var.workspace_id
  instance_id  = var.instance_id
  type         = "obs"
}
```

### Enable LTS log dump

```hcl
variable "workspace_id" {}
variable "instance_id" {}
variable "log_group_id" {}
variable "log_group_name" {}
variable "log_stream_id" {}
variable "log_stream_name" {}

resource "huaweicloud_dataarts_dataservice_instance_log_dump" "test" {
  workspace_id    = var.workspace_id
  instance_id     = var.instance_id
  type            = "lts"
  log_group_id    = var.log_group_id
  log_group_name  = var.log_group_name
  log_stream_id   = var.log_stream_id
  log_stream_name = var.log_stream_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the log dump is located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the log dump belongs.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the exclusive cluster to which the log
  dump belongs.

* `type` - (Required, String, NonUpdatable) Specifies the type of the log dump to be configured.  
  The valid values are as follows:
  + **obs**
  + **lts**

  -> When OBS log dump is enabled, the system will automatically create an OBS bucket if it does not exist, the bucket
   name consists of `dlm` and the current project ID, separated by a hyphen (-). e.g. `dlm-{project_id}`.

* `log_group_id` - (Optional, String, NonUpdatable) Specifies the ID of the LTS log group.  
  This parameter is required and available when the `type` is **lts**.

* `log_group_name` - (Optional, String, NonUpdatable) Specifies the name of the LTS log group.  
  This parameter is required and available when the `type` is **lts**.

* `log_stream_id` - (Optional, String, NonUpdatable) Specifies the ID of the LTS log stream.  
  This parameter is required and available when the `type` is **lts**.

* `log_stream_name` - (Optional, String, NonUpdatable) Specifies the name of the LTS log stream.  
  This parameter is required and available when the `type` is **lts**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using `workspace_id` and `instance_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_dataservice_instance_log_dump.test <workspace_id>/<instance_id>
```
