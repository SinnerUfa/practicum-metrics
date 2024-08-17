package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SinnerUfa/practicum-metric/internal/codes"
	"github.com/SinnerUfa/practicum-metric/internal/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const MAX_CONNS int = 10

const (
	CREATE_TABLE_COUNTERS string = `CREATE TABLE IF NOT EXISTS counters (
                                        cnt_name character varying NOT NULL,
                                        cnt_value bigint NOT NULL,
                                        CONSTRAINT cnt_name_prim PRIMARY KEY (cnt_name)
                                    );`

	INSERT_INTO_COUNTERS string = ` INSERT INTO counters ( cnt_name, cnt_value ) VALUES ( @name, @value )
                                    ON CONFLICT ON CONSTRAINT cnt_name_prim DO
                                    UPDATE SET cnt_value = counters.cnt_value + EXCLUDED.cnt_value;`

	SELECT_ALL_COUNTERS  string = `SELECT counters.cnt_name, counters.cnt_value FROM counters;`
	SELECT_NAME_COUNTERS string = `SELECT counters.cnt_value FROM counters WHERE counters.cnt_name = @name LIMIT 1;`

	CREATE_TABLE_GAUGES string = `CREATE TABLE IF NOT EXISTS gauges (
                                    gau_name character varying NOT NULL,
                                    gau_value double precision NOT NULL,
                                    CONSTRAINT gau_name_prim PRIMARY KEY (gau_name)
                                );`

	INSERT_INTO_GAUGES string = `INSERT INTO gauges (gau_name, gau_value) VALUES (@name, @value)
                                ON CONFLICT ON CONSTRAINT gau_name_prim DO
                                UPDATE SET gau_value = EXCLUDED.gau_value;`

	SELECT_ALL_GAUGES  string = `SELECT gauges.gau_name, gauges.gau_value FROM gauges;`
	SELECT_NAME_GAUGES string = `SELECT gauges.gau_value FROM gauges WHERE gauges.gau_name = @name LIMIT 1;`
)

var dbQueries map[string]string = map[string]string{
	"InsertCounters":     INSERT_INTO_COUNTERS,
	"SelectCounters":     SELECT_ALL_COUNTERS,
	"SelectNameCounters": SELECT_NAME_COUNTERS,
	"InsertGauges":       INSERT_INTO_GAUGES,
	"SelectGauges":       SELECT_ALL_GAUGES,
	"SelectNameGauges":   SELECT_NAME_GAUGES,
}

type Database struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func New(ctx context.Context, dsn string) (*Database, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.Join(err, db.Close())
	}
	db.SetMaxOpenConns(MAX_CONNS)
	db.SetMaxIdleConns(MAX_CONNS)
	db.SetConnMaxIdleTime(4 * time.Minute)
	db.SetConnMaxLifetime(15 * time.Minute)

	if _, err := db.Exec(CREATE_TABLE_COUNTERS); err != nil {
		return nil, errors.Join(err, db.Close())
	}
	if _, err := db.Exec(CREATE_TABLE_GAUGES); err != nil {
		return nil, errors.Join(err, db.Close())
	}
	DB := &Database{db: db}
	for name, query := range dbQueries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, errors.Join(err, DB.Close())
		}
		DB.stmts[name] = stmt
	}

	return DB, nil
}

func (d *Database) SetContext(ctx context.Context, m metrics.Metric) error {
	if m.Name == "" {
		return codes.ErrRepMetricNotSupported
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
		v, ok := m.Value.Int64()
		if !ok {
			return codes.ErrRepParseInt
		}
		tx, err := d.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		txStmt := tx.Stmt(d.stmts["InsertCounters"])
		_, err = txStmt.ExecContext(ctx, sql.Named("name", m.Name), sql.Named("value", v))
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
		return tx.Commit()
	case metrics.MetricTypeGauge:
		v, ok := m.Value.Float64()
		if !ok {
			return codes.ErrRepParseFloat
		}
		tx, err := d.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		txStmt := tx.Stmt(d.stmts["InsertGauges"])
		_, err = txStmt.ExecContext(ctx, sql.Named("name", m.Name), sql.Named("value", v))
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
		return tx.Commit()
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (d *Database) GetContext(ctx context.Context, m *metrics.Metric) error {
	if m.Name == "ping" && m.Type == "" && m.Value.IsString() && m.Value.String() == "ping" {
		return d.db.PingContext(ctx)
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
		var value int64
		r := d.stmts["SelectNameCounters"].QueryRowContext(ctx, sql.Named("name", m.Name))
		err := r.Scan(&value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return codes.ErrRepNotFound
			}
			return err
		}
		m.Value = metrics.Int(value)
	case metrics.MetricTypeGauge:
		var value float64
		r := d.stmts["SelectNameGauges"].QueryRowContext(ctx, sql.Named("name", m.Name))
		err := r.Scan(&value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return codes.ErrRepNotFound
			}
			return err
		}
		m.Value = metrics.Float(value)
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (d *Database) SetListContext(ctx context.Context, in []metrics.Metric) error {
	var err error

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txStmtCnt := tx.Stmt(d.stmts["InsertCounters"])
	txStmtGau := tx.Stmt(d.stmts["InsertGauges"])

	for _, m := range in {
		if m.Name == "" {
			err = codes.ErrRepMetricNotSupported
			break
		}
		switch m.Type {
		case metrics.MetricTypeCounter:
			v, ok := m.Value.Int64()
			if !ok {
				err = codes.ErrRepParseInt
				break
			}
			_, err = txStmtCnt.ExecContext(ctx, sql.Named("name", m.Name), sql.Named("value", v))
		case metrics.MetricTypeGauge:
			v, ok := m.Value.Float64()
			if !ok {
				err = codes.ErrRepParseFloat
				break
			}
			_, err = txStmtGau.ExecContext(ctx, sql.Named("name", m.Name), sql.Named("value", v))
		default:
			err = codes.ErrRepMetricNotSupported
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (d *Database) GetListContext(ctx context.Context) (out []metrics.Metric, err error) {
	rCnts, err := d.stmts["SelectCounters"].QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	rGaus, err := d.stmts["SelectGauges"].QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var (
		name             string
		errCnts, errGaus error
	)
	for rCnts.Next() {
		var value int64

		errCnts = rCnts.Scan(&name, &value)
		if errCnts != nil {
			break
		}

		out = append(out, metrics.Metric{Name: name, Value: metrics.Int(value), Type: metrics.MetricTypeCounter})
	}
	for rGaus.Next() {
		var value float64

		errGaus = rGaus.Scan(&name, &value)
		if errGaus != nil {
			break
		}

		out = append(out, metrics.Metric{Name: name, Value: metrics.Float(value), Type: metrics.MetricTypeGauge})
	}

	errCnts = errors.Join(errCnts, rCnts.Err())
	if !errors.Is(errCnts, sql.ErrNoRows) {
		err = errors.Join(err, errCnts)
	}

	errGaus = errors.Join(errGaus, rGaus.Err())
	if !errors.Is(errGaus, sql.ErrNoRows) {
		err = errors.Join(err, errGaus)
	}

	if err != nil {
		return nil, errors.Join(err, rCnts.Close(), rGaus.Close())
	}
	rCnts.Close()
	rGaus.Close()
	return out, nil
}

func (d *Database) Close() error {
	var err error
	for _, m := range d.stmts {
		err = errors.Join(err, m.Close())
	}
	return errors.Join(err, d.db.Close())
}
