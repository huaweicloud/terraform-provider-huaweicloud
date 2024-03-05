package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Render writes a Markdown formatted Schema definition to the specified writer.
func Render(w io.Writer, res *schema.Resource, name string) error {
	writeHeaderBlock(w, name)

	if err := writeArgumentBlock(w, res.Schema); err != nil {
		return fmt.Errorf("unable to render argument schema: %s", err)
	}

	if err := writeReadOnlyBlock(w, res.Schema); err != nil {
		return fmt.Errorf("unable to render attribute schema: %s", err)
	}

	writeTailBlock(w, res, name)
	return nil
}

type nestedType struct {
	anchorID string
	path     []string
	block    *schema.Schema
}

func writeBlockType(w io.Writer, path []string, block *schema.Schema, readOnly bool) ([]nestedType, error) {
	name := path[len(path)-1]

	_, err := io.WriteString(w, "* `"+name+"` - ")
	if err != nil {
		return nil, err
	}

	if readOnly {
		err = WriteReadOnlyAttributeDescription(w, block)
	} else {
		err = WriteAttributeDescription(w, block)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to write block description for %q: %w", name, err)
	}

	var anchorID string
	if readOnly {
		anchorID = "attrblock--" + strings.Join(path, "--")
	} else {
		anchorID = "block--" + strings.Join(path, "--")
	}
	nt := nestedType{
		anchorID: anchorID,
		path:     path,
		block:    block,
	}

	seeBelow := fmt.Sprintf("  The [%s](#%s) structure is documented below.\n\n", name, anchorID)
	_, err = io.WriteString(w, seeBelow)
	if err != nil {
		return nil, err
	}

	return []nestedType{nt}, nil
}

func writeHeaderBlock(w io.Writer, name string) {
	var backquote byte = 96
	headerStr := fmt.Sprintf(`
---
subcategory: "xxxx"
---

# %[1]s

<!--
please add the description of %[1]s
  + For resource: Manages xxx resource within HuaweiCloud.
  + For data source: Use this data source to get the list of xxx.
-->

## Example Usage

<!-- please add the usage of %[1]s -->
%[2]s%[2]s%[2]shcl

%[2]s%[2]s%[2]s

`, name, string(backquote))

	if _, err := io.WriteString(w, headerStr); err != nil {
		fmt.Printf("unable to render header: %s", err)
	}
}

func writeTailBlock(w io.Writer, res *schema.Resource, name string) {
	var backquote byte = 96
	if res.Timeouts != nil {
		timeoutsStr := "## Timeouts\n\n" +
			"This resource provides the following timeouts configuration options:\n\n"

		if res.Timeouts.Create != nil {
			timeoutsStr += fmt.Sprintf("* `create` - Default is %s.\n", res.Timeouts.Create)
		}

		if res.Timeouts.Update != nil {
			timeoutsStr += fmt.Sprintf("* `update` - Default is %s.\n", res.Timeouts.Update)
		}

		if res.Timeouts.Delete != nil {
			timeoutsStr += fmt.Sprintf("* `delete` - Default is %s.\n", res.Timeouts.Delete)
		}

		timeoutsStr = strings.ReplaceAll(timeoutsStr, "m0s.", " minutes.")
		if _, err := io.WriteString(w, timeoutsStr); err != nil {
			fmt.Printf("unable to render Timeouts: %s", err)
		}
	}

	if res.Importer != nil {
		importerStr := fmt.Sprintf(`
## Import

The xxx can be imported using %[1]sid%[1]s, e.g.

%[1]s%[1]s%[1]sbash
$ terraform import %[2]s.test <id>
%[1]s%[1]s%[1]s
`, string(backquote), name)

		if _, err := io.WriteString(w, importerStr); err != nil {
			fmt.Printf("unable to render Import: %s", err)
		}

		importIgnore := fmt.Sprintf(`
<!--
Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: %[1]sparam1%[1]s, %[1]sparam2%[1]s, %[1]sparam3%[1]s ...
It is generally recommended running %[1]sterraform plan%[1]s after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

%[1]s%[1]s%[1]shcl
resource "%s" "test" {
    ...

  lifecycle {
    ignore_changes = [
      param1, param2,
    ]
  }
}
%[1]s%[1]s%[1]s
-->
`, string(backquote), name)

		if _, err := io.WriteString(w, importIgnore); err != nil {
			fmt.Printf("unable to render Import ignore: %s", err)
		}
	}
}

