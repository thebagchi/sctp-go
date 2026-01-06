package per

import (
	"encoding/asn1"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/thebagchi/sctp-go/gnbsim/per/bitbuffer"
)

// Decoder represents a PER decoder for bit-level decoding
type Decoder struct {
	codec   *bitbuffer.Codec
	aligned bool
}

// NewDecoder creates a new PER decoder from a byte slice
// aligned: true for APER (Aligned PER), false for UPER (Unaligned PER)
func NewDecoder(data []byte, aligned bool) *Decoder {
	return &Decoder{
		codec:   bitbuffer.CreateReader(data),
		aligned: aligned,
	}
}

// Align aligns the decoder to byte boundary (only for aligned mode)
func (d *Decoder) Align() error {
	if d.aligned {
		return d.codec.Advance()
	}
	return nil
}

// ReadInt decodes an integer value
// ITU-T X.691 Section 10.2: Integer
func (d *Decoder) ReadInt(lb, ub *int64, extensible bool) (int64, error) {
	if lb != nil && ub != nil {
		if extensible {
			extended, err := d.ReadExtensionBit()
			if err != nil {
				return 0, err
			}

			if extended {
				return d.readUnconstrainedInt()
			}
		}

		valueRange := *ub - *lb + 1

		value, err := d.readConstrainedValue(valueRange)
		if err != nil {
			return 0, err
		}

		return *lb + int64(value), nil

	} else if lb != nil {
		value, err := d.readSemiConstrainedValue(*lb)
		if err != nil {
			return 0, err
		}
		return value, nil

	} else {
		return d.readUnconstrainedInt()
	}
}

// ReadBool decodes a boolean value (single bit)
// ITU-T X.691 Section 10.1: Boolean
func (d *Decoder) ReadBool() (bool, error) {
	value, err := d.codec.Read(1)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

// ReadExtensionBit decodes the extension presence bit
// ITU-T X.691 Clauses 10.9, 11.5, 27-30: Extension presence bit (0 = no extension, 1 = extension present)
func (d *Decoder) ReadExtensionBit() (bool, error) {
	value, err := d.codec.Read(1)
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

// ReadFloat decodes a floating-point value
// ITU-T X.691 Section 10.5: Real
func (d *Decoder) ReadFloat() (float64, error) {
	bits, err := d.codec.Read(64)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bits), nil
}

// ReadNull decodes NULL (no decoding required)
// ITU-T X.691 Section 10.0: Null
func (d *Decoder) ReadNull() error {
	return nil
}

// ReadObjectIdentifier decodes OBJECT IDENTIFIER
// ITU-T X.691 Section 10.2: Object Identifier
// ObjectIdentifier is encoded as an OCTET STRING containing the BER encoding
func (d *Decoder) ReadObjectIdentifier() (asn1.ObjectIdentifier, error) {
	// Read as OCTET STRING (which contains BER-encoded OID)
	encoded, err := d.ReadOctetString(nil, nil, false)
	if err != nil {
		return nil, fmt.Errorf("failed to read ObjectIdentifier octet string: %w", err)
	}

	// Unmarshal the BER-encoded OID
	var oid asn1.ObjectIdentifier
	rest, err := asn1.Unmarshal(encoded, &oid)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ObjectIdentifier: %w", err)
	}

	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after ObjectIdentifier unmarshal")
	}

	return oid, nil
}

// ReadChoice decodes the choice index for a CHOICE type
// ITU-T X.691 Section 11.4: CHOICE
// ReadChoiceIndex decodes a CHOICE index (returns the selected alternative index)
func (d *Decoder) ReadChoiceIndex(numChoices int, extensible bool) (int, error) {
	if extensible {
		extended, err := d.ReadExtensionBit()
		if err != nil {
			return 0, err
		}

		if extended {
			value, err := d.readNormallySmallValue()
			if err != nil {
				return 0, err
			}
			return int(value), nil
		}
	}

	if numChoices <= 1 {
		return 0, nil
	}

	valueRange := int64(numChoices)
	index, err := d.readConstrainedValue(valueRange)
	if err != nil {
		return 0, err
	}

	if index >= uint64(numChoices) {
		return 0, fmt.Errorf("ReadChoice: decoded index %d is out of range [0, %d)", index, numChoices)
	}

	return int(index), nil
}

