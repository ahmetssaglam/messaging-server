package configs

import (
	"messaging-server/internal/models"
	pkgUtils "messaging-server/pkg/utils"
)

var PostgresConfig = models.PostgresConfigStruct{
	ConnStr: pkgUtils.GetEnvStr("POSTGRES_URI", ""),
}
