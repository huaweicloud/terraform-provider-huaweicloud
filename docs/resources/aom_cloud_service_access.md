---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_cloud_service_access"
description:  |-
  Manages an AOM cloud service access resource within HuaweiCloud.
---

# huaweicloud_aom_cloud_service_access

Manages an AOM cloud service access resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "service" {}

resource "huaweicloud_aom_cloud_service_access" "test" {
  instance_id           = var.instance_id
  service               = var.service
  tag_sync              = "auto"
  enterprise_project_id = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the AOM prometheus instance ID.
  Changing this parameter will create a new resource.

* `service` - (Required, String, ForceNew) Specifies the service name.
  Changing this parameter will create a new resource.

* `tag_sync` - (Required, String) Specifies whether tags are automatically synchronized.
  Valid values are **auto** and **manual**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the instance belongs.
  Defaults to **0**. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID and the format is `<instance_id>/<service>`.

## Import

The AOM cloud service access resource can be imported using `instance_id` and `service` separated by a slash e.g.

```bash
$ terraform import huaweicloud_aom_cloud_service_access.test <instance_id>/<service>
```
