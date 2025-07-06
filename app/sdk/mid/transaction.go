package mid

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"service/app/sdk/errs"
	"service/business/sdk/sqldb"
	"service/foundation/logger"
	"service/foundation/web"
)

func BeginCommitRollback(log *logger.Logger, bgn sqldb.Beginner) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			hasCommitted := false

			log.Info(ctx, "BEGIN TRANSACTION")
			tx, err := bgn.Begin()

			if err != nil {
				return errs.Newf(errs.Internal, "BEGIN TRANSACTION: %s", err)
			}

			defer func() {
				if !hasCommitted {
					log.Info(ctx, "ROLLBACK TRANSACTION")
				}
				if err := tx.Rollback(); err != nil {
					if errors.Is(err, sql.ErrTxDone) {
						return
					}
					log.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
				}
			}()

			ctx = setTran(ctx, tx)
			resp := next(ctx, r)

			// if errs.isError(resp) != nil {
			// 	return resp
			// }

			log.Info(ctx, "COMMIT TRANSACTION")
			if err := tx.Commit(); err != nil {
				return errs.Newf(errs.Internal, "COMMIT TRANSACTION: %s", err)
			}

			hasCommitted = true

			return resp
		}
		return h
	}
	return m
}
