package ethereum

import (
	"encoding/gob"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mitchellh/go-homedir"
	"math/big"
	"os"
	"path/filepath"
)

var SmartPoolDir = getSmartPoolDir()
var CounterFile = getCounterFile()
var SharesFile = getSharesFile()

func getSmartPoolDir() string {
	result, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(result, ".smartpool")
}

func getCounterFile() string {
	return filepath.Join(SmartPoolDir, "counter")
}

func getSharesFile() string {
	return filepath.Join(SmartPoolDir, "active_shares")
}

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (*FileStorage) PersistLatestCounter(counter *big.Int) error {
	err := os.MkdirAll(SmartPoolDir, 0766)
	if err != nil {
		return err
	}
	f, err := os.Create(CounterFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(counter)
	return err
}

func (*FileStorage) LoadLatestCounter() (*big.Int, error) {
	counter := big.NewInt(0)
	f, err := os.Open(CounterFile)
	if err != nil {
		return counter, err
	}
	dec := gob.NewDecoder(f)
	err = dec.Decode(counter)
	return counter, err
}

type gobShare struct {
	BlockHeader     *types.Header
	Nonce           types.BlockNonce
	MixDigest       common.Hash
	ShareDifficulty *big.Int
	MinerAddress    string
	SolutionState   int
}

func (*FileStorage) PersistActiveShares(shares []*Share) error {
	err := os.MkdirAll(SmartPoolDir, 0766)
	if err != nil {
		return err
	}
	f, err := os.Create(SharesFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	gobShares := []gobShare{}
	for _, s := range shares {
		gobShares = append(gobShares, gobShare{
			s.BlockHeader(),
			s.nonce,
			s.mixDigest,
			s.shareDifficulty,
			s.minerAddress,
			s.SolutionState,
		})
	}
	err = enc.Encode(gobShares)
	return err
}

func (*FileStorage) LoadActiveShares() ([]*Share, error) {
	shares := []*Share{}
	gobShares := []gobShare{}
	f, err := os.Open(SharesFile)
	if err != nil {
		return shares, err
	}
	dec := gob.NewDecoder(f)
	err = dec.Decode(&gobShares)
	if err != nil {
		return shares, err
	}
	for _, gobShare := range gobShares {
		shares = append(shares, &Share{
			gobShare.BlockHeader,
			gobShare.Nonce,
			gobShare.MixDigest,
			gobShare.ShareDifficulty,
			gobShare.MinerAddress,
			gobShare.SolutionState,
			nil,
		})
	}
	return shares, nil
}
