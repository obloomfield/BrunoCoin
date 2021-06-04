package pkg

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, John Roy,
 *	Parker Ljung
 *
 */

// ChkBlk (CheckBlock) validates a block based on multiple
// conditions.
// To be valid:
// The block must be syntactically (ChkBlkSyn), semantically
// (ChkBlkSem), and configurally (ChkBlkConf) valid.
// Each transaction on the block must be syntactically (ChkTxSyn),
// semantically (ChkTxSem), and configurally (ChkTxConf) valid.
// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// Inputs:
// b *block.Block the block to be checked for validity
// Returns:
// bool True if the block is valid. false
// otherwise
// TODO:
// to be valid

// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// The block's size must be less than or equal to the largest
// allowed block size.
// The block hash must be less than the difficulty target.
// The block's first transaction must be of type Coinbase.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// t.IsCoinbase()
// b.SatisfiesPOW(...)
// n.Conf.MxBlkSz
// b.Sz()
// n.Chain.ChkChainsUTXO(...)
func (n *Node) ChkBlk(b *block.Block) bool {
	//testing all nil cases, I think some of these should be important :P:
	if b == nil || b.Transactions == nil || len(b.Transactions) == 0 || b.Transactions[0] == nil || b.Transactions[0].SumOutputs() == 0 {
		return false
	}
	isFirstCoinbase := b.Transactions[0].IsCoinbase()
	isBlockSizeValid := b.Sz() < n.Conf.MxBlkSz
	isPOWSatisfied := b.SatisfiesPOW(b.Hdr.DiffTarg)
	isValidUTXO := n.Chain.ChkChainsUTXO(b.Transactions, b.Hdr.PrvBlkHsh)
	return isFirstCoinbase && isBlockSizeValid && isPOWSatisfied && isValidUTXO
}

// ChkTx (CheckTransaction) validates a transaction.
// Inputs:
// t *tx.Transaction the transaction to be checked for validity
// Returns:
// bool True if the transaction is syntactically valid. false
// otherwise
// TODO:
// to be valid:

// The transaction's inputs and outputs must not be empty.
// The transaction's output amounts must be larger than 0.
// The sum of the transaction's inputs must be larger
// than the sum of the transaction's outputs.
// The transaction must not double spend any UTXO.
// The unlocking script on each of the transaction's
// inputs must successfully unlock each of the corresponding
// UTXO.
// The transaction must not be larger than the
// maximum allowed block size.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// note: let u be a transaction output object
// n.Conf.MxBlkSz
// t.Sz()
// u.IsUnlckd(...)
// n.Chain.GetUTXO(...)
// n.Chain.IsInvalidInput(...)
// t.SumInputs()
// t.SumOutputs()
func (n *Node) ChkTx(t *tx.Transaction) bool {
	noEmptyInputs := len(t.Inputs) > 0 && len(t.Outputs) > 0
	inSumLarger := t.SumInputs() > t.SumOutputs()
	smallerThanMax := t.Sz() <= n.Conf.MxBlkSz
	TxiTxoValidity := checkValidPerTX(t.Inputs, t.Outputs, n)

	return noEmptyInputs && inSumLarger && smallerThanMax && TxiTxoValidity
}

func checkValidPerTX(txis []*txi.TransactionInput, txos []*txo.TransactionOutput, n *Node) bool {
	for _, txi := range txis {
		//no double spend issue || unlocking script issue
		if n.Chain.IsInvalidInput(txi) || !n.Chain.GetUTXO(txi).IsUnlckd(txi.UnlockingScript) {
			return false
		}
	}
	for _, txo := range txos {
		// amount isn't 0 or less
		if txo.Amount <= 0 {
			return false
		}
	}
	return true
}
