package types


// FaultMessage 模拟 protobuf 生成的 msgpb.FaultMessage
type FaultMessage struct {
	Hint string
}


func (f *FaultMessage) GetHint() string {
	if f == nil {
		return ""
	}
	return f.Hint
}

func (f *FaultMessage) SetHint(s string) *FaultMessage {
    return &FaultMessage{
        Hint: s,
    }
}
func (f *FaultMessage) Clone() *FaultMessage {
    return &FaultMessage{
        Hint: f.Hint,
    }
}