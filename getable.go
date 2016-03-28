package core

import "github.com/wolfgarnet/REST"

type Getable interface  {
	Get(token string) REST.Node
}