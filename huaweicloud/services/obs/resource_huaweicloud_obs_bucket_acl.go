package obs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type GrantType string

const (
	GrantOwner           GrantType = "OWNER"
	GrantPublic          GrantType = "PUBLIC"
	GrantLogDeliveryUser GrantType = "LOG_DELIVERY_USER"
	GrantAccount         GrantType = "ACCOUNT"
)

// @API OBS PUT ?acl
// @API OBS GET ?acl
func ResourceOBSBucketAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOBSBucketAclCreate,
		ReadContext:   resourceOBSBucketAclRead,
		UpdateContext: resourceOBSBucketAclCreate,
		DeleteContext: resourceOBSBucketAclDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the bucket to which to set the acl.`,
			},
			"owner_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bucket owner permission.`,
			},
			"public_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Description: `Specifies the public permission.`,
			},
			"log_delivery_user_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Description: `Specifies the log delivery user permission.`,
			},
			"account_permission": {
				Type:        schema.TypeSet,
				Elem:        accountPermissionSchema(),
				Optional:    true,
				Description: `Specifies the account permissions.`,
			},
		},
	}
}

func permissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_to_bucket": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to bucket. Valid values are **READ** and **WRITE**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ", "WRITE",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func accountPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the account id to authorize. The account id cannot be the bucket owner, 
and must be unique.`,
			},
			"access_to_bucket": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to bucket. Valid values are **READ** and **WRITE**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ", "WRITE",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func buildBucketAccessesFromRawMap(rawMap map[string]interface{}) []string {
	var accesses []string
	if accessArray, ok := rawMap["access_to_bucket"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}
	if accessArray, ok := rawMap["access_to_acl"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}
	return accesses
}

func buildUserTypeGrants(accesses []string, accountID string) []obs.Grant {
	userGrants := make([]obs.Grant, len(accesses))
	for i, access := range accesses {
		userGrants[i] = obs.Grant{
			Permission: obs.PermissionType(access),
			Grantee: obs.Grantee{
				Type: obs.GranteeUser,
				ID:   accountID,
			},
		}
	}
	return userGrants
}

func buildGroupTypeGrants(accesses []string, uriType obs.GroupUriType) []obs.Grant {
	groupGrants := make([]obs.Grant, len(accesses))
	for i, access := range accesses {
		groupGrants[i] = obs.Grant{
			Permission: obs.PermissionType(access),
			Grantee: obs.Grantee{
				Type: obs.GranteeGroup,
				URI:  uriType,
			},
		}
	}
	return groupGrants
}

func findCurrentOwnerGrant(obsClient *obs.ObsClient, d *schema.ResourceData, cfg *config.Config) (*obs.Grant, error) {
	bucket := d.Get("bucket").(string)
	output, err := obsClient.GetBucketAcl(bucket)
	if err != nil {
		return nil, getObsError("Error retrieving OBS bucket current acl", bucket, err)
	}
	for _, grant := range output.Grants {
		grantee := grant.Grantee
		if grantee.Type == obs.GranteeUser && grantee.ID == cfg.DomainID {
			return &grant, nil
		}
	}
	return nil, fmt.Errorf("%s: cannot find owner grant from current grants", bucket)
}

func buildObsBucketAclGrantsParam(obsClient *obs.ObsClient, d *schema.ResourceData,
	cfg *config.Config) ([]obs.Grant, error) {
	var grants []obs.Grant
	if rawArray, ok := d.Get("owner_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			accesses := buildBucketAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The owner permission types: %v.", accesses)
			grants = append(grants, buildUserTypeGrants(accesses, cfg.DomainID)...)
		}
	}

	if len(grants) == 0 {
		// the owner permission is empty, read the current owner permission. make the owner has permissions
		grant, err := findCurrentOwnerGrant(obsClient, d, cfg)
		if err != nil {
			return nil, err
		}
		grants = append(grants, *grant)
	}

	permissions := d.Get("account_permission").(*schema.Set)
	for _, raw := range permissions.List() {
		if rawMap, rawOk := raw.(map[string]interface{}); rawOk {
			accesses := buildBucketAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The account permission types: %v.", accesses)
			grants = append(grants, buildUserTypeGrants(accesses, rawMap["account_id"].(string))...)
		}
	}

	if rawArray, ok := d.Get("public_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			accesses := buildBucketAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The public permission types: %v.", accesses)
			grants = append(grants, buildGroupTypeGrants(accesses, obs.GroupAllUsers)...)
		}
	}

	if rawArray, ok := d.Get("log_delivery_user_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			accesses := buildBucketAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The log delivery user permission types: %v.", accesses)
			grants = append(grants, buildGroupTypeGrants(accesses, obs.GroupLogDelivery)...)
		}
	}
	return grants, nil
}

func resourceOBSBucketAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	grantParam, err := buildObsBucketAclGrantsParam(obsClient, d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	params := &obs.SetBucketAclInput{
		Bucket: bucket,
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = grantParam
	_, err = obsClient.SetBucketAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error creating OBS bucket acl", bucket, err))
	}

	d.SetId(bucket)
	return resourceOBSBucketAclRead(ctx, d, meta)
}

func flattenGrantsByPermissionType(grants []obs.Grant, cfg *config.Config) map[GrantType][]obs.Grant {
	var ownerGrants []obs.Grant
	var accountGrants []obs.Grant
	var logDeliveryUserGrants []obs.Grant
	var publicGrants []obs.Grant

	for _, grant := range grants {
		grantee := grant.Grantee
		if grantee.Type == obs.GranteeUser && grantee.ID == cfg.DomainID {
			// owner grants
			ownerGrants = append(ownerGrants, grant)
			continue
		}
		if grantee.Type == obs.GranteeUser && grantee.ID != cfg.DomainID {
			// account grants
			accountGrants = append(accountGrants, grant)
			continue
		}

		granteeURI := obs.GroupUriType(parseGranteeURI(grantee.URI))
		if grantee.Type == obs.GranteeGroup && granteeURI == obs.GroupLogDelivery {
			// log delivery user grants
			logDeliveryUserGrants = append(logDeliveryUserGrants, grant)
			continue
		}
		if grantee.Type == obs.GranteeGroup && granteeURI == obs.GroupAllUsers {
			// public grants
			publicGrants = append(publicGrants, grant)
		}
	}
	return map[GrantType][]obs.Grant{
		GrantOwner:           ownerGrants,
		GrantAccount:         accountGrants,
		GrantLogDeliveryUser: logDeliveryUserGrants,
		GrantPublic:          publicGrants,
	}
}

// parseGranteeURI use to parse uri. For example: http://acs.amazonaws.com/groups/global/AllUsers
func parseGranteeURI(granteeURI obs.GroupUriType) string {
	uri := string(granteeURI)
	if len(uri) == 0 {
		return ""
	}
	uriArray := strings.Split(uri, "/")
	if len(uriArray) == 0 {
		return ""
	}
	return uriArray[len(uriArray)-1]
}

func flattenAccessesFromGrant(grant obs.Grant) (bucketAccesses []string, aclAccesses []string) {
	switch grant.Permission {
	case obs.PermissionRead:
		bucketAccesses = []string{"READ"}
	case obs.PermissionWrite:
		bucketAccesses = []string{"WRITE"}
	case obs.PermissionReadAcp:
		aclAccesses = []string{"READ_ACP"}
	case obs.PermissionWriteAcp:
		aclAccesses = []string{"WRITE_ACP"}
	case obs.PermissionFullControl:
		// permission type `PermissionFullControl` means domain id has all permission
		bucketAccesses = []string{"READ", "WRITE"}
		aclAccesses = []string{"READ_ACP", "WRITE_ACP"}
	default:
		log.Printf("[WARN] The grant permission: %s not support.", grant.Permission)
	}
	return
}

func flattenPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	var accessToBucket []string
	var accessToAcl []string
	for _, grant := range grants {
		bucketAccesses, aclAccesses := flattenAccessesFromGrant(grant)
		accessToBucket = append(accessToBucket, bucketAccesses...)
		accessToAcl = append(accessToAcl, aclAccesses...)
	}
	if len(accessToBucket) == 0 && len(accessToAcl) == 0 {
		return nil
	}
	ownerPermissionMap := map[string]interface{}{
		"access_to_bucket": accessToBucket,
		"access_to_acl":    accessToAcl,
	}
	return []map[string]interface{}{ownerPermissionMap}
}

func flattenAccountPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	// accountIDSet stores accountID set
	accountIDSet := make(map[string]bool)
	// accessToBucketMap stores bucket accesses, the key is accountID
	accessToBucketMap := make(map[string][]string)
	// accessToAclMap stores acl accesses, the key is accountID
	accessToAclMap := make(map[string][]string)
	for _, grant := range grants {
		accountID := grant.Grantee.ID
		// append new accesses.
		bucketAccesses, aclAccesses := flattenAccessesFromGrant(grant)
		accessToBucketMap[accountID] = append(accessToBucketMap[accountID], bucketAccesses...)
		accessToAclMap[accountID] = append(accessToAclMap[accountID], aclAccesses...)
		accountIDSet[accountID] = true
	}

	var m []map[string]interface{}
	for id := range accountIDSet {
		m = append(m, map[string]interface{}{
			"access_to_bucket": accessToBucketMap[id],
			"access_to_acl":    accessToAclMap[id],
			"account_id":       id,
		})
	}
	return m
}

func resourceOBSBucketAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketAcl(d.Id())
	if err != nil {
		return diag.FromErr(getObsError("Error retrieving OBS bucket acl", d.Id(), err))
	}

	permissionTypeMap := flattenGrantsByPermissionType(output.Grants, cfg)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", d.Id()),
		d.Set("owner_permission", flattenPermission(permissionTypeMap[GrantOwner])),
		d.Set("public_permission", flattenPermission(permissionTypeMap[GrantPublic])),
		d.Set("log_delivery_user_permission", flattenPermission(permissionTypeMap[GrantLogDeliveryUser])),
		d.Set("account_permission", flattenAccountPermission(permissionTypeMap[GrantAccount])),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket acl fields: %s", err)
	}
	return nil
}

func resourceOBSBucketAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	// the operation to delete acl is to set the owner full permissions
	params := &obs.SetBucketAclInput{
		Bucket: d.Id(),
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = []obs.Grant{
		{
			Grantee: obs.Grantee{
				Type: obs.GranteeUser,
				ID:   cfg.DomainID,
			},
			Permission: obs.PermissionFullControl,
		},
	}
	_, err = obsClient.SetBucketAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting OBS bucket acl", d.Id(), err))
	}
	return nil
}
