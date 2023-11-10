// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package idna

// This file contains definitions for interpreting the trie value of the idna
// trie generated by "go run gen*.go". It is shared by both the generator
// program and the resultant package. Sharing is achieved by the generator
// copying gen_trieval.go to trieval.go and changing what's above this comment.

// info holds information from the IDNA mapping table for a single rune. It is
// the value returned by a trie lookup. In most cases, all information fits in
// a 16-bit value. For mappings, this value may contain an index into a slice
// with the mapped string. Such mappings can consist of the actual mapped value
// or an XOR pattern to be applied to the bytes of the UTF8 encoding of the
// input rune. This technique is used by the cases packages and reduces the
// table size significantly.
//
// The per-rune values have the following format:
//
//	if mapped {
//	  if inlinedXOR {
//	    15..13 inline XOR marker
//	    12..11 unused
//	    10..3  inline XOR mask
//	  } else {
//	    15..3  index into xor or mapping table
//	  }
//	} else {
//	    15..14 unused
//	    13     mayNeedNorm
//	    12..11 attributes
//	    10..8  joining type
//	     7..3  category type
//	}
//	   2  use xor pattern
//	1..0  mapped category
//
// See the definitions below for a more detailed description of the various
// bits.
type info uint16

const (
	catSmallMask = 0x3
	catBigMask   = 0xF8
	indexShift   = 3
	xorBit       = 0x4    // interpret the index as an xor pattern
	inlineXOR    = 0xE000 // These bits are set if the XOR pattern is inlined.

	joinShift = 8
	joinMask  = 0x07

	// Attributes
	attributesMask = 0x1800
	viramaModifier = 0x1800
	modifier       = 0x1000
	rtl            = 0x0800

	mayNeedNorm = 0x2000
)

// A category corresponds to a category defined in the IDNA mapping table.
type category uint16

const (
	unknown              category = 0 // not currently defined in unicode.
	mapped               category = 1
	disallowedSTD3Mapped category = 2
	deviation            category = 3
)

const (
	valid               category = 0x08
	validNV8            category = 0x18
	validXV8            category = 0x28
	disallowed          category = 0x40
	disallowedSTD3Valid category = 0x80
	ignored             category = 0xC0
)

// join types and additional rune information
const (
	joiningL = (iota + 1)
	joiningD
	joiningT
	joiningR

	//the following types are derived during processing
	joinZWJ
	joinZWNJ
	joinVirama
	numJoinTypes
)

func (c info) isMapped() bool {
	return c&0x3 != 0
}

func (c info) category() category {
	small := c & catSmallMask
	if small != 0 {
		return category(small)
	}
	return category(c & catBigMask)
}

func (c info) joinType() info {
	if c.isMapped() {
		return 0
	}
	return (c >> joinShift) & joinMask
}

func (c info) isModifier() bool {
	return c&(modifier|catSmallMask) == modifier
}

func (c info) isViramaModifier() bool {
	return c&(attributesMask|catSmallMask) == viramaModifier
}
