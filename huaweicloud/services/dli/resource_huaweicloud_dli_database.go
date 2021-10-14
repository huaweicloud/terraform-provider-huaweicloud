package dli

import (
	"context"
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/databases"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceDliSqlDatabaseV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDliSqlDatabaseV1Create,
		ReadContext:   ResourceDliSqlDatabaseV1Read,
		UpdateContext: ResourceDliSqlDatabaseV1Update,
		DeleteContext: ResourceDliSqlDatabaseV1Delete,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][\w_]{0,127}$`),
						"The name consists of 1 to 128 characters, starting with a letter or digit. "+
							"Only letters, digits and underscores (_) are allowed."),
					validation.StringMatch(regexp.MustCompile(`[A-Za-z_]`), "The name cannot be all digits."),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceDliSqlDatabaseV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	opts := databases.CreateOpts{
		Name:                dbName,
		Description:         d.Get("description").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, config),
	}
	_, err = databases.Create(c, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI database")
	}
	d.SetId(dbName)

	return ResourceDliSqlDatabaseV1Read(ctx, d, meta)
}

func setDliSqlDatabaseV1Parameters(d *schema.ResourceData, resp databases.Database) error {
	mErr := multierror.Append(nil,
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("owner", resp.Owner),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func GetDliSqlDatabaseByName(c *golangsdk.ServiceClient, dbName string) (databases.Database, error) {
	resp, err := databases.List(c, databases.ListOpts{
		Keyword: dbName, // Fuzzy matching.
	})
	if err != nil {
		return databases.Database{}, fmtp.Errorf("Error getting database: %s", err)
	}

	if len(resp.Databases) < 1 {
		return databases.Database{}, fmtp.Errorf("Unable to find the specified database (%s): %s", dbName, err)
	}
	for _, db := range resp.Databases {
		if db.Name == dbName {
			return db, nil
		}
	}

	return databases.Database{}, fmtp.Errorf("Only find some databases with this character (%s) in the name", dbName)
}

func ResourceDliSqlDatabaseV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	db, err := GetDliSqlDatabaseByName(c, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Error getting SQL database: %s", err)
	}

	err = setDliSqlDatabaseV1Parameters(d, db)
	if err != nil {
		return fmtp.DiagErrorf("An error occurred during resource parameter setting for SQL database: %s", err)
	}
	return nil
}

func ResourceDliSqlDatabaseV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}

	_, err = databases.UpdateDBOwner(c, d.Id(), databases.UpdateDBOwnerOpts{
		NewOwner: d.Get("owner").(string),
	})
	if err != nil {
		return fmtp.DiagErrorf("Error updating SQL database owner: %s", err)
	}

	return ResourceDliSqlDatabaseV1Read(ctx, d, meta)
}

func ResourceDliSqlDatabaseV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud DLI v1 client: %s", err)
	}
	err = databases.Delete(c, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting SQL database: %s", err)
	}
	d.SetId("")
	return nil
}
