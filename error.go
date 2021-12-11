package grpcsteps

const (
	// ErrInvalidGRPCMethod indicates that the service method is not in SERVICE/METHOD format.
	ErrInvalidGRPCMethod err = `invalid grpc method`
	// ErrGRPCServiceNotFound indicates that the service is not found.
	ErrGRPCServiceNotFound err = `grpc service not found`
	// ErrGRPCMethodNotFound indicates that the service method is not found.
	ErrGRPCMethodNotFound err = `grpc method not found`
	// ErrGRPCMethodNotSupported indicates that the service method is not supported.
	ErrGRPCMethodNotSupported err = `grpc method not supported`
)

type err string

// Error returns the error string.
func (e err) Error() string {
	return string(e)
}
