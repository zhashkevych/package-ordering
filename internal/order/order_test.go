package order

import (
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	// Edge example: packs {23,31,53}; amount: 500_000
	result := CalculatePacks(500000, []int{23, 31, 53})
	expected := map[int]int{23: 2, 31: 7, 53: 9429}

	if !compareMaps(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// Simple: 1, 250
	result = CalculatePacks(1, []int{250, 500})
	expected = map[int]int{250: 1}

	if !compareMaps(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

	// Simple: 501, {250,500}
	result = CalculatePacks(501, []int{250, 500})
	expected = map[int]int{500: 1, 250: 1}

	if !compareMaps(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}

}

func compareMaps(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