// ReadBitString decodes a bit string with optional size constraints
// ITU-T X.691 Section 10.7: Bit String
func (d *Decoder) ReadBitString(lb, ub *int64, extensible bool) (asn1.BitString, error) {
	var nbits uint

	if extensible {
		extended, err := d.ReadExtensionBit()
		if err != nil {
			return asn1.BitString{}, err
		}

		if extended {
			value, err := d.readNormallySmallValue()
			if err != nil {
				return asn1.BitString{}, err
			}
			nbits = uint(value)
			// Align before contents in APER if bitLen >= 16 bits
			if d.aligned && nbits >= 16 {
				if err := d.Align(); err != nil {
					return asn1.BitString{}, err
				}
			}
			bits, err := d.readBits(nbits)
			return asn1.BitString{Bytes: bits, BitLength: int(nbits)}, err
		}
	}

	if lb != nil && ub != nil {
		sizeRange := *ub - *lb + 1

		size, err := d.readConstrainedValue(sizeRange)
		if err != nil {
			return asn1.BitString{}, err
		}

		nbits = uint(int64(size) + *lb)

	} else if lb != nil {
		size, err := d.readSemiConstrainedValue(*lb)
		if err != nil {
			return asn1.BitString{}, err
		}
		nbits = uint(size)

	} else {
		size, err := d.readLength()
		if err != nil {
			return asn1.BitString{}, err
		}
		nbits = uint(size)
	}

	// Align before contents in APER if bitLen >= 16 bits (common threshold for byte alignment)
	if d.aligned && nbits >= 16 {
		if err := d.Align(); err != nil {
			return asn1.BitString{}, err
		}
	}

	value, err := d.readBits(nbits)
	if err != nil {
		return asn1.BitString{}, err
	}

	return asn1.BitString{Bytes: value, BitLength: int(nbits)}, nil
}

// ReadOctetString decodes an octet string (byte sequence) with optional size constraints
// ITU-T X.691 Section 10.8: Octet String
func (d *Decoder) ReadOctetString(lb, ub *int64, extensible bool) ([]byte, error) {
	var size uint64

	if extensible {
		extended, err := d.ReadExtensionBit()
		if err != nil {
			return nil, err
		}
		if extended {
			temp, err := d.readNormallySmallValue()
			if err != nil {
				return nil, err
			}
			size = temp
			if d.aligned {
				if err := d.Align(); err != nil {
					return nil, err
				}
			}
			return d.codec.ReadBytes(int(size))
		}
	}

	if lb != nil && ub != nil {
		rangeSize := *ub - *lb + 1
		decoded, err := d.readConstrainedValue(rangeSize)
		if err != nil {
			return nil, err
		}
		size = uint64(int64(decoded) + *lb)
	} else if lb != nil {
		decoded, err := d.readSemiConstrainedValue(*lb)
		if err != nil {
			return nil, err
		}
		size = uint64(decoded)
	} else {
		temp, err := d.readLength()
		if err != nil {
			return nil, err
		}
		size = temp
	}

	// Align before contents in APER (X.691 Section 10.8)
	if d.aligned {
		if err := d.Align(); err != nil {
			return nil, err
		}
	}
	return d.codec.ReadBytes(int(size))
}

