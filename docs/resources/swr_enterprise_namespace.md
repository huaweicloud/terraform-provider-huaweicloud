---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_namespace"
description: |-
  Manages a SWR enterprise namespace resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_namespace

Manages a SWR enterprise namespace resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "huaweicloud_swr_enterprise_namespace" "test" {
  instance_id = var.instance_id
  name        = var.name

  metadata {
    public = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `name` - (Required, String) Specifies the namespace name.

* `metadata` - (Required, List) Specifies the metadata.
  The [metadata](#block--metadata) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the instance namespace.

<a name="block--metadata"></a>
The `metadata` block supports:

* `public` - (Required, String) Specifies whether the namespace is public.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `namespace_id` - Indicates the namespace ID.

* `repo_count` - Indicates the repo count of the namespace.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The namespace can be imported using `instance_id` and `name`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_namespace.test <instance_id>/<name>
```
