package xerrors_test

import (
	"bytes"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

// beware, the lines in this method are required to test frame wrapping logic
// any formatting changes will break tests and will require changes to the expected line constants
func encodeScenarios() []encodeScenario {
	return []encodeScenario{
		{
			name:              "nonWrapped",
			err:               xerrors.New("msg"),
			expectedBasicOut:  "msg",
			expectedDetailOut: "msg",
		},
		{
			name:              "singleWrapped",
			err:               xerrors.Wrap("wrapping_msg", xerrors.New("cause_msg")),
			expectedBasicOut:  "wrapping_msg: cause_msg",
			expectedDetailOut: "wrapping_msg(xerrors_test.encodeScenarios:default_test.go:24): cause_msg",
		},
		{
			name: "doubleWrapped",
			err: xerrors.Wrap(
				"wrapping_msg_2",
				xerrors.Wrap("wrapping_msg_1", xerrors.New("cause_msg")),
			),
			expectedBasicOut:  "wrapping_msg_2: wrapping_msg_1: cause_msg",
			expectedDetailOut: "wrapping_msg_2(xerrors_test.encodeScenarios:default_test.go:30): wrapping_msg_1(xerrors_test.encodeScenarios:default_test.go:32): cause_msg",
		},
	}
}

type encodeScenario struct {
	name              string
	err               error
	expectedBasicOut  string
	expectedDetailOut string
}

func expectedBasicOutput(scenario encodeScenario) string {
	return scenario.expectedBasicOut
}

func expectedDetailOutput(scenario encodeScenario) string {
	return scenario.expectedDetailOut
}

func testEncode(t *testing.T, encode func(error) string, outputReader func(encodeScenario) string) {
	scenarios := encodeScenarios()

	t.Run("individual", func(t *testing.T) {
		for _, scenario := range scenarios {
			scenario := scenario

			t.Run(scenario.name, func(t *testing.T) {
				if out, expectedOut := encode(scenario.err), outputReader(scenario); out != expectedOut {
					t.Fatalf("expected %q got %q", expectedOut, out)
				}
			})
		}
	})

	t.Run("stress", func(t *testing.T) {
		const stressReps = 10000

		for _, scenario := range scenarios {
			scenario := scenario

			t.Run(scenario.name, func(t *testing.T) {
				wg := sync.WaitGroup{}
				var errorCount int32

				wg.Add(stressReps)
				for i := 0; i < stressReps; i++ {
					go func() {
						if out, expectedOut := encode(scenario.err), outputReader(scenario); out != expectedOut {
							atomic.AddInt32(&errorCount, 1)
						}
						wg.Done()
					}()
				}

				wg.Wait()

				if errorCount != 0 {
					t.Fatalf("expecting no async related errors, found %d/%d", errorCount, stressReps)
				}
			})
		}
	})
}

func TestString(t *testing.T) {
	testEncode(t, xerrors.String, expectedBasicOutput)
}

func TestDetailString(t *testing.T) {
	testEncode(t, xerrors.DetailString, expectedDetailOutput)
}

func testByteEncode(t *testing.T, encode func(error) []byte, outputReader func(encodeScenario) string) {
	noise := bytes.Repeat([]byte{'-'}, 100)

	testEncode(
		t,
		func(err error) string {
			b := encode(err)
			out := string(b)
			copy(b, noise) // writing to the output []byte to prove there is no Pool related conflict
			return out
		},
		outputReader,
	)
}

func TestBytes(t *testing.T) {
	testByteEncode(t, xerrors.Bytes, expectedBasicOutput)
}

func TestDetailBytes(t *testing.T) {
	testByteEncode(t, xerrors.DetailBytes, expectedDetailOutput)
}

func BenchmarkString(b *testing.B) {
	scenarios := encodeScenarios()

	for _, scenario := range scenarios {
		scenario := scenario

		b.ResetTimer()

		b.Run(scenario.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// not bothering to check the output, already covered in the tests
				_ = xerrors.String(scenario.err)
			}
		})
	}
}

func BenchmarkDetailString(b *testing.B) {
	scenarios := encodeScenarios()

	for _, scenario := range scenarios {
		scenario := scenario

		b.ResetTimer()

		b.Run(scenario.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// not bothering to check the output, already covered in the tests
				_ = xerrors.DetailString(scenario.err)
			}
		})
	}
}
