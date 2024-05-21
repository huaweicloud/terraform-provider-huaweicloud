---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_authorizations"
description: |-
  Use this data source to get the list of CC authorizations.
---

# huaweicloud_cc_authorizations

Use this data source to get the list of CC authorizations.

## Example Usage

```hcl
variable authorization_id {}

data "huaweicloud_cc_authorizations" "test" {
  authorization_id = var.authorization_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `authorization_id` - (Optional, String) Specifies the ID of the cross-account authorization.

* `name` - (Optional, String) Specifies the name of the cross-account authorization.

* `description` - (Optional, String) Specifies the description of the cross-account authorization.

* `cloud_connection_domain_id` - (Optional, String) Specifies the account ID that the cloud connection belongs to.

* `cloud_connection_id` - (Optional, String) Specifies the cloud connection ID.

* `instance_id` - (Optional, String) Specifies the network instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorizations` - The authorized instances.

  The [authorizations](#authorizations_struct) structure is documented below.

<a name="authorizations_struct"></a>
The `authorizations` block supports:

* `id` - The ID of the cross-account authorization.

* `name` - The name of the cross-account authorization.

* `description` - The description of the cross-account authorization.

* `status` - The authorization status.

* `instance_type` - The type of an authorized network instance.

* `cloud_connection_domain_id` - The account ID that the cloud connection belongs to.

* `instance_id` - The network instance ID.

* `cloud_connection_id` - The cloud connection ID.

* `domain_id` - The ID of the account that the network instance belongs to.

* `project_id` - The project ID of the network instance.

* `region_id` - The region ID of the network instance.

* `created_at` - Time when the resource was created.

* `updated_at` - Time when the resource was updated.
