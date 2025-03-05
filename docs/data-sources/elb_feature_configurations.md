---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_feature_configurations"
description: |-
  Use this data source to get the list of feature configurations of ELB of a tenant.
---

# huaweicloud_elb_feature_configurations

Use this data source to get the list of feature configurations of ELB of a tenant.

## Example Usage

```hcl
data "huaweicloud_elb_feature_configurations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `feature` - (Optional, String) Specifies the feature name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configs` - Indicates the feature configuration list.

  The [configs](#configs_struct) structure is documented below.

<a name="configs_struct"></a>
The `configs` block supports:

* `id` - Indicates the ID of the configuration.

* `feature` - Indicates the feature name.

* `type` - Indicates the type of the feature configuration value.

* `value` - Indicates the feature configuration value.
  For example, the value **true** or **false** indicates that the feature is enabled or disabled.
  The feature value of the quota is an integer, indicating that the quota is limited.

* `switch` - Indicates whether to enable feature configuration.
  The value can be:
  + **true**: The feature configuration has taken effect.
  + **false**: The feature configuration does not take effect.

* `service` - Indicates the service. The value is fixed at **ELB**.

* `description` - Indicates the feature configuration description.

* `caller` - Indicates the configuration creator.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.
