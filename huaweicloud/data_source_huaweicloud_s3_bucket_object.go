package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceS3BucketObject() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceS3BucketObjectRead,
		DeprecationMessage: "use huaweicloud obs resource instead",

		Schema: map[string]*schema.Schema{
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cache_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_disposition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sse_kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"website_redirect_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceS3BucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	input := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	if v, ok := d.GetOk("range"); ok {
		input.Range = aws.String(v.(string))
	}
	if v, ok := d.GetOk("version_id"); ok {
		input.VersionId = aws.String(v.(string))
	}

	versionText := ""
	uniqueId := bucket + "/" + key
	if v, ok := d.GetOk("version_id"); ok {
		versionText = fmt.Sprintf(" of version %q", v.(string))
		uniqueId += "@" + v.(string)
	}

	log.Printf("[DEBUG] Reading S3 object: %s", input)
	out, err := conn.HeadObject(&input)
	if err != nil {
		return fmt.Errorf("Failed getting S3 object: %s Bucket: %q Object: %q", err, bucket, key)
	}
	if out.DeleteMarker != nil && *out.DeleteMarker == true {
		return fmt.Errorf("Requested S3 object %q%s has been deleted",
			bucket+key, versionText)
	}

	log.Printf("[DEBUG] Received S3 object: %s", out)

	d.SetId(uniqueId)

	d.Set("cache_control", out.CacheControl)
	d.Set("content_disposition", out.ContentDisposition)
	d.Set("content_encoding", out.ContentEncoding)
	d.Set("content_language", out.ContentLanguage)
	d.Set("content_length", out.ContentLength)
	d.Set("content_type", out.ContentType)
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(*out.ETag, `"`))
	d.Set("expiration", out.Expiration)
	d.Set("expires", out.Expires)
	d.Set("last_modified", out.LastModified.Format(time.RFC1123))
	if err := d.Set("metadata", pointersMapToStringList(out.Metadata)); err != nil {
		return fmt.Errorf("[DEBUG] Error saving metadata to state for HuaweiCloud S3 object (%s): %s", d.Id(), err)
	}
	d.Set("server_side_encryption", out.ServerSideEncryption)
	d.Set("sse_kms_key_id", out.SSEKMSKeyId)
	d.Set("version_id", out.VersionId)
	d.Set("website_redirect_location", out.WebsiteRedirectLocation)

	if isContentTypeAllowed(out.ContentType) {
		input := s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
		if v, ok := d.GetOk("range"); ok {
			input.Range = aws.String(v.(string))
		}
		if out.VersionId != nil {
			input.VersionId = out.VersionId
		}
		out, err := conn.GetObject(&input)
		if err != nil {
			return fmt.Errorf("Failed getting S3 object: %s", err)
		}

		buf := new(bytes.Buffer)
		bytesRead, err := buf.ReadFrom(out.Body)
		if err != nil {
			return fmt.Errorf("Failed reading content of S3 object (%s): %s",
				uniqueId, err)
		}
		log.Printf("[INFO] Saving %d bytes from S3 object %s", bytesRead, uniqueId)
		d.Set("body", buf.String())
	} else {
		contentType := ""
		if out.ContentType == nil {
			contentType = "<EMPTY>"
		} else {
			contentType = *out.ContentType
		}

		log.Printf("[INFO] Ignoring body of S3 object %s with Content-Type %q",
			uniqueId, contentType)
	}

	return nil
}

// This is to prevent potential issues w/ binary files
// and generally unprintable characters
// See https://github.com/hashicorp/terraform/pull/3858#issuecomment-156856738
func isContentTypeAllowed(contentType *string) bool {
	if contentType == nil {
		return false
	}

	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^text/.+"),
		regexp.MustCompile("^application/json$"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(*contentType) {
			return true
		}
	}

	return false
}
