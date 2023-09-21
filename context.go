package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("result: ", val)
	fmt.Println("took: ", time.Since(start))
}

// func fetchUserData(ctx context.Context, userID int) (int, error) {
// 	val, err := fetchThirdPartyStuffWhichCanBeSlow()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return val, nil
// }

// Struct is defined as per the return type of fetchUserData()
type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
    
	//since go routines returns nothing so create a channel of Response return type
	respch := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		//value and err pushed to Response struct and then to respch channel
		respch <- Response{
			value: val,
			err:   err,
		}
	}()

	//Accessing go routines in for select pattern
	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("fetching error from third party took too long to respond")
			// respch channel values are pushed into resp variable
		case resp := <-respch:
			return resp.value, resp.err
		}
	}
}

func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	time.Sleep(time.Millisecond * 500)
	return 666, nil
}
