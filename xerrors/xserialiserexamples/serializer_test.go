package xserialiserexamples_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors/xserialiserexamples"
)

func TestSerializers(t *testing.T) {
	type serializerOutputs struct {
		colonBasicSerialised    string
		colonDetailSerialised   string
		basicKeyValueSerialised string
		jsonKeyValueSerialised  string
		frameOnlySerialised     string
	}

	scenarios := []struct {
		name            string
		err             error
		expectedOutputs serializerOutputs
	}{
		{
			name: "nonWrapped",
			err:  xerrors.New("msg"),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "msg",
				colonDetailSerialised:   "msg",
				basicKeyValueSerialised: "?-msg",
				jsonKeyValueSerialised:  `{"unknown_0":"msg"}`,
				frameOnlySerialised:     "",
			},
		},
		{
			name: "singleWrapped",
			err: xerrors.Wrap(
				"wrapping_msg",
				xerrors.New("cause_msg"),
				xerrors.OmitFrame(),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "wrapping_msg: cause_msg",
				colonDetailSerialised:   "wrapping_msg: cause_msg",
				basicKeyValueSerialised: "?-wrapping_msg ?-cause_msg",
				jsonKeyValueSerialised:  `{"unknown_1":"wrapping_msg","unknown_0":"cause_msg"}`,
				frameOnlySerialised:     "",
			},
		},
		{
			name: "doubleWrapped",
			err: xerrors.Wrap(
				"wrapping_msg_2",
				xerrors.Wrap(
					"wrapping_msg_1",
					xerrors.New("cause_msg"),
					xerrors.OmitFrame(),
				),
				xerrors.OmitFrame(),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "wrapping_msg_2: wrapping_msg_1: cause_msg",
				colonDetailSerialised:   "wrapping_msg_2: wrapping_msg_1: cause_msg",
				basicKeyValueSerialised: "?-wrapping_msg_2 ?-wrapping_msg_1 ?-cause_msg",
				jsonKeyValueSerialised:  `{"unknown_2":"wrapping_msg_2","unknown_1":"wrapping_msg_1","unknown_0":"cause_msg"}`,
				frameOnlySerialised:     "",
			},
		},
		{
			name: "basicKeyValueError",
			err: xserialiserexamples.NewBasicKeyValueError(
				[][2]string{{"key_foo", "value_foo"}, {"key_bar", "value_bar"}},
				xerrors.New("msg"),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "key_foo-value_foo key_bar-value_bar: msg",
				colonDetailSerialised:   "key_foo-value_foo key_bar-value_bar: msg",
				basicKeyValueSerialised: "key_foo-value_foo key_bar-value_bar ?-msg",
				jsonKeyValueSerialised:  `{"key_foo":"value_foo","key_bar":"value_bar","unknown_0":"msg"}`,
				frameOnlySerialised:     "",
			},
		},
		{
			name: "jsonKeyValueError",
			err: xserialiserexamples.NewJSONKeyValueError(
				[][2]string{{"key_foo", "value_foo"}, {"key_bar", "value_bar"}},
				xerrors.New("msg"),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    `{"key_foo":"value_foo","key_bar":"value_bar"}: msg`,
				colonDetailSerialised:   `{"key_foo":"value_foo","key_bar":"value_bar"}: msg`,
				basicKeyValueSerialised: "key_foo-value_foo key_bar-value_bar ?-msg",
				jsonKeyValueSerialised:  `{"key_foo":"value_foo","key_bar":"value_bar","unknown_0":"msg"}`,
				frameOnlySerialised:     "",
			},
		},
		{
			name: "singleWrappedWithFrame",
			err: xerrors.Wrap(
				"wrapping_msg",
				newPretendFrameError(
					"my/pkg/foobar.myMethod",
					"/my/home/my/gopath/src/my/pkg/foobar/myfile.go",
					100,
					xerrors.New("cause_msg"),
				),
				xerrors.OmitFrame(),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "wrapping_msg: cause_msg",
				colonDetailSerialised:   "wrapping_msg(foobar.myMethod:myfile.go:100): cause_msg",
				basicKeyValueSerialised: "?-wrapping_msg ?-cause_msg",
				jsonKeyValueSerialised:  `{"unknown_1":"wrapping_msg","unknown_0":"cause_msg"}`,
				frameOnlySerialised:     "my/pkg/foobar.myMethod:/my/home/my/gopath/src/my/pkg/foobar/myfile.go:100",
			},
		},
		{
			name: "doubleWrappedWithFrame",
			err: xerrors.Wrap(
				"wrapping_msg_2",
				newPretendFrameError(
					"my/pkg/foobar2.myMethod2",
					"/my/home/my/gopath/src/my/pkg/foobar2/myfile2.go",
					200,
					xerrors.Wrap(
						"wrapping_msg_1",
						newPretendFrameError(
							"my/pkg/foobar1.myMethod1",
							"/my/home/my/gopath/src/my/pkg/foobar1/myfile1.go",
							100,
							xerrors.New("cause_msg"),
						),
						xerrors.OmitFrame(),
					),
				),
				xerrors.OmitFrame(),
			),
			expectedOutputs: serializerOutputs{
				colonBasicSerialised:    "wrapping_msg_2: wrapping_msg_1: cause_msg",
				colonDetailSerialised:   "wrapping_msg_2(foobar2.myMethod2:myfile2.go:200): wrapping_msg_1(foobar1.myMethod1:myfile1.go:100): cause_msg",
				basicKeyValueSerialised: "?-wrapping_msg_2 ?-wrapping_msg_1 ?-cause_msg",
				jsonKeyValueSerialised:  `{"unknown_2":"wrapping_msg_2","unknown_1":"wrapping_msg_1","unknown_0":"cause_msg"}`,
				frameOnlySerialised: "my/pkg/foobar2.myMethod2:/my/home/my/gopath/src/my/pkg/foobar2/myfile2.go:200" +
					"\n\t" + "my/pkg/foobar1.myMethod1:/my/home/my/gopath/src/my/pkg/foobar1/myfile1.go:100",
			},
		},
	}

	formatterScenarios := []struct {
		serializerName    string
		serializerFactory func() xerrors.Serializer
		expectedOutput    func(serializerOutputs) string
	}{
		{
			serializerName:    "colonBasicSerialised",
			serializerFactory: xerrors.NewColonBasicSerializer,
			expectedOutput:    func(outs serializerOutputs) string { return outs.colonBasicSerialised },
		},
		{
			serializerName:    "colonDetailSerialised",
			serializerFactory: xerrors.NewColonDetailedSerializer,
			expectedOutput:    func(outs serializerOutputs) string { return outs.colonDetailSerialised },
		},
		{
			serializerName:    "basicKeyValueSerialised",
			serializerFactory: xserialiserexamples.NewBasicKeyValueSerializer,
			expectedOutput:    func(outs serializerOutputs) string { return outs.basicKeyValueSerialised },
		},
		{
			serializerName:    "jsonKeyValueSerialised",
			serializerFactory: xserialiserexamples.NewJSONKeyValueSerializer,
			expectedOutput:    func(outs serializerOutputs) string { return outs.jsonKeyValueSerialised },
		},
		{
			serializerName:    "frameOnlySerialised",
			serializerFactory: xserialiserexamples.NewFrameOnlySerializer,
			expectedOutput:    func(outs serializerOutputs) string { return outs.frameOnlySerialised },
		},
	}

	for _, formatterScenario := range formatterScenarios {
		formatterScenario := formatterScenario

		t.Run(formatterScenario.serializerName, func(t *testing.T) {
			printer := xerrors.NewPrinter(formatterScenario.serializerFactory)

			for _, scenario := range scenarios {
				scenario := scenario

				t.Run(scenario.name, func(t *testing.T) {
					buf := bytes.Buffer{}
					if err := printer.Write(&buf, scenario.err); err != nil {
						t.Fatalf("error serialising error: %s", err)
					}

					expectedOut := formatterScenario.expectedOutput(scenario.expectedOutputs)
					if buf.String() != expectedOut {
						t.Fatalf("mismatched output, expected %q got %q", expectedOut, buf.String())
					}
				})
			}
		})
	}
}

type pretendFrameError struct {
	function, file string
	line           int
	xerrors.Wrapping
}

func (err *pretendFrameError) Error() string {
	out := ""
	if err.function != "" {
		out += err.function + ":"
	}
	if err.file != "" {
		out += err.file + ":" + strconv.Itoa(err.line) + " "
	}
	if err.function != "" || err.file != "" {
		out = out[:len(out)-1]
	}
	return out
}

func (err *pretendFrameError) FrameLocation() (string, string, int) {
	return err.function, err.file, err.line
}

func newPretendFrameError(function, file string, line int, err error) error {
	return &pretendFrameError{
		function: function,
		file:     file,
		line:     line,
		Wrapping: xerrors.NewWrapping(err, xerrors.OmitFrame()),
	}
}

var _ xerrors.FrameError = (*pretendFrameError)(nil)
