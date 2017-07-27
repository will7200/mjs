package main

/*
func byJobTime(a, b interface{}) int {
	c1 := a.(job.Job)
	c2 := b.(job.Job)

	switch {
	case c1.NextRunAt.After(c2.NextRunAt):
		return 1
	case c2.NextRunAt.Before(c2.NextRunAt):
		return -1
	default:
		return 0
	}
}

func withint64(a, b interface{}) int {
	c1 := a.(int64)
	c2 := b.(int64)

	switch {
	case c1 > c2:
		return 1
	case c1 < c2:
		return -1
	default:
		return 0
	}
}
func Clone(a, b interface{}) {

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}
func main() {
	tree := rbt.NewWith(withint64)
	j := &job.Job{Name: "Test", NextRunAt: time.Now()}
	jj := &job.Job{Name: "Test3", NextRunAt: time.Now().Add(5 * time.Second)}
	tt := &job.Job{}
	Clone(j, tt)
	tt.Name = "Test2"
	tt.NextRunAt = tt.NextRunAt.Add(10 * time.Second)
	tree.Put(j.NextRunAt.Unix(), j)
	tree.Put(tt.NextRunAt.Unix(), tt)
	tree.Put(jj.NextRunAt.Unix(), jj)
	fmt.Println(tree)
	a := tree.Left()
	fmt.Printf("%+v\n", a.Value)
	b := tree.Right()
	c, _ := tree.Floor(int64(1))
	d, _ := tree.Ceiling(int64(1))
	_, _, _, _ = a, b, c, d
}
*/
