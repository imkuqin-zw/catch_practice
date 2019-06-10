package service

import "shop/final_consistency/models"

//存储预发送消息
func (s *Service) StoreMsgWaitingConfirm(msg *models.TransactionMsg) {

}

//确认预发送消息并发送消息
func (s *Service) ConfirmAndSendMessage() {

}

//查询并处理超时的预发送消息
func (s *Service) DealTimeOutMsgWaitingConfirm() {

}

//func (s *Service) DealUnSend