// ReadCharacterString decodes a character string with optional size constraints
// ITU-T X.691 Section 10.9-10.12: Restricted Character Strings
func (d *Decoder) ReadCharacterString(lb, ub *int64, extensible bool) (string, error) {
	value, err := d.ReadOctetString(lb, ub, extensible)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// ReadEnumeration decodes an enumeration value with optional constraints
// ITU-T X.691 Section 10.6: Enumerated
func (d *Decoder) ReadEnumeration(numValues uint64, extensible bool) (uint64, error) {
	if extensible {
		extended, err := d.ReadExtensionBit()
		if err != nil {
			return 0, err
		}

		if extended {
			return d.readNormallySmallValue()
		}
	}

	return d.readConstrainedValue(int64(numValues))
}

// readConstrainedValue decodes a constrained whole number
// ITU-T X.691 Section 10.2: Whole Number Constraints
func (d *Decoder) readConstrainedValue(valueRange int64) (uint64, error) {
	if valueRange <= 0 {
		return 0, fmt.Errorf("value range must be positive")
	}

	if valueRange <= 255 {
		bitsNeeded := bitsRequired(uint64(valueRange - 1))
		if bitsNeeded == 0 {
			bitsNeeded = 1
		}
		return d.codec.Read(uint8(bitsNeeded))
	}

	if d.aligned {
		if err := d.codec.Align(); err != nil {
			return 0, err
		}
	}

	if valueRange == 256 {
		return d.codec.Read(8)
	}

	if valueRange <= 0x10000 {
		return d.codec.Read(16)
	}

	// For valueRange > 65536: length + value (X.691 10.2.5)
	// Length is encoded as constrained whole number with range 1-8 bytes
	temp, err := d.readConstrainedValue(8)
	if err != nil {
		return 0, err
	}
	length := temp + 1
	if length > 8 {
		return 0, fmt.Errorf("large constrained integer length %d exceeds maximum 8 bytes", length)
	}

	value, err := d.codec.Read(uint8(length * 8))
	if err != nil {
		return 0, err
	}
	return value, nil
}

// readSemiConstrainedValue decodes a semi-constrained whole number
func (d *Decoder) readSemiConstrainedValue(lb int64) (int64, error) {
	if d.aligned {
		if err := d.codec.Align(); err != nil {
			return 0, err
		}
	}

	length, err := d.readLength()
	if err != nil {
		return 0, err
	}

	value, err := d.codec.Read(uint8(length * 8))
	if err != nil {
		return 0, err
	}

	return lb + int64(value), nil
}

// readUnconstrainedInt decodes an unconstrained integer
func (d *Decoder) readUnconstrainedInt() (int64, error) {
	if d.aligned {
		if err := d.codec.Align(); err != nil {
			return 0, err
		}
	}

	length, err := d.readLength()
	if err != nil {
		return 0, err
	}

	if length > 8 {
		return 0, fmt.Errorf("integer length %d exceeds maximum 8 bytes", length)
	}

	value, err := d.codec.Read(uint8(length * 8))
	if err != nil {
		return 0, err
	}

	return int64(value), nil
}

// readNormallySmallValue decodes a normally small non-negative whole number
func (d *Decoder) readNormallySmallValue() (uint64, error) {
	gt64, err := d.codec.Read(1)
	if err != nil {
		return 0, err
	}

	if gt64 == 0 {
		val, err := d.codec.Read(6)
		if err != nil {
			return 0, err
		}
		return val, nil
	}

	length, err := d.readLength()
	if err != nil {
		return 0, err
	}

	return d.codec.Read(uint8(length * 8))
}

// readLength decodes a length determinant
// ITU-T X.691 Section 10.9.3.6–10.9.3.8: Length Determinant with Fragmentation
func (d *Decoder) readLength() (uint64, error) {
	if d.aligned {
		if err := d.Align(); err != nil {
			return 0, err
		}
	}

	first, err := d.codec.Read(8)
	if err != nil {
		return 0, err
	}

	if first < 128 {
		return first, nil
	}
	if first < 0xC0 {
		second, err := d.codec.Read(8)
		if err != nil {
			return 0, err
		}
		return uint64((first&0x3F)<<8 | second), nil
	}

	// Fragmentation
	var (
		total   = uint64(0)
		current = first
	)

	for current >= 0xC0 {
		if d.aligned {
			if err := d.Align(); err != nil {
				return 0, err
			}
		}

		switch current {
		case 0xC0:
			total += 16384
		case 0xC1:
			total += 32768
		case 0xC2:
			total += 49152
		case 0xC3:
			total += 65536
		default:
			return 0, fmt.Errorf("invalid fragment indicator 0x%02X", current)
		}

		current, err = d.codec.Read(8)
		if err != nil {
			return 0, err
		}
	}

	// Final length
	if current < 128 {
		total += current
		return total, nil
	}
	if current < 0xC0 {
		second, err := d.codec.Read(8)
		if err != nil {
			return 0, err
		}
		total += uint64((current&0x3F)<<8 | second)
		return total, nil
	}
	return 0, fmt.Errorf("invalid final length 0x%02X", current)
}

// readBits reads a bit string from the decoder
func (d *Decoder) readBits(numBits uint) ([]byte, error) {
	if numBits == 0 {
		return []byte{}, nil
	}

	var (
		result = make([]byte, (numBits+7)/8)
		nbytes = numBits / 8
		rbits  = numBits % 8
	)

	for i := range nbytes {
		b, err := d.codec.Read(8)
		if err != nil {
			return nil, err
		}
		result[i] = byte(b)
	}

	if rbits > 0 {
		b, err := d.codec.Read(uint8(rbits))
		if err != nil {
			return nil, err
		}
		result[nbytes] = byte(b) << (8 - rbits)
	}

	return result, nil
}

// ReadSequence decodes a SEQUENCE using reflection
func (d *Decoder) ReadSequence(value any) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("ReadSequence requires a struct")
	}
	return d.decodeSequence(v)
}

