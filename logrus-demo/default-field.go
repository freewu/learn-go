package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	requestLogger := log.WithFields(log.Fields{"request_id": "request_id", "user_ip": "user_ip"})
	requestLogger.Info("something happened on that request") // time="2021-09-26T16:21:00+08:00" level=info msg="something happened on that request" request_id=request_id user_ip=user_ip
	requestLogger.Warn("something not great happened") // time="2021-09-26T16:21:00+08:00" level=warning msg="something not great happened" request_id=request_id user_ip=user_ip

	/*
	默认为 time / msg / level 三个字段默认添加
	1 time. The timestamp when the entry was created.
	2 msg. The logging message passed to {Info,Warn,Error,Fatal,Panic} after the AddFields call. E.g. Failed to send event.
	3 level. The logging level. E.g. info.
	 */

}
