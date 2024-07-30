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
  instance_id = var.instance_id
  service     = var.service
  tag_sync    = "auto"

  tags {
    sync   = true
    key    = "key"
    values = ["value1", "value2"]
  }
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

* `tags` - (Required, List) Specifies the tags list.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `sync` - (Required, Bool) Specifies whether tag is synchronized.

* `key` - (Required, String) Specifies the tag key.

* `values` - (Required, List) Specifies the tag values list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID and the format is `<instance_id>/<service>`.

## Import

The AOM cloud service access resource can be imported using `instance_id` and `service` separated by a slash e.g.

```bash
$ terraform import huaweicloud_aom_cloud_service_access.test <instance_id>/<service>
```
