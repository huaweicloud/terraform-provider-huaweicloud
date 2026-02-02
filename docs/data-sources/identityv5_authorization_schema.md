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
data "huaweicloud_identityv5_authorization_schema" "test" {
  service_code = "iam"
}
```

## Argument Reference

The following arguments are supported:

* `service_code` - (Required, String) Specifies the service name abbreviation to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `actions` - The list of authorization items supported by the cloud service.  
  The [actions](#v5_authorization_schema_actions) structure is documented below.

* `conditions` - The list of condition keys supported by the cloud service.  
  The [conditions](#v5_authorization_schema_conditions) structure is documented below.

* `operations` - The list of operations supported by the cloud service.  
  The [operations](#v5_authorization_schema_operations) structure is documented below.

* `resources` - The list of resources supported by the cloud service.  
  The [resources](#v5_authorization_schema_resources) structure is documented below.

* `version` - The version number of the service authorization summary.

<a name="v5_authorization_schema_actions"></a>
The `actions` block supports:

* `name` - The authorization item name.

* `permission_only` - Whether the authorization item is only used as a permission point and does not
  correspond to any operation.

* `resources` - The list of resources associated with the authorization item, used to define resource-level
  permissions for the authorization item.  
  The [resources](#v5_authorization_schema_actions_resources) structure is documented below.

* `access_level` - The access level granted when using this authorization item in a policy.

* `aliases` - The list of authorization item aliases.  
  It is used to accommodate scenarios where authorization items are renamed or split into new authorization items.

* `condition_keys` - The service custom conditional attribute list and some global attribute list supported by the
  authorization items and are independent of the resources.

* `description` - The description of the authorization item.  
  The [description](#v5_authorization_schema_description) structure is documented below.

<a name="v5_authorization_schema_conditions"></a>
The `conditions` block supports:

* `description` - The description of the condition key.  
  The [description](#v5_authorization_schema_description) structure is documented below.

* `key` - The condition key name.

* `multi_valued` - Whether the condition value is multi-valued.

* `value_type` - The data type of the condition value.

<a name="v5_authorization_schema_operations"></a>
The `operations` block supports:

* `dependent_actions` - The other authorization item list that this operation may require.

* `operation_action` - The action of the operation.

* `operation_id` - The OpenAPI operation identifier.

<a name="v5_authorization_schema_resources"></a>
The `resources` block supports:

* `urn_template` - The uniform resource name template for the resource.

* `type_name` - The type name of the resource.

<a name="v5_authorization_schema_actions_resources"></a>
The `resources` block supports:

* `condition_keys` - The service custom conditional attribute list and some global attribute list supported
  by the authorization item and resource.  
  It only takes effect when both the authorization item and resource match.

* `required` - Whether the resource type is mandatory for this authorization item.  
  It means the authorization item definitely involves operations on this type of resource.

* `urn_template` - The uniform resource name template for the resource.

<a name="v5_authorization_schema_description"></a>
The `description` block supports:

* `en_us` - The English description.

* `zh_cn` - The Chinese description.
