---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_feature_configurations"
description: |-
  Use this data source to get the list of feature configurations of a load balancer.
---

# huaweicloud_elb_loadbalancer_feature_configurations

Use this data source to get the list of feature configurations of a load balancer.

## Example Usage

```hcl
variable "loadbalancer_id" {}

data "huaweicloud_elb_loadbalancer_feature_configurations" "test" {
  loadbalancer_id = var.loadbalancer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `loadbalancer_id` - (Required, String) Specifies the load balancer ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `features` - Specifies the load balancer feature information list.

  The [features](#features_struct) structure is documented below.

<a name="features_struct"></a>
The `features` block supports:

* `feature` - Specifies the feature name.

* `type` - Specifies the type of the feature configuration value.

* `value` - Specifies the feature value.
  For example, the value **true** or **false** indicates that the feature is enabled or disabled.
  The feature value of the quota is an integer, indicating that the quota is limited.
