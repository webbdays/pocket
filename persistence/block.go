package persistence

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/pokt-network/pocket/persistence/kvstore"
	"github.com/pokt-network/pocket/persistence/types"
	"github.com/pokt-network/pocket/shared/codec"
	"github.com/pokt-network/pocket/shared/crypto"
	"github.com/pokt-network/pocket/shared/modules"
)

// OPTIMIZE(team): get from blockstore or keep in memory
func (p PostgresContext) GetLatestBlockHeight() (latestHeight uint64, err error) {
	ctx, tx, err := p.GetCtxAndTx()
	if err != nil {
		return 0, err
	}

	err = tx.QueryRow(ctx, types.GetLatestBlockHeightQuery()).Scan(&latestHeight)
	return
}

// OPTIMIZE: get from blockstore or keep in cache/memory
func (p PostgresContext) GetBlockHash(height int64) ([]byte, error) {
	ctx, tx, err := p.GetCtxAndTx()
	if err != nil {
		return nil, err
	}

	var hexHash string
	err = tx.QueryRow(ctx, types.GetBlockHashQuery(height)).Scan(&hexHash)
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(hexHash)
}

func (p PostgresContext) GetHeight() (int64, error) {
	return p.Height, nil
}

func (p PostgresContext) TransactionExists(transactionHash string) (bool, error) {
	hash, err := hex.DecodeString(transactionHash)
	if err != nil {
		return false, err
	}
	res, err := p.txIndexer.GetByHash(hash)
	if res == nil {
		// check for not found
		if err != nil && err.Error() == kvstore.BadgerKeyNotFoundError {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (p PostgresContext) StoreTransaction(txResult modules.TxResult) error {
	return p.txIndexer.Index(txResult)
}

// Creates a block protobuf object using the schema defined in the persistence module
func (p *PostgresContext) prepareBlock(proposerAddr []byte, quorumCert []byte) (*types.Block, error) {
	var prevHash []byte
	if p.Height == 0 {
		prevHash = []byte("")
	} else {
		var err error
		prevHash, err = p.GetBlockHash(p.Height - 1)
		if err != nil {
			return nil, err
		}
	}

	// The order (descending) is important here since it is used to comprise the hash in the block
	txResults, err := p.txIndexer.GetByHeight(p.Height, false)
	if err != nil {
		return nil, err
	}

	txs := make([]byte, 0)
	for _, txResult := range txResults {
		txHash, err := txResult.Hash()
		if err != nil {
			return nil, err
		}
		txs = append(txs, txHash...)
	}

	block := &types.Block{
		Height:            uint64(p.Height),
		Hash:              hex.EncodeToString(p.currentStateHash),
		PrevHash:          hex.EncodeToString(prevHash),
		ProposerAddress:   proposerAddr,
		QuorumCertificate: quorumCert,
		TransactionsHash:  crypto.SHA3Hash(txs), // TODO: Externalize this elsewhere?
	}

	return block, nil
}

// Inserts the block into the postgres database
func (p *PostgresContext) insertBlock(block *types.Block) error {
	ctx, tx, err := p.GetCtxAndTx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, types.InsertBlockQuery(block.Height, block.Hash, block.ProposerAddress, block.QuorumCertificate))
	return err
}

// Stores the block in the key-value store
func (p PostgresContext) storeBlock(block *types.Block) error {
	blockBz, err := codec.GetCodec().Marshal(block)
	if err != nil {
		return err
	}
	return p.blockStore.Set(heightToBytes(p.Height), blockBz)
}

func heightToBytes(height int64) []byte {
	heightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(heightBytes, uint64(height))
	return heightBytes
}