func writeArgumentBlock(w io.Writer, blocks map[string]*schema.Schema) error {
	_, err := io.WriteString(w, "## Argument Reference\n\n"+
		"The following arguments are supported:\n\n")
	if err != nil {
		return err
	}

	if att, ok := blocks["region"]; ok {
		err := writeRegionAttribute(w, att)
		if err != nil {
			return fmt.Errorf("failed to write region attribute: %s", err)
		}

		delete(blocks, "region")
	}

	return writeBlockChildren(w, nil, blocks)
}

func writeReadOnlyBlock(w io.Writer, blocks map[string]*schema.Schema) error {
	_, err := io.WriteString(w, "## Attribute Reference\n\n"+
		"In addition to all arguments above, the following attributes are exported:\n\n")
	if err != nil {
		return err
	}

	err = writeIDAttribute(w)
	if err != nil {
		return fmt.Errorf("failed to write ID attribute: %s", err)
	}

	delete(blocks, "id")
	return writeReadOnlyChildren(w, nil, blocks)
}

func writeBlockChildren(w io.Writer, parents []string, blocks map[string]*schema.Schema) error {
	groups := [2][]string{}
	for n, childAtt := range blocks {
		if childAtt.Deprecated != "" {
			continue
		}

		if childAttributeIsRequired(childAtt) {
			groups[0] = append(groups[0], n)
		} else if childAttributeIsOptional(childAtt) {
			groups[1] = append(groups[1], n)
		}
	}

	// For each characteristic group
	//   If Attribute
	//     Write out summary including characteristic and type (if primitive type or collection of primitives)
	//     If NestedAttribute type, Object type or collection of Objects, add to list of nested types
	//   ElseIf Block
	//     Write out summary including characteristic
	//     Add block to list of nested types
	//   End
	// End
	// For each nested type:
	//   Write out heading
	//   If Block
	//     Recursively call this function (writeBlockChildren)
	//   ElseIf Object
	//     Call writeObjectChildren, which
	//       For each Object Attribute
	//         Write out summary including characteristic and type (if primitive type or collection of primitives)
	//         If Object type or collection of Objects, add to list of nested types
	//       End
	//       Recursively do nested type functionality
	//   ElseIf NestedAttribute
	//     Call writeNestedAttributeChildren, which
	//       For each nested Attribute
	//         Write out summary including characteristic and type (if primitive type or collection of primitives)
	//         If NestedAttribute type, Object type or collection of Objects, add to list of nested types
	//       End
	//       Recursively do nested type functionality
	//   End
	// End
	nestedTypes := []nestedType{}
	for i := range groups {
		sortedNames := groups[i]
		if len(sortedNames) == 0 {
			continue
		}
		sort.Strings(sortedNames)

		for _, name := range sortedNames {
			path := parents
			path = append(path, name)

			childAttr := blocks[name]
			if childAttributeIsBlock(childAttr) {
				nt, err := writeBlockType(w, path, childAttr, false)
				if err != nil {
					return fmt.Errorf("unable to render block %q: %w", name, err)
				}

				nestedTypes = append(nestedTypes, nt...)
			} else {
				err := writeAttribute(w, path, childAttr)
				if err != nil {
					return fmt.Errorf("unable to render attribute %q: %w", name, err)
				}
			}
		}
	}

	err := writeNestedTypes(w, nestedTypes)
	if err != nil {
		return err
	}

	return nil
}

