package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	pb "github.com/bryutus/protobuf-tutorial/tutorial"
	"github.com/golang/protobuf/proto"
)

func promptForAddress(r io.Reader) (*pb.Person, error) {
	// A protocol buffer can be created like any struct.
	p := &pb.Person{}

	rd := bufio.NewReader(r)
	fmt.Print("Enter person ID number:\n")
	// An int32 field in the .proto file is represented as an int32 field
	// in the generated Go struct.
	if _, err := fmt.Fscanf(rd, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Printf("Enter name:\n")
	name, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	// A string field in the .proto file results in a string field in Go.
	// We trim the whitespacee because rd.ReadString includes the trailing
	// newline character in its output
	p.Name = strings.TrimSpace(name)

	fmt.Printf("Enter email address (blank for none):\n")
	email, err := rd.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter a phone number (or leave blank to finish):\n")
		phone, err := rd.ReadString('\n')
		if err != nil {
			break
		}

		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}

		// The PhoneNumber message type is nested within the Person
		// message in the .proto file. This results in a Go struct
		// named using name of the parent prefixed to the name of
		// the nested message. Just s with pb.Person, it can be
		// created like any other struct.
		pn := &pb.Person_PhoneNumber{
			Number: phone,
		}

		// A repeated proto field mapps to a slice field in Go. We can
		// append to it like any other slice.
		p.Phones = append(p.Phones, pn)
	}

	return p, nil
}

// Main reads the entire address book from a file, adds one person based on
// user input, then writes it back out to the same file.
func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	// Read the existing address book.
	in, err := ioutil.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found. Creating new file.\n", fname)
		} else {
			log.Fatalln("Error reading file:", err)
		}
	}

	// [START marshal_proto]
	book := &pb.AddressBook{}
	// [START_EXCLUDE]
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	// Add an address.
	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("Error with address:", err)
	}
	book.People = append(book.People, addr)
	// [End_EXCUTE]
}
