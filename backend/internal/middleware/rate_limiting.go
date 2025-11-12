package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor armazena o limiter e a última vez que foi visto.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Mapa para armazenar os limiters por IP
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// cleanUpVisitors remove entradas antigas do mapa periodicamente.
func init() {
	go cleanUpVisitors()
}

func cleanUpVisitors() {
	for {
		time.Sleep(1 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimitMiddleware é o middleware que aplica o limite de taxa.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtém o endereço IP do requisitante.
		// Em produção real, você deve priorizar o 'X-Forwarded-For' se estiver atrás de um proxy.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Erro interno", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			// Cria um novo limiter para este IP.
			// Permite 5 requisições por segundo (Limit), com um "burst" (pico) de 10.
			limiter := rate.NewLimiter(2, 10) // Ajuste: 2 req/seg, burst de 5
			visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
			v = visitors[ip]
		}

		v.lastSeen = time.Now()
		mu.Unlock()

		// Verifica se o requisitante (IP) excedeu o limite.
		if !v.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests) // Retorna 429
			return
		}

		// Se não excedeu, passa para o próximo handler (o login).
		next.ServeHTTP(w, r)
	})
}
