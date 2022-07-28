package persistence

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/pokt-network/pocket/persistence/schema"
)

// OPTIMIZE(team): get from blockstore or keep in memory
func (p PostgresContext) GetLatestBlockHeight() (latestHeight uint64, err error) {
	ctx, txn, err := p.DB.GetCtxAndTxn()
	if err != nil {
		return 0, err
	}

	err = txn.QueryRow(ctx, schema.GetLatestBlockHeightQuery()).Scan(&latestHeight)
	return
}

func (p PostgresContext) GetBlockHash(height int64) ([]byte, error) {
	ctx, txn, err := p.DB.GetCtxAndTxn()
	if err != nil {
		return nil, err
	}

	var hexHash string
	err = txn.QueryRow(ctx, schema.GetBlockHashQuery(height)).Scan(&hexHash)
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(hexHash)
}

func (p PostgresContext) NewSavePoint(bytes []byte) error {
	log.Println("Block - NewSavePoint not implemented")
	return nil
}

func (p PostgresContext) RollbackToSavePoint(bytes []byte) error {
	log.Println("TODO: Block - RollbackToSavePoint not implemented")
	return nil
}

func (p PostgresContext) AppHash() ([]byte, error) {
	log.Println("TODO: Block - AppHash not implemented")
	return []byte("A real app hash, I am not"), nil
}

func (p PostgresContext) Reset() error {
	log.Println("TODO: Block - Reset not implemented")
	return nil
}

func (p PostgresContext) Commit() error {
	p.DB.Tx.Commit(context.TODO())
	return nil
}

func (p PostgresContext) Release() {
	p.DB.Tx.Rollback(context.TODO())
}

func (p PostgresContext) GetHeight() (int64, error) {
	return p.Height, nil
}

func (p PostgresContext) TransactionExists(transactionHash string) bool {
	log.Println("TODO: Block - TransactionExists not implemented")
	return true
}
