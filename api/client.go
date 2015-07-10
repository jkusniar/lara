package api

import (
	"fmt"
	"log"
)

// TODO is there a need for these types/functions to be public (UpperCase) ?
// They are called only through reflection

type CreateClientRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r *CreateClientRequest) String() string {
	return fmt.Sprint("firstName: ", r.FirstName, ", lastName: ", r.LastName)
}

type CreateClientResponse struct {
	Id string `json:"id"`
}

func (r *CreateClientResponse) String() string {
	return fmt.Sprint("Id: ", r.Id)
}

func CreateClient(req *CreateClientRequest) (*CreateClientResponse, error) {
	log.Println("Creating client ", req.FirstName, ", ", req.LastName)
	return &CreateClientResponse{"f47ac10b-58cc-4372-a567-0e02b2c3d479"}, nil
}
