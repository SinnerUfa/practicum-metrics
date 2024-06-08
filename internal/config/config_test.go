package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Load(t *testing.T) {

	t.Run("TypeLoadTest", func(t *testing.T) {
		type testArg struct {
			name  string
			value string
		}
		type testValue struct {
			A uint   `env:"A" flag:"A"`
			B int    `env:"B" flag:"-"`
			C string `env:"-" flag:"C"`
			D uint   `env:""flag:""`
			E uint   `env:"-" flag:"-"`
			F uint
		}

		type test struct {
			args  []testArg
			flags []string
			want  testValue
		}

		tests := []test{
			{
				[]testArg{
					{"A", "11"}, {"B", "21"}, {"C", "311"}, {"D", "41"}, {"E", "51"}, {"F", "61"},
				},
				[]string{
					"-A", "11", "-C", "31",
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
			err := Load(&in, test.flags)
			if assert.NoError(t, err) {
				assert.Equal(t, in, test.want)
			}
		}

	})
}
