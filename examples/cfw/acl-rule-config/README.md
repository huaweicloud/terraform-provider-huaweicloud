# Create ACL rules with address groups, service groups and domain name groups

This example provides best practice code for using Terraform to create ACL rules in HuaweiCloud CFW service,
including IP address groups, service groups, domain name groups and ACL rules with different matching strategies.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A Cloud Firewall instance already created

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CFW firewall is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `fw_instance_id` - The firewall instance ID
* `address_group_name` - The name of the IP address group
* `service_group_name` - The name of the service group
* `domain_name_group_name` - The name of the domain name group
* `acl_rule_ip_name` - The name of the IP-based ACL rule
* `acl_rule_domain_name` - The name of the domain-based ACL rule
* `acl_rule_group_name` - The name of the group-based ACL rule

#### Optional Variables

* `address_group_description` - The description of the IP address group (default: "")
* `service_group_description` - The description of the service group (default: "")
* `domain_name_group_type` - The type of the domain name group (default: 0)
* `domain_name_group_description` - The description of the domain name group (default: "")
* `domain_name_group_domains` - The list of domain names in the domain name group (default: [])
* `acl_rule_ip_description` - The description of the IP-based ACL rule (default: "")
* `acl_rule_type` - The ACL rule type (default: 0)
* `acl_rule_address_type` - The ACL rule address type (default: 0)
* `acl_rule_action_type` - The ACL rule action type (default: 0)
* `acl_rule_long_connect_enable` - Whether to enable persistent connections (default: 0)
* `acl_rule_status` - The ACL rule status (default: 1)
* `acl_rule_applications` - The application list of the ACL rule (default: ["HTTPS"])
* `acl_rule_source_addresses` - The source IP address list (default: ["1.1.1.1"])
* `acl_rule_destination_addresses` - The destination IP address list (default: ["1.1.1.2"])
* `acl_rule_custom_service_protocol` - The protocol type of the custom service (default: 6)
* `acl_rule_custom_service_source_port` - The source port of the custom service (default: "81")
* `acl_rule_custom_service_dest_port` - The destination port of the custom service (default: "82")
* `acl_rule_domain_description` - The description of the domain-based ACL rule (default: "")
* `acl_rule_domain_direction` - The direction of the domain-based ACL rule (default: 1)
* `acl_rule_destination_domain_address_name` - The destination domain address name (default: "*.baidu.com")
* `acl_rule_group_description` - The description of the group-based ACL rule (default: "")
* `acl_rule_service_group_protocol` - The protocol type used by the service group (default: 6)
* `tags` - The key/value pairs to associate with the resources

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  fw_instance_id         = "your_firewall_instance_id"
  address_group_name     = "your_address_group_name"
  service_group_name     = "your_service_group_name"
  domain_name_group_name = "your_domain_name_group_name"
  acl_rule_ip_name       = "your_acl_rule_ip_name"
  acl_rule_domain_name   = "your_acl_rule_domain_name"
  acl_rule_group_name    = "your_acl_rule_group_name"
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
* The `fw_instance_id` can be obtained from the CFW console or the `huaweicloud_cfw_firewalls` data source
* The `object_id` is automatically obtained from the firewall data source
* ACL rules are evaluated in order from top to bottom
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.71.2 |
