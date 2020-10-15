/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	pb "helloworld"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	//address     = ":50051"
	address     = "dist29:50051"
	defaultName = "Bro"
	clientName  = "CAMIONES"
)

func getInput() string {
	fmt.Println("Inserte nombre: ")
	var input string
	fmt.Scanln(&input)
	return input
}

//items contiene info acerca de un producto
type Items struct {
	id    string
	tipo  string
	valor string
	src   string
	dest  string
	reply string
	date  string
}

// Envio retorna si el envio se hace o no se hace
func Envio() bool {
	in := []int{0, 1, 1, 1, 1}
	randomIndex := rand.Intn(len(in))
	pick := in[randomIndex]

	if pick == 1 {
		return true
	}
	return false
}

func realizarEnvio(c pb.GreeterClient, tipo string) {

	// esto dentro del codigo de camiones
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.TrakingRequest(ctx, orden)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	if tipo == "retail" {

	} else if tipo == "pyme" {

	}

}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	//waiting Time
	fmt.Println("\nIngrese Tiempo de Envio de Cada Paquete")
	waitingTime, _ := strconv.Atoi(getInput(2))
	fmt.Printf("\nTiempo: %d\n", waitingTime)

	// Contact the server and print out its response.

	for {
		realizarEnvio(c, "retail")
		time.Sleep(3 * time.Second)
		realizarEnvio(c, "retail")
		time.Sleep(3 * time.Second)
		realizarEnvio(c, "pyme")
		time.Sleep(3 * time.Second)
	}
}
