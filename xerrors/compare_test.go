package xerrors_test

import (
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

type barError struct {
	xerrors.Wrapping
}

func (barError) Error() string { return "bar" }

func TestSimilar(t *testing.T) {
	scenarios := []struct {
		name          string
		err1          error
		err2          error
		expectedEqual bool
	}{
		{
			name:          "equalNew",
			err1:          xerrors.New("foo"),
			err2:          xerrors.New("foo"),
			expectedEqual: true,
		},
		{
			name:          "differentNew",
			err1:          xerrors.New("foo"),
			err2:          xerrors.New("bar"),
			expectedEqual: false,
		},
		{
			name:          "equalWrapWithoutFrames",
			err1:          xerrors.Wrap("bar", xerrors.New("foo"), xerrors.OmitFrame()),
			err2:          xerrors.Wrap("bar", xerrors.New("foo"), xerrors.OmitFrame()),
			expectedEqual: true,
		},
		{
			name:          "similarWrapWithFrames",
			err1:          xerrors.Wrap("bar", xerrors.New("foo")),
			err2:          xerrors.Wrap("bar", xerrors.New("foo")),
			expectedEqual: true,
		},
		{
			name:          "differentWrapInner",
			err1:          xerrors.Wrap("bar", xerrors.New("foo1"), xerrors.OmitFrame()),
			err2:          xerrors.Wrap("bar", xerrors.New("foo2"), xerrors.OmitFrame()),
			expectedEqual: false,
		},
		{
			name:          "differentWrapOuter",
			err1:          xerrors.Wrap("bar1", xerrors.New("foo"), xerrors.OmitFrame()),
			err2:          xerrors.Wrap("bar2", xerrors.New("foo"), xerrors.OmitFrame()),
			expectedEqual: false,
		},
		{
			name:          "differentTypesSameError",
			err1:          xerrors.Wrap("foo", nil, xerrors.OmitFrame()),
			err2:          xerrors.New("foo"),
			expectedEqual: false,
		},
		{
			name:          "mismatchedOrder",
			err1:          xerrors.Wrap("foo", xerrors.Wrap("bar", nil, xerrors.OmitFrame()), xerrors.OmitFrame()),
			err2:          xerrors.Wrap("bar", xerrors.Wrap("foo", nil, xerrors.OmitFrame()), xerrors.OmitFrame()),
			expectedEqual: false,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario

		t.Run(scenario.name, func(t *testing.T) {
			if equal := xerrors.Similar(scenario.err1, scenario.err2); equal != scenario.expectedEqual {
				t.Fatalf("mismatched output, expected %t got %t", scenario.expectedEqual, equal)
			}

			if !xerrors.Similar(scenario.err1, scenario.err1) {
				t.Fatal("an error must always be similar to itself")
			}

			if !xerrors.Similar(scenario.err2, scenario.err2) {
				t.Fatal("an error must always be similar to itself")
			}

			if scenario.expectedEqual {
				if !xerrors.Contains(scenario.err1, scenario.err2) || !xerrors.Contains(scenario.err2, scenario.err1) {
					t.Fatal("similar errors must also be contained by each other")
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	scenarios := []struct {
		name             string
		err1             error
		err2             error
		expectedContains bool
	}{
		{
			name:             "containedSentinel",
			err1:             xerrors.Wrap("bar", xerrors.New("foo")),
			err2:             xerrors.New("foo"),
			expectedContains: true,
		},
		{
			name:             "differentSentinel",
			err1:             xerrors.Wrap("bar", xerrors.New("foo")),
			err2:             xerrors.New("bar"),
			expectedContains: false,
		},
		{
			name:             "containedWrapWithoutFrames",
			err1:             xerrors.Wrap("bar", xerrors.New("foo"), xerrors.OmitFrame()),
			err2:             xerrors.Wrap("bar", nil, xerrors.OmitFrame()),
			expectedContains: true,
		},
		{
			name:             "containedWrapWithFrames",
			err1:             xerrors.Wrap("bar", xerrors.New("foo")),
			err2:             xerrors.Wrap("bar", nil),
			expectedContains: true,
		},
		{
			name:             "differentWrap",
			err1:             xerrors.Wrap("bar1", xerrors.New("foo"), xerrors.OmitFrame()),
			err2:             xerrors.Wrap("bar2", nil, xerrors.OmitFrame()),
			expectedContains: false,
		},
		{
			name:             "mismatchedOrder",
			err1:             xerrors.Wrap("foo", xerrors.Wrap("bar", nil, xerrors.OmitFrame()), xerrors.OmitFrame()),
			err2:             xerrors.Wrap("bar", xerrors.Wrap("foo", nil, xerrors.OmitFrame()), xerrors.OmitFrame()),
			expectedContains: false,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario

		t.Run(scenario.name, func(t *testing.T) {
			if contains := xerrors.Contains(scenario.err1, scenario.err2); contains != scenario.expectedContains {
				t.Fatalf("mismatched output, expected %t got %t", scenario.expectedContains, contains)
			}
		})
	}
}