func writeReadOnlyChildren(w io.Writer, parents []string, blocks map[string]*schema.Schema) error {
	names := []string{}
	for n, childAtt := range blocks {
		if childAtt.Deprecated != "" {
			continue
		}
		if !childAttributeIsReadOnly(childAtt) {
			continue
		}
		names = append(names, n)
	}

	if len(names) == 0 {
		return nil
	}

	sortedNames := names
	sort.Strings(sortedNames)

	nestedTypes := []nestedType{}
	for _, name := range sortedNames {
		path := parents
		path = append(path, name)

		if childAtt, ok := blocks[name]; ok {
			if childAttributeIsBlock(childAtt) {
				nt, err := writeBlockType(w, path, childAtt, true)
				if err != nil {
					return fmt.Errorf("unable to render block %q: %w", name, err)
				}

				nestedTypes = append(nestedTypes, nt...)
			} else {
				err := writeReadOnlyAttribute(w, path, childAtt)
				if err != nil {
					return fmt.Errorf("unable to render attribute %q: %w", name, err)
				}
			}
		}
	}

	err := writeReadOnlyNestedTypes(w, nestedTypes)
	if err != nil {
		return err
	}

	return nil
}

func writeNestedTypes(w io.Writer, nestedTypes []nestedType) error {
	for _, nt := range nestedTypes {
		_, err := io.WriteString(w, "<a name=\""+nt.anchorID+"\"></a>\n")
		if err != nil {
			return err
		}

		name := nt.path[len(nt.path)-1]
		_, err = io.WriteString(w, fmt.Sprintf("The `%s` block supports:\n\n", name))
		if err != nil {
			return err
		}

		var nestBlock map[string]*schema.Schema
		if nt.block.Elem != nil {
			if v, ok := nt.block.Elem.(*schema.Resource); ok {
				nestBlock = v.Schema
			}
		}
		if nestBlock == nil {
			return fmt.Errorf("TODO: not yet supported")
		}

		err = writeBlockChildren(w, nt.path, nestBlock)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeReadOnlyNestedTypes(w io.Writer, nestedTypes []nestedType) error {
	for _, nt := range nestedTypes {
		_, err := io.WriteString(w, "<a name=\""+nt.anchorID+"\"></a>\n")
		if err != nil {
			return err
		}

		name := nt.path[len(nt.path)-1]
		_, err = io.WriteString(w, fmt.Sprintf("The `%s` block supports:\n\n", name))
		if err != nil {
			return err
		}

		var nestBlock map[string]*schema.Schema
		if nt.block.Elem != nil {
			if v, ok := nt.block.Elem.(*schema.Resource); ok {
				nestBlock = v.Schema
			}
		}
		if nestBlock == nil {
			return fmt.Errorf("TODO: not yet supported")
		}

		err = writeReadOnlyChildren(w, nt.path, nestBlock)
		if err != nil {
			return err
		}
	}

	return nil
}

func childAttributeIsBlock(att *schema.Schema) bool {
	if att.Elem != nil {
		if _, ok := att.Elem.(*schema.Resource); ok {
			return true
		}
	}
	return false
}

func childAttributeIsRequired(att *schema.Schema) bool {
	return att.Required
}

func childAttributeIsOptional(att *schema.Schema) bool {
	return att.Optional
}

// Read-only is computed but not optional.
func childAttributeIsReadOnly(att *schema.Schema) bool {
	if att.Elem != nil {
		if nestBlock, ok := att.Elem.(*schema.Resource); ok {
			for _, subAtt := range nestBlock.Schema {
				if implChildAttributeIsReadOnly(subAtt) {
					return true
				}
			}
		}
		return implChildAttributeIsReadOnly(att)
	}

	return implChildAttributeIsReadOnly(att)
}

func implChildAttributeIsReadOnly(att *schema.Schema) bool {
	// these shouldn't be able to be required, but just in case
	return att.Computed && !att.Optional && !att.Required
}
