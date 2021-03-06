package lambda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Int(t *testing.T) {
	var testData = []struct {
		expr   string
		dst    []int
		expect []int
	}{
		{
			expr:   `(p) =>  p % 2 == 1`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{1, 3, 5, 7},
		},
		{
			expr:   `(v) =>  (v&1) == 1`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{1, 3, 5, 7},
		},
		{
			expr:   `(v) =>  v&1 == 0`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{2, 4, 6},
		},
		{
			expr:   `(v) =>  v+v<= 6 `,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{1, 2, 3},
		},
		{
			expr:   `(v) =>  (v<<2)>>1 == 8`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{4},
		},
		{
			expr:   `(v) =>  v%2 == 0`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{2, 4, 6},
		},
		{
			expr:   `(v) =>  v > 5`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{6, 7},
		},
		{
			expr:   `(v) =>  v > 5+1`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{7},
		},
		{
			expr:   `(v) =>  v*2 == v+3`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{3},
		},
		{
			expr:   `(v) =>  v > 1+2+3/(0+1)`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{7},
		},
		{
			expr:   `(v) =>  v <= 1+2+3/(0+1)`,
			dst:    []int{1, 2, 3, 4, 5, 6, 7},
			expect: []int{1, 2, 3, 4, 5, 6},
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		mstate := Filter(tc.expr)
		if !ass.NoError(mstate.err, "%s", tc.expr) {
			t.FailNow()
		}
		result := mstate.Call(tc.dst)
		if !ass.NoError(result.err, "%s", tc.expr) {
			t.FailNow()
		}
		inf, err := result.Interface()
		if !ass.NoError(err, "%s", tc.expr) {
			t.FailNow()
		}
		ans, ok := inf.([]int)
		if !ass.True(ok) {
			t.FailNow()
		}
		ass.Equal(tc.expect, ans, "expr=%s", tc.expr)
	}
}

func TestFilter_Float64(t *testing.T) {
	var testData = []struct {
		expr   string
		dst    []float64
		expect []float64
	}{
		{
			expr:   `(v) =>  v*2 == v+3`,
			dst:    []float64{1, 2, 3, 4, 5, 6, 7},
			expect: []float64{3},
		},
		{
			expr:   `(v) =>  v*2 >= 5`,
			dst:    []float64{1, 2, 3, 4, 5, 6, 7},
			expect: []float64{3, 4, 5, 6, 7},
		},
		{
			expr:   `(v) =>  v*2 >= 5+1+1`,
			dst:    []float64{1, 2, 3, 4, 5, 6, 7},
			expect: []float64{4, 5, 6, 7},
		},
		{
			expr:   `(v) =>  v*2 >= 5+1+1*((2+1*10000)*0+1)`,
			dst:    []float64{1, 2, 3, 4, 5, 6, 7},
			expect: []float64{4, 5, 6, 7},
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		mstate := Filter(tc.expr)
		if !ass.NoError(mstate.err, "%s", tc.expr) {
			t.FailNow()
		}
		result := mstate.Call(tc.dst)
		if !ass.NoError(result.err, "%s", tc.expr) {
			t.FailNow()
		}
		inf, err := result.Interface()
		if !ass.NoError(err, "%s", tc.expr) {
			t.FailNow()
		}
		ans, ok := inf.([]float64)
		if !ass.True(ok) {
			t.FailNow()
		}
		ass.Equal(tc.expect, ans, "expr=%s", tc.expr)
	}
}

type Student struct {
	Age  int
	Name string
}

func TestFilter_Struct(t *testing.T) {
	var students = []Student{
		Student{
			Name: "deen",
			Age:  24,
		},
		Student{
			Name: "bob",
			Age:  22,
		},
		Student{
			Name: "alice",
			Age:  23,
		},
		Student{
			Name: "tom",
			Age:  25,
		},
		Student{
			Name: "jerry",
			Age:  20,
		},
	}
	var testData = []struct {
		expr   string
		expect []int
	}{
		{
			expr:   `(v) => v.Age+2 > 24+1`,
			expect: []int{0, 3},
		},
		{
			expr:   `(v) => v.Age >= 23`,
			expect: []int{0, 2, 3},
		},
		{
			expr:   `(v) => v.Age < 23`,
			expect: []int{1, 4},
		},
		{
			expr:   `(v) => v.Age < 23 || v.Name == "tom"`,
			expect: []int{1, 3, 4},
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		mstate := Filter(tc.expr)
		if !ass.NoError(mstate.err, "%s", tc.expr) {
			t.FailNow()
		}
		result := mstate.Call(students)
		if !ass.NotNil(result, "%s", tc.expr) {
			t.FailNow()
		}
		if !ass.NoError(result.err, "%s", tc.expr) {
			t.FailNow()
		}
		inf, err := result.Interface()
		if !ass.NoError(err, "%s", tc.expr) {
			t.FailNow()
		}
		ans, ok := inf.([]Student)
		if !ass.True(ok) {
			t.FailNow()
		}
		var expectArr []Student
		for _, idx := range tc.expect {
			expectArr = append(expectArr, students[idx])
		}
		ass.Equal(expectArr, ans, "expr=%s", tc.expr)
	}
}

func TestFilter_Pointer(t *testing.T) {
	var students = []*Student{
		&Student{
			Name: "deen",
			Age:  24,
		},
		&Student{
			Name: "bob",
			Age:  22,
		},
		&Student{
			Name: "alice",
			Age:  23,
		},
		&Student{
			Name: "tom",
			Age:  25,
		},
		&Student{
			Name: "jerry",
			Age:  20,
		},
	}
	var testData = []struct {
		expr   string
		expect []int
	}{
		{
			expr:   `(v) => v.Age+2 > 24+1`,
			expect: []int{0, 3},
		},
		{
			expr:   `(v) => v.Age >= 23`,
			expect: []int{0, 2, 3},
		},
		{
			expr:   `(v) => v.Age < 23`,
			expect: []int{1, 4},
		},
	}
	ass := assert.New(t)
	for _, tc := range testData {
		mstate := Filter(tc.expr)
		if !ass.NoError(mstate.err, "%s", tc.expr) {
			t.FailNow()
		}
		result := mstate.Call(students)
		if !ass.NotNil(result, "%s", tc.expr) {
			t.FailNow()
		}
		if !ass.NoError(result.err, "%s", tc.expr) {
			t.FailNow()
		}
		inf, err := result.Interface()
		if !ass.NoError(err, "%s", tc.expr) {
			t.FailNow()
		}
		ans, ok := inf.([]*Student)
		if !ass.True(ok) {
			t.FailNow()
		}
		var expectArr []*Student
		for _, idx := range tc.expect {
			expectArr = append(expectArr, students[idx])
		}
		ass.Equal(expectArr, ans, "expr=%s", tc.expr)
	}
}