// decodeSequence decodes a struct (SEQUENCE or CHOICE)
func (d *Decoder) decodeSequence(v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", v.Kind())
	}

	t := v.Type()
	numFields := t.NumField()

	// Determine if extensible and where extension starts
	var extensible bool
	var extensionStart int = numFields

	for i := range numFields {
		field := t.Field(i)
		if field.Name == "_" && field.Type.Kind() == reflect.Struct && field.Type.NumField() == 0 {
			tag := field.Tag.Get(TAG_KEY)
			opts := parseTag(tag)
			if opts.ext {
				extensible = true
				extensionStart = i + 1
				break
			}
		}
	}

	// Count optional fields for presence bitmap (root only)
	var optionalCount int
	for i := 0; i < extensionStart; i++ {
		field := t.Field(i)
		if field.Name == "_" {
			continue
		}
		tag := field.Tag.Get(TAG_KEY)
		opts := parseTag(tag)
		if opts.opt || field.Type.Kind() == reflect.Ptr {
			optionalCount++
		}
	}

	// Read optional presence bitmap (X.691 11.5.1)
	// Per X.691 11.5.1: N bits where N = number of optional fields
	optionalPresent := make([]bool, optionalCount)
	if optionalCount > 0 {
		// Read bitmap bits directly (no length encoding for preamble)
		bitmapBytes, err := d.readBits(uint(optionalCount))
		if err != nil {
			return err
		}
		for i := 0; i < optionalCount; i++ {
			byteIdx := i / 8
			bitIdx := uint(7 - (i % 8))
			optionalPresent[i] = (bitmapBytes[byteIdx] & (1 << bitIdx)) != 0
		}
	}

	// Read extension bit if extensible
	var hasExtension bool
	if extensible {
		extBit, err := d.ReadExtensionBit()
		if err != nil {
			return err
		}
		hasExtension = extBit
	}

	// Decode root fields
	optIdx := 0
	for i := 0; i < extensionStart; i++ {
		field := t.Field(i)
		if field.Name == "_" {
			continue
		}
		fv := v.Field(i)
		tag := field.Tag.Get(TAG_KEY)
		opts := parseTag(tag)

		isOptional := opts.opt || field.Type.Kind() == reflect.Ptr
		if isOptional {
			if optIdx >= len(optionalPresent) {
				return fmt.Errorf("optional bitmap mismatch")
			}
			if !optionalPresent[optIdx] {
				optIdx++
				continue // absent
			}
			optIdx++
		}

		if err := d.decodeField(field, fv); err != nil {
			return err
		}
	}

	// Extension handling: read extension bitmap and decode extension fields as open types
	if hasExtension && extensionStart > 0 && extensionStart < numFields {
		// Collect extension fields defined in this struct
		var extensionFields []int
		for i := extensionStart; i < numFields; i++ {
			field := t.Field(i)
			if field.Name == "_" {
				continue
			}
			tag := field.Tag.Get(TAG_KEY)
			if tag != "" {
				extensionFields = append(extensionFields, i)
			}
		}

		// Read extension bitmap from the stream
		// The bitmap indicates which extension fields (from the encoded message) are present
		// Per X.691 11.5.3: Extension bitmap is length-prefixed
		if d.aligned {
			if err := d.Align(); err != nil {
				return err
			}
		}
		extBitmapLen, err := d.readLength()
		if err != nil {
			return err
		}

		// Read extension bitmap bits as BIT STRING
		// May be larger than what we have fields for - we need to read all bits to stay in sync
		extBitmap := make([]bool, extBitmapLen)
		if extBitmapLen > 0 {
			// Read bitmap bits directly (no length encoding for preamble within the open type)
			bitmapBytes, err := d.readBits(uint(extBitmapLen))
			if err != nil {
				return err
			}
			for i := uint64(0); i < extBitmapLen; i++ {
				byteIdx := i / 8
				bitIdx := uint(7 - (i % 8))
				extBitmap[i] = (bitmapBytes[byteIdx] & (1 << bitIdx)) != 0
			}

			// Decode extension fields
			// We iterate through all bits in the bitmap, decoding each present field
			for i := uint64(0); i < extBitmapLen; i++ {
				if extBitmap[i] {
					// Align before open type (X.691 11.5.3)
					if d.aligned {
						if err := d.Align(); err != nil {
							return err
						}
					}

					// Read open type (unconstrained octet string) for every present field
					extBytes, err := d.ReadOctetString(nil, nil, false)
					if err != nil {
						return err
					}

					// Check if we have this field defined
					if i < uint64(len(extensionFields)) {
						fieldIdx := extensionFields[i]
						field := t.Field(fieldIdx)
						fv := v.Field(fieldIdx)

						// Decode extension value from the open type bytes
						tempDec := NewDecoder(extBytes, d.aligned)
						if err := tempDec.decodeField(field, fv); err != nil {
							return err
						}
					}
					// If field doesn't exist in our struct (from a newer version),
					// we've already consumed the open type data, so just skip it
				}
			}
		}
	}

	return nil
}

