package slicer

import (
	"crumbl/utils"
	"errors"
	"math"
	"math/rand"
	"strings"
)

const (
	// MAX_SLICES ...
	MAX_SLICES = 4 // The owner of the data + 3 trustees is optimal as of this version

	// MAX_DELTA is the maximum allowed deltaMax in the system as of this version
	MAX_DELTA = 5

	// MIN_INPUT_SIZE ...
	MIN_INPUT_SIZE = 8 // Input below 8 characters must be left-padded

	// MIN_SLICE_SIZE ...
	MIN_SLICE_SIZE = 2
)

//--- TYPES

// Slicer ...
type Slicer struct {
	NumberOfSlices int
	DeltaMax       int
}

// Slice ...
type Slice string

//--- METHODS

// Apply returns the slices from the passed data
func (s Slicer) Apply(data string) (slices []Slice, err error) {
	splits, err := s.Split(data)
	if err != nil {
		return
	}
	fixedLength := int(math.Floor(float64(len(data)/s.NumberOfSlices))) + s.DeltaMax
	for _, split := range splits {
		slice := utils.LeftPad(split, fixedLength)
		slices = append(slices, Slice(slice))
	}
	if len(slices) != s.NumberOfSlices {
		err = errors.New("wrong number of slices")
		return
	}
	return
}

// Unapply rebuild the original data from the slices
func (s Slicer) Unapply(slices []Slice) (data string, err error) {
	if len(slices) == 0 {
		err = errors.New("impossible to use empty slices")
		return
	}
	var splits []string
	for _, slice := range slices {
		splits = append(splits, utils.Unpad(string(slice)))
	}
	data = strings.Join(splits, "")
	return
}

// Split plits the passed data using a mask
func (s *Slicer) Split(data string) (splits []string, err error) {
	masks, err := s.buildSplitMask(len(data), SeedFor((data)))
	if err != nil {
		return
	}
	for _, m := range masks {
		splits = append(splits, data[m.Start:m.End])
	}
	return
}

// GetDeltaMax ...
func GetDeltaMax(dataLength int, numberOfSlices int) int {
	sliceSize := dataLength / numberOfSlices
	if dataLength <= MIN_INPUT_SIZE || sliceSize <= MIN_SLICE_SIZE {
		return 0
	}
	var deltaMax int
	for delta := 1; delta < MAX_DELTA+1; delta++ {
		deltaMax = delta
		if delta < 2*(sliceSize-MIN_SLICE_SIZE) {
			continue
		} else {
			break
		}
	}
	return deltaMax
}

type mask struct {
	Start int
	End   int
}

func (s *Slicer) buildSplitMask(dataLength int, seed Seed) (masks []mask, err error) {
	dl := float64(dataLength)
	nos := float64(s.NumberOfSlices)
	dm := float64(s.DeltaMax)
	averageSliceLength := math.Floor(dl / nos)
	minLen := math.Max(averageSliceLength-math.Floor(dm/2), math.Floor(dl/(nos+1)+1))
	maxLen := math.Min(averageSliceLength+math.Floor(dm/2), math.Ceil(dl/(nos-1)-1))
	delta := math.Min(dm, maxLen-minLen)
	length := 0
	rand.Seed(int64(seed))
	for dataLength > 0 {
		randomNum := math.Min(math.Floor(rand.Float64()*(math.Min(maxLen, dl)+1-minLen)+minLen), float64(dataLength))
		if randomNum == 0 {
			continue
		}
		b := math.Floor((dl - randomNum) / minLen)
		r := math.Floor(float64((dataLength - int(randomNum)) % int(minLen)))
		if r <= b*delta {
			m := mask{
				Start: length,
				End:   int(math.Min(dl, float64(length)+randomNum)),
			}
			masks = append(masks, m)
			length += int(randomNum)
			dataLength -= int(randomNum)
		}
	}
	if len(masks) == 0 {
		err = errors.New("unable to build split masks")
		return
	}
	return
}