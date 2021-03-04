package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"sync"
	"time"
)

//熔断器
//最大并发数

type Product struct {
	ID int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r<6 {
		time.Sleep(time.Second*5)
	}
	return Product{
		ID:    101,
		Title: "go",
		Price: 12,
	},nil

}

//超时之后会显示一个推荐商品
func RecProduct() (Product, error) {
	return Product{
		ID:    102,
		Title: "推荐商品",
		Price: 22,
	},nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	configA := hystrix.CommandConfig{
		Timeout:                4000,
	}//timeout 设置为4s，则不会出现hystrix：timeout

	hystrix.ConfigureCommand("get_prod",configA)
	hystrix.ConfigureCommand("get_prod2", hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,//支持最大并发数为5
		RequestVolumeThreshold: 20,//有20个请求才进行错误百分比计算
		SleepWindow:            5,//熔断器：默认关闭，请求次数异常超过设定比例则打开，打开后直接执行降级函数，半开为定期打开
		   						  //过了5s则尝试服务是否可用。默认5s
		ErrorPercentThreshold:  20,//超过20%熔断器打开，错误百分比，默认为50%
	})
	circuit, _, _ := hystrix.GetCircuit("get_prod")

	resultChan := make(chan Product, 1)

	wg := sync.WaitGroup{}
	//设置最大并发数
	for i:=0; i<10; i++ {
		wg.Add(1)
		defer wg.Done()
		go func() {
			errs := hystrix.Do("get_prod2", func() error {
				p, _ := getProduct() //这里会随机延迟三秒
				resultChan <- p
				//fmt.Println(p)
				return nil
			}, func(err error) error { //降级函数
				//fmt.Println(RecProduct())
				rcp, err := RecProduct()
				resultChan <- rcp //推荐商品
				return nil
			})
			//主要是此方法
			errs2 := hystrix.Go("get_prod", func() error {
				p, _ := getProduct() //这里会随机延迟三秒
				resultChan <- p
				//fmt.Println(p)
				return nil
			}, func(err error) error { //降级函数
				//fmt.Println(RecProduct())
				//rcp, err := RecProduct()
				//resultChan <- rcp	//推荐商品

				return errors.New("prod errors")
			})

			select {
			case getProd := <-resultChan:
				fmt.Println(getProd)
			case err := <-errs2:
				fmt.Println(err)
			}
			if errs != nil {
				fmt.Println(errs)
			}
			if errs2 != nil {
				fmt.Println(errs2)
			}
			fmt.Println(circuit.IsOpen())
			time.Sleep(time.Second * 1)
		}()
	}
	wg.Wait()
}

/*
get_prod
如果r<6
hystrix: timeout
hystrix: timeout
{101 go 12}
{101 go 12}
{101 go 12}
hystrix: timeout

 */


/*
get_prod2,超时会显示推荐商品

{102 推荐商品 22} <nil>
{101 go 12}
{102 推荐商品 22} <nil>
{101 go 12}
{102 推荐商品 22} <nil>
{101 go 12}
{101 go 12}

 */