package testdatas

import "github.com/nasjp/tiny-gotests/internal/processor/testdata/sample"

type Sample struct{}

func Minimum() {}

func ParamOne(a int) {}

func ParamTwo(a int, b int) {}

func ParamInternalStruct(a Sample) {}

func ParamInternalPointer(a *Sample) {}

func ParamExternalStruct(a sample.Sample) {}

func ParamExternalPointer(a *sample.Sample) {}

func ResultOne() int { return 0 }

func ResultTwo() (int, int) { return 0, 0 }

func ResultErr() error { return nil }

func ResultOneAndErr() (int, error) { return 0, nil }

func ResultInternalStruct() Sample { return Sample{} }

func ResultInternalPointer() *Sample { return nil }

func ResultExternalStruct() sample.Sample { return sample.Sample{} }

func ResultExternalPointer() *sample.Sample { return nil }

func (s Sample) ReceiverStruct() {}

func (s *Sample) ReceiverPointer() {}
