---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_resource_infos"
description: |-
  Use this data source to get the resource distribution information.
---

# huaweicloud_cpcs_resource_infos

Use this data source to get the resource distribution information.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
data "huaweicloud_cpcs_resource_infos" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cloud_service` - The istribution of dcloud-native cryptographic service resources.
  The [cloud_service](#resource_infos_cloud_service) structure is documented below.

* `ccsp_service` - The distribution of password service resources.
  The [ccsp_service](#resource_infos_ccsp_service) structure is documented below.

* `vsm` - The distribution of cloud password virtual machine resources.
  The [vsm](#resource_infos_vsm) structure is documented below.

* `app` - The application information.
  The [app](#resource_infos_application) structure is documented below.

* `kms` - The distribution of KMS key resources.
  The [kms](#resource_infos_kms) structure is documented below.

<a name="resource_infos_cloud_service"></a>
The `cloud_service` block supports:

* `service_num` - The current number of cloud-native password services enabled.

* `resource_num` - The total number of resource instances of cloud-native cryptographic services.

* `resource_distribution` - The distribution of cloud-native cryptographic service resource instances.
  The [resource_distribution](#resource_distribution_struct) structure is documented below.

<a name="resource_infos_ccsp_service"></a>
The `ccsp_service` block supports:

* `cluster_num` - The number of password service clusters.

* `instance_num` - The number of password service instances.

* `instance_quota` - The quota number of password service instances that can be created.

* `instance_distribution` - The distribution of password service instances by service type.
  The [instance_distribution](#instance_distribution_struct) structure is documented below.

<a name="resource_infos_vsm"></a>
The `vsm` block supports:

* `cluster_num` - The number of VSM clusters.

* `cpcs_cluster_num` - The number of VSM clusters created and managed by CPCS.

* `instance_num` - The total number of VSM instances owned.

* `cpcs_instance_num` - The number of VSM instances created and managed by CPCS.

* `instance_quota` - The number of VSM instance quotas allocated.

<a name="resource_infos_application"></a>
The `app` block supports:

* `app_num` - The number of simple applications created in CPCS.

<a name="resource_infos_kms"></a>
The `kms` block supports:

* `total_num` - The number of KMS keys.

* `result` - The quantity distribution of KMS key types.
  The [result](#result_struct) structure is documented below.

<a name="resource_distribution_struct"></a>
The `resource_distribution` block supports:

* `kms` - The number of keys of the KMS service.

<a name="instance_distribution_struct"></a>
The `instance_distribution` block supports:

* `encrypt_decrypt` - The number of encryption and decryption service instances.

* `sign_verify` - The number of signature verification service instances.

* `kms` - The number of KMS service instances.

* `timestamp` - The number of timestamp service instances.

* `colla_sign` - The number of instances of collaborative signature services.

* `otp` - The number of dynamic password service instances.

* `db_encrypt` - The number of database encryption service instances.

* `file_encrypt` - The number of file encryption service instances.

* `digit_seal` - The number of instances of electronic signature and seal services.

* `ssl_vpn` - The number of SSL and VPN service instances.

<a name="result_struct"></a>
The `result` block supports:

* `aes_256` - The number of keys for the AES_256 algorithm.

* `sm4` - The number of keys for the SM4 algorithm.

* `rsa_2048` - The number of keys for the RSA_2048 algorithm.

* `rsa_3072` - The number of keys for the RSA_3072 algorithm.

* `rsa_4096` - The number of keys for the RSA_4096 algorithm.

* `ec_p256` - The number of keys for the EC_P256 algorithm.

* `ec_p384` - The number of keys for the EC_P384 algorithm.

* `sm2` - The number of keys for the SM2 algorithm.
