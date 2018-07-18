package stacks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// CreatedStack represents the object extracted from a Create operation.
type CreatedStack struct {
	ID    string           `json:"id"`
	Links []golangsdk.Link `json:"links"`
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	golangsdk.Result
}

// Extract returns a pointer to a CreatedStack object and is called after a
// Create operation.
func (r CreateResult) Extract() (*CreatedStack, error) {
	var s struct {
		CreatedStack *CreatedStack `json:"stack"`
	}
	err := r.ExtractInto(&s)
	return s.CreatedStack, err
}

// StackPage is a pagination.Pager that is returned from a call to the List function.
type StackPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a ListResult contains no Stacks.
func (r StackPage) IsEmpty() (bool, error) {
	stacks, err := ExtractStacks(r)
	return len(stacks) == 0, err
}

// ListedStack represents an element in the slice extracted from a List operation.
type ListedStack struct {
	CreationTime time.Time        `json:"-"`
	Description  string           `json:"description"`
	ID           string           `json:"id"`
	Links        []golangsdk.Link `json:"links"`
	Name         string           `json:"stack_name"`
	Status       string           `json:"stack_status"`
	StatusReason string           `json:"stack_status_reason"`
	UpdatedTime  time.Time        `json:"-"`
}

func (r *ListedStack) UnmarshalJSON(b []byte) error {
	type tmp ListedStack
	var s struct {
		tmp
		CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
		UpdatedTime  golangsdk.JSONRFC3339NoZ `json:"updated_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ListedStack(s.tmp)

	r.CreationTime = time.Time(s.CreationTime)
	r.UpdatedTime = time.Time(s.UpdatedTime)

	return nil
}

// ExtractStacks extracts and returns a slice of ListedStack. It is used while iterating
// over a stacks.List call.
func ExtractStacks(r pagination.Page) ([]ListedStack, error) {
	var s struct {
		ListedStacks []ListedStack `json:"stacks"`
	}
	err := (r.(StackPage)).ExtractInto(&s)
	return s.ListedStacks, err
}

// RetrievedStack represents the object extracted from a Get operation.
// RetrievedStack represents the object extracted from a Get operation.
type RetrievedStack struct {
	Capabilities        []interface{}     `json:"capabilities"`
	CreationTime        time.Time         `json:"-"`
	Description         string            `json:"description"`
	DisableRollback     bool              `json:"disable_rollback"`
	ID                  string            `json:"id"`
	TenantId            string            `json:"tenant_id"`
	Links               []golangsdk.Link  `json:"links"`
	NotificationTopics  []interface{}     `json:"notification_topics"`
	Outputs             []*Output         `json:"outputs"`
	Parameters          map[string]string `json:"parameters"`
	Name                string            `json:"stack_name"`
	Status              string            `json:"stack_status"`
	StatusReason        string            `json:"stack_status_reason"`
	Tags                []string          `json:"tags"`
	TemplateDescription string            `json:"template_description"`
	Timeout             int               `json:"timeout_mins"`
	UpdatedTime         time.Time         `json:"-"`
}

// The Output data type.
type Output struct {
	_ struct{} `type:"structure"`

	// User defined description associated with the output.
	Description string `json:"description"`

	// The name of the export associated with the output.
	//ExportName *string `json:"name"`

	// The key associated with the output.
	OutputKey *string `json:"output_key"`

	// The value associated with the output.
	OutputValue *string `json:"output_value"`
}

// String returns the string representation
func (s Output) String() string {
	return Prettify(s)
}

// GoString returns the string representation
func (s Output) GoString() string {
	return s.String()
}

// SetDescription sets the Description field's value.
func (s *Output) SetDescription(v string) *Output {
	s.Description = v
	return s
}

// SetOutputKey sets the OutputKey field's value.
func (s *Output) SetOutputKey(v string) *Output {
	s.OutputKey = &v
	return s
}

// SetOutputValue sets the OutputValue field's value.
func (s *Output) SetOutputValue(v string) *Output {
	s.OutputValue = &v
	return s
}

// RetrievedStack represents the object extracted from a Get operation.
func (r *RetrievedStack) UnmarshalJSON(b []byte) error {
	type tmp RetrievedStack
	var s struct {
		tmp
		CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
		UpdatedTime  golangsdk.JSONRFC3339NoZ `json:"updated_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = RetrievedStack(s.tmp)

	r.CreationTime = time.Time(s.CreationTime)
	r.UpdatedTime = time.Time(s.UpdatedTime)

	return nil
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	golangsdk.Result
}

// Extract returns a pointer to a CreatedStack object and is called after a
// Create operation.
func (r GetResult) Extract() (*RetrievedStack, error) {
	var s struct {
		Stack *RetrievedStack `json:"stack"`
	}
	err := r.ExtractInto(&s)
	return s.Stack, err
}

// UpdateResult represents the result of a Update operation.
type UpdateResult struct {
	golangsdk.ErrResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}

// Prettify returns the string representation of a value.
func Prettify(i interface{}) string {
	var buf bytes.Buffer
	prettify(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// prettify will recursively walk value v to build a textual
// representation of the value.
func prettify(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		strtype := v.Type().String()
		if strtype == "time.Time" {
			fmt.Fprintf(buf, "%s", v.Interface())
			break
		} else if strings.HasPrefix(strtype, "io.") {
			buf.WriteString("<buffer>")
			break
		}

		buf.WriteString("{\n")

		names := []string{}
		for i := 0; i < v.Type().NumField(); i++ {
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) {
				continue // ignore unexported fields
			}
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice || f.Kind() == reflect.Map) && f.IsNil() {
				continue // ignore unset fields
			}
			names = append(names, name)
		}

		for i, n := range names {
			val := v.FieldByName(n)
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(n + ": ")
			prettify(val, indent+2, buf)

			if i < len(names)-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	case reflect.Slice:
		strtype := v.Type().String()
		if strtype == "[]uint8" {
			fmt.Fprintf(buf, "<binary> len %d", v.Len())
			break
		}

		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			prettify(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("{\n")

		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			prettify(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	default:
		if !v.IsValid() {
			fmt.Fprint(buf, "<invalid value>")
			return
		}
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		case io.ReadSeeker, io.Reader:
			format = "buffer(%p)"
		}
		fmt.Fprintf(buf, format, v.Interface())
	}
}
