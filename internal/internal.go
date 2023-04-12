package internal

import (
	"time"
)

const (
	// app constants
	REISSUE_LIFO      = "LIFO"
	REISSUE_FIFO      = "FIFO"
	REISSUE_LOWFIRST  = "Lowest-First"
	REISSUE_HIGHFIRST = "Highest-First"
	REISSUE_RANDOM    = "Random"

	INCREMENT_NEXT     = "Next"
	INCREMENT_NEXTEVEN = "Next-Even"
	INCREMENT_NEXTODD  = "Next-Odd"
	INCREMENT_STEP     = "Step"

	// application default config (so YAML config can be optional)
	RECLAIMKEYS         = false
	REISSUEPOLICY       = REISSUE_FIFO
	INDEXPAD            = "0"
	MININDEX            = 0
	MAXINDEX            = 0
	INCREMENTPOLICE     = INCREMENT_NEXT
	INCREMENTSTEP       = 1
	REQUIRECONFIRMATION = false
	CONFIRMDEADLINE     = 3600 * time.Second

	CTRUE = "true"

	// config
	CONFIGFILE = "keymaker.yaml"
	APIVERSION = "/api/1.0"

	// server
	SRV_HOST = ""
	SRV_PORT = "18652"

	// tls configuration
	TLS_FOLDER    = "certs"
	TLS_ORG       = "Spoon Boy"
	TLS_VALID_FOR = 365 * 24 * time.Hour
)

// StoreReadWriter interface defines the behaviour of a store
type StoreReadWriter interface {
	InitSequence(seqPrefix SeqPrefix, config SeqConfig, begin int) error
	Read(seqPrefix SeqPrefix) (Sequence, error)
	Write(seqPrefix SeqPrefix, payload Sequence) error
}

// SeqPrefix is string type it values of which hold a sequence
// prefix (the part used to identify the sequence group
type SeqPrefix string

// AppConfig holds the config we need to pass about
type AppConfig map[string]string

// Reservation is an ID and timestamp reserved that must be confirmed with
// an additional request to comfirm the sequence is being used
// Reservations which are not confirmed will be reclaimed
// periodically using and made available for reuse if
// configuration supports this
type Reservation struct {
	Index     int
	Timestamp time.Time
}

// Sequence represents everything we need to capture in order
// to maintain a sequence with its configuration metadata as
// well as use it making no assumptions about how data is stored
type Sequence struct {
	Config    SeqConfig     `json:"config"`
	Last      int           `json:"last"`
	Reclaimed []int         `json:"reclaimed,omitempty"`
	Reserved  []Reservation `json:"reserved,omitempty"`
}

// SeqConfig stores the sequence configuration if not is supplied
// when the sequence is initialised the configuration supplied in
// the yaml config file will be applied, if no yaml configuration
// is supplied we use the constant defined defaults
type SeqConfig struct {
	ReclaimKeys     bool          `json:"reclaimKeys,omitempty"`
	ReissuePolicy   string        `json:"reissuePolicy,omitempty"`
	IndexPad        string        `json:"indexPad,omitempty"`
	MinIndex        int           `json:"minIndex,omitempty"`
	MaxIndex        int           `json:"maxIndex,omitempty"`
	IncrementPolicy string        `json:"incrementPolicy,omitempty"`
	IncrementStep   int           `json:"incrementStep,omitempty"`
	RequireConfirm  bool          `json:"requireConfirm,omitempty"`
	ConfirmDeadline time.Duration `json:"confirmDeadline,omitempty"`
}
