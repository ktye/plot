package plot

import (
	"encoding/json"
	"fmt"
)

// json encode/decode complex vectors [z,z,..]->[r,i,r,i,..]
//
// marshal/unmarshal     Line.C
// marshal CaptionColumn
// unmarshal Caption     caption column data must be converted, complex columns can be guessed if not on the first column (todo minmax)

//json encode/decode complex vectors as float [z,z,..]->[r,i,r,i,..]

func (l *Line) MarshalJSON() ([]byte, error) {
	type L Line
	return json.Marshal(&struct {
		C []float64
		*L
	}{
		C: fz(l.C),
		L: (*L)(l),
	})
}
func (l *Line) UnmarshalJSON(b []byte) error {
	type L Line
	x := &struct {
		*L
		C []float64
	}{
		L: (*L)(l),
	}
	if e := json.Unmarshal(b, &x); e != nil {
		return e
	}
	x.L.C = zf(x.C)
	return nil
}

func (c *CaptionColumn) MarshalJSON() ([]byte, error) {
	type C CaptionColumn
	var data interface{}
	data = c.Data
	if z, o := c.Data.([]complex128); o {
		data = fz(z)
	}
	return json.Marshal(&struct {
		Data interface{}
		*C
	}{
		Data: data,
		C:    (*C)(c),
	})
}
func (c *Caption) UnmarshalJSON(b []byte) error { //if a []float64 column has twice the length as the first, it's complex
	type C Caption
	x := &struct {
		*C
	}{
		C: (*C)(c),
	}
	if e := json.Unmarshal(b, &x); e != nil {
		return e
	}
	// data unmarshals as []interface{}, it should be []string|[]int|[]float|[]complex128
	for i := range c.Columns {
		if c.Columns[i].Data != nil {
			if v, o := c.Columns[i].Data.([]interface{}); o {
				col, e := autocol(v)
				if e != nil {
					return e
				}
				c.Columns[i].Data = col
			}
		}
	}

	if c != nil && len(c.Columns) > 1 {
		n := c.Rows()
		for i := range c.Columns {
			if c.Columns[i].Data != nil {
				if f, o := c.Columns[i].Data.([]float64); o {
					if len(f) == 2*n {
						c.Columns[i].Data = zf(f)
					}
				}
			}
		}
	}
	return nil
}

func autocol(v []interface{}) (interface{}, error) { //convert caption column data after json unmarshal
	if v == nil {
		return nil, nil
	}
	if len(v) == 0 {
		return []string{}, nil
	}
	n := len(v)
	switch t := v[0].(type) {
	case string:
		r := make([]string, n)
		for i := range v {
			r[i] = v[i].(string)
		}
		return r, nil
	case int:
		r := make([]int, n)
		for i := range v {
			r[i] = v[i].(int)
		}
		return r, nil
	case float64:
		r := make([]float64, n)
		for i := range v {
			r[i] = v[i].(float64)
		}
		return r, nil
	default:
		return nil, fmt.Errorf("cannot unmarshal caption column type: %T", t)
	}
}

func fz(Z []complex128) []float64 {
	if Z == nil {
		return nil
	}
	r := make([]float64, 2*len(Z))
	for i, z := range Z {
		r[2*i], r[2*i+1] = real(z), imag(z)
	}
	return r
}
func zf(F []float64) []complex128 {
	if F == nil {
		return nil
	}
	r := make([]complex128, len(F)/2)
	for i := range r {
		r[i] = complex(F[2*i], F[2*i+1])
	}
	return r
}

/* shadow a specific field (embed&alias)
https://choly.ca/post/go-json-marshalling/
type MyUser struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"lastSeen"`
}
func (u *MyUser) MarshalJSON() ([]byte, error) {
	type Alias MyUser
	return json.Marshal(&struct {
		LastSeen int64 `json:"lastSeen"`
		*Alias
	}{
		LastSeen: u.LastSeen.Unix(),
		Alias:    (*Alias)(u),
	})
}
func (u *MyUser) UnmarshalJSON(data []byte) error {
	type Alias MyUser
	aux := &struct {
		LastSeen int64 `json:"lastSeen"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.LastSeen = time.Unix(aux.LastSeen, 0)
	return nil
}
*/
