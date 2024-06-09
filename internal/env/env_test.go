package env

import (
	"os"
	"testing"

	codes "github.com/SinnerUfa/practicum-metric/internal/err_codes"
	"github.com/stretchr/testify/assert"
)

func Test_Load(t *testing.T) {

	t.Run("TypeInTest", func(t *testing.T) {
		type test struct {
			in   any
			want error
		}
		a := int(1)
		b := float32(1)
		c := []int{0, 1, 2}
		type tmp struct {
			in int
		}
		d := tmp{20}

		tests := []test{
			{
				a,
				codes.ErrEnvNoAcsess,
			},
			{
				b,
				codes.ErrEnvNoAcsess,
			},
			{
				c,
				codes.ErrEnvNoAcsess,
			},
			{
				d,
				codes.ErrEnvNoAcsess,
			},
			{
				&a,
				codes.ErrEnvNotStructure,
			},
			{
				&b,
				codes.ErrEnvNotStructure,
			},
			{
				&c,
				codes.ErrEnvNotStructure,
			},
		}

		for _, test := range tests {
			err := Load(test.in)
			if assert.Error(t, err) {
				assert.Equal(t, test.want, err)
			}
		}
	})
	t.Run("TypeFieldTest", func(t *testing.T) {
		type test struct {
			in   any
			want error
		}
		type tmp0 struct {
			in int `env:"IN"`
		}
		type tmp1 struct {
			in float32 `env:"IN"`
		}
		type tmp2 struct {
			In float32 `env:"IN"`
		}
		a := tmp0{20}
		b := tmp1{20}
		c := tmp2{20}
		os.Setenv("IN", "0")
		tests := []test{
			{
				&a,
				codes.ErrEnvFieldNotSet,
			},
			{
				&b,
				codes.ErrEnvFieldNotSet,
			},
			{
				&c,
				codes.ErrFlgFieldNotSupported,
			},
		}
		for _, test := range tests {
			err := Load(test.in)
			if assert.Error(t, err) {
				assert.Equal(t, test.want, err)
			}
		}
	})

	t.Run("TypeLoadTest", func(t *testing.T) {
		type testArg struct {
			name  string
			value string
		}
		type testValue struct {
			A uint   `env:"A"`
			B int    `env:"B"`
			C string `env:"C"`
			D uint   `env:""`
			E uint   `env:"-"`
			F uint
		}

		type test struct {
			args []testArg
			want testValue
		}

		tests := []test{
			{
				[]testArg{
					{"A", "11"}, {"B", "21"}, {"C", "31"}, {"D", "41"}, {"E", "51"}, {"F", "61"},
				},
				testValue{
					11, 21, "31", 4, 5, 6,
				},
			},
		}
		for _, test := range tests {
			in := testValue{
				A: 1,
				B: -2,
				C: "3",
				D: 4,
				E: 5,
				F: 6,
			}

			for _, a := range test.args {
				os.Setenv(a.name, a.value)
			}
			err := Load(&in)
			if assert.NoError(t, err) {
				assert.Equal(t, in, test.want)
			}
		}

	})
}
