---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_openid_jwks"
description: |-
  Use this data source to get the openID JWKS of a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_cluster_openid_jwks

Use this data source to get the openID JWKS of a CCE cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_cluster_openid_jwks" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster certificate. If omitted, the
  provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID which the cluster certificate in.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `keys` - The list of keys.
  The [keys](#CCECluster_jwks_keys) structure is documented below.

<a name="CCECluster_jwks_keys"></a>
The `keys` block supports:

* `use` - How the key was meant to be used; sig represents the signature.

* `kty` - The family of cryptographic algorithms used with the key.

* `kid` - The unique identifier for the key.

* `alg` - The specific cryptographic algorithm used with the key.

* `n` - The modulus for the RSA public key.

* `e` - The exponent for the RSA public key.
