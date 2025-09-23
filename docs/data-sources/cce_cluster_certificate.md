---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_certificate"
description: ""
---

# huaweicloud_cce_cluster_certificate

Use this data source to get the certificate of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_cluster_certificate" "test" {
  cluster_id = var.cluster_id
  duration   = 30
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster certificate. If omitted, the
  provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID which the cluster certificate in.

* `duration` - (Required, Int) Specifies the duration of the cluster certificate. The unit is days. The valid value in
  [1, 1827]. If the input value is -1, it will use the maximum 1827 as `duration` value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.

* `current_context` - The current context of the cluster certificate.

* `clusters` - The clusters information of the cluster certificate.
  The [clusters](#CCECluster_clusters) structure is documented below.

* `users` - The users information of cluster the certificate.
  The [users](#CCECluster_users) structure is documented below.

* `contexts` - The contexts information of the cluster certificate.
  The [contexts](#CCECluster_contexts) structure is documented below.

<a name="CCECluster_clusters"></a>
The `clusters` block supports:

* `name` - The cluster name of the cluster certificate.

* `server` - The server address of the cluster certificate.

* `certificate_authority_data` - The certificate authority data of the cluster certificate.

* `insecure_skip_tls_verify` - Whether insecure skip tls verify of the cluster certificate.

<a name="CCECluster_users"></a>
The `users` block supports:

* `name` - The user name of the cluster certificate. The value is fixed to `user`.

* `client_certificate_data` - The client certificate data of the cluster certificate.

* `client_key_data` - The client key data of the cluster certificate.

<a name="CCECluster_contexts"></a>
The `contexts` block supports:

* `name` - The context name of the cluster certificate.

* `cluster` - The context cluster of the cluster certificate.

* `user` - The context user of the cluster certificate.
