package grpcsteps

// ErrInvalidGRPCMethod indicates that the service method is not in SERVICE/METHOD format.
const ErrInvalidGRPCMethod err = "invalid grpc method"

type err string

// Error returns the error string.
func (e err) Error() string {
	return string(e)
}
