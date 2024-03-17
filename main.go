package main

import (
	"fmt"
	"sync"
	"time"
)

type Beans struct {
	weightGrams int
	// indicate some state change? create a new type?
	ground bool // Added to indicate if the beans are ground or not
}

type Grinder struct {
	gramsPerSecond int
	busy           bool
	mu             sync.Mutex // To control access to the grinder
}

type Brewer struct {
	// assume we have unlimited water, but we can only run a certain amount of water per second into our brewer + beans
	ouncesWaterPerSecond int
	busy                 bool
	mu                   sync.Mutex // To control access to the brewer
	maxWaterPerSecond    int
}

type Order struct {
	id                   int // maybe uuid
	ouncesOfCoffeeWanted int
	coffeeStrength       int // grams of beans per ounce of coffee (default to 2g)
}

type CoffeeShop struct {
	grinders                 []*Grinder
	brewers                  []*Brewer
	totalAmountUngroundBeans int
	baristas                 int
	orders                   chan Order // channel to manage orders
}

type Coffee struct {
	// should hold size maybe?
	size int
}

func (g *Grinder) Grind(beans Beans) Beans {
	// how long should it take this function to complete?
	// i.e. time.Sleep(XXX)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.busy = true
	defer func() {
		g.busy = false
	}()

	// Calculate grinding time
	gramsToGrind := beans.weightGrams
	timeToGrind := time.Duration(gramsToGrind/g.gramsPerSecond) * time.Second
	time.Sleep(timeToGrind)
	return Beans{weightGrams: beans.weightGrams, ground: true}
}

func (b *Brewer) Brew(beans Beans) Coffee {
	// how long should it take this function to complete?
	// i.e. time.Sleep(YYY)

	b.mu.Lock()
	defer b.mu.Unlock()
	b.busy = true
	defer func() {
		b.busy = false
	}()
	// Calculate brewing time, assume we need 6 ounces of water for every 12 grams of beans
	ouncesOfWater := 6 * (beans.weightGrams / 12)
	timeToBrew := time.Duration(ouncesOfWater/b.ouncesWaterPerSecond) * time.Second
	time.Sleep(timeToBrew)
	return Coffee{size: ouncesOfWater}
}

func NewCoffeeShop(grinders []*Grinder, brewers []*Brewer, baristas int) *CoffeeShop {
	return &CoffeeShop{
		grinders: grinders,
		brewers:  brewers,
		baristas: baristas,
		orders:   make(chan Order, 100), // set capacity to hold 100 incoming orders
	}
}

func (cs *CoffeeShop) MakeCoffee(order Order) {
	fmt.Printf("Placing order ID %d for %d ounces of coffee onto the channel\n", order.id, order.ouncesOfCoffeeWanted)
	go func() {
		cs.orders <- order
	}()
}

func (cs *CoffeeShop) Start() {
	for i := 0; i < cs.baristas; i++ {
		go func(id int) {
			for order := range cs.orders {
				// process order
				fmt.Printf("Barista %d is processing order %d\n", id, order.id)
				// find available grinder
				for _, grinder := range cs.grinders {
					if !grinder.busy {
						ungroundBeans := Beans{weightGrams: order.ouncesOfCoffeeWanted * order.coffeeStrength}
						groundBeans := grinder.Grind(ungroundBeans)
						// find available brewer
						for _, brewer := range cs.brewers {
							if !brewer.busy {
								coffee := brewer.Brew(groundBeans)
								fmt.Printf("Order %d completed: %d ounces of coffee\n", order.id, coffee.size)
								break
							}
						}
						break
					}
				}

			}
		}(i)
	}
}

func main() {
	// Premise: we want to model a coffee shop. An order comes in, and then with a limited amount of grinders and
	// brewers (each of which can be "busy"): we must grind un-ground beans, take the resulting ground beans, and then
	// brew them into liquid coffee. We need to coordinate the work when grinders and/or brewers are busy doing work
	// already. What Go data structure(s) might help us coordinate the steps: order -> grinder -> brewer -> coffee?
	//
	// Some of the struct types and their functions need to be filled in properly. It may be helpful to finish the
	// Grinder impl, and then Brewer impl each, and then see how things all fit together inside CoffeeShop afterward.

	//b := Beans{weightGrams: 10}
	g1 := &Grinder{gramsPerSecond: 5}
	g2 := &Grinder{gramsPerSecond: 3}
	g3 := &Grinder{gramsPerSecond: 12}

	b1 := &Brewer{ouncesWaterPerSecond: 100}
	b2 := &Brewer{ouncesWaterPerSecond: 25}

	cs := NewCoffeeShop([]*Grinder{g1, g2, g3}, []*Brewer{b1, b2}, 2)

	cs.Start()

	numCustomers := 10
	for i := 0; i < numCustomers; i++ {
		// in parallel, all at once, make calls to MakeCoffee
		fmt.Printf("i: %d \n", i)
		orderId := i
		go func() {
			cs.MakeCoffee(Order{
				id:                   orderId + 1,
				ouncesOfCoffeeWanted: 12,
				coffeeStrength:       2,
			})
		}()
	}

	time.Sleep(10 * time.Second) // Wait some time for orders to be processed
	// Issues with the above
	// 1. Assumes that we have unlimited amounts of grinders and brewers.
	//		- How do we build in logic that takes into account that a given Grinder or Brewer is busy?
	// 2. Does not take into account that brewers must be used after grinders are done.
	// 		- Making a coffee needs to be done sequentially: find an open grinder, grind the beans, find an open brewer,
	//		  brew the ground beans into coffee.
	// 3. A lot of assumptions (i.e. 2 grams needed for 1 ounce of coffee) are left as comments in the code.
	// 		- How can we make these assumptions configurable, so that our coffee shop can serve let's say different
	//		  strengths of coffee via the Order that is placed (i.e. 5 grams of beans to make 1 ounce of coffee)?
}
