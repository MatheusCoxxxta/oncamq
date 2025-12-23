# OnçaMQ


(PT-BR)

Busquei algumas opções para consumir filas do BullMQ com Golang pra ganhar throughput e poder do para as atividades CPU-intesive de um produto Node.js. A proposta era simples: o app Node.js já está adicionando jobs a fila e lidando com o consumo deles. Agora o consumo ficaria a cargo de um worker Golang.

Bem rapidamente entendi que não tem uma solução oficial, até existem algumas libs que fazem o processo completo que o BullMQ faz, mas foge do meu ponto. Achei também algumas pessoas que tomaram a iniciativa de criar um consumer para o padrão BullMQ com Golang bem recentemente, para preencher essa lacuna. Como não tem uma solução fortemente utilizada pela comunidade, e não me dei tão bem com a DX de algumas soluções, resolvi fazer a minha.

(en-US)

I looked into some options for consuming BullMQ queues with Golang to gain throughput and CPU power for the Node.js product's intensive tasks. The idea was simple: the Node.js app is already adding jobs to the queue and handling some consumption. Now, the consumption would be handled by a Golang worker.

Pretty quickly I realized there’s no official solution. There are some libraries that replicate BullMQ’s full workflow, but that wasn’t my goal. I also found a few people who recently took the initiative to create a Golang consumer for the BullMQ pattern to fill this gap. Since there’s no widely adopted solution in the community, and I didn’t have the best experience with the DX of some options, I decided to build my own.

O uso é simples:

- Importações necessárias:

```golang
import (
	oncamq "github.com/MatheusCoxxxta/oncamq/worker"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)
```

- Crie uma instancia de Redis

```golang
func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6370",
		Password: "redis",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
}
```

- Crie a instancia de worker que será utilizada para o consumo, passando a instancia de redis que será consumida, a fila e os handlers que serão utilizados para disparar ações pré-definidas de acordo com o nome do job.

```golang
	emailQueueWorker := oncamq.Worker{
		Instance: redisClient,
		Queue:    "emailQueue",
		Handlers: oncamq.Handlers{
			"firstAcess": SendFirstMail,
		},
	}
```

- Adicione essa instancia de worker na chamada para iniciar o worker

```golang

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6370",
		Password: "redis",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	emailQueueWorker := oncamq.Worker{
		Instance: redisClient,
		Queue:    "emailQueue",
		Handlers: oncamq.Handlers{
			"firstAcess": SendFirstMail,
		},
	}

	oncamq.StartWorker(emailQueueWorker)
}
```

- Para usar workers para filas diferentes de forma concorrente, use uma Goroutine

```golang

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6370",
		Password: "redis",
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	emailQueueWorker := oncamq.Worker{
		Instance: redisClient,
		Queue:    "emailQueue",
		Handlers: oncamq.Handlers{
			"firstAcess": SendFirstMail,
		},
	}

	paymentQueueWorker := oncamq.Worker{
		Instance: redisClient,
		Queue:    "paymentQueue",
		Handlers: oncamq.Handlers{
			"createCustomer":   CreateCustomer,
			"startTransaction": StartTransaction,
		},
	}

	go oncamq.StartWorker(emailQueueWorker)
	go oncamq.StartWorker(paymentQueueWorker)

    select {}
}
```
