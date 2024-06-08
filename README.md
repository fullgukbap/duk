# duk
Elegant TCP IP Input Ouput library

# Support feature
server
- Event-driven TCP IP communication 
- Interval Broadcast
- Sugar coding with Hooks 

# Example Code
```go
func main() {
	app := duk.New()

    // interval broadcast할 변수 등록
	app.Register(30, struct{}{})

	app.OnConnect(func(id string, conn net.Conn) {
		fmt.Println("새로운 연결!")
	})

	app.OnDisconnect(func(id string, err error) {
		fmt.Println("끊김!")
	})

	app.On("welcome", func(c *duk.Ctx) {
		payload := struct{ Name string }{}
		fmt.Println("사용자 요청 값: ", c.Parse(&payload))
	})

	log.Fatal(app.Listen(":8080"))
}
```
