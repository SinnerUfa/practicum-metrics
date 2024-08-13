package config

import (
	"testing"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	"github.com/stretchr/testify/assert"
)

func Test_LoadFlags(t *testing.T) {

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
				codes.ErrFlgNoAcsess,
			},
			{
				b,
				codes.ErrFlgNoAcsess,
			},
			{
				c,
				codes.ErrFlgNoAcsess,
			},
			{
				d,
				codes.ErrFlgNoAcsess,
			},
			{
				&a,
				codes.ErrFlgNotStructure,
			},
			{
				&b,
				codes.ErrFlgNotStructure,
			},
			{
				&c,
				codes.ErrFlgNotStructure,
			},
		}

		for _, test := range tests {
			err := LoadFlags(test.in, []string{})
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
			in int `flag:"in"`
		}
		type tmp1 struct {
			in float32 `flag:"in"`
		}
		type tmp2 struct {
			In float32 `flag:"in"`
		}
		a := tmp0{20}
		b := tmp1{20}
		c := tmp2{20}

		tests := []test{
			{
				&a,
				codes.ErrFlgFieldNotSet,
			},
			{
				&b,
				codes.ErrFlgFieldNotSupported,
			},
			{
				&c,
				codes.ErrFlgFieldNotSupported,
			},
		}
		for _, test := range tests {
			err := LoadFlags(test.in, []string{})
			if assert.Error(t, err) {
				assert.Equal(t, test.want, err)
			}
		}
	})

	t.Run("TypeLoadTest", func(t *testing.T) {
		type testValue struct {
			A uint   `flag:"A"`
			B int    `flag:"B"`
			G bool   `flag:"G"`
			C string `flag:"C"`
			D uint   `flag:""`
			E uint   `flag:"-"`
			F uint
			J bool `flag:"J"`
			K bool `flag:"K"`
		}

		type test struct {
			args []string
			want testValue
		}

		tests := []test{
			{
				[]string{
					"-A", "11", "-B", "21", "-C", "31", "-G", "false", "-J", "true", "-K", "",
				},
				testValue{
					A: 11, B: 21, C: "31", D: 4, E: 5, F: 6, G: false, J: true, K: false,
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
				G: true,
				J: false,
				K: true,
			}

			err := LoadFlags(&in, test.args)
			if assert.NoError(t, err) {
				assert.Equal(t, in, test.want)
			}
		}

	})
}