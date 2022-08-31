package pagination

import (
	"fmt"

	"github.com/chnsz/golangsdk"
)

// the query type that can be supported
const (
	Marker   string = "marker" // limit + marker
	Offset          = "offset" // limit + offset
	PageSize        = "page"   // pagesize + page
)

// QueryOpts is the operation for ListAllItems
// currently, you can change the default marker field with `MarkerField`
type QueryOpts struct {
	MarkerField string
}

// ListAllItems is the method to get all pages from initialURL
func ListAllItems(client *golangsdk.ServiceClient, qType string, initialURL string,
	opts *QueryOpts) (interface{}, error) {

	var queryPage Pager
	switch qType {
	case Marker:
		pageBuilder := func(r PageResult) Page {
			p := MarkerPageBase{PageResult: r}
			if opts != nil && opts.MarkerField != "" {
				p.setMarkerField(opts.MarkerField)
			}
			p.Owner = p
			return p
		}

		queryPage = NewPager(client, initialURL, pageBuilder)
	case Offset:
		pageBuilder := func(r PageResult) Page {
			return OffsetPageBase{PageResult: r}
		}

		queryPage = NewPager(client, initialURL, pageBuilder)
	case PageSize:
		pageBuilder := func(r PageResult) Page {
			return PageSizeBase{PageResult: r}
		}

		queryPage = NewPager(client, initialURL, pageBuilder)
	default:
		err := golangsdk.ErrUnexpectedType{}
		err.Expected = fmt.Sprintf("%s/%s/%s", Marker, Offset, PageSize)
		err.Actual = qType
		return nil, err
	}

	pages, err := queryPage.AllPages()
	if err != nil {
		return nil, err
	}
	return pages.GetBody(), nil
}
