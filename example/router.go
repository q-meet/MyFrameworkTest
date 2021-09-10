package example

import (
	"fmt"
	"frames/binding"
	"frames/httprouter"
	"log"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\\n", ps.ByName("name"))
}

// 写一个中间件
func timeMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		timeStart := time.Now()

		next(w, r, ps)

		timeElapsed := time.Since(timeStart)
		fmt.Println("处理请求耗时 : ", timeElapsed)
	}
}

// Booking contains binded and validated data.
type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func getBookable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var b Booking
	fmt.Println(ps)
	if err := binding.Query.Bind(r, &b); err == nil {
		fmt.Println(b)
		fmt.Fprintf(w, "message Booking dates are valid!")
	} else {
		fmt.Println(b)
		fmt.Fprintf(w, "error!"+err.Error())
	}
}

func Router(){

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/bookable", getBookable)
	router.GET("/bookable1", timeMiddleware(getBookable))

	log.Fatal(http.ListenAndServe(":8085", router))
}
