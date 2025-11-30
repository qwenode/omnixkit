package tyy

import "github.com/qwenode/omnixkit/examples/kitfault/types"

type Response struct {
    Data         string
    FaultMessage *types.FaultMessage
}
func (r *Response) GetFaultMessage() *types.FaultMessage {
	return r.FaultMessage
}

func (r *Response) SetFaultMessage(msg *types.FaultMessage) {
	r.FaultMessage = msg
}