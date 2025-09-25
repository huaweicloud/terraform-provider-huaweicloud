---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_domain_route_update"
description: |-
  Manages a WAF domain route update resource within HuaweiCloud.
---

# huaweicloud_waf_domain_route_update

Manages a WAF domain route update resource within HuaweiCloud.

-> This resource is only a one-time action resource using to update WAF domain route. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "domain_id" {}

resource "huaweicloud_waf_domain_route_update" "test" {
  instance_id = var.domain_id

  routes {
    name  = "example_route"
    cname = "example_cname"

    servers {
      back_protocol = "HTTP"
      address       = "192.168.1.1"
      port          = 80
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this setting will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of WAF domain.

* `routes` - (Required, List, NonUpdatable) Specifies the list of route configurations.

  The [routes](#domain_route_update_routes) structure is documented below.

<a name="domain_route_update_routes"></a>
The `routes` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the WAF cluster.

* `servers` - (Required, List, NonUpdatable) Specifies the list of protected domain source site server information.

  The [servers](#domain_route_update_servers) structure is documented below.

* `cname` - (Optional, String, NonUpdatable) Specifies the cname suffix of the WAF cluster.

<a name="domain_route_update_servers"></a>
The `servers` block supports:

* `back_protocol` - (Optional, String, NonUpdatable) Specifies the protocol for WAF to forward client requests to the
  protected domain origin server. The valid values are **HTTP** and **HTTPS**.

* `address` - (Optional, String, NonUpdatable) Specifies the IP address of the source server for client access.

* `port` - (Optional, Int, NonUpdatable) Specifies the business port for WAF to forward client requests to the source
  service.

## Attribute Reference

The following attributes are exported:

* `id` - The resource ID, which is the same as the `instance_id`.