// ReadSequenceOf decodes SEQUENCE OF (slice/array)
func (d *Decoder) ReadSequenceOf(elements any, lb, ub *int64, extensible bool) error {
	v := reflect.ValueOf(elements)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("ReadSequenceOf requires pointer to slice/array")
	}
	slice := v.Elem()
	if slice.Kind() != reflect.Slice && slice.Kind() != reflect.Array {
		return fmt.Errorf("ReadSequenceOf requires slice or array")
	}

	// Decode length
	var length uint64
	if lb != nil && ub != nil {
		rangeSize := *ub - *lb + 1
		if extensible {
			extended, err := d.ReadExtensionBit()
			if err != nil {
				return err
			}
			if extended {
				val, err := d.readNormallySmallValue()
				if err != nil {
					return err
				}
				length = val
			} else {
				val, err := d.readConstrainedValue(rangeSize)
				if err != nil {
					return err
				}
				length = val + uint64(*lb)
			}
		} else {
			val, err := d.readConstrainedValue(rangeSize)
			if err != nil {
				return err
			}
			length = val + uint64(*lb)
		}
	} else if lb != nil {
		val, err := d.readSemiConstrainedValue(*lb)
		if err != nil {
			return err
		}
		length = uint64(val)
	} else {
		val, err := d.readLength()
		if err != nil {
			return err
		}
		length = val
	}

	// Resize slice if needed
	if slice.Kind() == reflect.Slice {
		slice.Set(reflect.MakeSlice(slice.Type(), int(length), int(length)))
	} else if int(length) != slice.Len() {
		return fmt.Errorf("fixed array length mismatch: expected %d, got %d", slice.Len(), length)
	}

	// Decode elements
	for i := 0; i < int(length); i++ {
		elem := slice.Index(i)
		if err := d.decodeValue(elem.Addr().Interface()); err != nil {
			return err
		}
	}
	return nil
}

