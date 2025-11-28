package kitvalidate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

type ValidationError struct {
	Field   string
	Message string
}

type ErrorDetailBuilder interface {
	BuildErrorDetail(errors []ValidationError) (*connect.ErrorDetail, error)
}

type ErrorDetailBuilderFunc func(errors []ValidationError) (*connect.ErrorDetail, error)

func (f ErrorDetailBuilderFunc) BuildErrorDetail(errors []ValidationError) (*connect.ErrorDetail, error) {
	return f(errors)
}

// An Option configures an [Interceptor].
type Option interface {
	apply(*Interceptor)
}

// WithValidator configures the [Interceptor] to use a customized
// [protovalidate.Validator]. See [protovalidate.ValidatorOption] for the range
// of available customizations.
func WithValidator(validator protovalidate.Validator) Option {
	return optionFunc(
		func(i *Interceptor) {
			i.validator = validator
		},
	)
}

// WithErrorDetailBuilder configures the [Interceptor] to use a custom
// [ErrorDetailBuilder] for building error details from validation errors.
func WithErrorDetailBuilder(builder ErrorDetailBuilder) Option {
	return optionFunc(
		func(i *Interceptor) {
			i.errorDetailBuilder = builder
		},
	)
}

type Interceptor struct {
	validator          protovalidate.Validator
	errorDetailBuilder ErrorDetailBuilder
}



// 使用示例:
// 创建带有自定义 ErrorDetailBuilder 的 interceptor
//interceptor, err := kitrpc.NewInterceptor(
//	kitrpc.WithErrorDetailBuilder(NewValidateErrorDetailBuilder()),
//)
//if err != nil {
//	// 处理错误
//}
// 然后在 connect handler 中使用这个 interceptor

// func NewValidateErrorDetailBuilder() kitrpc.ErrorDetailBuilder {
//	return kitrpc.ErrorDetailBuilderFunc(func(errors []kitrpc.ValidationError) (*connect.ErrorDetail, error) {
//		fieldErr := &msgpb.ValidateMessages{
//			Fields: make([]*msgpb.ValidateMessage, 0, len(errors)),
//		}
//		for _, e := range errors {
//			fieldErr.Fields = append(fieldErr.Fields, &msgpb.ValidateMessage{
//				Field:   e.Field,
//				Message: e.Message,
//			})
//		}
//		return connect.NewErrorDetail(fieldErr)
//	})
// }
func NewInterceptor(opts ...Option) (*Interceptor, error) {
	var interceptor Interceptor
	for _, opt := range opts {
		opt.apply(&interceptor)
	}

	if interceptor.validator == nil {
		validator, err := protovalidate.New()
		if err != nil {
			return nil, fmt.Errorf("construct validator: %w", err)
		}
		interceptor.validator = validator
	}

	return &interceptor, nil
}

// WrapUnary implements connect.Interceptor.
func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := validate(i.validator, i.errorDetailBuilder, req.Any()); err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

// WrapStreamingClient implements connect.Interceptor.
func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return &streamingClientInterceptor{
			StreamingClientConn:  next(ctx, spec),
			validator:            i.validator,
			errorDetailBuilder:   i.errorDetailBuilder,
		}
	}
}

// WrapStreamingHandler implements connect.Interceptor.
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(
			ctx, &streamingHandlerInterceptor{
				StreamingHandlerConn: conn,
				validator:            i.validator,
				errorDetailBuilder:   i.errorDetailBuilder,
			},
		)
	}
}

type streamingClientInterceptor struct {
	connect.StreamingClientConn

	validator          protovalidate.Validator
	errorDetailBuilder ErrorDetailBuilder
}

func (s *streamingClientInterceptor) Send(msg any) error {
	if err := validate(s.validator, s.errorDetailBuilder, msg); err != nil {
		return err
	}
	return s.StreamingClientConn.Send(msg)
}

type streamingHandlerInterceptor struct {
	connect.StreamingHandlerConn

	validator          protovalidate.Validator
	errorDetailBuilder ErrorDetailBuilder
}

func (s *streamingHandlerInterceptor) Receive(msg any) error {
	if err := s.StreamingHandlerConn.Receive(msg); err != nil {
		return err
	}
	return validate(s.validator, s.errorDetailBuilder, msg)
}

type optionFunc func(*Interceptor)

func (f optionFunc) apply(i *Interceptor) { f(i) }

func validate(validator protovalidate.Validator, builder ErrorDetailBuilder, msg any) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("expected proto.Message, got %T", msg)
	}
	err := validator.Validate(protoMsg)
	if err == nil {
		return nil
	}

	connectErr := connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed"))

	if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
		var validationErrors []ValidationError
		for _, violation := range validationErr.ToProto().Violations {
			// 处理字段级验证错误
			if violation.GetField() != nil && len(violation.GetField().GetElements()) > 0 {
				for _, element := range violation.GetField().GetElements() {
					validationErrors = append(
						validationErrors, ValidationError{
							Field:   element.GetFieldName(),
							Message: violation.GetMessage(),
						},
					)
					break
				}
			} else {
				// 处理消息级验证错误（如 CEL 表达式）
				validationErrors = append(
					validationErrors, ValidationError{
						Field:   "", // 消息级错误没有具体字段
						Message: violation.GetMessage(),
					},
				)
			}
		}

		if len(validationErrors) > 0 {
			if builder != nil {
				if detail, err := builder.BuildErrorDetail(validationErrors); err == nil && detail != nil {
					connectErr.AddDetail(detail)
				}
			} else {
				var errMsgs []string
				for _, e := range validationErrors {
					if e.Field != "" {
						errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
					} else {
						errMsgs = append(errMsgs, e.Message)
					}
				}
				connectErr = connect.NewError(connect.CodeInvalidArgument, errors.New(strings.Join(errMsgs, "; ")))
			}
		}
	} else {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connectErr
}
