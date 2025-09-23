---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_datasource_connection"
description: ""
---

# huaweicloud_dli_datasource_connection

Manages a DLI datasource **enhanced** connection resource within HuaweiCloud.  

## Example Usage

```hcl
  variable "name" {}
  variable "vpc_id" {}
  variable "subnet_id" {}

  resource "huaweicloud_dli_datasource_connection" "test" {
    name      = var.name
    vpc_id    = var.vpc_id
    subnet_id = var.subnet_id
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) The name of a datasource connection.  
  The name can contain `1` to `64` characters, only letters, digits and underscores (_) are allowed.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) The VPC ID of the service to be connected.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) The subnet ID of the service to be connected.

  Changing this parameter will create a new resource.

* `route_table_id` - (Optional, String, ForceNew) The route table ID associated with the subnet of the service to be connected.

  Changing this parameter will create a new resource.

* `queues` - (Optional, List) List of queue names that are available for datasource connections.

* `hosts` - (Optional, List) The user-defined host information. A maximum of 20,000 records are supported.
The [Host](#datasourceConnection_Host) structure is documented below.

* `routes` - (Optional, List) List of routes.

The [Route](#datasourceConnection_Route) structure is documented below.

* `tags` - (Optional, Map, ForceNew) The key/value pairs to associate with the datasource connection.

  Changing this parameter will create a new resource.

<a name="datasourceConnection_Host"></a>
The `Host` block supports:

* `name` - (Required, String) The user-defined host name.  
  The valid length is limited from `1` to `128`, only letters, digits, hyphens (-) and underscores (_) are allowed.
  And the name must be start with a letter.

* `ip` - (Required, String) IPv4 address of the host.

<a name="datasourceConnection_Route"></a>
The `Route` block supports:

* `name` - (Required, String) The route name.  
  The valid length is limited from `1` to `64`.

* `cidr` - (Required, String) The CIDR of the route.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The connection status.  
  The options are as follows:
    + **ACTIVE**: The datasource connection is activated.
    + **DELETED**: The datasource connection is deleted.

## Import

The DLI datasource connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dli_datasource_connection.test 0ce123456a00f2591fabc00385ff1234
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `tags`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dli_datasource_connection" "test" {
    ...

    lifecycle {
      ignore_changes = [
        tags,
      ]
    }
}
