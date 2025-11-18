package kitrpc

import (
    "errors"

    "connectrpc.com/connect"
)

func NewInvalidArgument(message string) *connect.Error {
    return NewInvalidArgumentErr(errors.New(message))
}
func NewInvalidArgumentErr(err error) *connect.Error {
    return connect.NewError(connect.CodeInvalidArgument, err)
}

func NewNotFound(message string) *connect.Error {
    return NewNotFoundErr(errors.New(message))
}
func NewNotFoundErr(err error) *connect.Error {
    return connect.NewError(connect.CodeNotFound, err)
}

func NewAlreadyExists(message string) *connect.Error {
    return NewAlreadyExistsErr(errors.New(message))
}
func NewAlreadyExistsErr(err error) *connect.Error {
    return connect.NewError(connect.CodeAlreadyExists, err)
}

func NewPermissionDenied(message string) *connect.Error {
    return NewPermissionDeniedErr(errors.New(message))
}
func NewPermissionDeniedErr(err error) *connect.Error {
    return connect.NewError(connect.CodePermissionDenied, err)
}

func NewUnauthenticated(message string) *connect.Error {
    return NewUnauthenticatedErr(errors.New(message))
}
func NewUnauthenticatedErr(err error) *connect.Error {
    return connect.NewError(connect.CodeUnauthenticated, err)
}

func NewInternal(message string) *connect.Error {
    return NewInternalErr(errors.New(message))
}
func NewInternalErr(err error) *connect.Error {
    return connect.NewError(connect.CodeInternal, err)
}

func NewUnimplemented(message string) *connect.Error {
    return NewUnimplementedErr(errors.New(message))
}
func NewUnimplementedErr(err error) *connect.Error {
    return connect.NewError(connect.CodeUnimplemented, err)
}

func NewUnavailable(message string) *connect.Error {
    return NewUnavailableErr(errors.New(message))
}
func NewUnavailableErr(err error) *connect.Error {
    return connect.NewError(connect.CodeUnavailable, err)
}

func NewDeadlineExceeded(message string) *connect.Error {
    return NewDeadlineExceededErr(errors.New(message))
}
func NewDeadlineExceededErr(err error) *connect.Error {
    return connect.NewError(connect.CodeDeadlineExceeded, err)
}

func NewCanceled(message string) *connect.Error {
    return NewCanceledErr(errors.New(message))
}
func NewCanceledErr(err error) *connect.Error {
    return connect.NewError(connect.CodeCanceled, err)
}

func NewUnknown(message string) *connect.Error {
    return NewUnknownErr(errors.New(message))
}
func NewUnknownErr(err error) *connect.Error {
    return connect.NewError(connect.CodeUnknown, err)
}

func NewFailedPrecondition(message string) *connect.Error {
    return NewFailedPreconditionErr(errors.New(message))
}
func NewFailedPreconditionErr(err error) *connect.Error {
    return connect.NewError(connect.CodeFailedPrecondition, err)
}

func NewAborted(message string) *connect.Error {
    return NewAbortedErr(errors.New(message))
}
func NewAbortedErr(err error) *connect.Error {
    return connect.NewError(connect.CodeAborted, err)
}

func NewOutOfRange(message string) *connect.Error {
    return NewOutOfRangeErr(errors.New(message))
}
func NewOutOfRangeErr(err error) *connect.Error {
    return connect.NewError(connect.CodeOutOfRange, err)
}

func NewDataLoss(message string) *connect.Error {
    return NewDataLossErr(errors.New(message))
}
func NewDataLossErr(err error) *connect.Error {
    return connect.NewError(connect.CodeDataLoss, err)
}

func NewResourceExhausted(message string) *connect.Error {
    return NewResourceExhaustedErr(errors.New(message))
}
func NewResourceExhaustedErr(err error) *connect.Error {
    return connect.NewError(connect.CodeResourceExhausted, err)
}
