package care

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapcare struct {
	source  error
	message string
	fields  []zap.Field
}

func New(message string, fields ...zap.Field) error {
	return zapcare{
		source: errors.New(message),
		fields: fields,
	}
}

func Of(err error, message string, fields ...zap.Field) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(zapcare); ok {
		if len(fields) == 0 {
			err.source = errors.Wrap(err.source, message)
			return err
		}

		if len(err.fields) == 0 {
			return withFields(errors.Wrap(err.source, message), fields...)
		}

		return zapcare{
			source:  err,
			fields:  fields,
			message: message,
		}
	}

	return withFields(errors.Wrap(err, message), fields...)
}

func (e zapcare) Error() string {
	if e.message != "" {
		return e.message + ": " + e.source.Error()
	}

	return e.source.Error()
}

func (e zapcare) Unwrap() error {
	return e.source
}

func withFields(err error, fields ...zap.Field) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(zapcare); ok {
		err.fields = append(err.fields, fields...)
		return err
	}

	return zapcare{
		source: err,
		fields: fields,
	}
}

func ToZap(err error) zap.Field {
	return toNamedZap("error", err)
}

func toNamedZap(key string, err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}

	if e, ok := err.(zapcore.ObjectMarshaler); ok {
		return zap.Object(key, e)
	}

	return zap.NamedError(key, err)
}

func (e zapcare) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range e.fields {
		f.AddTo(enc)
	}

	if e.message != "" {
		zap.String("message", e.message).AddTo(enc)
	}

	ToZap(e.source).AddTo(enc)
	return nil
}

func (e zapcare) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.source)
			if e.message != "" {
				fmt.Fprint(s, e.message)
			}
			return
		}
		fallthrough
	case 's', 'q':
		fmt.Fprint(s, e.Error())
	}
}

type Wrapper struct {
	fields []zap.Field
}

func (w Wrapper) New(message string, fields ...zap.Field) error {
	return New(message, append(fields, w.fields...)...)
}

func (w Wrapper) Of(err error, message string, fields ...zap.Field) error {
	return Of(err, message, append(fields, w.fields...)...)
}

func With(fields ...zap.Field) Wrapper {
	return Wrapper{fields: fields}
}

func (w Wrapper) With(fields ...zap.Field) Wrapper {
	return Wrapper{fields: append(w.fields, fields...)}
}
