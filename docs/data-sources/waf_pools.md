---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_pools"
description: |-
  Use this data source to get the list of WAF pools within HuaweiCloud.
---

# huaweicloud_waf_pools

Use this data source to get the list of WAF pools within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_pools" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID used for filtering.
  If you need to query the resources bound to all enterprise projects of the current user, set this parameter to
  **all_granted_eps**. The default value is **0**, indicating the default enterprise project.

* `name` - (Optional, String) Specifies the pool name used for filtering. The `name` uses fuzzy matching.

* `type` - (Optional, List) Specifies the pool type list used for filtering.  
  The valid values are as follows:
  + **elb**: basic ELB type.
  + **elb-v2**: ELB-v2 type.
  + **elb-shadow**: SaaS ELB type.
  + **standard-container**: reverse proxy dedicated engine group (in-cloud, dedicated for tenants).
  + **standard-cloud**: reverse proxy dedicated engine group (in-cloud).
  + **standard**: reverse proxy dedicated engine group (out-of-cloud).
  + **detector-cloud**: bypass detection dedicated engine group (in-cloud).
  + **detector**: bypass detection dedicated engine group (out-of-cloud).

* `vpc_id` - (Optional, String) Specifies the VPC ID associated with the pool.

* `detail` - (Optional, Bool) Specifies whether to query the detailed information of the pools.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The pool list.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The pool ID.

* `name` - The pool name.

* `region` - The region where the pool is located.

* `type` - The pool type.

* `vpc_id` - The VPC ID associated with the pool.

* `description` - The description of the pool.

* `hosts` - The list of protected domain names associated with the pool.

  The [hosts](#hosts_instances_struct) structure is documented below.

* `instances` - The list of engine instances associated with the pool.

  The [instances](#hosts_instances_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID associated with the pool.

* `create_time` - The creation time of the pool, in milliseconds.

<a name="hosts_instances_struct"></a>
The `hosts`/`instances` block supports:

* `id` - The resource ID.

* `name` - The resource name.

* `service_ip` - The engine instance IP.
