package main

import (
	"context"
	"fmt"
    "math/rand"
	"sync"
	"time"
)

type Order struct{
    Id int
    Amount float64
    Description string
}

func processOrder(ctx context.Context ,workerId int ,orderChannel <-chan Order , unprocessedChannel chan <- Order , wg *sync.WaitGroup , workerPool chan struct{} , processed *int){
    defer wg.Done()

    for {
        select {
        case order , ok := <- orderChannel:
            if !ok {
                return
            }
            
            fmt.Printf("Pedido %d recebido\n" , order.Id)
            
            select {
            case workerPool <- struct{}{}:
            case <-ctx.Done():
                fmt.Printf("Worker %d cancelado antes de processar o pedido %d\n", workerId, order.Id)
                unprocessedChannel<-order
                return
            }
            
            fmt.Printf("Worker: %d\t Procesando o pedido ID:%d\n" ,workerId, order.Id)


            processingSuccess := process(order.Id)

            if !processingSuccess {
                fmt.Printf("Worker: %d\t falhou ao processar o pedido ID:%d\n" ,workerId, order.Id)
                unprocessedChannel<-order
                <-workerPool
                continue
            }
            
            fmt.Printf("Pedido ID:%d processado\n" , order.Id)
            *processed++
            <-workerPool

            
        case <-ctx.Done():
            fmt.Println("Worker cancelado")
            return
        }
       
    }
}

func process(id int) bool {
    if id%2 == 0{
        return true
    }
    return false
}

func handleProcessOrders(orderChannel , unprocessedChannel chan Order,workerPool chan struct{} , items []Order , processed *int){
    var wg sync.WaitGroup

    ctx , cancel := context.WithTimeout(context.Background(),15*time.Second)
    defer cancel()
    
    for i := 0; i <= 5 ; i++{
        wg.Add(1)
        go processOrder(ctx ,i ,orderChannel,unprocessedChannel,&wg ,workerPool , processed)
    }

    for _ , v := range items{
        orderChannel<-v
        time.Sleep(1*time.Second)
    }

    close(orderChannel)
    wg.Wait()

    close(unprocessedChannel)
}

func runOrders(items []Order , processed *int) []Order{
    workers := 3
    workerPool := make(chan struct{},workers)
    orderChannel := make(chan Order,10)
    unprocessedChannel := make(chan Order,10)
    var unprocessedItems []Order
    handleProcessOrders(orderChannel,unprocessedChannel,workerPool,items,processed)
    
    for order := range unprocessedChannel {
        fmt.Printf("Id do pedido: %d\n" ,order.Id)
        unprocessedItems = append(unprocessedItems , order)
    }
    return unprocessedItems;
}

func reprocessOrders(items []Order , processed *int) {
    maxRetries := 5
    retries:= 0
   
    for len(items) > 0 && retries < maxRetries {
        delay := 1<<retries * time.Second
        time.Sleep(delay)
        fmt.Println("REPROCESSANDO PEDIDOS .... ",delay)
        
        for i ,v := range items{
            fmt.Printf("Pedido %d\n",v.Id )
            numero := rand.Intn(100)
            items[i].Id=numero
        }

        items = runOrders(items,processed)
        

        retries++
       
    }
}

func seedItems(items *[]Order){
    for i:= 0 ; i < 15 ; i++{
        *items = append(*items , Order{Id:i,Amount: float64(i) * 10.0,Description:fmt.Sprintf("Pedido %d" , i)})
    }
}

func main(){
    timeStart := time.Now()
    var items []Order
    processed := 0
    seedItems(&items)
    fmt.Println("Total de pedidos enviados: ",len(items))
    items = runOrders(items , &processed)

    reprocessOrders(items ,&processed)

    fmt.Println("Pedidos finalizados" , processed , "tempo decorrido: " ,time.Since(timeStart))
}
