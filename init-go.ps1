New-Item -ItemType Directory -Force -Path api-gateway/cmd/server, api-gateway/internal/config, api-gateway/internal/delivery/http, api-gateway/internal/delivery/mqtt, api-gateway/internal/domain, api-gateway/internal/repository, api-gateway/internal/usecase, api-gateway/pkg/logger

Set-Location api-gateway
go mod init github.com/smartfactory/api-gateway
go get -u github.com/gin-gonic/gin github.com/eclipse/paho.mqtt.golang github.com/jmoiron/sqlx github.com/lib/pq github.com/influxdata/influxdb-client-go/v2
