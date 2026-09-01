---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_pool"
description: |-
  Manages a WAF pool resource within HuaweiCloud.
---

# huaweicloud_waf_pool

Manages a WAF pool resource within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

## Example Usage

```hcl
variable "pool_name" {}
variable "type" {}
variable "vpc_id" {}

resource "huaweicloud_waf_pool" "test" {
  name   = var.pool_name
  type   = var.type
  vpc_id = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF pool resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the pool name. The maximum length is `256` characters. Only digits,
  letters, underscores (_), hyphens (-) and dots (.) are allowed.

* `type` - (Required, String, NonUpdatable) Specifies the pool type. Valid values are:
  + **elb**: Basic ELB type.
  + **elb-v2**: ELB-v2 type.
  + **elb-shadow**: SaaS ELB type.
  + **standard-container**: Reverse proxy dedicated engine group (in-cloud, dedicated for tenants).
  + **standard-cloud**: Reverse proxy dedicated engine group (in-cloud).
  + **standard**: Reverse proxy dedicated engine group (out-of-cloud).
  + **detector-cloud**: Bypass detection dedicated engine group (in-cloud).
  + **detector**: Bypass detection dedicated engine group (out-of-cloud).

* `vpc_id` - (Required, String, NonUpdatable) Specifies the VPC ID associated with the pool.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the pool.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID of WAF pool.
  This parameter is valid only when the enterprise project is enabled.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the pools under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The pool ID.

* `hosts` - The list of protected domain names associated with the pool.

  The [hosts](#pool_Hosts_Instances) structure is documented below.

* `instances` - The list of engine instances associated with the pool.

  The [instances](#pool_Hosts_Instances) structure is documented below.

* `create_time` - The creation time of the pool, in milliseconds.

<a name="pool_Hosts_Instances"></a>
The `hosts` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `service_ip` - The engine instance IP.

## Import

The WAF pool resource can be imported using `id` and `enterprise_project_id`, separated by a slash, e.g.

### Import resource under the default enterprise project

```bash
$ terraform import huaweicloud_waf_pool.test <id>/0
```

### Import resource from non default enterprise project

```bash
$ terraform import huaweicloud_waf_pool.test <id>/<enterprise_project_id>
```
