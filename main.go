package main

import (
	"fmt"
	"sync"
	"time"
)

func main(){
	var wg sync.WaitGroup
	wg.Add(2)

	sl1Signal := make(chan bool)
	sl2Signal := make(chan bool)

	go LaunchSite1(sl1Signal,sl2Signal,&wg)
	go LaunchSite2(sl1Signal,sl2Signal,&wg)

	wg.Wait()
}

func LaunchSite1(sl1Signal,sl2Signal chan bool, wg *sync.WaitGroup){
	defer wg.Done()

	satId := 1
	for satId <= 500 {
		var satGroup sync.WaitGroup
		satGroup.Add(4)
		for i:=0;i<4;i++{

			go func(satId int){
				defer satGroup.Done()

				for countDown:=10;countDown>=0;countDown--{
					time.Sleep(time.Second)
					fmt.Println("Satellite ",satId, " from SL1 countdown: ",countDown)
				}

				fmt.Println("******** Satellite ",satId, " Launched from SL1 ********")
			}(satId)

			satId++
		}
		satGroup.Wait()
		satId += 4
		sl2Signal <- true // sending SL2 signal that SL1 has launched its group of 4 Satellites
		<- sl1Signal
	}
}

func LaunchSite2(sl1Signal,sl2Signal chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	satId := 5
	for satId <= 500 {
		<- sl2Signal
		var satGroup sync.WaitGroup
		satGroup.Add(4)
		for i:=0;i<4;i++{

			go func(satId int){
				defer satGroup.Done()

				for countDown:=10;countDown>=0;countDown--{
					time.Sleep(time.Second)
					fmt.Println("Satellite ",satId, " from SL2 countdown: ",countDown)
				}

				fmt.Println("******** Satellite ",satId, " Launched from SL2 ********")
			}(satId)

			satId++
		}
		satGroup.Wait()
		satId += 4
		sl1Signal <- true // sending SL1 signal that SL2 has launched its group of 4 Satellites
	}
	close(sl1Signal)
	<-sl2Signal
}
