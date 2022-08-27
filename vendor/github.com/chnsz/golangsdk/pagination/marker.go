package pagination

import (
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk"
)

// MarkerPage is a stricter Page interface that describes additional functionality required for use with NewMarkerPager.
// For convenience, embed the MarkedPageBase struct.
type MarkerPage interface {
	Page

	// LastMarker returns the last "marker" value on this page.
	LastMarker() (string, error)
}

// MarkerPageBase is a page in a collection that's paginated by "limit" and "marker" query parameters.
type MarkerPageBase struct {
	PageResult

	// Owner is a reference to the embedding struct.
	Owner MarkerPage
	// markerField is the field/key to get the next marker, the default value is "id"
	markerField string
}

// NextPageURL generates the URL for the page of results after this one.
func (current MarkerPageBase) NextPageURL() (string, error) {
	currentURL := current.URL

	mark, err := current.Owner.LastMarker()
	if err != nil {
		return "", err
	}

	if mark == "" {
		return "", nil
	}

	q := currentURL.Query()
	q.Set("marker", mark)
	currentURL.RawQuery = q.Encode()

	return currentURL.String(), nil
}

// IsEmpty satisifies the IsEmpty method of the Page interface
func (current MarkerPageBase) IsEmpty() (bool, error) {
	if pb, ok := current.Body.(map[string]interface{}); ok {
		for k, v := range pb {
			// ignore xxx_links
			if !strings.HasSuffix(k, "links") {
				// check the field's type. we only want []interface{} (which is really []map[string]interface{})
				switch vt := v.(type) {
				case []interface{}:
					return len(vt) == 0, nil
				}
			}
		}
	}

	if b, ok := current.Body.([]interface{}); ok {
		return len(b) == 0, nil
	}

	err := golangsdk.ErrUnexpectedType{}
	err.Expected = "map[string]interface{}/[]interface{}"
	err.Actual = fmt.Sprintf("%T", current.Body)
	return true, err
}

// GetBody returns the linked page's body. This method is needed to satisfy the
// Page interface.
func (current MarkerPageBase) GetBody() interface{} {
	return current.Body
}

// LastMarker method returns the last ID in a page.
func (current MarkerPageBase) LastMarker() (string, error) {
	var pageItems []interface{}

	switch pb := current.Body.(type) {
	case map[string]interface{}:
		for k, v := range pb {
			// ignore xxx_links
			if !strings.HasSuffix(k, "links") {
				// check the field's type. we only want []interface{} (which is really []map[string]interface{})
				switch vt := v.(type) {
				case []interface{}:
					pageItems = vt
					break
				}
			}
		}
	case []interface{}:
		pageItems = pb
	default:
		err := golangsdk.ErrUnexpectedType{}
		err.Expected = "map[string]interface{}/[]interface{}"
		err.Actual = fmt.Sprintf("%T", pb)
		return "", err
	}

	if len(pageItems) == 0 {
		return "", nil
	}

	lastItem := pageItems[len(pageItems)-1]
	field := current.getMarkerField()
	lastID, err := searchField(lastItem, field)
	if err != nil {
		return "", err
	}

	// check the marker field
	if lastID != nil {
		if id, ok := lastID.(string); ok && id != "" {
			return id, nil
		}
	}

	return "", fmt.Errorf("can not find '%s' in item", field)
}

func (current *MarkerPageBase) getMarkerField() string {
	if current.markerField == "" {
		return "id"
	}
	return current.markerField
}

func (current *MarkerPageBase) setMarkerField(field string) {
	if field != "" {
		current.markerField = field
	}
	return
}

func searchField(root interface{}, field string) (interface{}, error) {
	keys := strings.Split(field, ".")
	current := root
	for _, k := range keys {
		if currentMap, ok := current.(map[string]interface{}); ok {
			if v, exist := currentMap[k]; exist {
				current = v
			} else {
				return nil, fmt.Errorf("can not find '%s' in item", field)
			}
		} else {
			return nil, fmt.Errorf("item is not a map but %T", current)
		}
	}

	return current, nil
}
