---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_authorization_schema"
description: |-
  Use this data source to get the list of IAM authorization schemas within HuaweiCloud.
---

# huaweicloud_identityv5_authorization_schema

Use this data source to get the list of IAM authorization schemas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_authorization_schema" "schema" {
  service_code = "iam"
}
```

## Argument Reference

The following arguments are supported:

* `service_code` - (Required, String) Specifies the service name abbreviation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `actions` - Indicates the list of authorization items supported by the cloud service.
  The [actions](#Authorization_action) structure is documented below.

* `conditions` - Indicates the list of condition keys supported by the cloud service.
  The [conditions](#Authorization_condition) structure is documented below.

* `operations` - Indicates the list of operations supported by the cloud service.
  The [operations](#Authorization_operation) structure is documented below.

* `resources` - Indicates the list of resources supported by the cloud service.
  The [resources](#Authorization_op_resource) structure is documented below.

* `version` - Indicates version number of the service authorization summary.

<a name="Authorization_action"></a>
The `actions` block supports:

* `name` - Indicates the authorization item name.

* `permission_only` - Indicates whether the authorization item is only used as a permission point and does not
  correspond to any operation.

* `resources` - Indicates the list of resources associated with the authorization item, used to define resource-level
  permissions for the authorization item.
  The [resources](#Authorization_resource) structure is documented below.

* `condition_keys` - Indicates service custom condition attributes and some global attributes for the authorization
  item and resource, which only take effect when both the authorization item and resource match.

* `access_level` - Indicates access level granted when using this authorization item in a policy.

* `aliases` - Indicates list of authorization item aliases, used to accommodate scenarios where authorization items
  are renamed or split into new authorization items.

* `condition_keys` - Indicates service custom condition attributes and some global attributes supported by the
  authorization item, which are unrelated to resources.

* `description` - Indicates description of the authorization item.
  The [description](#Authorization_description) structure is documented below.

<a name="Authorization_condition"></a>
The `conditions` block supports:

* `description` - Indicates description of the condition key.
  The [description](#Authorization_condition_description) structure is documented below.

* `key` - Indicates the condition key name.

* `multi_valued` - Indicates whether the condition value is multi-valued.

* `value_type` - Indicates data type of the condition value.

<a name="Authorization_operation"></a>
The `operations` block supports:

* `dependent_actions` - Indicates other authorization items that this operation may require.

* `operation_action` - Indicates the action of the operation.

* `operation_id` - Indicates OpenAPI operation identifier.

<a name="Authorization_op_resource"></a>
The `resources` block supports:

* `urn_template` - Indicates the uniform resource name template for the resource.

* `type_name` - Indicates the type name of the resource.

<a name="Authorization_resource"></a>
The `resources` block supports:

* `required` - Indicates identifies whether the resource type is mandatory for this authorization item,
  meaning the authorization item definitely involves operations on this type of resource.

* `urn_template` - Indicates the uniform resource name template for the resource.

<a name="Authorization_description"></a>
The `description` block supports:

* `en_us` - Indicates English description.

* `zh_cn` - Indicates Chinese description.

<a name="Authorization_condition_description"></a>
The `description` block supports:

* `en_us` - Indicates English description.

* `zh_cn` - Indicates Chinese description.
