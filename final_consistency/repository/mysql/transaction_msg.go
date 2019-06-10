package mysql

import (
	"shop/final_consistency/models"
)

func (r *RepoMysql) GetTransMsgById(id uint64) (*models.TransactionMsg, error) {
	m := &models.TransactionMsg{Id: id}
	if err := r.GetById(db, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *RepoMysql) InsertTransMsg(m *models.TransactionMsg) error {
	return r.insert(db, m)
}
