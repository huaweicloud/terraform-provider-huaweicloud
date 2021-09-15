package thesaurus

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// LoadCustomThesaurusReq This is a auto create Body Object
type LoadThesaurusReq struct {
	// OBS bucket where word dictionary files are stored. The bucket must be a standard storage or infrequently
	// accessed storage and cannot be the archived storage.
	BucketName string `json:"bucketName" required:"true"`
	// Main word dictionary file object, which must be a text file encoded in UTF-8 without BOM. Each line contains
	// one sub-word. The maximum file size is 100 MB.
	MainObject string `json:"mainObject,omitempty"`
	// Stop word dictionary file object, which must be a text file encoded in UTF-8 without BOM. Each line contains
	// one sub-word. The maximum file size is 20 MB.
	StopObject string `json:"stopObject,omitempty"`
	// Synonym word dictionary file, which must be a text file encoded in UTF-8 without BOM. Each line contains one
	// group of sub-words. The maximum file size is 20 MB.
	SynonymObject string `json:"synonymObject,omitempty"`
}

func Get(c *golangsdk.ServiceClient, clusterId string) (*ThesaurusStatusResp, error) {
	var rst ThesaurusStatusResp
	_, err := c.Get(getURL(c, clusterId), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		return &rst, nil
	}
	return nil, err
}

// LoadIKThesaurus
func Load(c *golangsdk.ServiceClient, clusterId string, opts LoadThesaurusReq) *golangsdk.ErrResult {
	var r golangsdk.ErrResult

	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return &r
	}

	_, r.Err = c.Post(loadURL(c, clusterId), b, nil, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}

// DeleteIKThesaurus
func Delete(c *golangsdk.ServiceClient, clusterId string) *golangsdk.ErrResult {
	var r golangsdk.ErrResult
	_, r.Err = c.Delete(deleteURL(c, clusterId), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &r
}
