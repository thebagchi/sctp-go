package per

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Example usage:
//
//		------------------ ENUMERATED ------------------
//		type StatusType uint64
//
//		const (
//			StatusIdle       StatusType = iota // 0: System is idle
//			StatusConnecting                   // 1: Establishing connection
//			StatusActive                       // 2: System is active
//			StatusError                        // 3: Error state
//			// Future extension values would start at 4 or higher
//		)
//
//		------------------ CHOICE ------------------
//		type BandwidthConfig struct {
//			Value int64 `per:"lb=1,ub=10000"` // Bandwidth value in Mbps (1-10000)
//		}
//		type ActionChoice struct {
//			// Root alternatives (indices 0–2)
//			Start  *bool            `per:"choice=0"`   // Start operation
//			Stop   *bool            `per:"choice=1"`   // Stop operation
//			Config *BandwidthConfig `per:"choice=2"`   // Configure bandwidth
//			_      struct{}         `per:"ext"`        // Extension marker
//			Pause  *bool            `per:"choice=3"`   // Pause operation (extension)
//			Resume *bool            `per:"choice=4"`   // Resume operation (extension)
//			Reboot *bool            `per:"choice=5"`   // Reboot system (extension)
//		}
//
//
//		------------------ SEQUENCE OF (List) ------------------
//		type Channel struct {
//			ID    int64 `per:"lb=1,ub=1000"` // Channel identifier (1-1000)
//			Power int64 `per:"lb=0,ub=100"`  // Power level in dBm (0-100)
//		}
//		type ChannelList struct {
//			FixedChannels      [8]Channel `per:"lb=8,ub=8"`            // Exactly 8 channels required
//			VariableChannels   []Channel  `per:"lb=0,ub=16"`           // 0-16 channels allowed
//			ExtensibleChannels []Channel  `per:"lb=1,ub=5,ext"`        // 1-5 in root, extensible
//		}
//
//		------------------ FULL CONSOLIDATED EXAMPLE ------------------
//		type FullMessage struct {
//			MessageID     int64            `per:"lb=1,ub=100000"`       // Unique message identifier (1-100000)
//			Priority      *int64           `per:"lb=0,ub=10"`           // Message priority level (0-10, optional)
//			Flags         asn1.BitString   `per:"lb=32,ub=32"`          // Fixed 32-bit flag set
//			Mask          asn1.BitString   `per:"lb=0,ub=64"`           // Variable mask up to 64 bits
//			Payload       []byte           `per:"lb=0,ub=1024"`         // Binary payload data (0-1024 bytes)
//			Status        StatusType       `per:"enum=4,ext"`           // Current system status (4 root values + extensible)
//			Action        ActionChoice     `per:"ext"`                  // Action to perform (extensible choice)
//			Channels      []Channel        `per:"lb=1,ub=8,ext"`        // List of channels (1-8 in root, extensible)
//			_             struct{}         `per:"ext"`                  // Extension marker for the SEQUENCE
//			FutureVersion *int64           `per:""`                     // Future version field (extension)
//			Experimental  *BandwidthConfig `per:""`                     // Experimental bandwidth config (extension)
//		}

const (
	TAG_KEY            = "per"
	TAG_OPTIONAL       = "opt"
	TAG_EXTENSIBLE     = "ext"
	TAG_PREFIX_CHOICE  = "choice="
	TAG_PREFIX_ENUM    = "enum="
	TAG_PREFIX_LB      = "lb="
	TAG_PREFIX_UB      = "ub="
	TAG_PREFIX_DEFAULT = "def="
)

type Tag struct {
	opt    bool    // field is OPTIONAL
	def    *string // default value (as string)
	ext    bool    // type has ... (extensible)
	lb     *int64  // lower bound
	ub     *int64  // upper bound
	choice *int    // CHOICE alternative index
	enum   *uint64 // ENUMERATED root value count
}

func parseTag(tag string) *Tag {
	opt := &Tag{}

	if tag == "" {
		return opt
	}

	for part := range strings.SplitSeq(tag, ",") {
		part = strings.TrimSpace(part)

		switch {
		case part == "opt":
			opt.opt = true

		case part == "ext":
			opt.ext = true

		case strings.HasPrefix(part, "choice="):
			if s := strings.TrimPrefix(part, "choice="); s != "" {
				if v, err := strconv.Atoi(s); err == nil {
					opt.choice = &v
				}
			}

		case strings.HasPrefix(part, "enum="):
			if s := strings.TrimPrefix(part, "enum="); s != "" {
				if v, err := strconv.ParseUint(s, 10, 64); err == nil {
					opt.enum = &v
				}
			}

		case strings.HasPrefix(part, "lb="):
			if s := strings.TrimPrefix(part, "lb="); s != "" {
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					opt.lb = &v
				}
			}

		case strings.HasPrefix(part, "ub="):
			if s := strings.TrimPrefix(part, "ub="); s != "" {
				if v, err := strconv.ParseInt(s, 10, 64); err == nil {
					opt.ub = &v
				}
			}

		case strings.HasPrefix(part, "def="):
			val := strings.Trim(strings.TrimPrefix(part, "def="), `"`)
			opt.def = &val
		}
	}

	return opt
}

// TagCache provides thread-safe caching of parsed struct field tags
// Maps from (struct type, field index) to parsed Tag
type TagCache struct {
	mtx   sync.RWMutex
	cache map[reflect.Type]map[int]*Tag
}

var tagCache = &TagCache{
	cache: make(map[reflect.Type]map[int]*Tag),
}

// GetFieldTag returns the parsed tag for a struct field, using a cache to avoid re-parsing
// Lookup key is (struct type, field index) to ensure unique identification
func GetFieldTag(field reflect.StructField, parent reflect.Type, index int) *Tag {
	tagCache.mtx.RLock()
	if typeCache, exists := tagCache.cache[parent]; exists {
		if tag, exists := typeCache[index]; exists {
			tagCache.mtx.RUnlock()
			return tag
		}
	}
	tagCache.mtx.RUnlock()

	// Not in cache, parse and cache it
	tagStr := field.Tag.Get(TAG_KEY)
	parsed := parseTag(tagStr)

	tagCache.mtx.Lock()
	if _, exists := tagCache.cache[parent]; !exists {
		tagCache.cache[parent] = make(map[int]*Tag)
	}
	tagCache.cache[parent][index] = parsed
	tagCache.mtx.Unlock()

	return parsed
}
