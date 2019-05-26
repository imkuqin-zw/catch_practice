package handler

import z "shop/plugins/zap"

var log *z.Logger

func Init() {
	log = z.GetLogger()
}
