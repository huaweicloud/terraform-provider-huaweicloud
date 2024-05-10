---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_connection_associate"
description: ""
---

# huaweicloud_dli_datasource_connection_associate

Using this resource to associate the elastic resource pools to the DLI datasource **enhanced** connection within
HuaweiCloud.

-> A connection can only have one resource.

## Example Usage

### Grant queue binding permission to the datasource enhanced connection

```hcl
variable "connection_id" {}
variable "associated_pool_names" {
  type = list(string)
}

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = connection_id
  elastic_resource_pools = var.associated_pool_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `connection_id` - (Required, String, ForceNew) Specifies the ID of the datasource **enhanced** connection to be
  associated.  
  Changing this parameter will create a new resource.

* `elastic_resource_pools` - (Required, List) Specifies the list of elastic resoruce pool names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also `connection_id`.

## Import

The associate relationship can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dli_datasource_connection_associate.test <id>
```
