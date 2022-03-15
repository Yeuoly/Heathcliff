package types

type Buffer struct {
	Buf []byte
	N   int
}

func (c *Buffer) Takeover(buf []byte, n int) {
	c.Buf = buf
	c.N = n
}

func (c *Buffer) Malloc(n int) {
	if n <= 0 {
		return
	}

	if c.Buf == nil {
		c.Buf = make([]byte, n)
		c.N = n
	} else {
		if n > c.N {
			c.Buf = make([]byte, n)
		}
	}
}
