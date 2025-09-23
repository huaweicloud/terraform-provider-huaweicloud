---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_cluster_certificate"
description: |-
  Use this data source to get the certificate of a CCE Autopilot cluster within HuaweiCloud.
---

# huaweicloud_cce_autopilot_cluster_certificate

Use this data source to get the certificate of a CCE Autopilot cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_autopilot_cluster_certificate" "test" {
  cluster_id = var.cluster_id
  duration   = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID to which the cluster certificate belongs.

* `duration` - (Required, Int) Specifies the duration of the cluster certificate.
  The unit is days. The valid value in [1, 1825]. If the input value is -1,
  it will use the maximum 1825 as `duration` value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.

* `clusters` - The clusters information of the cluster certificate.

  The [clusters](#clusters_struct) structure is documented below.

* `users` - The users information of cluster the certificate.

  The [users](#users_struct) structure is documented below.

* `contexts` - The contexts information of the cluster certificate.

  The [contexts](#contexts_struct) structure is documented below.

* `current_context` - The current context of the cluster certificate.

<a name="clusters_struct"></a>
The `clusters` block supports:

* `name` - The cluster name of the cluster certificate.

* `cluster` - The cluster information.

  The [cluster](#clusters_cluster_struct) structure is documented below.

<a name="clusters_cluster_struct"></a>
The `cluster` block supports:

* `server` - The server address of the cluster certificate.

* `certificate_authority_data` - The certificate authority data of the cluster certificate.

* `insecure_skip_tls_verify` - Whether insecure skip tls verify of the cluster certificate.

<a name="users_struct"></a>
The `users` block supports:

* `name` - The user name of the cluster certificate.
  The value is fixed to `user`.

* `user` - The user information.

  The [user](#users_user_struct) structure is documented below.

<a name="users_user_struct"></a>
The `user` block supports:

* `client_certificate_data` - The client certificate data of the cluster certificate.

* `client_key_data` - The client key data of the cluster certificate.

<a name="contexts_struct"></a>
The `contexts` block supports:

* `name` - The context name of the cluster certificate.

* `context` - The user information.

  The [context](#contexts_context_struct) structure is documented below.

<a name="contexts_context_struct"></a>
The `context` block supports:

* `cluster` - The context cluster of the cluster certificate.

* `user` - The context user of the cluster certificate.
