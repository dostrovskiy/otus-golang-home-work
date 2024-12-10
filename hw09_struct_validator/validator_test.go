package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Cat struct {
		Age   int    `validate:"min:2|max:15"`
		Color string `validate:"in:red,green,blue"`
		Name  string `validate:"len:6|regexp:^\\p{L}+$"`
	}

	Movie struct {
		ID    string `validate:"len:36"`
		Title string
		Stars []User `validate:"nested"`
		Staff []User
	}

	WrongRegExp struct {
		Name string `validate:"len:4|regexp:^\\p{L+$"`
	}

	WrongRegExpSlice struct {
		Names []string `validate:"len:4|regexp:^\\g+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Cat{Age: 7, Color: "red", Name: "Barsik"},
			nil,
		},
		{
			Response{Code: 404},
			nil,
		},
		{
			App{Version: "12345"},
			nil,
		},
		{
			Token{Header: []byte("header"), Payload: []byte("payload"), Signature: []byte("signature")},
			nil,
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-426655440000",
				Name:   "Иван Иванович",
				Age:    50,
				Email:  "2K3iM@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			nil,
		},
		{
			Movie{
				ID:    "123e4567-e89b-12d3-a456-426655440000",
				Title: "Pulp Fiction",
				Stars: []User{
					{
						ID:     "123e4567-e89b-12d3-a456-426655440000",
						Name:   "Bruce Willis",
						Age:    39,
						Email:  "willis@google.com",
						Role:   "stuff",
						Phones: []string{"12345678901", "10987654321"},
					},
					{
						ID:     "123e4567-e89b-12d3-a456-426655440000",
						Name:   "John Travolta",
						Age:    40,
						Email:  "travolta@google.com",
						Role:   "stuff",
						Phones: []string{"12345678901", "10987654321"},
					},
				},
				Staff: []User{
					{
						ID:     "123e4567-e89b-12d3-a456-426655440000",
						Name:   "John Doe",
						Age:    50,
						Email:  "doe@mail.com",
						Role:   "stuff",
						Phones: []string{"12345678901", "10987654321"},
					},
				},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			Cat{Age: 1, Color: "red", Name: "BaRs1k3"},
			ValidationErrors{
				{
					Field: "Age",
					Err:   MinIntValidationError(2),
				},
				{
					Field: "Name",
					Err:   LenStrValidationError(6),
				},
				{
					Field: "Name",
					Err:   RegExpStrValidationError("^\\p{L}+$"),
				},
			},
		},
		{
			Response{Code: 1},
			ValidationErrors{
				{
					Field: "Code",
					Err:   InIntValidationError([]int64{200, 404, 500}),
				},
			},
		},
		{
			App{Version: "1"},
			ValidationErrors{
				{
					Field: "Version",
					Err:   LenStrValidationError(5),
				},
			},
		},
		{
			User{
				ID:     "123e4567-e89b-12d3-a456-42665544000",
				Name:   "Иван Иванович",
				Age:    50,
				Email:  "2K3iM@example.com.",
				Role:   "user",
				Phones: []string{"123456789012", "210987654321"},
			},
			ValidationErrors{
				{
					Field: "ID",
					Err:   LenStrValidationError(36),
				},
				{
					Field: "Email",
					Err:   RegExpStrValidationError("^\\w+@\\w+\\.\\w+$"),
				},
				{
					Field: "Role",
					Err:   InStrValidationError([]string{"admin", "stuff"}),
				},
				{
					Field: "Phones",
					Err:   LenStrValidationError(11),
				},
				{
					Field: "Phones",
					Err:   LenStrValidationError(11),
				},
			},
		},
		{
			Movie{
				ID:    "123e4567-e89b-12d3-a456-426655440000",
				Title: "Pulp Fiction",
				Stars: []User{
					{
						ID:     "123e4567-e89b-12d3-a456-426655440000",
						Name:   "Bruce Willis",
						Age:    39,
						Email:  "willis@google.com.",
						Role:   "starr",
						Phones: []string{"12345678901", "10987654321"},
					},
					{
						ID:     "123e4567-e89b-12d3-a456-426655440000",
						Name:   "John Travolta",
						Age:    40,
						Email:  "travolta@google.com",
						Role:   "stuff",
						Phones: []string{"12345678901", "109876543210"},
					},
				},
				Staff: []User{
					{
						ID:     "123e4567-e89b-12d3-a456-42665544",
						Name:   "John Doe",
						Age:    50,
						Email:  "doe@mail.com",
						Role:   "stuff",
						Phones: []string{"12345678901", "109876543212"},
					},
				},
			},
			ValidationErrors{
				{
					Field: "Email",
					Err:   RegExpStrValidationError("^\\w+@\\w+\\.\\w+$"),
				},
				{
					Field: "Role",
					Err:   InStrValidationError([]string{"admin", "stuff"}),
				},
				{
					Field: "Phones",
					Err:   LenStrValidationError(11),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.Truef(t, errors.As(err, &ValidationErrors{}), "actual error %q", err)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestInvalidErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			WrongRegExp{Name: "rose"},
			ErrRuleSyntax(""),
		},
		{
			WrongRegExpSlice{Names: []string{"rose", "iris", "lily"}},
			ErrRuleSyntax(""),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorContainsf(t, err, "error parsing regexp", "Expected containing 'error parsing regexp', actual: %v", err)
		})
	}
}

type (
	Dumb struct {
		Numbers []string `validate:"regexp:^\\w+$"`
	}
)

func BenchmarkValidateRegexpStringSlice(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			nums := []string{}
			for i := 0; i < size; i++ {
				nums = append(nums, fmt.Sprintf("%d", i))
			}
			dumb := Dumb{nums}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Validate(dumb)
			}
		})
	}
}

type (
	People struct {
		Names []string `validate:"len:4|in:John,Paul,Lily,Rose"`
	}
)

func BenchmarkValidateLenInStringSlice(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			names := []string{}
			for i := 0; i < size; i++ {
				names = append(names, "George")
			}
			dumb := Dumb{names}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Validate(dumb)
			}
		})
	}
}

type (
	Nursery struct {
		Ages []int `validate:"min:3|max:7000|in:3,40,500,600,7000"`
	}
)

func BenchmarkValidateMinMaxInIntSlice(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			ages := []int{}
			for i := 0; i < size; i++ {
				ages = append(ages, rand.Intn(100))
			}
			sunshine := Nursery{ages}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Validate(sunshine)
			}
		})
	}
}

type (
	Company struct {
		Employees []User `validate:"nested"`
	}
)

func BenchmarkValidateStructSlice(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			employees := []User{}
			for i := 0; i < size; i++ {
				employees = append(employees, User{
					ID:     uuid.NewString(),
					Name:   fmt.Sprintf("User %d", i),
					Age:    rand.Intn(100),
					Email:  fmt.Sprintf("user%d@domain.com", i),
					Role:   "stuff",
					Phones: []string{"12345678901", "109876543210"},
				})
			}
			noname := Company{employees}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Validate(noname)
			}
		})
	}
}
