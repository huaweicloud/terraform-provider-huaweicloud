package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func writeRegionAttribute(w io.Writer, att *schema.Schema) error {
	path := []string{"region"}

	att.Optional = true
	att.ForceNew = true
	if att.Description == "" {
		desc := "Specifies the region in which to create the resource.\n"
		if isData {
			desc = "Specifies the region in which to query the data source.\n"
		}
		att.Description = desc + "  If omitted, the provider-level region will be used."
	}
	return writeAttribute(w, path, att)
}

func writeIDAttribute(w io.Writer) error {
	path := []string{"id"}

	desc := "The resource ID."
	if isData {
		desc = "The data source ID."
	}

	att := &schema.Schema{
		Type:        schema.TypeString,
		Description: desc,
	}
	return writeReadOnlyAttribute(w, path, att)
}

func writeAttribute(w io.Writer, path []string, att *schema.Schema) error {
	name := path[len(path)-1]

	_, err := io.WriteString(w, "* `"+name+"` - ")
	if err != nil {
		return err
	}

	err = WriteAttributeDescription(w, att)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func writeReadOnlyAttribute(w io.Writer, path []string, att *schema.Schema) error {
	name := path[len(path)-1]

	_, err := io.WriteString(w, "* `"+name+"` - ")
	if err != nil {
		return err
	}

	err = WriteReadOnlyAttributeDescription(w, att)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "\n")
	if err != nil {
		return err
	}

	return nil
}

func WriteAttributeDescription(w io.Writer, att *schema.Schema) error {
	if att.Deprecated != "" {
		return nil
	}

	_, err := io.WriteString(w, "(")
	if err != nil {
		return err
	}

	if childAttributeIsRequired(att) {
		_, err = io.WriteString(w, "Required, ")
		if err != nil {
			return err
		}
	} else if childAttributeIsOptional(att) {
		_, err = io.WriteString(w, "Optional, ")
		if err != nil {
			return err
		}
	}

	err = WriteType(w, att.Type)
	if err != nil {
		return err
	}

	if att.ForceNew && !isData {
		_, err := io.WriteString(w, ", ForceNew")
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(w, ")")
	if err != nil {
		return err
	}

	desc := normalizeDescription(att.Description)
	if desc != "" {
		_, err = io.WriteString(w, fmt.Sprintf(" %s\n", desc))
	} else {
		_, err = io.WriteString(w, " <!-- please add the description of the argument -->\n")
	}
	if err != nil {
		return err
	}

	if att.ForceNew && !isData {
		_, err := io.WriteString(w, "  Changing this creates a new resource.\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteReadOnlyAttributeDescription(w io.Writer, att *schema.Schema) error {
	var err error
	if att.Deprecated != "" {
		return nil
	}

	desc := normalizeDescription(att.Description)
	if desc != "" {
		_, err = io.WriteString(w, fmt.Sprintf("%s\n", desc))
	} else {
		_, err = io.WriteString(w, "<!-- please add the description of the attribute -->\n")
	}

	return err
}

func normalizeDescription(desc string) string {
	prefix := "schema:"
	if !strings.HasPrefix(desc, prefix) {
		return desc
	}

	extras := strings.SplitN(desc, ";", 2)
	if len(extras) > 1 {
		return strings.Trim(extras[1], " ")
	}

	return ""
}

func WriteType(w io.Writer, ty schema.ValueType) error {
	var valueType string
	switch ty {
	case schema.TypeBool:
		valueType = "Bool"
	case schema.TypeInt:
		valueType = "Int"
	case schema.TypeFloat:
		valueType = "Float"
	case schema.TypeString:
		valueType = "String"
	case schema.TypeList, schema.TypeSet:
		valueType = "List"
	case schema.TypeMap:
		valueType = "Map"
	default:
		return fmt.Errorf("unexpected type %s", ty.String())
	}

	_, err := io.WriteString(w, valueType)
	return err
}
