package main

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pkg/errors"
)

// PreampLine is the line for preamp.
const PreampLine = "Preamp: %f dB"

// EQLine              index        freq       gain    quality
const EQLine = `Filter %d: ON PK Fc %d Hz Gain %f dB Q %f`

type Band struct {
	Type      string  `json:"type"`        // Bell
	Mode      string  `json:"mode"`        // RLC (BT)
	Slope     string  `json:"slope"`       // x1
	Solo      bool    `json:"solo,string"` // false
	Mute      bool    `json:"mute,string"` // false
	Gain      float64 `json:"gain,string"`
	Frequency uint16  `json:"frequency,string"`
	Quality   float64 `json:"q,string"`
}

func NewBand() Band {
	return Band{
		Type:  "Bell",
		Mode:  "RLC (BT)",
		Slope: "1x",
	}
}

func ParsePreamp(preampLine string) (preamp float64, err error) {
	_, err = fmt.Sscanf(preampLine, PreampLine, &preamp)
	return
}

func ParseBand(eqLine string) (band Band, i int, err error) {
	band = NewBand()

	_, err = fmt.Sscanf(
		eqLine, EQLine,
		&i, &band.Frequency, &band.Gain, &band.Quality,
	)
	if err != nil {
		return Band{}, 0, err
	}

	// 0 index.
	i--

	return
}

type EqualizerJSON struct {
	Output struct {
		Equalizer struct {
			State         bool    `json:"state,string"`
			Mode          string  `json:"mode"` // IIR
			NumBands      int     `json:"num-bands,string"`
			SplitChannels bool    `json:"split-channels,string"`
			InputGain     float64 `json:"input-gain,omitempty"`

			Left  map[string]Band `json:"left"`
			Right map[string]Band `json:"right"`
		} `json:"equalizer"`
	} `json:"output"`
}

// NewIIREqualizer returns a new single-channel IIR equalizer that's activated.
func NewIIREqualizer() EqualizerJSON {
	output := EqualizerJSON{}
	output.Output.Equalizer.State = true
	output.Output.Equalizer.Mode = "IIR"
	output.Output.Equalizer.Left = map[string]Band{}
	output.Output.Equalizer.Right = output.Output.Equalizer.Left
	return output
}

func (eq *EqualizerJSON) AddBand(i int, b Band) {
	eq.Output.Equalizer.Left[fmt.Sprintf("band%d", i)] = b
	eq.Output.Equalizer.NumBands = len(eq.Output.Equalizer.Left)
}

func (eq *EqualizerJSON) SetPreamp(preamp float64) {
	eq.Output.Equalizer.InputGain = preamp
}

func (eq EqualizerJSON) Patch(original []byte) ([]byte, error) {
	newEQ, err := json.Marshal(eq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal new EQ")
	}

	return jsonpatch.MergePatch(original, newEQ)
}
