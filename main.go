package main
import (
    "log"
    "net/http"
)

const (
    addr = "127.0.0.1:4000"
    resumePath = "/"
)

func resumeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "Метод запрещён!", 405)
        return
    }

    log.Println(r.URL.Query().Get("id"))
    w.Write([]byte("Hello from constellation"))
}


func main() {
    mux := http.NewServeMux()
    mux.HandleFunc(resumePath, resumeHandler)

    log.Println("Запуск веб-сервера по адресу", addr)
    err := http.ListenAndServe(addr, mux)
    if err != nil {
        log.Fatal(err)
    }
}
