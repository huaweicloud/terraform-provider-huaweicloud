# Create a VPN connection example

This example provides best practice code for using Terraform to create a connection in HuaweiCloud VPN service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The name of the VPN connection.
* `gateway_id` -The VPN gateway ID.
* `gateway_ip` - The VPN gateway IP ID.
* `vpn_type` - The connection type.
* `customer_gateway_id` - The customer gateway ID.
* `psk` - The pre-shared key.

#### Optional Variables

* `peer_subnets` - The CIDR list of customer subnets.
* `tunnel_local_address` - The local tunnel address.
* `tunnel_peer_address` - The peer tunnel address.
* `enable_nqa` - Whether to enable NQA check.
* `ikepolicy` - The IKE policy configurations.
  + `authentication_algorithm` - The authentication algorithm.
  + `encryption_algorithm` - The encryption algorithm.
  + `ike_version` - The IKE negotiation version.
  + `lifetime_seconds` - The life cycle of SA in seconds.
  + `local_id_type` - The local ID type.
  + `local_id` - The local ID.
  + `peer_id_type` - The peer ID type.
  + `peer_id` - The peer ID.
  + `phase1_negotiation_mode` - The negotiation mode.
  + `authentication_method` - The authentication method during IKE negotiation.
  + `dh_group` - The DH group used for key exchange in phase 1.
  + `dpd` - The dead peer detection (DPD) object.
    - `timeout` - The interval for retransmitting DPD packets.
    - `interval` - The DPD idle timeout period.
    - `msg` - The format of DPD packets.
* `ipsecpolicy` - The IPsec policy configurations.
  + `authentication_algorithm` - The authentication algorithm.
  + `encryption_algorithm` - The encryption algorithm.
  + `pfs` - The DH key group used by PFS.
  + `lifetime_seconds` - The lifecycle time of Ipsec tunnel in seconds.
  + `transform_protocol` - The transform protocol.
  + `encapsulation_mode` - The encapsulation mod.
* `policy_rules` - The policy rules.
  + `rule_index` - The rule index.
  + `destination` - The list of destination CIDRs.
  + `source` - The source CIDR.
* `tags` - The tags of the VPN connection.
* `ha_role` - The mode of the VPN connection.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpn_gateway_az_flavor          = "tf_test_vpn_connection"
  vpn_gateway_az_attachment_type = "tf_test_vpn_connection"
  vpc_name                       = "tf_test_vpn_connection"
  subnet_name                    = "tf_test_vpn_connection"
  vpn_gateway_name               = "tf_test_vpn_connection"
  vpn_customer_gateway_name      = "tf_test_vpn_connection"
  vpn_customer_gateway_id_value  = "tf_test_vpn_connection"
  vpn_connection_name            = "tf_test_vpn_connection"
  vpn_connection_peer_subnets    = ["tf_test_vpn_connection"]
  vpn_connection_vpn_type        = "tf_test_vpn_connection"
  vpn_connection_psk             = "tf_test_vpn_connection"
  vpn_connection_enable_nqa      = true
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.69.1 |
