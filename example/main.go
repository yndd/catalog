package main

import (
	"fmt"

	"github.com/yndd/catalog"
	_ "github.com/yndd/catalog/vendors/all"

	targetv1 "github.com/yndd/target/apis/target/v1"
)

func main() {
	c := catalog.Default
	for _, f := range c.List() {
		fmt.Printf("%+v\n", f)
	}
	fn, err := c.Get(catalog.Key{
		Name:      "configure_lldp",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeNokiaSRL,
		Platform:  "",
		SwVersion: "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", fn)
}
