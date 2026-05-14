---
subcategory: "CCE"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_pod_identity_associations"
description: |-
  Use this data source to get the pod identity associations of a specified CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_pod_identity_associations

Use this data source to get the pod identity associations of a specified CCE cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_pod_identity_associations" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the pod identity associations.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `associations` - The list of pod identity associations.
  The [associations](#cce_pod_identity_associations_attr) structure is documented below.

<a name="cce_pod_identity_associations_attr"></a>
The `associations` block supports:

* `uid` - The UID of the pod identity association.

* `cluster_id` - The ID of the cluster that the pod identity association belongs to.

* `namespace` - The namespace of the service account associated with a pod identity association.

* `service_account` - The name of the service account associated with a pod identity association.

* `agency_name` - The name of the agency associated with a pod identity association.

* `created_at` - The time when a pod identity association was created.

* `updated_at` - The time when a pod identity association was last updated.

* `tags` - The key/value pairs to associate with the pod identity association.
