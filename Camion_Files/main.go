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
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	//address     = ":50051"
	address     = "dist29:50052"
	defaultName = "Bro"
	clientName  = "CAMIONES"
)

//items contiene info acerca de un producto
type Items struct {
	id    string
	tipo  string
	valor string
	src   string
	dest  string
	reply string
	date  string
	atts  string
}

/*
	Envio()
	funcion para hacer un envio con probabilidad de succes de 0.8
	Input: nada
	returns: bool, estado de envio
*/
func Envio() bool {
	in := []int{0, 1, 1, 1, 1}
	randomIndex := rand.Intn(len(in))
	pick := in[randomIndex]

	if pick == 1 {
		return true
	}
	return false
}

/*
	getTime()
	Obtiene el tiempo de la maquina
	Input: no tiene
	returns: fecha actual, string
*/
func getTime() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

/*
	realizarEnvio()
	dado una orden se realiza el envio, haciendo registro de este
	Input: pb.GreeterClient c, string tipo, int intentoTime, *os.File f
	returns: nada
*/
func realizarEnvio(c pb.GreeterClient, tipo string, intentoTime int, f *os.File) {

	// esto dentro del codigo de camiones

	// PEDIR Y RECIBIR UN PAQUETE
	dat := &pb.DeliveryRequest{
		R: tipo,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	received, err := c.SendInformation(ctx, dat)
	if err != nil {
		log.Fatalf("\ncould not greet: %v", err)
	}

	if received.OrderID != "" {
		var try bool

		var intento int
		var Nintentos int

		var IntentoFinal string
		var enviado bool = false
		var newEstado string = "En Bodega"

		var precio string = received.GetProductValue()
		value, _ := strconv.Atoi(precio)
		received.GetOrderID()

		fmt.Println("Orden Code: %s", received.GetOrderID())
		fmt.Println("Estado: En Bodega")
		fmt.Println("Saliendo!!!")
		fmt.Println("Estado: En camino")
		newEstado = "En Camino"

		if tipo == "retail" {

			fmt.Println("Realizo pedido de retail")
			for intento = 1; intento < 4; intento++ {
				//hace cosas

				try = Envio()
				fmt.Println("Intento enviar")
				if try {
					IntentoFinal = strconv.Itoa(intento)
					newEstado = "Recibido"
					fmt.Println("Envio Realizado")
					fmt.Println("Nuevo estado: Recibido")
					enviado = true
					break
				}
				// ACTUALIZAR ESTADO PAQUETE
				dat := &pb.StatusResponse{
					TrackingCode: received.OrderID,
					Status:       newEstado,
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				received1, err := c.TrackingStatusUpdate(ctx, dat)
				if err != nil {
					log.Fatalf("\ncould not greet with retail: %v%s", err, received1)
				}

				// tiempo de espera despues de un envio
				time.Sleep(time.Duration(intentoTime) * time.Second)
			}
			if !try && enviado == false {
				IntentoFinal = "3"
				newEstado = "No Recibido"
				fmt.Println("Fallo al enviar")
				fmt.Println("Estado: No Recibido")
			}
		} else { //pyme
			fmt.Print("\nRealizo pedido de pyme")
			if value < 10 {
				Nintentos = 1 // 1 base + 0 extra
			} else if value < 20 {
				Nintentos = 2 // 1 base + 1 extra
			} else { // mayor a 20
				Nintentos = 3 // 1 base + 2 extra
			}

			for intento = 1; intento < Nintentos+1; intento++ {
				//hace cosas

				try = Envio()

				if try {
					IntentoFinal = strconv.Itoa(intento)
					newEstado = "Recibido"
					fmt.Println("Envio Realizado")
					fmt.Println("Nuevo estado: Recibido")
					enviado = true
					break
				}
				// ACTUALIZAR ESTADO PAQUETE
				dat := &pb.StatusResponse{
					TrackingCode: received.OrderID,
					Status:       newEstado,
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				received1, err := c.TrackingStatusUpdate(ctx, dat)
				if err != nil {
					log.Fatalf("\ncould not greet with pyme: %v%s", err, received1)
				}

				// tiempo de espera despues de un envio
				time.Sleep(time.Duration(intentoTime) * time.Second)
			}
			if !try && enviado == false {
				IntentoFinal = strconv.Itoa(Nintentos)
				newEstado = "No Recibido"
				fmt.Println("Fallo al enviar")
				fmt.Println("Estado: No Recibido")
			}

		}

		received.Attempts = IntentoFinal
		fmt.Printf("\nNumero de intentos: %s\n", IntentoFinal)

		// agregar numero de intentos
		orderUpdate := &pb.StatusResponse{
			TrackingCode: received.OrderID,
			Status:       newEstado,
			Attempts:     IntentoFinal,
		}
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		m, err := c.TrackingStatusFinal(ctx, orderUpdate)
		if err != nil {
			log.Fatalf("\ncould not greet at the end: %v\n\tTrackingcode: %s\n\tStatus: %s%s\n", err, received.OrderID, newEstado, m)
		}

		t := getTime()
		if newEstado != "Envio Realizado" {
			t = "0"
		}
		toFile := fmt.Sprintf("%s,%s,%s,%s,%s,%s\n", received.OrderID, received.ProductValue, received.Src, received.Dest, IntentoFinal, t)
		_, err = f.WriteString(toFile)
		check(err)

	} else {
		fmt.Println("\nNo hay ordenes pendientes")
	}

}

/*
	camion()
	helper que usa la coneccion con logistica y hace envios de paquetes en base a intervalos de tiempo (pedidoTime e intentoTime), ademas lleva el registro de entregas en un archivo csv
	Input: conexion c, string tipo, int intento time, int pedidoTime, archivo f
	returns: nada
*/
func camion(c pb.GreeterClient, tipo string, intentoTime int, pedidoTime int, f *os.File) {

	realizarEnvio(c, tipo, intentoTime, f)

	// tiempo de espera despues de un envio
	time.Sleep(time.Duration(pedidoTime) * time.Second)

	realizarEnvio(c, tipo, intentoTime, f)

}

/*
	getInput()
	pide input de un usuario segun corresponda la situacion y retorna la eleccion
	Input: int x
	returns: string input, eleccion
*/
func getInput(x int) string {
	if x == 1 {
		fmt.Print("\nIngrese ID producto a ordenar: ")
	} else if x == 3 {
		fmt.Print("\nIngrese ID producto a Realizar Seguimiento: ")
	} else {
		fmt.Print("\nSeleccion: ")
	}
	var input string
	fmt.Scanln(&input)
	return input
}

/*
	check()
	comprueba un error y lo muestra
	Input: error e
	returns: nada
*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("\ndid not connect with server: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	//waiting Time 1
	fmt.Println("\nIngrese tiempo entre cada intento de envio de paquete")
	intentoTime, _ := strconv.Atoi(getInput(2))
	fmt.Printf("\nTiempo: %d\n", intentoTime)

	//waiting Time 2
	fmt.Println("\nIngrese tiempo de envio entre cada pedido")
	pedidoTime, _ := strconv.Atoi(getInput(2))
	fmt.Printf("\nTiempo: %d\n", pedidoTime)

	// Contact the server and print out its response.

	// open csvs
	f, er := os.Create("registry_truck_retail_1.csv")
	check(er)
	f1, err1 := os.Create("registry_truck_retail_2.csv")
	check(err1)
	f2, err2 := os.Create("registry_truck_pyme_1.csv")
	check(err2)

	toFile := "OrderID,ProductValue,Origen,Destino,Intentos,hora\n"
	_, e := f.WriteString(toFile)
	check(e)
	_, e1 := f1.WriteString(toFile)
	check(e1)
	_, e2 := f2.WriteString(toFile)
	check(e2)

	for {
		go camion(c, "retail", intentoTime, pedidoTime, f)
		time.Sleep(time.Duration(intentoTime+pedidoTime) * time.Second)

		go camion(c, "retail", intentoTime, pedidoTime, f1)
		time.Sleep(time.Duration(intentoTime+pedidoTime) * time.Second)

		camion(c, "pyme", intentoTime, pedidoTime, f2)
		time.Sleep(time.Duration(intentoTime+pedidoTime) * time.Second)
	}

	defer f.Close()
	defer f1.Close()
	defer f2.Close()
}
