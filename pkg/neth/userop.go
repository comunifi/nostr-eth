package neth

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Allowed function signatures
	FuncSigSingle             = crypto.Keccak256([]byte("execute(address,uint256,bytes)"))[:4]
	FuncSigBatch              = crypto.Keccak256([]byte("executeBatch(address[],uint256[],bytes[])"))[:4]
	FuncSigSafeExecFromModule = crypto.Keccak256([]byte("execTransactionFromModule(address,uint256,bytes,uint8)"))[:4]
)

type UserOp struct {
	Sender               common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce                *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode             []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData             []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"         mapstructure:"callGasLimit"         validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"         mapstructure:"maxFeePerGas"         validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature            []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
}

// MarshalJSON returns a JSON encoding of the UserOperation.
func (op *UserOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Sender               string `json:"sender"`
		Nonce                string `json:"nonce"`
		InitCode             string `json:"initCode"`
		CallData             string `json:"callData"`
		CallGasLimit         string `json:"callGasLimit"`
		VerificationGasLimit string `json:"verificationGasLimit"`
		PreVerificationGas   string `json:"preVerificationGas"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		PaymasterAndData     string `json:"paymasterAndData"`
		Signature            string `json:"signature"`
	}{
		Sender:               op.Sender.String(),
		Nonce:                hexutil.EncodeBig(op.Nonce),
		InitCode:             hexutil.Encode(op.InitCode),
		CallData:             hexutil.Encode(op.CallData),
		CallGasLimit:         hexutil.EncodeBig(op.CallGasLimit),
		VerificationGasLimit: hexutil.EncodeBig(op.VerificationGasLimit),
		PreVerificationGas:   hexutil.EncodeBig(op.PreVerificationGas),
		MaxFeePerGas:         hexutil.EncodeBig(op.MaxFeePerGas),
		MaxPriorityFeePerGas: hexutil.EncodeBig(op.MaxPriorityFeePerGas),
		PaymasterAndData:     hexutil.Encode(op.PaymasterAndData),
		Signature:            hexutil.Encode(op.Signature),
	})
}

// UnmarshalJSON parses a JSON encoding of the UserOperation.
func (op *UserOp) UnmarshalJSON(input []byte) error {
	type Alias struct {
		Sender               string `json:"sender"`
		Nonce                string `json:"nonce"`
		InitCode             string `json:"initCode"`
		CallData             string `json:"callData"`
		CallGasLimit         string `json:"callGasLimit"`
		VerificationGasLimit string `json:"verificationGasLimit"`
		PreVerificationGas   string `json:"preVerificationGas"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		PaymasterAndData     string `json:"paymasterAndData"`
		Signature            string `json:"signature"`
	}

	aux := &Alias{}
	if err := json.Unmarshal(input, aux); err != nil {
		return err
	}

	op.Sender = common.HexToAddress(aux.Sender)
	op.Nonce, _ = hexutil.DecodeBig(aux.Nonce)
	op.InitCode, _ = hexutil.Decode(aux.InitCode)
	op.CallData, _ = hexutil.Decode(aux.CallData)
	op.CallGasLimit, _ = hexutil.DecodeBig(aux.CallGasLimit)
	op.VerificationGasLimit, _ = hexutil.DecodeBig(aux.VerificationGasLimit)
	op.PreVerificationGas, _ = hexutil.DecodeBig(aux.PreVerificationGas)
	op.MaxFeePerGas, _ = hexutil.DecodeBig(aux.MaxFeePerGas)
	op.MaxPriorityFeePerGas, _ = hexutil.DecodeBig(aux.MaxPriorityFeePerGas)
	op.PaymasterAndData, _ = hexutil.Decode(aux.PaymasterAndData)
	op.Signature, _ = hexutil.Decode(aux.Signature)

	return nil
}

func (u *UserOp) Copy() UserOp {
	copy := UserOp{
		Sender:               u.Sender,
		Nonce:                new(big.Int).Set(u.Nonce),
		InitCode:             append([]byte(nil), u.InitCode...),
		CallData:             append([]byte(nil), u.CallData...),
		CallGasLimit:         new(big.Int).Set(u.CallGasLimit),
		VerificationGasLimit: new(big.Int).Set(u.VerificationGasLimit),
		PreVerificationGas:   new(big.Int).Set(u.PreVerificationGas),
		MaxFeePerGas:         new(big.Int).Set(u.MaxFeePerGas),
		MaxPriorityFeePerGas: new(big.Int).Set(u.MaxPriorityFeePerGas),
		PaymasterAndData:     append([]byte(nil), u.PaymasterAndData...),
		Signature:            append([]byte(nil), u.Signature...),
	}

	return copy
}

func (u *UserOp) GetHash(chainID *big.Int) string {
	// ABI encode only chainID, sender, and nonce
	packed := make([]byte, 0, 96) // Pre-allocate for 3 * 32 bytes

	// Encode chainID (uint256 - 32 bytes)
	chainIDPadded := make([]byte, 32)
	chainID.FillBytes(chainIDPadded)
	packed = append(packed, chainIDPadded...)

	// Encode sender (address - 32 bytes)
	senderPadded := make([]byte, 32)
	copy(senderPadded[12:], u.Sender.Bytes())
	packed = append(packed, senderPadded...)

	// Encode nonce (uint256 - 32 bytes)
	noncePadded := make([]byte, 32)
	u.Nonce.FillBytes(noncePadded)
	packed = append(packed, noncePadded...)

	// Return the keccak256 hash of the packed data
	return crypto.Keccak256Hash(packed).Hex()
}
