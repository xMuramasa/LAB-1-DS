// reciver / consummer
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/streadway/amqp"
)

// Items2
// Struct para formatear el struct item y leer como json
type Items2 struct {
	Id          string `json:"id"`
	Order_type  string `json:"order_type"`
	Order_value string `json:"order_value"`
	Tracking    string `json:"tracking"`
	Status      string `json:"status"`
	Atts        string `json:"atts"`
}

//struct Balance, struct que guarda los datos a escribir  de una orden
type Balance struct {
	Id       string
	Tracking string
	tipo     string
	Atts     float64
	total    float64
	ganancia float64
	perdida  float64
}

/*
	failOnError()
	comprueba un mensaje de error y lo muestra
	Input: error, string
	returns: nada
*/
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

/*
	SetupCloseHandler()
	se ocupa de mostrar todo cuando hay ctr+c
	Input: nada
	returns: nada
*/
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		fmt.Printf("\nBALANCE GENERAL:\n")
		fmt.Printf("GANANCIAS: %f, PERDIDAS: %f, ENVIOS TOTALES: %d, ENVIOS NO ENTREGADOS: %d, ENVIOS ENTREGADOS: %d\n", gananciasTotal, perdidasTotal, enviosTotales, enviosNoEntregados, enviosEntregados)
		os.Exit(0)
	}()
}

var gananciasTotal = 0.0
var perdidasTotal = 0.0
var enviosEntregados = 0
var enviosNoEntregados = 0
var enviosTotales = 0

func main() {
	// Inicaimos conexion
	conn, err := amqp.Dial("amqp://test:test@10.6.40.169:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// se abre el canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// creacion de cola
	q, err := ch.QueueDeclare(
		"hello-queue", // name
		false,         // durable
		false,         // delete when usused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//se reciven msg del la cola
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//bloquea la ejecucion del main.go hasta que recibe un valor
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)

			var reading Items2

			err = json.Unmarshal([]byte(d.Body), &reading)
			if err != nil {
				log.Fatalf("oh shoiit: %v", err)
			}

			var balance Balance

			//valor del producto
			valorProducto, _ := strconv.ParseFloat(reading.Order_value, 64)

			//intentos del producto
			intentos, _ := strconv.ParseFloat(reading.Atts, 64)

			//pedidas del producto
			var perdidas float64 = 10 * float64(intentos-1.0)

			// valor a ganar aka BALANCE
			var total float64
			total = math.Abs(valorProducto - perdidas)

			//Asignacion base
			balance.Id = reading.Id
			balance.Tracking = reading.Tracking
			balance.tipo = reading.Order_type
			balance.Atts = intentos

			if reading.Status == "Recibido" {
				balance.total = total
				balance.ganancia = valorProducto
				balance.perdida = perdidas
				enviosEntregados = enviosEntregados + 1
			} else {
				//No Recibido
				if reading.Order_type == "Normal" {
					// Normal
					balance.total = total
					balance.ganancia = 0.0
					balance.perdida = perdidas
				} else if reading.Order_type == "prioritario" {
					//Prioritario
					newValue := valorProducto * 0.3

					balance.total = math.Abs(newValue - perdidas)
					balance.ganancia = newValue
					balance.perdida = perdidas
				} else {
					//Retail
					balance.total = total
					balance.ganancia = valorProducto
					balance.perdida = valorProducto
				}
				enviosNoEntregados = enviosNoEntregados + 1
			}
			gananciasTotal = gananciasTotal + balance.ganancia
			perdidasTotal = perdidasTotal + balance.perdida
			fmt.Println(balance)
			// Guardar Registro PAPOPE
			enviosTotales = enviosTotales + 1
		}
	}()
	SetupCloseHandler()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
