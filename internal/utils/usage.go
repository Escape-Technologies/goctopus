package utils

import (
	"flag"
	"fmt"
)

func PrintUsage() {
	PrintASCII()
	fmt.Println("Usage: goctopus [options] [addresses]")
	fmt.Println("[addresses]: A list of addresses to fingerprint, comma separated.\nAddresses can be in the form of http://example.com/graphql or example.com.\n If an input file is specified, this argument is ignored.")
	fmt.Println("[options]:")
	flag.PrintDefaults()
}
