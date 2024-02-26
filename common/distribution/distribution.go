package distribution

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Used for marshalling and unmarshalling big integers.
type BigInt struct {
	*big.Int
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" {
		return nil
	}
	var z big.Int
	_, ok := z.SetString(string(p), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = &z
	return nil
}

type Distribution struct {
	data *orderedmap.OrderedMap[gethcommon.Address, *orderedmap.OrderedMap[gethcommon.Address, *BigInt]]
}

func NewDistribution() *Distribution {
	data := orderedmap.New[gethcommon.Address, *orderedmap.OrderedMap[gethcommon.Address, *BigInt]]()
	return &Distribution{
		data: data,
	}
}

// Set sets the value for a given address.
func (d *Distribution) Set(address, token gethcommon.Address, amount *big.Int) {
	allocatedTokens, found := d.data.Get(address)
	if !found {
		allocatedTokens = orderedmap.New[gethcommon.Address, *BigInt]()
		d.data.Set(address, allocatedTokens)
	}
	allocatedTokens.Set(token, &BigInt{Int: amount})
}

// Get gets the value for a given address.
func (d *Distribution) Get(address, token gethcommon.Address) *big.Int {
	allocatedTokens, found := d.data.Get(address)
	if !found {
		return big.NewInt(0)
	}
	amount, found := allocatedTokens.Get(token)
	if !found {
		return big.NewInt(0)
	}
	return amount.Int
}

// Add adds the other distribution to this distribution.
// assumes other is non nil
func (d *Distribution) Add(other *Distribution) {
	for accountPair := other.data.Oldest(); accountPair != nil; accountPair = accountPair.Next() {
		address := accountPair.Key
		for tokenPair := accountPair.Value.Oldest(); tokenPair != nil; tokenPair = tokenPair.Next() {
			token := tokenPair.Key
			amount := tokenPair.Value
			currentAmount := d.Get(address, token)
			d.Set(address, token, currentAmount.Add(currentAmount, amount.Int))
		}
	}
}

func (d *Distribution) GetNumLeaves() int {
	numLeaves := 0
	for accountPair := d.data.Oldest(); accountPair != nil; accountPair = accountPair.Next() {
		numLeaves += accountPair.Value.Len()
	}
	return numLeaves
}

func (d *Distribution) MarshalJSON() ([]byte, error) {
	return d.data.MarshalJSON()
}

func (d *Distribution) UnmarshalJSON(p []byte) error {
	data := orderedmap.New[gethcommon.Address, *orderedmap.OrderedMap[gethcommon.Address, *BigInt]]()
	err := data.UnmarshalJSON(p)
	if err != nil {
		return err
	}
	d.data = data
	return nil
}

// Merklizes the distribution and returns the account tree and the token trees.
func (d *Distribution) Merklize() (*merkletree.MerkleTree, []*merkletree.MerkleTree, error) {
	tokenTrees := make([]*merkletree.MerkleTree, 0)

	// todo: parallelize this
	accountLeafs := make([][]byte, d.data.Len())

	for accountPair := d.data.Oldest(); accountPair != nil; accountPair = accountPair.Next() {
		address := accountPair.Key
		// fetch the leafs for the tokens for this account
		tokenLeafs := make([][]byte, accountPair.Value.Len())
		for tokenPair := accountPair.Value.Oldest(); tokenPair != nil; tokenPair = tokenPair.Next() {
			token := tokenPair.Key
			amount := tokenPair.Value
			tokenLeafs = append(tokenLeafs, encodeTokenLeaf(token, amount.Int))
		}

		// create a merkle tree for the tokens for this account
		tokenTree, err := merkletree.NewTree(
			merkletree.WithData(tokenLeafs),
			merkletree.WithHashType(keccak256.New()),
		)
		if err != nil {
			return nil, nil, err
		}
		tokenTrees = append(tokenTrees, tokenTree)

		// append the root to the list of account leafs
		accountRoot := tokenTree.Root()
		accountLeafs = append(accountLeafs, encodeAccountLeaf(address, accountRoot))
	}

	accountTree, err := merkletree.NewTree(
		merkletree.WithData(accountLeafs),
		merkletree.WithHashType(keccak256.New()),
	)
	if err != nil {
		return nil, nil, err
	}

	return accountTree, tokenTrees, nil
}

// encodeAccountLeaf encodes an account leaf for a token distribution.
func encodeAccountLeaf(account gethcommon.Address, accountRoot []byte) []byte {
	// (account || accountRoot)
	return append(account.Bytes(), accountRoot[:]...)
}

// encodeTokenLeaf encodes a token leaf for a token distribution.
func encodeTokenLeaf(token gethcommon.Address, amount *big.Int) []byte {
	// todo: handle this better
	amountU256, _ := uint256.FromBig(amount)
	amountBytes := amountU256.Bytes32()
	// (token || amount)
	return append(token.Bytes(), amountBytes[:]...)
}
