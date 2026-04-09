---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_exception_rules"
description: |-
  Use this data source to query the exception rules under the DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_exception_rules

Use this data source to query the exception rules under the DWS cluster within HuaweiCloud.

## Example Usage

### Query all exception rules in the cluster

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_exception_rules" "test" {
  cluster_id = var.cluster_id
}
```

### Query the exception rules with a fuzzy word in the cluster

```hcl
variable "cluster_id" {}
variable "fuzzy_name_word" {}

data "huaweicloud_dws_cluster_exception_rules" "test" {
  cluster_id = var.cluster_id
  rule_name  = var.fuzzy_name_word
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the exception rules are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster to which the exception rules belong.

* `rule_name` - (Optional, String) Specifies the name of the exception rule to be queried, which is the fuzzy query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - The list of the exception rules that matched filter parameters.  
  The [rules](#dws_cluster_exception_rules) structure is documented below.

<a name="dws_cluster_exception_rules"></a>
The `rules` block supports:

* `name` - The name of the exception rule.

* `configurations` - The configuration items of the exception rule.  
  The type of configuration attribute is `map`.
