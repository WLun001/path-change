package main

// copied from https://github.com/tektoncd/triggers/blob/main/pkg/apis/triggers/v1beta1/interceptor_types.go

type InterceptorRequest struct {
	// Body is the incoming HTTP event body. We use a "string" representation of the JSON body
	// in order to preserve the body exactly as it was sent (including spaces etc.). This is necessary
	// for some interceptors e.g. GitHub for validating the body with a signature. While []byte can also
	// store an exact representation of the body, `json.Marshal` will compact []byte to a base64 encoded
	// string which means that we will lose the spaces any time we marshal this struct.
	Body string `json:"body,omitempty"`

	// Header are the headers for the incoming HTTP event
	Header map[string][]string `json:"header,omitempty"`

	// Extensions are extra values that are added by previous interceptors in a chain
	Extensions map[string]interface{} `json:"extensions,omitempty"`

	// InterceptorParams are the user specified params for interceptor in the Trigger
	InterceptorParams map[string]interface{} `json:"interceptor_params,omitempty"`

	// Context contains additional metadata about the event being processed
	Context *TriggerContext `json:"context"`
}

type TriggerContext struct {
	// EventURL is the URL of the incoming event
	EventURL string `json:"event_url,omitempty"`
	// EventID is a unique ID assigned by Triggers to each event
	EventID string `json:"event_id,omitempty"`
	// TriggerID is of the form namespace/$ns/triggers/$name
	TriggerID string `json:"trigger_id,omitempty"`
}

type InterceptorResponse struct {
	// Extensions are additional fields that is added to the interceptor event.
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	// Continue indicates if the EventListener should continue processing the Trigger or not
	Continue bool `json:"continue"` // Don't add omitempty -- it  will remove the continue field when the value is false.
	// Status is an Error status containing details on any interceptor processing errors
	Status Status `json:"status"`
}

type Status struct {
	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code Code `json:"code,omitempty"`
	// A developer-facing error message, which should be in English.
	Message string `json:"message,omitempty"`
}

// copied from https://pkg.go.dev/google.golang.org/grpc/codes
type Code uint32

const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	//
	// The gRPC framework will generate this error code when cancellation
	// is requested.
	Canceled Code = 1

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	//
	// The gRPC framework will generate this error code in the above two
	// mentioned cases.
	Unknown Code = 2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	//
	// This error code will not be generated by the gRPC framework.
	InvalidArgument Code = 3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	//
	// The gRPC framework will generate this error code when the deadline is
	// exceeded.
	DeadlineExceeded Code = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	//
	// This error code will not be generated by the gRPC framework.
	NotFound Code = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	//
	// This error code will not be generated by the gRPC framework.
	AlreadyExists Code = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	//
	// This error code will not be generated by the gRPC core framework,
	// but expect authentication middleware to use it.
	PermissionDenied Code = 7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	//
	// This error code will be generated by the gRPC framework in
	// out-of-memory and server overload situations, or when a message is
	// larger than the configured maximum size.
	ResourceExhausted Code = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	//
	// This error code will not be generated by the gRPC framework.
	FailedPrecondition Code = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	//
	// This error code will not be generated by the gRPC framework.
	Aborted Code = 10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	//
	// This error code will not be generated by the gRPC framework.
	OutOfRange Code = 11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	//
	// This error code will be generated by the gRPC framework. Most
	// commonly, you will see this error code when a method implementation
	// is missing on the server. It can also be generated for unknown
	// compression algorithms or a disagreement as to whether an RPC should
	// be streaming.
	Unimplemented Code = 12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	//
	// This error code will be generated by the gRPC framework in several
	// internal error conditions.
	Internal Code = 13

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	//
	// This error code will be generated by the gRPC framework during
	// abrupt shutdown of a server process or network connection.
	Unavailable Code = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	//
	// This error code will not be generated by the gRPC framework.
	DataLoss Code = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	//
	// The gRPC framework will generate this error code when the
	// authentication metadata is invalid or a Credentials callback fails,
	// but also expect authentication middleware to generate it.
	Unauthenticated Code = 16

	_maxCode = 17
)
