package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type Order struct{
    Id int
    Amount float64
    Description string
}

func processOrder(ctx context.Context ,workerId int ,orderChannel <-chan Order , wg *sync.WaitGroup , workerPool chan struct{}){
    defer wg.Done()

    for {
        select {
        case order , ok := <- orderChannel:
            if !ok {
                return
            }
            select {
            case workerPool <- struct{}{}:
            case <-ctx.Done():
                fmt.Printf("Worker %d cancelado antes de processar o pedido %d\n", workerId, order.Id)
                return
            }
            fmt.Printf("Worker: %d\t Procesando o pedido ID:%d\n" ,workerId, order.Id)
            select {
            case <-time.After(2 * time.Second):
                fmt.Println("Pedido processado")
            case <- ctx.Done():
                fmt.Printf("worker cancelado enquanto processava %d\n",order.Id)
                return
            }
            <-workerPool
        case <-ctx.Done():
            fmt.Println("Worker cancelado")
            return
        }
    }
}

func main(){
    workers := 3
    orderChannel := make(chan Order,10)
    workerPool := make(chan struct{},workers)
    var wg sync.WaitGroup

    ctx , cancel := context.WithTimeout(context.Background(),2*time.Second)
    defer cancel()

    for i := 0; i <= 5 ; i++{
        wg.Add(1)
        go processOrder(ctx ,i ,orderChannel,&wg ,workerPool)
    }

    for i:= 0 ; i < 15 ; i++{
        fmt.Printf("Pedido %d recebido\n" , i)
        orderChannel<-Order{Id:i,Amount: float64(i) * 10.0,Description:fmt.Sprintf("Pedido %d" , i)}
    }

    close(orderChannel)
    wg.Wait()
    fmt.Println("Pedidos finalizados")
}
