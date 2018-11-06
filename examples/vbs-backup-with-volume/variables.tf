variable "user_name" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_username"
  description = "The Username to login with."
}

variable "domain_name" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_domainname"
  description = "The Name of the Domain to scope to (Identity v3)."
}

variable "tenant_name" {
  default = "cn-north-1"
  description = "The Name of the Tenant (Identity v2) or Project (Identity v3) to login with."
}

variable "region" {
  default = "cn-north-1"
  description = "The region of the HuaweiCloud cloud to use."
}

variable "password" {
  # If you don't fill this in, you will be prompted for it
  #default = "your_password'
  description = "The Password to login with."
}

variable "auth_url" {
  default = "https://iam.cn-north-1.myhwclouds.com"
  description = "The Identity authentication URL."
}

### VM (Instance) Settings
variable "instance_count" {
  default = "1"
}

variable "disk_size_gb" {
  default = "0"
}

### Project Settings
variable "project" {
  default = "terraform"
}