// ReadChoice decodes a CHOICE using the _ extensible marker
// Extension alternatives are decoded as open types (length + raw PER bits)
func (d *Decoder) ReadChoice(value any) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("ReadChoice requires pointer to struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("ReadChoice requires struct")
	}

	t := v.Type()

	// Find root count and detect extensibility
	rootCount := 0
	extensible := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "_" && field.Type.Kind() == reflect.Struct {
			tag := field.Tag.Get(TAG_KEY)
			opts := parseTag(tag)
			if opts.ext {
				extensible = true
				break
			}
		}
		if tag := field.Tag.Get(TAG_KEY); strings.Contains(tag, "choice=") {
			rootCount++
		}
	}

	// Read choice index using ReadChoiceIndex
	index, err := d.ReadChoiceIndex(rootCount, extensible)
	if err != nil {
		return err
	}

	// Determine if this is an extension
	extended := extensible && index >= rootCount

	// For extensions, we need to read the open type data first
	// This ensures we consume it from the stream even if the field doesn't exist in our struct
	var openTypeData []byte
	if extended {
		// Align before open type (X.691 11.5.3)
		if d.aligned {
			if err := d.Align(); err != nil {
				return err
			}
		}

		var err error
		openTypeData, err = d.ReadOctetString(nil, nil, false)
		if err != nil {
			return err
		}
	}

	// Find field with matching choice index
	found := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(TAG_KEY)
		opts := parseTag(tag)
		if opts.choice == nil {
			continue
		}
		if *opts.choice == index {
			fv := v.Field(i)
			if fv.Kind() == reflect.Pointer {
				fv.Set(reflect.New(fv.Type().Elem()))
				fv = fv.Elem()
			}

			if extended {
				// Extension alternative - decode from open type
				tempDec := NewDecoder(openTypeData, d.aligned)
				if err := tempDec.decodeValue(fv.Addr().Interface()); err != nil {
					return err
				}
			} else {
				// Root alternative - decode value directly
				if err := d.decodeValue(fv.Addr().Interface()); err != nil {
					return err
				}
			}

			found = true
			break
		}
	}

	// If extension field was not found in struct, it's from a newer version - already consumed from stream
	if !found && extended {
		return nil
	}

	if !found {
		return fmt.Errorf("unknown choice index %d", index)
	}
	return nil
}

// decodeField decodes a single field
func (d *Decoder) decodeField(field reflect.StructField, value reflect.Value) error {
	tag := field.Tag.Get(TAG_KEY)
	opts := parseTag(tag)

	switch field.Type.Kind() {
	case reflect.Int64:
		val, err := d.ReadInt(opts.lb, opts.ub, opts.ext)
		if err != nil {
			return err
		}
		value.SetInt(val)
		return nil

	case reflect.Bool:
		val, err := d.ReadBool()
		if err != nil {
			return err
		}
		value.SetBool(val)
		return nil

	case reflect.String:
		val, err := d.ReadCharacterString(opts.lb, opts.ub, opts.ext)
		if err != nil {
			return err
		}
		value.SetString(val)
		return nil

	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			val, err := d.ReadOctetString(opts.lb, opts.ub, opts.ext)
			if err != nil {
				return err
			}
			value.SetBytes(val)
			return nil
		}
		// SEQUENCE OF
		return d.ReadSequenceOf(value.Addr().Interface(), opts.lb, opts.ub, opts.ext)

	case reflect.Array:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			val, err := d.ReadOctetString(opts.lb, opts.ub, opts.ext)
			if err != nil {
				return err
			}
			reflect.Copy(value, reflect.ValueOf(val))
			return nil
		}
		return d.ReadSequenceOf(value.Addr().Interface(), opts.lb, opts.ub, opts.ext)

	case reflect.Struct:
		if opts.choice != nil {
			return d.ReadChoice(value.Addr().Interface())
		}
		return d.ReadSequence(value.Addr().Interface())

	default:
		return fmt.Errorf("unsupported field type %s", field.Type)
	}
}

// decodeValue helper
func (d *Decoder) decodeValue(value any) error {
	if value == nil {
		return nil
	}
	return d.ReadSequence(value)
}
