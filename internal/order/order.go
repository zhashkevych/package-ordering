package order

import (
	"sort"
)

// PackResult maps pack size to how many of that pack are used
type PackResult map[int]int

// CalculatePacks returns the optimal pack breakdown for an order.
// packSizes does not need to be sorted.
func CalculatePacks(amount int, packSizes []int) PackResult {
	// Sort packSizes largest to smallest
	sorted := append([]int{}, packSizes...)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))

	type state struct {
		overflow int
		packs    int
		packSet  PackResult
	}

	// Large enough to cover any practical case
	maxA := amount + sorted[0]
	dp := make([]*state, maxA+1)
	dp[0] = &state{overflow: -amount, packs: 0, packSet: make(PackResult)}

	for _, sz := range sorted {
		for i := sz; i <= maxA; i++ {
			prev := dp[i-sz]
			if prev == nil {
				continue
			}
			currOverflow := prev.overflow + sz
			currPacks := prev.packs + 1

			// If this is the first way to reach 'i' or a strictly better way, take it
			if dp[i] == nil ||
				currOverflow < dp[i].overflow ||
				(currOverflow == dp[i].overflow && currPacks < dp[i].packs) {
				// Copy and add this pack
				newSet := make(PackResult)
				for k, v := range prev.packSet {
					newSet[k] = v
				}
				newSet[sz]++
				dp[i] = &state{overflow: currOverflow, packs: currPacks, packSet: newSet}
			}
		}
	}

	// Best answer for requested amount or more (overflow >= 0)
	best := (*state)(nil)
	for i := amount; i <= maxA; i++ {
		if dp[i] != nil && dp[i].overflow >= 0 &&
			(best == nil || dp[i].overflow < best.overflow ||
				(dp[i].overflow == best.overflow && dp[i].packs < best.packs)) {
			best = dp[i]
		}
	}

	if best != nil {
		return best.packSet
	}
	return nil
}
