package data

import (
	"encoding/json"
	"time"

	cjson "github.com/tent/canonical-json-go"
)

type SignedTimestamp struct {
	Signatures []Signature
	Signed     Timestamp
	Dirty      bool
}

type Timestamp struct {
	Type    string    `json:"_type"`
	Version int       `json:"version"`
	Expires time.Time `json:"expires"`
	Meta    Files     `json:"meta"`
}

func NewTimestamp() (*SignedTimestamp, error) {
	return &SignedTimestamp{
		Signatures: make([]Signature, 0),
		Signed: Timestamp{
			Type:    TUFTypes["timestamp"],
			Version: 0,
			Expires: DefaultExpires("timestamp"),
			Meta:    make(Files),
		},
	}, nil
}

func (ts SignedTimestamp) ToSigned() (*Signed, error) {
	s, err := cjson.Marshal(ts.Signed)
	if err != nil {
		return nil, err
	}
	signed := json.RawMessage{}
	err = signed.UnmarshalJSON(s)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(ts.Signatures))
	copy(sigs, ts.Signatures)
	return &Signed{
		Signatures: sigs,
		Signed:     signed,
	}, nil
}

func TimestampFromSigned(s *Signed) (*SignedTimestamp, error) {
	ts := Timestamp{}
	err := json.Unmarshal(s.Signed, &ts)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(s.Signatures))
	copy(sigs, s.Signatures)
	return &SignedTimestamp{
		Signatures: sigs,
		Signed:     ts,
	}, nil
}
