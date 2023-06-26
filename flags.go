package main

type sliceFlag []string

func (i *sliceFlag) String() string {
	return "my string representation"
}

func (i *sliceFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
