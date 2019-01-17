package main

import (
	"encoding/json"
	"log"
	"time"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

/**
 * Cria uma pool de conexao com o redis
*/
func getPool() *redis.Pool {
	if(pool == nil){
        pool = &redis.Pool{
			MaxIdle: 3,
			IdleTimeout: 240 * time.Second,
			Dial: func () (redis.Conn, error) { return redis.Dial("tcp", redisUrlServer) },
		  }
	}

	return pool	
}

/**
 * Adiciona cabeçalhos de resposta http
*/
func JsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/**
 * Serviço responsável por retornar se o cliente pode ser impactado
 * Retorna true - caso ele possa ser impactado
 * Retorna false - caso ele não possa ser impactado no momento
*/
func IsClienteImpactadoEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	
	clienteId := IsClienteImpactado( params["id"] )
	
	JsonResponse(w, http.StatusOK, clienteId)	
}

/**
 * Serviço responsável por remover um cliente da noventena
 * Recebe o id do cliente como parametro e o remove do cache
 * 
*/
func DeleteClienteImpactadoEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	DelValue( params["id"] )
	
	JsonResponse(w, http.StatusOK, true)
}

/**
 * Verifica se um cliente pode ser impactado pela pesquisa. 
 * Retorna true, caso o cliente possa ser impactado e salva este cliente no redis(cache)
 * Retorna false, caso o cliente ja tenha sido impactado anteriormente, dentro do período de 90 dias
*/
func IsClienteImpactado(idCliente string) bool {
	clienteInNoventena := GetChaveFromRedis(idCliente)
	/**
	 * Verifica no cache, se o cliente foi impactado nos últimos 90 dias (está em noventena).
	 * Estando o cliente, armazenado no cache, ele está em estado de noventena.
	 * Retorna false, o cliente não pode ser impactado novamente
	*/
	if len(clienteInNoventena) > 0 {
        return false
	}else{
		/**
		 * Cliente não está em estado de noventena. Não foi impactado nos últimos 90 dias.
		 * Armazena o cliente no cache, para que ele não possa ser impactado novamente(noventena) e retorn true
		*/
		SetChaveRedis(idCliente, idCliente, noventena)
        return true
	}
	
}

/**
 * Obtem o valor da chave armazenada no Redis (Cache). 
 * Se a chave não existir no cache, o nil será retornado
*/
func GetChaveFromRedis(chave string) string {
	conn := getPool().Get()
	value, _ := redis.String(conn.Do("GET", chave))
	defer conn.Close()
    return value
}

/**
 * Armazena a chave no Redis (Cache) 
 * e define um tempo de expiração para a chave. 
 * 
*/
func SetChaveRedis(chave string, valor string, ttl int)  {
	conn := getPool().Get()
	conn.Do("SETEX", chave, ttl, valor)
	defer conn.Close()
}

/**
 * Remove a chave e seu valor armazenados no Redis (Cache). 
 * 
*/
func DelValue(chave string) bool {
	conn := getPool().Get()
	conn.Do("DEL", chave)
	defer conn.Close()
	return true
}

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/clientes/{id}", IsClienteImpactadoEndpoint).Methods("GET")
	r.HandleFunc("/clientes/{id}", DeleteClienteImpactadoEndpoint).Methods("DELETE")
	
	if err := http.ListenAndServe( serverPort, r); err != nil {
		log.Fatal(err)
	}
	
}
