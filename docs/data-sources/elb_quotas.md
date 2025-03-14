---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_quotas"
description: |-
  Use this data source to get the list of quotas of ELB and related resources in a specific project.
---

# huaweicloud_elb_quotas

Use this data source to get the list of quotas of ELB and related resources in a specific project.

## Example Usage

```hcl
data "huaweicloud_elb_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `project_id` - Indicates the project ID.

* `loadbalancer` - Indicates the load balancer quota.
  + If the value is greater than or equal to **0**, it indicates the load balancer quota.
  + If the value is **-1**, the quota is not limited.

* `listener` - Indicates the listener quota.
  + If the value is greater than or equal to **0**, it indicates the listener quota.
  + If the value is **-1**, the quota is not limited.

* `l7policy` - Indicates the forwarding policy quota.
  + If the value is greater than or equal to **0**, it indicates the forwarding policy quota.
  + If the value is **-1**, the quota is not limited.

* `pool` - Indicates the backend server group quota.
  + If the value is greater than or equal to **0**, it indicates the backend server group quota.
  + If the value is **-1**, the quota is not limited.

* `member` - Indicates the backend server quota.
  + If the value is greater than or equal to **0**, it indicates the backend server quota.
  + If the value is **-1**, the quota is not limited.

* `security_policy` - Indicates the custom security policy quota.
  + If the value is greater than or equal to **0**, it indicates the custom security policy quota.
  + If the value is **-1**, the quota is not limited.

* `ipgroup` - Indicates the IP address group quota.
  + If the value is greater than or equal to **0**, it indicates the IP address group quota.
  + If the value is **-1**, the quota is not limited.

* `ipgroup_max_length` - Indicates the maximum number of IP addresses that can be added to an IP address group.
  + If the value is greater than or equal to **0**, it indicates the IP address quota.
  + If the value is **-1**, the quota is not limited.

* `healthmonitor` - Indicates the health check quota.
  + If the value is greater than or equal to **0**, it indicates the health check quota.
  + If the value is **-1**, the quota is not limited.

* `certificate` - Indicates the certificate quota.
  + If the value is greater than or equal to **0**, it indicates the certificate quota.
  + If the value is **-1**, the quota is not limited.

* `ipgroup_bindings` - Indicates the maximum number of listeners that can be associated with an IP address group.
  + If the value is greater than or equal to 0, it indicates the maximum number of listeners that can be associated with
    an IP address group.
  + If the value is **-1**, the quota is not limited.

* `listeners_per_loadbalancer` - Indicates the maximum number of listeners that can be associated with a load balancer.
  + If the value is greater than or equal to **0**, it indicates the current quota.
  + **-1** indicates that the quota is not limited.
  
  -> **NOTE:** The maximum number of listeners that can be added to each load balancer is not limited, but it is recommended
  that the listeners not exceed the default quota.

* `listeners_per_pool` - Indicates the maximum number of listeners that can be associated with a backend server group.
  + If the value is greater than or equal to **0**, it indicates the current quota.
  + **-1** indicates that the quota is not limited.

* `l7policies_per_listener` - Indicates the maximum number of forwarding policies that can be configured for a listener.
  + If the value is greater than or equal to **0**, it indicates the forwarding policy quota.
  + **-1** indicates that the quota is not limited.

* `ipgroups_per_listener` - Indicates the maximum number of IP address groups that can be associated with a listener.
  + If the value is greater than or equal to **0**, it indicates the IP address group quota.
  + **-1** indicates that the quota is not limited.

* `members_per_pool` - Indicates the maximum number of backend servers in a backend server group.
  + If the value is greater than or equal to **0**, it indicates the backend server quota.
  + If the value is **-1**, the quota is not limited.

* `pools_per_l7policy` - Indicates the maximum number of backend server groups that can be used by a forwarding policy.
  + If the value is greater than or equal to **0**, it indicates the backend server group quota.
  + **-1** indicates that the quota is not limited.

* `condition_per_policy` - Indicates the maximum number of forwarding rules per forwarding policy.
  + If the value is greater than or equal to **0**, it indicates the current quota.
  + **-1** indicates that the quota is not limited.
