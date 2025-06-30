package protoutil

import "google.golang.org/protobuf/types/known/wrapperspb"

func UnwrapString(w *wrapperspb.StringValue) *string {
	if w == nil {
		return nil
	}
	return &w.Value
}

func UnwrapBool(w *wrapperspb.BoolValue) *bool {
	if w == nil {
		return nil
	}
	return &w.Value
}

func UnwrapUint64(w *wrapperspb.UInt64Value) *uint64 {
	if w == nil {
		return nil
	}
	return &w.Value
}

func UnwrapInt64(w *wrapperspb.Int64Value) *int64 {
	if w == nil {
		return nil
	}
	return &w.Value
}

func UnwrapDouble(w *wrapperspb.DoubleValue) *float64 {
	if w == nil {
		return nil
	}
	return &w.Value
}
