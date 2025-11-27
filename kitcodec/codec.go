package kitcodec

import (
    "bytes"
    "encoding/json"
    "fmt"

    "connectrpc.com/connect"

    "google.golang.org/protobuf/encoding/protojson"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/runtime/protoiface"
)

const (
    codecNameProto           = "proto"
    codecNameJSON            = "json"
    codecNameJSONCharsetUTF8 = codecNameJSON + "; charset=utf-8"
)

// Codec marshals structs (typically generated from a schema) to and from bytes.
type Codec interface {
    // Name returns the name of the Codec.
    //
    // This may be used as part of the Content-Type within HTTP. For example,
    // with gRPC this is the content subtype, so "application/grpc+proto" will
    // map to the Codec with name "proto".
    //
    // Names must not be empty.
    Name() string
    // Marshal marshals the given message.
    //
    // Marshal may expect a specific type of message, and will error if this type
    // is not given.
    Marshal(any) ([]byte, error)
    // Unmarshal unmarshals the given message.
    //
    // Unmarshal may expect a specific type of message, and will error if this
    // type is not given.
    Unmarshal([]byte, any) error
}

// marshalAppender is an extension to Codec for appending to a byte slice.
type marshalAppender interface {
    Codec

    // MarshalAppend marshals the given message and appends it to the given
    // byte slice.
    //
    // MarshalAppend may expect a specific type of message, and will error if
    // this type is not given.
    MarshalAppend([]byte, any) ([]byte, error)
}

// stableCodec is an extension to Codec for serializing with stable output.
type stableCodec interface {
    Codec

    // MarshalStable marshals the given message with stable field ordering.
    //
    // MarshalStable should return the same output for a given input. Although
    // it is not guaranteed to be canonicalized, the marshalling routine for
    // MarshalStable will opt for the most normalized output available for a
    // given serialization.
    //
    // For practical reasons, it is possible for MarshalStable to return two
    // different results for two inputs considered to be "equal" in their own
    // domain, and it may change in the future with codec updates, but for
    // any given concrete value and any given version, it should return the
    // same output.
    MarshalStable(any) ([]byte, error)

    // IsBinary returns true if the marshalled data is binary for this codec.
    //
    // If this function returns false, the data returned from Marshal and
    // MarshalStable are considered valid text and may be used in contexts
    // where text is expected.
    IsBinary() bool
}

type protoJSONCodec struct {
    name string
}

var _ Codec = (*protoJSONCodec)(nil)

func WithProtoJSON() connect.Option {
    return connect.WithCodec(&protoJSONCodec{codecNameJSON})
}
func (c *protoJSONCodec) Name() string { return c.name }

func (c *protoJSONCodec) Marshal(message any) ([]byte, error) {
    protoMessage, ok := message.(proto.Message)
    if !ok {
        return nil, errNotProto(message)
    }
    return protojson.MarshalOptions{}.Marshal(protoMessage)
}

func (c *protoJSONCodec) MarshalAppend(dst []byte, message any) ([]byte, error) {
    protoMessage, ok := message.(proto.Message)
    if !ok {
        return nil, errNotProto(message)
    }
    return protojson.MarshalOptions{}.MarshalAppend(dst, protoMessage)
}

func (c *protoJSONCodec) Unmarshal(binary []byte, message any) error {
    protoMessage, ok := message.(proto.Message)
    if !ok {
        return errNotProto(message)
    }
    if len(binary) == 0 {
        binary = []byte("{}")
    }
    // Discard unknown fields so clients and servers aren't forced to always use
    // exactly the same version of the schema.
    options := protojson.UnmarshalOptions{DiscardUnknown: true}
    err := options.Unmarshal(binary, protoMessage)
    if err != nil {
        return fmt.Errorf("unmarshal into %T: %w", message, err)
    }
    return nil
}

func (c *protoJSONCodec) MarshalStable(message any) ([]byte, error) {
    // protojson does not offer a "deterministic" field ordering, but fields
    // are still ordered consistently by their index. However, protojson can
    // output inconsistent whitespace for some reason, therefore it is
    // suggested to use a formatter to ensure consistent formatting.
    // https://github.com/golang/protobuf/issues/1373
    messageJSON, err := c.Marshal(message)
    if err != nil {
        return nil, err
    }
    compactedJSON := bytes.NewBuffer(messageJSON[:0])
    if err = json.Compact(compactedJSON, messageJSON); err != nil {
        return nil, err
    }
    return compactedJSON.Bytes(), nil
}

func (c *protoJSONCodec) IsBinary() bool {
    return false
}

// readOnlyCodecs is a read-only interface to a map of named codecs.
type readOnlyCodecs interface {
    // Get gets the Codec with the given name.
    Get(string) Codec
    // Protobuf gets the user-supplied protobuf codec, falling back to the default
    // implementation if necessary.
    //
    // This is helpful in the gRPC protocol, where the wire protocol requires
    // marshaling protobuf structs to binary even if the RPC procedures were
    // generated from a different IDL.
    Protobuf() Codec
    // Names returns a copy of the registered codec names. The returned slice is
    // safe for the caller to mutate.
    Names() []string
}

func errNotProto(message any) error {
    if _, ok := message.(protoiface.MessageV1); ok {
        return fmt.Errorf("%T uses protobuf,but use v2/v3", message)
    }
    return fmt.Errorf("%T doesn't implement proto.Message", message)
}
