package metal

import (
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func dberr(err error) error {
	if err == nil {
		return nil
	}

	var sqliteErr *sqlite.Error

	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			return status.Error(codes.AlreadyExists, sqliteErr.Error())
		case sqlite3.SQLITE_NOTFOUND:
			return status.Error(codes.NotFound, sqliteErr.Error())
		case sqlite3.SQLITE_READONLY:
			return status.Error(codes.PermissionDenied, sqliteErr.Error())
		case sqlite3.SQLITE_BUSY:
			return status.Error(codes.Unavailable, sqliteErr.Error())
		case sqlite3.SQLITE_IOERR:
			return status.Error(codes.DataLoss, sqliteErr.Error())
		default:
			return status.Error(codes.Internal, sqliteErr.Error())
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}

var (
	errMissingAppliance     = status.Error(codes.InvalidArgument, "appliance not specified")
	errMissingAttribute     = status.Error(codes.InvalidArgument, "attribute not specified")
	errMissingCluster       = status.Error(codes.InvalidArgument, "cluster not specified")
	errMissingEnvironment   = status.Error(codes.InvalidArgument, "environment not specified")
	errMissingHost          = status.Error(codes.InvalidArgument, "host not specified")
	errMissingInterface     = status.Error(codes.InvalidArgument, "inteface not specified")
	errMissingMake          = status.Error(codes.InvalidArgument, "make not specified")
	errMissingModel         = status.Error(codes.InvalidArgument, "model not specified")
	errMissingNetwork       = status.Error(codes.InvalidArgument, "network not specified")
	errMissingRack          = status.Error(codes.InvalidArgument, "rack not specified")
	errMissingZone          = status.Error(codes.InvalidArgument, "zone not specified")
	errNameOrGlob           = status.Error(codes.InvalidArgument, "must provide either name or glob")
	errInvalidHostType      = status.Error(codes.InvalidArgument, "invalid host type")
	errInvalidInterfaceType = status.Error(codes.InvalidArgument, "invalid interface type")
	errInvalidBondMode      = status.Error(codes.InvalidArgument, "invalid bond mode")
	errInvalidArchitecture  = status.Error(codes.InvalidArgument, "invalid architecture")
)
