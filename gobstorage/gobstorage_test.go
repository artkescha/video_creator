package gobstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	filename := "configuration"

	save1 := 1
	_ = Save(filename, save1)
	var load1 int
	_ = Load(filename, &load1)
	assert.Equal(t, save1, load1)

	save2 := []int{1, 2, 3}
	_ = Save(filename, save2)
	var load2 []int
	_ = Load(filename, &load2)
	assert.Error(t, Load(filename, load2))
	assert.Equal(t, save2, load2)

	save2 = append(save2, 4)
	_ = Save(filename, save2)
	_ = Load(filename, &load2)
	assert.Error(t, Load(filename, load2))
	assert.Equal(t, save2, load2)

	save3 := T{"11", 1}
	_ = Save(filename, save3)
	var load3 T
	_ = Load(filename, &load3)
	assert.Error(t, Load(filename, load3))
	assert.Equal(t, save3, load3)

	save4 := []T{{"11", 1}, {"йц", 2}}
	_ = Save(filename, save4)
	var load4 []T
	_ = Load(filename, &load4)
	assert.Error(t, Load(filename, load4))
	assert.Equal(t, save4, load4)

	save5 := []*T{{"11", 1}, {"йц", 2}, {"qwerty", 3}}
	_ = Save(filename, save5)
	var load5 []*T
	_ = Load(filename, &load5)
	assert.Error(t, Load(filename, load5))
	assert.Equal(t, save5, load5)

	str1 := "11"
	str2 := "йц"
	save6 := []TT{{&str1, 1}, {&str2, 2}, {nil, 3}}
	_ = Save(filename, save6)
	var load6 []TT
	_ = Load(filename, &load6)
	assert.Error(t, Load(filename, load6))
	assert.Equal(t, save6, load6)
}

type T struct {
	Str string
	D   int
}

type TT struct {
	Str *string `json:"-" bson:"-"`
	D   int     `json:"isActive" bson:"is_active"`
}
