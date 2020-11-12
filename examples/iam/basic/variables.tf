variable "iden_user_name_default" {
  type    = string
  default = "iden_user_default"
}

variable "iden_users" {
    description = "List of identity users."
    type = list(object({
        name        = string
        description = string
    }))
    default = [
        {name = "iden_user_A", description = "This is a identity user."},
        {name = "iden_user_B", description = "This is a identity user."}
    ]
}